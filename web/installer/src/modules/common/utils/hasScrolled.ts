/**
 * 判断目标元素是否出现了滚动条
 */

export const hasScrolled = (selector: string) => {
  let target = document.getElementById(selector);

  if (target && target.scrollHeight && target.clientHeight) {
    return target.scrollHeight > target.clientHeight;
  } else {
    return false;
  }
};
