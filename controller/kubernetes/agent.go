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

func UpdateAgentCM(clientset *kubernetes.Clientset, targetIP string, action string, targetIPs ...string) error {
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
		if strings.Contains(targetIP, "svc-") {
			innerConfig.Ips = append(innerConfig.Ips, targetIPs...)
			addHeadlessSet(targetIP, &innerConfig, targetIPs)
		} else {
			innerConfig.Ips = append(innerConfig.Ips, targetIP)
		}
	case "delete":
		innerConfig.Ips = slices.DeleteFunc(innerConfig.Ips, func(ip string) bool {
			return ip == targetIP
		})

	default:
		utils.HandelError(err, "KRC9041", "unsupported mutation action")
	}

	innerConfig.Ips = uniqueElements(innerConfig.Ips)
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

func addHeadlessSet(name string, config *store.AgentConfig, ips []string) {
	if config.HeadlessMap == nil {
		config.HeadlessMap = make(map[string]string)
	}
	for _, ip := range ips {
		config.HeadlessMap[ip] = name
	}
}

func uniqueElements(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, val := range slice {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}
