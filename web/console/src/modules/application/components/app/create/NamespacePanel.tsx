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
import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../AppContainer';
import { router } from '../../../router';
import { FormPanel } from '@tencent/ff-component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class NamespacePanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions, appCreation } = this.props;
    /** 拉取集群列表 */
    //不要保存filter旧数据
    actions.cluster.list.reset();
    actions.cluster.list.applyFilter();
  }

  render() {
    let { actions, route, appCreation, clusterList, namespaceList } = this.props;
    let action = actions.app.create.addAppWorkflow;

    return (
      <React.Fragment>
        <FormPanel.Item
          label={t('运行集群')}
          vkey="spec.targetCluster"
          select={{
            showRefreshBtn: true,
            value: appCreation.spec ? appCreation.spec.targetCluster : '',
            model: clusterList,
            action: actions.cluster.list,
            valueField: x => x.metadata.name,
            displayField: x => `${x.metadata.name}(${x.spec.displayName})`,
            onChange: value => {
              actions.cluster.list.selectCluster(value);
              actions.app.create.updateCreationState({
                metadata: Object.assign({}, appCreation.metadata, {
                  namespace: ''
                }),
                spec: Object.assign({}, appCreation.spec, {
                  targetCluster: value
                })
              });
            }
          }}
        ></FormPanel.Item>
        <FormPanel.Item
          label={t('命名空间')}
          vkey="metadata.namespace"
          select={{
            showRefreshBtn: true,
            value: appCreation.metadata && appCreation.metadata.namespace ? appCreation.metadata.namespace : '',
            model: namespaceList,
            action: actions.namespace.list,
            valueField: x => x.metadata.name,
            displayField: x => `${x.metadata.name}`,
            onChange: value => {
              actions.namespace.list.selectNamespace(value);
              actions.app.create.updateCreationState({
                metadata: Object.assign({}, appCreation.metadata, {
                  namespace: value
                })
              });
            }
          }}
        ></FormPanel.Item>
      </React.Fragment>
    );
  }
}
