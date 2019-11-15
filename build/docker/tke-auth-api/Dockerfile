FROM BASE_IMAGE

ENV TZ=Asia/Shanghai
RUN apk add tzdata && cp /usr/share/zoneinfo/$TZ /etc/localtime

RUN echo "hosts: files dns" >> /etc/nsswitch.conf

WORKDIR /app
ADD tke-auth-api /app/bin/
COPY web/ /app/web/auth

ENTRYPOINT ["/app/bin/tke-auth-api"]
