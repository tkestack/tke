import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../ChartApp';
import { FormPanel } from '@tencent/ff-component';
import { Button, Tabs, TabPanel, Card, Bubble, Icon, ContentView, Drawer, Tree, Row, Col } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { YamlEditorPanel } from '../../../../common/components';
import { ChartTreeFile } from '@src/modules/registry/models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface AppCreateState {
  content?: string;
  selectedVersion?: string;
}

@connect(state => state, mapDispatchToProps)
export class FileTreePanel extends React.Component<RootProps, AppCreateState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      content: '',
      selectedVersion: ''
    };
  }

  render() {
    let { actions, chartEditor, chartInfo, route } = this.props;

    const versionOptions = chartEditor
      ? chartEditor.status.versions.map(v => {
          return {
            text: v.version,
            value: v.version
          };
        })
      : [];

    const getNode = (id: string, tree: ChartTreeFile) => {
      if (!tree) {
        return undefined;
      }
      if (tree.fullPath === id) {
        return tree;
      }
      for (let i = 0; i < tree.children.length; i++) {
        let node = getNode(id, tree.children[i]);
        if (node) {
          return node;
        }
      }
      return undefined;
    };

    const lookupTree = (tree: ChartTreeFile) => {
      if (!tree) {
        return {};
      }
      let node = {
        id: tree.fullPath,
        content: tree.name,
        expandable: tree.children.length > 0,
        selectable: false,
        children: []
      };
      if (tree.children.length > 0) {
        node.children = tree.children.map(c => {
          return lookupTree(c);
        });
      } else {
        delete node.children;
      }
      return node;
    };
    let treeData = [];
    if (chartInfo.object.data) {
      treeData.push(lookupTree(chartInfo.object.data.fileTree));
    }
    return (
      <ContentView>
        <ContentView.Body>
          <Card>
            <Card.Body>
              <FormPanel isNeedCard={false}>
                <FormPanel.Item
                  label={t('Chart版本')}
                  select={{
                    value: this.state.selectedVersion
                      ? this.state.selectedVersion
                      : (chartEditor && chartEditor.selectedVersion && chartEditor.selectedVersion.version) || '',
                    valueField: 'value',
                    displayField: 'text',
                    options: versionOptions,
                    onChange: value => {
                      this.setState({ selectedVersion: value, content: '' });
                      //加载文件
                      if (chartEditor) {
                        actions.chart.detail.chartInfo.applyFilter({
                          cluster: '',
                          namespace: '',
                          metadata: {
                            namespace: chartEditor.metadata.namespace,
                            name: chartEditor.metadata.name
                          },
                          chartVersion: value,
                          projectID: route.queries['prj']
                        });
                      }
                    }
                  }}
                />
                <FormPanel.Item
                  label={
                    <Tree
                      selectable
                      data={treeData}
                      activable={true}
                      defaultExpandedIds={
                        chartInfo.object.data && chartInfo.object.data.fileTree
                          ? [chartInfo.object.data.fileTree.fullPath]
                          : ['']
                      }
                      onActive={(selectedIds, context) => {
                        let node = getNode(selectedIds[0], chartInfo.object.data.fileTree);
                        console.log(node);
                        if (node) {
                          this.setState({ content: node.data });
                        }
                      }}
                    />
                  }
                >
                  <YamlEditorPanel config={this.state.content} />
                </FormPanel.Item>
              </FormPanel>
            </Card.Body>
          </Card>
        </ContentView.Body>
      </ContentView>
    );
  }
}
