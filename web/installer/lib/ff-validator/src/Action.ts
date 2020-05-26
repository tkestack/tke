import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { getValidatorActionType } from './ActionType';
import {
    FieldConfig, ModelTypeEnum, Rule, RuleTypeEnum, ValidateMethodOptions, ValidateSchema,
    ValidationIns, ValidatorModel, ValidatorStatusEnum
} from './Model';
import { cloneDeep, formatVkeyRegExp } from './utils';
import { initValidator, Validation } from './Validation';

/**
 * 根据fields，即用户的自定义校验规则配置 以及 用户传入的vKey，获得初始需要校验的vKeys的数组
 * @param fields: FieldConfig[] 用户自定义的校验规则配置
 * @param vKey: string | string[] 用户调用 validate 时传入的 vKey
 * @return FieldConfig[]  用户自定义校验规则当中符合规则的校验项
 */
const getValidateConfig = (fields: FieldConfig[], vKey?: string | string[]): SchemaVkey[] => {
  let originalVKey = vKey instanceof Array ? vKey : vKey ? [vKey] : [];
  // 对vKey进行格式化，即传入的vKey是 labels[1].key => labels[].key，用于匹配用户自定义的校验规则
  let formatVKeys = originalVKey.map(item => item.replace(formatVkeyRegExp, '[]'));
  // 最终命中的自定义配置项规则的vKey
  let hitSchemaVkey = [];
  // 如果设置了vKey，则只校验符合vKey的rules
  if (formatVKeys.length) {
    formatVKeys.forEach((item, index) => {
      let hasDefined = fields
        .filter(field => field.vKey.startsWith(item))
        .map(field => {
          // 需要将传入的vKey 和 当前的field的vKey进行结合，如field.vKey为containers[].name，传入的vKey为containers[0]，则结合为 containets[0].name
          let originalVKeyArr = originalVKey[index].split('.');
          let fieldVkeyArr = field.vKey.split('.');
          fieldVkeyArr.splice(0, originalVKeyArr.length, ...originalVKeyArr);
          return fieldVkeyArr.join('.');
        });
      hitSchemaVkey.push(...hasDefined);
    });
  } else {
    hitSchemaVkey = fields.map(field => field.vKey);
  }
  return hitSchemaVkey.map(schemaVkey => ({
    schemaVkey,
    value: ''
  }));
};

interface SchemaVkey {
  /** 校验的vKey，如 containers[0].name */
  schemaVkey: string;

  /** 根据可用path，获取到当前校验key的值 */
  value: any;
}

/**
 * 对需要校验的形如 SchemaVkey 的数组进行展开，将包含通配符的进行平铺
 * @param schemaVkey: SchemaVKey 根据用户定义和传入key获取的初始schemaVkey
 * @param getStoreValue: (path: string, store: any) => any  根据传入的path 和 store，获得对应的值
 * @return 返回一个将所有通配 [] 都展开的数组
 */
const expandSchemaVkey = (schemaVkeyConfig: SchemaVkey, store: any): SchemaVkey[] => {
  let finalSchemaVkeyArr: SchemaVkey[] = [];
  let { schemaVkey } = schemaVkeyConfig;
  // 判断当前 schemaVkey是否包含通配符 []，并且记录其位置
  let firstBracketOnlyIndex = schemaVkey.indexOf('[]');
  if (firstBracketOnlyIndex > -1) {
    // 如果能找到通配符 []，说明需要对其进行展开，如containers[0].labels[] => containers[0].labels[0]……
    let useablePath = schemaVkey.slice(0, firstBracketOnlyIndex);
    let currentStore: any[] = getStoreValue(useablePath, store);
    let restPartExceptBracket = schemaVkey.slice(firstBracketOnlyIndex + 2);
    currentStore.forEach((item, index) => {
      finalSchemaVkeyArr.push({
        schemaVkey: `${useablePath}[${index}]${restPartExceptBracket}`,
        value: ''
      });
    });
  }
  return finalSchemaVkeyArr;
};

