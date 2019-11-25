import { Identifiable } from '@tencent/qcloud-lib';

interface Base extends Identifiable {
  /**配置文件Id */
  configId?: string;

  /**配置文件名称 */
  name?: string;

  /**创建时间 */
  createdAt?: string;
}

export interface Config extends Base {
  /**版本数量 */
  totalCount?: number;

  /**修改时间 */
  updatedAt?: string;
}

export interface ConfigFilter {
  /**搜索字段 */
  search?: string;
}

export interface Version extends Base {
  /**版本名称 */
  version?: string;

  /**配置数据 */
  data?: string;

  /**描述 */
  description?: string;
}

export interface VersionFilter {
  /**配置文件Id */
  configId?: string;

  /**地域 */
  regionId?: number | string;
}

export interface Variable extends Identifiable {
  /**变量名称 */
  key?: string;

  /**变量值 */
  value?: string;

  /**变量类型 */
  type?: string;

  /**变量名是否是规范化变量 */
  isLegal?: boolean;
}

export interface VariableFilter {
  /**配置文件Id */
  configId?: string;

  /**版本名称 */
  version?: string;
}
