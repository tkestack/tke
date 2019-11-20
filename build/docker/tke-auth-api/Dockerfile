FROM alpine:3.10

RUN echo "hosts: files dns" >> /etc/nsswitch.conf

WORKDIR /app
ADD tke-auth-api /app/bin/
COPY web/ /app/web/auth

ENTRYPOINT ["/app/bin/tke-auth-api"]
