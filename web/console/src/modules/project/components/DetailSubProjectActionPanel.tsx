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
import { PlatformTypeEnum } from '../constants/Config';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class DetailSubProjectActionPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { route, actions } = this.props;

    actions.detail.project.applyPolling({ parentProject: route.queries['projectId'] });
  }

  componentWillUnmount() {
    let { route, actions } = this.props;
    actions.detail.project.clearPolling();
    actions.detail.project.performSearch('');
  }

  render() {
    let { actions, project, projectDetail, route, platformType, userManagedProjects } = this.props;
    let enableOp =
      platformType === PlatformTypeEnum.Manager ||
      (platformType === PlatformTypeEnum.Business &&
        userManagedProjects.list.data.records.find(
          item => item.name === (projectDetail ? projectDetail.metadata.name : null)
        ));
    return (
      <div className="tc-action-grid">
        <Justify
          left={
            enableOp ? (
              <React.Fragment>
                <Button
                  type="primary"
                  onClick={() => {
                    actions.project.inputParentPorject(
                      projectDetail ? projectDetail.metadata.name : route.queries['projectId']
                    );
                    router.navigate({ sub: 'create' });
                  }}
                >
                  {t('新建子业务')}
                </Button>

                <Button
                  type="primary"
                  onClick={() => {
                    actions.project.addExistMultiProject.start([]);
                  }}
                >
                  {t('添加已有业务')}
                </Button>
              </React.Fragment>
            ) : (
              <></>
            )
          }
          right={
            <SearchBox
              value={project.query.keyword || ''}
              onChange={actions.detail.project.changeKeyword}
              onSearch={actions.detail.project.performSearch}
              placeholder={t('请输入业务名称')}
            />
          }
        />
      </div>
    );
  }
}
