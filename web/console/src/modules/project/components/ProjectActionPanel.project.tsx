import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Justify, SearchBox } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { WorkflowDialog } from '../../common/components';
import { allActions } from '../actions';
import { projectActions } from '../actions/projectActions';
import { Manager } from '../models';
import { router } from '../router';
import { CreateProjectPanel } from './CreateProjectPanel';
import { EditProjectManagerPanel } from './EditProjectManagerPanel';
import { RootProps } from './ProjectApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ProjectActionPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    actions.project.poll({});
    actions.project.projectUserInfo.applyFilter({});
  }
  componentWillUnmount() {
    const { actions } = this.props;
    actions.project.clearPolling();
    actions.project.performSearch('');
  }
  render() {
    let { actions, project } = this.props;

    return (
      <div className="tc-action-grid">
        <Justify
          right={
            <SearchBox
              value={project.query.keyword || ''}
              onChange={actions.project.changeKeyword}
              onSearch={actions.project.performSearch}
              placeholder={t('请输入业务名称')}
            />
          }
        />
      </div>
    );
  }
}
