import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';
import { WorkflowDialog } from '../../../../common/components';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Text } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourceTappPodDeleteDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot, region } = this.props,
      { resourceDetailState } = subRoot,
      { podSelection, removeTappPodFlow } = resourceDetailState;

    return (
      <WorkflowDialog
        caption={t('删除pod')}
        workflow={removeTappPodFlow}
        action={actions.workflow.removeTappPod}
        params={region.selection ? region.selection.value : ''}
        targets={removeTappPodFlow.targets}
        // preAction={() => {
        //   actions.resourceDetail.pod.podSelect([]);
        // }}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div className="docker-dialog jiqun">
            <p>
              <strong>{t('您确认删除以下pod吗？')}</strong>
            </p>
            <p style={{ maxWidth: '550px' }}>
              {podSelection.map((item, index) => (
                <Text key={index} style={{ marginRight: '10px', wordBreak: 'break-all' }}>
                  {item.metadata.name}
                </Text>
              ))}
            </p>
            <Text theme="danger">{t('pod删除后将不可恢复，请提前备份好数据。')}</Text>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
