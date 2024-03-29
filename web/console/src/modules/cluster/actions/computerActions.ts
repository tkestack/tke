/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { MachineStatus } from './../constants/Config';
import { createFFListActions, extend, RecordSet, uuid } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { initValidator, ResourceInfo } from '../../common/models';
import { cloneDeep } from '../../common/utils/cloneDeep';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { Computer, DialogNameEnum, Resource, RootState } from '../models';
import { ComputerFilter, ComputerLabel } from '../models/Computer';
import * as WebAPI from '../WebAPI';
import { dialogActions } from './dialogActions';

type GetState = () => RootState;

/** 节点列表的actions */
const FFModelComputerActions = createFFListActions<Computer, ComputerFilter>({
  actionName: FFReduxActionName.COMPUTER,
  fetcher: async (query, getState: GetState) => {
    let { clusterVersion } = getState();
    let nodeInfo: ResourceInfo = resourceConfig(clusterVersion)['node'],
      nodeItems = await WebAPI.fetchResourceList(query, {
        resourceInfo: nodeInfo,
        isContinue: true,
      });
    //将machine资源和node资源关联起来
    nodeItems.records = nodeItems.records.map((item) => {
      if (
        Object.keys(item.metadata.labels).findIndex((key) => key.indexOf('node-role.kubernetes.io/master') !== -1) !==
        -1
      ) {
        item.metadata.role = 'Master&Etcd';
      } else {
        item.metadata.role = 'Worker';
      }
      let phase;
      if (item.status.conditions) {
        let nodeStatus = item.status.conditions.find((item) => item.type === 'Ready');
        phase = nodeStatus.status === 'True' ? 'Running' : 'Failed';
      }
      item.status.phase = phase;
      return item;
    });
    return nodeItems;
  },
  getRecord: (getState: GetState) => {
    return getState().subRoot.computerState.computer;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { computer } = getState().subRoot.computerState;
    if (record.data.recordCount) {
      dispatch(FFModelComputerActions.select(record.data.records[0]));
    }
    if (record.data.records.filter((item) => item.status.phase !== 'Running').length === 0) {
      dispatch(FFModelComputerActions.clearPolling());
    }
  },
});

const FFModelMachineActions = createFFListActions<Computer, ComputerFilter>({
  actionName: FFReduxActionName.MACHINE,
  fetcher: async (query, getState: GetState) => {
    let { clusterVersion } = getState();
    let k8sQueryObj = {
      fieldSelector: {
        'spec.clusterName': query.filter.clusterId ? query.filter.clusterId : undefined,
        'status.phase': MachineStatus.Initializing,
      },
    };
    let machinesInfo: ResourceInfo = resourceConfig(clusterVersion).machines,
      machineItems = await WebAPI.fetchResourceList(query, {
        resourceInfo: machinesInfo,
        k8sQueryObj,
        isContinue: true,
      });
    machineItems.records = machineItems.records.map((item) => ({
      id: uuid(),
      metadata: {
        name: item.spec.ip,
        creationTimestamp: item.metadata.creationTimestamp,
        role: 'Worker',
      },
      spec: {
        machineName: item.metadata.name,
        podCIDR: '-',
      },
      status: {
        capacity: {
          cpu: 0,
          memory: 0,
        },
        conditions: item.status.conditions,
        addresses: [],
        phase: item.status.phase,
      },
    }));
    return machineItems;
  },
  getRecord: (getState: GetState) => {
    return getState().subRoot.computerState.machine;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { machine } = getState().subRoot.computerState;
    if (getState().dialogState[DialogNameEnum.computerStatusDialog]) {
      let finder = record.data.records.find(
        (item) => item.metadata.name === (machine.selection && machine.selection.metadata.name)
      );
      if (finder) {
        dispatch(FFModelMachineActions.select(finder));
      } else {
        dispatch(dialogActions.updateDialogState(DialogNameEnum.computerStatusDialog));
      }
    } else if (record.data.recordCount) {
      dispatch(FFModelMachineActions.select(record.data.records[0]));
    }
    if (record.data.records.length === 0) {
      dispatch(FFModelMachineActions.clearPolling());
    }
  },
});

