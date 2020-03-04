function parseDate(dateStr) {
  if (/^(\d{4})[-\s\.,\/]*(\d+)[-\s\.,\/]*(\d+)(?:\s*(\d+)\:(\d+):(\d+))?$/.test(dateStr)) {
    let year = +RegExp.$1;
    let month = +RegExp.$2;
    let day = +RegExp.$3;
    let hour = +RegExp.$4 || 0;
    let minute = +RegExp.$5 || 0;
    let second = +RegExp.$6 || 0;
    return [year, +month - 1, day, hour, minute, second];
  }
  return null;
}

export function getDateUTC(str) {
  return Date.UTC.apply(Date, parseDate(str));
}
