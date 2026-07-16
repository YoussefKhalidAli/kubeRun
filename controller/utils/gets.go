package utils

import "fmt"

func GetShadowName(name string) string {
	return fmt.Sprintf("shadow-%v", name)
}
