import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ExternalLink, Justify, Text } from '@tencent/tea-component';

import { SelectList } from '../../common/components';
import { cloneDeep } from '../../common/utils';
import { allActions } from '../actions';
import { RootProps } from './LogStashApp';
import { FormPanel } from '@tencent/ff-component/src';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class LogStashHeadPanel extends React.Component<RootProps, any> {
  componentDidMount() {
    const { actions, regionList, isOpenLogStash, route, namespaceSelection } = this.props;

    // let { clusterId } = route.queries;
    let byProject = window.location.href.includes('tkestack-project');
    if (byProject) {
      // 获取业务列表
      actions.cluster.initProjectList();
    } else {
      if (regionList.data.recordCount === 0) {
        // 第一次进入或者刷新了页面 需要进行地域的拉取,拉取之后会进行集群列表的拉取
        actions.region.applyFilter({});
      }
    }

    // 如果是在列表页，要发起获取LogList的操作
    // 此处是因为，在create | update 当中，已经进行region 和 cluster的拉取，需要更新一下logList列表的更新
    if (isOpenLogStash) {
      actions.logDaemonset.fetch();
    }
  }

  render() {
    let { actions, clusterList, clusterSelection, projectList, projectSelection } = this.props;

    const getSelector = (filterType) => {
      if (filterType === 'byProject') {
        let projectListOptions = projectList.map((p, index) => ({
          text: p.displayName,
          value: p.name
        }));

        return (<>
          <FormPanel.InlineText>{t('业务：')}</FormPanel.InlineText>
          <FormPanel.Select
            options={projectListOptions}
            value={projectSelection}
            onChange={value => {
              actions.cluster.selectProject(value);
            }}
          ></FormPanel.Select>
        </>);
      }

      //渲染集群列表selectList选择项
      const selectClusterList = cloneDeep(clusterList);
      selectClusterList.data.records = clusterList.data.records.map(cluster => {
        return { clusterId: cluster.metadata.name, clusterName: cluster.spec.displayName };
      });

      return (<>
        <Text theme="label" className="text tea-mr-2n tea-ml-1n" verticalAlign="middle">
          {t('集群')}
        </Text>
        <SelectList
          mode="select"
          recordData={selectClusterList}
          value={clusterSelection[0] ? clusterSelection[0].metadata.name : ''}
          valueField="clusterId"
          textFields={['clusterId', 'clusterName']}
          textFormat={`\${clusterId} (\${clusterName})`}
          align="start"
          className="tc-15-select m"
          style={{ display: 'inline-block', lineHeight: '29px', verticalAlign: '-2px' }}
          onSelect={actions.cluster.selectCluster}
          onRetry={actions.cluster.fetch}
          isUnshiftDefaultItem={false}
        />
      </>);
    };

    return (
      <Justify
        left={
          <React.Fragment>
            <h2 className="tea-h2">{t('日志采集')}</h2>
            {getSelector(window.location.href.includes('/tkestack-project') ? 'byProject' : 'byCluster')}
          </React.Fragment>
        }
      />
    );
  }
}
