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

ENV ANSIBLE=/ansible-playbook
ENV WORKSPACE=/data
ENV HOST_WORKSPACE=""

COPY ansible-playbook/* /ansible-playbook/
COPY initial.sh /initial.sh
COPY docker-compose.yml /docker-compose.yml
COPY dist/main /main

WORKDIR ${ANSIBLE}

CMD /main -ansible=${ANSIBLE} -workspace=${WORKSPACE} -host-workspace=${HOST_WORKSPACE}