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

	resource.Mux.Lock()
	resourceKind := resource.Resource
	name := resource.ResourceName
	namespace := resource.Namespace
	resource.Mux.Unlock()

	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: count,
		},
	}

	if resourceKind == "Deployment" {
		_, err := clientset.AppsV1().Deployments(resource.Namespace).UpdateScale(
			context.TODO(), name, scale, metav1.UpdateOptions{},
		)
		if err != nil {
			utils.HandelError(err, "KRC9060", fmt.Sprintf("Couldn't scale deployment %v", name))
		}
	}

	if count == 0 {
		go resource.Server.Start()
		PatchService(resource, count)
		resource.Mux.Lock()
		resource.Status = "Asleep"
		resource.Mux.Unlock()
	} else {
		waitForPodReady(resource)
		resource.Server.Proxy.Store("http://" + destIp[0])
		PatchService(resource, count)
		time.Sleep(5 * time.Second)
		resource.Server.Signal.Unlock()
		resource.Mux.Lock()
		resource.Status = "Awake"
		resource.Mux.Unlock()
	}
}

func PatchService(resource *store.TargetDto, count int32) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientset := GetClientset()

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

func waitForPodReady(resource *store.TargetDto) {
	clientset := GetClientset()

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
