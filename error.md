# This file outlines the error codes That appear inside KubeRun.

## Errors follow a convention for easier debugging. All error codes look like `KRYXXX`, Where `Y` is the type of error and `XXX` is the number.

---

## 1. Administrative errors:

### All administrative errors are caused by a misconfiguration in KubeRun setup and start with `KR0`

* **KR0403:** This error indicates KubeRun doesn't have the permissions it needs to run. KubeRun tries its best to get all permissions it needs, it should work out the box. There might be something in your environment blocking it. KubeRun requiers the following permissions:
  * **Agent:**
    1. **CAP_NET_ADMIN:** The KubeRun agent requires this to bind to Netlink conntrack multicast groups. This is how KubeRun Knows when the last time a service was accessed in your environment and based on that, whether to scale to 0 or not.
* **KR0404:** This error indicates KubeRun didn't find a necessary subsystem/process it needs to run. KubeRun assumes your system has:
  * **Agent:**
    1. **nf_conntrac:** The KubeRun agent requires this to estalish a connection to Netlink conntrack multicast groups. This is how KubeRun Knows when the last time a service was accessed in your environment and based on that, whether to scale to 0 or not.
