package deployment

import (
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"kuberun.com/controller/client"
	"kuberun.com/controller/utils"
)

func LabelDeplyment(ctx context.Context, resourceNamespace string, deployment *v1.Deployment, clusterIP string) {
	clientset := client.GetClientset()

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		latest, err := clientset.AppsV1().Deployments(resourceNamespace).Get(ctx, deployment.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if latest.Labels == nil {
			latest.Labels = make(map[string]string)
		}
		latest.Labels["kuberun/clusterIP"] = clusterIP
		latest.Labels["kuberun/run"] = "true"

		_, updateErr := clientset.AppsV1().Deployments(resourceNamespace).Update(
			ctx,
			latest,
			metav1.UpdateOptions{},
		)
		return updateErr
	})
	if err != nil {
		utils.HandelError(err, "KRC1443M", fmt.Sprintf("Couldn't update deployment %v after retrying", deployment.Name))
	}
}
