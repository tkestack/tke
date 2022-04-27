import Paint from "./paint";
import Tooltip, { TooltipContentType } from "./tooltip";
import Event from "./event";
import * as utils from "./utils";
import { COLORS, LINE, CHART } from "./theme";
import { FormatStringNoHTMLSharp, GRAPH } from "./utils";
import { OverlayModel, Model, SeriesModel } from "./model";

/**
 * 设置options后要先进行相关参数的计算，在 setOptions 方法中进行。
 * 把计算流程集中在 setOptions 进行统一管理。
 * 整个图表的核心Graph分为两大块：数据计算和canvas绘图，分别以 draw 和 calculate 为前缀。
 *
 */
export default abstract class Graph {
  protected _options: any = {
    activeColor: "#006eff",
    disabledColor: "#CCC",
    gridColor: "#eee",
    axisColor: "#000",
    auxiliaryLineColor: "#819AA4",
    fontColor: "rgb(124, 134, 142)",
    colorTheme: COLORS.Types.Default,
    // 图表显示基本配置
    width: 600,
    height: 400,
    paddingHorizontal: 35, // 水平空白间隔
    paddingVertical: 40, // 垂直空白间隔
    font: "normal 10px Verdana, Helvetica, Arial, sans-serif",
    labelScale: 10, // X轴默认显示label个数
    gridScale: 5, // Y轴刻度个数
    showLegend: true,
    showYAxis: false,
    activeHover: false,
    hoverPrecision: 2, // hover折线的准确率
    labelAlignCenter: true, // label 默认居中
    showTooltip: true, // 是否使用 tooltip 显示数据
    isKilobyteFormat: false, // tooltip或Y轴刻度是否显示文件大小格式的值
    overlay: false, // 标记图表是否叠加数据
    isSeriesTime: false, // 标记是否为时序图表
    isSeries: false, // 标记是否为连续图表

    max: 0, // 用于连续图表设置宽度起点
    min: 0, // 用于连续图表设置宽度终点
    hoverPointData: null, // function, 用于处理鼠标hover后显示的数据
    loading: false,
    reload: null,
    title: "", // 图表名称
    tooltipLabels: null, // tooltip支持标签显示
    yAxis: [], // y轴字段显示
    unit: "",

    // 图表数据
    labels: [],
    series: [], // 数据结构：Array<DataType>
    colors: null
  };

  protected _model: Model | OverlayModel | SeriesModel;
  protected _titleElement: HTMLDivElement;
  protected _reloadElement: HTMLDivElement;

  protected _drawHoverLine = false; // 在 onMouseMove 中鼠标hover时，标记是否有折线被hover，mainPanel 需要被重绘
  protected _showColorPoint = true; // 如果折线只用单色显示，在不需要在tooltip中显示线颜色
  protected _parentContainer: any;
  protected _container: HTMLDivElement;
  protected _mainPanel: Paint; // 显示主要的静态图表
  protected _auxiliaryPanel: Paint; // 辅助层，显示动态图画，例如显示鼠标悬浮时的标准线
  protected _tooltip: Tooltip; // 显示辅助线上点的信息
  protected _event: Event;

  protected _chartPosition: {
    right: 0;
    left: 0;
    top: 0;
    bottom: 0;
  };

  protected _chartSize: {
    width: 0;
    height: 0;
  };

  // 根据不同chart类型进行绘画, 子类需要重写父类方法
  abstract draw(): void;

  get DefaultColor() {
    return COLORS.Types.Default;
  }

  protected constructor(id: string, options) {
    this._parentContainer = document.getElementById(id);
    if (!this._parentContainer) {
      return;
    }

    this.setOptions(options);

    // 新建一个div，用于设置两个 canvas 图层
    this._container = document.createElement("div");
    this._container.style.cssText = `width: ${this._options.width}px; 
      height: ${this._options.height}px;
      position: relative; 
      text-align: left;
      background-color: #fff;
      cursor: auto;`;

    // 新建两个 canvas 图层
    this._mainPanel = new Paint(
      "",
      this._options.width,
      this._options.height,
      "position: absolute; user-select: none; z-index: 1"
    );
    // 统一设置字体居中
    this._mainPanel.context2d.textAlign = "center";
    this._container.appendChild(this._mainPanel.canvas);
    this._auxiliaryPanel = new Paint(
      "",
      this._options.width,
      this._options.height,
      "position: absolute; -webkit-tap-highlight-color: transparent; user-select: none; cursor: default; z-index: 2"
    );
    this._auxiliaryPanel.context2d.textAlign = "center";
    this._container.appendChild(this._auxiliaryPanel.canvas);

    // 显示选择点信息
    this._tooltip = new Tooltip(
      () => {
        this.clearAuxiliaryInfo();
      },
      params => {
        this.highlightLine(params.legend);
      }
    );
    this._container.appendChild(this._tooltip.entry);

    // 将新建的 div 填充到父容器中
    this._parentContainer.innerHTML = ""; // 清空容器 div 中的元素
    this._parentContainer.appendChild(this._container);

    // 注册事件对象，只有辅助层才有事件响应，所以只需要注册 _auxiliaryPanel
    this._event = new Event(this._auxiliaryPanel);
    this.addListeners();
  }

