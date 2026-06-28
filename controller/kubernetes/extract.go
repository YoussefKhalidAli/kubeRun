package kubernetes

import (
	"context"
	"errors"
	"slices"
	"time"

	"go.yaml.in/yaml/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/server"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func addService(svc corev1.ServiceSpec, metadata metav1.ObjectMeta, clientset kubernetes.Interface) {
	resourceName, resource := FindResource(clientset, svc.Selector, metadata.Namespace)
	store.Targets[svc.ClusterIP] = &store.TargetDto{
		LastAccessed: time.Now(),
		ResourceName: resourceName,
		Namespace:    metadata.Namespace,
		Resource:     resource,
		ServiceName:  metadata.Name,
		Server:       server.New(),
		IsSleep:      false,
		ServicePorts: mapServicePorts(svc.Ports),
		SelectorMap:  svc.Selector,
	}
	updateAgentCM(clientset, svc.ClusterIP, "add")

	store.PrintTargets()
}

func deleteService(clusterIP string, clientset kubernetes.Interface) {

	updateAgentCM(clientset, clusterIP, "delete")
	delete(store.Targets, clusterIP)
	store.PrintTargets()
}

func ParseService(clientset kubernetes.Interface, obj any, operation string) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		err := errors.New("Not a service")
		utils.HandelError(err, "KRC9030", "The returned object is not a kubernetes service")
	}

	isRun := filterAnnotations(svc.ObjectMeta.Annotations)

	if isRun && operation == "add" {
		addService(svc.Spec, svc.ObjectMeta, clientset)
	} else if operation == "delete" {
		deleteService(svc.Spec.ClusterIP, clientset)
	}
}

func filterAnnotations(anns map[string]string) bool {
	run := false
	if anns[store.RunAnnotation] == "true" {
		run = true
	}
	return run
}

func mapServicePorts(portsMap []corev1.ServicePort) *[]int {
	targetPortsMap := make([]int, len(portsMap))
	for index, port := range portsMap {
		targetPortsMap[index] = port.TargetPort.IntValue()
	}
	return &targetPortsMap
}

func updateAgentCM(clientset kubernetes.Interface, targetIP string, action string) error {
	ctx := context.Background()

	cm, err := clientset.CoreV1().ConfigMaps(store.KubeRunNamespace).Get(ctx, store.KubeRunAgentConfigName, metav1.GetOptions{})
	if err != nil {
		utils.HandelError(err, "KRC0404", "failed to get configmap")
	}

	if cm.Data == nil || cm.Data["config.yml"] == "" {
		utils.HandelError(err, "KRC0404", "configmap data or config.yml key is missing")
	}

	var innerConfig store.AgentConfig
	err = yaml.Unmarshal([]byte(cm.Data["config.yml"]), &innerConfig)
	if err != nil {
		utils.HandelError(err, "KRC9040", "failed to unmarshal nested yml payload")
	}

	switch action {
	case "add":
		if !slices.Contains(innerConfig.Ips, targetIP) {
			innerConfig.Ips = append(innerConfig.Ips, targetIP)
		}
	case "delete":
		innerConfig.Ips = slices.DeleteFunc(innerConfig.Ips, func(ip string) bool {
			return ip == targetIP
		})
	default:
		utils.HandelError(err, "KRC9041", "unsupported mutation action")
	}

	updatedBytes, err := yaml.Marshal(&innerConfig)
	if err != nil {
		utils.HandelError(err, "KRC9042", "failed to marshal updated config payload")
	}

	cm.Data["config.yml"] = string(updatedBytes)

	_, err = clientset.CoreV1().ConfigMaps(store.KubeRunNamespace).Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		utils.HandelError(err, "KRC1440", "failed to marshal updated config payload")
	}

	return nil
}
