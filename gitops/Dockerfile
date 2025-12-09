# Use the official ArgoCD image so we get the correct 'argocd-cmp-server' binary
FROM quay.io/argoproj/argocd:v3.2.1

# Switch to root to install tools
USER root

# Install SOPS
# (We use a specific version to ensure stability)
RUN apt-get update && \
    apt-get install -y wget && \
    wget -O /usr/local/bin/sops https://github.com/getsops/sops/releases/download/v3.8.1/sops-v3.8.1.linux.amd64 && \
    chmod +x /usr/local/bin/sops && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Switch back to the argocd user
USER 999

# The entrypoint is already set correctly in the base image, 
# but we need to run the cmp-server specifically for the sidecar
ENTRYPOINT [ "/var/run/argocd/argocd-cmp-server" ]
