import * as React from "react";
import { Icon, Table, StatusTip, TableAddon, Button, Form, Select, TagSearchBox } from "tea-component";
import { TagSearchBox as SingleTagSearchBox , } from "tea-component";
import Chart from "../../charts/index";
import { ModelType } from "../../core/model";
import * as utils from "../../core/utils";
import PureChart from "./PureChart";
import { STORE } from "../helper";
import { Kilobyte, QUERY } from "../constants";
import { CHART_PANEL, NameValueType } from "../core";
import * as languages from "../../i18n";



const version = (window as any).VERSION || "zh";
const language = languages[version];

const { scrollable } = Table.addons;

/**
 * 生成系统聚合方式选项
 */
const AggregationOptions = Object.keys(QUERY.Aggregation).map(key => {
  return {
    value: key,
    text: QUERY.Aggregation[key]
  };
});

function FormatTagsParam(tagDimensions, conditions, tagConditions) {
  // 参考 TagSearchBox 返回的参数
  const dimensionItems = tagDimensions.map(item => item.attr.value);
  // condition 支持多条件查询
  const conditionItems = [].concat(
    conditions,
    tagConditions.map(item => [item.attr.value, "in", item.values.map(value => value.name)])
  );
  return { dimensionItems, conditionItems };
}

interface FilterTableChartProps {
  chartId: string;
  width?: number;
  height?: number;
  table: string;
  startTime: Date;
  endTime: Date;
  metric: CHART_PANEL.MetricType;
  tooltipLabels: (value, from) => string;
  dimensions: Array<NameValueType>;
  conditions: Array<Array<any>>;
  periodOptions: Array<{ value: string; text: string }>;
}

interface FilterTableChartState {
  loading: boolean;
  aggregation: string;
  tagDimensions: Array<any>;
  tagConditions: Array<any>;
  period: string;
  labels: Array<string>;
  lines: Array<ModelType>;
}

export class FilterTableChart extends React.Component<FilterTableChartProps, FilterTableChartState> {
  constructor(props) {
    super(props);
    let tagDimensions = [];
    let tagConditions = [];
    // 判断 metric 是否有 storeKey，
    // 处理 defaultGroupBy, defaultConditions 搜索条件初始化 default* 的数据结构：
    if (props.metric.storeKey) {
      const { defaultGroupBy, defaultConditions } = STORE.Get(props.metric.storeKey, {
        defaultGroupBy: [],
        defaultConditions: []
      });
      tagDimensions = defaultGroupBy.map(item => {
        const { groupBy, values } = item;
        return {
          attr: {
            type: "onlyKey",
            key: groupBy.value,
            name: groupBy.name,
            value: groupBy.value
          },
          values
        };
      });
      tagConditions = defaultConditions.map(item => {
        const { groupBy, values } = item;
        return {
          attr: {
            type: "input",
            key: groupBy.value,
            name: groupBy.name,
            value: groupBy.value
          },
          values
        };
      });
    }
    this.state = {
      loading: false,
      tagDimensions: tagDimensions,
      tagConditions: tagConditions,
      aggregation: AggregationOptions[0].value,
      period: props.periodOptions[0].value,
      labels: [],
      lines: []
    };
  }

  componentWillMount() {
    // 加载数据
    this.fetchChartData();
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    let hasUpdate = false;
    let hasFetch = false;
    const { table, startTime, endTime, metric, conditions, periodOptions } = this.props;
    if (
      JSON.stringify(startTime) !== JSON.stringify(prevProps.startTime) ||
      JSON.stringify(endTime) !== JSON.stringify(prevProps.endTime) ||
      JSON.stringify(metric) !== JSON.stringify(prevProps.metric) ||
      JSON.stringify(conditions) !== JSON.stringify(prevProps.conditions)
    ) {
      hasFetch = true;
    }

    // 更新 period 选项
    let period = this.state.period;
    if (JSON.stringify(periodOptions) !== JSON.stringify(prevProps.periodOptions)) {
      period = periodOptions[0].value;
      this.setState({ period });
      hasUpdate = true;
    }

    const { dimensionItems, conditionItems } = FormatTagsParam(
      this.state.tagDimensions,
      conditions,
      this.state.tagConditions
    );
    if (JSON.stringify(prevState) !== JSON.stringify(this.state)) {
      hasUpdate = true;
    }

    if (hasFetch) {
      this.requestChartData(table, startTime, endTime, metric, dimensionItems, conditionItems, period);
    }

    return hasUpdate || hasFetch;
  }

