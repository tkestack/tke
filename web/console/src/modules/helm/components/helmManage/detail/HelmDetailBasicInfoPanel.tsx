import * as React from 'react';
import { RootProps } from '../../HelmApp';
import { FormPanel } from '@tencent/ff-component';
import { dateFormatter } from '../../../../../../helpers';
import { Switch, Modal, Table, TableColumn, Text, ContentView, Card } from '@tea/component';
import { HelmResource } from '../../../models';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import { ResourceUrl } from '../../../constants/Config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

interface State {
  showYamlDialog?: boolean;
  yaml?: string;
}
export class HelmDetailBasicInfoPanel extends React.Component<RootProps, State> {
  state = {
    showYamlDialog: false,
    yaml: ''
  };
  showYaml(yaml) {
    this.setState({
      showYamlDialog: true,
      yaml
    });
  }
  _renderYamlDialog() {
    const cancel = () => this.setState({ showYamlDialog: false });

    return (
      <Modal visible={true} caption={t('查看YAML')} onClose={cancel} size={960} disableEscape={true}>
        <Modal.Body>
          <CodeMirror
            value={this.state.yaml}
            options={{
              lineNumbers: true,
              mode: 'yaml',
              theme: 'monokai',
              readOnly: true,
              lineWrapping: true, // 自动换行
              styleActiveLine: true // 当前行背景高亮
            }}
          />
        </Modal.Body>
      </Modal>
    );
  }

  click(resource) {
    let { route } = this.props;
    let rid = route.queries['rid'],
      clusterId = route.queries['clusterId'];

    open(ResourceUrl[resource.kind](rid, clusterId));
  }
  _renderTablePanel() {
    let {
      actions,
      detailState: { helm },

      route
    } = this.props;

    const columns: TableColumn<HelmResource>[] = [
      {
        key: 'name',
        header: t('资源'),
        width: '40%',
        render: x => (
          <Text parent="div" overflow>
            {/* {ResourceUrl[x.kind] ? (
              <a href="javascript:void(0);" onClick={e => this.click(x)}>
                {x.name}
              </a>
            ) : ( */}
            {x.name}
          </Text>
        )
      },
      {
        key: 'kind',
        header: t('类型'),
        width: '40%',
        render: x => (
          <Text parent="div" overflow>
            {x.kind}
          </Text>
        )
      },
      {
        key: 'operation',
        header: t('操作'),
        width: '20%',
        render: x => (
          <a href="javascript:void(0)" onClick={e => this.showYaml(x.yaml)}>
            {t('查看YAML')}
          </a>
        )
      }
    ];

    return <Table columns={columns} records={helm.resources} />;
  }

  renderValueYaml() {
    return (
      <CodeMirror
        className={'codeMirrorHeight'}
        value={this.props.detailState.helm.configYaml}
        options={{
          lineNumbers: true,
          mode: 'yaml',
          theme: 'monokai',
          readOnly: true,
          lineWrapping: true, // 自动换行
          styleActiveLine: true // 当前行背景高亮
        }}
      />
    );
  }

  setRefresh(isRefresh: boolean) {
    this.props.actions.detail.setRefresh(isRefresh);
  }

  render() {
    let { actions, detailState, route } = this.props,
      { helm, isRefresh } = detailState;
    if (!helm) {
      return <noscript />;
    }

    return (
      <ContentView>
        <ContentView.Body>
          <FormPanel title={t('基本信息')}>
            <FormPanel.Item text label={t('Helm名称')}>
              {helm.name}
            </FormPanel.Item>
            <FormPanel.Item text label={t('所属命名空间')}>
              {helm.namespace}
            </FormPanel.Item>
            <FormPanel.Item text label={t('Helm描述')}>
              {helm.chart_metadata.description}
            </FormPanel.Item>
            <FormPanel.Item text label={t('首次部署时间')}>
              {dateFormatter(new Date(helm.info.first_deployed), 'YYYY-MM-DD HH:mm:ss')}
            </FormPanel.Item>
            <FormPanel.Item text label={t('最后一次部署时间')}>
              {dateFormatter(new Date(helm.info.last_deployed), 'YYYY-MM-DD HH:mm:ss')}
            </FormPanel.Item>
          </FormPanel>

          {helm.resources.length > 0 && (
            <Card>
              <Card.Body title={t('资源列表')}>{this._renderTablePanel()}</Card.Body>
            </Card>
          )}
          {helm.configYaml && (
            <Card>
              <Card.Body title={t('自定义参数列表')} className={'height200'}>
                {/* {this._renderValueTablePanel()} */}
                {this.renderValueYaml()}
              </Card.Body>
            </Card>
          )}
          <Card>
            <Card.Body
              title={t('资源状态')}
              operation={
                <React.Fragment>
                  <span
                    className="descript-text"
                    style={{
                      display: 'inline-block',
                      verticalAlign: 'middle',
                      marginRight: '10px',
                      fontSize: '12px'
                    }}
                  >
                    {t('自动刷新')}
                  </span>
                  <Switch value={isRefresh} onChange={checked => this.setRefresh(checked)} className="mr20" />
                </React.Fragment>
              }
            >
              {/* <div style={{ height: 30 }}>
                  <div className="col" style={{ float: 'left' }}>
                    <h3>{t('资源状态')}</h3>
                  </div>

                  <div className="col" style={{ float: 'right' }}>
                    <span
                      className="descript-text"
                      style={{
                        display: 'inline-block',
                        verticalAlign: 'middle',
                        marginRight: '10px',
                        fontSize: '12px'
                      }}
                    >
                      {t('自动刷新')}
                    </span>
                    <Switch value={isRefresh} onChange={checked => this.setRefresh(checked)} className="mr20" />
                  </div>
                </div> */}
              <pre style={{ backgroundColor: 'black', color: 'white', maxHeight: 600, overflowY: 'scroll' }}>
                {helm.info.status.resources}
              </pre>
            </Card.Body>
          </Card>
          {this.state.showYamlDialog && this._renderYamlDialog()}
        </ContentView.Body>
      </ContentView>
    );
  }
}