  // 添加 canvas 事件
  addListeners() {
    this._container.addEventListener("mousemove", this.onMouseMove.bind(this));
    this._container.addEventListener("mouseout", this.onMouseOut.bind(this));
    this._container.addEventListener("click", this.onMouseDown.bind(this));
  }

  isPointInActualArea(x: number, y: number): boolean {
    return (
      x < this._chartPosition.right &&
      x > this._chartPosition.left &&
      y > this._chartPosition.top &&
      y < this._chartPosition.bottom
    );
  }

  /**
   * 高亮 legend 名称的折线
   * @param {string} legend
   */
  highlightLine(legend: string) {
    this._options.series.forEach(line => {
      line.hover = line.legend === utils.FormatStringNoHTMLSharp(legend);
    });
    // 重绘折线
    this.drawHighlightLine();
  }

  /**
   * 对数值进行格式化显示
   * @param {number} value
   * @returns {string}
   */
  formatValue(value: number): string {
    if (!value || isNaN(value)) {
      return `${value}`;
    }
    value = utils.MATH.Strip(value, 6);
    return `${value}`;
  }

  // 设置图标配置项
  setOptions(options) {
    this._options = { ...this._options, ...options };
    // 用户设置了Y轴的值，就设置了刻度数
    this._options.gridScale =
      this._options.yAxis.length > 0
        ? this._options.yAxis.length
        : this._options.gridScale;
    // 显示视图的范围点
    this._chartPosition = {
      left: this._options.paddingHorizontal as any,
      right: (this._options.width - this._options.paddingHorizontal) as any,
      top: (this._options.paddingVertical * 2) as any,
      bottom: (this._options.height - this._options.paddingVertical) as any
    };
    // 视图大小
    this._chartSize = {
      width: (this._chartPosition.right - this._chartPosition.left) as any,
      height: (this._chartPosition.bottom - this._chartPosition.top) as any
    };
    // 当有数据变化是重新计算 canvas 各参数
    if (options.hasOwnProperty("series")) {
      this._options.series = JSON.parse(JSON.stringify(options.series));
      const colors =
        this._options.colors ||
        COLORS.Values[this._options.colorTheme] ||
        COLORS.Values[COLORS.Types.Default];
      // tooltip 显示折线点颜色
      this._showColorPoint = !this._options.overlay && colors.length !== 1;
      /**
       * 初始化每个line的显示属性
       */
      this._options.series.forEach((line, index) => {
        line.legend = FormatStringNoHTMLSharp(line.legend);
        line.hover = false; // hover 则高亮曲线
        line.show = true;
        line.color = line.disable
          ? this._options.disabledColor
          : colors[index % colors.length];
      });
      // 根据 line.disable 进行排序，disable先画线，显示的曲线不会被disable曲线覆盖
      this._options.series.sort((x, y) =>
        x.disable === y.disable ? 0 : x.disable ? -1 : 1
      );
    }
    const modelOptions = {
      chartPosition: this._chartPosition,
      chartSize: this._chartSize,
      labelAlignCenter: this._options.labelAlignCenter,
      yGridNum: this._options.gridScale,
      isKilobyteFormat: this._options.isKilobyteFormat,
      sequence: JSON.parse(JSON.stringify(this._options.yAxis)) // 非空数组的话用于计算数据点的位置, 复制数组防止污染外部数据
    };
    this.initModel(modelOptions);
    this.clearAuxiliaryInfo();
  }

