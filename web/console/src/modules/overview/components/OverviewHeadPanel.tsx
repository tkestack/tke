import * as React from 'react';

import { t } from '@tencent/tea-app/lib/i18n';
import { Justify, ExternalLink } from '@tea/component';

import { RootProps } from './OverviewApp';

export class OverviewHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    this.props.actions.clusterOverActions.applyFilter({});
  }
  render() {
    let { actions } = this.props;

    return (
      <Justify
        left={<h2>{t('概览')}</h2>}
        right={
          <ExternalLink href={'https://github.com/tkestack/tke/tree/master/docs/guide/zh-CN'}>
            容器服务帮助手册
          </ExternalLink>
        }
      ></Justify>
    );
  }
}
