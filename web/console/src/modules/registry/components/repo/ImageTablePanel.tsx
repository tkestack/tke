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
import { Image, Tag } from '../../models';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

interface ImageTableState {
  expandedKeys?: string[];
  showUsageGuideline?: boolean;
}

export class ImageTablePanel extends React.Component<RootProps, any> {
  state = {
    expandedKeys: [],
    showUsageGuideline: false
  };

  componentDidMount() {
    this.props.actions.image.applyFilter({
      namespace: this.props.route.queries['ns']
    });
  }

  render() {
    return (
      <React.Fragment>
        <div className="tc-action-grid">
          <Justify
            left={
              <React.Fragment>
                {/* <Button
                  type="primary"
                  onClick={() => {
                    let urlParams = router.resolve(this.props.route);
                    router.navigate(
                      Object.assign({}, urlParams, { sub: 'repo', mode: 'icreate' }),
                      this.props.route.queries
                    );
                  }}
                >
                  {t('新建')}
                </Button> */}
                <Button
                  type="primary"
                  onClick={() => {
                    this.setState({ showUsageGuideline: true });
                  }}
                >
                  {t('镜像上传指引')}
                </Button>
              </React.Fragment>
            }
          ></Justify>
        </div>
        {this._renderTablePanel()}
        {this._renderDeleteImageDialog()}
        {this._renderUsageGuideDialog()}
      </React.Fragment>
    );
  }

  private _renderTablePanel() {
    const columns: TableColumn<Image>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (x: Image) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.name}</span>
          </Text>
        )
      },
      {
        key: 'desc',
        header: t('描述'),
        render: (x: Image) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.displayName}</span>
          </Text>
        )
      },
      {
        key: 'visibility',
        header: t('权限类型'),
        render: (x: Image) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.visibility === 'Public' ? t('公有') : t('私有')}</span>
          </Text>
        )
      },
      {
        key: 'pullCount',
        header: t('下载次数'),
        render: (x: Image) => (
          <Text parent="div" overflow>
            <span className="text">{x.status.pullCount}</span>
          </Text>
        )
      },
      {
        key: 'settings',
        header: '操作',
        width: 100,
        render: image => (
          <React.Fragment>
            <Button
              type="link"
              onClick={() => {
                this.props.actions.image.deleteImage.start([image]);
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
        emptyTips={<div className="text-center">{t('镜像列表为空')}</div>}
        listModel={{
          list: this.props.image.list,
          query: Object.assign({}, this.props.image.query, {
            filter: {
              namespace: this.props.route.queries['ns']
            }
          })
        }}
        actionOptions={this.props.actions.image}
        addons={[
          expandable({
            expandedKeys: this.state.expandedKeys,
            onExpandedKeysChange: keys => this.setState({ expandedKeys: keys }),
            render: (x: Image) => {
              return (
                <table>
                  <tbody>
                    <tr className="tc-detail-row">
                      <td colSpan={columns.length - 1}>
                        <ImageTagTable
                          tags={x.status.tags || []}
                          imageName={x.spec.name}
                          docRegUrl={this.props.dockerRegistryUrl.data}
                          nsName={this.props.route.queries['nsName']}
                        />
                      </td>
                      <td />
                    </tr>
                  </tbody>
                </table>
              );
            },
            gapCell: 1
          })
        ]}
      />
    );
  }

  private _renderDeleteImageDialog() {
    const { actions, deleteImage } = this.props;
    return (
      <WorkflowDialog
        caption={t('删除镜像')}
        workflow={deleteImage}
        action={actions.image.deleteImage}
        targets={deleteImage.targets}
        postAction={() => {}}
        params={{}}
      >
        <DialogBodyLayout>
          <p className="til tea-text-overflow">
            <strong className="tip-top">
              {t('确定要删除该镜像：{{imageName}} 么？', {
                imageName: deleteImage.targets ? deleteImage.targets[0].spec.name : ''
              })}
            </strong>
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
        caption={t('镜像上传指引')}
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
                <strong>拉取 Registry 中指定镜像</strong>
              </p>
              <code>
                <Clip target="#fetchRegistry" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="fetchRegistry">{`sudo docker pull ${this.props.dockerRegistryUrl.data}/${this.props.route.queries['nsName']}/nginx:latest`}</p>
              </code>
              <p className="text-weak">
                {t('以拉取 {{namespace}} 命名空间下名为 nginx 的镜像仓库内版本为 latest 的容器镜像为例', {
                  namespace: this.props.route.queries['nsName']
                })}
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
                <p id="pushRegistry2">{`sudo docker tag nginx:latest ${this.props.dockerRegistryUrl.data}/${this.props.route.queries['nsName']}/nginx:latest`}</p>
              </code>
              <br />
              <code>
                <Clip target="#pushRegistry3" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="pushRegistry3">{`sudo docker push ${this.props.dockerRegistryUrl.data}/${this.props.route.queries['nsName']}/nginx:latest`}</p>
              </code>
              <p className="text-weak">
                <Trans>
                  以推送本地最新版本 nginx 镜像到容器镜像服务内 {this.props.route.queries['nsName']}/nginx 镜像仓库为例
                </Trans>
              </p>
            </li>
          </ul>
        </div>
      </TipDialog>
    );
  }
}

interface TagProps {
  tags?: Tag[];
  imageName?: string;
  docRegUrl?: string;
  nsName?: string;
}

class ImageTagTable extends React.Component<TagProps, any> {
  render() {
    return (
      <Table
        records={this.props.tags}
        columns={[
          {
            key: 'name',
            header: '版本名称',
            render: (x: Tag) => (
              <Text parent="div" overflow>
                <span className="text">{x.name}</span>
              </Text>
            )
          },
          {
            key: 'digest',
            header: '数字摘要',
            render: (x: Tag) => (
              <Text parent="div" overflow>
                <span className="text">{x.digest}</span>
              </Text>
            )
          },
          {
            key: 'createdTime',
            header: '创建时间',
            render: (x: Tag) => (
              <Text parent="div" overflow>
                <span className="text">{dateFormatter(new Date(x.timeCreated), 'YYYY-MM-DD HH:mm:ss')}</span>
              </Text>
            )
          },
          {
            key: 'addr',
            header: '路径',
            render: (x: Tag) => (
              <Text parent="div" overflow>
                <Text
                  className="text"
                  overflow
                  style={{ width: '80%' }}
                  id={`_${x.digest.substring(x.digest.length - 10, x.digest.length)}`}
                >
                  {`${this.props.docRegUrl}/${this.props.nsName}/${this.props.imageName}:${x.name}`}
                </Text>{' '}
                <Clip target={`#_${x.digest.substring(x.digest.length - 10, x.digest.length)}`} />
              </Text>
            )
          }
        ]}
      />
    );
  }
}
