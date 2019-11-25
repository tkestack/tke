export function getOffsetDate(dateParams: Date, day: number) {
  let date = new Date(dateParams.valueOf());
  date.setDate(date.getDate() + day);
  return date;
}
