/** 获取当前控制台modules的 域名映射表 */
export enum ConsoleModuleEnum {
  /** tke-apiserver 版本 */
  PLATFORM = 'platform',

  /** 业务的版本详情 */
  Business = 'business',

  /** 通知模块 */
  Notify = 'notify',

  /** 告警模块 */
  Monitor = 'monitor',

  /** 镜像仓库 */
  Registry = 'registry',

  /** 日志模块 */
  LogAgent = 'logagent',

  /** 认证模块 */
  Auth = 'auth',

  /** 审计模块 */
  Audit = 'audit',

  /** Helm应用模块 */
  Application = 'application'
}
