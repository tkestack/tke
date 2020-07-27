FROM shsmu/alpine:3.1.2

WORKDIR /data
RUN mkdir -p /data/build/helm

COPY ./*_helm /data/build/helm
COPY Dockerstart /start
RUN chmod +x /start
