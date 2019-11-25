import * as React from 'react';
import { SearchBox, Table } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../actions';
import { RootProps } from './PersistentEventApp';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component/lib/justify';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ClusterActionPanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, cluster } = this.props;

    return (
      <Table.ActionPanel>
        <Justify
          left={<React.Fragment />}
          right={
            <SearchBox
              value={cluster.query.keyword || ''}
              onChange={actions.cluster.changeKeyword}
              onSearch={actions.cluster.performSearch}
              onClear={() => {
                actions.cluster.performSearch('');
              }}
              placeholder={t('请输入集群ID')}
            />
          }
        />
      </Table.ActionPanel>
    );
  }
}
