import * as React from 'react';
import { connect } from 'react-redux';

import { ContentView } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';

import { allActions } from '../../actions';
import { RootState } from '../../models';
import { router } from '../../router';
import { StrategyActionPanel } from './StrategyActionPanel';
import { StrategyDetailsPanel } from './StrategyDetailsPanel';
import { StrategyHeadPanel } from './StrategyHeadPanel';
import { StrategyTablePanel } from './StrategyTablePanel';

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class StrategyApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    actions.strategy.poll();
  }
  render() {
    let { route } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;
    return (
      <React.Fragment>
        {sub ? (
          <ContentView>
            <ContentView.Header>
              <StrategyHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <StrategyDetailsPanel />
            </ContentView.Body>
          </ContentView>
        ) : (
          <ContentView>
            <ContentView.Header>
              <StrategyHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <StrategyActionPanel />
              <StrategyTablePanel />
            </ContentView.Body>
          </ContentView>
        )}
      </React.Fragment>
    );
  }
}
