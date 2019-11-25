import { QueryState } from '@tencent/qcloud-redux-query';
import { Log, LogFilter } from './../../common/models';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { RecordSet } from '../../common/models';
import { RouteState } from '../../../../helpers';

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 控制台日志列表 */
  logList?: FetcherState<RecordSet<Log>>;

  /** 控制台日志列表查询 */
  logQuery?: QueryState<LogFilter>;

  /** 控制台日志选择 */
  logSelection?: Array<Log>;
}
