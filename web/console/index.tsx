import { i18n } from '@tea/app';
import { translation } from '@i18n/translation';
// 国际化工具的初始化
i18n.init({ translation });
import * as React from 'react';
import { Entry, insertCSS } from '@tencent/ff-redux';
import { Cluster } from './src/modules/cluster';
import { Project } from './src/modules/project';
import { Registry } from './src/modules/registry';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Wrapper, PlatformTypeEnum } from './Wrapper';
import { Addon } from './src/modules/addon';
import { Uam } from './src/modules/uam';
import { PersistentEvent } from './src/modules/persistentEvent';
import { AlarmPolicy } from './src/modules/alarmPolicy';
import { Notify } from './src/modules/notify';
import { LogStash } from './src/modules/logStash';
import { Helm } from './src/modules/helm';
import { TipDialog } from './src/modules/common';
import { Button, Alert, Text } from '@tencent/tea-component';
import { Init_Forbiddent_Config } from './helpers/reduceNetwork';

// 公有云的图表组件为异步加载，这里为了减少路径配置，还是保留为同步加载，预先import即可变成不split
import '@tencent/tchart/build/ChartsComponents';
import { BlankPage } from './blankPage';

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

/** ============================== start 容器服务 模块 ================================= */
Entry.register({
  businessKey: 'tkestack',

  routes: {
    /**
     * @url https://{{domain}}/tkestack
     */
    index: {
      title: 'TKEStack',
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Cluster />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/cluster
     */
    cluster: {
      title: t('集群管理 - TKEStack'),
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
      title: t('业务管理 - TKEStack'),
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
      title: t('扩展组件 - TKEStack'),
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
      title: t('组织资源 - TKEStack'),
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
      title: t('访问管理 - TKEStack'),
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
      title: t('告警设置 - TKEStack'),
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
      title: t('通知设置 - TKEStack'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <Notify />
        </Wrapper>
      )
    },

    /**
     * @url https://{{domain}}/tkestack/helm
     */
    helm: {
      title: t('Helm 应用 - TKEStack'),
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
      title: t('日志采集 - TKEStack'),
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
      title: t('事件持久化 - TKEStack'),
      container: (
        <Wrapper platformType={PlatformTypeEnum.Manager}>
          <ForbiddentDialog />
          <PersistentEvent />
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
    }
  }
});
/** ============================== end 容器服务 模块 ================================= */
