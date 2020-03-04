/**
 * 检查一个 IP 地址是否为合法的 IP 地址
 * */
export function isValidIPAddress(ipAddress: string) {
  if (!ipAddress) return false;

  let newIpAddress = ipAddress.trim();

  const segments = newIpAddress.split('.');
  if (segments.length !== 4) return false;

  return segments.reduce((prev, curr) => {
    const value = parseInt(curr, 10);
    return prev && value >= 0 && value <= 255 && String(value) === curr;
  }, true);
}
