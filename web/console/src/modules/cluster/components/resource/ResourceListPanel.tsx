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
import { TipInfo } from 'src/modules/common';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Alert, ContentView, ExternalLink, Justify, Text } from '@tencent/tea-component';

import { isEmpty } from '../../../common/utils';
import { allActions } from '../../actions';
import { SubRouter } from '../../models';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';
import { ClusterDetailPanel } from './clusterInfomation/ClusterDetail';
import { BatchDrainComputerDialog } from './nodeManage/BatchDrainComputerDialog';
import { BatchTurnOnScheduleComputerDialog } from './nodeManage/BatchTurnOnScheduleComputerDialog';
import { BatchUnScheduleComputerDialog } from './nodeManage/BatchUnScheduleComputerDialog';
import { ComputerActionPanel } from './nodeManage/ComputerActionPanel';
import { ComputerStatusDialog } from './nodeManage/ComputerStatusDialog';
import { ComputerTablePanel } from './nodeManage/ComputerTablePanel';
import { DeleteComputerDialog } from './nodeManage/DeleteComputerDialog';
import { UpdateNodeLabelDialog } from './nodeManage/UpdateNodeLabelDialog';
import { UpdateNodeTaintDialog } from './nodeManage/UpdateNodeTaintDialog';
import { ResourceHeaderPanel } from './ResourceHeaderPanel';
import { ResourceSidebarPanel } from './ResourceSidebarPanel';
import { ResourceActionPanel } from './resourceTableOperation/ResourceActionPanel';
import { ResourceDeleteDialog } from './resourceTableOperation/ResourceDeleteDialog';
import { ResourceEventPanel } from './resourceTableOperation/ResourceEventPanel';
import { ResourceLogPanel } from './resourceTableOperation/ResourceLogPanel';
import { ResourceTablePanel } from './resourceTableOperation/ResourceTablePanel';
import { HPAPanel } from '../scale/hpa';
import { CronHpaPanel } from '../scale/cronhpa';
import { VMListPanel, SnapshotTablePanel } from './virtual-machine';

const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

export interface ResourceListPanelProps extends RootProps {
  /** subRouterList */
  subRouterList: SubRouter[];
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceListPanel extends React.Component<ResourceListPanelProps, {}> {
  render() {
    const { subRoot, route, subRouterList, actions } = this.props,
      urlParams = router.resolve(route),
      { resourceInfo } = subRoot;
    let content: JSX.Element;
    let headTitle = '';
    const resource = urlParams['resourceName'];
    // 判断应该展示什么组件
    switch (resource) {
      case 'info':
        content = <ClusterDetailPanel {...this.props} />;
        headTitle = t('基础信息');
        break;

      case 'node':
        content = (
          <React.Fragment>
            <ComputerActionPanel />
            <ComputerTablePanel />
            <BatchUnScheduleComputerDialog />
            <BatchTurnOnScheduleComputerDialog />
            <UpdateNodeLabelDialog {...this.props} />
            <UpdateNodeTaintDialog />
            <DeleteComputerDialog {...this.props} />
            <ComputerStatusDialog
              dialogState={this.props.dialogState}
              machine={this.props.subRoot.computerState.machine}
            />
            <BatchDrainComputerDialog {...this.props} />
            <div id="ComputerMonitorPanel" />
          </React.Fragment>
        );
        headTitle = t('节点列表');
        break;

      case 'log':
        content = <ResourceLogPanel />;
        headTitle = t('日志');
        break;

      case 'event':
        content = <ResourceEventPanel />;
        headTitle = t('事件');
        break;

      case 'hpa':
        content = <HPAPanel />;
        headTitle = t('HorizontalPodAutoscaler');
        break;

      case 'cronhpa':
        content = <CronHpaPanel />;
        headTitle = t('CronHPA');
        break;

      case 'virtual-machine':
        content = <VMListPanel route={route} />;
        headTitle = 'Virtual Machine';
        break;

      default:
        content = isEmpty(resourceInfo) ? (
          loadingElement
        ) : (
          <React.Fragment>
            <ResourceActionPanel />
            <ResourceTablePanel />
          </React.Fragment>
        );

        headTitle = resourceInfo.headTitle;

        break;
    }

    return (
      <React.Fragment>
        <ContentView>
          <ContentView.Header>
            <ResourceHeaderPanel />
          </ContentView.Header>
          <ContentView.Body
            sidebar={<ResourceSidebarPanel route={route} actions={actions} subRouterList={subRouterList} />}
          >
            <ContentView>
              <ContentView.Header>
                <Justify
                  left={
                    <React.Fragment>
                      {resource === 'tapp' && (
                        <Alert style={{ marginBottom: '8px' }}>
                          <Text verticalAlign="middle">
                            {t(
                              'TApp是腾讯云自研的一种workload类型，支持有/无状态的应用类型，可进行Pod级别的指定删除、原地升级、挂载独立数据盘等操作，'
                            )}
                          </Text>
                          {/* <ExternalLink>了解更多</ExternalLink> */}
                        </Alert>
                      )}
                      <h2 className="tea-h2">{headTitle || ''}</h2>
                    </React.Fragment>
                  }
                />
              </ContentView.Header>
              <ContentView.Body>{content}</ContentView.Body>
            </ContentView>
          </ContentView.Body>
        </ContentView>
        <ResourceDeleteDialog />
      </React.Fragment>
    );
  }
}
