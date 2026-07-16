# kubeRun Helm Chart

kubeRun is a scale-to-zero operator for Kubernetes that tracks connection state with conntrack and scales workloads down and back up.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

## Installing the Chart

### Option 1: From the Helm repository

Add the kubeRun Helm repository and install from it:

```bash
helm repo add kuberun https://youssefkhalidali.github.io/kubeRun/
helm repo update
helm install kuberun kuberun/kuberun
```

### Option 2: From a local checkout

To install the chart with the release name `kuberun` in the `default` namespace:

```bash
helm install kuberun ./charts/kuberun
```

If installing to a different namespace (e.g. `kuberun`), you must specify both `--namespace` and `--set namespace` to match:

```bash
helm install kuberun ./charts/kuberun --namespace kuberun --set namespace=kuberun
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `namespace` | string | `"default"` | Namespace to install resources. Must match helm install `--namespace` |
| `imagePullSecrets` | list | `[]` | Secrets for pulling docker images |
| `rbac.create` | bool | `true` | Create ClusterRole and ClusterRoleBinding resources |
| `controller.image.repository` | string | `"youssefkali/kuberun-controller"` | Controller image repository |
| `controller.image.tag` | string | `"v0.4.52"` | Controller image tag |
| `controller.image.pullPolicy` | string | `"Always"` | Controller image pull policy |
| `controller.syncMinutes` | int | `15` | Sync duration minutes. `store.SyncTime = syncMinutes * time.Minute / 2` |
| `controller.resources.requests.cpu` | string | `"50m"` | Controller requested CPU |
| `controller.resources.requests.memory` | string | `"64Mi"` | Controller requested memory |
| `controller.resources.limits.cpu` | string | `"200m"` | Controller CPU limit |
| `controller.resources.limits.memory` | string | `"128Mi"` | Controller memory limit (scales with cluster/resource count) |
| `controller.nodeSelector` | object | `{}` | Node selector for controller pod |
| `controller.tolerations` | list | `[]` | Tolerations for controller pod |
| `controller.affinity` | object | `{}` | Affinity for controller pod |
| `controller.podAnnotations` | object | `{}` | Extra annotations for controller pod |
| `controller.podLabels` | object | `{}` | Extra labels for controller pod |
| `controller.serviceAccount.create` | bool | `true` | Create ServiceAccount for the controller |
| `controller.serviceAccount.name` | string | `""` | Override controller ServiceAccount name |
| `controller.extraEnv` | list | `[]` | Extra environment variables to append to controller pod |
| `agent.image.repository` | string | `"youssefkali/kuberun-agent"` | Agent image repository |
| `agent.image.tag` | string | `"v0.3.0"` | Agent image tag |
| `agent.image.pullPolicy` | string | `"Always"` | Agent image pull policy |
| `agent.configMapName` | string | `"kuberun-agent-config"` | ConfigMap name for agent config |
| `agent.config.update` | bool | `false` | Enable GroupCTUpdate netlink group for established connection updates |
| `agent.resources.requests.cpu` | string | `"10m"` | Agent requested CPU |
| `agent.resources.requests.memory` | string | `"16Mi"` | Agent requested memory |
| `agent.resources.limits.cpu` | string | `"50m"` | Agent CPU limit |
| `agent.resources.limits.memory` | string | `"32Mi"` | Agent memory limit |
| `agent.extraTolerations` | list | `[]` | Additive tolerations for agent pods |
| `agent.nodeSelector` | object | `{}` | Node selector for agent pods |
| `agent.affinity` | object | `{}` | Affinity for agent pods |
| `agent.podAnnotations` | object | `{}` | Extra annotations for agent pods |
| `agent.podLabels` | object | `{}` | Extra labels for agent pods |
