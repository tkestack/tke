import { FormPanel } from '@tencent/ff-component';
import { isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { bindActionCreators, deepClone } from '@tencent/qcloud-lib';
import { t } from '@tencent/tea-app/lib/i18n';
import { Alert, Button, Modal } from '@tencent/tea-component';
import * as React from 'react';
import { connect } from 'react-redux';
import { getWorkflowError } from '../../common';
import { allActions } from '../actions';
import { namespaceActions } from '../actions/namespaceActions';
import { resourceLimitTypeToText, resourceTypeToUnit } from '../constants/Config';
import { ProjectResourceLimit } from '../models/Project';
import { router } from '../router';
import { CreateProjectResourceLimitPanel } from './CreateProjectResourceLimitPanel';
import { RootProps } from './ProjectApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class CreateNamespacePanel extends React.Component<RootProps, {}> {
  state = {
    isShowDialog: false
  };
  componentDidMount() {
    let { actions, project, manager } = this.props;
    actions.cluster.applyFilter({});
    if (project.list.data.recordCount === 0) {
      actions.project.applyFilter({});
    }
  }

  componentWillUnmount() {
    this.props.actions.namespace.clearEdition();
  }
  _handleSubmit() {
    let { actions, namespaceEdition, project, route } = this.props;
    actions.namespace.validateNamespaceEdition();
    if (namespaceActions._validateNamespaceEdition(namespaceEdition)) {
      actions.namespace.createNamespace.start([namespaceEdition], {
        projectId: project.selections[0] ? project.selections[0].metadata.name : route.queries['projectId']
      });
      actions.namespace.createNamespace.perform();
    }
  }
  formatResourceLimit(resourceLimit: ProjectResourceLimit[]) {
    let content = resourceLimit.map((item, index) => (
      <FormPanel.Text key={index}>{`${resourceLimitTypeToText[item.type]}:${item.value}${
        resourceTypeToUnit[item.type]
      }`}</FormPanel.Text>
    ));
    return content;
  }

  render() {
    let { namespaceEdition, actions, cluster, project, createNamespace, route } = this.props;

    let projectSelection = project.selections[0] ? project.selections[0] : null;

    let finalClusterList = deepClone(cluster);
    //筛选出project中的集群
    if (projectSelection) {
      let projectClusterList = projectSelection.spec.clusters ? Object.keys(projectSelection.spec.clusters) : [];
      finalClusterList.list.data.records = finalClusterList.list.data.records.filter(
        item => projectClusterList.indexOf(item.clusterId + '') !== -1
      );
      finalClusterList.list.data.recordCount = finalClusterList.list.data.records.length;
    }

    let failed = createNamespace.operationState === OperationState.Done && !isSuccessWorkflow(createNamespace);

    return (
      <FormPanel>
        <FormPanel.Item
          message={t(
            '最长48个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾，名称不能以"kube-"开头'
          )}
          text
          label={t('名称')}
          validator={namespaceEdition.v_namespaceName}
          errorTipsStyle="Icon"
          input={{
            value: namespaceEdition.namespaceName,
            onChange: actions.namespace.inputNamespaceName,
            onBlur: () => {
              actions.namespace.validateNamespaceName();
            }
          }}
        />
        <FormPanel.Item text label={t('业务')}>
          {projectSelection ? (
            <React.Fragment>
              <FormPanel.InlineText>
                {t(projectSelection.metadata.name + '(' + projectSelection.spec.displayName + ')')}
              </FormPanel.InlineText>
            </React.Fragment>
          ) : (
            <noscript />
          )}
        </FormPanel.Item>
        <FormPanel.Item
          label={t('集群')}
          validator={namespaceEdition.v_clusterName}
          errorTipsStyle="Icon"
          select={{
            model: finalClusterList,
            action: actions.cluster,
            value: namespaceEdition.clusterName,
            valueField: x => x.clusterId,
            displayField: x => `${x.clusterId}(${x.clusterName})`,
            onChange: value => {
              actions.namespace.selectCluster(value);
              actions.namespace.validateClusterName();
            }
          }}
        />
        <FormPanel.Item label={'资源限制'}>
          {this.formatResourceLimit(namespaceEdition.resourceLimits)}
          <Button
            disabled={namespaceEdition.clusterName === ''}
            icon={'pencil'}
            onClick={() => {
              this.setState({
                isShowDialog: true
              });
            }}
          ></Button>
        </FormPanel.Item>
        <FormPanel.Footer>
          <React.Fragment>
            <Button
              type="primary"
              disabled={createNamespace.operationState === OperationState.Performing}
              onClick={this._handleSubmit.bind(this)}
            >
              {failed ? t('重试') : t('完成')}
            </Button>
            <Button
              type="weak"
              onClick={() => {
                actions.namespace.clearEdition();
                router.navigate({ sub: 'detail', tab: 'namespace' }, route.queries);
              }}
            >
              {t('取消')}
            </Button>
            {failed ? (
              <Alert
                type="error"
                style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
              >
                {getWorkflowError(createNamespace)}
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
    const { actions, project, namespaceEdition } = this.props;
    let { isShowDialog } = this.state;
    let projectSelection = project.selections[0] ? project.selections[0] : null;

    let clusterName = namespaceEdition.clusterName;

    let resourceLimits = projectSelection && clusterName ? projectSelection.spec.clusters[clusterName].hard : {};
    return (
      <Modal visible={isShowDialog} caption={t('编辑资源限制')} onClose={() => this.setState({ isShowDialog: false })}>
        <CreateProjectResourceLimitPanel
          parentResourceLimits={resourceLimits}
          onCancel={() => {
            this.setState({ isShowDialog: false });
          }}
          resourceLimits={namespaceEdition.resourceLimits}
          onSubmit={requestLimits => {
            actions.namespace.updateNamespaceResourceLimit(requestLimits);
          }}
        />
      </Modal>
    );
  }
}
