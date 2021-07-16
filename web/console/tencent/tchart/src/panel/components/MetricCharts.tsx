import * as React from "react";
import Chart from "charts/index";
import { Kilobyte } from "../constants";
import { CHART_PANEL } from "../core";


interface MetricChartsProps {
  className: string;
  style: React.CSSProperties;
  loading: boolean;
  min: number;
  max: number;
  seriesGroup: Array<CHART_PANEL.ChartDataType>;
  reload?: Function;
}

interface MetricChartsState {
  chartEndIndex: number;
}

export default class MetricCharts extends React.Component<MetricChartsProps, MetricChartsState> {
  ChartPageSize = 3; // panel 为 list 状态时每次滚动加载 chart 数目
  private _id: string = "MetricCharts";
  private _charts = [];
  private _chartHeight = 388;
  private _chartWidth = 726;

  constructor(props) {
    super(props);
    this.state = {
      chartEndIndex: this.ChartPageSize
    };
  }

  componentDidMount() {
    this.scrollUpdateCharts(this.props, this.state.chartEndIndex);
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (nextProps.seriesGroup.length !== this.props.seriesGroup.length
      || nextState.chartEndIndex !== this.state.chartEndIndex) {
      this.scrollUpdateCharts(nextProps, nextState.chartEndIndex);
      return true;
    }
    if (JSON.stringify(nextProps.seriesGroup) !== JSON.stringify(this.props.seriesGroup)
      || nextProps.min !== this.props.min
      || nextProps.max !== this.props.max
      || nextProps.loading !== this.props.loading) {
      this.updateCharts(nextProps, nextState.chartEndIndex);
      return true
    }
    return false;
  }

  updateCharts(props, chartEndIndex) {
    const {loading, min, max, seriesGroup, reload} = props;
    const length = Math.min(seriesGroup.length, chartEndIndex);
    for (let i = 0; i < length; i++) {
      const series = seriesGroup[i].lines;
      const title = seriesGroup[i].title;
      const labels = seriesGroup[i].labels;
      const field = seriesGroup[i].field;
      const tooltipLabels = seriesGroup[i].tooltipLabels;
      if (this._charts[i]) {
        this._charts[i].setType(field.chartType, {
          loading,
          min,
          max,
          title,
          labels,
          series,
          reload,
          tooltipLabels,
          colorTheme: field.colorTheme,
          colors: field.colors,
          yAxis: field.scale || [],
          isKilobyteFormat: field.thousands === Kilobyte,
          unit: field.unit
        });
      } else {
        this._charts.push(
          new Chart(
            `${this._id}_${i}`,
            {
              width: this._chartWidth,
              height: this._chartHeight,
              paddingHorizontal: 15, // 水平空白间隔
              paddingVertical: 30, // 垂直空白间隔
              isSeriesTime: true,
              showLegend: false,
              isKilobyteFormat: field.thousands === Kilobyte,
              title,
              min,
              max,
              tooltipLabels,
              colorTheme: field.colorTheme,
              colors: field.colors,
              yAxis: field.scale || [],
              unit: field.unit,
              loading,
              labels,
              series,
              reload
            },
            field.chartType
          )
        );
      }
    }
  }

  scrollUpdateCharts(props, chartEndIndex) {
    const {loading, min, max, seriesGroup, reload} = props;
    const startIndex = this._charts.length;
    const length = Math.min(seriesGroup.length, chartEndIndex);
    for (let i = startIndex; i < length; i++) {
      const series = seriesGroup[i].lines;
      const title = seriesGroup[i].title;
      const labels = seriesGroup[i].labels;
      const field = seriesGroup[i].field;
      const tooltipLabels = seriesGroup[i].tooltipLabels;
      this._charts.push(
        new Chart(
          `${this._id}_${i}`,
          {
            width: this._chartWidth,
            height: this._chartHeight,
            paddingHorizontal: 15, // 水平空白间隔
            paddingVertical: 30, // 垂直空白间隔
            isSeriesTime: true,
            showLegend: false,
            isKilobyteFormat: field.thousands === Kilobyte,
            title,
            min,
            max,
            tooltipLabels,
            yAxis: field.scale || [],
            unit: field.unit,
            loading,
            labels,
            series,
            reload
          },
          field.chartType
        )
      );
    }
  }

  render() {
    const {className, style} = this.props;

    return (
      <div
        className={className}
        style={{...style}}
        onScroll={(e) => {
          // 图表滚动加载
          let element = e.target as any;
          // 向下滚动加载
          // 提前加载
          const index = this.state.chartEndIndex - 1;
          if (Math.floor(this._chartHeight * index - element.scrollTop) <= element.clientHeight) {
            if (this.state.chartEndIndex < this.props.seriesGroup.length) {
              const nextChartEndIndex = this.state.chartEndIndex + this.ChartPageSize;
              this.setState({chartEndIndex: nextChartEndIndex});
            }
          }
        }}
      >
        {
          this.props.seriesGroup.map((item, index) => {
            const id = `${this._id}_${index}`;
            return (
              <div
                key={id}
                style={{width: "100%", height: this._chartHeight}}
              >
                <div id={id}></div>
              </div>
            )
          })
        }
      </div>
    )
  }
}