package utils

import "fmt"

func GetShadowName(name string) string {
	return fmt.Sprintf("shadow-%v", name)
}

func GetHeadlessServiceKey(name string) string {
	return fmt.Sprintf("svc-%v", name)
}
