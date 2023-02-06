FROM registry.access.redhat.com/ubi8-minimal

ENV SVC_NAME svc

ENV APP_ROOT=/opt/app-root \
    DATA_ROOT=/opt/data \
    HOME=/opt/app-root/src \
    PATH=/opt/app-root/src/bin:/opt/app-root/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

USER root

VOLUME ${DATA_ROOT}
RUN mkdir -p ${DATA_ROOT} && \
    mkdir -p ${HOME}/bin

COPY $SVC_NAME ${HOME}/bin/svc

WORKDIR ${HOME}
EXPOSE 80 443

CMD ["svc"]
