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
	http.Post(controllerUrl, "text/plain", strings.NewReader(ip))
}
