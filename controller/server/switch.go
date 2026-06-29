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
	Port   int
	mux    *http.ServeMux
	server *http.Server
	Proxy  atomic.Value
	Signal sync.RWMutex
}

var Switches int = 0

func New() *Switch {
	mux := http.NewServeMux()
	port := Switches + 4445
	sw := &Switch{
		Port: port,
		mux:  mux,
	}
	sw.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}
	sw.Signal.Lock()
	Switches++
	sw.mux.HandleFunc("/", sw.SwitchHandler)
	return sw
}

func (sw *Switch) Start() {
	fmt.Printf("switch number %v listener booted successfully\n", sw.Port)
	err := sw.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		utils.HandelError(err, "KRC9019", fmt.Sprintf("Couldn't boot up switch number %v server.", sw.Port))
	}
}

func (sw *Switch) SwitchHandler(w http.ResponseWriter, r *http.Request) {
	sw.Signal.RLock()
	defer sw.Signal.RUnlock()

	println("Unlocked")
	proxyDestination, _ := sw.Proxy.Load().(string)
	println("proxyDestination", proxyDestination)

	proxyReq(w, r, proxyDestination)
	sw.server.Shutdown(context.TODO())
	sw.Signal.Lock()
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
