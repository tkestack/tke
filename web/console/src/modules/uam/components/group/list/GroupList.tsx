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
import { HeaderPanel } from './HeaderPanel';
import { ActionPanel } from './ActionPanel';
import { TablePanel } from './TablePanel';
import { RootProps } from '../GroupPanel';
import { ContentView, Card, Justify, Icon } from '@tea/component';

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect((state) => state, mapDispatchToProps)
export class GroupList extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    const { actions } = this.props;
    /** 取消轮询 */
    actions.group.list.clearPolling();
  }

  componentDidMount() {
    const { actions } = this.props;
    /** 拉取用户组列表 */
    actions.group.list.poll();
  }

  render() {
    return (
      <React.Fragment>
        <ActionPanel />
        <TablePanel />
        {/*<ContentView>*/}
        {/*  <ContentView.Header>*/}
        {/*    <HeaderPanel />*/}
        {/*  </ContentView.Header>*/}
        {/*  <ContentView.Body>*/}
        {/*    <ActionPanel />*/}
        {/*    <TablePanel />*/}
        {/*  </ContentView.Body>*/}
        {/*</ContentView>*/}
      </React.Fragment>
    );
  }
}
