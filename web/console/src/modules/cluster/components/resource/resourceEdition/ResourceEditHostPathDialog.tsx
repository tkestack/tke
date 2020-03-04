import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Modal } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField } from '../../../../common/components';
import { FormLayout } from '../../../../common/layouts';
import { allActions } from '../../../actions';
import { validatorActions } from '../../../actions/validatorActions';
import { VolumeItem } from '../../../models';
import { RootProps } from '../../ClusterApp';

/** hostPath检查类型的列表 */
const HostPathTypeList = {
  NoChecks: {
    value: 'NoChecks',
    tip: t('不对指定路径做任何检查')
  },
  DirectoryOrCreate: {
    value: 'DirectoryOrCreate',
    tip: t('如果给定路径中不存在任何内容，则将根据需要创建一个空目录，且权限设置为0755')
  },
  Directory: {
    value: 'Directory',
    tip: t('目录必须存在于指定路径中')
  },
  FileOrCreate: {
    value: 'FileOrCreate',
    tip: t('如果给定路径中不存在任何内容，则会根据需要创建一个空文件，且权限设置为0644')
  },
  File: {
    value: 'File',
    tip: t('文件必须存在于指定路径中')
  },
  Socket: {
    value: 'Socket',
    tip: t('UNIX socket 必须存在于指定路径中')
  },
  CharDevice: {
    value: 'CharDevice',
    tip: t('字符设备必须存在于指定路径中')
  },
  BlockDevice: {
    value: 'BlockDevice',
    tip: t('块设备必须存在于指定路径中')
  }
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceEditHostPathDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot } = this.props,
      { isShowHostPathDialog, currentEditingVolumeId, volumes } = subRoot.workloadEdit;

    // 如果不需要展示HostPath的弹窗配置
    if (!isShowHostPathDialog) {
      return <noscript />;
    }

    let currentVolume: VolumeItem = volumes.find(item => item.id === currentEditingVolumeId);

    let { hostPath, v_hostPath, hostPathType } = currentVolume;

    const cancel = () => {
      actions.editWorkload.toggleHostPathDialog();
    };

    const perform = () => {
      actions.validate.workload.validateHostPath(currentVolume.hostPath, currentEditingVolumeId);
      // 确定需要确认当前的hostPath是否正确
      if (validatorActions.workload._validateHostPath(hostPath).status === 1) {
        actions.editWorkload.toggleHostPathDialog();
      }
    };

    let hostPathTypeOptions = Object.keys(HostPathTypeList).map((item, index) => (
      <option key={index} value={item}>
        {item}
      </option>
    ));

    return (
      <Modal visible={true} caption={t('设置主机路径')} onClose={cancel} disableEscape={true} size={650}>
        <Modal.Body>
          <FormLayout>
            <div className="param-box server-update add">
              <ul className="form-list jiqun fixed-layout">
                <FormItem label={t('主机路径')}>
                  <InputField
                    type="text"
                    placeholder={t('如: /data/dev')}
                    tipMode="popup"
                    validator={v_hostPath}
                    value={hostPath}
                    onChange={data => {
                      actions.editWorkload.updateVolume({ hostPath: data }, currentEditingVolumeId);
                    }}
                    onBlur={data => {
                      actions.validator.workload.validateHostPath(data, currentEditingVolumeId);
                    }}
                  />
                </FormItem>
                <FormItem label={t('检查类型')}>
                  <div className="form-unit" style={{ marginTop: '5px' }}>
                    <select
                      className="tc-15-select m"
                      value={hostPathType}
                      onChange={e =>
                        actions.editWorkload.updateVolume({ hostPathType: e.target.value }, currentEditingVolumeId)
                      }
                    >
                      {hostPathTypeOptions}
                    </select>
                    <p className="text-label">{HostPathTypeList[hostPathType].tip}</p>
                  </div>
                </FormItem>
              </ul>
            </div>
          </FormLayout>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={perform}>
            {t('确认')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
