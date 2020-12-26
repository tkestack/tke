import * as React from 'react';

import { Button, ExternalLink, Form, Input, Segment, Switch, Text } from '@tencent/tea-component';

import { getValidateStatus } from '../../common/utils';
import { validateActions } from '../actions/validateActions';
import { Machine } from '../models';
import { Arg } from '../models/RootState';
import { CIDR } from './CIDR';
import { EditingItem } from './EditingItem';
import { RootProps } from './InstallerApp';
import { ListItem } from './ListItem';

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
        remove: actions.installer.removeDockerExtraArgs,
        keyPlaceholder: 'dns-search',
        valuePlaceholder: 'example.com'
      },
      apiServerExtraArgs: {
        update: actions.installer.updateApiServerExtraArgs,
        remove: actions.installer.removeApiServerExtraArgs,
        keyPlaceholder: 'etcd-prefix',
        valuePlaceholder: '/k8s-registry'
      },
      controllerManagerExtraArgs: {
        update: actions.installer.updateControllerManagerExtraArgs,
        remove: actions.installer.removeControllerManagerExtraArgs,
        keyPlaceholder: 'secure-port',
        valuePlaceholder: '10257'
      },
      schedulerExtraArgs: {
        update: actions.installer.updateSchedulerExtraArgs,
        remove: actions.installer.removeSchedulerExtraArgs,
        keyPlaceholder: 'port',
        valuePlaceholder: '10251'
      },
      kubeletExtraArgs: {
        update: actions.installer.updateKubeletExtraArgs,
        remove: actions.installer.removeKubeletExtraArgs,
        keyPlaceholder: 'config',
        valuePlaceholder: '/etc/kubernetes/kubelet.config'
      }
    };

    return args.map((arg, index) => (
      <section style={{ marginBottom: '10px' }} key={index}>
        <Input
          size="s"
          value={arg.key}
          onChange={value => actionsMap[type].update({ key: value }, arg.id)}
          placeholder={actionsMap[type].keyPlaceholder}
        />
        <Text style={{ margin: '0px 10px', fontSize: '12px' }}>=</Text>
        <Input
          value={arg.value}
          onChange={value => actionsMap[type].update({ value: value }, arg.id)}
          placeholder={actionsMap[type].valuePlaceholder}
        />
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
            message={
              editState.v_networkDevice.message ||
              '设置集群节点使用的网卡，Galaxy插件的floating IP功能会使用该网卡做桥接，如无特殊情况，一般为eth0'
            }
          >
            <Input
              value={editState.networkDevice}
              onChange={value => actions.installer.updateEdit({ networkDevice: value })}
            />
          </Form.Item>
          <Form.Item
            label="GPU类型"
            message={
              editState.gpuType === 'none' ? (
                '选择GPU类型后，平台将自动为节点安装相应的GPU驱动和运行时工具'
              ) : editState.gpuType === 'Virtual' ? (
                <>
                  平台会自动为集群安装
                  <ExternalLink
                    href={
                      'https://github.com/tkestack/docs/blob/master/docs/zh/%E4%BA%A7%E5%93%81%E7%89%B9%E8%89%B2%E5%8A%9F%E8%83%BD/GPUManager.md'
                    }
                  >
                    GPUManager
                  </ExternalLink>
                  扩展组件
                </>
              ) : (
                <>
                  平台会自动为集群安装
                  <ExternalLink href={'https://github.com/NVIDIA/k8s-device-plugin'}>nvidia-k8s-plugin</ExternalLink>
                </>
              )
            }
          >
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
                      为docker(19.03.14)运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://docs.docker.com/engine/reference/commandline/run/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addDockerExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                  <Form.Item label="kube-apiserver设置">
                    {this.renderArg(editState.apiServerExtraArgs, 'apiServerExtraArgs')}
                    <Form.Text>
                      为kube-apiserver(1.18.3)运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addApiServerExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                  <Form.Item label="kube-controller-manager设置">
                    {this.renderArg(editState.controllerManagerExtraArgs, 'controllerManagerExtraArgs')}
                    <Form.Text>
                      为kube-controller-manager(1.18.3)运行设置自定义参数，默认不需要添加，详细请参考
                      <ExternalLink href="https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/">
                        帮助文档
                      </ExternalLink>
                    </Form.Text>
                    <Button type="link" onClick={() => actions.installer.addControllerManagerExtraArgs()}>
                      添加
                    </Button>
                  </Form.Item>
                  <Form.Item label="kube-scheduler设置">
                    {this.renderArg(editState.schedulerExtraArgs, 'schedulerExtraArgs')}
                    <Form.Text>
                      为kube-scheduler(1.18.3)运行设置自定义参数，默认不需要添加，详细请参考
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
                      为kubelet(1.18.3)运行设置自定义参数，默认不需要添加，详细请参考
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
