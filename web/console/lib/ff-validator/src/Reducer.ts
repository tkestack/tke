import { ValidateSchema } from './Model';
import { Reducer } from 'redux';
import { reduceToPayload } from './utils/reduceToPayload';
import { getValidatorActionType } from './ActionType';

/**
 * 获取校验的Reducer，注入到reducer当中
 * @param userDefinedSchema: Schema[] 用户自定义的校验组件
 * @return subReducer
 */
export const createValidatorReducer = (userDefinedSchema: ValidateSchema): Reducer => {
  let formKey = userDefinedSchema.formKey;
  return reduceToPayload(getValidatorActionType(formKey), {});
};
