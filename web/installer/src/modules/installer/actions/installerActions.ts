import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { extend, uuid } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { RootState, EditState } from '../models';
import { initMachine, initArg } from '../reducers/initState';
import { cloneDeep, isEmpty } from '../../common/utils';
import { generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { initValidation } from '../../../../helpers/Validator';
import { Base64 } from 'js-base64';
import { validateActions } from './validateActions';

type GetState = () => RootState;

const fetchClusterActions = generateFetcherActionCreator({
  actionType: ActionType.FetchCluster,
  fetcher: (getState: GetState) => {
    return WebAPI.fetchCluster();
  },
  finish: (dispatch, getState) => {
    const { cluster } = getState();
    switch (cluster.data.record['progress']['status']) {
      case 'Unknown':
        dispatch(restActions.stepNext('step1'));
        break;
      case 'Doing':
      case 'Success':
      case 'Failed':
        dispatch(restActions.stepNext('step10'));
        break;
      default:
        dispatch(restActions.stepNext('step1'));
    }
  }
});

const fetchProgressActions = generateFetcherActionCreator({
  actionType: ActionType.FetchProgress,
  fetcher: async (getState: GetState) => {
    let rsp = await WebAPI.fetchProgress();
    if (rsp.record['status'] !== 'Doing') {
      clearInterval(window['pollProgress']);
    }
    return rsp;
  }
});

const workflowActions = {
  createCluster: generateWorkflowActionCreator<EditState, void>({
    actionType: ActionType.CreateCluster,
    workflowStateLocator: (state: RootState) => state.createCluster,
    operationExecutor: WebAPI.createCluster,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { createCluster } = getState();
        if (isSuccessWorkflow(createCluster)) {
          dispatch(installerActions.createCluster.reset());
          dispatch(installerActions.stepNext('step10'));
          dispatch(installerActions.poll());
        }
      }
    }
  })
};

