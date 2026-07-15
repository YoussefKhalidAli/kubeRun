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
	Shadow       string
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
	jsonData, err := json.MarshalIndent(Targets, "", "  ")
	if err != nil {
		fmt.Printf("Error printing map: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}
