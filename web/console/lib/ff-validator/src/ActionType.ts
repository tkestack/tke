/**
 * 获取ActionType
 * @param formKey: string 表单的名称
 * @param actionName: string  触发的Action的名称
 * @return uppercase action
 */
export const getValidatorActionType = (formKey: string) => {
  return `${formKey}_VALIDATION`.toUpperCase();
};
