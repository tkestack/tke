@tencent/tea-app
================

`@tencent/tea-app` 为接入腾讯云控制台的业务提供下面两个能力：

- 向控制台注册业务模块入口
- 为业务模块提供控制台接口调用

## 开始使用

业务应优先考虑接入 tea 框架，获取完整的开发服务配套，点击[这里开始](http://tapd.oa.com/tcp_access/markdown_wikis/view/#1020399462008761129)。

如果业务出于**特殊情况**，需要单独使用，可独立安装此包：

```
npm i @tencent/tea-app --save
```

为了可以通过 `@tea/app` 导入，需要在编译的 Webpack 配置中指定 alias（使用 Tea 框架生成的，无需关注）：

```js
{
  alias: {
    "@tea/app": path.resolve(APP_ROOT, "./node_modules/@tencent/tea-app"
  }
}
```

## 模块入口注册

使用 `app.routes()` 注册指定路由的模块入口。

```js
import { app } from '@tea/app';

app.routes({
  'cvm': CvmIndexPage,
  'cvm/list': CvmListPage,
  'cvm/detail': CvmDetailPage,
});
```

入口模块最多支持两级路由，要支持模块内部路由，请参考[这篇文章](http://tapd.oa.com/tcp_access/markdown_wikis/view/#1020399462008691883)。

## 菜单配置

菜单配置请移至 http://yehe.isd.com/buffet 编辑。

## 样式配置

样式配置请移至 http://yehe.isd.com/buffet 编辑。

## 用户信息

### `app.user.curent()`

异步获取当前用户信息。

```js
const user = await app.user.current();
```

用户信息数据结构如下：

```ts
export interface AppUserData {
  /** 是否为主账号 */
  isOwner: boolean;

  /** 当前用户登录的 UIN */
  loginUin: number;

  /** 当前用户登录的主账号 UIN */
  ownerUin: number;

  /** 当前用户登录的主账号 APPID */
  appId: number;

  /** 当前用户的实名认证信息，如果未实名认证，此字段为 null */
  identity: AppUserIdentityInfo | null;

  /** 用户昵称 */
  nickName?: string;

  /** 用户标识名称，含开发商信息 */
  displayName?: string;
}
```

### `app.user.getAntiCSRFToken()`

获取当前登录用户的反 CSRF 凭据，用于网络请求中防御跨站脚本攻击。

### `app.user.checkWhitelist()`

异步检查用户是否在指定的白名单中。

**注意**：这个方法不会缓存白名单查询结果，如果不想重复查询，业务请自行缓存。

```js
// 如在白名单内，返回当前 ownerUin，否则返回 0
const ownerUin = await app.user.checkWhitelist('CLB_NEW_CONSOLE');
if (ownerUin) {
  // ...白名单操作
}
```

### `app.user.checkWhitelistBatch()`

批量检查用户是否在指定的白名单中。

**注意**：这个方法不会缓存白名单查询结果，如果不想重复查询，业务请自行缓存。

```js
const keys = ['CLB_NEW_CONSOLE', 'CLB_NEW_CONSOLE2'];
const result = await app.user.checkWhitelistBatch(keys);
const [k1, k2] = keys;
if (result[k1] || result[k2]) {
  // 命中任意白名单
}
```

### `app.user.getLastRegionId()`

获取用户最后使用的 regionId

- 如果不存在，则返回 -1
- 该数据基于用户当前的 ownerUin 存储在 localStorage 中
- 基于上一点，该数据可能会被其它业务修改，所以使用前应该先校验合法性

### `app.user.setLastRegionId()`

设置用户最后一次访问的 regionId，该数据可以被其它业务使用。

### `app.user.getLastProjectId()`

获取用户最后使用的 projectId

- 如果不存在，则返回 `-1`
- 该数据基于用户当前的 ownerUin 存储在 localStorage 中
- 基于上一点，该数据可能会被其它业务修改，所以使用前应该先校验合法性

### `app.user.setLastProjectId()`

设置用户最后一次访问的 projectId，该数据可以被其它业务使用。

### `app.user.getPermitedProjectInfo()`

异步获取用户有权限的项目信息。

```js
const { isShowAll, projects } = await app.user.getPermitedProjectInfo();
```

返回数据结构如下：

```ts
interface PermitedProjectInfo {
  /**
   * 当前用户是否有具备查看所有项目的权限
   */
  isShowAll: boolean;

  /**
   * 当前用户有权限的项目列表
   */
  projects: ProjectItem[];
}

interface ProjectItem {
  /**
   * 项目 ID，为 0 表示默认项目
   */
  projectId: number;

  /**
   * 项目名称
   */
  projectName: string;
}
```

### `app.user.getPermitedProjectList()`

同上，不过直接返回 `projects`：

```js
const projects = await app.user.getPermitedProjectList();

for (let { projectId, projectName } of projects) {
  console.log(projectId, projectName);
}
```

### `app.user.login()`

清除用户当前登录态，并弹出登录对话框。

### `app.user.logout()`

清除用户当前登录态。

### `app.user.on("invalidate", callback)`

用户登录态失效时触发：

```js
app.user.on("invalidate", (event) => {
  console.log(event);
});
```

其中 `event` 参数结构如下：

```ts
interface InvalidateEventArgs {
  source: "accountChanged" | "logout";
}
```

## 云 API 请求

> 请求 云 API v3 可使用 `app.capi.requestV3` 方法，该方法已自动填充 v3 相关参数

```js
const response = await app.capi.request({
  regionId: 1,
  serviceType: 'cvm',
  cmd: 'DescribeInstance',
  data: { Offset: 0, Limit: 20 }
});
```

方法签名：`app.capi.request(body, options): Promise<any>`

其中，`body` 和 `options` 参数结构如下：

```ts
interface RequestBody {
  /**
   * 请求的云 API 地域
   */
  regionId: number;

  /**
   * 请求的云 API 业务
   */
  serviceType: string;

  /**
   * 请求的云 API 名称
   */
  cmd: string;

  /**
   * 请求的云 API 数据
   */
  data?: any;
}

interface RequestOptions {

  /**
   * 是否使用安全的临时密钥 API 方案，建议使用 true
   * @default true
   */
  secure?: boolean;

  /**
   * 使用的云 API 版本，该参数配合 secure 使用
   *
   *   - `secure == false`，该参数无意义
   *   - `secure == true && version = 1`，将使用老的临时密钥服务进行密钥申请，否则使用新的密钥服务
   *   - `secure == true && version = 3`，使用云 API v3 域名请求，不同地域域名不同
   */
  version?: number;

  /**
   * 是否将客户端 IP 附加在云 API 的 `clientIP` 参数中
   */
  withClientIP?: boolean;

  /**
   * 是否将客户端 UA 附加在云 API 的 `clientUA` 参数中
   */
  withClientUA?: boolean;

  /**
   * 是否在顶部显示接口错误
   * 默认为 true，会提示云 API 调用错误信息，如果自己处理异常，请设置该配置为 false
   */
  tipErr?: boolean;
}
```

## MFA

### `app.mfa.verify()`

校验 MFA：调用敏感接口前，需要调用 MFA 让用户完成认证。

```js
// 发起 MFA 校验
const mfaPassed = await app.mfa.verify('cvm:DestroyInstance');

if (!mfaPassed) {
  // 校验取消，跳过后续业务
  return;
}

// 校验完成，调用云 API
const result = await app.capi.request({
  serviceType: 'cvm',
  cmd: 'DestroyInstance',
  // ...
});
```

### `app.mfa.request()`

校验 MFA 后调用云 API，使用 `app.mfa.request()` 具备失败重新校验的能力。

```js
const result = await app.mfa.request({
  regionId: 1,
  serviceType: 'cvm',
  cmd: 'DestroyInstance',
  data: {
    instanceId: 'ins-a5d3ccw8c'
  }
}, {
  onMFAError: error => {
    // 碰到 MFA 的错误请进行重试逻辑，业务可以自己限制重试次数
    return error.retry();
  }
})
```

> 注意：*如果已经使用了 `app.mfa.verify()` 方法进行 MFA 校验，则无需再使用该方法发起 API 请求，直接使用*

## 顶部提示

`tips` 下包含 `success`、`error` 和 `loading` 方法，每个方法都支持两种传参方式。

### 成功提示

#### API

`duration` 默认值 4000ms

```
app.tips.success(message, duration?);
```
```
app.tips.success({ message, duration? });
```

#### 示例

```js
app.tips.success('操作成功');
```

### 错误提示

#### API

`duration` 默认值 4000ms

```
app.tips.error(message, duration?);
```
```
app.tips.error({ message, duration? });
```

#### 示例

```js
app.tips.error('操作失败');
```

### loading 提示

`loading` 可使用其返回值中的 `stop` 方法手动停止提示。

**注意：当前 loading 提示持续不会超过 `5000ms`，超时后自动销毁。**

#### API

`message` 默认值 `"正在加载..."`（已处理国际化）

`duration` 默认值 `4000ms`

```
app.tips.loading(message?, duration?);
```
```
app.tips.loading({ message?, duration? });
```

#### 示例

```js
app.tips.loading();
```

```js
async function someAsyncTask() {
  const loadingTip = app.tips.loading({ duration: 5000 });

  try {
    await doAsyncTask();
  } catch (err) {
  } finally {
    loadingTip.stop();
  }
}
```

## 文档标题

### 配置模块文档标题

```js
app.routes({
  'cvm': {
    title: '云主机-首页',
    // component: CvmIndexPage,
    render: () => <CvmIndexPage />
  }
});
```

### 动态设置文档标题

提供 React Hooks `useDocumentTitle()` 让业务设置文档标题。

```js
function CvmDetail({ cvmName }) {
  useDocumentTitle(cvmName);
  return (
    <div>...</div>
  );
}
```

## 数据上报

由于上报数据占用 Insight 集群的存储/计算资源，当前采取下面的策略限制资源：

- 10 秒内调用次数不能超过 10 条，超过的部分丢弃并且控制台输出告警（最多输出一次）
- 10 秒内调用次数如果超过 100 条，则后续所有上报全部丢弃，预防不小心的高频上报

### 点击流上报

```js
app.insight.reportHotTag("TAG-NAME")
```

统计数据可在 http://yehe.isd.com/console-insight/hottag 查询

### 自定义上报

自定义上报限制：

- 总字段个数（`stringFields` + `integerFields`）不能超过 20 个
- 字符串字段的长度，不能超过 4096

```js
app.insight.stat({
  ns: 'cvm',
  event: 'restart',
  stringFields: {
    instance: 'ins-1d7x8C3s'
  },
  integerFields: {
    cost: 1000
  }
});
```

上报流水可以在 [Kibana](http://platform.cloud.oa.com/app/kibana#/discover/AWzdnPcWA6vLf9CW0Tg3) 查询，请自行修改 `d.ns` 和 `d.stat` 条件。

## SDK

2019 年 4 月份开始，控制台正式支持 SDK 的接入，详情请参考[此文档](http://tapd.oa.com/tcp_access/markdown_wikis/view/#1020399462008859945)。

### 注册 SDK

对于使用 `tea-cli` 生成的项目，下面的使用代码会自动生成：

```js
app.sdk.register('demo-sdk', () => {
  // SDK 导出的方法
  return {
    hello() {
      console.log('Hello from demo-sdk');
    }
  }
});
```

### 引用 SDK

要引用在控制台已经登记的 SDK，使用下列 API：

```js
try {
  const demoSDK = await app.sdk.use('demo-sdk');
  demoSDK.hello();
}
catch (err) {
  console.error('加载 SDK 失败：' + err.message);
}
```
