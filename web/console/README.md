# tke-stack 控制台开发指南
tke-stack开发阶段调试主要方式是将浏览器环境中的js请求代理到本地来实现的
## 安装配置[whistle](https://wproxy.org/whistle/)
1. `npm i whistle -g`
2. 访问[http://127.0.0.1:8899/](http://127.0.0.1:8899/)
3. 安装证书，开启https代理，[详情](https://wproxy.org/whistle/webui/https.html)
4. 添加代理规则
  ```yaml
  ## tke-stack
  /\/static\/js\/project\.js$/	http://localhost:8090/js/app.js
  /\/static\/js\/index2\.js$/	http://localhost:8090/js/app.js
  /\/static\/js\/installer\.js$/	http://localhost:8088/js/app.js
  ```
## 安装chrome插件[SwitchyOmega](https://chrome.google.com/webstore/detail/proxy-switchyomega/padekgcemlokbadohgkifijomclgjgif)
1. 安装插件，需要翻墙，自行解决
2. 进入SwitchyOmega配置页面，新增情景模式`whistle`  ，配置如下：  
- 代理协议： http  
- 代理服务器： 127.0.0.1  
- 代理端口： 8899  
## 调试
1. 运行本地项目：
- 平台管理： `npm run dev`
- 业务管理： `npm run dev_project`
2. 进入已部署好的tke-stack控制台页面
3. `SwitchyOmega`选择情景模式为`whistle`
4. 刷新页面，检查项目运行前端页面是否已是本地开发修改过后的页面，简单的方法就是F12查看react-devtool中的componnets中的组件名称是否已经变为组件真实名称