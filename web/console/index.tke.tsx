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
import './i18n';
import { Entry, insertCSS } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Alert, Button, Text } from '@tencent/tea-component';
import * as React from 'react';
import { Wrapper } from './Wrapper';
import { Init_Forbiddent_Config } from './helpers/reduceNetwork';
import { Addon } from './src/modules/addon';
import { AlarmRecord } from './src/modules/alarmRecord';
import { Audit } from './src/modules/audit';
import { TipDialog } from './src/modules/common';
import { PersistentEvent } from './src/modules/persistentEvent';

// 公有云的图表组件为异步加载，这里为了减少路径配置，还是保留为同步加载，预先import即可变成不split
// import '@tencent/tchart/build/ChartsComponents';
import { Overview } from '@src/modules/overview';
import { VNCPage } from '@src/modules/vnc';
import { BlankPage } from './blankPage';
import { PlatformTypeEnum } from './config';
import { MiddlewareAppContainer } from './tencent/paas-midleware';

// const ClusterPromise = import(/* webpackPrefetch: true */ './src/modules/cluster');
// const Cluster = React.lazy(() => ClusterPromise);
import Cluster from './src/modules/cluster';

// const UamPromise = import(/* webpackPrefetch: true */ './src/modules/uam');
// const Uam = React.lazy(() => UamPromise);
import Uam from './src/modules/uam';

// const RegistryPromise = import(/* webpackPrefetch: true */ './src/modules/registry');
// const Registry = React.lazy(() => RegistryPromise);
import Registry from './src/modules/registry';

// const LogStashPromise = import(/* webpackPrefetch: true */ './src/modules/logStash');
// const LogStash = React.lazy(() => LogStashPromise);
import LogStash from './src/modules/logStash';

// const ProjectPromise = import(/* webpackPrefetch: true */ './src/modules/project');
// const Project = React.lazy(() => ProjectPromise);
import Project from './src/modules/project';

// const HelmPromise = import(/* webpackPrefetch: true */ './src/modules/helm');
// const Helm = React.lazy(() => HelmPromise);
import Helm from './src/modules/helm';

// const ApplicationPromise = import(/* webpackPrefetch: true */ './src/modules/application');
// const Application = React.lazy(() => ApplicationPromise);
import Application from './src/modules/application';

// const AlarmPolicyPromise = import(/* webpackPrefetch: true */ './src/modules/alarmPolicy');
// const AlarmPolicy = React.lazy(() => AlarmPolicyPromise);
import AlarmPolicy from './src/modules/alarmPolicy';

// const NotifyPromise = import(/* webpackPrefetch: true */ './src/modules/notify');
// const Notify = React.lazy(() => NotifyPromise);
import Notify from './src/modules/notify';

import { getCustomConfig } from '@config';

const Title = getCustomConfig()?.title ?? 'TKEStack';

insertCSS(
  'hidden-checkbox',
  `.hidden-checkbox .tea-form-check { display : none }
`
);

insertCSS(
  'singleDatePicker',
  `.tc-15-calendar-i-pre-m:hover span { display : block }
   .tc-15-calendar-i-pre-m span { display : none }
   .tc-15-calendar-i-next-m:hover span { display : block }
   .tc-15-calendar-i-next-m span { display : none }
`
);

insertCSS(
  'myTagSearchBox',
  `.myTagSearchBox{ width:100% !important; background-color: #fff; }
   .myTagSearchBox .tea-search__inner{ width:100% !important; }
   .myTagSearchBox .tea-tag-group{ float:left; }
   .myTagSearchBox .tea-search__tips{ float:left; }
   .myTagSearchBox.is-active .tea-search__tips{ display:none; }
`
);
insertCSS(
  'back-link',
  `
  .back-link { margin-right:24px; }
`
);

/** ======= hack外层，使其能够接受参数，触发state变化，进行侧边栏的更新 ======== */
interface TempWrapperProps {
  /** 当前业务的businessKey */
  businessKey: string;

