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

import { bindActionCreators, FetchState } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Icon, Justify, Select, Text } from '@tencent/tea-component';

import { allActions } from '../../../actions';
import { router } from '../../../router';
import { RootProps } from '../../HelmApp';

const routerSea = seajs.require('router');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HelmHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const {
        route,
        actions,
        listState: { region, cluster }
      } = this.props,
      urlParams = router.resolve(route);
    // 这里对从resouce列表返回后，判断当前的状态
    let isNeedFetchRegion = region.list.data.recordCount ? false : true;
    isNeedFetchRegion && actions.region.applyFilter({});
    actions.cluster.applyFilter({});

    if (region.selection && cluster.selection && (!route.queries['rid'] || !route.queries['clusterId'])) {
      router.navigate(
        urlParams,
        Object.assign({}, route.queries, { rid: region.selection.value, clusterId: cluster.selection.metadata.name })
      );
    }
  }
  onSelect(clusterId: string) {
    let {
      actions,
      listState: { cluster }
    } = this.props;

    let item = cluster.list.data.records.find(item => item.clusterId === clusterId);
    if (item) {
      actions.cluster.select(item);
    }
  }
  render() {
    let {
      actions,
      listState: { cluster },
      namespaceList,
      namespaceSelection
    } = this.props;
    let clusterContent: React.ReactNode;

    if (cluster.list.fetched !== true || cluster.list.fetchState === FetchState.Fetching) {
      clusterContent = <Icon type="loading" />;
    } else {
      let clusterOptions = cluster.list.data.records.map(item => ({
        value: item.metadata.name,
        text: `${item.metadata.name}(${item.spec.displayName ? item.spec.displayName : '未命名'})`
      }));
      let namespaceOptions = namespaceList.data.records.map((p, index) => ({
        text: p.displayName,
        value: p.name
      }));
      clusterContent = (
        <>
          <Select
            size="m"
            options={clusterOptions}
            value={cluster.selection ? cluster.selection.metadata.name : null}
            onChange={value => {
              let finder = cluster.list.data.records.find(item => item.metadata.name === value);
              actions.cluster.selectCluster(finder);
            }}
            placeholder={cluster.list.data.recordCount ? t('请选择集群') : t('当前地域暂无集群')}
          />
          <Text reset verticalAlign={'middle'}>
            namespace
          </Text>
          <Select
            options={namespaceOptions}
            value={namespaceSelection}
            onChange={value => actions.namespace.selectNamespace(value)}
          />
        </>
      );
    }

    return (
      <Justify
        left={
          <React.Fragment>
            <h2>{t('Helm应用')}</h2>
            <Text theme="label" className="text tea-mr-2n tea-ml-1n" verticalAlign="middle">
              {t('集群')}
            </Text>
            {clusterContent}
          </React.Fragment>
        }
      />
    );
  }
}
