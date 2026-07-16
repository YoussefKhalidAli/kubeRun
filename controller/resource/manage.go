package resource

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/deployment"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

var logger = utils.Logger.With("module", "resource")

func FindResource(clientset *kubernetes.Clientset, selectorMap map[string]string, resourceNamespace string, clusterIP string) (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2)
	defer cancel()

	labelSet := labels.Set(selectorMap)
	selector := labelSet.AsSelector()

	listOptions := metav1.ListOptions{
		LabelSelector: selector.String(),
	}

	var resourceName, resourceKind string = "", ""

	wait.PollUntilContextTimeout(ctx, 2*time.Second, time.Minute/2, true, func(pollCtx context.Context) (bool, error) {
		pods, err := clientset.CoreV1().Pods(resourceNamespace).List(pollCtx, listOptions)
		if err != nil {
			return false, nil
		}

		if len(pods.Items) == 0 {
			return false, nil
		}

		targetPod := pods.Items[0]
		podOwner := targetPod.ObjectMeta.OwnerReferences[0]

		if podOwner.Kind == "ReplicaSet" {
			replicaSet, err := clientset.AppsV1().ReplicaSets(resourceNamespace).Get(pollCtx, podOwner.Name, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}

			if len(replicaSet.ObjectMeta.OwnerReferences) > 0 {
				dep, err := clientset.AppsV1().Deployments(resourceNamespace).Get(pollCtx, replicaSet.ObjectMeta.OwnerReferences[0].Name, metav1.GetOptions{})
				if err != nil {
					return false, nil
				}
				deployment.LabelDeplyment(ctx, resourceNamespace, dep, clusterIP)
				resourceName = replicaSet.ObjectMeta.OwnerReferences[0].Name
				resourceKind = replicaSet.ObjectMeta.OwnerReferences[0].Kind
			} else {
				resourceName = podOwner.Name
				resourceKind = podOwner.Kind
			}
		} else {
			resourceName = podOwner.Name
			resourceKind = podOwner.Kind
		}

		return true, nil
	})

	return resourceName, resourceKind
}

func DeleteResource(clusterIP string) {
	target, ok := store.Targets[clusterIP]
	if !ok || target == nil {
		return
	}
	target.Mux.Lock()
	target.Resource = ""
	target.ResourceName = ""

	if target.Status == "Asleep" {
		for _, server := range target.Servers {
			server.Kill()
		}
	}

	target.Mux.Unlock()
	logger.Info("deleted resource", "cluster_ip", clusterIP)
}
