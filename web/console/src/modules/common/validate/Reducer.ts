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

import { combineReducers, Reducer } from 'redux';

import { reduceToPayload } from '@tencent/ff-redux';

import { initValidator } from '../models';
import { getValidationActionType } from './ActionType';
import { ValidateSchema } from './Model';

/**
 * 获取校验的Reducer，注入到reducer当中
 * @param userDefinedSchema: Schema[] 用户自定义的校验组件
 * @return subReducer
 */
export const generateValidateReducer = (userDefinedSchema: ValidateSchema): Reducer => {
  return (state, action) => {
    let validateReducer = {},
      formKey = userDefinedSchema.formKey;
    userDefinedSchema.fields.forEach(ins => {
      let fieldKey = ins.vKey;
      validateReducer[fieldKey] = reduceToPayload(getValidationActionType(formKey, fieldKey), initValidator);
    });
    return combineReducers(validateReducer)(state, action);
  };
};
