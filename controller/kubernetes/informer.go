package kubernetes

import (
	"flag"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

var clientset *kubernetes.Clientset

func connect() {
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

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = store.RunLabel
	}))

	serviceInformer(factory)
	deploymentInformer(factory)

	stopChan := make(chan struct{})
	defer close(stopChan)
	factory.Start(stopChan)
	println("we a go")
	<-stopChan

}

func serviceInformer(factory informers.SharedInformerFactory) {

	serviceInformer := factory.Core().V1().Services().Informer()

	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			svc := obj.(*corev1.Service)
			AddService(svc.Spec, svc.ObjectMeta, clientset)
		},
		DeleteFunc: func(obj any) {
			svc := obj.(*corev1.Service)
			DeleteService(svc.Spec.ClusterIP, clientset)
		},
		UpdateFunc: func(_ any, obj any) {
			svc := obj.(*corev1.Service)
			UpdateService(svc.Spec.ClusterIP, svc)
		},
	})
}

func deploymentInformer(factory informers.SharedInformerFactory) {

	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj any) {
			depClusterIP := obj.(*metav1.ObjectMeta).Labels["kuberun/clusterIP"]
			DeleteResource(depClusterIP)
		},
	})
}

func GetClientset() *kubernetes.Clientset {
	return clientset
}
