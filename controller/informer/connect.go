package informer

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"kuberun.com/controller/client"
	"kuberun.com/controller/store"
)

func Connect() {
	clientset := client.Connect()

	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
		opts.LabelSelector = store.RunLabel
	}))

	serviceInformer(factory)
	deploymentInformer(factory)

	statefulsetInformer(factory)
	endpointSlicesInformer(factory)

	stopChan := make(chan struct{})
	defer close(stopChan)
	factory.Start(stopChan)
	println("we a go")
	<-stopChan

}
