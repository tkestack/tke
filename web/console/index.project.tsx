import { app, i18n } from '@tencent/tea-app';
import { translation } from '@i18n/translation';
import ReactDOM from 'react-dom';
// 国际化工具的初始化
i18n.init({ translation });
import * as React from 'react';
import { Entry, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Application } from './src/modules/cluster/index.project';
import { Wrapper, PlatformTypeEnum } from './Wrapper';
import { Registry } from '@src/modules/registry';
import { Init_Forbiddent_Config } from '@helper/reduceNetwork';
import { TipDialog } from '@src/modules/common';
import { Button, Alert, Text } from '@tencent/tea-component';
import { Helm } from '@src/modules/helm';
import { LogStash } from '@src/modules/logStash';
import { PersistentEvent } from '@src/modules/persistentEvent';
import { AlarmPolicy } from '@src/modules/alarmPolicy';
import { Notify } from '@src/modules/notify';
import { Application as App } from '@src/modules/application';
// 公有云的图表组件为异步加载，这里为了减少路径配置，还是保留为同步加载，预先import即可变成不split
import '@tencent/tchart/build/ChartsComponents';
import { Project } from '@src/modules/project';
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
    let { isShow, message } = this.state.forbiddentConfig;
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
