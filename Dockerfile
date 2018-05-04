FROM softleader/ansible
MAINTAINER softleader.com.tw

RUN curl -L https://github.com/docker/compose/releases/download/1.21.0/docker-compose-$(uname -s)-$(uname -m) -o /docker-compose && \
	curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash && \
	helm init -c

ENV CAPTAIN_KUBE=/captain-kube
ENV PLAYBOOKS=${CAPTAIN_KUBE}/playbooks

VOLUME /tmp

# 這邊跟 Captain Kube 的設定或控制有關, 不會在 runtime 參考到的
COPY docs/initial.sh /initial.sh
COPY docs/upgrade.sh /upgrade.sh
COPY docs/docker-compose.yml /docker-compose.yml
COPY docs/daemon.yaml /daemon.yaml

# 這邊是 Captain Kube runtime 使用
COPY docs/playbooks/ ${PLAYBOOKS}/
COPY dist/main ${CAPTAIN_KUBE}/main
COPY templates/ ${CAPTAIN_KUBE}/templates/
COPY static/ ${CAPTAIN_KUBE}/static/

WORKDIR ${CAPTAIN_KUBE}

CMD ${CAPTAIN_KUBE}/main -playbooks=${PLAYBOOKS} -workdir=/data