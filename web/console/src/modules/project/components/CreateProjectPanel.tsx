import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators, deepClone } from '@tencent/qcloud-lib';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Icon, Select, Text, Button, Alert, Modal } from '@tencent/tea-component';

import { FormItem, FormPanel, InputField, LinkButton } from '../../common/components';
import { FormLayout } from '../../common/layouts';
import { allActions } from '../actions';
import { resourceLimitTypeToText, resourceTypeToUnit } from '../constants/Config';
import { CreateProjectResourceLimitPanel } from './CreateProjectResourceLimitPanel';
import { EditProjectManagerPanel } from './EditProjectManagerPanel';
import { RootProps } from './ProjectApp';
import { getWorkflowError } from '../../common';
import { router } from '../router';
import { projectActions } from '../actions/projectActions';
import { OperationState, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { ProjectResourceLimit } from '../models/Project';

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
    let { projectEdition, actions, project, route, createProject, cluster } = this.props;

    let projectListOpions = project.list.data.records.map(item => {
      return { text: `${item.metadata.name}(${item.spec.displayName})`, value: item.metadata.name };
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
        <FormPanel.Item label={t('业务成员')}>
          <div style={{ width: 600 }}>
            <EditProjectManagerPanel {...this.props} />
          </div>
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
          }}
        />
      </Modal>
    );
  }
}
