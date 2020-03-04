import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../actions';
import { RootProps } from './ApplicationApp.project';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ApplicationHeadPanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, projectList, projectSelection, cluster } = this.props;
    /** 渲染业务列表 */
    let projectListOptions = projectList.map((p, index) => ({
      text: p.displayName,
      value: p.name
    }));
    return (
      <div className="manage-area-title secondary-title">
        <h2 style={{ float: 'left' }}>应用</h2>
        <FormPanel.InlineText>{t('业务：')}</FormPanel.InlineText>
        <FormPanel.Select
          options={projectListOptions}
          value={projectSelection}
          onChange={value => {
            actions.projectNamespace.selectProject(value);
          }}
        ></FormPanel.Select>
      </div>
    );
  }
}
