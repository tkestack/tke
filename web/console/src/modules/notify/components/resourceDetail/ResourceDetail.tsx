import * as React from 'react';
import { RootProps } from '../NotifyApp';
import { LinkButton, FormPanel } from '../../../common/components';
import { router } from '../../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon } from '@tencent/tea-component';

export class ResourceDetail extends React.Component<RootProps, {}> {
  render() {
    let { actions, route } = this.props;
    let id = route.queries.resourceIns;
    let urlParams = router.resolve(route);
    let resource = this.props[urlParams.resourceName] || this.props.channel;
    let ins = resource.list.data.records.find(ins => ins.metadata.name === id);
    return ins ? this.renderIns(ins) : <Icon type="loading" />;
  }

  renderIns(ins) {
    let { actions, route } = this.props;
    let urlParams = router.resolve(route);
    return (
      <FormPanel
        title={t('基本信息')}
        operation={
          <LinkButton
            onClick={() => {
              router.navigate({ ...urlParams, mode: 'update' }, route.queries);
            }}
          >
            {t('编辑')}
          </LinkButton>
        }
      >
        <FormPanel.Item text label={t('名称')}>
          {ins.spec.displayName}
        </FormPanel.Item>
        {this.renderMore(ins)}
      </FormPanel>
    );
  }

  renderMore(ins) {
    return <React.Fragment />;
  }
}
