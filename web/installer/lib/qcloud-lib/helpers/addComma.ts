/**
 * 位数字添加逗号分隔
 *
 * @param {number} x 要添加分割符的数字
 * @see http://stackoverflow.com/questions/2901102/how-to-print-a-number-with-commas-as-thousands-separators-in-javascript
 */
export function addComma(x: number): string {
  if (!x) return '0';
  let parts = x.toString().split('.');
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  return parts.join('.');
}
