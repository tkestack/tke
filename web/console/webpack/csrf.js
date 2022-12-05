function DJBHash(e) {
  let t = 5381;
  for (let n = 0, r = e.length; n < r; ++n) {
    t += (t << 5) + e.charCodeAt(n);
  }

  return t & 0x7fffffff;
}

function createCsrfCode(cookieStr) {
  const cookie = cookieStr.split('; ').reduce((all, item) => {
    const [key, value] = item.split('=');

    return {
      ...all,
      [key]: value
    };
  }, {});

  const tkeCookie = cookie?.['tke'];

  if (!tkeCookie) return 0;

  const code = DJBHash(tkeCookie);

  return code;
}

module.exports = {
  createCsrfCode
};
