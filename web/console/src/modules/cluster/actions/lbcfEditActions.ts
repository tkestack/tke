/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { FFReduxActionName } from './../constants/Config';
import { initStringArray } from './../constants/initState';
import { KeyValue } from 'src/modules/common';

import { deepClone, uuid, createFFListActions, RecordSet } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config/resourceConfig';
import * as ActionType from '../constants/ActionType';
import { initLbcfBackGroupEdition, initLbcfBGPort, initSelector } from '../constants/initState';
import { RootState, BackendGroup } from '../models';
import { ResourceFilter, Resource } from '../models/ResourceOption';
import { CLB, Selector } from '../models/ServiceEdit';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { validateLbcfActions } from './validateLbcfActions';
import { BackendType } from '../constants/Config';
import { GameBackgroupEdition } from '../models/LbcfEdit';

type GetState = () => RootState;

const lbcfDriverFFReduxAction = createFFListActions<Resource, ResourceFilter>({
  actionName: FFReduxActionName.LBCF_DRIVER,
  fetcher: async (query, getState: GetState, fetchOptions) => {
    let resourceInfo = resourceConfig(getState().clusterVersion).lbcf_driver;
    let kubesystemResponse = await WebAPI.fetchResourceList(
      Object.assign({}, query, { filter: Object.assign({}, query.filter, { namespace: 'kube-system' }) }),
      { resourceInfo }
    );
    let response = await WebAPI.fetchResourceList(query, { resourceInfo });
    let resourceList = kubesystemResponse.records.concat(response.records);
    const result: RecordSet<Resource> = {
      recordCount: resourceList.length,
      records: resourceList
    };
    return result;
  },
  getRecord: (getState: GetState) => {
    return getState().subRoot.lbcfEdit.driver;
  },
  selectFirst: true
});

