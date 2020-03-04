/// <reference path="constants.d.ts" />
/// <reference path="seajs.d.ts" />
/// <reference path="manager.d.ts" />
/// <reference path="appUtil.d.ts" />
/// <reference path="qccomponent.d.ts" />
/// <reference path="net.d.ts" />
/// <reference path="router.d.ts" />
/// <reference path="eventtarget.d.ts" />
/// <reference path="tips.d.ts" />

declare namespace nmc {
  interface Require {
    (id: string): any;
    (id: "router"): nmc.Router;
    (id: "qccomponent"): nmc.Bee;
    (id: "tips"): nmc.Tips;
    (id: "net"): nmc.Net;
    (id: "$"): JQueryStatic;
    (id: "manager"): nmc.Manager;
    (id: "appUtil"): nmc.AppUtil;
    (id: "config/constants"): nmc.Constants;
    async(modules: string[], callback: Function);
  }
  export function render(html: string, moduleName: string): void;
}
