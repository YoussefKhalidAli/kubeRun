package store

import (
	"encoding/json"
	"sync"
	"time"

	"kuberun.com/controller/server"
)

type TargetStatus string

const (
	Awake    TargetStatus = "Awake"
	Waking   TargetStatus = "Waking"
	Asleep   TargetStatus = "Asleep"
	Sleeping TargetStatus = "Sleeping"
)

type TargetDto struct {
	LastAccessed time.Time
	ResourceName string
	Namespace    string
	Resource     string
	ServiceName  string
	Status       TargetStatus
	UpdateMarker string
	Mux          sync.Mutex       `json:"-"`
	Servers      []*server.Switch `json:"-"`
	Endpoints    []string
	ServicePorts *[]int
	SelectorMap  map[string]string
}

var Targets map[string]*TargetDto

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

func PrintTargets() {
	jsonData, err := json.Marshal(Targets)
	if err != nil {
		logger.Error("error marshaling targets map", "error", err)
		return
	}
	logger.Debug("current targets", "targets_json", string(jsonData))
}
