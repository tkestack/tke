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

package prometheus

import (
	"fmt"
	"strings"
)

func scapeConfigForPrometheus() string {
	cfgStr := `
    # Use kubelet_running_pod_count to get kube node labels
    - job_name: 'kubernetes-nodes'
      scrape_timeout: 60s
      scheme: https
      tls_config:
        ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        insecure_skip_verify: true
      bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
      kubernetes_sd_configs:
      - role: node
      relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_node_label_(.+)
      metric_relabel_configs:
      - source_labels: [ __name__ ]
        regex: 'kubelet_running_pod_count'
        action: keep
      - regex: (__name__|instance|node_role_kubernetes_io_master)
        action: labelkeep
      - source_labels: [ __name__ ]
        target_label: "node_role"
        replacement: "Node"
      - source_labels: [node_role_kubernetes_io_master]
        regex: "true"
        target_label: "node_role"
        replacement: "Master"
      - source_labels: [instance]
        target_label: "node"
      - regex: "instance|node_role_kubernetes_io_master"
        action: labeldrop

    - job_name: 'kubernetes-cadvisor'
      scrape_timeout: 60s
      scheme: https
      tls_config:
        ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        insecure_skip_verify: true
      bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
      kubernetes_sd_configs:
      - role: node
      relabel_configs:
      - action: labelmap
        regex: __meta_kubernetes_node_label_(.+)
      - target_label: __metrics_path__
        replacement: /metrics/cadvisor
      metric_relabel_configs:
      - source_labels: [ __name__ ]
        regex: 'container_fs_writes_bytes_total|container_fs_reads_bytes_total|container_fs_writes_total|container_fs_reads_total|container_cpu_usage_seconds_total|container_memory_usage_bytes|container_memory_cache|container_network_receive_bytes_total|container_network_transmit_bytes_total|container_network_receive_packets_total|container_network_transmit_packets_total'
        action: keep
      - regex: (__name__|container_name|pod_name|namespace|cpu|interface|device)
        action: labelkeep
      - source_labels: [pod_name]
        regex: "^$"
        action: drop
      - source_labels: [container_name]
        regex: "^$"
        action: drop
      - source_labels: [container_name]
        regex: "POD"
        target_label: container_name
        replacement: "pause"
      - source_labels: [id]
        regex: "/kubepods/(.*)pod(.*)/(.*)"
        target_label: container_id
        replacement: $2

    - job_name: 'tke-service-endpoints'
      scrape_timeout: 60s
      kubernetes_sd_configs:
      - role: endpoints
      relabel_configs:
      - source_labels: [__meta_kubernetes_service_annotation_tke_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (https?)
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
      metric_relabel_configs:
      - source_labels: [ __name__ ]
        regex: 'container_gpu_utilization|container_request_gpu_utilization|container_gpu_memory_total|container_request_gpu_memory|kube_node_status_allocatable|kube_node_status_capacity|kube_node_status_allocatable_cpu_cores|kube_node_status_allocatable_memory_bytes|kube_job_status_failed|kube_statefulset_status_replicas_ready|kube_statefulset_replicas|kube_daemonset_status_number_unavailable|kube_deployment_status_replicas_unavailable|kube_pod_labels|kube_pod_info|kube_pod_status_ready|kube_pod_container_status_restarts_total|kube_pod_container_resource_requests|kube_pod_container_resource_limits|kube_node_status_condition|kube_node_status_capacity_cpu_cores|kube_node_status_capacity_memory_bytes|kube_replicaset_owner|kube_namespace_labels'
        action: keep
      - source_labels: [created_by_kind]
        action: replace
        target_label: workload_kind
      - source_labels: [created_by_name]
        action: replace
        target_label: "workload_name"
      - source_labels: [pod]
        regex: (.+)
        action: replace
        target_label: "pod_name"
      - source_labels: [container]
        regex: (.+)
        action: replace
        target_label: "container_name"
      - source_labels: [label_tke_cloud_tencent_com_projectName]
        action: replace
        target_label: "project_name"
      - regex: "created_by_kind|created_by_name|pod|job|uid|pod_ip|host_ip|instance|__meta_kubernetes_namespace|__meta_kubernetes_service_name|__meta_kubernetes_service_label_(.+)|owner_is_controller|container"
        action: labeldrop

    - job_name: 'kubernetes-service-endpoints'
      scrape_timeout: 60s
      kubernetes_sd_configs:
      - role: endpoints
      relabel_configs:
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (https?)
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
      - source_labels: [__meta_kubernetes_service_label_kubernetes_io_name]
        regex: CoreDNS
        action: drop
      metric_relabel_configs:

    - job_name: 'kubernetes-pods'
      scrape_timeout: 60s
      honor_labels: false
      kubernetes_sd_configs:
      - role: pod
      tls_config:
        insecure_skip_verify: true
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (.+)
      metric_relabel_configs:

    - job_name: 'tke-pods'
      scrape_timeout: 60s
      honor_labels: false
      kubernetes_sd_configs:
      - role: pod
      tls_config:
        insecure_skip_verify: true
        ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
      bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_tke_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: namespace
      - source_labels: [__meta_kubernetes_pod_name]
        action: drop
        regex: etcd.+
      - source_labels: [__meta_kubernetes_pod_name]
        action: replace
        target_label: pod_name
      - source_labels: [__meta_kubernetes_pod_node_name]
        action: replace
        target_label: node
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (.+)
      metric_relabel_configs:
      - source_labels: [ __name__ ]
        regex: 'scheduler_e2e_scheduling_latency_microseconds_sum|scheduler_e2e_scheduling_latency_microseconds_count|apiserver_request_latencies_summary_count|apiserver_request_latencies_summary_sum|node_sockstat_TCP_inuse|node_network_transmit_bytes|node_network_receive_bytes|node_filesystem_size|node_filesystem_avail|node_disk_bytes_written|node_disk_bytes_read|node_disk_writes_completed|node_disk_reads_completed'
        action: keep
      - regex: "instance|job|pod_name|namespace|scope|subresource"
        action: labeldrop

    - job_name: 'tke-etcd'
      scrape_timeout: 60s
      honor_labels: false
      scheme: https
      kubernetes_sd_configs:
      - role: pod
      tls_config:
        insecure_skip_verify: true
        ca_file: /etc/prometheus/secrets/prometheus-etcd/etcd-ca.crt
        cert_file: /etc/prometheus/secrets/prometheus-etcd/etcd-client.crt
        key_file: /etc/prometheus/secrets/prometheus-etcd/etcd-client.key
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_tke_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: namespace
      - source_labels: [__meta_kubernetes_pod_name]
        action: keep
        regex: etcd.+
      - source_labels: [__meta_kubernetes_pod_name]
        action: replace
        target_label: pod_name
      - source_labels: [__meta_kubernetes_pod_node_name]
        action: replace
        target_label: node
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (.+)
      metric_relabel_configs:
      - source_labels: [ __name__ ]
        regex: 'etcd_server_leader_changes_seen_total|etcd_debugging_mvcc_db_total_size_in_bytes'
        action: keep
      - regex: "instance|job|pod_name|namespace|scope|subresource"
        action: labeldrop
`
	return cfgStr
}

