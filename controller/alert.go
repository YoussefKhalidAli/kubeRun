package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func Alert() {

	http.HandleFunc("/alert", alertHandler)

	logger.Info("alert listener started", "addr", ":4444")

	err := http.ListenAndServe(":4444", nil)
	if err != nil {
		utils.HandelError(err, "KRC9011H", "Couldn't boot up alert server.")
	}
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	ip, err := io.ReadAll(r.Body)
	if err != nil {
		utils.HandelError(err, "KRC9010M", "Couldn't parse alert body.")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	ipStr := string(ip)
	logger.Info("alert received", "ip", ipStr)

	target, ok := store.Targets[ipStr]
	if !ok || target == nil {
		utils.HandelError(fmt.Errorf("alert from unknown IP: %s", ipStr), "KRC1448L", "alert from unknown IP")
		http.Error(w, "Unknown IP", http.StatusBadRequest)
		return
	}

	target.Mux.Lock()
	target.LastAccessed = time.Now()
	target.Mux.Unlock()

}
