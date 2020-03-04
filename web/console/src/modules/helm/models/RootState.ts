import { Namespace } from './Namespace';
import { ResourceFilter, Resource } from './../../common/models/Resource';
import { WorkflowState } from '@tencent/ff-redux';
import { RecordSet } from '@tencent/qcloud-lib';
import { HelmCreation, DetailState, ListState } from './';
import { RouteState } from '../../../../helpers';
import { FetcherState, QueryState } from '@tencent/ff-redux';

export interface RootState {
  /** 路由 */
  route?: RouteState;

  listState?: ListState;

  detailState?: DetailState;

  /**创建Helm所选参数*/
  helmCreation?: HelmCreation;

  /**是否显示提示 */
  isShowTips?: boolean;

  /**业务侧逻辑*/

  /** namespace列表 */
  namespaceList?: FetcherState<RecordSet<Namespace>>;

  /** namespace查询条件 */
  namespaceQuery?: QueryState<ResourceFilter>;

  /** namespace selection */
  namespaceSelection?: string;

  /** namespacesetQuery */
  projectNamespaceQuery?: QueryState<ResourceFilter>;

  /** namespaceset */
  projectNamespaceList?: FetcherState<RecordSet<Resource>>;

  /** projectList */
  projectList?: any[];

  /** projectSelection */
  projectSelection?: string;
  /**end */
}
