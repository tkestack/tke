import { FormPanel } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import * as React from 'react';
import { RootProps } from './ProjectApp';

export class EditProjectNamePanel extends React.Component<RootProps, {}> {
  render() {
    let { projectEdition, actions } = this.props;

    return (
      <FormPanel isNeedCard={false}>
        <FormPanel.Item
          label={t('业务名称')}
          validator={projectEdition.v_displayName}
          errorTipsStyle="Icon"
          message={t('业务名称不能超过63个字符')}
          input={{
            value: projectEdition.displayName,
            onChange: actions.project.inputProjectName,
            placeholder: '请输入业务名称',
            onBlur: e => {
              actions.project.validateDisplayName(e.target.value);
            }
          }}
        />
      </FormPanel>
    );
  }
}
