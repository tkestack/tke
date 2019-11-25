export const getCidrMask = (cidr: string): string => {
  if (!cidr) {
    return '';
  } else {
    let splits = cidr.split('/');
    return splits.length >= 2 ? splits[1] : '';
  }
};
