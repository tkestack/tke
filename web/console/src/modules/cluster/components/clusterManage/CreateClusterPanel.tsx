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

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Button, ContentView } from '@tencent/tea-component';

import { resourceConfig } from '../../../../../config';
import { getWorkflowError, InputField, ResourceInfo, TipInfo } from '../../../../modules/common';
import { allActions } from '../../actions';
import { validateClusterCreationAction } from '../../actions/validateClusterCreationAction';
import { CreateResource } from '../../models';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';
import { ClusterSubpageHeaderPanel } from './ClusterSubpageHeaderPanel';
import { KubeconfigFileParse } from './KubeconfigFileParse';
import { AsUserExtraInput } from './AsUserExtraInput';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class CreateClusterPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    const { actions } = this.props;
    const action = actions.workflow.createCluster;
    action.reset();
    actions.clusterCreation.clearClusterCreationState();
  }

  render() {
    let { actions, clusterCreationState, createClusterFlow, route } = this.props,
      {
        v_apiServer,
        v_certFile,
        v_name,
        v_token,
        apiServer,
        certFile,
        name,
        token,
        clientCert,
        clientKey,
        username,
        as,
        asUserExtra,
        v_asUserExtra
      } = clusterCreationState;
    const workflow = createClusterFlow;
    const action = actions.workflow.createCluster;
    const clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }

      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }

      router.navigate({}, { rid: route.queries['rid'] });
    };

    function transAsUserExtra(asUserExtra: { key: string; value: string }[]) {
      if (!asUserExtra || asUserExtra.length <= 0) return undefined;

      return asUserExtra.reduce(
        (all, { key, value }) => ({
          ...all,
          [key]: all?.[key] ? `${all?.[key]},${value}` : value
        }),
        {}
      );
    }

    const perform = () => {
      actions.validate.clusterCreation.validateclusterCreationState();
      if (validateClusterCreationAction._validateclusterCreationState(clusterCreationState)) {
        const tempName = apiServer.substring(8);
        const tempSplit = tempName.split(':');
        let host = tempSplit[0];
        let path = '',
          port = '';
        if (host.indexOf('/') !== -1) {
          const index = host.indexOf('/');
          path = host.substring(index);
          host = host.substring(0, index);
          port = '443';
        } else {
          port = tempSplit[1] ? tempSplit[1].split('/')[0] : '443';
          if (tempSplit[1] && tempSplit[1].indexOf('/') !== -1) {
            path = tempSplit[1] ? tempSplit[1].substring(tempSplit[1].indexOf('/')) : '';
          }
        }
        let certIsBase64;
        try {
          const certOrigin = window.atob(clusterCreationState.certFile);
          certIsBase64 = window.btoa(certOrigin) === clusterCreationState.certFile;
        } catch {
          certIsBase64 = false;
        }
        const data = {
          kind: 'Cluster',
          apiVersion: `${clusterInfo.group}/${clusterInfo.version}`,
          metadata: {
            generateName: 'cls'
          },
          spec: {
            displayName: clusterCreationState.name,
            type: 'Imported'
          },
          status: {
            addresses: [
              {
                host: host,
                type: 'Advertise',
                port: +port,
                path: path
              }
            ],
            credential: {
              caCert: certIsBase64 ? clusterCreationState.certFile : window.btoa(clusterCreationState.certFile),
              clientCert: clusterCreationState.clientCert || undefined,
              clientKey: clusterCreationState.clientKey || undefined,
              token: clusterCreationState.token || undefined,
              username: clusterCreationState.username || undefined,
              as: clusterCreationState.as || undefined,
              'as-user-extra': transAsUserExtra(asUserExtra)
            }
          }
        };

        const createClusterData: CreateResource[] = [
          {
            id: uuid(),
            resourceInfo: clusterInfo,
            mode: 'create',
            jsonData: JSON.stringify(data)
          }
        ];
        action.start(createClusterData);
        action.perform();
      }
    };
    function parseKubeconfigSuccess(params) {
      actions.clusterCreation.updateClusterCreationState({ ...params });
    }

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    return (
      <ContentView>
        <ContentView.Header>
          <ClusterSubpageHeaderPanel />
        </ContentView.Header>
        <ContentView.Body>
          <FormPanel>
            <FormPanel.Item label={t('名称')}>
              <InputField
                type="text"
                value={name}
                placeholder={t('请输入集群名称')}
                tipMode="popup"
                validator={v_name}
                tip={t('最长60个字符')}
                onChange={value => actions.clusterCreation.updateClusterCreationState({ name: value })}
                onBlur={actions.validate.clusterCreation.validateClusterName}
              />
            </FormPanel.Item>
            <FormPanel.Item label="KubeConfig File">
              <KubeconfigFileParse onSuccess={parseKubeconfigSuccess} />
            </FormPanel.Item>
            <FormPanel.Item label="API Server">
              <InputField
                type="text"
                value={apiServer}
                style={{ marginRight: '5px' }}
                placeholder={t('请输入 域名 或 ip地址')}
                tipMode="popup"
                validator={v_apiServer}
                onChange={value => actions.clusterCreation.updateClusterCreationState({ apiServer: value })}
                onBlur={actions.validate.clusterCreation.validateApiServer}
              />
            </FormPanel.Item>
            <FormPanel.Item label="CertFile">
              <InputField
                type="textarea"
                value={certFile}
                placeholder={t('请输入CertFile')}
                tipMode="popup"
                validator={v_certFile}
                onChange={value => actions.clusterCreation.updateClusterCreationState({ certFile: value })}
                onBlur={actions.validate.clusterCreation.validateCertfile}
              />
            </FormPanel.Item>
            <FormPanel.Item label="Token">
              <InputField
                type="textarea"
                value={token}
                placeholder={t('请输入Token')}
                tipMode="popup"
                validator={v_token}
                onChange={value => actions.clusterCreation.updateClusterCreationState({ token: value })}
                onBlur={actions.validate.clusterCreation.validateToken}
              />
            </FormPanel.Item>
            <FormPanel.Item label="Client-Certificate">
              <InputField
                type="textarea"
                value={clientCert}
                placeholder={t('请输入Client-Certificate')}
                tipMode="popup"
                onChange={value => actions.clusterCreation.updateClusterCreationState({ clientCert: value })}
              />
            </FormPanel.Item>
            <FormPanel.Item label="Client-Key">
              <InputField
                type="textarea"
                value={clientKey}
                placeholder={t('请输入Client-Key')}
                tipMode="popup"
                onChange={value => actions.clusterCreation.updateClusterCreationState({ clientKey: value })}
              />
            </FormPanel.Item>

            <FormPanel.Item label="username">
              <InputField
                type="textarea"
                value={username}
                placeholder={t('请输入username')}
                tipMode="popup"
                onChange={value => actions.clusterCreation.updateClusterCreationState({ username: value })}
              />
            </FormPanel.Item>

            <FormPanel.Item label="as">
              <InputField
                type="textarea"
                value={as}
                placeholder={t('请输入as')}
                tipMode="popup"
                onChange={value => actions.clusterCreation.updateClusterCreationState({ as: value })}
              />
            </FormPanel.Item>

            <FormPanel.Item label="as-user-extra">
              <AsUserExtraInput
                data={asUserExtra}
                onChange={value => actions.clusterCreation.updateClusterCreationState({ asUserExtra: value })}
              />
            </FormPanel.Item>

            <FormPanel.Footer>
              <React.Fragment>
                <Button
                  className="m"
                  type="primary"
                  disabled={workflow.operationState === OperationState.Performing}
                  onClick={perform}
                >
                  {failed ? t('重试') : t('提交')}
                </Button>
                <Button type="weak" onClick={cancel}>
                  取消
                </Button>
                <TipInfo type="error" isForm isShow={failed}>
                  {getWorkflowError(workflow)}
                </TipInfo>
              </React.Fragment>
            </FormPanel.Footer>
          </FormPanel>
        </ContentView.Body>
      </ContentView>
    );
  }
}
