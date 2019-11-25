import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { SubHeaderPanel } from './SubHeaderPanel';
import { router } from '../../../router';
import { UpdateServiceAccessTypePanel } from '../resourceTableOperation/UpdateServiceAccessTypePanel';
import { UpdateWorkloadRegistryPanel } from '../resourceTableOperation/UpdateWorkloadRegistryPanel';
import { UpdateWorkloadPodNumPanel } from '../resourceTableOperation/UpdateWorkloadPodNumPanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { EditLbcfBackGroupPanel } from './EditLbcfBackGroupPanel';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UpdateResourcePanel extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let headTitle = '';

    // 更新资源所需展示的页面
    let content: JSX.Element;

    // 判断当前的资源
    let resourceType = urlParams['resourceName'],
      updateType = urlParams['tab'];

    if (resourceType === 'svc' && updateType === 'modifyType') {
      content = <UpdateServiceAccessTypePanel />;
      headTitle = t('更新访问方式');
    } else if (
      (resourceType === 'deployment' || resourceType === 'statefulset' || resourceType === 'daemonset') &&
      updateType === 'modifyRegistry'
    ) {
      content = <UpdateWorkloadRegistryPanel />;
      headTitle = t('滚动更新镜像');
    } else if (resourceType === 'tapp' && updateType === 'modifyRegistry') {
      content = <UpdateWorkloadRegistryPanel />;
      headTitle = t('更新镜像');
    } else if ((resourceType === 'deployment' || resourceType === 'tapp') && updateType === 'modifyPod') {
      content = <UpdateWorkloadPodNumPanel />;
      headTitle = t('更新实例数量');
    } else if (resourceType === 'lbcf' && updateType === 'createBG') {
      content = <EditLbcfBackGroupPanel />;
      headTitle = t('配置后端负载');
    } else if (resourceType === 'lbcf' && updateType === 'updateBG') {
      content = <EditLbcfBackGroupPanel />;
      headTitle = t('更新后端负载');
    }

    return (
      <div className="manage-area">
        <SubHeaderPanel headTitle={headTitle} />
        {content}
      </div>
    );
  }
}
