package informer

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"kuberun.com/controller/client"
	"kuberun.com/controller/deployment"
	"kuberun.com/controller/resource"
	"kuberun.com/controller/utils"
)

func deploymentInformer(factory informers.SharedInformerFactory) {
	clientset := client.GetClientset()

	deploymentInformer := factory.Apps().V1().Deployments().Informer()

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj any) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			dep, ok := obj.(*appsv1.Deployment)
			if !ok || dep == nil {
				return
			}
			depClusterIP := ""
			if dep.Labels != nil {
				depClusterIP = dep.Labels["kuberun/clusterIP"]
			}
			logger.Info("deleting deployment", "cluster_ip", depClusterIP)
			current, err := clientset.AppsV1().Deployments(dep.Namespace).Get(
				ctx, dep.Name, metav1.GetOptions{},
			)

			switch {
			case errors.IsNotFound(err) || current.UID != dep.UID:
				resource.DeleteResource(depClusterIP)
			case err != nil:
				utils.HandelError(err, "KRC9022M", "Couldn't verify deployment deletion")
			default:
				deployment.LabelDeplyment(ctx, dep.Namespace, dep, depClusterIP)
			}
		},
	})
}
