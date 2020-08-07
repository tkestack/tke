import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Card, Icon, Select, Switch, TableColumn, Text } from '@tea/component';
import { bindActionCreators, FetchState, insertCSS } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../../../helpers';
import { Clip, FormItem, GridTable } from '../../../../common/components';
import { allActions } from '../../../actions';
import { Event } from '../../../models';
import { RootProps } from '../../ClusterApp';

insertCSS(
  'ResourceEventPanel',
  `
.tc-15-table-fixed-body.tc-15-table-rowhover {
    overflow: inherit;
}
`
);

const workloadTypeList = [
  {
    value: '',
    label: t('全部类型')
  },
  {
    value: 'cronjob',
    label: 'CronJob'
  },
  {
    value: 'daemonset',
    label: 'DaemonSet'
  },
  {
    value: 'deployment',
    label: 'Deployment'
  },
  {
    value: 'ingress',
    label: 'Ingress'
  },
  {
    value: 'job',
    label: 'Job'
  },
  {
    value: 'node',
    label: 'Node'
  },
  {
    value: 'pods',
    label: 'Pod'
  },
  {
    value: 'pv',
    label: 'PersistentVolume'
  },
  {
    value: 'pvc',
    label: 'PersistentVolumeClaim'
  },
  {
    value: 'sc',
    label: 'StorageClass'
  },
  {
    value: 'statefulset',
    label: 'StatefulSet'
  },
  {
    value: 'svc',
    label: 'Service'
  },
  {
    value: 'tapp',
    label: 'TApp'
  }
];
interface ResourceEventPanelState {
  /** 是否需要继续判断是否选择第一个命名空间 */
  isNeedReceive?: boolean;
}
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceEventPanel extends React.Component<RootProps, ResourceEventPanelState> {
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
    let { actions, route } = this.props;
    // 停止轮询
    actions.resourceEvent.fetch({ noCache: true });
    actions.resourceEvent.clearPollEvent();
    actions.resourceEvent.workload.selectWorkload('');
    actions.resourceEvent.workload.selectWorkloadType('');
  }

  render() {
    return (
      <React.Fragment>
        {this._renderEventFilterBar()}
        {this._renderEventTablePanel()}
      </React.Fragment>
    );
  }

  /** 展示条件筛选的部分 */
  private _renderEventFilterBar() {
    let { subRoot, namespaceList } = this.props,
      { namespaceSelection, workloadType, workloadList, workloadSelection, isAutoRenew } = subRoot.resourceEventOption;

    // 展示命名空间的选择列表
    let namespaceOptions = namespaceList.data.records.map(n => {
      return {
        value: n.name,
        text: n.displayName
      };
    });

    // 展示workloadType的选择列表
    let workloadTypeOptions = workloadTypeList.map(w => ({
      value: w.value,
      text: w.label
    }));

    // 展示workloadList的选择列表
    let workloadListOptions = workloadList.data.records.map(w => {
      return {
        value: w.metadata.name,
        text: w.metadata.name
      };
    });

    /** 加载中的样式 */
    let loadingElement: JSX.Element = (
      <div style={{ display: 'inline-block' }}>
        <i className="n-loading-icon" />
        &nbsp; <span className="text">{t('加载中...')}</span>
      </div>
    );

    return (
      <Card>
        <Card.Body title="条件筛选">
          <div className="param-box server-update add">
            <div className="param-bd">
              <ul className="form-list fixed-layout">
                <FormItem label={t('命名空间')}>
                  {namespaceList.fetched !== true || namespaceList.fetchState === FetchState.Fetching ? (
                    <Icon type="loading" />
                  ) : (
                    <Select
                      size="m"
                      options={namespaceOptions}
                      value={namespaceSelection}
                      onChange={value => {
                        this._handleSelectForNamespace(value);
                      }}
                    />
                  )}
                </FormItem>
                <FormItem label={t('类型')}>
                  <Select
                    size="m"
                    options={workloadTypeOptions}
                    value={workloadType}
                    onChange={value => {
                      this._handleSelectForWorkloadType(value);
                    }}
                  />
                </FormItem>
                <FormItem label={t('名称')}>
                  {workloadType === '' ? (
                    <p className="text-label">{t('请先选择类型和Namespace')}</p>
                  ) : workloadList.fetchState === FetchState.Fetching ? (
                    loadingElement
                  ) : (
                    <Select
                      size="m"
                      options={workloadListOptions}
                      value={workloadSelection}
                      onChange={value => {
                        this._handleSelectForWorkload(value);
                      }}
                    />
                  )}
                </FormItem>
                <FormItem label={t('其余选项')}>
                  <span
                    className="descript-text"
                    style={{ display: 'inline-block', verticalAlign: 'middle', marginRight: '10px', fontSize: '12px' }}
                  >
                    {t('自动刷新')}
                  </span>
                  <Switch value={isAutoRenew} onChange={checked => this._handleSwitch(checked)} className="mr20" />
                </FormItem>
              </ul>
            </div>
          </div>
        </Card.Body>
      </Card>
    );
  }

  /** 处理namespace的选择 */
  private _handleSelectForNamespace(namespace: string) {
    let { actions, subRoot } = this.props,
      { workloadType, workloadSelection } = subRoot.resourceEventOption;
    actions.resourceEvent.selectNamespace(namespace);

    // 这里需要去拉取数据
    actions.resourceEvent.fetchEventData(workloadType, namespace, workloadSelection);
  }

  /** 处理workloadType的选择 */
  private _handleSelectForWorkloadType(type: string) {
    let { actions, subRoot } = this.props,
      { namespaceSelection } = subRoot.resourceEventOption;
    actions.resourceEvent.workload.selectWorkloadType(type);

    // 切换类型需要清空原来的workload的选择项
    actions.resourceEvent.workload.selectWorkload('');

    // 拉取相对应的事件的数据
    actions.resourceEvent.fetchEventData(type, namespaceSelection, '');
  }

  /** 处理workload的选择 */
  private _handleSelectForWorkload(workload: string) {
    let { actions, subRoot } = this.props,
      { namespaceSelection, workloadType } = subRoot.resourceEventOption;
    actions.resourceEvent.workload.selectWorkload(workload);

    // 拉取相对应的事件的数据
    actions.resourceEvent.fetchEventData(workloadType, namespaceSelection, workload);
  }

  /** 处理自动刷新按钮 */
  private _handleSwitch(isChecked: boolean) {
    let { actions, subRoot } = this.props,
      { namespaceSelection, workloadType, workloadSelection } = subRoot.resourceEventOption;

    actions.resourceEvent.toggleAutoRenew();
    if (!isChecked) {
      actions.resourceEvent.clearPollEvent();
    } else {
      // 进行日志的拉取
      actions.resourceEvent.fetchEventData(workloadType, namespaceSelection, workloadSelection, true);
    }
  }

  /** 展示事件列表 */
  private _renderEventTablePanel() {
    let { actions, subRoot } = this.props,
      { eventList, eventQuery } = subRoot.resourceEventOption;

    /** 处理时间 */
    const reduceTime = (time: string) => {
      let [first, second] = dateFormatter(new Date(time), 'YYYY-MM-DD HH:mm:ss').split(' ');

      return <Text>{`${first} ${second}`}</Text>;
    };

    const columns: TableColumn<Event>[] = [
      {
        key: 'firstTime',
        header: t('首次出现时间'),
        width: '10%',
        render: x => reduceTime(x.firstTimestamp)
      },
      {
        key: 'lastTime',
        header: t('最后出现时间'),
        width: '10%',
        render: x => reduceTime(x.lastTimestamp)
      },
      {
        key: 'type',
        header: t('级别'),
        width: '8%',
        render: x => (
          <div>
            <p className={classnames('text-overflow', { 'text-danger': x.type === 'Warning' })}>{x.type}</p>
          </div>
        )
      },
      {
        key: 'resourceType',
        header: t('资源类型'),
        width: '8%',
        render: x => (
          <div>
            <p title={x.involvedObject.kind} className="text-overflow">
              {x.involvedObject.kind}
            </p>
          </div>
        )
      },
      {
        key: 'name',
        header: t('资源名称'),
        width: '12%',
        render: x => (
          <div>
            <span id={'eventName' + x.id} title={x.metadata.name} className="text-overflow m-width">
              {x.metadata.name}
            </span>
            <Clip target={'#eventName' + x.id} className="hover-icon" />
          </div>
        )
      },
      {
        key: 'content',
        header: t('内容'),
        width: '12%',
        render: x => (
          <Bubble placement="bottom" content={x.reason || null}>
            <Text parent="div" overflow>
              {x.reason}
            </Text>
          </Bubble>
        )
      },
      {
        key: 'desp',
        header: t('详细描述'),
        width: '15%',
        render: x => (
          <Bubble placement="bottom" content={x.message || null}>
            <Text parent="div" overflow>
              {x.message}
            </Text>
          </Bubble>
        )
      },
      {
        key: 'count',
        header: t('出现次数'),
        width: '6%',
        render: x => (
          <Text parent="div" overflow>
            {x.count}
          </Text>
        )
      }
    ];

    let emptyTips: JSX.Element = <div>{t('事件列表为空')}</div>;

    return (
      <GridTable
        columns={columns}
        emptyTips={emptyTips}
        listModel={{
          list: eventList,
          query: eventQuery
        }}
        actionOptions={actions.resourceEvent}
      />
    );
  }
}
