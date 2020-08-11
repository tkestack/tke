import {
  ReduxAction,
  extend,
  generateWorkflowActionCreator,
  OperationTrigger,
  isSuccessWorkflow,
  getWorkflowStatistics,
  createFFObjectActions,
  uuid
} from '@tencent/ff-redux';
import {
  RootState,
  ChartFilter,
  ChartEditor,
  Chart,
  ChartInfo,
  ChartInfoFilter,
  ChartDetailFilter,
  ChartVersionFilter,
  ChartVersion,
  RemovedChartVersion,
  ChartTreeFile
} from '../../models';
import * as ActionTypes from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { initChartEditorState } from '../../constants/initState';
import { ChartValidateSchema } from '../../constants/ChartValidateConfig';
import { router } from '../../router';
import { createValidatorActions, getValidatorActionType } from '@tencent/ff-validator';
import { cluster } from '@config/resource/k8sConfig';

type GetState = () => RootState;

/**
 * 修改Chart
 */
const updateChartWorkflow = generateWorkflowActionCreator<Chart, ChartDetailFilter>({
  actionType: ActionTypes.UpdateChart,
  workflowStateLocator: (state: RootState) => state.chartUpdateWorkflow,
  operationExecutor: WebAPI.updateChart,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      if (isSuccessWorkflow(getState().chartUpdateWorkflow)) {
        //表示编辑模式结束
        let { chartEditor } = getState();
        dispatch({
          type: ActionTypes.UpdateChartEditorState,
          payload: Object.assign({}, chartEditor, { v_editing: false })
        });
        const params = getState().chartUpdateWorkflow.params;
        /** 重新获取最新数据，从而Detail可以连续编辑且使用到最新的resourceVersion */
        dispatch(detailActions.poll(params));
      }
      /** 结束工作流 */
      dispatch(detailActions.updateChartWorkflow.reset());
    }
  }
});

/**
 * 删除模板
 */
const removeChartVersionWorkflow = generateWorkflowActionCreator<ChartVersion, ChartVersionFilter>({
  actionType: ActionTypes.RemoveChartVersion,
  workflowStateLocator: (state: RootState) => state.chartVersionRemoveWorkflow,
  operationExecutor: WebAPI.deleteChartVersion,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      //在列表页删除的动作，因此直接重新拉取一次数据
      const params = getState().chartVersionRemoveWorkflow.params;
      dispatch(detailActions.poll(params.chartDetailFilter));
      /** 结束工作流 */
      dispatch(detailActions.removeChartVersionWorkflow.reset());
    }
  }
});

/**
 * TODO: 获取模板
 */
