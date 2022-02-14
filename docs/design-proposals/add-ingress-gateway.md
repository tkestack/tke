# Add Ingress Gateway

## Why we need ingress gateway

### Gateway Rules

Since `tke-gateway` has occupied `80` and `443` port in global cluster, and all rules is writen as go code, if we want to export some service through `tke-gateway`, we have to rebuild `tke-gateway`. If we have an `nginx ingress` before `tke-gateway`, we cloud set `tke-service` as default backend, and set more rules for more service through `nginx ingress`


### Edge Case

For some edge cases, some nodes cannot access internal address of global cluster. But some services need thsat all nodes can access cluster internal IP: built-in `influxdb` is master0 node's 8086 port, and auth is 31138.

For edge cases, we can export these services through `nginx ingress` in `7-layer` networks, and we can set `LB` or `public IP` for our `ingress` to make it accessabel for edge nodes.

### Reduce TLS/SSL Maintenance

If we set `nginx ingress` as all service exporter, we cloud maintain TLS/SSL and hostname in one object.


## How to do it


### Allow NodePort in 80/443 port

Add `--service-cluster-ip-range` in `/etc/kubernetes/manifests/kube-apiserver.yaml`.

### Disabel hostNetwork mode and NodePort service

Disabel `hostNetwork` in `tke-gateway` and `influxdb`, and export these service through clusterIP.

Transfore `tke-auth` `NodePort` service to `ClusterIP` service.

### Prepare ingress controller chart

Please check https://github.com/kubernetes/ingress-nginx/tree/main/charts/ingress-nginx

We should set `ingress-nginx-controller service` as `NodePort` and export through 80/443.

Our chart should have an default ingress object:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: gateway
  namespace: tke
spec:
  defaultBackend:
    service:
      name: tke-gateway
      port:
        number: 80
  rules:
    # rules for other service like influxdb and tke-auth
  ingressClassName: nginx
```
