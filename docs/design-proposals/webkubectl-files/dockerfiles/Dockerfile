FROM busybox
RUN wget -O /bin/kubectl "https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl" \
&& chmod +x /bin/kubectl
CMD ["/bin/sh", "-c", "while true; do sleep 30; done;"]
