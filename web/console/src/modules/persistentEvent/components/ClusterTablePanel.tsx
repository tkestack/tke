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
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Icon } from '@tea/component';
import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { Resource } from '../../common';
import { Clip, LinkButton } from '../../common/components';
import { includes } from '../../common/utils';
import { clsRegionMap } from '../../logStash/constants/Config';
import { allActions } from '../actions';
import { isNeedPollPE, peStatus } from '../constants/Config';
import { router } from '../router';
import { RootProps } from './PersistentEventApp';

const routerSea = seajs.require('router');

/** 加载中的样式 */
const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ClusterTablePanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    // 清空轮询条件
    actions.pe.clearPollEvent();
  }

  render() {
    return this._renderTablePanel();
  }

  /** 渲染列表 */
  private _renderTablePanel() {
    let { actions, cluster, route } = this.props;

    const columns: TablePanelColumnProps<Resource>[] = [
      {
        key: 'name',
        header: t('ID/名称'),
        width: '10%',
        render: x => (
          <div>
            <span id={`clusterId${x.id}`}>
              <a
                className="text-overflow"
                title={x.metadata.name}
                href="javascript:;"
                onClick={() => {
                  this._handleClickForCluster(x);
                }}
              >
                {x.metadata.name}
              </a>
            </span>
            <Clip target={`#clusterId${x.id}`} className="hover-icon" />
            <div className="sl-editor-name">
              <span className="text-overflow m-width" title={x.spec.displayName}>
                {x.spec.displayName}
              </span>
            </div>
          </div>
        )
      },
      {
        key: 'peStatus',
        header: t('状态'),
        width: '10%',
        render: x => this._getPEStatus(x)
      },
      {
        key: 'store',
        header: t('存储端'),
        width: '10%',
        render: x => this._getStoreType(x)
      },
      {
        key: 'storeInfo',
        header: t('存储对象'),
        width: '20%',
        render: x => this._getStoreInfo(x)
      },
      {
        key: 'operation',
        header: t('操作'),
        width: '10%',
        render: x => this._renderOprationCell(x)
      }
    ];

    let emptyTips: JSX.Element = (
      <div className="text-center">
        <Trans>
          集群列表为空，您可以
          <a href="javascript:;" onClick={() => routerSea.navigate(`tke/cluster`)}>
            [新建一个集群]
          </a>
        </Trans>
      </div>
    );

    return <TablePanel columns={columns} emptyTips={emptyTips} model={cluster} action={actions.cluster} />;
  }

  /** 获取当前集群的PE的开启状态 */
  private _getPEStatus(cluster: Resource) {
    let { peList } = this.props;

    let clusterName = cluster.metadata ? cluster.metadata.name : '';

    let peInfo = peList.data.records.find(item => item.spec.clusterName === clusterName);

    let isNeedLoadingIcon = false;
    if (peInfo && peInfo.status.phase) {
      isNeedLoadingIcon = includes(isNeedPollPE, (peInfo.status.phase as string).toLowerCase()) ? true : false;
    }

    return (
      <div className={peInfo && peStatus[peInfo.status.phase] ? peStatus[peInfo.status.phase].classname : ''}>
        <span style={{ verticalAlign: 'middle' }}>
          {peInfo ? peStatus[(peInfo.status.phase as string).toLowerCase()].text : t('未开启')}
        </span>
        {peInfo && peInfo.status.phase === 'running' && <i className="n-success-icon" style={{ marginLeft: '5px' }} />}
        {peInfo && peInfo.status.reason && (
          <Bubble placement="right" content={peInfo.status.reason || null}>
            <div className="tc-15-bubble-icon">
              <i className="tc-icon n-error-icon" />
            </div>
          </Bubble>
        )}
        {isNeedLoadingIcon && <Icon type="loading" className="tea-ml-1n" />}
      </div>
    );
  }

  /** 获取集群开启的持久化存储端的类型 */
  private _getStoreType(cluster: Resource) {
    let { peList } = this.props;

    let content = '';
    let clusterName = cluster.metadata ? cluster.metadata.name : '';
    let peInfo = peList.data.records.find(item => item.spec.clusterName === clusterName);
    if (peInfo) {
      let storeType = Object.keys(peInfo.spec.persistentBackEnd)[0];
      content = storeType === 'cls' ? t('日志服务CLS') : 'Elasticsearch';
    } else {
      content = '-';
    }
    return <div>{content}</div>;
  }

  /** 获取集群开启的持久化存储对象的信息 */
  private _getStoreInfo(cluster: Resource) {
    let { peList, route } = this.props;

    let content: JSX.Element;
    let clusterName = cluster.metadata && cluster.metadata.name;
    let peInfo = peList.data.records.find(item => item.spec.clusterName === clusterName);
    if (peInfo) {
      let storeType = Object.keys(peInfo.spec.persistentBackEnd)[0];
      let backEndInfo = peInfo.spec.persistentBackEnd[storeType];
      if (storeType === 'cls') {
        content = (
          <div>
            <div className="text-overflow m-width" style={{ maxWidth: '90%' }}>
              <span style={{ verticalAlign: 'middle' }}>{t('日志集')}</span>
              <a
                href="javascript:;"
                style={{ textDecoration: 'none' }}
                onClick={() => {
                  // 这里是cls那边的实现，把内容写在localStorage当中
                  localStorage.setItem('cls_logset', backEndInfo.logSetId);
                  routerSea.navigate(`/cls/logset/desc?region=${clsRegionMap[route.queries['rid']]}`);
                }}
              >{`( ${backEndInfo.logSetId} )`}</a>
            </div>
            <div className="sl-editor-name text-overflow m-width">
              <span style={{ verticalAlign: 'middle' }}>{t('日志主题')}</span>
              <a
                style={{ textDecoration: 'none' }}
                href="javascript:;"
                onClick={() => {
                  // 这里是cls的实现，把内容写在localStorage当中
                  localStorage.setItem('cls_logset', backEndInfo.logSetId);
                  localStorage.setItem('cls_topic', backEndInfo.topicId);
                  routerSea.navigate(`/cls/logset/config?region=${clsRegionMap[route.queries['rid']]}`);
                }}
              >{`( ${backEndInfo.topicId} )`}</a>
            </div>
          </div>
        );
      } else {
        content = (
          <div>
            <div className="text-overflow m-width">
              <span>{t('ES地址')}</span>
              <span>{`( ${backEndInfo.scheme || 'http'}://${backEndInfo.ip}:${backEndInfo.port} )`}</span>
            </div>
            <div className="sl-editor-name text-overflow m-width">
              <span>{t('索引')}</span>
              <span>{`( ${backEndInfo.indexName || 'fluentd'} )`}</span>
            </div>
          </div>
        );
      }
    } else {
      content = <div>-</div>;
    }
    return content;
  }

  /** 集群id的处理 */
  private _handleClickForCluster(cluster: Resource) {
    let { route } = this.props;
    let clusterId = cluster.metadata.name;
    let navigatePath = `/tkestack/cluster/sub/list/basic/info?rid=${route.queries['rid']}&clusterId=${clusterId}`;

    // 进行路由的跳转
    routerSea.navigate(navigatePath);
  }

  /** 渲染操作按钮 */
  private _renderOprationCell(cluster: Resource) {
    let { actions, route, peList } = this.props;

    /** 设置按钮 */
    const renderEditButton = () => {
      let disabled = false,
        errorTip = '';

      /**
       * 用于判断是否已经开通过persistentEvent了，cluster.status.AddOns里面是否有 persistentEvent
       * 条件2: peList当中能够找到当前的集群
       */
      let isHasCreatePE =
        (cluster.status.addOns && cluster.status.addOns.persistentEvent) ||
        peList.data.records.find(item => item.spec.clusterName === cluster.metadata.name)
          ? true
          : false;

      return (
        <LinkButton
          tipDirection={'right'}
          errorTip={errorTip}
          disabled={disabled}
          onClick={() => {
            if (!disabled) {
              // 选择当前的cluster
              actions.cluster.selectCluster(cluster);

              let urlParams = router.resolve(route);
              router.navigate(
                Object.assign({}, urlParams, { mode: isHasCreatePE ? 'update' : 'create' }),
                Object.assign({}, route.queries, { clusterId: cluster.metadata.name })
              );
            }
          }}
        >
          {isHasCreatePE ? t('更新设置') : t('设置')}
        </LinkButton>
      );
    };
    return <div>{renderEditButton()}</div>;
  }
}
