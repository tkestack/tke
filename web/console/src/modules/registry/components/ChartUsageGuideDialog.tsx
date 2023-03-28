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
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { RootState } from '../models';
import { allActions } from '../actions';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Clip, GridTable, TipDialog, WorkflowDialog } from '../../common/components';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

export interface RootProps extends RootState {
  actions?: typeof allActions;
}
interface ChartUsageGuideDialogProps extends RootProps {
  showDialog: boolean;
  chartGroupName: string;
  registryUrl: string;
  username: string;
  onClose: Function;
}

interface ChartUsageGuideDialogState extends RootProps {
  showDialog?: boolean;
}

@connect(state => state, mapDispatchToProps)
export class ChartUsageGuideDialog extends React.Component<ChartUsageGuideDialogProps, ChartUsageGuideDialogState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      showDialog: this.props.showDialog
    };
  }

  componentWillReceiveProps(nextProps: ChartUsageGuideDialogProps) {
    const { showDialog } = nextProps;
    if (showDialog !== this.props.showDialog) {
      this.setState({
        showDialog: showDialog
      });
    }
  }

  render() {
    return (
      <TipDialog
        isShow={this.state.showDialog}
        width={680}
        caption={t('Chart 上传指引')}
        cancelAction={() => {
          this.setState({ showDialog: false });
          this.props.onClose && this.props.onClose();
        }}
        performAction={() => {
          this.setState({ showDialog: false });
          this.props.onClose && this.props.onClose();
        }}
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
                  <a href="https://helm.sh/docs/intro/quickstart/" target="_blank" rel="noreferrer">
                    安装 Helm
                  </a>
                  .{' '}
                </Trans>
              </p>
              <code>
                <Clip target="#installHelm" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="installHelm">{`curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | DESIRED_VERSION=v3.6.2 bash`}</p>
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
                <p id="addTkeRepo">{`helm repo add ${this.props.chartGroupName} http://${this.props.registryUrl}/chart/${this.props.chartGroupName} --username ${this.props.username} --password [访问凭证] `}</p>
              </code>
              <p className="text-weak">
                <Trans>
                  获取有效访问凭证信息，请前往
                  <a
                    href="javascript:;"
                    onClick={() => {
                      const urlParams = router.resolve(this.props.route);
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
                <p id="installHelmPush">{`helm plugin install https://github.com/chartmuseum/helm-push`}</p>
              </code>
              <p className="text-weak">
                <Trans>
                  如安装失败，可以手动下载后解压到$HOME/.local/share/helm/plugins/helm-push，解压路径可以通过helm env查看
                </Trans>
              </p>
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
                <p id="pushHelmDir">{`helm cm-push ./myapp ${this.props.chartGroupName}`}</p>
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
                <p id="pushHelmTar">{`helm cm-push myapp-1.0.1.tgz ${this.props.chartGroupName}`}</p>
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
                <p id="downloadChart">{`helm pull ${this.props.chartGroupName}/myapp`}</p>
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
                <p id="downloadSChart">{`helm pull ${this.props.chartGroupName}/myapp --version 1.0.1`}</p>
              </code>
            </li>
          </ul>
        </div>
      </TipDialog>
    );
  }
}
