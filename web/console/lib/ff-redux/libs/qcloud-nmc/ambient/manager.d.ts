declare namespace nmc {
  // region-id to region-name
  interface RegionMap {
    [id: number]: string;
  }

  // whitelist map
  interface WhitelistMap {
    [key: string]: string[];
  }

  interface Project {
    projectId: number;
    ownerUin: number;
    appid: number;
    name: string;
    creatorUin: number;
    srcPlat: string;
    srcAppid: number;
    status: number;
    createTime: string;
    isDefault: number;
    info: string;
  }

  interface PermProject {
    projectId: number;
    ownerUin: number;
    appid: number;
    name: string;
    creatorUin: number;
    srcPlat: string;
    srcAppid: number;
    status: number;
    createTime: string;
    isDefault: number;
    info: string;
  }

  interface UserInfo {
    nick: string;
    uin: string;
    ownerUin: number;
    skey: string;
    isOwner: boolean;
    permProjects: PermProject[];
    isProjectUser: boolean;
    isShowAllProject: boolean;
    globalManageResourcePerm: boolean;
    globalManageFinancePerm: boolean;
    globalManageHumanResourcePerm: boolean;
    isGlobalUser: boolean;
    isAgent: boolean;
    isAgentClient: boolean;
    isSpread: boolean;
    canManageCloudResource: boolean;
  }

  interface CurInfo {
    name: string;
    type: number;
    area: number;
    id_card_type: number;
    id_card: string;
    organization_code: string;
    authenticateType: number;
  }

  interface OwnerInfo {
    user_type: string;
    uin: string;
    owner_uin: number;
    type: number;
    area: number;
    name: string;
    id_card: string;
    id_card_type: number;
    id_card_url: any;
    id_card_url_new: string;
    id_card_cache_url: any;
    id_card_cache_url_new: string;
    entry_card_cache_url_new: string;
    entry_card_url_new: string;
    contact: string;
    tel: string;
    organization_code: string;
    cur_organization_code: string;
    authenticateType: number;
    cur_authenticateType: number;
    mail: string;
    cur_mail: string;
    mail_pass: number;
    cur_mail_pass: number;
    register_status: string;
    devcheck_pass: number;
    cur_devcheck_state: number;
    bank_account: any;
    bank_number: any;
    perfection: number;
    wanIpTime: string;
    wanRestrict: number;
    accredit: number;
    devcheck_msg: string;
    cur_devcheck_msg: string;
    isModify: number;
    cur_info: CurInfo;
    bizInfo: string;
    addr: string;
    auth_method: number;
  }

  interface CommonData {
    projects: Project[];
    isShowAllProject: boolean;
    canManageCloudResource: boolean;
    globalManageFinancePerm: boolean;
    globalManageResourcePerm: boolean;
    globalManageHumanResourcePerm: boolean;
    isOwner: boolean;
    isProjectUser: boolean;
    isAgent: boolean;
    isAgentClient: boolean;
    isSpread: boolean;
    ownerUin: number;
    userInfo: UserInfo;
    ownerInfo: OwnerInfo;
    appId: number;
  }

  interface Manager {
    queryRegion(
      callback: (region: RegionMap) => any,
      fail?: (...args: any[]) => any
    );
    queryWhiteList(
      options: { whiteKey: string[] },
      callback: (result: WhitelistMap) => any
    );
    getProjects(callback: Function): void;
    getAllWhiteList(callback: (all: WhitelistMap) => any);
    getComData(callback: (data: CommonData) => any);
  }
}
