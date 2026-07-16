## A table of `kuberun` labels. All labels are `kuberun/key: value`

|key|value(s)|purpose|
| -------- | ------- | ------- |
|operator|`controller`, `agent`, or `shadow`|Marks a `kuberun` operator |
|run|"true"|Tells `kuberun` controller to manage service and it's resource. |
|key|`ClusterIP` or `svc-<serviceName>`|Tells `kuberun` controller where this service's data is stored. |
|accessed|A date|Tells `kuberun` controller when this service was last accessed. |
|status|`Asleep` or `Awake`|Tells `kuberun` controller whether this service is at 0 or 1 replicas. |
