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
import { TablePanel as CTablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { App } from '../../../models';
import { RootProps } from '../AppContainer';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { TipInfo } from '../../../../common/components/';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class TablePanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, appList, route } = this.props;

    const columns: TableColumn<App>[] = [
      {
        key: 'name',
        header: t('应用名'),
        render: (x: App) => (
          <Text parent="div" overflow>
            <a
              href="javascript:;"
              onClick={e => {
                let q = Object.assign({}, route.queries, { app: x.metadata.name, appName: x.spec.name });
                router.navigate({ sub: 'app', mode: 'detail' }, q);
              }}
            >
              {x.spec.name || '-'}
            </a>
            &nbsp;
            {(x.status['phase'] !== 'Succeeded' ||
              !x.status['releaseStatus'] ||
              x.status['releaseStatus'] !== 'deployed' ||
              (x.status['observedGeneration'] && x.status['observedGeneration'] < x.metadata['generation'])) && (
              <Icon
                type="loading"
                title={x.status['phase'] + (x.status['message'] ? '(' + x.status['message'] + ')' : '')}
              />
            )}
          </Text>
        )
      },
      {
        key: 'releaseStatus',
        header: t('状态'),
        render: (x: App) => (
          <Text parent="div" className={x.status.releaseStatus === 'deployed' ? 'text-success' : 'text-danger'}>
            {x.status.releaseStatus || '-'}
          </Text>
        )
      },
      {
        key: 'revision',
        header: t('版本号'),
        render: (x: App) => <Text parent="div">{x.status.revision || '-'}</Text>
      },
      {
        key: 'type',
        header: t('类型'),
        render: (x: App) => <Text parent="div">{x.spec.type || '-'}</Text>
      },
      {
        key: 'chartgroup',
        header: t('Chart仓库'),
        render: (x: App) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.chart ? x.spec.chart.chartGroupName : ''}</span>
          </Text>
        )
      },
      {
        key: 'chart',
        header: t('Chart'),
        render: (x: App) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.chart ? x.spec.chart.chartName : ''}</span>
          </Text>
        )
      },
      {
        key: 'chartversion',
        header: t('Chart版本'),
        render: (x: App) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.chart ? x.spec.chart.chartVersion : ''}</span>
          </Text>
        )
      },
      {
        key: 'releaseLastUpdated',
        header: t('更新时间'),
        render: (x: App) => (
          <Text parent="div">
            {x.status.releaseLastUpdated
              ? dateFormat(new Date(x.status.releaseLastUpdated), 'yyyy-MM-dd hh:mm:ss')
              : '-'}
          </Text>
        )
      },
      { key: 'operation', header: t('操作'), render: app => this._renderOperationCell(app) }
    ];

    return (
      <React.Fragment>
        <CTablePanel
          recordKey={record => {
            return record.metadata.name;
          }}
          columns={columns}
          model={appList}
          action={actions.app.list}
          rowDisabled={record => record.status['phase'] === 'Terminating'}
          emptyTips={emptyTips}
          isNeedPagination={true}
          bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
        />
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (app: App) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton onClick={() => this._removeApp(app)}>{t('删除')}</LinkButton>
      </React.Fragment>
    );
  };

  _removeApp = async (app: App) => {
    let { actions } = this.props;
    const yes = await Modal.confirm({
      message: t('确定删除应用：') + `${app.spec.name}？`,
      description: <p className="text-danger">{t('删除该应用后，相关数据将永久删除，请谨慎操作。')}</p>,
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.app.list.removeAppWorkflow.start([app]);
      actions.app.list.removeAppWorkflow.perform();
    }
  };
}