  children: React.ReactNode;
}

class TempWrapper extends React.Component<TempWrapperProps, any> {
  render() {
    return <Wrapper platformType={PlatformTypeEnum.Manager}>{this.props.children}</Wrapper>;
  }
}
/** ======= hack外层，使其能够接受参数，触发state变化，进行侧边栏的更新 ======== */

/** ========================== 展示没有权限弹窗 ================================ */
export let changeForbiddentConfig;

interface ForbiddentDialogState {
  forbiddentConfig: { isShow: boolean; message: string };
}

class ForbiddentDialog extends React.Component<any, ForbiddentDialogState> {
  constructor(props: any) {
    super(props);
    this.state = {
      forbiddentConfig: Init_Forbiddent_Config
    };

    changeForbiddentConfig = (config: { isShow: boolean; message: string }) => {
      this.setState({ forbiddentConfig: config });
    };
  }

  render() {
    const { isShow, message } = this.state.forbiddentConfig;
    return (
      <TipDialog
        isShow={isShow}
        caption="系统提示"
        cancelAction={() => {
          this._closeDialog();
        }}
        footerButton={
          <Button
            type="weak"
            onClick={() => {
              this._closeDialog();
            }}
          >
            关闭
          </Button>
        }
      >
        <React.Fragment>
          <Alert type="info">{message}</Alert>
          <Text>请联系平台管理员进行相关资源的调整</Text>
        </React.Fragment>
      </TipDialog>
    );
  }

  private _closeDialog() {
    changeForbiddentConfig(Init_Forbiddent_Config);
  }
}

/** ========================== 展示没有权限弹窗 ================================ */

/** ============================== start 容器服务 模块 ================================= */
Entry.register({
  businessKey: 'tkestack',

  routes: {
    /**
     * @url https://{{domain}}/tkestack
     */
    index: {
      title: Title,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Overview />
        </Wrapper>
      )
    },
    /**
     * @url https://{{domain}}/tkestack/overview
     */
    overview: {
      title: `${t('概览')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Overview />
        </Wrapper>
      )
    },
    /**
     * @url https://{{domain}}/tkestack/cluster
     */
    cluster: {
      title: `${t('集群管理')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Cluster />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/project
     */
    project: {
      title: `${t('业务管理')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Project />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/addon
     */
    addon: {
      title: `${t('扩展组件')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Addon />
        </Wrapper>
      )
    },

    /**
     * @url https://dev.console.tke.com/tke/regitry
     */
    registry: {
      title: `${t('组织资源')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Registry />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/access
     */
    uam: {
      title: `${t('访问管理')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Uam />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/alarm
     */
    alarm: {
      title: `${t('告警设置')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <AlarmPolicy />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/notify
     */
    notify: {
      title: `${t('通知设置')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Notify />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/alarm-record
     */
    'alarm-record': {
      title: `${t('告警记录')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <AlarmRecord />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/application
     */
    application: {
      title: `${t('Helm 应用')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Application />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/helm
     */
    helm: {
      title: `${t('Helm2 应用')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Helm />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domian}}/tkestack/log
     */
    log: {
      title: `${t('日志采集')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <LogStash />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/persistent-event
     */
    'persistent-event': {
      title: `${t('事件持久化')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <PersistentEvent />
        </Wrapper>
      )
    },
    /**
     * @url https://{{domain}}/tkestack/middleware
     */
    middleware: {
      title: `${t('中间件列表')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <MiddlewareAppContainer platform="tkeStack" />
        </Wrapper>
      )
    },
    /**
     * @url https://{{domain}}/tkestack/audit
     */
    audit: {
      title: `${t('审计记录')} - ${Title}`,
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Audit />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/blank
     */
    blank: {
      title: t('暂无权限'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager} sideBar={false}>
          <ForbiddentDialog />
          <BlankPage />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/vnc
     */
    vnc: {
      title: t('VNC'),
      container: <VNCPage />
    }
  }
});
/** ============================== end 容器服务 模块 ================================= */
