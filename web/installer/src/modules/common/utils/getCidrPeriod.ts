export const getCidrPeriod = (cidr: string): string => {
  if (!cidr) {
    return '';
  } else {
    let splits = cidr.split('.');
    return splits.length >= 1 ? splits[0] : '';
  }
};
