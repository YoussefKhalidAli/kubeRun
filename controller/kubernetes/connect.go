package kubernetes

import (
	"flag"
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
		var kubeconfig *string = flag.String("kubeconfig", filepath.Join("/home/usef", ".kube", "config"), "(optional) absolute path to the kubeconfig file")
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

	factory := informers.NewSharedInformerFactory(clientset, time.Duration(utils.SyncTime*600))

	serviceInformer := factory.Core().V1().Services().Informer()

	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{AddFunc: AddService, DeleteFunc: DeleteService})
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	println("we a go")
	<-stopCh
}
