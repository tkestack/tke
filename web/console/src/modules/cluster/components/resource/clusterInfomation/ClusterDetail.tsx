/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
import { RootProps } from '../../ClusterApp';
import { KubectlDialog } from '../../KubectlDialog';
import { ClusterDetailBasicInfoPanel } from './ClusterDetailBasicInfoPanel';
import { UpdateClusterAllocationRatioDialog } from './UpdateClusterAllocationRatioDialog';
import { ClusterPlugInfoPanel } from './ClusterPlugInfoPanel';

const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });
@connect(state => state, mapDispatchToProps)
export class ClusterDetailPanel extends React.Component<RootProps, {}> {
  render() {
    return (
      <React.Fragment>
        <ClusterDetailBasicInfoPanel {...this.props} />
        <ClusterPlugInfoPanel {...this.props} />
        <UpdateClusterAllocationRatioDialog />
        <KubectlDialog {...this.props} />
      </React.Fragment>
    );
  }
}
