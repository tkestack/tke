/**
 * 获取components、info等基础组件的版本，不会更改
 */
export const basicServerVersion = {
  basicUrl: 'apis',
  group: 'gateway.tkestack.io',
  version: 'v1'
};

/**
 * 这里是后台的api Server 的Version的版本
 * 对于不是集群内的资源，而是CRD的，其版本由tke-apiserver的版本决定
 */
export const apiServerVersion = {
  basicUrl: 'apis',
  group: 'platform.tkestack.io',
  version: 'v1'
};

/**
 * 业务的Server版本
 * 根据tke的版本进行变化
 */
export const businessServerVersion = {
  basicUrl: 'apis',
  group: 'business.tkestack.io',
  version: 'v1'
};

/**
 * 通知、告警的Server版本
 * 根据tke的版本进行变化
 */
export const notifyServerVersion = {
  basicUrl: 'apis',
  group: 'notify.tkestack.io',
  version: 'v1'
};

/**
 * 认证模块
 * 用户管理、策略管理
 */
export const authServerVersion = {
  basicUrl: 'apis',
  group: 'auth.tkestack.io',
  version: 'v1'
};

/**
 * 告警模块
 */
export const monitorServerVersion = {
  basicUrl: 'apis',
  group: 'monitor.tkestack.io',
  version: 'v1'
};

/**
 * 审计模块
 */
export const auditServerVersion = {
  basicUrl: 'apis',
  group: 'audit.tkestack.io',
  version: 'v1'
};

/**
 * 新日志模块
 */
export const logAgentServiceVersion = {
  basicUrl: 'apis',
  group: 'logagent.tkestack.io',
  version: 'v1'
};
 
/**
 * 仓库模块
 */
export const registryServerVersion = {
  basicUrl: 'apis',
  group: 'registry.tkestack.io',
  version: 'v1'
};

/**
 * 应用模块
 */
export const applicationServerVersion = {
  basicUrl: 'apis',
  group: 'application.tkestack.io',
  version: 'v1'
};
