package main

import (
	"errors"
	"log"
	"os"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	KubeRunController string   `yaml:"kube_run_controller"`
	Update            bool     `yaml:"update"`
	Ips               []string `yaml:"ips"`
}

var Config AgentConfig
var configPath = "/etc/agent-config/config.yml"

func LoadVariables() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch errno {
		case syscall.ENOMEM:
			HandelError(err, "KR0012", "Out of Kernel Memory")
		case syscall.EMFILE:
			HandelError(err, "KR0024", "Too Many Active Watcher Instances")
		case syscall.ENFILE:
			HandelError(err, "KR0023", "System-Wide File Descriptor Exhaustion")
		default:
			HandelError(err, "KR9010", "Generic System Initialization Failure")
		}
	} else {
		HandelError(err, "KR9011", "Unknown Watcher Error")
	}

	readVariables()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				readVariables()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				HandelError(err, "PlaceHolder", "_")
			}
		}
	}()

	err = watcher.Add(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return watcher

}

func readVariables() {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		HandelError(err, "KR0404", "Config file not found in /etc/agent-config/config.yml")
	}
	yaml.Unmarshal(bytes, &Config)
}
