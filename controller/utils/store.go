package utils

import "time"

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
