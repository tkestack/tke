import { ResourceTable } from './ResourceTable';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TablePanelColumnProps } from '@tencent/ff-component';
import { Resource } from '../../../common';
import { router } from '../../router';
const TypeMapping = {
  smtp: t('邮件'),
  text: t('邮件'),
  tencentCloudSMS: t('短信'),
  wechat: t('微信公众号'),
  webhook: t('webhook'),
};

export class ResourceTableChannel extends ResourceTable {
  getColumns(): TablePanelColumnProps<Resource>[] {
    const resourceName = getThisResourceName.call(this);
    function getThisResourceName() {
      const urlParams = router.resolve(this.props.route);
      return urlParams['resourceName'] || '';
    }

    return [
      {
        key: 'type',
        header: t('类型'),
        render: resource => {
          let channel = getChannel.call(this);
          function getChannel() {
            switch (resourceName) {
              case 'channel':
                const channel = resource;
                return channel;
              case 'template':
                const template = resource;
                const channelName = template.metadata.namespace;
                const channelList = this.props.channel.list.data.records;
                function getChannelRecord(name) {
                  return channelList.find(channel => channel.metadata.name === name);
                }
                return getChannelRecord(channelName);
            }
          }

          if (!channel || !channel.spec) {
            return '-';
          }

          function getTypeDesc(channelSpec) {
            for (const type in TypeMapping) {
              if (channelSpec.hasOwnProperty(type)) {
                return TypeMapping[type];
              }
            }
            return '-';
          }

          return getTypeDesc(channel.spec);
        }
      }
    ];
  }
}
