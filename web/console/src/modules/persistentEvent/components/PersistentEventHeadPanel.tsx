import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';

import { allActions } from '../actions';
import { RootProps } from './PersistentEventApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class PersistentEventHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions, region, route } = this.props;

    // 这里对从创建界面返回之后，判断当前的状态
    let isNeedFetchRegion = region.list.data.recordCount ? false : true;
    isNeedFetchRegion && actions.region.applyFilter({});
    !isNeedFetchRegion && actions.cluster.applyFilter({ regionId: +route.queries['rid'] });
  }

  render() {
    return (
      <React.Fragment>
        <Justify
          left={
            <React.Fragment>
              <h2 style={{ float: 'left' }}>{t('事件持久化')}</h2>
            </React.Fragment>
          }
        />
      </React.Fragment>
    );
  }
}
