import * as React from 'react';
import { RootProps } from './ClusterApp';
import { Modal, Button, Switch, Alert, Text } from '@tea/component';
import { FormPanel, Clip } from '../../common/components';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { allActions } from '../actions';
import { connect } from 'react-redux';
import { Cluster } from '../../common/models';
import { downloadCrt } from '../../../../helpers';
import { DialogNameEnum } from '../models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class KubectlDialog extends React.Component<RootProps, any> {
  render() {
    let { dialogState, actions, cluster, clustercredential } = this.props,
      { kuberctlDialog } = dialogState;

    let clusterInfo: Cluster = cluster.selection;

    const cancel = () => {
      actions.cluster.clearClustercredential();
      actions.dialog.updateDialogState(DialogNameEnum.kuberctlDialog);
    };

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
                  onClick={e => downloadCrt(clusterInfo.status.credential.caCert)}
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
                      width: '475px',
                      whiteSpace: 'pre-wrap',
                      overflow: 'auto',
                      height: '300px'
                    }}
                  >
                    {clustercredential.caCert && window.atob(clustercredential.caCert)}
                  </p>
                </div>
              </div>
            </FormPanel.Item>
          </FormPanel>
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
