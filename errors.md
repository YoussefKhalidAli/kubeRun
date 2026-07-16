# KubeRun Error Codes Reference

This document catalogs and explains the error codes that appear in KubeRun.

## Error Code Convention

To simplify debugging, every error code follows the convention: **`KRZYXXXS`**

* **`KR`**: Prefix for KubeRun.
* **`Z`**: Location where the error occurred:
  * **`A`**: Agent
  * **`C`**: Controller
* **`Y`**: Type of the error:
  * **`0`**: Administrative/Setup (misconfigurations)
  * **`1`**: Kubernetes connection / networking
  * **`9`**: Internal bug / system error
* **`XXX`**: Unique 3-digit error identifier.
* **`S`**: Severity tier:
  * **`H` (High / Panic)**: Unrecoverable failures. Mapped to `slog.Error` level and then halts execution (panics).
  * **`M` (Medium / Log & Continue)**: Recoverable failures. Mapped to `slog.Warn` level; execution resumes.
  * **`L` (Low / Suppressed)**: Expected, benign, or transient network conditions. Mapped to `slog.Debug` level (suppressed in default production configs but visible when `LOG_LEVEL=debug`).

## Logging Integration

All error handling via `utils.HandelError` is logged as a structured message with the following fields:
* `error_code`: The error code string (e.g. `KRA0012H`).
* `severity`: The severity character parsed from the code (`H`, `M`, or `L`).
* `details`: Additional contextual details about the error.
* `error`: The underlying error message (`err.Error()`).

---

## 1. Administrative / Setup Errors (Type `0`)

These errors indicate configuration or privilege issues on the host system or cluster.

| Error Code | Severity | Description |
| :--- | :---: | :--- |
| **KRA0012H** | High | The host operating system experienced severe memory starvation and could not allocate RAM to initialize core kernel monitoring structures. |
| **KRA0023H** | High | The host machine ran out of system file allocations (file descriptors) due to high open file density. |
| **KRA0024H** | High | The user account reached the maximum number of concurrent file-monitoring queues (`fsnotify`) allowed by the Linux kernel. |
| **KRA0403H** | High | The KubeRun Agent lacks required permissions (`CAP_NET_ADMIN`) to bind to Netlink conntrack multicast groups. |
| **KRA0404H** | High | The KubeRun Agent cannot find the required `nf_conntrack` kernel subsystem on the host node. |
| **KRC0404H** | High | The KubeRun Controller could not locate the required `agent-config` ConfigMap. Check the installation manifests. |
| **KRA0405H** | High | The KubeRun Agent configuration file was not found in `/etc/agent-config/config.yml`. |

---

## 2. Kubernetes Connection / Network Errors (Type `1`)

These errors occur when KubeRun interacts with the Kubernetes API server or makes networking requests to other cluster components.

| Error Code | Severity | Description |
| :--- | :---: | :--- |
| **KRC1440M** | Medium | The controller failed to update the `agent-config` ConfigMap. |
| **KRC1441M** | Medium | The controller failed to update or replace a Service definition in Kubernetes. |
| **KRC1442M** | Medium | The controller failed to resolve Agent hostnames via DNS lookup. |
| **KRC1443M** | Medium | The controller failed to update Deployment or StatefulSet labels after retrying. |
| **KRC1444M** | Medium | The controller failed to recreate a service during the traffic routing replacement. |
| **KRC1445M** | Medium | The controller failed to create the shadow service. |
| **KRC1446M** | Medium | The controller failed to send an HTTP POST update payload to the Agent. |
| **KRC1447M** | Medium | The controller failed to delete an old service during replacement. |
| **KRC1448L** | Low | An alert was received from an unknown IP address (benign/transient mismatch). |
| **KRC1448M** | Medium | The requested target IP/key was not found in the controller's active target registry. |
| **KRC1449M** | Medium | The controller failed to get a service definition from the Kubernetes API. |
| **KRC1450M** | Medium | The controller failed to delete the shadow service. |
| **KRC1451L** | Low | The controller failed to send an HTTP POST request to turn off/kill a dynamic traffic switch. |
| **KRC1452M** | Medium | The controller failed to update a Service's metadata labels after multiple retries. |
| **KRA1453M** | Medium | The Agent failed to send a traffic conntrack alert HTTP POST to the controller. |

---

## 3. KubeRun Internal Bug / System Errors (Type `9`)

These errors indicate internal software errors or unexpected environment failures during runtime.

| Error Code | Severity | Description |
| :--- | :---: | :--- |
| **KRA9010H** | High | Generic system initialization watcher failure at Agent startup. |
| **KRA9011H** | High | Unknown watcher error during Agent startup. |
| **KRA9012M** | Medium | `fsnotify` reported an error while monitoring the configuration file at runtime. |
| **KRA9013H** | High | The Agent failed to unmarshal its YAML configuration payload. |
| **KRA9014H** | High | The Agent's background Netlink connection listener failed to bind or subscribe to conntrack multicast group events. |
| **KRA9020H** | High | The Agent updates listener HTTP server failed to bind or start. |
| **KRA9021M** | Medium | The Agent failed to read or parse the HTTP request body in the update endpoint. |
| **KRA9022H** | High | Unexpected conntrack Dial failure (other than permission or module missing errors) when establishing the Netlink connection. |
| **KRC9010M** | Medium | The controller failed to read or parse the HTTP alert payload body. |
| **KRC9011H** | High | The controller alert HTTP server failed to bind or start. |
| **KRC9019M** | Medium | A dynamically spawned switch HTTP server failed to start. |
| **KRC9020H** | High | The controller failed to load the cluster configuration (In-Cluster or local kubeconfig). |
| **KRC9021H** | High | The controller failed to initialize the Kubernetes client clientset. |
| **KRC9022M** | Medium | The controller failed to verify a Deployment deletion. |
| **KRC9023M** | Medium | The controller failed to verify a StatefulSet deletion. |
| **KRC9040H** | High | The controller failed to unmarshal nested YAML data inside the `agent-config` ConfigMap. |
| **KRC9041M** | Medium | The controller encountered an unsupported mutation action in the ConfigMap updater. |
| **KRC9042H** | High | The controller failed to marshal updated configuration data back into YAML. |
| **KRC9060M** | Medium | The controller failed to scale a resource (Deployment or StatefulSet) via the scale API. |
| **KRC9061M** | Medium | The controller failed to list pods while waiting for them to transition to the Ready state. |
