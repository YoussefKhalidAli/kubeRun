package main

import (
	"time"
)

type TargetDto struct {
	LastAccessed time.Time
	Name         string
	Namespace    string
	Resource     string
}

var Targets map[string]*TargetDto

func main() {
	Targets = make(map[string]*TargetDto)
	Targets["10.110.91.113"] = &TargetDto{
		LastAccessed: time.Now().Add(-24 * time.Hour),
		Name:         "auth-service-pod",
		Namespace:    "production",
		Resource:     "pods",
	}

	Targets["10.106.160.35"] = &TargetDto{
		LastAccessed: time.Now().Add(-24 * time.Hour),
		Name:         "payment-gateway-svc",
		Namespace:    "staging",
		Resource:     "services",
	}
	Alert()

}
