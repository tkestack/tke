import * as React from 'react';
import {
  ConsoleModuleMapProps,
  reduceK8sRestfulPath,
  Method,
  reduceNetworkRequest,
  setConsoleAPIAddress
} from './helpers';
import { ResourceInfo, RequestParams } from './src/modules/common/models';
import { resourceConfig } from './config';
import { isEmpty } from './src/modules/common/utils';
import * as classnames from 'classnames';
import { Button, Icon, Text, Bubble, NavMenu, List } from '@tencent/tea-component';

// @ts-ignore
const routerSea = seajs.require('router');

/**平台管理员,业务成员,游客,未初始化 */
enum UserType {
  admin = 'admin',
  member = 'member',
  other = 'other',
  init = 'init'
}

/** 获取当前控制台modules的 域名映射表 */
export enum ConsoleModuleEnum {
  /** tke-apiserver 版本 */
  PLATFORM = 'platform',

  /** 业务的版本详情 */
  Business = 'business',

  /** 通知模块 */
  Notify = 'notify',

  /** 告警模块 */
  Monitor = 'monitor',

  /** 镜像仓库 */
  Registry = 'registry',

  /** 日志模块 */
  LogAgent = 'logagent',

  /** 认证模块 */
  Auth = 'auth',

  /** 审计模块 */
  Audit = 'audit'
}

export enum PlatformTypeEnum {
  /** 平台 */
  Manager = 'manager',

  /** 业务 */
  Business = 'business'
}

interface RouterConfig {
  /** 导航的路由 */
  url?: string;

  /** 当前路由的名称 */
  title: string;

  /** 依赖于平台组件的安装 */
  watchModule: ConsoleModuleEnum | ConsoleModuleEnum[];

  /** 是否为小标题 */
  isTitle?: boolean;

  /** 二级下拉列表的配置 */
  subRouterConfig?: RouterConfig[];
}

/** 基础的侧边栏导航栏配置 */
const commonRouterConfig: RouterConfig[] = [
  {
    url: '/tkestack/overview',
    title: '概览',
    watchModule: ConsoleModuleEnum.Monitor
  },
  {
    url: '/tkestack/cluster',
    title: '集群管理',
    watchModule: ConsoleModuleEnum.PLATFORM
  },
  {
    url: '/tkestack/project',
    title: '业务管理',
    watchModule: ConsoleModuleEnum.Business
  },
  {
    url: '/tkestack/addon',
    title: '扩展组件',
    watchModule: ConsoleModuleEnum.PLATFORM
  },
  {
    title: '组织资源',
    watchModule: [ConsoleModuleEnum.Registry, ConsoleModuleEnum.Auth],
    subRouterConfig: [
      {
        url: '/tkestack/registry/repo',
        title: '镜像仓库管理',
        watchModule: ConsoleModuleEnum.Registry
      },
      {
        url: '/tkestack/registry/chart',
        title: 'Chart包仓库管理',
        watchModule: ConsoleModuleEnum.Registry
      },
      {
        url: '/tkestack/registry/apikey',
        title: '访问凭证',
        watchModule: ConsoleModuleEnum.Auth
      }
    ]
  },
  {
    title: '访问管理',
    watchModule: [ConsoleModuleEnum.Auth],
    subRouterConfig: [
      {
        url: '/tkestack/uam/user',
        title: '用户管理',
        watchModule: ConsoleModuleEnum.Auth
      },
      {
        url: '/tkestack/uam/strategy',
        title: '策略管理',
        watchModule: ConsoleModuleEnum.Auth
      }
    ]
  },
  {
    title: '监控&告警',
    watchModule: [ConsoleModuleEnum.Monitor, ConsoleModuleEnum.Notify],
    subRouterConfig: [
      {
        url: '/tkestack/alarm',
        title: '告警设置',
        watchModule: ConsoleModuleEnum.Monitor
      },
      {
        url: '/tkestack/notify',
        title: '通知设置',
        watchModule: ConsoleModuleEnum.Notify
      },
      {
        url: '/tkestack/alarm-record',
        title: '告警记录',
        watchModule: ConsoleModuleEnum.Notify
      }
    ]
  },
  {
    title: '运维中心',
    watchModule: [ConsoleModuleEnum.PLATFORM, ConsoleModuleEnum.Audit, ConsoleModuleEnum.LogAgent],
    subRouterConfig: [
      {
        url: '/tkestack/helm',
        title: 'Helm应用',
        watchModule: ConsoleModuleEnum.PLATFORM
      },
      {
        url: '/tkestack/log',
        title: '日志采集',
        watchModule: ConsoleModuleEnum.LogAgent
      },
      {
        url: '/tkestack/persistent-event',
        title: '事件持久化',
        watchModule: ConsoleModuleEnum.PLATFORM
      },
      {
        url: '/tkestack/audit',
        title: '审计记录',
        watchModule: ConsoleModuleEnum.Audit
      }
    ]
  }
];

