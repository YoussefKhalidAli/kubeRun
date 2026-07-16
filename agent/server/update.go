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

	logger.Info("update listener started", "addr", ":4443")

	err := http.ListenAndServe(":4443", nil)
	if err != nil {
		utils.HandelError(err, "KRA9020H", "Couldn't boot up update listener.")
	}

}

func updatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	ip, err := io.ReadAll(r.Body)
	if err != nil {
		utils.HandelError(err, "KRA9021M", "Couldn't parse update body.")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	logger.Info("received ip update", "ip", string(ip))
	configMu.Lock()
	store.Config.Ips = append(store.Config.Ips, string(ip))
	configMu.Unlock()
}
