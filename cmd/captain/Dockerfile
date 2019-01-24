FROM ubuntu:bionic

RUN apt-get update && \
    apt-get install -y gnupg2 curl apt-transport-https && \
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
    echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list && \
    apt-get update && \
    apt-get install -y kubectl && \
    apt-get autoclean && \
	curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash && \
	helm init -c && \
	curl -fsSL https://clis.ng.bluemix.net/install/linux | bash

COPY icp-linux-amd64.tgz /tmp
RUN tar -C /tmp -xzf /tmp/icp-linux-amd64.tgz && \
	bx plugin install /tmp/icp-linux-amd64 && \
	rm -rf /tmp/*

WORKDIR /app
ADD captain /app/captain

ENTRYPOINT ["./captain"]
