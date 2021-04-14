import * as compareVersions from 'compare-versions';

export const compareVersion = (firstVersion: string, secondVersion: string) => {
  // 判断是否包含“-”， 因为项目这边：1.20.4-tke.1 版本号是大于1.20.4的
  const [firstVersionPart1, firstVersionPart2] = firstVersion.split('-');
  const [secondVersionPart1, secondVersionPart2] = secondVersion.split('-');

  const compareResult = compareVersions(firstVersionPart1, secondVersionPart1);

  // 第一部分可以判断出大小
  if (compareResult !== 0) return compareResult;

  // 第二部分都不存在
  if (firstVersionPart2 === undefined && secondVersionPart2 === undefined) return 0;

  // 第二部分也相同
  if (firstVersionPart2 === secondVersionPart2) return 0;

  // firstVersionPart2不存在
  if (firstVersionPart2 === undefined) return -1;

  // secondVersionPart2不存在
  if (secondVersionPart2 === undefined) return 1;

  // 第二部分都存在
  return compareVersions(firstVersion, secondVersion);
};
