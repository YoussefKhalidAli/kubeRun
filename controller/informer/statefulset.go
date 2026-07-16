package informer

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"kuberun.com/controller/client"
	"kuberun.com/controller/resource"
	"kuberun.com/controller/statefulset"
	"kuberun.com/controller/utils"
)

func statefulsetInformer(factory informers.SharedInformerFactory) {
	clientset := client.GetClientset()

	statefulsetInformer := factory.Apps().V1().StatefulSets().Informer()

	statefulsetInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj any) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sts, ok := obj.(*appsv1.StatefulSet)
			if !ok || sts == nil {
				return
			}
			stsClusterIP := ""
			if sts.Labels != nil {
				stsClusterIP = sts.Labels["kuberun/clusterIP"]
			}
			println("Deleting statefulset with clusterIP: ", stsClusterIP)
			current, err := clientset.AppsV1().StatefulSets(sts.Namespace).Get(
				ctx, sts.Name, metav1.GetOptions{},
			)

			switch {
			case errors.IsNotFound(err) || current.UID != sts.UID:
				resource.DeleteResource(stsClusterIP)
			case err != nil:
				utils.HandelError(err, "KRC9023M", "Couldn't verify statefulset deletion")
			default:
				statefulset.LabelStatefulSet(ctx, sts.Namespace, sts, stsClusterIP)
			}
		},
	})
}
