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
import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Radio, Select, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { OperationState, bindActionCreators, insertCSS, isSuccessWorkflow, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../../config/resourceConfig';
import { reduceNs } from '../../../../../../helpers';
import { FormItem, InputField, TipInfo } from '../../../../common/components';
import { FixedFormLayout, FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { ResourceInfo } from '../../../../common/models';
import { getWorkflowError, isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validateWorkloadActions } from '../../../actions/validateWorkloadActions';
import {
  NodeAbnormalStrategy,
  ResourceTypeList,
  RestartPolicyTypeList,
  WorkloadNetworkTypeEnum,
  affinityType
} from '../../../constants/Config';
import {
  Computer,
  ContainerEnv,
  ContainerItem,
  CreateResource,
  DifferentInterfaceResourceOperation,
  HealthCheckItem,
  HpaEditJSONYaml,
  MetricOption,
  ServiceEditJSONYaml,
  VolumeItem,
  WorkloadEditJSONYaml
} from '../../../models';
import { AffinityRule } from '../../../models/WorkloadEdit';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { EditResourceAdvancedPanel } from './EditResourceAdvancedPanel';
import { EditResourceContainerNumPanel } from './EditResourceContainerNumPanel';
import { EditResourceContainerPanel } from './EditResourceContainerPanel';
import { EditResourceLabelPanel } from './EditResourceLabelPanel';
import { EditResourceVolumePanel } from './EditResourceVolumePanel';
import { EditServiceAdvanceSettingPanel } from './EditServiceAdvanceSettingPanel';
import { EditServiceCommunicationPanel } from './EditServiceCommunicationPanel';
import { ReduceServiceAnnotations, ReduceServiceJSONData, ReduceServicePorts } from './EditServicePanel';
import { EditServicePortMapPanel } from './EditServicePortMapPanel';
import { ResourceEditHostPathDialog } from './ResourceEditHostPathDialog';
import { ResourceSelectConfigDialog } from './ResourceSelectConfigDialog';

/** service YAML当中的type映射 */
const serviceTypeMap = {
  LoadBalancer: 'LoadBalancer',
  ClusterIP: 'ClusterIP',
  NodePort: 'NodePort',
  SvcLBTypeInner: 'LoadBalancer',
  ExternalName: 'ExternalName'
};

insertCSS(
  'EditResourceVisualizationPanel',
  `
.specific .tc-15-radio-wrap{
    margin-left: 0px;
    display: block;
    margin-bottom: 12px;
}
.specific .tc-15-radio-wrap:last-child{
    margin-bottom: 0;
}
`
);

interface EditResourceVisualizationPanelState {
  /** isOpenAdvanced */
  isOpenAdvancedSetting?: boolean;

  /**service访问设置高级设置 */
  isOpenServiceAdvancedSetting?: boolean;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceVisualizationPanel extends React.Component<RootProps, EditResourceVisualizationPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isOpenAdvancedSetting: false,
      isOpenServiceAdvancedSetting: false
    };
  }

  componentDidMount() {
    let { actions, route, cluster } = this.props,
      urlParams = router.resolve(route);

    // 初始化namespace
    actions.editWorkload.selectNamespace(route.queries['np']);

    // 初始化新建workload的类型
    actions.editWorkload.selectResourceType(urlParams['resourceName']);

    // 判断使用是否可以使用gpu
    actions.editWorkload.isCanUseGpu();

    //判断是否可显示tapp
    actions.editWorkload.isCanUseTapp();

    //初始化是否有超售比
    if (cluster.selection && cluster.selection.spec.properties && cluster.selection.spec.properties.oversoldRatio) {
      const oversoldRatio = cluster.selection.spec.properties.oversoldRatio;
      actions.editWorkload.initOversoldRatio(oversoldRatio);
    }
  }

  //在当前页面直接刷新时需要判断超售比前后状态
  componentWillReceiveProps(nextProps: RootProps) {
    const { actions, cluster } = nextProps;
    if (nextProps.cluster.selection && !this.props.cluster.selection) {
      if (cluster.selection.spec.properties && cluster.selection.spec.properties.oversoldRatio) {
        const oversoldRatio = cluster.selection.spec.properties.oversoldRatio;
        actions.editWorkload.initOversoldRatio(oversoldRatio);
      }
    }
  }

  componentWillUnmount() {
    const { actions } = this.props;
    // 清除workloadEdit当中的所有数据，避免缓存
    actions.editWorkload.clearWorkloadEdit();
    // 如果同时创建Service，需要清空ServiceEdit当中的信息
    actions.editSerivce.clearServiceEdit();
    // 清除两个flow的设置
    actions.workflow.modifyResource.reset();
    actions.workflow.applyResource.reset();
  }

  render() {
    const { actions, subRoot, namespaceList, route, cluster } = this.props,
      urlParams = router.resolve(route),
      { workloadEdit, modifyResourceFlow, applyResourceFlow, serviceEdit, isNeedExistedLb, addons } = subRoot,
      {
        isCreateService,
        v_workloadName,
        workloadName,
        description,
        v_description,
        namespace,
        v_namespace,
        workloadType,
        cronSchedule,
        v_cronSchedule,
        completion,
        v_completion,
        parallelism,
        v_parallelism,
        restartPolicy,
        imagePullSecrets,
        configEdit,
        nodeAbnormalMigratePolicy,
        isCanUseTapp
      } = workloadEdit,
      { secretList } = configEdit;

    // 是否开启高级设置
    const isOpenAdvanced = this.state.isOpenAdvancedSetting;

    /** 渲染 重启策略列表 */
    const restartOptions = RestartPolicyTypeList.map((item, index) => (
      <option key={index} value={item.value}>
        {item.label}
      </option>
    ));

    /**
     * 渲染imagePullSecret的列表
     * pre: secret的类型为 kubernetes.io/dockercfg
     */
    const finalSecretList = secretList.data.records.filter(item => item.type === 'kubernetes.io/dockercfg');
    const secretListOptions = finalSecretList.map((item, index) => (
      <option key={index} value={item.metadata.name}>
        {item.metadata.name}
      </option>
    ));
    secretListOptions.unshift(
      <option key={uuid()} value="">
        {t('请选择dockercfg类型的Secret')}
      </option>
    );

    const failed =
      (modifyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyResourceFlow)) ||
      (applyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(applyResourceFlow));

    // 判断是否deployment 或者 statefulset
    const isDeploymentOrStateful =
      workloadType === 'deployment' || workloadType === 'statefulset' || workloadType === 'tapp';

    const namespaceOptions = namespaceList.data.records.map(item => ({
      value: item.name,
      text: item.displayName
    }));

    const finalResourceTypeList = [];
    ResourceTypeList.forEach(list => {
      if (list.value !== 'tapp' || isCanUseTapp) {
        finalResourceTypeList.push(list);
      } else if (list.value === 'tapp' && addons['TappController'] !== undefined) {
        finalResourceTypeList.push(list);
      }
    });

    return (
      <MainBodyLayout>
        <FormLayout>
          <div className="param-box server-update add">
            <ul
              className="form-list jiqun fixed-layout"
              style={isDeploymentOrStateful ? {} : { paddingBottom: '50px' }}
            >
              <FormItem label={t('工作负载名')}>
                <InputField
                  type="text"
                  placeholder={t('请输入Workload名称')}
                  tipMode="popup"
                  validator={v_workloadName}
                  value={workloadName}
                  tip={t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
                  onChange={actions.editWorkload.inputWorkloadName}
                  onBlur={actions.validate.workload.validateWorkloadName}
                />
              </FormItem>
              <FormItem label={t('描述')}>
                <InputField
                  type="textarea"
                  placeholder={t('请输入描述信息，不超过1000个字符')}
                  tipMode="popup"
                  validator={v_description}
                  value={description}
                  onChange={actions.editWorkload.inputWorkloadDesp}
                  onBlur={actions.validate.workload.validateWorkloadDesp}
                />
              </FormItem>

              <EditResourceLabelPanel />
              <FormItem label={t('集群')}>
                <FormPanel.InlineText style={{ marginLeft: 0 }}>
                  {cluster.selection ? cluster.selection.metadata.name : route.queries['clusterId']}
                </FormPanel.InlineText>
              </FormItem>
              <FormItem label={t('命名空间')}>
                <Select
                  size="m"
                  options={namespaceOptions}
                  value={namespace}
                  onChange={value => {
                    actions.editWorkload.selectNamespace(value);
                    actions.namespace.selectNamespace(value);
                  }}
                />
              </FormItem>
              <FormItem label={t('类型')}>
                <div
                  className="form-unit specific"
                  style={{ backgroundColor: '#f2f2f2', padding: '15px 10px 10px', maxWidth: '350px' }}
                >
                  <Radio.Group value={workloadType} onChange={value => this._handleResourceTypeSelect(value)}>
                    {finalResourceTypeList.map((item, rIndex) => {
                      return (
                        <Radio key={rIndex} name={item.value}>
                          {item.label}
                        </Radio>
                      );
                    })}
                  </Radio.Group>
                </div>
              </FormItem>
              <FormItem label={t('节点异常策略')} isShow={workloadType === 'tapp'}>
                <Select
                  size="m"
                  options={NodeAbnormalStrategy}
                  value={nodeAbnormalMigratePolicy}
                  onChange={value => {
                    actions.editWorkload.selectNodeAbnormalMigratePolicy(value);
                  }}
                />
                <Text parent="p" theme="label" style={{ marginTop: '8px' }}>
                  {nodeAbnormalMigratePolicy === 'true'
                    ? t('迁移，调度策略与Deployment一致，Pod会迁移到新的节点')
                    : t('不迁移，调度策略与StatefulSel一致，异常pod不会被迁移')}
                </Text>
              </FormItem>
              <FormItem label={t('执行策略')} isShow={workloadType === 'cronjob'}>
                <div className={classnames('form-unit', { 'is-error': v_cronSchedule.status === 2 })}>
                  <InputField
                    type="text"
                    style={{ width: '340px' }}
                    placeholder={t('请输入执行策略，如: 0 0 2 1 *')}
                    tipMode="popup"
                    validator={v_cronSchedule}
                    value={cronSchedule}
                    onChange={actions.editWorkload.inputCronjobSchedule}
                    onBlur={actions.validate.workload.validateCronSchedule}
                  />
                </div>
              </FormItem>
              <FormItem label={t('Job设置')} isShow={workloadType === 'job' || workloadType === 'cronjob'}>
                <FixedFormLayout style={{ width: '310px', paddingTop: '5px' }}>
                  <FormItem label={t('重复执行次数')} tips={t('该Job下的Pod需要重复执行次数')}>
                    <InputField
                      type="text"
                      style={{ width: '150px' }}
                      tipMode="popup"
                      validator={v_completion}
                      value={completion}
                      onChange={actions.editWorkload.inputJobCompletion}
                      onBlur={actions.validate.workload.validateJobCompletion}
                    />
                  </FormItem>
                  <FormItem label={t('并行度')} tips={t('该Job下Pod并行执行的数量')}>
                    <InputField
                      type="text"
                      style={{ width: '150px' }}
                      tipMode="popup"
                      validator={v_parallelism}
                      value={parallelism}
                      onChange={actions.editWorkload.inputJobParallelism}
                      onBlur={actions.validate.workload.validateJobParallel}
                    />
                  </FormItem>
                  <FormItem
                    label={t('失败重启策略')}
                    tips={t(
                      'Pod下容器异常推出后的重启策略， Never：不重启容器，直至Pod下所有容器退出; OnFailure : Pod继续运行，容器将重新启动'
                    )}
                  >
                    <select
                      className="tc-15-select m"
                      style={{ marginRight: '6px', minWidth: '150px' }}
                      value={restartPolicy}
                      onChange={e => {
                        actions.editWorkload.selectRestartPolicy(e.target.value);
                      }}
                    >
                      {restartOptions}
                    </select>
                  </FormItem>
                </FixedFormLayout>
              </FormItem>

              <EditResourceVolumePanel />

              <EditResourceContainerPanel />

              <EditResourceContainerNumPanel />

              <EditResourceAdvancedPanel isOpenAdvanced={isOpenAdvanced} />

              <a
                href="javascript:;"
                className="more-links-btn"
                onClick={() => {
                  this.setState({ isOpenAdvancedSetting: !isOpenAdvanced });
                  /**默认选择指定节点*/
                  if (!isOpenAdvanced === false) {
                    actions.editWorkload.selectNodeSelectType(affinityType.unset);
                  }
                }}
              >
                {isOpenAdvanced ? t('隐藏高级设置') : t('显示高级设置')}
              </a>

              <li className="pure-text-row fixed">
                <div className="form-input">
                  <Button
                    className="mr10"
                    type="primary"
                    disabled={
                      modifyResourceFlow.operationState === OperationState.Performing ||
                      applyResourceFlow.operationState === OperationState.Performing
                    }
                    onClick={this._handleSubmit.bind(this)}
                  >
                    {failed ? t('重试') : t('创建Workload')}
                  </Button>
                  <Button onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}>
                    {t('取消')}
                  </Button>
                  <TipInfo isShow={failed} type="error" isForm>
                    {getWorkflowError(modifyResourceFlow) || getWorkflowError(applyResourceFlow)}
                  </TipInfo>
                </div>
              </li>
            </ul>

            {isDeploymentOrStateful && (
              <div style={{ paddingBottom: '50px' }}>
                <hr className="hr-mod" />
                <div className="param-hd">
                  <h3>{t('访问设置（Service）')}</h3>
                </div>
                <div className="param-bd">
                  <ul className="form-list fixed-layout jiqun">
                    <FormItem label="Service">
                      <div className="form-unit">
                        <label className="form-ctrl-label">
                          <input
                            type="checkbox"
                            className="tc-15-checkbox"
                            checked={isCreateService}
                            style={{ verticalAlign: 'middle' }}
                            onChange={() => actions.editWorkload.isCreateService(isCreateService)}
                          />
                          {t('启用')}
                        </label>
                      </div>
                    </FormItem>

                    <EditServiceCommunicationPanel
                      communicationType={serviceEdit.communicationType}
                      communicationSelectAction={actions.editSerivce.selectCommunicationType}
                      isOpenHeadless={serviceEdit.isOpenHeadless}
                      toggleHeadlessAction={actions.editSerivce.isOpenHeadless}
                      isShow={isCreateService}
                    />

                    <EditServicePortMapPanel
                      addPortMap={actions.editSerivce.addPortMap}
                      deletePortMap={actions.editSerivce.deletePortMap}
                      communicationType={serviceEdit.communicationType}
                      portsMap={serviceEdit.portsMap}
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
                      isShow={isCreateService}
                    />
                  </ul>
                  <div>
                    <EditServiceAdvanceSettingPanel
                      communicationType={serviceEdit.communicationType}
                      isShow={this.state.isOpenServiceAdvancedSetting}
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
                      onClick={e =>
                        this.setState({ isOpenServiceAdvancedSetting: !this.state.isOpenServiceAdvancedSetting })
                      }
                    >
                      <span style={{ verticalAlign: 'middle' }}>
                        {this.state.isOpenServiceAdvancedSetting ? t('隐藏高级设置') : t('显示高级设置')}
                      </span>
                    </a>
                  </div>
                </div>
              </div>
            )}
          </div>
        </FormLayout>

        <ResourceSelectConfigDialog />
        <ResourceEditHostPathDialog />
      </MainBodyLayout>
    );
  }

  /** 生成 workload类型的radio列表 */
  private _handleResourceTypeSelect(resourceType: string) {
    const { actions } = this.props;
    actions.editWorkload.selectResourceType(resourceType);
  }

  /** 处理提交请求 */
  /* eslint-disable */
  private _handleSubmit() {
    let { actions, subRoot, route, region, clusterVersion } = this.props,
      { mode, workloadEdit, serviceEdit } = subRoot;

    actions.validate.workload.validateWorkloadEdit();

    if (validateWorkloadActions._validateWorkloadEdit(workloadEdit, serviceEdit)) {
      let {
        isCreateService,
        workloadType,
        scaleType,
        minReplicas,
        workloadLabels,
        workloadName,
        containerNum,
        namespace,
        restartPolicy,
        completion,
        parallelism,
        volumes,
        description,
        containers,
        isNeedContainerNum,
        cronSchedule,
        imagePullSecrets,
        nodeAffinityType,
        nodeAffinityRule,
        computer,
        workloadAnnotations,
        nodeAbnormalMigratePolicy,
        networkType,
        floatingIPReleasePolicy,
        oversoldRatio,
        isOpenCronHpa
      } = workloadEdit;

      // 当前该资源的具体配置
      let workloadResourceInfo: ResourceInfo = resourceConfig(clusterVersion)[workloadType];

      let finalVolumes = volumes;
      let volumesInfo = this._reduceVolumes(finalVolumes);

      // 进行容器的相关数据拼接
      let containersInfo = this._reduceContainers(containers, volumes, { oversoldRatio, networkType });

      // 进行容器的labels的数据拼接，默认有一个 qcloud-app: workload的名称，很懂监控等都用qcloud-app的标签
      let labelsInfo = { 'qcloud-app': workloadName };
      workloadLabels.forEach(label => {
        labelsInfo[label.labelKey] = label.labelValue;
      });

      // selector的相关配置信息
      let selectorContent = {
        matchLabels: labelsInfo
      };

      // 描述信息放在annotaitions
      let annotations = {};
      if (description) {
        annotations['description'] = description;
      }

      // 判断当前的工作负载类型
      let isCronJobs = workloadType === 'cronjob',
        isJobs = workloadType === 'job',
        isCronJobOrCronJob = isCronJobs || isJobs,
        isStatefulset = workloadType === 'statefulset',
        isDeployment = workloadType === 'deployment',
        isTapp = workloadType === 'tapp';

      // 判断当前的实例数量是否为hpa类型
      let isAutoScale = scaleType === 'autoScale';

      // spec当中的 restartPolicy，job || cronjob的重启策略不能为 always，给用户选择他的重启策略
      let finalRestartPolicy = isCronJobOrCronJob ? restartPolicy : 'Always';

      // node亲和性调度的相关信息
      let affinityInfo =
        nodeAffinityType !== affinityType.unset
          ? this._reduceNodeAffinityInfo(nodeAffinityType, nodeAffinityRule, computer.selections)
          : '';

      // 如果选择了网络模式，需要把网络模式写在annotations当中
      let templateAnnotations = {};

      //将annotation赋值到template中
      if (workloadAnnotations.length) {
        workloadAnnotations.forEach(annotation => {
          annotations[annotation.labelKey] = annotation.labelValue;
          templateAnnotations[annotation.labelKey] = annotation.labelValue;
        });
      }
      if (networkType) {
        if (networkType === WorkloadNetworkTypeEnum.Nat || networkType === WorkloadNetworkTypeEnum.Overlay) {
          templateAnnotations['k8s.v1.cni.cncf.io/networks'] = 'galaxy-flannel';
        } else if (networkType === WorkloadNetworkTypeEnum.FloatingIP) {
          templateAnnotations['k8s.v1.cni.cncf.io/networks'] = 'galaxy-k8s-vlan';
          templateAnnotations['k8s.v1.cni.galaxy.io/release-policy'] =
            floatingIPReleasePolicy === 'always' ? '' : floatingIPReleasePolicy;
        }
      }

      // 日志目录信息放在 annotations 中
      let logInfo = containers.reduce((prev, item) => {
        if (isEmpty(item.logPath)) {
          return prev;
        }
        prev[item.name] = item.logDir + item.logPath;
        return prev;
      }, {});
      if (!isEmpty(logInfo)) {
        templateAnnotations['log.tke.cloud.tencent.com/log-dir'] = JSON.stringify(logInfo);
      }

      // template的内容，因为cronJob是放在 jobTemplate当中
      let templateContent = {
        metadata: {
          labels: labelsInfo,
          annotations: isEmpty(templateAnnotations) ? undefined : templateAnnotations
        },
        spec: {
          volumes: volumesInfo.length ? volumesInfo : undefined,
          containers: containersInfo,
          restartPolicy: finalRestartPolicy,
          imagePullSecrets: imagePullSecrets.length
            ? imagePullSecrets.map(item => ({
                name: item.secretName
              }))
            : undefined,
          affinity: affinityInfo ? affinityInfo : undefined,
          hostNetwork: networkType === WorkloadNetworkTypeEnum.Host ? true : undefined
        }
      };

      // cronjob的独有的配置
      let jobTemplateContent = {
        spec: {
          template: templateContent,
          completions: +completion,
          parallelism: +parallelism
        }
      };

      // 构建创建workload的json的格式，model的定义 https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#jobspec-v1-batch
      let jsonData: WorkloadEditJSONYaml = {
        kind: workloadResourceInfo.headTitle,
        apiVersion: (workloadResourceInfo.group ? workloadResourceInfo.group + '/' : '') + workloadResourceInfo.version,
        metadata: {
          name: workloadName,
          namespace: reduceNs(namespace),
          labels: labelsInfo,
          annotations: isEmpty(annotations) ? undefined : annotations
        },
        spec: {
          replicas: isNeedContainerNum ? (isAutoScale ? +minReplicas : +containerNum) : undefined,
          serviceName: isStatefulset && serviceEdit.isOpenHeadless ? workloadName : undefined,
          schedule: isCronJobs ? cronSchedule : undefined,
          template: !isCronJobs ? templateContent : undefined,
          jobTemplate: isCronJobs ? jobTemplateContent : undefined,
          selector: !isCronJobOrCronJob ? selectorContent : undefined,
          completions: isJobs ? +completion : undefined,
          parallelism: isJobs ? +parallelism : undefined,
          forceDeletePod: isTapp ? (nodeAbnormalMigratePolicy === 'true' ? true : false) : undefined
        }
      };

      /**
       * ========================== 此处是同时创建Service ==========================
       * pre: deployment || statefulset || tapp
       */
      let serviceJsonData =
        isCreateService && (isDeployment || isStatefulset || isTapp)
          ? JSON.stringify(this._reduceServiceData(labelsInfo))
          : '';

      /**
       * ========================== 此处是同时创建hpa ==========================
       * pre: deployment || statefulset || tapp
       */
      let hpaJsonData =
        (isDeployment || isTapp || isStatefulset) && isAutoScale ? JSON.stringify(this._reduceHpaData()) : '';

      /**
       *  ========================== 此处是同时创建cronhpa ==========================
       * pre: deployment || statefulset || tapp
       */
      let cronhpaJsonData =
        (isDeployment || isTapp || isStatefulset) && isOpenCronHpa ? JSON.stringify(this._reduceCronHpaData()) : '';

      /** 最终传过去的json的数据 */
      let finalJSON = serviceJsonData + hpaJsonData + cronhpaJsonData + JSON.stringify(jsonData);

      let resource: CreateResource = {
        id: uuid(),
        resourceInfo: workloadResourceInfo,
        mode,
        namespace: namespace,
        clusterId: route.queries['clusterId'],
        jsonData: finalJSON
      };
      /**Tapp的请求需要创建多个不同接口的资源,所以需要进行特殊处理 */
      if (isTapp) {
        let resources: CreateResource[] = [];
        let differentInterfaceResourceOperation: DifferentInterfaceResourceOperation[] = [];
        if (isAutoScale) {
          /**创建hpa资源 */
          resources.push({
            id: uuid(),
            resourceInfo: resourceConfig(clusterVersion)['hpa'],
            mode,
            namespace: namespace,
            clusterId: route.queries['clusterId'],
            jsonData: hpaJsonData
          });
          /**不需要operation 所以传入的值为{},需要传值 */
          differentInterfaceResourceOperation.push({});
        }
        /** 创建Service */
        if (serviceJsonData) {
          resources.push({
            id: uuid(),
            resourceInfo: resourceConfig(clusterVersion)['svc'],
            mode,
            namespace,
            clusterId: route.queries['clusterId'],
            jsonData: serviceJsonData
          });
        }
        /** 创建Cronhpa资源 */
        if (cronhpaJsonData) {
          resources.push({
            id: uuid(),
            resourceInfo: resourceConfig(clusterVersion)['cronhpa'],
            mode,
            namespace,
            clusterId: route.queries['clusterId'],
            jsonData: cronhpaJsonData
          });
        }

        /**创建TAPP资源 */
        resources.push({
          id: uuid(),
          resourceInfo: workloadResourceInfo,
          mode,
          namespace: namespace,
          clusterId: route.queries['clusterId'],
          jsonData: JSON.stringify(jsonData)
        });

        differentInterfaceResourceOperation.push({});
        actions.workflow.applyDifferentInterfaceResource.start(resources, differentInterfaceResourceOperation);
        actions.workflow.applyDifferentInterfaceResource.perform();
      } else {
        if ((isDeployment || isStatefulset) && (isCreateService || isAutoScale || isOpenCronHpa)) {
          actions.workflow.applyResource.start([resource], region.selection.value);
          actions.workflow.applyResource.perform();
        } else {
          actions.workflow.modifyResource.start([resource], region.selection.value);
          actions.workflow.modifyResource.perform();
        }
      }
    }
  }
  /* eslint-enable */

  /** 处理cronhpa的相关信息 */
  private _reduceCronHpaData() {
    const { subRoot, clusterVersion } = this.props,
      { cronMetrics, workloadType, workloadName, namespace } = subRoot.workloadEdit;

    const resourceInfo = resourceConfig(clusterVersion)['cronhpa'],
      ResourceInfo = resourceConfig(clusterVersion)[workloadType];

    const isTapp = workloadType === 'tapp' ? true : false;

    const jsonData = {
      kind: resourceInfo.headTitle,
      apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
      metadata: {
        name: workloadName,
        namespace: reduceNs(namespace)
      },
      spec: {
        scaleTargetRef: {
          apiVersion: isTapp ? 'apps.tkestack.io/v1' : ResourceInfo.group + '/' + ResourceInfo.version,
          kind: ResourceInfo.headTitle,
          name: workloadName
        },
        crons: cronMetrics.map(metric => {
          return {
            schedule: metric.crontab,
            targetReplicas: +metric.targetReplicas
          };
        })
      }
    };

    return JSON.parse(JSON.stringify(jsonData));
  }

  /** 处理hpa的相关信息 */
  private _reduceHpaData() {
    const { subRoot, clusterVersion } = this.props,
      { minReplicas, maxReplicas, workloadName, namespace, metrics, workloadType } = subRoot.workloadEdit;

    // hpa的相关配置信息

    const resourceInfo = resourceConfig(clusterVersion)['hpa'],
      ResourceInfo = resourceConfig(clusterVersion)[workloadType];

    const isTapp = workloadType === 'tapp' ? true : false;
    // 处理hpa的metrics
    const metricsInfo: MetricOption[] = metrics.map(item => {
      let tmp: MetricOption;
      if (
        item.type === 'cpuUtilization' ||
        item.type === 'cpuAverage' ||
        item.type === 'memoryUtilization' ||
        item.type === 'memoryAverage'
      ) {
        tmp = {
          type: 'Resource',
          resource: {
            name: item.type === 'cpuUtilization' || item.type === 'cpuAverage' ? 'cpu' : 'memory',
            targetAverageUtilization:
              item.type === 'cpuUtilization' || item.type === 'memoryUtilization' ? +item.value : undefined,
            targetAverageValue:
              item.type === 'cpuAverage' || item.type === 'memoryAverage'
                ? item.type === 'cpuAverage'
                  ? +item.value * 1000 + 'm'
                  : item.value + 'Mi'
                : undefined
          }
        };
      } else if (item.type === 'inBandwidth' || item.type === 'outBandwidth') {
        tmp = {
          type: 'Pods',
          pods: {
            metricName: item.type === 'inBandwidth' ? 'pod_in_bandwidth' : 'pod_out_bandwidth',
            targetAverageValue: item.value
          }
        };
      }
      return tmp;
    });

    const jsonData: HpaEditJSONYaml = {
      kind: resourceInfo.headTitle,
      apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
      metadata: {
        name: workloadName,
        namespace: reduceNs(namespace),
        labels: {
          'qcloud-app': workloadName
        }
      },
      spec: {
        minReplicas: +minReplicas,
        maxReplicas: +maxReplicas,
        metrics: metricsInfo,
        scaleTargetRef: {
          apiVersion: isTapp ? 'apps.tkestack.io/v1' : ResourceInfo.group + '/' + ResourceInfo.version,
          kind: ResourceInfo.headTitle,
          name: workloadName
        }
      }
    };

    return JSON.parse(JSON.stringify(jsonData));
  }

  /**
   * 处理Service的相关信息
   * @param labelInfo: any selector当中的信息
   */
  private _reduceServiceData(labelInfo: any) {
    const { route, subRoot, clusterVersion } = this.props,
      { workloadEdit, serviceEdit } = subRoot,
      { portsMap, communicationType, isOpenHeadless } = serviceEdit;

    // svc的相关配置信息
    const resourceInfo = resourceConfig(clusterVersion)['svc'];

    // vpc内访问、购买lb带宽等，都放置在annotations里面实现
    const annotations = ReduceServiceAnnotations(serviceEdit, route.queries['clusterId']);

    // 构建端口映射
    const ports = ReduceServicePorts(portsMap, communicationType);

    const sessionConfig = {
      externalTrafficPolicy: serviceEdit.externalTrafficPolicy,
      sessionAffinity: serviceEdit.sessionAffinity,
      sessionAffinityTimeout: serviceEdit.sessionAffinityTimeout
    };

    const jsonData: ServiceEditJSONYaml = ReduceServiceJSONData({
      resourceInfo,
      ports,
      annotations,
      selectorObj: labelInfo,
      namespace: workloadEdit.namespace,
      communicationType,
      serviceName: workloadEdit.workloadName,
      isOpenHeadless,
      sessionConfig
    });

    return jsonData;
  }

  /** 处理数据卷的相关信息 */
  private _reduceVolumes(volumes: VolumeItem[]) {
    let volumesInfo = [];
    volumesInfo = volumes.map(volume => {
      const volumeItem = {
        name: volume.name
      };

      if (volume.volumeType === 'emptyDir') {
        // 如果使用临时路径，目前当中的配置不需要提供给用户
        volumeItem['emptyDir'] = {};
      } else if (volume.volumeType === 'hostPath') {
        // type 默认为 ""  不需要填入
        volumeItem['hostPath'] = {
          path: volume.hostPath,
          type: volume.hostPathType !== 'NoChecks' ? volume.hostPathType : undefined
        };
      } else if (volume.volumeType === 'nfsDisk') {
        const [server, path] = volume.nfsPath.split(':');
        volumeItem['nfs'] = {
          path,
          server
        };
      } else if (volume.volumeType === 'configMap') {
        volumeItem['configMap'] = {};
        volumeItem['configMap']['name'] = volume.configName;

        // 如果有configKey，则说明只需要特定的key
        if (volume.configKey) {
          volumeItem['configMap']['items'] = volume.configKey.map(c => ({
            key: c.configKey,
            mode: c.mode ? parseInt(c.mode, 8) : undefined,
            path: c.path
          }));
        }
      } else if (volume.volumeType === 'secret') {
        volumeItem['secret'] = {};
        volumeItem['secret']['secretName'] = volume.secretName;

        // 如果有 secretKey，则说明只需要特定的key
        if (volume.secretKey) {
          volumeItem['secret']['items'] = volume.secretKey.map(c => ({
            key: c.configKey,
            mode: c.mode ? parseInt(c.mode, 8) : undefined,
            path: c.path
          }));
        }
      } else if (volume.volumeType === 'pvc') {
        volumeItem['persistentVolumeClaim'] = {
          claimName: volume.pvcSelection
        };
      }
      return volumeItem;
    });
    return volumesInfo;
  }

  /** 处理container相关的配置项 */
  private _reduceContainers(containers: ContainerItem[], volumes: VolumeItem[], extraOption?: any) {
    let containersInfo = [];
    const { oversoldRatio, networkType } = extraOption;
    containersInfo = containers.map(c => {
      const containerItem = {
        name: c.name,
        image: c.registry + ':' + (c.tag ? c.tag : 'latest'),
        imagePullPolicy: c.imagePullPolicy
      };

      // 挂载点的相关配置
      if (c.mounts.length && volumes.length) {
        containerItem['volumeMounts'] = c.mounts.map(m => {
          return {
            mountPath: m.mountPath,
            subPath: m.mountSubPath ? m.mountSubPath : undefined,
            name: m.volume,
            readOnly: m.mode === 'rw' ? undefined : true
          };
        });
      }

      // request/limit的相关配置request
      let cpuLimit = c.cpuLimit.find(cpu => cpu.type === 'limit').value,
        cpuRequest = c.cpuLimit.find(cpu => cpu.type === 'request').value,
        memLimit = c.memLimit.find(mem => mem.type === 'limit').value,
        memRequest = c.memLimit.find(mem => mem.type === 'request').value;

      if (oversoldRatio.cpu) {
        cpuRequest = (+cpuLimit * 1.0) / +oversoldRatio.cpu + '';
      }
      if (oversoldRatio.memory) {
        memRequest = Math.ceil((+memLimit * 1.0) / +oversoldRatio.memory) + '';
      }
      containerItem['resources'] = {};
      // !!!注意：如果设置了gpu，需要在limits里面设定
      if (
        cpuLimit !== '' ||
        memLimit !== '' ||
        +c.gpu > 0 ||
        +c.gpuMem > 0 ||
        +c.gpuCore > 0 ||
        networkType === WorkloadNetworkTypeEnum.FloatingIP
      ) {
        containerItem['resources'] = {
          limits: {
            cpu: cpuLimit ? cpuLimit : undefined,
            memory: memLimit ? memLimit + 'Mi' : undefined,
            'nvidia.com/gpu': +c.gpu > 0 ? c.gpu + '' : undefined,
            'tencent.com/vcuda-core': +c.gpuCore ? +c.gpuCore * 100 : undefined,
            'tencent.com/vcuda-memory': +c.gpuMem ? +c.gpuMem : undefined,
            'tke.cloud.tencent.com/eni-ip': networkType === WorkloadNetworkTypeEnum.FloatingIP ? '1' : undefined
          }
        };
      }
      if (
        cpuRequest !== '' ||
        memRequest !== '' ||
        +c.gpuMem > 0 ||
        +c.gpuCore > 0 ||
        networkType === WorkloadNetworkTypeEnum.FloatingIP
      ) {
        containerItem['resources'] = Object.assign({}, containerItem['resources'], {
          requests: {
            cpu: cpuRequest ? cpuRequest : undefined,
            memory: memRequest ? memRequest + 'Mi' : undefined,
            'tencent.com/vcuda-core': +c.gpuCore ? +c.gpuCore * 100 : undefined,
            'tencent.com/vcuda-memory': +c.gpuMem ? +c.gpuMem : undefined,
            'tke.cloud.tencent.com/eni-ip': networkType === WorkloadNetworkTypeEnum.FloatingIP ? '1' : undefined
          }
        });
      }

      containerItem['env'] = [];
      c.envItems.forEach(env => {
        const envItem = {
          name: env.name
        };

        if (env.type === ContainerEnv.EnvTypeEnum.UserDefined) {
          envItem['value'] = env.value;
        } else if (
          env.type === ContainerEnv.EnvTypeEnum.SecretKeyRef ||
          env.type === ContainerEnv.EnvTypeEnum.ConfigMapRef
        ) {
          const isSecret = env.type === ContainerEnv.EnvTypeEnum.SecretKeyRef;
          const keyRef = {
            key: isSecret ? env.secretDataKey : env.configMapDataKey,
            name: isSecret ? env.secretName : env.configMapName,
            optional: false
          };

          envItem['valueFrom'] = {
            [isSecret ? 'secretKeyRef' : 'configMapKeyRef']: keyRef
          };
        } else if (env.type === ContainerEnv.EnvTypeEnum.FieldRef) {
          envItem['valueFrom'] = {
            fieldRef: {
              apiVersion: env.apiVersion,
              fieldPath: env.fieldName
            }
          };
        } else if (env.type === ContainerEnv.EnvTypeEnum.ResourceFieldRef) {
          envItem['valueFrom'] = {
            resourceFieldRef: {
              containerName: c.name,
              resource: env.resourceFieldName,
              divisor: env.divisor
            }
          };
        }
        containerItem['env'].push(envItem);
      });

      // 如果有工作目录
      if (c.workingDir) {
        containerItem['workingDir'] = c.workingDir;
      }

      // 如果有运行命令 command: string[]
      if (c.cmd) {
        containerItem['command'] = c.cmd.trim().split('\n');
      }

      // 如果有运行参数
      if (c.arg) {
        containerItem['args'] = c.arg.trim().split('\n');
      }

      // 特权级容器
      if (c.privileged) {
        containerItem['securityContext'] = {
          privileged: c.privileged
        };
      }

      // 增加权限集
      if (!isEmpty(c.addCapabilities)) {
        if (isEmpty(containerItem['securityContext'])) {
          containerItem['securityContext'] = {};
        }
        if (isEmpty(containerItem['securityContext']['capabilities'])) {
          containerItem['securityContext']['capabilities'] = {};
        }
        containerItem['securityContext']['capabilities']['add'] = c.addCapabilities;
      }

      // 删除权限集
      if (!isEmpty(c.dropCapabilities)) {
        if (isEmpty(containerItem['securityContext'])) {
          containerItem['securityContext'] = {};
        }
        if (isEmpty(containerItem['securityContext']['capabilities'])) {
          containerItem['securityContext']['capabilities'] = {};
        }
        containerItem['securityContext']['capabilities']['drop'] = c.dropCapabilities;
      }

      // 存活检查
      const reduceHealthCheck = (healthCheckItem: HealthCheckItem) => {
        const healthItem = {
          failureThreshold: +healthCheckItem.unhealthThreshold,
          successThreshold: +healthCheckItem.healthThreshold,
          initialDelaySeconds: healthCheckItem.delayTime ? +healthCheckItem.delayTime : undefined,
          timeoutSeconds: healthCheckItem.timeOut ? +healthCheckItem.timeOut : undefined,
          periodSeconds: healthCheckItem.intervalTime ? +healthCheckItem.intervalTime : undefined
        };

        if (healthCheckItem.checkMethod === 'methodTcp') {
          healthItem['tcpSocket'] = {
            port: +healthCheckItem.port
          };
        } else if (healthCheckItem.checkMethod === 'methodHttp') {
          healthItem['httpGet'] = {
            path: healthCheckItem.path,
            port: +healthCheckItem.port,
            scheme: healthCheckItem.protocol
          };
        } else if (healthCheckItem.checkMethod === 'methodCmd') {
          healthItem['exec'] = {
            command: healthCheckItem.cmd.split('\n').map(item => item.trim())
          };
        }

        return healthItem;
      };

      if (c.healthCheck.isOpenLiveCheck) {
        const healthCheckItem = c.healthCheck.liveCheck;
        containerItem['livenessProbe'] = reduceHealthCheck(healthCheckItem);
      }

      if (c.healthCheck.isOpenReadyCheck) {
        const healthCheckItem = c.healthCheck.readyCheck;
        containerItem['readinessProbe'] = reduceHealthCheck(healthCheckItem);
      }

      return JSON.parse(JSON.stringify(containerItem));
    });
    return containersInfo;
  }
  /** 处理亲和性调度的相关信息 */
  private _reduceNodeAffinityInfo(nodeAffinityType: string, nodeAffinityRule: AffinityRule, nodeSelection: Computer[]) {
    const affinityInfo = {};
    if (nodeAffinityType === affinityType.node) {
      const nodeSelector = nodeSelection.map(node => {
        return node.metadata.labels['kubernetes.io/hostname'];
      });
      affinityInfo['requiredDuringSchedulingIgnoredDuringExecution'] = {
        nodeSelectorTerms: [
          {
            matchExpressions: [
              {
                key: 'kubernetes.io/hostname',
                operator: 'In',
                values: nodeSelector
              }
            ]
          }
        ]
      };
    } else if (nodeAffinityType === affinityType.rule) {
      affinityInfo['preferredDuringSchedulingIgnoredDuringExecution'] = nodeAffinityRule.preferredExecution[0]
        .preference.matchExpressions.length
        ? [
            {
              preference: {
                matchExpressions: nodeAffinityRule.preferredExecution[0].preference.matchExpressions.map(rule => {
                  return {
                    key: rule.key,
                    operator: rule.operator,
                    values: rule.values ? rule.values.split(';') : undefined
                  };
                })
              },
              weight: 1
            }
          ]
        : undefined;
      affinityInfo['requiredDuringSchedulingIgnoredDuringExecution'] = nodeAffinityRule.requiredExecution[0]
        .matchExpressions.length
        ? {
            nodeSelectorTerms: [
              {
                matchExpressions: nodeAffinityRule.requiredExecution[0].matchExpressions.map(rule => {
                  return {
                    key: rule.key,
                    operator: rule.operator,
                    values: rule.values ? rule.values.split(';') : undefined
                  };
                })
              }
            ]
          }
        : undefined;
    }
    return { nodeAffinity: affinityInfo };
  }
}