/** 基础的侧边栏导航栏配置 */
const businessCommonRouterConfig: RouterConfig[] = [
  {
    url: '/tkestack-project/application',
    title: '应用管理',
    watchModule: ConsoleModuleEnum.Business
  },
  {
    url: '/tkestack-project/helm',
    title: 'Helm应用',
    watchModule: ConsoleModuleEnum.PLATFORM
  },
  {
    url: '/tkestack-project/project',
    title: '业务管理',
    watchModule: ConsoleModuleEnum.Business
  },
  {
    title: '组织资源',
    watchModule: [ConsoleModuleEnum.Registry, ConsoleModuleEnum.Auth],
    subRouterConfig: [
      {
        url: '/tkestack-project/registry/repo',
        title: '仓库管理',
        watchModule: ConsoleModuleEnum.Registry
      },
      {
        url: '/tkestack-project/registry/apikey',
        title: '访问凭证',
        watchModule: ConsoleModuleEnum.Auth
      }
    ]
  },
  {
    title: '监控&告警',
    watchModule: [ConsoleModuleEnum.Monitor, ConsoleModuleEnum.Notify],
    subRouterConfig: [
      {
        url: '/tkestack-project/alarm',
        title: '告警设置',
        watchModule: ConsoleModuleEnum.Monitor
      },
      {
        url: '/tkestack-project/notify',
        title: '通知设置',
        watchModule: ConsoleModuleEnum.Notify
      }
    ]
  },
  {
    title: '运维中心',
    watchModule: [ConsoleModuleEnum.LogAgent],
    subRouterConfig: [
      {
        url: '/tkestack-project/log',
        title: '日志采集',
        watchModule: ConsoleModuleEnum.LogAgent
      }
    ]
  }
];

interface ConsoleWrapperProps {
  /** 平台侧业务侧 */
  platformType: PlatformTypeEnum;

  /** 是否需要侧边导航栏 */
  sideBar?: boolean;
}

interface UserInfo {
  extra: any;
  groups: string[];
  name: string;
}

interface Project {
  id: string;
  name: string;
}

interface ConsoleWrapperState {
  /** 当前选中的路由 */
  selected: string;

  /** 当前折叠的路由 */
  toggleName: string;

  /** 用户的id */
  userInfo: UserInfo;

  /** 平台切换按钮选项展开 */
  isShowPlatformSwitch: boolean;

  /** 控制台的api映射 */
  consoleApiMap: ConsoleModuleMapProps;

  /** 该用户是否为平台管理员,业务成员,游客 */
  userType: UserType;

  /**该用户负责的业务 */
  projects: Project[];

  /** 是否展示user的下拉框 */
  isShowUserDropdown: boolean;

  /** 最终展示的路由 */
  routerConfig: RouterConfig[];

  /** 判断二级路由是否开启 */
  asideRouterSelect: {
    index: number;
    isShow: boolean;
  };
}

export class Wrapper extends React.Component<ConsoleWrapperProps, ConsoleWrapperState> {
  constructor(props: ConsoleWrapperProps) {
    super(props);
    this.state = {
      selected: location.pathname.split('/').slice(0, 4).join('/'),
      toggleName: '',
      userInfo: {
        extra: '',
        groups: [],
        name: ''
      },
      consoleApiMap: window['modules'] || {},
      isShowPlatformSwitch: false,
      userType: UserType.init,
      projects: [],
      isShowUserDropdown: false,
      routerConfig: [],
      asideRouterSelect: {
        index: -1,
        isShow: false
      }
    };
  }

  async componentWillMount() {
    await this.getConsoleModule();
    this.state.userInfo.name === '' && (await this.getUserInfo());
    this.state.userType === UserType.init && (await this.getUserProjectInfo());
  }

