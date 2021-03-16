## TKEStack升级到1.19 api遗留问题

### 1. basic authentication被废弃问题

https://github.com/tkestack/tke/issues/1079

#### 1.1. 背景

1.19 api中k8s不再提供basicauthen/password的接口。

#### 1.2. 现状

registry目前从httprequest中获取basicauth信息，通过获取的信息实现了k8s basicauth的接口，进行认证。

#### 1.3. 解决方案

k8s中认证接口还保留httprequest方式，registry可将携带认证信息的httprequest作为入参数，实现k8s的httprequest接口，在接口实现中从request中获取basicauth认证信息进行认证。

### 2. ListWatchUntil被废弃问题

https://github.com/tkestack/tke/issues/1080

#### 2.1. 背景

https://github.com/kubernetes/kubernetes/pull/90855

#### 2.2. 现状

pkg/util/apiclient中的公共函数中用到了k8s watch包中的ListWatchUntil监听事件。

#### 2.3. 解决方案

根据社区的意见，https://github.com/kubernetes/kubernetes/pull/90855，直接使用watch包中的util代替，并将intialResourceVersion设置为默认值1。

### 3. TargetRAMMB参数从apiserver中移除问题

https://github.com/tkestack/tke/issues/1081

#### 3.1. 背景

1.19后apiserver启动参数中移除了TargetRAMMB。

#### 3.2. 现状

当前平台组件apiserver启动时会读取TargetRAMMB参数，以此参数作为etcd的缓存参数。

#### 3.3. 解决方案

将此参数从apiserver中移除，etcd的缓存参数使用DefaultWatchCacheSize。

### 4. 自定义资源storage实现需要table convertor问题

https://github.com/tkestack/tke/issues/1082

#### 4.1. 背景

1.19版本后storage不会再自动添加default table convertor，入参需要明确声明。

#### 4.2. 现状

目前平台storage实现都未使用特别声明table convertor。

#### 4.3. 解决方案

明确声明使用default table convertor，以兼容1.19的api。

### 5. internalversion被停用问题

https://github.com/tkestack/tke/issues/1083

#### 5.1. 背景

1.19后k8s不会再将internalversion自动转化为v1/beta1等实际版本api，需要明确访问的api版本。

#### 5.2. 现状

platform-api中通过proxy代理到目标集群资源，代理时其中listOption用到了internalversion导致与1.19不兼容。

#### 5.3. 解决方案

platform-api在将请求代理到目标集群api-server前将涉及到使用internalversion的listOption转换为v1的listOption。

### 6. TAPP/CronHPA在前端list时无法正常显示

https://github.com/tkestack/tke/issues/1086

#### 6.1. 背景

根据TAPP/CronHPA的报错，可参考社区类似的案例https://github.com/kubernetes/kubernetes/issues/94688。

#### 6.2. 现状

当平台组件api不升级到1.19时，TAPP/CronHPA前端list展示没有问题，但是在平台组件api升级到1.19后前端list展示将不能正常显示。

#### 6.3. 解决方案

根据维护TAPP/CronHPA到团队反馈，当前这两个组件尚未兼容1.19的api，需要这两个组件完成对1.19 api对兼容工作。
