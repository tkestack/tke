import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { OperationState, isSuccessWorkflow, FetchState } from '@tencent/ff-redux';
import { Button } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';
import { CreateResource, ServicePorts, ServiceEditJSONYaml } from '../../../models';
import { FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { getWorkflowError, isEmpty } from '../../../../common/utils';
import { TipInfo } from '../../../../common/components';
import { EditServiceCommunicationPanel } from '../resourceEdition/EditServiceCommunicationPanel';
import { EditServicePortMapPanel } from '../resourceEdition/EditServicePortMapPanel';
import { router } from '../../../router';
import { validateServiceActions } from '../../../actions/validateServiceActions';
import { ServiceTypeMap } from '../resourceEdition/EditServicePanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { EditServiceAdvanceSettingPanel } from '../resourceEdition/EditServiceAdvanceSettingPanel';
import { SessionAffinity } from '../../../constants/Config';

/** 加载中的样式 */
const loadingElement = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);
interface UpdateServiceAccessTypePanelState {
  isOpenAdvancedSetting?: boolean;
}
let mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UpdateServiceAccessTypePanel extends React.Component<RootProps, UpdateServiceAccessTypePanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isOpenAdvancedSetting: false
    };
  }
  componentDidMount() {
    let { actions, subRoot, route } = this.props,
      { resourceList } = subRoot.resourceOption;

    // 这里是从列表页进入的时候，需要去初始化 serviceEdit当中的内容，如果是直接在当前页面刷新的话，会去拉取列表，在fetchResource之后，会初始化
    if (resourceList.data.recordCount) {
      let finder = resourceList.data.records.find(item => item.metadata.name === route.queries['resourceIns']);
      finder && actions.editSerivce.initServiceEditForUpdate(finder);
    }
  }

  componentWillUnmount() {
    let { actions } = this.props;
    actions.editSerivce.clearServiceEdit();
    actions.workflow.updateResourcePart.reset();
  }

  render() {
    let { actions, subRoot, route, cluster } = this.props,
      urlParams = router.resolve(route),
      { updateResourcePart, serviceEdit, resourceOption, isNeedExistedLb } = subRoot,
      { resourceList, resourceSelection } = resourceOption,
      { isOpenHeadless, communicationType, portsMap } = serviceEdit;

    let failed = updateResourcePart.operationState === OperationState.Done && !isSuccessWorkflow(updateResourcePart);

    return (
      <MainBodyLayout>
        <FormLayout>
          {resourceList.fetched !== true || resourceList.fetchState === FetchState.Fetching ? (
            loadingElement
          ) : (
            <div className="param-box server-update add">
              <ul className="form-list jiqun fixed-layout">
                <EditServiceCommunicationPanel
                  communicationType={communicationType}
                  communicationSelectAction={actions.editSerivce.selectCommunicationType}
                  isOpenHeadless={isOpenHeadless}
                  toggleHeadlessAction={actions.editSerivce.isOpenHeadless}
                  isDisabledChangeCommunicationType={
                    resourceSelection[0] && resourceSelection[0].spec.clusterIP === 'None' ? true : false
                  }
                  isDisabledToggleHeadless={true}
                />

                <EditServicePortMapPanel
                  addPortMap={actions.editSerivce.addPortMap}
                  deletePortMap={actions.editSerivce.deletePortMap}
                  communicationType={communicationType}
                  portsMap={portsMap}
                  updatePortMap={(obj: any, pId: string) => {
                    actions.editSerivce.updatePortMap(obj, pId);
                  }}
                  validatePortProtocol={(value: any, pId: string) => {
                    actions.validate.service.validatePortProtocol(value, pId);
                  }}
                  validateTargetPort={(value: any, pId: string) => {
                    actions.validate.service.validateTargetPort(value, pId);
                  }}
                  validateNodePort={(value: any, pId: string) => {
                    actions.validate.service.validateNodePort(value, pId);
                  }}
                  validateServicePort={(value: any, pId: string) => {
                    actions.validate.service.validateServicePort(value, pId);
                  }}
                />
                <li className="pure-text-row fixed">
                  <div className="form-input">
                    <Button
                      className="mr10"
                      type="primary"
                      disabled={updateResourcePart.operationState === OperationState.Performing}
                      onClick={this._handleSubmit.bind(this)}
                    >
                      {failed ? t('重试') : t('更新访问方式')}
                    </Button>
                    <Button
                      onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}
                    >
                      {t('取消')}
                    </Button>
                    <TipInfo
                      isShow={failed}
                      className="error"
                      style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px' }}
                    >
                      {getWorkflowError(updateResourcePart)}
                    </TipInfo>
                  </div>
                </li>
              </ul>
              <div>
                <hr className="hr-mod" />
                <EditServiceAdvanceSettingPanel
                  isShow={this.state.isOpenAdvancedSetting}
                  communicationType={communicationType}
                  validatesessionAffinityTimeout={actions.validate.service.validatesessionAffinityTimeout}
                  chooseChoosesessionAffinityMode={actions.editSerivce.chooseChoosesessionAffinityMode}
                  inputsessionAffinityTimeout={actions.editSerivce.inputsessionAffinityTimeout}
                  chooseExternalTrafficPolicyMode={actions.editSerivce.chooseExternalTrafficPolicyMode}
                  externalTrafficPolicy={serviceEdit.externalTrafficPolicy}
                  sessionAffinity={serviceEdit.sessionAffinity}
                  sessionAffinityTimeout={serviceEdit.sessionAffinityTimeout}
                  v_sessionAffinityTimeout={serviceEdit.v_sessionAffinityTimeout}
                />
                <a
                  href="javascript:;"
                  className="more-links-btn"
                  style={{ marginLeft: '-5px', marginBottom: '10px', display: 'inline-block' }}
                  onClick={e => this.setState({ isOpenAdvancedSetting: !this.state.isOpenAdvancedSetting })}
                >
                  <span style={{ verticalAlign: 'middle' }}>
                    {this.state.isOpenAdvancedSetting ? t('隐藏高级设置') : t('显示高级设置')}
                  </span>
                </a>
              </div>
            </div>
          )}
        </FormLayout>
      </MainBodyLayout>
    );
  }

  /** 处理提交请求 */
  private _handleSubmit() {
    let { actions, subRoot, route } = this.props,
      { resourceInfo, mode, serviceEdit, resourceOption } = subRoot,
      { resourceSelection } = resourceOption;

    actions.validate.service.validateUpdateServiceAccessEdit();

    if (validateServiceActions._validateUpdateServiceAccessEdit(serviceEdit)) {
      let {
        portsMap,
        communicationType,
        isOpenHeadless,
        externalTrafficPolicy,
        sessionAffinity,
        sessionAffinityTimeout
      } = serviceEdit;

      // 构建端口映射
      let ports = portsMap.map(port => {
        let tmp: ServicePorts = {
          name: port.protocol.toLocaleLowerCase() + '-' + port.targetPort + '-' + port.port,
          nodePort: communicationType !== 'ClusterIP' && port.nodePort ? +port.nodePort : null,
          port: +port.port,
          targetPort: +port.targetPort,
          protocol: port.protocol
        };

        return tmp;
      });

      // 构建创建service 的json的格式
      let jsonData: ServiceEditJSONYaml = {
        spec: {
          clusterIP: isOpenHeadless ? 'None' : undefined,
          type: ServiceTypeMap[communicationType],
          ports: ports,
          externalTrafficPolicy: communicationType !== 'ClusterIP' ? externalTrafficPolicy : undefined,
          sessionAffinity: sessionAffinity,
          sessionAffinityConfig:
            sessionAffinity === SessionAffinity.ClientIP
              ? {
                  clientIP: {
                    timeoutSeconds: +sessionAffinityTimeout
                  }
                }
              : undefined
        }
      };

      // 去除当中不需要的数据
      jsonData = JSON.parse(JSON.stringify(jsonData));

      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        mode,
        namespace: route.queries['np'],
        clusterId: route.queries['clusterId'],
        resourceIns: route.queries['resourceIns'],
        jsonData: JSON.stringify(jsonData),
        isStrategic: false
      };

      actions.workflow.updateResourcePart.start([resource], +route.queries['rid']);
      actions.workflow.updateResourcePart.perform();
    }
  }
}
