import * as React from "react";
import { Toolbar } from "../components/toolbar";
import MetricCharts from "../components/MetricCharts";
import InstanceList, { ColumnType } from "../components/InstanceList";
import { TIME_PICKER, QUERY } from "../constants";
import { NameValueType, CHART_PANEL } from "../core";
import { Period, TransformField } from "../helper";
import { Select } from "tea-component";

require("./Instances.less");

export interface ChartInstancesPanelProps {
  width?: number;
  height?: number;
  tables: Array<CHART_PANEL.TableType>;
  groupBy: Array<NameValueType>;
  instance: {
    // 类table中column和data
    columns: Array<ColumnType>;
    list: Array<any>;
  };
  projectId: string;
  platformType: any;
  children?: React.ReactNode;
}

interface ChartInstancesPanelState {
  loading: boolean;
  startTime: Date;
  endTime: Date;
  seriesGroup: Array<CHART_PANEL.ChartDataType>;
  aggregation: string;
  instanceList?: Array<any>;
  periodOptions: Array<{ value: string; text: string }>;
  period: string;
}

export class ChartInstancesPanel extends React.Component<
  ChartInstancesPanelProps,
  ChartInstancesPanelState
> {
  private _instances = {}; // group by 查询返回数据中的不同类别的groupBy的值

  constructor(props) {
    super(props);
    const endTime = new Date();
    const startTime = new Date(endTime.getTime() - 1000 * 60 * 60);
    const periodOptions = CHART_PANEL.GeneratePeriodOptions(startTime, endTime);
    this.state = {
      loading: false,
      startTime: TIME_PICKER.Tabs[0].from as any,
      endTime: TIME_PICKER.Tabs[0].to as any,
      seriesGroup: [],
      instanceList: props.instance.list,
      aggregation: Object.keys(QUERY.Aggregation)[0],
      periodOptions: periodOptions,
      period: periodOptions[0].value,
    };
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    let hasUpdate = false;
    let hasFetch = false;
    const { tables, groupBy, instance } = this.props;
    let instanceList = [];
    if (instance) {
      instanceList = this.props.instance.list;
      if (JSON.stringify(instance) !== JSON.stringify(prevProps.instance)) {
        hasUpdate = true;
      }
    }
    if (
      JSON.stringify(groupBy) !== JSON.stringify(prevProps.groupBy) ||
      JSON.stringify(tables) !== JSON.stringify(prevProps.tables)
    ) {
      hasFetch = true;
    }

    if (hasFetch) {
      this.fetchDataByTables({ instances: instanceList });
    } else if (hasUpdate) {
      instanceList = this.combineInstancesWithGroupBy(instance.list);
      this.setState({ instanceList });
    }

    if (JSON.stringify(prevState) !== JSON.stringify(this.state)) {
      hasUpdate = true;
    }

    return hasUpdate || hasFetch;
  }

  /**
   * 根据key值判断chart折线是否隐藏
   * @param instanceList
   * @param key
   * @returns {boolean}
   */
  checkLineIsDisable(instanceList = [], key = "") {
    if (instanceList.length === 0) {
      return false;
    }
    const instance = instanceList.find((instance) => {
      const instanceKey = CHART_PANEL.GenerateRowKey(
        this.props.groupBy.map((groupBy) => instance[groupBy.value])
      );
      return instanceKey.indexOf(key) !== -1;
    });
    if (!instance || !instance.isChecked) {
      return true;
    }
    return false;
  }

  combineInstancesWithGroupBy(instances) {
    // this.state.instanceList 为用户通过接口传入的值 例如结构 { workload_kind: "Deployment", isChecked: true },
    let instanceList = instances || [].concat(this.state.instanceList);
    /**
     * 根据 groupByData 字典中key的值，查找instances中是否匹配
     * 对instances和groupBy做交集
     */
    Object.keys(this._instances).forEach((instanceKey) => {
      let instance = instanceList.find((instance) => {
        // 通过groupBy的key在 this.state.instanceList 的值中获取值列表,
        // 如：this.props.groupBy=[workload_kind], instance={workload_kind: "Deployment", isChecked: true}
        // groupByValues 为 [Deployment]
        const groupByValues = this.props.groupBy.map(
          (groupByItem) => instance[groupByItem.value]
        );
        const rowKey = CHART_PANEL.GenerateRowKey(groupByValues);
        return instanceKey === rowKey;
      });
      if (!instance) {
        // 用户配置的instances值通过groupBy查询找不到时，则新增instance
        instanceList.push(this._instances[instanceKey]);
      }
    });
    return instanceList;
  }

  /**
   * 获取空对象
   */
  getEmptyData(fields = []) {
    let seriesGroup = [];
    const tables = fields.length > 0 ? [{ fields }] : this.props.tables;
    tables.forEach((tableInfo) => {
      tableInfo.fields.forEach((field) => {
        seriesGroup.push({
          field,
          title: field.alias,
          labels: [],
          lines: [],
        });
      });
    });
    return seriesGroup;
  }

  onChangeQueryTime(startTime, endTime) {
    if (this.state.loading) {
      return;
    }
    const periodOptions = CHART_PANEL.GeneratePeriodOptions(startTime, endTime);
    const period = periodOptions[0].value;
    this.setState({ startTime, endTime, periodOptions, period }, () => {
      // 选择时间段后重新请求数据
      this.fetchDataByTables({ instances: this.state.instanceList });
    });
  }

  onChangePeriod(period) {
    this.setState({ period }, () => {
      // 选择时间粒度后重新请求数据
      this.fetchDataByTables({ instances: this.state.instanceList });
    });
  }

  /**
   * 勾选instances。更新到state中
   */
  onCheckInstances = (instanceList: Array<any>) => {
    this.setState({ instanceList }, () => {
      const seriesGroupTemp = JSON.parse(
        JSON.stringify(this.state.seriesGroup)
      );
      const seriesGroup = seriesGroupTemp.map(
        (series: CHART_PANEL.ChartDataType) => {
          // series 是每一个图表的数据
          series.lines.forEach((line) => {
            // line 是每个图表中每条线的数据
            line.disable = this.checkLineIsDisable(
              this.state.instanceList,
              line.legend
            );
          });
          return series;
        }
      );
      // 对象深拷贝有问题
      this.setState({ seriesGroup });
    });
  };

  fetchDataByTables(options: any = {}) {
    this._instances = {}; // 请求时清除原来数据加载的 _instances 对象
    // 初始状态
    this.setState({
      seriesGroup: this.getEmptyData(),
      instanceList: this.props.instance.list,
      loading: true,
    });

    const { tables } = this.props;
    const requestPromises = tables.map((tableInfo) => {
      return this.requestData(tableInfo);
    });
    Promise.all(requestPromises).then((values) => {
      const seriesGroup = [].concat(...values);
      // combineInstancesWithGroupBy 合并 this._instances 与 instances
      let instanceList = this.combineInstancesWithGroupBy(
        this.props.instance.list
      );
      this.setState({ seriesGroup, instanceList, loading: false });
    });
  }

  /**
   * 发起查询数据请求
   * @param {string} table
   * @param {Array<any>} fields
   * @param {Array<any>} groupBy
   * @param {Array<any>} conditions
   * @param {string} period
   */
  requestData(tableInfo: any = {}) {
    const { fields, conditions, table } = tableInfo;
    const { startTime, endTime, period } = this.state;

    const groupByItems = this.props.groupBy.map((item) => item.value) || [];
    // 查询的时间粒度
    // const period = Period(startTime, endTime) as any;

    return CHART_PANEL.RequestData({
      table: table,
      startTime,
      endTime,
      fields: fields,
      dimensions: groupByItems,
      conditions: conditions,
      period: period,
    })
      .then((res) => {
        const { columns, data } = res;
        /**
         * 根据 groupBy 条件对数据做聚合，已 groupBy 为key进行存储
         */
        const dimensionsData = CHART_PANEL.AggregateDataByDimensions(
          groupByItems,
          columns,
          data
        );

        /**
         * 根据 groupBy 数据生成默认 instances
         */
        const columnsInfo = groupByItems.map((value) => {
          return {
            index: columns.indexOf(value),
            value,
          };
        });
        Object.keys(dimensionsData).forEach((groupByKey) => {
          //从第一行的数据获取group by的值
          const row = dimensionsData[groupByKey][0];
          // 缓存数据中groupBy联合主键的instance list
          if (!this._instances[groupByKey]) {
            let instance = {}; // 初始化 instance 对象
            columnsInfo.forEach((column) => {
              const key = column.value;
              // key 为groupBy的参数，row[column.index]对应groupBy的值 如{workload_kind: "Deployment"}
              instance[key] = row[column.index];
            });
            // 如果传入的instance list 为空数组是，数据中的instance对象为勾选状态
            instance["isChecked"] =
              this.props.instance && this.props.instance.list.length > 0
                ? false
                : true;
            this._instances[groupByKey] = instance;
          }
        });

        /**
         * 生成 labels
         */
        const timestampIndex = columns.indexOf(`timestamp(${period}s)`);
        // 默认表格第一列为时间序列, 返回的时间列表可能不是startTime开始，以endTime结束，需要补帧
        let labels = (Object.values(dimensionsData)[0] || ([] as any)).map(
          (item) => item[timestampIndex]
        );
        labels = CHART_PANEL.OffsetTimeSeries(
          startTime,
          endTime,
          parseInt(period),
          labels
        );

        /**
         * 生成 series
         */
        let seriesGroup = [];
        // 获取需要显示的field的数据，每个field是一个图表。
        fields.forEach((field) => {
          const fieldIndex = CHART_PANEL.FindMetricIndex(columns, field.expr); // 数据列下标

          const valueTransform = field.valueTransform || ((value) => value);
          // 每个图表中显示groupBy类别一个field -> chart
          let lines = [];
          Object.keys(dimensionsData).forEach((groupByKey) => {
            // 每个groupBy的数据集， fieldsIndex的index获取groupBy该下标的值
            const sourceValue = dimensionsData[groupByKey];
            const data = {};
            sourceValue.forEach((row) => {
              let value = row[fieldIndex];
              // data 记录 时间戳row[timestampIndex] 对应 value
              data[row[timestampIndex]] = valueTransform(value);
            });

            let line = CHART_PANEL.GenerateLine(
              groupByKey,
              data,
              this.checkLineIsDisable(this.props.instance.list, groupByKey)
            );
            lines.push(line);
          });

          // 每个chart 的数据存储到 seriesGroup
          seriesGroup.push({
            labels,
            field,
            lines,
            title: field.alias,
            tooltipLabels: (value, from) => {
              return field.valueLabels
                ? field.valueLabels(value, from)
                : `${TransformField(value || 0, field.thousands, 3)}`;
            },
          } as CHART_PANEL.ChartDataType);
        });
        return seriesGroup;
      })
      .catch((e) => {
        return this.getEmptyData(fields);
      });
  }

  render() {
    const {
      loading,
      seriesGroup,
      instanceList,
      startTime,
      endTime,
      periodOptions,
      period,
    } = this.state;

    return (
      <div className="tc-15-rich-dialog-bd">
        <Toolbar
          duration={{ from: this.state.startTime, to: this.state.endTime }}
          onChangeTime={this.onChangeQueryTime.bind(this)}
        >
          <span
            style={{
              fontSize: 12,
              display: "inline-block",
              verticalAlign: "middle",
            }}
          >
            统计粒度：
          </span>
          <Select
            type="native"
            size="s"
            appearence="button"
            options={periodOptions}
            value={period}
            onChange={(value) => {
              this.onChangePeriod(value);
            }}
            placeholder="请选择"
          />
          {this.props.children}
        </Toolbar>

        <div className="monitor-dialog-data-box">
          <div className="tc-g">
            <div className="tc-g-u-3-4">
              <MetricCharts
                className="monitor-chart-grid"
                style={{
                  maxHeight: "calc(100vh - 200px)",
                  height: "668px",
                  overflowY: "auto",
                  overflowX: "hidden",
                  display: "block",
                  border: "1px solid #ddd",
                }}
                loading={loading}
                min={typeof startTime === "string" ? 0 : startTime.getTime()}
                max={typeof endTime === "string" ? 0 : endTime.getTime()}
                seriesGroup={seriesGroup}
                reload={this.fetchDataByTables.bind(this, {
                  instances: this.state.instanceList,
                })}
              />
            </div>

            <div className="tc-g-u-1-4">
              <div className="tc-15-mod-selector-tb">
                <div className="tc-15-option-cell options-left">
                  <div className="tc-15-option-bd">
                    <InstanceList
                      className="tc-15-option-box tc-scroll"
                      style={{
                        maxHeight: "calc(100vh - 200px)",
                        height: "668px",
                      }}
                      columns={this.props.instance.columns}
                      list={instanceList}
                      update={this.onCheckInstances}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
