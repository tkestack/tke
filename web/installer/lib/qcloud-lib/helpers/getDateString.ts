function pad(num: number) {
  let r = String(num);
  if (r.length === 1) {
    r = `0${r}`;
  }
  return r;
}

export function getDateString(date: Date) {
  return [date.getFullYear(), pad(date.getMonth() + 1), pad(date.getDate())].join('');
}
