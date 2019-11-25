const prefix = 'dynamic-inserted-css-';

/**
 * insert css to the current page
 */
export function insertCSS(id: string, cssText: string) {
  let style: HTMLStyleElement;
  style = document.getElementById(prefix + id) as HTMLStyleElement;
  if (!style) {
    style = document.createElement('style');
    style.id = prefix + id;

    // IE8/9 can not use document.head
    document.getElementsByTagName('head')[0].appendChild(style);
  }
  if (style.textContent !== cssText) {
    style.textContent = cssText;
  }
  return style;
}
