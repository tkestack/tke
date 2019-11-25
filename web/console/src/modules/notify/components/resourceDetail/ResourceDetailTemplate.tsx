import * as React from 'react';
import { ResourceDetail } from './ResourceDetail';
import { FormPanel, LinkButton } from '../../../common/components';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Text } from '@tencent/tea-component';
import { router } from '../../router';

export class ResourceDetailTempalte extends ResourceDetail {
  componentDidMount() {
    this.props.actions.resource.channel.fetch({});
  }

  renderMore(ins) {
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('渠道')}>
          {this.renderChannel(ins)}
        </FormPanel.Item>
        {['text', 'tencentCloudSMS', 'wechat']
          .filter(key => ins.spec[key])
          .map(key => {
            return Object.keys(ins.spec[key]).map(property => (
              <FormPanel.Item text key={property} label={property}>
                {ins.spec[key][property] || '-'}
              </FormPanel.Item>
            ));
          })}
      </React.Fragment>
    );
  }

  renderChannel(x) {
    let channelId = x.metadata.namespace;
    let channel = this.props.channel.list.data.records.find(channel => channel.metadata.name === channelId);
    let { route } = this.props;
    let urlParams = router.resolve(route);
    return (
      <React.Fragment>
        <Text>
          <LinkButton
            onClick={() => {
              router.navigate(
                { ...urlParams, mode: 'detail', resourceName: 'channel' },
                { ...route.queries, resourceIns: channelId }
              );
            }}
            className="tea-text-overflow"
          >
            {channelId}
          </LinkButton>
        </Text>
        {channel && <Text theme="weak">({channel.spec.displayName})</Text>}
      </React.Fragment>
    );
  }
}
