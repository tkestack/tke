import * as React from 'react';

import { t } from '@tencent/tea-app/lib/i18n';

import { RootProps } from './OverviewApp';

export class OverviewHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    this.props.actions.clusterOverActions.applyFilter({});
  }
  render() {
    let { actions } = this.props;

    return <h2>{t('概览')}</h2>;
  }
}
