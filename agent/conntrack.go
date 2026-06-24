package main

func main() {
	// kubeRunNamespace := os.Getenv("NAMESPACE")

	c := KernelListener()
	defer c.Close()
	Filter()
}
