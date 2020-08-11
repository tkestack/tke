import { FormPanel } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';
import * as React from 'react';
import { RootProps } from './AlarmPolicyApp';
export class AlarmPolicyHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;

    actions.projectNamespace.initProjectList();
  }

  render() {
    let { actions, projectList, namespaceList, projectSelection, namespaceSelection, cluster } = this.props;

    let projectListOptions = projectList.map((p, index) => ({
      text: p.displayName,
      value: p.name
    }));
    let namespaceOptions = namespaceList.data.records.map((p, index) => ({
      text: p.displayName,
      value: p.name
    }));
    return (
      <Justify
        left={
          <div style={{ lineHeight: '28px' }}>
            <h2 style={{ float: 'left' }}>{t('告警设置')}</h2>
            <FormPanel.InlineText>{t('项目：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={t('业务')}
              options={projectListOptions}
              value={projectSelection}
              onChange={value => {
                actions.projectNamespace.selectProject(value);
              }}
            ></FormPanel.Select>
            <FormPanel.InlineText>{t('namespace：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={'namespace'}
              options={namespaceOptions}
              value={namespaceSelection}
              onChange={value => actions.namespace.selectNamespace(value)}
            />
          </div>
        }
      />
    );
  }
}
