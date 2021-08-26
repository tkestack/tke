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

import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Button, Select } from '@tea/component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, SelectList, TipInfo } from '../../../../common/components';
import { FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { ResourceInfo } from '../../../../common/models';
import { getWorkflowError, isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validateServiceActions } from '../../../actions/validateServiceActions';
import { SessionAffinity } from '../../../constants/Config';
import { CreateResource, PortMap, ServiceEdit, ServiceEditJSONYaml, ServicePorts } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { EditServiceAdvanceSettingPanel } from './EditServiceAdvanceSettingPanel';
import { EditServiceCommunicationPanel } from './EditServiceCommunicationPanel';
import { EditServicePortMapPanel } from './EditServicePortMapPanel';
import { EditServiceWorkloadDialog } from './EditServiceWorkloadDialog';
import { reduceNs } from '../../../../../../helpers';

/** service YAML当中的type映射 */
export const ServiceTypeMap = {
  LoadBalancer: 'LoadBalancer',
  ClusterIP: 'ClusterIP',
  NodePort: 'NodePort',
  SvcLBTypeInner: 'LoadBalancer',
  ExternalName: 'ExternalName'
};

/**
 * 处理端口映射
 * @param ports: Portmap
 */
export const ReduceServicePorts = (portsMap: PortMap[], communicationType: string) => {
  let isNotClusterIP = communicationType !== 'ClusterIP';

  return portsMap.map(port => {
    let tmp: ServicePorts = {
      name: port.protocol.toLocaleLowerCase() + '-' + port.targetPort + '-' + port.port,
      nodePort: isNotClusterIP && port.nodePort ? +port.nodePort : undefined,
      port: +port.port,
      targetPort: +port.targetPort,
      protocol: port.protocol
    };

    return tmp;
  });
};

/**
 * 处理annotation的内容，pc内访问、购买lb带宽等，都放置在annotations里面实现
 * @param serviceEdit: ServiceEdit  服务的编辑信息
 * @param clusterId: string 当前的集群Id
 */
export const ReduceServiceAnnotations = (serviceEdit: ServiceEdit, clusterId: string) => {
  let { description } = serviceEdit;

  let annotations = {};

  if (description) {
    annotations['description'] = description;
  }

  return annotations;
};

/**
 * 处理Service的json格式
 * @params resourceInfo: ResourceInfo  当前资源的配置信息
 * @params ports: any[] 端口的信息
 * @params annotations: any 的信息
 * @params selectorObj: any
 * @params namespace: string 命名空间
 * @params communicationType: string 当前的访问方式
 * @params serviceName: string 服务的名称
 * @params isOpenHeadless: boolean 是否开启headless Service
 */
export const ReduceServiceJSONData = (dataObj: {
  resourceInfo: ResourceInfo;
  ports: any[];
  annotations: any;
  selectorObj: any;
  namespace: string;
  communicationType: string;
  serviceName: string;
  isOpenHeadless: boolean;
  sessionConfig: any;
}) => {
  let {
    resourceInfo,
    ports,
    annotations,
    selectorObj,
    namespace,
    communicationType,
    serviceName,
    isOpenHeadless,
    sessionConfig
  } = dataObj;

  let jsonData: ServiceEditJSONYaml = {
    kind: resourceInfo.headTitle,
    apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
    metadata: {
      name: serviceName,
      namespace: reduceNs(namespace),
      annotations: isEmpty(annotations) ? undefined : annotations
    },
    spec: {
      clusterIP: isOpenHeadless ? 'None' : undefined,
      type: ServiceTypeMap[communicationType],
      ports: ports,
      selector: isEmpty(selectorObj) ? undefined : selectorObj,
      externalTrafficPolicy: communicationType !== 'ClusterIP' ? sessionConfig.externalTrafficPolicy : undefined,
      sessionAffinity: sessionConfig.sessionAffinity,
      sessionAffinityConfig:
        sessionConfig.sessionAffinity === SessionAffinity.ClientIP
          ? {
              clientIP: {
                timeoutSeconds: +sessionConfig.sessionAffinityTimeout
              }
            }
          : undefined
    }
  };

  // 去除当中不需要的数据
  return JSON.parse(JSON.stringify(jsonData));
};

interface EditServicePanelState {
  isOpenAdvancedSetting?: boolean;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditServicePanel extends React.Component<RootProps, EditServicePanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isOpenAdvancedSetting: false
    };
  }
  componentWillUnmount() {
    let { actions } = this.props;
    actions.editSerivce.clearServiceEdit();
  }

  componentDidMount() {
    let { actions, route } = this.props;
    // 初始化namespace
    actions.editSerivce.selectNamespace(route.queries['np']);
  }

  render() {
    return (
      <MainBodyLayout>
        <FormLayout>
          <div className="param-box server-update add">
            {this._renderBasicInfo()}
            {this._renderServiceSetting()}
            {this._renderAdvancePanel()}
            {this._renderBindWorkload()}
          </div>
        </FormLayout>

        <EditServiceWorkloadDialog />
      </MainBodyLayout>
    );
  }

  /** 基本信息的填写区域 */
  private _renderBasicInfo() {
    let { actions, subRoot, namespaceList } = this.props,
      { serviceEdit } = subRoot;

    let namespaceOptions = namespaceList.data.records.map(item => ({
      value: item.name,
      text: item.displayName
    }));

    return (
      <div>
        <div className="param-hd">
          <h3>{t('基本信息')}</h3>
        </div>
        <div className="param-bd">
          <ul className="form-list fixed-layout">
            <FormItem label={t('服务名称')}>
              <div className={classnames('form-unit', { 'is-error': serviceEdit.v_serviceName.status === 2 })}>
                <InputField
                  type="text"
                  placeholder={t('请输入服务名称')}
                  tipMode="popup"
                  validator={serviceEdit.v_serviceName}
                  value={serviceEdit.serviceName}
                  tip={t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
                  onChange={actions.editSerivce.inputServiceName}
                  onBlur={actions.validate.service.validateServiceName}
                />
              </div>
            </FormItem>
            <FormItem label={t('描述')}>
              <div className={classnames('form-unit', { 'is-error': serviceEdit.v_description.status === 2 })}>
                <InputField
                  type="textarea"
                  placeholder={t('请输入描述信息，不超过1000个字符')}
                  tipMode="popup"
                  validator={serviceEdit.v_description}
                  value={serviceEdit.description}
                  onChange={actions.editSerivce.inputServiceDesp}
                  onBlur={actions.validate.service.validateServiceDesp}
                />
              </div>
            </FormItem>
            <FormItem label={t('命名空间')}>
              <div className={classnames('form-unit', { 'is-error': serviceEdit.v_namespace.status === 2 })}>
                <Select
                  size="m"
                  options={namespaceOptions}
                  value={serviceEdit.namespace}
                  onChange={value => {
                    actions.editSerivce.selectNamespace(value);
                    actions.namespace.selectNamespace(value);
                  }}
                />
              </div>
            </FormItem>
          </ul>
        </div>
      </div>
    );
  }

  /** Service设置信息的填写 */
  private _renderServiceSetting() {
    let { subRoot, actions, cluster } = this.props,
      { serviceEdit, isNeedExistedLb } = subRoot,
      { communicationType, portsMap, isOpenHeadless } = serviceEdit;

    return (
      <div>
        <hr className="hr-mod" />
        <div className="param-hd">
          <h3>{t('访问设置(Service)')}</h3>
        </div>
        <div className="param-bd">
          <ul className="form-list fixed-layout">
            <EditServiceCommunicationPanel
              communicationType={communicationType}
              communicationSelectAction={actions.editSerivce.selectCommunicationType}
              isOpenHeadless={isOpenHeadless}
              toggleHeadlessAction={actions.editSerivce.isOpenHeadless}
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
            {/* <EditServiceAdvanceSetting> */}
          </ul>
        </div>
      </div>
    );
  }

  private _renderAdvancePanel() {
    let { actions, subRoot } = this.props,
      { serviceEdit } = subRoot;
    return (
      <div>
        {this.state.isOpenAdvancedSetting && (
          <div className="param-hd">
            <h3>
              <Trans>
                高级设置<span className="text-label">（选填）</span>
              </Trans>
            </h3>
          </div>
        )}
        <div className="param-bd" style={{ marginBottom: '0px' }}>
          <EditServiceAdvanceSettingPanel
            communicationType={serviceEdit.communicationType}
            isShow={this.state.isOpenAdvancedSetting}
            validatesessionAffinityTimeout={actions.validate.service.validatesessionAffinityTimeout}
            chooseChoosesessionAffinityMode={actions.editSerivce.chooseChoosesessionAffinityMode}
            inputsessionAffinityTimeout={actions.editSerivce.inputsessionAffinityTimeout}
            chooseExternalTrafficPolicyMode={actions.editSerivce.chooseExternalTrafficPolicyMode}
            externalTrafficPolicy={serviceEdit.externalTrafficPolicy}
            sessionAffinity={serviceEdit.sessionAffinity}
            sessionAffinityTimeout={serviceEdit.sessionAffinityTimeout}
            v_sessionAffinityTimeout={serviceEdit.v_sessionAffinityTimeout}
          />
        </div>
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
    );
  }

  /** Service的 deployment绑定 */
  private _renderBindWorkload() {
    let { actions, subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { modifyResourceFlow, serviceEdit } = subRoot,
      { selector } = serviceEdit;

    let failed = modifyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyResourceFlow);

    return (
      <div>
        <hr className="hr-mod" />
        <div className="param-hd">
          <h3>
            {t('Workload绑定')}
            <span className="text-label">{t('（选填）')}</span>
          </h3>
        </div>
        <div className="param-bd" style={{ paddingBottom: '25px' }}>
          <ul className="form-list fixed-layout jiqun">
            <FormItem label="Selectors">
              {selector.map((s, index) => {
                return (
                  <div key={index} className="form-unit" style={{ marginBottom: '5px' }}>
                    <div className={s.v_key.status === 2 ? 'is-error' : ''} style={{ display: 'inline-block' }}>
                      <Bubble placement="bottom" content={s.v_key.status === 2 ? s.v_key.message : null}>
                        <input
                          type="text"
                          className="tc-15-input-text m"
                          placeholder="key"
                          value={s.key}
                          onChange={e => actions.editSerivce.updateSelectorConfig({ key: e.target.value }, s.id + '')}
                          onBlur={e => {
                            actions.validate.service.validateSelectorContent({ key: e.target.value }, s.id + '');
                          }}
                        />
                      </Bubble>
                    </div>
                    <span style={{ margin: '0 5px' }} className="inline-help-text">
                      =
                    </span>
                    <div
                      className={s.v_value.status === 2 ? 'is-error' : ''}
                      style={{ display: 'inline-block', marginRight: '10px' }}
                    >
                      <Bubble placement="bottom" content={s.v_value.status === 2 ? s.v_value.message : null}>
                        <input
                          type="text"
                          className="tc-15-input-text m"
                          placeholder="value"
                          value={s.value}
                          onChange={e => {
                            actions.editSerivce.updateSelectorConfig({ value: e.target.value }, s.id + '');
                          }}
                          onBlur={e => {
                            actions.validate.service.validateSelectorContent({ value: e.target.value }, s.id + '');
                          }}
                        />
                      </Bubble>
                    </div>
                    <a
                      href="javascript:;"
                      onClick={() => {
                        actions.editSerivce.deleteSelectorContent(s.id + '');
                      }}
                    >
                      <i className="icon-cancel-icon" />
                    </a>
                  </div>
                );
              })}
              <a href="javascript:;" className="more-links-btn" onClick={actions.editSerivce.addSelector}>
                {t('添加')}
              </a>
              <span style={{ verticalAlign: '1px' }}> | </span>
              <a
                href="javascript:;"
                className="more-links-btn"
                onClick={actions.editSerivce.workload.toggleIsShowWorkloadDialog}
              >
                {t('引用Workload')}
              </a>
            </FormItem>
            <li className="pure-text-row fixed">
              <div className="form-input">
                <Button
                  className="mr10"
                  type="primary"
                  disabled={modifyResourceFlow.operationState === OperationState.Performing}
                  onClick={this._handleSubmit.bind(this)}
                >
                  {failed ? t('重试') : t('创建服务')}
                </Button>
                <Button onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}>
                  {t('取消')}
                </Button>
                <TipInfo
                  isShow={failed}
                  className="error"
                  style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
                >
                  {getWorkflowError(modifyResourceFlow)}
                </TipInfo>
              </div>
            </li>
          </ul>
        </div>
      </div>
    );
  }

  /** 处理提交请求 */
  private _handleSubmit() {
    let { actions, subRoot, route, region } = this.props,
      { resourceInfo, mode, serviceEdit } = subRoot;

    actions.validate.service.validateServiceEdit();

    if (validateServiceActions._validateServiceEdit(serviceEdit)) {
      let { portsMap, communicationType, selector, namespace, serviceName, isOpenHeadless } = serviceEdit;

      // 构建端口映射
      let ports = ReduceServicePorts(portsMap, communicationType);

      // vpc内访问、购买lb带宽等，都放置在annotations里面实现
      let annotations = ReduceServiceAnnotations(serviceEdit, route.queries['clusterId']);

      // selector
      let selectorObj = {};
      if (selector.length) {
        selector.forEach(s => {
          selectorObj[s.key] = s.value;
        });
      }

      let sessionConfig = {
        externalTrafficPolicy: serviceEdit.externalTrafficPolicy,
        sessionAffinity: serviceEdit.sessionAffinity,
        sessionAffinityTimeout: serviceEdit.sessionAffinityTimeout
      };
      // 构建创建service 的json的格式
      let jsonData: ServiceEditJSONYaml = ReduceServiceJSONData({
        resourceInfo,
        ports,
        annotations,
        selectorObj,
        namespace,
        communicationType,
        serviceName,
        isOpenHeadless,
        sessionConfig
      });

      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        mode,
        namespace: namespace,
        clusterId: route.queries['clusterId'],
        jsonData: JSON.stringify(jsonData)
      };

      actions.workflow.modifyResource.start([resource], region.selection.value);
      actions.workflow.modifyResource.perform();
    }
  }
}
