import * as React from 'react';
import { connect } from 'react-redux';
import { TablePanel as CTablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators, OperationState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { ChartGroup } from '../../../models';
import { RootProps } from '../ChartGroupApp';
import { ChartUsageGuideDialog } from '../../ChartUsageGuideDialog';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface ChartUsageGuideDialogState extends RootProps {
  showChartUsageGuideDialog?: boolean;
  chartGroupName?: string;
  registryUrl?: string;
}

@connect(state => state, mapDispatchToProps)
export class TablePanel extends React.Component<RootProps, ChartUsageGuideDialogState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      showChartUsageGuideDialog: false,
      chartGroupName: '',
      registryUrl: ''
    };
  }

  render() {
    let { actions, chartGroupList, chartGroupRepoUpdateWorkflow, userInfo } = this.props;
    const isImported = (x: ChartGroup): boolean => {
      if (x.spec.type === 'Imported') {
        return true;
      }
      return false;
    };

    const columns: TableColumn<ChartGroup>[] = [
      {
        key: 'name',
        header: t('仓库名'),
        render: (x: ChartGroup) => (
          <Text parent="div" overflow>
            <a
              href="javascript:;"
              onClick={e => {
                router.navigate(
                  { sub: 'chartgroup', mode: 'detail' },
                  {
                    name: x.spec.name,
                    cg: x.metadata.name,
                    prj: x.spec.projects && x.spec.projects.length > 0 ? x.spec.projects[0] : ''
                  }
                );
              }}
            >
              {x.spec.name || '-'}
            </a>
            {x.status['phase'] === 'Terminating' && <Icon type="loading" />}
            {chartGroupRepoUpdateWorkflow.operationState === OperationState.Performing &&
              chartGroupRepoUpdateWorkflow.targets &&
              chartGroupRepoUpdateWorkflow.targets.length > 0 &&
              chartGroupRepoUpdateWorkflow.targets[0].metadata.name === x.metadata.name && <Icon type="loading" />}
          </Text>
        )
      },
      {
        key: 'description',
        header: t('描述'),
        render: (x: ChartGroup) => <Text parent="div">{x.spec.description || '-'}</Text>
      },
      {
        key: 'type',
        header: t('类型'),
        render: (x: ChartGroup) => {
          if (!x.spec.type) {
            return <Text parent="div">-</Text>;
          }
          switch (x.spec.type) {
            case 'SelfBuilt':
              return <Text parent="div">{t('自建')}</Text>;
            case 'Imported':
              return <Text parent="div">{t('导入')}</Text>;
            case 'System':
              return <Text parent="div">{t('平台')}</Text>;
            default:
              return <Text parent="div">-</Text>;
          }
        }
      },
      {
        key: 'visibility',
        header: t('权限范围'),
        render: (x: ChartGroup) => {
          if (!x.spec.visibility) {
            return <Text parent="div">-</Text>;
          }
          switch (x.spec.visibility) {
            case 'Public':
              return <Text parent="div">{t('公共')}</Text>;
            case 'User':
              return <Text parent="div">{t('指定用户')}</Text>;
            case 'Project':
              return <Text parent="div">{t('指定业务')}</Text>;
            default:
              return <Text parent="div">-</Text>;
          }
        }
      },
      {
        key: 'chartCount',
        header: t('Chart包数'),
        render: (x: ChartGroup) => (
          <Text parent="div" overflow>
            <span className="text">{x.status.chartCount}</span>
          </Text>
        )
      },
      { key: 'operation', header: t('操作'), render: x => this._renderOperationCell(x, isImported(x)) }
    ];

    return (
      <React.Fragment>
        <CTablePanel
          recordKey={record => {
            return record.metadata.name;
          }}
          columns={columns}
          model={chartGroupList}
          action={actions.chartGroup.list}
          rowDisabled={record => record.status['phase'] === 'Terminating'}
          emptyTips={emptyTips}
          isNeedPagination={true}
          bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
        />
        <ChartUsageGuideDialog
          showDialog={this.state.showChartUsageGuideDialog}
          chartGroupName={this.state.chartGroupName}
          registryUrl={this.state.registryUrl}
          username={userInfo ? userInfo.name : 'tkestack'}
          onClose={() => {
            this.setState({
              showChartUsageGuideDialog: false
            });
          }}
        />
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (chartGroup: ChartGroup, imported: boolean) => {
    let { actions, dockerRegistryUrl } = this.props;
    return (
      <React.Fragment>
        {!imported && (
          <LinkButton
            onClick={() => {
              this.setState({
                showChartUsageGuideDialog: true,
                chartGroupName: chartGroup.spec.name,
                registryUrl: dockerRegistryUrl.data
              });
            }}
          >
            <Trans>上传指引</Trans>
          </LinkButton>
        )}
        {imported && (
          <LinkButton onClick={() => this._updateChartGroup(chartGroup)}>
            <Trans>同步仓库</Trans>
          </LinkButton>
        )}
        <LinkButton onClick={() => this._removeChartGroup(chartGroup)}>
          <Trans>删除</Trans>
        </LinkButton>
      </React.Fragment>
    );
  };

  _removeChartGroup = async (chartGroup: ChartGroup) => {
    let { actions } = this.props;
    const yes = await Modal.confirm({
      message: t('确定删除仓库：') + `${chartGroup.spec.displayName}？`,
      description: <p className="text-danger">{t('删除该仓库后，相关数据将永久删除，请谨慎操作。')}</p>,
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.chartGroup.list.removeChartGroupWorkflow.start([chartGroup], {
        repoType: 'all'
      });
      actions.chartGroup.list.removeChartGroupWorkflow.perform();
    }
  };

  _updateChartGroup = async (chartGroup: ChartGroup) => {
    let { actions } = this.props;
    const yes = await Modal.confirm({
      message: t('确定同步仓库：') + `${chartGroup.spec.displayName}？`,
      description: <p className="text-danger">{t('同步该仓库后，将会更新所有chart的版本信息。')}</p>,
      okText: t('同步'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.chartGroup.list.repoUpdateChartGroupWorkflow.start([chartGroup]);
      actions.chartGroup.list.repoUpdateChartGroupWorkflow.perform();
    }
  };
}
