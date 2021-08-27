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

import { Bubble, Button, Table } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component/lib/justify';

import { allActions } from '../../../actions';
import { ClusterHelmStatus } from '../../../constants/Config';
import { router } from '../../../router';
import { RootProps } from '../../HelmApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HelmActionPanel extends React.Component<RootProps, {}> {
  render() {
    let {
      actions,
      listState: { helmList, helmQuery, region, clusterHelmStatus },
      route
    } = this.props;
    return (
      <Table.ActionPanel>
        <Justify
          left={
            <Bubble
              placement="right"
              content={clusterHelmStatus.code !== ClusterHelmStatus.RUNNING ? '请先开通Helm应用' : null}
            >
              <Button
                type="primary"
                disabled={clusterHelmStatus.code !== ClusterHelmStatus.RUNNING}
                onClick={() => router.navigate({ sub: 'create' }, route.queries)}
              >
                {t('新建')}
              </Button>
            </Bubble>
          }
          right={<React.Fragment />}
        />
      </Table.ActionPanel>
    );
  }
}
