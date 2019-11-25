export enum DialogNameEnum {
  /** 集群凭证 弹窗 */
  kuberctlDialog = 'kuberctlDialog',

  /** 集群状态详情 弹窗 */
  clusterStatusDialog = 'clusterStatusDialog',

  /** 节点状态详情 弹窗 */

  computerStatusDialog = 'computerStatusDialog'
}

export interface DialogState {
  /** 是否展示 集群凭证 弹窗 */
  [DialogNameEnum.kuberctlDialog]: boolean;

  /** 是否展示集群状态的弹窗 */
  [DialogNameEnum.clusterStatusDialog]: boolean;

  /** 是否展示节点状态的弹窗 */
  [DialogNameEnum.computerStatusDialog]: boolean;
}
