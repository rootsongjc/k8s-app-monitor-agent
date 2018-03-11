FROM alpine

ADD k8s-app-monitor-agent /usr/bin/

EXPOSE 8888

ENV APP_PORT 3000
ENV SERVICE_NAME localhost

ENTRYPOINT /usr/bin/k8s-app-monitor-agent
