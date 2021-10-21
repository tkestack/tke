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
import { Justify } from '@tencent/tea-component';

import { allActions } from '../actions';
import { router } from '../router';
import { RootProps } from './LogStashApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class LogStashSubHeadPanel extends React.Component<RootProps, any> {
  componentDidMount() {
    let { actions, route, regionList, projectList } = this.props,
      { clusterId, stashName, namespace } = route.queries;
    let urlParams = router.resolve(route);
    let mode = urlParams['mode'],
      isCreate = mode === 'create',
      isUpdate = mode === 'update',
      isDetail = mode === 'detail';

    //刷不刷新页面在datail模式下都是要log的信息
    if (isDetail) {
      actions.log.fetchSpecificLog(stashName, clusterId, namespace, mode);
    }
    let byProject = window.location.href.includes('tkestack-project');
    if (byProject && projectList.length === 0) {
      actions.cluster.initProjectList();
    } else if (!byProject && regionList.data.recordCount === 0) {
      // 进行地域的拉取
      //如果用户刷新了页面，则需要重新获取
      actions.region.fetch();
    } else {
      //非detail模式下才需要获取namespace信息
      //拉取log详细信息
      if (isUpdate) {
        actions.log.fetchSpecificLog(stashName, clusterId, namespace, 'update');
      }
      //创建页面下帮助用户自动选择namesapce
      if (isCreate) {
        actions.namespace.autoSelectNamespaceForCreate();
      }
    }
  }

  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let { mode } = urlParams;
    let title = '';
    switch (mode) {
      case 'create':
        title = t('新建日志采集规则');
        break;
      case 'update':
        title = t('编辑日志采集规则');
        break;
      case 'detail':
        title = route.queries['stashName'];
        break;
      default:
        title = '';
    }

    return (
      <Justify
        left={
          <React.Fragment>
            <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
              <i className="btn-back-icon" />
              {t('返回')}
            </a>
            <span className="line-icon">|</span>
            <h2 className="tea-h2">{title}</h2>
          </React.Fragment>
        }
      />
    );
  }

  /** 回退按钮 */
  private goBack() {
    let { route } = this.props;
    let newRouteQueies = JSON.parse(
      JSON.stringify(Object.assign({}, route.queries, { stashName: undefined, namespace: undefined }))
    );
    router.navigate({}, newRouteQueies);
  }
}
