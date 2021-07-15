import moment from "moment";

const DAY = 24 * 60 * 60;
const HOUR = 60 * 60;
const MINUTE = 60;

export let minPeriod = 60;
declare var process : {
  env: {
    NODE_ENV: string
  }
}

export let durationsByPeriod = {
  // [period(s)] : day
  10: 1,
  60: 15,
  300: 31,
  3600: 62,
  86400: 186
};
/**
 * 根据时间间隔计算查询时间粒度(秒)
 * @param {Date} from
 * @param {Date} to
 * @returns {string}
 */
export function Period(from: Date, to: Date): number {
  const range = moment(to).diff(moment(from));
  let startTime = moment(from);
  let startTimeDiffFromNow = -startTime.diff(undefined);
  if (startTimeDiffFromNow < durationsByPeriod[10] * DAY * 1000) {
    if (range <= 8 * HOUR * 1000) return 1 * MINUTE; // 最小
  }

  if (startTimeDiffFromNow < durationsByPeriod[MINUTE] * DAY * 1000) {
    if (range <= 2 * DAY * 1000) return 1 * MINUTE;
  }

  if (startTimeDiffFromNow < durationsByPeriod[5 * MINUTE] * DAY * 1000) {
    // 因后台存储，临时规则改为5 分钟 288 个点， 即3天
    // if (range <= 10 * DAY * 1000) return 5 * MINUTE;
    if (range <= 3 * DAY * 1000) return 5 * MINUTE;
  }

  if (startTimeDiffFromNow < durationsByPeriod[HOUR] * DAY * 1000) {
    if (range <= 120 * DAY * 1000) return HOUR;
  }

  if (startTimeDiffFromNow < durationsByPeriod[DAY] * DAY * 1000) {
    return DAY;
  }

  return DAY * 31;
}

export function TimeFormat(from: Date, to: Date): string {
  const range = moment(to).diff(moment(from), "hours");
  return range <= 24 ? "HH:mm" : "MM-DD HH:mm";
}

const UNITS = ["", "K", "M", "G", "T", "P"];

/**
 * 进行单位换算
 * @param {number} value
 * @param {number} thousands
 * @param {number} toFixed
 */
export function TransformField(_value: number, thousands, toFixed = 3, units = UNITS) {
  let value = _value;
  let isValueDefined = !isNaN(value) && value !== null;
  if (!isValueDefined) return "-";

  let unitBase = units[0];
  let i = units.indexOf(unitBase);
  if (isValueDefined && thousands) {
    while (i < units.length && value / thousands > 1) {
      value /= thousands;
      ++i;
    }
    unitBase = units[i] || "";
  }
  let svalue;
  if (value > 1) {
    svalue = value.toFixed(toFixed);
    svalue = svalue.replace(/0+$/, "");
    svalue = svalue.replace(/\.$/, "");
  } else if (value) {
    // 如果数值很小，保留toFixed位有效数字
    let tens = 0;
    let v = Math.abs(value);
    while (v < 1) {
      v *= 10;
      ++tens;
    }
    svalue = value.toFixed(tens + toFixed - 1);
    svalue = svalue.replace(/0+$/, "");
    svalue = svalue.replace(/\.$/, "");
  } else {
    svalue = value;
  }
  return String(svalue) + (value !== 0 ? unitBase : "");
}


export namespace STORE {
  export function Set(key, data) {
    localStorage.setItem(key, JSON.stringify(data));
  }

  export function Get(key, defaultValue = null) {
    try {
      return JSON.parse(localStorage.getItem(key)) || defaultValue;
    } catch (e) {
      return defaultValue;
    }
  }
}
