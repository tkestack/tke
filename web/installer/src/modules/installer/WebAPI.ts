import axios from 'axios';
import { Base64 } from 'js-base64';

import { OperationResult } from '@tencent/ff-redux';

import { Record } from '../common/models';
import { EditState } from './models';

const host = location.host;
export async function fetchCluster() {
  let cluster = {};
  try {
    const rsp = await axios.get(`http://${host}/api/cluster/global`);

    if (rsp.data) {
      cluster = rsp.data;
    }
  } catch (e) {
    if (e.response.data.code === 400) {
      cluster = {};
    }
  }

  let progress = await fetchProgress();

  const result: Record<any> = {
    record: {
      config: cluster,
      progress: progress.record
    }
  };

  return result;
}

export async function verifyLicense(license: string) {
  let verify = {};
  try {
    const rsp = await axios.post(`http://${host}/api/license/verify`, {
      license
    });

    if (!rsp.data.code) {
      verify = rsp.data;
    }
  } catch (e) {
    verify = {};
  }

  return verify;
}

// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error: error ? error.message : '' }));
  }
  return [{ success: !error, target: target as T, error: error ? error.message : '' }];
}

export async function createCluster(edits: Array<EditState>) {
  // 先把host用分号分割的展平为多个

  edits[0].machines = edits[0].machines.reduce(
    (all, { host, ...restMachine }) => all.concat(host.split(';').map(ip => ({ host: ip, ...restMachine }))),
    []
  );

  try {
    // 先校验机器连通性
    let requests = edits[0].machines.map(m => {
      const req = async () => {
        let params = {
          host: m.host,
          port: +m.port,
          user: m.user,
          password: Base64.encode(m.password)
        };

        if (m.authWay === 'password') {
          params['password'] = Base64.encode(m.password);
        } else {
          if (m.password) {
            params['privatePassword'] = Base64.encode(m.password);
          }
          params['privateKey'] = Base64.encode(m.cert);
        }

        let rsp = await axios.post(`http://${host}/api/ssh`, params);
        return rsp;
      };

      return req()
        .then(
          res => {
            if (res['data']['status'] === 'OK') {
              return {
                [m.host]: {
                  code: 0,
                  message: 'OK'
                }
              };
            } else {
              return {
                [m.host]: {
                  code: res['response']['data'].code,
                  message: res['response']['data'].message
                }
              };
            }
          },
          res => {
            return {
              [m.host]: {
                code: res['response']['data'].code,
                message: res['response']['data'].message
              }
            };
          }
        )
        .catch(error => {
          return {
            [m.host]: {
              code: -1,
              message: error.message
            }
          };
        });
    });

    let results = await Promise.all(requests),
      re = results.reduce((prev, next) => Object.assign({}, prev, next), {});
    let isVerified = true,
      error = null;
    for (let key in re) {
      if (re[key]['code'] !== 0) {
        isVerified = false;
        error = re[key];
      }
    }

    if (!isVerified) {
      return operationResult(edits, error);
    } else {
      let machines = [];
      edits[0].machines.forEach(m => {
        m.host.split(';').forEach(ip => {
          let mac = {
            ip: ip,
            port: +m.port,
            username: m.user
          };
          if (m.authWay === 'password') {
            mac['password'] = Base64.encode(m.password);
          } else {
            if (m.password) {
              mac['privatePassword'] = Base64.encode(m.password);
            }
            mac['privateKey'] = Base64.encode(m.cert);
          }
          machines.push(mac);
        });
      });
      let params = {
        cluster: {
          apiVersion: 'platform.tkestack.io/v1',
          kind: 'Cluster',
          spec: {
            networkDevice: edits[0].networkDevice,
            features: {
              enableMetricsServer: true
            },
            dockerExtraArgs: edits[0].dockerExtraArgs.reduce((prev, next) => {
              if (next.key) {
                return Object.assign({}, prev, { [next.key]: next.value });
              } else {
                return Object.assign({}, prev);
              }
            }, {}),
            kubeletExtraArgs: edits[0].kubeletExtraArgs.reduce((prev, next) => {
              if (next.key) {
                return Object.assign({}, prev, { [next.key]: next.value });
              } else {
                return Object.assign({}, prev);
              }
            }, {}),
            apiServerExtraArgs: edits[0].apiServerExtraArgs.reduce((prev, next) => {
              if (next.key) {
                return Object.assign({}, prev, { [next.key]: next.value });
              } else {
                return Object.assign({}, prev);
              }
            }, {}),
            controllerManagerExtraArgs: edits[0].controllerManagerExtraArgs.reduce((prev, next) => {
              if (next.key) {
                return Object.assign({}, prev, { [next.key]: next.value });
              } else {
                return Object.assign({}, prev);
              }
            }, {}),
            schedulerExtraArgs: edits[0].schedulerExtraArgs.reduce((prev, next) => {
              if (next.key) {
                return Object.assign({}, prev, { [next.key]: next.value });
              } else {
                return Object.assign({}, prev);
              }
            }, {}),
            clusterCIDR: edits[0].cidr,
            properties: {
              maxClusterServiceNum: +edits[0].serviceNumLimit,
              maxNodePodNum: +edits[0].podNumLimit
            },
            type: 'Baremetal',
            machines: machines
          }
        },
        config: {
          basic: {
            username: edits[0].username,
            password: Base64.encode(edits[0].password)
          }
        }
      };

      // GPU设置
      if (edits[0].gpuType !== 'none') {
        params.cluster.spec.features['gpuType'] = edits[0].gpuType;
      }

      // 认证模块设置
      if (edits[0].authType === 'tke') {
        params.config['auth'] = {
          tke: {}
        };
      } else if (edits[0].authType === 'oidc') {
        params.config['auth'] = {
          oidc: {
            issuerURL: edits[0].issuerURL,
            clientID: edits[0].clientID,
            caCert: Base64.encode(edits[0].caCert)
          }
        };
      } else {
        // 未开启认证
      }

      // 镜像模块设置
      if (edits[0].repoType === 'tke') {
        params.config['registry'] = {
          tke: {
            domain: edits[0].repoSuffix
          }
        };
      } else {
        params.config['registry'] = {
          thirdParty: {
            domain: edits[0].repoAddress,
            namespace: edits[0].repoNamespace,
            username: edits[0].repoUser,
            password: Base64.encode(edits[0].repoPassword)
          }
        };
      }

      // 应用商店
      if (edits[0].application) {
        params.config['application'] = {};
      }

      // 业务模块设置
      if (edits[0].openBusiness) {
        params.config['business'] = {};
      } else {
        //未开启业务模块
      }

      if (edits[0].openAudit) {
        params.config['audit'] = {
          elasticSearch: {
            address: edits[0].auditEsUrl,
            username: edits[0].auditEsUsername ? edits[0].auditEsUsername : undefined,
            password: edits[0].auditEsPassword ? Base64.encode(edits[0].auditEsPassword) : undefined,
            reserveDays: edits[0].auditEsReserveDays
          }
        };
        params.config['audit'] = JSON.parse(JSON.stringify(params.config['audit']));
      }

      // 监控模块设置
      if (edits[0].monitorType === 'es') {
        params.config['monitor'] = {
          es: {
            url: edits[0].esUrl,
            username: edits[0].esUsername,
            password: Base64.encode(edits[0].esPassword)
          }
        };
      } else if (edits[0].monitorType === 'external-influxdb') {
        params.config['monitor'] = {
          influxDB: {
            external: {
              url: edits[0].influxDBUrl,
              username: edits[0].influxDBUsername,
              password: Base64.encode(edits[0].influxDBPassword)
            }
          }
        };
      } else if (edits[0].monitorType === 'tke-influxdb') {
        params.config['monitor'] = {
          influxDB: {
            local: {}
          }
        };
      } else {
        // 未开启监控模块
      }

      params.config['logagent'] = {};

      // 高可用设置
      if (edits[0].haType === 'tke') {
        params.config['ha'] = {
          tke: {
            vip: edits[0].haTkeVip
          }
        };
      } else if (edits[0].haType === 'thirdParty') {
        params.config['ha'] = {
          thirdParty: {
            vip: edits[0].haThirdVip,
            vport: +edits[0].haThirdVipPort || '6443'
          }
        };
      } else {
        // 不高可用
      }

      //控制台设置
      if (edits[0].openConsole) {
        params.config['gateway'] = {
          domain: edits[0].consoleDomain
        };
        if (edits[0].certType === 'selfSigned') {
          params.config['gateway']['cert'] = {
            selfSigned: {}
          };
        } else {
          params.config['gateway']['cert'] = {
            thirdParty: {
              certificate: Base64.encode(edits[0].certificate),
              privateKey: Base64.encode(edits[0].privateKey)
            }
          };
        }
      }

      const rsp = await axios.post(`http://${host}/api/cluster`, params);
      if (rsp.data) {
        return operationResult(edits);
      } else {
        return operationResult(edits, rsp.data);
      }
    }
  } catch (e) {
    return operationResult(edits, e.response.data);
  }
}

export async function fetchProgress() {
  let re = {};
  try {
    const rsp = await axios.get(`http://${host}/api/cluster/global/progress`);
    //const rsp = await axios.get(`http://${host}/api/cluster/global/retry`);
    if (rsp.data) {
      re = rsp.data;
    }
  } catch (e) {
    if (e.response.data.code === 400) {
      re = {};
    }
  }

  const result: Record<any> = {
    record: re
  };

  return result;
}
