import * as React from 'react';
import { RootProps } from './AddonApp';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../actions';
import { connect } from 'react-redux';
import { WorkflowDialog, CreateResource, ResourceInfo } from '../../common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ResourceNameMap } from '../constants/Config';
import { resourceConfig } from '../../../../config';
import { Text } from '@tencent/tea-component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class AddonDeleteDialog extends React.Component<RootProps, any> {
  render() {
    let { actions, route, deleteResourceFlow, openAddon, clusterVersion } = this.props;

    let { rid, clusterId } = route.queries;

    // 需要提交的数据
    let resource: CreateResource;

    if (openAddon.selection) {
      let selection = openAddon.selection;
      let resourceName = ResourceNameMap[selection.spec.type]
        ? ResourceNameMap[selection.spec.type]
        : selection.spec.type;
      let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[resourceName];
      resource = {
        id: uuid(),
        resourceInfo,
        clusterId,
        resourceIns: selection.metadata.name
      };
    }

    return (
      <WorkflowDialog
        caption={t('删除扩展组件')}
        workflow={deleteResourceFlow}
        action={actions.workflow.deleteResource}
        params={+rid}
        targets={[resource]}
        isDisabledConfirm={openAddon.selection ? false : true}
      >
        <Text theme="strong" parent="p">
          {t('您确定要删除 {{addonName}} 扩展组件吗？', {
            addonName: openAddon.selection ? openAddon.selection.spec.type : '-'
          })}
        </Text>
      </WorkflowDialog>
    );
  }
}
