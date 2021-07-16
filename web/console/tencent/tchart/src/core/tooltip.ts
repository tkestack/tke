import {HoverStatus} from './event';


// 只针对图表使用
export interface TooltipContentType {
  legend: string,
  color: string,
  hover: boolean,
  value: any,
  label: any,
}

export default class Tooltip {
  private MaxNumOfData = 8; // 同屏下最多显示数据条目数
  protected _self: HTMLDivElement;
  protected _css: HTMLStyleElement;
  private _closeBtn: any;
  private _tooltipDataTbody: any;
  private _fixed: boolean = false;

  protected _closeEvent: any;
  protected _hoverEvent: any;

  private _showColorPoint: boolean = false;
  private _x: number;
  private _y: number;

  constructor(closeEvent, hoverEvent) {
    this._closeEvent = closeEvent;
    this._hoverEvent = hoverEvent;
    this._self = document.createElement("div");
    this._self.style.cssText = 'position: absolute; ' +
      'background-color: white; ' +
      'opacity: 0.9; ' +
      'color: #000000; ' +
      'padding: 0.3rem 0.5rem; ' +
      'border-radius: 3px; ' +
      'z-index: 10; ' +
      'visibility: hidden; ' +
      'font-size: 12px; ' +
      'box-shadow: rgba(33, 33, 33, 0.2) 0px 1px 2px; ' +
      'line-height: 1.2em; ' +
      'top: 0px;';
    this._css = document.createElement("style");
    this._css.innerHTML = '.TChart_close-btn {' +
      '  position: absolute;' +
      '  right: 10px;' +
      '  top: 5px;' +
      '  display: inline-block;' +
      '  width: 17px;' +
      '  height: 17px;' +
      '  opacity: 0.3;' +
      '}' +
      '.TChart_close-btn:hover {' +
      '  opacity: 1;' +
      '}' +
      '.TChart_close-btn:before, .TChart_close-btn:after {' +
      '  position: absolute;' +
      '  left: 15px;' +
      '  content: \' \';' +
      '  height: 15px;' +
      '  width: 2px;' +
      '  background-color: #333;' +
      '}' +
      '.TChart_close-btn:before {' +
      '  transform: rotate(45deg);' +
      '}' +
      '.TChart_close-btn:after {' +
      '  transform: rotate(-45deg);' +
      '}' +
      '.TChart_tooltip-data:hover {' +
      '  color: #006eff !important;' +
      '}' +
      '.TChart_tooltip-data.hover {' +
      '  color: #006eff !important;' +
      '}';
    document.getElementsByTagName("head")[0].appendChild(this._css);
  }

  get entry(): HTMLDivElement {
    return this._self;
  }

  get width(): number {
    return this._self.offsetWidth;
  }

  get height(): number {
    return this._self.offsetHeight;
  }

  get X(): number {
    return this._x;
  }

  get Y(): number {
    return this._y;
  }

  get fixed() {
    return this._fixed;
  }

  isPointInTooltipArea(x: number, y: number): boolean {
    return (
      x < this.X + this.width &&
      x > this.X &&
      y > this.Y &&
      y < this.Y + this.height
    );
  }

  setInformation(title: string, content: Array<TooltipContentType>, showColorPoint: boolean): void {
    this._self.innerHTML = `
                  <table style="${content.length < this.MaxNumOfData ? '' : 'height:130px;'} display: flex; flex-direction: column;">
                    <thead style="margin-bottom: 4px; display: block; position: relative; ">
                      <tr>
                        <td style="display: inline-block; font-size: 12px; height: 20px; line-height: 20px">${title}</td>
                        <td style="display: inline-block; margin-left: 10px; width: 20px;height: 20px"></td>
                      </tr>
                     </thead>
                   </table>`;

    if (this._fixed) {
      this._closeBtn = document.createElement("i");
      this._closeBtn.setAttribute("class", "TChart_close-btn");
      this._self.getElementsByTagName("td")[1].appendChild(this._closeBtn);
      this._closeBtn && (this._closeBtn.onclick = (e) => {
        e.stopPropagation();
        this.close();
        this._closeEvent && this._closeEvent();
      });
    }

    this._tooltipDataTbody = document.createElement("tbody");
    this._tooltipDataTbody.style = "flex: 1;overflow: auto;";
    this._showColorPoint = showColorPoint;
    this._tooltipDataTbody.innerHTML = content.map(item => {
      return `<tr style="text-align: right;" class="TChart_tooltip-data ${item.hover? 'hover': ''}" data="${item.legend}">
                ${
                  this._showColorPoint
                    ? `<td><span style="display: block; margin-right: 4px; border-radius: 50%; width: 8px; height: 8px; background: ${item.color};"></span></td>`
                    : ""
                }
                <td style="display: block;text-align: left; min-width: 100px; margin-right: 8px;">${item.legend} </td>
                <td style="text-align: right;">${item.label}</td> 
              </tr>`;
    }).join("");;
    this._self.firstElementChild.appendChild(this._tooltipDataTbody);
    // 注册 tooltip 数据条 hover 事件
    if (this._tooltipDataTbody) {
      this._tooltipDataTbody.onmouseover = (e) => {
        e.stopPropagation();
        let legend = "";
        if (e.target.parentElement.hasAttribute("data")) {
          legend = e.target.parentElement.getAttribute("data");
        }
        this._hoverEvent && this._hoverEvent({legend});
      };

      this._tooltipDataTbody.onmouseout = (e) => {
        e.stopPropagation();
        this._hoverEvent && this._hoverEvent({legend: ""});
      };
    }
  }

  /**
   * 关闭tooltip，取消固定状态
   */
  close() {
    this.setFixed(false);
    this.hidden();
  }

  show(x: number, y: number): void {
    this._x = x;
    this._y = y;
    this._self.style.visibility = "visible";
    this._self.style.transform = `translate(${x}px, ${y}px)`;
  }

  hidden(): void {
    this._self.style.visibility = "hidden";
  }

  setFixed(fixed) {
    this._fixed = fixed;
  }

}