  //获取用户信息包括用户业务信息
  async getUserInfo() {
    let infoResourceInfo: ResourceInfo = resourceConfig()['info'];
    let url = reduceK8sRestfulPath({ resourceInfo: infoResourceInfo });
    let params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      let response = await reduceNetworkRequest(params);
      this.setState({
        userInfo: response.data
      });
    } catch (error) {}
  }

  /**
   * 获取当前版本支持的模块，如 是否有tcr
   */
  async getConsoleModule() {
    let moduleResourceInfo: ResourceInfo = resourceConfig()['module'];
    let url = reduceK8sRestfulPath({ resourceInfo: moduleResourceInfo });
    let params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      let consoleApiMap;
      if (isEmpty(this.state.consoleApiMap)) {
        let response = await reduceNetworkRequest(params);
        consoleApiMap = response.data.components;

        // 设置全局的变量，console的值
        setConsoleAPIAddress(consoleApiMap);
        this.setState({ consoleApiMap });
      } else {
        consoleApiMap = this.state.consoleApiMap;
      }

      // 进行路由的更新
      let moduleKeys = Object.keys(consoleApiMap);
      let initRouterConfig =
        this.props.platformType === PlatformTypeEnum.Business ? businessCommonRouterConfig : commonRouterConfig;
      let currentRouterConfig: RouterConfig[] = initRouterConfig.filter((routerConfig, index) => {
        if (Array.isArray(routerConfig.watchModule)) {
          return routerConfig.watchModule.some(item => moduleKeys.includes(item));
        }
        return moduleKeys.includes(routerConfig.watchModule);
      });

      // 过滤二级路由
      currentRouterConfig.forEach(routerConfig => {
        let subRouterConfig = routerConfig.subRouterConfig;
        if (subRouterConfig) {
          // 重写subRouterConfig属性
          routerConfig.subRouterConfig = subRouterConfig.filter(subRouterConf => {
            if (Array.isArray(subRouterConf.watchModule)) {
              return subRouterConf.watchModule.some(item => moduleKeys.includes(item));
            }
            return moduleKeys.includes(subRouterConf.watchModule);
          });
        }
      });

      // 二级路由信息的初始化
      let subRouterIndex = -1;
      // 进行二级路由信息的初始化
      currentRouterConfig.forEach((routerConfig, index) => {
        // 进行二级路由信息的初始化
        if (subRouterIndex < 0 && routerConfig.subRouterConfig) {
          let subRouterUrl = routerConfig.subRouterConfig.map(item => item.url);
          if (subRouterUrl.includes(this.state.selected)) {
            subRouterIndex = index;
          }
        }
      });

      this.setState({
        routerConfig: currentRouterConfig,
        asideRouterSelect: {
          index: subRouterIndex,
          isShow: subRouterIndex > -1
        }
      });
    } catch (error) {}
  }

  //获取用户信息包括用户业务信息
  async getUserProjectInfo() {
    let userResourceInfo: ResourceInfo = resourceConfig().portal;
    let url = reduceK8sRestfulPath({ resourceInfo: userResourceInfo });
    let params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      let response = await reduceNetworkRequest(params);
      if (response.code === 0) {
        let projects = Object.keys(response.data.projects).map(key => {
          return {
            id: key,
            name: response.data.projects[key]
          };
        });
        let userType = response.data.administrator
          ? UserType.admin
          : projects.length !== 0
          ? UserType.member
          : UserType.other;
        this.setState({
          userType,
          projects
        });
        if (userType === UserType.member && this.props.platformType === PlatformTypeEnum.Manager) {
          location.href = location.origin + '/tkestack-project/application';
        } else if (
          userType === UserType.admin &&
          projects.length === 0 &&
          this.props.platformType === PlatformTypeEnum.Business
        ) {
          location.href = location.origin + '/tkestack';
        } else if (userType === UserType.other && window.location.pathname.indexOf('tkestack/blank') === -1) {
          window.location.pathname = 'tkestack/blank';
        }
      }
    } catch (error) {}
  }

  // 退出页面
  async userLogout() {
    let logoutInfo: ResourceInfo = resourceConfig().logout;
    let url = reduceK8sRestfulPath({ resourceInfo: logoutInfo });
    let params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      let response = await reduceNetworkRequest(params);
    } catch (error) {}
  }

  /**进行路由的跳转 */
  onNav(path: string) {
    this.setState({ selected: path });
    routerSea.navigate(path);
  }

  /** 选择折叠的路由 */
  onToggleName(name: string) {
    this.setState({ toggleName: name });
  }

  _handleHoverForFlatformSwitch(isShow: boolean = false) {
    this.setState({
      isShowPlatformSwitch: isShow
    });
  }

  render() {
    let query = window.location.search;
    let finalContent: React.ReactNode;

    if (isEmpty(this.state.consoleApiMap)) {
      finalContent = <noscript />;
    } else {
      let { sideBar = true } = this.props;
      finalContent = (
        <React.Fragment>
          {this._renderTopBar(query)}

          <div className="qc-animation-empty container container-tke2-cluster" id="tkestack" style={{ left: 0 }}>
            {sideBar && this._renderSideBar(query)}

            <div id="appArea" className="main" style={sideBar ? {} : { left: 0 }}>
              {this.props.children}
            </div>
          </div>
        </React.Fragment>
      );
    }
    return finalContent;
  }

  /**
   * 展示顶部导航栏
   */
  private _renderTopBar(query: string) {
    return (
      <NavMenu
        left={
          <React.Fragment>
            <NavMenu.Item type="logo">
              <img src="/static/icon/logo.svg" style={{ height: '30px' }} alt="logo" />
            </NavMenu.Item>
          </React.Fragment>
        }
        right={
          <React.Fragment>
            <NavMenu.Item
              type="dropdown"
              overlay={() => (
                <List type="option">
                  <List.Item
                    onClick={async () => {
                      await this.userLogout();
                      location.reload();
                    }}
                  >
                    退出
                  </List.Item>
                </List>
              )}
            >
              {this.state.userInfo.name}
            </NavMenu.Item>
          </React.Fragment>
        }
      />
    );
  }

  /**
   * 展示侧边导航栏
   */
  private _renderSideBar(query: string) {
    let { platformType } = this.props;
    let { userType, projects } = this.state;
    let routerConfig: RouterConfig[] = this.state.routerConfig;
    return (
      <div className="aside qc-aside-new">
        <div className="qc-aside-area">
          <div className="qc-aside-area-main">
            <h2 className="qc-aside-headline">
              <Text verticalAlign="middle">{platformType === PlatformTypeEnum.Manager ? '平台管理' : '业务管理'}</Text>
              {userType === UserType.admin && projects.length ? (
                <Bubble
                  content={platformType === PlatformTypeEnum.Manager ? '切换至业务管理控制台' : '切换至平台管理控制台'}
                  placement="right"
                >
                  <Icon
                    type="convertip"
                    className="tea-ml-2n"
                    style={{ verticalAlign: '-9px' }}
                    onClick={() => {
                      location.href =
                        location.origin +
                        (platformType === PlatformTypeEnum.Manager ? '/tkestack-project' : '/tkestack');
                    }}
                  />
                </Bubble>
              ) : (
                <noscript />
              )}
            </h2>
            <ul className="qc-aside-list def-scroll keyboard-focus-obj">
              {routerConfig.map((routerIns, index) => {
                let routerContent: React.ReactNode;
                if (routerIns.isTitle) {
                  routerContent = (
                    <li className="qc-aside-title">
                      <span>{routerIns.title}</span>
                    </li>
                  );
                } else {
                  let isSelected =
                    this.state.selected.includes(routerIns.url) ||
                    (this.state.selected.split('/').length <= 2 && index === 0);

                  /** 需要判断当前路由设置是否为二级路由设置 */
                  if (routerIns.subRouterConfig) {
                    let subRouterUrl = routerIns.subRouterConfig.map(item => item.url);
                    isSelected = subRouterUrl.includes(this.state.selected);
                    let selectedIndex = subRouterUrl.findIndex(item => item === this.state.selected);
                    let { index: asideIndex, isShow } = this.state.asideRouterSelect;

                    routerContent = (
                      <li
                        key={index}
                        className={asideIndex === index && isShow ? 'qc-aside-select qc-aside-child-select' : ''}
                      >
                        <a
                          style={{ paddingLeft: '24px' }}
                          className="qc-aside-level-1"
                          href="javascript:;"
                          onClick={() => {
                            this.setState({
                              asideRouterSelect: {
                                index,
                                isShow: asideIndex !== index ? true : !isShow
                              }
                            });
                          }}
                        >
                          <Text style={{ marginLeft: 0 }}>{routerIns.title}</Text>
                          <i className="qc-aside-up-icon" />
                        </a>
                        <ul className="qc-aside-subitem">
                          {routerIns.subRouterConfig.map((subRouter, subIndex) => {
                            return (
                              <li key={subIndex}>
                                <a
                                  className={classnames('qc-aside-level-2', {
                                    'qc-aside-select': selectedIndex === subIndex
                                  })}
                                  href="javascript:;"
                                  onClick={() => {
                                    if (selectedIndex !== subIndex) {
                                      this.onNav(subRouter.url);
                                    }
                                  }}
                                  target="_self"
                                >
                                  <span>{subRouter.title}</span>
                                </a>
                              </li>
                            );
                          })}
                        </ul>
                      </li>
                    );
                  } else {
                    routerContent = (
                      <li key={index}>
                        <a
                          style={{ paddingLeft: '24px' }}
                          className={classnames('qc-aside-level-1', {
                            'qc-aside-select': isSelected
                          })}
                          href="javascript:;"
                          onClick={e => {
                            if (!isSelected) {
                              // 这里需要区分是否为别的业务，如果是别的业务，是进行业务的跳转
                              if (this.props.platformType === PlatformTypeEnum.Manager) {
                                this.onNav(routerIns.url);
                              } else {
                                this.onNav(routerIns.url + query);
                              }
                            }
                          }}
                          target="_self"
                        >
                          <span style={isSelected ? { marginLeft: 0, color: '#4093ff' } : { marginLeft: 0 }}>
                            {routerIns.title}
                          </span>
                        </a>
                      </li>
                    );
                  }
                }
                return routerContent;
              })}
            </ul>
          </div>
        </div>
      </div>
    );
  }
}
