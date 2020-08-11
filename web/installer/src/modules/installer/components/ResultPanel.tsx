import { Base64 } from 'js-base64';
import * as React from 'react';

import { insertCSS } from '@tencent/ff-redux';
import { Alert, Button, Card, ExternalLink, Modal } from '@tencent/tea-component';

import { downloadCrt } from '../../../../helpers';
import { Clip, CodeMirrorEditor } from '../../common/components';
import { RootProps } from './InstallerApp';

insertCSS(
  'ResultPanelStyle',
  `
.rich-textarea .copy-btn {
  padding: 1px 3px;
  display: inline-block;
  width: 24px;
}
`
);

interface ResultPanelState {
  isShowertDialog?: boolean;
}
export class ResultPanel extends React.Component<RootProps, ResultPanelState> {
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
    let preStatus = this.props.clusterProgress.data.record['status'],
      nextStatus = nextProps.clusterProgress.data.record['status'];
    if ((!preStatus || preStatus === 'Doing') && nextStatus === 'Success') {
      this.setState({ isShowertDialog: true });
    }
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
    const { clusterProgress } = this.props;
    return (
      <div style={{ maxWidth: '1000px', minHeight: '600px', margin: '0 auto' }}>
        <h2 style={{ margin: '40px 0px', fontWeight: 600 }}>TKE Enterprise 安装初始化</h2>
        <Card>
          <Card.Body>
            <li>
              <CodeMirrorEditor
                isShowHeader={false}
                value={clusterProgress.data.record['data']}
                lineNumbers={true}
                readOnly={true}
                height={500}
                mode="javascript"
              />

              <div style={{ marginTop: '80px' }}>
                {clusterProgress.data.record['status'] === 'Success' ? (
                  <Button className="mr10" onClick={() => this.setState({ isShowertDialog: true })}>
                    查看指引
                  </Button>
                ) : (
                  <noscript />
                )}

                <Button
                  disabled={clusterProgress.data.record['status'] !== 'Success'}
                  onClick={() => {
                    window.location.href = clusterProgress.data.record['url'];
                  }}
                >
                  {clusterProgress.data.record['status'] === 'Doing' ? (
                    <span>
                      <i className="n-loading-icon" />
                      安装中...
                    </span>
                  ) : clusterProgress.data.record['status'] === 'Success' ? (
                    <span>访问TKE控制台</span>
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
                <Modal
                  visible={this.state.isShowertDialog}
                  style={{ width: '540px' }}
                  caption="操作指引"
                  onClose={() => this.setState({ isShowertDialog: false })}
                >
                  <Modal.Body>
                    <h3>初始化完成</h3>
                    <p style={{ margin: '20px 0' }}>请按如下指引完成操作：</p>
                    <div className="tc-processes-vertical ">
                      <ol>
                        <li className="current">
                          <div className="pv-content">
                            <div className="content-inner" style={{ display: 'block' }}>
                              <h4>配置host</h4>
                              <p>
                                请配置本地解析记录（Linux: /etc/hosts， Windows: c:\windows\system32\drivers\etc\hosts）
                              </p>
                              <div className="rich-textarea hide-number">
                                <Clip target={'#host'} className="copy-btn">
                                  复制
                                </Clip>
                                <div className="rich-content" contentEditable={false}>
                                  <pre
                                    className="rich-text"
                                    id="host"
                                    style={{
                                      width: '480px',
                                      marginTop: '0px',
                                      marginBottom: '0px',
                                      whiteSpace: 'pre-wrap',
                                      overflow: 'auto',
                                      height: '200px'
                                    }}
                                  >
                                    {this.getHost(
                                      clusterProgress.data.record['hosts'],
                                      clusterProgress.data.record['servers']
                                    )}
                                  </pre>
                                </div>
                              </div>
                            </div>
                          </div>
                        </li>
                        <li className="current">
                          <div className="pv-content">
                            <div className="content-inner" style={{ display: 'block' }}>
                              <h4>导入根证书</h4>
                              {clusterProgress.data.record['caCert'] ? (
                                <p>请导入根证书到系统“受信任的根证书颁发机构”中：</p>
                              ) : (
                                <p>您已选择使用自有证书，请跳过此步</p>
                              )}
                              {clusterProgress.data.record['caCert'] ? (
                                <p className="op-area">
                                  <div className="rich-textarea hide-number">
                                    <Clip target={'#caCert'} className="copy-btn">
                                      复制
                                    </Clip>
                                    <Button
                                      type="link"
                                      onClick={e =>
                                        downloadCrt(Base64.decode(clusterProgress.data.record['caCert']), 'caCert.crt')
                                      }
                                      className="copy-btn"
                                      style={{ right: '40px' }}
                                    >
                                      下载
                                    </Button>
                                    <div className="rich-content" contentEditable={false}>
                                      <pre
                                        className="rich-text"
                                        id="caCert"
                                        style={{
                                          width: '480px',
                                          marginTop: '0px',
                                          marginBottom: '0px',
                                          whiteSpace: 'pre-wrap',
                                          overflow: 'auto',
                                          height: '200px'
                                        }}
                                      >
                                        {Base64.decode(clusterProgress.data.record['caCert'])}
                                      </pre>
                                    </div>
                                  </div>
                                </p>
                              ) : (
                                <noscript />
                              )}
                            </div>
                          </div>
                        </li>
                        <li className="current">
                          <div className="pv-content">
                            <div className="content-inner" style={{ display: 'block' }}>
                              <h4>访问控制台</h4>
                              <p>
                                现在可访问TKE Enterprise控制台：
                                <ExternalLink href={clusterProgress.data.record['url']}>
                                  {clusterProgress.data.record['url']}
                                </ExternalLink>
                              </p>
                              <p>
                                用户名：
                                {clusterProgress.data.record['username']}
                                <p>
                                  密&nbsp;&nbsp;&nbsp;码：
                                  {clusterProgress.data.record['password']
                                    ? Base64.decode(clusterProgress.data.record['password'])
                                    : ''}
                                </p>
                              </p>
                            </div>
                          </div>
                        </li>
                      </ol>
                    </div>
                  </Modal.Body>
                </Modal>
              </div>
            </li>
          </Card.Body>
        </Card>
      </div>
    );
  }
}
