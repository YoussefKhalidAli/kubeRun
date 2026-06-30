package kubernetes

func Kubernetes() {
	go SyncLoop()
	connect()
}
