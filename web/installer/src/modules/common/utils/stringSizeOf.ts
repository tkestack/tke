/**
 * 计算字符串所占的内存字节数，默认使用 utf-8的形式进行编码
 * @param userScript:string
 * @param charset:编码方式
 * @returns length  length = bytes
 */
export const stringSizeOf = (userScript: string, encode: 'utf-8' | 'utf-16' = 'utf-8') => {
  let totalSize = 0,
    charset = encode ? encode.toLowerCase() : '';

  if (charset === 'utf-16') {
    userScript.split('').forEach((item, index) => {
      let charCode = userScript.charCodeAt(index);
      if (charCode <= 0xffff) {
        totalSize += 2;
      } else {
        totalSize += 4;
      }
    });
  } else {
    userScript.split('').forEach((item, index) => {
      let charCode = userScript.charCodeAt(index);
      if (charCode <= 0x007f) {
        totalSize += 1;
      } else if (charCode <= 0x07ff) {
        totalSize += 2;
      } else if (charCode <= 0xffff) {
        // utf-8 在范围 0xD800 - 0xDFFF不存在任何字符
        totalSize += 3;
      } else {
        totalSize += 4;
      }
    });
  }

  return totalSize;
};
