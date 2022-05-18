export function isInIframe() {
  try {
    return window.self !== window.top;
  } catch (e) {
    return true;
  }
}
