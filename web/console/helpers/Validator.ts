/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { ControllerFieldState, UseFormStateReturn } from 'react-hook-form';
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

  const result = !!value;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}不能为空`
  };
}

function validateMinLength(value, rule: Rule) {
  const result = value.length >= rule.minLength;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}长度不能小于${rule.minLength}位`
  };
}

function validateMaxLength(value, rule: Rule) {
  const result = value.length <= rule.maxLength;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}长度不能大于${rule.maxLength}位`
  };
}

function validateMinValue(value, rule: Rule) {
  const result = value >= rule.minValue;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}值不能小于${rule.minValue}`
  };
}

function validateMaxValue(value, rule: Rule) {
  const result = value <= rule.maxValue;
  return {
    status: result ? 1 : 2,
    message: result ? '' : `${rule.label}值不能大于${rule.maxValue}`
  };
}

function validateReg(value, rule: Rule) {
  const result = rule.reg.test(value);
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

export function getReactHookFormStatusWithMessage({
  fieldState,
  formState
}: {
  fieldState: ControllerFieldState;
  formState: UseFormStateReturn<any>;
}): {
  status?: 'error' | 'success';
  message?: string;
} {
  if (!fieldState.isTouched && !fieldState.isDirty && !formState.isSubmitted) {
    return {};
  }

  return fieldState.invalid
    ? {
        status: 'error',
        message: fieldState?.error?.message ?? ''
      }
    : {
        status: 'success'
      };
}
