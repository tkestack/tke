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
import { TkeVersion } from '@/src/modules/common/components/tke-version';
import { PermissionProvider, checkCustomVisible } from '@common/components/permission-provider';
import { getCustomConfig } from '@config';
import { ConsoleModuleEnum } from '@config/platform';
import { insertCSS } from '@tencent/ff-redux';
import * as React from 'react';
import { ExternalLink, Layout, List, Menu, NavMenu, StatusTip } from 'tea-component';
import 'tea-component/dist/tea.css';
import { PlatformTypeEnum, resourceConfig } from './config';
import {
  ConsoleModuleMapProps,
  Method,
  reduceK8sRestfulPath,
  reduceNetworkRequest,
  setConsoleAPIAddress,
  isInIframe
} from './helpers';
import { RequestParams, ResourceInfo } from './src/modules/common/models';
import { isEmpty } from './src/modules/common/utils';

require('promise.prototype.finally').shim();

insertCSS(
  'custom_theme_menu',
  `
  ._custom_theme_menu .tea-menu__list li.is-selected>.tea-menu__item {
    background: #006eff;
  }
`
);

const { LoadingTip } = StatusTip;

const routerSea = seajs.require('router');

/**平台管理员,业务成员,游客,未初始化 */
enum UserType {
  admin = 'admin',
  member = 'member',
  other = 'other',
  init = 'init'
}

export interface IPlatformContext {
  type: PlatformTypeEnum;
}

export const PlatformContext = React.createContext<IPlatformContext>({ type: PlatformTypeEnum.Manager });

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

  icon?: [string, string];

  key?: string;
}

console.log('svg----->', require('./public/static/icon/overview.svg'));

