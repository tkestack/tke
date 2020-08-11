import * as React from 'react';
import { connect } from 'react-redux';
import { LinkButton, emptyTips } from '../../../../common/components';
import { Table, TableColumn, Text, Modal, Card, Bubble, Icon, ContentView } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { History } from '../../../models';
import { RootProps } from '../AppContainer';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import { dateFormat } from '../../../../../../helpers/dateUtil';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HistoryTablePanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, historyList, route } = this.props;
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
              <Table
                recordKey={record => {
                  return record.id.toString();
                }}
                records={historyList.histories}
                columns={columns}
              />
            </Card.Body>
          </Card>
        </ContentView.Body>
      </ContentView>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (app: History) => {
    if (app.status === 'deployed') {
      return false;
    }
    return (
      <React.Fragment>
        <LinkButton onClick={() => this._rollbackApp(app)}>{t('回滚')}</LinkButton>
      </React.Fragment>
    );
  };

  _rollbackApp = async (app: History) => {
    let { actions } = this.props;
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
