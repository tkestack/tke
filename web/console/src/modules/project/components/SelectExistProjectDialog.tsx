import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanelTransferTable, FormPanelTransferTableTableProps } from '@tencent/ff-component';
import { bindActionCreators, FetchState, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Alert, Button, Text, Modal } from '@tencent/tea-component';

import { allActions } from '../actions';
import { RootProps } from './ProjectApp';
import { Project } from '../models';
import { projectStatus } from '../constants/Config';
import { getWorkflowError } from '@src/modules/common';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class SelectExistProjectDialog extends React.Component<RootProps, {}> {
  render() {
    const { project, actions, addExistMultiProject, projectDetail, route } = this.props;

    // 参数配置
    const selectorProps: FormPanelTransferTableTableProps<Project> = {
      tagSearch: {
        minWidth: 300,
        attributes: [{ type: 'input', key: 'projectId', name: '业务名称' }]
      },
      /** 要供选择的数据 */
      model: project,

      action: actions.project,

      isNeedScollLoding: false,

      rowDisabled: item => {
        return (
          !!item.spec.parentProjectName ||
          item.metadata.name === projectDetail.metadata.name ||
          item.status.phase !== 'Active'
        );
      },

      /** 选择器标题 */
      title: t('当前有以下业务'),

      columns: [
        {
          key: 'name',
          width: '70%',
          header: t('ID/名称'),
          render: x => (
            <div>
              <Text parent="div" overflow>
                {`${x.metadata.name}(${x.spec.displayName ? x.spec.displayName : '未命名'})`}
              </Text>
            </div>
          )
        },
        {
          key: 'phase',
          width: '30%',
          header: t('状态'),
          render: x => (
            <React.Fragment>
              <Text parent="div" overflow theme={projectStatus[x.status.phase]}>
                {x.status.phase}
              </Text>
            </React.Fragment>
          )
        }
      ],
      recordKey: 'id'
    };
    const cancel = () => {
      actions.project.clearSelection();

      if (addExistMultiProject.operationState === OperationState.Done) {
        actions.project.addExistMultiProject.reset();
      }
      if (addExistMultiProject.operationState === OperationState.Started) {
        actions.project.addExistMultiProject.cancel();
      }
    };

    let failed =
      addExistMultiProject.operationState === OperationState.Done && !isSuccessWorkflow(addExistMultiProject);
    return (
      <Modal
        size={700}
        visible={addExistMultiProject.operationState !== OperationState.Pending}
        caption={t('添加已有业务')}
        onClose={() => cancel()}
      >
        <FormPanelTransferTable<Project> {...selectorProps} />;
        <React.Fragment>
          <Button
            type="primary"
            style={{ margin: '0px 5px 0px 40px' }}
            onClick={() => {
              if (project.selections.length !== 0) {
                actions.project.addExistMultiProject.start(
                  project.selections,
                  projectDetail ? projectDetail.metadata.name : route.queries['projectId']
                );
                actions.project.addExistMultiProject.perform();
              }
            }}
          >
            {failed ? t('重试') : t('完成')}
          </Button>
          <Button
            type="weak"
            onClick={() => {
              cancel();
            }}
          >
            {t('取消')}
          </Button>
          {failed ? (
            <Alert
              type="error"
              style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
            >
              {getWorkflowError(addExistMultiProject)}
            </Alert>
          ) : (
            <noscript />
          )}
        </React.Fragment>
      </Modal>
    );
  }
}
