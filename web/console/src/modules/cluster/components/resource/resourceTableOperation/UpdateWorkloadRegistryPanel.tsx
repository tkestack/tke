import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Radio, Select } from '@tea/component';
import { bindActionCreators, FetchState, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, TipInfo } from '../../../../common/components';
import { FixedFormLayout, FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { getWorkflowError } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validateWorkloadActions } from '../../../actions/validateWorkloadActions';
import { CreateResource, WorkloadEditJSONYaml } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

/** 加载中的样式 */
const loadingElement = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

/** 资源的更新方式列表 */
const updateTypeList = [
  {
    label: t('滚动更新（推荐）'),
    value: 'RollingUpdate',
    tip: t('对实例进行逐个更新，这种方式可以让您不中断业务实现对服务的更新')
  },
  {
    label: 'OnDelete',
    value: 'OnDelete',
    tip: t('手动删除实例时触发更新')
  }
];

const deploymentUpdateTypeList = [
  {
    label: t('滚动更新（推荐）'),
    value: 'RollingUpdate',
    tip: t('对实例进行逐个更新，这种方式可以让您不中断业务实现对服务的更新')
  },
  {
    label: t('快速更新'),
    value: 'Recreate',
    tip: t('直接关闭所有实例，启动相同数量的新实例')
  }
];

/** 滚动更新的策略，先删除pod 还是先创建pod */
const rollingUpdateTypeList = [
  {
    value: 'createPod',
    label: t('启动新的Pod,停止旧的Pod'),
    tip: t('请确认集群有足够的CPU和内存用于启动新的Pod, 否则可能导致集群崩溃')
  },
  {
    value: 'destroyPod',
    label: t('停止旧的Pod，启动新的Pod')
  },
  {
    value: 'userDefined',
    label: t('自定义')
  }
];

