import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../config';
import { Cluster, ResourceInfo, WorkflowDialog } from '../../../common';
import { allActions } from '../../actions';
import { CreateResource } from '../../models';
import { RootProps } from '../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ClusterDeleteDialog extends React.Component<RootProps, any> {
  render() {
    let { deleteClusterFlow, actions, region, cluster } = this.props;

    let target = deleteClusterFlow.targets && deleteClusterFlow.targets[0] ? deleteClusterFlow.targets[0] : null;

    // 需要提交的数据
    let clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    let resourceIns = target && (target as Cluster).metadata ? (target as Cluster).metadata.name : '';

    let resource: CreateResource = {
      id: uuid(),
      resourceInfo: clusterInfo,
      resourceIns
    };

    return (
      <WorkflowDialog
        caption={t('删除集群')}
        workflow={deleteClusterFlow}
        action={actions.workflow.deleteCluster}
        params={region.selection ? region.selection.value : ''}
        targets={[resource]}
        isDisabledConfirm={resourceIns ? false : true}
      >
        <React.Fragment>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <p style={{ wordWrap: 'break-word' }}>
              <strong>
                {t('您确定要删除集群：{{resourceIns}} 吗？', {
                  resourceIns
                })}
              </strong>
            </p>
          </div>
        </React.Fragment>
      </WorkflowDialog>
    );
  }
}
