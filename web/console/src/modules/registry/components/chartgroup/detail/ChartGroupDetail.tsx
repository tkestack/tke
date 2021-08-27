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
import { allActions } from '../../../actions';
import { RootProps } from '../ChartGroupApp';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { ContentView } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ChartGroupDetail extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.chartGroup.detail.updateChartGroupWorkflow.reset();
    actions.chartGroup.detail.clearEditorState();
    actions.chartGroup.detail.clearValidatorState();

    actions.user.associate.clearUserAssociation();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 查询具体仓库，从而Detail可以用到 */
    actions.chartGroup.detail.fetchChartGroup({ name: route.queries['cg'], projectID: route.queries['prj'] });

    /** 获取具备权限的业务列表 */
    actions.project.list.fetch();
    /** 拉取用户信息 */
    actions.user.detail.fetchUserInfo();
    /** 拉取用户列表 */
    actions.user.associate.userList.performSearch('');
  }

  render() {
    return (
      <React.Fragment>
        <ContentView>
          <ContentView.Header>
            <HeaderPanel />
          </ContentView.Header>
          <ContentView.Body>
            <BaseInfoPanel />
          </ContentView.Body>
        </ContentView>
      </React.Fragment>
    );
  }
}
