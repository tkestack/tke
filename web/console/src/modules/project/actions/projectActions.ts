import { detailActions } from './detailActions';
import { ProjectUserMap } from './../models/Project';
import { FFReduxActionName } from './../constants/Config';
import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import {
  createFFListActions,
  deepClone,
  extend,
  generateWorkflowActionCreator,
  isSuccessWorkflow,
  OperationTrigger,
  uuid,
  createFFObjectActions
} from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { initValidator } from '../../common/models/Validation';
import * as ActionType from '../constants/ActionType';
import { initProjectEdition, initProjectResourceLimit, resourceTypeToUnit } from '../constants/Config';
import { Project, ProjectEdition, ProjectFilter, RootState } from '../models';
import { Manager } from '../models/Manager';
import { ProjectResourceLimit } from '../models/Project';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

const FFModelProjectActions = createFFListActions<Project, ProjectFilter>({
  actionName: 'project',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchProjectList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().project;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState(),
      urlParams = router.resolve(route);
    if (urlParams['sub'] === 'createNS') {
      let finder;
      if (route.queries['projectId']) {
        finder = record.data.records.find(item => item.metadata.name === route.queries['projectId']);
        finder = finder ? finder : record.data.recordCount !== 0 ? record.data.records[0] : null;
      }
      if (finder) {
        dispatch(FFModelProjectActions.selects([finder]));
      }
    }
    if (record.data.records.filter(item => item.status.phase !== 'Active').length === 0) {
      dispatch(FFModelProjectActions.clearPolling());
    }
  }
});

const FFObjectProjectUserInfoActions = createFFObjectActions<ProjectUserMap, ProjectFilter>({
  actionName: FFReduxActionName.ProjectUserInfo,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchProjectUserInfo(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().projectUserInfo;
  }
});

