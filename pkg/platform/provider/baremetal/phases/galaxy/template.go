/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package galaxy

const (
	//GalaxyDaemonsetTemplate decoded as galaxy daemonset
	GalaxyDaemonsetTemplate = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    app: galaxy
  name: galaxy
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app: galaxy
  template:
    metadata:
      labels:
        app: galaxy
    spec:
      serviceAccountName: galaxy
      hostNetwork: true
      hostPID: true
      containers:
      - image: galaxy:1.0.0-alpha
        command: ["/bin/sh"]
        args: ["-c", "cp -p /etc/cni/net.d/00-galaxy.conf /host/etc/cni/net.d/; cp -p /opt/cni/bin/* /host/opt/cni/bin/; /usr/bin/galaxy --logtostderr=true --v=3"]
        name: galaxy
        resources:
          requests:
            cpu: 100m
            memory: 200Mi
        securityContext:
          privileged: true
        volumeMounts:
        - name: galaxy-run
          mountPath: /var/run/galaxy/
        - name: flannel-run
          mountPath: /run/flannel
        - name: kube-config
          mountPath: /host/etc/kubernetes/
        - name: galaxy-log
          mountPath: /data/galaxy/logs
        - name: galaxy-etc
          mountPath: /etc/galaxy
        - name: cni-config
          mountPath: /host/etc/cni/net.d/
        - name: cni-bin
          mountPath: /host/opt/cni/bin
        - name: cni-etc
          mountPath: /etc/cni/net.d
        - name: cni-state
          mountPath: /var/lib/cni
        - name: docker-sock
          mountPath: /run/docker.sock
      tolerations:
      - effect: NoSchedule
        operator: Exists
      terminationGracePeriodSeconds: 30
      volumes:
      - name: galaxy-run
        hostPath:
          path: /var/run/galaxy
      - name: flannel-run
        hostPath:
          path: /run/flannel
      - name: kube-config
        hostPath:
          path: /etc/kubernetes/
      - name: cni-bin-dir
        hostPath:
          path: /opt/cni/bin
      - name: galaxy-log
        emptyDir: {}
      - configMap:
          defaultMode: 420
          name: galaxy-etc
        name: galaxy-etc
      - name: cni-config
        hostPath:
          path: /etc/cni/net.d/
      - name: cni-bin
        hostPath:
          path: /opt/cni/bin
      - name: cni-state
        hostPath:
          path: /var/lib/cni
      - configMap:
          defaultMode: 420
          name: cni-etc
        name: cni-etc
      - name: docker-sock
        hostPath:
          path: /run/docker.sock
`

	//GalaxyCM decoded as galaxy & cni configMap
	GalaxyCM = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: galaxy-etc
  namespace: kube-system
data:
  galaxy.json: |
    {
      "NetworkConf":[
        {"type":"tke-route-eni","eni":"eth1","routeTable":1},
        {"type":"galaxy-flannel", "delegate":{"type":"galaxy-veth"},"subnetFile":"/run/flannel/subnet.env"},
        {"type":"galaxy-k8s-vlan", "device":"{{ .DeviceName }}", "default_bridge_name": "br0"},
        {"type": "galaxy-k8s-sriov", "device": "{{ .DeviceName }}", "vf_num": 10}
      ],
      "DefaultNetworks": ["galaxy-flannel"]
    }
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cni-etc
  namespace: kube-system
data:
  00-galaxy.conf: |
    {
      "type": "galaxy-sdn",
      "capabilities": {"portMappings": true}
    }
`

	//FlannelDaemonset decoded as flannel daemonset
	FlannelDaemonset = `
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: kube-flannel-ds-amd64
  namespace: kube-system
  labels:
    k8s-app: flannel
spec:
  selector:
    matchLabels:
      k8s-app: flannel
  template:
    metadata:
      labels:
        k8s-app: flannel
    spec:
      hostNetwork: true
      nodeSelector:
        beta.kubernetes.io/arch: amd64
      tolerations:
      - operator: Exists
        effect: NoSchedule
      serviceAccountName: flannel
      containers:
      - name: kube-flannel
        image: {{ .Image }}
        command:
        - /opt/bin/flanneld
        args:
        - --ip-masq
        - --kube-subnet-mgr
        resources:
          requests:
            cpu: "100m"
            memory: "50Mi"
          limits:
            cpu: "100m"
            memory: "256Mi"
        securityContext:
          privileged: true
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        volumeMounts:
        - name: run
          mountPath: /run
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      volumes:
        - name: run
          hostPath:
            path: /run
        - name: cni
          hostPath:
            path: /etc/cni/net.d
        - name: flannel-cfg
          configMap:
            name: kube-flannel-cfg
`

	//FlannelCM decoded as flannel configMap
	FlannelCM = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-flannel-cfg
  namespace: kube-system
  labels:
    tier: node
    app: flannel
data:
  cni-conf.json: |
    {
      "name": "cbr0",
      "plugins": [
        {
          "type": "flannel",
          "delegate": {
            "hairpinMode": true,
            "isDefaultGateway": true
          }
        },
        {
          "type": "portmap",
          "capabilities": {
            "portMappings": true
          }
        }
      ]
    }
  net-conf.json: |
    {
      "Network": "{{ .Network }}",
      "Backend": {
        "Type": "vxlan"
      }
    }
`
)
