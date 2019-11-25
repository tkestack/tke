import { ComputerLabel, ComputerTaint } from './../models/Computer';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import { validateServiceActions } from './validateServiceActions';
import { validateNamespaceActions } from './validateNamespaceActions';
import { validateWorkloadActions } from './validateWorkloadActions';
import { validateSecretActions } from './validateSecretActinos';
import { validateCMActions } from './validateCMActions';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { validateClusterCreationAction } from './validateClusterCreationAction';
import { validateCreateICAction } from './validateCreateICAction';
import { validateLbcfActions } from './validateLbcfActions';
import { AllocationRatioEdition } from '../models/AllocationRatioEdition';

type GetState = () => RootState;

export const validatorActions = {
  service: validateServiceActions,
  namespace: validateNamespaceActions,
  workload: validateWorkloadActions,
  secret: validateSecretActions,
  cm: validateCMActions,
  clusterCreation: validateClusterCreationAction,
  createIC: validateCreateICAction,
  lbcf: validateLbcfActions,

  _validateComputerLabelValue(value, type, disabled) {
    let reg = /^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/,
      status = 0,
      message = '',
      content = type === 'key' ? t('Label名') : t('Label值');
    if (disabled) {
      status = 1;
      message = '';
    } else {
      if (!value) {
        status = 2;
        message = content + t('不能为空');
      } else if (value.length > 63) {
        status = 2;
        message = content + t('长度不能超过63个字符');
      } else if (!reg.test(value)) {
        status = 2;
        message = content + t('格式不正确');
      } else if (type === 'key' && value.includes('kubernetes')) {
        status = 2;
        message = content + t('格式不正确');
      } else {
        status = 1;
        message = '';
      }
    }

    return { status, message };
  },
  validateComputerLabelValue(Id) {
    return (dispatch, getState) => {
      let { labelEdition } = getState().subRoot.computerState,
        { labels } = labelEdition,
        eIndex = labels.findIndex(e => e.id === Id);
      let result = validatorActions._validateComputerLabelValue(labels[eIndex].value, 'value', labels[eIndex].disabled);
      labels[eIndex].v_value = result;
      dispatch({
        type: ActionType.UpdateLabelEdition,
        payload: Object.assign({}, labelEdition, { labels })
      });
    };
  },
  validateComputerLabelKey(Id) {
    return (dispatch, getState) => {
      let { labelEdition } = getState().subRoot.computerState,
        { labels } = labelEdition,
        eIndex = labels.findIndex(e => e.id === Id);
      let result = validatorActions._validateComputerLabelValue(labels[eIndex].key, 'key', labels[eIndex].disabled);
      labels[eIndex].v_key = result;
      dispatch({
        type: ActionType.UpdateLabelEdition,
        payload: Object.assign({}, labelEdition, { labels })
      });
    };
  },
  _validateAllComputerLabel(labels: ComputerLabel[]) {
    let result = true;
    labels.forEach(label => {
      result =
        result &&
        validatorActions._validateComputerLabelValue(label.key, 'key', label.disabled).status === 1 &&
        validatorActions._validateComputerLabelValue(label.value, 'value', label.disabled).status === 1;
    });
    return result;
  },
  validateAllComputerLabel() {
    return (dispatch, getState) => {
      let { labelEdition } = getState().subRoot.computerState,
        { labels } = labelEdition;
      labels.forEach(label => {
        dispatch(validatorActions.validateComputerLabelKey(label.id));
        dispatch(validatorActions.validateComputerLabelValue(label.id));
      });
    };
  },
  _validateComputerTaintValue(value, type, disabled) {
    let reg = /^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/,
      status = 0,
      message = '',
      content = type === 'key' ? t('Taint名') : t('Taint值');
    if (disabled) {
      status = 1;
      message = '';
    } else if (type === 'value' && value === '') {
      //taint 值可以为空
      status = 1;
      message = '';
    } else {
      if (!value) {
        status = 2;
        message = content + t('不能为空');
      } else if (value.length > 63) {
        status = 2;
        message = content + t('长度不能超过63个字符');
      } else if (!reg.test(value)) {
        status = 2;
        message = content + t('格式不正确');
      } else if (type === 'key' && value.includes('kubernetes')) {
        status = 2;
        message = content + t('格式不正确');
      } else {
        status = 1;
        message = '';
      }
    }

    return { status, message };
  },
  validateComputerTaintValue(Id) {
    return (dispatch, getState: GetState) => {
      let { taintEdition } = getState().subRoot.computerState,
        { taints } = taintEdition,
        eIndex = taints.findIndex(e => e.id === Id);
      let result = validatorActions._validateComputerTaintValue(taints[eIndex].value, 'value', taints[eIndex].disabled);
      taints[eIndex].v_value = result;
      dispatch({
        type: ActionType.UpdateTaintEdition,
        payload: Object.assign({}, taintEdition, { taints })
      });
    };
  },
  validateComputerTaintKey(Id) {
    return (dispatch, getState) => {
      let { taintEdition } = getState().subRoot.computerState,
        { taints } = taintEdition,
        eIndex = taints.findIndex(e => e.id === Id);
      let result = validatorActions._validateComputerTaintValue(taints[eIndex].key, 'key', taints[eIndex].disabled);
      taints[eIndex].v_key = result;
      dispatch({
        type: ActionType.UpdateTaintEdition,
        payload: Object.assign({}, taintEdition, { taints })
      });
    };
  },
  _validateAllComputerTaint(taints: ComputerTaint[]) {
    let result = true;
    taints.forEach(taint => {
      result =
        result &&
        validatorActions._validateComputerTaintValue(taint.key, 'key', taint.disabled).status === 1 &&
        validatorActions._validateComputerTaintValue(taint.value, 'value', taint.disabled).status === 1;
    });
    return result;
  },
  validateAllComputerTaint() {
    return (dispatch, getState) => {
      let { taintEdition } = getState().subRoot.computerState,
        { taints } = taintEdition;
      taints.forEach(taint => {
        dispatch(validatorActions.validateComputerTaintKey(taint.id));
        dispatch(validatorActions.validateComputerTaintValue(taint.id));
      });
    };
  },

  _validateClusterAllocationRatio(value: number) {
    let reg = /^\d+(\.\d{1,3})?$/,
      status = 0,
      message = '';

    if (isNaN(value)) {
      status = 2;
      message = t('数据格式不正确，超售比只能是数字，且只能精确到0.01');
    } else if (value === 0) {
      status = 2;
      message = t('超售比不能为空或者不能为0');
    } else if (!reg.test(value + '')) {
      status = 2;
      message = t('数据格式不正确，超售比只能是数字，且只能精确到0.01');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },
  validateClusterAllocationRatio(type?: string, value?: string) {
    return (dispatch, getState: GetState) => {
      let { clusterAllocationRatioEdition } = getState().subRoot;
      let result = validatorActions._validateClusterAllocationRatio(+value);
      let obj;
      if (type === 'cpu') {
        obj = {
          v_cpuRatio: result
        };
      } else {
        obj = {
          v_memoryRatio: result
        };
      }
      dispatch({
        type: ActionType.UpdateClusterAllocationRatioEdition,
        payload: Object.assign({}, clusterAllocationRatioEdition, obj)
      });
    };
  },
  validateAllClusterAllocationRatio() {
    return (dispatch, getState: GetState) => {
      let {
        clusterAllocationRatioEdition: { isUseCpu, isUseMemory, cpuRatio, memoryRatio }
      } = getState().subRoot;
      isUseCpu && dispatch(validatorActions.validateClusterAllocationRatio('cpu', cpuRatio));
      isUseMemory && dispatch(validatorActions.validateClusterAllocationRatio('memory', memoryRatio));
    };
  },
  _validateAllClusterAllocationRatio(clusterAllocationRatioEdition: AllocationRatioEdition) {
    let result = true;
    let { isUseCpu, isUseMemory, memoryRatio, cpuRatio } = clusterAllocationRatioEdition;
    if (isUseCpu) {
      result = result && validatorActions._validateClusterAllocationRatio(+cpuRatio).status === 1;
    }
    if (isUseMemory) {
      result = result && validatorActions._validateClusterAllocationRatio(+memoryRatio).status === 1;
    }
    return result;
  }
};
