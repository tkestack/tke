export function versionBigThanOrEqual(v1: string, v2: string): boolean {
  const [_, v1_part2] = v1.split('.');
  const [__, v2_part2] = v2.split('.');

  const len = Math.max(v1_part2.length, v2_part2.length);

  const v1Number = parseFloat(v1) * len;
  const v2Number = parseFloat(v2) * len;

  return v1Number >= v2Number;
}