  fetchChartData() {
    const { table, startTime, endTime, metric, conditions } = this.props;
    const { tagDimensions, tagConditions, period } = this.state;
    const { dimensionItems, conditionItems } = FormatTagsParam(tagDimensions, conditions, tagConditions);
    this.requestChartData(table, startTime, endTime, metric, dimensionItems, conditionItems, period);
  }

  requestChartData(table, startTime, endTime, metric, dimensions, conditions, period) {
    this.setState({ loading: true, labels: [], lines: [] });
    // 生成请求字段
    const aggregationFields = [{ ...metric, expr: `${this.state.aggregation}(${metric.expr})` }];
    CHART_PANEL.RequestData({
      table,
      startTime,
      endTime,
      fields: aggregationFields,
      dimensions: dimensions,
      conditions: conditions,
      period: period
    })
      .then(res => {
        const { columns, data } = res;
        // 根据维度聚合返回数据
        const dimensionsData = CHART_PANEL.AggregateDataByDimensions(dimensions, columns, data);
        /**
         * 生成 labels
         */
        const timestampIndex = columns.indexOf(`timestamp(${period}s)`);
        // 默认表格第一列为时间序列, 返回的时间列表可能不是startTime开始，以endTime结束，需要补帧
        let labels = (Object.values(dimensionsData)[0] || ([] as any)).map(item => item[timestampIndex]);
        labels = CHART_PANEL.OffsetTimeSeries(startTime, endTime, parseInt(period), labels);
        /**
         * 生成 lines
         */
        const metricIndex = CHART_PANEL.FindMetricIndex(columns, `${metric.expr}_${this.state.aggregation}`); // 数据列下标
        const valueTransform = metric.valueTransform || (value => value);
        // 每个图表中显示groupBy类别一个field -> chart
        let lines = [];
        Object.keys(dimensionsData).forEach(dimensionKey => {
          // 每个dimension的数据集， fieldsIndex获取该dimension列下标的值
          const sourceValue = dimensionsData[dimensionKey];
          const data = {};
          sourceValue.forEach(row => {
            let value = row[metricIndex];
            // data 记录 时间戳row[timestampIndex] 对应 value
            data[row[timestampIndex]] = valueTransform(value);
          });

          let line = CHART_PANEL.GenerateLine(dimensionKey, data);
          lines.push(line);
        });

        this.setState({ loading: false, labels, lines });
      })
      .catch(e => {
        this.setState({ loading: false });
      });
  }

  updateState(options) {
    this.setState({ ...options }, () => {
      this.fetchChartData();
    });
  }

  saveTagsByStoreKey(saveKey: "defaultGroupBy" | "defaultConditions", tags) {
    if (this.props.metric.storeKey) {
      let defaultValues = STORE.Get(this.props.metric.storeKey, { defaultGroupBy: [], defaultConditions: [] });
      defaultValues[saveKey] = tags.map(item => {
        return {
          groupBy: item.attr,
          values: item.values
        };
      });
      STORE.Set(this.props.metric.storeKey, defaultValues);
    }
  }

