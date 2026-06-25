package utils

import "time"

type TargetDto struct {
	LastAccessed time.Time
	Name         string
	Namespace    string
	Resource     string
}

var Targets map[string]*TargetDto
var SyncTime int = 10
