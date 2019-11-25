import { ValidateSchema, RuleTypeEnum, Rule, ValidationIns, FieldConfig, ValidatorStatusEnum } from './Model';
import { getValidationActionType } from './ActionType';
import { Validation, initValidator } from '../models';

/**
 * 构造表单校验全局实例
 * @param userDefinedSchema: ValidateSchema<T>  表单校验的完整配置
 * @param getData: () => any  获取store的方法
 * @param getValidatorData: () => any 获取Validation[]的结果，指定从store的什么地方取
 * @param dispatch: Redux.Dispatch  Redux的触发方式
 * @return;
 */
export function generateValidationIns(options: {
  userDefinedSchema: ValidateSchema;
  getData: () => any;
  getValidatorData: () => any;
  dispatch: Redux.Dispatch;
}): ValidationIns {
  let { userDefinedSchema, getData, getValidatorData, dispatch } = options;

  /**
   * 触发校验的action
   * @param vKey?: string 非必传  传入只校验对应vKey的值
   */
  const validate = (vKey?: string | string[]): void => {
    let finalStore = getData();
    let finalValidateField: FieldConfig[];
    let fields: FieldConfig[] = userDefinedSchema.fields;

    let finalVKeys = vKey instanceof Array ? vKey : vKey ? [vKey] : [];
    /** 如果设置了vKey，则只校验vKey的rules */
    if (finalVKeys.length) {
      let hasDefined = fields.filter(item => finalVKeys.includes(item.vKey));
      finalValidateField = hasDefined.length ? hasDefined : [];
    } else {
      finalValidateField = fields;
    }

    finalValidateField.forEach(field => {
      let result: Validation = initValidator;
      for (let rule of field.rules) {
        let ruleType = typeof rule === 'string' ? rule : rule.type;
        result = validateMethod({
          store: finalStore,
          value: finalStore[field.vKey],
          label: field.label,
          rule: rule as Rule
        })[ruleType]();
        // 自己自定义的rules当中，只要有一个不符合，就不需要继续向下检查
        if (result.status !== ValidatorStatusEnum.Success) break;
      }
      // 触发对应的Action，由Reducer修改对应的store的值
      dispatch({
        type: getValidationActionType(userDefinedSchema.formKey, field.vKey),
        payload: result
      });

      /**
       * @pre 只有传入了vKey，才需要进行递归校验，因为不传vKey，默认都是校验所有选项，不需要进行额外的校验，不进行递归
       * @pre 如果传入的key包含需要校验dependent的，也不需要进行递归
       * 需要校验当前field是否被其他field所依赖，如input2依赖于input1的选项，则校验input1，需要重新校验input2
       */
      if (finalVKeys.length) {
        let dependentFields: string[] = fields
          .filter(item => !finalVKeys.includes(item.vKey) && item.dependentKey === field.vKey)
          .map(item => item.vKey);
        dependentFields.length && validate(dependentFields);
      }
    });
  };

  /**
   * 获得校验结果，以 Validation的形式返回
   * @param vKey: string  非必传，如果传入，只返回vKey的校验结果
   * @return Validation : Validation[]
   */
  const getValue = (vKey?: string): Validation | Validation[] => {
    let finalStore = getValidatorData();
    let finalResult: Validation[] = [];
    if (vKey) {
      let specificValidator = finalStore[vKey];
      finalResult.push(specificValidator ? specificValidator : { status: ValidatorStatusEnum.Failed, message: '' });
    } else {
      for (let key in finalStore) {
        finalResult.push(finalStore[key]);
      }
    }
    return vKey ? finalResult[0] : finalResult;
  };

  /**
   * 获得校验结果，以 boolean 的形式返回
   * @param vKey: string  非必传，如果传入，只返回vKey的校验结果
   * @return boolean
   */
  const isValid = (vKey?: string): boolean => {
    let finalStore = getValidatorData();
    let finalResult: boolean = true;
    if (vKey) {
      finalResult = finalStore[vKey] && finalStore[vKey].status === ValidatorStatusEnum.Success ? true : false;
    } else {
      for (let key in finalStore) {
        finalResult = finalResult && (finalStore[key] as Validation).status === ValidatorStatusEnum.Success;
        if (!finalResult) break;
      }
    }
    return finalResult;
  };

  return {
    validate,
    getValue,
    isValid
  };
}

