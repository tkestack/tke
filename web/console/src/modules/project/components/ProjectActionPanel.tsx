import * as React from 'react';
import { Button, SearchBox, Justify } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../actions';
import { RootProps } from './ProjectApp';
import { WorkflowDialog } from '../../common/components';
import { CreateProjectPanel } from './CreateProjectPanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { projectActions } from '../actions/projectActions';
import { EditProjectManagerPanel } from './EditProjectManagerPanel';
import { Manager } from '../models';
import { router } from '../router';

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
