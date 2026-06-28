package kubernetes

import (
	"context"
	"fmt"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func ScaleResource(resource *store.TargetDto, count int32) {
	clientset := GetClientset()

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
		resource.IsSleep = false
		resource.Server.Signal <- true
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
			svc.Spec.Ports[index].TargetPort = intstr.FromInt(4444)
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
