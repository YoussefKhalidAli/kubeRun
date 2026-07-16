package scale

import (
	"context"
	"fmt"
	"strings"
	"time"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kuberun.com/controller/client"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func ScaleResource(key string, count int32) {
	clientset := client.GetClientset()

	resource, ok := store.Targets[key]
	if !ok || resource == nil {
		utils.HandelError(fmt.Errorf("target %s not found in store", key), "KRC1448M", "target not found in store")
		return
	}
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

	var err error
	switch resourceKind {
	case "Deployment":
		_, err = clientset.AppsV1().Deployments(resource.Namespace).UpdateScale(
			context.TODO(), name, scale, metav1.UpdateOptions{},
		)
	case "StatefulSet":
		_, err = clientset.AppsV1().StatefulSets(resource.Namespace).UpdateScale(
			context.TODO(), name, scale, metav1.UpdateOptions{})
	}

	if err != nil {
		utils.HandelError(err, "KRC9060M", fmt.Sprintf("Couldn't scale deployment %v", name))
		return
	}

	if count == 0 {
		for _, server := range resource.Servers {
			go server.Start()
		}

		PatchService(key, count)
		resource.Mux.Lock()
		resource.Status = "Asleep"
		resource.Mux.Unlock()
	} else {
		podIP := WaitForPodReady(resource)

		if strings.Contains(key, "svc-") {
			for _, server := range resource.Servers {
				server.Proxy.Store("http://" + podIP)
			}

		} else if podIP != "" {
			for _, server := range resource.Servers {
				server.Proxy.Store("http://" + key)
			}
		}

		PatchService(key, count)
		time.Sleep(time.Second)
		for _, server := range resource.Servers {
			server.Signal.Unlock()
		}
		resource.Mux.Lock()
		resource.Status = "Awake"
		resource.Mux.Unlock()
	}
}
