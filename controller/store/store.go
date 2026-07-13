package store

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	Mux          sync.Mutex       `json:"-"`
	Servers      []*server.Switch `json:"-"`
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
var syncMinutes time.Duration = getSyncMinutes()
var SyncTime = syncMinutes * time.Minute / 2
var KubeRunNamespace = getEnvString("KUBERUN_NAMESPACE", "default")
var KubeRunAgentConfigName = getEnvString("KUBERUN_AGENT_CONFIG_NAME", "kuberun-agent-config")
var KubeRunAgent = getKubeRunAgent()

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

func getSyncMinutes() time.Duration {
	const defaultMinutes = 15

	raw := os.Getenv("SYNC_MINUTES")
	if raw == "" {
		return defaultMinutes
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		fmt.Printf("Invalid SYNC_MINUTES value %q, falling back to default of %v minute(s)\n", raw, defaultMinutes)
		return defaultMinutes
	}

	return time.Duration(parsed)
}

func getEnvString(key, defaultValue string) string {
	val, set := os.LookupEnv(key)
	if !set {
		return defaultValue
	}
	if val == "" {
		fmt.Printf("Warning: environment variable %s is empty-string-but-set, falling back to default %q\n", key, defaultValue)
		return defaultValue
	}
	return val
}

func getKubeRunAgent() string {
	agentServiceName := getEnvString("KUBERUN_AGENT_SERVICE_NAME", "kuberun-agent")
	return fmt.Sprintf("%s.%s.svc.cluster.local", agentServiceName, KubeRunNamespace)
}
