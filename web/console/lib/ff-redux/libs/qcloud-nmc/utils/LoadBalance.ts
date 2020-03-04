export interface VipInfoList {
  vip: string;
  ispId: number;
}

export interface HealthStatus {
  ip: string;
  protocol: string;
  port: string;
  vport: number;
  healthStatus: number;
}

export interface BindDevice {
  deviceLanIp: string;
  deviceWanIp: string;
  uuid: string;
  alias: string;
  vpcId: number;
}

export interface LoadBalance {
  id: number;
  appId: number;
  LBId: number;
  LBName: string;
  LBType: string;
  LBDomain: string;
  vip: string[];
  status: number;
  projectId: number;
  sessionExpire: number;
  vpcId: number;
  subnetId: number;
  desState: number;
  addTimestamp: string;
  vipInfoList: VipInfoList[];
  healthStatus: HealthStatus[];
  bindDevice: BindDevice[];
  netName: string;
}

interface LBResponse {
  code: number;
  data: {
    LBList: LoadBalance[];
    page: number;
    count: number;
    totalNum: number;
  };
}

const net = seajs.require('net');

async function fetchLoadBalanceList(data: any): Promise<LoadBalance[]> {
  return new Promise<LoadBalance[]>((resolve, reject) => {
    net.send(
      { method: 'GET', url: '/cgi/lb?action=queryLB' },
      {
        data,
        cb: (response: LBResponse) => {
          if (response.code === 0) {
            const list = response.data.LBList;
            list.forEach(x => {
              x.id = x.LBId;
            });
            resolve(list);
          } else {
            reject(new Error('拉取负载均衡失败：' + JSON.stringify(response)));
          }
        }
      }
    );
  });
}

export async function loadLoadLalanceList(regionId: string | number, projectId: number) {
  return fetchLoadBalanceList({
    count: 999999,
    page: 1,
    projectId,
    regionId
  });
}
