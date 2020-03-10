import { ValidatorStatusEnum, ValidatorModel } from '../Model';
import { Validation } from '../Validation';

/**
 * 获得校验结果，以 boolean 的形式返回
 * @param validatorState: ValidatorModel 需要校验的validator
 * @return boolean
 */
export const isValid = (validatorState: ValidatorModel): boolean => {
  let finalResult = true;
  for (let key in validatorState) {
    finalResult =
      finalResult && validatorState[key].status !== ValidatorStatusEnum.Failed;
    if (!finalResult) break;
  }
  return finalResult;
};
