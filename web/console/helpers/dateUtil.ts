/**
 * 日期格式化函数
 * @param date 日期
 * @param fmt 格式
 */
export function dateFormat(date: Date, format: string): string {
  let fmt = format;

  let o = {
    'M+': date.getMonth() + 1, //月份
    'd+': date.getDate(), //日
    'h+': date.getHours(), //小时
    'm+': date.getMinutes(), //分
    's+': date.getSeconds(), //秒
    'q+': Math.floor((date.getMonth() + 3) / 3), //季度
    S: date.getMilliseconds() //毫秒
  };
  if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length));
  for (let k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      fmt = fmt.replace(RegExp.$1, RegExp.$1.length === 1 ? o[k] : ('00' + o[k]).substr(('' + o[k]).length));
    }
  }
  return fmt;
}

/**
 * 计算时间差，返回相差的天数
 * @param startDateStr 起始日期字符串
 * @param endDateStr 结束日期字符串
 */
export function countDownDate(startDateStr: string, endDateStr?: string): number {
  let startDate = new Date(startDateStr).getTime(),
    endDate = endDateStr ? new Date(endDateStr).getTime() : Date.now(),
    diffValue = endDate - startDate;

  return Math.ceil(diffValue / (1000 * 3600 * 24));
}

/**
 * 计算剩余时间，返回hh:dd:ss，精确到秒
 * @param startDateTimeStr 起始时间字符串
 * @param endDateTimeStr 结束日期字符串
 */

function formatTime(str: string | number) {
  return ('00' + str).substr(str.toString().length);
}
export function countRestTime(startTimeStr: string, totalTimeStr?: string): string {
  let startTime = new Date(startTimeStr).getTime(),
    endTime = Date.now(),
    diffTime = endTime - startTime,
    totalTime = parseInt(totalTimeStr) * 3600 * 1000,
    restTime = Math.floor((totalTime - diffTime) / 1000);
  if (restTime <= 0) {
    return '00:00:00';
  } else {
    let hours = Math.floor(restTime / 3600),
      minutes = Math.floor((restTime % 3600) / 60),
      seconds = restTime % 60;

    return formatTime(hours) + ':' + formatTime(minutes) + ':' + formatTime(seconds);
  }
}

/**
 * 截取日期
 * @param dateStr 日期字符串
 */
export function splitDate(dateStr: string) {
  if (!dateStr) {
    return '';
  }
  let splitStr = dateStr.split(' ');
  return splitStr ? splitStr[0] : '';
}

/**
 * 格式化持续时间
 * @param seconds
 * @return string
 */
export function humanizeDuration(initSecons: number) {
  let seconds = initSecons;

  if (seconds < 0) {
    return 'N/A';
  }

  let result = '';
  if (seconds > 24 * 3600) {
    let days = Math.floor(seconds / (24 * 3600));
    result += `${days}天`;
    seconds -= days * (24 * 3600);
  }

  if (seconds > 3600) {
    let hours = Math.floor(seconds / 3600);
    result += `${hours}时`;
    seconds -= hours * 3600;
  }
  if (seconds > 60) {
    let minutes = Math.floor(seconds / 60);
    result += `${minutes}分`;
    seconds -= minutes * 60;
  }
  result += `${seconds}秒`;
  return result;
}
