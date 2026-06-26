package kubernetes

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"kuberun.com/controller/utils"
)

func AddService(obj any) {
	svc, _, _ := parseService(obj)

	fmt.Println("add", svc)
	fmt.Println("clusterIp", svc.ClusterIP)
	println("Added-----------------------------------------------------------")
}

func DeleteService(obj any) {
	svc, _, _ := parseService(obj)
	fmt.Println("del", svc)
	fmt.Println("clusterIp", svc.ClusterIP)
	println("Deleted-----------------------------------------------------------")
}

func parseService(obj any) (corev1.ServiceSpec, string, bool) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		err := errors.New("Not a service")
		utils.HandelError(err, "KRC9030", "The returned object is not a kubernetes service")
	}

	isRun := filterAnnotations(svc.ObjectMeta.Annotations)
	ns := svc.ObjectMeta.Namespace

	fmt.Println("IsRun", isRun)
	return svc.Spec, ns, isRun
}

func filterAnnotations(anns map[string]string) bool {
	run := false
	if anns[utils.RunAnnotation] == "true" {
		run = true
	}
	return run
}
