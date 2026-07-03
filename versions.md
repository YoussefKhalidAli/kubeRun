|operator|current version|
| -------- | ------- |
|controller| 0.3.2|
|agent| 0.2.0|

## Controller
- 0.3.2: 
    - Dynamically Add/Remove deployments exposed by managed services
      - Before: Controller needs the service o be created after or alongside registered services. If the service is created before the controller, it will not be able to find it.
      - Now: Controller partly registers the service and waits for it's exposed ddeployments before starting the scale proccess. Deleting the deployment resets service back to partially managed.
    - Handle clusterIP changes
      - Before: Controller only updates the data of the service.
      - Now: Controller updates the data of the service and also replaces the entire entry if the clusterIP of the service is changed

## Agent
- 0.2.0: Add instant ip updates
  - Before: Agent only updates the ip of the service when the agent cm is updated. This had a 1-2 minute delay.
  - Now: Agent updates the ip of the service instantly when it changes. This is done by sending a request from controller to all agents as soon as a new service is registered.
  - Note: Agent still primarly uses the cm to manage it's variables, the new `update` server just allows for instant updates untill the agent registers the updated cm.
