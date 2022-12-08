import Cookies from 'js-cookie';
import SparkMD5 from 'spark-md5';

let CSRF_TOKEN = null;

export function createCSRFHeader() {
  if (CSRF_TOKEN === null) {
    const tkeCookie = Cookies.get('tke') ?? '';
    CSRF_TOKEN = SparkMD5.hash(tkeCookie);
  }

  return {
    'X-CSRF-TOKEN': CSRF_TOKEN
  };
}
