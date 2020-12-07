import * as React from 'react';
import { connect } from 'react-redux';
import { IsShowLoginDialog } from 'src/modules/cluster/constants/ActionType';
import { TappGrayUpdateEditItem } from 'src/modules/cluster/models/ResourceDetailState';

import { Bubble, Button, Justify, Modal, Table, TagSearchBox } from '@tea/component';
import { bindActionCreators, insertCSS, uuid } from '@tencent/ff-redux';
import { ChartInstancesPanel } from '@tencent/tchart';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../../config';
import { initValidator } from '../../../../../../src/modules/common';
import { cloneDeep, isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { CreateResource, PodFilterInNode } from '../../../models';
import { containerMonitorFields, MonitorPanelProps, podMonitorFields } from '../../../models/MonitorPanel';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { IsInNodeManageDetail } from './ResourceDetail';
import { ResourceGrayUpgradeDialog } from './ResourceGrayUpgradeDialog';
import { ResourceTappPodDeleteDialog } from './ResourceTappPodDeleteDialog';
import { reduceNs } from '@helper';

/** k8s pod的状态值 */
const PodPhase = ['Pending', 'Running', 'Succeeded', 'Failed', 'Unknown'];

/** tagSearch的key的映射名称 */
const TagSearchKeyMap = {
  podName: t('Pod名称'),
  phase: t('状态'),
  namespace: t('命名空间')
};

interface ResourcePodActionPanelState {
  /** searchbox的 */
  searchBoxValues?: any[];

  monitorPanelProps?: MonitorPanelProps;
}

/**渲染Tapp的操作方式 灰度升级按钮 */

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourcePodActionPanel extends React.Component<RootProps, ResourcePodActionPanelState> {
  constructor(props: RootProps) {
    super(props);
    let { route, subRoot } = props,
      urlParams = router.resolve(route),
      { podFilterInNode } = subRoot.resourceDetailState;

    // 判断是否在nodeManage的pod列表里面
    const isInNodeManage = IsInNodeManageDetail(urlParams['type']);
    let initSearchBoxValues = [];
    // 需要初始化展示的数据，因为有可能几个tab之间切来切去
    if (isInNodeManage && !isEmpty(podFilterInNode)) {
      const keys = Object.keys(podFilterInNode);
      initSearchBoxValues = keys.map(item => {
        return {
          attr: {
            // type: 'single',
            key: item,
            name: TagSearchKeyMap[item]
            // values: []
          },
          values: [{ name: podFilterInNode[item] }]
        };
      });
    }
    this.state = {
      searchBoxValues: initSearchBoxValues
    };
  }

  /**渲染pod管理页面的详情操作按钮 */

  render() {
    const { route } = this.props,
      urlParams = router.resolve(route);

    const type = urlParams['type'];
    const isInNodeManage = IsInNodeManageDetail(type);
    const { resourceName } = this.props.subRoot;

    let monitorButton = null;
    monitorButton = type === 'resource' || isInNodeManage ? this._renderMonitorButton() : null;
    return (
      <Table.ActionPanel>
        <Justify
          left={
            <React.Fragment>
              {monitorButton}
              {resourceName === 'tapp' ? this._renderTappOperationBar() : null}
            </React.Fragment>
          }
          right={this._renderTagSearchBox()}
        />
        {this.state && this.state.monitorPanelProps && (
          <Modal
            visible={true}
            caption={this.state.monitorPanelProps.title}
            onClose={() => this.setState({ monitorPanelProps: undefined })}
            size={1050}
          >
            <Modal.Body>
              <ChartInstancesPanel
                tables={this.state.monitorPanelProps.tables}
                groupBy={this.state.monitorPanelProps.groupBy}
                instance={this.state.monitorPanelProps.instance}
              >
                {this.state.monitorPanelProps.headerExtraDOM}
              </ChartInstancesPanel>
            </Modal.Body>
          </Modal>
        )}
      </Table.ActionPanel>
    );
  }

  private _renderTappOperationBar() {
    const { podSelection } = this.props.subRoot.resourceDetailState;
    return (
      <React.Fragment>
        <Bubble content={podSelection.length <= 0 ? t('请先选择容器') : null} />
        <Button
          type="primary"
          disabled={podSelection.length <= 0}
          onClick={this._handleClickForTappPodRemove.bind(this)}
        >
          {t('删除')}
        </Button>
        <Bubble content={podSelection.length <= 0 ? t('请先选择容器') : null} placement="right">
          <Button
            type="primary"
            disabled={podSelection.length <= 0}
            onClick={this._handleClickForUpdateGrayTapp.bind(this)}
          >
            {t('灰度升级')}
          </Button>
        </Bubble>
      </React.Fragment>
    );
  }
  /**
   * pre: 在Node详情页内
   * 生成搜索框
   */
  private _renderTagSearchBox() {
    const { route } = this.props,
      urlParams = router.resolve(route);

    const type = urlParams['type'];
    const isInNodeManage = IsInNodeManageDetail(type);
    const { subRoot, namespaceList } = this.props,
      { podList } = subRoot.resourceDetailState;

    /** podFilter的选择项 */
    const podNameValues: any[] = [];
    const podPhaseValues: any[] = PodPhase.map(item => ({
      key: 'phase',
      name: item
    }));
    const namespaceValues: any[] = [];

    podList.data.records.forEach(item => {
      // 获取podName的集合
      podNameValues.push({
        key: 'podName',
        name: item.metadata.name
      });
    });

    namespaceList.data.records.forEach(item => {
      namespaceValues.push({
        key: 'namespace',
        name: item.name
      });
    });

    // tagSearch的过滤选项
    const attributes = [
      {
        type: 'input',
        key: 'podName',
        name: t('Pod名称')
      },
      {
        type: 'single',
        key: 'phase',
        name: t('状态'),
        values: podPhaseValues
      },

      {
        type: 'input',
        key: 'ip',
        name: t('实例IP')
      }
    ].concat(
      isInNodeManage
        ? {
            type: 'single',
            key: 'namespace',
            name: t('命名空间'),
            values: namespaceValues
          }
        : []
    );

    return (
      <div style={{ width: 600, float: 'right' }}>
        <TagSearchBox
          className="myTagSearchBox"
          attributes={attributes}
          // minWidth={420}
          value={this.state.searchBoxValues}
          onChange={tags => {
            this._handleClickForTagSearch(tags);
          }}
        />
      </div>
    );
  }

  /** 搜索框的操作，不同的搜索进行相对应的操作 */
  private _handleClickForTagSearch(tags: any[]) {
    const { actions } = this.props;

    // 这里是控制tagSearch的展示
    this.setState({
      searchBoxValues: tags
    });

    const podFilter: PodFilterInNode = {};
    tags.forEach(item => {
      const attrKey = item.attr ? item.attr.key : 'podName';
      podFilter[attrKey] = item.values[0]['name'];
    });
    actions.resourceDetail.pod.updatePodFilterInNode(podFilter);
  }

  private _handleClickForTappPodRemove() {
    const { podSelection } = this.props.subRoot.resourceDetailState;
    const { clusterVersion, route } = this.props;
    const podResourceInfo = resourceConfig(clusterVersion)['tapp'];

    //删除tapp pod 需要修改tapp的yaml上的spec.statuses字段，具体为{'pod名字':killed}
    const statuses = {};
    podSelection.forEach(pod => {
      const indexName = pod.metadata.name.split(pod.metadata.generateName)[1];
      if (indexName) {
        statuses[indexName] = 'Killed';
      }
    });
    const jsonYaml = {
      spec: {
        statuses
      }
    };
    // 需要提交的数据
    const resource: CreateResource = {
      id: uuid(),
      resourceInfo: podResourceInfo,
      namespace: route.queries['np'],
      clusterId: route.queries['clusterId'],
      resourceIns: route.queries['resourceIns'],
      jsonData: JSON.stringify(jsonYaml),
      isStrategic: false
    };
    this.props.actions.workflow.removeTappPod.start([resource]);
  }

  private _handleClickForUpdateGrayTapp() {
    const { podSelection } = this.props.subRoot.resourceDetailState;
    const items: TappGrayUpdateEditItem[] = podSelection.map(pod => {
      return {
        name: pod.metadata.name,
        generateName: pod.metadata.generateName,
        containers: pod.spec.containers.map(container => {
          let index = container.image.lastIndexOf(':'),
            imageName,
            imageTag;
          if (index !== -1) {
            imageName = container.image.slice(0, index);
            imageTag = container.image.slice(index + 1);
          } else {
            imageName = container.image;
            imageTag = '';
          }
          return {
            name: container.name,
            imageName,
            imageTag,
            v_imageName: initValidator
          };
        })
      };
    });
    this.props.actions.resourceDetail.pod.initTappGrayUpdate(items);
    this.props.actions.workflow.updateGrayTapp.start();
  }
  /** render监控按钮 */
  private _renderMonitorButton() {
    return (
      <Button
        onClick={() => {
          this._handleMonitor('podMonitor');
        }}
        type="primary"
      >
        {t('监控')}
      </Button>
    );
  }

  /** 处理监控的相关操作 */
  private _handleMonitor(monitorType?: string, pod_name?: string) {
    const { route, subRoot } = this.props,
      urlParams = router.resolve(route),
      { resourceDetailState } = subRoot;

    // 判断是否在node详情页当中
    const isInNodeManage = IsInNodeManageDetail(urlParams['type']);

    const containerById = {};
    for (const container of resourceDetailState.containerList) {
      containerById[container.id] = container;
    }

    const monitorPanelProps =
      monitorType === 'podMonitor'
        ? {
            tables: [
              {
                table: 'k8s_pod',
                conditions: [
                  ['tke_cluster_instance_id', '=', route.queries['clusterId']],
                  ...(isInNodeManage
                    ? [['node', '=', route.queries['resourceIns'] || '']]
                    : [
                        ['workload_name', '=', route.queries['resourceIns'] || ''],
                        ['namespace', '=', reduceNs(route.queries['np'] || 'default')]
                      ])
                ],
                fields: podMonitorFields
              }
            ],
            groupBy: [{ value: 'pod_name' }],
            instance: {
              columns: [
                {
                  key: 'pod_name',
                  name: t('Pod名称')
                }
              ],
              list: resourceDetailState.podList.data.records.map(item => ({
                pod_name: item.metadata.name,
                isChecked:
                  !resourceDetailState.podSelection.length ||
                  resourceDetailState.podSelection.find(selected => selected.metadata.name === item.metadata.name)
              }))
            }
          }
        : {
            tables: [
              {
                table: 'k8s_container',
                conditions: [
                  ['tke_cluster_instance_id', '=', route.queries['clusterId']],
                  [
                    'pod_name',
                    '=',
                    pod_name ||
                      (resourceDetailState.podList.data.records[0]
                        ? resourceDetailState.podList.data.records[0].metadata.name
                        : '')
                  ],
                  ...(isInNodeManage ? [] : [['namespace', '=', reduceNs(route.queries['np'] || 'default')]])
                ],
                fields: containerMonitorFields
              }
            ],
            groupBy: [{ value: 'container_name' }],
            instance: {
              columns: [
                {
                  key: 'container_name',
                  name: t('容器名称')
                }
              ],
              list: []
            }
          };

    this.setState({
      monitorPanelProps: {
        ...monitorPanelProps,
        title: t('Pod监控'),
        headerExtraDOM: (
          <ul className="form-list">
            <li>
              <div className="form-label">
                <label>{t('对比维度')}</label>
              </div>
              <div className="form-input">
                <div className="form-unit">
                  <div className="tc-15-rich-radio" role="radiogroup">
                    {[
                      { label: 'Pod', key: 'podMonitor' },
                      { label: 'Container', key: 'containerMonitor' }
                    ].map(item => (
                      <button
                        key={item.key}
                        onClick={e => this._handleMonitor(item.key, '')}
                        className={'tc-15-btn m ' + (monitorType === item.key ? 'checked' : '')}
                        role="radio"
                      >
                        {item.label}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            </li>
            {monitorType === 'containerMonitor' && (
              <li>
                <div className="form-label">
                  <label>{t('所属Pod')}</label>
                </div>
                <select
                  className="tc-15-select m"
                  onChange={e => {
                    this._handleMonitor(monitorType, e.target.value);
                  }}
                >
                  {resourceDetailState.podList.data.records.map(pod => (
                    <option key={pod.metadata.name} value={pod.metadata.name}>
                      {pod.metadata.name}
                    </option>
                  ))}
                </select>
              </li>
            )}
          </ul>
        )
      } as MonitorPanelProps
    });
  }
}
