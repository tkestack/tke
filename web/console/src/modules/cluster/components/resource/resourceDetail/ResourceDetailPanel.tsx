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
import { Bubble, Button, Card, PopConfirm, Table, TableColumn, TabPanel, Tabs, Text } from '@tea/component';
import { stylize } from '@tea/component/table/addons/stylize';
import { bindActionCreators, OperationState, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';
import { resourceConfig } from '../../../../../../config';
import { dateFormatter } from '../../../../../../helpers';
import { HeadBubble, ListItem } from '../../../../common/components';
import { DetailLayout } from '../../../../common/layouts';
import { DetailDisplayFieldProps, DetailInfoProps } from '../../../../common/models';
import { isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { ExternalTrafficPolicy, ResourceStatus, SessionAffinity } from '../../../constants/Config';
import { BackendGroup, BackendRecord, CreateResource, LbcfResource, PortMap, Resource, RuleMap } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

interface ResourceDetailPanelState {
  /** 当前选择的tab */
  tabName?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceDetailPanel extends React.Component<RootProps, ResourceDetailPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      tabName: ''
    };
  }

  render() {
    const { subRoot } = this.props,
      { resourceName, resourceDetailState } = subRoot,
      { resourceDetailInfo } = resourceDetailState;

    const istapp = resourceName === 'tapp';
    const resourceIns = resourceDetailInfo.selection;

    if (istapp && resourceIns) {
      //tapp 需要展示灰度升级的container信息
      //如果有灰度升级项，需要展示在container当中
      const { templates, templatePool } = resourceIns.spec;
      if (templates) {
        const extraContainers = Object.values(templates);
        extraContainers.forEach((item: string) => {
          if (templatePool[item]) {
            //修改提示信息
            const grayUpdateContainers = templatePool[item].spec.containers.map(c => {
              return {
                ...c,
                name: `${c.name} (灰度升级${item})`
              };
            });

            resourceIns.spec.template.spec.containers.push(...grayUpdateContainers);
          }
        });
      }
    }

    return resourceIns === null ? (
      <noscript />
    ) : (
      <div>
        {this._renderBasicInfo(resourceIns)}
        {resourceName === 'ingress' && this._renderRules(resourceIns)}
        {resourceName === 'svc' && this._renderAdvancedInfo(resourceIns)}
        {resourceName === 'lbcf' && this._renderBackGroup(resourceIns)}
        {this._renderVolumes(resourceIns)}
        {this._renderContainer(resourceIns)}
      </div>
    );
  }

  /** 展示基础数据 */
  private _renderBasicInfo(resourceIns: Resource) {
    const { subRoot } = this.props,
      { resourceInfo } = subRoot;

    const basicInfo = resourceInfo.detailField.detailInfo.info ? resourceInfo.detailField.detailInfo.info : {};
    const blockKeys = Object.keys(basicInfo);
    let content: JSX.Element;

    if (blockKeys.length) {
      // 这里需要去遍历 info里面的  metadata 和 status当中的信息

      const showContentObj: any[] = [];
      let showContentArr: any[] = [];

      blockKeys.forEach((blockKey, index) => {
        const detailInfoField = basicInfo[blockKey].dataField[0].split('.');
        // 需要展示的详情的数据
        const detailInfo = this._getFinalData(detailInfoField, resourceIns);
        const displayField = Object.keys(basicInfo[blockKey].displayField).length
          ? basicInfo[blockKey].displayField
          : {};

        // 需要展示的字段名
        const showField = Object.keys(displayField);

        showField.forEach((item, showIndex) => {
          const fieldInfo = displayField[item];
          /**
           * 这里是去判断annotations里面，其为 *.*.*. 可以有很多 .
           * eg: storageclass.beta.kubernetes.io/is-default-class
           * 那么这里的 数组 ['annotations', 'storageclass.beta.kubernetes.io/is-default-class']
           */
          const dataFieldIns =
            fieldInfo.dataField[0] !== ''
              ? fieldInfo.dataField.length > 1
                ? fieldInfo.dataField
                : fieldInfo.dataField[0].split('.')
              : [];
          let showData = this._getFinalData(dataFieldIns, detailInfo);

          // 这里是要去判断noExsit的展示
          showData = showData === '' ? fieldInfo.noExsitedValue : showData;

          // 这里根据displayField当中搞得格式去展示相对应的showContent
          const showElement = this._renderFormItem({ showData, fieldInfo, detailInfo });
          // 这里需要对齐进行排序
          showContentObj.push({
            order: fieldInfo.order,
            item: showElement
          });
        });
      });

      showContentArr = showContentObj.sort((prev, next) => +prev.order - +next.order).map(showObj => showObj.item);

      content = (
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('基本信息')}</h3>
            </div>
            <div className="param-bd">
              <ul className="item-descr-list">{showContentArr}</ul>
            </div>
          </div>
        </DetailLayout>
      );
    } else {
      content = <noscript />;
    }

    return content;
  }

  /** 展示转发规则 rules */
  private _renderRules(resourceIns: Resource) {
    const { subRoot } = this.props,
      { resourceInfo } = subRoot;

    let content: JSX.Element;
    // 展示转发规则
    const showContentArr: any[] = [];

    const isQcloudIngress =
      resourceIns &&
      resourceIns.metadata &&
      resourceIns.metadata.annotations &&
      resourceIns.metadata.annotations['kubernetes.io/ingress.class'] === 'qcloud'
        ? true
        : false;

    // 这里是为了腾讯云专用的Ingress的展示，转发配置是放在annotaions里面的
    if (isQcloudIngress) {
      const annotations = resourceIns.metadata.annotations;
      const showData = {
        http: annotations['kubernetes.io/ingress.http-rules'] || '',
        https: annotations['kubernetes.io/ingress.https-rules'] || ''
      };

      const showContent = this._renderRulesItem(showData, true, resourceIns);
      showContentArr.push(showContent);
    } else {
      const basicInfo = resourceInfo.detailField.detailInfo.rules ? resourceInfo.detailField.detailInfo.rules : {};
      const blockKeys = Object.keys(basicInfo);

      if (blockKeys.length) {
        blockKeys.forEach((blockKey, index) => {
          const detailInfoField = basicInfo[blockKey].dataField[0].split('.');
          //需要展示的详情的数据
          const detailInfo = this._getFinalData(detailInfoField, resourceIns);
          const displayField = Object.keys(basicInfo[blockKey].displayField).length
            ? basicInfo[blockKey].displayField
            : {};

          // 需要展示的字段名
          const showField = Object.keys(displayField);

          showField.forEach((item, showIndex) => {
            const fieldInfo = displayField[item];
            const dataFieldIns = fieldInfo.dataField[0] !== '' ? fieldInfo.dataField[0].split('.') : [];
            let showData = this._getFinalData(dataFieldIns, detailInfo);

            // 这里是要去判断noExist的展示
            showData = showData === '' ? [] : showData;

            let showContent;
            // 这里根据displayField当中搞的格式去展示相对应的showContent
            // showContent = this._renderFormItem({ showData, fieldInfo, detailInfo });
            showContent = this._renderRulesItem(showData);
            showContentArr.push(showContent);
          });
        });
      }
    }

    if (showContentArr.length) {
      content = (
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('转发配置')}</h3>
            </div>
            <div className="param-bd">
              <ul className="item-descr-list">{showContentArr}</ul>
            </div>
          </div>
        </DetailLayout>
      );
    } else {
      content = <noscript />;
    }

    return content;
  }

  /** 展示数据卷，数据卷这里特殊处理 */
  private _renderVolumes(resourceIns: Resource) {
    const { subRoot } = this.props,
      { resourceInfo } = subRoot;

    const detailInfo = resourceInfo.detailField.detailInfo.volume ? resourceInfo.detailField.detailInfo.volume : {};
    const blockKeys = Object.keys(detailInfo);
    let content: JSX.Element;

    if (blockKeys.length) {
      // 目前此区域只展示 volumes，故只需要获取volumes
      const detailInfoField = detailInfo['volumes'].dataField[0].split('.');
      const volumns = this._getFinalData(detailInfoField, resourceIns);

      // 展示数据卷资源的信息
      const renderVolumeInfo = (volumn, type: string) => {
        let content = '';
        if (type === 'configMap') {
          content = volumn[type]['name'] + t('（资源名称）');
        } else if (type === 'secret') {
          content = volumn[type]['secretName'] + t('（资源名称）');
        } else if (type === 'hostPath') {
          const validateType = volumn[type]['type'] ? volumn[type]['type'] : t('不校验');
          content = volumn[type]['path'] + t('（主机路径）') + validateType + t('（路径检查类型）');
        } else if (type === 'nfs') {
          content = volumn[type]['server'] + ':' + volumn[type]['path'] + t('（NFS路径）');
        } else if (type === 'persistentVolumeClaim') {
          content = volumn[type]['claimName'] + t('（资源名称）');
        } else {
          content = '';
        }
        return content;
      };

      content = (
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('数据卷（Volumes）')}</h3>
            </div>
            <div className="param-bd">
              {volumns.length === 0 ? (
                <span className="text-label">{t('暂无数据卷')}</span>
              ) : (
                <ul className="item-descr-list">
                  {volumns.map((volumn, index) => {
                    /**
                     * 返回的volumes格式固定
                     * {
                     *      name: ***
                     *      [volumesType]: any
                     * }
                     * volumesType即数据卷的类型
                     */
                    const keys = Object.keys(volumn);
                    const label = keys.filter(item => item !== 'name')[0]; // 这里是因为volumes 当中的数据卷的类型是不固定的

                    return (
                      <ListItem key={index} label={label}>
                        <span className="text">{volumn['name'] + t('（卷名称） ')}</span>
                        <span className="text">{renderVolumeInfo(volumn, label)}</span>
                      </ListItem>
                    );
                  })}
                </ul>
              )}
            </div>
          </div>
        </DetailLayout>
      );
    } else {
      content = <noscript />;
    }

    return content;
  }

  /**展示service高级设置 */
  private _renderAdvancedInfo(resourceIns: Resource) {
    const { subRoot } = this.props,
      { resourceInfo } = subRoot;

    const detailInfo = resourceInfo.detailField.detailInfo.advancedInfo
      ? resourceInfo.detailField.detailInfo.advancedInfo
      : {};
    const blockKeys = Object.keys(detailInfo);
    let content: JSX.Element;

    if (blockKeys.length) {
      const advancedInfoField = detailInfo['spec'].dataField[0].split('.');
      const advancedInfo = this._getFinalData(advancedInfoField, resourceIns);

      content = (
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('高级设置')}</h3>
            </div>
            <div className="param-bd">
              <ul className="item-descr-list">
                <ListItem label="ExternalTrafficPolicy">
                  {advancedInfo.externalTrafficPolicy
                    ? advancedInfo.externalTrafficPolicy
                    : ExternalTrafficPolicy.Cluster}
                </ListItem>
                <ListItem label="Session Affinity">
                  {advancedInfo.sessionAffinity ? advancedInfo.sessionAffinity : SessionAffinity.None}
                </ListItem>
                <ListItem
                  label={t('最大会话保持时间')}
                  isShow={advancedInfo.sessionAffinity === SessionAffinity.ClientIP}
                >
                  {advancedInfo.sessionAffinityConfig &&
                  advancedInfo.sessionAffinityConfig.clientIP &&
                  advancedInfo.sessionAffinityConfig.clientIP.timeoutSeconds
                    ? advancedInfo.sessionAffinityConfig.clientIP.timeoutSeconds + t('秒')
                    : '-'}
                </ListItem>
              </ul>
            </div>
          </div>
        </DetailLayout>
      );
    } else {
      content = <noscript />;
    }

    return content;
  }
  /** 展示容器 */
  private _renderBackGroup(resourceIns: LbcfResource) {
    const { subRoot } = this.props,
      { resourceInfo } = subRoot;

    const detailInfo = resourceInfo.detailField.detailInfo.backGroup
      ? resourceInfo.detailField.detailInfo.backGroup
      : {};
    const blockKeys = Object.keys(detailInfo);
    let content: JSX.Element;

    if (blockKeys.length) {
      // 目前此区域只展示backGroups，故只需要获取backGroups
      const detailInfoField = detailInfo['backGroups'].dataField[0].split('.');
      const backGroups: BackendGroup[] = this._getFinalData(detailInfoField, resourceIns);

      const tabs = backGroups
        ? backGroups.map((item, index) => {
            const tab = {
              id: item['name'] + index,
              label: item['name']
            };
            return tab;
          })
        : [];

      let selected = tabs[0];
      if (this.state.tabName) {
        const finder = tabs.find(x => x.id === this.state.tabName);
        if (finder) {
          selected = finder;
        }
      }

      content = (
        <Card>
          <Card.Header>
            <h3>{t('后端负载(BackendGroup)')}</h3>
          </Card.Header>
          <Card.Body>
            <Tabs
              tabs={tabs}
              activeId={selected ? selected.id : ''}
              onActive={tab => {
                this.setState({ tabName: tab.id });
              }}
              // className="tc-15-tab tc-15-tab-alt noMarginTop"
            >
              {this._renderContainerBody(backGroups, detailInfo['backGroups'])}
            </Tabs>
          </Card.Body>
        </Card>
      );
    } else {
      content = <noscript />;
    }

    return content;
  }

  /** 展示容器 */
  private _renderContainer(resourceIns: Resource) {
    const { subRoot } = this.props,
      { resourceInfo } = subRoot;

    const detailInfo = resourceInfo.detailField.detailInfo.container
      ? resourceInfo.detailField.detailInfo.container
      : {};
    const blockKeys = Object.keys(detailInfo);
    let content: JSX.Element;

    if (blockKeys.length) {
      // 目前此区域只展示 volumes，故只需要获取volumes
      const detailInfoField = detailInfo['containers'].dataField[0].split('.');
      const containers = this._getFinalData(detailInfoField, resourceIns);

      const tabs = containers
        ? containers.map((item, index) => {
            const tab = {
              id: item['name'] + index,
              label: item['name']
            };
            return tab;
          })
        : [];

      let selected = tabs[0] ? tabs[0] : {};
      if (this.state.tabName) {
        const finder = tabs.find(x => x.id === this.state.tabName);
        if (finder) {
          selected = finder;
        }
      }

      content = (
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('容器（Containers）')}</h3>
            </div>
            <div className="param-bd docker-param">
              <Tabs
                tabs={tabs}
                activeId={selected.id ? selected.id : ''}
                onActive={tab => {
                  this.setState({ tabName: tab.id });
                }}
                className="tc-15-tab tc-15-tab-alt noMarginTop"
              >
                {this._renderContainerBody(containers, detailInfo['containers'])}
              </Tabs>
            </div>
          </div>
        </DetailLayout>
      );
    } else {
      content = <noscript />;
    }

    return content;
  }

  /** 生成container内部的展示 */
  private _renderContainerBody(containers: any[], detailInfo: DetailInfoProps) {
    // 配置文件当中的detailInfo 下的 container
    const displayField = Object.keys(detailInfo).length ? detailInfo.displayField : {};
    // 需要展示的字段名
    const showField = Object.keys(displayField);

    return containers
      ? containers.map((container, index) => (
          /* eslint-disable */
          <TabPanel key={index} id={container['name'] + index}>
            <ul className="item-descr-list">
              {showField.map((item, showIndex) => {
                let fieldInfo = displayField[item];
                // 这里主要是因为 一些key自己就包含 . 所以使用不同的写法
                let dataFieldIns =
                  fieldInfo.dataField.length > 1 ? fieldInfo.dataField : fieldInfo.dataField[0].split('.');
                let showData = this._getFinalData(dataFieldIns, container);

                if (showData !== '') {
                  // 这里根据每种格式去生成相对应的showContent
                  let showElement = this._renderFormItem({ showData, fieldInfo });
                  return showElement;
                } else {
                  return <noscript key={showIndex} />;
                }
              })}
            </ul>
          </TabPanel>
        ))
      : [];
  }

  /** 获取最终的数据 */
  private _getFinalData(detailInfoField: string[], dataResource) {
    let result = dataResource;

    for (let index = 0; index < detailInfoField.length; index++) {
      // 如果result不为一个object，则遍历结束
      if (typeof result !== 'object') {
        break;
      }
      result = result[detailInfoField[index]];
    }

    return result || '';
  }

  /** 生成相对应的格式内容 */
  private _renderFormItem(options: { showData: any; fieldInfo: DetailDisplayFieldProps; detailInfo?: any }) {
    let { showData, fieldInfo, detailInfo } = options;
    let { resourceName } = this.props.subRoot;
    let { actions } = this.props;
    let showContent,
      isShowListItem = true;

    // 当showData为 - 的时候，即 此时找不到值，则原样返回即可
    if (fieldInfo.dataFormat === 'text' || showData === '-') {
      // 最为普通的直接返回showData展示即可
      showContent = <p>{showData}</p>;
    } else if (fieldInfo.dataFormat === 'status') {
      // 展示资源的状态
      let statusMap = ResourceStatus[resourceName];
      showContent = statusMap ? (
        <p className={classnames('', statusMap[showData] && statusMap[showData].classname)}>
          {(statusMap[showData] && statusMap[showData].text) || '-'}
        </p>
      ) : (
        <p>showData</p>
      );
    } else if (fieldInfo.dataFormat === 'time') {
      // 时间需要进行处理为 YYYY-MM-DD HH:mm:ss 的格式
      showContent = <p>{dateFormatter(new Date(showData), 'YYYY-MM-DD HH:mm:ss')}</p>;
    } else if (fieldInfo.dataFormat === 'labels') {
      // label的展示为 key:value、key:value
      let showLabels = '';
      let keys = Object.keys(showData);
      keys.forEach((item, index) => {
        showLabels += item + '：' + showData[item];
        if (index !== keys.length - 1) {
          showLabels += '、';
        }
      });
      showContent = <p>{showLabels}</p>;
    } else if (fieldInfo.dataFormat === 'keyvalue') {
      // keyvalue的展示为一行一个
      // key1:value1
      // key2:value2

      showContent = Object.keys(showData).map((item, index) => {
        return (
          <p key={index}>
            {item}:{showData[item]}
          </p>
        );
      });
    } else if (fieldInfo.dataFormat === 'ip') {
      showContent = (
        <p>
          <span id={fieldInfo.label}>{showData}</span>
          {showData === 'None' && fieldInfo.extraInfo && <span>{` (${fieldInfo.extraInfo})`}</span>}
        </p>
      );
    } else if (fieldInfo.dataFormat === 'rules') {
      showContent = this._renderRulesItem(showData);
    } else if (fieldInfo.dataFormat === 'ports') {
      showContent = this._renderPortsItem(showData, detailInfo);
    } else if (fieldInfo.dataFormat === 'pods') {
      showContent = this._renderPodsItem(showData, fieldInfo, detailInfo);
    } else if (fieldInfo.dataFormat === 'replicas') {
      // 展示副本的相关
      showContent = this._renderReplicasItem(showData, fieldInfo, detailInfo);
    } else if (fieldInfo.dataFormat === 'array') {
      // command 是一个数组，取第一项进行展示
      showContent = showData.map((data, index) => {
        return <p key={index}>{data}</p>;
      });
    } else if (fieldInfo.dataFormat === 'env') {
      // env 环境变量的展示 是 a = b
      showContent = showData.map((data, index) => {
        if (data['value']) {
          return <p key={index}>{`${data['name']}=${data['value']}`}</p>;
        } else if (data['valueFrom']) {
          let showKey = Object.keys(data['valueFrom'])[0],
            refData = data['valueFrom'][showKey];
          return (
            <div key={index}>
              <span className="text">{`${data['name']}=${showKey}`}</span>
              <Bubble
                placement="left"
                content={
                  <React.Fragment>
                    {Object.entries(refData ?? {}).map(([key, value]) => (
                      <>
                        <p>{`${t('名称：')}${key}`}</p>
                        <p>{`Key：${value}`}</p>
                      </>
                    ))}
                  </React.Fragment>
                }
              >
                <i className="plaint-icon" style={{ marginLeft: '5px' }} />
              </Bubble>
            </div>
          );
        } else {
          return <p key={index}>{`${data['name']}`}</p>;
        }
      });
    } else if (fieldInfo.dataFormat === 'volume') {
      showContent = showData.map((data, index) => {
        return (
          <p key={index}>
            <span className="text text-label">{t('数据卷名称: ')}</span>
            <span className="text">{`${data['name']} `}</span>
            <span className="text text-label">{t('目标路径: ')}</span>
            <span className="text">{`${data['mountPath']} `}</span>
            <span className="text text-label">{t('挂载子路径: ')}</span>
            <span className="text">{`${data['subPath'] ? data['subPath'] : t('未设置，默认全覆盖目标路径')} `}</span>
          </p>
        );
      });
    } else if (fieldInfo.dataFormat === 'probe') {
      showContent = this._renderProbeItem(showData);
    } else if (fieldInfo.dataFormat === 'mapText') {
      let { mapTextConfig } = fieldInfo;
      showContent = (
        <Text parent="div" overflow>
          {mapTextConfig[showData]}
        </Text>
      );
    } else if (fieldInfo.dataFormat === 'gameBGPort') {
      showContent = this._reduceGameBackendGroupPort(showData);
    } else if (fieldInfo.dataFormat === 'operator') {
      let {
        clusterVersion,
        namespaceSelection,
        route,
        subRoot: {
          deleteResourceFlow,
          detailResourceOption: { detailDeleteResourceSelection }
        }
      } = this.props;
      //detail页面删除backendGroup
      showContent = (
        <PopConfirm
          title="确定要删除后端负载配置？"
          message="删除后，后端负载配置将不再生效"
          visible={deleteResourceFlow.operationState !== OperationState.Pending}
          footer={
            <>
              <Button
                type="link"
                onClick={() => {
                  let bgResourceInfo = resourceConfig(clusterVersion).lbcf_bg;
                  let resourceIns = detailDeleteResourceSelection;

                  let resource: CreateResource = {
                    id: uuid(),
                    resourceInfo: bgResourceInfo,
                    namespace: namespaceSelection,
                    clusterId: route.queries['clusterId'],
                    resourceIns
                  };
                  actions.workflow.deleteResource.start([resource]);
                  actions.workflow.deleteResource.perform();
                }}
              >
                删除
              </Button>
              <Button
                type="text"
                onClick={() => {
                  if (deleteResourceFlow.operationState === OperationState.Done) {
                    actions.workflow.deleteResource.reset();
                  }
                  if (deleteResourceFlow.operationState === OperationState.Started) {
                    actions.workflow.deleteResource.cancel();
                  }
                }}
              >
                取消
              </Button>
            </>
          }
          placement="top-start"
        >
          <Button
            type={'link'}
            onClick={() => {
              actions.resource.selectDetailDeleteResouceIns(showData);
              actions.workflow.deleteResource.start([]);
            }}
          >
            {t('删除')}
          </Button>
        </PopConfirm>
      );
    } else if (fieldInfo.dataFormat === 'backendRecords') {
      showContent = this._renderbackendRecordsItem(showData);
    } else if (fieldInfo.dataFormat === 'forceDeletePod') {
      showContent = <Text>{showData ? '迁移' : '不迁移'}</Text>;
    }

    return (
      <ListItem key={uuid()} isShow={isShowListItem} label={fieldInfo.label}>
        {showContent}
      </ListItem>
    );
  }

  /** 展示存活检查、就绪检查 */
  private _renderProbeItem(showData: any) {
    let showContent = [];

    let checkWayContent = [];
    // 存活检查、就绪检查里面的 检查方式需要判断一下，其他指标都是number类型
    let keys = showData ? Object.keys(showData) : [];
    keys.forEach(key => {
      let typeofData = Object.prototype.toString.call(showData[key]);
      if (typeofData !== '[object Object]') {
        showContent.push(<p>{`${key}：${showData[key]}`}</p>);
      } else {
        // 检查方法的实际配置
        let checkInfo = showData[key];
        let checkInfoKeys = Object.keys(checkInfo);
        // 检查方法下还有一些内容，如port、path等
        checkInfoKeys.forEach(checkItemKey => {
          let content;
          if (typeof checkInfo[checkItemKey] === 'object') {
            content = <p>{`${checkItemKey}：${checkInfo[checkItemKey][0]}`}</p>;
          } else {
            content = <p>{`${checkItemKey}：${checkInfo[checkItemKey]}`}</p>;
          }
          checkWayContent.push(content);
        });

        showContent.push(<p>{`inspection method：${key}`}</p>);
      }
    });
    // 将 请求方式内部的东西都平铺到要展示的内容上
    let finalContent = [...checkWayContent, ...showContent];
    let finalShowContent = (
      <Bubble placement="left" content={finalContent || null}>
        <a href="javascript:;" style={{ textDecoration: 'none', fontSize: '14px' }}>
          <span style={{ verticalAlign: 'middle' }}>Detail</span>
          <i style={{ verticalAlign: 'middle' }} className="plaint-icon" />
        </a>
      </Bubble>
    );

    return finalShowContent;
  }

  /** 展示 replicas */
  private _renderReplicasItem(showData: any, fieldInfo: DetailDisplayFieldProps, detailInfo?: any) {
    let showContent;
    if (!isEmpty(fieldInfo.subDisplayField)) {
      let displayField = fieldInfo.subDisplayField;
      let displayKeys = Object.keys(displayField);

      let replicasShowObj = {};
      displayKeys.forEach(item => {
        let dataField = displayField[item].dataField[0].split('.');
        let data = this._getFinalData(dataField, detailInfo);
        replicasShowObj[displayField[item]['label']] = data === '' ? displayField[item].noExsitedValue : data;
      });

      // 这里是需要去展示status当中的replicas的各种状态
      let subDisplayContent = '';
      let replicasObjKeys = Object.keys(replicasShowObj);
      replicasObjKeys.forEach((sub, index) => {
        subDisplayContent += sub + ':' + replicasShowObj[sub];
        if (index !== replicasObjKeys.length - 1) {
          subDisplayContent += '、';
        }
      });
      showContent = <p>{t('期望Pod数量:') + `${showData}（${subDisplayContent}）`}</p>;
    } else {
      showContent = <p>{showData}</p>;
    }

    return showContent;
  }

  /** 展示pod的状态 */
  private _renderPodsItem(showData: any, fieldInfo: DetailDisplayFieldProps, detailInfo?: any) {
    let showContent;
    // pods的展示形式为 Pod status（desired: 0, ready: 0, succeed: 0, failed: 0）
    if (!isEmpty(fieldInfo.subDisplayField)) {
      let displayField = fieldInfo.subDisplayField;
      let displayKeys = Object.keys(displayField);

      let podsStatus = {};
      displayKeys.forEach(item => {
        let dataField = displayField[item].dataField[0].split('.');
        let data = this._getFinalData(dataField, detailInfo);
        podsStatus[displayField[item]['label']] = data === '' ? displayField[item].noExsitedValue : data;
      });

      // 这里是需要去展示在pod status当中的各种状态
      let subDisplayContent = '';
      let podsStatusKeys = Object.keys(podsStatus);
      podsStatusKeys.forEach((sub, index) => {
        subDisplayContent += sub + ':' + podsStatus[sub];
        if (index !== podsStatusKeys.length - 1) {
          subDisplayContent += '，';
        }
      });
      showContent = <p>{subDisplayContent}</p>;
    } else {
      showContent = <p>unknow</p>;
    }

    return showContent;
  }

  /** 展示端口映射相关 */
  private _renderPortsItem(showData: any, detailInfo?: any) {
    let showContent;

    let columns: TableColumn<PortMap>[] = [
      {
        key: 'protocol',
        header: column => {
          return (
            <HeadBubble align="start" title={t('协议')} text={t('使用公网/内网负载均衡时，TCP和UDP协议不能混合使用')} />
          );
        },
        render: x => (
          <Text parent="div" overflow>
            {x.protocol}
          </Text>
        )
      },
      {
        key: 'port',
        header: t('容器端口'),
        render: x => (
          <Text parent="div" overflow>
            {x.targetPort}
          </Text>
        )
      },
      {
        key: 'nodePort',
        header: t('主机端口'),
        render: x => (
          <Text parent="div" overflow>
            {x.nodePort}
          </Text>
        )
      },
      {
        key: 'targetPort',
        header: t('服务端口'),
        render: x => (
          <Text parent="div" overflow>
            {x.port}
          </Text>
        )
      }
    ];

    // 这里需要判断，因为clusterIP的类型，没有NodePort
    if (detailInfo.type !== 'NodePort') {
      let nodeIndex = columns.findIndex(c => c.key === 'nodePort');
      columns.splice(nodeIndex, 1);
    }

    let records: PortMap[] = showData.map(item => ({
      id: uuid(),
      protocol: item.protocol || '-',
      nodePort: item.nodePort || '-',
      port: item.port || '-',
      targetPort: item.targetPort || '-'
    }));

    return (
      <Table
        columns={columns}
        records={records}
        addons={[
          stylize({
            style: { overflow: 'visible', maxWidth: '800px' }
          })
        ]}
      />
    );
  }

  /** 展示转发规则 */
  private _renderRulesItem(showData: any, isQcloudIngress: boolean = false, resource?: Resource) {
    let columns: TableColumn<RuleMap>[] = [
      {
        key: 'port',
        header: t('监听端口'),
        render: x => (
          <Text parent="div" overflow>
            {x.protocol === 'http' ? '80' : '443'}
          </Text>
        )
      },
      {
        key: 'host',
        header: column => {
          return (
            <HeadBubble
              title={t('域名')}
              text={t(
                '非通配的域名支持的字符集 a-z 0-9 . -; 通配的域名，目前只支持 *.example.com的形式，且单个域名中只支持 * 出现一次'
              )}
            />
          );
        },
        render: x => (
          <div>
            {isQcloudIngress ? (
              <Text parent="div" overflow>
                {x.host
                  ? x.host
                  : resource.status.loadBalancer.ingress
                  ? resource.status.loadBalancer.ingress[0].ip
                  : '-'}
              </Text>
            ) : (
              <Text parent="div" overflow>
                {x.host}
              </Text>
            )}
          </div>
        )
      },
      {
        key: 'path',
        header: t('URL路径'),
        render: x => (
          <Text parent="div" overflow>
            {x.path}
          </Text>
        )
      },
      {
        key: 'serviceName',
        header: column => <HeadBubble title={t('后端服务')} text={t('不支持配置访问方式为不启用的服务')} />,
        render: x => (
          <Text parent="div" overflow>
            {x.serviceName}
          </Text>
        )
      },
      {
        key: 'servicePort',
        header: t('服务端口'),
        render: x => (
          <Text parent="div" overflow>
            {x.servicePort}
          </Text>
        )
      }
    ];

    let records: RuleMap[] = [];

    if (isQcloudIngress) {
      columns.unshift({
        key: 'protocol',
        header: t('协议'),
        render: x => (
          <Text parent="div" overflow>
            {x.protocol}
          </Text>
        )
      });

      let httpRules = showData['http'] === '' || showData['http'] === 'null' ? [] : JSON.parse(showData['http']),
        httpsRules = showData['https'] === '' || showData['https'] === 'null' ? [] : JSON.parse(showData['https']);

      httpRules.forEach(item => {
        let tmp: RuleMap = {
          id: uuid(),
          protocol: 'http',
          path: item['path'],
          serviceName: item['backend']['serviceName'],
          servicePort: item['backend']['servicePort'],
          host: item['host']
        };
        records.push(tmp);
      });

      httpsRules.forEach(item => {
        let tmp: RuleMap = {
          id: uuid(),
          protocol: 'https',
          path: item['path'],
          serviceName: item['backend']['serviceName'],
          servicePort: item['backend']['servicePort'],
          host: item['host']
        };
        records.push(tmp);
      });
    } else {
      records = showData.map(rule => {
        let pathInfo = rule['http']['paths'][0],
          path = pathInfo['path'] || '-',
          serviceName = pathInfo['backend']['serviceName'],
          servicePort = pathInfo['backend']['servicePort'];

        let tmp: RuleMap = {
          id: uuid(),
          path,
          serviceName,
          servicePort,
          host: rule['host']
        };

        return tmp;
      });
    }

    return (
      <Table
        columns={columns}
        records={records}
        addons={[
          stylize({
            style: { overflow: 'visible' }
          })
        ]}
      />
    );
  }
  private _reduceGameBackendGroupPort(showData: any) {
    return (
      <Text parent="div" overflow>
        {showData['protocol']}:{showData['portNumber']}
      </Text>
    );
  }
  /** 展示backendRecord */
  private _renderbackendRecordsItem(showData: any) {
    let columns: TableColumn<BackendRecord>[] = [
      {
        key: 'name',
        header: t('名称'),
        width: '20%',
        render: x => (
          <Text parent="div" overflow>
            {x.name}
          </Text>
        )
      },
      {
        key: 'backendAddr',
        header: t('backendAddr'),
        width: '40%',
        render: x => (
          <Bubble content={x.backendAddr} placement={'top'}>
            <Text parent="div" overflow>
              {x.backendAddr}
            </Text>
          </Bubble>
        )
      },
      {
        key: 'condition',
        width: '20%',
        header: t('condition'),
        render: x => {
          let registerItem = x.conditions.filter(condition => condition.type === 'Registered');
          return (
            <Text parent="div" overflow>
              {registerItem.length ? `Registered:${registerItem[0].status}` : 'Registered:false'}
            </Text>
          );
        }
      },
      {
        key: 'operator',
        width: '20%',
        header: t('操作'),
        render: x => {
          return (
            <Button
              type={'link'}
              onClick={() => {
                let { route, actions } = this.props;
                let urlParams = router.resolve(route);
                actions.resource.initDetailResourceName('lbcf_br', x.name);
                router.navigate(Object.assign(urlParams, { tab: 'event' }), route.queries);
              }}
            >
              查看事件
            </Button>
          );
        }
      }
    ];

    return (
      <Table
        columns={columns}
        records={showData}
        addons={[
          stylize({
            style: { overflow: 'visible' }
          })
        ]}
      />
    );
  }
}
