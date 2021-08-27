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
import { Button, ContentView, Justify, Table, TableColumn, Text } from '@tencent/tea-component';

import { dateFormatter } from '../../../../../helpers';
import { GridTable, WorkflowDialog, TipDialog, Clip } from '../../../common/components';
import { DialogBodyLayout } from '../../../common/layouts';
import { ApiKey } from '../../models';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

interface ApiKeyState {
  showUsageGuideline: boolean;
}

export class ApiKeyTablePanel extends React.Component<RootProps, ApiKeyState> {
  state = {
    showUsageGuideline: false
  };

  componentDidMount() {
    this.props.actions.apiKey.fetch();
  }

  render() {
    return (
      <ContentView>
        <ContentView.Header>
          <Justify left={<h2>{t('访问凭证')}</h2>} />;
        </ContentView.Header>
        <ContentView.Body>
          <div className="tc-action-grid">
            <Justify
              left={
                <React.Fragment>
                  <Button
                    type="primary"
                    onClick={() => {
                      let urlParams = router.resolve(this.props.route);
                      router.navigate(Object.assign({}, urlParams, { sub: 'apikey', mode: 'create' }), {});
                    }}
                  >
                    {t('新建')}
                  </Button>
                  <Button
                    type="weak"
                    onClick={() => {
                      this.setState({ showUsageGuideline: true });
                    }}
                  >
                    {t('使用指引')}
                  </Button>
                </React.Fragment>
              }
            ></Justify>
          </div>
          {this._renderTablePanel()}
          {this._renderDeleteApiKeyDialog()}
          {this._renderToggleKeyStatusDialog()}
          {this._renderUsageGuideDialog()}
        </ContentView.Body>
      </ContentView>
    );
  }

