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
import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import { resourceLimitTypeToText, resourceTypeToUnit } from '@src/modules/project/constants/Config';
import { Bubble, Text } from '@tea/component';
import { bindActionCreators, FetchState } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../../../helpers';
import { ListItem } from '../../../../common/components';
import { DetailLayout } from '../../../../common/layouts';
import { allActions } from '../../../actions';
import { ResourceStatus } from '../../../constants/Config';
import { RootProps } from '../../ClusterApp';

const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceNamespaceDetailPanel extends React.Component<RootProps, {}> {
  render() {
    let { subRoot } = this.props,
      { resourceOption, resourceDetailState } = subRoot,
      { resourceDetailInfo } = resourceDetailState;

    let resourceIns = resourceDetailInfo.selection;

    let statusMap = ResourceStatus['np'];

    let isNeedLoading =
      resourceDetailInfo.list.fetched !== true ||
      resourceDetailInfo.list.fetchState === FetchState.Fetching ||
      resourceIns === null;

    return isNeedLoading ? (
      loadingElement
    ) : (
      <div>
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('基本信息')}</h3>
            </div>
            <div className="param-bd">
              <ul className="item-descr-list">
                <ListItem label={t('名称')}>{resourceIns.metadata.name}</ListItem>

                <ListItem label={t('状态')}>
                  <p
                    className={classnames(
                      '',
                      statusMap[resourceIns.status.phase] && statusMap[resourceIns.status.phase].classname
                    )}
                  >
                    {(statusMap[resourceIns.status.phase] && statusMap[resourceIns.status.phase].text) || '-'}
                  </p>
                </ListItem>

                <ListItem label={t('描述')}>
                  {resourceIns.metadata.annotations ? resourceIns.metadata.annotations.description : '-'}
                </ListItem>

                <ListItem label={t('创建时间')}>
                  {dateFormatter(new Date(resourceIns.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
                </ListItem>
                {/* <ListItem label={t('资源限制')}>{this._reduceResourceLimit(resourceIns.spec.hard)}</ListItem>
                <ListItem label={t('已使用资源')}>{this._reduceResourceLimit(resourceIns.status.used)}</ListItem> */}
              </ul>
            </div>
          </div>
        </DetailLayout>
      </div>
    );
  }
  private _reduceResourceLimit(showData) {
    let resourceLimitKeys = showData ? Object.keys(showData) : [];
    let content = resourceLimitKeys.map((item, index) => (
      <Text parent="p" key={index}>{`${resourceLimitTypeToText[item]}:${
        resourceTypeToUnit[item] === 'MiB'
          ? valueLabels1024(showData[item], K8SUNIT.Mi)
          : valueLabels1000(showData[item], K8SUNIT.unit)
      }${resourceTypeToUnit[item]}`}</Text>
    ));
    return (
      <Bubble placement="left" content={content}>
        {content.filter((item, index) => index < 3)}
      </Bubble>
    );
  }
}