const fetchChartVersionFileActions = createFFObjectActions<any, ChartVersionFilter>({
  actionName: ActionTypes.ChartVersionFile,
  fetcher: async query => {
    let response = await WebAPI.fetchChartVersionFile(query.filter);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartVersionFile;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {}
});

/**
 * 获取chart详情
 */
const fetchChartInfoActions = createFFObjectActions<ChartInfo, ChartInfoFilter>({
  actionName: ActionTypes.ChartInfo,
  fetcher: async (query, getState: GetState) => {
    let response: ChartInfo = await WebAPI.fetchChartInfo(query.filter);
    const findNode = (tree: ChartTreeFile, path: string): ChartTreeFile => {
      for (let i = 0; i < tree.children.length; i++) {
        let c = tree.children[i];
        if (c.fullPath === path) {
          return c;
        }
        if (path.indexOf(c.fullPath) === 0) {
          let f = findNode(c, path);
          if (f) {
            return f;
          }
        }
      }
      return undefined;
    };

    const { route, chartInfo } = getState();
    let tree: ChartTreeFile = { name: route.queries['chartName'], fullPath: '', data: '', children: [] };
    if (response && response.spec.rawFiles) {
      Object.keys(response.spec.rawFiles).forEach(path => {
        if (!path.includes('/')) {
          tree.children.push({
            fullPath: path,
            name: path,
            data: response.spec.rawFiles[path],
            children: []
          });
        } else {
          let p = '',
            c = path;
          let lastNode: ChartTreeFile = tree;
          while (true) {
            p = p === '' ? c.substring(0, c.indexOf('/')) : p + '/' + c.substring(0, c.indexOf('/'));
            c = c.substring(c.indexOf('/') + 1);
            let f = findNode(tree, p);
            //找到父节点
            if (f) {
              lastNode = f;
              //不是叶子节点，继续
              if (c.includes('/')) {
                continue;
              }
              //叶子节点
              f.children.push({
                fullPath: path,
                name: c,
                data: response.spec.rawFiles[path],
                children: []
              });
              break;
            } else {
              //未找到父节点
              let newNode = {
                fullPath: p,
                name: p.includes('/') ? p.substring(p.lastIndexOf('/') + 1) : p,
                data: '',
                children: []
              };
              lastNode.children.push(newNode);
              lastNode = newNode;

              //不是叶子节点，继续
              if (c.includes('/')) {
                continue;
              }
              //叶子节点
              lastNode.children.push({
                fullPath: path,
                name: c,
                data: response.spec.rawFiles[path],
                children: []
              });
              break;
            }
          }
        }
      });
    }
    response.fileTree = tree;
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartInfo;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {}
});

/**
 * 获取chart详情
 */
const fetchChartActions = createFFObjectActions<Chart, ChartDetailFilter>({
  actionName: ActionTypes.Chart,
  fetcher: async query => {
    let response = await WebAPI.fetchChart(query.filter);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartDetail;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let sorted = [];
    let selectedVersion: ChartVersion = { id: uuid() };
    if (record.data) {
      sorted = record.data.status.versions.sort((a, b) => {
        let oDate1 = new Date(a.timeCreated);
        let oDate2 = new Date(b.timeCreated);
        return oDate1.getTime() > oDate2.getTime() ? -1 : 1;
      });
      selectedVersion = sorted[0];
    }

    let editor: ChartEditor = record.data;
    editor.sortedVersions = sorted;
    editor.selectedVersion = selectedVersion;
    dispatch({
      type: ActionTypes.UpdateChartEditorState,
      payload: editor
    });

    let { route } = getState(),
      urlParam = router.resolve(route);
    if (urlParam['tab'] && urlParam['tab'] === 'file') {
      //请求文件目录树
      dispatch(
        detailActions.chartInfo.applyFilter({
          cluster: '',
          namespace: '',
          metadata: {
            namespace: editor.metadata.namespace,
            name: editor.metadata.name
          },
          chartVersion: selectedVersion.version,
          projectID: route.queries['prj']
        })
      );
    }

    //如果removedChartVersions列表的元素已经不在最新列表中，则删除该元素
    if (record.data) {
      let { removedChartVersions } = getState();
      let vs = [];
      removedChartVersions.versions.forEach(v => {
        if (v.namespace === record.data.metadata.namespace && v.name === record.data.metadata.name) {
          if (
            record.data.status.versions.find(x => {
              return x.version === v.version;
            }) !== undefined
          ) {
            vs.push(v);
          }
        } else {
          vs.push(v);
        }
      });
      dispatch({
        type: ActionTypes.RemovedChartVersions,
        payload: { versions: vs }
      });
    }
  }
});

const restActions = {
  updateChartWorkflow,
  removeChartVersionWorkflow,

  /** 轮询操作 */
  poll: (filter: ChartDetailFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        detailActions.polling({
          delayTime: 5000,
          filter: filter
        })
      );
    };
  },

  validator: createValidatorActions({
    userDefinedSchema: ChartValidateSchema,
    validateStateLocator: (store: RootState) => {
      return store.chartEditor;
    },
    validatorStateLocation: (store: RootState) => {
      return store.chartValidator;
    }
  }),

  /** 更新状态 */
  updateEditorState: obj => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { chartEditor } = getState();
      dispatch({
        type: ActionTypes.UpdateChartEditorState,
        payload: Object.assign({}, chartEditor, obj)
      });
    };
  },

  /** 离开更新页面，清除Editor当中的内容 */
  clearEditorState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdateChartEditorState,
      payload: initChartEditorState
    };
  },

  /** 离开创建页面，清除Validator当中的内容 */
  clearValidatorState: (): ReduxAction<any> => {
    return {
      type: getValidatorActionType(ChartValidateSchema.formKey),
      payload: {}
    };
  },

  /** 记录已删除版本到state */
  addRemovedChartVersion: (version: RemovedChartVersion) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { removedChartVersions } = getState();
      const found = removedChartVersions.versions.find(x => {
        return x.name === version.name && x.namespace === version.namespace && x.version === version.version;
      });
      if (!found) {
        removedChartVersions.versions.push(version);
        dispatch({
          type: ActionTypes.RemovedChartVersions,
          payload: removedChartVersions
        });
      }
    };
  }
};
export const detailActions = extend({}, fetchChartActions, { chartInfo: fetchChartInfoActions }, restActions);
