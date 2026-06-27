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
var configDir = "/etc/agent-config"
var configPath = configDir + "/config.yml"

func LoadVariables() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {

		var errno syscall.Errno
		if errors.As(err, &errno) {
			switch errno {
			case syscall.ENOMEM:
				HandelError(err, "KRA0012", "Out of Kernel Memory")
			case syscall.EMFILE:
				HandelError(err, "KRA0024", "Too Many Active Watcher Instances")
			case syscall.ENFILE:
				HandelError(err, "KRA0023", "System-Wide File Descriptor Exhaustion")
			default:
				HandelError(err, "KRA9010", "Generic System Initialization Failure")
			}
		} else {
			HandelError(err, "KRA9011", "Unknown Watcher Error")
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
				log.Println("event:", event)
				readVariables()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				HandelError(err, "KRA9012", "fsnotify reported an error while watching the config file")
			}
		}
	}()

	err = watcher.Add(configDir)
	if err != nil {
		log.Fatal(err)
	}

	return watcher

}

func readVariables() {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		HandelError(err, "KRA0404", "Config file not found in /etc/agent-config/config.yml")
	}
	yaml.Unmarshal(bytes, &Config)
}
