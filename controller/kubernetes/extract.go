package kubernetes

import (
	"errors"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kuberun.com/controller/utils"
)

func addService(svc corev1.ServiceSpec, metadata v1.ObjectMeta) {
	utils.Targets[svc.ClusterIP] = &utils.TargetDto{
		LastAccessed: time.Now(),
		ResourceName: metadata.Name,
		Namespace:    metadata.Namespace,
		Resource:     "placeholderRs",
	}

	utils.PrintTargets()
}

func deleteService(clusterIP string) {

	delete(utils.Targets, clusterIP)
	utils.PrintTargets()
	// TODO: Delte from agent config map
}

func ParseService(obj any, operation string) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		err := errors.New("Not a service")
		utils.HandelError(err, "KRC9030", "The returned object is not a kubernetes service")
	}

	isRun := filterAnnotations(svc.ObjectMeta.Annotations)

	if isRun && operation == "add" {
		addService(svc.Spec, svc.ObjectMeta)
	} else if operation == "delete" {
		deleteService(svc.Spec.ClusterIP)
	}
}

func filterAnnotations(anns map[string]string) bool {
	run := false
	if anns[utils.RunAnnotation] == "true" {
		run = true
	}
	return run
}
