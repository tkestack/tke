FROM BASE_IMAGE

RUN echo "hosts: files dns" >> /etc/nsswitch.conf

WORKDIR /app
ADD tke-logagent-api /app/bin/
ENTRYPOINT ["/app/bin/tke-logagent-api"]
