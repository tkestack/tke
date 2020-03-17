import { extend, generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow, createFFListActions } from '@tencent/ff-redux';
import { RootState, CommonUserAssociation, UserPlain, CommonUserFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import {
  initCommonUserAssociationState
} from '../../constants/initState';

type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchUserActions = createFFListActions<UserPlain, CommonUserFilter>({
  actionName: ActionTypes.UserPlainList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchCommonUserList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().userPlainList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {

  }
});

/**
 * 列表操作
 */
const fetchUserAssociatedActions = createFFListActions<UserPlain, CommonUserFilter>({
  actionName: ActionTypes.CommonUserAssociatedList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchCommonUserAssociatedList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().commonUserAssociatedList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    // 拉取关联列表之后，更新关联用户面板WorkflowDialog会用到的CommonUserAssociation状态数据
    dispatch({
      type: ActionTypes.UpdateCommonUserAssociation,
      payload: Object.assign({}, getState().commonUserAssociation, {
        users: record.data.records,
        originUsers: record.data.records,
        addUsers: [],
        removeUsers: [],
      })
    });
  }
});

/**
 * 关联用户操作
 */
const associateUserWorkflow = generateWorkflowActionCreator<CommonUserAssociation, CommonUserFilter>({
  actionType: ActionTypes.CommonAssociateUser,
  workflowStateLocator: (state: RootState) => state.commonAssociateUserWorkflow,
  operationExecutor: WebAPI.commonAssociateUser,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      let { commonAssociateUserWorkflow, commonUserAssociation, commonUserFilter } = getState();
      if (isSuccessWorkflow(commonAssociateUserWorkflow)) {
        /** 回调函数 */
        commonUserFilter.callback && commonUserFilter.callback();
        /** 开始解绑工作流 */
        if (commonUserAssociation.removeUsers.length > 0) {
          dispatch(associateActions.disassociateUserWorkflow.start([commonUserAssociation], commonUserFilter));
          dispatch(associateActions.disassociateUserWorkflow.perform());
        } else {
          /** 重新加载关联数据 */
          dispatch(associateActions.userAssociatedList.applyFilter(commonUserFilter));
        }
      }
      /** 结束工作流 */
      dispatch(associateActions.associateUserWorkflow.reset());
    }
  }
});

/**
 * 解绑用户操作
 */
const disassociateUserWorkflow = generateWorkflowActionCreator<CommonUserAssociation, CommonUserFilter>({
  actionType: ActionTypes.CommonDisassociateUser,
  workflowStateLocator: (state: RootState) => state.commonDisassociateUserWorkflow,
  operationExecutor: WebAPI.commonDisassociateUser,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      /** 解绑之后重新加载关联数据 */
      let { commonDisassociateUserWorkflow, commonUserFilter } = getState();
      if (isSuccessWorkflow(commonDisassociateUserWorkflow)) {
        /** 回调函数 */
        commonUserFilter.callback && commonUserFilter.callback();
        dispatch(associateActions.userAssociatedList.applyFilter(commonUserFilter));
      }
      /** 结束工作流 */
      dispatch(associateActions.disassociateUserWorkflow.reset());
    }
  }
});

const restActions = {
  associateUserWorkflow,
  disassociateUserWorkflow,

  /** 设置用户过滤器 */
  setupUserFilter: (filter: CommonUserFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdateCommonUserFilter,
        payload: Object.assign({}, getState().commonUserFilter, filter),
      });
    };
  },

  /** 选中用户，根据原始数据计算将添加的用户和将删除的用户 */
  selectUser: (users: UserPlain[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      /** 选中关联用户，则更新关联用户面板WorkflowDialog会用到的CommonUserAssociation状态数据 */
      /** 比对计算出新增和删除的用户，originUsers是指原先绑定的用户 */
      const { originUsers } = getState().commonUserAssociation;
      const getDifferenceSet = (arr1, arr2) => {
        let a1 = arr1.map(JSON.stringify);
        let a2 = arr2.map(JSON.stringify);
        return a1.concat(a2).filter(v => !a1.includes(v) || !a2.includes(v)).map(JSON.parse);
      };
      let allUsers = users.concat(originUsers);
      let removeUsers = getDifferenceSet(users, allUsers);
      let addUsers = getDifferenceSet(originUsers, allUsers);
      dispatch({
        type: ActionTypes.UpdateCommonUserAssociation,
        payload: Object.assign({}, getState().commonUserAssociation, {
          users: users,
          addUsers: addUsers,
          removeUsers: removeUsers
        })
      });
    };
  },

  /** 清除用户关联状态数据 */
  clearUserAssociation: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdateCommonUserAssociation,
        payload: initCommonUserAssociationState
      });
    };
  }
};

export const associateActions = extend({},
  {
    userList: fetchUserActions,
    userAssociatedList: fetchUserAssociatedActions
  },
  restActions);
