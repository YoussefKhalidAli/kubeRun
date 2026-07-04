package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"

	"kuberun.com/controller/utils"
)

type Switch struct {
	Port         int
	mux          *http.ServeMux
	server       *http.Server
	serverMu     sync.Mutex
	Proxy        atomic.Value
	Signal       sync.RWMutex
	shutdownOnce sync.Once
	ScaleUp      func()
}

var Switches int = 0

func New() *Switch {
	port := Switches + 4445
	sw := &Switch{
		Port: port,
	}
	sw.Signal.Lock()
	Switches++
	return sw
}

func (sw *Switch) Start() {
	sw.serverMu.Lock()
	mux := http.NewServeMux()
	mux.HandleFunc("/", sw.SwitchHandler)
	sw.mux = mux
	sw.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", sw.Port),
		Handler: mux,
	}
	server := sw.server
	sw.shutdownOnce = sync.Once{}
	sw.serverMu.Unlock()

	fmt.Printf("switch number %v listener booted successfully\n", sw.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		utils.HandelError(err, "KRC9019", fmt.Sprintf("Couldn't boot up switch number %v server.", sw.Port))
	}
}

func (sw *Switch) Stop() {
	sw.serverMu.Lock()
	server := sw.server
	once := &sw.shutdownOnce
	sw.serverMu.Unlock()

	once.Do(func() {
		if server != nil {
			server.Shutdown(context.Background())
		}
		sw.Signal.Lock()
	})
}

func (sw *Switch) SwitchHandler(w http.ResponseWriter, r *http.Request) {
	sw.ScaleUp()
	sw.Signal.RLock()
	defer sw.Signal.RUnlock()

	println("Unlocked")
	proxyDestination, ok := sw.Proxy.Load().(string)

	if ok {
		println("proxyDestination", proxyDestination)
		proxyReq(w, r, proxyDestination)
	} else {
		message := []byte("Couldn't find pod to proxy")
		w.Write(message)
	}

	go sw.Stop()
}

func proxyReq(w http.ResponseWriter, req *http.Request, destination string) {
	target, err := url.Parse(destination)
	if err != nil {
		http.Error(w, "invalid destination", http.StatusBadGateway)
		return
	}

	_, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		if target.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}

	target.Host = net.JoinHostPort(target.Hostname(), port)
	proxy := httputil.NewSingleHostReverseProxy(target)
	req.URL.Host = target.Host
	req.URL.Scheme = target.Scheme
	req.Host = target.Host
	proxy.ServeHTTP(w, req)
}
