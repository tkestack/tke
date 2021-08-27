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
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Justify } from '@tencent/tea-component';

import { allActions } from '../actions';
import { router } from '../router';
import { RootProps } from './AddonApp';

/** 面包屑 扩展组件的类型展示 */
const addonTypeNameMap = {
  helm: 'Helm'
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class AddonDetailHeadPanel extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props;

    // 回到集群列表即 ''，回到集群列表
    let breadHeader = this._reduceBreadCrumbs();

    return (
      <Justify
        left={
          <React.Fragment>
            <Icon
              style={{ cursor: 'pointer' }}
              type="btnback"
              className="tea-mr-1n"
              onClick={() => {
                this._handleClickForTurnBack();
              }}
            />
            {breadHeader}
          </React.Fragment>
        }
      />
    );
  }

  /** 生成面包屑导航 /组件(地域)/集群ID(集群名)/... */
  private _reduceBreadCrumbs() {
    let { route, cluster } = this.props,
      urlParams = router.resolve(route);

    let { type } = urlParams;
    let { rid, clusterId, resourceIns } = route.queries;

    let firstBreadName = `组件`;
    let secondBreadName = `${clusterId}(${cluster.selection ? cluster.selection.spec.displayName : '-'})`;
    let thirdBreadName = `${addonTypeNameMap[type] ? addonTypeNameMap[type] : type}:${resourceIns}`; // 如果不存在映射，则展示原本的
    let breads: any[] = [firstBreadName, secondBreadName, thirdBreadName];

    let content: React.ReactNode;

    content = (
      <ol className="breadcrumb">
        {breads.map((bread, index) => {
          return (
            <li key={index}>
              {index !== 0 ? (
                <span>{bread}</span>
              ) : (
                <a
                  href="javascript:;"
                  onClick={e => {
                    index === 0 && this._handleClickForTurnBack();
                  }}
                >
                  {bread}
                </a>
              )}
            </li>
          );
        })}
      </ol>
    );

    return content;
  }

  /** 返回处理 */
  private _handleClickForTurnBack() {
    let { route } = this.props;
    let { rid, clusterId } = route.queries;
    let newRouteQuery = { rid, clusterId };
    router.navigate({}, newRouteQuery);
  }
}
