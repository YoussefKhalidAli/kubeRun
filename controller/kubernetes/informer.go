package kubernetes

import (
	"context"
	"flag"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	statefulsetInformer(factory)
	endpointSlicesInformer(factory)

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
			svc, ok := obj.(*corev1.Service)
			if !ok {
				panic("couldn't convert object to service")
			}
			key := svc.Spec.ClusterIP
			if key == "None" {
				key = GetHeadlessServiceKey(svc.ObjectMeta.Name)
			}
			DeleteTarget(clientset, key)
			store.PrintTargets()
		},
		UpdateFunc: func(old any, obj any) {
			svc := obj.(*corev1.Service)
			UpdateService(svc.Spec.ClusterIP, svc, old.(*corev1.Service), clientset)
		},
	})
}

func deploymentInformer(factory informers.SharedInformerFactory) {

	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj any) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			dep, ok := obj.(*appsv1.Deployment)
			depClusterIP := dep.Labels["kuberun/clusterIP"]
			if !ok {
				return
			}
			println("Deleting deployment with clusterIP: ", depClusterIP)
			current, err := clientset.AppsV1().Deployments(dep.Namespace).Get(
				ctx, dep.Name, metav1.GetOptions{},
			)

			switch {
			case errors.IsNotFound(err) || current.UID != dep.UID:
				DeleteResource(depClusterIP)
			case err != nil:
				utils.HandelError(err, "KRC9022", "Couldn't verify deployment deletion")
			default:
				LabelDeplyment(ctx, dep.Namespace, dep, depClusterIP)
			}
		},
	})
}

func statefulsetInformer(factory informers.SharedInformerFactory) {

	statefulsetInformer := factory.Apps().V1().StatefulSets().Informer()

	statefulsetInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj any) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sts, ok := obj.(*appsv1.StatefulSet)
			stsClusterIP := sts.Labels["kuberun/clusterIP"]
			if !ok {
				return
			}
			println("Deleting statefulset with clusterIP: ", stsClusterIP)
			current, err := clientset.AppsV1().StatefulSets(sts.Namespace).Get(
				ctx, sts.Name, metav1.GetOptions{},
			)

			switch {
			case errors.IsNotFound(err) || current.UID != sts.UID:
				DeleteResource(stsClusterIP)
			case err != nil:
				utils.HandelError(err, "KRC9022", "Couldn't verify deployment deletion")
			default:
				LabelStatefulSet(ctx, sts.Namespace, sts, stsClusterIP)
			}
		},
	})
}

func endpointSlicesInformer(factory informers.SharedInformerFactory) {

	endpointSlicesInformer := factory.Discovery().V1().EndpointSlices().Informer()

	endpointSlicesInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(_ any, obj any) {
			println("Got endpoint")
			slice := obj.(*discoveryv1.EndpointSlice)
			owner := slice.ObjectMeta.OwnerReferences[0].Name
			endpoints := slice.Endpoints
			AddSlice(owner, endpoints)
		},
	})
}

func GetClientset() *kubernetes.Clientset {
	return clientset
}
