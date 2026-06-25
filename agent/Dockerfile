FROM alpine:3.20

# Run as root because inotify and network monitoring require elevated host capabilities
USER root

# Copy the compiled static binary into the image
COPY kuberun-agent /usr/local/bin/kube-agent

# Set the binary as the startup command
ENTRYPOINT ["/usr/local/bin/kube-agent"]
