package server

import (
	"fmt"
	"net/http"
	"strings"

	"kuberun.com/agent/store"
)

func Alert(ip string) {
	controllerUrl := fmt.Sprintf("http://%v/alert", store.Config.KubeRunController)
	println(controllerUrl)

	target := ip
	if val, exists := store.Config.HeadlessMap[ip]; exists {
		target = val
	}

	http.Post(controllerUrl, "text/plain", strings.NewReader(target))
}
