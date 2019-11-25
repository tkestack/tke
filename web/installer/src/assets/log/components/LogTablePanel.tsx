import * as React from 'react';
import { TablePanelGeneric, TablePanelColumnGeneric, TablePanelSmartTip, Pagination } from '@tencent/qcloud-component';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { RootProps } from '../components/LogApp';
import { Log } from '../../common/models';
import { TableLayout } from '../../common/layouts';
import { ParamsPanel } from './ParamsPanel';
import { TipDialog } from '../../common/components';
import { regionMapList } from '../../../../config/region';
import { formatTime } from '../../common/utils';

interface LogTablePanelState {
  isShowRequestParamsDialog?: boolean;
  isShowResponseParamsDialog?: boolean;
  params?: string;
}

// 白名单策略列表的相关配置
type Column = TablePanelColumnGeneric<Log>;
const TablePanel = TablePanelGeneric as new () => TablePanelGeneric<Column, Log>;

export class LogTablePanel extends React.Component<RootProps, LogTablePanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isShowRequestParamsDialog: false,
      isShowResponseParamsDialog: false,
      params: ''
    };
  }

  render() {
    return (
      <TableLayout>
        {this._renderTablePanel()}
        {this._renderPagination()}
        {this._renderRequestParamsDialog()}
        {this._renderResponseParamsDialog()}
      </TableLayout>
    );
  }

  switchRequestParamsShow(isShow: boolean, params: string) {
    let formatParams = params ? JSON.stringify(JSON.parse(params), null, 2) : '';
    this.setState({ isShowRequestParamsDialog: isShow, params: formatParams });
  }

  switchResponseParamsShow(isShow: boolean, params: string) {
    let formatParams = params ? JSON.stringify(JSON.parse(params), null, 2) : '';
    this.setState({ isShowResponseParamsDialog: isShow, params: formatParams });
  }

  private _renderTablePanel() {
    let { actions, logList, logQuery } = this.props;

    const columns: Column[] = [
      {
        id: 'seqId',
        headTitle: 'seqId',
        width: '20%',
        bodyCell: x => (
          <div>
            <p className="text-overflow">{x.seqId}</p>
          </div>
        )
      },
      {
        id: 'time',
        headTitle: '请求时间',
        width: '12%',
        bodyCell: x => (
          <div>
            <p className="text-overflow">{formatTime(x.time)}</p>
          </div>
        )
      },
      {
        id: 'appId',
        headTitle: 'AppID',
        width: '7%',
        bodyCell: x => (
          <div>
            <p className="text-overflow">{x.appId}</p>
          </div>
        )
      },
      {
        id: 'ownerUin',
        headTitle: 'OwnerUin',
        width: '8%',
        bodyCell: x => (
          <div>
            <p className="text-overflow">{x.ownerUin}</p>
          </div>
        )
      },
      {
        id: 'uin',
        headTitle: 'UIN',
        width: '8%',
        bodyCell: x => (
          <div>
            <p className="text-overflow">{x.uin}</p>
          </div>
        )
      },
      {
        id: 'action',
        headTitle: '接口',
        width: '13%',
        bodyCell: x => (
          <div>
            <p className="text-overflow" title={x.action}>
              {x.action}
            </p>
          </div>
        )
      },
      {
        id: 'body',
        headTitle: '请求参数',
        width: '6%',
        bodyCell: x => (
          <div>
            <i
              className="icon-log"
              style={{ cursor: 'pointer' }}
              data-title="查看参数"
              data-logviewer
              onClick={() => {
                this.switchRequestParamsShow(true, x.body);
              }}
            />
          </div>
        )
      },
      {
        id: 'returnCode',
        headTitle: '返回码',
        width: '5%',
        bodyCell: x => (
          <div>
            <p className="text-overflow" title={x.returnCode as string}>
              {x.returnCode}
            </p>
          </div>
        )
      },
      {
        id: 'result',
        headTitle: '返回结果',
        width: '6%',
        bodyCell: x => (
          <div>
            <i
              className="icon-log"
              style={{ cursor: 'pointer' }}
              data-title="查看结果"
              data-logviewer
              onClick={() => {
                this.switchResponseParamsShow(true, x.result);
              }}
            />
          </div>
        )
      },
      {
        id: 'ip',
        headTitle: '请求IP',
        width: '8%',
        bodyCell: x => (
          <div>
            <p className="text-overflow">{x.ip}</p>
          </div>
        )
      },
      {
        id: 'timeCost',
        headTitle: '耗时',
        width: '6%',
        bodyCell: x => (
          <div>
            <p className="text-overflow">{x.timeCost}ms</p>
          </div>
        )
      }
    ];

    const smartTip = TablePanelSmartTip.render({
      fetcher: logList,
      query: logQuery,
      onClearSearch: () => actions.performSearch(''),
      onRetry: () => actions.fetch(),
      enableLoading: true,
      emptyTips:
        logList.fetchState === FetchState.Fetching ? (
          <div>
            <i className="n-loading-icon" />
            &nbsp; <span className="text">加载中...</span>
          </div>
        ) : (
          <div className="text-center">无日志数据</div>
        )
    });

    return (
      <div>
        <TablePanel columns={columns} records={logList.data.records} topTip={smartTip} />
      </div>
    );
  }

  /** 渲染分页组件 */
  private _renderPagination() {
    let { logList, logQuery, actions } = this.props,
      { pageIndex, pageSize } = logQuery.paging;

    return (
      <Pagination
        pageIndex={pageIndex}
        pageSize={pageSize}
        minPageSize={100}
        pageSizeInterval={100}
        maxPageSize={500}
        recordCount={logList.data.recordCount}
        onPagingChange={query => {
          if (query.pageIndex > Math.ceil(logList.data.recordCount / query.pageSize)) {
            query.pageIndex = 1;
          }
          actions.changePaging(query);
        }}
      />
    );
  }

  private _renderRequestParamsDialog() {
    let { isShowRequestParamsDialog, params } = this.state;
    return (
      <TipDialog
        isShow={isShowRequestParamsDialog}
        caption="请求参数"
        body={<ParamsPanel params={params} />}
        cancelAction={() => {
          this.switchRequestParamsShow(false, '');
        }}
        performAction={() => {
          this.switchRequestParamsShow(false, '');
        }}
      />
    );
  }

  private _renderResponseParamsDialog() {
    let { isShowResponseParamsDialog, params } = this.state;
    return (
      <TipDialog
        isShow={isShowResponseParamsDialog}
        caption="返回参数"
        body={<ParamsPanel params={params} />}
        cancelAction={() => {
          this.switchResponseParamsShow(false, '');
        }}
        performAction={() => {
          this.switchResponseParamsShow(false, '');
        }}
      />
    );
  }
}
