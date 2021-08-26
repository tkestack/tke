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

import { Bubble, Button, Modal, Table, TableColumn, Text, Checkbox } from '@tea/component';
import { stylize } from '@tea/component/table/addons/stylize';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../../../actions';
import { ContainerStatusMap, podRemoteShellOptions } from '../../../constants/Config';
import { PodContainer } from '../../../models';
import { RootProps } from '../../ClusterApp';
import { FormPanel } from '@tencent/ff-component';

interface ResourcePodRemoteLoginDialogState {
  /** 是否使用用户自定义的shell */
  isUserDefined?: boolean;

  /** 用户选择的shell */
  shellSelected?: string;

  /** 用户自己输入的shell */
  userDefinedShell?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourcePodRemoteLoginDialog extends React.Component<RootProps, ResourcePodRemoteLoginDialogState> {
  constructor(props) {
    super(props);
    this.state = {
      isUserDefined: false,
      shellSelected: '/bin/bash',
      userDefinedShell: ''
    };
  }

  render() {
    let { actions, subRoot, route } = this.props,
      { resourceDetailState } = subRoot,
      { podSelection, isShowLoginDialog } = resourceDetailState;

    if (!isShowLoginDialog) {
      return <noscript />;
    }

    const cancel = () => {
      // 关闭 远程登录的弹窗
      actions.resourceDetail.pod.toggleLoginDialog();
      // 置空当前的pod的选项
      actions.resourceDetail.pod.podSelect([]);
    };

    let containers = podSelection[0].spec.containers,
      containerStatus = podSelection[0].status.containerStatuses;

    // 容器登录的web-console的网址
    let loginUrl = '';
    // 容器登录的命名空间
    let namespace: string = podSelection[0] ? podSelection[0].metadata.namespace : 'default';

    let columns: TableColumn<PodContainer>[] = [
      {
        key: 'name',
        header: t('容器名称'),
        width: '45%',
        render: x => (
          <div className="sl-editor-name">
            <span className="text-overflow m-width" title={x.name}>
              {x.name}
            </span>
          </div>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '25%',
        render: x => this._reduceContainerStatus(x, containerStatus)
      },
      {
        key: 'description',
        header: t('操作'),
        width: '25%',
        render: x => (
          <div className="text-left">
            <div className="sl-editor-name">
              <a
                href={
                  '/webtty.html?clusterName=' +
                  route.queries['clusterId'] +
                  '&projectName=' +
                  (route.queries['projectName'] || '') +
                  '&podName=' +
                  podSelection[0].metadata.name +
                  '&containerName=' +
                  (x.name ? x.name : '') +
                  '&namespace=' +
                  namespace +
                  `&command=${this.state.isUserDefined ? this.state.userDefinedShell : this.state.shellSelected}`
                }
                target="_blank"
              >
                {t('登录')}
              </a>
            </div>
          </div>
        )
      }
    ];

    let containersLength = containers.length;
    return (
      <Modal visible={true} caption={t('容器登录')} onClose={cancel} disableEscape={true}>
        <Modal.Body>
          <Text>
            <Trans count={containersLength}>
              该实例下共有<strong className="text-warning">{{ containersLength }}个</strong>容器
            </Trans>
          </Text>
          <Table
            columns={columns}
            records={containers}
            addons={[
              stylize({
                className: 'ovm-dialog-tablepanel',
                bodyStyle: { overflowY: 'auto', height: 160 }
              })
            ]}
          />
          <FormPanel isNeedCard={false} className="tea-mt-2n">
            <FormPanel.Item label={t('执行命令行')} text>
              <Checkbox
                value={this.state.isUserDefined}
                onChange={value => this.setState({ isUserDefined: value })}
                className="tea-mb-1n"
              >
                {t('使用自定义执行命令行')}
              </Checkbox>
              {this.state.isUserDefined ? (
                <FormPanel.Input
                  value={this.state.userDefinedShell}
                  onChange={value => this.setState({ userDefinedShell: value })}
                  placeholder="eg: /bin/bash"
                />
              ) : (
                <FormPanel.Select
                  options={podRemoteShellOptions}
                  value={this.state.shellSelected}
                  onChange={value => this.setState({ shellSelected: value })}
                />
              )}
            </FormPanel.Item>
          </FormPanel>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }

  /** 处理容器的状态 */
  private _reduceContainerStatus(container: PodContainer, containerStatus: any[]) {
    let finder = containerStatus ? containerStatus.find(c => c.name === container.name) : undefined,
      statusKey = finder && Object.keys(finder.state)[0];

    return (
      <div>
        <span
          className={classnames(
            'text-overflow',
            ContainerStatusMap[statusKey] && ContainerStatusMap[statusKey].classname
          )}
        >
          {ContainerStatusMap[statusKey] ? ContainerStatusMap[statusKey].text : '-'}
        </span>
        {statusKey && statusKey !== 'running' && (
          <Bubble placement="right" content={finder.state[statusKey].reason || null}>
            <div className="tc-15-bubble-icon">
              <i className="tc-icon icon-what" />
            </div>
          </Bubble>
        )}
      </div>
    );
  }
}
