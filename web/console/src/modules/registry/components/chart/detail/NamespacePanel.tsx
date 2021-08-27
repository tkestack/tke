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
import { RootProps } from '../ChartApp';
import { Button, Bubble, Icon } from '@tencent/tea-component';
import { router } from '../../../router';
import { FormPanel } from '@tencent/ff-component';
import { ChartInfoFilter } from '../../../models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface NamespaceProps extends RootProps {
  chartInfoFilter: ChartInfoFilter;
}

@connect(state => state, mapDispatchToProps)
export class NamespacePanel extends React.Component<NamespaceProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    /** 拉取集群列表 */
    //不要保存filter旧数据
    actions.cluster.list.reset();
    actions.cluster.list.applyFilter({ chartInfoFilter: this.props.chartInfoFilter });
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
            value: appCreation.spec ? appCreation.spec.targetCluster : '',
            model: clusterList,
            action: actions.cluster.list,
            valueField: x => x.metadata.name,
            displayField: x => `${x.metadata.name}(${x.spec.displayName})`,
            onChange: value => {
              //选择集群时不回调请求chartInfo
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
            value: appCreation.metadata && appCreation.metadata.namespace ? appCreation.metadata.namespace : '',
            model: namespaceList,
            action: actions.namespace.list,
            valueField: x => x.metadata.name,
            displayField: x => `${x.metadata.name}`,
            onChange: value => {
              //选择命名空间时不回调请求chartInfo
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
