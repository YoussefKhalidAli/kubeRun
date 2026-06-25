package main

import (
	"io"
	"net/http"
)

func Alert() {

	http.HandleFunc("/alert", alertHandler)

	println("Alert listener booted successfully")

	err := http.ListenAndServe(":4444", nil)
	if err != nil {
		HandelError(err, "KRC9011", "Couldn't boot up http server.")
	}
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, err := io.ReadAll(r.Body)
	if err != nil {
		HandelError(err, "KRC9010", "Couldn't parse alert body.")
	}
	defer r.Body.Close()

	println(string(ip))

}
