type CompareType = 'gt' | 'lt' | 'ge' | 'le' | 'eq';
/**
 * 判断集群版本是否符合要求
 * @param clusterVersion: string  当前的集群版本
 * @param targetClusterVersion: string  目标对比集群版本
 * @param type: lt | gt | ge | le | eq 对比的类型，默认为ge
 */
export const satisfyClusterVersion = (clusterVersion = '', targetClusterVersion = '', type: CompareType = 'ge') => {
  const [major, version] = (clusterVersion ?? '').split('.');
  const [minMajor, minVersion] = (targetClusterVersion ?? '').split('.');

  const majorNum = +major || 0,
    versionNum = +version || 0,
    minMajorNum = +minMajor || 0,
    minVersionNum = +minVersion || 0;

  if (type === 'ge') {
    return majorNum > minMajorNum || (majorNum === minMajorNum && versionNum >= minVersionNum);
  } else if (type === 'gt') {
    return majorNum > minMajorNum || (majorNum === minMajorNum && versionNum > minVersionNum);
  } else if (type === 'lt') {
    return majorNum < minMajorNum || (majorNum === minMajorNum && versionNum < minVersionNum);
  } else if (type === 'le') {
    return majorNum < minMajorNum || (majorNum === minMajorNum && versionNum <= minVersionNum);
  } else if (type === 'eq') {
    return majorNum === minMajorNum && versionNum === minVersionNum;
  }
};
