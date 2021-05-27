import './i18n';
import ReactDOM from 'react-dom';
import * as React from 'react';
import { Entry, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Wrapper, PlatformTypeEnum } from './Wrapper';
import { Init_Forbiddent_Config } from '@helper/reduceNetwork';
import { TipDialog } from '@src/modules/common';
import { Button, Alert, Text } from '@tencent/tea-component';
import { PersistentEvent } from '@src/modules/persistentEvent';
// 公有云的图表组件为异步加载，这里为了减少路径配置，还是保留为同步加载，预先import即可变成不split
import '@tencent/tchart/build/ChartsComponents';

const ApplicationPromise = import(/* webpackPrefetch: true */ './src/modules/cluster/index.project');
const Application = React.lazy(() => ApplicationPromise);

const RegistryPromise = import(/* webpackPrefetch: true */ './src/modules/registry');
const Registry = React.lazy(() => RegistryPromise);

const HelmPromise = import(/* webpackPrefetch: true */ './src/modules/helm');
const Helm = React.lazy(() => HelmPromise);

const LogStashPromise = import(/* webpackPrefetch: true */ './src/modules/logStash');
const LogStash = React.lazy(() => LogStashPromise);

const AlarmPolicyPromise = import(/* webpackPrefetch: true */ './src/modules/alarmPolicy');
const AlarmPolicy = React.lazy(() => AlarmPolicyPromise);

const NotifyPromise = import(/* webpackPrefetch: true */ './src/modules/notify');
const Notify = React.lazy(() => NotifyPromise);

const AppPromise = import(/* webpackPrefetch: true */ './src/modules/application');
const App = React.lazy(() => AppPromise);

const ProjectPromise = import(/* webpackPrefetch: true */ './src/modules/project');
const Project = React.lazy(() => ProjectPromise);

insertCSS(
  'myTagSearchBox',
  `.myTagSearchBox{ width:100% !important; background-color: #fff; }
   .myTagSearchBox .tea-search__inner{ width:100% !important; }
   .myTagSearchBox .tea-tag-group{ float:left; }
   .myTagSearchBox .tea-search__tips{ float:left; }
   .myTagSearchBox.is-active .tea-search__tips{ display:none; }
`
);

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

Entry.register({
  businessKey: 'tkestack-project',

  routes: {
    /**
     * @url https://{{domain}}//tkestack-project/application
     */
    index: {
      title: '应用管理 - TKEStack业务侧',
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <Application />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}//tkestack-project/application
     */
    application: {
      title: '应用管理 - TKEStack业务侧',
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <Application />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack-project/project
     */
    project: {
      title: t('业务管理 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <Project />
        </Wrapper>
      )
    },

    /**
     * @urlhttps://{{domain}}//tkestack-project/regitry
     */
    registry: {
      title: t('组织资源 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <Registry />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/alarm
     */
    alarm: {
      title: t('告警设置 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <AlarmPolicy />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/notify
     */
    notify: {
      title: t('通知设置 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <Notify />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/helm-application
     */
    app: {
      title: t('Helm 应用 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <App />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/helm
     */
    helm: {
      title: t('Helm2 应用 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <Helm />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domian}}/tkestack/log
     */
    log: {
      title: t('日志采集 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <LogStash />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/persistent-event
     */
    'persistent-event': {
      title: t('事件持久化 - TKEStack业务侧'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Business}>
          <ForbiddentDialog />
          <PersistentEvent />
        </Wrapper>
      )
    }
  }
});
