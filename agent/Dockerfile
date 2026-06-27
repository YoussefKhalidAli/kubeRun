FROM alpine:3.20

USER root

COPY kuberun-agent /usr/local/bin/kube-agent
RUN chmod +x /usr/local/bin/kube-agent

ENTRYPOINT ["/usr/local/bin/kube-agent"]
