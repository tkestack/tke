# TKEStack Common

TKEStack Common 是前端控制台的基础库依赖库的相关配置，包含路由、业务配置、样式表等。

## 主要功能

提供基础的依赖库、静态资源、样式表等。
目前框架主要支持业务路由的设置以及业务的加载

## 如何配置新增业务模块

查看 `index.html` 文件，查找到如下的代码

```javascript
// 此处为配置js的路由
g_config_data.jsRouter_zh = {
  tke: {
    content: '/static/js/index2.js', // 表示业务加载的js的路径
    isUpdateLteIE9: '',
    component: 'tea2.0' // 表示业务所使用的组件，tea1.0、tea2.0，所加载的组件不同
  }
};

// 此处为配置css的相关配置文件
window.g_config_data.css_config = {
  tke: ['/static/css/business/tkestack/tkestack.css']
};
```

配置当中的 `keyName` 由业务注册的 `businessKey` 决定。
