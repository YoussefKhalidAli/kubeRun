package store

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"kuberun.com/controller/utils"
)

var logger = utils.Logger.With("module", "store")

type AgentConfig struct {
	KubeRunController string            `yaml:"kube_run_controller"`
	Update            bool              `yaml:"update"`
	Ips               []string          `yaml:"ips"`
	HeadlessMap       map[string]string `yaml:"headless_map"`
}

// Configs
var syncMinutes time.Duration = getSyncMinutes()
var SyncTime = syncMinutes * time.Minute / 2
var KubeRunNamespace = getEnvString("KUBERUN_NAMESPACE", "default")
var KubeRunAgentConfigName = getEnvString("KUBERUN_AGENT_CONFIG_NAME", "kuberun-agent-config")
var KubeRunAgent = getKubeRunAgent()
var KubeRunPodIp string

// Labels
var RunLabel = "kuberun/run=true"

func getSyncMinutes() time.Duration {
	const defaultMinutes = 15

	raw := os.Getenv("SYNC_MINUTES")
	if raw == "" {
		return defaultMinutes
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		logger.Warn("invalid sync minutes configuration", "value", raw, "fallback_default", defaultMinutes)
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
		logger.Warn("empty environment variable fallback", "key", key, "fallback_default", defaultValue)
		return defaultValue
	}
	return val
}

func getKubeRunAgent() string {
	agentServiceName := getEnvString("KUBERUN_AGENT_SERVICE_NAME", "kuberun-agent")
	return fmt.Sprintf("%s.%s.svc.cluster.local", agentServiceName, KubeRunNamespace)
}
