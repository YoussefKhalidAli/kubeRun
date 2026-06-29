package store

import (
	"encoding/json"
	"fmt"
	"time"

	"kuberun.com/controller/server"
)

type TargetDto struct {
	LastAccessed time.Time
	ResourceName string
	Namespace    string
	Resource     string
	ServiceName  string
	Status       string
	Server       *server.Switch
	ServicePorts *[]int
	SelectorMap  map[string]string
}

type AgentConfig struct {
	KubeRunController string   `yaml:"kube_run_controller"`
	Update            bool     `yaml:"update"`
	Ips               []string `yaml:"ips"`
}

var Targets map[string]*TargetDto

// Configs
var syncMinutes time.Duration = 1
var SyncTime time.Duration = syncMinutes * time.Minute / 2
var KubeRunNamespace string = "default"
var KubeRunAgentConfigName string = "kuberun-agent-config"

// Annotations
var RunLabel string = "kuberun.com/run"

func PrintTargets() {
	jsonData, err := json.MarshalIndent(Targets, "", "  ")
	if err != nil {
		fmt.Printf("Error printing map: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}
