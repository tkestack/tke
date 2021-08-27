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

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, deepClone, isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Alert, Bubble, Button, Icon, Modal, Text } from '@tencent/tea-component';

import { getWorkflowError, RequestParams, ResourceInfo } from '../../common';
import { allActions } from '../actions';
import { projectActions } from '../actions/projectActions';
import { resourceLimitTypeToText, resourceTypeToUnit, PlatformTypeEnum } from '../constants/Config';
import { ProjectResourceLimit } from '../models/Project';
import { router } from '../router';
import { CreateProjectResourceLimitPanel } from './CreateProjectResourceLimitPanel';
import { EditProjectManagerPanel } from './EditProjectManagerPanel';
import { RootProps } from './ProjectApp';
import { resourceConfig } from '@config/resourceConfig';
import { reduceK8sRestfulPath } from '@helper/urlUtil';
import { Method, reduceNetworkRequest } from '@helper/reduceNetwork';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class CreateProjectPanel extends React.Component<
  RootProps,
  { currentClusterIndex: number; isShowDialog: boolean }
> {
  state = {
    currentClusterIndex: 0,
    isShowDialog: false
  };
  componentDidMount() {
    let { actions, project, manager } = this.props;
    actions.cluster.applyFilter({});
    if (project.list.data.recordCount === 0) {
      actions.project.applyFilter({});
    }
    if (manager.list.data.recordCount === 0) {
      actions.manager.applyFilter({});
    }
    this.getUserInfo();
  }

  //获取用户信息包括用户业务信息
  async getUserInfo() {
    let { actions } = this.props;
    let infoResourceInfo: ResourceInfo = resourceConfig()['info'];
    let url = reduceK8sRestfulPath({ resourceInfo: infoResourceInfo });
    let params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      let response = await reduceNetworkRequest(params);
      let loginUserInfo = {
        id: '',
        name: '',
        displayName: ''
      };
      if (!response.code) {
        const { uid, name, extra } = response.data;
        loginUserInfo = {
          id: uid,
          name,
          displayName: extra.displayname ? extra.displayname[0] : ''
        };
      }
      actions.project.selectManager([loginUserInfo]);
    } catch (error) {}
  }

  formatResourceLimit(resourceLimit: ProjectResourceLimit[]) {
    let content = resourceLimit.map((item, index) => (
      <Text parent="p" key={index}>{`${resourceLimitTypeToText[item.type]}:${item.value}${
        resourceTypeToUnit[item.type]
      }`}</Text>
    ));
    return content;
  }

  _handleSubmit() {
    let { actions, projectEdition } = this.props;
    actions.project.validateProjection();
    if (projectActions._validateProjection(projectEdition)) {
      actions.project.createProject.start([projectEdition]);
      actions.project.createProject.perform();
    }
  }

  render() {
    let { projectEdition, actions, project, route, createProject, cluster, platformType } = this.props;

    let projectListOpions = project.list.data.records.map(item => {
      return { text: `${item.metadata.name}(${item.spec.displayName})`, value: item.metadata.name };
    });

    projectListOpions.unshift({
      text: '无上级业务',
      value: ''
    });

    let finalClusterList = deepClone(cluster);

    let parentProjectSelection = projectEdition.parentProject
      ? project.list.data.records.find(item => item.metadata.name === projectEdition.parentProject)
      : null;
    //筛选出project中的集群
    if (parentProjectSelection) {
      let parentClusterList = parentProjectSelection.spec.clusters
        ? Object.keys(parentProjectSelection.spec.clusters)
        : [];
      finalClusterList.list.data.records = finalClusterList.list.data.records.filter(
        item => parentClusterList.indexOf(item.clusterId + '') !== -1
      );
      finalClusterList.list.data.recordCount = finalClusterList.list.data.records.length;
    }

    let failed = createProject.operationState === OperationState.Done && !isSuccessWorkflow(createProject);
    return (
      <FormPanel>
        <FormPanel.Item
          label={t('业务名称')}
          errorTipsStyle="Icon"
          message={t('业务名称不能超过63个字符')}
          validator={projectEdition.v_displayName}
          input={{
            value: projectEdition.displayName,
            onChange: value => actions.project.inputProjectName(value),
            onBlur: e => {
              actions.project.validateDisplayName(e.target.value);
            }
          }}
        />
        <FormPanel.Item label={t('业务管理员')}>
          <div style={{ width: 600 }}>
            <EditProjectManagerPanel {...this.props} />
          </div>
          {(!projectEdition.members || projectEdition.members.length === 0) && (
            <Text theme="danger" style={{ fontSize: '12px' }}>
              需要至少选择一个责任人
            </Text>
          )}
        </FormPanel.Item>
        <FormPanel.Item label={t('集群')}>
          {projectEdition.clusters.map((item, index) => {
            let resourceLimitContent = this.formatResourceLimit(item.resourceLimits);
            return (
              <React.Fragment key={index}>
                <div style={{ marginBottom: 5 }} className={item.v_name.status === 2 ? 'is-error' : ''}>
                  <Bubble placement="top" content={item.v_name.status === 2 ? <p>{item.v_name.message}</p> : null}>
                    <div style={{ display: 'inline-block' }}>
                      <FormPanel.Select
                        label={t('集群')}
                        value={item.name}
                        model={finalClusterList}
                        action={actions.cluster}
                        valueField={x => x.clusterId}
                        displayField={x => `${x.clusterId}(${x.clusterName})`}
                        onChange={clusterId => {
                          actions.project.updateClusters(index, clusterId);
                          actions.project.validateClustersName(index);
                        }}
                        style={{ marginRight: 5 }}
                      ></FormPanel.Select>
                    </div>
                  </Bubble>
                  <Button
                    type={'link'}
                    disabled={item.name === ''}
                    onClick={() =>
                      this.setState({
                        isShowDialog: true,
                        currentClusterIndex: index
                      })
                    }
                  >
                    {t('填写资源限制')}
                  </Button>
                  {resourceLimitContent.length && (
                    <Bubble content={resourceLimitContent}>
                      <Icon type="detail" />
                    </Bubble>
                  )}
                  <Button
                    icon={'close'}
                    disabled={projectEdition.clusters.length === 1}
                    onClick={() => {
                      actions.project.deleteClusters(index);
                    }}
                  />
                </div>
              </React.Fragment>
            );
          })}
          <Button type={'link'} onClick={() => actions.project.addClusters()}>
            {t('新增集群')}
          </Button>
        </FormPanel.Item>

        <FormPanel.Item label={t('上级业务')}>
          <FormPanel.Select
            label={t('上级业务')}
            options={projectListOpions}
            value={projectEdition.parentProject}
            disabled={platformType === PlatformTypeEnum.Business}
            onChange={value => {
              actions.project.inputParentPorject(value);
            }}
          />
        </FormPanel.Item>
        <FormPanel.Footer>
          <React.Fragment>
            <Button
              type="primary"
              disabled={createProject.operationState === OperationState.Performing}
              onClick={this._handleSubmit.bind(this)}
            >
              {failed ? t('重试') : t('完成')}
            </Button>
            <Button
              type="weak"
              onClick={() => {
                actions.project.clearEdition();
                router.navigate({}, route.queries);
              }}
            >
              {t('取消')}
            </Button>
            {failed ? (
              <Alert
                type="error"
                style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
              >
                {getWorkflowError(createProject)}
              </Alert>
            ) : (
              <noscript />
            )}
          </React.Fragment>
        </FormPanel.Footer>
        {this._renderEditProjectLimitDialog()}
      </FormPanel>
    );
  }

  private _renderEditProjectLimitDialog() {
    const { actions, project, projectEdition } = this.props;
    let { currentClusterIndex, isShowDialog } = this.state;
    let parentProjectSelection = projectEdition.parentProject
      ? project.list.data.records.find(item => item.metadata.name === projectEdition.parentProject)
      : null;
    let clusterName = projectEdition.clusters[currentClusterIndex].name;

    let parentResourceLimits =
      parentProjectSelection && clusterName ? parentProjectSelection.spec.clusters[clusterName].hard : {};

    return (
      <Modal visible={isShowDialog} caption={t('编辑资源限制')} onClose={() => this.setState({ isShowDialog: false })}>
        <CreateProjectResourceLimitPanel
          parentResourceLimits={parentResourceLimits}
          onCancel={() => {
            this.setState({ isShowDialog: false, currentClusterIndex: 0 });
          }}
          resourceLimits={projectEdition.clusters[currentClusterIndex].resourceLimits}
          onSubmit={requestLimits => {
            actions.project.updateClustersLimit(currentClusterIndex, requestLimits);
            this.setState({ isShowDialog: false, currentClusterIndex: 0 });
          }}
        />
      </Modal>
    );
  }
}
