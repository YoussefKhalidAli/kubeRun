FROM alpine:3.20

USER root

COPY kuberun-controller /usr/local/bin/kube-controller
RUN chmod +x /usr/local/bin/kube-controller

ENTRYPOINT ["/usr/local/bin/kube-controller"]
