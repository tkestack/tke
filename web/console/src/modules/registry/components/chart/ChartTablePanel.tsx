import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Justify, Table, TableColumn, Text } from '@tencent/tea-component';
import { expandable } from '@tea/component/table/addons/expandable';

import { dateFormatter } from '../../../../../helpers';
import { Clip, GridTable, TipDialog, WorkflowDialog } from '../../../common/components';
import { DialogBodyLayout } from '../../../common/layouts';
import { ChartIns } from '../../models';
import { router } from '../../router';
import { RootProps } from '../RegistryApp';

export class ChartTablePanel extends React.Component<RootProps, any> {
  state = {
    showUsageGuideline: false
  };

  componentDidMount() {
    this.props.actions.chartIns.applyFilter({
      chartgroup: this.props.route.queries['cg']
    });
  }

  render() {
    return (
      <React.Fragment>
        <div className="tc-action-grid">
          <Justify
            left={
              <React.Fragment>
                <Button
                  type="primary"
                  onClick={() => {
                    this.setState({ showUsageGuideline: true });
                  }}
                >
                  {t('Chart上传指引')}
                </Button>
              </React.Fragment>
            }
          ></Justify>
        </div>
        {this._renderTablePanel()}
        {this._renderUsageGuideDialog()}
      </React.Fragment>
    );
  }

  private _renderTablePanel() {
    const columns: TableColumn<ChartIns>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.name}</span>
          </Text>
        )
      },
      {
        key: 'desc',
        header: t('描述'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.displayName}</span>
          </Text>
        )
      },
      {
        key: 'visibility',
        header: t('权限类型'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.spec.visibility === 'Public' ? t('公有') : t('私有')}</span>
          </Text>
        )
      },
      {
        key: 'pullCount',
        header: t('下载次数'),
        render: (x: ChartIns) => (
          <Text parent="div" overflow>
            <span className="text">{x.status.pullCount}</span>
          </Text>
        )
      },
      {
        key: 'settings',
        header: '操作',
        width: 100,
        render: () => '-'
      }
    ];

    return (
      <GridTable
        columns={columns}
        emptyTips={<div className="text-center">{t('Chart列表为空')}</div>}
        listModel={{
          list: this.props.chartIns.list,
          query: Object.assign({}, this.props.chartIns.query, {
            filter: {
              chartgroup: this.props.route.queries['cg']
            }
          })
        }}
        actionOptions={this.props.actions.chartIns}
      />
    );
  }

  private _renderUsageGuideDialog() {
    return (
      <TipDialog
        isShow={this.state.showUsageGuideline}
        width={680}
        caption={t('Chart 上传指引')}
        cancelAction={() => this.setState({ showUsageGuideline: false })}
        performAction={() => this.setState({ showUsageGuideline: false })}
      >
        <div className="mirroring-box" style={{ marginTop: '0px' }}>
          <ul className="mirroring-upload-list">
            <li>
              <p>
                <strong>
                  <Trans>登录</Trans> TKEStack Docker Registry
                </strong>
              </p>
              <code>
                <Clip target="#loginDocker" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="loginDocker">{`sudo docker login -u tkestack -p [访问凭证]`}</p>
              </code>
            </li>
          </ul>
        </div>
      </TipDialog>
    );
  }
}
