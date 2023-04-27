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

import { Bubble, Button, SearchBox, Table } from '@tea/component';
import { Justify } from '@tea/component/justify';
import { bindActionCreators } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { PermissionProvider } from '@src/modules/common/components';
import { dateFormatter, downloadCsv } from '../../../../../helpers';
import { Cluster } from '../../../common';
import { allActions } from '../../actions';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';

console.log('clusteraction----->', t('新建独立集群'));

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ClusterActionPanel extends React.Component<RootProps, any> {
  downloadHandle(clusters: Cluster[]) {
    const rows = [];
    const head = ['ID', t('集群状态'), t('K8S版本'), t('集群类型'), t('创建时间')];

    clusters.forEach((item: Cluster) => {
      const row = [
        item.metadata.name,
        item.status.phase,
        item.status.version,
        item.spec.type,
        dateFormatter(new Date(item.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')
      ];
      rows.push(row);
    });

    downloadCsv(rows, head, `tke-clusterList${Date.now()}.csv`);
  }

  render() {
    const { actions, cluster, route } = this.props;

    const bubbleContent = null;

    return (
      <Table.ActionPanel>
        <Justify
          left={
            <PermissionProvider value="platform.cluster.create_import_button">
              <Bubble placement="right" content={bubbleContent}>
                <Button
                  type="primary"
                  onClick={() => router.navigate({ sub: 'create' }, { rid: route.queries['rid'] })}
                >
                  {t('导入集群')}
                </Button>
                <Button
                  type="primary"
                  onClick={() => router.navigate({ sub: 'createIC' }, { rid: route.queries['rid'] })}
                >
                  {t('新建独立集群')}
                </Button>
              </Bubble>
            </PermissionProvider>
          }
          right={
            <React.Fragment>
              <SearchBox
                value={cluster.query.keyword || ''}
                onChange={actions.cluster.changeKeyword}
                onSearch={actions.cluster.performSearch}
                onClear={() => {
                  actions.cluster.performSearch('');
                }}
                placeholder={t('请输入集群ID')}
              />
              <Button
                icon="download"
                title={t('导出全部')}
                onClick={() => this.downloadHandle(cluster.list.data.records)}
              />
            </React.Fragment>
          }
        />
      </Table.ActionPanel>
    );
  }
}
