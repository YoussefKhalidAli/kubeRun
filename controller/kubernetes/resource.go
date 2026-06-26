package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func FindResource(clientset kubernetes.Interface, selectorMap map[string]string, resourceNamespace string) (string, string) {
	ctx := context.Background()

	labelSet := labels.Set(selectorMap)
	selector := labelSet.AsSelector()

	listOptions := metav1.ListOptions{
		LabelSelector: selector.String(),
	}

	pods, err := clientset.CoreV1().Pods(resourceNamespace).List(ctx, listOptions)
	if err != nil {
		// idk
	}

	podOwner := pods.Items[0].ObjectMeta.OwnerReferences[0]
	var resource metav1.OwnerReference
	if podOwner.Kind == "ReplicaSet" {
		replicaSet, err := clientset.AppsV1().ReplicaSets(resourceNamespace).Get(ctx, podOwner.Name, metav1.GetOptions{})
		if err != nil {
			// later
		}
		resource = replicaSet.ObjectMeta.OwnerReferences[0]
	} else {
		resource = podOwner
	}

	if err != nil {
		// idk
	}
	return resource.Name, resource.Kind
}
