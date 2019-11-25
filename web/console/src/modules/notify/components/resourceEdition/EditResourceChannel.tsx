import * as React from 'react';
import { t } from '@tencent/tea-app/lib/i18n';
import { Form } from '@tea/component/form';
import { Input } from '@tea/component/input';
import { channelSchema } from '../../schema/channelSchema';
import { onChange } from '../../schema/schemaUtil';
import { Radio } from '@tencent/tea-component';
import { EditResource } from './EditResource';

export class EditResourceChannel extends EditResource {
  getSchema() {
    return channelSchema;
  }

  renderForm() {
    let resource = channelSchema;
    resource = this.state.resource;
    return (
      <Form>
        <Form.Item label={t('名称')} required>
          <Input
            placeholder={t('请填写名称')}
            value={resource.properties.spec.properties.displayName.value}
            onChange={onChange(resource.properties.spec.properties.displayName)}
          />
        </Form.Item>

        <Form.Item label={t('渠道')}>
          <Radio.Group value={resource.properties.spec.pick} onChange={onChange(resource.properties.spec)}>
            <Radio name="smtp">{t('邮件')}</Radio>
            <Radio name="tencentCloudSMS">{t('短信')}</Radio>
            <Radio name="wechat">{t('微信公众号')}</Radio>
          </Radio.Group>
        </Form.Item>
        {this.renderFields(resource.properties.spec.properties[resource.properties.spec['pick']])}
      </Form>
    );
  }
}
