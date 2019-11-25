export function uuid() {
  // 最多使用 1e6，否则 IE toString() 会出来指数表示法
  let timeLead = 1e6;
  return (+new Date() * timeLead + Math.random() * --timeLead).toString(36);
}
