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

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Card, ContentView, Justify, TableColumn, Text } from '@tencent/tea-component';

import { Clip, GridTable, TipDialog, WorkflowDialog } from '../../../common/components';
import { DialogBodyLayout } from '../../../common/layouts';
import { Chart } from '../../models';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

export class ChartGroupTablePanel extends React.Component<RootProps, any> {
  componentDidMount() {
    this.props.actions.charts.fetch();
  }

  render() {
    return (
      <ContentView>
        <ContentView.Header>
          <Justify left={<h2>{t('Chart包命名空间')}</h2>} />;
        </ContentView.Header>
        <ContentView.Body>
          {
            /// #if tke
            <div className="tc-action-grid">
              <Justify
                left={
                  <React.Fragment>
                    <Button
                      type="primary"
                      onClick={() => {
                        let urlParams = router.resolve(this.props.route);
                        router.navigate(Object.assign({}, urlParams, { sub: 'chart', mode: 'create' }), {});
                      }}
                    >
                      {t('新建')}
                    </Button>
                  </React.Fragment>
                }
              />
            </div>
            /// #endif
          }
          {this._renderTablePanel()}
          {this._renderDeleteChartDialog()}
        </ContentView.Body>
      </ContentView>
    );
  }

  private _renderTablePanel() {
    const columns: TableColumn<Chart>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (x: Chart) => (
          <Text parent="div" overflow>
            <a
              title={x.spec.name}
              href="javascript:;"
              onClick={() => {
                let urlParams = router.resolve(this.props.route);
                router.navigate(Object.assign({}, urlParams, { sub: 'chart', mode: 'detail', tab: 'charts' }), {
                  cg: x.metadata.name,
                  cgName: x.spec.name
                });
              }}
              className="tea-text-overflow"
            >
              {x.spec.name}
            </a>
          </Text>
        )
      },
      {
        key: 'desc',
        header: t('描述'),
        render: (x: Chart) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.displayName}</span>
          </Text>
        )
      },
      {
        key: 'visibility',
        header: t('权限类型'),
        render: (x: Chart) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.visibility === 'Public' ? t('公有') : t('私有')}</span>
          </Text>
        )
      },
      {
        key: 'chartCount',
        header: t('Chart包数'),
        render: (x: Chart) => (
          <Text parent="div" overflow>
            <span className="text">{x.status.versions ? x.status.versions.length : 0}</span>
          </Text>
        )
      },
      {
        key: 'settings',
        header: '操作',
        width: 100,
        render: chart => (
          <React.Fragment>
            <Button
              type="link"
              onClick={() => {
                this.props.actions.charts.deleteChart.start([chart]);
              }}
            >
              {t('删除')}
            </Button>
          </React.Fragment>
        )
      }
    ];

    return (
      <GridTable
        columns={columns}
        emptyTips={<div className="text-center">{t('Chart包命名空间列表为空')}</div>}
        listModel={{
          list: this.props.chart.list,
          query: this.props.chart.query
        }}
        actionOptions={this.props.actions.charts}
      />
    );
  }

  private _renderDeleteChartDialog() {
    const { actions, deleteChart } = this.props;
    return (
      <WorkflowDialog
        caption={t('删除Chart包命名空间')}
        workflow={deleteChart}
        action={actions.charts.deleteChart}
        targets={deleteChart.targets}
        postAction={() => {}}
        params={{}}
      >
        <DialogBodyLayout>
          <p className="til tea-text-overflow">
            <strong className="tip-top">
              {t('确定要删除Chart包命名空间：{{chartName}} 么？', {
                chartName: deleteChart.targets ? deleteChart.targets[0].spec.name : ''
              })}
            </strong>
          </p>
          <p className="text-danger">{t('删除该Chart包命名空间后，相关数据将永久删除，请谨慎操作。')}</p>
        </DialogBodyLayout>
      </WorkflowDialog>
    );
  }
}
