import Page from './page';

export default {
  register({businessKey, routes}) {
    for (const moduleKey in routes) {
      if (routes.hasOwnProperty(moduleKey)) {
        const moduleConfig = routes[moduleKey];
        const modulePath = moduleConfig.path || `modules/${businessKey}/${moduleKey}/${moduleKey}`;

        // window.define 使用的 sea.js 的 define 挂载到 window 对象上
        (window as any).define(modulePath, (require) => {
          return new Page({
            businessKey,
            id: moduleKey,
            title: moduleConfig.title,
            require,
            ...moduleConfig,
            requires: moduleConfig.requires
          });
        });
      }
    }
  }
}