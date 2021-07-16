import * as React from "react";
import Entry from "./entry";
import ChartPanel, { ChartFilterPanel, ChartInstancesPanel, ColorTypes } from "../panel/index";
import { ChartType, ModelType } from "../charts/index";

class Chart extends React.Component {
  state = {
    table: "k8s_pod",
    fields: [
      {
        expr: "sum(k8s_pod_cpu_core_used)",
        unit: "%",
        thousands: 1024,
        alias: "CPU使用率"
      },
      {
        expr: "sum(k8s_pod_rate_cpu_core_used_limit)",
        unit: "%",
        thousands: 1,
        alias: "CPU限制",
        chartType: ChartType.Area
      },
      {
        expr: "sum(k8s_pod_rate_cpu_core_used_limit)",
        unit: "%",
        thousands: 1024,
        alias: "CPU限制"
      },
      {
        expr: "mean(k8s_pod_network_transmit_bytes)",
        alias: "网络出带宽",
        unit: "bps",
        thousands: 1024
      }
    ],
    conditions: [],
    groupBy: [
      {
        name: "workload_kind",
        value: "workload_kind"
      }
    ],
    instance: {
      columns: [
        {
          key: "workload_kind",
          name: "ID/名称"
        }
      ],
      list: [
        {
          workload_kind: "Deployment",
          isChecked: true
        },
        {
          workload_kind: "StatefulSet",
          isChecked: true
        },
        {
          isChecked: false,
          workload_kind: "DaemonSet"
        }
      ]
    }
  };

  constructor(props) {
    super(props);
    const linkTEA = document.createElement("link");
    linkTEA.rel = "stylesheet";
    linkTEA.href = "https://imgcache.qq.com/open_proj/proj_qcloud_v2/tea-style/dist/tea.me81529eb.min.css";
    document.head.appendChild(linkTEA);
    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = "https://imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/global/css/v1/global-201812261637.css";
    document.head.appendChild(link);
  }

  onReload() {
    this.setState({
      table: "k8s_pod",
      conditions: [["tke_cluster_instance_id", "=", "cls-lmp5edcu"], ["unInstanceId", "=", "ins-19402kvw"]],
      fields: [
        {
          expr: "max(k8s_pod_status_ready)",
          alias: "异常状态",
          chartType: "area",
          colors: ["#006eff"],
          valueLabels: {
            "0": "正常",
            "1": '<span class="text-danger">异常</span>'
          }
        },
        {
          expr: "mean(k8s_pod_cpu_core_used)",
          alias: "CPU使用量",
          unit: "核"
        },
        {
          expr: "mean(k8s_pod_rate_cpu_core_used_node)",
          alias: "CPU利用率(占主机)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_rate_cpu_core_used_request)",
          alias: "CPU利用率(占Request)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_rate_cpu_core_used_limit)",
          alias: "CPU利用率(占Limit)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_mem_usage_bytes)",
          alias: "内存使用量",
          unit: "B",
          thousands: 1024
        },
        {
          expr: "mean(k8s_pod_mem_no_cache_bytes)",
          alias: "内存使用量(不包含cache)",
          unit: "B",
          thousands: 1024
        },
        {
          expr: "mean(k8s_pod_rate_mem_usage_node)",
          alias: "内存利用率(占主机)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_rate_mem_no_cache_node)",
          alias: "内存利用率(占主机,不包含cache)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_rate_mem_usage_request)",
          alias: "内存利用率(占Request)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_rate_mem_no_cache_request)",
          alias: "内存利用率(占Request,不包含cache)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_rate_mem_usage_limit)",
          alias: "内存利用率(占Limit)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_rate_mem_no_cache_limit)",
          alias: "内存利用率(占Limit,不含cache)",
          unit: "%"
        },
        {
          expr: "mean(k8s_pod_network_receive_bytes)",
          alias: "网络入带宽",
          unit: "bps",
          thousands: 1024
        },
        {
          expr: "mean(k8s_pod_network_transmit_bytes)",
          alias: "网络出带宽",
          unit: "bps",
          thousands: 1024
        },
        {
          expr: "mean(k8s_pod_network_receive_bytes)",
          alias: "网络入流量",
          unit: "B",
          thousands: 1024
        },
        {
          expr: "mean(k8s_pod_network_transmit_bytes)",
          alias: "网络出流量",
          unit: "B",
          thousands: 1024
        },
        {
          expr: "mean(k8s_pod_network_receive_packets)",
          alias: "网络入包量",
          unit: "个/s"
        },
        {
          expr: "mean(k8s_pod_network_transmit_packets)",
          alias: "网络出包量",
          unit: "个/s"
        }
      ],
      groupBy: ["pod_name"],
      instance: {
        columns: [
          {
            key: "pod_name",
            name: "Pod名称"
          }
        ],
        list: []
      }
    });
  }

