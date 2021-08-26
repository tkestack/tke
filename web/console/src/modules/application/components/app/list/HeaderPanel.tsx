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

import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../AppContainer';
import { FormPanel } from '@tencent/ff-component';
import { namespace } from '@config/resource/k8sConfig';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HeaderPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    /** 拉取应用列表 */
    // actions.app.list.poll();
    //不要保存filter旧数据
    actions.cluster.list.reset();
    actions.cluster.list.applyFilter();
  }

  render() {
    let { clusterList, namespaceList, actions, route } = this.props;
    let urlParam = router.resolve(route);
    const { mode } = urlParam;
    return (
      <Justify
        left={
          <React.Fragment>
            <h2>{t('应用管理')}</h2>
            <FormPanel.InlineText>{t('集群：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={t('集群')}
              model={clusterList}
              action={actions.cluster.list}
              value={clusterList.selection ? clusterList.selection.metadata.name : ''}
              onChange={value => {
                actions.cluster.list.selectCluster(value);
              }}
              valueField={x => x.metadata.name}
              displayField={x => `${x.spec.displayName}`}
            ></FormPanel.Select>
            <FormPanel.InlineText>{t('命名空间：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={'命名空间'}
              model={namespaceList}
              action={actions.namespace.list}
              value={namespaceList.selection ? namespaceList.selection.metadata.name : ''}
              valueField={x => x.metadata.name}
              displayField={x => `${x.metadata.name}`}
              onChange={value => {
                actions.namespace.list.selectNamespace(value);
              }}
            />
          </React.Fragment>
        }
      />
    );
  }
}