  // 重新设置图表的大小
  public setChartSize(width: number, height: number) {
    this.setOptions({ width, height });
    if (this._container) {
      this._container.style.width = `${width}px`;
      this._container.style.height = `${height}px`;
    }
    this._mainPanel && this._mainPanel.setSize(width, height);
    this._auxiliaryPanel && this._auxiliaryPanel.setSize(width, height);
    this.draw();
  }

  initModel(modelOptions) {
    // 状态模式
    if (this._options.isSeries) {
      this._options.max =
        this._options.max === 0
          ? this._options.labels[this._options.labels.length - 1]
          : this._options.max;
      this._options.min =
        this._options.min === 0 ? this._options.labels[0] : this._options.min;
      this._model = new SeriesModel(
        {
          ...modelOptions,
          pieceOfLabel: { max: this._options.max, min: this._options.min },
          overlay: this._options.overlay
        },
        this._options.labels,
        this._options.series
      );
    } else if (this._options.overlay) {
      this._model = new OverlayModel(
        modelOptions,
        this._options.labels,
        this._options.series
      );
    } else {
      this._model = new Model(
        modelOptions,
        this._options.labels,
        this._options.series
      );
    }
  }

  /**
   * 当要显示的label数大于默认设置的label数时可以进行简化，
   * @returns {boolean}
   */
  hasSimplyLabel(): boolean {
    return this._options.labels.length > this._options.labelScale;
  }

  /**
   * 鼠标悬浮事件，辅助 panel 不可用
   */
  isEventUnable(): boolean {
    return this._options.loading || this._options.series.length === 0;
  }

  /**
   * 清除鼠标悬浮时显示的辅助信息
   */
  clearAuxiliaryInfo() {
    this._tooltip && this._tooltip.close();
    this._auxiliaryPanel &&
      this._auxiliaryPanel.clearRect(
        0,
        0,
        this._options.width,
        this._options.height
      );
  }

  drawLoading() {
    if (this._options.loading) {
      this._auxiliaryPanel.drawSpinner();
    } else {
      this._auxiliaryPanel.removeSpinner();
    }
  }

  drawTitle() {
    if (!this._options.title) {
      return;
    }

    if (!this._titleElement) {
      this._titleElement = GRAPH.CreateDivElement(
        this._options.paddingVertical,
        this._chartPosition.left,
        this._options.axisColor
      );
      this._container.appendChild(this._titleElement);
    }

    const title = this._options.title || "";
    const unit = this._options.unit ? `（${this._options.unit}）` : "";
    this._titleElement.innerHTML = `
        ${title}
        <span style="color: ${this._options.fontColor}">${unit}</span>
    `;
  }

  /**
   * 绘画坐标轴 x，y
   */
  drawAxis() {
    // draw xAxis
    this._mainPanel.drawLine(
      this._chartPosition.left,
      this._chartPosition.bottom,
      this._chartPosition.right,
      this._chartPosition.bottom,
      this._options.axisColor
    );

    // showYAxis 控制是否显示Y轴
    if (
      this._options.showYAxis &&
      this._model.MaxValue &&
      this._model.MaxValue !== 0
    ) {
      // draw yAxis
      this._mainPanel.drawLine(
        this._chartPosition.left,
        this._chartPosition.top,
        this._chartPosition.left,
        this._chartPosition.bottom,
        this._options.axisColor
      );
    }
  }

  drawLegends() {
    if (!this._options.showLegend || this._options.overlay) {
      // 叠加图状态不显示legend
      this._event.removeAllMouseDownEvent();
      return;
    }
    const markRectSize = { width: 10, height: 8 };
    // 默认每个legend的间隔
    const legendGap = 10;
    let legendX = this._options.paddingHorizontal;
    let legendY = this._options.paddingVertical * 1.5;
    this._options.series.forEach((line, index) => {
      const legendWidth =
        this._mainPanel.measureText(line.legend, this._options.font) +
        markRectSize.width;
      // 判断是否需要换行
      if (legendX + legendWidth > this._chartPosition.right) {
        legendX = this._options.paddingHorizontal;
        legendY += markRectSize.height * 1.5;
      }
      // 绘画色块
      const upperLeftCornerX = legendX;
      const upperLeftCornerY = legendY - markRectSize.height;
      this._mainPanel.drawRect(
        upperLeftCornerX,
        upperLeftCornerY,
        markRectSize.width,
        markRectSize.height,
        line.show ? line.color : this._options.disabledColor
      );
      // 绘画直线名称
      this._mainPanel.drawText(
        line.legend,
        legendX + markRectSize.width + 2,
        legendY,
        line.show ? this._options.fontColor : this._options.disabledColor,
        this._options.font,
        "left"
      );

      //注册legend点击事件
      const eventArea = {
        top: upperLeftCornerY,
        bottom: legendY,
        left: legendX,
        right: legendX + legendWidth
      };
      // 之前存在事件，则移除
      this._event.removeMouseDownEvent(eventArea);
      // 注册事件
      this._event.registerMouseDownEvent(eventArea, () => {
        line.show = !line.show;
        if (this._options.overlay) {
          // 叠加态需要重新计算点的位置
          this._model.setOptions(
            {},
            this._options.labels,
            this._options.series
          );
        }
        this.draw();
      });

      legendX += legendWidth + legendGap;
    });
  }

