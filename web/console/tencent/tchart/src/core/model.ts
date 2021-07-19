import { MATH } from './utils';


export interface ModelType {
  legend: string;
  show?: boolean;
  yPos?: Array<number>;
  disable?: boolean;
  data: {[x: string]: number};
}

const DefaultMaxValue = 5;

const InitialState = {
  maxValue: DefaultMaxValue,            // y轴最大值
  minValue: 0,                          // y轴最小值
  sequence: [],
  scaleValue: {exponent: 0, value: 0},  // y轴刻度间距值
  xAxisTickMarkGap: 0,
  xAxisTickMark: []
};

/**
 * 图表的数据模型，处理传入的数据
 */
abstract class BaseModel {

  protected state = {...InitialState};

  protected _options = {
    overlay: false,
    pieceOfLabel: { // 用于连续性图表类型, 记录开始和结束点
      min: 0,
      max: 0
    },
    labelAlignCenter: false,
    isKilobyteFormat: false,
    sequence: [],   // 配置的序列
    yGridNum: 0,
    chartPosition: {
      right: 0,
      left: 0,
      top: 0,
      bottom: 0
    },
    chartSize: {
      width: 0,
      height: 0
    }
  };

  /**
   * 子类继承计算最大、最小值
   */
  protected abstract calculateMaxAndMinSourceValue(series: Array<ModelType>): {maximum: number, minimum: number};

  public abstract calculatePoints(labels: Array<any>, series: Array<ModelType>);

  protected constructor(options, label: Array<any>, series: Array<ModelType>) {
    this.setOptions(options, label, series);
  }

  public setOptions(options, labels: Array<any>, series: Array<ModelType>) {
    this._options = {
      ...this._options,
      ...options
    };
    this.calculatePoints(labels, series);
  }

  get ScaleValue() {
    return this.state.scaleValue.value * (10 ** this.state.scaleValue.exponent);
  }

  // ScaleValue 基数
  get ScaleValueCardinal() {
    return this.state.scaleValue.value;
  }

  // ScaleValue 指数
  get ScaleValueExponent() {
    return this.state.scaleValue.exponent;
  }

  get MaxValue() {
    return this.state.maxValue;
  }

  get XAxisTickMarkGap() {
    return this.state.xAxisTickMarkGap;
  }

  get XAxisTickMark() {
    return this.state.xAxisTickMark;
  }

  get Sequence() {
    return this.state.sequence;
  }

  protected resetState() {
    this.state = {...InitialState};
  }

  protected setState(states) {
    this.state = {...this.state, ...states};
  }

  protected arrayMaxAndMinValue(sources: Array<number>): {maximum: number, minimum: number} {
    let maximum = 0, minimum = 0;
    if (sources.length > 0) {
      maximum = Math.max(...sources);
      const min = Math.min(...sources);
      if (min < 0) {
        minimum = min;
      }
    }
    return {maximum, minimum};
  }

  /**
   * 根据数据源计算y轴显示整数最大值
   */
  protected calculateScaleAndMaxValue(maximum: number, minimum: number): void {
    // 处理数值为0的边界情况
    if (maximum === 0) {
      // 如果最大值为0，则使用默认值
      maximum = DefaultMaxValue;
    }

    // 用户自定义刻度大小
    if (this._options.sequence.length > 0 && this._options.sequence[this._options.sequence.length - 1] > maximum) {
      const sequence = this._options.sequence;
      const gridValue = sequence[sequence.length - 1] / sequence.length;
      this.setState({
        scaleValue: {exponent: 0, value: gridValue},
        maxValue: sequence[sequence.length - 1],
        minValue: minimum,
        sequence: sequence.reverse()
      });
      return;
    }

    const {gridValue, sequence} = this._options.isKilobyteFormat
      ? MATH.ArithmeticSequence(minimum, maximum, this._options.yGridNum)
      : MATH.ScaleSequence(minimum, maximum, this._options.yGridNum);
    this.setState({
      scaleValue: {exponent: 0, value: gridValue},
      maxValue: sequence[sequence.length - 1],
      minValue: minimum,
      sequence: sequence.reverse()
    });
  }