interface UpdateWorkloadRegistryPanelState {
  /** 当前编辑的container的Id */
  currentEditingContainerId?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UpdateWorkloadRegistryPanel extends React.Component<RootProps, UpdateWorkloadRegistryPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      currentEditingContainerId: ''
    };
  }

  componentDidMount() {
    let { actions, subRoot, route } = this.props,
      { ffResourceList } = subRoot.resourceOption;

    // 这里是从列表页进入的时候，需要去初始化 workloadEdit当中的内容，如果是直接在当前页面刷新的话，会去拉取列表，在fetchResource之后，会初始化
    if (ffResourceList.list.data.recordCount) {
      const finder = ffResourceList.list.data.records.find(item => item.metadata.name === route.queries['resourceIns']);
      finder && actions.editWorkload.initWorkloadEditForUpdateRegistry(finder);
    }
  }

  componentWillUnmount() {
    const { actions } = this.props;
    actions.editWorkload.clearWorkloadEdit();
    actions.workflow.updateResourcePart.reset();
  }

  render() {
    let { actions, subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { workloadEdit, updateResourcePart, resourceOption } = subRoot,
      { ffResourceList } = resourceOption,
      {
        partition,
        v_partition,
        rollingUpdateStrategy,
        resourceUpdateType,
        v_minReadySeconds,
        minReadySeconds,
        batchSize,
        v_batchSize,
        maxSurge,
        v_maxSurge,
        maxUnavailable,
        v_maxUnavailable
      } = workloadEdit;

    /** 当前滚动更新镜像的资源类型 */
    const workloadType = urlParams['resourceName'],
      isDeployment = workloadType === 'deployment',
      isStatefulset = workloadType === 'statefulset',
      isDaemonset = workloadType === 'daemonset',
      isTapp = workloadType === 'tapp';

    /** 渲染更新方式 的列表 */
    const finalUpdateTypeList = isDeployment ? deploymentUpdateTypeList : updateTypeList;
    const updateTypeOptions = finalUpdateTypeList.map(item => ({
      value: item.value,
      text: item.label
    }));
    const updateTypeTips = finalUpdateTypeList.find(item => item.value === resourceUpdateType).tip;

    /** 滚动更新的策略 选择项 */
    const finder = rollingUpdateTypeList.find(c => c.value === rollingUpdateStrategy),
      tip = finder ? finder.tip : '';

    const failed = updateResourcePart.operationState === OperationState.Done && !isSuccessWorkflow(updateResourcePart);

    return (
      <MainBodyLayout>
        <FormLayout>
          {ffResourceList.list.fetched !== true || ffResourceList.list.fetchState === FetchState.Fetching ? (
            loadingElement
          ) : (
            <div className="param-box server-update add">
              <ul className="form-list jiqun fixed-layout">
                <FormItem label={t('更新方式')} isShow={!isTapp}>
                  <Select
                    style={{ width: '150px' }}
                    options={updateTypeOptions}
                    value={resourceUpdateType}
                    onChange={value => {
                      actions.editWorkload.changeResourceUpdateType(value);
                    }}
                  />
                  <p className="text-label">{updateTypeTips}</p>
                </FormItem>
                <FormItem
                  label={t('更新间隔')}
                  isShow={resourceUpdateType === 'RollingUpdate' && !isStatefulset && !isTapp}
                >
                  <InputField
                    type="text"
                    tipMode="popup"
                    style={{ maxWidth: '80px' }}
                    ops={t('秒')}
                    validator={v_minReadySeconds}
                    value={minReadySeconds}
                    onChange={actions.editWorkload.inputMinReadySeconds}
                    onBlur={actions.validate.workload.validateMinReadySeconds}
                  />
                </FormItem>
                <FormItem
                  label={t('更新策略')}
                  isShow={resourceUpdateType === 'RollingUpdate' && isDeployment && !isTapp}
                >
                  <div className="form-unit">
                    <Radio.Group
                      value={rollingUpdateStrategy}
                      onChange={value => actions.editWorkload.changeRollingUpdateStrategy(value)}
                      // style={{ fontSize: '12px', display: 'inline-block' }}
                    >
                      {rollingUpdateTypeList.map((item, rIndex) => {
                        return (
                          <Radio key={rIndex} name={item.value}>
                            {item.label}
                          </Radio>
                        );
                      })}
                    </Radio.Group>
                    {tip && <p className="text-label">{tip}</p>}
                  </div>
                </FormItem>

                <FormItem label={t('策略配置')} isShow={resourceUpdateType === 'RollingUpdate' && !isTapp}>
                  <div className="run-docker-box" style={isStatefulset ? {} : { paddingBottom: '1px' }}>
                    <div className="edit-param-list">
                      <div className="param-box">
                        <div className="param-bd" style={{ marginBottom: '0' }}>
                          <ul
                            className="form-list fixed-layout"
                            style={isStatefulset ? { marginTop: '0' } : { marginTop: '0', paddingBottom: '0' }}
                          >
                            <FormItem label="Pods" isShow={isDeployment && rollingUpdateStrategy !== 'userDefined'}>
                              <InputField
                                type="text"
                                tipMode="popup"
                                placeholder={t('0、正整数或者正百分数（default: 25%）')}
                                validator={v_batchSize}
                                value={batchSize}
                                onChange={actions.editWorkload.inputBatchSize}
                                onBlur={actions.validate.workload.validateBatchSize}
                              />
                              <p className="text-label">{t('Pod将批量启动或停止')}</p>
                            </FormItem>
                            <FormItem label="MaxSurge" isShow={isDeployment && rollingUpdateStrategy === 'userDefined'}>
                              <InputField
                                type="text"
                                tipMode="popup"
                                style={{ minWidth: '300px' }}
                                placeholder={t('0、正整数或者正百分数（default: 25%）')}
                                validator={v_maxSurge}
                                value={maxSurge}
                                onChange={actions.editWorkload.inputMaxSurge}
                                onBlur={actions.validate.workload.validateMaxSurge}
                              />
                              <p className="text-label">{t('允许超出所需规模的最大Pod数量')}</p>
                            </FormItem>
                            <FormItem
                              label="MaxUnavailable"
                              isShow={(isDeployment && rollingUpdateStrategy === 'userDefined') || isDaemonset}
                            >
                              <InputField
                                type="text"
                                tipMode="popup"
                                style={{ minWidth: '300px' }}
                                placeholder={t('0、正整数或者正百分数（default: 25%）')}
                                validator={v_maxUnavailable}
                                value={maxUnavailable}
                                onChange={actions.editWorkload.inputMaxUnavaiable}
                                onBlur={actions.validate.workload.validateMaxUnavaiable}
                              />
                              <p className="text-label">{t('允许最大不可用的Pod数量')}</p>
                            </FormItem>
                            <FormItem label="Partition" isShow={isStatefulset}>
                              <InputField
                                type="text"
                                tipMode="popup"
                                style={{ minWidth: '300px' }}
                                placeholder={t('0或者正整数（default: 0）')}
                                validator={v_partition}
                                value={partition}
                                onChange={actions.editWorkload.inputPartition}
                                onBlur={actions.validate.workload.validatePartition}
                              />
                            </FormItem>
                          </ul>
                        </div>
                      </div>
                    </div>
                  </div>
                </FormItem>

                <FormItem label={t('MaxUnavailabl（个）')} isShow={isTapp}>
                  <InputField
                    type="text"
                    style={{ minWidth: '300px' }}
                    placeholder={t('正整数')}
                    tipMode="popup"
                    validator={v_maxUnavailable}
                    value={maxUnavailable}
                    onChange={actions.editWorkload.inputMaxUnavaiable}
                    onBlur={() => actions.validate.workload.validateMaxUnavaiable(true)}
                  />
                  <p className="text-label">{t('允许最大不可用数量，可限制更新并发数')}</p>
                </FormItem>
                {this._renderContainerInfo()}

                <li className="pure-text-row fixed">
                  <div className="form-input">
                    <Button
                      className="mr10"
                      type="primary"
                      disabled={updateResourcePart.operationState === OperationState.Performing}
                      onClick={this._handleSubmit.bind(this)}
                    >
                      {failed ? t('重试') : t('完成')}
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
            </div>
          )}
        </FormLayout>
      </MainBodyLayout>
    );
  }

  /** 处理提交请求 */
  private _handleSubmit() {
    let { actions, subRoot, route } = this.props,
      { resourceInfo, mode, workloadEdit, resourceOption } = subRoot,
      { ffResourceList } = resourceOption,
      targetResource = ffResourceList.selection;

    actions.validate.workload.validateUpdateRegistryEdit();

    if (validateWorkloadActions._validateUpdateRegistryEdit(workloadEdit)) {
      const {
        minReadySeconds,
        containers,
        workloadType,
        partition,
        resourceUpdateType,
        rollingUpdateStrategy,
        maxSurge,
        maxUnavailable,
        batchSize
      } = workloadEdit;

      // 当前的资源的类型
      const isStatefulset = workloadType === 'statefulset',
        isDeployment = workloadType === 'deployment',
        isDaemonset = workloadType === 'daemonset',
        isTapp = workloadType === 'tapp';
      // 当前镜像的更新方式
      const isRollingUpdate = resourceUpdateType === 'RollingUpdate';

      // 获取deployment滚动更新的内容
      let deploymentRollingUpdateContent = {};
      if (isDeployment && isRollingUpdate) {
        deploymentRollingUpdateContent = {
          maxSurge:
            rollingUpdateStrategy === 'userDefined'
              ? +maxSurge
              : rollingUpdateStrategy === 'createPod'
              ? +batchSize
              : 0,
          maxUnavailable:
            rollingUpdateStrategy === 'userDefined'
              ? +maxUnavailable
              : rollingUpdateStrategy === 'createPod'
              ? 0
              : +batchSize
        };
      }

      // statefulset滚动更新的内容
      let statefulsetRollingUpdateContent = {};
      if (isStatefulset && isRollingUpdate) {
        statefulsetRollingUpdateContent = {
          partition: partition ? +partition : 0
        };
      }

      // daemonset滚动更新的内容
      let daemonsetRollingUpdateContent = {};
      if (isDaemonset && isRollingUpdate) {
        daemonsetRollingUpdateContent = maxUnavailable
          ? {
              maxUnavailable: +maxUnavailable
            }
          : undefined;
      }

      //tapp滚动更新的内容
      let tappRollingUpdateContent = {};
      if (isTapp) {
        tappRollingUpdateContent = {
          template: 'default',
          maxUnavailable: +maxUnavailable
        };
      }
      // 获取容器的内容
      const containersInfo = containers.map(c => {
        const targetContainer = targetResource.spec.template.spec.containers.find(e => e.name === c.name);
        return !isTapp //对于tapp来说，由于更新方式是merge（非strategy-merge），所以需要带上原本的container之前的内容 不然之前的内容会被清空掉
          ? {
              name: c.name,
              image: c.registry + ':' + (c.tag ? c.tag : 'latest')
            }
          : {
              ...targetContainer,
              name: c.name,
              image: c.registry + ':' + (c.tag ? c.tag : 'latest')
            };
      });

      /**
       * https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/#rollingupdatedeployment-v1-apps
       * 构建创建 workload的json的格式
       */
      let jsonData: WorkloadEditJSONYaml = {
        spec: {
          minReadySeconds: isStatefulset || isTapp ? undefined : isRollingUpdate ? +minReadySeconds : 0,
          strategy: isDeployment
            ? {
                type: resourceUpdateType,
                rollingUpdate: isRollingUpdate ? deploymentRollingUpdateContent : null
              }
            : undefined,
          updateStrategy: !isDeployment
            ? isTapp
              ? tappRollingUpdateContent
              : {
                  type: resourceUpdateType,
                  rollingUpdate: isRollingUpdate
                    ? isStatefulset
                      ? statefulsetRollingUpdateContent
                      : daemonsetRollingUpdateContent
                    : null
                }
            : undefined,
          templates: isTapp ? null : undefined,
          template: {
            spec: {
              containers: containersInfo
            }
          }
        }
      };

      // 去除当中不需要的数据
      jsonData = JSON.parse(JSON.stringify(jsonData));

      const resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        mode,
        namespace: route.queries['np'],
        clusterId: route.queries['clusterId'],
        resourceIns: route.queries['resourceIns'],
        jsonData: JSON.stringify(jsonData),
        isStrategic: !isTapp //Tapp并不支持Strategic-merge-patch的方式
      };

      actions.workflow.updateResourcePart.start([resource], +route.queries['rid']);
      actions.workflow.updateResourcePart.perform();
    }
  }

  /** 展示容器的相关信息 */
  private _renderContainerInfo() {
    const { actions, subRoot } = this.props,
      { workloadEdit, resourceOption } = subRoot,
      { ffResourceList } = resourceOption,
      { containers } = workloadEdit;

    const loadingElement: JSX.Element = (
      <div>
        <i className="n-loading-icon" />
        &nbsp; <span className="text">{t('加载中...')}</span>
      </div>
    );

    return (
      <FormItem label={t('容器')}>
        {ffResourceList.list.fetched === false || ffResourceList.list.fetchState === FetchState.Fetching
          ? loadingElement
          : containers.map((container, index) => {
              const cKey = container.id + '';

              return (
                <FixedFormLayout key={index} isRemoveUlMarginTop={true}>
                  <FormItem label={t('名称')}>
                    <p className="text">{container.name}</p>
                  </FormItem>
                  <FormItem label={t('镜像')}>
                    <div className={classnames('form-unit', { 'is-error': container.v_registry.status === 2 })}>
                      <input
                        type="text"
                        className="tc-15-input-text m mr10"
                        style={{ minWidth: 300 }}
                        value={container.registry}
                        onChange={e => actions.editWorkload.updateContainer({ registry: e.target.value }, cKey)}
                        onBlur={e => actions.validate.workload.validateRegistrySelection(e.target.value, cKey)}
                      />
                      {/* <Button type="link">{t('选择镜像')}</Button> */}
                    </div>
                  </FormItem>
                  <FormItem label={t('镜像版本（Tag）')} className="tag-mod">
                    <div className="tc-15-autocomplete xl">
                      <input
                        style={{ minWidth: 300 }}
                        placeholder={t('不填默认为latest')}
                        type="text"
                        className="tc-15-input-text m"
                        value={container.tag}
                        onChange={e => {
                          actions.editWorkload.updateContainer({ tag: e.target.value }, cKey);
                        }}
                      />
                    </div>
                  </FormItem>
                </FixedFormLayout>
              );
            })}
      </FormItem>
    );
  }
}
