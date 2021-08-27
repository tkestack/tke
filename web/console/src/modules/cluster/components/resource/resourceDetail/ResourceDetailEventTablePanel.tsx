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

import { Bubble, Pagination, TableColumn, Text } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../../../helpers';
import { Clip, GridTable } from '../../../../common/components';
import { TableLayout } from '../../../../common/layouts';
import { allActions } from '../../../actions';
import { Event } from '../../../models';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceDetailEventTablePanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions } = this.props;
    actions.resourceDetail.event.poll();
  }

  componentWillUnmount() {
    let { actions } = this.props;
    // 清除轮询
    actions.resourceDetail.event.clearPolling();
    // 清除所有的搜索条件
    actions.resourceDetail.event.reset();
    // 清空eventLis的data
    actions.resourceDetail.event.fetch({ noCache: true });
  }

  render() {
    return this._renderTablePanel();
  }

  /** 列表 */
  private _renderTablePanel() {
    let { subRoot, actions } = this.props,
      { resourceDetailState } = subRoot,
      { event } = resourceDetailState;

    /** 处理时间 */
    const reduceTime = (time: string) => {
      let [first, second] = dateFormatter(new Date(time), 'YYYY-MM-DD HH:mm:ss').split(' ');

      return (
        <React.Fragment>
          <Text>{`${first} ${second}`}</Text>
        </React.Fragment>
      );
    };

    const columns: TableColumn<Event>[] = [
      {
        key: 'firstTime',
        header: t('首次出现时间'),
        width: '10%',
        render: x => reduceTime(x.firstTimestamp)
      },
      {
        key: 'lastTime',
        header: t('最后出现时间'),
        width: '10%',
        render: x => reduceTime(x.lastTimestamp)
      },
      {
        key: 'type',
        header: t('级别'),
        width: '8%',
        render: x => (
          <div>
            <p className={classnames('text-overflow', { 'text-danger': x.type === 'Warning' })}>{x.type}</p>
          </div>
        )
      },
      {
        key: 'resourceType',
        header: t('资源类型'),
        width: '8%',
        render: x => (
          <div>
            <p title={x.involvedObject.kind} className="text-overflow">
              {x.involvedObject.kind}
            </p>
          </div>
        )
      },
      {
        key: 'name',
        header: t('资源名称'),
        width: '12%',
        render: x => (
          <div>
            <span id={'eventName' + x.id} title={x.metadata.name} className="text-overflow m-width">
              {x.metadata.name}
            </span>
            <Clip target={'#eventName' + x.id} />
          </div>
        )
      },
      {
        key: 'content',
        header: t('内容'),
        width: '12%',
        render: x => (
          <Bubble placement="bottom" content={x.reason || null}>
            <Text parent="div" overflow>
              {x.reason}
            </Text>
          </Bubble>
        )
      },
      {
        key: 'desp',
        header: t('详细描述'),
        width: '15%',
        render: x => (
          <Bubble placement="bottom" content={x.message || null}>
            <Text parent="div" overflow>
              {x.message}
            </Text>
          </Bubble>
        )
      },
      {
        key: 'count',
        header: t('出现次数'),
        width: '6%',
        render: x => (
          <div>
            <Text parent="div" overflow>
              {x.count}
            </Text>
          </div>
        )
      }
    ];

    return (
      <GridTable
        columns={columns}
        emptyTips={<div>{t('事件列表为空')}</div>}
        listModel={event}
        actionOptions={actions.resourceDetail.event}
      />
    );
  }
}
