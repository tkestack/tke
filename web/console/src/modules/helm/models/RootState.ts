import { FetcherState, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { Resource, ResourceFilter } from '../../common/models/Resource';
import { DetailState, HelmCreation, ListState } from './';
import { Namespace } from './Namespace';

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
