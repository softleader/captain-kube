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
		sshpass \
	&& rm -rf /var/cache/apk/* && \
	ls /usr/share/zoneinfo && \
	cp /usr/share/zoneinfo/Asia/Taipei /etc/localtime && \
	echo "Asia/Taipei" > /etc/timezone

ENV CAPTAIN_KUBE=/captain-kube
ENV PLAYBOOKS=${CAPTAIN_KUBE}/playbooks
ENV HOST_WORKSPACE=""

COPY initial.sh /initial.sh
COPY docs/playbooks/* ${PLAYBOOKS}/
COPY docker-compose.yml ${CAPTAIN_KUBE}/docker-compose.yml
COPY dist/main ${CAPTAIN_KUBE}/main
COPY templates/* ${CAPTAIN_KUBE}/templates/
COPY static/* ${CAPTAIN_KUBE}/static/

WORKDIR ${CAPTAIN_KUBE}

CMD ${CAPTAIN_KUBE}/main -playbooks=${PLAYBOOKS} -workspace=/data -host-workspace=${HOST_WORKSPACE}