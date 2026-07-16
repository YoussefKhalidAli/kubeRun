package store

import (
	"errors"
	"os"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"kuberun.com/agent/utils"
)

var logger = utils.Logger.With("module", "store")

type AgentConfig struct {
	KubeRunController string            `yaml:"kube_run_controller"`
	Update            bool              `yaml:"update"`
	Ips               []string          `yaml:"ips"`
	HeadlessMap       map[string]string `yaml:"headless_map"`
}

var Config AgentConfig
var configDir = "/etc/agent-config"
var configPath = configDir + "/config.yml"

func LoadVariables() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {

		var errno syscall.Errno
		if errors.As(err, &errno) {
			switch errno {
			case syscall.ENOMEM:
				utils.HandelError(err, "KRA0012H", "Out of Kernel Memory")
			case syscall.EMFILE:
				utils.HandelError(err, "KRA0024H", "Too Many Active Watcher Instances")
			case syscall.ENFILE:
				utils.HandelError(err, "KRA0023H", "System-Wide File Descriptor Exhaustion")
			default:
				utils.HandelError(err, "KRA9010H", "Generic System Initialization Failure")
			}
		} else {
			utils.HandelError(err, "KRA9011H", "Unknown Watcher Error")
		}
	}

	readVariables()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Info("config file event", "event", event.String())
				readVariables()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.HandelError(err, "KRA9012M", "fsnotify reported an error while watching the config file")
			}
		}
	}()

	err = watcher.Add(configDir)
	if err != nil {
		utils.HandelError(err, "KRA9010H", "Failed to add path to config file watcher")
	}

	return watcher

}

func readVariables() {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		utils.HandelError(err, "KRA0405H", "Config file not found in /etc/agent-config/config.yml")
		return
	}
	err = yaml.Unmarshal(bytes, &Config)
	if err != nil {
		utils.HandelError(err, "KRA9013H", "Failed to unmarshal agent config")
	}
}
