/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem } from '../../../../common/components';
import { isEmpty, uniq } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { VolumeMountModeList } from '../../../constants/Config';
import { VolumeItem } from '../../../models';
import { RootProps } from '../../ClusterApp';

const filterValidMountVolumes = (volumes: VolumeItem[]) => {
  let validVolumes = volumes.length
    ? volumes.filter(v => {
        return (
          (v.volumeType === 'emptyDir' && v.name) ||
          (v.volumeType === 'hostPath' && v.name && v.hostPath) ||
          (v.volumeType === 'nfsDisk' && v.name && v.nfsPath) ||
          (v.volumeType === 'configMap' && v.name && v.configName) ||
          (v.volumeType === 'secret' && v.name && v.secretName) ||
          (v.volumeType === 'pvc' && v.name && v.pvcSelection)
        );
      })
    : [];
  return validVolumes;
};

interface ContainerMountItemProps extends RootProps {
  cKey: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerMountItem extends React.Component<ContainerMountItemProps, {}> {
  render() {
    let isAvaiable = this._hasAvailableVolume();

    return isAvaiable ? (
      <FormItem label={t('挂载点')} tips={t('设置数据卷挂载到容器中的路径')} isPureText={true}>
        <div className="form-unit is-success">{this._renderMountList()}</div>
      </FormItem>
    ) : (
      <noscript />
    );
  }

  /** 校验是否有可挂在的数据盘 */
  private _hasAvailableVolume() {
    let { subRoot } = this.props,
      { workloadEdit } = subRoot,
      { volumes } = workloadEdit;

    let filters = uniq(filterValidMountVolumes(volumes), 'name');

    return !isEmpty(filters);
  }

  /** 渲染挂载的配置 */
  private _renderMountList() {
    let { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers, volumes } = workloadEdit;

    let container = containers.find(c => c.id === cKey),
      canAdd = isEmpty(container.mounts.filter(x => !x.volume || !x.mountPath)),
      self = this;

    /** 渲染数据卷列表 */
    let filters = uniq(filterValidMountVolumes(volumes), 'name');
    let volumeOptions = filters.map((f, index) => (
      <option key={index} value={f.name}>
        {f.name}
      </option>
    ));
    volumeOptions.unshift(
      <option key={-1} value="">
        {t('请选择')}
      </option>
    );

    /** 渲染权限列表 */
    let modeOptions = VolumeMountModeList.map((v, index) => (
      <option key={index} value={v.value}>
        {v.label}
      </option>
    ));

    let mountList = container.mounts.map((mount, index) => {
      return (
        <div
          key={index}
          className={classnames('code-list', {
            'is-error': mount.v_mountPath.status === 2 || mount.v_volume.status === 2
          })}
        >
          <Bubble placeholder="bottom" content={mount.v_volume.status === 2 ? mount.v_volume.message : null}>
            <select
              className="tc-15-select m"
              style={{ minWidth: '90px', marginRight: '5px' }}
              value={mount.volume}
              onChange={e => self._handleVolumetMount(e.target.value, cKey, mount.id + '')}
            >
              {volumeOptions}
            </select>
          </Bubble>
          <Bubble placement="bottom" content={mount.v_mountPath.status === 2 ? mount.v_mountPath.message : null}>
            <input
              type="text"
              placeholder={t('目标路径，如:/mnt')}
              className="tc-15-input-text m"
              style={{ width: '128px', marginRight: '5px' }}
              value={mount.mountPath}
              onChange={e => actions.editWorkload.updateMount({ mountPath: e.target.value }, cKey, mount.id + '')}
              onBlur={e => actions.validate.workload.validateVolumeMountPath(mount.mountPath, cKey, mount.id + '')}
            />
          </Bubble>
          <Bubble
            placement="bottom"
            content={
              mount.mountSubPath === ''
                ? t(
                    '挂载子路径，为空时全覆盖目标路径，如目标路径为/var/www/html，挂载子路径为html时，仅将数据卷覆盖html目录'
                  )
                : null
            }
          >
            <input
              type="text"
              placeholder={t('挂载子路径')}
              className="tc-15-input-text m"
              style={{ width: '128px', marginRight: '5px' }}
              value={mount.mountSubPath}
              onInput={e => actions.editWorkload.updateMount({ mountSubPath: e.target.value }, cKey, mount.id + '')}
            />
          </Bubble>
          <select
            className="tc-15-select m"
            style={{ minWidth: '90px', marginRight: '5px' }}
            value={mount.mode}
            onChange={e => actions.editWorkload.updateMount({ mode: e.target.value }, cKey, mount.id + '')}
          >
            {modeOptions}
          </select>
          <span className="inline-help-text">
            <a href="javascript:;" onClick={() => actions.editWorkload.deleteMount(cKey, mount.id + '')}>
              <i className="icon-cancel-icon" />
            </a>
          </span>
        </div>
      );
    });

    return (
      <div>
        {mountList}
        {canAdd ? (
          <a href="javascript:;" className="more-links-btn" onClick={() => actions.editWorkload.addMount(cKey)}>
            {t('添加挂载点')}
          </a>
        ) : (
          <Bubble placement="left" content={t('请先完成待编辑项')}>
            <a href="javascript:;" className="more-links-btn disabled">
              {t('添加挂载点')}
            </a>
          </Bubble>
        )}
      </div>
    );
  }

  /** 选择挂在的数据卷 */
  private _handleVolumetMount(volumeName: string, cKey: string, mId: string) {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot;

    // 需要更新更新一下数据卷的挂载状态
    let finder = workloadEdit.volumes.find(v => v.name === volumeName);
    if (finder) {
      actions.editWorkload.updateVolume({ isMounted: true }, finder.id + '');
    }
    // 更新容器的挂载项的状态
    actions.editWorkload.updateMount({ volume: volumeName }, cKey, mId);
    actions.validate.workload.validateVolumeMount(volumeName, cKey, mId);
  }
}
