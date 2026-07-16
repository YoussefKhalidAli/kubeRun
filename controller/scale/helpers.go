package scale

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"kuberun.com/controller/client"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func PatchService(key string, count int32) {
	clientset := client.GetClientset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resource, ok := store.Targets[key]
	if !ok || resource == nil {
		utils.HandelError(fmt.Errorf("target %s not found in store", key), "KRC1448M", "target not found in store")
		return
	}

	svc, err := replaceService(key, clientset, ctx)
	if err != nil {
		utils.HandelError(err, "KRC1441M", fmt.Sprintf("Couldn't update service %v", resource.ServiceName))
		return
	}

	resource.Mux.Lock()
	resource.UpdateMarker = svc.ResourceVersion
	resource.Mux.Unlock()
}

func replaceService(key string, clientset *kubernetes.Clientset, ctx context.Context) (*corev1.Service, error) {
	var emptySvc *corev1.Service

	resource, ok := store.Targets[key]
	if !ok || resource == nil {
		return emptySvc, fmt.Errorf("target %s not found in store", key)
	}

	resource.Mux.Lock()
	name := resource.ServiceName
	namespace := resource.Namespace
	resource.Mux.Unlock()

	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		utils.HandelError(err, "KRC1449M", fmt.Sprintf("Couldn't find service %v", name))
		return emptySvc, err
	}

	if svc.Spec.Type == corev1.ServiceTypeExternalName {
		svc.Spec.Type = corev1.ServiceTypeClusterIP
		svc.Spec.ExternalName = ""

		if corev1.ServiceType(svc.ObjectMeta.Labels["kuberun/type"]) == corev1.ServiceTypeClusterIP {
			svc.Spec.ClusterIP = key
			svc.Spec.ClusterIPs = []string{key}
		} else {
			svc.Spec.ClusterIP = "None"
			svc.Spec.ClusterIPs = nil
		}

	} else {
		svc.Spec.Type = corev1.ServiceTypeExternalName
		svc.Spec.ClusterIP = ""
		svc.Spec.ClusterIPs = nil

		svc.Spec.IPFamilies = nil
		svc.Spec.IPFamilyPolicy = nil

		resource.Mux.Lock()
		shadowDNSName := fmt.Sprintf("%v.%v.svc.cluster.local", utils.GetShadowName(resource.ServiceName), store.KubeRunNamespace)
		resource.Mux.Unlock()

		svc.Spec.ExternalName = shadowDNSName
	}

	svc.ResourceVersion = ""
	svc.UID = ""
	svc.CreationTimestamp = metav1.Time{}
	svc.Generation = 0
	svc.Status = corev1.ServiceStatus{}

	err = clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		utils.HandelError(err, "KRC1447M", fmt.Sprintf("Couldn't delete old service %v", name))
		return emptySvc, err
	}

	created, err := clientset.CoreV1().Services(namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		utils.HandelError(err, "KRC1444M", fmt.Sprintf("Couldn't recreate service %v", name))
		return emptySvc, err
	}

	return created, nil
}

func WaitForPodReady(resource *store.TargetDto) string {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeout := time.After(time.Minute)

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
		select {
		case <-ticker.C:
			pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: labelSelector,
			})
			if err != nil {
				utils.HandelError(err, "KRC9061M", fmt.Sprintf("Couldn't get pods for %v", name))
				return ""
			}
			for _, pod := range pods.Items {
				for _, condition := range pod.Status.Conditions {
					if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
						return pod.Status.PodIP
					}
				}
			}
		case <-timeout:
			return ""
		}
	}
}
