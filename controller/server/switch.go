package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"kuberun.com/controller/utils"
)

type Switch struct {
	Port   int
	mux    *http.ServeMux
	server *http.Server
	Proxy  string
	Signal sync.RWMutex
}

var Switches []*Switch

func New() *Switch {
	mux := http.NewServeMux()
	port := len(Switches) + 4445
	sw := &Switch{
		Port:  port,
		mux:   mux,
		Proxy: "",
	}
	sw.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}
	sw.Signal.Lock()
	Switches = append(Switches, sw)
	sw.mux.HandleFunc("/", sw.SwitchHandler)
	return sw
}

func (sw *Switch) Start() {
	fmt.Printf("switch number %v listener booted successfully\n", sw.Port)
	err := sw.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		utils.HandelError(err, "KRC9011", fmt.Sprintf("Couldn't boot up switch number %v server.", sw.Port))
	}
}

func (sw *Switch) SwitchHandler(w http.ResponseWriter, r *http.Request) {
	sw.Signal.RLock()
	defer sw.Signal.RUnlock()

	proxyReq(w, r, sw.Proxy)
	sw.server.Shutdown(context.TODO())
	sw.Signal.Lock()
}

func proxyReq(w http.ResponseWriter, req *http.Request, destination string) {
	target, err := url.Parse(destination)
	if err != nil {
		http.Error(w, "invalid destination", http.StatusBadGateway)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	req.URL.Host = target.Host
	req.URL.Scheme = target.Scheme
	req.Host = target.Host

	proxy.ServeHTTP(w, req)
}
