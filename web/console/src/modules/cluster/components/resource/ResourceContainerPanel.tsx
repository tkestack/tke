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

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { cloneDeep, isEmpty } from '../../../common/utils';
import { allActions } from '../../actions';
import { SubRouter } from '../../models';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';
import { ResourceDetail } from './resourceDetail/ResourceDetail';
import { EditResourcePanel } from './resourceEdition/EditResourcePanel';
import { UpdateResourcePanel } from './resourceEdition/UpdateResourcePanel';
import { ResourceListPanel } from './ResourceListPanel';
import { HPAPanel } from '@src/modules/cluster/components/scale/hpa';
import { CronHpaPanel } from '@src/modules/cluster/components/scale/cronhpa';
import { VMDetailPanel, SnapshotTablePanel } from './virtual-machine';

interface ResourceContainerPanelState {
  /** 共享锁 */
  isUnlocked?: boolean;

  /** 菜单栏列表 */
  finalSubRouterList?: SubRouter[];
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceContainerPanel extends React.Component<RootProps, ResourceContainerPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isUnlocked: true,
      finalSubRouterList: []
    };
  }

  componentWillUnmount() {
    const { actions } = this.props;
    // 离开集群内的页面的时候，需要清空subRoot当中的所有信息
    actions.resource.clearSubRoot();
  }

  componentDidMount() {
    const { actions, region, cluster, subRoot, route, clusterInfoList } = this.props;
    const { rid, clusterId } = route.queries;
    // 这里去拉取侧边栏的配置，侧边路由
    actions.subRouter.applyFilter({ clusterId, module: 'cluster' });

    // 这里是确保当前的cluster的相关信息是已经存在了，如果没有，则进行集群的拉取
    const isNeedFetchRegion = region.list.data.recordCount ? false : true;
    isNeedFetchRegion && actions.region.applyFilter({});

    /**
     * 如果需要进行地域的拉取，也需要进行clusterInfo的拉取，并初始化集群版本
     * 如果直接在集群详情页进行刷新的话，进行特定集群详情的拉取，从集群列表跳转过来的时候，会自动初始化集群版本
     */
    isNeedFetchRegion &&
      clusterInfoList.data.recordCount === 0 &&
      actions.cluster.clusterInfo.applyFilter({
        regionId: rid,
        specificName: clusterId
      });

    !isNeedFetchRegion &&
      cluster.list.data.recordCount === 0 &&
      actions.cluster.applyFilter({ regionId: region.selection.value });
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let { route, actions, namespaceSelection, subRoot } = nextProps,
      newUrlParam = router.resolve(route),
      { mode, resourceInfo, subRouterList, addons } = subRoot;
    const newMode = newUrlParam['mode'];
    const newResourceName = newUrlParam['resourceName'];
    const oldMode = this.props.subRoot.mode;

    if (newMode !== '' && oldMode !== newMode && newMode !== mode) {
      actions.resource.selectMode(newMode);
      // 这里是判断回退动作，取消动作等的时候，回到list页面，需要重新拉取一下，激活一下轮训的状态等
      if (newUrlParam['sub'] === 'sub' && !isEmpty(resourceInfo) && newMode === 'list') {
        actions.resource.resetPaging();
        actions.resource.poll();
      }

      // newUrlParam['sub'] === 'sub' && !isEmpty(resourceInfo) && newMode === 'list' && newResourceName !== 'hpa' && actions.resource.poll();
    }

    /** =================== 这里是判断二级菜单路由的配置 ====================== */
    // 判断路由是否拉取完毕
    if (subRouterList.fetched && this.props.subRoot.subRouterList.fetched !== subRouterList.fetched) {
      this.setState({
        finalSubRouterList: subRouterList.data.records
      });
    }

    if (Object.keys(addons).length !== 0 && Object.keys(this.props.subRoot.addons).length === 0) {
      const addonRouterConfig: {
        [props: string]: { routerIndex: number; routerConfig: { name: string; path: string } };
      } = {
        lbcf: {
          routerIndex: 0,
          routerConfig: {
            name: t('负载均衡'),
            path: 'lbcf'
          }
        },
        tapp: {
          routerIndex: 0,
          routerConfig: {
            name: t('TApp'),
            path: 'tapp'
          }
        }
      };
      const newRouterList: SubRouter[] = cloneDeep(this.state.finalSubRouterList);
      newRouterList.forEach((item, index) => {
        if (item.path === 'service') {
          addonRouterConfig['lbcf'].routerIndex = index;
        } else if (item.path === 'resource') {
          addonRouterConfig['tapp'].routerIndex = index;
        }
      });

      const keys = Object.keys(addons);

      keys.forEach(key => {
        if (key === 'LBCF') {
          const lbcfConfig = addonRouterConfig['lbcf'];
          newRouterList[lbcfConfig.routerIndex].sub.push(lbcfConfig.routerConfig);
        } else if (key === 'TappController') {
          const tappConfig = addonRouterConfig['tapp'];
          newRouterList[tappConfig.routerIndex].sub.push(tappConfig.routerConfig);
        }
      });
      this.setState({ finalSubRouterList: newRouterList });
    }
    /** =================== 这里是判断二级菜单路由的配置 ====================== */
  }

  render() {
    const { route } = this.props,
      urlParam = router.resolve(route);
    const { mode, resourceName } = urlParam;
    let content: JSX.Element;

    // 截断hpa和cronhpa的页面逻辑到scale模块
    if (mode !== 'list' && (resourceName === 'hpa' || resourceName === 'cronhpa')) {
      if (resourceName === 'hpa') {
        return <HPAPanel />;
      } else if (resourceName === 'cronhpa') {
        return <CronHpaPanel />;
      }
    } else if (mode === 'detail' && resourceName === 'virtual-machine') {
      return <VMDetailPanel />;
    } else if (mode === 'snapshot' && resourceName === 'virtual-machine') {
      return <SnapshotTablePanel route={route} />;
    } else {
      // 判断应该展示什么组件
      switch (mode) {
        case 'list':
          content = <ResourceListPanel subRouterList={this.state.finalSubRouterList} />;
          break;

        case 'detail':
          content = <ResourceDetail />;
          break;

        case 'create':
        case 'modify':
        case 'apply':
          content = <EditResourcePanel {...this.props} />;
          break;

        case 'update':
          content = <UpdateResourcePanel />;
          break;

        default:
          content = <ResourceListPanel subRouterList={this.state.finalSubRouterList} />;
          break;
      }

      return content;
    }
  }
}
