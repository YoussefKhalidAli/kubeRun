package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func ScaleResource(resource *store.TargetDto, count int32, destIp ...string) {
	clientset := GetClientset()
	fmt.Printf("Scaling %v to %v", resource, count)

	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.ResourceName,
			Namespace: resource.Namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: count,
		},
	}

	patchService(resource, count)

	if resource.Resource == "Deployment" {
		_, err := clientset.AppsV1().Deployments(resource.Namespace).UpdateScale(
			context.TODO(), resource.ResourceName, scale, metav1.UpdateOptions{},
		)
		if err != nil {
			utils.HandelError(err, "KRC9060", fmt.Sprintf("Couldn't scale deployment %v", resource.ResourceName))
		}
	}
	if count == 0 {
		resource.IsSleep = true
		go resource.Server.Start()
	} else {
		waitForPodReady(resource)
		resource.IsSleep = false
		resource.Server.Proxy = destIp[0]
		resource.Server.Signal.Unlock()
	}
}

func patchService(resource *store.TargetDto, count int32) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientset := GetClientset()

	svc, err := clientset.CoreV1().Services(resource.Namespace).Get(ctx, resource.ServiceName, metav1.GetOptions{})
	if err != nil {
		utils.HandelError(err, "KRC1441", fmt.Sprintf("Couldn't find service %v", resource.ServiceName))
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

	_, err = clientset.CoreV1().Services(resource.Namespace).Update(
		ctx,
		svc,
		metav1.UpdateOptions{},
	)
	if err != nil {
		utils.HandelError(err, "KRC1441", fmt.Sprintf("Couldn't update service %v", resource.ServiceName))
	}
}

func waitForPodReady(resource *store.TargetDto) {
	clientset := GetClientset()

	labels := []string{}
	for k, v := range resource.SelectorMap {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}
	labelSelector := strings.Join(labels, ",")

	for {
		pods, err := clientset.CoreV1().Pods(resource.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			utils.HandelError(err, "KRC9061", fmt.Sprintf("Couldn't get pods for %v", resource.ResourceName))
			return
		}
		for _, pod := range pods.Items {
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
					return
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
