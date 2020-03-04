import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';

import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';
import { KubectlDialog } from '../../KubectlDialog';
import { ClusterDetailBasicInfoPanel } from './ClusterDetailBasicInfoPanel';
import { UpdateClusterAllocationRatioDialog } from './UpdateClusterAllocationRatioDialog';

const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });
@connect(state => state, mapDispatchToProps)
export class ClusterDetailPanel extends React.Component<RootProps, {}> {
  render() {
    return (
      <React.Fragment>
        <ClusterDetailBasicInfoPanel {...this.props} />
        <UpdateClusterAllocationRatioDialog />
        <KubectlDialog {...this.props} />
      </React.Fragment>
    );
  }
}
