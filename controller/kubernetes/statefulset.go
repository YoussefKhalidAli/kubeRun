package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"kuberun.com/controller/utils"
)

func LabelStatefulSet(ctx context.Context, resourceNamespace string, statefulset *v1.StatefulSet, clusterIP string) {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		latest, err := clientset.AppsV1().StatefulSets(resourceNamespace).Get(ctx, statefulset.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if latest.Labels == nil {
			latest.Labels = make(map[string]string)
		}
		latest.Labels["kuberun/clusterIP"] = clusterIP
		latest.Labels["kuberun/run"] = "true"

		_, updateErr := clientset.AppsV1().StatefulSets(resourceNamespace).Update(
			ctx,
			latest,
			metav1.UpdateOptions{},
		)
		return updateErr
	})
	if err != nil {
		utils.HandelError(err, "KRC1443", fmt.Sprintf("Couldn't update statefulset %v after retrying", statefulset.Name))
	}
}
