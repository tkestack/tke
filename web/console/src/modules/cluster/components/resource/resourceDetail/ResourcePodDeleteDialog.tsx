import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';
import { WorkflowDialog } from '../../../../common/components';
import { CreateResource } from '../../../models';
import { resourceConfig } from '../../../../../../config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourcePodDeleteDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, subRoot, region, clusterVersion } = this.props,
      { resourceDetailState } = subRoot,
      { podSelection, deletePodFlow } = resourceDetailState;

    let podResourceInfo = resourceConfig(clusterVersion)['pods'];
    let deleteResourceIns = podSelection[0] ? podSelection[0].metadata.name : '';
    let namespace = podSelection[0] ? podSelection[0].metadata.namespace : 'default';

    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo: podResourceInfo,
      namespace,
      clusterId: route.queries['clusterId'],
      resourceIns: deleteResourceIns
    };

    return (
      <WorkflowDialog
        caption={t('销毁实例')}
        workflow={deletePodFlow}
        action={actions.workflow.deletePod}
        params={region.selection ? region.selection.value : ''}
        targets={[resource]}
        preAction={() => {
          actions.resourceDetail.pod.podSelect([]);
        }}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong>{t('您确定要销毁实例{{ deleteResourceIns }}吗？', { deleteResourceIns })}</strong>
              </p>
              <div className="block-help-text">{t('实例销毁重建后将不可恢复，请提前备份好数据。')}</div>
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