export const lbcfEditActions = {
  driver: lbcfDriverFFReduxAction,
  /** 输入名称 */
  inputLbcfName: (name: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Gate_Name,
        payload: name
      });
    };
  },
  selectLbcfNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);

      dispatch({
        type: ActionType.Gate_Namespace,
        payload: namespace
      });
      dispatch(validateLbcfActions.validateLbcfNamespace());
    };
  },
  selectConfig: (config: KeyValue[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Lbcf_Config,
        payload: config
      });
    };
  },
  selectArgs: (args: KeyValue[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.Lbcf_Args,
        payload: args
      });
    };
  },

  updateLbcfBGPort: (backGroupId: string, id: string, object: any) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();
      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { ports } = backGroupEdition;
      let index = ports.findIndex(item => item.id === id);
      let key = Object.keys(object);
      if (key.length) {
        ports[index][key[0]] = object[key[0]];
      }
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  addLbcfBGPort: (backGroupId: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { ports } = backGroupEdition;

      ports.push(Object.assign({}, initLbcfBGPort, { id: uuid() }));
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  deleteLbcfBGPort: (backGroupId: string, id: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { ports } = backGroupEdition;

      let index = ports.findIndex(item => item.id === id);
      ports.splice(index, 1);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  updateLbcfBGLabels: (backGroupId: string, id: string, object: any) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { labels } = backGroupEdition;
      let index = labels.findIndex(item => item.id === id);
      labels[index] = Object.assign({}, labels[index], object);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  addLbcfBGLabels: (backGroupId: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { labels } = backGroupEdition;

      labels.push(Object.assign({}, initSelector, { id: uuid() }));
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  deleteLbcfBGLabels: (backGroupId: string, id: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { labels } = backGroupEdition;
      let index = labels.findIndex(item => item.id === id);
      labels.splice(index, 1);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  updateLbcfBGAddress: (backGroupId: string, id: string, object: any) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { staticAddress } = backGroupEdition;
      let index = staticAddress.findIndex(item => item.id === id);
      staticAddress[index] = Object.assign({}, staticAddress[index], object);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  addLbcfBGAddress: (backGroupId: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { staticAddress } = backGroupEdition;

      staticAddress.push(Object.assign({}, initStringArray, { id: uuid() }));
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  deleteLbcfBGAddress: (backGroupId: string, id: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { staticAddress } = backGroupEdition;
      let index = staticAddress.findIndex(item => item.id === id);
      staticAddress.splice(index, 1);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  updateLbcfBGPodName: (backGroupId: string, id: string, object: any) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { byName } = backGroupEdition;
      let index = byName.findIndex(item => item.id === id);
      byName[index] = Object.assign({}, byName[index], object);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  addLbcfBGPodName: (backGroupId: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { byName } = backGroupEdition;

      byName.push(Object.assign({}, initStringArray, { id: uuid() }));
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  deleteLbcfBGPodName: (backGroupId: string, id: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      let { byName } = backGroupEdition;
      let index = byName.findIndex(item => item.id === id);
      byName.splice(index, 1);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  initLbcfBGLabels: (backGroupId: string, labels: Selector[]) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);

      let labelArray = [];
      labels.forEach(label => {
        labelArray.push({
          id: uuid(),
          key: label.key,
          value: label.value,
          v_key: { status: 1, message: '' },
          v_value: { status: 1, message: '' }
        });
      });
      backGroupEdition.labels = labelArray;
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  addLbcfBackGroup: () => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);

      newBackGroupEdition.push(Object.assign({}, initLbcfBackGroupEdition, { id: uuid() }));
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  deleteLbcfBackGroup: (backGroupId: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();
      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let index = newBackGroupEdition.findIndex(item => item.id === backGroupId);

      newBackGroupEdition.splice(index, 1);
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  inputLbcfBackGroupName: (backGroupId: string, name: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      backGroupEdition.name = name;
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  changeBackgroupEditStatus(backGroupId: string, onEdit: boolean) {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();
      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      backGroupEdition.onEdit = onEdit;
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  initGameBGEdition: (backGroups: BackendGroup[]) => {
    return async (dispatch, getState: GetState) => {
      let backGroupEditions: GameBackgroupEdition[] = [];
      backGroups.forEach(backGroup => {
        if (backGroup.pods) {
          let keys = backGroup.pods.labels ? Object.keys(backGroup.pods.labels) : [];
          backGroupEditions.push(
            Object.assign({}, initLbcfBackGroupEdition, {
              id: uuid(),
              name: backGroup.name,
              backgroupType: BackendType.Pods,
              ports: [
                Object.assign({}, initLbcfBGPort, {
                  id: uuid(),
                  portNumber: backGroup.pods.port.portNumber,
                  protocol: backGroup.pods.port.protocol
                })
              ],
              labels: keys.map(key => {
                return Object.assign({}, initSelector, {
                  id: uuid(),
                  key: key,
                  value: backGroup.pods.labels[key]
                });
              }),
              byName: backGroup.pods.byName
                ? backGroup.pods.byName.map(name => {
                    return Object.assign({}, initStringArray, {
                      id: uuid(),
                      value: name
                    });
                  })
                : []
            })
          );
        } else if (backGroup.service) {
          let keys = backGroup.service.nodeSelector ? Object.keys(backGroup.service.nodeSelector) : [];
          backGroupEditions.push(
            Object.assign({}, initLbcfBackGroupEdition, {
              id: uuid(),
              name: backGroup.name,
              backgroupType: BackendType.Service,
              ports: [
                Object.assign({}, initLbcfBGPort, {
                  id: uuid(),
                  portNumber: backGroup.service.port.portNumber,
                  protocol: backGroup.service.port.protocol
                })
              ],
              labels: keys.map(key => {
                return Object.assign({}, initSelector, {
                  id: uuid(),
                  key: key,
                  value: backGroup.service.nodeSelector[key]
                });
              }),
              serviceName: backGroup.service.name
            })
          );
        } else {
          backGroupEditions.push(
            Object.assign({}, initLbcfBackGroupEdition, {
              id: uuid(),
              name: backGroup.name,
              backgroupType: BackendType.Static,
              staticAddress: backGroup.static
                ? backGroup.static.map(name => {
                    return Object.assign({}, initStringArray, {
                      id: uuid(),
                      value: name
                    });
                  })
                : []
            })
          );
        }
      });
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: backGroupEditions
      });
    };
  },

  inputLbcfBackGroupType: (backGroupId: string, type: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      backGroupEdition.backgroupType = type;
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  inputLbcfBackGroupServiceName: (backGroupId: string, serviceName: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: {
          lbcfEdit: { lbcfBackGroupEditions }
        }
      } = getState();

      let newBackGroupEdition = deepClone(lbcfBackGroupEditions);
      let backGroupEdition = newBackGroupEdition.find(item => item.id === backGroupId);
      backGroupEdition.serviceName = serviceName;
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: newBackGroupEdition
      });
    };
  },

  clearEdition: () => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.ClearLbcfEdit
      });
    };
  }
};