  // 计算 label 在 x轴上的位置,数据点的x坐标需要使用
  protected calculateXAxisTickMark(labelsNum: number) {
    const {labelAlignCenter, chartPosition, chartSize} = this._options;
    this.state.xAxisTickMark = [];
    this.state.xAxisTickMarkGap = chartSize.width / labelsNum;
    // 计算标签居中间隔
    const centerMove = labelAlignCenter ? this.state.xAxisTickMarkGap / 2 : 0;
    // 计算标签居中间隔
    for (let i = 0; i < labelsNum; i++) {
      this.state.xAxisTickMark.push(chartPosition.left + i * this.state.xAxisTickMarkGap + centerMove)
    }
  }

}


/**
 * 常规计算
 */
export class Model extends BaseModel {

  constructor(options, labels: Array<any>, series: Array<ModelType>) {
    super(options, labels, series);
  }

  protected calculateMaxAndMinSourceValue(series: Array<ModelType>): {maximum: number, minimum: number} {
    // 返回每个对象数据的数据集，即为二维数组
    const data2dArray = series.map(item => Object.values(item.data));
    // 转换为一维数组 [].concat(...arr2d)，求出所有数据的最大值
    const dataFromAllSources = [].concat(...data2dArray);

    return this.arrayMaxAndMinValue(dataFromAllSources);
  }

  /**
   * 核心对外API
   * @param {Array<DataType>} series
   */
  public calculatePoints(labels: Array<any>, series: Array<ModelType>) {
    const {chartPosition, chartSize} = this._options;
    if (labels.length === 0) {
      this.resetState();
      return;
    }
    // X 轴刻度线坐标数组
    this.calculateXAxisTickMark(labels.length);
    const {maximum, minimum} = this.calculateMaxAndMinSourceValue(series);
    this.calculateScaleAndMaxValue(maximum, minimum);

    const {maxValue, xAxisTickMark} = this.state;

    // 每单位数值在y轴的高度
    const pointYGap = maxValue != 0 ? chartSize.height / maxValue : 0;
    series.forEach((line: ModelType) => {
      // 初始化属性值 yPos 和 show
      line.yPos = [];
      labels.forEach((label, index) => {
        const value = line.data[label];
        if (value !== null && !isNaN(value)) {
          const pointY = chartPosition.bottom - Math.round(pointYGap * value);
          // 在 canvas actual 区域内绘画
          line.yPos.push(pointY);
        }
        else {
          line.yPos.push(null);
        }
      });
    });
  }
}


/**
 * 叠加计算
 */
export class OverlayModel extends BaseModel {

  constructor(options, labels: Array<any>, series: Array<ModelType>) {
    super(options, labels, series);
  }

  protected calculateMaxAndMinSourceValue(series: Array<ModelType>): {maximum: number, minimum: number} {
    // 返回每组数据的数据集，即为二维数组
    const data2dArray = series.map(item => Object.values(item.data));
    const lineNum = data2dArray.length;
    if (lineNum > 0) {
      let combineArray = [];
      data2dArray[0].forEach((value, index) => {
        let sum = 0;
        // 计算各个节点不同线的总和
        for (let j = 0; j < lineNum; j++) {
          sum += data2dArray[j][index] || 0;
        }
        combineArray.push(sum);
      });

      return this.arrayMaxAndMinValue(combineArray);
    }

    return this.arrayMaxAndMinValue([]);
  }

