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
import classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';
import { CreateResource } from 'src/modules/cluster/models';

import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Text } from '@tencent/tea-component';

import { resourceConfig } from '../../../../config';
import { dateFormatter } from '../../../../helpers';
import { Clip, LinkButton } from '../../common/components';
import { allActions } from '../actions';
import { collectorStatus, logModeMap } from '../constants/Config';
import { Log } from '../models';
import { router } from '../router';
import { isCanCreateLogStash } from './LogStashActionPanel';
import { RootProps } from './LogStashApp';
import { HOST_LOG_INPUT_PATH_PREFIX } from '../constants/Config';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class LogStashTablePanel extends React.Component<RootProps, any> {
  render() {
    return <React.Fragment>{this._renderTablePanel()}</React.Fragment>;
  }

  /** 展示Table的内容 */
  private _renderTablePanel() {
    let { actions, logQuery, logList, route, clusterSelection, isOpenLogStash, isDaemonsetNormal } = this.props,
      urlParams = router.resolve(route);

    const columns: TablePanelColumnProps<Log>[] = [
      {
        key: 'name',
        header: t('名称'),
        width: '15%',
        render: x => (
          <React.Fragment>
            <Text overflow>
              <a
                id={`logStashId`}
                title={x.metadata.name}
                href="javascript:;"
                onClick={() => {
                  router.navigate(
                    Object.assign({}, urlParams, { mode: 'detail' }),
                    Object.assign({}, route.queries, {
                      stashName: x.metadata.name,
                      clusterId: route.queries['clusterId'],
                      rid: route.queries['rid'],
                      namespace: x.metadata.namespace
                    })
                  );
                }}
              >
                {x.metadata.name}
              </a>
            </Text>
            <Clip target={`#logStashId`} className="hover-icon" />
          </React.Fragment>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '15%',
        render: x => this._getDaemonsetStatus()
      },
      {
        key: 'logType',
        header: t('类型'),
        width: '15%',
        render: x => {
          let type = x?.spec?.input?.type;
          if (x?.spec?.input?.host_log_input?.path?.includes(HOST_LOG_INPUT_PATH_PREFIX)) {
            type = 'pod-log';
          }

          return <Text overflow>{logModeMap?.[type] ?? '-'}</Text>;
        }
      },
      {
        key: 'Namespace',
        header: t('命名空间'),
        width: '15%',
        render: x => <Text overflow>{x.metadata.namespace}</Text>
      },
      {
        key: 'createdAt',
        header: t('创建时间'),
        width: '15%',
        render: x => <Text>{dateFormatter(new Date(x.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}</Text>
      }
    ];

    const { canCreate, tip } = isCanCreateLogStash(
      clusterSelection[0],
      logList.data.records,
      isDaemonsetNormal,
      isOpenLogStash
    );

    const emptyTips: JSX.Element = (
      <React.Fragment>
        <Trans>
          <Text verticalAlign="middle">{t('您选择的该集群的日志采集规则列表为空，您可以')}</Text>
          <Bubble content={!canCreate ? tip : null}>
            <LinkButton disabled={!canCreate} onClick={this._handleCreate.bind(this)}>
              {t('新建一个日志采集规则')}
            </LinkButton>
          </Bubble>
        </Trans>
      </React.Fragment>
    );

    return (
      <TablePanel
        columns={columns}
        isNeedPagination={false}
        action={actions.log}
        model={{
          list: logList,
          query: logQuery
        }}
        emptyTips={emptyTips}
        getOperations={x => this._renderOperationCell(x)}
        operationsWidth={300}
      />
    );
  }

  /** 处理日志采集的操作 */
  private _handleCreate() {
    let { actions, isOpenLogStash, route, clusterSelection } = this.props,
      urlParams = router.resolve(route);

    if ((clusterSelection && clusterSelection[0] && clusterSelection[0].spec.logAgentName) || isOpenLogStash) {
      router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
    } else {
      actions.workflow.authorizeOpenLog.start();
    }
  }

  /** 获取当前的Daemonset的状态 */
  private _getDaemonsetStatus() {
    const { isDaemonsetNormal } = this.props;
    let content: JSX.Element;

    content = (
      <div>
        <span
          className={classnames(
            'text-overflow',
            collectorStatus[isDaemonsetNormal.phase] && collectorStatus[isDaemonsetNormal.phase].classname
          )}
        >
          {collectorStatus[isDaemonsetNormal.phase] ? collectorStatus[isDaemonsetNormal.phase].text : '-'}
        </span>
        <Bubble content={isDaemonsetNormal.reason ? isDaemonsetNormal.reason : null} />
      </div>
    );

    return content;
  }

  /** 操作按钮 */
  private _renderOperationCell(logStash: Log) {
    let { actions, route, clusterVersion, clusterSelection } = this.props,
      urlParams = router.resolve(route);
    const logAgentName = (clusterSelection && clusterSelection[0] && clusterSelection[0].spec.logAgentName) || '';

    // 编辑日志采集器规则的按钮
    const renderEditButton = () => {
      return (
        <LinkButton
          key={logStash.id + 'update'}
          tipDirection={'right'}
          onClick={() => {
            router.navigate(
              Object.assign({}, urlParams, { mode: 'update' }),
              Object.assign({}, route.queries, {
                stashName: logStash.metadata.name,
                clusterId: route.queries['clusterId'],
                rid: route.queries['rid'],
                namespace: logStash.metadata.namespace
              })
            );
          }}
        >
          {t('编辑收集规则')}
        </LinkButton>
      );
    };

    const renderDeleteButton = () => {
      return (
        <LinkButton
          key={logStash.id + 'delete'}
          onClick={() => {
            const resource: CreateResource = {
              id: uuid(),
              namespace: logStash.metadata.namespace,
              clusterId: route.queries['clusterId'],
              logAgentName,
              resourceIns: logStash.metadata.name,
              resourceInfo: resourceConfig(clusterVersion)['logcs']
            };
            actions.workflow.inlineDeleteLog.start([resource], +route.queries['rid']);
          }}
        >
          {t('删除')}
        </LinkButton>
      );
    };

    const btns = [];
    btns.push(renderEditButton());
    btns.push(renderDeleteButton());

    return btns;
  }
}
