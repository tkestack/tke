import * as React from "react";
import { Card, Form, Select } from "tea-component";
import { Toolbar } from "../components/toolbar";
import PureChart from "../components/PureChart";
import { TransformField } from "../helper";
import { TIME_PICKER, QUERY } from "../constants";
import { CHART_PANEL, NameValueType } from "../core";

export interface ChartPanelProps {
  width?: number;
  height?: number;
  tables: Array<CHART_PANEL.TableType>;
  groupBy: Array<NameValueType>;
  conditions?: Array<Array<string>>;
}

export interface ChartPanelState {
  loading: boolean;
  startTime: Date;
  endTime: Date;
  seriesGroup: Array<CHART_PANEL.ChartDataType>;
  periodOptions: Array<{ value: string; text: string }>;
  aggregation: string;
  period: string;
}

/**
 * 基础 Panel，提供数据请求功能
 */
export class ChartPanel extends React.Component<ChartPanelProps, ChartPanelState> {
  _id: string = "ChartPanel";

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
      aggregation: Object.keys(QUERY.Aggregation)[0],
      periodOptions: periodOptions,
      period: periodOptions[0].value
    };
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    let hasUpdate = false;
    let hasFetch = false;
    const { tables, groupBy } = this.props;
    if (
      JSON.stringify(groupBy) !== JSON.stringify(prevProps.groupBy) ||
      JSON.stringify(tables) !== JSON.stringify(prevProps.tables)
    ) {
      hasFetch = true;
    }

    if (hasFetch) {
      this.fetchDataByTables();
    }

    if (JSON.stringify(prevState) !== JSON.stringify(this.state)) {
      hasUpdate = true;
    }

    return hasUpdate || hasFetch;
  }

  onChangeQueryTime(startTime, endTime) {
    if (this.state.loading) {
      return;
    }
    const periodOptions = CHART_PANEL.GeneratePeriodOptions(startTime, endTime);
    const period = periodOptions[0].value;
    this.setState({ startTime, endTime, periodOptions, period }, () => {
      // 选择时间段后重新请求数据
      this.fetchDataByTables();
    });
  }

  onChangePeriod(period) {
    this.setState({ period }, () => {
      // 选择时间粒度后重新请求数据
      this.fetchDataByTables();
    });
  }

  getEmptyData(fields = []) {
    let seriesGroup = [];
    const tables = fields.length > 0 ? [{ fields }] : this.props.tables;
    tables.forEach(tableInfo => {
      tableInfo.fields.forEach(field => {
        seriesGroup.push({
          field,
          title: field.alias,
          labels: [],
          lines: []
        });
      });
    });
    return seriesGroup;
  }

  fetchDataByTables() {
    // 初始状态
    this.setState({ loading: true, seriesGroup: this.getEmptyData() });
    const { tables } = this.props;

    const requestPromises = tables.map(tableInfo => {
      return this.requestData(tableInfo);
    });
    Promise.all(requestPromises).then(values => {
      const seriesGroup = [].concat(...values);
      this.setState({ seriesGroup, loading: false });
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
  requestData(options: any = {}) {
    const { fields, conditions, table } = options;
    const groupBy = options.groupBy || this.props.groupBy;
    const { startTime, endTime, period } = this.state;
    const groupByItems = groupBy.map(item => item.value) || [];

    return CHART_PANEL.RequestData({
      table: table,
      startTime,
      endTime,
      fields,
      dimensions: groupByItems,
      conditions: conditions,
      period: period
    })
      .then(res => {
        const { columns, data } = res;

        /**
         * 根据 groupBy 条件对数据做聚合，已 groupBy 为key进行存储
         */
        let dimensionsData = CHART_PANEL.AggregateDataByDimensions(groupByItems, columns, data);

        /**
         * 生成 labels
         */
        const timestampIndex = columns.indexOf(`timestamp(${period}s)`);
        // 默认表格第一列为时间序列, 返回的时间列表可能不是startTime开始，以endTime结束，需要补帧
        let labels = (Object.values(dimensionsData)[0] || ([] as any)).map(item => item[timestampIndex]);
        labels = CHART_PANEL.OffsetTimeSeries(startTime, endTime, parseInt(period), labels);

        /**
         * 生成 series
         */
        let seriesGroup = [];
        // 获取需要显示的field的数据，每个field是一个图表。
        fields.forEach(field => {
          const fieldIndex = CHART_PANEL.FindMetricIndex(columns, field.expr); // 数据列下标

          const valueTransform = field.valueTransform || (value => value);
          // 每个图表中显示groupBy类别一个field -> chart
          let lines = [];
          Object.keys(dimensionsData).forEach(groupByKey => {
            // 每个groupBy的数据集， fieldsIndex的index获取groupBy该下标的值
            const sourceValue = dimensionsData[groupByKey];
            const data = {};
            sourceValue.forEach(row => {
              let value = row[fieldIndex];
              // data 记录 时间戳row[timestampIndex] 对应 value
              data[row[timestampIndex]] = valueTransform(value);
            });

            let line = CHART_PANEL.GenerateLine(groupByKey, data);
            lines.push(line);
          });

          const chart: CHART_PANEL.ChartDataType = {
            title: field.alias,
            labels: labels,
            field: field,
            tooltipLabels: (value, from = "tooltip") => {
              return field.valueLabels
                ? field.valueLabels(value, from)
                : `${TransformField(value || 0, field.thousands, 3)}`;
            },
            lines: lines
          };

          // 每个chart 的数据存储到 seriesGroup
          seriesGroup.push(chart);
        });
        return seriesGroup;
      })
      .catch(e => {
        return this.getEmptyData(fields);
      });
  }

  render() {
    const { loading, seriesGroup, periodOptions, period, startTime, endTime } = this.state;

    return (
      <div className="tea-tabs__tabpanel">
        <Toolbar
          style={{ marginLeft: 20 }}
          duration={{ from: this.state.startTime, to: this.state.endTime }}
          onChangeTime={this.onChangeQueryTime.bind(this)}
        >
          <span style={{ fontSize: 12, display: "inline-block", verticalAlign: "middle" }}>统计粒度：</span>
          <Select
            type="native"
            size="s"
            appearence="button"
            options={periodOptions}
            value={period}
            onChange={value => {
              this.onChangePeriod(value);
            }}
            placeholder="请选择"
          />
        </Toolbar>
        {seriesGroup.map((chart, index) => {
          const key = `${this._id}_${index}`;
          return (
            <div key={key} style={{ marginBottom: 10, padding: "0px 10px" }}>
              <PureChart
                id={key}
                loading={loading}
                min={typeof startTime === "string" ? 0 : startTime.getTime()}
                max={typeof endTime === "string" ? 0 : endTime.getTime()}
                width={this.props.width}
                heigth={this.props.height}
                title={chart.title}
                field={chart.field}
                labels={chart.labels}
                unit={chart.field.unit}
                tooltipLabels={chart.tooltipLabels}
                lines={chart.lines}
                reload={this.fetchDataByTables.bind(this)}
              />
            </div>
          );
        })}
      </div>
    );
  }
}
