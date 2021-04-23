import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Text } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormatter, formatMemory } from '../../../../../../helpers';
import { Clip, ListItem } from '../../../../common/components';
import { DetailLayout } from '../../../../common/layouts';
import { allActions } from '../../../actions';
import { Computer } from '../../../models';
import { RootProps } from '../../ClusterApp';
import { ReduceRequest } from './ResourcePodPanel';

const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceNodeDetailPanel extends React.Component<RootProps, {}> {
  render() {
    const { subRoot, route } = this.props,
      { resourceDetailState } = subRoot,
      { resourceDetailInfo } = resourceDetailState;

    // 当前选中的node节点
    const resourceIns = resourceDetailInfo.selection;

    // 当前的地域
    const regionId = route.queries['rid'];

    // 获取当前机器的配置

    return resourceIns === undefined ? (
      loadingElement
    ) : (
      <React.Fragment>
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('主机信息')}</h3>
            </div>
            <div className="param-bd">
              <ul className="item-descr-list">
                <ListItem label={t('节点名')}>
                  <Text>{resourceIns.metadata.name}</Text>
                </ListItem>

                {this._renderNodeStatus(resourceIns.status.conditions)}

                {this._renderComputerConfig(resourceIns?.status?.capacity)}

                {this._renderIPAddress(resourceIns.status.addresses)}

                <ListItem label={t('操作系统')}>
                  <span className="text">{resourceIns.status.nodeInfo['osImage'] || '-'}</span>
                </ListItem>

                <ListItem label={t('创建时间')}>
                  <span className="text">
                    {dateFormatter(new Date(resourceIns.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
                  </span>
                </ListItem>
              </ul>
            </div>
          </div>
        </DetailLayout>

        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('Kubernetes信息')}</h3>
            </div>
            <div className="param-bd">
              <ul className="item-descr-list">
                {this._renderKvData(resourceIns.metadata.labels, 'Kubernetes Labels')}

                {this._renderKvData(resourceIns.metadata.annotations, 'Kubernetes Annotations')}

                {resourceIns.spec.taints && this._renderNodeTaint(resourceIns.spec.taints, 'Kubernetes Taints')}
                {/* {this._renderAlreadyAllocatedResource(computerInfo, isNeedLoadingComputerInfo)} */}

                {this._renderAllocatableResource(resourceIns.status.allocatable)}

                <ListItem label={t('podCIDR')}>
                  <span className="text">{resourceIns.spec['podCIDR'] || '-'}</span>
                </ListItem>

                <ListItem label={t('容器运行时版本')}>
                  <span className="text">{resourceIns.status.nodeInfo['containerRuntimeVersion'] || '-'}</span>
                </ListItem>

                <ListItem label={t('Kubelet版本')}>
                  <span className="text">{resourceIns.status.nodeInfo['kubeletVersion'] || '-'}</span>
                </ListItem>

                <ListItem label={t('KubeProxy版本')}>
                  <span className="text">{resourceIns.status.nodeInfo['kubeProxyVersion'] || '-'}</span>
                </ListItem>
              </ul>
            </div>
          </div>
        </DetailLayout>
      </React.Fragment>
    );
  }

  /** 处理节点的状态展示 */
  private _renderNodeStatus(conditions: any[]) {
    const nodeCondition = conditions.find(item => item.type === 'Ready');

    const isNodeReady = nodeCondition.status === 'True' ? true : false;

    return (
      <ListItem label={t('状态')}>
        <Text theme={isNodeReady ? 'success' : 'danger'}>{isNodeReady ? '健康' : '异常'}</Text>
        {!isNodeReady && (
          <Bubble placement="bottom" content={nodeCondition.reason || t('未知异常')}>
            <i className="plaint-icon" />
          </Bubble>
        )}
      </ListItem>
    );
  }

  /** 处理ip地址的展示 */
  private _renderIPAddress(address: any[]) {
    let externalIP = '',
      internalIP = '';

    address.forEach(item => {
      if (item.type === 'InternalIP') {
        internalIP = item.address;
      } else if (item.type === 'ExternalIP') {
        externalIP = item.address;
      }
    });

    return (
      <ListItem label={t('IP地址')}>
        <div>
          <span className="text" id="detailExternalId">
            {externalIP || '-'}
          </span>
          <span className="text">{t(' (外网)')}</span>
          <Clip target={`#detailExternalId`} />
        </div>
        <div>
          <span className="text" id="detailInternalId">
            {internalIP || '-'}
          </span>
          <span className="text">{t(' (内网)')}</span>
          <Clip target={`#detailInternalId`} />
        </div>
      </ListItem>
    );
  }

  /** 处理key: value的展示
   * @param showData: string  展示的数据
   * @param label: string ListItem展示的数据
   */
  private _renderKvData(showData: any, label: string) {
    const keys = Object.keys(showData);

    return (
      <ListItem label={label}>
        {keys.length ? (
          keys.map((item, index) => {
            return <p key={index} className="text">{`${item}：${showData[item]}`}</p>;
          })
        ) : (
          <p className="text">{t('无')}</p>
        )}
      </ListItem>
    );
  }

  /** 处理Taint的展示
   * @param showData: string  展示的数据
   * @param label: string ListItem展示的数据
   */
  private _renderNodeTaint(showData: any, label: string) {
    return (
      <ListItem label={label}>
        {showData.length ? (
          showData.map((item, index) => {
            return (
              <p key={index} className="text">{`${item.key}${item.value ? '=' + item.value : ''}：${item.effect}`}</p>
            );
          })
        ) : (
          <p className="text">{t('无')}</p>
        )}
      </ListItem>
    );
  }

  /** 展示总可分配资源
   * @param label:string  listItem的label战术
   * @param allocatable:{}  需要处理的具体数据
   */
  private _renderAllocatableResource(allocatable: { cpu: string; memory: string }) {
    const finalCpu = ReduceRequest('cpu', allocatable),
      finalMem = (ReduceRequest('memory', allocatable) / 1024).toFixed(2);

    return (
      <ListItem label={t('总可分配资源')}>
        <span className="text">
          <span className="text-label">{`CPU: `}</span>
          <span>
            {t('{{count}} 核，', {
              count: finalCpu
            })}
          </span>
        </span>
        <span className="text">
          <span className="text-label">{t('内存: ')}</span>
          <span>{`${finalMem} GB`}</span>
        </span>
      </ListItem>
    );
  }

  /** 展示机器的配置 */
  private _renderComputerConfig(capacityConfig: any) {
    const capacity = {
      cpu: capacityConfig?.cpu,
      memory: capacityConfig?.memory
    };

    return (
      <ListItem label={t('配置')}>
        <Text verticalAlign="middle" theme="label">{`CPU: `}</Text>
        <Text verticalAlign="middle">{capacityConfig?.cpu ?? '-'} 核</Text>
        <Text verticalAlign="middle" theme="label">
          {t('内存: ')}
        </Text>
        <Text verticalAlign="middle">{formatMemory(capacity?.memory ?? '0', 'Gi')}</Text>
        <Text verticalAlign="middle" theme="label">{`Pods: `}</Text>
        <Text verticalAlign="middle">{capacityConfig?.pods || 0}</Text>
      </ListItem>
    );
  }
}
