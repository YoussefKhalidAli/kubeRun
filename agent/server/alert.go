package server

import (
	"fmt"
	"net/http"
	"strings"

	"kuberun.com/agent/store"
	"kuberun.com/agent/utils"
)

var logger = utils.Logger.With("module", "server")

func Alert(ip string) {
	controllerUrl := fmt.Sprintf("http://%v/alert", store.Config.KubeRunController)

	target := ip
	if val, exists := store.Config.HeadlessMap[ip]; exists {
		target = val
	}

	logger.Debug("sending alert", "controller_url", controllerUrl, "ip", ip, "target", target)

	resp, err := http.Post(controllerUrl, "text/plain", strings.NewReader(target))
	if err != nil {
		utils.HandelError(err, "KRA1453M", "Failed to send alert to controller")
		return
	}
	resp.Body.Close()
}
