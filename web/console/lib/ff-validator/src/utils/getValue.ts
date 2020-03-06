import { ValidatorModel, ValidatorStatusEnum } from '../Model';
import { Validation } from '../Validation';

/**
 * 获得校验结果，以 Validation的形式返回
 * @param validatorState: ValidatorModel  需要校验的validator
 * @param vKey: string | string[]  非必传，如果传入，只返回vKey的校验结果
 * @return Validation[]
 */
export const getValue = (options: {
  validatorState: ValidatorModel;
  vKey?: string | string[];
}): Validation[] => {
  let { validatorState, vKey } = options;
  let finalResult: Validation[] = [];
  if (vKey) {
    let finalVKeys = vKey instanceof Array ? vKey : [vKey];
    finalVKeys.forEach(keyName => {
      let specificValidator = validatorState[keyName];
      finalResult.push(
        specificValidator
          ? specificValidator
          : { status: ValidatorStatusEnum.Failed, message: '' }
      );
    });
  } else {
    for (let key in validatorState) {
      finalResult.push(validatorState[key]);
    }
  }
  return finalResult;
};
