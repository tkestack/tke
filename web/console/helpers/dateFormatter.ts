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
export function dateFormatter(date: Date, format: string) {
  const o = {
    /**
     * 完整年份
     * @example 2015 2016 2017 2018
     */
    YYYY() {
      return date.getFullYear().toString();
    },

    /**
     * 年份后两位
     * @example 15 16 17 18
     */
    YY() {
      return this.YYYY().slice(-2);
    },

    /**
     * 月份，保持两位数
     * @example 01 02 03 .... 11 12
     */
    MM() {
      return leftPad(this.M(), 2);
    },

    /**
     * 月份
     * @example 1 2 3 .... 11 12
     */
    M() {
      return (date.getMonth() + 1).toString();
    },

    /**
     * 每月中的日期，保持两位数
     * @example 01 02 03 .... 30 31
     */
    DD() {
      return leftPad(this.D(), 2);
    },

    /**
     * 每月中的日期
     * @example 1 2 3 .... 30 31
     */
    D() {
      return date.getDate().toString();
    },

    /**
     * 小时，24 小时制，保持两位数
     * @example 00 01 02 .... 22 23
     */
    HH() {
      return leftPad(this.H(), 2);
    },

    /**
     * 小时，24 小时制
     * @example 0 1 2 .... 22 23
     */
    H() {
      return date.getHours().toString();
    },

    /**
     * 小时，12 小时制，保持两位数
     * @example 00 01 02 .... 22 23
     */
    hh() {
      return leftPad(this.h(), 2);
    },

    /**
     * 小时，12 小时制
     * @example 0 1 2 .... 22 23
     */
    h() {
      const h = (date.getHours() % 12).toString();
      return h === '0' ? '12' : h;
    },

    /**
     * 分钟，保持两位数
     * @example 00 01 02 .... 59 60
     */
    mm() {
      return leftPad(this.m(), 2);
    },

    /**
     * 分钟
     * @example 0 1 2 .... 59 60
     */
    m() {
      return date.getMinutes().toString();
    },

    /**
     * 秒，保持两位数
     * @example 00 01 02 .... 59 60
     */
    ss() {
      return leftPad(this.s(), 2);
    },

    /**
     * 秒
     * @example 0 1 2 .... 59 60
     */
    s() {
      return date.getSeconds().toString();
    }
  };

  return Object.keys(o).reduce((pre, cur) => {
    return pre.replace(new RegExp(cur), match => {
      /* eslint-disable */
      return o[match].call(o);
      /* eslint-enable */
    });
  }, format);
}

function leftPad(num: number | string, width: number, c: number | string = '0'): string {
  const numStr = num.toString();
  const padWidth = width - numStr.length;
  return padWidth > 0 ? new Array(padWidth + 1).join(c.toString()) + numStr : numStr;
}

/**
 * 计算两个日期的时间差
 */
export function dateCompare(date1, inputDate2) {
  let date2 = typeof inputDate2 !== 'undefined' ? inputDate2 : new Date();

  //时间差的毫秒数
  let date3 = date2.getTime() - date1.getTime();

  //计算出相差天数
  let days = Math.floor(date3 / (24 * 3600 * 1000));

  //计算出小时数
  let leave1 = date3 % (24 * 3600 * 1000); //计算天数后剩余的毫秒数
  let hours = Math.floor(leave1 / (3600 * 1000));

  //计算相差分钟数
  let leave2 = leave1 % (3600 * 1000); //计算小时数后剩余的毫秒数
  let minutes = Math.floor(leave2 / (60 * 1000));

  //计算相差秒数
  let leave3 = leave2 % (60 * 1000); //计算分钟数后剩余的毫秒数
  let seconds = Math.round(leave3 / 1000);

  let ret = '';
  if (days > 0) {
    ret += days + '天';
  }
  if (days > 0 || hours > 0) {
    ret += hours + '小时';
  }
  if (days > 0 || hours > 0 || minutes > 0) {
    ret += minutes + '分钟';
  }
  if (days > 0 || hours > 0 || minutes > 0 || seconds > 0) {
    ret += seconds + '秒';
  }

  return ret;
}
