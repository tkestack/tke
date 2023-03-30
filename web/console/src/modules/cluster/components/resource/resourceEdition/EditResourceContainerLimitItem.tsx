/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';

interface ResourceItemProps extends RootProps {
  cKey: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerLimitItem extends React.Component<ResourceItemProps, {}> {
  render() {
    const {
      subRoot: {
        workloadEdit: { oversoldRatio }
      }
    } = this.props;
    const hasCpuRatio = oversoldRatio.cpu ? true : false;
    const hasMemoryRatio = oversoldRatio.memory ? true : false;
    return (
      <FormItem label={t('CPU/内存限制')}>
        <div className="form-input limit resource" style={{ paddingBottom: '0' }}>
          <div className="tc-input-group-wrap s">
            <p className="top-tip">
              <span>{t('CPU限制')}</span>
              <span>{t('内存限制')}</span>
            </p>
            {this._renderCpuItem(hasCpuRatio)}

            {this._renderMemItem(hasMemoryRatio)}
            <p className="form-input-help">
              <Trans>
                Request用于预分配资源,当集群中的节点没有request所要求的资源数量时,容器会创建失败。
                <br />
                Limit用于设置容器使用资源的最大上限,避免异常情况下节点资源消耗过多。
                <br />
                设置超售比后，此处只能设置CPU/内存Limit值，request值通过Limit和超售比计算得出。
              </Trans>
              <br />
            </p>
          </div>
        </div>
      </FormItem>
    );
  }

  /** cpu的相关选项 */
  private _renderCpuItem(hasCpuRatio?: boolean) {
    const { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    const container = containers.find(c => c.id === cKey);

    return container.cpuLimit.map((cpu, index) => {
      const partition = cpu.type === 'request' ? '-' : t('核');
      const cpuGroupClassName = cpu.type === 'request' || hasCpuRatio ? 'tc-input-group' : 'tc-input-group mr15';
      if (hasCpuRatio && cpu.type === 'request') {
        return <noscript />;
      } else {
        return (
          <div
            key={index}
            className={classnames(cpuGroupClassName, { 'is-error': cpu.v_value.status === 2 })}
            style={{ marginRight: hasCpuRatio ? '156px' : '' }}
          >
            <span className="tc-input-group-addon">{cpu.type}</span>
            <Bubble placement="bottom" content={cpu.v_value.status === 2 ? <p>{cpu.v_value.message}</p> : null}>
              <input
                type="text"
                className="tc-15-input-text m"
                style={{ width: '60px' }}
                placeholder={t('不限制')}
                value={cpu.value}
                onChange={e => actions.editWorkload.updateCpuLimit({ value: e.target.value }, cKey, cpu.id + '')}
                onBlur={e => actions.validate.workload.validateAllCpuLimit(container)}
              />
            </Bubble>
            <span className="inline-help-text" style={{ marginRight: '5px' }}>
              {partition}
            </span>
          </div>
        );
      }
    });
  }

  /** mem的相关选项 */
  private _renderMemItem(hasMemoryRatio?: boolean) {
    const { actions, subRoot, cKey } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    const container = containers.find(item => item.id === cKey);

    // 当limit不填是保持和request一致灰色提示
    const _renderPlaceholder = (): string => {
      let placeholder = t('不限制');
      if (container.memLimit[0].value && container.memLimit[0].v_value.status === 1) {
        placeholder = container.memLimit[0].value + '';
      }
      return placeholder;
    };

    return container.memLimit.map((mem, index) => {
      const partition = mem.type === 'request' ? '-' : 'MiB';
      if (hasMemoryRatio && mem.type === 'request') {
        return <noscript />;
      } else {
        return (
          <div
            key={index}
            className={classnames('tc-input-group', { 'is-error': mem.v_value.status === 2 })}
            style={{ marginRight: hasMemoryRatio ? '130px' : '' }}
          >
            <span className={classnames('tc-input-group-addon')}>{mem.type}</span>
            <Bubble placement="bottom" content={mem.v_value.status === 2 ? <p>{mem.v_value.message}</p> : null}>
              <input
                type="text"
                className="tc-15-input-text m"
                style={{ width: '60px' }}
                placeholder={mem.type === 'request' ? t('不限制') : _renderPlaceholder()}
                value={mem.value}
                onChange={e => actions.editWorkload.updateMemLimit({ value: e.target.value }, cKey, mem.id + '')}
                onBlur={e => actions.validate.workload.validateAllMemLimit(container)}
              />
            </Bubble>
            <span className="inline-help-text" style={{ marginRight: '5px' }}>
              {partition}
            </span>
          </div>
        );
      }
    });
  }
}
