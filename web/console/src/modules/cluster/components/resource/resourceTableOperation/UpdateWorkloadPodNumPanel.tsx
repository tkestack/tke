import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { OperationState, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { Button } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';
import { getWorkflowError } from '../../../../common/utils';
import { MainBodyLayout, FormLayout } from '../../../../common/layouts';
import { FormItem, TipInfo } from '../../../../common/components';
import { EditResourceContainerNumPanel } from '../resourceEdition/EditResourceContainerNumPanel';
import { router } from '../../../router';
import { validateWorkloadActions } from '../../../actions/validateWorkloadActions';
import { HpaEditJSONYaml, WorkloadEditJSONYaml, MetricOption, CreateResource, HpaMetrics } from '../../../models';
import { resourceConfig } from '../../../../../../config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** 加载中的样式 */
const loadingElement = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

let mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class UpdateWorkloadPodNumPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    // 清空workloadEdit的编辑项
    actions.editWorkload.clearWorkloadEdit();
  }

  componentDidMount() {
    let { actions, route, subRoot } = this.props,
      { resourceList } = subRoot.resourceOption;

    let { np: namespace, rid, clusterId, resourceIns } = route.queries;

    // 这里去拉取hpa的信息
    actions.editWorkload.hpa.applyFilter({
      namespace,
      regionId: +rid,
      clusterId,
      specificName: resourceIns
    });

    // 这里是从列表页进入的时候，需要去初始化 workloadEdit当中的内容，如果是直接在当前页面刷新的话，会去拉取列表，在fetchResource之后，会初始化
    if (resourceList.data.recordCount) {
      let finder = resourceList.data.records.find(item => item.metadata.name === resourceIns);
      finder && actions.editWorkload.updateContainerNum(finder.spec.replicas);
    }
  }

  render() {
    let { subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { resourceOption, updateResourcePart } = subRoot,
      { resourceSelection, resourceList } = resourceOption;

    let failed = updateResourcePart.operationState === OperationState.Done && !isSuccessWorkflow(updateResourcePart);

    return (
      <MainBodyLayout>
        <FormLayout>
          {resourceList.fetched !== true || resourceList.fetchState === FetchState.Fetching ? (
            loadingElement
          ) : (
            <div className="param-box server-update add">
              <ul className="form-list jiqun fixed-layout">
                <FormItem label={t('当前实例数量')}>
                  {(resourceSelection[0] && resourceSelection[0].spec && resourceSelection[0].spec.replicas) || 0}
                </FormItem>

                <EditResourceContainerNumPanel />

                <li className="pure-text-row fixed">
                  <Button
                    className="mr10"
                    type="primary"
                    disabled={updateResourcePart.operationState === OperationState.Performing}
                    onClick={this._handleSubmit.bind(this)}
                  >
                    {failed ? t('重试') : t('更新实例数量')}
                  </Button>
                  <Button onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}>
                    {t('取消')}
                  </Button>
                  <TipInfo
                    isShow={failed}
                    className="error"
                    style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px' }}
                  >
                    {getWorkflowError(updateResourcePart)}
                  </TipInfo>
                </li>
              </ul>
            </div>
          )}
        </FormLayout>
      </MainBodyLayout>
    );
  }

  /** 处理提交请求
   * 1. 如果是更新实例数量，只更新deployment，然后删除对应的hpa
   * 2. 更新hpa，只需要修改hpa的信息，deployment不需要做更改，如果原来没有创建hpa，则需要创建hpa
   */
  private _handleSubmit() {
    let { actions, subRoot, route, clusterVersion } = this.props,
      { workloadEdit, resourceName } = subRoot;

    let isManual = workloadEdit.scaleType === 'manualScale',
      isHasHpa = workloadEdit.hpaList.data.recordCount > 0 ? true : false,
      resourceInfo = resourceConfig(clusterVersion)[resourceName],
      hpaResourceInfo = resourceConfig(clusterVersion)['hpa'];

    actions.validate.workload.validatePodNumEdit();

    if (validateWorkloadActions._validatePodNumEdit(workloadEdit)) {
      let jsonData: WorkloadEditJSONYaml;
      let hpaJsonData: HpaEditJSONYaml;

      // 获取当前的ns 和 资源名称
      let { np: namespace, resourceIns: workloadName, clusterId, rid } = route.queries;

      if (isManual) {
        jsonData = {
          spec: {
            replicas: +workloadEdit.containerNum
          }
        };
      } else {
        let { minReplicas, maxReplicas, metrics } = workloadEdit;

        // 处理hpa的metrics
        let metricsInfo = this._reduceHpaMetrics(metrics);

        // 这里需要判断当前的hpa是否是已经创建好的，还是新的
        if (isHasHpa) {
          hpaJsonData = {
            spec: {
              minReplicas: +minReplicas,
              maxReplicas: +maxReplicas,
              metrics: metricsInfo
            }
          };
        } else {
          hpaJsonData = {
            kind: hpaResourceInfo.headTitle,
            apiVersion: (hpaResourceInfo.group ? hpaResourceInfo.group + '/' : '') + hpaResourceInfo.version,
            metadata: {
              name: workloadName,
              namespace,
              labels: {
                'qcloud-app': workloadName
              }
            },
            spec: {
              minReplicas: +minReplicas,
              maxReplicas: +maxReplicas,
              metrics: metricsInfo,
              scaleTargetRef: {
                apiVersion: resourceInfo.group + '/' + resourceInfo.version,
                kind: resourceInfo.headTitle,
                name: workloadName
              }
            }
          };
        }
      }

      // 当前的编辑模式
      let mode = isHasHpa || isManual ? 'update' : 'create';

      let resource: CreateResource = {
        id: uuid(),
        resourceInfo: isManual ? resourceInfo : hpaResourceInfo,
        mode,
        namespace,
        clusterId,
        resourceIns: mode === 'update' ? workloadName : '',
        isStrategic: resourceName === 'tapp' ? false : true,
        jsonData: workloadEdit.scaleType === 'autoScale' ? JSON.stringify(hpaJsonData) : JSON.stringify(jsonData)
      };

      // 如果是手动更新,，并且原本有hpa的信息，还需要删除相对应的hpa
      if (isManual && isHasHpa) {
        let hpaResource: CreateResource = {
          id: uuid(),
          resourceInfo: hpaResourceInfo,
          namespace,
          clusterId,
          resourceIns: workloadName
        };
        actions.workflow.deleteResource.start([hpaResource], +rid);
        actions.workflow.deleteResource.perform();
      }

      // 如果是手动 或者 hpa（原来已经存在，则为更新）
      if (isManual || (!isManual && isHasHpa)) {
        actions.workflow.updateResourcePart.start([resource], +rid);
        actions.workflow.updateResourcePart.perform();
      } else {
        actions.workflow.modifyResource.start([resource], +rid);
        actions.workflow.modifyResource.perform();
      }
    }
  }

  /** 处理hpa的metrics的信息 */
  private _reduceHpaMetrics(metrics: HpaMetrics[]) {
    return metrics.map(item => {
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
  }
}
