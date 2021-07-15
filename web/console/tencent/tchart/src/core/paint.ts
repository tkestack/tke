import { CHART } from 'core/theme';
import { GRAPH } from './utils';


export default class Paint {

  protected _canvas: HTMLCanvasElement;
  protected _ctx: CanvasRenderingContext2D;

  protected _lineDashConfig = [5, 5];

  protected _spinnerAnimation;


  constructor(id: string, width: number, height: number, style: string) {
    if (id) {
      this._canvas = <HTMLCanvasElement> document.getElementById(id);
    }
    else {
      this._canvas = <HTMLCanvasElement> document.createElement("canvas");
    }
    this._ctx = this._canvas.getContext('2d');
    this._canvas.style.cssText = style;
    this.setSize(width, height);
  }

  get canvas() {
    return this._canvas;
  }

  get bounds() {
    return this._canvas.getBoundingClientRect();
  }

  get context2d() {
    return this._ctx;
  }

  setSize(width: number, height: number) {
    GRAPH.ScaleCanvas(this._canvas, this._ctx, width, height);
  }

  setTextAlign(textAlign = "center") {
    (this.context2d as any).textAlign = textAlign;
  }

  setLineDash(lineDash: Array<number>) {
    this._lineDashConfig = lineDash;
  }

  addEventListener(event: string, callback: EventListener) {
    this._canvas.addEventListener(event, callback, false);
  }

  clearRect(startX: number, startY: number, endX: number, endY: number) {
    this._ctx.clearRect(startX, startY, endX, endY);
  }

  /**
   * 绘画线
   * @param {number} startX
   * @param {number} startY
   * @param {number} endX
   * @param {number} endY
   * @param {string} color
   */
  drawLine(startX: number, startY: number, endX: number, endY: number, color: string) {
    this._ctx.save();
    this._ctx.strokeStyle = color;
    this._ctx.beginPath();
    this._ctx.moveTo(startX, startY);
    this._ctx.lineTo(endX, endY);
    this._ctx.stroke();
    this._ctx.restore();
  }

  drawDashLine(startX: number, startY: number, endX: number, endY: number, color: string) {
    this._ctx.save();
    this._ctx.strokeStyle = color;
    this._ctx.beginPath();
    this._ctx.setLineDash(this._lineDashConfig);
    this._ctx.moveTo(startX, startY);
    this._ctx.lineTo(endX, endY);
    this._ctx.stroke();
    this._ctx.restore();
  }

  /**
   *
   * @param {Array<number>} xPos 为 labels 的 x 轴坐标数组
   * @param {object} points
   * @param {string} color
   */
  drawPolyLine(xPos: Array<number>, yPos: Array<number>, color: string, lineWidth: number = 1) {
    this._ctx.save();
    // draw ploy line
    this._ctx.strokeStyle = color;
    this._ctx.beginPath();
    let isContinuous = false; // 用来标记前面的点是否存在，存在就画线，不存在则需要移动后判断是否要画点
    this._ctx.lineWidth = lineWidth;

    xPos.forEach((x, index) => {
      const y = yPos[index];
      if (isNaN(y) || y === null) {
        if (isContinuous) {
          // 优化性能，不能在lineTo下进行stroke，会导致画图缓慢
          // 通过判断上传数据存在，只触发一次stroke
          this._ctx.stroke();
        }
        // y 坐标点不存在，则下个点为移动坐标
        isContinuous = false;
        return;
      }
      else {
        if (isContinuous) {
          this._ctx.lineTo(x, y);
        }
        else {
          this._ctx.moveTo(x, y);
          // 画点：如果下一个点不存在，则没有连续性，需要画点
          const yNext = yPos[index + 1];
          if (isNaN(yNext) || yNext === null) {
            this._ctx.fillStyle = color;
            this._ctx.beginPath();
            this._ctx.arc(x, y, 2, 0, Math.PI * 2, true);
            this._ctx.closePath();
            this._ctx.fill();
          }
        }
        isContinuous = true;
      }
    });

    this._ctx.stroke();
    this._ctx.restore();
  }

