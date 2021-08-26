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

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Justify, Table, TableColumn, Text } from '@tencent/tea-component';
import { expandable } from '@tea/component/table/addons/expandable';

import { dateFormatter } from '../../../../../helpers';
import { Clip, GridTable, TipDialog, WorkflowDialog } from '../../../common/components';
import { DialogBodyLayout } from '../../../common/layouts';
import { ChartIns } from '../../models';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

export class ChartTablePanel extends React.Component<RootProps, any> {
  state = {
    showUsageGuideline: false
  };

  componentDidMount() {
    this.props.actions.chartIns.applyFilter({
      chartgroup: this.props.route.queries['cg']
    });
  }

  render() {
    return (
      <React.Fragment>
        <div className="tc-action-grid">
          <Justify
            left={
              <React.Fragment>
                <Button
                  type="primary"
                  onClick={() => {
                    this.setState({ showUsageGuideline: true });
                  }}
                >
                  {t('Chart上传指引')}
                </Button>
              </React.Fragment>
            }
          ></Justify>
        </div>
        {this._renderTablePanel()}
        {this._renderUsageGuideDialog()}
      </React.Fragment>
    );
  }

  private _renderTablePanel() {
    const columns: TableColumn<ChartIns>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.name}</span>
          </Text>
        )
      },
      {
        key: 'desc',
        header: t('描述'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.displayName}</span>
          </Text>
        )
      },
      {
        key: 'visibility',
        header: t('权限类型'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.visibility === 'Public' ? t('公有') : t('私有')}</span>
          </Text>
        )
      },
      {
        key: 'pullCount',
        header: t('下载次数'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.status.pullCount}</span>
          </Text>
        )
      },
      {
        key: 'name',
        header: t('地址'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <Text className="text" overflow style={{ width: '80%' }} id={`_${x.spec.name}`}>
              {`http://${this.props.dockerRegistryUrl.data}/chart/${this.props.route.queries['cgName']}/${x.spec.name}.tgz`}
            </Text>{' '}
            <Clip target={`#_${x.spec.name}`} />
          </Text>
        )
      }
    ];

    return (
      <GridTable
        columns={columns}
        emptyTips={<div className="text-center">{t('Chart列表为空')}</div>}
        listModel={{
          list: this.props.chartIns.list,
          query: Object.assign({}, this.props.chartIns.query, {
            filter: {
              chartgroup: this.props.route.queries['cg']
            }
          })
        }}
        actionOptions={this.props.actions.chartIns}
      />
    );
  }

  private _renderUsageGuideDialog() {
    return (
      <TipDialog
        isShow={this.state.showUsageGuideline}
        width={680}
        caption={t('Chart 上传指引')}
        cancelAction={() => this.setState({ showUsageGuideline: false })}
        performAction={() => this.setState({ showUsageGuideline: false })}
      >
        <div className="mirroring-box" style={{ marginTop: '0px' }}>
          <ul className="mirroring-upload-list">
            <li>
              <p>
                <strong>
                  <Trans>前置条件</Trans>
                </strong>
              </p>
            </li>
            <li>
              <p>
                <Trans>
                  本地安装 Helm 客户端, 更多可查看{' '}
                  <a href="https://helm.sh/docs/intro/quickstart/" target="_blank">
                    安装 Helm
                  </a>
                  .{' '}
                </Trans>
              </p>
              <code>
                <Clip target="#installHelm" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="installHelm">{`$ curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>本地 Helm 客户端添加 TKEStack 的 repo.</Trans>
              </p>
              <code>
                <Clip target="#addTkeRepo" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="addTkeRepo">{`helm repo add ${this.props.route.queries['cgName']} http://${this.props.dockerRegistryUrl.data}/chart/${this.props.route.queries['cgName']} --username tkestack --password [访问凭证] `}</p>
              </code>
              <p className="text-weak">
                <Trans>
                  获取有效访问凭证信息，请前往
                  <a
                    href="javascript:;"
                    onClick={() => {
                      let urlParams = router.resolve(this.props.route);
                      router.navigate(Object.assign({}, urlParams, { sub: 'apikey', mode: '', tab: '' }), {});
                    }}
                  >
                    [访问凭证]
                  </a>
                  管理。
                </Trans>
              </p>
            </li>
            <li>
              <p>
                <Trans>安装 helm-push 插件</Trans>
              </p>
              <code>
                <Clip target="#installHelmPush" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="installHelmPush">{`$ helm plugin install https://github.com/chartmuseum/helm-push`}</p>
              </code>
            </li>
            <li>
              <p>
                <strong>
                  <Trans>上传Helm Chart</Trans>
                </strong>
              </p>
            </li>
            <li>
              <p>
                <Trans>上传文件夹</Trans>
              </p>
              <code>
                <Clip target="#pushHelmDir" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="pushHelmDir">{`$ helm push ./myapp ${this.props.route.queries['cgName']}`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>上传压缩包</Trans>
              </p>
              <code>
                <Clip target="#pushHelmTar" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="pushHelmTar">{`$ helm push myapp-1.0.1.tgz ${this.props.route.queries['cgName']}`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>下载最新版本</Trans>
              </p>
              <code>
                <Clip target="#downloadChart" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="downloadChart">{`$ helm fetch ${this.props.route.queries['cgName']}/myapp`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>下载指定版本</Trans>
              </p>
              <code>
                <Clip target="#downloadSChart" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="downloadSChart">{`$ helm fetch ${this.props.route.queries['cgName']}/myapp --version 1.0.1`}</p>
              </code>
            </li>
          </ul>
        </div>
      </TipDialog>
    );
  }
}