  render() {
    const { loading, tagDimensions, tagConditions, period, aggregation, labels, lines } = this.state;
    const { dimensions, periodOptions, chartId, width, height, metric, startTime, endTime } = this.props;
    // 生成维度和筛选条件选项
    const dimensionOptions = dimensions.map(item => {
      return {
        type: "single",
        key: item.value,
        name: item.name,
        value: item.value
      };
    });
    const conditionOptions = dimensions.map(item => {
      return {
        type: "input",
        key: item.value,
        name: item.name,
        value: item.value
      };
    });

    return (
      <div style={{ width: "100%" }}>
        <Form className="tea-form--vertical tea-form--inline tea-mt-2n">
          <Form.Item label="维度" style={{ width: "calc(100% - 300px)" }}>
            <SingleTagSearchBox
              minWidth={"100%"}
              attributes={dimensionOptions as any}
              tips={"请选择维度"}
              value={tagDimensions}
              onChange={tags => {
                if (dimensionOptions.length === 0) {
                  return;
                }
                const tagDimensions = tags.map(item => {
                  return {
                    attr: item.attr,
                    values: item.values
                  };
                });
                // 如果指标传递了storeKey，则缓存用户维度的选项，方便用户下次操作
                this.saveTagsByStoreKey("defaultGroupBy", tags);
                this.updateState({ tagDimensions });
              }}
            />
            
          </Form.Item>
          <Form.Item label="统计粒度">
            <Select
              style={{ paddingTop: 6 }}
              type="simulate"
              appearence="button"
              options={periodOptions}
              value={period}
              onChange={value => {
                this.updateState({ period: value });
              }}
              placeholder="请选择"
            />
          </Form.Item>
          <Form.Item label="统计方式">
            <Select
              style={{ paddingTop: 6 }}
              type="simulate"
              appearence="button"
              options={AggregationOptions}
              value={aggregation}
              onChange={value => {
                this.updateState({ aggregation: value });
              }}
              placeholder="请选择"
            />
          </Form.Item>
        </Form>
        <Form className="tea-form--vertical tea-form--inline tea-mt-2n">
          <Form.Item label="筛选条件" style={{ width: "100%" }}>
            <TagSearchBox
              minWidth={"100%"}
              attributes={conditionOptions as any}
              value={tagConditions}
              onChange={tags => {
                if (conditionOptions.length === 0) {
                  return;
                }
                const tagConditions = tags.map(item => {
                  return {
                    attr: item.attr,
                    values: item.values
                  };
                });
                this.saveTagsByStoreKey("defaultConditions", tags);
                this.updateState({ tagConditions });
              }}
            />
          </Form.Item>
        </Form>

        <TableChart
          id={chartId}
          loading={loading}
          min={typeof startTime === "string" ? 0 : startTime.getTime()}
          max={typeof endTime === "string" ? 0 : endTime.getTime()}
          width={width}
          height={height}
          title={metric.alias}
          field={metric}
          labels={labels}
          unit={metric.unit}
          tooltipLabels={metric.tooltipLabels}
          lines={lines}
          reload={this.fetchChartData.bind(this)}
        />
      </div>
    );
  }
}

interface TableChartProps extends CHART_PANEL.ChartDataType {
  id: string;
  width?: number;
  height?: number;
  loading: boolean;
  reload: Function;
}

interface TableChartState {
  selectedLabel: number;
  selectedRowKeys: Array<string>;
}

/**
 * 表格型图表
 * hover 图表时数据通过 table 显示，隐藏 tooltip
 */
class TableChart extends PureChart<TableChartProps, TableChartState> {
  _scrollAnchorRefs = {};

  constructor(props) {
    super(props);
    this.state = {
      selectedLabel: 0,
      selectedRowKeys: []
    };
  }

  componentDidMount() {
    const { loading, min, max, title, reload, labels, lines, field, tooltipLabels } = this.props;
    this._chartEntry = new Chart(
      this.props.id,
      {
        width: this.chartWidth,
        height: this.chartHeight,
        paddingHorizontal: 15, // 水平空白间隔
        paddingVertical: 30, // 垂直空白间隔
        isSeriesTime: true,
        showLegend: false,
        showTooltip: false,
        hoverPointData: this.mouseOverChart.bind(this),
        loading,
        min,
        max,
        title,
        labels,
        reload,
        tooltipLabels,
        colorTheme: field.colorTheme,
        colors: field.colors,
        yAxis: field.scale || [],
        series: lines,
        isKilobyteFormat: field.thousands === Kilobyte,
        unit: field.unit
      },
      field.chartType
    );
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (JSON.stringify(nextProps) !== JSON.stringify(this.props)) {
      const { loading, min, max, title, reload, labels, lines, field, tooltipLabels } = nextProps;
      this._chartEntry &&
        this._chartEntry.setType(field.chartType, {
          width: this.chartWidth,
          height: this.chartHeight,
          loading,
          min,
          max,
          title,
          labels,
          reload,
          tooltipLabels,
          colorTheme: field.colorTheme,
          colors: field.colors,
          yAxis: field.scale || [],
          series: lines,
          isKilobyteFormat: field.thousands === Kilobyte,
          unit: field.unit
        });
      return true;
    }
    if (JSON.stringify(nextState) !== JSON.stringify(this.state)) {
      return true;
    }
    return false;
  }

