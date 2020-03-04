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
