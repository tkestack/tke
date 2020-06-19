import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import { resourceLimitTypeToText, resourceTypeToUnit } from '@src/modules/project/constants/Config';
import { Bubble, TableColumn, Text } from '@tea/component';
import { selectable } from '@tea/component/table/addons/selectable';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../../../helpers';
import { Clip, HeadBubble, LinkButton } from '../../../../common/components';
import { DisplayFiledProps, OperatorProps } from '../../../../common/models';
import { includes, isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { ResourceLoadingIcon, ResourceStatus } from '../../../constants/Config';
import { Resource } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

/** 判断resource是否需要展示loading状态
 * @param resourceName: string  资源的名称，如deployment
 * @param item: Resource 当前实例
 */
export const IsResourceShowLoadingIcon = (resourceName: string, item: Resource) => {
  if (resourceName === 'np' && includes(ResourceLoadingIcon['npDelete'], (item.status && item.status.phase) || '')) {
    return true;
  } else if (resourceName === 'deployment') {
    // 判断readyReplicas的前提是 replicas是有的，而不是空，空即代表当前没有pod，不需要继续轮训
    return item.status.readyReplicas === undefined && item.status.replicas
      ? true
      : +item.status.readyReplicas < +item.status.replicas
      ? true
      : false;
  } else if (resourceName === 'svc') {
    let type = item.spec && item.spec.type,
      isClusterIP = type === 'ClusterIP',
      isNodePort = type === 'NodePort';
    return (isClusterIP && item.status && item.status.loadBalancer && item.status.loadBalancer.ingress === undefined) ||
      (isNodePort && item.status && item.status.loadBalancer && item.status.loadBalancer.ingress === undefined) ||
      (!isClusterIP && !isNodePort && item.status && item.status.loadBalancer && item.status.loadBalancer.ingress)
      ? false
      : true;
  } else if (resourceName === 'ingress') {
    // 此处需要判断是否为nginx-ingress，如果为nginx-ingress的话，则不需要进行轮询了
    let isNginxIngress =
      item['metadata']['annotations'] &&
      item['metadata']['annotations']['kubernetes.io/ingress.class'] === 'nginx-ingress'
        ? true
        : false;

    // 判断是否为qcloud-loadbalance
    let isQcloud =
      item.metadata.annotations && item.metadata.annotations['kubernetes.io/ingress.qcloud-loadbalance-id'];

    return isNginxIngress
      ? false
      : isQcloud
      ? item.metadata.annotations &&
        item.metadata.annotations['kubernetes.io/ingress.qcloud-loadbalance-id'] &&
        item.status.loadBalancer &&
        item.status.loadBalancer.ingress
        ? false
        : true
      : false;
  } else if (resourceName === 'pvc') {
    return item.status === undefined ? true : item.status.phase === 'Pending' ? true : false;
  } else if (resourceName === 'statefulset') {
    return item.status.replicas === undefined ? true : +item.status.replicas < +item.spec.replicas ? true : false;
  } else if (resourceName === 'daemonset') {
    return item.status.currentNumberScheduled === undefined
      ? true
      : +item.status.currentNumberScheduled < +item.status.desiredNumberScheduled
      ? true
      : false;
  } else if (resourceName === 'tapp') {
    /**待定 todo ing */
    return item.status.replicas === undefined
      ? true
      : +item.status.readyReplicas < +item.status.replicas
      ? true
      : false;
  }
  return false;
};

/** loading的样式 */

const loadingElement: JSX.Element = (
  <i style={{ verticalAlign: 'middle', marginLeft: '5px' }} className="n-loading-icon" />
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ResourceTablePanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    // 离开页面的话，清空当前的轮询操作
    actions.resource.clearPolling();
    // 离开页面的话，清空当前的多选
    actions.resource.selectMultipleResource([]);
  }

  render() {
    return this._renderTablePanel();
  }

  /** 展示普通的text */
  private _reduceText(showData: any, fieldInfo: DisplayFiledProps, resource: Resource, clipId: string) {
    let showContent;
    if (fieldInfo.isLink) {
      if (fieldInfo.isClip) {
        showContent = (
          <Bubble
            content={
              <Text parent="p" nowrap={false}>
                {showData}
              </Text>
            }
          >
            <Text overflow>
              <a
                id={clipId}
                href="javascript:;"
                onClick={e => {
                  this._handleClickForNavigate(resource);
                }}
              >
                {showData}
              </a>
            </Text>
            <Clip target={`#${clipId}`} />
          </Bubble>
        );
      } else {
        showContent = (
          <Bubble
            content={
              <Text parent="p" nowrap={false}>
                {showData}
              </Text>
            }
          >
            <Text overflow>
              <a
                href="javascript:;"
                onClick={e => {
                  this._handleClickForNavigate(resource);
                }}
              >
                {showData}
              </a>
            </Text>
          </Bubble>
        );
      }
    } else {
      if (fieldInfo.isClip) {
        showContent = (
          <Bubble
            content={
              <Text parent="p" nowrap={false}>
                {showData}
              </Text>
            }
          >
            <Text overflow id={clipId}>
              {showData}
            </Text>
            <Clip target={`#${clipId}`} />
          </Bubble>
        );
      } else {
        showContent = (
          <Bubble
            content={
              <Text parent="p" nowrap={false}>
                {showData}
              </Text>
            }
          >
            <Text overflow>{showData}</Text>
          </Bubble>
        );
      }
    }
    return showContent;
  }

  /** 展示labels的形式 */
  private _reduceLabelsForData(labels, direction: 'top' | 'bottom') {
    let showData = '',
      keys,
      isNoLabels = false;

    // 如果不是数组，showData就是Labels本身
    if (typeof labels === 'string') {
      showData = labels;
      isNoLabels = true;
    } else {
      keys = Object.keys(labels);
      keys.forEach((item, index) => {
        showData += item + ':' + labels[item];
        if (index !== keys.length - 1) {
          showData += '、';
        }
      });
    }

    return (
      <Bubble
        placement={direction}
        content={!isNoLabels ? keys.map((label, index) => <p key={index}>{`${label}:${labels[label]}`}</p>) : null}
      >
        <Text overflow>{showData}</Text>
      </Bubble>
    );
  }

  /** 获取操作按钮列表
   * @param resource: 每一行的数据本身
   * @param fieldInfo: 配置文件
   */
  private _renderOperationCell(resource: Resource, fieldInfo: DisplayFiledProps) {
    let { route, actions, subRoot, namespaceSelection } = this.props,
      urlParams = router.resolve(route),
      { clusterId } = route.queries,
      { resourceOption, resourceName } = subRoot,
      { ffResourceList } = resourceOption;

    // 操作列表的list
    let operatorList = fieldInfo.operatorList;
    // 更多按钮的 pop方向
    let resourceIndex = ffResourceList.list.data.records.findIndex(c => c.id === resource.id);
    let direction: 'down' | 'up' =
      resourceIndex < ffResourceList.list.data.recordCount - 2 || ffResourceList.list.data.recordCount < 4
        ? 'down'
        : 'up';

    /** 编辑yaml的按钮 */
    const renderModifyButton = (operator: OperatorProps) => {
      let disabled = false,
        errorTip = '';

      if (namespaceSelection === 'kube-system') {
        disabled = true;
        errorTip = t('当前Namespace下的资源不可编辑YAML，如需查看YAML，请前往详情页');
      } else if (resourceName === 'svc' && namespaceSelection !== 'kube-system') {
        //当资源为 servcie的时候，编辑YAML按钮的一些操作是不允许的
        disabled = resource['metadata']['name'] === 'kubernetes';
        errorTip = t('系统默认的Service不可进行此操作');
      }

      return (
        <LinkButton
          tipDirection={'left'}
          errorTip={errorTip}
          disabled={disabled}
          onClick={() => {
            if (!disabled) {
              actions.resource.select(resource);
              router.navigate(
                Object.assign({}, urlParams, { mode: 'modify' }),
                Object.assign({}, route.queries, {
                  resourceIns: resource.metadata.name
                })
              );
            }
          }}
        >
          {operator.name}
        </LinkButton>
      );
    };

    /** 删除的按钮 */
    const renderDeleteButton = (operator: OperatorProps) => {
      let disabled = false;
      let errorTip = '';

      //当资源为命名空间的时候，删除按钮的一些操作
      if (clusterId === 'cls-wbwpj79f') {
        disabled = true;
        errorTip = t('global集群资源不可删除');
      } else if (resourceName === 'np') {
        disabled = resource['metadata']['name'].indexOf('kube-') >= 0 || resource['metadata']['name'] === 'default';
        errorTip = t('命名空间不可删除');
      } else if (resourceName === 'svc' && namespaceSelection !== 'kube-system') {
        //当资源为 servcie的时候，删除按钮的一些操作是不允许的
        disabled = resource['metadata']['name'] === 'kubernetes';
        errorTip = t('系统默认的Service不可删除');
      } else {
        disabled = namespaceSelection === 'kube-system';
        errorTip = t('当前Namespace下的资源不可删除');
      }

      return (
        <LinkButton
          tipDirection={'left'}
          errorTip={errorTip}
          disabled={disabled}
          onClick={() => {
            if (!disabled) {
              actions.resource.selectDeleteResource([resource]);
              actions.workflow.deleteResource.start([]);
            }
          }}
        >
          {operator.name}
        </LinkButton>
      );
    };

    /**
     * 适用的操作
     * 1. 更新访问方式 —— Service
     * 2. 更新转发配置 —— Ingress，注意：nginx-ingress 不支持该操作
     * 3. 滚动更新镜像 —— Deployment、StatefulSet、Daemonset
     */
    const renderUpdateResourcePart = (operator: OperatorProps) => {
      let disabled = false,
        errorTip = '';

      // 这里是一些操作的信息的判断条件，disable的属性
      if (operator.actionType !== 'modifyPod' && namespaceSelection === 'kube-system') {
        disabled = true;
        errorTip = t('当前Namespace下的不可进行此操作');
      } else if (resourceName === 'svc' && namespaceSelection !== 'kube-system') {
        //当资源为 servcie的时候，删除按钮的一些操作是不允许的
        disabled = resource['metadata']['name'] === 'kubernetes';
        errorTip = t('系统默认的Service不可进行此操作');
      } else if (
        resourceName === 'ingress' &&
        resource['metadata']['annotations'] &&
        resource['metadata']['annotations']['kubernetes.io/ingress.class'] === 'nginx-ingress'
      ) {
        // 当ingress 问 nginx-ingress的时候，不支持可视化更新转发配置
        disabled = true;
        errorTip = t('nginx-ingress暂不支持此操作');
      }

      return (
        <LinkButton
          tipDirection={'left'}
          errorTip={errorTip}
          disabled={disabled}
          onClick={() => {
            if (!disabled) {
              actions.resource.select(resource);
              router.navigate(
                Object.assign({}, urlParams, {
                  mode: 'update',
                  tab: operator.actionType
                }),
                Object.assign({}, route.queries, {
                  resourceIns: resource.metadata.name
                })
              );
            }
          }}
        >
          {operator.name}
        </LinkButton>
      );
    };

    let btns = [];
    operatorList.forEach(operatorItem => {
      if (operatorItem.actionType === 'modify') {
        btns.push(renderModifyButton(operatorItem));
      } else if (operatorItem.actionType === 'delete') {
        btns.push(renderDeleteButton(operatorItem));
      } else if (
        operatorItem.actionType === 'modifyPod' ||
        operatorItem.actionType === 'modifyRule' ||
        operatorItem.actionType === 'modifyType' ||
        operatorItem.actionType === 'modifyRegistry' ||
        operatorItem.actionType === 'createBG' ||
        operatorItem.actionType === 'updateBG'
      ) {
        btns.push(renderUpdateResourcePart(operatorItem));
      }
    });
    return btns;
  }

  /** 展示ip的内容 */
  private _reduceIPCell(ipInfo: any, clipId: string, resource: Resource) {
    let { resourceName } = this.props.subRoot;
    let ipArray = ipInfo;
    // 如果ipArray 不是一个数组
    if (typeof ipArray !== 'object') {
      ipArray = [ipArray];
    }

    let isNginxIngress = false;
    // 此处需要判断是否为nginx-ingress
    if (resourceName === 'ingress') {
      isNginxIngress =
        resource['metadata']['annotations'] &&
        resource['metadata']['annotations']['kubernetes.io/ingress.class'] === 'nginx-ingress'
          ? true
          : false;
    }

    let content: JSX.Element;

    if (isNginxIngress) {
      content = <div className="ip-cell">-</div>;
    } else {
      let [clusterIP, ingressIP] = ipArray;
      let isShowLoading = IsResourceShowLoadingIcon(resourceName, resource);
      content = (
        <div className="ip-cell">
          {ingressIP && (
            <div className="sl-editor-name">
              <span className="text-overflow m-width" title={t('负载均衡IP：') + ingressIP} id={`${clipId}ingress`}>
                {ingressIP}
              </span>
              {resource.spec && resource.spec.type !== 'LoadBalancer' && isShowLoading && (
                <Bubble placement="bottom" content={t('删除中')}>
                  <i className="icon-what" />
                </Bubble>
              )}
              {isShowLoading && loadingElement}
              {!isShowLoading && ingressIP !== '-' && <Clip target={`#${clipId}ingress`} />}
            </div>
          )}
          <div>
            <span className="text-overflow m-width" title={t('服务IP：') + clusterIP} id={`${clipId}cluster`}>
              {clusterIP}
            </span>
            {clusterIP === '-' && isShowLoading && loadingElement}
            {clusterIP !== '-' && <Clip target={`#${clipId}cluster`} />}
          </div>
        </div>
      );
    }
    return content;
  }

  /** 展示status */
  private _reduceStatus(showData: any, resource: Resource) {
    let { resourceName } = this.props.subRoot;

    let statusMap = ResourceStatus[resourceName];

    return (
      <div>
        {statusMap ? (
          <p className={classnames('text-overflow', statusMap[showData] && statusMap[showData].classname)}>
            <span style={{ verticalAlign: 'middle' }}>{(statusMap[showData] && statusMap[showData].text) || '-'}</span>
            {IsResourceShowLoadingIcon(resourceName, resource) && loadingElement}
          </p>
        ) : (
          <Text parent="div" overflow>
            -
          </Text>
        )}
      </div>
    );
  }

  /** 展示映射的字段 */
  private _reduceMapText(showData: any, fieldInfo: DisplayFiledProps) {
    let { mapTextConfig } = fieldInfo;

    return (
      <Text parent="div" overflow>
        {mapTextConfig[showData] || fieldInfo.noExsitedValue}
      </Text>
    );
  }

  /** 展示副本的相关 */
  private _reduceReplicas(showData: any, resource: Resource) {
    let { resourceName } = this.props.subRoot;

    return (
      <Text parent="div" overflow>
        <span style={{ verticalAlign: 'middle' }}>{`${showData[0]}/${showData[1]}`}</span>
        {resource.status !== undefined && IsResourceShowLoadingIcon(resourceName, resource) && loadingElement}
      </Text>
    );
  }

  /** 展示ingress的后端服务 */
  private _reduceIngressRule_tke(showData: any, resource: Resource) {
    let httpRules =
        showData['kubernetes.io/ingress.http-rules'] && showData['kubernetes.io/ingress.http-rules'] !== 'null'
          ? JSON.parse(showData['kubernetes.io/ingress.http-rules'])
          : [],
      httpsRules =
        showData['kubernetes.io/ingress.https-rules'] && showData['kubernetes.io/ingress.https-rules'] !== 'null'
          ? JSON.parse(showData['kubernetes.io/ingress.https-rules'])
          : [];

    httpRules = httpRules.map(item => Object.assign({}, item, { protocol: 'http' }));
    httpsRules = httpsRules.map(item => Object.assign({}, item, { protocol: 'https' }));

    const getDomain = rule => {
      return `${rule.protocol}://${
        rule.host ? rule.host : resource.status.loadBalancer.ingress ? resource.status.loadBalancer.ingress[0].ip : '-'
      }${rule.path}`;
    };

    let finalRules = [...httpRules, ...httpsRules];

    let finalRulesLength = finalRules.length;
    return finalRules.length ? (
      <Bubble
        placement="top"
        content={finalRules.map((rule, index) => (
          <p key={index}>
            <span style={{ verticalAlign: 'middle' }}>{getDomain(rule)}</span>
            <span style={{ verticalAlign: 'middle' }}>{`-->`}</span>
            <span style={{ verticalAlign: 'middle' }}>{rule.backend.serviceName + ':' + rule.backend.servicePort}</span>
          </p>
        ))}
      >
        <p className="text-overflow" style={{ fontSize: '12px' }}>
          <span style={{ verticalAlign: 'middle' }}>{getDomain(finalRules[0])}</span>
          <span style={{ verticalAlign: 'middle' }}>{`-->`}</span>
          <span style={{ verticalAlign: 'middle' }}>
            {finalRules[0].backend.serviceName + ':' + finalRules[0].backend.servicePort}
          </span>
        </p>
        {finalRules.length > 1 && (
          <p className="text">
            <a href="javascript:;">
              <Trans count={finalRulesLength}>等{{ finalRulesLength }}条转发规则</Trans>
            </a>
          </p>
        )}
      </Bubble>
    ) : (
      <p className="text-overflow text">{t('无')}</p>
    );
  }

  /** 展示ingress的后端服务 */
  private _reduceIngressRule_standalone(showData: any) {
    let httpRules = showData !== '-' ? showData : [];
    let finalRules = httpRules.map(item => {
      return {
        protocol: 'http',
        host: item.host,
        path: item.http.paths[0].path || '',
        backend: item.http.paths[0].backend
      };
    });

    const getDomain = rule => {
      return `${rule.protocol}://${rule.host}${rule.path}`;
    };

    let finalRulesLength = finalRules.length;
    return finalRules.length ? (
      <Bubble
        placement="top"
        content={finalRules.map((rule, index) => (
          <p key={index}>
            <span style={{ verticalAlign: 'middle' }}>{getDomain(finalRules[0])}</span>
            <span style={{ verticalAlign: 'middle' }}>{`-->`}</span>
            <span style={{ verticalAlign: 'middle' }}>
              {finalRules[0].backend.serviceName + ':' + finalRules[0].backend.servicePort}
            </span>
          </p>
        ))}
      >
        <p className="text-overflow" style={{ fontSize: '12px' }}>
          <span style={{ verticalAlign: 'middle' }}>{getDomain(finalRules[0])}</span>
          <span style={{ verticalAlign: 'middle' }}>{`-->`}</span>
          <span style={{ verticalAlign: 'middle' }}>
            {finalRules[0].backend.serviceName + ':' + finalRules[0].backend.servicePort}
          </span>
        </p>
        {finalRules.length > 1 && (
          <p className="text">
            <a href="javascript:;">
              <Trans count={finalRulesLength}>等{{ finalRulesLength }}条转发规则</Trans>
            </a>
          </p>
        )}
      </Bubble>
    ) : (
      <p className="text-overflow text">{t('无')}</p>
    );
  }
  private _reducebackendGroups(showData) {
    let backendGroups = showData,
      backendGroupsLength = backendGroups !== '-' ? backendGroups.length : 0;
    return backendGroupsLength ? (
      <Bubble
        placement="right"
        content={backendGroups.map((backendGroup, index) => (
          <p key={index}>
            <span style={{ verticalAlign: 'middle' }}>{backendGroup.name}</span>
          </p>
        ))}
      >
        <p className="text-overflow" style={{ fontSize: '12px' }}>
          <span style={{ verticalAlign: 'middle' }}>{backendGroups[0].name}</span>
        </p>
        {backendGroupsLength > 1 && (
          <p className="text">
            <a href="javascript:;">
              <Trans count={backendGroupsLength}>等{{ backendGroupsLength }}条后端配置</Trans>
            </a>
          </p>
        )}
      </Bubble>
    ) : (
      <p className="text-overflow text">{t('无')}</p>
    );
  }

  /** 展示时间 */
  private _reduceTime(showData: any, direction: 'bottom' | 'top') {
    let time = dateFormatter(new Date(showData), 'YYYY-MM-DD HH:mm:ss');

    let [year, currentTime] = time.split(' ');

    return (
      <Bubble placement="left" content={time || null}>
        <p className="text-overflow">{year}</p>
        <p className="sl-editor-name text-overflow">{currentTime}</p>
      </Bubble>
    );
  }

  private _reduceResourceLimit(showData) {
    let resourceLimitKeys = showData !== '-' ? Object.keys(showData) : [];
    let content = resourceLimitKeys.map((item, index) => (
      <Text parent="p" key={index}>{`${resourceLimitTypeToText[item]}:${
        resourceTypeToUnit[item] === 'MiB'
          ? valueLabels1024(showData[item], K8SUNIT.Mi)
          : valueLabels1000(showData[item], K8SUNIT.unit)
      }${resourceTypeToUnit[item]}`}</Text>
    ));
    return (
      <Bubble placement="left" content={content}>
        {content.filter((item, index) => index < 3)}
      </Bubble>
    );
  }
  /** 获取最终展示的数据 */
  private _getFinalData(dataFieldIns, resource: Resource) {
    let result = resource;

    for (let index = 0; index < dataFieldIns.length; index++) {
      // 如果result不为一个 Object，则遍历结束
      if (typeof result !== 'object') {
        break;
      }
      result = result[dataFieldIns[index]]; // 这里做一下处理，防止因为配错找不到
    }

    // 返回空值，是因为如果不存在值，则使用配置文件的默认值
    return result || '';
  }

  /** 根据 fieldInfo的 dataFormat来决定显示的bodyCell的具体内容 */
  private _renderBodyCell(resource: Resource, fieldInfo: DisplayFiledProps, clipId: string) {
    let { subRoot } = this.props,
      { resourceOption } = subRoot,
      { ffResourceList } = resourceOption;

    let content;

    // fieldInfo当中的 dataField是一个数组，可以同时输入多个值
    let showData: any = [];
    fieldInfo.dataField.forEach(item => {
      let dataFieldIns = item.split('.');
      let data: any = this._getFinalData(dataFieldIns, resource);
      // 如果返回的为 '' ，即找不到这个对象，则使用配置文件所设定的默认值
      showData.push(data === '' ? fieldInfo.noExsitedValue : data);
    });

    showData = showData.length === 1 ? showData[0] : showData;

    // 这里是当列表有 bubble等情况的时候，判断当前行属于第几行
    let resourceIndex = ffResourceList.list.data.records.findIndex(item => item.id === resource.id);
    let direction: 'top' | 'bottom' =
      ffResourceList.list.data.recordCount < 4 || resourceIndex < ffResourceList.list.data.recordCount - 2
        ? 'top'
        : 'bottom';

    if (fieldInfo.dataFormat === 'text') {
      content = this._reduceText(showData, fieldInfo, resource, clipId);
    } else if (fieldInfo.dataFormat === 'labels') {
      content = this._reduceLabelsForData(showData, direction);
    } else if (fieldInfo.dataFormat === 'time') {
      content = this._reduceTime(showData, direction);
    } else if (fieldInfo.dataFormat === 'ip') {
      content = this._reduceIPCell(showData, clipId, resource);
    } else if (fieldInfo.dataFormat === 'status') {
      content = this._reduceStatus(showData, resource);
    } else if (fieldInfo.dataFormat === 'mapText') {
      content = this._reduceMapText(showData, fieldInfo);
    } else if (fieldInfo.dataFormat === 'replicas') {
      content = this._reduceReplicas(showData, resource);
    } else if (fieldInfo.dataFormat === 'ingressRule') {
      content = this._reduceIngressRule_standalone(showData);
    } else if (fieldInfo.dataFormat === 'backendGroups') {
      content = this._reducebackendGroups(showData);
    } else if (fieldInfo.dataFormat === 'resourceLimit') {
      content = this._reduceResourceLimit(showData);
    } else {
      content = this._reduceText(showData, fieldInfo, resource, clipId);
    }

    return content;
  }

  /** 生成table的表格信息 */
  private _renderTablePanel() {
    let { actions, subRoot } = this.props,
      { resourceOption, resourceInfo, resourceName } = subRoot,
      { ffResourceList, resourceMultipleSelection } = resourceOption;

    let addons = [];

    let displayField = !isEmpty(resourceInfo) && resourceInfo.displayField ? resourceInfo.displayField : {};
    // 根据 displayField当中的key来决定展示什么内容
    let showField = [];
    Object.keys(displayField).forEach(item => {
      let fieldInfo = displayField[item];

      // 操作的按钮现在都换成在tablePanel当中去展示
      if (fieldInfo.dataFormat === 'operator') return;

      if (fieldInfo.dataFormat === 'checker') {
        addons.push(
          selectable({
            value: resourceMultipleSelection.map(item => item.id as string),
            onChange: keys => {
              actions.resource.selectMultipleResource(
                ffResourceList.list.data.records.filter(item => keys.indexOf(item.id as string) !== -1)
              );
            }
          })
        );
        return;
      }
      let columnInfo: TableColumn<Resource> = {
        key: item + uuid(),
        header: fieldInfo.headTitle,
        width: fieldInfo.width,
        render: x => this._renderBodyCell(x, fieldInfo, item + uuid())
      };

      if (fieldInfo.headCell) {
        let style: React.CSSProperties = { display: 'block' };

        let headBubbleText = (
          <span style={style}>
            {fieldInfo.headCell.map((item, index) => (
              <span key={index} className="text" style={style}>
                {item}
              </span>
            ))}
          </span>
        );
        // columnInfo['headCell'] = <HeadBubble title={fieldInfo.headTitle} text={headBubbleText} />;
        columnInfo.header = column => <HeadBubble title={fieldInfo.headTitle} text={headBubbleText} />;
      }

      // return columnInfo;
      showField.push(columnInfo);
    });

    const columns: TableColumn<Resource>[] = showField;

    return (
      <TablePanel
        columns={columns}
        operationsWidth={240}
        getOperations={x =>
          this._renderOperationCell(
            x,
            Object.values(displayField).find(fieldInfo => fieldInfo.dataFormat === 'operator')
          )
        }
        action={actions.resource}
        model={ffResourceList}
        emptyTips={t('您选择的该资源的列表为空，您可以切换到其他命名空间')}
        addons={addons}
        rowDisabled={record => {
          if (resourceName === 'np') {
            return IsResourceShowLoadingIcon(resourceName, record);
          } else {
            return false;
          }
        }}
        isNeedContinuePagination={true}
        onRetry={() => {
          actions.resource.resetPaging();
        }}
      />
    );
  }

  /** 链接的跳转 */
  private _handleClickForNavigate(resource: Resource) {
    let { actions, route } = this.props,
      urlParams = router.resolve(route);

    // 选择当前的具体的resouce
    actions.resource.select(resource);
    // 进行路由的跳转
    router.navigate(
      Object.assign({}, urlParams, { mode: 'detail' }),
      Object.assign({}, route.queries, { resourceIns: resource.metadata.name })
    );
  }
}
