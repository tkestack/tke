/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

// 计算数据盘条 之间的 距离 和数值，主要取四个点： 开始、1/5、1/2、结束
export const getRangeSliderSection = function(diskMaxSize: number, rangeWidth: number = 600) {
  return [
    {
      value: 0,
      width: 0
    },
    {
      value: Math.floor(diskMaxSize / 5),
      width: Math.floor(rangeWidth / 5)
    },
    {
      value: Math.floor(diskMaxSize / 2),
      width: Math.floor(rangeWidth / 2)
    },
    {
      value: diskMaxSize,
      width: rangeWidth
    }
  ];
};

export const getSliderMarks = function(diskMinSize: number, diskMaxSize: number, unit: string) {
  return [
    {
      value: diskMinSize,
      label: diskMinSize + unit
    },
    {
      value: Math.floor(diskMaxSize / 5),
      label: Math.floor(diskMaxSize / 5) + unit
    },
    {
      value: Math.floor(diskMaxSize / 2),
      label: Math.floor(diskMaxSize / 2) + unit
    },
    {
      value: diskMaxSize,
      label: diskMaxSize + unit
    }
  ];
};
