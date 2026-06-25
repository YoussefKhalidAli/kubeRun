package main

import (
	"fmt"
	"net/http"
	"strings"
)

func Alert(ip string) {
	controllerUrl := fmt.Sprintf("http://%v/alert", Config.KubeRunController)
	println(controllerUrl)
	http.Post(controllerUrl, "text/plain", strings.NewReader(ip))
}
