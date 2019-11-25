/**
 * 获取cookies当中的某个字段
 * @params name: string cookies当中的字段值
 */
export const getCookie = (name: string) => {
  let reg = new RegExp('(?:^|;+|\\s+)' + name + '=([^;]*)'),
    match = document.cookie.match(reg);

  return !match ? '' : match[1];
};
