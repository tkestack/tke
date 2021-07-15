import Axios from "axios";
declare var devmode;
declare var seajs;
const iaas = seajs.require("models/iaas");

/**
 * CGI上报的code的映射
 */
const CGICodeEnum = {
  /** 参数错误 */
  InvalidParameter: 600,
  /** 参数取值错误 */
  InvalidParameterValue: 601,
  /** 缺少参数错误 */
  MissingParameter: 602,
  /** 未知参数错误 */
  UnknownParameter: 603,
  /** 接口不存在 */
  InvalidAction: 604,
  /** 请求次数超过了频率限制 */
  RequestLimitExceeded: 605,
  /** 接口版本不存在 */
  NoSuchVersion: 606,
  /** 接口不支持所传地域 */
  UnsupportedRegion: 607,
  /** 资源不存在 */
  ResourceNotFound: 608,
  /** 超过配额限制 */
  LimitExceeded: 609,
  /** 资源不可用 */
  ResourceUnavailable: 610,
  /** 资源不足 */
  ResourceInsufficient: 611,
  /** 操作失败 */
  FailedOperation: 612,
  /** 资源被占用 */
  ResourceInUse: 613,
  /** CAM签名/鉴权错误 */
  AuthFailure: 700,
  /** 未授权操作 */
  UnauthorizedOperation: 701,
  /** 操作不支持 */
  UnsupportedOperation: 702,
  /** 内部错误 */
  InternalError: 800,
  /** 内部错误超时 */
  InternalTimeout: 801
};

/**
 * 云api的相关操作
 */
export async function sendCapiRequest(
  serviceType: string,
  apiName: string,
  data?: any,
  regionId?: number,
  userOpts?: any
) {
  let params = {
    serviceType: serviceType,
    action: apiName,
    regionId: regionId,
    data: data
  };
  let opts = userOpts || {};
  opts = Object.assign({ camTips: false }, opts);
  if (opts && opts.clientTimeout) {
    params["clientTimeout"] = opts.clientTimeout;
  }

  // 统计开始时间
  let startTime = new Date();
  try {
    let result = await iaas.apiRequest(params, opts);

    // 统计成功的时间，进行数据的上报
    let duration = new Date().getTime() - startTime.getTime();
    sendCGITime({
      interfaceName: apiName,
      duration,
      isSuccess: true,
      moduleName: serviceType,
      code: 0,
      cgiParams: data,
      originCode: result.code
    });
    return result;
  } catch (error) {
    // 统计失败的时间，进行数据的上报
    let duration = new Date().getTime() - startTime.getTime();
    let code = 0;
    // 如果是内部错误，并且超时，则code 为 801，否则为其他
    if (duration > 5000 && CGICodeEnum[error.code] === CGICodeEnum.InternalError) {
      code = CGICodeEnum.InternalTimeout;
    } else {
      code = +error.code ? +error.code : CGICodeEnum[error.code] ? CGICodeEnum[error.code] : 4000;
    }
    sendCGITime({
      interfaceName: apiName,
      duration,
      isSuccess: false,
      moduleName: serviceType,
      code,
      originCode: error.code,
      cgiParams: data
    });
    throw error;
  }
}

/** huatuo上报的一些基本属性值 */
const BASIC_URL = "//report.huatuo.qq.com/code.cgi?appid=20488&platform=pc&domain=console.cloud.tencent.com&apn=unknow";

/** hubble上报的一些基本属性值 */
// const HUBBLE_BASIC_URL = '//report.hubble.qq.com/nv.cgi?__logname=data_0_278&key=cgi,type,code,time';
const HUBBLE_BASIC_URL = "__logname=data_0_277&key=cgi,type,code,time,module";
const HUBBLE_PV_BASIC_URL = "__logname=data_0_284&key=path,time,module";

let report_info = { url: "", count: 0 }; // 上报的数据存储，url 为当前已拼接的内容，count为已经拼接的数量
let report_hubble = { url: "", count: 0 };
let report_hubble_pv = { url: "", count: 0 };

let LastReportTimeOut: any = null;

/** 批量上报的格式内容 */
interface ReportInfo {
  cgi: string;
  code: number;
  type: number;
  time: number;
  /** 当前上报的URL */
  url: string;
  /** 当前已经上报的数量 */
  count: number;
}
/** 将接口、页面进行批量上报，使用http/1.1协议，故需要进行批量上报 */
const sendBatchInterfaceTimeToHuatuo = () => {
  // 如果在开发模式下，不进行数据的上报
  if (devmode.getStatus() === "on") {
    report_info = { url: "", count: 0 };
    report_hubble = { url: "", count: 0 };
    report_hubble_pv = { url: "", count: 0 };
  } else {
    let { url } = report_info as ReportInfo;
    http_img_sender()(`${BASIC_URL}&key=cgi,type,code,time${url}`);
    report_info = { url: "", count: 0 };

    let { url: hubbleUrl } = report_hubble as ReportInfo;
    let hubbleReportUrl = "//report.hubble.qq.com/nv.cgi";

    try {
      if (hubbleUrl.length) {
        const blob: Blob = new Blob([`${HUBBLE_BASIC_URL}${hubbleUrl}`], { type: "application/x-www-form-urlencoded" });
        // 需要判断是否支持sendBeacon方法
        if (navigator && navigator.sendBeacon) {
          navigator.sendBeacon(hubbleReportUrl, blob);
        } else {
          Axios({
            method: "post",
            url: hubbleReportUrl,
            data: blob
          });
        }
        report_hubble = { url: "", count: 0 };
      }

      // 上报pv
      let { url: pvUrl } = report_hubble_pv as ReportInfo;
      if (pvUrl.length) {
        const pvBlob: Blob = new Blob([`${HUBBLE_PV_BASIC_URL}${pvUrl}`], {
          type: "application/x-www-form-urlencoded"
        });
        if (navigator && navigator.sendBeacon) {
          navigator.sendBeacon(hubbleReportUrl, pvBlob);
        } else {
          Axios({
            method: "post",
            url: hubbleReportUrl,
            data: pvBlob
          });
        }
        report_hubble_pv = { url: "", count: 0 };
      }
    } catch (error) {
      // do nothing，不要影响业务逻辑
    }
  }
};

