import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators, FetchState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Justify, Select, Text } from '@tencent/tea-component';

import { allActions } from '../actions';
import { RootProps } from './AddonApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class AddonHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    // 进行地域的拉取，每次重新进入addon列表页面，都需要拉取最新的地域 和 集群的信息，并且需要重新轮询cluster下的addon列表
    actions.region.applyFilter({});
  }

  componentWillUnmount() {
    const { actions } = this.props;
    // 离开的时候 要清除定时器，不要浪费请求
    actions.cluster.addon.clearPolling();
  }

  render() {
    const { actions, cluster } = this.props;

    let clusterContent: React.ReactNode;

    if (cluster.list.fetched !== true || cluster.list.fetchState === FetchState.Fetching) {
      clusterContent = <Icon type="loading" />;
    } else {
      const clusterOptions = cluster.list.data.records.map(item => ({
        value: item.metadata.name,
        text: `${item.metadata.name}(${item.spec.displayName ? item.spec.displayName : '未命名'})`
      }));

      clusterContent = (
        <Select
          type="simulate"
          appearence="button"
          searchable
          size="m"
          options={clusterOptions}
          value={cluster.selection ? cluster.selection.metadata.name : null}
          onChange={value => {
            const finder = cluster.list.data.records.find(item => item.metadata.name === value);
            actions.cluster.selectCluster(finder);
          }}
          placeholder={cluster.list.data.recordCount ? t('请选择集群') : t('当前地域暂无集群')}
        />
      );
    }

    return (
      <Justify
        left={
          <React.Fragment>
            <h2>{t('扩展组件')}</h2>
            <Text theme="label" className="text tea-mr-2n tea-ml-1n" verticalAlign="middle">
              {t('集群')}
            </Text>
            {clusterContent}
          </React.Fragment>
        }
      />
    );
  }
}
