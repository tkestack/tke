import * as React from 'react';
import { RootProps } from './InstallerApp';
import { Button, Input, Form, Segment, Switch, Text, ExternalLink } from '@tencent/tea-component';
import { CIDR } from './CIDR';
import { EditingItem } from './EditingItem';
import { ListItem } from './ListItem';
import { Machine } from '../models';
import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';
import { Arg } from '../models/RootState';

export class Step3 extends React.Component<RootProps> {
  state = {
    openAdvanced: false
  };

  renderMachines(machines: Array<Machine>) {
    return machines.map(m => {
      return m.status === 'editing' ? (
        <EditingItem {...this.props} id={m.id} />
      ) : (
        <ListItem {...this.props} id={m.id} />
      );
    });
  }

  renderArg(args: Array<Arg>, type: string) {
    const { actions } = this.props;

    const actionsMap = {
      dockerExtraArgs: {
        update: actions.installer.updateDockerExtraArgs,
        remove: actions.installer.removeDockerExtraArgs
      },
      apiServerExtraArgs: {
        update: actions.installer.updateApiServerExtraArgs,
        remove: actions.installer.removeApiServerExtraArgs
      },
      controllerManagerExtraArgs: {
        update: actions.installer.updateControllerManagerExtraArgs,
        remove: actions.installer.removeControllerManagerExtraArgs
      },
      schedulerExtraArgs: {
        update: actions.installer.updateSchedulerExtraArgs,
        remove: actions.installer.removeSchedulerExtraArgs
      },
      kubeletExtraArgs: {
        update: actions.installer.updateKubeletExtraArgs,
        remove: actions.installer.removeKubeletExtraArgs
      }
    };

    return args.map(arg => (
      <section style={{ marginBottom: '10px' }}>
        <Input size="s" value={arg.key} onChange={value => actionsMap[type].update({ key: value }, arg.id)} />
        <Text style={{ margin: '0px 10px', fontSize: '12px' }}>=</Text>
        <Input value={arg.value} onChange={value => actionsMap[type].update({ value: value }, arg.id)} />
        <Button type="link" onClick={() => actionsMap[type].remove(arg.id)}>
          <i className="icon-cancel-icon" />
        </Button>
      </section>
    ));
  }

  render() {
    const { actions, editState, step } = this.props;
    return step === 'step3' ? (
      <section>
        <Form>
          <Form.Item
            label="网卡名称"
            required
            status={getValidateStatus(editState.v_networkDevice)}
            message={editState.v_networkDevice.message || '设置集群使用的网卡，如无特殊情况，一般为eth0'}
          >
            <Input
              value={editState.networkDevice}
              onChange={value => actions.installer.updateEdit({ networkDevice: value })}
            />
          </Form.Item>
          <Form.Item label="GPU类型" message="选择GPU类型后，平台将自动为节点安装相应的GPU驱动和运行时工具">
            <Segment
              value={editState.gpuType}
              options={[
                { text: '不使用', value: 'none' },
                { text: 'Virtual', value: 'Virtual' },
                { text: 'Physical', value: 'Physical' }
              ]}
              onChange={value => actions.installer.updateEdit({ gpuType: value })}
            />
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
          <Form.Item label="master节点">
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
          <Form.Item label="高级设置">
            <Switch value={this.state.openAdvanced} onChange={value => this.setState({ openAdvanced: value })}></Switch>
            {this.state.openAdvanced ? (
              <div className="run-docker-box" style={{ width: '100%', marginTop: '10px' }}>
                <Form>
                  <Form.Item label="docker设置">
                    {this.renderArg(editState.dockerExtraArgs, 'dockerExtraArgs')}
                    <Form.Text>
                      为docker运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://docs.docker.com/engine/reference/commandline/run/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addDockerExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                  <Form.Item label="kube-apiserver设置">
                    {this.renderArg(editState.apiServerExtraArgs, 'apiServerExtraArgs')}}
                    <Form.Text>
                      为kube-apiserver运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addApiServerExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                  <Form.Item label="kube-controller-manager设置">
                    {this.renderArg(editState.controllerManagerExtraArgs, 'controllerManagerExtraArgs')}}
                    <Form.Text>
                      为kube-controller-manager运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addControllerManagerExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                  <Form.Item label="kube-scheduler设置">
                    {this.renderArg(editState.schedulerExtraArgs, 'schedulerExtraArgs')}}
                    <Form.Text>
                      为kube-scheduler运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addSchedulerExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                  <Form.Item label="kubelet设置">
                    {this.renderArg(editState.kubeletExtraArgs, 'kubeletExtraArgs')}
                    <Form.Text>
                      为kubelet运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addKubeletExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                </Form>
              </div>
            ) : (
              <noscript />
            )}
          </Form.Item>
        </Form>
        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button style={{ marginRight: '10px' }} type="weak" onClick={() => actions.installer.stepNext('step2')}>
            上一步
          </Button>
          <Button
            type="primary"
            onClick={() => {
              actions.validate.validateStep3(editState);
              if (validateActions._validateStep3(editState)) {
                actions.installer.stepNext('step4');
              }
            }}
          >
            下一步
          </Button>
        </Form.Action>
      </section>
    ) : (
      <noscript></noscript>
    );
  }
}