  render() {
    const { table, fields, groupBy, conditions, instance } = this.state;
    return (
      <div>
        <button onClick={this.onReload.bind(this)}>Reload</button>
        <ChartFilterPanel
          tables={[
            {
              table: table,
              fields: [
                {
                  expr: "k8s_pod_cpu_core_used",
                  unit: "%",
                  thousands: 1024,
                  alias: "CPU使用率",
                  storeKey: "k8s_pod_cpu_core_used"
                },
                {
                  expr: "k8s_pod_rate_cpu_core_used_limit",
                  unit: undefined,
                  thousands: 1,
                  alias: "CPU限制",
                  chartType: ChartType.Area
                },
                {
                  expr: "mean(k8s_pod_network_receive_bytes)",
                  alias: "网络入带宽",
                  unit: "bps",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_network_transmit_bytes)",
                  alias: "网络出带宽",
                  unit: "bps",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_network_receive_bytes)",
                  alias: "网络入流量",
                  unit: "B",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_network_transmit_bytes)",
                  alias: "网络出流量",
                  unit: "B",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_mem_usage_bytes)",
                  alias: "内存使用量",
                  unit: "B",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_mem_no_cache_bytes)",
                  alias: "内存使用量(不包含cache)",
                  unit: "B",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_network_receive_bytes)",
                  alias: "网络入流量",
                  unit: "B",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_network_transmit_bytes)",
                  alias: "网络出流量",
                  unit: "B",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_mem_usage_bytes)",
                  alias: "内存使用量",
                  unit: "B",
                  thousands: 1024
                },
                {
                  expr: "mean(k8s_pod_mem_no_cache_bytes)",
                  alias: "内存使用量(不包含cache)",
                  unit: "B",
                  thousands: 1024
                }
              ],
              conditions: []
            }
          ]}
          groupBy={[
            {
              name: "workload_kind",
              value: "workload_kind"
            },
            {
              name: "pod_name",
              value: "pod_name"
            }
          ]}
          //          conditions={ [
          //            {
          //              name: 'tke_cluster_instance_id',
          //              value: 'tke_cluster_instance_id',
          //            },
          //            {
          //              name: 'unInstanceId',
          //              value: 'unInstanceId',
          //            },
          //          ] }
        />
        <ChartPanel
          tables={[
            {
              table: "k8s_pod",
              fields: [
                {
                  expr: "sum(k8s_pod_cpu_core_used)",
                  unit: "%",
                  thousands: 1024,
                  alias: "CPU使用率",
                  colorTheme: ColorTypes.Multi,
                },
                {
                  expr: "sum(k8s_pod_rate_cpu_core_used_limit)",
                  thousands: 1,
                  alias: "CPU限制",
                  chartType: ChartType.Area,
                  scale: [0, 45000, 90000, 135000, 180000, 225000]
                }
              ],
              conditions: []
            },
            {
              table: "k8s_pod",
              fields: [
                {
                  expr: "sum(k8s_pod_cpu_core_used)",
                  unit: "%",
                  thousands: 1024,
                  alias: "CPU使用率"
                },
                {
                  expr: "sum(k8s_pod_rate_cpu_core_used_limit)",
                  thousands: 1,
                  alias: "CPU限制",
                  chartType: ChartType.Area
                }
              ],
              conditions: []
            }
          ]}
          groupBy={groupBy}
          conditions={conditions}
        />
        <ChartInstancesPanel
          tables={[
            {
              table,
              fields,
              conditions
            }
          ]}
          groupBy={groupBy}
          instance={instance}
        />
        <br />
      </div>
    );
  }
}

Entry.register({
  businessKey: location.pathname.split("/")[1],
  routes: {
    index: {
      title: "概览 - argus",
      Component: Chart
    },

    explore: {
      title: "数据探索 - argus",
      Component: Chart
    }
  }
});