  drawLabels() {
    if (this._model.XAxisTickMark.length === 0) {
      return;
    }
    const hasSimplyLabel = this.hasSimplyLabel();
    // labelGapNum 标记隔多少个label显示一次
    const labelGapNum = hasSimplyLabel
      ? Math.ceil(this._options.labels.length / this._options.labelScale)
      : 1;
    const tickMarkGap = this._options.labelAlignCenter
      ? this._model.XAxisTickMarkGap / 2
      : this._model.XAxisTickMarkGap;
    const textY =
      this._chartPosition.bottom + this._options.paddingVertical / 2;
    // 获取需要显示的label数组
    let visibleLabels = this._options.labels.filter(
      (label, index) => index % labelGapNum === 0
    );

    if (this._options.isSeriesTime) {
      // 处理时序化的标签，对隔天的标签要做日期化处理
      visibleLabels = utils.TIME.FormatSeriesTime(visibleLabels);
    }

    // 使用间隔数据 labelGapNum 进行for循环
    for (let i = 0; i < visibleLabels.length; i++) {
      const label = visibleLabels[i];
      const xAxisLabelScale = this._model.XAxisTickMark[i * labelGapNum];
      // 简化label状态下不需要对 tickMark 进行位移
      const xAxisTickMark = hasSimplyLabel
        ? xAxisLabelScale
        : xAxisLabelScale + tickMarkGap;
      // 刻度线
      this._mainPanel.drawLine(
        xAxisTickMark,
        this._chartPosition.bottom,
        xAxisTickMark,
        this._chartPosition.bottom + 5,
        this._options.axisColor
      );
      this._mainPanel.drawText(
        label,
        xAxisLabelScale,
        textY,
        this._options.fontColor,
        this._options.font
      );
    }
  }

  drawSeriesLabels() {
    if (this._options.min === 0 && this._options.max === 0) {
      return;
    }
    // 需要满足现实标签的空间
    const labelNum = Math.floor(this._options.width / CHART.LabelGap);
    let visibleLabels = utils.TIME.GenerateSeriesTimeVisibleLabels(
      this._options.min,
      this._options.max,
      labelNum
    );
    const distance = this._options.max - this._options.min;
    const textY =
      this._chartPosition.bottom + this._options.paddingVertical / 2;

    // 使用间隔数据 labelGapNum 进行for循环
    for (let i = 0; i < visibleLabels.length; i++) {
      const label = visibleLabels[i];
      const xAxisTickMark =
        this._chartPosition.left +
        (this._chartSize.width * (label - this._options.min)) / distance;
      // 刻度线
      this._mainPanel.drawLine(
        xAxisTickMark,
        this._chartPosition.bottom,
        xAxisTickMark,
        this._chartPosition.bottom + 5,
        this._options.axisColor
      );
      // 刻度值
      this._mainPanel.drawText(
        utils.TIME.FormatTime(label),
        xAxisTickMark,
        textY,
        this._options.fontColor,
        this._options.font
      );
    }
  }

