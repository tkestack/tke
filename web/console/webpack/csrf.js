const SparkMD5 = require('spark-md5');

function parseCookie(cookieStr) {
  return cookieStr.split('; ').reduce((all, item) => {
    const [key, value] = item.split('=');

    return {
      ...all,
      [key]: value
    };
  }, {});
}

function createCSRFHeader(cookieStr) {
  const cookie = parseCookie(cookieStr);

  const tkeCookie = cookie?.['tke'] ?? '';

  const token = SparkMD5.hash(tkeCookie);

  return {
    'X-CSRF-TOKEN': token
  };
}

module.exports = {
  createCSRFHeader
};
