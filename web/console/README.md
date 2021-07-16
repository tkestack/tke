# tke-stack 控制台开发指南

## 配置`server.config.js`

`package.json`同级目录新建`server.config.js` ，内容如下：

```js
module.exports = {
  Host: 'http://www.example.com',
  Cookie: ''
};
```

- 配置 Host 为已安装的控制台的访问地址;
- 配置 Cookie 为已安装的控制台的登录后的 cookie；

## 启动服务

- 平台管理为`npm run start-tke`;
- 业务管理为`npm run start-project`

## 快速替换运服务器端的静态文件

主要是为了快速查看前端控制台的改动，而不用把整个项目完整 build 部署一次

1. `npm run build`打包项目；
2. 将`build`目录 ftp 传输到 global 节点所在的机器上,并改名为`assets`目录
3. 以下命令均在 global 节点所在机器上执行；
4. 执行`kubectl get po -n tke |grep gateway`找到 gateway 的名称，自行替换下方命令 gateway 名称；
5. 执行 `kubectl cp assets tke-gateway-tjtc4://app/ -n tke` ，用新打包的文件替换掉老的文件；
