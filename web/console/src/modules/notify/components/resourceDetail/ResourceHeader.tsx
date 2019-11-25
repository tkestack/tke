import * as React from 'react';
import { RootProps } from '../NotifyApp';
import { router } from '../../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { resourceConfig } from '../../../../../config';
import { Icon, Justify } from '@tencent/tea-component';

const rc = resourceConfig();

export class ResourceHeader extends React.Component<RootProps, {}> {
  goBack() {
    let { route } = this.props;
    let urlParams = router.resolve(route);
    router.navigate({ ...urlParams, mode: 'list' }, {});
  }

  render() {
    let { route } = this.props;
    let urlParams = router.resolve(route);
    let title = '';
    let resourceInfo = rc[urlParams.resourceName] || rc.channel;

    switch (urlParams['mode']) {
      case 'create':
        title = t('新建{{headTitle}}', resourceInfo);
        break;
      case 'update':
        title = t('更新{{headTitle}}', resourceInfo);
        break;
      case 'copy':
        title = t('复制{{headTitle}}', resourceInfo);
        break;
      case 'detail':
        title = t('{{headTitle}}详情', resourceInfo);
    }

    return (
      <React.Fragment>
        <Justify
          left={
            <React.Fragment>
              <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
                <Icon type="btnback" />
                {t('返回')}
              </a>
              <h2>{title}</h2>
            </React.Fragment>
          }
        />
      </React.Fragment>
    );
  }
}
