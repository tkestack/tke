export namespace MATH {
  /**
   * 解决浮点误差
   * @param {Number} num
   * @param {number} precision
   * @returns {number}
   * @constructor
   */
  export function Strip(num: Number, precision = 12) {
    return +parseFloat(num.toPrecision(precision));
  }

  export function ScientificNotation(value: number, exponent: number) {
    if (exponent > 0) {
      return `${value}E${exponent}`;
    } else if (exponent === 0) {
      return `${value}`;
    } else {
      return `${value}e${exponent}`;
    }
  }

  /**
   * 将数据转换为文件格式, 1024 = 2 ** 10
   * @param {number} size
   * @returns {string}
   */
  export function FormatFileSize(size: number, toFixed: number = 3): string {
    const Kilobyte = 1024;
    const fileSizeUnits = ["", "K", "M", "G", "T", "P", "E", "Z", "Y"];
    let exponent = 0;
    const maxExponent = fileSizeUnits.length - 1;
    while (size >= Kilobyte) {
      if (exponent >= maxExponent) {
        break;
      }
      size = size / Kilobyte;
      exponent += 1;
    }
    return `${Strip(size, toFixed)}${fileSizeUnits[exponent]}`;
  }

  /**
   * 根据最小值，最大值和间隔数，计算对应的10进制转2进制的等差数组
   * @param {number} minimum 最小值
   * @param {number} maximum 最大值
   * @param {number} gridNum 间隔数
   */
  export function ArithmeticSequence(
    minimum: number,
    maximum: number,
    gridNum: number = 5
  ): { gridValue: number; sequence: Array<number> } {
    if (gridNum === 0) {
      return { gridValue: 0, sequence: [] };
    }
    if (minimum === maximum && gridNum > 0) {
      return { gridValue: 0, sequence: [] };
    }

    // if minimum > maximum, reverse them
    if (minimum > maximum) {
      const temp = minimum;
      minimum = maximum;
      maximum = temp;
    }
    const step = Math.max(0, gridNum);

    // 计算数值间隔
    const gridGap = (maximum - minimum) / step;
    const power = Math.floor(Math.log2(gridGap)); // power 是以2为底的指数，2**power <= gridGap
    const ratio = gridGap / Math.pow(2, power);
    // 参考 d3.array.ticks, 原为10进制，改为2进制，我也不理解
    let multiple = 1;
    if (ratio >= Math.sqrt(50)) {
      multiple = 10;
    } else if (ratio >= Math.sqrt(10)) {
      multiple = 5;
    } else if (ratio >= Math.sqrt(2)) {
      multiple = 2;
    }
    // 等差值
    const gridValue = power >= 0 ? multiple * Math.pow(2, power) : -Math.pow(2, -power) / multiple;

    if (gridValue === 0 || !isFinite(gridValue)) {
      return { gridValue: 0, sequence: [] };
    }

    let ticks: Array<number>,
      n,
      i = 0;
    if (gridValue > 0) {
      const start = Math.floor(minimum / gridValue);
      const stop = Math.ceil(maximum / gridValue);
      ticks = new Array<number>((n = Math.ceil(stop - start + 1)));
      while (i < n) {
        ticks[i] = (start + i) * gridValue;
        i += 1;
      }
    } else {
      const start = Math.ceil(minimum * gridValue);
      const stop = Math.floor(maximum * gridValue);
      ticks = new Array<number>((n = Math.ceil(start - stop + 1)));
      while (i < n) {
        ticks[i] = (start - i) / gridValue;
        i += 1;
      }
    }

    return { gridValue: i > 0 ? ticks[1] - ticks[0] : ticks[0], sequence: ticks };
  }

  export function ScaleSequence(minimum: number, maximum: number, gridNum: number = 5) {
    const averageValue = maximum / gridNum;

    // 将最大值转换为 个位数 * 10的指数
    // 指数级
    let exponent = 0;
    let scaleValue = averageValue;
    if (-1 < scaleValue && scaleValue < 1) {
      // 小于零，乘以10计算到个位
      while (scaleValue < 1) {
        scaleValue = scaleValue * 10;
        exponent -= 1;
      }
    } else if (scaleValue > 10) {
      // 大于十，除以10计算到个位
      while (scaleValue > 10) {
        scaleValue = scaleValue / 10;
        exponent += 1;
      }
    }
    // 向上取整
    scaleValue = Math.ceil(scaleValue);

    const gridValue = scaleValue * 10 ** exponent;
    // 还未优化，以后支持负值显示
    let sequence = [];
    for (let i = 0; i <= gridNum; i++) {
      sequence.push(gridValue * i);
    }
    return { gridValue, sequence };
  }
}

export namespace TIME {
  export const DateFormat = {
    fullDateTime: "YYYY-MM-dd hh:mm",
    dateTime: "MM-dd hh:mm",
    time: "hh:mm",
    day: "MM-dd"
  };

  export function IsSameDay(dateLeft, dateRight) {
    let dateToCheck = dateLeft;
    let actualDate = dateRight;
    if (!(dateLeft instanceof Date) || !(dateRight instanceof Date)) {
      dateToCheck = new Date(dateLeft);
      actualDate = new Date(dateRight);
    }
    return (
      dateToCheck.getFullYear() === actualDate.getFullYear() &&
      dateToCheck.getMonth() === actualDate.getMonth() &&
      dateToCheck.getDate() === actualDate.getDate()
    );
  }

