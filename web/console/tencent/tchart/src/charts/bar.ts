import Graph from "core/graph";


/**
 * 只有叠加状态
 */
export default class BarChart extends Graph {

  constructor(id: string, options) {
    super(id, {...options, overlay: true});
  }

  setOptions(options) {
    super.setOptions({...options, overlay: true});
  }

  draw(): void {
    this._mainPanel.clearRect(0, 0, this._options.width, this._options.height);

    this.drawLoading();

    this.drawTitle();
    // 绘画坐标轴
    this.drawAxis();
    // 绘画等高线
    this.drawGrid();
    this.drawLabels();
    this.drawLegends();
    this.drawEmptyData();

    this.drawBar();
  }

}