package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"kuberun.com/controller/kubernetes"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func Alert() {

	http.HandleFunc("/alert", alertHandler)

	println("Alert listener booted successfully")

	err := http.ListenAndServe(":4444", nil)
	if err != nil {
		utils.HandelError(err, "KRC9011", "Couldn't boot up alert server.")
	}
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, err := io.ReadAll(r.Body)
	if err != nil {
		utils.HandelError(err, "KRC9010", "Couldn't parse alert body.")
	}
	defer r.Body.Close()

	target := store.Targets[string(ip)]

	target.LastAccessed = time.Now()
	if target.IsSleep {
		kubernetes.ScaleResource(target, 1, string(ip))
	}
	fmt.Printf("Hit %v", target)
}