  export function Format(date: Date, format: string) {
    let dateTmp = date;
    if (!(dateTmp instanceof Date)) {
      dateTmp = new Date(date);
    }
    const dateValues = [
      dateTmp.getFullYear(),
      dateTmp.getMonth() + 1,
      dateTmp.getDate(),
      dateTmp.getHours(),
      dateTmp.getMinutes(),
      dateTmp.getSeconds()
    ];
    const timeFormat = format
      .replace(/(YYYY|yyyy)/, dateValues[0].toString())
      .replace(/(MM)/, dateValues[1] < 10 ? `0${dateValues[1]}` : `${dateValues[1]}`)
      .replace(/(DD|dd)/, dateValues[2] < 10 ? `0${dateValues[2]}` : `${dateValues[2]}`)
      .replace(/(hh)/, dateValues[3] < 10 ? `0${dateValues[3]}` : `${dateValues[3]}`)
      .replace(/(mm)/, dateValues[4] < 10 ? `0${dateValues[4]}` : `${dateValues[4]}`)
      .replace(/(ss)/, dateValues[5] < 10 ? `0${dateValues[5]}` : `${dateValues[5]}`);

    return timeFormat;
  }

  export function FormatTime(time) {
    const date = new Date(time);
    const hour = date.getHours();
    const minute = date.getMinutes();
    if (hour === 0 && minute === 0) {
      return Format(date, DateFormat.day);
    }
    return Format(date, DateFormat.time);
  }

  export function FormatSeriesTime(seriesTime = []) {
    let format = DateFormat.time;
    const labelNum = seriesTime.length;
    return seriesTime.map((label, index) => {
      if (index + 1 === labelNum) {
        return Format(new Date(label), format);
      }
      if (!IsSameDay(label, seriesTime[index + 1])) {
        format = DateFormat.day;
        return Format(new Date(label), format);
      }
      const item = Format(new Date(label), format);
      format = DateFormat.time;
      return item;
    });
  }

  export function GenerateSeriesTimeVisibleLabels(min: number, max: number, labelNum: number = 12): Array<number> {
    if (max < min) {
      const temp = max;
      max = min;
      min = temp;
    }
    const defaultGaps = [300000, 600000, 900000, 1800000, 3600000, 7200000, 21600000, 43200000, 86400000]; // 单位毫秒， 5，10，15，30，60，120，360，720, 1440 分钟
    const actualGap = Math.floor((max - min) / labelNum); // labelNum 默认显示十二个label
    let gap = defaultGaps.find(gap => gap >= actualGap) || defaultGaps[defaultGaps.length - 1];


    const firstDate = new Date(min);
    let timestamp = 0;
    if (gap >= 86400000) {
      // 大于一天的间隔
      const hour = firstDate.getHours();
      const endDate = new Date(max);
      // 天数大于labelNum显示数，需要调整gap大小
      const diffDays = Math.round(Math.abs((firstDate.getTime() - endDate.getTime()) / 86400000));
      if (diffDays > labelNum) {
        gap = Math.ceil(diffDays / labelNum) * 86400000;
      }
      const remainder = hour % 24; // 小时余数
      timestamp = firstDate.setHours(remainder, 0, 0);
    } else if (gap >= 3600000) {
      let hour = firstDate.getHours();
      const gapHour = gap / 3600000; // 60分钟
      // 分钟余数
      const remainder = hour % gapHour;
      // 起始点取整点
      hour = hour + gapHour - remainder;
      timestamp = firstDate.setHours(hour, 0, 0);
    } else {
      let minute = firstDate.getMinutes();
      const gapMinute = gap / 60000; //一分钟
      // 分钟余数
      const remainder = minute % gapMinute;
      minute = minute + gapMinute - remainder;
      timestamp = firstDate.setMinutes(minute, 0);
    }

    let visibleLabels = [];
    while (timestamp < max) {
      visibleLabels.push(timestamp);
      timestamp += gap;
    }
    return visibleLabels;
  }
}

export function FormatStringNoHTMLSharp(str: string) {
  let str_tmp = str;
  if (str_tmp.indexOf("<") !== -1 && str_tmp.indexOf(">") !== -1) {
    str_tmp = str_tmp.replace("<", "");
    str_tmp = str_tmp.replace(">", "");
  }
  return str_tmp;
}

export namespace GRAPH {
  /**
   * 解决在 retina 屏幕下显示模糊
   * @param canvas
   * @param context
   * @param customWidth
   * @param customHeight
   * @constructor
   */
  export function ScaleCanvas(canvas, context, customWidth, customHeight) {
    if (!canvas || !context) {
      throw new Error("Must pass in `canvas` and `context`.");
    }

    const width = customWidth || canvas.width || canvas.clientWidth;
    const height = customHeight || canvas.height || canvas.clientHeight;

    const ratio = window.devicePixelRatio || 1;
    canvas.width = Math.round(width * ratio);
    canvas.height = Math.round(height * ratio);
    canvas.style.width = width + "px";
    canvas.style.height = height + "px";
    context.scale(ratio, ratio);
    return ratio;
  }

  export function CreateDivElement(top: number, left: number, color: string) {
    const element = document.createElement("div");
    element.style.cssText = `
          z-index: 3;
          color: ${color};
          position: absolute; 
          top: ${top}px;
          left: ${left}px;
          transform: translateY(-50%);
          text-align: left;`;
    return element;
  }
}
