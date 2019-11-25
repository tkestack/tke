import { Validation } from 'src/modules/common';

export interface Rule {
  /**标签名 */
  label?: string;

  /**是否为必填 */
  required?: boolean;

  /**校验正则 */
  reg?: RegExp;

  /**最小长度 */
  minLength?: number;

  /**最大长度 */
  maxLength?: number;

  /**最小值 */
  minValue?: number;

  /**最大值 */
  maxValue?: number;
}

function validateRequired(value, rule: Rule) {
  if (typeof value === 'number') {
    return { status: 1, message: '' };
  }

  let result = !!value;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}不能为空`
  };
}

function validateMinLength(value, rule: Rule) {
  let result = value.length >= rule.minLength;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}长度不能小于${rule.minLength}位`
  };
}

function validateMaxLength(value, rule: Rule) {
  let result = value.length <= rule.maxLength;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}长度不能大于${rule.maxLength}位`
  };
}

function validateMinValue(value, rule: Rule) {
  let result = value >= rule.minValue;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}值不能小于${rule.minValue}`
  };
}

function validateMaxValue(value, rule: Rule) {
  let result = value <= rule.maxValue;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}值不能大于${rule.maxValue}`
  };
}

function validateReg(value, rule: Rule) {
  let result = rule.reg.test(value);
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}格式不正确`
  };
}

export function Validate(value, rule: Rule): Validation {
  let result = {
    status: 0,
    message: ''
  };

  if (!rule) {
    return result;
  }

  if (rule.required) {
    result = validateRequired(value, rule);
  } else {
    if (value === '' || value === undefined) {
      return { status: 1, message: '' };
    }
  }

  if (result.status !== 2 && (rule.minLength || rule.minLength === 0)) {
    result = validateMinLength(value, rule);
  }

  if (result.status !== 2 && (rule.maxLength || rule.minLength === 0)) {
    result = validateMaxLength(value, rule);
  }

  if (result.status !== 2 && rule.reg) {
    result = validateReg(value, rule);
  }

  if (result.status !== 2 && rule.minValue) {
    result = validateMinValue(value, rule);
  }

  if (result.status !== 2 && rule.maxValue) {
    result = validateMaxValue(value, rule);
  }

  if (result.status !== 2) {
    result = {
      status: 1,
      message: ''
    };
  }

  return result;
}

export function isValidateSuccess(validates: Validation[]) {
  let result = true;
  validates.forEach(v => {
    result = result && v.status !== 2;
  });
  return result;
}