  /**
   * 核心对外API
   * @param {Array<DataType>} series
   */
  public calculatePoints(labels: Array<any>, series: Array<ModelType>) {
    const {chartPosition, chartSize} = this._options;
    if (labels.length === 0) {
      this.resetState();
      return;
    }
    // X 轴刻度线坐标数组
    this.calculateXAxisTickMark(labels.length);
    const {maximum, minimum} = this.calculateMaxAndMinSourceValue(series);
    this.calculateScaleAndMaxValue(maximum, minimum);

    const {maxValue, xAxisTickMark} = this.state;

    let previousHeight = labels.map(i => 0);
    // 每单位数值在y轴的高度
    const pointYGap = maxValue != 0 ? chartSize.height / maxValue : 0;
    series.forEach((line: ModelType) => {
      // 初始化属性值 yPos 和 show
      line.yPos = [];
      labels.forEach((label, index) => {
        const value = line.data[label] || 0;
        const pointY = chartPosition.bottom - Math.round(pointYGap * (value + (previousHeight[index] || 0)));
        // 在 canvas actual 区域内绘画
        line.yPos.push(pointY);
        previousHeight[index] = previousHeight[index] ? previousHeight[index] + value : value;
      });
    });
  }
}


/**
 * 时序图根据x点计算位置
 */
export class SeriesModel extends BaseModel {

  constructor(options, labels: Array<any>, series: Array<ModelType>) {
    super(options, labels, series);
  }

  protected calculateMaxAndMinSourceValue(series: Array<ModelType>): {maximum: number; minimum: number} {
    // 返回每个对象数据的数据集，即为二维数组
    const data2dArray = series.map(item => Object.values(item.data));

    // 返回每组数据的数据集，即为二维数组
    const lineNum = data2dArray.length;
    if (lineNum > 0) {
      let combineArray = [];
      if (this._options.overlay) {
        data2dArray[0].forEach((value, index) => {
          let sum = 0;
          // 计算各个节点不同线的总和
          for (let j = 0; j < lineNum; j++) {
            sum += data2dArray[j][index] || 0;
          }
          combineArray.push(sum);
        });
      }
      else {
        // 转换为一维数组 [].concat(...arr2d)，求出所有数据的最大值
        combineArray = [].concat(...data2dArray);
      }
      return this.arrayMaxAndMinValue(combineArray);
    }

    return this.arrayMaxAndMinValue([]);
  }

  calculatePoints(labels: Array<any>, series: Array<ModelType>) {
    const {chartPosition, pieceOfLabel, chartSize} = this._options;
    if (labels.length === 0) {
      this.resetState();
      return;
    }
    const {maximum, minimum} = this.calculateMaxAndMinSourceValue(series);
    this.calculateScaleAndMaxValue(maximum, minimum);

    const {maxValue} = this.state;
    let xAxisTickMark = [];
    const distance = pieceOfLabel.max - pieceOfLabel.min;
    labels.forEach(label => {
      const pointX = chartPosition.left + (label - pieceOfLabel.min) * chartSize.width / distance;
      xAxisTickMark.push(pointX);
    });
    this.setState({xAxisTickMark});

    let previousHeight = labels.map(i => 0);
    // 每单位数值在y轴的高度
    const pointYGap = maxValue != 0 ? chartSize.height / maxValue : 0;
    if (this._options.overlay) {
      series.forEach((line: ModelType) => {
        // 初始化属性值 yPos 和 show
        line.yPos = [];
        labels.forEach((label, index) => {
          const value = line.data[label];
          previousHeight[index] = value !== null && !isNaN(value) ? previousHeight[index] + value : previousHeight[index];
          const pointY = chartPosition.bottom - Math.round(pointYGap * previousHeight[index]);
          // 在 canvas actual 区域内绘画
          line.yPos.push(pointY);
        });
      });
    }
    else {
      series.forEach((line: ModelType) => {
        // 初始化属性值 yPos 和 show
        line.yPos = [];
        labels.forEach((label, index) => {
          const value = line.data[label];
//          const pointX = xAxisTickMark[index];
          if (value !== null && !isNaN(value)) {
            const pointY = chartPosition.bottom - Math.round(pointYGap * value);
            // 在 canvas actual 区域内绘画
            line.yPos.push(pointY);
          }
          else {
            line.yPos.push(null);
          }
        });
      });
    }
  }
}