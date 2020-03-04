import * as React from 'react';
import { ResourceDetail } from './ResourceDetail';
import { FormPanel } from '@tencent/ff-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ResourceTableReceiver } from '../resourceList/ResourceTableReceiver';

export class ResourceDetailReceiverGroup extends ResourceDetail {
  componentDidMount() {
    this.props.actions.resource.receiver.fetch({});
  }

  renderMore(ins) {
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('接收人')}>
          <ResourceTableReceiver {...this.props} resourceName={'receiver'} onlyTable bordered />
        </FormPanel.Item>
      </React.Fragment>
    );
  }
}
