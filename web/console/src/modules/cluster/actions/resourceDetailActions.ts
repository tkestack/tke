import { ReduxAction, extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { resourceDetailEventActions } from './resourceDetailEventActions';
import { resourceRsActions } from './resourceRsActions';
import { resourcePodActions } from './resourcePodActions';
import { resourcePodLogActions } from './resourcePodLogActions';

const ReduceSecretDataForPsw = (dataInfo: string) => {
  let jsonData = JSON.parse(window.atob(dataInfo)),
    jsonKeys = Object.keys(jsonData)[0];
  return jsonData[jsonKeys]['password'] || '';
};

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

export const resourceDetailActions = {
  /** 获取事件的相关操作 */
  event: resourceDetailEventActions,

  /** 获取修订历史版本的相关操作 */
  rs: resourceRsActions,

  /** 获取pod列表的相关操作 */
  pod: resourcePodActions,

  /** 获取日志的相关操作 */
  log: resourcePodLogActions,

  /** 获取资源，如果deployment的 yaml文件 */
  fetchResourceYaml: generateFetcherActionCreator({
    actionType: ActionType.FetchYaml,
    fetcher: async (getState: GetState, fetchOptions, dispatch) => {
      let { route, subRoot, namespaceSelection } = getState(),
        {
          resourceInfo,
          resourceOption,
          detailResourceOption: { detailResourceInfo, detailResourceSelection }
        } = subRoot,
        { resourceSelection } = resourceOption;
      if (resourceInfo.requestType.useDetailInfo) {
        let response = await WebAPI.fetchResourceYaml(
          detailResourceSelection,
          detailResourceInfo,
          namespaceSelection,
          route.queries['clusterId'],
          +route.queries['rid']
        );
        return response;
      } else {
        let response = await WebAPI.fetchResourceYaml(
          resourceSelection,
          resourceInfo,
          namespaceSelection,
          route.queries['clusterId'],
          +route.queries['rid']
        );
        return response;
      }
    }
  }),

  /** 离开详情页，清除detail的详情 */
  clearDetail: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearResourceDetail
    };
  }
};