/** 基础的侧边栏导航栏配置 */
const commonRouterConfig: RouterConfig[] = [
  {
    url: '/tkestack/overview',
    title: '概览',
    watchModule: ConsoleModuleEnum.Monitor,
    icon: [require('./public/static/icon/overview.svg'), require('./public/static/icon/overview-hover.svg')]
  },
  {
    url: '/tkestack/cluster',
    title: '集群管理',
    watchModule: ConsoleModuleEnum.PLATFORM,
    icon: [require('./public/static/icon/cluster.svg'), require('./public/static/icon/cluster-hover.svg')]
  },
  ...(checkCustomVisible('platform.project')
    ? [
        {
          url: '/tkestack/project',
          title: '业务管理',
          watchModule: ConsoleModuleEnum.Business,
          key: 'project'
        }
      ]
    : []),
  ...(checkCustomVisible('platform.addon')
    ? [
        {
          url: '/tkestack/addon',
          title: '扩展组件',
          watchModule: ConsoleModuleEnum.PLATFORM,
          key: 'addon'
        }
      ]
    : []),
  {
    title: '组织资源',
    icon: [require('./public/static/icon/registry.svg'), require('./public/static/icon/registry-hover.svg')],
    watchModule: [ConsoleModuleEnum.Registry, ConsoleModuleEnum.Auth],
    subRouterConfig: [
      {
        url: '/tkestack/registry/repo',
        title: '镜像仓库管理',
        watchModule: ConsoleModuleEnum.Registry
      },
      // {
      //   url: '/tkestack/registry/chartgroup',
      //   title: 'Helm仓库',
      //   watchModule: ConsoleModuleEnum.Registry
      // },
      {
        url: '/tkestack/registry/chart',
        title: 'Helm模板',
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
    icon: [require('./public/static/icon/uam.svg'), require('./public/static/icon/uam-hover.svg')],
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
    icon: [require('./public/static/icon/alarm.svg'), require('./public/static/icon/alarm-hover.svg')],
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
    icon: [require('./public/static/icon/app.svg'), require('./public/static/icon/app-hover.svg')],
    watchModule: [
      ConsoleModuleEnum.Application,
      ConsoleModuleEnum.PLATFORM,
      ConsoleModuleEnum.Audit,
      ConsoleModuleEnum.LogAgent
    ],
    subRouterConfig: [
      {
        url: '/tkestack/application/app',
        title: 'Helm应用',
        watchModule: ConsoleModuleEnum.Application
      },
      // {
      //   url: '/tkestack/helm',
      //   title: 'Helm2应用',
      //   watchModule: ConsoleModuleEnum.PLATFORM
      // },
      {
        url: '/tkestack/log',
        title: '日志采集',
        watchModule: ConsoleModuleEnum.LogAgent
      },
      ...(checkCustomVisible('platform.persistent-event')
        ? [
            {
              url: '/tkestack/persistent-event',
              title: '事件持久化',
              watchModule: ConsoleModuleEnum.PLATFORM
            }
          ]
        : []),
      {
        url: '/tkestack/audit',
        title: '审计记录',
        watchModule: ConsoleModuleEnum.Audit
      }
    ]
  },
  {
    icon: [require('./public/static/icon/data-service.svg'), require('./public/static/icon/data-service-hover.svg')],
    title: '数据服务',
    watchModule: [ConsoleModuleEnum.Middleware],
    subRouterConfig: [
      {
        url: '/tkestack/middleware',
        title: '中间件列表',
        watchModule: ConsoleModuleEnum.Middleware
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
        title: '镜像仓库管理',
        watchModule: ConsoleModuleEnum.Registry
      },
      // {
      //   url: '/tkestack-project/registry/chartgroup',
      //   title: 'Helm仓库',
      //   watchModule: ConsoleModuleEnum.Registry
      // },
      {
        url: '/tkestack-project/registry/chart',
        title: 'Helm模板',
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
    title: '运维管理',
    watchModule: [ConsoleModuleEnum.Application, ConsoleModuleEnum.LogAgent, ConsoleModuleEnum.PLATFORM],
    subRouterConfig: [
      {
        url: '/tkestack-project/app/app',
        title: 'Helm应用',
        watchModule: ConsoleModuleEnum.Application
      },
      // {
      //   url: '/tkestack-project/helm',
      //   title: 'Helm2应用',
      //   watchModule: ConsoleModuleEnum.PLATFORM
      // },
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

  children: React.ReactNode;
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
    const infoResourceInfo: ResourceInfo = resourceConfig()['info'];
    const url = reduceK8sRestfulPath({ resourceInfo: infoResourceInfo });
    const params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      const response = await reduceNetworkRequest(params);
      console.log(response?.data);
      this.setState({
        userInfo: response.data
      });
    } catch (error) {}
  }

  /**
   * 获取当前版本支持的模块，如 是否有tcr
   */
  async getConsoleModule() {
    const moduleResourceInfo: ResourceInfo = resourceConfig()['module'];
    const url = reduceK8sRestfulPath({ resourceInfo: moduleResourceInfo });
    const params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      let consoleApiMap;
      if (isEmpty(this.state.consoleApiMap)) {
        const response = await reduceNetworkRequest(params);
        consoleApiMap = response.data.components;

        // 设置全局的变量，console的值
        setConsoleAPIAddress(consoleApiMap);
        this.setState({ consoleApiMap });
      } else {
        consoleApiMap = this.state.consoleApiMap;
      }

      // 进行路由的更新
      const moduleKeys = Object.keys(consoleApiMap);
      const initRouterConfig =
        this.props.platformType === PlatformTypeEnum.Business ? businessCommonRouterConfig : commonRouterConfig;
      const currentRouterConfig: RouterConfig[] = initRouterConfig.filter((routerConfig, index) => {
        if (Array.isArray(routerConfig.watchModule)) {
          return routerConfig.watchModule.some(item => moduleKeys.includes(item));
        }
        return moduleKeys.includes(routerConfig.watchModule);
      });

      // 过滤二级路由
      currentRouterConfig.forEach(routerConfig => {
        const subRouterConfig = routerConfig.subRouterConfig;
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
          const subRouterUrl = routerConfig.subRouterConfig.map(item => item.url);
          if (subRouterUrl.includes(this.state.selected)) {
            subRouterIndex = index;
          }
        }
      });
      //追加数据服务菜单
      if (!currentRouterConfig?.every(item => item?.watchModule?.includes(ConsoleModuleEnum.Middleware))) {
        currentRouterConfig.push({
          icon: [
            require('./public/static/icon/data-service.svg'),
            require('./public/static/icon/data-service-hover.svg')
          ],
          title: '数据服务',
          watchModule: [ConsoleModuleEnum.Middleware],
          subRouterConfig: [
            {
              url: '/tkestack/middleware',
              title: '中间件列表',
              watchModule: ConsoleModuleEnum.Middleware
            }
          ]
        });
      }
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
    const userResourceInfo: ResourceInfo = resourceConfig().portal;
    const url = reduceK8sRestfulPath({ resourceInfo: userResourceInfo });
    const params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      const response = await reduceNetworkRequest(params);
      if (response.code === 0) {
        const projects = Object.keys(response.data.projects).map(key => {
          return {
            id: key,
            name: response.data.projects[key]
          };
        });
        const userType = response.data.administrator
          ? UserType.admin
          : projects.length !== 0
          ? UserType.member
          : UserType.other;
        this.setState({
          userType,
          projects
        });
        const isInBlankPage = window.location.pathname.indexOf('tkestack/blank') !== -1;
        if (userType === UserType.member && this.props.platformType === PlatformTypeEnum.Manager) {
          location.href = location.origin + '/tkestack-project/application';
        } else if (
          userType === UserType.admin &&
          projects.length === 0 &&
          this.props.platformType === PlatformTypeEnum.Business
        ) {
          location.href = location.origin + '/tkestack';
        } else if (userType === UserType.other && !isInBlankPage) {
          window.location.pathname = 'tkestack/blank';
        } else if (isInBlankPage) {
          if (userType === UserType.admin) {
            location.href = location.origin + '/tkestack';
          } else if (userType === UserType.member) {
            location.href = location.origin + '/tkestack-project/application';
          }
        }
      }
    } catch (error) {}
  }

  // 退出页面
  async userLogout() {
    const logoutInfo: ResourceInfo = resourceConfig().logout;
    const url = reduceK8sRestfulPath({ resourceInfo: logoutInfo });
    const params: RequestParams = {
      method: Method.get,
      url
    };
    try {
      const response = await reduceNetworkRequest(params);
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

  _handleHoverForFlatformSwitch(isShow = false) {
    this.setState({
      isShowPlatformSwitch: isShow
    });
  }

  render() {
    const query = window.location.search;
    let finalContent: React.ReactNode;

    if (isEmpty(this.state.consoleApiMap)) {
      finalContent = <noscript />;
    } else {
      const { sideBar = true } = this.props;
      finalContent = (
        <React.Fragment>
          {this._renderTopBar(query)}

          <div
            className="qc-animation-empty container container-tke2-cluster"
            id="tkestack"
            style={{ left: 0, top: isInIframe() ? 0 : undefined }}
          >
            {sideBar && this._renderSideBar(query)}

            <div id="appArea" className="main" style={sideBar ? {} : { left: 0 }}>
              {this.props.children}
            </div>
          </div>
        </React.Fragment>
      );
    }
    return (
      <PlatformContext.Provider value={{ type: this.props.platformType }}>
        <React.Suspense
          fallback={
            <div
              style={{
                width: '100vw',
                height: '100vh',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center'
              }}
            >
              <LoadingTip />
            </div>
          }
        >
          {finalContent}
        </React.Suspense>
      </PlatformContext.Provider>
    );
  }

  /**
   * 展示顶部导航栏
   */
  private _renderTopBar(query: string) {
    return (
      <NavMenu
        left={
          <React.Fragment>
            <NavMenu.Item>
              <img
                src={`/static/icon/${getCustomConfig()?.logoDir ?? 'default'}/logo.svg`}
                style={{ height: '30px' }}
                alt="logo"
              />
            </NavMenu.Item>
          </React.Fragment>
        }
        right={
          <React.Fragment>
            <PermissionProvider value="platform.overview.help">
              <NavMenu.Item>
                <ExternalLink href={'https://tkestack.github.io/docs/'}>容器服务帮助手册</ExternalLink>
              </NavMenu.Item>
            </PermissionProvider>

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

  private _renderSideBar(query: string) {
    const title = this.props.platformType === PlatformTypeEnum.Manager ? 'TKE 平台管理' : '业务管理';
    const routerList = this.state.routerConfig;
    const selected = this.state.selected;
    const openedIndex = this.state.asideRouterSelect.index;

    return (
      <Layout.Sider>
        <Menu theme="dark" className="_custom_theme_menu" title={title}>
          {routerList.map(({ url, title, subRouterConfig, icon }, index) => {
            if (subRouterConfig) {
              return (
                <Menu.SubMenu
                  title={title}
                  icon={icon}
                  opened={openedIndex === index}
                  onOpenedChange={opened => {
                    this.setState({
                      asideRouterSelect: {
                        index: opened ? index : -1,
                        isShow: this.state.asideRouterSelect.isShow
                      }
                    });
                  }}
                >
                  {subRouterConfig?.map(({ title, url, icon }) => (
                    <Menu.Item
                      title={title}
                      icon={icon}
                      selected={url === selected}
                      onClick={() => {
                        if (url === selected) return;

                        this.onNav(url);
                      }}
                    />
                  ))}
                </Menu.SubMenu>
              );
            } else {
              const isSelected = selected.includes(url) || (selected.split('/').length <= 2 && index === 0);

              return (
                <Menu.Item
                  title={title}
                  icon={icon}
                  selected={isSelected}
                  onClick={() => {
                    if (isSelected) return;

                    if (this.props.platformType === PlatformTypeEnum.Manager) {
                      this.onNav(url);
                    } else {
                      this.onNav(url + query);
                    }
                  }}
                />
              );
            }
          })}
        </Menu>

        <TkeVersion />
      </Layout.Sider>
    );
  }
}
