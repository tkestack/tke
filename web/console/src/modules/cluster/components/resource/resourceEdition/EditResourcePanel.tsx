import * as React from 'react';
import { Button, Icon } from '@tea/component';
import { OperationState, isSuccessWorkflow, FetchState } from '@tencent/ff-redux';
import { TipInfo } from '../../../../common/components';
import { RootProps } from '../../ClusterApp';
import { uuid } from '@tencent/qcloud-lib';
import { typeMapName, resourceConfig } from '../../../../../../config';
import { MainBodyLayout, FormLayout } from '../../../../common/layouts';
import { CreateResource } from '../../../models';
import { getWorkflowError, isEmpty } from '../../../../common/utils';
import { router } from '../../../router';
import { SubHeaderPanel } from './SubHeaderPanel';
import { YamlEditorPanel } from '../YamlEditorPanel';
import { EditServicePanel } from './EditServicePanel';
import { EditNamespacePanel } from './EditNamespacePanel';
import { EditResourceVisualizationPanel } from './EditResourceVisualizationPanel';
import { EditSecretPanel } from './EditSecretPanel';
import { EditConfigMapPanel } from './EditConfigMapPanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { CreateComputerPanel } from '../nodeManage/CreateComputerPanel';
import { EditLbcfPanel } from './EditLbcfPanel';

interface EditResourcePanelState {
  /** edited data */
  config?: string;
  parseError?;
}

