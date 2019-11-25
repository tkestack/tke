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

package template

const (
	GPUManagerDaemonSetTemplate = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  namespace: kube-system
spec:
  updateStrategy:
    type: RollingUpdate
  minReadySeconds: 10
  selector:
    matchLabels:
      name: gpu-manager-ds
  template:
    metadata:
      # This annotation is deprecated. Kept here for backward compatibility
      # See https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ""
      labels:
        name: gpu-manager-ds
    spec:
      serviceAccount: gpu-manager
      tolerations:
        # This toleration is deprecated. Kept here for backward compatibility
        # See https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/
        - key: CriticalAddonsOnly
          operator: Exists
        - key: tencent.com/vcuda-core
          operator: Exists
          effect: NoSchedule
      # Mark this pod as a critical add-on; when enabled, the critical add-on
      # scheduler reserves resources for critical add-on pods so that they can
      # be rescheduled after a failure.
      # See https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/
      priorityClassName: "system-node-critical"
      # only run node hash gpu device
      nodeSelector:
        nvidia-device-enable: enable
      hostPID: true
      containers:
        - name: gpu-manager
          securityContext:
            privileged: true
          ports:
            - containerPort: 5678
          resources:
            requests:
              cpu: "1"
              memory: 1Gi
            limits:
              cpu: "2"
              memory: 2Gi
          volumeMounts:
            - name: device-plugin
              mountPath: /var/lib/kubelet/device-plugins
            - name: vdriver
              mountPath: /etc/gpu-manager/vdriver
            - name: vmdata
              mountPath: /etc/gpu-manager/vm
            - name: log
              mountPath: /var/log/gpu-manager
            - name: rundir
              mountPath: /var/run
              readOnly: true
            - name: cgroup
              mountPath: /sys/fs/cgroup
              readOnly: true
            - name: usr-directory
              mountPath: /usr/local/host
              readOnly: true
          env:
            - name: LOG_LEVEL
              value: "3"
            - name: EXTRA_FLAGS
              value: "--incluster-mode=true"
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
      volumes:
        - name: device-plugin
          hostPath:
            type: Directory
            path: /var/lib/kubelet/device-plugins
        - name: vmdata
          hostPath:
            type: DirectoryOrCreate
            path: /etc/gpu-manager/vm
        - name: vdriver
          hostPath:
            type: DirectoryOrCreate
            path: /etc/gpu-manager/vdriver
        - name: log
          hostPath:
            type: DirectoryOrCreate
            path: /etc/gpu-manager/log
        - name: rundir
          hostPath:
            type: Directory
            path: /var/run
        - name: cgroup
          hostPath:
            type: Directory
            path: /sys/fs/cgroup
        # We have to mount /usr directory instead of specified library path, because of non-existing
        # problem for different distro
        - name: usr-directory
          hostPath:
            type: Directory
            path: /usr
`
	ImagePatchTemplate = `[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"%s"}]`
	DeploymentTemplate = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  namespace: kube-system
  labels:
    k8s-app: gpu-quota-admission
  name: gpu-quota-admission
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: gpu-quota-admission
  strategy:
    type: Recreate
  template:
    metadata:
      namespace: kube-system
      labels:
        k8s-app: gpu-quota-admission
    spec:
      tolerations:
      - key: "node-role.kubernetes.io/master" 
        value: ""
        effect: "NoSchedule"
      affinity:
        nodeAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 1
            preference:
              matchExpressions:
              - key: node-role.kubernetes.io/master
                operator: Exists
      serviceAccount: gpu-manager
      initContainers:
      - image: busybox
        name: init-kube-config
        command: ['sh', '-c',' mkdir -p /etc/kubernetes/ && cp /root/gpu-quota-admission/gpu-quota-admission.config /etc/kubernetes/']
        volumeMounts:
        - mountPath: /root/gpu-quota-admission/
          name: config
        securityContext:
          privileged: true
      containers:
      - name: gpu-quota-admission
        resources:
          requests:
            cpu: "1"
            memory: 1Gi
          limits:
            cpu: "2"
            memory: 2Gi
        env:
        - name: LOG_LEVEL
          value: "4"
        - name: EXTRA_FLAGS
          value: "--incluster-mode=true" 
        ports:
        - containerPort: 3456
        volumeMounts:
        - mountPath: /root/gpu-quota-admission/
          name: config
      dnsPolicy: ClusterFirstWithHostNet
      priority: 2000000000
      priorityClassName: system-cluster-critical
      volumes:
        - configMap:
            defaultMode: 420
            name: gpu-quota-admission
          name: config
`
	ServiceAccountTemplate = `
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: kube-system
  name: gpu-manager
`
	RoleTemplate = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: gpu-manager
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'
`
	RoleBindingTemplate = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: gpu-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: gpu-manager
subjects:
- kind: ServiceAccount
  name: gpu-manager
  namespace: kube-system
`
	MetricServiceTemplate = `
apiVersion: v1
kind: Service
metadata:
  name: gpu-manager-metric
  namespace: kube-system
  annotations:
    tke.prometheus.io/scrape: "true"
    prometheus.io/scheme: "http"
    prometheus.io/port: "5678"
    prometheus.io/path: "/metric"
  labels:
    kubernetes.io/cluster-service: "true"
spec:
  clusterIP: None
  ports:
    - name: metrics
      port: 5678
      protocol: TCP
      targetPort: 5678
  selector:
    name: gpu-manager-ds
`
	ConfigMapTemplate = `
apiVersion: v1
data:
  gpu-quota-admission.config: |
   {
        "QuotaConfigMapName": "gpuquota",
        "QuotaConfigMapNamespace": "kube-system",
        "GPUModelLabel": "gaia.tencent.com/gpu-model",
        "GPUPoolLabel": "gaia.tencent.com/gpu-pool"
    }
kind: ConfigMap
metadata:
  name: gpu-quota-admission
  namespace: kube-system
`
	ServiceTemplate = `
kind: Service
apiVersion: v1
metadata:
  namespace: kube-system
  name: gpu-quota-admission
spec:
  selector:
    k8s-app: gpu-quota-admission
  ports:
    - protocol: TCP
      port: 3456
      targetPort: 3456
`
)