/**
 * 根据path 和 store，获得当前的值
 * @param path: string  需要在store树当中递归的path
 * @param store: any  原始数据值
 */
const getStoreValue = (path: string, store: any): any => {
  let pathArr = path.split(/[\[\]\.]/).filter(item => item !== '');
  let currentStore = store;
  pathArr.forEach(item => {
    currentStore = currentStore[item];
  });
  return currentStore;
};

/**
 * 根据valueField获取最终的处理的值
 * @param field: FieldConfig  校验的配置
 * @return value: any
 */
const getValueByValueField = (field: FieldConfig, value: any) => {
  let finalValue: any;
  if (field.valueField) {
    if (field.valueField instanceof Function) {
      finalValue = field.valueField(value);
    } else if (typeof field.valueField === 'string') {
      finalValue = getValueByPath(field.valueField, value);
    }
  } else {
    finalValue = value;
  }
  return finalValue;
};

/**
 * 根据 string的值来获取对应的数据
 * @param path  对应的path，如 a.b.c等
 * @param store 原始的数据
 */
const getValueByPath = (path: string, store: any) => {
  let finalValue = store;
  let pathArr = path.split('.');
  pathArr.forEach(key => {
    // 这里是需要去过滤的，因为如果前面的已经是undefined了，再取值，有可能会报错的
    finalValue = finalValue ? finalValue[key] : null;
  });
  return finalValue;
};

/**
 * 执行具体的校验
 * @param field: FieldConfig  校验的配置项
 * @param finalValue: 最终值
 * @param finalStore: 当前校验的表单项的值，即 validateStateLocator
 * @return Validation
 */
const execValidate = (field: FieldConfig, value: any, store: any, extraStore: any): Validation => {
  let result: Validation = initValidator;
  // 如果有设置前置条件，则校验前置条件符合后，才进行真正的rules的校验，如果不通过，会初始化Validator
  if (field.condition && !field.condition(value, store)) {
    result = {
      status: ValidatorStatusEnum.Init,
      message: ''
    };
  } else {
    for (let rule of field.rules) {
      let ruleType = typeof rule === 'string' ? rule : rule.type;
      result = validateMethod({
        store,
        value,
        label: field.label,
        rule: rule as Rule,
        modelType: field.modelType,
        extraStore
      })[ruleType]();
      // 自己自定义的rules当中，只要有一个不符合，就不需要继续向下检查
      if (result.status !== ValidatorStatusEnum.Success) break;
    }
  }
  return result;
};

/**
 * 构造表单校验全局实例
 * @param userDefinedSchema: ValidateSchema<T>  表单校验的完整配置
 * @param validateStateLocator: (state: any) => any 需要进行校验的表单数据
 * @param validatorStateLocator: (state: any) => any 表单校验结果的存储store
 * @param extraValidateStateLocatorPath: string[] 适用于依赖外部store的情况，传入store的位置，从RootState
 * @return;
 */
