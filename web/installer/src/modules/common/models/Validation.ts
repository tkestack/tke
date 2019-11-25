export interface Validation {
  /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
  status?: number;

  /**结果描述 */
  message?: string;
}

export const initValidation = {
  /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
  status: 0,

  /**结果描述 */
  message: ''
};
