import { t } from '@tencent/tea-app/lib/i18n';

export const canNotOperateCluster = {
  clusterStatus: ['Initializing', 'Running', 'Failed', 'Upgrading', 'Terminating']
};

/** 集群状态 */
export const clusterStatus = {
  Initializing: {
    text: t('创建中'),
    classname: 'text-restart'
  },
  Running: {
    text: t('运行中'),
    classname: 'text-success'
  },
  Terminating: {
    text: t('删除中'),
    classname: 'text-restart'
  },
  Scaling: {
    text: t('规模调整中'),
    classname: 'text-restart'
  },
  Upgrading: {
    text: t('升级中'),
    classname: 'text-restart'
  },
  Failed: {
    text: t('异常'),
    classname: 'text-danger'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};
