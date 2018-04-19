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

COPY anible/* /anible/
COPY initial.sh /initial.sh
COPY dist/main /main

WORKDIR /anible

CMD /main