|operator|current version|
| -------- | ------- |
|controller| 0.3.2|
|agent| 0.2.0|

## Controller
- 0.3.2: 
    - Dynamically Add/Remove deployments exposed by managed services
      - Before: Controller needs the service to be created after or alongside registered services. If the service is created before the controller, it will not be able to find it.
      - Now: Controller partly registers the service and waits for it's exposed ddeployments before starting the scale proccess. Deleting the deployment resets service back to partially managed.
    - Handle clusterIP changes
      - Before: Controller only updates the data of the service.
      - Now: Controller updates the data of the service and also replaces the entire entry if the clusterIP of the service is changed
- 0.4.0: 
    - Dynamically Add/Remove statefulsets exposed by managed services
      - Before: Controller only watches deployments
      - Now: Controller watches both deployments and statefulsets
- 0.4.1: 
    - Set custome key for statefulset services
      - Before: Statefulsets use headless services so there is not `clusterIP`, so the key became `None`.
      - Now: Sets the key to `svc-<svcName>` for headless services.
      - Note: This allows scaling statefulsets.
- 0.4.2: 
    - Watch endpoint slices created by statefulsets' headless services.
      - Before: No way to actually watch statefulsets' headless services since they don't have clusterIPs.
      - Now: Watch all pod IPs using the endpoint slice.
- 0.4.3: 
    - Add headless services map.
      - Before: No way to alert controller of a hit on headless service.
      - Now: controller updates `headless_map` in the agent config. The agent uses this to translate pod IPs to target keys when dealing with headless services.
- 0.4.4: 
    - Scale statefulsets.
      - Before: Only scale deployments.
      - Now: Scale statefulsets.
- 0.4.41: 
    - Replace service.
      - Before: Problems with directing headless services to `kuberun-controller` due to port mismatch.
      - Now: Replace headless services to normal ClusterIP services during scale to 0.
    - Add key label.
      - Before: Key was inferred based on `ClusterIP`, which causes issues when temperarly replacing headless services to clusterIP services.
      - Now: Add the key to the service data as a label `kuberun/key`.
- 0.4.42: 
    - Handle switch ports.
      - Before: First service gets port 4445, following services get +1. No way to track unused ports (Deleted services).
      - Now: First service gets port 200, following services get +1. If a service is deleted that port is released. Port `4444` is always reserved.
- 0.4.43: 
    - Wait for target creation before adding slice.
      - Before: 2 second wait before adding slice to check if target was created.
      - Now: Check if target was created before adding slice.
      - Note: The check tries to find target every second and times out after *5* seconds. 
- 0.4.5: 
    - Switch per port.
      - Before: 1 switch per service. This caused misdirection when using browsers to access k8s resources.
      - Now: 1 switch per port. This way, every switch know where to direct traffic without need to check req port.
---

## Agent
- 0.2.0: Add instant ip updates
    - Before: Agent only updates the ip of the service when the agent cm is updated. This had a 1-2 minute delay.
    - Now: Agent updates the ip of the service instantly when it changes. This is done by sending a request from controller to all agents as soon as a new service is registered.
    - Note: Agent still primarly uses the cm to manage it's variables, the new `update` server just allows for instant updates untill the agent registers the updated cm.
- 0.3.0: Add headless services map.
    - Before: Agent Sends hit IP to controller.
    - Now: Agent sends the key of the hit resource, be it a `ClusterIP` or the key of a headless service.
