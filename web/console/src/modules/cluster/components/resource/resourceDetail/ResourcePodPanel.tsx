import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';
import { Bubble, Icon, TableColumn, Text } from '@tea/component';
import { autotip } from '@tea/component/table/addons/autotip';
import { expandable } from '@tea/component/table/addons/expandable';
import { selectable } from '@tea/component/table/addons/selectable';
import { stylize } from '@tea/component/table/addons/stylize';
import { bindActionCreators, OperationState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormatter, RouteState } from '../../../../../../helpers';
import { Clip, GridTable, HeadBubble, LinkButton } from '../../../../common/components';
import { TableLayout } from '../../../../common/layouts';
import { isEmpty } from '../../../../common/utils';
import { execColumnWidth } from '../../../../common/utils/tea_adapter';
import { allActions } from '../../../actions';
import { ContainerStatusMap } from '../../../constants/Config';
import { Pod, PodContainer, ResourceFilter } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { IsInNodeManageDetail } from './ResourceDetail';
import { PodTabel } from './ResourcePodTable';
const moment = require('moment');
moment.locale('zh-CN');

/** 获取containerId，去掉前缀 docker:// */
export function reduceContainerId(containerStatus: any[], containerName: string) {
  if (containerStatus) {
    const finder = containerStatus.find(c => c.name === containerName);
    return finder.containerID ? finder.containerID.replace('docker://', '') : '-';
  }
}

/** 处理cpu request的值 */
export const ReduceRequest = (type: 'cpu' | 'memory', requests: { cpu: string; memory: string }) => {
  let finalRequest = 0;
  const requestValue: string = !isEmpty(requests) && requests[type] ? requests[type].toLowerCase() : '0';
  if (type === 'cpu') {
    // cpu 以 核作为单位
    if (requestValue.includes('m')) {
      finalRequest = parseInt(requestValue) / 1000;
    } else {
      finalRequest = parseInt(requestValue);
    }
  } else {
    // memory 需要转换为 Mib的单位
    if (requestValue.includes('k')) {
      finalRequest = parseInt(requestValue) / 1024;
    } else if (requestValue.includes('g')) {
      finalRequest = parseInt(requestValue) * 1024;
    } else {
      finalRequest = parseInt(requestValue);
    }
  }
  return finalRequest;
};

/** 判断当前的pod列表是否需要进行轮询，用于列表的轮询 */
export const IsPodShowLoadingIcon = (item: Pod) => {
  let isNeedShowLoading = false;
  if (item.status.phase !== 'Running' && item.status.phase !== 'Succeeded') {
    isNeedShowLoading = true;
  } else if (item.status.conditions && item.status.phase !== 'Succeeded') {
    item.status.conditions.forEach(item => {
      if (item.status !== 'True') {
        isNeedShowLoading = true;
      }
    });
  }
  return isNeedShowLoading;
};

