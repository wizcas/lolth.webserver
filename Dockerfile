FROM scratch
ARG tag=anomalous
ARG commit_id=anomalous
COPY artifact/ca-certificates.crt /etc/ssl/certs/
COPY artifact/webserver /
ENTRYPOINT ["/webserver"]
LABEL org.opencontainers.image.version=${tag}
LABEL org.opencontainers.image.revision=${commit_id}
