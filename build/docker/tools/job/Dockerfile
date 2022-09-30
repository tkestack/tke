FROM python:3.6-alpine
ENV HOME="/" \
    OS_ARCH="amd64" \
    OS_NAME="linux" \
    APP_VERSION="1.22.13"

RUN mkdir -p /tmp/kubectl/ && cd /tmp/kubectl && \
    if [ ! -f kubectl-$APP_VERSION-1-$OS_NAME-$OS_ARCH-debian-11.tar.gz ]; then \
      wget https://downloads.bitnami.com/files/stacksmith/kubectl-$APP_VERSION-1-$OS_NAME-$OS_ARCH-debian-11.tar.gz ; \
    fi && \
    export CHECKSUM=$(sha256sum kubectl-$APP_VERSION-1-$OS_NAME-$OS_ARCH-debian-11.tar.gz) && \
    echo "$CHECKSUM" | sha256sum -c - && \
    tar -zxf kubectl-$APP_VERSION-1-$OS_NAME-$OS_ARCH-debian-11.tar.gz -C /bin/ --strip-components=4 && \
    rm -rf /tmp/kubectl && \
    apk add --update curl && \
    apk add --no-cache bash && \
    rm -rf /var/cache/apk/*

ENTRYPOINT [ "kubectl" ]
CMD [ "--help" ]
