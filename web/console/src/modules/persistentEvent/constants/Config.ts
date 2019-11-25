import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** ========================= start FFRedux的相关配置 ======================== */
export const FFReduxActionName = {
  REGION: 'region',
  CLUSTER: 'cluster'
};

/** 事件轮询的列表 */
export const PollEventName = {
  peList: 'pollPeList'
};

/** 是否需要轮询事件持久列表 */
export const isNeedPollPE = ['initializing', 'reinitializing', 'checking', 'failed'];

/** 事件持久化状态 */
export const peStatus = {
  initializing: {
    text: t('初始化'),
    classname: 'text-restart'
  },
  reinitializing: {
    text: t('重新初始化'),
    classname: 'text-restart'
  },
  running: {
    text: t('已开启'),
    classname: 'text-success'
  },
  failed: {
    text: t('失败'),
    classname: 'text-danger'
  },
  checking: {
    text: t('检查中'),
    classname: 'text-restart'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};
