import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, ExternalLink, Icon, Text } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, LinkButton } from '../../../../common/components';
import { allActions } from '../../../actions';
import { validateWorkloadActions } from '../../../actions/validateWorkloadActions';
import { VolumeTypeList } from '../../../constants/Config';
import { VolumeItem } from '../../../models';
import { RootProps } from '../../ClusterApp';

const AccessModeNameMap = {
  ReadWriteOnce: t('单机读写')
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceVolumePanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { volumes } = workloadEdit;

    let canAdd = validateWorkloadActions._canAddVolume(volumes);

    return (
      <FormItem className="vm data-mod" label={t('数据卷（选填）')}>
        <div className="form-unit">
          {volumes.length ? (
            <div className="tc-15-table-panel" style={{ maxWidth: '730px', marginTop: '2px' }}>
              <div className="tc-15-table-fixed-body">
                <table className="tc-15-table-box tc-15-table-rowhover">
                  <colgroup>
                    <col />
                    <col />
                    <col />
                    <col style={{ width: '13%' }} />
                  </colgroup>
                  <tbody>{this._renderVolumeContent()}</tbody>
                </table>
              </div>
            </div>
          ) : (
            <noscript />
          )}
        </div>
        <LinkButton disabled={!canAdd} tip="" errorTip={t('请先完成待编辑项')} onClick={actions.editWorkload.addVolume}>
          {t('添加数据卷')}
        </LinkButton>
        <p className="text-label">
          <span style={{ verticalAlign: '-1px' }}>
            {t(
              '为容器提供存储，目前支持临时路径、主机路径、云硬盘数据卷、文件存储NFS、配置文件、PVC，还需挂载到容器的指定路径中。'
            )}
          </span>
        </p>
      </FormItem>
    );
  }

  /** 切换数据卷类型的时候，进行一些数据的更新 */
  private _handleVolumeTypeSelect(volumeType: string, vId: string) {
    let { actions } = this.props;
    // 更新每一个volume的 数据卷类型
    actions.editWorkload.updateVolume({ volumeType }, vId);
  }

  /** 渲染表格当中的具体内容 */
  private _renderVolumeContent() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { volumes, containers, workloadType } = workloadEdit;

    /** 渲染数据卷类型列表，只有Statefulset有使用新的pvc */
    let finalVolumeTypeList = VolumeTypeList;
    let volumeTypeOptions = finalVolumeTypeList.map((v, vIndex) => (
      <option key={vIndex} value={v.value}>
        {v.label}
      </option>
    ));

    // 获取正确的上下文context
    let self = this;

    return volumes.map(volume => {
      // 这里需要去判断，如果数据卷已经被挂在到containers上面了
      let canDelete = true;
      containers.forEach(c => {
        if (volume.name !== '') {
          canDelete = canDelete && !c.mounts.find(m => m.volume === volume.name);
        }
      });

      // 是否需要展示未挂在提示
      let isNeedAlertError = !volume.isMounted;

      let volumeId = volume.id + '';

      return (
        <tr
          key={volume.id}
          className={classnames('tr-hover run-docker-box')}
          style={isNeedAlertError ? { border: '1px solid red' } : {}}
        >
          <td>
            <div>
              <span className="text-overflow">
                <select
                  className="tc-15-select m"
                  value={volume.volumeType}
                  onChange={e => self._handleVolumeTypeSelect(e.target.value, volumeId)}
                >
                  {volumeTypeOptions}
                </select>
              </span>
            </div>
          </td>
          <td>
            <div className="form-unit">
              {isNeedAlertError && (
                <div className="tc-15-bubble tc-15-bubble-bottom" style={{ marginTop: '-60px' }}>
                  <div className="tc-15-bubble-inner">
                    <p>{t('该数据卷还未挂载，请进行挂载或删除')}</p>
                  </div>
                </div>
              )}
              <InputField
                type="text"
                placeholder={t('名称，如：vol')}
                tipMode="popup"
                validator={volume.v_name}
                value={volume.name}
                onChange={data => {
                  actions.editWorkload.updateVolume({ name: data }, volumeId);
                }}
                onBlur={data => {
                  actions.validate.workload.validateVolumeName(data, volumeId);
                }}
              />
            </div>
          </td>
          <td>
            {volume.volumeType === 'hostPath' && this._renderHostPath(volume)}
            {volume.volumeType === 'nfsDisk' && this._renderNfsDisk(volume)}
            {volume.volumeType === 'configMap' && this._renderConfigMap(volume)}
            {volume.volumeType === 'secret' && this._renderSecret(volume)}
            {volume.volumeType === 'pvc' && this._renderExistedPVC(volume)}
          </td>
          <td>
            <LinkButton
              disabled={!canDelete}
              tip={t('删除')}
              errorTip={t('该数据卷已被挂载，不可删除')}
              onClick={() => actions.editWorkload.deleteVolume(volumeId)}
            >
              <i className="icon-cancel-icon" />
            </LinkButton>
          </td>
        </tr>
      );
    });
  }

  /** hostPath 主机路径类型 */
  private _renderHostPath(volume: VolumeItem) {
    return (
      <div className="form-unit">
        <Bubble
          placement="bottom"
          content={
            volume.hostPath ? (
              <React.Fragment>
                <p>{t('主机路径：') + volume.hostPath}</p>
                <p>{t('检查类型：') + volume.hostPathType}</p>
              </React.Fragment>
            ) : null
          }
        >
          {volume.hostPath ? (
            <div style={{ display: 'inline-block' }}>
              <Text theme="text" verticalAlign="middle">
                {t('主机路径配置')}
              </Text>
              <i className="plaint-icon" style={{ verticalAlign: 'middle', marginRight: '5px' }} />
            </div>
          ) : (
            <Text theme="text" verticalAlign="middle">
              {t('暂未设置主机路径')}
            </Text>
          )}
        </Bubble>

        <a
          href="javascript:;"
          onClick={() => {
            this._handleHostPathConfig(volume.id + '');
          }}
        >
          {volume.hostPath ? t('重新设置') : t('设置主机路径')}
        </a>
      </div>
    );
  }

  /** 设置主机路径的相关操作 */
  private _handleHostPathConfig(volumeId: string) {
    let { actions } = this.props;
    actions.editWorkload.toggleHostPathDialog();
    actions.editWorkload.changeCurrentEditingVolumeId(volumeId);
  }

  /** nfsDisk的展示 */
  private _renderNfsDisk(volume: VolumeItem) {
    let { actions } = this.props;

    return (
      <div className={classnames('form-unit', { 'is-error': volume.v_nfsPath.status === 2 })}>
        <Bubble placement="bottom" content={volume.v_nfsPath.status === 2 ? <p>{volume.v_nfsPath.message}</p> : null}>
          <input
            type="text"
            placeholder={t('NFS路径 如：127.0.0.1:/dir')}
            className="tc-15-input-text m"
            style={{ width: '170px' }}
            value={volume.nfsPath}
            onChange={e => actions.editWorkload.updateVolume({ nfsPath: e.target.value }, volume.id + '')}
            onBlur={e => actions.validate.workload.validateNfsPath(e.target.value, volume.id + '')}
          />
        </Bubble>
        <Bubble content={'请确保节点当中已经安装 nfs-utils包，才可正常使用nfs数据盘'}>
          <Icon className="tea-ml-1n" type="help" />
        </Bubble>
      </div>
    );
  }

  /** configMap的展示 */
  private _renderConfigMap(volume: VolumeItem) {
    return (
      <div className="form-unit">
        <Bubble
          placement="bottom"
          content={
            volume.configKey.length
              ? volume.configKey.map((item, index) => {
                  return <p key={index}>{`${item.configKey}、${item.path}、${item.mode}`}</p>;
                })
              : null
          }
        >
          {volume.configName ? (
            <span className="text" title={volume.configName} style={{ verticalAlign: 'middle', maxWidth: '90px' }}>
              <span className="text" style={{ display: 'block' }}>
                {volume.configName}
              </span>
              <span className="text" style={{ display: 'block' }}>
                {volume.configKey.length ? t('指定部分Key') : t('全部Key')}
              </span>
            </span>
          ) : (
            <span className="text" style={{ verticalAlign: 'middle' }}>
              {t('暂未选择ConfigMap')}
            </span>
          )}
          {volume.configKey.length !== 0 && (
            <i className="plaint-icon" style={{ verticalAlign: 'middle', marginRight: '5px' }} />
          )}
        </Bubble>
        <a
          href="javascript:;"
          onClick={() => {
            this._handleConfigMapOrSecretSelect(volume);
          }}
        >
          {volume.configKey.length ? t('重新选择') : t('选择配置项')}
        </a>
      </div>
    );
  }

  /** configMap 和 secret的选择操作 */
  private _handleConfigMapOrSecretSelect(volume: VolumeItem, isSecret: boolean = false) {
    let { actions } = this.props;
    actions.editWorkload.toggleConfigDialog();
    actions.editWorkload.changeCurrentEditingVolumeId(volume.id + '');

    // 这里需要去判断当前操作的是 configmap还是secret
    let keyLength = isSecret ? volume.secretKey.length : volume.configKey.length;
    actions.editWorkload.config.changeKeyType(keyLength ? 'optional' : 'all');
  }

  /** secret的展示 */
  private _renderSecret(volume: VolumeItem) {
    return (
      <div className="form-unit">
        <Bubble
          placement="bottom"
          content={
            volume.secretKey.length
              ? volume.secretKey.map((item, index) => {
                  return <p key={index}>{`${item.configKey}、${item.path}、${item.mode}`}</p>;
                })
              : null
          }
        >
          {volume.secretName ? (
            <span className="text" title={volume.secretName} style={{ verticalAlign: 'middle', maxWidth: '90px' }}>
              <span className="text" style={{ display: 'block' }}>
                {volume.secretName}
              </span>
              <span className="text" style={{ display: 'block' }}>
                {volume.secretKey.length ? t('指定部分Key') : t('全部Key')}
              </span>
            </span>
          ) : (
            <span className="text" style={{ verticalAlign: 'middle' }}>
              {t('暂未选择Secret')}
            </span>
          )}
          {volume.secretKey.length !== 0 && (
            <i className="plaint-icon" style={{ verticalAlign: 'middle', marginRight: '5px' }} />
          )}
        </Bubble>
        <a
          href="javascript:;"
          onClick={() => {
            this._handleConfigMapOrSecretSelect(volume, true);
          }}
        >
          {volume.secretKey.length ? t('重新选择') : t('选择Secret')}
        </a>
      </div>
    );
  }

  /** pvc的相关选项 */
  private _renderExistedPVC(volume: VolumeItem) {
    let { subRoot } = this.props,
      { pvcList } = subRoot.workloadEdit;

    // 渲染 pvc选择的列表
    let pvcOptions = pvcList.data.recordCount
      ? pvcList.data.records.map((p, index) => {
          return (
            <option key={index} value={p.metadata.name}>
              {p.metadata.name}
            </option>
          );
        })
      : [];
    pvcOptions.unshift(
      <option key={-1} value="">
        {pvcList.data.recordCount ? t('请选择PVC') : t('无可用PVC')}
      </option>
    );

    return (
      <div className={classnames('form-unit', { 'is-error': volume.v_pvcSelection.status === 2 })}>
        <Bubble placement="left" content={volume.v_pvcSelection.status === 2 ? volume.v_pvcSelection.message : null}>
          <select
            className="tc-15-select m"
            style={{ maxWidth: '180px' }}
            value={volume.pvcSelection}
            onChange={e => {
              this._handleExistedPvcSelect(e.target.value, volume.id + '');
            }}
          >
            {pvcOptions}
          </select>
        </Bubble>
      </div>
    );
  }

  /** 选择已有的pvc的相关操作 */
  private _handleExistedPvcSelect(pvcName: string, vId: string) {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { pvcList } = workloadEdit;

    let pvcInfo = pvcList.data.records.find(item => item.metadata.name === pvcName);
    actions.editWorkload.updateVolume({ pvcSelection: pvcName }, vId);
    // 校验当前的pvc选择是否合法
    actions.validate.workload.validateVolumePvc(pvcName, vId);
  }
}
