import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';
import { WorkflowDialog } from '../../../../common/components';
import { CreateResource, RsEditJSONYaml } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourceRollbackDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, subRoot, region } = this.props,
      { resourceInfo, resourceDetailState } = subRoot,
      { rollbackResourceFlow, rsSelection } = resourceDetailState;

    let rsVersion = rsSelection[0] ? +rsSelection[0].metadata.annotations['deployment.kubernetes.io/revision'] : 0;

    let jsonData: RsEditJSONYaml = {
      kind: 'DeploymentRollback',
      apiVersion: 'extensions/v1beta1',
      name: route.queries['resourceIns'],
      rollbackTo: {
        revision: rsVersion
      }
    };

    let resourceIns = route.queries['resourceIns'];

    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo,
      namespace: route.queries['np'],
      clusterId: route.queries['clusterId'],
      resourceIns,
      jsonData: JSON.stringify(jsonData)
    };

    return (
      <WorkflowDialog
        caption={t('回滚资源')}
        workflow={rollbackResourceFlow}
        action={actions.workflow.rollbackResource}
        params={region.selection ? region.selection.value : ''}
        targets={[resource]}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong>
                  {t('您确定要回滚{{headTitle}}：{{resourceIns}} 至 版本v{{rsVersion}}吗？', {
                    headTitle: resourceInfo.headTitle,
                    resourceIns,
                    rsVersion
                  })}
                </strong>
              </p>
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