  // 图表鼠标悬浮时回调函数
  mouseOverChart(params) {
    const { xAxisTickMarkIndex, mousePosition, content } = params;
    // 筛选出被hover的折线
    const selectedRowKeys = content.filter(item => item.hover).map(item => item.legend);
    // table 中对应数据进行滚动显示
    if (selectedRowKeys.length > 0) {
      const element = this._scrollAnchorRefs[selectedRowKeys[selectedRowKeys.length - 1]].current;
      // firefox,IE 不支持 scrollIntoViewIfNeeded
      if (element.scrollIntoViewIfNeeded) {
        element.scrollIntoViewIfNeeded(false);
      } else {
        element.scrollIntoView({ block: "end", behavior: "smooth" });
      }
    }
    this.setState({ selectedLabel: xAxisTickMarkIndex, selectedRowKeys });
  }

  // 根据鼠标选中的label，获取各折线对应的点，显示在表格中
  getTableRecords(label: string, lines: Array<any>) {
    const tableRecords = lines.map(line => {
      return {
        key: utils.FormatStringNoHTMLSharp(line.legend),
        value: line.data[label]
      };
    });
    tableRecords.sort((a, b) => b.value - a.value);
    return tableRecords;
  }

  // hover 在表格的事件
  hoverTableEvent(legend: string) {
    this._chartEntry && this._chartEntry.highlightLine(legend);
  }

  render() {
    const { loading, labels, lines } = this.props;
    const tableRecords = this.getTableRecords(labels[this.state.selectedLabel], lines);

    return (
      <div style={{ width: "100%" }}>
        {/* 图表div，id为图表容器标识 */}
        <div style={{ width: "100%" }} id={this.props.id} />

        <div className="tea-justify-grid tea-mt-2n">
          <div className="tea-justify-grid__col tea-justify-grid__col--left">
            <h3 className="tea-h3" style={{ fontSize: 14 }}>
              {language.DataDetail} &nbsp;&nbsp;
              {labels[this.state.selectedLabel] &&
                utils.TIME.Format(labels[this.state.selectedLabel] as any, utils.TIME.DateFormat.fullDateTime)}
            </h3>
          </div>
        </div>
        <div style={{ width: "calc(100% - 40px)", margin: "0 auto" }}>
          <Table
            columns={[
              {
                key: "key",
                header: "Source",
                render: record => {
                  this._scrollAnchorRefs[record.key] = React.createRef();
                  return <span ref={this._scrollAnchorRefs[record.key]}>{record.key}</span>;
                }
              },
              {
                key: "value",
                header: "Value"
              }
            ]}
            records={tableRecords}
            recordKey={"key"}
            rowClassName={record => (this.state.selectedRowKeys.indexOf(record.key) !== -1 ? "is-selected" : "")}
            topTip={
              loading && (
                <StatusTip
                  // @ts-ignore
                  status={loading ? "loading" : "none"}
                />
              )
            }
            addons={[
              // 支持表格滚动，高度超过 180 开始显示滚动条
              scrollable({
                maxHeight: 180,
                onScrollBottom: () => {}
              }),
              hoverRowEvent({ event: this.hoverTableEvent.bind(this) })
            ]}
          />
        </div>
      </div>
    );
  }
}

/**
 * 注册 table hover 事件
 */
function hoverRowEvent(options: any): TableAddon {
  return {
    onInjectRow: next => (...args) => {
      const result = next(...args);
      const [record, rowKey] = args;
      return {
        ...result,
        row: React.cloneElement(result.row, {
          onMouseOver: e => {
            options.event(rowKey);
          },
          onMouseLeave: e => {
            options.event("");
          }
        })
      };
    }
  };
}
