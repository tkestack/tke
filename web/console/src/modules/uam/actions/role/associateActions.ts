import { extend, generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow, createFFListActions } from '@tencent/ff-redux';
import { RootState, RoleAssociation, RolePlain, RoleFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import {
  initRoleAssociationState
} from '../../constants/initState';

type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchRoleActions = createFFListActions<RolePlain, RoleFilter>({
  actionName: ActionTypes.RolePlainList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchRolePlainList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().rolePlainList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {

  }
});

/**
 * 列表操作
 */
const fetchRoleAssociatedActions = createFFListActions<RolePlain, RoleFilter>({
  actionName: ActionTypes.RoleAssociatedList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchRoleAssociatedList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().roleAssociatedList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    // 拉取关联列表之后，更新关联角色面板WorkflowDialog会用到的RoleAssociation状态数据
    dispatch({
      type: ActionTypes.UpdateRoleAssociation,
      payload: Object.assign({}, getState().roleAssociation, {
        roles: record.data.records,
        originRoles: record.data.records,
        addRoles: [],
        removeRoles: [],
      })
    });
  }
});

/**
 * 关联操作
 */
const associateRoleWorkflow = generateWorkflowActionCreator<RoleAssociation, RoleFilter>({
  actionType: ActionTypes.AssociateRole,
  workflowStateLocator: (state: RootState) => state.associateRoleWorkflow,
  operationExecutor: WebAPI.associateRole,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      let { associateRoleWorkflow, roleAssociation, roleFilter } = getState();
      if (isSuccessWorkflow(associateRoleWorkflow)) {
        /** 回调函数 */
        roleFilter.callback && roleFilter.callback();
        /** 开始解绑工作流 */
        if (roleAssociation.removeRoles.length > 0) {
          dispatch(associateActions.disassociateRoleWorkflow.start([roleAssociation], roleFilter));
          dispatch(associateActions.disassociateRoleWorkflow.perform());
        } else {
          /** 重新加载关联数据 */
          dispatch(associateActions.roleAssociatedList.applyFilter(roleFilter));
        }
      }
      /** 结束工作流 */
      dispatch(associateActions.associateRoleWorkflow.reset());
    }
  }
});

/**
 * 解绑操作
 */
const disassociateRoleWorkflow = generateWorkflowActionCreator<RoleAssociation, RoleFilter>({
  actionType: ActionTypes.DisassociateRole,
  workflowStateLocator: (state: RootState) => state.disassociateRoleWorkflow,
  operationExecutor: WebAPI.disassociateRole,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      /** 解绑之后重新加载关联数据 */
      let { disassociateRoleWorkflow, roleFilter } = getState();
      if (isSuccessWorkflow(disassociateRoleWorkflow)) {
        /** 回调函数 */
        roleFilter.callback && roleFilter.callback();
        dispatch(associateActions.roleAssociatedList.applyFilter(roleFilter));
      }
      /** 结束工作流 */
      dispatch(associateActions.disassociateRoleWorkflow.reset());
    }
  }
});

const restActions = {
  associateRoleWorkflow,
  disassociateRoleWorkflow,

  /** 设置角色过滤器 */
  setupRoleFilter: (filter: RoleFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdateRoleFilter,
        payload: Object.assign({}, getState().roleFilter, filter),
      });
    };
  },

  /** 选中角色，根据原始数据计算将添加的角色和将删除的角色 */
  selectRole: (roles: RolePlain[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      /** 选中关联角色，则更新关联角色面板WorkflowDialog会用到的RoleAssociation状态数据 */
      /** 比对计算出新增和删除的角色，originRoles是指原先绑定的角色 */
      const { originRoles } = getState().roleAssociation;
      const getDifferenceSet = (arr1, arr2) => {
        let a1 = arr1.map(JSON.stringify);
        let a2 = arr2.map(JSON.stringify);
        return a1.concat(a2).filter(v => !a1.includes(v) || !a2.includes(v)).map(JSON.parse);
      };
      let allRoles = roles.concat(originRoles);
      let removeRoles = getDifferenceSet(roles, allRoles);
      let addRoles = getDifferenceSet(originRoles, allRoles);
      dispatch({
        type: ActionTypes.UpdateRoleAssociation,
        payload: Object.assign({}, getState().roleAssociation, {
          roles: roles,
          addRoles: addRoles,
          removeRoles: removeRoles
        })
      });
    };
  },

  /** 清除角色关联状态数据 */
  clearRoleAssociation: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdateRoleAssociation,
        payload: initRoleAssociationState
      });
    };
  }
};

export const associateActions = extend({},
  {
    roleList: fetchRoleActions,
    roleAssociatedList: fetchRoleAssociatedActions
  },
  restActions);
