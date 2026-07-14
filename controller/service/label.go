package service

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func labelService(clientset *kubernetes.Clientset, name string, namespace string, key string, labelKey string, labelValue string) {

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		ctx := context.Background()
		latest, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if latest.Labels == nil {
			latest.Labels = make(map[string]string)
		}

		latest.Labels[labelKey] = labelValue

		updated, updateErr := clientset.CoreV1().Services(namespace).Update(
			ctx,
			latest,
			metav1.UpdateOptions{},
		)

		target := store.Targets[key]
		target.Mux.Lock()
		target.UpdateMarker = updated.ResourceVersion
		target.Mux.Unlock()

		return updateErr
	})
	if err != nil {
		utils.HandelError(err, "KRC1443", fmt.Sprintf("Couldn't update statefulset %v after retrying", name))
	}
}