  /**
   * 保证 previousPoints 与 currentPoints 起始位置和结束位置一致
   * @param {object} previousPoints
   * @param {object} currentPoints
   * @param {string} color
   * @param {number} alpha
   */
  drawArea(previousPoints: object, currentPoints: object, color: string, alpha: number) {
    function sortNumber(a, b) {
      return Number(a) - Number(b);
    }
//    Chrome Opera 的 JavaScript 解析引擎遵循的是新版 ECMA-262 第五版规范。因此，使用 for-in 语句遍历对象属性时遍历书序并非属性构建顺序。
//    而 IE6 IE7 IE8 Firefox Safari 的 JavaScript 解析引擎遵循的是较老的 ECMA-262 第三版规范，属性遍历顺序由属性构建的顺序决定。
//
//    Chrome Opera 中使用 for-in 语句遍历对象属性时会遵循一个规律：
//    它们会先提取所有 key 的 parseFloat 值为非负整数的属性，然后根据数字顺序对属性排序首先遍历出来，然后按照对象定义的顺序遍历余下的所有属性。
    let xPos = Object.keys(currentPoints).sort(sortNumber);

    // previousPoints 的起始或结束位置没有数据，则返回
    if (isNaN(previousPoints[xPos[0]]) && isNaN(previousPoints[xPos[xPos.length - 1]])) {
      return;
    }

    this._ctx.save();
    this._ctx.fillStyle = color;
    this._ctx.globalAlpha = alpha;
    this._ctx.beginPath();
    const startX = Number(xPos[0]);
    // 绘画区域
    this._ctx.moveTo(startX, previousPoints[startX]);
    xPos.forEach(x => {
      this._ctx.lineTo(Number(x), currentPoints[x]);
    });

    xPos = Object.keys(previousPoints).sort(sortNumber);
    for (let i = xPos.length - 1; i >= 0; i--) {
      const x = xPos[i];
      this._ctx.lineTo(Number(x), previousPoints[x]);
    }
    this._ctx.closePath();
    this._ctx.fill();
    this._ctx.restore();
  }

  drawPoints(points: object, color: string) {
    this._ctx.save();
    const endAngle = Math.PI * 2;
    let xPos = Object.keys(points);
    xPos.forEach(x => {
      this._ctx.fillStyle = "#fff";
      this._ctx.beginPath();
      this._ctx.arc(Number(x), points[x], 6, 0, endAngle, true);
      this._ctx.closePath();
      this._ctx.fill();

      this._ctx.fillStyle = color;
      this._ctx.beginPath();
      this._ctx.arc(Number(x), points[x], 4, 0, endAngle, true);
      this._ctx.closePath();
      this._ctx.fill();
    });
    this._ctx.restore();
  }

  /**
   * 绘画柱
   * @param {number} upperLeftCornerX
   * @param {number} upperLeftCornerY
   * @param {number} width
   * @param {number} height
   * @param {string} color
   */
  drawRect(upperLeftCornerX: number, upperLeftCornerY: number, width: number, height: number, color: string, alpha: number = 1.0) {
    this._ctx.save();
    this._ctx.fillStyle = color;
    this._ctx.globalAlpha = alpha;
    this._ctx.fillRect(upperLeftCornerX, upperLeftCornerY, width, height);
    this._ctx.restore();
  }

  /**
   * @param {string} text
   * @param {number} x
   * @param {number} y
   * @param {string} color
   * @param {string} font
   */
  drawText(text: string, x: number, y: number, color: string, font: string, textAlign: string = "center") {
    this.setTextAlign(textAlign);
    this._ctx.save();
    this._ctx.fillStyle = color;
    this._ctx.font = font;
    this._ctx.fillText(text, x, y);
    this._ctx.restore();
    this.setTextAlign();
  }

  /**
   * drawSpinner 动画函数
   */
  drawSpinner() {
    // spinner 参数, 参数也可以写为对象的属性，但是这不符合最小化作用范围
    const _degrees = new Date();
    const _offset = 16;
    const width = Number(this._canvas.style.width.replace('px', ''));
    const height = Number(this._canvas.style.height.replace('px', ''));
    const moveX = width / 120;
    const lineX = width / 60;
    const lineWidth = width / 500;

    function spinnerAnimation() {
      this._spinnerAnimation = window.requestAnimationFrame(spinnerAnimation.bind(this));

      const rotation = parseInt((((new Date() as any - (_degrees as any)) / 1000) * _offset) as any) / _offset;
      this._ctx.save();
      this._ctx.clearRect(0, 0, width, height);
      this._ctx.translate(width / 2, height / 2);
      this._ctx.rotate(Math.PI * 2 * rotation);
      for (let i = 0; i < _offset; i ++) {
        this._ctx.beginPath();
        this._ctx.rotate(Math.PI * 2 / _offset);
        this._ctx.moveTo(moveX, 0);
        this._ctx.lineTo(lineX, 0);
        this._ctx.lineWidth = lineWidth;
        this._ctx.strokeStyle = "rgba(0, 111, 250," + i / _offset + ")";
        this._ctx.stroke();
      }
      this._ctx.restore();
    };
    // 防止有重复的spinner
    this.removeSpinner();
    spinnerAnimation.apply(this);
    // 安全机制，用于在限制时间(6秒)内关闭drawSpinner (云API5秒请求超时)
    setTimeout(() => {
      this.removeSpinner();
    }, CHART.SpinnerTime);
  }

  removeSpinner() {
    this._spinnerAnimation && window.cancelAnimationFrame(this._spinnerAnimation);
    if (this._spinnerAnimation) {
      this._ctx.clearRect(0, 0, this._canvas.width, this._canvas.height);
    }
    this._spinnerAnimation = null;
  }

  /**
   * 返回
   * @param {string} text
   * @returns {number}
   */
  measureText(text: string, font: string) {
    this._ctx.font = font;
    return this._ctx.measureText(text).width;
  }

}