export class EditResourcePanel extends React.Component<RootProps, EditResourcePanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      config: ''
    };
  }

  componentDidMount() {
    let { subRoot, actions, route } = this.props,
      urlParams = router.resolve(route),
      { resourceInfo } = subRoot;

    // 编辑模式下，需要进行yaml的拉取，已经resourceInfo的前提下，进行拉取
    urlParams['mode'] === 'modify' && !isEmpty(resourceInfo) && actions.resourceDetail.fetchResourceYaml.fetch();
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let oldYamlList = this.props.subRoot.resourceDetailState.yamlList,
      newYamlList = nextProps.subRoot.resourceDetailState.yamlList,
      resourceSelection = this.props.subRoot.resourceOption.resourceSelection,
      newResourceSelection = nextProps.subRoot.resourceOption.resourceSelection,
      newMode = nextProps.subRoot.mode;

    if (oldYamlList.data.recordCount === 0 && newYamlList.data.recordCount) {
      this.setState({
        config: newYamlList.data.records[0] || ''
      });
    } else if (newMode === 'modify' && resourceSelection.length === 0 && newResourceSelection.length) {
      // 这里是判断在 编辑yaml界面，直接刷新页面，需要去拉取yaml
      nextProps.actions.resourceDetail.fetchResourceYaml.fetch();
    }
  }

  componentWillUnmount() {
    let { actions } = this.props;
    actions.resourceDetail.clearDetail();
    // 清除modifyResource workflow的信息
    actions.workflow.modifyResource.reset();
    // 清除applyResourceFlow workflow的信息
    actions.workflow.applyResource.reset();
  }

  render() {
    let { subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { mode } = subRoot;

    // 创建页面所需要展示的界面
    let content: JSX.Element;

    // 如果 模式为 modify，只提供 Yaml创建
    let resourceType = urlParams['resourceName'];
    let kind = urlParams['type'];
    let headTitle = resourceConfig()[resourceType] ? typeMapName[mode] + resourceConfig()[resourceType].headTitle : '';

    if (resourceType === 'svc' && mode === 'create') {
      content = <EditServicePanel />;
    } else if (kind === 'resource' && mode === 'create') {
      content = <EditResourceVisualizationPanel />;
      headTitle = typeMapName[mode] + 'Workload';
      // } else if (resourceType === 'ingress' && mode === 'create') {
      //   content = <EditIngressPanel />;
    } else if (resourceType === 'np' && mode === 'create') {
      content = <EditNamespacePanel />;
    } else if (resourceType === 'secret' && mode === 'create') {
      content = <EditSecretPanel />;
    } else if (resourceType === 'configmap' && mode === 'create') {
      content = <EditConfigMapPanel />;
    } else if (resourceType === 'node' && mode === 'create') {
      content = <CreateComputerPanel />;
      headTitle = t('添加节点');
    } else if (resourceType === 'lbcf' && mode === 'create') {
      content = <EditLbcfPanel />;
      headTitle = t('新建负载均衡');
    } else if (mode === 'apply') {
      content = this._editResourceYaml();
      headTitle = t('YAML创建资源');
    } else {
      content = this._editResourceYaml();
    }

    return (
      <div className="manage-area">
        <SubHeaderPanel headTitle={headTitle} />
        {content}
      </div>
    );
  }

  /** 一般资源，没有可视化创建的，通过Yaml来创建的 */
  _editResourceYaml() {
    let { subRoot } = this.props,
      { modifyResourceFlow, mode, resourceDetailState, applyResourceFlow } = subRoot,
      { yamlList } = resourceDetailState;

    let failed = modifyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyResourceFlow);
    let isNeedLoading =
      mode !== 'apply' &&
      mode !== 'create' &&
      (yamlList.fetched !== true || yamlList.fetchState === FetchState.Fetching);

    // 创建多种资源的错误判断
    let applyFailed = applyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(applyResourceFlow);

    return (
      <MainBodyLayout>
        {isNeedLoading ? (
          <FormLayout>
            <Icon type="loading" />
          </FormLayout>
        ) : (
          <FormLayout style={{ marginBottom: '50px' }}>
            <YamlEditorPanel config={this.state.config} handleInputForEditor={this._handleForInputEditor.bind(this)} />
            <ul className="form-list jiqun">
              <li className="pure-text-row fixed">
                <div className="form-input">
                  {(mode === 'create' || mode === 'modify') && (
                    <Button
                      className="mr10"
                      type="primary"
                      title={t('完成')}
                      disabled={modifyResourceFlow.operationState === OperationState.Performing}
                      onClick={() => {
                        this._handleSubmit();
                      }}
                    >
                      {t('完成')}
                    </Button>
                  )}
                  {mode === 'apply' && (
                    <Button
                      className="mr10"
                      type="primary"
                      title={t('完成')}
                      disabled={applyResourceFlow.operationState === OperationState.Performing}
                      onClick={() => {
                        this._handleSubmit();
                      }}
                    >
                      {applyFailed || failed ? t('重试') : t('完成')}
                    </Button>
                  )}
                  <Button title={t('取消')} onClick={this.goBack.bind(this)}>
                    {t('取消')}
                  </Button>
                  <TipInfo isShow={mode === 'apply' ? applyFailed : failed} type="error" isForm>
                    {mode === 'apply' ? getWorkflowError(applyResourceFlow) : getWorkflowError(modifyResourceFlow)}
                  </TipInfo>
                </div>
              </li>
            </ul>
          </FormLayout>
        )}
      </MainBodyLayout>
    );
  }

  /** 处理提交请求 */
  _handleSubmit() {
    let { actions, subRoot, namespaceSelection, region, route } = this.props,
      { resourceInfo, mode, resourceOption } = subRoot,
      { resourceSelection } = resourceOption;

    if (this.state.config !== '') {
      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        mode,
        namespace: namespaceSelection,
        yamlData: this.state.config,
        clusterId: route.queries['clusterId'],
        resourceIns: mode === 'modify' ? resourceSelection[0].metadata.name : ''
      };

      if (mode === 'apply') {
        actions.workflow.applyResource.start([resource], region.selection.value);
        actions.workflow.applyResource.perform();
      } else {
        actions.workflow.modifyResource.start([resource], region.selection.value);
        actions.workflow.modifyResource.perform();
      }
    }
  }

  goBack() {
    let { actions, route, namespaceSelection } = this.props,
      urlParam = router.resolve(route);
    // 回到列表处
    router.navigate(Object.assign({}, urlParam, { mode: 'list' }), route.queries);
  }

  /** yaml的编辑回显操作 */
  private _handleForInputEditor(config: string) {
    // 这里这么写的原因是为了解决 codemirrow 的光标的Bug
    this.setState({ config });
  }
}
