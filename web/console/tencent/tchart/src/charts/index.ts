import LineChart from "./line";
import AreaChart from "./area";
import BarChart from "./bar";
import SeriesChart from "./series";
import { COLORS } from "../core/theme";


export const ColorTypes = COLORS.Types;

export { ModelType } from "../core/model";
/**
 * 图表类型
 */
export const ChartType = {
  Line: "line",
  Area: "area",
  Bar: "bar",
  Series: "series"
};

export default class Chart {
  protected _id;
  protected _options = {};
  private _instance;

  constructor(id: string, options, chartType = ChartType.Line) {
    this._id = id;
    this.setType(chartType, options);
  }

  get ID() {
    return this._id;
  }

  setType(chartType, options) {
    let _class = null;
    switch (chartType) {
      case ChartType.Line:
        _class = LineChart;
        break;
      case ChartType.Area:
        _class = AreaChart;
        break;
      case ChartType.Bar:
        _class = BarChart;
        break;
      case ChartType.Series:
        _class = SeriesChart;
        break;
      default:
        _class = SeriesChart;
    }
    // 判断是否需要改变类型
    if (!this._instance || !(this._instance instanceof _class) || !this._instance._parentContainer) {
      this._options = {...this._options, ...options};
      this._instance = new _class(this._id, this._options);
      this.draw();
    } else if (this._instance._parentContainer) {
      this.setOptions(options);
    }
  }

  setOptions(options) {
    this._options = {...this._options, ...options};
    this._instance.setOptions(options);
    this.draw();
  }

  draw() {
    if (this._instance._parentContainer) {
      requestAnimationFrame(() => {
        this._instance.draw();
      });
    }
  }

  public highlightLine(legend: string) {
    this._instance && this._instance.highlightLine(legend);
  }

  public setSize(width: number, height: number) {
    this._instance && this._instance.setChartSize(width, height);
  }

}
