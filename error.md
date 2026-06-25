# This file outlines the error codes That appear inside KubeRun.

## Errors follow a convention for easier debugging. All error codes look like `KRZYXXX`, Where `Z` is where the error happened(Controller (C), Agent (A)), `Y` is the type of error, and `XXX` is the number.

---

## 1. Administrative errors:

### All administrative errors are caused by a misconfiguration in KubeRun setup and start with `KR0`

* **KRA0012:** his error surfaces when the host operating system is experiencing severe memory starvation and cannot allocate a small, non-swappable slice of RAM to initialize the core kernel monitoring structures.
* **KRA0023:** This error indicates that the entire host machine has completely run out of system file allocations due to a massive density of open files or network connections across all active containers.
* **KRA0024:** This error triggers when a single user account reaches the maximum number of concurrent file-monitoring queues allowed by the Linux kernel.
* **KRA0403:** This error indicates KubeRun doesn't have the permissions it needs to run. KubeRun tries its best to get all permissions it needs, it should work out the box. There might be something in your environment blocking it. KubeRun requiers the following permissions:
  * **Agent:**
    1. **CAP_NET_ADMIN:** The KubeRun agent requires this to bind to Netlink conntrack multicast groups. This is how KubeRun Knows when the last time a service was accessed in your environment and based on that, whether to scale to 0 or not.
* **KRA0404:** This error indicates KubeRun didn't find a necessary subsystem/process it needs to run. KubeRun assumes your system has:
  * **Agent:**
    1. **nf_conntrac:** The KubeRun agent requires this to estalish a connection to Netlink conntrack multicast groups. This is how KubeRun Knows when the last time a service was accessed in your environment and based on that, whether to scale to 0 or not.

---

## 2. KubeRun errors:

### These errors are probably caused by an issue/bug in the KubeRun source code. If you recieved an error starting with `KR9`, please [open an issue](https://github.com/YoussefKhalidAli/kubeRun/issues/new) with all the details.
