import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Input } from '@tencent/tea-component';

import { allActions } from '../actions';
import { RootProps } from './AddonApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditPersistentEventPanel extends React.Component<RootProps, any> {
  render() {
    let { editAddon, route, actions } = this.props,
      { peEdit } = editAddon,
      { esAddress, indexName, v_esAddress, v_indexName } = peEdit;

    let { rid } = route.queries;

    return (
      <React.Fragment>
        <FormPanel.Item validator={v_esAddress} label={t('Elasticsearch地址')} errorTipsStyle="Bubble">
          <Input
            value={esAddress}
            onChange={value => {
              actions.editAddon.pe.inputEsAddress(value);
            }}
            placeholder="eg: http://190.0.0.1:9200"
            onBlur={actions.validator.validateEsAddress}
          />
        </FormPanel.Item>

        <FormPanel.Item
          validator={v_indexName}
          label={t('索引')}
          errorTipsStyle="Bubble"
          message={t('最长60个字符，只能包含小写字母、数字及分隔符("-"、"_"、"+")，且必须以小写字母开头')}
        >
          <Input
            value={indexName}
            onChange={value => {
              actions.editAddon.pe.inputIndexName(value);
            }}
            placeholder="eg: fluentd"
            onBlur={actions.validator.validateIndexName}
          />
        </FormPanel.Item>
      </React.Fragment>
    );
  }
}