  /**
   * 绘画等分线，且标记值
   * @param {number} scaleValue 单位刻度值
   */
  drawGrid() {
    if (this._model.ScaleValue === 0) {
      return;
    }
    const gridGap = Math.ceil(
      this._chartSize.height / (this._model.Sequence.length - 1)
    );
    for (let i = 0; i < this._model.Sequence.length; i++) {
      // 从上到下画等高线
      const gridLineY = this._chartPosition.top + i * gridGap;

      // 在图表区域内画线
      if (gridLineY < this._chartPosition.bottom) {
        this._mainPanel.drawLine(
          this._chartPosition.left,
          gridLineY,
          this._chartPosition.right,
          gridLineY,
          this._options.gridColor
        );
      }
      if (gridLineY + this._chartSize.width <= this._chartPosition.bottom) {
        // 根据奇偶绘画不同色块
        this._mainPanel.drawRect(
          this._chartPosition.left,
          gridLineY,
          this._chartSize.width,
          gridGap,
          i % 2 === 0 ? "#F9F9F9" : "#FFF"
        );
      }

      // 写刻度值,跟tooltip一直，Y轴和tooltip可以用凑通过tooltipLabels方法进行显示
      let gridScaleValueText = this._options.tooltipLabels
        ? this._options.tooltipLabels(this._model.Sequence[i], "yAxis")
        : this.formatValue(this._model.Sequence[i]);
      // text位置超过 this._chartPosition.left - textAndYAxisGap，则在canvas画布0位置上又对齐，
      const textAndYAxisGap = 10;
      const textWidth = this._mainPanel.measureText(
        gridScaleValueText,
        this._options.font
      );
      const textX = this._chartPosition.left - textAndYAxisGap - textWidth;
      this._mainPanel.drawText(
        gridScaleValueText,
        textX > 0 ? textX : 0,
        gridLineY,
        this._options.fontColor,
        this._options.font,
        "left"
      );
    }
  }

  /**
   * 数据为空时提醒无数据
   */
  drawEmptyData() {
    this._reloadElement && (this._reloadElement.style.display = "none");
    if (this._options.loading) {
      return;
    }
    if (!this._reloadElement) {
      this._reloadElement = GRAPH.CreateDivElement(
        this._options.height / 2,
        this._options.width / 2,
        this._options.fontColor
      );
      this._reloadElement.style.display = "none";
      this._reloadElement.style.transform = "translate(-50%, -50%)";
      const noDataText = document.createElement("span");
      noDataText.innerText = "暂无数据, ";
      this._reloadElement.appendChild(noDataText);
      const reloadText = document.createElement("span");
      reloadText.innerText = "重新加载";
      reloadText.style.color = this._options.activeColor;
      reloadText.onclick = () => {
        this._options.reload();
      };
      this._reloadElement.appendChild(reloadText);
      this._container.appendChild(this._reloadElement);
    }
    if (
      this._options.series.length === 0 &&
      this._options.labels.length === 0
    ) {
      if (this._options.labels.length === 0) {
        this._reloadElement && (this._reloadElement.style.display = "block");
      } else {
        this._mainPanel.drawText(
          "无数据",
          this._options.width / 2,
          this._options.height / 2,
          this._options.fontColor,
          "13px Verdana, Helvetica, Arial, sans-serif"
        );
      }
    }
  }

  /**
   * 高亮状态时需要重绘折线
   */
  drawHighlightLine() {
    if (this._options.activeHover) {
      this.draw();
    }
  }

  /**
   * 当鼠标悬浮于图表时显示辅助信息
   * @param {Array<{data: object; color: string}>} showPoints 显示点数据
   * @param {{x: number; y: number}} mousePosition                 辅助线位置
   */
  drawAuxiliaryInfo(
    showPoints: Array<{ point: object; color: string }>,
    mousePosition: { x: number; y: number }
  ) {
    // 画纵轴辅助线
    this._auxiliaryPanel.drawDashLine(
      mousePosition.x,
      this._chartPosition.top,
      mousePosition.x,
      this._chartPosition.bottom,
      this._options.auxiliaryLineColor
    );

    // 绘画数据点
    showPoints.forEach(item => {
      this._auxiliaryPanel.drawPoints(item.point, item.color);
    });
  }

  /**
   * 显示 tooltip 信息
   * @param {boolean} fixedTooltip       是否可以设置可固定tooltip
   * @param {{x: number; y: number}} position
   * @param {{lebel: string; tooltipContent: Array<TooltipContentType>}} toolTipInfo tooltip信息
   */
  drawTooltip(
    fixedTooltip: boolean,
    position: { x: number; y: number },
    toolTipInfo: { label: string; content: Array<TooltipContentType> }
  ) {
    // 绘画 tooltip
    this._tooltip.setFixed(fixedTooltip);
    this._tooltip.setInformation(
      toolTipInfo.label,
      toolTipInfo.content,
      this._showColorPoint
    );
    // 计算能显示tooltip的x坐标, y坐标
    const toolTipX =
      position.x < this._options.width / 2
        ? position.x + 10
        : position.x - 10 - this._tooltip.width;
    const toolTipY =
      position.y < this._options.height / 2
        ? position.y + 10
        : position.y - 10 - this._tooltip.height;
    this._tooltip.show(toolTipX, toolTipY);
  }

