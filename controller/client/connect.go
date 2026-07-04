package client

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"kuberun.com/controller/utils"
)

var clientset *kubernetes.Clientset

func Connect() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		// home, _ := os.UserHomeDir()
		var kubeconfig *string = flag.String("kubeconfig", filepath.Join("/home/ubuntu/", ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			utils.HandelError(err, "KRC9020", "Controller couldn't create cluster config.")
		}
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		utils.HandelError(err, "KRC9021", "Controller couldn't establish clientset")
	}

	return clientset
}

func GetClientset() *kubernetes.Clientset {
	return clientset
}
