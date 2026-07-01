package kubernetes

import (
	"context"
	"net"
	"net/http"
	"slices"
	"strings"

	"go.yaml.in/yaml/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kuberun.com/controller/store"
	"kuberun.com/controller/utils"
)

func UpdateAgentCM(clientset *kubernetes.Clientset, targetIP string, action string) error {
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

func UpdateAgents(ip string) {
	endpoints, err := net.LookupHost(store.KubeRunAgent)
	if err != nil {
		utils.HandelError(err, "KRC1442", "failed to find agents")
	}
	for _, endpoint := range endpoints {
		_, err := http.Post("http://"+endpoint+":4443/update", "text/plain", strings.NewReader(ip))
		if err != nil {
			utils.HandelError(err, "KRC1442", "failed to update agents")
		}
	}

}
