import * as React from 'react';
import { RootProps } from './InstallerApp';
import { Button, Bubble, Switch, Card, Alert, List, Form, Segment, Input } from '@tencent/tea-component';
import { EditingItem } from './EditingItem';
import { CIDR } from './CIDR';
import { ListItem } from './ListItem';
import { Machine } from '../models';
import { validateActions } from '../actions/validateActions';
import { getWorkflowError, getValidateStatus } from '../../common/utils';
import { OperationState, isSuccessWorkflow } from '@tencent/ff-redux';

export class FormPanel extends React.Component<RootProps> {
  renderMachines(machines: Array<Machine>) {
    return machines.map(m => {
      return m.status === 'editing' ? (
        <EditingItem {...this.props} id={m.id} />
      ) : (
        <ListItem {...this.props} id={m.id} />
      );
    });
  }
  render() {
    const { editState, actions, createCluster } = this.props,
      { machines } = editState;

    const machine = machines.find(m => m.status === 'editing');
    const canAdd = machine ? validateActions._validateMachine(machine) : true;
    let failed = createCluster.operationState === OperationState.Done && !isSuccessWorkflow(createCluster);
    return (
      <div style={{ maxWidth: '1000px', minHeight: '600px', margin: '0 auto' }} className="server-update">
        <h2 style={{ margin: '40px 0px', fontWeight: 600 }}>TKE Enterprise 安装初始化</h2>
        <Card>
          <Card.Body>
            <Alert>
              <h4 style={{ marginBottom: 10 }}>注意事项</h4>
              <p>
                安装程序将按照选定的目标机器安装管控kubernetes集群，并且根据选定的功能部署tke企业版的服务组件以及控制台，请注意：
              </p>
              <List type="number">
                <List.Item>
                  目标机器只支持centos、ubuntu两种，建议内核4.14以上，64位，并拥有8核16G内存100G硬盘以上的可用资源；
                </List.Item>
                <List.Item>目标机器请保证未自行安装docker、kubernetes等组件；</List.Item>
              </List>
            </Alert>

            <Form layout="fixed">
              <Form.Item label="目标机器">
                {this.renderMachines(editState.machines)}
                <div style={{ width: '100%' }}>
                  <Button
                    type="link"
                    style={{
                      width: '100%',
                      border: '1px dashed #ddd',
                      display: 'block',
                      height: '30px',
                      lineHeight: '30px',
                      padding: '0 20px',
                      boxSizing: 'border-box',
                      textAlign: 'center'
                    }}
                    onClick={() => actions.installer.addMachine()}
                  >
                    添加机器
                  </Button>
                </div>
              </Form.Item>
              <Form.Item label="容器网络">
                <div className="run-docker-box cidr">
                  <div className="edit-param-list">
                    <div className="param-box">
                      <div className="param-bd">
                        <CIDR
                          parts={['192', '172', '10']}
                          value={editState.cidr}
                          maxNodePodNum={editState.podNumLimit}
                          maxClusterServiceNum={editState.serviceNumLimit}
                          minMaskCode="14"
                          maxMaskCode="19"
                          onChange={(cidr, podNumLimit, serviceNumLimit) =>
                            actions.installer.updateEdit({
                              cidr,
                              podNumLimit,
                              serviceNumLimit
                            })
                          }
                          {...this.props}
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </Form.Item>
              <Form.Item label="镜像仓库">
                <Segment
                  options={[
                    { text: '本机仓库', value: 'local' },
                    { text: '远程仓库', value: 'remote' }
                  ]}
                  value={editState.repoType}
                  onChange={value => actions.installer.updateEdit({ repoType: value })}
                />

                {editState.repoType === 'remote' ? (
                  <p className="text">
                    使用远程仓库统作为初始镜像仓库，目标机器将从远程仓库拉取TKE依赖docker镜像，须保证目标机器可访问远程仓库。
                  </p>
                ) : (
                  <p className="text">使用TKE Enterprise 自带的镜像仓库</p>
                )}
                {editState.repoType === 'remote' ? (
                  <div className="run-docker-box" style={{ marginTop: '10px' }}>
                    <Form>
                      <Form.Item
                        label="访问地址"
                        required
                        status={getValidateStatus(editState.v_repoAddress)}
                        message={editState.v_repoAddress.message}
                      >
                        <Input
                          value={editState.repoAddress}
                          onChange={repoAddress => actions.installer.updateEdit({ repoAddress })}
                        />
                      </Form.Item>
                      <Form.Item
                        label="用户名"
                        required
                        status={getValidateStatus(editState.v_repoUser)}
                        message={editState.v_repoUser.message}
                      >
                        <Input
                          value={editState.repoUser}
                          onChange={repoUser => actions.installer.updateEdit({ repoUser })}
                        />
                      </Form.Item>
                      <Form.Item
                        label="密码"
                        required
                        status={getValidateStatus(editState.v_repoPassword)}
                        message={editState.v_repoPassword.message}
                      >
                        <Input
                          type="password"
                          value={editState.repoPassword}
                          onChange={repoPassword => actions.installer.updateEdit({ repoPassword })}
                        />
                      </Form.Item>
                    </Form>
                  </div>
                ) : (
                  <noscript />
                )}
              </Form.Item>
              <Form.Item
                label="域名后缀"
                required
                status={getValidateStatus(editState.v_domain)}
                message={editState.v_domain.message}
              >
                <Input value={editState.domain} onChange={domain => actions.installer.updateEdit({ domain })} />
              </Form.Item>
              <Form.Item label="是否使用已有证书">
                <Switch
                  value={editState.isUseCert}
                  onChange={checked => actions.installer.updateEdit({ isUseCert: checked })}
                />
              </Form.Item>
              <Form.Item
                label="证书(Certificate)"
                style={{ display: editState.isUseCert ? 'table-row' : 'none' }}
                required
                status={getValidateStatus(editState.v_certificate)}
                message={editState.v_certificate.message}
              >
                <Input
                  multiline
                  style={{ width: '400px' }}
                  value={editState.certificate}
                  onChange={certificate => actions.installer.updateEdit({ certificate })}
                />
              </Form.Item>
              <Form.Item
                label="私钥(PrivateKey)"
                style={{ display: editState.isUseCert ? 'table-row' : 'none' }}
                required
                status={getValidateStatus(editState.v_privateKey)}
                message={editState.v_privateKey.message}
              >
                <Input
                  multiline
                  style={{ width: '400px' }}
                  value={editState.privateKey}
                  onChange={privateKey => actions.installer.updateEdit({ privateKey })}
                />
              </Form.Item>
            </Form>
            <Form.Action>
              <Button
                className="tc-15-btn m mr10"
                onClick={() => {
                  actions.installer.stepNext(2);
                }}
              >
                上一步
              </Button>
              <Bubble content={!validateActions._validateEdit(editState) ? '请先完成编辑项' : ''}>
                <Button
                  disabled={createCluster.operationState === OperationState.Performing}
                  type="primary"
                  onClick={() => {
                    actions.validate.validateEdit(editState);
                    if (validateActions._validateEdit(editState)) {
                      actions.installer.createCluster.start([editState]);
                      actions.installer.createCluster.perform();
                    }
                  }}
                >
                  {createCluster.operationState === OperationState.Performing ? (
                    <i className="n-loading-icon" />
                  ) : (
                    <noscript />
                  )}
                  提交
                </Button>
              </Bubble>
              {failed && (
                <Alert
                  type="error"
                  style={{
                    display: 'inline-block',
                    marginTop: '10px',
                    marginBottom: '0px'
                  }}
                >
                  {getWorkflowError(createCluster)}
                </Alert>
              )}
            </Form.Action>
          </Card.Body>
        </Card>
      </div>
    );
  }
}
