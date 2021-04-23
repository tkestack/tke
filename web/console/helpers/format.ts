// 所有单位都换算成byte
const memoryMap1000 = ['K', 'M', 'G', 'T', 'P', 'E'].reduce(
  (all, key, index) => Object.assign(all, { [key]: Math.pow(1000, index + 1) }),
  {}
);

const memoryMap1024 = ['KI', 'MI', 'GI', 'TI', 'PI', 'EI'].reduce(
  (all, key, index) => Object.assign(all, { [key]: Math.pow(1024, index + 1) }),
  {}
);

const memoryMap = Object.assign(memoryMap1000, memoryMap1024);

export const formatMemory = (
  memory: string,
  finalUnit: 'K' | 'M' | 'G' | 'T' | 'P' | 'E' | 'Ki' | 'Mi' | 'Gi' | 'Ti' | 'Pi' | 'Ei'
) => {
  const unit = memory.toUpperCase().match(/[KMGTPEI]+/)?.[0] ?? 'MI';

  const memoryNum = parseInt(memory);

  return `${((memoryNum * memoryMap[unit]) / memoryMap[finalUnit.toUpperCase()]).toLocaleString()} ${finalUnit}`;
};
