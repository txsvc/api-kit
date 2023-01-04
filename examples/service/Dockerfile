FROM registry.access.redhat.com/ubi8-minimal

ENV SVC_NAME svc

ENV APP_ROOT=/opt/app-root \
    DATA_ROOT=/opt/data \
    HOME=/opt/app-root/src \
    PATH=/opt/app-root/src/bin:/opt/app-root/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

USER root

VOLUME ${DATA_ROOT}
RUN mkdir -p ${DATA_ROOT} && \
    chown -R 1001:0 ${DATA_ROOT}

RUN mkdir -p ${HOME}/bin && \
    chown -R 1001:0 ${APP_ROOT}

COPY $SVC_NAME ${HOME}/bin/svc

RUN chown -R 1001:0 ${APP_ROOT} && \
    chmod -R ug+rwx ${APP_ROOT}

WORKDIR ${HOME}
EXPOSE 8080

USER 1001

CMD ["svc"]
