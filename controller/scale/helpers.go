package scale

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
	"kuberun.com/controller/client"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func PatchService(resource *store.TargetDto, count int32) {
	clientset := client.GetClientset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resource.Mux.Lock()
	name := resource.ServiceName
	namespace := resource.Namespace
	resource.Mux.Unlock()

	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		utils.HandelError(err, "KRC1441", fmt.Sprintf("Couldn't find service %v", name))
	}

	if count == 0 {
		svc.Spec.Selector = map[string]string{
			"KubeRun": "Controller",
		}
		for index, _ := range svc.Spec.Ports {
			svc.Spec.Ports[index].TargetPort = intstr.FromInt(resource.Server.Port)
		}
	} else {
		svc.Spec.Selector = resource.SelectorMap
		servicePorts := *(resource.ServicePorts)
		for index, _ := range svc.Spec.Ports {
			svc.Spec.Ports[index].TargetPort = intstr.FromInt(servicePorts[index])
		}
	}

	resource.Mux.Lock()
	updated, err := clientset.CoreV1().Services(namespace).Update(
		ctx,
		svc,
		metav1.UpdateOptions{},
	)

	if err != nil {
		utils.HandelError(err, "KRC1441", fmt.Sprintf("Couldn't update service %v", name))
	}

	resource.UpdateMarker = updated.ResourceVersion
	resource.Mux.Unlock()
}

func WaitForPodReady(resource *store.TargetDto) string {
	clientset := client.GetClientset()

	resource.Mux.Lock()
	selectors := resource.SelectorMap
	name := resource.ResourceName
	namespace := resource.Namespace
	resource.Mux.Unlock()

	labels := []string{}
	for k, v := range selectors {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}
	labelSelector := strings.Join(labels, ",")

	for {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			utils.HandelError(err, "KRC9061", fmt.Sprintf("Couldn't get pods for %v", name))
			return ""
		}
		for _, pod := range pods.Items {
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
					return pod.Status.PodIP
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
