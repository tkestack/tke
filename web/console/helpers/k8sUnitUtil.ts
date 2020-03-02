export enum K8SUNIT {
  m = 'm',
  unit = 'unit',
  K = 'k',
  M = 'M',
  G = 'G',
  T = 'T',
  P = 'P',
  Ki = 'Ki',
  Mi = 'Mi',
  Gi = 'Gi',
  Ti = 'Ti',
  Pi = 'Pi'
}

export function valueLabels1000(value, targetUnit) {
  return transformField(
    value,
    1000,
    3,
    [K8SUNIT.unit, K8SUNIT.K, K8SUNIT.M, K8SUNIT.G, K8SUNIT.T, K8SUNIT.P],
    targetUnit
  );
}

export function valueLabels1024(value, targetUnit) {
  return transformField(
    value,
    1024,
    3,
    [K8SUNIT.unit, K8SUNIT.Ki, K8SUNIT.Mi, K8SUNIT.Gi, K8SUNIT.Ti, K8SUNIT.Pi],
    targetUnit
  );
}

const UNITS = [K8SUNIT.unit, K8SUNIT.Ki, K8SUNIT.Mi, K8SUNIT.Gi, K8SUNIT.Ti, K8SUNIT.Pi];

/**
 * 进行单位换算
 * 实现k8s数值各单位之间的相互转换
 * @param {string} value
 * @param {number} thousands
 * @param {number} toFixed
 */
export function transformField(_value: string, thousands, toFixed = 3, units = UNITS, targetUnit: K8SUNIT) {
  let reg = /^(\d+(\.\d{1,2})?)([A-Za-z]+)?$/;
  let value;
  let unitBase;
  if (reg.test(_value)) {
    [value, unitBase] = [+RegExp.$1, RegExp.$3];
    if (unitBase === '') {
      unitBase = K8SUNIT.unit;
    }
  } else {
    return '0';
  }

  let i = units.indexOf(unitBase),
    targetI = units.indexOf(targetUnit);
  if (thousands) {
    if (targetI >= i) {
      while (i < targetI) {
        value /= thousands;
        ++i;
      }
    } else {
      while (targetI < i) {
        value *= thousands;
        ++targetI;
      }
    }
  }
  let svalue;
  if (value > 1) {
    svalue = value.toFixed(toFixed);
    svalue = svalue.replace(/0+$/, '');
    svalue = svalue.replace(/\.$/, '');
  } else if (value) {
    // 如果数值很小，保留toFixed位有效数字
    let tens = 0;
    let v = Math.abs(value);
    while (v < 1) {
      v *= 10;
      ++tens;
    }
    svalue = value.toFixed(tens + toFixed - 1);
    svalue = svalue.replace(/0+$/, '');
    svalue = svalue.replace(/\.$/, '');
  } else {
    svalue = value;
  }
  return String(svalue);
}
