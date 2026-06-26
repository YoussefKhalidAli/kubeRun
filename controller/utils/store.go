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
}

var Targets map[string]*TargetDto

// Configs
var SyncTime int = 10

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