const restActions = {
  stepNext: (step: string) => {
    return (dispatch, getState) => {
      dispatch({
        type: ActionType.StepNext,
        payload: step
      });
    };
  },
  addMachine: () => {
    return (dispatch, getState) => {
      const { editState } = getState();
      const machine = editState.machines.find(m => m.status === 'editing');
      const canAdd = machine ? validateActions._validateMachine(machine) : true;
      if (machine) {
        dispatch(validateActions.validateMachine(machine));
      }

      if (canAdd) {
        const machines = cloneDeep(editState.machines).map(m => Object.assign({}, m, { status: 'edited' }));

        machines.push(Object.assign({}, initMachine, { id: uuid() }));

        dispatch({
          type: ActionType.UpdateEdit,
          payload: { machines }
        });
      }
    };
  },

  addDockerExtraArgs: () => {
    return (dispatch, getState) => {
      const { editState } = getState(),
        args = cloneDeep(editState.dockerExtraArgs);

      args.push(Object.assign({}, initArg, { id: uuid() }));
      dispatch({
        type: ActionType.UpdateEdit,
        payload: { dockerExtraArgs: args }
      });
    };
  },

  removeDockerExtraArgs: (id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.dockerExtraArgs.forEach(m => {
        if (m.id !== id) {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { dockerExtraArgs: args }
      });
    };
  },

  updateDockerExtraArgs: (obj: any, id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.dockerExtraArgs.forEach(m => {
        if (m.id === id) {
          args.push(Object.assign(m, obj));
        } else {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { dockerExtraArgs: args }
      });
    };
  },

  addKubeletExtraArgs: () => {
    return (dispatch, getState) => {
      const { editState } = getState(),
        args = cloneDeep(editState.kubeletExtraArgs);

      args.push(Object.assign({}, initArg, { id: uuid() }));
      dispatch({
        type: ActionType.UpdateEdit,
        payload: { kubeletExtraArgs: args }
      });
    };
  },

  removeKubeletExtraArgs: (id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.kubeletExtraArgs.forEach(m => {
        if (m.id !== id) {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { kubeletExtraArgs: args }
      });
    };
  },

  updateKubeletExtraArgs: (obj: any, id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.kubeletExtraArgs.forEach(m => {
        if (m.id === id) {
          args.push(Object.assign(m, obj));
        } else {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { kubeletExtraArgs: args }
      });
    };
  },

  addApiServerExtraArgs: () => {
    return (dispatch, getState) => {
      const { editState } = getState(),
        args = cloneDeep(editState.apiServerExtraArgs);

      args.push(Object.assign({}, initArg, { id: uuid() }));
      dispatch({
        type: ActionType.UpdateEdit,
        payload: { apiServerExtraArgs: args }
      });
    };
  },

  removeApiServerExtraArgs: (id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.apiServerExtraArgs.forEach(m => {
        if (m.id !== id) {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { apiServerExtraArgs: args }
      });
    };
  },

  updateApiServerExtraArgs: (obj: any, id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.apiServerExtraArgs.forEach(m => {
        if (m.id === id) {
          args.push(Object.assign(m, obj));
        } else {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { apiServerExtraArgs: args }
      });
    };
  },

  addControllerManagerExtraArgs: () => {
    return (dispatch, getState) => {
      const { editState } = getState(),
        args = cloneDeep(editState.controllerManagerExtraArgs);

      args.push(Object.assign({}, initArg, { id: uuid() }));
      dispatch({
        type: ActionType.UpdateEdit,
        payload: { controllerManagerExtraArgs: args }
      });
    };
  },

  removeControllerManagerExtraArgs: (id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.controllerManagerExtraArgs.forEach(m => {
        if (m.id !== id) {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { controllerManagerExtraArgs: args }
      });
    };
  },

  updateControllerManagerExtraArgs: (obj: any, id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.controllerManagerExtraArgs.forEach(m => {
        if (m.id === id) {
          args.push(Object.assign(m, obj));
        } else {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { controllerManagerExtraArgs: args }
      });
    };
  },

  addSchedulerExtraArgs: () => {
    return (dispatch, getState) => {
      const { editState } = getState(),
        args = cloneDeep(editState.schedulerExtraArgs);

      args.push(Object.assign({}, initArg, { id: uuid() }));
      dispatch({
        type: ActionType.UpdateEdit,
        payload: { schedulerExtraArgs: args }
      });
    };
  },

  removeSchedulerExtraArgs: (id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.schedulerExtraArgs.forEach(m => {
        if (m.id !== id) {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { schedulerExtraArgs: args }
      });
    };
  },

  updateSchedulerExtraArgs: (obj: any, id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        args = [];

      editState.schedulerExtraArgs.forEach(m => {
        if (m.id === id) {
          args.push(Object.assign(m, obj));
        } else {
          args.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { schedulerExtraArgs: args }
      });
    };
  },

  removeMachine: (id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        machines = [];

      editState.machines.forEach(m => {
        if (m.id !== id) {
          machines.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { machines }
      });
    };
  },

  updateMachine: (obj: any, id: string | number) => {
    return (dispatch, getState) => {
      let { editState } = getState(),
        machines = [];

      editState.machines.forEach(m => {
        if (m.id === id) {
          machines.push(Object.assign(m, obj));
        } else {
          machines.push(m);
        }
      });

      dispatch({
        type: ActionType.UpdateEdit,
        payload: { machines }
      });
    };
  },

  updateEdit: (obj: any) => {
    return dispatch => {
      dispatch({
        type: ActionType.UpdateEdit,
        payload: obj
      });
    };
  },

  updateEditFromCluster: () => {
    return (dispatch, getState) => {
      const { config } = getState().cluster.data.record;
      let cluster = {
        machines: config.cluster.spec.machines.map(m => {
          return {
            status: 'edited',
            host: m.ip,
            v_host: initValidation,
            port: m.port,
            v_port: initValidation,
            authWay: m.privateKey ? 'cert' : 'password',
            user: m.username,
            v_user: initValidation,
            password: m.privateKey ? Base64.decode(m.privatePassword) : Base64.decode(m.password),
            v_password: initValidation,
            cert: m.privateKey || '',
            v_cert: initValidation
          };
        }),
        cidr: config.cluster.spec.clusterCIDR,
        podNumLimit: config.cluster.spec.properties.maxNodePodNum,
        serviceNumLimit: config.cluster.spec.properties.maxClusterServiceNum,
        repoType: config.Config.registry ? 'remote' : 'local',
        repoAddress: config.Config.registry ? config.Config.registry.server : '',
        v_repoAddress: initValidation,
        repoUser: config.Config.registry ? config.Config.registry.username : '',
        v_repoUser: initValidation,
        repoPassword: config.Config.registry ? Base64.decode(config.Config.registry.password) : '',
        v_repoPassword: initValidation,
        domain: config.Config.dnsDomain,
        v_domain: initValidation,
        license: ''
      };
      dispatch({
        type: ActionType.UpdateEdit,
        payload: cluster
      });
    };
  },

  poll: () => {
    return dispatch => {
      dispatch(installerActions.progress.fetch());

      clearInterval(window['pollProgress']);
      window['pollProgress'] = setInterval(() => {
        dispatch(restActions.poll());
      }, 3000);
    };
  }
};

export const installerActions = extend(
  {},
  {
    cluster: fetchClusterActions,
    progress: fetchProgressActions
  },
  workflowActions,
  restActions
);
