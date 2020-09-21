import * as ActionType from '../constants/ActionType';
import { EditState, Machine, RootState } from '../models';
import { installerActions } from './installerActions';

type GetState = () => RootState;

export const validateActions = {
  _validateUsername(username: string) {
    let status = 0,
      message = '';

    //验证用户名
    if (!username) {
      status = 2;
      message = '用户名不能为空';
    } else if (!/^[a-z]*$/.test(username)) {
      status = 2;
      message = '用户名只支持小写字母';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateUsername(username?: string) {
    return dispatch => {
      const v_username = validateActions._validateUsername(username);
      dispatch(installerActions.updateEdit({ v_username }));
    };
  },

  _validatePassword(password: string) {
    let status = 0,
      message = '';

    //验证密码
    if (!password) {
      status = 2;
      message = '密码不能为空';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validatePassword(password?: string) {
    return dispatch => {
      const v_password = validateActions._validatePassword(password);
      dispatch(installerActions.updateEdit({ v_password }));
    };
  },

  _validateConfirmPassword(confirmPassword: string, password: string) {
    let status = 0,
      message = '';

    //验证确认密码
    if (!confirmPassword) {
      status = 2;
      message = '确认密码不能为空';
    } else if (confirmPassword !== password) {
      status = 2;
      message = '两次密码输入不一致';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateConfirmPassword(confirmPassword: string, password: string) {
    return dispatch => {
      const v_confirmPassword = validateActions._validateConfirmPassword(confirmPassword, password);
      dispatch(installerActions.updateEdit({ v_confirmPassword }));
    };
  },

  _validateTkeVip(vip: string, haType: string) {
    let status = 0,
      message = '';

    //验证tke vip地址
    if (haType === 'tke') {
      if (!vip) {
        status = 2;
        message = 'VIP地址不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateTkeVip(vip: string, haType: string) {
    return dispatch => {
      const v_haTkeVip = validateActions._validateTkeVip(vip, haType);
      dispatch(installerActions.updateEdit({ v_haTkeVip }));
    };
  },

  _validateThirdVip(vip: string, haType: string) {
    let status = 0,
      message = '';

    //验证tke vip地址
    if (haType === 'thirdParty') {
      if (!vip) {
        status = 2;
        message = 'VIP地址不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateThirdVip(vip: string, haType: string) {
    return dispatch => {
      const v_haThirdVip = validateActions._validateThirdVip(vip, haType);
      dispatch(installerActions.updateEdit({ v_haThirdVip }));
    };
  },

  _validateStep2(editState: EditState) {
    let result =
      validateActions._validateUsername(editState.username).status === 1 &&
      validateActions._validatePassword(editState.password).status === 1 &&
      validateActions._validateConfirmPassword(editState.confirmPassword, editState.password).status === 1 &&
      validateActions._validateTkeVip(editState.haTkeVip, editState.haType).status === 1 &&
      validateActions._validateThirdVip(editState.haThirdVip, editState.haType).status === 1;

    return result;
  },

  validateStep2(editState: EditState) {
    return dispatch => {
      const v_username = validateActions._validateUsername(editState.username),
        v_password = validateActions._validatePassword(editState.password),
        v_confirmPassword = validateActions._validateConfirmPassword(editState.confirmPassword, editState.password),
        v_haTkeVip = validateActions._validateTkeVip(editState.haTkeVip, editState.haType),
        v_haThirdVip = validateActions._validateThirdVip(editState.haThirdVip, editState.haType);

      dispatch(installerActions.updateEdit({ v_username, v_password, v_confirmPassword, v_haTkeVip, v_haThirdVip }));
    };
  },

  _validateNetworkDevice(networkDevice: string) {
    let status = 0,
      message = '';

    //验证网卡名称
    if (!networkDevice) {
      status = 2;
      message = '网卡名称不能为空';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateNetworkDevice(networkDevice?: string) {
    return dispatch => {
      const v_networkDevice = validateActions._validateNetworkDevice(networkDevice);
      dispatch(installerActions.updateEdit({ v_networkDevice }));
    };
  },

  _validateHost(host: string) {
    let reg = /^\d{1,3}(\.\d{1,3}){3}(;\d{1,3}(\.\d{1,3}){3})*$/,
      status = 0,
      message = '';

    //验证访问地址
    if (!host) {
      status = 2;
      message = '访问地址不能为空';
    } else if (!reg.test(host)) {
      status = 2;
      message = '访问地址格式不正确';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateHost(host?: string, key?: string) {
    return dispatch => {
      const v_host = validateActions._validateHost(host);
      dispatch(installerActions.updateMachine({ v_host }, key));
    };
  },

  _validatePort(port: string) {
    let reg = /^\d*$/,
      status = 0,
      message = '';

    //验证访问端口
    if (!port) {
      status = 2;
      message = '访问端口不能为空';
    } else if (!reg.test(port)) {
      status = 2;
      message = '访问端口格式不正确';
    } else if (+port < 1 || +port > 65535) {
      status = 2;
      message = '访问端口范围不正确';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validatePort(port?: string, key?: string) {
    return dispatch => {
      const v_port = validateActions._validatePort(port);
      dispatch(installerActions.updateMachine({ v_port }, key));
    };
  },

  _validateUser(user: string) {
    let status = 0,
      message = '';

    //验证用户名
    if (!user) {
      status = 2;
      message = '用户名不能为空';
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateUser(user?: string, key?: string) {
    return dispatch => {
      const v_user = validateActions._validateUser(user);
      dispatch(installerActions.updateMachine({ v_user }, key));
    };
  },

  _validateMachinePassword(password: string, way?: string) {
    let status = 0,
      message = '';

    //验证密码
    if (way === 'password') {
      if (!password) {
        status = 2;
        message = '密码不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateMachinePassword(password?: string, way?: string, key?: string) {
    return dispatch => {
      const v_password = validateActions._validateMachinePassword(password, way);
      dispatch(installerActions.updateMachine({ v_password }, key));
    };
  },

  _validateCert(cert: string, way?: string) {
    let status = 0,
      message = '';

    //验证证书
    if (way === 'cert') {
      if (!cert) {
        status = 2;
        message = '证书不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateCert(cert?: string, way?: string, key?: string) {
    return dispatch => {
      const v_cert = validateActions._validateCert(cert);
      dispatch(installerActions.updateMachine({ v_cert }, key));
    };
  },

  _validateMachine(machine: Machine) {
    let result =
      validateActions._validateHost(machine.host).status === 1 &&
      validateActions._validatePort(machine.port).status === 1 &&
      validateActions._validateUser(machine.user).status === 1 &&
      validateActions._validateMachinePassword(machine.password, machine.authWay).status === 1 &&
      validateActions._validateCert(machine.cert, machine.authWay).status === 1;

    return result;
  },

  validateMachine(machine: Machine) {
    return dispatch => {
      const v_host = validateActions._validateHost(machine.host),
        v_port = validateActions._validatePort(machine.port),
        v_user = validateActions._validateUser(machine.user),
        v_password = validateActions._validateMachinePassword(machine.password, machine.authWay),
        v_cert = validateActions._validateCert(machine.cert, machine.authWay);

      dispatch(installerActions.updateMachine({ v_host, v_port, v_user, v_password, v_cert }, machine.id));
    };
  },

  _validateAllMachines(machines: Array<Machine>) {
    let result = true;
    machines.forEach(m => {
      result = result && validateActions._validateMachine(m);
    });

    return result;
  },

  validateAllMachines(machines: Array<Machine>) {
    return dispatch => {
      machines.forEach(m => {
        dispatch(validateActions.validateMachine(m));
      });
    };
  },

  _validateStep3(editState: EditState) {
    let result =
      validateActions._validateNetworkDevice(editState.networkDevice).status === 1 &&
      validateActions._validateAllMachines(editState.machines);

    return result;
  },

  validateStep3(editState: EditState) {
    return dispatch => {
      const v_networkDevice = validateActions._validateNetworkDevice(editState.networkDevice);

      dispatch(installerActions.updateEdit({ v_networkDevice }));
      dispatch(validateActions.validateAllMachines(editState.machines));
    };
  },

  _validateTenantID(tenantID: string, authType?: string) {
    let status = 0,
      message = '';

    //验证默认租户名
    if (authType === 'tke') {
      if (!tenantID) {
        status = 2;
        message = '默认租户名不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateTenantID(tenantID?: string, authType?: string) {
    return dispatch => {
      const v_tenantID = validateActions._validateTenantID(tenantID, authType);
      dispatch(installerActions.updateEdit({ v_tenantID }));
    };
  },

  _validateIssuerUrl(issuerUrl: string, authType?: string) {
    let status = 0,
      message = '';

    //验证IssuerUrl
    if (authType === 'oidc') {
      if (!issuerUrl) {
        status = 2;
        message = 'IssuerUrl不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateIssuerUrl(issuerUrl?: string, authType?: string) {
    return dispatch => {
      const v_issuerURL = validateActions._validateIssuerUrl(issuerUrl, authType);
      dispatch(installerActions.updateEdit({ v_issuerURL }));
    };
  },

  _validateClientID(clientID: string, authType?: string) {
    let status = 0,
      message = '';

    //验证ClientID
    if (authType === 'oidc') {
      if (!clientID) {
        status = 2;
        message = 'ClientID不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateClientID(clientID?: string, authType?: string) {
    return dispatch => {
      const v_clientID = validateActions._validateClientID(clientID, authType);
      dispatch(installerActions.updateEdit({ v_clientID }));
    };
  },

  _validateCaCert(caCert: string, authType?: string) {
    let status = 0,
      message = '';

    //验证CaCert
    if (authType === 'oidc') {
      if (!caCert) {
        status = 2;
        message = 'CA证书不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateCaCert(caCert?: string, authType?: string) {
    return dispatch => {
      const v_caCert = validateActions._validateCaCert(caCert, authType);
      dispatch(installerActions.updateEdit({ v_caCert }));
    };
  },

  _validateStep4(editState: EditState) {
    let result =
      validateActions._validateIssuerUrl(editState.issuerURL, editState.authType).status === 1 &&
      validateActions._validateClientID(editState.clientID, editState.authType).status === 1 &&
      validateActions._validateCaCert(editState.caCert, editState.authType).status === 1;

    return result;
  },

  validateStep4(editState: EditState) {
    return dispatch => {
      const v_issuerURL = validateActions._validateIssuerUrl(editState.issuerURL, editState.authType),
        v_clientID = validateActions._validateClientID(editState.clientID, editState.authType),
        v_caCert = validateActions._validateCaCert(editState.caCert, editState.authType);

      dispatch(installerActions.updateEdit({ v_issuerURL, v_clientID, v_caCert }));
    };
  },

  _validateRepoTenantID(repoTenantID: string, repoType?: string) {
    let status = 0,
      message = '';

    //验证默认租户名
    if (repoType === 'tke') {
      if (!repoTenantID) {
        status = 2;
        message = '默认租户名不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateRepoTenantID(repoTenantID?: string, authType?: string) {
    return dispatch => {
      const v_repoTenantID = validateActions._validateRepoTenantID(repoTenantID, authType);
      dispatch(installerActions.updateEdit({ v_repoTenantID }));
    };
  },

  _validateRepoSuffix(repoSuffix: string, repoType?: string) {
    let status = 0,
      message = '';

    //验证域名后缀
    if (repoType === 'tke') {
      if (!repoSuffix) {
        status = 2;
        message = '域名后缀不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateRepoSuffix(repoSuffix?: string, repoType?: string) {
    return dispatch => {
      const v_repoSuffix = validateActions._validateRepoSuffix(repoSuffix, repoType);
      dispatch(installerActions.updateEdit({ v_repoSuffix }));
    };
  },

  _validateRepoAddress(repoAddress: string, repoType?: string) {
    let status = 0,
      message = '';

    //验证仓库地址
    if (repoType === 'thirdParty') {
      if (!repoAddress) {
        status = 2;
        message = '域名后缀不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateRepoAddress(repoAddress?: string, repoType?: string) {
    return dispatch => {
      const v_repoAddress = validateActions._validateRepoAddress(repoAddress, repoType);
      dispatch(installerActions.updateEdit({ v_repoAddress }));
    };
  },

  _validateRepoNamespace(repoNamespace: string, repoType?: string) {
    let status = 0,
      message = '';

    //验证命名空间
    if (repoType === 'thirdParty') {
      if (!repoNamespace) {
        status = 2;
        message = '命名空间不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateRepoNamespace(repoNamespace?: string, repoType?: string) {
    return dispatch => {
      const v_repoNamespace = validateActions._validateRepoNamespace(repoNamespace, repoType);
      dispatch(installerActions.updateEdit({ v_repoNamespace }));
    };
  },

  _validateRepoUser(repoUser: string, repoType?: string) {
    let status = 0,
      message = '';

    //验证命名空间
    if (repoType === 'thirdParty') {
      if (!repoUser) {
        status = 2;
        message = '用户名不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateRepoUser(repoUser?: string, repoType?: string) {
    return dispatch => {
      const v_repoUser = validateActions._validateRepoUser(repoUser, repoType);
      dispatch(installerActions.updateEdit({ v_repoUser }));
    };
  },

  _validateRepoPassword(repoPassword: string, repoType?: string) {
    let status = 0,
      message = '';

    //验证密码
    if (repoType === 'thirdParty') {
      if (!repoPassword) {
        status = 2;
        message = '密码不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateRepoPassword(repoPassword?: string, repoType?: string) {
    return dispatch => {
      const v_repoPassword = validateActions._validateRepoUser(repoPassword, repoType);
      dispatch(installerActions.updateEdit({ v_repoPassword }));
    };
  },

  _validateStep5(editState: EditState) {
    let result =
      validateActions._validateRepoSuffix(editState.repoSuffix, editState.repoType).status === 1 &&
      validateActions._validateRepoAddress(editState.repoAddress, editState.repoType).status === 1 &&
      validateActions._validateRepoNamespace(editState.repoNamespace, editState.repoType).status === 1 &&
      validateActions._validateRepoUser(editState.repoUser, editState.repoType).status === 1 &&
      validateActions._validateRepoPassword(editState.repoPassword, editState.repoType).status === 1;

    return result;
  },

  validateStep5(editState: EditState) {
    return dispatch => {
      const v_repoSuffix = validateActions._validateRepoSuffix(editState.repoSuffix, editState.repoType),
        v_repoAddress = validateActions._validateRepoAddress(editState.repoAddress, editState.repoType),
        v_repoNamespace = validateActions._validateRepoNamespace(editState.repoNamespace, editState.repoType),
        v_repoUser = validateActions._validateRepoUser(editState.repoUser, editState.repoType),
        v_repoPassword = validateActions._validateRepoPassword(editState.repoPassword, editState.repoType);

      dispatch(
        installerActions.updateEdit({
          v_repoSuffix,
          v_repoAddress,
          v_repoNamespace,
          v_repoUser,
          v_repoPassword
        })
      );
    };
  },

  _validateStep6(editState: EditState) {
    let result = true;
    if (editState.openAudit) {
      result =
        validateActions._validateESUrl(editState.auditEsUrl, 'es').status === 1 &&
        validateActions._validateESUsername(editState.auditEsUsername, 'es').status === 1 &&
        validateActions._validateESPassword(editState.auditEsPassword, 'es').status === 1;
    }

    return result;
  },

  validateStep6(editState: EditState) {
    return dispatch => {
      const v_auditEsUrl = validateActions._validateESUrl(editState.auditEsUrl, 'es'),
        v_auditEsUsername = validateActions._validateESUsername(editState.auditEsUsername, 'es'),
        v_auditEsPassword = validateActions._validateESPassword(editState.auditEsPassword, 'es');

      dispatch(
        installerActions.updateEdit({
          v_auditEsUrl,
          v_auditEsUsername,
          v_auditEsPassword
        })
      );
    };
  },

  _validateESUrl(esUrl: string, monitorType?: string) {
    let status = 0,
      message = '';

    //验证ES地址
    if (monitorType === 'es') {
      if (!esUrl) {
        status = 2;
        message = 'ES地址不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateESUrl(esUrl?: string, monitorType?: string) {
    return dispatch => {
      const v_esUrl = validateActions._validateESUrl(esUrl, monitorType);
      dispatch(installerActions.updateEdit({ v_esUrl }));
    };
  },

  _validateESUsername(esUsername: string, monitorType?: string) {
    let status = 0,
      message = '';

    //验证ES用户名
    if (monitorType === 'es') {
      if (!esUsername) {
        status = 1;
        message = '';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateESUsername(esUsername?: string, monitorType?: string) {
    return dispatch => {
      const v_esUsername = validateActions._validateESUsername(esUsername, monitorType);
      dispatch(installerActions.updateEdit({ v_esUsername }));
    };
  },

  _validateESPassword(esPassword: string, monitorType?: string) {
    let status = 0,
      message = '';

    //验证ES密码
    if (monitorType === 'es') {
      if (!esPassword) {
        status = 1;
        message = '';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateESPassword(esPassword?: string, monitorType?: string) {
    return dispatch => {
      const v_esPassword = validateActions._validateESPassword(esPassword, monitorType);
      dispatch(installerActions.updateEdit({ v_esPassword }));
    };
  },

  _validateInfluxDBUrl(influxDBUrl: string, monitorType?: string) {
    let status = 0,
      message = '';

    //验证InfluxDB地址
    if (monitorType === 'external-influxdb') {
      if (!influxDBUrl) {
        status = 2;
        message = 'InfluxDB地址不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateInfluxDBUrl(influxDBUrl?: string, monitorType?: string) {
    return dispatch => {
      const v_influxDBUrl = validateActions._validateInfluxDBUrl(influxDBUrl, monitorType);
      dispatch(installerActions.updateEdit({ v_influxDBUrl }));
    };
  },

  _validateInfluxDBUsername(influxDBUsername: string, monitorType?: string) {
    let status = 0,
      message = '';

    //验证InfluxDB用户名
    if (monitorType === 'external-influxdb') {
      if (!influxDBUsername) {
        status = 2;
        message = 'InfluxDB用户名不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateInfluxDBUsername(influxDBUsername?: string, monitorType?: string) {
    return dispatch => {
      const v_influxDBUsername = validateActions._validateInfluxDBUsername(influxDBUsername, monitorType);
      dispatch(installerActions.updateEdit({ v_influxDBUsername }));
    };
  },

  _validateInfluxDBPassword(influxDBPassword: string, monitorType?: string) {
    let status = 0,
      message = '';

    //验证InfluxDB密码
    if (monitorType === 'external-influxdb') {
      if (!influxDBPassword) {
        status = 2;
        message = 'InfluxDB密码不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
    }

    return { status, message };
  },

  validateInfluxDBPassword(influxDBPassword?: string, monitorType?: string) {
    return dispatch => {
      const v_influxDBPassword = validateActions._validateInfluxDBPassword(influxDBPassword, monitorType);
      dispatch(installerActions.updateEdit({ v_influxDBPassword }));
    };
  },

  _validateStep7(editState: EditState) {
    let result =
      validateActions._validateESUrl(editState.esUrl, editState.monitorType).status === 1 &&
      validateActions._validateESUsername(
        editState.esUsername,
        editState.esUsername || editState.esPassword ? editState.monitorType : null
      ).status === 1 &&
      validateActions._validateESPassword(
        editState.esPassword,
        editState.esUsername || editState.esPassword ? editState.monitorType : null
      ).status === 1 &&
      validateActions._validateInfluxDBUrl(editState.influxDBUrl, editState.monitorType).status === 1 &&
      validateActions._validateInfluxDBUsername(editState.influxDBUsername, editState.monitorType).status === 1 &&
      validateActions._validateInfluxDBPassword(editState.influxDBPassword, editState.monitorType).status === 1;

    return result;
  },

  validateStep7(editState: EditState) {
    return dispatch => {
      const v_esUrl = validateActions._validateESUrl(editState.esUrl, editState.monitorType),
        v_esUsername = validateActions._validateESUsername(editState.esUsername, editState.monitorType),
        v_esPassword = validateActions._validateESPassword(editState.esPassword, editState.monitorType),
        v_influxDBUrl = validateActions._validateInfluxDBUrl(editState.influxDBUrl, editState.monitorType),
        v_influxDBUsername = validateActions._validateInfluxDBUsername(
          editState.influxDBUsername,
          editState.monitorType
        ),
        v_influxDBPassword = validateActions._validateInfluxDBPassword(
          editState.influxDBPassword,
          editState.monitorType
        );

      dispatch(
        installerActions.updateEdit({
          v_esUrl,
          v_esUsername,
          v_esPassword,
          v_influxDBUrl,
          v_influxDBUsername,
          v_influxDBPassword
        })
      );
    };
  },

  _validateDomain(domain: string, openConsole: boolean) {
    let reg = /^(([a-zA-Z0-9_-])+(\.)?)*$/,
      status = 0,
      message = '';

    //验证企业域名
    if (openConsole) {
      if (!domain) {
        status = 2;
        message = '域名后缀不能为空';
      } else if (!reg.test(domain)) {
        status = 2;
        message = '域名后缀格式不正确';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateDomain(domain: string, openConsole: boolean) {
    return dispatch => {
      const v_consoleDomain = validateActions._validateDomain(domain, openConsole);
      dispatch(installerActions.updateEdit({ v_consoleDomain }));
    };
  },

  _validateCertificate(certificate: string, openConsole: boolean, certType: string) {
    let status = 0,
      message = '';

    //验证自有证书
    if (openConsole && certType === 'thirdParty') {
      if (!certificate) {
        status = 2;
        message = '证书不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateCertificate(certificate: string, openConsole: boolean, certType: string) {
    return dispatch => {
      const v_certificate = validateActions._validateCertificate(certificate, openConsole, certType);
      dispatch(installerActions.updateEdit({ v_certificate }));
    };
  },

  _validatePrivateKey(privateKey: string, openConsole: boolean, certType: string) {
    let status = 0,
      message = '';

    //验证私钥
    if (openConsole && certType === 'thirdParty') {
      if (!privateKey) {
        status = 2;
        message = '私钥不能为空';
      } else {
        status = 1;
        message = '';
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validatePrivateKey(privateKey: string, openConsole: boolean, certType: string) {
    return dispatch => {
      const v_privateKey = validateActions._validateCertificate(privateKey, openConsole, certType);
      dispatch(installerActions.updateEdit({ v_privateKey }));
    };
  },

  _validateStep8(editState?: EditState) {
    const result =
      validateActions._validateCertificate(editState.certificate, editState.openConsole, editState.certType).status ===
        1 &&
      validateActions._validatePrivateKey(editState.privateKey, editState.openConsole, editState.certType).status === 1;

    return result;
  },

  validateStep8(editState?: EditState) {
    return dispatch => {
      const v_certificate = validateActions._validateCertificate(
          editState.certificate,
          editState.openConsole,
          editState.certType
        ),
        v_privateKey = validateActions._validateCertificate(
          editState.privateKey,
          editState.openConsole,
          editState.certType
        );
      dispatch(
        installerActions.updateEdit({
          v_certificate,
          v_privateKey
        })
      );
    };
  }
};
