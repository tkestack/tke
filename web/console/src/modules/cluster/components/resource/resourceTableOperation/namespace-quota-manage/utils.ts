import { MemoryUnits, MemoryExchangeRate, MemoryUnitEnum, UnitTypeEnum } from './constants';

export function getObjectPropByPath({ path, record }: { path: Array<string | number>; record: object }): string | null {
  return path.reduce((pre, path) => pre?.[path], record) ?? null;
}

export function getPercentWithStatus({ used, all }: { used: number; all: number }) {
  let percent = all && (used / all) * 100;

  percent = Math.min(percent, 100);

  const status: 'danger' | 'default' = percent > 50 ? 'danger' : 'default';

  return {
    percent,
    status
  };
}

/**
 * 获取memory单位, 没有单位的为Byte
 */
export function getUnitByMemory(value: string) {
  return MemoryUnits.find(unit => value?.includes(unit)) ?? MemoryUnitEnum.Byte;
}

// 内存单位转换
export function transMemoryUnit(value: string, unit: MemoryUnitEnum) {
  if (value === null) return 0;

  const originUnit = getUnitByMemory(value);

  const bytesValue = parseFloat(value) * MemoryExchangeRate[originUnit];

  return bytesValue / MemoryExchangeRate[unit];
}

// cpu 单位转换
export function transCpuUnit(value: string) {
  if (value === null) return 0;

  if (value?.includes('m')) {
    return parseFloat(value) / 1000;
  } else {
    return parseFloat(value);
  }
}

// 其他类型的转换
export function transOthers(value: string) {
  if (value === null) return 0;

  return parseFloat(value);
}

export function getMemoryValueWithUnit(m: string) {
  const unit = getUnitByMemory(m);
  const value = parseFloat(m);

  return [value, unit];
}

// 比较带单位的memory
export function compareMemory(value1: string, value2: string) {
  const value1ToByte = transMemoryUnit(value1, MemoryUnitEnum.Byte);

  const value2ToByte = transMemoryUnit(value2, MemoryUnitEnum.Byte);

  return value1ToByte - value2ToByte;
}

// 比较cpu
export function compareCpu(value1: string, value2: string) {
  const value1ToDefault = transCpuUnit(value1);
  const value2ToDefault = transCpuUnit(value2);

  return value1ToDefault - value2ToDefault;
}

// 比较其他的类型
export function compareOther(value1: string, value2: string) {
  const value1ToDefault = transOthers(value1);
  const value2ToDefault = transOthers(value2);

  return value1ToDefault - value2ToDefault;
}