/** 拼接huatuo上报url的方法
 * @param cgi: 上报的cgi接口
 * @param code: 上报的code
 * @param type: 接口请求是否成功 1: 成功, 2: 失败
 * @param time: number  接口耗时
 * @param isHubble: boolean 是否为上报到Hubble上
 */
const buildReportUrl = (options: {
  cgi: string;
  code: number | string;
  type: number;
  time: number;
  isHubble: boolean;
  moduleName: string;
  isPageView: boolean;
}) => {
  let { cgi, code, type, time, isHubble = false, moduleName, isPageView } = options;
  let currentIndex = isPageView ? ++report_hubble_pv.count : isHubble ? ++report_hubble.count : ++report_info.count;

  if (isPageView) {
    let segment = `&${currentIndex}_1=${cgi}&${currentIndex}_2=${time}&${currentIndex}_3=${moduleName}`;
    // 保存在Hubble PV当中的信息
    report_hubble_pv.url += segment;
  } else {
    let segment: string = `&${currentIndex}_1=${cgi}&${currentIndex}_2=${type}&${currentIndex}_3=${code}&${currentIndex}_4=${time}`;
    if (isHubble) {
      let hubbleSegment = `${segment}&${currentIndex}_5=${moduleName}`;
      // 保存在Hubble当中的信息
      report_hubble.url += hubbleSegment;
    } else {
      report_info.url += segment;
    }
  }
};

/**
 * 统一创建上报的方法
 * @param isReportImmediate: boolean  是否立即上报
 */
const createReporter = (isReportImmediate: boolean = false) => {
  if (isReportImmediate) {
    // 如果需要立即上报，则清除之前的timer，并且先进行上报
    clearTimeout(LastReportTimeOut);
    LastReportTimeOut = null;

    // 上报到Huatuo
    sendBatchInterfaceTimeToHuatuo();
  } else {
    if (!LastReportTimeOut) {
      LastReportTimeOut = setTimeout(() => {
        sendBatchInterfaceTimeToHuatuo();
        LastReportTimeOut = null;
      }, 8000);
    }
  }
};

/**
 * CGI上报进行批量上报，将数据存储在localStorage当中
 * @param prefixName: string  前缀名称
 * @param interfaceName: string 接口名称
 * @param duration: number 上报的时间
 * @param isSuccess: boolean  是否成功
 * @param code: number  失败的话，需要传入code，默认是0
 */
const storeReportInfoInLocal = (options: {
  interfaceName: string;
  isHubble?: boolean;
  prefixName?: string;
  duration?: number;
  isSuccess?: boolean;
  code?: number | string;
  moduleName?: string;
  isPageView?: boolean;
}) => {
  let {
    prefixName = "",
    interfaceName,
    duration = 1,
    isSuccess = true,
    code = 0,
    isHubble = false,
    moduleName = "",
    isPageView = false
  } = options;

  // 存储此次的上报信息
  buildReportUrl({
    cgi: encodeURIComponent(prefixName + interfaceName),
    code,
    type: isSuccess ? 1 : 2,
    time: duration,
    isHubble,
    moduleName: encodeURIComponent(moduleName),
    isPageView
  });
  let { url } = report_info as ReportInfo;
  createReporter(url.length > 1700);
};

interface YunCGIParams {
  /** 请求头 */
  Accept?: string;
  /** 请求路径 */
  Path?: string;
}

/**
 * 上报接口的测速时间
 * @param interfaceName: string 接口名称
 * @param duration: number  接口测速时间
 * @param isSuccess: boolean  是否成功调用
 * @param moduleName: string  当前接口的模块名
 * @param code: number  当前返回的code
 * @param originCode: string  当前返回的实际的code
 */
function sendCGITime(options: {
  interfaceName: string;
  duration: number;
  isSuccess: boolean;
  moduleName: string;
  code: number;
  cgiParams: YunCGIParams;
  originCode: string;
}) {
  let { interfaceName, duration, isSuccess, moduleName, originCode, code = 0 } = options;

  // 存储到 localStorage当中
  storeReportInfoInLocal({
    prefixName: `${moduleName}/`,
    interfaceName: interfaceName,
    duration,
    isSuccess,
    code
  });

  // 上报到Hubble
  storeReportInfoInLocal({
    prefixName: `${moduleName}/`,
    interfaceName: interfaceName,
    duration,
    isSuccess,
    code: originCode,
    isHubble: true,
    moduleName
  });
}

/**
 * 用图片进行数据的上报
 */
function http_img_sender() {
  let img = new Image();
  let sender = function(src) {
    // 完成、错误或者中断之后，解除绑定，消灭对象
    img.onload = img.onerror = img.onabort = function() {
      img.onload = img.onerror = img.onabort = null;
      img = null;
    };
    img.src = src;
  };
  return sender;
}
