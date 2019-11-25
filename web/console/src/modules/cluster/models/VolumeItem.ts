import { Identifiable } from '@tencent/qcloud-lib';
import { Validation } from '../../common/models';

export interface VolumeItem extends Identifiable {
  /** 卷类型 */
  volumeType: string;
  v_volumeType: Validation;

  /** 卷名称 */
  name: string;
  v_name: Validation;

  /** 源路径 */
  hostPathType: string;
  hostPath: string;
  v_hostPath: Validation;

  /** nfs路径 */
  nfsPath: string;
  v_nfsPath: Validation;

  /** 配置项 */
  configKey: ConfigItems[];
  configName: string;

  /** secret的相关配置 */
  secretKey: ConfigItems[];
  secretName: string;

  /** pvc的选择 */
  pvcSelection: string;
  v_pvcSelection: Validation;

  /** 新创建的pvc */
  newPvcName: string;
  pvcEditInfo: PvcEditInfo;

  /** 当前数据卷是否被挂载 */
  isMounted: boolean;
}

export interface PvcEditInfo {
  /** accessMode */
  accessMode: string;

  /** storageClassName */
  storageClassName: string;

  /** storage */
  storage: string;
}

export interface ConfigItems extends Identifiable {
  /** config的Key */
  configKey?: string;
  v_configKey?: Validation;

  /** 配置的子路径 */
  path?: string;
  v_path?: Validation;

  /** 当前mode */
  mode?: string;
  v_mode?: Validation;
}
