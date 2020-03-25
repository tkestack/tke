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
    actions.manager.applyFilter({});
    actions.manager.fetchAdminstratorInfo();
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
          left={
            <React.Fragment>
              <Button
                type="primary"
                onClick={() => {
                  router.navigate({ sub: 'create' });
                }}
              >
                {t('新建业务')}
              </Button>
              <Button
                type="primary"
                onClick={() => {
                  actions.manager.initAdminstrator();
                  actions.manager.modifyAdminstrator.start();
                }}
              >
                {t('设置管理员')}
              </Button>
            </React.Fragment>
          }
          right={
            <SearchBox
              value={project.query.keyword || ''}
              onChange={actions.project.changeKeyword}
              onSearch={actions.project.performSearch}
              placeholder={t('请输入业务名称')}
            />
          }
        />
        {this._renderEditAdminstratorDialog()}
      </div>
    );
  }

  private _renderEditAdminstratorDialog() {
    const { actions, projectEdition, modifyAdminstrator } = this.props;
    return (
      <WorkflowDialog
        caption={t('编辑管理员')}
        workflow={modifyAdminstrator}
        action={actions.manager.modifyAdminstrator}
        targets={[projectEdition]}
        params={{}}
        postAction={() => {
          actions.project.clearEdition();
        }}
        width={700}
      >
        <EditProjectManagerPanel
          {...this.props}
          rowDisabled={(record: Manager) => {
            return record.name === 'admin';
          }}
        />
      </WorkflowDialog>
    );
  }
}
