import * as React from 'react';
import { RootProps } from './PersistentEventApp';
import { OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../actions';
import { connect } from 'react-redux';
import { WorkflowDialog } from '../../common/components';
import { CreateResource } from '../models';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class PersistentEventDeleteDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, peSelection, resourceInfo, deletePeFlow } = this.props;

    let resourceIns: string = peSelection[0] ? peSelection[0].metadata.name : '';
    let clusterId = route.queries['clusterId'];

    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo,
      mode: 'delete',
      clusterId,
      resourceIns
    };

    return (
      <WorkflowDialog
        caption={t('删除资源')}
        workflow={deletePeFlow}
        action={actions.workflow.deletePeFlow}
        params={route.queries['rid']}
        targets={[resource]}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong>
                  {t('您确定要删除当前集群 {{clusterId}} 的 {{headTitle}} 资源吗？', {
                    clusterId,
                    headTitle: resourceInfo.headTitle
                  })}
                </strong>
              </p>
              <div className="block-help-text text-danger">{t('该资源下所有Pods将一并销毁，请提前备份好数据。')}</div>
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
