package server

import (
	"io"
	"net/http"
	"sync"

	"kuberun.com/agent/store"
	"kuberun.com/agent/utils"
)

var configMu sync.Mutex

func Updates() {
	http.HandleFunc("/update", updatesHandler)

	println("Booted up update listener")

	err := http.ListenAndServe(":4443", nil)
	if err != nil {
		utils.HandelError(err, "KRA9020", "Couldn't boot up update listener.")
	}

}

func updatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, err := io.ReadAll(r.Body)
	if err != nil {
		utils.HandelError(err, "KRA9021", "Couldn't parse update body.")
	}
	defer r.Body.Close()
	println("Recieved: ", string(ip))
	configMu.Lock()
	store.Config.Ips = append(store.Config.Ips, string(ip))
	configMu.Unlock()
}
