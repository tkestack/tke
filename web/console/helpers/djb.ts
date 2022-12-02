import Cookies from 'js-cookie';

export function DJBHash(e: string) {
  let t = 5381;
  for (let n = 0, r = e.length; n < r; ++n) {
    t += (t << 5) + e.charCodeAt(n);
  }

  return t & 0x7fffffff;
}

export function createCsrfCode() {
  const tkeCookie = Cookies.get('tke');

  if (!tkeCookie) return 0;

  return DJBHash(tkeCookie);
}
