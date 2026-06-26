package kubernetes

import (
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kuberun.com/controller/utils"
)

func AddService(obj any) {
	svc, metadata, isRun := parseService(obj)

	if !isRun {
		return
	}

	utils.Targets[svc.ClusterIP] = &utils.TargetDto{
		LastAccessed: time.Now(),
		ResourceName: metadata.Name,
		Namespace:    metadata.Namespace,
		Resource:     "placeholderRs",
	}
}

func DeleteService(obj any) {
	svc, _, isRun := parseService(obj)

	if !isRun {
		delete(utils.Targets, svc.ClusterIP)
	}
	// TODO: Delte from agent config map
}

func parseService(obj any) (corev1.ServiceSpec, v1.ObjectMeta, bool) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		err := errors.New("Not a service")
		utils.HandelError(err, "KRC9030", "The returned object is not a kubernetes service")
	}

	isRun := filterAnnotations(svc.ObjectMeta.Annotations)

	fmt.Println("IsRun", isRun)
	return svc.Spec, svc.ObjectMeta, isRun
}

func filterAnnotations(anns map[string]string) bool {
	run := false
	if anns[utils.RunAnnotation] == "true" {
		run = true
	}
	return run
}
