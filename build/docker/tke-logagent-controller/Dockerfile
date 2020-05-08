FROM BASE_IMAGE

RUN echo "hosts: files dns" >> /etc/nsswitch.conf

WORKDIR /app
ADD tke-logagent-controller /app/bin/
ENTRYPOINT ["/app/bin/tke-logagent-controller"]
