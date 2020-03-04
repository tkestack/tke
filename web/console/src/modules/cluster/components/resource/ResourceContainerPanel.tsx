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
    let { actions } = this.props;
    // 离开集群内的页面的时候，需要清空subRoot当中的所有信息
    actions.resource.clearSubRoot();
  }

  componentDidMount() {
    let { actions, region, cluster, subRoot, route, clusterInfoList } = this.props;
    let { rid, clusterId } = route.queries;
    // 这里去拉取侧边栏的配置，侧边路由
    actions.subRouter.applyFilter({});

    // 这里是确保当前的cluster的相关信息是已经存在了，如果没有，则进行集群的拉取
    let isNeedFetchRegion = region.list.data.recordCount ? false : true;
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
    let newMode = newUrlParam['mode'];
    let oldMode = this.props.subRoot.mode;

    if (newMode !== '' && oldMode !== newMode && newMode !== mode) {
      actions.resource.selectMode(newMode);
      // 这里是判断回退动作，取消动作等的时候，回到list页面，需要重新拉取一下，激活一下轮训的状态等
      newUrlParam['sub'] === 'sub' &&
        !isEmpty(resourceInfo) &&
        newMode === 'list' &&
        actions.resource.poll({
          namespace: namespaceSelection,
          clusterId: route.queries['clusterId'],
          regionId: +route.queries['rid']
        });
    }

    /** =================== 这里是判断二级菜单路由的配置 ====================== */
    // 判断路由是否拉取完毕
    if (subRouterList.fetched && this.props.subRoot.subRouterList.fetched !== subRouterList.fetched) {
      this.setState({
        finalSubRouterList: subRouterList.data.records
      });
    }

    if (Object.keys(addons).length !== 0 && Object.keys(this.props.subRoot.addons).length === 0) {
      let addonRouterConfig: {
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
      let newRouterList: SubRouter[] = cloneDeep(this.state.finalSubRouterList);
      newRouterList.forEach((item, index) => {
        if (item.path === 'service') {
          addonRouterConfig['lbcf'].routerIndex = index;
        } else if (item.path === 'resource') {
          addonRouterConfig['tapp'].routerIndex = index;
        }
      });

      let keys = Object.keys(addons);

      keys.forEach(key => {
        if (key === 'LBCF') {
          let lbcfConfig = addonRouterConfig['lbcf'];
          newRouterList[lbcfConfig.routerIndex].sub.push(lbcfConfig.routerConfig);
        } else if (key === 'TappController') {
          let tappConfig = addonRouterConfig['tapp'];
          newRouterList[tappConfig.routerIndex].sub.push(tappConfig.routerConfig);
        }
      });
      this.setState({ finalSubRouterList: newRouterList });
    }
    /** =================== 这里是判断二级菜单路由的配置 ====================== */
  }

  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);

    let content: JSX.Element;

    let mode = urlParam['mode'];
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