export function createValidatorActions(options: {
  userDefinedSchema: ValidateSchema;
  validateStateLocator: (state: any) => any;
  validatorStateLocation: (state: any) => any;
  extraValidateStateLocatorPath?: string[];
}): ValidationIns {
  let { userDefinedSchema, validateStateLocator, validatorStateLocation, extraValidateStateLocatorPath = [] } = options;

  /**
   * 触发校验的action
   * @param vKey: string 非必传  传入只校验对应vKey的值
   * @param callback: (validatorResult: ValidatorModel) => void 提供回调返回所有校验结果
   */
  const validate = (vKey?: string | string[], callback?: (validateResult: ValidatorModel) => void) => {
    return (dispatch: Redux.Dispatch, getState: () => any) => {
      let store = getState();

      // 判断是否传入了除校验数据项内的额外数据，如workloadEdit 外层的serviceEdit
      let extraStore = {};
      if (extraValidateStateLocatorPath.length) {
        extraValidateStateLocatorPath.forEach((path, index) => {
          extraStore[index] = getValueByPath(path, store);
        });
      }

      let finalStore = validateStateLocator(store);
      let validatorResult: ValidatorModel;
      if (vKey) {
        validatorResult = cloneDeep(validatorStateLocation(store));
      } else {
        validatorResult = {};
      }

      let fields: FieldConfig[] = userDefinedSchema.fields;

      /**
       * 根据传入的vKey，获取用户自定义的校验规则配置中，符合规则的校验规则配置项，并且会将用户传入的vKey 和 配置的vKey进行结合
       * 如 传入的为 containers[0]，符合规则的校验规则有 containers[].name、containers[].labels[]
       * 则返回会是 [containers[0].name、containers[0].labels[]]
       * @returns 获得初始需要进行校验的配置项，结构为
       * {
       *  key: containers[0].name,
       *  value: ''
       * }
       */
      let schemaVkeyArr = getValidateConfig(fields, vKey);

      /**
       * 循环去处理 schemaVkeyArr当中的数据，直到所有vKey都没有 通配符为止
       */
      for (let i = 0; i < schemaVkeyArr.length; i++) {
        let schemaVkeyConfig = schemaVkeyArr[i];
        let { schemaVkey } = schemaVkeyConfig;
        // 判断当前schemaVkey是否包含通配符 []
        let isContainBracketOnly = schemaVkey.includes('[]');
        if (isContainBracketOnly) {
          let expandSchemaVkeyArr = expandSchemaVkey(schemaVkeyConfig, finalStore);
          expandSchemaVkeyArr.length && schemaVkeyArr.push(...expandSchemaVkeyArr);
        } else {
          schemaVkeyArr[i] = {
            schemaVkey,
            value: getStoreValue(schemaVkey, finalStore)
          };
        }
      }

      /**
       * 过滤掉schemaVkeyArr当中包含 通配符 [] 的项
       * @returns 获得最终需要执行校验的vKeys的
       */
      schemaVkeyArr = schemaVkeyArr.filter(schemaVkeyConfig => !schemaVkeyConfig.schemaVkey.includes('[]'));

      /**
       * 进行校验规则的校验
       */
      schemaVkeyArr.forEach(schemaVkeyConfig => {
        let { schemaVkey, value } = schemaVkeyConfig;
        let field: FieldConfig = fields.find(field => field.vKey === schemaVkey.replace(formatVkeyRegExp, '[]'));
        // 如果能找到对应的校验规则配置项，说明存在该项的校验配置，需要执行校验
        if (field) {
          // valueField 是在取到当前指定vKey的值之后，当前值还是一个复杂类型，如 FFReduxModel，这时候需要指定valueField来判断如何取值
          let finalValue: any = getValueByValueField(field, value);

          // 执行校验的校验结果
          let result = execValidate(field, finalValue, finalStore, extraStore);

          // 将结果储存在 validatorResult当中，用以触发action 和 传入到callback当中
          validatorResult[schemaVkey] = result;

          /**
           * @pre 只有传入了vKey，才需要进行递归校验，因为不传vKey，默认都是校验所有选项，不需要进行额外的校验，不进行递归
           * @pre 传入的vKey当中，必须不包含 主动watch的配置项，才需要触发 由于 field 的校验，而触发的依赖校验
           * A --(watch)--> B，当B变化的时候，A也需要校验
           */
          if (vKey) {
            let finalVKeys = vKey instanceof Array ? vKey : vKey ? [vKey] : [];
            let watchFields: string[] = fields
              .filter(watchField => {
                return (
                  !finalVKeys.some(final => watchField.vKey.includes(final.replace(formatVkeyRegExp, '[]'))) &&
                  watchField.watchKey === field.vKey
                );
              })
              .map(item => item.vKey);
            watchFields.length && validate(watchFields);
          }
        }
      });

      /**
       * 触发对应的Action，由Reducer修改对应的store的值，如果当前校验项是一个数组，默认的store是 {}，保存扁平化 { specificKey: Validation }
       * 这里的 specificKey，能够跑到这里来，说明vKey都是很具体的了
       */
      dispatch({
        type: getValidatorActionType(userDefinedSchema.formKey),
        payload: validatorResult
      });

      /**
       * 如果使用方传入了callback，则将validatorResult暴露出去
       */
      callback && callback(validatorResult);
    };
  };

  return {
    validate
  };
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
  let { store, value, label, rule, modelType = ModelTypeEnum.Normal, extraStore } = options;
  let { errorTip, limit, customFunc } = rule;

  const validateRequired = (): Validation => {
    let result = false;
    if (modelType === ModelTypeEnum.FFRedux) {
      result = !!value.selection || value.selections.length > 0;
    } else {
      result = !!value;
    }

    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}不能为空', { label })
    };
  };

  const validateMinLength = (): Validation => {
    let minLength = limit ? +limit : 0;
    let valueLength = value ? value.length : 0;
    let result = valueLength >= minLength;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}至少包含{{count}}个字符', { label, count: minLength })
    };
  };

  const validateMaxLength = (): Validation => {
    let maxLength = limit ? +limit : Number.MAX_SAFE_INTEGER;
    let valueLength = value ? value.length : 0;
    let result = valueLength <= maxLength;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}不能超过{{count}}个字符', { label, count: maxLength })
    };
  };

  const validateMinValue = (): Validation => {
    let minValue = limit ? +limit : 0;
    let result = !Number.isNaN(value) && value >= minValue;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}最小值为{{minValue}}', { label, minValue })
    };
  };

  const validateMaxValue = (): Validation => {
    let maxValue = limit ? +limit : Number.MAX_SAFE_INTEGER;
    let result = !Number.isNaN(value) && value <= maxValue;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}最大值为{{maxValue}}', { label, maxValue })
    };
  };

  const validateRegExp = (): Validation => {
    let result = limit ? limit.test(value) : false;
    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}格式不正确', { label })
    };
  };

  const validateMinSelect = (): Validation => {
    let minValue = limit ? +limit : 0;
    let result = false;
    if (modelType === ModelTypeEnum.FFRedux) {
      result = value.selections.length >= minValue;
    } else {
      let isArray = value instanceof Array;
      result = isArray ? value.length >= minValue : false;
    }

    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}至少选择{{count}}项', { label, count: minValue })
    };
  };

  const validateMaxSelect = (): Validation => {
    let maxValue = limit ? +limit : Number.MAX_SAFE_INTEGER;
    let result = false;
    if (modelType === ModelTypeEnum.FFRedux) {
      result = value.selections.length <= maxValue;
    } else {
      let isArray = value instanceof Array;
      result = isArray ? value.length <= maxValue : false;
    }

    return {
      status: result ? ValidatorStatusEnum.Success : ValidatorStatusEnum.Failed,
      message: result ? '' : errorTip ? errorTip : t('{{label}}至多选择{{count}}项', { label, count: maxValue })
    };
  };

  const validateCustom = (): Validation => {
    return customFunc ? customFunc(value, store, extraStore) : { status: ValidatorStatusEnum.Success, message: '' };
  };

  return {
    [RuleTypeEnum.isRequire]: validateRequired,
    [RuleTypeEnum.maxValue]: validateMaxValue,
    [RuleTypeEnum.minValue]: validateMinValue,
    [RuleTypeEnum.minLength]: validateMinLength,
    [RuleTypeEnum.maxLength]: validateMaxLength,
    [RuleTypeEnum.regExp]: validateRegExp,
    [RuleTypeEnum.custom]: validateCustom,
    [RuleTypeEnum.minSelect]: validateMinSelect,
    [RuleTypeEnum.maxSelect]: validateMaxSelect
  };
};
