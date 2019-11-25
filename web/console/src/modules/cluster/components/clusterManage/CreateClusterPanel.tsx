import * as React from 'react';
import { RootProps } from '../ClusterApp';
import { ContentView, Form, Button, Card, Justify, Input } from '@tencent/tea-component';
import { ClusterSubpageHeaderPanel } from './ClusterSubpageHeaderPanel';
import { OperationState, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { router } from '../../router';
import { CreateResource } from '../../models';
import { connect } from 'react-redux';
import { validateClusterCreationAction } from '../../actions/validateClusterCreationAction';
import { resourceConfig } from '../../../../../config';
import { uuid } from '../../../../../lib/_util';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { allActions } from '../../actions';
import { ResourceInfo, InputField, getWorkflowError, TipInfo, FormPanel } from '../../../../modules/common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class CreateClusterPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    let action = actions.workflow.createCluster;
    action.reset();
    actions.clusterCreation.clearClusterCreationState();
  }

  render() {
    let { actions, clusterCreationState, createClusterFlow, route } = this.props,
      {
        v_apiServer,
        v_certFile,
        port,
        v_port,
        v_name,
        v_token,
        apiServer,
        certFile,
        name,
        token
      } = clusterCreationState;
    const workflow = createClusterFlow;
    const action = actions.workflow.createCluster;
    let clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }

      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }

      router.navigate({}, { rid: route.queries['rid'] });
    };

    const perform = () => {
      actions.validate.clusterCreation.validateclusterCreationState();
      if (validateClusterCreationAction._validateclusterCreationState(clusterCreationState)) {
        let certIsBase64;
        try {
          let certOrigin = window.atob(clusterCreationState.certFile);
          certIsBase64 = window.btoa(certOrigin) === clusterCreationState.certFile;
        } catch {
          certIsBase64 = false;
        }
        let data = {
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
                host: clusterCreationState.apiServer,
                type: 'Advertise',
                port: +clusterCreationState.port
              }
            ],
            credential: {
              caCert: certIsBase64 ? clusterCreationState.certFile : window.btoa(clusterCreationState.certFile)
            }
          }
        };
        if (clusterCreationState.token) {
          data.status.credential['token'] = clusterCreationState.token;
        }

        let createClusterData: CreateResource[] = [
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
              <InputField
                type="text"
                value={port}
                placeholder={t('请输入 port')}
                tipMode="popup"
                validator={v_port}
                onChange={value => actions.clusterCreation.updateClusterCreationState({ port: value })}
                onBlur={actions.validate.clusterCreation.validatePort}
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
