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
import { LinkButton, YamlEditorPanel } from '../../../../common/components';
import { Table, TableColumn, Text, Modal, Card, Justify, Button, ContentView } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { History } from '../../../models';
import { RootProps } from '../AppContainer';
import { selectable } from '@tea/component/table/addons/selectable';
import { dateFormat } from '../../../../../../helpers/dateUtil';

const tips = seajs.require('tips');
const jsDiff = require('diff');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface State {
  showYamlDialog?: boolean;
  yaml?: string;
  revisions?: string[];
  selectedHistories?: History[];
}
@connect(state => state, mapDispatchToProps)
export class HistoryTablePanel extends React.Component<RootProps, State> {
  state = {
    showYamlDialog: false,
    yaml: '',
    revisions: [],
    selectedHistories: []
  };
  showYaml(yaml, revisions) {
    this.setState({
      showYamlDialog: true,
      yaml,
      revisions
    });
  }

  _renderYamlDialog() {
    const cancel = () => this.setState({ showYamlDialog: false, yaml: '' });
    const title =
      this.state.revisions && this.state.revisions.length === 1
        ? '版本: ' + this.state.revisions[0]
        : '版本比对: ' + this.state.revisions[0] + '和' + this.state.revisions[1];
    return (
      <Modal visible={true} caption={t(title)} onClose={cancel} size={700} disableEscape={true}>
        <Modal.Body>
          <YamlEditorPanel readOnly={true} config={this.state.yaml} />
        </Modal.Body>
      </Modal>
    );
  }

  render() {
    const { actions, historyList, route } = this.props;
    const columns: TableColumn<History>[] = [
      {
        key: 'name',
        header: t('应用名'),
        render: (x: History) => (
          <Text parent="div" overflow>
            {(x.involvedObject && x.involvedObject.spec && x.involvedObject.spec.name) || '-'}
          </Text>
        )
      },
      {
        key: 'chart',
        header: t('Chart'),
        render: (x: History) => <Text parent="div">{x.chart || '-'}</Text>
      },
      {
        key: 'revision',
        header: t('版本'),
        render: (x: History) => <Text parent="div">{x.revision || '-'}</Text>
      },
      {
        key: 'status',
        header: t('状态'),
        render: (x: History) => {
          return x.status === 'deployed' ? (
            <Text parent="div" className="tea-text-success tea-align-middle">
              {t('运行中')}
            </Text>
          ) : (
            <Text parent="div">{t('已废弃')}</Text>
          );
        }
      },
      {
        key: 'description',
        header: t('描述'),
        render: (x: History) => <Text parent="div">{x.description || '-'}</Text>
      },
      {
        key: 'updated',
        header: t('更新时间'),
        render: (x: History) => (
          <Text parent="div">{x.updated ? dateFormat(new Date(x.updated), 'yyyy-MM-dd hh:mm:ss') : '-'}</Text>
        )
      },
      { key: 'operation', header: t('操作'), render: app => this._renderOperationCell(app) }
    ];

    return (
      <ContentView>
        <ContentView.Body>
          <Card>
            <Card.Body>
              <Table.ActionPanel>
                <Justify
                  left={
                    <Button
                      type="primary"
                      onClick={e => {
                        e.preventDefault();
                        if (!this.state.selectedHistories || this.state.selectedHistories.length !== 2) {
                          tips.error('请选择两个比对版本', 2000);
                          return;
                        }
                        //context的数目会影响比对行数的显示
                        const diff = jsDiff.createTwoFilesPatch(
                          t('版本: ') + this.state.selectedHistories[1].revision,
                          t('版本: ') + this.state.selectedHistories[0].revision,
                          this.state.selectedHistories[1].manifest,
                          this.state.selectedHistories[0].manifest,
                          '',
                          '',
                          { context: 0 }
                        );
                        this.showYaml(diff, [
                          this.state.selectedHistories[0].revision,
                          this.state.selectedHistories[1].revision
                        ]);
                      }}
                    >
                      {t('参数比对')}
                    </Button>
                  }
                />
              </Table.ActionPanel>
              <Table
                recordKey={record => {
                  return record.id.toString();
                }}
                records={historyList.histories}
                columns={columns}
                addons={[
                  selectable({
                    value: this.state.selectedHistories.map(item => item.id as string),
                    onChange: keys => {
                      this.setState({
                        selectedHistories: historyList.histories.filter(item => keys.indexOf(item.id as string) > -1)
                      });
                    }
                  })
                ]}
              />
            </Card.Body>
          </Card>
          {this.state.showYamlDialog && this._renderYamlDialog()}
        </ContentView.Body>
      </ContentView>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (app: History) => {
    if (app.status === 'deployed') {
      return (
        <React.Fragment>
          <LinkButton
            onClick={() => {
              this.showYaml(app.manifest, [app.revision]);
            }}
          >
            {t('参数')}
          </LinkButton>
        </React.Fragment>
      );
    }
    return (
      <React.Fragment>
        <LinkButton
          onClick={() => {
            this.showYaml(app.manifest, [app.revision]);
          }}
        >
          {t('参数')}
        </LinkButton>
        <LinkButton onClick={() => this._rollbackApp(app)}>{t('回滚')}</LinkButton>
      </React.Fragment>
    );
  };

  _rollbackApp = async (app: History) => {
    const { actions } = this.props;
    const yes = await Modal.confirm({
      message:
        t('确定回滚应用：') +
        `${(app.involvedObject && app.involvedObject.spec && app.involvedObject.spec.name) || '-'}` +
        t('到版本：') +
        `${app.revision}` +
        '？',
      description: <p className="text-danger">{t('请谨慎操作。')}</p>,
      okText: t('回滚'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.app.history.rollbackAppWorkflow.start([app]);
      actions.app.history.rollbackAppWorkflow.perform();
    }
  };
}
