import { Validation, EnvItem, HealthCheck, MountItem } from './';
import { Identifiable } from '@tencent/qcloud-lib';

export interface ContainerItem extends Identifiable {
  /**状态，用于标识该条数据的编辑状态 */
  status: string;

  /**容器名称 */
  name?: string;
  v_name?: Validation;

  /**镜像 */
  repository?: string;
  v_repository?: Validation;

  /**地域 */
  regionId?: string | number;

  /**版本 */
  tag?: string;
  v_tag?: Validation;

  /**是否打开高级设置 */
  isOpenAdvancedSetting: boolean;

  /**CPU Request*/
  cpuRequest?: number;
  v_cpuRequest?: Validation;

  /**CPU Limit*/
  cpuLimit?: number;
  v_cpuLimit?: Validation;

  /**内存Request*/
  memRequest?: number;
  v_memRequest?: Validation;

  /**内存Limit*/
  memLimit?: number;
  v_memLimit?: Validation;

  /**运行命令 */
  cmd?: string;
  v_cmd?: Validation;

  /**运行参数 */
  arg: string;
  v_arg: Validation;

  /**环境变量 */
  envs?: EnvItem[];

  /**挂载点 */
  mounts?: MountItem[];

  /**健康检查 */
  healthCheck?: HealthCheck;

  /**工作目录 */
  workingDir?: string;
  v_workingDir?: Validation;

  /**是否是特权级容器 */
  privileged?: boolean;

  gpu?: number;
}

export interface CpuLimitItem extends Identifiable {
  /**类型 */
  type?: string;

  /**限制值 */
  value?: string;
  v_value?: Validation;
}
