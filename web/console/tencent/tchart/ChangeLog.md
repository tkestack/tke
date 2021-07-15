# V1.1.0
### Feature
* 支持groupBy，conditions，period，aggregation查询，选中数据通过表格显示。需求 [TKE Mesh](http://tapd.oa.com/cloudMonitor/prong/stories/view/1020393412064629839)

# v1.1.1
### Chore
* 优化当groupBy为空时，table source 显示 --
* 监控项，单位名称前字段为黑色，单位为灰色

#V1.1.2
### Feature
* ChartFilterPanelState 支持 Resize 功能

### Chore
* 图表 title 文本可复制

# V1.1.3
### Feature
* 增加新的 PureChart 支持，有 Resize能力  需求 [监控图表单实例](http://tapd.oa.com/cloudMonitor/prong/stories/view/1020393412064769375)

### Chore
* 统一 panel 继承关系，base.tsx 为基础 panel，提供统一访问数据接口。
* 修改组件传入数据结构，统一使用tables数组进行请求参数配置

# V1.1.4
### Chore
* 优化chart的 Resize，当宽度小于设定值是，高度做等比例缩小

### Fix
* ChartFilterPanel 中conditions取值来自于传入的 groupBy

# V1.1.5
### Fix
* Chart 重新加载错误
* 单位为 undefined 时 tooltip 显示错误

### Chore
* 重新加载使用div展示
* ChartPanel 和 ChartFilterPanel 支持传递 width 和 height 调整大小

# V1.1.6
### Fix
* ChartInstancesPanel 重新加载是被置空
* ChartPanel 的 aggregateDataByGroupby 方法对 columns , data 两个参数做预处理
* 优化 label 隔天time的格式化

# V1.1.7
### Fix
* 优化 tooltip tr 之间没有间隔
* ChartFilterPanel中图表没有reload功能

# V1.1.8
### Fix
* 修复 一个table数据为空时导致所有的table都为空
* 优化 tooltip 无数据显示

# V1.2.1
### Feature
* ChartFilterPanel 支持多field 查询， tab切换与下拉懒加载
* request 支持解析 influxDB
* ChartInstancesPanel和ChartPanel支持下拉懒加载

### Fix
* ChartInstancesPanel 修复当instance中columns对多groupBy是显示折线错误

# V1.2.2
### Chore
* ChartPanel 添加统计粒度查询功能

# V1.2.3
### Fix
* FindMetricIndex修复column index查找bug
* ChartPanel 将统计粒度select格式设置为native

# V1.2.7
### Chore
* core.ts 调整labels补帧条件

# V1.3.0
### Feature
* 添加连续时序图表

# V1.3.4
### Feature
* 支持国际化

# V1.3.5
### Fix
* 修复单点无绘画

# V1.3.6
### Chore
* 适配后端无数据返回情况

# V1.3.9
### Fix
* 图表无数据是labels为空bug
### Chore
* 融合版tke获取动态监控地址
* 兼容influxdb返回series没有拼接聚合方式

# V1.3.10
### Fix
* 图表labels计算死循环bug

# V1.3.15
### Feature
支持Y轴刻度自定义
### Fix
* 鼠标hover后xAxisTickMarkIndex下标越界
* ChartInstancesPanel 图表更新时显示有未更新图表

# V1.3.17
### Fix
chart yAxis 没有复制传递的数值，导致修改外部数据

# V1.3.18
### Fix
yAxis 1024格式刻度时近似错误

# V1.3.19
### Fix
连续两个点无画线

### Chore
画点线操作平频繁，使用requestAnimationFrame提升性能
固定住的tooltip，鼠标移动到其它图表，不清除
日期选择限制范围30天
时间和粒度的关系保持和公有云一致

# V1.3.20
### Chore
Y轴和tooltip可以通过tooltipLabels控制显示

# V1.3.22
### Fix
PureChart 高度显示异常

# V1.3.24
### Chore
统一 tooltipLabels 处理Y轴和tooltip数字显示
修改GeneratePeriodOptions粒度单位显示错误

# V1.3.25
### Fix
修复设置Y轴的数值，没有初始化相应刻度数
修复图表列预占位后滚动函数计算错误

# V1.3.27
### Fix
instance chart 闪屏
instance chart 勾选新逻辑，以传入的instance list 的checked状态为主，当instance list 为空数组是，折线为显示状态，即disabled为false

# V1.3.28
### Chore
drawSpinner 设置安全机制，方式图表在loading状态被关闭后requestAnimationFrame还在执行
优化 MetricCharts 图表加载性能

# V1.4.0
### Feature
优化 paint 的核心方法 drawPolyLine，使用 Array 替代 Object 数据结构。
新增图片折线颜色配置参数：colorTheme 传入 ColorTypes 类型，colors 传入颜色字符串列表，优先使用colors。


