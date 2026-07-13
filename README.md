Build a production-ready Helm chart for kubeRun at kuberun/charts. Base it on the existing raw
manifests in k8s/ (agent-config.yml, agent.yml, controller-permissions.yml, controller.yml) and
the env vars added in controller v0.4.52 (KUBERUN_NAMESPACE, KUBERUN_AGENT_CONFIG_NAME,
KUBERUN_AGENT_SERVICE_NAME, SYNC_MINUTES).

IMPORTANT CONSTRAINTS (do not violate these — the app breaks or becomes unsafe otherwise):

1. Single release per cluster. kubeRun's controller runs as a single replica with in-memory state
   and no leader election, and its ClusterRole/ClusterRoleBinding are cluster-scoped by design (it
   watches Services/Deployments/StatefulSets across ALL namespaces, not just its own). Do not add
   a replicaCount value — hardcode replicas: 1 in the Deployment template with the existing
   "DO NOT increase replicas" comment carried over from controller.yml.

2. Fixed resource names, decoupled from Helm release name. The controller Service must always be
   named exactly "kuberun-controller" and the agent Service must always be named exactly
   "kuberun-agent", regardless of .Release.Name or nameOverride/fullnameOverride. This is required
   because:
     a) controller/sync.go, controller/service/manage.go, and controller/informer/slice.go contain
        hardcoded string checks against the literal "kuberun-controller" to avoid the controller
        trying to track/scale itself. If the Service or Deployment name were templated to something
        else, the controller could attempt to scale itself, which is a serious failure mode. Do not
        templatize these Go-level checks or the resource names they depend on.
     b) The agent's Go code builds its controller endpoint from
        <service-name>.<namespace>.svc.cluster.local using KUBERUN_AGENT_SERVICE_NAME.
   nameOverride/fullnameOverride, if you add them at all, may only affect ancillary resources
   (ServiceAccount name, ConfigMap name if desired) — never the controller/agent Service or
   workload names.

3. Namespace is now dynamic via KUBERUN_NAMESPACE (added in controller v0.4.52), but the chart
   should still default all resources to install into "default" and should keep the
   ClusterRoleBinding's ServiceAccount subject namespace, the ConfigMap's data, and the env vars
   all deriving from ONE values.yml field (e.g. .Values.namespace), so they can never drift
   independently. Do not silently allow the chart to be installed into a different namespace than
   .Values.namespace resolves to — if the user sets a namespace via `--namespace` on the Helm CLI
   that differs from .Values.namespace, either use .Values.namespace consistently everywhere and
   document this clearly in NOTES.txt/helm.md, or add a `{{ fail }}` check in _helpers.tpl if they
   diverge. Pick whichever approach is more idiomatic Helm and explain your choice.

4. Wire the new v0.4.52 env vars into the controller Deployment template:
   - KUBERUN_NAMESPACE           <- .Values.namespace
   - KUBERUN_AGENT_CONFIG_NAME   <- .Values.agent.configMapName (default "kuberun-agent-config")
   - KUBERUN_AGENT_SERVICE_NAME  <- .Values.agent.serviceName (default "kuberun-agent", but per
                                     constraint 2 this should really just be a fixed literal
                                     "kuberun-agent" — expose it as a value only if you also update
                                     the ConfigMap/Service templates to consistently use the same
                                     value, otherwise hardcode it and explain why in a comment)
   Also update the ConfigMap's `kube_run_controller:` field to be templated from the same
   controller-service-name + namespace values used in constraint 2/3, instead of the hardcoded
   string that's currently in k8s/agent-config.yml.

5. Preserve the DaemonSet's required agent capabilities: hostNetwork: true, privileged: true,
   NET_ADMIN + SYS_ADMIN capabilities, and the control-plane/master tolerations, on by default.
   Expose agent.extraTolerations, agent.nodeSelector, agent.affinity as additive/overridable
   values without removing the required defaults.

Chart structure:
- Chart.yml (name: kuberun, description from README.md, appVersion "0.4.52", home/sources
  pointing at the GitHub repo, keywords: kubernetes, scale-to-zero, autoscaling, operator,
  conntrack)
- values.yml
- values.schema.json (basic validation: pullPolicy enum, syncMinutes positive integer, namespace
  as a non-empty string)
- templates/_helpers.tpl
- templates/serviceaccount.yml
- templates/clusterrole.yml
- templates/clusterrolebinding.yml
- templates/configmap-agent-config.yml
- templates/daemonset-agent.yml
- templates/service-agent.yml
- templates/deployment-controller.yml
- templates/service-controller.yml
- templates/NOTES.txt
- README.md (auto-generated values table, standard for Artifact Hub)
- .helmignore

values.yml should expose:

Shared:
- namespace (default "default") — single source of truth per constraint 3
- imagePullSecrets (default [])
- rbac.create (default true)

Controller:
- controller.image.repository (default youssefkali/kuberun-controller)
- controller.image.tag (default "v0.4.52")
- controller.image.pullPolicy (default Always)
- controller.syncMinutes (default 15) — in values.yml comments AND in helm.md, explain the real
  math: store.SyncTime = syncMinutes * time.Minute / 2, so syncMinutes=15 means workloads scale to
  zero after ~7.5 minutes of no traffic, not 15
- controller.resources.requests (default cpu: 50m, memory: 64Mi)
- controller.resources.limits (default cpu: 200m, memory: 128Mi) — add a values.yml comment that
  memory scales with cluster size / number of kuberun/run=true resources tracked, tune from
  observed usage
- controller.nodeSelector / tolerations / affinity (default {})
- controller.podAnnotations / podLabels (default {})
- controller.serviceAccount.create (default true) and .name override
- controller.extraEnv (list)

Agent:
- agent.image.repository (default youssefkali/kuberun-agent)
- agent.image.tag (default "v0.3.0")
- agent.image.pullPolicy (default Always)
- agent.config.update (default false) — document in values.yml comment: enables GroupCTUpdate
  netlink group, also fires on established-connection updates not just new connections
- agent.resources.requests (default cpu: 10m, memory: 16Mi)
- agent.resources.limits (default cpu: 50m, memory: 32Mi)
- agent.extraTolerations (default [], appended to the required baseline tolerations, never
  replacing them)
- agent.nodeSelector / affinity (default {})
- agent.podAnnotations / podLabels (default {})

After building the chart:
1. Run `helm lint ./charts/kuberun` and `helm template ./charts/kuberun` — fix all errors/warnings.
2. Run `helm template ./charts/kuberun --set namespace=foo` and confirm the ConfigMap's
   kube_run_controller field, the ClusterRoleBinding subject namespace, and all env vars correctly
   reflect "foo" consistently everywhere — paste this rendered output for me to review.
3. Diff the default-values rendered templates against the original k8s/*.yml manifests and call
   out any unintentional differences (labels, selectors, annotations).
4. Confirm resource requests/limits are attached to both the controller Deployment and agent
   DaemonSet containers.

Do not write helm.md yet — I'll ask for that separately once I've reviewed the chart output.
