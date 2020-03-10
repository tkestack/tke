import { ValidatorModel } from '../Model';

/**
 * 获得校验结果，以 Validation的形式返回
 * @param validatorState: ValidatorModel  需要校验的validator
 * @return 返回所有指定的key的错误信息的集合
 */
export const getMessage = (validatorState: ValidatorModel): string[] => {
  let finalResult = [];

  for (let key in validatorState) {
    let hasMessage = validatorState[key].message;
    hasMessage && finalResult.push(hasMessage);
  }
  return finalResult;
};