interface ValidateMethodOptions {
  /** 传入的全局的state */
  store?: any;

  /** 需被校验的值 */
  value: any;

  /** 校验项的label */
  label: string;

  /** 校验的规则 */
  rule: Rule;
}

/**
 * 实际校验的方式
 * @param options: ValidateMethodOptions  需要传入的校验选项
 */
const validateMethod = (
  options: ValidateMethodOptions
): {
  [props: string]: () => Validation;
} => {
  let { store, value, label, rule } = options;

  const validateRequired = (): Validation => {
    let result = !!value;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}不能为空`
    };
  };

  const validateMinLength = (): Validation => {
    let minLength = rule.limit ? +rule.limit : 0;
    let result = value && value.length >= minLength;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}至少包含${rule.limit}个字符`
    };
  };

  const validateMaxLength = (): Validation => {
    let maxLength = rule.limit ? +rule.limit : undefined;
    let result = maxLength ? value && value.length <= maxLength : true;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}不能超过${maxLength}个字符`
    };
  };

  const validateMinValue = (): Validation => {
    let minValue = rule.limit ? +rule.limit : 0;
    let result = !Number.isNaN(value) && value >= minValue;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}最小值为${minValue}`
    };
  };

  const validateMaxValue = (): Validation => {
    let maxValue = rule.limit ? +rule.limit : Number.MAX_SAFE_INTEGER;
    let result = !Number.isNaN(value) && value <= maxValue;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}最大值为${maxValue}`
    };
  };

  const validateRegExp = (): Validation => {
    let result = rule.limit ? rule.limit.test(value) : false;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}格式不正确`
    };
  };

  const validateMinCheckBoxCount = (): Validation => {
    let minValue = rule.limit ? +rule.limit : 0;
    let result = !Number.isNaN(value) && value >= minValue;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}至少选择${minValue}项`
    };
  };

  const validateMaxCheckBoxCount = (): Validation => {
    let maxValue = rule.limit ? +rule.limit : Number.MAX_SAFE_INTEGER;
    let result = !Number.isNaN(value) && value <= maxValue;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : rule.errorTip ? rule.errorTip : `${label}至多选择${maxValue}项`
    };
  };

  const validateCustom = (): Validation => {
    return rule.customFunc ? rule.customFunc(value, store) : { status: ValidatorStatusEnum.Success, message: '' };
  };

  return {
    [RuleTypeEnum.isRequire]: validateRequired,
    [RuleTypeEnum.maxValue]: validateMaxValue,
    [RuleTypeEnum.minValue]: validateMinValue,
    [RuleTypeEnum.minLength]: validateMinLength,
    [RuleTypeEnum.maxLength]: validateMaxLength,
    [RuleTypeEnum.regExp]: validateRegExp,
    [RuleTypeEnum.custom]: validateCustom,
    [RuleTypeEnum.minCheckBoxCount]: validateMinCheckBoxCount,
    [RuleTypeEnum.maxCheckBoxCount]: validateMaxCheckBoxCount
  };
};

interface UniqValidateMethodOptions {
  /** 传入的全局的state */
  store?: any;

  /** 校验项的label */
  label: string;

  /** rules的配置 */
  rules: (Rule | string)[];
}

/**
 * 提供单独的校验选项
 * @param value: any  需要被校验的值
 * @param options: ValidateMethodOptions  需要传入的校验项
 */
export function validateValue(value: any, options: UniqValidateMethodOptions): Validation {
  let { rules, ...restOptions } = options;
  let result: Validation = initValidator;
  for (let rule of rules) {
    let ruleType = typeof rule === 'string' ? rule : rule.type;
    result = validateMethod(Object.assign({}, restOptions, { rule: rule as Rule, value }))[ruleType]();
    // 自己自定义的rules当中，只要有一个不符合，就不需要向下检查了
    if (result.status !== ValidatorStatusEnum.Success) break;
  }
  return result;
}
