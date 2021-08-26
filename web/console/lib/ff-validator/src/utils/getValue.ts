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

import { ValidatorModel, ValidatorStatusEnum } from '../Model';
import { Validation } from '../Validation';

/**
 * 获得校验结果，以 Validation的形式返回
 * @param validatorState: ValidatorModel  需要校验的validator
 * @param vKey: string | string[]  非必传，如果传入，只返回vKey的校验结果
 * @return Validation[]
 */
export const getValue = (options: { validatorState: ValidatorModel; vKey?: string | string[] }): Validation[] => {
  let { validatorState, vKey } = options;
  let finalResult: Validation[] = [];
  if (vKey) {
    let finalVKeys = vKey instanceof Array ? vKey : [vKey];
    finalVKeys.forEach(keyName => {
      let specificValidator = validatorState[keyName];
      finalResult.push(specificValidator ? specificValidator : { status: ValidatorStatusEnum.Failed, message: '' });
    });
  } else {
    for (let key in validatorState) {
      finalResult.push(validatorState[key]);
    }
  }
  return finalResult;
};