const restActions = {
  projectUserInfo: FFObjectProjectUserInfoActions,

  poll: (filter?: ProjectFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { project } = getState();
      dispatch(
        FFModelProjectActions.polling({
          filter: filter || project.query.filter,
          delayTime: 10000
        })
      );
    };
  },
  /** 创建业务 */
  createProject: generateWorkflowActionCreator<ProjectEdition, void>({
    actionType: ActionType.CreateProject,
    workflowStateLocator: (state: RootState) => state.createProject,
    operationExecutor: WebAPI.editProject,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { createProject, route } = getState();
        if (isSuccessWorkflow(createProject)) {
          dispatch(restActions.createProject.reset());
          dispatch(restActions.clearEdition());
          dispatch(projectActions.poll());
          router.navigate({});
        }
      }
    }
  }),

  /** 编辑业务名称 */
  editProjectName: generateWorkflowActionCreator<ProjectEdition, void>({
    actionType: ActionType.EditProjectName,
    workflowStateLocator: (state: RootState) => state.editProjectName,
    operationExecutor: WebAPI.editProject,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { editProjectName, route } = getState(),
          urlParams = router.resolve(route);
        if (isSuccessWorkflow(editProjectName)) {
          dispatch(restActions.editProjectName.reset());
          dispatch(restActions.clearEdition());
          if (urlParams['sub'] === 'detail') {
            dispatch(projectActions.fetchDetail(route.queries['projectId']));
          } else {
            dispatch(projectActions.poll());
          }
        }
      }
    }
  }),

  /** 编辑业务负责人 */
  editProjectManager: generateWorkflowActionCreator<ProjectEdition, void>({
    actionType: ActionType.EditProjectManager,
    workflowStateLocator: (state: RootState) => state.editProjectManager,
    operationExecutor: WebAPI.editProject,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { editProjectManager, route } = getState(),
          urlParams = router.resolve(route);
        if (isSuccessWorkflow(editProjectManager)) {
          dispatch(restActions.editProjectManager.reset());
          dispatch(restActions.clearEdition());
          if (urlParams['sub'] === 'detail') {
            dispatch(projectActions.fetchDetail(route.queries['projectId']));
          } else {
            dispatch(projectActions.poll());
          }
        }
      }
    }
  }),

  /** 编辑业务描述 */
  editProjecResourceLimit: generateWorkflowActionCreator<ProjectEdition, void>({
    actionType: ActionType.EditProjecResourceLimit,
    workflowStateLocator: (state: RootState) => state.editProjecResourceLimit,
    operationExecutor: WebAPI.editProject,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { editProjecResourceLimit, route } = getState(),
          urlParams = router.resolve(route);
        if (isSuccessWorkflow(editProjecResourceLimit)) {
          dispatch(restActions.editProjecResourceLimit.reset());
          dispatch(restActions.clearEdition());
          dispatch(projectActions.fetchDetail(route.queries['projectId']));
        }
      }
    }
  }),

  /** 编辑业务描述 */
  addExistMultiProject: generateWorkflowActionCreator<Project, string>({
    actionType: ActionType.AddExistMultiProject,
    workflowStateLocator: (state: RootState) => state.addExistMultiProject,
    operationExecutor: WebAPI.addExistMultiProject,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { addExistMultiProject, route } = getState(),
          urlParams = router.resolve(route);
        if (isSuccessWorkflow(addExistMultiProject)) {
          dispatch(restActions.addExistMultiProject.reset());
          dispatch(projectActions.clearSelection());
          dispatch(detailActions.project.applyPolling(route.queries['projectId']));
        }
      }
    }
  }),

  /** 编辑业务描述 */
  deleteParentProject: generateWorkflowActionCreator<Project, string>({
    actionType: ActionType.DeleteParentProject,
    workflowStateLocator: (state: RootState) => state.deleteParentProject,
    operationExecutor: WebAPI.deleteParentProject,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { deleteParentProject, route } = getState(),
          urlParams = router.resolve(route);
        if (isSuccessWorkflow(deleteParentProject)) {
          dispatch(restActions.deleteParentProject.reset());
          dispatch(detailActions.project.clearSelection());
          dispatch(detailActions.project.applyPolling(route.queries['projectId']));
        }
      }
    }
  }),

  /** 删除业务 */
  deleteProject: generateWorkflowActionCreator<Project, string>({
    actionType: ActionType.DeleteProject,
    workflowStateLocator: (state: RootState) => state.deleteProject,
    operationExecutor: WebAPI.deleteProject,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { deleteProject, route } = getState();
        if (isSuccessWorkflow(deleteProject)) {
          dispatch(restActions.deleteProject.reset());
          dispatch(projectActions.poll());
        }
      }
    }
  }),

  /**拉取业务详情 */
  fetchDetail: (projectId?: string) => {
    return async dispatch => {
      let project = await WebAPI.fetchProjectDetail(projectId);
      dispatch({
        type: ActionType.ProjectDetail,
        payload: project
      });
    };
  },

  /** --begin编辑action */
  inputProjectName: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { displayName: value })
      });
    };
  },

  inputParentPorject: (projectId: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let clusters = getState().projectEdition.clusters;
      clusters = clusters.map(item => Object.assign({}, item, { name: '' }));
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { parentProject: projectId, clusters })
      });
    };
  },

  _validateDisplayName(name: string) {
    let status = 0,
      message = '';

    // 验证内存限制
    if (name === '') {
      status = 2;
      message = t('业务名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('业务名称长度不能超过63个字符');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateDisplayName(name: string) {
    return async (dispatch, getState: GetState) => {
      let result = projectActions._validateDisplayName(name);
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { v_displayName: result })
      });
    };
  },

  _validateProjection(projectEdition: ProjectEdition) {
    let ok = true && projectEdition.displayName !== '';
    projectEdition.clusters.forEach(cluster => {
      if (cluster.name === '') {
        ok = false;
      }
    });
    if (!projectEdition.members || projectEdition.members.length === 0) {
      ok = false;
    }
    return ok;
  },

  validateProjection() {
    return async (dispatch, getState: GetState) => {
      let { projectEdition } = getState();
      projectEdition.clusters.forEach((cluster, index) => {
        dispatch(projectActions.validateClustersName(index));
      });
      dispatch(projectActions.validateDisplayName(projectEdition.displayName));
    };
  },

  selectManager: (values: Array<Manager>) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { members: values })
      });
    };
  },

  initEdition: (project: Project) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        manager: { list }
      } = getState();
      let {
        spec: { clusters }
      } = project;
      let clusterKeys = clusters ? Object.keys(clusters) : [];
      let clustersInfo = clusterKeys.map(cluster => {
        let hard = clusters[cluster].hard;
        let hardInfo = hard
          ? Object.keys(hard).map(key => {
              let value = hard[key];
              /**CPU类 */
              if (resourceTypeToUnit[key] === '核' || resourceTypeToUnit[key] === '个') {
                value = valueLabels1000(value, K8SUNIT.unit);
              } else if (resourceTypeToUnit[key] === 'MiB') {
                value = valueLabels1024(value, K8SUNIT.Mi);
              }
              return Object.assign({}, initProjectResourceLimit, { type: key, id: uuid(), value });
            })
          : [];
        return {
          name: cluster,
          v_name: initValidator,
          resourceLimits: hardInfo
        };
      });
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: {
          id: project.id,
          resourceVersion: project.metadata.resourceVersion,
          displayName: project.spec.displayName,
          members: project.spec.members.map(item => {
            let finder = list.data.records.find(manager => manager.name === item);
            if (finder) {
              return finder;
            } else {
              return {
                name: item,
                displayName: '用户不存在'
              };
            }
          }),
          parentProject: project.spec.parentProjectName ? project.spec.parentProjectName : '',
          clusters: clustersInfo,
          status: project.status
        }
      });
    };
  },

  //更新project集群
  updateClusters: (index: number, clusterSelection: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { projectEdition } = getState(),
        { clusters } = projectEdition;
      let newClusters = deepClone(clusters);
      if (newClusters[index]) {
        newClusters[index].name = clusterSelection;
      }

      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { clusters: newClusters })
      });
    };
  },

  //删除project集群
  deleteClusters: (index: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { projectEdition } = getState(),
        { clusters } = projectEdition;
      let newClusters = deepClone(clusters);
      if (newClusters[index]) {
        newClusters.splice(index, 1);
      }

      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { clusters: newClusters })
      });
    };
  },

  //添加project集群
  addClusters: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { projectEdition } = getState(),
        { clusters } = projectEdition;
      let newClusters = deepClone(clusters);
      newClusters.push({
        name: '',
        v_name: initValidator,
        resourceLimits: []
      });

      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { clusters: newClusters })
      });
    };
  },

  updateClustersLimit: (index: number, resourceLimits: ProjectResourceLimit[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { projectEdition } = getState(),
        { clusters } = projectEdition;
      let newClusters = deepClone(clusters);
      if (newClusters[index]) {
        newClusters[index] = Object.assign({}, newClusters[index], { resourceLimits });
      }
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { clusters: newClusters })
      });
    };
  },

  validateClustersName: (index: number) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        projectEdition: { clusters }
      } = getState();
      let newClusters = deepClone(clusters);
      if (newClusters[index]) {
        let filter = newClusters.filter(item => item.name === newClusters[index].name);
        if (!newClusters[index].name) {
          newClusters[index].v_name = {
            status: 2,
            message: '集群名称不能为空'
          };
        } else if (filter.length > 1) {
          newClusters[index].v_name = {
            status: 2,
            message: '限制不能重复填写'
          };
        } else {
          newClusters[index].v_name = {
            status: 1,
            message: ''
          };
        }
      }
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: Object.assign({}, getState().projectEdition, { clusters: newClusters })
      });
    };
  },

  clearEdition: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateProjectEdition,
        payload: initProjectEdition
      });
    };
  }
  /** --end编辑action */
};

export const projectActions = extend({}, FFModelProjectActions, restActions);
