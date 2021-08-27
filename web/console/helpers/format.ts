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
