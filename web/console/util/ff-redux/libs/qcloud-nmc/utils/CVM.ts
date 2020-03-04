export interface Subnet {
  [index: number]: string;
}

export interface RelatedVpc {
  name: string;
  subnet: Subnet;
  vpcCIDR: string;
  unVpcId: string;
}

export interface CVM {
  id: string;
  appId: number;
  deviceId: number;
  tmpDeviceId: number;
  deviceAssetId: string;
  deviceLanIp: string;
  deviceWanIp: string;
  localEthDev: string;
  idcName: string;
  ispName: string;
  deviceImageId: number;
  uuid: string;
  deviceClass: string;
  osName: string;
  cpu: number;
  mem: number;
  disk: number;
  bandwidth: number;
  rootSize: number;
  swapSize: number;
  addTimeStamp: string;
  alias: string;
  runflag: number;
  status: number;
  lastOperation: string;
  updateTimestamp: string;
  deadline: string;
  isSafeIsolated: number;
  safeIsolatedInfo: string;
  projectId: number;
  deviceClassFlag: number;
  autoRenewFlag: number;
  isolateTime: string;
  customId: number;
  hypervisorUpdateFlag: number;
  vpcId: number;
  billId: any;
  rootType: number;
  zoneId: number;
  zoneName: string;
  subnetId: number;
  isVpcGateway: number;
  diskType: number;
  itemId: number;
  cvmPayMode: number;
  networkPayMode: number;
  uInstanceId: string;
  unImgId: string;
  imageType: string;
  cbsList: any[];
  hypervisor: number;
  operateAuth: number;
  subnet: string;
  relatedVpcId: number;
  statusDesc: string;
  deviceTypeDesc: string;
  netName: string;
  subnetName: string;
  relatedVpc: RelatedVpc;
  supportClassicLink: boolean;
  osValue: string;
  osType: string;
  defaultUser: string;
}

const net = seajs.require('net');

async function fetchCVMList(data: any): Promise<CVM[]> {
  return new Promise<CVM[]>((resolve, reject) => {
    net.send(
      { method: 'GET', url: '/cgi/cvm?action=getCvmList' },
      {
        data,
        cb: response => {
          if (response.code === 0) {
            const list = response.data.deviceList as CVM[];
            list.forEach(x => {
              x.id = x.uuid;
            });
            resolve(list);
          } else {
            reject(new Error('拉取主机失败：' + JSON.stringify(response)));
          }
        }
      }
    );
  });
}

export async function loadCVMList(regionId: string | number, projectId: number) {
  return fetchCVMList({
    count: 999999,
    offset: 0,
    projectId,
    regionId
  });
}

export async function findCVMByLanIp(regionId: string | number, lanIp: string) {
  return fetchCVMList({
    count: 1,
    offset: 0,
    projectId: -1,
    regionId,
    vagueIp: lanIp
  });
}
