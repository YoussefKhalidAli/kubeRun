package store

import (
	"encoding/json"
	"fmt"
	"sync"
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
	UpdateMarker string
	Mux          sync.Mutex     `json:"-"`
	Server       *server.Switch `json:"-"`
	Endpoints    []string
	ServicePorts *[]int
	SelectorMap  map[string]string
}

type AgentConfig struct {
	KubeRunController string            `yaml:"kube_run_controller"`
	Update            bool              `yaml:"update"`
	Ips               []string          `yaml:"ips"`
	HeadlessMap       map[string]string `yaml:"headless_map"`
}

var Targets map[string]*TargetDto

// Configs
var syncMinutes time.Duration = 1
var SyncTime = syncMinutes * time.Minute / 2
var KubeRunNamespace = "default"
var KubeRunAgentConfigName = "kuberun-agent-config"
var KubeRunAgent = "kuberun-agent.default.svc.cluster.local"

// Labels
var RunLabel = "kuberun/run=true"

func PrintTargets() {
	jsonData, err := json.MarshalIndent(Targets, "", "  ")
	if err != nil {
		fmt.Printf("Error printing map: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func (t *TargetDto) MarshalJSON() ([]byte, error) {
	type Alias TargetDto
	return json.Marshal(&struct {
		*Alias
		Mux    any `json:"Mux,omitempty"`    // Overwrite and ignore
		Server any `json:"Server,omitempty"` // Overwrite and ignore
	}{
		Alias: (*Alias)(t),
	})
}
