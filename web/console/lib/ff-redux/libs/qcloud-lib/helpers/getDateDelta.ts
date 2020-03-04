/**Parses string formatted as YYYY-MM-DD to a Date object.
 * If the supplied string does not match the format, an
 * invalid Date (value NaN) is returned.
 * @param {string} dateStringInRange format YYYY-MM-DD, with year in
 * range of 0000-9999, inclusive.
 * @return {Date} Date object representing the string.
 */
function parseISO8601(dateStringInRange) {
  let isoExp = /^\s*(\d{4})-(\d\d)-(\d\d)\s*/,
    date = new Date(NaN),
    month,
    parts = isoExp.exec(dateStringInRange);

  if (parts) {
    month = +parts[2];
    date.setFullYear(+parts[1], month - 1, +parts[3]);
    if (month !== date.getMonth() + 1) {
      date.setTime(NaN);
    }
  }
  return date;
}

const ONE_DAY = 24 * 60 * 60 * 1000;

export function getDateDelta(date1, date2) {
  return Math.floor(+parseISO8601(date2) - +parseISO8601(date1)) / ONE_DAY;
}
