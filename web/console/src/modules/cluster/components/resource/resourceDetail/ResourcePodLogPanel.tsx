import * as React from 'react';
import { connect } from 'react-redux';

import { Justify, Switch } from '@tea/component';
import { bindActionCreators, FetchState } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { DetailLayout } from '../../../../common/layouts';
import { allActions } from '../../../actions';
import { TailList } from '../../../constants/Config';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { YamlEditorPanel } from '../YamlEditorPanel';

// 加载中的样式
const loadingElement: JSX.Element = (
  <div style={{ verticalAlign: 'middle', display: 'inline-block' }}>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourcePodLogPanel extends React.Component<RootProps, {}> {
  componentWillMount() {
    let { actions, route, subRoot } = this.props,
      urlParams = router.resolve(route),
      { podList } = subRoot.resourceDetailState;

    /**
     * pre: podList的列表为空
     * 拉取pod列表，之所以在这里进行拉取，是因为查看日志的地方，需要有pod列表的信息
     */
    if (podList.data.recordCount === 0 && urlParams['type'] === 'resource' && urlParams['resourceName'] !== 'cronjob') {
      // 进行podList的拉取
      actions.resourceDetail.pod.poll({
        namespace: route.queries['np'],
        regionId: +route.queries['rid'],
        clusterId: route.queries['clusterId'],
        specificName: route.queries['resourceIns']
      });
    }
  }

  componentDidMount() {
    let { subRoot, actions } = this.props,
      { podName, containerName, tailLines } = subRoot.resourceDetailState.logOption;

    if (podName !== '' && containerName !== '') {
      actions.resourceDetail.log.handleFetchData(podName, containerName, tailLines);
    }
  }

  componentWillUnmount() {
    let { actions } = this.props;
    // 停止轮询
    actions.resourceDetail.log.clearPollLog();
  }

  render() {
    return (
      <DetailLayout>
        {this._renderFilterBar()}
        {this._renderLogContent()}
      </DetailLayout>
    );
  }

  /** 展示日志内容 */
  private _renderLogContent() {
    let { logList } = this.props.subRoot.resourceDetailState;

    let logContent = logList.data.recordCount ? logList.data.records[0] : t('暂无日志');

    return <YamlEditorPanel config={logContent} readOnly={true} isNeedRefreshContent={true} mode="text/x-sh" />;
  }

  /** 渲染日志的过滤条件 */
  private _renderFilterBar() {
    let { subRoot, actions } = this.props,
      { podList, logOption } = subRoot.resourceDetailState,
      { podName, containerName, tailLines, isAutoRenew } = logOption;

    // 判断是否需要展示loading态
    let isShowLoading = podList.fetched !== true || podList.fetchState === FetchState.Fetching;

    // 渲染pod选择的列表
    let podOptions = podList.data.recordCount
      ? podList.data.records.map((p, index) => {
          return (
            <option key={index} value={p.metadata.name}>
              {p.metadata.name}
            </option>
          );
        })
      : [];
    podOptions.unshift(
      <option key={-1} value="">
        {podList.data.recordCount ? t('请选择Pod') : t('暂无Pod')}
      </option>
    );

    // 渲染container的列表
    let containerList = podName ? podList.data.records.find(p => p.metadata.name === podName).spec.containers : [];
    let containerOptions = containerList.length
      ? containerList.map((c, index) => {
          return (
            <option key={index} value={c.name}>
              {c.name}
            </option>
          );
        })
      : [];
    containerOptions.unshift(
      <option key={-1} value="">
        {containerList.length ? t('请选择Container') : t('该Pod下暂无Container')}
      </option>
    );

    // 渲染拉取数据条数的
    let tailOptions = TailList.map((tail, index) => {
      let text = tail.label;
      return (
        <option key={index} value={tail.value}>
          {t('显示{{text}}', { text })}
        </option>
      );
    });

    return (
      <div className="tc-action-grid" style={{ marginTop: '0' }}>
        <Justify
          left={
            isShowLoading ? (
              loadingElement
            ) : (
              <React.Fragment>
                <select
                  className="tc-15-select m"
                  value={podName}
                  onChange={e => {
                    actions.resourceDetail.log.selectPod(e.target.value);
                  }}
                >
                  {podOptions}
                </select>
                <select
                  className="tc-15-select m tea-ml-2n"
                  disabled={podName === ''}
                  value={containerName}
                  onChange={e => {
                    actions.resourceDetail.log.selectContainer(e.target.value);
                  }}
                >
                  {containerOptions}
                </select>
                <select
                  className="tc-15-select m tea-ml-2n"
                  disabled={podName === '' || containerName === ''}
                  value={tailLines}
                  onChange={e => {
                    actions.resourceDetail.log.selectTailLine(e.target.value);
                  }}
                >
                  {tailOptions}
                </select>
              </React.Fragment>
            )
          }
          right={
            <React.Fragment>
              <span
                className="descript-text"
                style={{ display: 'inline-block', verticalAlign: 'middle', marginRight: '10px', fontSize: '12px' }}
              >
                {t('自动刷新')}
              </span>
              <Switch value={isAutoRenew} onChange={checked => this._handleSwitch(checked)} className="mr20" />
            </React.Fragment>
          }
        />
      </div>
    );
  }

  /** 开启、关闭自动刷新 */
  private _handleSwitch(isChecked: boolean) {
    let { actions, subRoot } = this.props,
      { podName, containerName, tailLines } = subRoot.resourceDetailState.logOption;

    if (!isChecked) {
      actions.resourceDetail.log.toggleAutoRenew();
      actions.resourceDetail.log.clearPollLog();
    } else {
      // 进行日志的拉取
      actions.resourceDetail.log.handleFetchData(podName, containerName, tailLines);
    }
  }
}
