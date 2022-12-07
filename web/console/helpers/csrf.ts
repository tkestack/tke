import Cookies from 'js-cookie';
import SparkMD5 from 'spark-md5';

export function createCSRFHeader() {
  const tkeCookie = Cookies.get('tke') ?? '';
  const token = SparkMD5.hash(tkeCookie);

  return {
    'X-CSRF-TOKEN': token
  };
}
