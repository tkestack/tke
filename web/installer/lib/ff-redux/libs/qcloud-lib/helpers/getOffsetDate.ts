export function getOffsetDate(date: Date, day: number) {
  let newDate = new Date(date.valueOf());
  newDate.setDate(date.getDate() + day);
  return date;
}
