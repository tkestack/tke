import { CHART_PANEL } from 'panel/core';
import * as React from 'react';
import Chart from "charts/index";
import { Kilobyte, CHART } from '../constants';


export interface PureChartProps extends CHART_PANEL.ChartDataType {
  id: string;
  width?: number;
  height?: number;
  style?: React.CSSProperties;
  loading: boolean;
  reload?: Function;
}

/**
 * 基础图表组件，支持缩放
 */
export default class PureChart<P extends PureChartProps, S> extends React.Component<P, S> {
  _width = 0;
  _height = 0;
  _chartEntry = null;

  constructor(props) {
    super(props);
  }

  get chartWidth() {
    if (this.props.width) {
      return this.props.width;
    }
    if (!document.getElementById(this.props.id)) {
      this._width = CHART.DefaultSize.width;
    } else {
      const width = document.getElementById(this.props.id).clientWidth;
      this._width = width > 0 ? width : CHART.DefaultSize.width;
    }
    return this._width;
  }

  get chartHeight() {
    if (this.props.height) {
      return this.props.height;
    }
    if (this._width < CHART.DefaultSize.width) {
      this._height = CHART.DefaultSize.height * this._width / CHART.DefaultSize.width;
    } else {
      this._height = CHART.DefaultSize.height;
    }
    return this._height;
  }

  componentWillMount() {
    window.addEventListener('resize', this.onResize);
  }

  componentDidMount() {
    const {loading, title, min, max, reload, labels, lines, field, tooltipLabels} = this.props;
    this._chartEntry = new Chart(this.props.id,
      {
        width: this.chartWidth,
        height: this.chartHeight,
        paddingHorizontal: 15, // 水平空白间隔
        paddingVertical: 30, // 垂直空白间隔
        isSeriesTime: true,
        showLegend: false,
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
      field.chartType);
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (JSON.stringify(nextProps) !== JSON.stringify(this.props)) {
      const {loading, title, min, max, reload, labels, lines, field, tooltipLabels} = nextProps;
      this._chartEntry && this._chartEntry.setType(field.chartType, {
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

  componentWillUnmount() {
    window.removeEventListener('resize', this.onResize);
  }

  onResize = () => {
    this._chartEntry && this._chartEntry.setSize(this.chartWidth, this.chartHeight);
  };

  render() {
    const {style = {}} = this.props;
    return (
      <div style={{width: '100%', height: CHART.DefaultSize.height, ...style}} id={this.props.id}>
      </div>
    )
  }
}