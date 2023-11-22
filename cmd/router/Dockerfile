FROM registry.access.redhat.com/ubi9-minimal

USER root

ENV SVC_NAME svc

ENV CADDYFILE /config/router/Caddyfile

EXPOSE 80 443 2019

# See https://caddyserver.com/docs/conventions#file-locations for details
ENV XDG_CONFIG_HOME=/config
ENV XDG_DATA_HOME=/data

RUN set -eux; \
	mkdir -p \
		/config/router \
		/data/router \
		/etc/router

# Copy extra files to the image
COPY ./run /root/usr/bin/

VOLUME /config
VOLUME /data

# copy the binary
COPY $SVC_NAME /usr/bin/svc

# make it executable
RUN chmod +x /usr/bin/svc

# Entrypoint
CMD ["/root/usr/bin/run"]