# kubeRun

kubeRun is a lightweight, zero-configuration **scale-to-zero** operator for Kubernetes. It automatically scales inactive workloads (Deployments and StatefulSets) down to 0 replicas to conserve resources and wakes them up instantly when new traffic is received.

## Why kubeRun?

In cloud-native development and testing environments, idle workloads waste massive amounts of CPU and memory. While scale-to-zero solutions exist, they often present two major drawbacks:
1. **Overkill Feature Sets:** Popular alternatives (such as Knative Serving or KEDA) are extremely heavy and full of features when all you need is a simple, reliable 0 → 1 → 0 scaling loop.
2. **Annoying Dev/Ops Overhead:** Existing tools typically require configuring specific API gateways or ingress layers (like Envoy, Istio, or Contour) and managing custom manifest formats. 

**kubeRun** was built to be transparent and lightweight. It operates directly on native Kubernetes Services, Deployments, and StatefulSets without requiring any specific gateways or custom resource definitions (CRDs).

---

## Architecture & How It Works

kubeRun is composed of two primary components:

1. **kubeRun Controller**: 
   - Runs as a standard Kubernetes Deployment.
   - Watches cluster services and matches them to their backing Deployments or StatefulSets.
   - Handles the idle detection and scales workloads down to 0 replicas when there is no traffic.
2. **kubeRun Agent**: 
   - Runs as a lightweight `DaemonSet` on every node in the cluster.
   - Operates with `NET_ADMIN` and `hostNetwork: true` to bind to Linux Netlink conntrack multicast groups (`nf_conntrack`).
   - Instantly detects new connection attempts to tracked service IPs at the kernel level and alerts the controller to scale the corresponding workload back up.

---

## Documentation

For detailed information on troubleshooting and version changes, see:
- [Error Codes Reference](error.md) — Comprehensive guide to troubleshooting administrative, Kubernetes connection, and internal KubeRun errors.
- [Version History & Changelog](versions.md) — Version tracking for the controller and agent.

---

## Upcoming Roadmap

Here are the features we are working on for upcoming releases:
* **Better Logging:** Structured, searchable, and configurable log outputs for production environments.
* **Better Error Handling:** Graceful recovery and reporting instead of the current development behavior (which panics on everything to ease debugging).
* **High Availability & Stateless Controller** Adding leader election for multi-replica setups and migrating in-memory tracking to Kubernetes resource annotations to ensure zero-downtime restarts.
* **Cilium CNI Support:** Native compatibility with Cilium-powered networking.
* **Helm Chart Deployment:** Bundling the entire system (permissions, configmaps, agent daemonset, and controller deployment) into a single, configurable Helm chart for easy installations.
