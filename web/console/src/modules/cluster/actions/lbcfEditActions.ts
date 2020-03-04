import { KeyValue } from 'src/modules/common';

import { deepClone, uuid } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config/resourceConfig';
import * as ActionType from '../constants/ActionType';
import { initLbcfBackGroupEdition, initLbcfBGPort, initSelector } from '../constants/initState';
import { RootState } from '../models';
import { Namespace } from '../models/Namespace';
import { ResourceFilter } from '../models/ResourceOption';
import { CLB, Selector } from '../models/ServiceEdit';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { validateLbcfActions } from './validateLbcfActions';

type GetState = () => RootState;

export const lbcfEditActions = {
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

  initGameBGEdition: backGroups => {
    return async (dispatch, getState: GetState) => {
      let backGroupEditions = [];
      backGroups.forEach(backGroup => {
        let keys = Object.keys(backGroup.labels);
        backGroupEditions.push(
          Object.assign({}, initLbcfBackGroupEdition, {
            id: uuid(),
            name: backGroup.name,
            v_name: { status: 1, message: '' },
            ports: [
              Object.assign({}, initLbcfBGPort, {
                id: uuid(),
                portNumber: backGroup.port.portNumber,
                protocol: backGroup.port.protocol
              })
            ],
            labels: keys.map(key => {
              return Object.assign({}, initSelector, {
                id: uuid(),
                key: key,
                value: backGroup.labels[key]
              });
            })
          })
        );
      });
      dispatch({
        type: ActionType.GBG_UpdateLbcfBackGroup,
        payload: backGroupEditions
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
