import { t } from '@tencent/tea-app/lib/i18n';
import { ResourceTableChannel } from './ResourceTableChannel';

export class ResourceTableTemplate extends ResourceTableChannel {
  componentDidMount() {
    this.props.actions.resource.channel.fetch({});
  }

  getColumns() {
    return [
      {
        key: 'channel',
        header: t('渠道'),
        render: x => {
          let channelId = x.metadata.namespace;
          let channel = this.props.channel.list.data.records.find(channel => channel.metadata.name === channelId);
          return this.renderNameColumn(channelId, channel && channel.spec.displayName, 'channel');
        }
      },
      ...super.getColumns()
    ];
  }
}
