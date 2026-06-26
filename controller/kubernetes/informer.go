package kubernetes

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"kuberun.com/controller/utils"
)

func connect() {
	config, err := rest.InClusterConfig()
	if err != nil {
		home, _ := os.UserHomeDir()
		var kubeconfig *string = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			utils.HandelError(err, "KRC9020", "Controller couldn't create cluster config.")
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.HandelError(err, "KRC9021", "Controller couldn't establish clientset")
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Duration(utils.SyncTime)*time.Minute-time.Minute)

	serviceInformer := factory.Core().V1().Services().Informer()

	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: func(obj any) { ParseService(clientset, obj, "add") }, DeleteFunc: func(obj any) { ParseService(clientset, obj, "delete") }})
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	println("we a go")
	<-stopCh
}
