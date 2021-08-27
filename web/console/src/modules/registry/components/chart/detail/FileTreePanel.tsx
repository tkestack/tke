/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
let deepEqual = require('deep-equal');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface AppCreateState {
  content?: string;
  defaultContent?: string;
  selectedVersion?: string;
  selectedTreeID?: string;
}

@connect(state => state, mapDispatchToProps)
export class FileTreePanel extends React.Component<RootProps, AppCreateState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      content: '',
      defaultContent: '',
      selectedVersion: '',
      selectedTreeID: 'values.yaml'
    };
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let { chartInfo } = nextProps;
    if (chartInfo.object.data) {
      let node = this.getNode('values.yaml', chartInfo.object.data.fileTree);
      if (node) {
        this.setState({ defaultContent: node.data });
      }
    }
  }

  getNode = (id: string, tree: ChartTreeFile) => {
    if (!tree) {
      return undefined;
    }
    if (tree.fullPath === id) {
      return tree;
    }
    for (let i = 0; i < tree.children.length; i++) {
      let node = this.getNode(id, tree.children[i]);
      if (node) {
        return node;
      }
    }
    return undefined;
  };

  lookupTree = (tree: ChartTreeFile) => {
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
        return this.lookupTree(c);
      });
    } else {
      delete node.children;
    }
    return node;
  };

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
    let treeData = [];
    if (chartInfo.object.data) {
      treeData.push(this.lookupTree(chartInfo.object.data.fileTree));
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
                      this.setState({ selectedVersion: value, content: '', selectedTreeID: 'values.yaml' });
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
                      defaultActiveIds={['values.yaml']}
                      activeIds={[this.state.selectedTreeID]}
                      defaultExpandedIds={
                        chartInfo.object.data && chartInfo.object.data.fileTree
                          ? [chartInfo.object.data.fileTree.fullPath]
                          : ['']
                      }
                      onActive={(selectedIds, context) => {
                        let node = this.getNode(selectedIds[0], chartInfo.object.data.fileTree);
                        if (node) {
                          this.setState({ content: node.data, selectedTreeID: selectedIds[0] });
                        }
                      }}
                    />
                  }
                >
                  <YamlEditorPanel
                    readOnly={true}
                    config={
                      this.state.selectedTreeID === 'values.yaml' ? this.state.defaultContent : this.state.content
                    }
                  />
                </FormPanel.Item>
              </FormPanel>
            </Card.Body>
          </Card>
        </ContentView.Body>
      </ContentView>
    );
  }
}
