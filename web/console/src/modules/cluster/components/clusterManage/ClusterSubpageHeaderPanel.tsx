import * as React from 'react';
import { RootProps } from '../ClusterApp';
import { router } from '../../router';
import { Justify, Icon } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export class ClusterSubpageHeaderPanel extends React.Component<RootProps, {}> {
  goBack() {
    history.back();
  }

  render() {
    let title = t('导入集群');

    return (
      <Justify
        left={
          <React.Fragment>
            <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
              <Icon type="btnback" />
              {t('返回')}
            </a>
            <span className="line-icon">|</span>
            <h2>{title}</h2>
          </React.Fragment>
        }
      />
    );
  }
}
