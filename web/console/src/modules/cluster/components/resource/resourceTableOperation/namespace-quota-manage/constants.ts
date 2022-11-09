import { t } from '@tencent/tea-app/lib/i18n';

export enum UnitTypeEnum {
  MEMORY,
  CPU,
  OTHER
}

export enum MemoryUnitEnum {
  Ki = 'Ki',
  Mi = 'Mi',
  Gi = 'Gi',
  Ti = 'Ti',
  Pi = 'Pi',
  Ei = 'Ei',

  K = 'K',
  M = 'M',
  G = 'G',
  T = 'T',
  P = 'P',
  E = 'E',

  Byte = 'Byte',

  //特殊的m,代表毫字节
  m = 'm'
}

export const DefaultUnitByType = {
  [UnitTypeEnum.MEMORY]: MemoryUnitEnum.Gi,
  [UnitTypeEnum.CPU]: t('核'),
  [UnitTypeEnum.OTHER]: t('个')
};

// 顺序很重要
export const MemoryUnits = [
  MemoryUnitEnum.Mi,
  MemoryUnitEnum.Gi,
  MemoryUnitEnum.Ti,
  MemoryUnitEnum.M,
  MemoryUnitEnum.G,
  MemoryUnitEnum.T,
  MemoryUnitEnum.m
];

/**
 * 内存單位 -> bytes 的滙率
 */
export const MemoryExchangeRate = {
  [MemoryUnitEnum.K]: 1000,
  [MemoryUnitEnum.M]: 1000 * 1000,
  [MemoryUnitEnum.G]: 1000 * 1000 * 1000,
  [MemoryUnitEnum.T]: 1000 * 1000 * 1000 * 1000,
  [MemoryUnitEnum.P]: 1000 * 1000 * 1000 * 1000 * 1000,
  [MemoryUnitEnum.E]: 1000 * 1000 * 1000 * 1000 * 1000 * 1000,

  [MemoryUnitEnum.Ki]: 1024,
  [MemoryUnitEnum.Mi]: 1024 * 1024,
  [MemoryUnitEnum.Gi]: 1024 * 1024 * 1024,
  [MemoryUnitEnum.Ti]: 1024 * 1024 * 1024 * 1024,
  [MemoryUnitEnum.Pi]: 1024 * 1024 * 1024 * 1024 * 1024,
  [MemoryUnitEnum.Ei]: 1024 * 1024 * 1024 * 1024 * 1024 * 1024,

  [MemoryUnitEnum.Byte]: 1,

  [MemoryUnitEnum.m]: 1 / 1000
};

export const MemoryUnitSelectOptions = [
  { text: MemoryUnitEnum.Mi, value: MemoryUnitEnum.Mi },
  { text: MemoryUnitEnum.Gi, value: MemoryUnitEnum.Gi },
  { text: MemoryUnitEnum.Ti, value: MemoryUnitEnum.Ti },

  { text: MemoryUnitEnum.M, value: MemoryUnitEnum.M },
  { text: MemoryUnitEnum.G, value: MemoryUnitEnum.G },
  { text: MemoryUnitEnum.T, value: MemoryUnitEnum.T }
];
