import * as ActionType from '../constants/ActionType';
import { RootState, SecretData, SecretEdit } from '../models';
import { cloneDeep } from '../../common/utils';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

type GetState = () => RootState;

export const validateSecretActions = {
  /** 校验名称是否正确 */
  _validateName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';
    // 验证secret名称
    if (!name) {
      status = 2;
      message = t('Secret名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('Secret不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('Secret名称格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateName() {
    return async (dispatch, getState: GetState) => {
      let { name } = getState().subRoot.secretEdit;

      const result = validateSecretActions._validateName(name);
      dispatch({
        type: ActionType.SecV_Name,
        payload: result
      });
    };
  },

  /** 校验变量名是否正确 */
  _validateKeyName(keyName: string, secretData: SecretData[]) {
    let status = 0,
      message = '',
      reg = /^[a-z]([-a-z0-9_\.]*[a-z0-9])?$/;

    if (!keyName) {
      status = 2;
      message = t('请输入变量名');
    } else if (!reg.test(keyName)) {
      status = 2;
      message = t('变量名格式不正确');
    } else if (secretData.filter(s => s.keyName === keyName).length > 1) {
      status = 2;
      message = t('变量名不可重复');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateKeyName(keyName: string, dataId: string) {
    return async (dispatch, getState: GetState) => {
      let secretData: SecretData[] = cloneDeep(getState().subRoot.secretEdit.data),
        dIndex = secretData.findIndex(item => item.id === dataId),
        result = validateSecretActions._validateKeyName(keyName, secretData);

      secretData[dIndex]['v_keyName'] = result;
      dispatch({
        type: ActionType.Sec_UpdateData,
        payload: secretData
      });
    };
  },

  _validateAllKeyName(secretData: SecretData[]) {
    let result = true;
    secretData.forEach(item => {
      result = result && validateSecretActions._validateKeyName(item.keyName, secretData).status === 1;
    });
    return result;
  },

  validateAllKeyName() {
    return async (dispatch, getState: GetState) => {
      getState().subRoot.secretEdit.data.forEach(d => {
        dispatch(validateSecretActions.validateKeyName(d.keyName, d.id + ''));
      });
    };
  },

  /** 校验变量值是否正确 */
  _validateKeyValue(value: string) {
    let status = 0,
      message = '',
      reg = /^[a-z]([-a-z0-9_\.]*[a-z0-9])?$/;

    if (!value) {
      status = 2;
      message = t('请输入变量值');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateKeyValue(value: string, dataId: string) {
    return async (dispatch, getState: GetState) => {
      let secretData: SecretData[] = cloneDeep(getState().subRoot.secretEdit.data),
        dIndex = secretData.findIndex(item => item.id === dataId),
        result = validateSecretActions._validateKeyValue(value);

      secretData[dIndex]['v_value'] = result;
      dispatch({
        type: ActionType.Sec_UpdateData,
        payload: secretData
      });
    };
  },

  _validateAllKeyValue(secretData: SecretData[]) {
    let result = true;
    secretData.forEach(item => {
      result = result && validateSecretActions._validateKeyValue(item.value).status === 1;
    });
    return result;
  },

  validateAllKeyValue() {
    return async (dispatch, getState: GetState) => {
      getState().subRoot.secretEdit.data.forEach(d => {
        dispatch(validateSecretActions.validateKeyValue(d.value, d.id + ''));
      });
    };
  },

  /** 校验第三方镜像仓库域名 */
  _validateThirdHubDomain(domain: string) {
    let reg = /^((https|http):\/\/)?(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$/i,
      status = 0,
      message = '';
    if (!domain) {
      status = 2;
      message = t('仓库域名不能为空');
    } else if (domain.length > 128) {
      status = 2;
      message = t('仓库域名长度不能超过128个字符');
    } else if (!reg.test(domain)) {
      status = 2;
      message = t('仓库域名格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateThirdHubDomain() {
    return async (dispatch, getState: GetState) => {
      let domain = getState().subRoot.secretEdit.domain;
      let result = validateSecretActions._validateThirdHubDomain(domain);
      dispatch({
        type: ActionType.SecV_Domain,
        payload: result
      });
    };
  },

  /** 校验第三方镜像仓库的用户名 */
  _validateThirdHubUsername(username: string) {
    let status = 0,
      message = '';
    if (!username) {
      status = 2;
      message = t('用户名不能为空');
    } else if (username.length > 64) {
      status = 2;
      message = t('用户名长度不能超过64个字符');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateThirdHubUsername() {
    return (dispatch, getState: GetState) => {
      let username = getState().subRoot.secretEdit.username;
      let result = validateSecretActions._validateThirdHubUsername(username);
      dispatch({
        type: ActionType.SecV_Username,
        payload: result
      });
    };
  },

  /** 校验第三方镜像仓库的密码 */
  _validateThirdHubPassword(password: string) {
    let status = 0,
      message = '';
    if (!password) {
      status = 2;
      message = t('密码不能为空');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateThirdHubPassword() {
    return (dispatch, getState: GetState) => {
      let password = getState().subRoot.secretEdit.password;
      let result = validateSecretActions._validateThirdHubPassword(password);
      dispatch({
        type: ActionType.SecV_Password,
        payload: result
      });
    };
  },

  /** 校验secret表单 */
  _validateSecretEdit(secrecEdit: SecretEdit) {
    let result = true;
    result = result && validateSecretActions._validateName(secrecEdit.name).status === 1;

    if (secrecEdit.secretType === 'Opaque') {
      result =
        result &&
        validateSecretActions._validateAllKeyName(secrecEdit.data) &&
        validateSecretActions._validateAllKeyValue(secrecEdit.data);
    } else if (secrecEdit.secretType === 'kubernetes.io/dockercfg') {
      result =
        result &&
        validateSecretActions._validateThirdHubDomain(secrecEdit.domain).status === 1 &&
        validateSecretActions._validateThirdHubUsername(secrecEdit.username).status === 1 &&
        validateSecretActions._validateThirdHubPassword(secrecEdit.password).status === 1;
    }

    return result;
  },

  validateSecretEdit() {
    return async (dispatch, getState: GetState) => {
      let { secretType } = getState().subRoot.secretEdit;

      dispatch(validateSecretActions.validateName());

      if (secretType === 'Opaque') {
        dispatch(validateSecretActions.validateAllKeyName());
        dispatch(validateSecretActions.validateAllKeyValue());
      } else if (secretType === 'kubernetes.io/dockercfg') {
        dispatch(validateSecretActions.validateThirdHubDomain());
        dispatch(validateSecretActions.validateThirdHubUsername());
        dispatch(validateSecretActions.validateThirdHubPassword());
      }
    };
  }
};
