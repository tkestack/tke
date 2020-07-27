type RegionId = string | number;
type ProjectId = string | number;

interface AppUtil {
  getRegionId: () => RegionId;
  setRegionId: (regionId: RegionId) => void;
  getProjectId: (skipAll?: boolean) => ProjectId;
  setProjectId: (ProjectId: ProjectId) => void;
  submitForm: (url: string, data: any, method?: string) => void;
}

const appUtil: AppUtil = seajs.require('appUtil');
const util = seajs.require('util');

/**
 * 获取当前默认的地域ID（存于localStorage中，nmc全局属性）
 */
export function getRegionId(): number {
  return +appUtil.getRegionId();
}

/**
 * 设置为默认地域（存于localStorage中，nmc全局属性）
 */
export function setRegionId(regionId: number | string) {
  appUtil.setRegionId(regionId);
}

/**
 * 提交表单，跳转到订单支付页
 *      @param url: String      请求的URL
 *      @param data: Object     JSON数据
 *      @param method: String   请求的方式，默认为 post
 */
export function submitForm(url: string, data: any, method?: string) {
  appUtil.submitForm(url, data, method);
}

/**
 * 获取当前默认的集群ID
 */
export function getClusterId(): string {
  let dId = '',
    rId = util.cookie.get('clusterId');
  if (!rId) {
    if (window.localStorage) {
      let locRid = localStorage[util.getUin() + '_clusterId'];
      rId = locRid || dId;
    } else {
      rId = dId;
    }
  }
  return rId;
}

/**
 * 设置默认集群
 */
export function setProjectName(projectId: string) {
  util.cookie.set('projectId', projectId);
}

export function getProjectName(): string {
  let rId = util.cookie.get('projectId');
  return rId;
}

/**
 * 设置默认集群
 */
export function setClusterId(clusterId: string) {
  util.cookie.set('clusterId', clusterId);
  if (window.localStorage) {
    localStorage[util.getUin() + '_clusterId'] = clusterId;
  }
}
/**
 * 获取当前默认的集群命名空间
 */
export function getClusterNamespace(): string {
  let dId = '',
    rId = util.cookie.get('clusterNamespace');
  if (!rId) {
    if (window.localStorage) {
      let locRid = localStorage[util.getUin() + '_clusterNamespace'];
      rId = locRid || dId;
    } else {
      rId = dId;
    }
  }
  return rId;
}

/**
 * 设置默认集群命名空间
 */
export function setClusterNamespace(clusterId: string) {
  util.cookie.set('clusterNamespace', clusterId);
  if (window.localStorage) {
    localStorage[util.getUin() + '_clusterNamespace'] = clusterId;
  }
}

/**
 * 获取当前请求ID
 */
export function getRequestId(): number {
  let dId = '',
    rId = util.cookie.get('requestId');
  if (!rId) {
    if (window.localStorage) {
      let locRid = localStorage[util.getUin() + '_requestId'];
      rId = locRid || dId;
    } else {
      rId = dId;
    }
  }
  return rId;
}

/**
 * 设置当前请求ID
 */
export function setRequestId(requestId: number, clusterId: string) {
  util.cookie.set('requestId', requestId);
  if (window.localStorage) {
    localStorage[util.getUin() + '_' + clusterId + '_requestId'] = requestId;
  }
}

/**
 * 获取当前默认的地域列表
 */
export function getRegionList(): string[] {
  let list = util.cookie.get('regionList') || '';
  if (!list) {
    if (window.localStorage) {
      list = localStorage[util.getUin() + '_regionList'] || '';
    } else {
      list = '';
    }
  }
  return list ? list.split('|') : [];
}

/**
 * 设置地域列表
 */
export function setRegionList(list: string[]) {
  let listStr = list.join('|');
  util.cookie.set('regionList', listStr);
  if (window.localStorage) {
    localStorage[util.getUin() + '_regionList'] = listStr;
  }
}

export function debounce(fn, delay: number, immediate?: boolean) {
  let timer = null;
  return function () {
    const that = this;
    const args = arguments; // 箭头函数不能用arguments参数
    if (timer) {
      clearTimeout(timer);
    }
    if (immediate) {
      // 首次立即响应，后边过了delay时间之后才能响应
      if (timer === null) {
        fn.apply(that, args);
      }

      // 过了delay时间，设置判定条件timer = null; 上边才可以响应
      timer = setTimeout(function () {
        timer = null;
      }, delay);
    } else {
      // 每次触发都要delay时间之后才响应
      timer = setTimeout(function () {
        fn.apply(that, args);
      }, delay);
    }
  };
}
