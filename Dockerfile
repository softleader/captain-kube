FROM alpine
MAINTAINER softleader.com.tw

RUN apk update && \
	apk --no-cache add \
		bash \
		vim \
		tree \
		curl \
		procps \
		tzdata \
		ansible \
		python \
		openssh \
		openssl \
		sshpass \
	&& rm -rf /var/cache/apk/* && \
	ls /usr/share/zoneinfo && \
	cp /usr/share/zoneinfo/Asia/Taipei /etc/localtime && \
	echo "Asia/Taipei" > /etc/timezone && \
	curl -L https://github.com/docker/compose/releases/download/1.21.0/docker-compose-$(uname -s)-$(uname -m) -o /docker-compose && \
	curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash

ENV CAPTAIN_KUBE=/captain-kube
ENV PLAYBOOKS=${CAPTAIN_KUBE}/playbooks

VOLUME /tmp

COPY docs/initial.sh /initial.sh
COPY docs/upgrade.sh /upgrade.sh
COPY docs/playbooks/ ${PLAYBOOKS}/
COPY dist/main ${CAPTAIN_KUBE}/main
COPY templates/ ${CAPTAIN_KUBE}/templates/
COPY static/ ${CAPTAIN_KUBE}/static/
COPY docs/docker-compose.yml ${CAPTAIN_KUBE}/docker-compose.yml

WORKDIR ${CAPTAIN_KUBE}

CMD ${CAPTAIN_KUBE}/main -playbooks=${PLAYBOOKS} -workdir=/data