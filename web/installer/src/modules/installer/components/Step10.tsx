import * as React from 'react';
import { RootProps } from './InstallerApp';
import { Button, Form, Alert, Modal, ExternalLink, Stepper, List, Text } from '@tencent/tea-component';
import { CodeMirrorEditor } from '../../common/components';
import { Base64 } from 'js-base64';
import { downloadCrt } from '../../../../helpers';

interface Step10State {
  isShowertDialog?: boolean;
}

export class Step10 extends React.Component<RootProps, Step10State> {
  state = {
    isShowertDialog: false
  };
  componentDidMount() {
    const { actions } = this.props;
    actions.installer.poll();
  }

  componentWillUnmount() {
    clearInterval(window['pollProgress']);
  }

  componentWillReceiveProps(nextProps) {
    // let preStatus = this.props.clusterProgress.data.record['status'],
    //   nextStatus = nextProps.clusterProgress.data.record['status'];
    // if ((!preStatus || preStatus === 'Doing') && nextStatus === 'Success') {
    //   this.setState({ isShowertDialog: true });
    // }
  }

  getHost(hostArr: string[], serverArr: string[]) {
    if (!hostArr || !hostArr.length || !serverArr || !serverArr.length) {
      return '';
    } else {
      let hosts = '';
      hostArr.forEach(h => {
        hosts += serverArr[0] + ' ' + h + '\n';
      });
      return hosts;
    }
  }
  render() {
    const { step } = this.props;
    const { clusterProgress } = this.props;

    let steps = [];
    if (clusterProgress.data.record['hosts'] && clusterProgress.data.record['hosts'].length) {
      steps.push({
        id: 'hosts',
        label: '配置域名解析',
        detail: (
          <section>
            <p>
              请配置如下DNS解析，临时访问也可以配置本地解析（Linux: /etc/hosts， Windows:
              c:\windows\system32\drivers\etc\hosts）
            </p>
            <div className="rich-textarea hide-number">
              <div className="rich-content" contentEditable={false}>
                <pre
                  className="rich-text"
                  id="host"
                  style={{
                    width: '520px',
                    marginTop: '0px',
                    marginBottom: '0px',
                    whiteSpace: 'pre-wrap',
                    overflow: 'auto'
                  }}
                >
                  {this.getHost(clusterProgress.data.record['hosts'], clusterProgress.data.record['servers'])}
                </pre>
              </div>
            </div>
          </section>
        )
      });
    }
    steps.push({
      id: 'access',
      label: '集群访问',
      detail: (
        <List type="number">
          <List.Item>
            访问TKE Stack控制台：
            <ExternalLink href={clusterProgress.data.record['url']}>{clusterProgress.data.record['url']}</ExternalLink>
            <List type="bullet">
              <List.Item>
                用户名：
                {clusterProgress.data.record['username']}
              </List.Item>
              <List.Item>
                密&nbsp;&nbsp;&nbsp;码：
                {clusterProgress.data.record['password'] ? Base64.decode(clusterProgress.data.record['password']) : ''}
              </List.Item>
            </List>
          </List.Item>
          <List.Item>
            通过kubeconfig 访问global集群：
            <List type="bullet">
              <List.Item>
                <Button
                  type="link"
                  onClick={() =>
                    downloadCrt(Base64.decode(clusterProgress.data.record['kubeconfig']), 'kubeconfig.txt')
                  }
                >
                  下载kubeconfig
                </Button>
              </List.Item>
              <Text theme="label">
                TKE
                global集群的访问凭证信息存储在：部署节点/opt/tke-installer/data/目录下，您可以在此目录再次查看以上信息，关于kubeclt
                的配置，请参考
                <ExternalLink href="https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/">
                  帮助文档
                </ExternalLink>
              </Text>
            </List>
          </List.Item>
        </List>
      )
    });
    return step === 'step10' ? (
      <section>
        <CodeMirrorEditor
          isShowHeader={false}
          value={clusterProgress.data.record['data']}
          lineNumbers={true}
          readOnly={true}
          height={500}
          mode="javascript"
        />
        <Modal
          visible={this.state.isShowertDialog}
          caption="操作指引"
          size="l"
          onClose={() => this.setState({ isShowertDialog: false })}
          disableEscape
        >
          <Modal.Body style={{ maxHeight: '600px', overflow: 'auto' }}>
            <p style={{ marginBottom: '20px' }}>请按如下指引进行操作：</p>
            <Stepper type="process-vertical-dot" steps={steps} />
          </Modal.Body>
        </Modal>

        <Form.Action style={{ position: 'absolute', bottom: '20px', left: '20px', width: '960px' }}>
          <Button
            disabled={clusterProgress.data.record['status'] !== 'Success'}
            onClick={() => this.setState({ isShowertDialog: true })}
          >
            {clusterProgress.data.record['status'] === 'Doing' ? (
              <span>
                <i className="n-loading-icon" />
                安装中...
              </span>
            ) : clusterProgress.data.record['status'] === 'Success' ? (
              <span>查看指引</span>
            ) : (
              <span>安装失败</span>
            )}
          </Button>
          {clusterProgress.data.record['status'] === 'Failed' && (
            <Alert
              type="error"
              style={{
                display: 'inline-block',
                marginTop: '0px',
                marginBottom: '0px',
                marginLeft: '20px'
              }}
            >
              安装失败
            </Alert>
          )}
        </Form.Action>
      </section>
    ) : (
      <noscript></noscript>
    );
  }
}
