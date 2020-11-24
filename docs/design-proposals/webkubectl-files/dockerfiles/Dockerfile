FROM busybox
RUN wget -O /bin/kubectl "https://storage.googleapis.com/kubernetes-release/release/$(wget -q -O - https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
&& chmod +x /bin/kubectl
CMD ["/bin/sh", "-c", "while true; do sleep 30; done;"]
