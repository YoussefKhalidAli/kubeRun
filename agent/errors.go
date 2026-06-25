package main

import (
	"fmt"
)

const (
	ColorReset = "\033[0m"
	ColorCyan  = "\033[36m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
)

func HandelError(err error, t string, details string) {
	fmt.Printf("Error: %v %v %v \n ", ColorRed, t, ColorReset)
	if details != "_" {
		fmt.Printf("Extra information: %v %v %v", ColorGreen, details, ColorReset)
	}
	fmt.Printf("Visit %v https://github.com/YoussefKhalidAli/kubeRun/blob/master/errors.md %v for more details about this error", ColorCyan, ColorReset)
	panic(err)
}
