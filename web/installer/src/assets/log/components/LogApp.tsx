
import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { RootState } from '../models/';
import { MainBodyLayout } from '../../common/layouts/';
import { LogHeadPanel } from './LogHeadPanel';
import { LogActionPanel } from "./LogActionPanel";
import { LogTablePanel } from './LogTablePanel';
import { logActions } from '../actions/logActions';
export interface RootProps extends RootState {
  actions?: typeof logActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: logActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class LogApp extends React.Component<RootProps, void> {
  render() {
    return (
      <div className="manage-area">
        <LogHeadPanel {...this.props} />
        <MainBodyLayout>
          <LogActionPanel {...this.props} />
          <LogTablePanel {...this.props} />
        </MainBodyLayout>
      </div>
    );
  }
}
