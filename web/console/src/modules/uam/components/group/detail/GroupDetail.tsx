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
import { RootProps } from '../GroupPanel';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { ContentView } from '@tea/component';

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect((state) => state, mapDispatchToProps)
export class GroupDetail extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.group.detail.updateGroupWorkflow.reset();
    actions.group.detail.clearEditorState();
    actions.group.detail.clearValidatorState();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 查询具体用户组，从而Detail可以用到 */
    actions.group.detail.fetchGroup({ name: route.queries['groupName'] });
  }

  render() {
    return <BaseInfoPanel />;
    // return (
    //   <React.Fragment>
    //     <ContentView>
    //       <ContentView.Header>
    //         <HeaderPanel />
    //       </ContentView.Header>
    //       <ContentView.Body>
    //         <BaseInfoPanel />
    //       </ContentView.Body>
    //     </ContentView>
    //   </React.Fragment>
    // );
  }
}
