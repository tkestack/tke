import * as React from 'react';
import { RootProps } from './NotifyApp';
import * as classnames from 'classnames';
import { router } from '../router';
import { BasicRouter, SubRouter } from '../../cluster/models';

export interface ResourceListProps extends RootProps {
  /** subRouterList */
  subRouterList: SubRouter[];
}

interface ResourceSidebarState {
  /** 当前选中的resouce */
  currentPath?: string;

  /** 是否触发二级导航栏 */
  isOpenSecondBar?: string;

  /** 判断触发那个二级导航栏 */
  secondBarPath?: string;
}

export class ResourceSidebar extends React.Component<ResourceListProps, ResourceSidebarState> {
  constructor(props: ResourceListProps) {
    super(props);
    let urlParams = router.resolve(props.route);
    this.state = {
      currentPath: urlParams['resourceName'] || 'channel',
      isOpenSecondBar: '',
      secondBarPath: ''
    };
  }

  componentWillReceiveProps(nextProps: ResourceListProps) {
    let { route, subRouterList } = nextProps,
      urlParams = router.resolve(route),
      { type = '', resourceName = '' } = urlParams;

    let oldSubRouterList = JSON.stringify(this.props.subRouterList),
      newSubRouterList = JSON.stringify(subRouterList);

    /**
     * 这里判断二级路由是否加载完毕，并且判断是否为当前的路由是否有变化，无变化则不需要判断
     * condition:
     * 1. 首次进入resourceContainer的时候，会去触发拉取二级菜单的配置，会触发oldSubRouterList !== newSubRouterList，进行第一次state初始化
     * 2. 由于请求subRouter之后，会在finish的时候，进行动态菜单的判断，所有 subRouterList是动态变化的，参考条件1，可以进行菜单栏的展开
     * 3. 后续的点击过程中，会触发 前者的判断条件，subRouter有了，url变化，动态对菜单进行更新
     */
    if (
      (subRouterList.length && (type !== this.state.secondBarPath || resourceName !== this.state.currentPath)) ||
      oldSubRouterList !== newSubRouterList
    ) {
      // 处理state的值的变化
      this.setState({
        isOpenSecondBar: subRouterList.find(item => item.path === type) ? type : '',
        secondBarPath: type,
        currentPath: resourceName || 'channel'
      });
    }
  }

  render() {
    return (
      <div className="secondary-aside">
        <div className="secondary-aside-area">
          <div className="secondary-aside-area-main">{this._renderBarList()}</div>
        </div>
      </div>
    );
  }

  /**
   * 生成二级导航路由主体部分
   */
  _renderBarList() {
    let { subRouterList } = this.props;

    return (
      <ul className="secondary-aside-list">
        {subRouterList.map((sidebar, index) => {
          if (sidebar.sub) {
            return (
              <li
                key={index}
                className={classnames('', {
                  'secondary-aside-select': this.state.isOpenSecondBar === sidebar.path
                })}
              >
                <a
                  href="javascript:;"
                  className="secondary-aside-level-1"
                  onClick={() => {
                    this._handleClickForFirstBar(sidebar.path, true);
                  }}
                >
                  <span>{sidebar.name}</span>
                  <i className="secondary-aside-up-icon" />
                </a>
                {this._renderSecondBarList(sidebar.sub, sidebar.path)}
              </li>
            );
          } else {
            return (
              <li key={index}>
                <a
                  href="javascript:;"
                  onClick={e => {
                    this._handleClickForFirstBar(sidebar.basicUrl, false, sidebar.path);
                    e.stopPropagation();
                  }}
                  className={classnames('secondary-aside-level-1', {
                    'secondary-aside-select': this.state.currentPath === sidebar.basicUrl
                  })}
                >
                  <span>{sidebar.name}</span>
                </a>
              </li>
            );
          }
        })}
      </ul>
    );
  }

  /**
   * 路由公共处理部分，以及数据数据请求
   * @param subSidebarPath  二级路由
   * @param sidebarPath   一级路由
   */
  private _handleDataFetcher(subSidebarPath: string, sidebarPath: string) {
    let { actions, route } = this.props,
      urlParams = router.resolve(route);

    // 避免重复点击，进行重复的操作
    if (urlParams['resourceName'] !== subSidebarPath || urlParams.mode !== 'list') {
      // 初始化resource fetcher的相关配置信息，因为多个resource用的是同一个fetcher
      actions.resource[subSidebarPath] && actions.resource[subSidebarPath].reset();

      // 路由的跳转
      router.navigate(
        Object.assign({}, urlParams, {
          mode: 'list',
          type: sidebarPath,
          resourceName: subSidebarPath
        }),
        Object.assign({}, route.queries)
      );
      actions.resource[subSidebarPath] && actions.resource[subSidebarPath].fetch();
    }
  }

  /**
   * 处理一级导航的操作
   * @param path  跳转的路由
   * @param isClickNested 是否点击了含有二级的路由
   */
  private _handleClickForFirstBar(path: string, isClickNested: boolean = false, sidebarPath?: string) {
    if (!isClickNested) {
      // 因为是非折叠的路由，所有openSecondBar置为空
      this.setState({
        isOpenSecondBar: ''
      });
      this._handleDataFetcher(path, sidebarPath);
    } else {
      this.setState({
        isOpenSecondBar: path
      });
    }
  }

  /**
   * 处理二级导航的操作
   * @param subSidebarpath 跳转的路由
   * @param sidebarPath   一级路由
   */
  private _handleClickForSecondBar(subSidebarpath: string, sidebarPath: string) {
    this.setState({
      currentPath: subSidebarpath
    });
    this._handleDataFetcher(subSidebarpath, sidebarPath);
  }

  /**
   * 生成二级导航栏
   */
  private _renderSecondBarList(subMenu: BasicRouter[], sidebarPath: string) {
    let subMenuList = subMenu.map((subSidebar, index) => {
      return (
        <li key={index}>
          <a
            href="javascript:;"
            onClick={e => {
              this._handleClickForSecondBar(subSidebar.path, sidebarPath);
              e.stopPropagation();
            }}
            className={classnames('secondary-aside-level-2', {
              'secondary-aside-select': this.state.currentPath === subSidebar.path
            })}
          >
            {subSidebar.name}
          </a>
        </li>
      );
    });

    return (
      <ul className="secondary-aside-subitem" style={{ paddingBottom: '0' }}>
        {subMenuList}
      </ul>
    );
  }
}
