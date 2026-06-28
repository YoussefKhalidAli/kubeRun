package server

import (
	"fmt"
	"net/http"

	"kuberun.com/controller/utils"
)

type Switch struct {
	Port   int
	mux    *http.ServeMux
	server *http.Server
	reqs   []*http.Request
	Signal chan bool
}

var Switches []*Switch

func New() *Switch {
	mux := http.NewServeMux()
	port := len(Switches) + 4445
	sw := &Switch{
		Port:   port,
		mux:    mux,
		reqs:   []*http.Request{},
		Signal: make(chan bool),
	}
	sw.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}
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

func (sw *Switch) Stop() {
	sw.server.Close()
}

func (sw *Switch) SwitchHandler(w http.ResponseWriter, r *http.Request) {
	sw.reqs = append(sw.reqs, r)

	release := <-sw.Signal

	if release {
		// proxy func
	}
	sw.Stop()
}

func proxyReq(req *http.Request) {

}
