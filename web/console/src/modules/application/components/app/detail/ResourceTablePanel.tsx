import * as React from 'react';
import { connect } from 'react-redux';
import { LinkButton, emptyTips } from '../../../../common/components';
import { Table, TableColumn, Text, Modal, Card, Bubble, Icon, ContentView } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { Resource } from '../../../models';
import { RootProps } from '../AppContainer';
import { UnControlled as CodeMirror } from 'react-codemirror2';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface State {
  showYamlDialog?: boolean;
  yaml?: string;
}
@connect(state => state, mapDispatchToProps)
export class ResourceTablePanel extends React.Component<RootProps, State> {
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
    const cancel = () => this.setState({ showYamlDialog: false, yaml: '' });
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

  render() {
    let { actions, resourceList, route } = this.props;
    const columns: TableColumn<Resource>[] = [
      {
        key: 'name',
        header: t('资源名'),
        render: (x: Resource) => (
          <Text parent="div" overflow>
            <a href="javascript:;" onClick={e => {}}>
              {x.metadata.name || '-'}
            </a>
          </Text>
        )
      },
      {
        key: 'namespace',
        header: t('命名空间'),
        render: (x: Resource) => <Text parent="div">{x.metadata.namespace || '-'}</Text>
      },
      {
        key: 'kind',
        header: t('类型'),
        render: (x: Resource) => <Text parent="div">{x.kind || '-'}</Text>
      },
      {
        key: 'operation',
        header: t('操作'),
        render: (x: Resource) => (
          <a href="javascript:void(0)" onClick={e => this.showYaml(x.yaml)}>
            {t('查看YAML')}
          </a>
        )
      }
    ];

    return (
      <ContentView>
        <ContentView.Body>
          <Card>
            <Card.Body>
              <Table
                recordKey={record => {
                  return record.id.toString();
                }}
                records={resourceList.resources}
                columns={columns}
              />
              {this.state.showYamlDialog && this._renderYamlDialog()}
            </Card.Body>
          </Card>
        </ContentView.Body>
      </ContentView>
    );
  }
}