func recordRulesForPrometheus() string {
	rules := fmt.Sprintf(`
groups:
- name: k8s-ag-data
  rules:
  - record: kube_node_labels
    expr: kubelet_running_pod_count*0 + 1

  - record: kube_node_status_capacity_gpu
    expr: sum by(node) (kube_node_status_capacity{resource="tencent_com_vcuda_core"})

  - record: kube_node_status_capacity_gpu_memory
    expr: sum by(node) (kube_node_status_capacity{resource="tencent_com_vcuda_memory"})

  - record: kube_node_status_allocatable_gpu
    expr: sum by(node) (kube_node_status_allocatable{resource="tencent_com_vcuda_core"})

  - record: kube_node_status_allocatable_gpu_memory
    expr: sum by(node) (kube_node_status_allocatable{resource="tencent_com_vcuda_memory"})

  - record: __pod_info1
    expr: kube_pod_info* on(node) group_left(node_role) kube_node_labels

  - record: __pod_info2
    expr:  label_replace(label_replace(__pod_info1{workload_kind="ReplicaSet"} * on (workload_name,namespace) group_left(owner_name, owner_kind) label_replace(kube_replicaset_owner,"workload_name","$1","replicaset","(.*)"),"workload_name","$1","owner_name","(.*)"),"workload_kind","$1","owner_kind","(.*)")  or on(pod_name,namesapce)  __pod_info1{workload_kind != "ReplicaSet"}

  - record: k8s_cluster_cpu_core_total
    expr: sum(kube_node_status_allocatable_cpu_cores * on(node) group_left kube_node_labels {node_role="Node"})

  - record: k8s_cluster_memory_total
    expr: sum(kube_node_status_allocatable_memory_bytes * on(node) group_left kube_node_labels {node_role="Node"})

  - record: k8s_cluster_gpu_total
    expr: sum(kube_node_status_allocatable_gpu * on(node) group_left kube_node_labels {node_role="Node"})

  - record: k8s_cluster_gpu_memory_total
    expr: sum(kube_node_status_allocatable_gpu_memory * on(node) group_left kube_node_labels {node_role="Node"})

  - record: k8s_container_cpu_core_used
    expr: sum(rate(container_cpu_usage_seconds_total[2m])) by (container_name, namespace, pod_name) * on(namespace, pod_name) group_left(workload_kind, workload_name, node, node_role)  __pod_info2

  - record: k8s_container_rate_cpu_core_used_request
    expr: k8s_container_cpu_core_used * 100 / on (pod_name,namespace,container_name)  group_left  kube_pod_container_resource_requests{resource="cpu"}

  - record: k8s_container_rate_cpu_core_used_limit
    expr: k8s_container_cpu_core_used * 100 / on (pod_name,namespace,container_name)  group_left  kube_pod_container_resource_limits{resource="cpu"}

  - record: k8s_container_rate_cpu_core_used_node
    expr: k8s_container_cpu_core_used * 100 / on(node) group_left  kube_node_status_capacity_cpu_cores

  - record: k8s_container_mem_usage_bytes
    expr: container_memory_usage_bytes * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_mem_no_cache_bytes
    expr: (container_memory_usage_bytes -  container_memory_cache)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_rate_mem_usage_request
    expr: k8s_container_mem_usage_bytes * 100 / on (pod_name,namespace,container_name)  group_left kube_pod_container_resource_requests{resource="memory"}

  - record: k8s_container_rate_mem_no_cache_request
    expr: k8s_container_mem_no_cache_bytes * 100 / on (pod_name,namespace,container_name)  group_left kube_pod_container_resource_requests{resource="memory"}

  - record: k8s_container_rate_mem_usage_limit
    expr: k8s_container_mem_usage_bytes * 100 / on (pod_name,namespace,container_name)  group_left kube_pod_container_resource_limits{resource="memory"}

  - record: k8s_container_rate_mem_no_cache_limit
    expr: k8s_container_mem_no_cache_bytes * 100 / on (pod_name,namespace,container_name)  group_left kube_pod_container_resource_limits{resource="memory"}

  - record: k8s_container_rate_mem_usage_node
    expr: k8s_container_mem_usage_bytes * 100 / on(node) group_left  kube_node_status_capacity_memory_bytes

  - record: k8s_container_rate_mem_no_cache_node
    expr: k8s_container_mem_no_cache_bytes * 100 / on(node) group_left  kube_node_status_capacity_memory_bytes

  - record: k8s_container_gpu_used
    expr: container_gpu_utilization{gpu="total"} * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role) __pod_info2

  - record: k8s_container_rate_gpu_used_request
    expr: k8s_container_gpu_used / on (pod_name,namespace,container_name) group_left container_request_gpu_utilization

  - record: k8s_container_rate_gpu_used_node
    expr: k8s_container_gpu_used * 100 / on(node) group_left kube_node_status_capacity_gpu

  - record: k8s_container_gpu_memory_used
    expr: container_gpu_memory_total{gpu_memory="total"} / 256 * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role) __pod_info2

  - record: k8s_container_rate_gpu_memory_used_request
    expr: k8s_container_gpu_memory_used * 100 / on (pod_name,namespace,container_name) group_left() (container_request_gpu_memory / 256)

  - record: k8s_container_rate_gpu_memory_used_node
    expr: k8s_container_gpu_memory_used * 100 / on(node) group_left() kube_node_status_capacity_gpu_memory

  - record: k8s_container_network_receive_bytes_bw
    expr: sum(rate(container_network_receive_bytes_total[2m])) without(interface)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_network_transmit_bytes_bw
    expr: sum(rate(container_network_transmit_bytes_total[2m])) without(interface)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_network_receive_bytes
    expr: sum(idelta(container_network_receive_bytes_total[2m])) without(interface)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_network_transmit_bytes
    expr: sum(idelta(container_network_transmit_bytes_total[2m])) without(interface)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_network_receive_packets
    expr: sum(rate(container_network_receive_packets_total[2m])) without(interface)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_network_transmit_packets
    expr: sum(rate(container_network_transmit_packets_total[2m])) without(interface)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_fs_read_bytes
    expr: sum(rate(container_fs_reads_bytes_total[2m])) without(device)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_fs_write_bytes
    expr: sum(rate(container_fs_writes_bytes_total[2m])) without(device)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_fs_read_times
    expr: sum(rate(container_fs_reads_total[2m])) without(device)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_container_fs_write_times
    expr: sum(rate(container_fs_writes_total[2m])) without(device)  * on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_pod_cpu_core_used
    expr: sum(k8s_container_cpu_core_used) without (container_name,container_id)

  - record: k8s_pod_rate_cpu_core_used_request
    expr: sum(k8s_container_cpu_core_used + on (container_name, pod_name, namespace) group_left kube_pod_container_resource_requests{resource="cpu"} * 0) without(container_name )   * 100  / on (pod_name,namespace)  group_left  sum(kube_pod_container_resource_requests{resource="cpu"})  without(container_name)

  - record: k8s_pod_rate_cpu_core_used_limit
    expr: sum(k8s_container_cpu_core_used + on (container_name, pod_name, namespace) group_left kube_pod_container_resource_limits{resource="cpu"} * 0) without(container_name )   * 100  / on (pod_name,namespace)  group_left  sum(kube_pod_container_resource_limits{resource="cpu"})  without(container_name)

  - record: k8s_pod_rate_cpu_core_used_node
    expr: k8s_pod_cpu_core_used *100 /  on(node) group_left  kube_node_status_capacity_cpu_cores

  - record: k8s_pod_mem_usage_bytes
    expr: sum(k8s_container_mem_usage_bytes) without (container_name,container_id)

  - record: k8s_pod_mem_no_cache_bytes
    expr: sum(k8s_container_mem_no_cache_bytes) without (container_name,container_id)

  - record: k8s_pod_rate_mem_usage_request
    expr: sum(k8s_container_mem_usage_bytes + on (container_name, pod_name, namespace) group_left kube_pod_container_resource_requests{resource="memory"} * 0) without(container_name )   * 100    / on (pod_name,namespace)  group_left  sum(kube_pod_container_resource_requests{resource="memory"}) without(container_name)

  - record: k8s_pod_rate_mem_no_cache_request
    expr: sum(k8s_container_mem_no_cache_bytes + on (container_name, pod_name, namespace) group_left kube_pod_container_resource_requests{resource="memory"} * 0) without(container_name )   * 100    / on (pod_name,namespace)  group_left  sum(kube_pod_container_resource_requests{resource="memory"}) without(container_name)

  - record: k8s_pod_rate_mem_usage_limit
    expr: sum(k8s_container_mem_usage_bytes + on (container_name, pod_name, namespace) group_left kube_pod_container_resource_limits{resource="memory"} * 0) without(container_name )   * 100    / on (pod_name,namespace)  group_left  sum(kube_pod_container_resource_limits{resource="memory"})  without(container_name)

  - record: k8s_pod_rate_mem_no_cache_limit
    expr: sum(k8s_container_mem_no_cache_bytes + on (container_name, pod_name, namespace) group_left kube_pod_container_resource_limits{resource="memory"} * 0) without(container_name )   * 100    / on (pod_name,namespace)  group_left  sum(kube_pod_container_resource_limits{resource="memory"})  without(container_name)

  - record: k8s_pod_rate_mem_usage_node
    expr: k8s_pod_mem_usage_bytes * 100  /  on(node) group_left  kube_node_status_capacity_memory_bytes

  - record: k8s_pod_rate_mem_no_cache_node
    expr: k8s_pod_mem_no_cache_bytes * 100 / on(node) group_left  kube_node_status_capacity_memory_bytes

  - record: k8s_pod_gpu_used
    expr: sum(k8s_container_gpu_used) without (container_name,container_id)

  - record: k8s_pod_gpu_request
    expr: sum(container_request_gpu_utilization * 100) without(container_name)

  - record: k8s_pod_rate_gpu_used_request
    expr: sum(k8s_container_gpu_used + on (container_name, pod_name, namespace) group_left container_request_gpu_utilization * 0) without(container_name) * 100 / on (pod_name,namespace) group_left k8s_pod_gpu_request

  - record: k8s_pod_rate_gpu_used_node
    expr: k8s_pod_gpu_used * 100  /  on(node) group_left  kube_node_status_capacity_gpu

  - record: k8s_pod_gpu_memory_used
    expr: sum(k8s_container_gpu_memory_used) without (container_name,container_id)

  - record: k8s_pod_gpu_memory_request
    expr: sum(container_request_gpu_memory / 256)  without(container_name)

  - record: k8s_pod_rate_gpu_memory_used_request
    expr: sum(k8s_container_gpu_memory_used + on (container_name, pod_name, namespace) group_left container_request_gpu_memory * 0) without(container_name) * 100  / on (pod_name,namespace) group_left k8s_pod_gpu_memory_request

  - record: k8s_pod_rate_gpu_memory_used_node
    expr: k8s_pod_gpu_memory_used * 100  /  on(node) group_left() kube_node_status_capacity_gpu_memory

  - record: k8s_pod_network_receive_bytes_bw
    expr: sum(k8s_container_network_receive_bytes_bw) without (container_name,container_id)

  - record: k8s_pod_network_transmit_bytes_bw
    expr: sum(k8s_container_network_transmit_bytes_bw) without (container_name,container_id)

  - record: k8s_pod_network_receive_bytes
    expr: sum(k8s_container_network_receive_bytes) without (container_name,container_id)

  - record: k8s_pod_network_transmit_bytes
    expr: sum(k8s_container_network_transmit_bytes) without (container_name,container_id)

  - record: k8s_pod_network_receive_packets
    expr: sum(k8s_container_network_receive_packets) without (container_name,container_id)

  - record: k8s_pod_network_transmit_packets
    expr: sum(k8s_container_network_transmit_packets) without (container_name,container_id)

  - record: k8s_pod_fs_read_bytes
    expr: sum(k8s_container_fs_read_bytes) without (container_name,container_id)

  - record: k8s_pod_fs_write_bytes
    expr: sum(k8s_container_fs_write_bytes) without (container_name,container_id)

  - record: k8s_pod_fs_read_times
    expr: sum(k8s_container_fs_read_times) without (container_name,container_id)

  - record: k8s_pod_fs_write_times
    expr: sum(k8s_container_fs_write_times) without (container_name,container_id)

  - record: k8s_pod_status_ready
    expr: sum(kube_pod_status_ready{condition="true"}) by (namespace,pod_name) *  on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_pod_restart_total
    expr: sum(idelta(kube_pod_container_status_restarts_total [2m])) by (namespace,pod_name) *  on(namespace, pod_name) group_left(workload_kind,workload_name,node, node_role)  __pod_info2

  - record: k8s_node_status_ready
    expr: max(kube_node_status_condition{condition="Ready", status="true"} * on (node) group_left(node_role)  kube_node_labels)  without(condition, status)

  - record: k8s_node_pod_restart_total
    expr: sum(k8s_pod_restart_total) without (pod_name,workload_kind,workload_name,namespace)

  - record: k8s_node_cpu_usage
    expr: sum(k8s_pod_cpu_core_used) without(namespace,pod_name,workload_kind,workload_name) *100 / on(node) group_left kube_node_status_capacity_cpu_cores

  - record: k8s_node_mem_usage
    expr: sum(k8s_pod_mem_usage_bytes) without(namespace,pod_name,workload_kind,workload_name) *100 / on(node) group_left kube_node_status_capacity_memory_bytes

  - record: k8s_node_gpu_usage
    expr: sum(k8s_pod_gpu_used) without(namespace,pod_name,workload_kind,workload_name) *100 / on(node) group_left kube_node_status_capacity_gpu

  - record: k8s_node_gpu_memory_usage
    expr: sum(k8s_pod_gpu_memory_used) without(namespace,pod_name,workload_kind,workload_name) *100 / on(node) group_left() kube_node_status_capacity_gpu_memory

  - record: k8s_node_fs_write_bytes
    expr: (sum by (node) (irate(node_disk_bytes_written[2m]))) *on(node) group_left(node_role) kube_node_labels

  - record: k8s_node_fs_read_bytes
    expr: (sum by (node) (irate(node_disk_bytes_read[2m])))*on(node) group_left(node_role) kube_node_labels

  - record: k8s_node_fs_write_times
    expr: (sum by (node) (irate(node_disk_writes_completed[2m])))*on(node) group_left(node_role) kube_node_labels

  - record: k8s_node_fs_read_times
    expr: (sum by (node) (irate(node_disk_reads_completed[2m])))*on(node) group_left(node_role) kube_node_labels

  - record: k8s_node_pod_num
    expr: count(k8s_pod_status_ready) without (pod_name,workload_kind,workload_name,namespace)

  - record: k8s_node_disk_space_rate
    expr: (100 - sum (node_filesystem_avail{fstype=~"ext3|ext4|xfs"}) by (node) / sum (node_filesystem_size{fstype=~"ext3|ext4|xfs"}) by (node) *100) *on(node) group_left(node_role) kube_node_labels

  - record: k8s_node_network_receive_bytes_bw
    expr: (sum by (node) (irate(node_network_receive_bytes{device!~"lo|veth(.*)|virb(.*)|docker(.*)|tunl(.*)|v-h(.*)|flannel(.*)"}[5m])))*on(node) group_left(node_role) kube_node_labels

  - record: k8s_node_network_transmit_bytes_bw
    expr: (sum by (node) (irate(node_network_transmit_bytes{device!~"lo|veth(.*)|virb(.*)|docker(.*)|tunl(.*)|v-h(.*)|flannel(.*)"}[5m])))*on(node) group_left(node_role) kube_node_labels

  - record: k8s_workload_abnormal
    expr: |-
            max(label_replace(
            label_replace(
            label_replace(
            kube_deployment_status_replicas_unavailable,
            "workload_kind","Deployment","","")
            ,"workload_name","$1","deployment","(.*)"),
            "__name__", "k8s_workload_abnormal", "__name__","(.*)") ) by (namespace, workload_name, workload_kind,__name__)
            or on (namespace,workload_name,workload_kind, __name__)
            max(label_replace(
            label_replace(
            label_replace(
            kube_daemonset_status_number_unavailable,
            "workload_kind","DaemonSet","","")
            ,"workload_name","$1","daemonset","(.*)"),
            "__name__", "k8s_workload_abnormal", "__name__","(.*)") ) by (namespace, workload_name, workload_kind,__name__)
            or on (namespace,workload_name,workload_kind, __name__)
            max(label_replace(
            label_replace(
            label_replace(
            (kube_statefulset_replicas - kube_statefulset_status_replicas_ready),
            "workload_kind","StatefulSet","","")
            ,"workload_name","$1","statefulset","(.*)"),
            "__name__", "k8s_workload_abnormal", "__name__","(.*)") ) by (namespace, workload_name, workload_kind,__name__)
            or on (namespace,workload_name,workload_kind, __name__)
            max(label_replace(
            label_replace(
            label_replace(
            (kube_job_status_failed),
            "workload_kind","Job","","")
            ,"workload_name","$1","job_name","(.*)"),
            "__name__", "k8s_workload_abnormal", "__name__","(.*)") ) by (namespace, workload_name, workload_kind,__name__)
            or on (namespace,workload_name,workload_kind, __name__)
            max(label_replace(
            label_replace(
            label_replace(
            (kube_cronjob_info * 0),
            "workload_kind","CronJob","","")
            ,"workload_name","","cronjob","(.*)"),
            "__name__", "k8s_workload_abnormal", "__name__","(.*)") ) by (namespace, workload_name, workload_kind,__name__)

  - record: k8s_workload_pod_restart_total
    expr: sum(k8s_pod_restart_total) by(namespace,workload_kind,workload_name)

  - record: k8s_workload_cpu_core_used
    expr: sum(k8s_pod_cpu_core_used) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_rate_cpu_core_used_cluster
    expr: k8s_workload_cpu_core_used * 100 / scalar(k8s_cluster_cpu_core_total)

  - record: k8s_workload_mem_usage_bytes
    expr: sum(k8s_pod_mem_usage_bytes) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_mem_no_cache_bytes
    expr: sum(k8s_pod_mem_no_cache_bytes) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_rate_mem_usage_bytes_cluster
    expr: k8s_workload_mem_usage_bytes  * 100 / scalar(k8s_cluster_memory_total)

  - record: k8s_workload_rate_mem_no_cache_cluster
    expr: k8s_workload_mem_no_cache_bytes * 100 / scalar(k8s_cluster_memory_total)

  - record: k8s_workload_network_receive_bytes_bw
    expr: sum(k8s_pod_network_receive_bytes_bw) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_network_transmit_bytes_bw
    expr: sum(k8s_pod_network_transmit_bytes_bw)  by(workload_name, workload_kind, namespace)

  - record: k8s_workload_network_receive_bytes
    expr: sum(k8s_pod_network_receive_bytes)  by(workload_name, workload_kind, namespace)

  - record: k8s_workload_network_transmit_bytes
    expr: sum(k8s_pod_network_transmit_bytes) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_network_receive_packets
    expr: sum(k8s_pod_network_receive_packets)  by(workload_name, workload_kind, namespace)

  - record: k8s_workload_network_transmit_packets
    expr: sum(k8s_pod_network_transmit_packets) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_fs_read_bytes
    expr: sum(k8s_pod_fs_read_bytes) by (workload_name, workload_kind, namespace)

  - record: k8s_workload_fs_write_bytes
    expr: sum(k8s_pod_fs_write_bytes) by (workload_name, workload_kind, namespace)

  - record: k8s_workload_fs_read_times
    expr: sum(k8s_pod_fs_read_times) by (workload_name, workload_kind, namespace)

  - record: k8s_workload_fs_write_times
    expr: sum(k8s_pod_fs_write_times) by (workload_name, workload_kind, namespace)

  - record: k8s_workload_gpu_used
    expr: sum(k8s_pod_gpu_used) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_rate_gpu_used_cluster
    expr: k8s_workload_gpu_used * 100 / scalar(k8s_cluster_gpu_total)

  - record: k8s_workload_gpu_memory_used
    expr: sum(k8s_pod_gpu_memory_used) by(workload_name, workload_kind, namespace)

  - record: k8s_workload_rate_gpu_memory_used_cluster
    expr: k8s_workload_gpu_memory_used * 100 / scalar(k8s_cluster_gpu_memory_total)

  - record: k8s_namespace_cpu_core_used
    expr: sum(k8s_pod_cpu_core_used) by (namespace)

  - record: k8s_namespace_rate_cpu_core_used_cluster
    expr: k8s_namespace_cpu_core_used * 100 / scalar(k8s_cluster_cpu_core_total)

  - record: k8s_namespace_mem_usage_bytes
    expr: sum(k8s_pod_mem_usage_bytes) by (namespace)

  - record: k8s_namespace_mem_no_cache_bytes
    expr: sum(k8s_pod_mem_no_cache_bytes) by (namespace)

  - record: k8s_namespace_rate_mem_usage_bytes_cluster
    expr: k8s_namespace_mem_usage_bytes * 100 / scalar(k8s_cluster_memory_total)

  - record: k8s_namespace_rate_mem_no_cache_cluster
    expr: k8s_namespace_mem_no_cache_bytes * 100 / scalar(k8s_cluster_memory_total)

  - record: k8s_namespace_network_receive_bytes_bw
    expr: sum(k8s_pod_network_receive_bytes_bw) by(namespace)

  - record: k8s_namespace_network_transmit_bytes_bw
    expr: sum(k8s_pod_network_transmit_bytes_bw) by(namespace)

  - record: k8s_namespace_network_receive_bytes
    expr: sum(k8s_pod_network_receive_bytes) by(namespace)

  - record: k8s_namespace_network_transmit_bytes
    expr: sum(k8s_pod_network_transmit_bytes) by(namespace)

  - record: k8s_namespace_network_receive_packets
    expr: sum(k8s_pod_network_receive_packets) by(namespace)

  - record: k8s_namespace_network_transmit_packets
    expr: sum(k8s_pod_network_transmit_packets) by(namespace)

  - record: k8s_namespace_fs_read_bytes
    expr: sum(k8s_workload_fs_read_bytes) by (namespace)

  - record: k8s_namespace_fs_write_bytes
    expr: sum(k8s_workload_fs_write_bytes) by (namespace)

  - record: k8s_namespace_fs_read_times
    expr: sum(k8s_workload_fs_read_times) by (namespace)

  - record: k8s_namespace_fs_write_times
    expr: sum(k8s_workload_fs_write_times) by (namespace)

  - record: k8s_namespace_gpu_used
    expr: sum(k8s_pod_gpu_used) by (namespace)

  - record: k8s_namespace_rate_gpu_used_cluster
    expr: k8s_namespace_gpu_used * 100 / scalar(k8s_cluster_gpu_total)

  - record: k8s_namespace_gpu_memory_used
    expr: sum(k8s_pod_gpu_memory_used) by (namespace)

  - record: k8s_namespace_rate_gpu_memory_used_cluster
    expr: k8s_namespace_gpu_memory_used * 100 / scalar(k8s_cluster_gpu_memory_total)

  - record: k8s_cluster_cpu_core_used
    expr:  sum(k8s_pod_cpu_core_used{node_role="Node"})

  - record: k8s_cluster_mem_usage_bytes
    expr:  sum(k8s_pod_mem_usage_bytes{node_role="Node"})

  - record: k8s_cluster_mem_no_cache_bytes
    expr: sum(k8s_pod_mem_no_cache_bytes{node_role="Node"})

  - record: k8s_cluster_rate_cpu_core_used_cluster
    expr: k8s_cluster_cpu_core_used  * 100 / scalar(k8s_cluster_cpu_core_total)

  - record: k8s_cluster_rate_cpu_core_request_cluster
    expr: sum(kube_pod_container_resource_requests{resource="cpu"} * on(node) group_left kube_node_labels {node_role="Node"} ) * 100 / scalar(k8s_cluster_cpu_core_total)

  - record: k8s_cluster_rate_mem_usage_bytes_cluster
    expr: k8s_cluster_mem_usage_bytes * 100 / scalar(k8s_cluster_memory_total)

  - record: k8s_cluster_rate_mem_no_cache_bytes_cluster
    expr: k8s_cluster_mem_no_cache_bytes * 100 / scalar(k8s_cluster_memory_total)

  - record: k8s_cluster_rate_mem_request_bytes_cluster
    expr: sum(kube_pod_container_resource_requests{resource="memory"} * on(node) group_left kube_node_labels {node_role="Node"} ) * 100 / scalar(k8s_cluster_memory_total)

  - record: k8s_cluster_network_receive_bytes_bw
    expr: sum(k8s_pod_network_receive_bytes_bw{node_role="Node"})

  - record: k8s_cluster_network_transmit_bytes_bw
    expr: sum(k8s_pod_network_transmit_bytes_bw{node_role="Node"})

  - record: k8s_cluster_network_receive_bytes
    expr: sum(k8s_pod_network_receive_bytes{node_role="Node"})

  - record: k8s_cluster_network_transmit_bytes
    expr: sum(k8s_pod_network_transmit_bytes{node_role="Node"})

  - record: k8s_cluster_network_receive_packets
    expr: sum(k8s_pod_network_receive_packets{node_role="Node"})

  - record: k8s_cluster_network_transmit_packets
    expr: sum(k8s_pod_network_transmit_packets{node_role="Node"})

  - record: k8s_cluster_fs_read_bytes
    expr: sum(k8s_pod_fs_read_bytes{node_role="Node"})

  - record: k8s_cluster_fs_write_bytes
    expr: sum(k8s_pod_fs_write_bytes{node_role="Node"})

  - record: k8s_cluster_fs_read_times
    expr: sum(k8s_pod_fs_read_times{node_role="Node"})

  - record: k8s_cluster_fs_write_times
    expr: sum(k8s_pod_fs_write_times{node_role="Node"})

  - record: k8s_cluster_gpu_used
    expr:  sum(k8s_pod_gpu_used{node_role="Node"})

  - record: k8s_cluster_rate_gpu_used_cluster
    expr: k8s_cluster_gpu_used  * 100 / scalar(k8s_cluster_gpu_total)

  - record: k8s_cluster_rate_gpu_request_cluster
    expr: sum(k8s_pod_gpu_request * on(node) group_left kube_node_labels {node_role="Node"}) * 100 / scalar(k8s_cluster_gpu_total)

  - record: k8s_cluster_gpu_memory_used
    expr:  sum(k8s_pod_gpu_memory_used{node_role="Node"})

  - record: k8s_cluster_rate_gpu_memory_used_cluster
    expr: k8s_cluster_gpu_memory_used  * 100 / scalar(k8s_cluster_gpu_memory_total)

  - record: k8s_cluster_rate_gpu_memory_request_cluster
    expr: sum(k8s_pod_gpu_memory_request * on(node) group_left kube_node_labels {node_role="Node"} ) * 100 / scalar(k8s_cluster_gpu_memory_total)

  - record: project_namespace_cpu_core_used
    expr: k8s_namespace_cpu_core_used* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_mem_usage_bytes
    expr: k8s_namespace_mem_usage_bytes* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_mem_no_cache_bytes
    expr: k8s_namespace_mem_no_cache_bytes* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_gpu_used
    expr: k8s_namespace_gpu_used* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_gpu_memory_used
    expr: k8s_namespace_gpu_memory_used* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_network_receive_bytes_bw
    expr: k8s_namespace_network_receive_bytes_bw* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_network_transmit_bytes_bw
    expr: k8s_namespace_network_transmit_bytes_bw* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_network_receive_bytes
    expr: k8s_namespace_network_receive_bytes* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_network_transmit_bytes
    expr: k8s_namespace_network_transmit_bytes* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_fs_read_bytes
    expr: k8s_namespace_fs_read_bytes* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_namespace_fs_write_bytes
    expr: k8s_namespace_fs_write_bytes* on(namespace) group_left(project_name) kube_namespace_labels

  - record: project_cluster_cpu_core_used
    expr: sum(project_namespace_cpu_core_used) by (project_name)

  - record: project_cluster_rate_cpu_core_used_cluster
    expr: project_cluster_cpu_core_used * 100 / scalar(k8s_cluster_cpu_core_total)

  - record: project_cluster_memory_usage_bytes
    expr: sum(project_namespace_mem_usage_bytes) by (project_name)

  - record: project_cluster_memory_no_cache_bytes
    expr: sum(project_namespace_mem_no_cache_bytes) by (project_name)

  - record: project_cluster_rate_memory_usage_bytes_cluster
    expr: project_cluster_memory_usage_bytes * 100 / scalar(k8s_cluster_memory_total)

  - record: project_cluster_rate_memory_no_cache_cluster
    expr: project_cluster_memory_no_cache_bytes * 100 / scalar(k8s_cluster_memory_total)

  - record: project_cluster_gpu_used
    expr: sum(project_namespace_gpu_used) by (project_name)

  - record: project_cluster_rate_gpu_used_cluster
    expr: project_cluster_gpu_used * 100 / scalar(k8s_cluster_gpu_total)

  - record: project_cluster_gpu_memory_used
    expr: sum(project_namespace_gpu_memory_used) by (project_name)

  - record: project_cluster_rate_gpu_memory_used_cluster
    expr: project_cluster_gpu_memory_used * 100 / scalar(k8s_cluster_gpu_memory_total)

  - record: project_cluster_network_receive_bytes_bw
    expr: sum(project_namespace_network_receive_bytes_bw) by (project_name)

  - record: project_cluster_network_transmit_bytes_bw
    expr: sum(project_namespace_network_transmit_bytes_bw) by (project_name)

  - record: project_cluster_network_receive_bytes
    expr: sum(project_namespace_network_receive_bytes) by (project_name)

  - record: project_cluster_network_transmit_bytes
    expr: sum(project_namespace_network_transmit_bytes) by (project_name)

  - record: project_cluster_fs_read_bytes
    expr: sum(project_namespace_fs_read_bytes) by (project_name)

  - record: project_cluster_fs_write_bytes
    expr: sum(project_namespace_fs_write_bytes) by (project_name)

  - record: k8s_component_apiserver_ready
    expr: up{instance=~"(.*)60001"} * on(node) group_left(node_role) kube_node_labels

  - record: k8s_component_etcd_ready
    expr: up{instance=~"(.*)2379"} * on(node) group_left(node_role) kube_node_labels

  - record: k8s_component_scheduler_ready
    expr: up{instance=~"(.*)10251"} * on(node) group_left(node_role) kube_node_labels

  - record: k8s_component_controller_manager_ready
    expr: up{instance=~"(.*)10252"} * on(node) group_left(node_role) kube_node_labels

  - record: k8s_component_apiserver_request_latency
    expr: sum(apiserver_request_latencies_summary_sum) by (node) / sum(apiserver_request_latencies_summary_count) by (node)

  - record: k8s_component_scheduler_scheduling_latency
    expr: sum(scheduler_e2e_scheduling_latency_microseconds_sum) by (node) / sum(scheduler_e2e_scheduling_latency_microseconds_count) by (node)
`)

	return rules
}

func configForAlertManager(webhookAddr string) string {
	config := fmt.Sprintf(`
    global:
      resolve_timeout: 5m

    route:
      group_by: ['alertname','alarmPolicyName','version']
      group_wait: 1s
      group_interval: 1s
      repeat_interval: 300s
      receiver: 'web.hook'
      routes:
      - match:
          service: app
        receiver: 'web.hook'

    receivers:
    - name: 'web.hook'
      webhook_configs:
      - url: '%s'
        http_config:
         tls_config:
          insecure_skip_verify: true
`, webhookAddr)

	return config
}

func configForPrometheusBeat(hosts []string, user, password string) string {
	var hostStrs []string
	for _, h := range hosts {
		hostStrs = append(hostStrs, fmt.Sprintf("\"%s\"", h))
	}
	config := fmt.Sprintf(`
    prometheusbeat:
      listen: ":8080"
      context: "/prometheus"
    output.elasticsearch:
      hosts: [%s]
      username: %s
      password: %s
`, strings.Join(hostStrs, ","), user, password)

	return config
}
