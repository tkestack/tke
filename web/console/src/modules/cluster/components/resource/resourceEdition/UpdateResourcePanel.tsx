/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../../../actions';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { UpdateServiceAccessTypePanel } from '../resourceTableOperation/UpdateServiceAccessTypePanel';
import { UpdateWorkloadPodNumPanel } from '../resourceTableOperation/UpdateWorkloadPodNumPanel';
import { UpdateWorkloadRegistryPanel } from '../resourceTableOperation/UpdateWorkloadRegistryPanel';
import { WorkloadUpdatePanel } from '../resourceTableOperation/workloadUpdate';
import { EditLbcfBackGroupPanel } from './EditLbcfBackGroupPanel';
import { SubHeaderPanel } from './SubHeaderPanel';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UpdateResourcePanel extends React.Component<RootProps, {}> {
  render() {
    const { route } = this.props,
      urlParams = router.resolve(route);

    let headTitle = '';

    // 更新资源所需展示的页面
    let content: JSX.Element;

    // 判断当前的资源
    const resourceType = urlParams['resourceName'],
      updateType = urlParams['tab'];

    const { clusterVersion } = this.props;

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
    } else if (
      ['deployment', 'statefulset', 'daemonset', 'cronjob'].includes(resourceType) &&
      ['modifyStrategy', 'modifyNodeAffinity'].includes(updateType)
    ) {
      return <WorkloadUpdatePanel kind={resourceType} updateType={updateType} clusterVersion={clusterVersion} />;
    }

    return (
      <div className="manage-area">
        <SubHeaderPanel headTitle={headTitle} />
        {content}
      </div>
    );
  }
}
