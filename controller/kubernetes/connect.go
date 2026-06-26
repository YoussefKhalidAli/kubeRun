package kubernetes

import "fmt"

func AddService(obj any) {
	fmt.Println("add", obj)
	println("Added-----------------------------------------------------------")
}

func DeleteService(obj any) {
	fmt.Println("del", obj)
	println("Deleted-----------------------------------------------------------")
}
