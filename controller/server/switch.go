package server

import (
	"fmt"
	"net/http"
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

	proxyReq(r, sw.Proxy)
	sw.server.Close()
	sw.Signal.Lock()
}

func proxyReq(req *http.Request, destination string) {

}