const restActions = {
  machine: FFModelMachineActions,

  poll: (filter?: ComputerFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        subRoot: {
          computerState: { computer },
        },
      } = getState();
      dispatch(
        FFModelComputerActions.polling({
          filter: filter || computer.query.filter,
          delayTime: 8000,
        })
      );
    };
  },

  pollMachine: (filter?: ComputerFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        subRoot: {
          computerState: { machine },
        },
      } = getState();
      dispatch(
        FFModelMachineActions.polling({
          filter: filter || machine.query.filter,
          delayTime: 8000,
        })
      );
    };
  },
  /**初始化label */
  initLabelEdition: (labels: { [props: string]: string }, computerName: string) => {
    return async (dispatch, getState: GetState) => {
      let { labelEdition } = getState().subRoot.computerState;
      let labelArray = [];
      let labelKeys = Object.keys(labels);
      labelKeys.forEach((key) => {
        let disabled = key.includes('kubernetes');
        let item = {
          value: labels[key],
          key: key,
          disabled,
        };
        if (disabled) {
          labelArray.unshift(item);
        } else {
          labelArray.push(item);
        }
      });
      dispatch({
        type: ActionType.UpdateLabelEdition,
        payload: Object.assign({}, labelEdition, {
          computerName,
          labels: labelArray.map((label) => {
            return Object.assign({}, label, { id: uuid(), v_key: initValidator, v_value: initValidator });
          }),
          originLabel: labels,
        }),
      });
    };
  },

  /** 新增label变量 */
  addLabel: () => {
    return async (dispatch, getState: GetState) => {
      let { labelEdition } = getState().subRoot.computerState;
      let labels = cloneDeep(labelEdition.labels);
      let newlabel = {
        id: uuid(),
        key: '',
        value: '',
      };

      labels.push(newlabel);
      dispatch({
        type: ActionType.UpdateLabelEdition,
        payload: Object.assign({}, labelEdition, { labels }),
      });
    };
  },

  /** 删除label变量 */
  deleteLabel: (Id: string) => {
    return async (dispatch, getState: GetState) => {
      let { labelEdition } = getState().subRoot.computerState;
      let labels = cloneDeep(labelEdition.labels),
        eIndex = labels.findIndex((e) => e.id === Id);

      labels.splice(eIndex, 1);
      dispatch({
        type: ActionType.UpdateLabelEdition,
        payload: Object.assign({}, labelEdition, { labels }),
      });
    };
  },

  /** 更新label变量 */
  updateLabel: (obj: any, Id: string) => {
    return async (dispatch, getState: GetState) => {
      let { labelEdition } = getState().subRoot.computerState;
      let labels = cloneDeep(labelEdition.labels),
        eIndex = labels.findIndex((e) => e.id === Id),
        objKeys = Object.keys(obj);

      objKeys.forEach((item) => {
        labels[eIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.UpdateLabelEdition,
        payload: Object.assign({}, labelEdition, { labels }),
      });
    };
  },
  /**初始化label */
  initTaintEdition: (taints: { [props: string]: string }[], computerName: string) => {
    return async (dispatch, getState: GetState) => {
      let { taintEdition } = getState().subRoot.computerState;
      if (!taints) {
        dispatch({
          type: ActionType.UpdateTaintEdition,
          payload: Object.assign({}, taintEdition, {
            computerName,
            taints: [
              { id: uuid(), v_key: initValidator, v_value: initValidator, key: '', value: '', effect: 'NoSchedule' },
            ],
          }),
        });
      } else {
        let taintsArray = [];
        taints.forEach((taint) => {
          let disabled = taint.key.indexOf('kubernetes') !== -1 ? true : false;
          let item = Object.assign({}, { value: '' }, taint, {
            id: uuid(),
            v_key: initValidator,
            v_value: initValidator,
            disabled,
          });
          if (disabled) {
            taintsArray.unshift(item);
          } else {
            taintsArray.push(item);
          }
        });
        dispatch({
          type: ActionType.UpdateTaintEdition,
          payload: Object.assign({}, taintEdition, {
            computerName,
            taints: taintsArray,
          }),
        });
      }
    };
  },

  /** 新增节点污点 */
  addTaint: () => {
    return async (dispatch, getState: GetState) => {
      let { taintEdition } = getState().subRoot.computerState;
      let taints = cloneDeep(taintEdition.taints);
      let newtaints = {
        id: uuid(),
        key: '',
        v_key: initValidator,
        value: '',
        v_value: initValidator,
        effect: 'NoSchedule',
      };

      taints.push(newtaints);
      dispatch({
        type: ActionType.UpdateTaintEdition,
        payload: Object.assign({}, taintEdition, { taints }),
      });
    };
  },

  /** 删除Taint变量 */
  deleteTaint: (Id: string) => {
    return async (dispatch, getState: GetState) => {
      let { taintEdition } = getState().subRoot.computerState;
      let taints = cloneDeep(taintEdition.taints),
        eIndex = taints.findIndex((e) => e.id === Id);

      taints.splice(eIndex, 1);
      dispatch({
        type: ActionType.UpdateTaintEdition,
        payload: Object.assign({}, taintEdition, { taints }),
      });
    };
  },

  /** 更新label变量 */
  updateTaint: (obj: any, Id: string) => {
    return async (dispatch, getState: GetState) => {
      let { taintEdition } = getState().subRoot.computerState;
      let taints = cloneDeep(taintEdition.taints),
        eIndex = taints.findIndex((e) => e.id === Id),
        objKeys = Object.keys(obj);

      objKeys.forEach((item) => {
        taints[eIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.UpdateTaintEdition,
        payload: Object.assign({}, taintEdition, { taints }),
      });
    };
  },

  showMachine: (isShow: boolean) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.IsShowMachine,
        payload: isShow,
      });
    };
  },

  fetchdeleteMachineResouceIns: (nodeName: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          computerState: {
            machine: { query },
          },
        },
        clusterVersion,
      } = getState();
      let k8sQueryObj = {
        fieldSelector: {
          'spec.clusterName': query.filter.clusterId ? query.filter.clusterId : undefined,
          'spec.ip': nodeName,
        },
      };
      let machinesInfo: ResourceInfo = resourceConfig(clusterVersion).machines,
        machineItems = await WebAPI.fetchResourceList(query, {
          resourceInfo: machinesInfo,
          k8sQueryObj,
          isContinue: true,
        });
      dispatch({
        type: ActionType.FetchDeleteMachineResouceIns,
        payload: machineItems.recordCount ? machineItems.records[0] : null,
      });
    };
  },
};

export const computerActions = extend(FFModelComputerActions, restActions);
