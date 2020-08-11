import * as React from 'react';
import { connect } from 'react-redux';

import { Card, Icon, Select, Switch } from '@tea/component';
import { bindActionCreators, FetchState } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { FormItem } from '../../../../common/components';
import { allActions } from '../../../actions';
import { TailList } from '../../../constants/Config';
import { RootProps } from '../../ClusterApp';
import { YamlEditorPanel } from '../YamlEditorPanel';

const workloadTypeList = [
  {
    value: 'deployment',
    label: 'Deployment'
  },
  {
    value: 'statefulset',
    label: 'StatefulSet'
  },
  {
    value: 'daemonset',
    label: 'DaemonSet'
  },
  {
    value: 'job',
    label: 'Job'
  },
  {
    value: 'tapp',
    label: 'TApp'
  }
];

const inlineDisplayStyle: React.CSSProperties = {
  marginRight: '6px',
  maxWidth: '180px'
};

interface ResourceLogPanelState {
  /** 是否需要继续判断是否选择第一个命名空间 */
  isNeedReceive?: boolean;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceLogPanel extends React.Component<RootProps, ResourceLogPanelState> {
  constructor(props) {
    super(props);
    this.state = {
      isNeedReceive: true
    };
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let { namespaceList } = nextProps;

    if (
      namespaceList.fetched === true &&
      namespaceList.data.records &&
      namespaceList.data.records[0] &&
      this.state.isNeedReceive
    ) {
      this.setState({ isNeedReceive: false });
      this._handleSelectForNamespace(namespaceList.data.records[0].name);
    }
  }

  componentWillUnmount() {
    let { actions } = this.props;
    // 停止轮询 和 关闭自动刷新按钮
    this._handleSwitch(false);
    actions.resourceLog.workload.fetch({ noCache: true });
    // 清除pod选项
    this._clearPodOption();
  }

  render() {
    let { logList } = this.props.subRoot.resourceLogOption;

    return (
      <React.Fragment>
        {this._renderLogFilterBar()}
        {this._renderLogContent()}
      </React.Fragment>
    );
  }

  /** 展示日志内容的部分 */
  private _renderLogContent() {
    let { logList } = this.props.subRoot.resourceLogOption;

    let logContent = logList.data.recordCount ? logList.data.records[0] : t('暂无日志');

    return (
      <Card>
        <Card.Body>
          <YamlEditorPanel config={logContent} readOnly={true} isNeedRefreshContent={true} mode="text/x-sh" />
        </Card.Body>
      </Card>
    );
  }

  /** 展示条件筛选的部分 */
  private _renderLogFilterBar() {
    let { subRoot, namespaceList } = this.props,
      { addons } = subRoot,
      {
        workloadType,
        isAutoRenew,
        namespaceSelection,
        workloadList,
        workloadSelection,
        podList,
        containerSelection,
        podSelection,
        tailLines
      } = subRoot.resourceLogOption;

    // 展示workloadType的选择列表
    let finalWorkloadTypeList = workloadTypeList.filter(item => {
      if (item.value !== 'tapp') {
        return true;
      } else {
        return addons['TappController'] !== undefined;
      }
    });

    let workloadTypeOptions = finalWorkloadTypeList.map((w, index) => ({
      value: w.value,
      text: w.label
    }));

    // 展示命名空间的选择列表
    let namespaceOptions = namespaceList.data.records.map(item => ({
      value: item.name,
      text: item.displayName
    }));

    // 展示workloadList的选择列表
    let workloadListOptions = workloadList.data.records.map(w => ({
      value: w.metadata.name,
      text: w.metadata.name
    }));

    // 展示podList的选择列表
    let podListOptions = podList.data.records.map(p => ({
      value: p.metadata.name,
      text: p.metadata.name
    }));

    // 展示container的选择列表
    let finder = podList.data.records.find(p => p.metadata.name === podSelection);
    let containerOptions = finder
      ? finder.spec.containers.map(c => ({
          value: c.name,
          text: c.name
        }))
      : [];

    // 展示拉取数据条数的选择列表
    let tailOptions = TailList.map((tail, index) => {
      let text = tail.label;
      return {
        value: tail.value,
        text: `${t('显示{{ text }}', { text })}`
      };
    });

    /** 加载中的样式 */
    let loadingElement: JSX.Element = (
      <div style={{ display: 'inline-block' }}>
        <i className="n-loading-icon" />
        &nbsp; <span className="text">{t('加载中...')}</span>
      </div>
    );

    // 判断pod选项是否需要loading
    let isPodOptionLoading =
      podList.fetchState === FetchState.Fetching || workloadList.fetchState === FetchState.Fetching;

    return (
      <Card>
        <Card.Body title={t('条件筛选')}>
          <div className="param-box server-update add">
            <div className="param-bd">
              <ul className="form-list fixed-layout">
                <FormItem label={t('工作负载选项')} tips={t('工作负载类型、命名空间、Workload实例')}>
                  {namespaceList.fetchState === FetchState.Fetching || namespaceList.fetched !== true ? (
                    <Icon type="loading" />
                  ) : (
                    <React.Fragment>
                      <Select
                        style={inlineDisplayStyle}
                        value={namespaceSelection}
                        options={namespaceOptions}
                        onChange={value => {
                          this._handleSelectForNamespace(value);
                        }}
                        placeholder={t('请选择命名空间')}
                      />
                      <Select
                        style={inlineDisplayStyle}
                        value={workloadType}
                        options={workloadTypeOptions}
                        onChange={value => {
                          this._handleSelectForWorkloadType(value);
                        }}
                        placeholder={t('请选择工作负载类型')}
                      />
                      {workloadList.fetchState === FetchState.Fetching ? (
                        <Icon type="loading" />
                      ) : (
                        <Select
                          style={inlineDisplayStyle}
                          value={workloadSelection}
                          options={workloadListOptions}
                          onChange={value => {
                            this._handleSelectForWorkload(value);
                          }}
                          placeholder={t('请选择工作负载')}
                        />
                      )}
                    </React.Fragment>
                  )}
                </FormItem>
                <FormItem label={t('Pod选项')} tips={t('Pod实例、Container实例')}>
                  {isPodOptionLoading ? (
                    <Icon type="loading" />
                  ) : (
                    <React.Fragment>
                      <Select
                        style={inlineDisplayStyle}
                        options={podListOptions}
                        value={podSelection}
                        onChange={value => {
                          this._handleSelectForPod(value);
                        }}
                        placeholder={t('请选择Pod实例')}
                      />
                      <Select
                        style={inlineDisplayStyle}
                        value={containerSelection}
                        options={containerOptions}
                        onChange={value => {
                          this._handleSelectForContainer(value);
                        }}
                        placeholder={t('请选择Container')}
                      />
                    </React.Fragment>
                  )}
                </FormItem>
                <FormItem label={t('其他选项')}>
                  <span
                    className="descript-text"
                    style={{ display: 'inline-block', verticalAlign: 'middle', marginRight: '10px', fontSize: '12px' }}
                  >
                    {t('自动刷新')}
                  </span>
                  <Switch value={isAutoRenew} onChange={checked => this._handleSwitch(checked)} />
                  <Select
                    className="tea-ml-3n"
                    style={inlineDisplayStyle}
                    value={tailLines}
                    options={tailOptions}
                    onChange={value => {
                      this._handleSelectForTailLine(value);
                    }}
                  />
                </FormItem>
              </ul>
            </div>
          </div>
        </Card.Body>
      </Card>
    );
  }

  /** workloadType的选择 */
  private _handleSelectForWorkloadType(type: string) {
    let { actions } = this.props;

    // 清除pod选项
    this._clearPodOption();

    actions.resourceLog.workload.selectWorkloadType(type);
    // 切换类型的时候，需要清空原来的日志数据，并且关闭自动刷新
    actions.resourceLog.closeAutoRenewAndClearLog();
  }

  /** workload的选择 */
  private _handleSelectForWorkload(workload: string) {
    let { actions } = this.props;

    // 清空pod选项
    this._clearPodOption();

    // 选择workload
    actions.resourceLog.workload.selectWorkload(workload);
    // 切换类型的时候，需要清空原来的日志数据，并且关闭自动刷新
    actions.resourceLog.closeAutoRenewAndClearLog();
  }

  /** 选择展示条目数据的操作 */
  private _handleSelectForTailLine(tailLine: string) {
    let { actions, subRoot } = this.props,
      { podSelection, containerSelection } = subRoot.resourceLogOption;

    actions.resourceLog.log.selectTailLine(tailLine);
    // 进行数据的拉取
    actions.resourceLog.fetchLogData(podSelection, containerSelection, tailLine);
  }

  /** namespace的选择 */
  private _handleSelectForNamespace(namespace: string) {
    let { actions } = this.props;
    // 清除pod选项
    this._clearPodOption();

    actions.resourceLog.selectNamespace(namespace);
    // 切换命名空间，需要清空原来的日志数据，并且关闭自动刷新
    actions.resourceLog.closeAutoRenewAndClearLog();
  }

  /** 开启、关闭自动刷新 */
  private _handleSwitch(isChecked: boolean) {
    let { actions, subRoot } = this.props,
      { podSelection, containerSelection, tailLines, isAutoRenew } = subRoot.resourceLogOption;

    if (!isChecked) {
      isAutoRenew && actions.resourceLog.log.toggleAutoRenew();
      actions.resourceLog.log.clearPollLog();
    } else {
      // 进行日志的拉取
      actions.resourceLog.fetchLogData(podSelection, containerSelection, tailLines);
    }
  }

  /** 选择pod之后，需要进行一些操作 */
  private _handleSelectForPod(podName: string) {
    let { actions } = this.props;

    actions.resourceLog.pod.selectPod(podName);
    this._handleSelectForContainer('');
  }

  /** 选择container的时候，拉取数据 */
  private _handleSelectForContainer(containerName: string) {
    let { actions } = this.props;

    // 选择container
    actions.resourceLog.pod.selectContainer(containerName);
  }

  /** 清除pod选项 */
  private _clearPodOption() {
    let { actions } = this.props;
    //
    actions.resourceLog.pod.fetch({ noCache: true });
    actions.resourceLog.pod.selectContainer('');
  }
}
