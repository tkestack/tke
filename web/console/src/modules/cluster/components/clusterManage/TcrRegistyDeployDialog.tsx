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

import { Button, Modal, Select, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import {
    initValidator, LinkButton, SelectList, TipInfo, Validation
} from '../../../../modules/common';
import { allActions } from '../../actions';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';

interface TcrRegistyDeployState {
  clusterSelection: string; //集群选择
  v_clusterSelection: Validation; //选择校验
  isShow: boolean; //是否显示对话显示框
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class TcrRegistyDeployDialog extends React.Component<RootProps, TcrRegistyDeployState> {
  constructor(props, content) {
    super(props, content);
    this.state = {
      clusterSelection: '',
      v_clusterSelection: initValidator,
      isShow: false
    };
  }

  _validateClusterSelection(cluster) {
    let status = 0;
    let message = '';
    if (!cluster) {
      status = 2;
      message = '请选择一个集群';
    } else {
      status = 1;
      message = '';
    }
    return {
      status,
      message
    };
  }

  validateClusterSelection() {
    let result = this._validateClusterSelection(this.state.clusterSelection);
    this.setState({
      v_clusterSelection: result
    });
  }

  selectCluster(cluster) {
    this.setState(
      {
        clusterSelection: cluster
      },
      this.validateClusterSelection
    );
  }
  crtateCluster() {
    let { route } = this.props;

    router.navigate(
      {
        sub: 'createIC'
      },
      {
        rid: route.queries['rid']
      }
    );
  }

  importCluster() {
    let { route } = this.props;

    router.navigate(
      {
        sub: 'create'
      },
      {
        rid: route.queries['rid']
      }
    );
  }
  cancel() {
    let { route } = this.props;
    let params = router.resolve(route);
    let newRouteQueies = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { registry: undefined })));
    router.navigate(params, newRouteQueies);
  }

  handleSubmit() {
    //跳转的到创建workload的界面
    this.validateClusterSelection();
    let registry = '',
      tag = '',
      name = '';

    let { clusterSelection } = this.state;

    if (this._validateClusterSelection(clusterSelection).status === 2) {
      return;
    }

    let { route, actions, cluster, subRoot } = this.props;
    let { containers } = subRoot.workloadEdit;

    //镜像名称
    let reg = /(.*)\/(.*):(.*)/g;

    let result = reg.exec(route.queries['registry']);

    tag = result && result[3];
    name = result && result[2];
    registry = result && result[1] + '/' + result[2];

    let selectCluster = cluster.list.data.records.find(item => item.metadata.name === clusterSelection);

    actions.cluster.selectCluster([selectCluster], true);

    //进行路由的跳转,跳转到工作负载的创建页面
    router.navigate(
      {
        sub: 'sub',
        mode: 'create',
        type: 'resource',
        resourceName: 'deployment' //默认为deployment
      },
      {
        clusterId: clusterSelection,
        rid: route.queries['rid'],
        np: 'default'
      }
    );

    actions.resource.initResourceInfoAndFetchData(true, 'deployment');

    let editingContainer = containers.find(c => c.status === 'editing');

    actions.editWorkload.updateContainer({ registry, tag, name }, editingContainer.id as string);
  }

  render() {
    let { clusterSelection, v_clusterSelection } = this.state;
    let { route } = this.props;

    let selectClusterList = this.props.cluster.list.data.records.map(item => ({
      value: item.metadata.name,
      text: `${item.metadata.name}（${item.spec.displayName || '-'}）`,
      disabled: item.status.phase.toLowerCase() !== 'running' //状态不为running的不可选择
    }));

    return (
      <Modal
        caption={t('创建工作负载')}
        disableEscape={true}
        visible={!!route.queries['registry'] && !!/(.*)\/(.*):(.*)/g.exec(route.queries['registry'])}
      >
        <Modal.Body>
          <FormPanel isNeedCard={false}>
            <FormPanel.Item
              validator={v_clusterSelection}
              label={t('集群')}
              select={{
                options: selectClusterList,
                value: clusterSelection,
                onChange: value => {
                  this.selectCluster(value);
                  this.validateClusterSelection();
                }
              }}
              message={
                !selectClusterList.length && (
                  <Text theme="danger">
                    {t('当前账号下无可用集群，请')}
                    <LinkButton onClick={this.crtateCluster.bind(this)}>{t('新建独立集群')}</LinkButton>
                    {t('或者')}
                    <LinkButton onClick={this.importCluster.bind(this)}>{t('导入集群')}</LinkButton>
                  </Text>
                )
              }
            />
          </FormPanel>
          <TipInfo type="success" style={{ marginTop: '20px ' }}>
            {t(`当前工作负载的实例内容器将默认使用该镜像：${route.queries['registry']}`)}
          </TipInfo>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={this.handleSubmit.bind(this)}>
            {t('新建')}
          </Button>
          <Button onClick={this.cancel.bind(this)}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
