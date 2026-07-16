package utils

import (
	"fmt"
	"strings"
)

const (
	ColorReset = "\033[0m"
	ColorCyan  = "\033[36m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
)

func HandelError(err error, t string, details string) {
	if !strings.Contains(t, "L") {
		fmt.Printf("Error: %v %v %v \n ", ColorRed, t, ColorReset)
		fmt.Printf("Extra information: %v %v %v \n", ColorCyan, details, ColorReset)
		fmt.Printf("Visit %v https://github.com/YoussefKhalidAli/kubeRun/blob/master/errors.md %v for more details about this error \n", ColorGreen, ColorReset)
	}
	if strings.Contains(t, "H") {
		panic(err)
	}
}
