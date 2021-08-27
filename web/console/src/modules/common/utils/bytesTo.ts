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
/**
 * Bytes -> GB/MB/KB/Bytes
 */
export function bytesTo(bytes, type = 'auto' as 'GB' | 'MB' | 'KB' | 'Bytes' | 'auto') {
  let result = bytes,
    _type = 'Bytes',
    _GB = 1000 * 1000 * 1000,
    _MB = 1000 * 1000,
    _KB = 1000;

  if (type === 'auto') {
    if (bytes >= _GB) {
      _type = 'GB';
    } else if (bytes >= _MB && bytes < _GB) {
      _type = 'MB';
    } else if (bytes < _MB && bytes >= _KB) {
      _type = 'KB';
    } else {
      _type = 'Bytes';
    }
  } else {
    _type = type;
  }

  switch (_type) {
    case 'GB': {
      result = Math.round(bytes / _GB).toFixed(1);
      break;
    }
    case 'MB': {
      result = Math.round(bytes / _MB).toFixed(1);
      break;
    }
    case 'KB': {
      result = Math.round(bytes / _KB).toFixed(1);
      break;
    }
    default: {
      result = bytes.toFixed(0);
      break;
    }
  }
  return `${result} ${_type}`;
}

const UNITS = ['', 'K', 'M', 'G', 'T', 'P'];

/**
 * 进行单位换算
 * @param {number} value
 * @param {number} thousands
 * @param {number} toFixed
 */
export function TransformField(_value: number, thousands, toFixed = 3, units = UNITS) {
  let value = _value;
  let isValueDefined = !isNaN(value) && value !== null;
  if (!isValueDefined) return '-';

  let unitBase = units[0];
  let i = units.indexOf(unitBase);
  if (isValueDefined && thousands) {
    while (i < units.length && value / thousands > 1) {
      value /= thousands;
      ++i;
    }
    unitBase = units[i] || '';
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
  return String(svalue) + (value !== 0 ? unitBase : '');
}
