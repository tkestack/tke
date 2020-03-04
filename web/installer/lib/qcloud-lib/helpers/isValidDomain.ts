export function isValidDomain(domain: string) {
  // regex from https://github.com/johnotander/domain-regex/blob/master/index.js
  return /^((?=[a-z0-9-]{1,63}\.)(xn--)?[a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,63}$/.test(domain);
}
