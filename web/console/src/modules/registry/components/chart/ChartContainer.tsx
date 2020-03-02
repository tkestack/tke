import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { router } from '../../router';
import { RootProps } from '../RegistryApp';
import { CreateChartPanel } from './CreateChartPanel';
import { ChartDetailPanel } from './ChartDetailPanel';
import { ChartGroupTablePanel } from './ChartGroupTablePanel';

export class ChartContainer extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);

    if (urlParam['sub'] === 'chart') {
      if (urlParam['mode'] === 'list') {
        return <ChartGroupTablePanel {...this.props} />;
      } else if (urlParam['mode'] === 'create') {
        return <CreateChartPanel {...this.props} />;
      } else if (urlParam['mode'] === 'detail' && urlParam['tab'] === 'charts') {
        return <ChartDetailPanel {...this.props} />;
      } else {
        return <ChartGroupTablePanel {...this.props} />;
      }
    } else {
      return <ChartGroupTablePanel {...this.props} />;
    }
  }
}
