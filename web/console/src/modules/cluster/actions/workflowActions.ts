import {
  extend,
  generateWorkflowActionCreator,
  isSuccessWorkflow,
  OperationHooks,
  OperationTrigger
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { CreateIC, CreateResource, DifferentInterfaceResourceOperation, RootState } from '../models';
import { AllocationRatioEdition } from '../models/AllocationRatioEdition';
import { Computer, ComputerLabelEdition, ComputerOperator } from '../models/Computer';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { clusterActions } from './clusterActions';
import { clusterCreationAction } from './clusterCreationAction';
import { computerActions } from './computerActions';
import { createICAction } from './createICAction';
import { resourceActions } from './resourceActions';
import { resourceDetailActions } from './resourceDetailActions';

type GetState = () => RootState;

/**
 * 集群列表刷新
 */
const fetchClusterAfterReset: OperationHooks = {
  [OperationTrigger.Reset]: dispatch => {
    dispatch(clusterActions.fetch());
    dispatch(clusterActions.selectCluster([]));
  }
};

/**
 * 节点列表刷新Action
 */
const fetchComputerAfterReset: OperationHooks = {
  [OperationTrigger.Reset]: dispatch => {
    dispatch(computerActions.poll());
    dispatch(computerActions.selects([]));
  }
};

/**
 * 操作流的Actions
 */
export const workflowActions = {
  /** 删除集群 Flow的操作 */
  deleteCluster: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.DeleteCluster,
    workflowStateLocator: (state: RootState) => state.deleteClusterFlow,
    operationExecutor: WebAPI.deleteResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { deleteClusterFlow } = getState();

        if (isSuccessWorkflow(deleteClusterFlow)) {
          dispatch(clusterActions.poll());
          dispatch(workflowActions.deleteCluster.reset());
        }
      }
    }
  }),

  /**
   * 编辑集群名称
   */
  modifyClusterName: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ModifyClusterNameWorkflow,
    workflowStateLocator: (state: RootState) => state.modifyClusterName,
    operationExecutor: WebAPI.modifyClusterName,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { modifyClusterName } = getState();

        if (isSuccessWorkflow(modifyClusterName)) {
          dispatch(workflowActions.modifyClusterName.reset());
          dispatch(clusterActions.applyFilter({}));
        }
      }
    }
  }),

  /** 批量UnSchedule */
  batchUnScheduleComputer: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.BatchUnScheduleComputer,
    workflowStateLocator: (state: RootState) => state.subRoot.computerState.batchUnScheduleComputer,
    operationExecutor: WebAPI.updateResourceIns,
    after: extend({}, fetchComputerAfterReset, {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let {
          subRoot: {
            computerState: { batchUnScheduleComputer }
          }
        } = getState();

        if (isSuccessWorkflow(batchUnScheduleComputer)) {
          dispatch(workflowActions.batchUnScheduleComputer.reset());
        }
      }
    })
  }),

  /** 批量 TurnOnSchedule */
  batchTurnOnScheduleComputer: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.BatchTurnOnSchedulingComputer,
    workflowStateLocator: (state: RootState) => state.subRoot.computerState.batchTurnOnSchedulingComputer,
    operationExecutor: WebAPI.updateResourceIns,
    after: extend({}, fetchComputerAfterReset, {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let {
          subRoot: {
            computerState: { batchTurnOnSchedulingComputer }
          },
          route
        } = getState();

        if (isSuccessWorkflow(batchTurnOnSchedulingComputer)) {
          dispatch(workflowActions.batchTurnOnScheduleComputer.reset());
        }
      }
    })
  }),

  /** 批量 驱逐 Node节点 */
  batchDrainComputer: generateWorkflowActionCreator<Computer, ComputerOperator>({
    actionType: ActionType.DrainComputer,
    workflowStateLocator: (state: RootState) => state.subRoot.computerState.drainComputer,
    operationExecutor: WebAPI.drainComputer,
    after: extend({}, fetchComputerAfterReset, {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { drainComputer } = getState().subRoot.computerState;
        if (isSuccessWorkflow(drainComputer)) {
          dispatch(workflowActions.batchDrainComputer.reset());
          dispatch(computerActions.selects([]));
        }
      }
    })
  }),

  /** 编辑 Node节点 */
  updateNodeLabel: generateWorkflowActionCreator<ComputerLabelEdition, ComputerOperator>({
    actionType: ActionType.UpdateNodeLabel,
    workflowStateLocator: (state: RootState) => state.subRoot.computerState.updateNodeLabel,
    operationExecutor: WebAPI.updateComputerLabel,
    after: extend({}, fetchComputerAfterReset, {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { updateNodeLabel } = getState().subRoot.computerState;
        if (isSuccessWorkflow(updateNodeLabel)) {
          dispatch(workflowActions.updateNodeLabel.reset());
        }
      }
    })
  }),

  /** 驱逐 Node节点 */
  updateNodeTaint: generateWorkflowActionCreator<ComputerLabelEdition, ComputerOperator>({
    actionType: ActionType.UpdateNodeTaint,
    workflowStateLocator: (state: RootState) => state.subRoot.computerState.updateNodeTaint,
    operationExecutor: WebAPI.updateComputerTaints,
    after: extend({}, fetchComputerAfterReset, {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { updateNodeTaint } = getState().subRoot.computerState;
        if (isSuccessWorkflow(updateNodeTaint)) {
          dispatch(workflowActions.updateNodeTaint.reset());
        }
      }
    })
  }),

  /** 更新超售比 */
  updateClusterAllocationRatio: generateWorkflowActionCreator<CreateResource, any>({
    actionType: ActionType.UpdateClusterAllocationRatio,
    workflowStateLocator: (state: RootState) => state.subRoot.updateClusterAllocationRatio,
    operationExecutor: WebAPI.updateResourceIns,
    after: extend({}, fetchClusterAfterReset, {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { updateClusterAllocationRatio } = getState().subRoot;
        if (isSuccessWorkflow(updateClusterAllocationRatio)) {
          dispatch(workflowActions.updateClusterAllocationRatio.reset());
        }
      }
    })
  }),

  /** 更新token*/
  updateClusterToken: generateWorkflowActionCreator<CreateResource, any>({
    actionType: ActionType.UpdateClusterToken,
    workflowStateLocator: (state: RootState) => state.updateClusterToken,
    operationExecutor: WebAPI.updateResourceIns,
    after: extend(
      {},
      {
        [OperationTrigger.Done]: (dispatch, getState) => {
          let { updateClusterToken } = getState();
          if (isSuccessWorkflow(updateClusterToken)) {
            dispatch(clusterActions.clearClustercredential());
            dispatch(workflowActions.updateClusterToken.reset());
          }
        }
      }
    )
  }),

  /** 创建多种资源的flow */
  applyResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ApplyResource,
    workflowStateLocator: (state: RootState) => state.subRoot.applyResourceFlow,
    operationExecutor: WebAPI.applyResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, namespaceSelection, route } = getState(),
          urlParams = router.resolve(route),
          { applyResourceFlow, workloadEdit, resourceName, namespaceEdit } = subRoot;
        if (isSuccessWorkflow(applyResourceFlow)) {
          // 如果是创建页面吗，即创建workload页面，则跳转至事件列表
          if (urlParams['mode'] === 'create') {
            if (resourceName !== 'np') {
              // 如果resourceType是config，则直接跳转到列表即可，因为config没有事件
              if (urlParams['type'] === 'config') {
                router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries);
              } else {
                // 这里是去更新resourceInfo的信息，不然跳转会出问题，拉取pod列表等的都会出错
                dispatch(resourceActions.initResourceInfo(workloadEdit.workloadType));
                router.navigate(
                  Object.assign({}, urlParams, {
                    resourceName: workloadEdit.workloadType,
                    mode: 'detail',
                    tab: 'event'
                  }),
                  Object.assign({}, route.queries, { resourceIns: workloadEdit.workloadName })
                );
              }
            } else {
              router.navigate(
                Object.assign({}, urlParams, {
                  mode: 'detail',
                  tab: 'nsInfo'
                }),
                Object.assign({}, route.queries, { resourceIns: namespaceEdit.name })
              );
            }
          } else {
            // 如果是使用yaml创建，则直接跳回列表页
            router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries);
          }
          dispatch(workflowActions.applyResource.reset());
          dispatch(resourceActions.poll());
        }
      }
    }
  }),

  /**创建多种资源的Flow 每种资源调用的接口不一样 比如tapp*/
  applyDifferentInterfaceResource: generateWorkflowActionCreator<CreateResource, DifferentInterfaceResourceOperation[]>(
    {
      actionType: ActionType.ApplyDifferentInterfaceResource,
      workflowStateLocator: (state: RootState) => state.subRoot.applyDifferentInterfaceResourceFlow,
      operationExecutor: WebAPI.applyDifferentInterfaceResource,
      after: {
        [OperationTrigger.Done]: (dispatch, getState: GetState) => {
          let { subRoot, namespaceSelection, route } = getState(),
            urlParams = router.resolve(route),
            { applyDifferentInterfaceResourceFlow, workloadEdit, resourceName, namespaceEdit } = subRoot;
          if (isSuccessWorkflow(applyDifferentInterfaceResourceFlow)) {
            // 如果是创建页面吗，即创建workload页面，则跳转至事件列表
            if (urlParams['mode'] === 'create') {
              if (resourceName !== 'np') {
                // 如果resourceType是config，则直接跳转到列表即可，因为config没有事件
                if (urlParams['type'] === 'config') {
                  router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries);
                } else {
                  // 这里是去更新resourceInfo的信息，不然跳转会出问题，拉取pod列表等的都会出错
                  dispatch(resourceActions.initResourceInfo(workloadEdit.workloadType));
                  router.navigate(
                    Object.assign({}, urlParams, {
                      resourceName: workloadEdit.workloadType,
                      mode: 'detail',
                      tab: 'event'
                    }),
                    Object.assign({}, route.queries, { resourceIns: workloadEdit.workloadName })
                  );
                }
              } else {
                router.navigate(
                  Object.assign({}, urlParams, {
                    mode: 'detail',
                    tab: 'nsInfo'
                  }),
                  Object.assign({}, route.queries, { resourceIns: namespaceEdit.name })
                );
              }
            } else {
              // 如果是使用yaml创建，则直接跳回列表页
              router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries);
            }
            dispatch(workflowActions.applyResource.reset());
            dispatch(resourceActions.poll());
          }
        }
      }
    }
  ),
  /** 创建、编辑resourceIns */
  modifyResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ModifyResource,
    workflowStateLocator: (state: RootState) => state.subRoot.modifyResourceFlow,
    operationExecutor: WebAPI.modifyResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, namespaceSelection, route } = getState(),
          urlParams = router.resolve(route),
          { modifyResourceFlow, workloadEdit } = subRoot;
        if (isSuccessWorkflow(modifyResourceFlow)) {
          // 跳转到resouce的列表界面，其余的都调到 事件列表
          let { type: resourceType, resourceName } = urlParams;

          if (
            resourceType === 'namespace' ||
            resourceType === 'config' ||
            resourceType === 'storage' ||
            resourceName === 'ingress'
          ) {
            router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries);
          } else {
            let target = modifyResourceFlow.targets[0],
              resourceIns =
                target.mode === 'create'
                  ? JSON.parse(modifyResourceFlow.targets[0].jsonData).metadata.name
                  : route.queries['resourceIns'],
              urlChangeParams = { mode: 'detail', tab: 'event' };

            // 这里去判断是因为创建workload的时候，有五种类型进行选择，所以需要去更改路由和 当前的resourceInfo的信息
            if (resourceType === 'resource' && target.mode === 'create') {
              dispatch(resourceActions.initResourceInfo(workloadEdit.workloadType));
              urlChangeParams = Object.assign({}, urlChangeParams, { resourceName: workloadEdit.workloadType });
            }

            router.navigate(
              Object.assign({}, urlParams, urlChangeParams),
              Object.assign({}, route.queries, { resourceIns })
            );
          }
          dispatch(workflowActions.modifyResource.reset());
          dispatch(resourceActions.poll());
        }
      }
    }
  }),

  /** 删除resourceIns */
  deleteResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.DeleteResource,
    workflowStateLocator: (state: RootState) => state.subRoot.deleteResourceFlow,
    operationExecutor: WebAPI.deleteResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, namespaceSelection, route } = getState(),
          { deleteResourceFlow } = subRoot;

        if (isSuccessWorkflow(deleteResourceFlow)) {
          dispatch(workflowActions.deleteResource.reset());
          dispatch(resourceActions.poll());
        }
      }
    }
  }),

  /** 删除pod */
  deletePod: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.DeletePod,
    workflowStateLocator: (state: RootState) => state.subRoot.resourceDetailState.deletePodFlow,
    operationExecutor: WebAPI.deleteResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, route } = getState(),
          { deletePodFlow } = subRoot.resourceDetailState;

        if (isSuccessWorkflow(deletePodFlow)) {
          dispatch(resourceDetailActions.pod.podSelect([]));
          dispatch(workflowActions.deletePod.reset());
          setTimeout(() => {
            dispatch(
              resourceDetailActions.pod.poll({
                namespace: route.queries['np'],
                regionId: +route.queries['rid'],
                clusterId: route.queries['clusterId'],
                specificName: route.queries['resourceIns']
              })
            );
          }, 1000);
        }
      }
    }
  }),

  /** 回滚resourceIns */
  rollbackResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.RollBackResource,
    workflowStateLocator: (state: RootState) => state.subRoot.resourceDetailState.rollbackResourceFlow,
    operationExecutor: WebAPI.rollbackResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, namespaceSelection, route } = getState(),
          { resourceDetailState } = subRoot,
          { rollbackResourceFlow } = resourceDetailState;

        if (isSuccessWorkflow(rollbackResourceFlow)) {
          dispatch(workflowActions.rollbackResource.reset());
          dispatch(
            resourceDetailActions.rs.applyFilter({
              namespace: namespaceSelection,
              clusterId: route.queries['clusterId'],
              regionId: +route.queries['rid']
            })
          );
        }
      }
    }
  }),

  /**删除指定的tapp下的pod */
  removeTappPod: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.RemoveTappPod,
    workflowStateLocator: (state: RootState) => state.subRoot.resourceDetailState.removeTappPodFlow,
    operationExecutor: WebAPI.updateResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { route } = getState();
        let { removeTappPodFlow } = getState().subRoot.resourceDetailState;

        if (isSuccessWorkflow(removeTappPodFlow)) {
          dispatch(workflowActions.removeTappPod.reset());
          dispatch(
            resourceDetailActions.pod.applyFilter({
              namespace: route.queries['np'],
              regionId: +route.queries['rid'],
              clusterId: route.queries['clusterId'],
              specificName: route.queries['resourceIns']
            })
          );
        }
      }
    }
  }),
  /**指定tapp容器实例灰度升级 */
  updateGrayTapp: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.UpdateGrayTapp,
    workflowStateLocator: (state: RootState) => state.subRoot.resourceDetailState.updateGrayTappFlow,
    operationExecutor: WebAPI.updateResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { route } = getState();
        let { updateGrayTappFlow } = getState().subRoot.resourceDetailState;

        if (isSuccessWorkflow(updateGrayTappFlow)) {
          dispatch(workflowActions.updateGrayTapp.reset());
          dispatch(resourceActions.fetch());
          dispatch(
            resourceDetailActions.pod.applyFilter({
              namespace: route.queries['np'],
              regionId: +route.queries['rid'],
              clusterId: route.queries['clusterId'],
              specificName: route.queries['resourceIns']
            })
          );
        }
      }
    }
  }),
  /**
   * 1. 更新访问方式 —— Service
   * 2. 更新转发配置 —— Ingress
   * 3. 滚动更新镜像 —— Deployment、StatefulSet、Daemonset
   */
  updateResourcePart: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.UpdateResourcePart,
    workflowStateLocator: (state: RootState) => state.subRoot.updateResourcePart,
    operationExecutor: WebAPI.updateResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, route } = getState(),
          urlParams = router.resolve(route),
          { updateResourcePart } = subRoot;

        if (isSuccessWorkflow(updateResourcePart)) {
          dispatch(workflowActions.updateResourcePart.reset());
          // 进行路由的跳转
          let urlChangeParams = JSON.parse(
            JSON.stringify(Object.assign({}, urlParams, { mode: 'list', tab: undefined }))
          );
          let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { resourceIns: undefined })));
          router.navigate(urlChangeParams, routeQueries);

          // 还需要进行资源列表的拉取
          dispatch(resourceActions.poll());
        }
      }
    }
  }),

  createCluster: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.CreateCluster,
    workflowStateLocator: (state: RootState) => state.createClusterFlow,
    operationExecutor: WebAPI.createImportClsutter,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { createClusterFlow, route } = getState();

        if (isSuccessWorkflow(createClusterFlow)) {
          router.navigate({}, { rid: route.queries['rid'] });
          dispatch(clusterCreationAction.clearClusterCreationState());
          dispatch(workflowActions.createCluster.reset());
        }
      }
    }
  }),

  createIC: generateWorkflowActionCreator<CreateIC, number>({
    actionType: ActionType.CreateIC,
    workflowStateLocator: (state: RootState) => state.createICWorkflow,
    operationExecutor: WebAPI.createIC,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let { createICWorkflow, route } = getState();

        if (isSuccessWorkflow(createICWorkflow)) {
          router.navigate({}, { rid: route.queries['rid'] });
          dispatch(createICAction.clear());
          dispatch(workflowActions.createIC.reset());
        }
      }
    }
  }),
  createComputer: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.CreateComputer,
    workflowStateLocator: (state: RootState) => state.subRoot.computerState.createComputerWorkflow,
    operationExecutor: WebAPI.modifyMultiResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let {
          subRoot: {
            computerState: { createComputerWorkflow }
          },
          route
        } = getState();

        if (isSuccessWorkflow(createComputerWorkflow)) {
          dispatch(workflowActions.createComputer.reset());
          router.navigate({ sub: 'sub', mode: 'list', type: 'nodeManange', resourceName: 'node' }, route.queries);
        }
      }
    }
  }),
  deleteComputer: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.DeleteComputer,
    workflowStateLocator: (state: RootState) => state.subRoot.computerState.deleteComputer,
    operationExecutor: WebAPI.deleteResourceIns,
    after: extend({}, fetchComputerAfterReset, {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
        let {
          subRoot: {
            computerState: { deleteComputer }
          },
          route
        } = getState();

        if (isSuccessWorkflow(deleteComputer)) {
          router.navigate({ sub: 'sub', mode: 'list', type: 'nodeManange', resourceName: 'node' }, route.queries);
          dispatch(workflowActions.deleteComputer.reset());
        }
      }
    })
  }),

  /**
   * 1. 更新访问方式 —— Service
   * 2. 更新转发配置 —— Ingress
   * 3. 滚动更新镜像 —— Deployment、StatefulSet、Daemonset
   */
  updateMultiResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.UpdateMultiResource,
    workflowStateLocator: (state: RootState) => state.subRoot.updateMultiResource,
    operationExecutor: WebAPI.updateMultiResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, route } = getState(),
          urlParams = router.resolve(route),
          { updateMultiResource } = subRoot;

        if (isSuccessWorkflow(updateMultiResource)) {
          dispatch(workflowActions.updateMultiResource.reset());
          // 进行路由的跳转
          let urlChangeParams = JSON.parse(
            JSON.stringify(Object.assign({}, urlParams, { mode: 'list', tab: undefined }))
          );
          let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { resourceIns: undefined })));
          router.navigate(urlChangeParams, routeQueries);

          // 还需要进行资源列表的拉取
          dispatch(resourceActions.poll());
        }
      }
    }
  }),

  /** 同时创建、编辑多个resourceIns 例如backendGroup*/
  modifyMultiResource: generateWorkflowActionCreator<CreateResource, number>({
    actionType: ActionType.ModifyMultiResource,
    workflowStateLocator: (state: RootState) => state.subRoot.modifyMultiResourceWorkflow,
    operationExecutor: WebAPI.modifyMultiResourceIns,
    after: {
      [OperationTrigger.Done]: (dispatch, getState: GetState) => {
        let { subRoot, namespaceSelection, route } = getState(),
          urlParams = router.resolve(route),
          { modifyMultiResourceWorkflow } = subRoot;
        if (isSuccessWorkflow(modifyMultiResourceWorkflow)) {
          // 跳转到 事件列表
          let resourceName = urlParams['resourceName'],
            resourceAction = urlParams['tab'];
          if (resourceName === 'lbcf' && resourceAction === 'createBG') {
            let resourceIns;

            resourceIns = route.queries['resourceIns'];
            let urlChangeParams = { mode: 'detail', tab: 'event' };

            router.navigate(
              Object.assign({}, urlParams, urlChangeParams),
              Object.assign({}, route.queries, { resourceIns })
            );
          }
          dispatch(workflowActions.modifyMultiResource.reset());
          dispatch(resourceActions.poll());
        }
      }
    }
  })
};
