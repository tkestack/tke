import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { initValidator } from '../../common/models';
export const initRegionInfo = {
  name: t('广州'),
  value: 1,
  area: t('华南地区')
};

/**
 * 创建集群步骤
 */
export const initSteps = [
  {
    id: 1,
    label: t('集群信息')
  },
  {
    id: 2,
    label: t('选择机型')
  },
  {
    id: 3,
    label: t('云服务器配置')
  }
];

/**
 * 修改名称
 */
export const initName = {
  clusterName: '',
  v_clusterName: initValidator
};

/**
 * 修改描述
 */
export const initDescription = {
  clusterDescription: '',
  v_clusterDescription: initValidator
};

export const initAsgCreation = {
  /**弹性伸缩组名称 */
  autoScalingGroupName: '',
  v_autoScalingGroupName: initValidator,
  v_sameGroupName: initValidator,

  /**参照实例Id */
  instanceId: '',
  v_instanceId: initValidator,

  /**弹性伸缩组 最小实例 */
  minSize: '',
  v_minSize: initValidator,

  /**弹性伸缩组 最大实例 */
  maxSize: '',
  v_maxSize: initValidator,

  /**登录方式 */
  loginMethod: 'password',

  /**实例密码 */
  password: '',
  v_password: initValidator,

  /**密钥Id */
  keyId: '',
  v_keyId: initValidator,

  needSecurityAgent: true,

  needMonitorAgent: true,

  isOpenAdvanced: false,

  userScript: '',

  v_userScript: initValidator,

  isUserScriptBase64: false,

  isUnSchedule: false,

  /**弹性伸缩组 标签 */
  label: []
};

export const initAsgGlobalSetting = {
  expander: 'random',

  /** 每次最大缩容的完全空闲的节点数量 */
  maxEmptyBulkDelete: '10',

  /** unready节点所占百分比，超过此比例后CA停止工作 */
  maxTotalUnreadyPercentage: '0',

  /** unready节点总数超过此数值时，CA停止工作 */
  okTotalUnreadyCount: '0',

  /** 扩容多久后开始考虑缩容条件, 默认值10分钟 */
  scaleDownDelay: '10',

  /** 是否开启了缩容 0 | 1 */
  scaleDownEnabled: '0',

  /** 节点连续多久空闲时被缩容，默认值10分钟 */
  scaleDownUnneededTime: '10',

  /** 集群中unready状态的节点持续多久后被缩容，默认值20分钟 */
  scaleDownUnreadyTime: '20',

  /** request占用比  小于 此值时， 尝试缩容 */
  scaleDownUtilizationThreshold: '0',

  /** 是否不缩容含有本地存储Pod的节点 0 | 1 */
  skipNodesWithLocalStorage: '0',

  /** 是否不缩容含有 kube-system 下的pod节点, 0 | 1 */
  skipNodesWithSystemPods: '0',

  unregisteredNodeRemovalTime: '0'
};

export const initEditAsgOption = Object.assign({}, initAsgGlobalSetting, {
  v_maxEmptyBulkDelete: initValidator,
  v_scaleDownUnneededTime: initValidator,
  v_scaleDownDelay: initValidator
});

/**第三方仓库初始化 */
export const initThirdHub = {
  clusterId: '',

  regionId: '',

  domain: '',

  v_domain: initValidator,

  username: '',

  v_username: initValidator,

  password: '',

  v_password: initValidator,

  isSelectAllCluster: true,

  namespaceSelection: [],

  v_namespaceSelection: initValidator
};

export const initProjectSelect = {
  selectProject: -1
};