  /**
   * 折线图画折线
   */
  drawLine() {
    // 画点线操作平频繁，使用requestAnimationFrame提升性能
    this._options.series.forEach((line, index) => {
      if (line.show) {
        const lineWidth = line.hover ? LINE.Width.active : LINE.Width.normal;
        this._mainPanel.drawPolyLine(
          this._model.XAxisTickMark,
          line.yPos,
          line.color,
          lineWidth
        );
      }
    });
  }

  /**
   * 面积图
   */
  drawArea() {
    this._options.series.forEach((line, index) => {
      const xPos = this._model.XAxisTickMark;
      // points被取消，需要对 drawArea 进行重写
      if (line.show) {
        const previousLinePoints =
          index > 0
            ? this._options.series[index - 1].points
            : {
                [xPos[0]]: this._chartPosition.bottom - 0.5,
                [xPos[xPos.length - 1]]: this._chartPosition.bottom - 0.5
              };
        this._mainPanel.drawArea(
          previousLinePoints,
          line.points,
          line.color,
          0.3
        );
        this._mainPanel.drawPolyLine(xPos, line.yPos, line.color);
      }
    });
  }

  /**
   * 面积图
   */
  drawAreaOnlyShowTopLine() {
    const xPos = this._model.XAxisTickMark;
    const color = this.DefaultColor;
    for (let i = this._options.series.length - 1; i >= 0; i--) {
      const line = this._options.series[i];
      if (line.show) {
        this._mainPanel.drawPolyLine(xPos, line.yPos, color);
        break;
      }
    }
  }

  /**
   * 条形图
   */
  drawBar() {
    const barWidth = this._model.XAxisTickMarkGap * 0.8;
    const marginLeft = barWidth / 2;
    let previousHeight = this._options.labels.map(i => 0);
    // 绘画数据折线
    this._options.series.forEach(line => {
      if (line.show) {
        this._model.XAxisTickMark.forEach((x, index) => {
          const y = line.yPos[index];
          const barHeight =
            this._chartPosition.bottom - y - (previousHeight[index] || 0);
          this._mainPanel.drawRect(
            x - marginLeft,
            y,
            barWidth,
            barHeight,
            line.color,
            0.8
          );
          previousHeight[index] = previousHeight[index]
            ? previousHeight[index] + barHeight
            : barHeight;
        });
      }
    });
  }

