FROM ubuntu:bionic

ARG GRPC_HEALTH_PROBE_VERSION=v0.3.1

RUN apt-get update && \
    apt-get install -y gnupg2 curl wget apt-transport-https && \
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list && \
    apt-get update && \
    apt-get install -y kubectl && \
    apt-get autoclean && \
	curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash && \
	helm init -c && \
	curl -fsSL https://clis.ng.bluemix.net/install/linux | bash && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY icp-linux-amd64.tgz /tmp
RUN tar -C /tmp -xzf /tmp/icp-linux-amd64.tgz && \
	bx plugin install /tmp/icp-linux-amd64 && \
	rm -rf /tmp/*

WORKDIR /app
ADD captain /app/captain

ENTRYPOINT ["./captain"]
