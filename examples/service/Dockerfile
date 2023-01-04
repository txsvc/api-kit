FROM registry.access.redhat.com/ubi8-minimal

ENV BUILD_NAME svc
#ENV PROJECT_ID podops
#ENV GOOGLE_APPLICATION_CREDENTIALS /config/credentials.json

EXPOSE 8080

# copy the binary
COPY $BUILD_NAME /usr/bin/svc

# ... and make it executable
RUN chmod +x /usr/bin/svc

USER 1001

ENTRYPOINT ["/usr/bin/svc" "-tls", "false"]