interface ResourcePodPanelState {
  /** 当前需要展开的pod列表 */
  // expanded?: string[];
  expandedKeys?: string[];
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourcePodPanel extends React.Component<RootProps, ResourcePodPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      expandedKeys: []
    };
  }

  componentWillUnmount() {
    const { actions } = this.props;
    actions.resourceDetail.pod.clearPollEvent();
  }

  render() {
    return this._renderTablePanel();
  }

  /** 展示pod列表的 */
  _renderTablePanel() {
    const { actions, subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { resourceDetailState } = subRoot,
      { podQuery, podSelection, podList } = resourceDetailState;

    // 根据type 来确定 pod列表需要展示什么信息，node详情下的pod 需要展示更多的信息
    const isInNodeManage = IsInNodeManageDetail(urlParams['type']);
    const columns: TableColumn<Pod>[] = this._reduceColumns(isInNodeManage, urlParams);

    const addons = [
      stylize({
        className: 'docker-table',
        bodyClassName: 'tc-15-table-panel tc-15-table-fixed-body'
      }),
      selectable({
        value: podSelection.map(item => item.id as string),
        onChange: keys => {
          actions.resourceDetail.pod.podSelect(
            podList.data.records.filter(item => keys.indexOf(item.id as string) !== -1)
          );
        }
      })
    ];
    addons.push(
      expandable({
        expandedKeys: this.state.expandedKeys,
        onExpandedKeysChange: keys => this.setState({ expandedKeys: keys }),
        render(x) {
          return (
            <tr className="tc-detail-row">
              <td />
              <td colSpan={columns.length - 1}>
                <ResourcePodContainerTable
                  containers={x.spec.containers}
                  containerStatus={x.status.containerStatuses}
                  podId={x.id + ''}
                  podList={podList.data.records}
                  route={route}
                  isInNodeManage={isInNodeManage}
                />
              </td>
            </tr>
          );
        },
        gapCell: 1
      })
    );

    return <PodTabel columns={execColumnWidth(columns)} addons={addons} {...this.props} />;
  }

  /** 列表需要展示的项 */
  private _reduceColumns(isInNodeManage = false, urlParams) {
    const columns: TableColumn<Pod>[] = [
      {
        key: 'name',
        header: t('实例名称'),
        width: '10%',
        render: x => (
          <Bubble content={x.metadata.name || null}>
            <span id={x.id + ''} className="text-overflow m-width" style={{ maxWidth: '72%' }}>
              <a
                href="javascript:;"
                title={x.metadata.name}
                className={classnames('expander', { expanded: this._isExpanded(x.id + '') })}
                onClick={() => this._toggle(x.id + '')}
              >
                {x.metadata.name}
              </a>
            </span>
            <Clip target={'#' + x.id} className="hover-icon" />
          </Bubble>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '12%',
        render: x => this._renderPodStatus(x)
      },
      {
        key: 'image',
        header: t('镜像'),
        render: x => (
          <ul>
            {x.spec.containers.map(c => (
              <li>{c.image}</li>
            ))}
          </ul>
        )
      },
      {
        key: 'hostIP',
        header: t('实例所在节点IP'),
        width: '10%',
        render: x => (
          <div>
            <span id={'nodeIp' + x.id} className="text-overflow m-width" title={x.status.hostIP || '-'}>
              {x.status.hostIP ? x.status.hostIP : '-'}
            </span>
            <Clip target={'#nodeIp' + x.id} className="hover-icon" isShow={!!x.status.hostIP} />
          </div>
        )
      },
      {
        key: 'podIP',
        header: t('实例IP'),
        width: '10%',
        render: x => (
          <div>
            <span
              id={'podIP' + x.id}
              className="text-overflow m-width"
              title={x.status.podIP}
              style={{ maxWidth: '72%' }}
            >
              {x.status.podIP ? x.status.podIP : '-'}
            </span>
            <Clip target={'#podIP' + x.id} className="hover-icon" isShow={!!x.status.podIP} />
          </div>
        )
      },
      {
        key: 'period',
        width: '10%',
        header: column => <HeadBubble autoflow={true} title={t('运行时间')} text={t('实例从启动至今的时间')} />,
        // render: x => this._getDays(x.status.startTime)
        render: x => (
          <Bubble content={`创建时间：${moment(x.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')}`}>
            <Text>{this.runningTime(x.status.startTime)}</Text>
          </Bubble>
        )
      },
      {
        key: 'operation',
        header: t('操作'),
        width: '12%',
        render: x => this._renderOperationCell(x)
      }
    ];

    if (isInNodeManage) {
      const nodeColumns: TableColumn<Pod>[] = [
        {
          key: 'cpuRequest',
          header: 'CPU Request',
          width: '10%',
          render: x => this._renderPodCPURequest(x.spec.containers)
        },
        {
          key: 'memRequest',
          header: t('内存 Request'),
          width: '10%',
          render: x => this._renderPodMemRequest(x.spec.containers)
        },
        {
          key: 'namespace',
          header: t('命名空间'),
          width: '8%',
          render: x => (
            <Bubble content={x.metadata.namespace || '-'}>
              <Text overflow>{x.metadata.namespace || '-'}</Text>
            </Bubble>
          )
        },
        {
          key: 'workload',
          header: t('所属工作负载'),
          width: '10%',
          render: x => this._renderOwnerWorkload(x.metadata.ownerReferences, urlParams, x.metadata.namespace)
        },
        {
          key: 'restartCount',
          width: '10%',
          header: column => (
            <React.Fragment>
              <Text verticalAlign="middle">{t('重启次数')}</Text>
              <Bubble content={t('实例中容器的重启次数之和')}>
                <Icon type="info" style={{ top: '1px' }} />
              </Bubble>
            </React.Fragment>
          ),
          render: x => this._renderPodRestartCount(x.status.containerStatuses || [])
        }
      ];
      columns.splice(4, 1, ...nodeColumns);
    }

    return columns;
  }

  /** 展示所属工作负载 */
  private _renderOwnerWorkload(ownerReferences: any, urlParams, namespace: string) {
    let { route, actions, subRoot } = this.props,
      { podQuery } = subRoot.resourceDetailState;
    const ownerInfo = ownerReferences ? ownerReferences[0] : undefined;
    const workloadType = ownerInfo ? (ownerInfo['kind'] === 'ReplicaSet' ? 'Deployment' : ownerInfo['kind']) : '-';
    let workloadName;
    // 如果workloadType 是 deployment的话，需要把 最后一个 - 后面的内容去除
    if (workloadType === 'Deployment') {
      const name: string = ownerInfo['name'];
      const splitName = name.split('-');
      workloadName = splitName.slice(0, splitName.length - 1).join('-');
    } else {
      workloadName = ownerInfo ? ownerInfo['name'] : '-';
    }

    const _handleClickForNagigate = () => {
      router.navigate(
        Object.assign({}, urlParams, { type: 'resource', resourceName: workloadType.toLowerCase() }),
        Object.assign({}, route.queries, { resourceIns: workloadName, np: namespace })
      );
      // 需要重新拉取一下内容，workload的列表
      actions.resource.initResourceInfoAndFetchData(true, workloadType.toLowerCase());
      // 重新进行pod列表的拉取
      actions.resourceDetail.pod.poll(
        Object.assign({}, podQuery.filter, {
          namespace,
          specificName: workloadName
        })
      );
    };

    return (
      <Bubble
        content={
          <React.Fragment>
            <p>
              <Text>{t('工作负载名称：')}</Text>
              {workloadName !== '-' ? (
                <a
                  href="javascript:;"
                  onClick={() => {
                    _handleClickForNagigate();
                  }}
                >
                  {ownerInfo ? `${workloadName}` : '-'}
                </a>
              ) : (
                <Text>{workloadName}</Text>
              )}
            </p>
            <p>
              <Text>{t('所属工作负载：')}</Text>
              <Text>{workloadType}</Text>
            </p>
          </React.Fragment>
        }
      >
        <span className="text text-overflow m-width">{ownerInfo ? workloadName : '-'}</span>
        <div className="sl-editor-name">
          <span className="text text-overflow text-label">{workloadType}</span>
        </div>
      </Bubble>
    );
  }

  /** 展示pod下的所有的container的Cpu request */
  private _renderPodCPURequest(containers: PodContainer[]) {
    let allRequest = 0;
    containers.forEach(item => {
      allRequest += ReduceRequest('cpu', item.resources.requests);
    });
    return (
      <Text parent="div" overflow>
        {allRequest > 0 ? t('{{allRequest}} 核', { allRequest }) : t('无限制')}
      </Text>
    );
  }

  /** 展示pod下所有的container的 mem Request */
  private _renderPodMemRequest(containers: PodContainer[]) {
    let allRequest = 0;
    containers.forEach(item => {
      allRequest += ReduceRequest('memory', item.resources.requests);
    });
    return (
      <Text parent="div" overflow>
        {allRequest > 0 ? `${allRequest} M` : t('无限制')}
      </Text>
    );
  }

  /** 展示pod下的所有container的重启次数 */
  private _renderPodRestartCount(containerStatuses: any[]) {
    let content: JSX.Element;
    let restartTime = 0;
    if (containerStatuses) {
      containerStatuses.forEach(item => {
        restartTime += item.restartCount;
      });
    }
    return (
      <Text parent="div" overflow>
        {t('{{restartTime}} 次', { restartTime })}
      </Text>
    );
  }

  /** 展示pod的状态 */
  private _renderPodStatus(pod: Pod) {
    // 判断当前的pod是否完全ok，即使是running，需要判断conditions的状态来确定
    let isStatusOK = pod.status.phase === 'Running' || pod.status.phase === 'Succeeded',
      isConditionOK = true;

    // 是否存在conditions的状态，就算是running状态下，也需要进行判断condition
    const isHasCondition: boolean = pod.status.conditions ? true : false;
    // Succeeded状态下，不需要判断condition，因为已经执行完了
    if (isHasCondition && isStatusOK && pod.status.phase !== 'Succeeded') {
      pod.status.conditions.forEach(item => {
        isConditionOK = isConditionOK && item.status === 'True';
      });
    }

    function getPodStatus(pod) {
      let result = pod.status.phase;
      if (pod.status.phase === 'Running' && pod.metadata.deletionTimestamp) {
        result = <Text theme="warning">Terminating</Text>;
      }
      if (pod.status.reason) {
        result = (
          <>
            {pod.status.phase}
            <Text theme="warning">({pod.status.reason})</Text>
          </>
        );
      }
      return result;
    }
    return (
      <div>
        <span
          className={classnames('text-overflow', {
            'text-success': isStatusOK && isConditionOK,
            'text-danger': isStatusOK && !isConditionOK
          })}
        >
          {getPodStatus.call(this, pod)}
        </span>
        {(!isStatusOK || !isConditionOK) && pod.status.conditions && (
          <Bubble
            placement="left"
            className="mr20"
            content={pod.status.conditions.map((item, index) => {
              return (
                <div key={index} style={{ marginBottom: '3px' }}>
                  <span
                    className="text-label"
                    style={{ verticalAlign: 'middle', width: '90px', maxWidth: '100px', display: 'inline-block' }}
                  >{`${item.type}`}</span>
                  {item.status === 'True' ? (
                    <i className="n-success-icon" />
                  ) : (
                    <span style={{ verticalAlign: 'middle' }} className="text-danger">
                      <i className="n-error-icon" style={{ verticalAlign: 'middle' }} />
                      <span style={{ verticalAlign: 'middle', marginLeft: '5px' }}>{item.reason || item.status}</span>
                    </span>
                  )}
                </div>
              );
            })}
          >
            <i className="plaint-icon" />
          </Bubble>
        )}
        {IsPodShowLoadingIcon(pod) && (
          <i style={{ verticalAlign: 'middle', marginLeft: '5px' }} className="n-loading-icon" />
        )}
      </div>
    );
  }

  /** 展示操作选项 */
  private _renderOperationCell(pod: Pod) {
    const { actions, subRoot } = this.props,
      { deletePodFlow } = subRoot.resourceDetailState;

    const renderDeleteButton = () => {
      const isDeleting = deletePodFlow.operationState === OperationState.Performing,
        loginDisabled = isDeleting || (pod.status.phase !== 'Running' && pod.status.phase !== 'Succeeded');

      return (
        <span>
          <LinkButton
            disabled={isDeleting}
            onClick={() => {
              !isDeleting && actions.workflow.deletePod.start([]);
              actions.resourceDetail.pod.podSelect([pod]);
            }}
          >
            {t('销毁重建')}
          </LinkButton>
          <LinkButton
            disabled={loginDisabled}
            errorTip={t('当前容器状态不可登录')}
            onClick={() => {
              !loginDisabled && this._toggleRemoteLoginDialog(pod);
            }}
          >
            {t('远程登录')}
          </LinkButton>
        </span>
      );
    };
    return <div>{renderDeleteButton()}</div>;
  }

  /** 处理远程登录的按钮 */
  private _toggleRemoteLoginDialog(pod: Pod) {
    const { actions } = this.props;
    actions.resourceDetail.pod.podSelect([pod]);
    actions.resourceDetail.pod.toggleLoginDialog();
  }

  /** 展示创建的时间 */
  private _getCreateTime(createTime: string) {
    const time = dateFormatter(new Date(createTime), 'YYYY-MM-DD HH:mm:ss');
    const [first, second] = time.split(' ');

    return <Text>{`${first} ${second}`}</Text>;
  }

  private runningTime = (startTime: string) => moment.duration(moment().diff(moment(startTime))).humanize(true);

  /** 运行时间 */
  private _getDays(startTime: string) {
    let content = '';

    const create = Date.parse(startTime);
    const now = Date.now();
    if (!startTime) {
      content = '-';
    } else if (now < create) {
      content = '0d 0h';
    } else {
      const days = Math.floor((now - create) / (1000 * 3600 * 24)) + 'd',
        hours = (Math.floor((now - create) / (1000 * 3600)) % 24) + 'h';
      content = days + ' ' + hours;
    }

    return (
      <Text parent="div" overflow>
        {content}
      </Text>
    );
  }

  /** 判断当前实例列表是否需要展开 */
  private _isExpanded(id: string) {
    return this.state.expandedKeys.indexOf(id) > -1;
  }

  /** 开关当前的列表展开与否 */
  private _toggle(id: string) {
    if (this._isExpanded(id)) {
      this._collapse(id);
    } else {
      this._expand(id);
    }
  }

  /** 增加展开的列表项 */
  private _expand(id: string) {
    this.setState({
      expandedKeys: [...this.state.expandedKeys, id]
    });
  }

  /** 删除展开的列表项 */
  private _collapse(id: string) {
    this.setState({
      // expanded: this.state.expanded.filter(x => x !== instanceId)
      expandedKeys: this.state.expandedKeys.filter(x => x !== id)
    });
  }
}

/** ====== 这里是展示 rowAppend进去的一些内容，主要是 container的信息 */
interface ResourcePodContainerTableProps {
  /** container */
  containers?: PodContainer[];

  containerStatus?: any[];

  /** 地域id */
  route?: RouteState;

  /** podId */
  podId?: string;

  /** podList */
  podList?: Pod[];

  /** isInNodeManage */
  isInNodeManage?: boolean;
}

class ResourcePodContainerTable extends React.Component<ResourcePodContainerTableProps, {}> {
  render() {
    const { containers, containerStatus = [], isInNodeManage = false } = this.props;

    let basicColgroup: JSX.Element;

    // node详情页的colGroup
    if (isInNodeManage) {
      basicColgroup = (
        <colgroup>
          <col />
          <col />
          <col width="25%" />
          <col />
          <col />
          <col />
          <col />
          <col />
          <col />
        </colgroup>
      );
    } else {
      basicColgroup = (
        <colgroup>
          <col />
          <col />
          <col width="40%" />
          <col width="10%" />
        </colgroup>
      );
    }

    return (
      <div className="tc-15-table-panel">
        <div className="tc-15-table-fixed-head">
          <table className="tc-15-table-box">
            {basicColgroup}
            <thead>
              <tr>
                <th>
                  <Text parent="div" overflow>
                    {t('容器名称')}
                  </Text>
                </th>
                <th>
                  <Text parent="div" overflow>
                    {t('容器ID')}
                  </Text>
                </th>
                <th>
                  <Text parent="div" overflow>
                    {t('镜像版本号')}
                  </Text>
                </th>

                {/* start 这里只有在node详情页里面的container列表才需要展示 */}
                {isInNodeManage && (
                  <th>
                    <Text parent="div" overflow>
                      CPU Request
                    </Text>
                  </th>
                )}
                {isInNodeManage && (
                  <th>
                    <Text parent="div" overflow>
                      CPU Limit
                    </Text>
                  </th>
                )}
                {isInNodeManage && (
                  <th>
                    <Text parent="div" overflow>
                      {t('内存 Request')}
                    </Text>
                  </th>
                )}
                {isInNodeManage && (
                  <th>
                    <Text parent="div" overflow>
                      {t('内存 Limit')}
                    </Text>
                  </th>
                )}
                {isInNodeManage && (
                  <th>
                    <Text parent="div" overflow>
                      {t('重启次数')}
                    </Text>
                  </th>
                )}
                {/* start 这里只有在node详情页里面的container列表才需要展示 */}

                <th>
                  <Text parent="div" overflow>
                    {t('状态')}
                  </Text>
                </th>
              </tr>
            </thead>
          </table>
        </div>
        <div className="tc-15-table-fixed-body">
          <table className="tc-15-table-box tc-15-table-rowhover">
            {basicColgroup}
            <tbody>
              {isEmpty(containers) ? (
                <tr>
                  <td colSpan={isInNodeManage ? 8 : 4} className="text-center">
                    {t('该实例下容器列表为空，请切换实例或更新该实例')}
                  </td>
                </tr>
              ) : (
                containers.map((container, index) => (
                  <tr key={index}>
                    <td>
                      <Text parent="div" overflow>
                        {container.name}
                      </Text>
                    </td>

                    <td>
                      <div>
                        <span id={'cId' + container.id} className="text-overflow m-width" style={{ maxWidth: '74%' }}>
                          {containerStatus.length ? reduceContainerId(containerStatus, container.name) : '-'}
                        </span>
                        <Clip target={'#cId' + container.id} className="hover-icon" />
                      </div>
                    </td>

                    <td>
                      <Bubble placement="bottom" style={{ width: '100%' }} content={container.image || null}>
                        <Text parent="div" overflow>
                          {container.image}
                        </Text>
                      </Bubble>
                    </td>

                    {/* start 这里只有在node详情里面的container列表才需要展示 */}
                    {isInNodeManage && (
                      <td>
                        <Text parent="div" overflow>
                          {this._renderCPUAndMemory(t('核'), ReduceRequest('cpu', container.resources.requests))}
                        </Text>
                      </td>
                    )}
                    {isInNodeManage && (
                      <td>
                        <Text parent="div" overflow>
                          {this._renderCPUAndMemory(t('核'), ReduceRequest('cpu', container.resources.limits))}
                        </Text>
                      </td>
                    )}
                    {isInNodeManage && (
                      <td>
                        <Text parent="div" overflow>
                          {this._renderCPUAndMemory('M', ReduceRequest('memory', container.resources.requests))}
                        </Text>
                      </td>
                    )}
                    {isInNodeManage && (
                      <td>
                        <Text parent="div" overflow>
                          {this._renderCPUAndMemory('M', ReduceRequest('memory', container.resources.limits))}
                        </Text>
                      </td>
                    )}
                    {isInNodeManage && (
                      <td>
                        <Text parent="div" overflow>
                          {t('{{count}} 次', {
                            count: containerStatus
                              ? 0
                              : containerStatus.find(item => item.name === container.name).restartCount
                          })}
                        </Text>
                      </td>
                    )}
                    {/* end 这里只有在node详情里面的container列表才需要展示 */}

                    <td>{this._reduceContainerStatus(container)}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    );
  }

  /** 展示cpu 和 memory的数据 */
  private _renderCPUAndMemory(unit: string, data: number) {
    return data > 0 ? data + unit : t('无限制');
  }

  // 展示容器的状态
  private _reduceContainerStatus(con: PodContainer) {
    const { containerStatus = [] } = this.props;
    const finder = containerStatus.find(c => c.name === con.name);
    const statusKey = finder && Object.keys(finder.state)[0];

    let content: JSX.Element;

    if (finder) {
      content = (
        <div>
          <span
            className={classnames(
              'text-overflow',
              ContainerStatusMap[statusKey] && ContainerStatusMap[statusKey].classname
            )}
          >
            {ContainerStatusMap[statusKey] ? ContainerStatusMap[statusKey].text : '-'}
          </span>
          {statusKey && statusKey !== 'running' && (
            <Bubble placement="right" className="mr20" content={finder.state[statusKey].reason || null}>
              <div className="tc-15-bubble-icon">
                <i className="tc-icon icon-what" />
              </div>
            </Bubble>
          )}
        </div>
      );
    } else {
      content = (
        <div>
          <span
            className={classnames(
              'text-overflow',
              ContainerStatusMap[statusKey] && ContainerStatusMap[statusKey].classname
            )}
          >
            {ContainerStatusMap[statusKey] ? ContainerStatusMap[statusKey].text : '-'}
          </span>
        </div>
      );
    }

    return content;
  }
}