  /**
   * 悬浮鼠标事件处理
   * @param e
   */
  onMouseMove(e, fixedTooltip: boolean = false): void {
    if (this.isEventUnable()) {
      return;
    }
    if (this._tooltip.fixed && !fixedTooltip) {
      // 当进入fixed tooltip状态时，onMouseMove 显示 tooltip 效果取消
      return;
    }

    const mouseX = (e as any).clientX - this._mainPanel.bounds.left;
    const mouseY = (e as any).clientY - this._mainPanel.bounds.top;
    this.clearAuxiliaryInfo();
    // 保证鼠标指针在有效区域内
    if (this.isPointInActualArea(mouseX, mouseY)) {
      let xAxisTickMarkIndex = 0;
      if (this._model.XAxisTickMarkGap === 0) {
        xAxisTickMarkIndex = this._model.XAxisTickMark.findIndex(
          x => x >= mouseX
        );
        if (xAxisTickMarkIndex === -1) {
          return;
        }
      } else {
        xAxisTickMarkIndex = Math.round(
          (mouseX - this._chartPosition.left) / this._model.XAxisTickMarkGap
        );
      }

      // 处理边界情况，当鼠标接近离开x轴宽度区域时，四舍五入会多1
      if (xAxisTickMarkIndex === this._options.labels.length) {
        xAxisTickMarkIndex -= 1;
      }

      const mousePosition = {
        x: this._model.XAxisTickMark[xAxisTickMarkIndex],
        y: mouseY
      };

      let showPoints = []; // 显示辅助线与折线交汇的点
      let content = []; // tooltip 要显示的数据内容
      let hasRedrawLine = false; // 标记是否调用 drawHighlightLine
      const previousPointX =
        this._model.XAxisTickMark[xAxisTickMarkIndex - 1] || 0;
      // 绘画数据点
      this._options.series.forEach((line, index) => {
        if (line.show) {
          const pointY = line.yPos[xAxisTickMarkIndex];

          if (pointY) {
            const label = this._options.labels[xAxisTickMarkIndex];
            if (this._options.overlay) {
              // 显示点的位置信息, 叠加态只显示最大的点（一个）
              if (index === this._options.series.length - 1) {
                showPoints.push({
                  point: { [mousePosition.x]: pointY },
                  color: this.DefaultColor
                });
              }
            } else {
              // 显示点的位置信息
              showPoints.push({
                point: { [mousePosition.x]: pointY },
                color: line.color
              });
            }

            // 不在叠加态、tooltip 未固定状态 且折线显示的情况下，才有悬浮高亮
            if (this._options.activeHover && !fixedTooltip && !line.disable) {
              // 判断鼠标是否 hover 在折线上，是则高亮折线
              let hasHover = false;
              const previousPointY = line.yPos[xAxisTickMarkIndex];
              if (previousPointY) {
                // 计算直线斜率 k=(y2-y1)/(x2-x1)
                const k =
                  (pointY - previousPointY) /
                  (mousePosition.x - previousPointX);
                const b = pointY - mousePosition.x * k;
                hasHover =
                  Math.abs(k * mouseX + b - mouseY) <
                  this._options.hoverPrecision;
              }
              line.hover = hasHover;
              if (hasHover) {
                hasRedrawLine = true; // 标记需要高亮重绘折线
              }
            } else {
              line.hover = false;
            }

            // 置灰折线不需要在tooltip中显示信息
            if (!line.disable) {
              // 判断是否有自定义显示 tooltipLabels
              const value = line.data[label];
              const valueStr = this._options.tooltipLabels
                ? this._options.tooltipLabels(value) || ""
                : this.formatValue(value);
              content.push({
                legend: line.legend,
                color: line.color,
                hover: line.hover,
                value: value,
                label: `${valueStr}${this._options.unit || ""}`
              } as TooltipContentType);
            }
          }
        }
      });
      // 根据 value 对 tooltip 显示数据做排序显示
      content.sort((a, b) => b.value - a.value);
      // 判断title 是否是时间序列类型，是则做时间格式化
      const title = this._options.isSeriesTime
        ? utils.TIME.Format(
            this._options.labels[xAxisTickMarkIndex],
            utils.TIME.DateFormat.fullDateTime
          )
        : this._options.labels[xAxisTickMarkIndex];
      const showInfo = {
        label: title,
        content
      };

      // 外部入口，处理显示的数据
      this._options.hoverPointData &&
        this._options.hoverPointData({
          xAxisTickMarkIndex,
          mousePosition,
          content
        });

      // 显示辅助面板信息
      this.drawAuxiliaryInfo(showPoints, mousePosition);
      if (this._options.showTooltip) {
        // 绘画 tooltip
        this.drawTooltip(fixedTooltip, mousePosition, showInfo);
      }
      // 判断是否与上次显示状态一致，不一致则需要重绘折线
      if (hasRedrawLine !== this._drawHoverLine) {
        this.drawHighlightLine();
        this._drawHoverLine = hasRedrawLine;
      }
    }
  }

  /**
   * 鼠标点击事件处理
   * @param e
   */
  onMouseDown(e): void {
    e.stopPropagation();
    const mouseX = (e as any).clientX - this._mainPanel.bounds.left;
    const mouseY = (e as any).clientY - this._mainPanel.bounds.top;
    if (this._tooltip.isPointInTooltipArea(mouseX, mouseY)) {
      return;
    }
    this.onMouseMove(e, true);
  }

  onMouseOut(e): void {
    e.stopPropagation();
    const mouseX = (e as any).clientX - this._mainPanel.bounds.left;
    const mouseY = (e as any).clientY - this._mainPanel.bounds.top;
    if (!this.isPointInActualArea(mouseX, mouseY)) {
      if (this._options.showTooltip && !this._tooltip.fixed) {
        // 显示tooltip情况下且tooltip不是固定状态，mouse out 才需 clear 面板
        this.clearAuxiliaryInfo();
      }
      // 外部入口，处理显示的数据
      this._options.hoverPointData &&
        this._options.hoverPointData({
          xAxisTickMarkIndex: 0,
          mousePosition: {},
          content: []
        });
    }
  }
}
