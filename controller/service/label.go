package service

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"kuberun.com/controller/store"
)

func labelService(clientset *kubernetes.Clientset, name string, namespace string, key string, labelKey []string, labelValue []string) error {

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		ctx := context.Background()
		latest, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if latest.Labels == nil {
			latest.Labels = make(map[string]string)
		}

		for index, _ := range labelKey {
			latest.Labels[fmt.Sprintf("kuberun/%v", labelKey[index])] = labelValue[index]
		}

		updated, updateErr := clientset.CoreV1().Services(namespace).Update(
			ctx,
			latest,
			metav1.UpdateOptions{},
		)
		if updateErr != nil {
			return updateErr
		}

		target, ok := store.Targets[key]
		if ok && target != nil {
			target.Mux.Lock()
			target.UpdateMarker = updated.ResourceVersion
			target.Mux.Unlock()
		}

		return nil
	})
	return err
}
