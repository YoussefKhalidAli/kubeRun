package kubernetes

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func addService(svc corev1.ServiceSpec, metadata metav1.ObjectMeta, clientset kubernetes.Interface) {
	resourceName, resource := FindResource(clientset, svc.Selector, metadata.Namespace)
	if resourceName == "kuberun-controller" || resource == "DaemonSet" {
		println("Found unmanagable resource. Skipping")
		return
	}

	CreateService(svc, metadata, resourceName, resource)

	UpdateAgentCM(clientset, svc.ClusterIP, "add")

	store.PrintTargets()
}

func deleteService(clusterIP string, clientset kubernetes.Interface) {

	UpdateAgentCM(clientset, clusterIP, "delete")
	store.Targets[clusterIP].Server.Stop()
	delete(store.Targets, clusterIP)

	store.PrintTargets()
}

func ParseService(clientset kubernetes.Interface, obj any, operation string) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		err := errors.New("Not a service")
		utils.HandelError(err, "KRC9030", "The returned object is not a kubernetes service")
	}

	isRun := FilterAnnotations(svc.ObjectMeta.Annotations)

	if isRun && operation == "add" {
		addService(svc.Spec, svc.ObjectMeta, clientset)
	} else if operation == "delete" {
		deleteService(svc.Spec.ClusterIP, clientset)
	}
}
