import * as React from 'react';
import { connect } from 'react-redux';
import { Cluster } from 'src/modules/common/models';

import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Text } from '@tencent/tea-component';

import { LinkButton } from '../../common/components';
import { allActions } from '../actions';
import { router } from '../router';
import { RootProps } from './LogStashApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class LogSettingTablePanel extends React.Component<RootProps, any> {
  componentDidMount(): void {
    let { actions } = this.props;
    actions.cluster.fetch();
  }

  render() {
    return <React.Fragment>{this._renderTablePanel()}</React.Fragment>;
  }

  /** 展示Table的内容 */
  private _renderTablePanel() {
    let { actions, clusterList, clusterQuery, route } = this.props,
      urlParams = router.resolve(route);

    const columns: TablePanelColumnProps<Cluster>[] = [
      {
        key: 'clusterId',
        header: t('集群ID/名称'),
        width: '25%',
        render: item => (
          <React.Fragment>
            <Text overflow>
              {item.metadata.name}
            </Text>
          </React.Fragment>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '25%',
        render: item => (item.spec.logAgentName ? <span>运行中</span> : <span>未开启</span>)
      },
      {
        key: 'logType',
        header: t('版本'),
        width: '25%',
        render: item => <Text overflow>{item.spec.version}</Text>
      }
    ];

    return (
      <TablePanel
        columns={columns}
        isNeedPagination={false}
        action={actions.cluster}
        model={{
          list: clusterList,
          query: clusterQuery,
        }}
        getOperations={x => this._renderOperationCell(x)}
        operationsWidth={300}
      />
    );
  }

  /** 处理开关日志采集组件的的操作 */
  private _handleSwitch(cluster: Cluster) {
    let { actions, route } = this.props;
    let enableLogAgent = !cluster.spec.logAgentName;
    if (enableLogAgent) {
      actions.cluster.enableLogAgent(cluster);
    } else {
      actions.cluster.disableLogAgent(cluster);
    }

    actions.cluster.applyFilter({});

    return;
  }

  /** 操作按钮 */
  private _renderOperationCell(cluster: Cluster) {
    let { actions, route } = this.props;

    // 编辑日志采集器规则的按钮
    const renderSwitchButton = () => {
      return (
        <LinkButton
          key={cluster.metadata.name + 'update'}
          tipDirection={'right'}
          onClick={() => {
            this._handleSwitch(cluster);
          }}
        >
          {!cluster.spec.logAgentName ? t('开启') : t('关闭')}
        </LinkButton>
      );
    };

    let btns = [];
    btns.push(renderSwitchButton());

    return btns;
  }
}
