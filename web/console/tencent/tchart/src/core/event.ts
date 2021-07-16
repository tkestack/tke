import Paint from './paint';


export const HoverStatus = {
  on: 'on',
  out: 'out',
};

export interface EventRectType {
  top: number,
  bottom: number,
  left: number,
  right: number,
}

export default class Event {

  protected _paint: Paint;
  protected _mouseDownEventList: Array<{ area: EventRectType, func: Function }> = [];

  constructor(paint: Paint) {
    this._paint = paint;
    this._paint.addEventListener("mousedown", this.onMouseDown.bind(this))
  }

  registerMouseDownEvent(rect: EventRectType, func: Function) {
    this._mouseDownEventList.push({
      area: rect,
      func
    });
  }

  removeAllMouseDownEvent() {
    this._mouseDownEventList = [];
  }

  removeMouseDownEvent(area: EventRectType) {
    this._mouseDownEventList = this._mouseDownEventList.filter(item => {
      if (item.area.top === area.top && item.area.bottom === area.bottom &&
        item.area.left === area.left && item.area.right === area.right) {
        return false;
      }
      return true;
    })
  }

  onMouseDown(e) {
    const mouseX = (e as any).clientX - this._paint.bounds.left;
    const mouseY = (e as any).clientY - this._paint.bounds.top;
    this._mouseDownEventList.forEach(event => {
      const { area, func } = event;
      if (mouseX > area.left && mouseX < area.right && mouseY > area.top && mouseY < area.bottom) {
        func();
      }
    });
  }

}