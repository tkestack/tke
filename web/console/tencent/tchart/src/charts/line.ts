import Graph from 'core/graph';


/**
 * activeHover: 支持鼠标悬浮时高亮折线
 */
export default class LineChart extends Graph {

  constructor(id: string, options) {
    super(id, {activeHover: true, ...options});
  }

  setOptions(options) {
    super.setOptions({...options});
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

    // 绘画数据折线
    this.drawLine();
  }

}