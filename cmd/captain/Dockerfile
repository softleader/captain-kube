FROM ubuntu:xenial

RUN apt-get update && \
    apt-get install -y curl && \
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