  private _renderTablePanel() {
    const columns: TableColumn<ApiKey>[] = [
      {
        key: 'name',
        header: t('凭证'),
        render: x => {
          let _startEight = x.spec.apiKey.substring(0, 8);
          let _endEight = x.spec.apiKey.substring(Math.max(x.spec.apiKey.length - 8, 0), x.spec.apiKey.length);
          return (
            <React.Fragment>
              <p className="tea-text-overflow" style={{ position: 'relative' }}>
                {`${_startEight}...${_endEight}`}{' '}
                <span id={`${_startEight}${_endEight}`} style={{ position: 'absolute', top: -1000 }}>
                  {x.spec.apiKey}
                </span>{' '}
                <Clip target={`#${_startEight}${_endEight}`} />
              </p>
            </React.Fragment>
          );
        }
      },
      {
        key: 'desc',
        header: t('描述'),
        render: x => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.description || '-'}</span>
          </Text>
        )
      },
      {
        key: 'issue_at',
        header: t('创建时间'),
        render: x => (
          <Text parent="div" overflow>
            <span className="text">{dateFormatter(new Date(x.spec.issue_at), 'YYYY-MM-DD HH:mm:ss')}</span>
          </Text>
        )
      },
      {
        key: 'expire_at',
        header: t('过期时间'),
        render: x => (
          <Text parent="div" overflow>
            <span className="text">{dateFormatter(new Date(x.spec.expire_at), 'YYYY-MM-DD HH:mm:ss')}</span>
          </Text>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        render: x => (
          <Text parent="div" overflow>
            <Text theme={x.status.disabled || x.status.expired ? 'danger' : 'success'}>
              {x.status.expired ? t('已过期') : x.status.disabled ? t('已禁用') : t('已启用')}
            </Text>
          </Text>
        )
      },
      {
        key: 'settings',
        header: '操作',
        width: 100,
        render: key => (
          <React.Fragment>
            <Button
              type="link"
              disabled={key.status.expired}
              onClick={() => {
                this.props.actions.apiKey.toggleKeyStatus.start([key]);
              }}
            >
              {key.status.disabled ? t('启用') : t('禁用')}
            </Button>
            <Button
              type="link"
              onClick={() => {
                this.props.actions.apiKey.deleteApiKey.start([key]);
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
        emptyTips={<div className="text-center">{t('访问凭证列表为空')}</div>}
        listModel={{
          list: this.props.apiKey.list,
          query: this.props.apiKey.query
        }}
        actionOptions={this.props.actions.apiKey}
      />
    );
  }

  private _renderDeleteApiKeyDialog() {
    const { actions, deleteApiKey } = this.props;
    return (
      <WorkflowDialog
        caption={t('删除访问凭证')}
        workflow={deleteApiKey}
        action={actions.apiKey.deleteApiKey}
        targets={deleteApiKey.targets}
        postAction={() => {}}
        params={{}}
        confirmMode={
          deleteApiKey.targets
            ? deleteApiKey.targets[0].status.expired
              ? null
              : {
                  label: t('访问凭证'),
                  value: deleteApiKey.targets ? deleteApiKey.targets[0].spec.apiKey : ''
                }
            : null
        }
      >
        <DialogBodyLayout>
          <p className="til">
            <strong className="tip-top">{t('确定要删除该访问凭证么？')}</strong>
          </p>
          <p className="text-danger">{t('删除该访问凭证后，该凭证将永久失效，请谨慎操作。')}</p>
        </DialogBodyLayout>
      </WorkflowDialog>
    );
  }

  private _renderToggleKeyStatusDialog() {
    const { actions, toggleKeyStatus } = this.props;
    return (
      <WorkflowDialog
        caption={t('修改凭证状态')}
        workflow={toggleKeyStatus}
        action={actions.apiKey.toggleKeyStatus}
        targets={toggleKeyStatus.targets}
        postAction={() => {}}
        params={{}}
      >
        <DialogBodyLayout>
          <p className="til">
            <strong className="tip-top">{t('确定要修改该凭证状态么？')}</strong>
          </p>
        </DialogBodyLayout>
      </WorkflowDialog>
    );
  }

  private _renderUsageGuideDialog() {
    return (
      <TipDialog
        isShow={this.state.showUsageGuideline}
        width={680}
        caption={t('使用指引')}
        cancelAction={() => this.setState({ showUsageGuideline: false })}
        performAction={() => this.setState({ showUsageGuideline: false })}
      >
        <div className="mirroring-box" style={{ marginTop: '0px' }}>
          <ul className="mirroring-upload-list">
            <li>
              <p>
                <strong>
                  <Trans>登录</Trans> TKEStack Docker Registry
                </strong>
              </p>
              <code>
                <Clip target="#loginDocker" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="loginDocker">{`sudo docker login -u tkestack -p [访问凭证] ${this.props.dockerRegistryUrl.data}`}</p>
              </code>
              <p className="text-weak">
                <Trans>[访问凭证]需要在有效期内而且是“启用”状态下，否则会无法登录，请重新选择有效凭证或新建。</Trans>
              </p>
            </li>
            <li>
              <p>
                <strong>
                  <Trans>拉取 Registry 中指定镜像</Trans>
                </strong>
              </p>
              <code>
                <Clip target="#fetchRegistry" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="fetchRegistry">{`sudo docker pull ${this.props.dockerRegistryUrl.data}/[命名空间]/nginx:latest`}</p>
              </code>
              <p className="text-weak">
                <Trans>以拉取命名空间下名为 nginx 的镜像仓库内版本为 latest 的容器镜像为例</Trans>
              </p>
            </li>
            <li>
              <p>
                <strong>
                  <Trans>推送本地镜像到 Registry 中</Trans>
                </strong>
              </p>
              <code>
                <Clip target="#pushRegistry2" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="pushRegistry2">{`sudo docker tag nginx:latest ${this.props.dockerRegistryUrl.data}/[命名空间]/nginx:latest`}</p>
              </code>
              <br />
              <code>
                <Clip target="#pushRegistry3" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="pushRegistry3">{`sudo docker push ${this.props.dockerRegistryUrl.data}/[命名空间]/nginx:latest`}</p>
              </code>
              <p className="text-weak">
                <Trans>以推送本地最新版本 nginx 镜像到容器镜像服务内 [命名空间]/nginx 镜像仓库为例</Trans>
              </p>
            </li>
          </ul>
        </div>
      </TipDialog>
    );
  }
}
