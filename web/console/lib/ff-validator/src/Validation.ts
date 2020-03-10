export interface Validation {
  /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
  status?: number;

  /**结果描述 */
  message?: string | React.ReactNode;

  /**
   * 返回的校验列表
   * 目前仅 CIDR 有使用
   */
  list?: any[];
}

export const initValidator = {
  /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
  status: 0,

  /**结果描述 */
  message: ''
};
