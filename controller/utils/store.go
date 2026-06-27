package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type TargetDto struct {
	LastAccessed time.Time
	ResourceName string
	Namespace    string
	Resource     string
	ServiceName  string
	IsSleep      bool
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
var SyncTime time.Duration = syncMinutes * time.Minute / 10
var KubeRunNamespace string = "default"
var KubeRunAgentConfigName string = "kuberun-agent-config"

// Annotations
var RunAnnotation string = "kuberun.com/run"

func PrintTargets() {
	jsonData, err := json.MarshalIndent(Targets, "", "  ")
	if err != nil {
		fmt.Printf("Error printing map: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}
