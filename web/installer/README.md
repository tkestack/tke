# installer部分开发说明
## 起手
1. `whistle`添加代理规则`/\/static\/js\/installer\.js$/	http://localhost:8088/js/app.js`
2. `npm run common_pre`: 准备依赖
3. `npm run dev`: dev开发模式，前端经过`whistle`代理进行开发调试

## 组件开发调试
1. `npm run storybook`：自行参考`storybook`开发文档