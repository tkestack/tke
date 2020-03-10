import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import {
    createFFListActions, extend, generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger,
    uuid
} from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import * as ActionType from '../constants/ActionType';
import {
    initNamespaceEdition, initProjectResourceLimit, resourceTypeToUnit
} from '../constants/Config';
import {
    Namespace, NamespaceEdition, NamespaceFilter, NamespaceOperator, RootState
} from '../models';
import { ProjectResourceLimit } from '../models/Project';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

const FFModelNamespaceActions = createFFListActions<Namespace, NamespaceFilter>({
  actionName: 'namespace',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchNamespaceList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().namespace;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (
      record.data.records.filter(item => item.status.phase !== 'Available' && item.status.phase !== 'Failed').length ===
      0
    ) {
      dispatch(FFModelNamespaceActions.clearPolling());
    }
  }
});

const restActions = {
  poll: (filter?: NamespaceFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { namespace } = getState();
      dispatch(
        FFModelNamespaceActions.polling({
          filter: filter || namespace.query.filter,
          delayTime: 8000
        })
      );
    };
  },

  /** 创建Namespace */
  createNamespace: generateWorkflowActionCreator<NamespaceEdition, NamespaceOperator>({
    actionType: ActionType.CreateNamespace,
    workflowStateLocator: (state: RootState) => state.createNamespace,
    operationExecutor: WebAPI.editNamespace,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { createNamespace, route } = getState();
        if (isSuccessWorkflow(createNamespace)) {
          router.navigate({ sub: 'detail', tab: 'namespace' }, route.queries);
          dispatch(restActions.createNamespace.reset());
          dispatch(restActions.clearEdition());
        }
      }
    }
  }),

  /** 创建Namespace */
  editNamespaceResourceLimit: generateWorkflowActionCreator<NamespaceEdition, NamespaceOperator>({
    actionType: ActionType.EditNamespaceResourceLimit,
    workflowStateLocator: (state: RootState) => state.editNamespaceResourceLimit,
    operationExecutor: WebAPI.editNamespace,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { editNamespaceResourceLimit, route } = getState();
        if (isSuccessWorkflow(editNamespaceResourceLimit)) {
          dispatch(namespaceActions.poll());
          dispatch(restActions.editNamespaceResourceLimit.reset());
          dispatch(restActions.clearEdition());
        }
      }
    }
  }),

  /** 删除Namespace */
  deleteNamespace: generateWorkflowActionCreator<Namespace, NamespaceOperator>({
    actionType: ActionType.DeleteNamespace,
    workflowStateLocator: (state: RootState) => state.deleteNamespace,
    operationExecutor: WebAPI.deleteNamespace,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { deleteNamespace, route } = getState();
        if (isSuccessWorkflow(deleteNamespace)) {
          dispatch(restActions.deleteNamespace.reset());
          dispatch(namespaceActions.poll());
        }
      }
    }
  }),

  select: (namespaces: Namespace[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(FFModelNamespaceActions.selects(namespaces));
    };
  },

  selectCluster: (value: string | number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateNamespaceEdition,
        payload: Object.assign({}, getState().namespaceEdition, { clusterName: value })
      });
    };
  },

  initNamespaceEdition: (namespace: Namespace) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let hardInfo = namespace.spec.hard
        ? Object.keys(namespace.spec.hard).map(key => {
            let value = namespace.spec.hard[key];
            /**CPU类 */
            /**CPU类 */
            if (resourceTypeToUnit[key] === '核' || resourceTypeToUnit[key] === '个') {
              value = parseFloat(valueLabels1000(value, K8SUNIT.unit)) + '';
            } else if (resourceTypeToUnit[key] === 'MiB') {
              value = parseFloat(valueLabels1024(value, K8SUNIT.Mi)) + '';
            }
            /**个数不需要转化 */
            return Object.assign({}, initProjectResourceLimit, { type: key, id: uuid(), value });
          })
        : [];
      dispatch({
        type: ActionType.UpdateNamespaceEdition,
        payload: {
          id: namespace.id,
          resourceVersion: namespace.metadata.resourceVersion,
          namespaceName: namespace.spec.namespace,
          clusterName: namespace.spec.clusterName,
          resourceLimits: hardInfo,
          status: namespace.status
        }
      });
    };
  },

  inputNamespaceName: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateNamespaceEdition,
        payload: Object.assign({}, getState().namespaceEdition, { namespaceName: value })
      });
    };
  },

  updateNamespaceResourceLimit: (resourceLimits: ProjectResourceLimit[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateNamespaceEdition,
        payload: Object.assign({}, getState().namespaceEdition, { resourceLimits: resourceLimits })
      });
    };
  },

  validateNamespaceName() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        namespaceEdition: { namespaceName }
      } = getState();
      let result = namespaceActions._validateNamespaceName(namespaceName);
      dispatch({
        type: ActionType.UpdateNamespaceEdition,
        payload: Object.assign({}, getState().namespaceEdition, { v_namespaceName: result })
      });
    };
  },

  validateClusterName() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        namespaceEdition: { clusterName }
      } = getState();
      let result;
      if (clusterName === '') {
        result = {
          status: 2,
          message: '集群不能为空'
        };
      } else {
        result = {
          status: 1,
          message: ''
        };
      }
      dispatch({
        type: ActionType.UpdateNamespaceEdition,
        payload: Object.assign({}, getState().namespaceEdition, { v_clusterName: result })
      });
    };
  },
  /**
   * 校验namespace名称是否正确
   */
  _validateNamespaceName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    // 验证ingress名称
    if (!name) {
      status = 2;
      message = t('Namespace名称不能为空');
    } else if (name.length > 48) {
      status = 2;
      message = t('Namespace名称不能超过48个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('Namespace名称格式不正确');
    } else if (name.startsWith('kube-')) {
      status = 2;
      message = t('Namespace名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  _validateNamespaceEdition(namespaceEdtion: NamespaceEdition) {
    let ok = true && namespaceEdtion.clusterName !== '';
    ok = ok && namespaceActions._validateNamespaceName(namespaceEdtion.namespaceName).status === 1;
    return ok;
  },

  validateNamespaceEdition() {
    return async (dispatch, getState: GetState) => {
      dispatch(namespaceActions.validateNamespaceName());
      dispatch(namespaceActions.validateClusterName());
    };
  },
  clearEdition: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateNamespaceEdition,
        payload: initNamespaceEdition
      });
    };
  }
};

export const namespaceActions = extend(FFModelNamespaceActions, restActions);
