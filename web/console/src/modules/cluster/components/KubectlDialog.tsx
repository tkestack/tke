import * as React from 'react';
import { connect } from 'react-redux';

import { Alert, Button, Modal, Switch, Text, ExternalLink } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { downloadCrt, downloadKubeconfig, getKubectlConfig } from '../../../../helpers';
import { Clip } from '../../common/components';
import { Cluster } from '../../common/models';
import { allActions } from '../actions';
import { DialogNameEnum } from '../models';
import { RootProps } from './ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class KubectlDialog extends React.Component<RootProps, any> {
  render() {
    let { dialogState, actions, cluster, clustercredential } = this.props,
      { kuberctlDialog } = dialogState;

    const clusterInfo: Cluster = cluster.selection;

    const cancel = () => {
      actions.cluster.clearClustercredential();
      actions.dialog.updateDialogState(DialogNameEnum.kuberctlDialog);
    };

    const clusterId = clusterInfo && clusterInfo.metadata.name;

    function getHost(clusterInfo: Cluster) {
      if (!clusterInfo) {
        return '';
      }

      const address =
        clusterInfo.status.addresses.find(({ type }) => type === 'Advertise') || clusterInfo.status.addresses[0];

      return `https://${address.host}:${address.port}${address.path}`;
    }

    const kubeconfig = getKubectlConfig({
      caCert: clustercredential.caCert,
      token: clustercredential.token,
      host: getHost(clusterInfo),
      clusterId: clustercredential.clusterName
    });

    return kuberctlDialog ? (
      <Modal visible={true} caption={t('集群凭证')} onClose={cancel} size={680} disableEscape={true}>
        <Modal.Body>
          <Alert>{t('安装Kubectl后，您可以通过用户名密码或集群CA证书登录到集群')}</Alert>
          <FormPanel isNeedCard={false}>
            <FormPanel.Item label={t('集群APIServer地址')} text>
              {clusterInfo.status.addresses.map((item, index) => {
                return (
                  <React.Fragment key={index}>
                    <Text theme="strong" id={`apiserver${index}`}>{`${item.host}:${item.port}`}</Text>
                    <Text theme="label">{` (${item.type})`}</Text>
                    <Clip target={`#apiserver${index}`} />
                  </React.Fragment>
                );
              })}
            </FormPanel.Item>
            <FormPanel.Item label={t('Token')}>
              <div className="rich-textarea hide-number">
                <Clip target={'#token'} className="copy-btn">
                  {t('复制')}
                </Clip>
                <div className="rich-content" contentEditable={false}>
                  <p
                    className="rich-text"
                    id="token"
                    style={{ width: '475px', whiteSpace: 'pre-wrap', height: '25px', marginTop: '0px' }}
                  >
                    {clustercredential.token}
                  </p>
                </div>
              </div>
            </FormPanel.Item>
            <FormPanel.Item label={t('集群CA证书')}>
              <div className="rich-textarea hide-number">
                <Clip target={'#certificationAuthority'} className="copy-btn">
                  {t('复制')}
                </Clip>
                <a
                  href="javascript:void(0)"
                  onClick={e => downloadCrt(clustercredential.caCert ? window.atob(clustercredential.caCert) : '')}
                  className="copy-btn"
                  style={{ right: '50px' }}
                >
                  {t('下载')}
                </a>
                <div className="rich-content" contentEditable={false}>
                  <p
                    className="rich-text"
                    id="certificationAuthority"
                    style={{
                      width: '480px',
                      whiteSpace: 'pre-wrap',
                      overflow: 'auto',
                      height: '64px'
                    }}
                  >
                    {clustercredential.caCert && window.atob(clustercredential.caCert)}
                  </p>
                </div>
              </div>
            </FormPanel.Item>

            <FormPanel.Item label="Kubeconfig">
              <div className="rich-textarea hide-number">
                <Clip target={'#Kubeconfig'} className="copy-btn">
                  {t('复制')}
                </Clip>
                <a
                  href="javascript:void(0)"
                  onClick={e => downloadKubeconfig(kubeconfig, `${clusterId}-config`)}
                  className="copy-btn"
                  style={{ right: '50px' }}
                >
                  {t('下载')}
                </a>
                <div className="rich-content" contentEditable={false}>
                  <p
                    className="rich-text"
                    id="Kubeconfig"
                    style={{
                      width: '480px',
                      whiteSpace: 'pre-wrap',
                      overflow: 'auto',
                      height: '64px'
                    }}
                  >
                    {kubeconfig}
                  </p>
                </div>
              </div>
            </FormPanel.Item>
          </FormPanel>

          <div
            style={{
              textAlign: 'left',
              borderTop: '1px solid #D1D2D3',
              paddingTop: '10px',
              marginTop: '10px',
              color: '#444'
            }}
          >
            <Trans>
              <h3 style={{ marginBottom: '1em' }}>通过Kubectl连接Kubernetes集群操作说明:</h3>
              <p style={{ marginBottom: '5px' }}>
                1.
                <ExternalLink href="https://kubernetes.io/zh/docs/tasks/tools/install-kubectl/">
                  安装和设置kubectl
                </ExternalLink>
                。
              </p>
              <p style={{ marginBottom: '5px' }}>2. 配置 Kubeconfig：</p>
              <ul>
                <li style={{ listStyle: 'disc', marginLeft: '15px' }}>
                  <p style={{ marginBottom: '5px' }}>
                    若当前访问客户端尚未配置任何集群的访问凭证，即 ~/.kube/config 内容为空，可直接复制上方 kubeconfig
                    访问凭证内容并粘贴入 ~/.kube/config 中。
                  </p>
                </li>
                <li style={{ listStyle: 'disc', marginLeft: '15px' }}>
                  <p style={{ marginBottom: '5px' }}>
                    若当前访问客户端已配置了其他集群的访问凭证，你可下载上方 kubeconfig
                    至指定位置，并执行以下指令以合并多个集群的 config。
                  </p>
                  <div className="rich-textarea hide-number" style={{ width: '100%' }}>
                    <div className="rich-content">
                      <Clip target={'#kubeconfig-merge'} className="copy-btn">
                        复制
                      </Clip>
                      <pre
                        className="rich-text"
                        id="kubeconfig-merge"
                        style={{
                          whiteSpace: 'pre-wrap',
                          overflow: 'auto'
                        }}
                      >
                        KUBECONFIG=~/.kube/config:~/Downloads/{{ clusterId }}-config kubectl config view --merge
                        --flatten &gt; ~/.kube/config
                        <br />
                        export KUBECONFIG=~/.kube/config
                      </pre>
                    </div>
                  </div>
                  <p style={{ marginBottom: '5px' }}>
                    其中，~/Downloads/{{ clusterId }}-config 为本集群的 kubeconfig
                    的文件路径，请替换为下载至本地后的实际路径。
                  </p>
                </li>
              </ul>
              <p style={{ marginBottom: '5px' }}>3. 访问 Kubernetes 集群：</p>
              <ul>
                <li style={{ marginLeft: '15px' }}>
                  <p style={{ marginBottom: '5px' }}>
                    完成 kubeconfig 配置后，执行以下指令查看并切换 context 以访问本集群：
                  </p>
                  <div className="rich-textarea hide-number" style={{ width: '100%' }}>
                    <div className="rich-content">
                      <Clip target={'#kubeconfig-visit'} className="copy-btn">
                        复制
                      </Clip>
                      <pre
                        className="rich-text"
                        id="kubeconfig-visit"
                        style={{
                          whiteSpace: 'pre-wrap',
                          overflow: 'auto'
                        }}
                      >
                        kubectl config get-contexts
                        <br />
                        kubectl config use-context {{ clusterId }}-context-default
                      </pre>
                    </div>
                  </div>
                  <p style={{ marginBottom: '5px' }}>
                    可执行 kubectl get nodes
                    测试集群是否可正常访问集群。若无法连接，请确保客户端地址和APIServer地址在同一网络环境下。
                  </p>
                </li>
              </ul>
            </Trans>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={cancel}>
            {t('关闭')}
          </Button>
        </Modal.Footer>
      </Modal>
    ) : (
      <noscript />
    );
  }
}
