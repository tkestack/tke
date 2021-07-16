"use strict";
var __values = (this && this.__values) || function (o) {
    var m = typeof Symbol === "function" && o[Symbol.iterator], i = 0;
    if (m) return m.call(o);
    return {
        next: function () {
            if (o && i >= o.length) o = void 0;
            return { value: o && o[i++], done: !o };
        }
    };
};
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
Object.defineProperty(exports, "__esModule", { value: true });
var warning = require("warning");
var entry_1 = require("./entry");
var config_1 = require("./config");
var i18n_1 = require("./i18n");
/**
 * 注册应用路由
 *
 * @param routemap 路由映射，每个路由对应一个路由定义
 *
 * @example
```js
app.routes({
  'cvm': {
    render: () => <CvmIndex />,
  },
  'cvm/overview': {
    render: () => <CvmOverview />,
  }
});
```
 *
 * @todo 支持 React-Hot
 */
function routes(routemap) {
    var e_1, _a, e_2, _b;
    var args = [];
    for (var _i = 1; _i < arguments.length; _i++) {
        args[_i - 1] = arguments[_i];
    }
    warning(args.length === 0, "导航菜单及样式配置现已移至 yehe，请移除 `app.routes()` 中 `menu` 及 `css` 配置");
    // i18n.init() 要求在 app.routes() 之前执行
    if (!i18n_1.getI18NInstance()) {
        throw new Error("未初始化国际化资源，请在 `app.route()` 之前执行 `i18n.init()`");
    }
    var controllerSet = new Set();
    try {
        for (var _c = __values(Object.entries(routemap)), _d = _c.next(); !_d.done; _d = _c.next()) {
            var _e = __read(_d.value, 2), route = _e[0], entry = _e[1];
            // 兼容开头包含 / 情况
            route = route.replace(/^\//, "");
            addEntry(route, entry);
        }
    }
    catch (e_1_1) { e_1 = { error: e_1_1 }; }
    finally {
        try {
            if (_d && !_d.done && (_a = _c.return)) _a.call(_c);
        }
        finally { if (e_1) throw e_1.error; }
    }
    try {
        // 提供默认 404 页面
        for (var controllerSet_1 = __values(controllerSet), controllerSet_1_1 = controllerSet_1.next(); !controllerSet_1_1.done; controllerSet_1_1 = controllerSet_1.next()) {
            var controller = controllerSet_1_1.value;
            var route404 = controller + "/404";
            if (!routemap[route404]) {
                addEntry(route404, 404);
            }
        }
    }
    catch (e_2_1) { e_2 = { error: e_2_1 }; }
    finally {
        try {
            if (controllerSet_1_1 && !controllerSet_1_1.done && (_b = controllerSet_1.return)) _b.call(controllerSet_1);
        }
        finally { if (e_2) throw e_2.error; }
    }
    function addEntry(route, entry) {
        if (!isLowerCase(route)) {
            throw new Error("\u8DEF\u7531 " + route + " \u975E\u6CD5\uFF1A\u8DEF\u7531\u53EA\u5141\u8BB8\u4F7F\u7528\u5C0F\u5199\u5B57\u6BCD\uFF01");
        }
        if (!entry) {
            throw new Error("\u8DEF\u7531 " + route + " \u975E\u6CD5\uFF1A\u8DEF\u7531\u5B9E\u4F8B\u4E3A\u7A7A\uFF01");
        }
        if (typeof entry !== "function" &&
            typeof entry !== "object" &&
            entry !== 404) {
            throw new Error("\u8DEF\u7531 " + route + " \u975E\u6CD5\uFF1A\u8DEF\u7531\u5B9E\u4F8B\u53EA\u80FD\u4E3A Function\u3001Object \u6216\u4F7F\u7528 404 \u8868\u793A\u9875\u9762\u4E0D\u5B58\u5728\uFF0C\u4F46\u662F\u4F20\u4E86 " + typeof entry + "\uFF01");
        }
        var _a = __read(route.split("/")), controller = _a[0], action = _a[1], rest = _a.slice(2);
        if (rest.length >= 1) {
            throw new Error("\u8DEF\u7531 " + route + " \u975E\u6CD5\uFF1A\u53EA\u5141\u8BB8\u6CE8\u518C\u4E00\u7EA7\u548C\u4E8C\u7EA7\u8DEF\u7531\n    \u66F4\u591A\u7EA7\u9700\u6C42\u8BF7\u4F7F\u7528\u6A21\u5757\u5185\u8DEF\u7531\uFF1Ahttp://tapd.oa.com/tcp_access/markdown_wikis/view/#1020399462008691883");
        }
        controllerSet.add(controller);
        // 没提供 action 的，使用 index
        action = action || "index";
        // 通过控制台的模块约定方式来注册模块
        window.define("modules/" + controller + "/" + action + "/" + action, config_1.config(entry_1.factory)(controller, action, entry));
    }
}
exports.routes = routes;
function isLowerCase(str) {
    return str === str.toLowerCase();
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoicm91dGVzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vc3JjL2NvcmUvcm91dGVzLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQSxpQ0FBbUM7QUFDbkMsaUNBQTRDO0FBQzVDLG1DQUFrQztBQUNsQywrQkFBeUM7QUFFekM7Ozs7Ozs7Ozs7Ozs7Ozs7OztHQWtCRztBQUNILFNBQWdCLE1BQU0sQ0FBQyxRQUFxQjs7SUFBRSxjQUFPO1NBQVAsVUFBTyxFQUFQLHFCQUFPLEVBQVAsSUFBTztRQUFQLDZCQUFPOztJQUNuRCxPQUFPLENBQ0wsSUFBSSxDQUFDLE1BQU0sS0FBSyxDQUFDLEVBQ2pCLDJEQUEyRCxDQUM1RCxDQUFDO0lBRUYsb0NBQW9DO0lBQ3BDLElBQUksQ0FBQyxzQkFBZSxFQUFFLEVBQUU7UUFDdEIsTUFBTSxJQUFJLEtBQUssQ0FDYiwrQ0FBK0MsQ0FDaEQsQ0FBQztLQUNIO0lBRUQsSUFBTSxhQUFhLEdBQUcsSUFBSSxHQUFHLEVBQVUsQ0FBQzs7UUFFeEMsS0FBMkIsSUFBQSxLQUFBLFNBQUEsTUFBTSxDQUFDLE9BQU8sQ0FBQyxRQUFRLENBQUMsQ0FBQSxnQkFBQSw0QkFBRTtZQUE1QyxJQUFBLHdCQUFjLEVBQWIsYUFBSyxFQUFFLGFBQUs7WUFDcEIsY0FBYztZQUNkLEtBQUssR0FBRyxLQUFLLENBQUMsT0FBTyxDQUFDLEtBQUssRUFBRSxFQUFFLENBQUMsQ0FBQztZQUNqQyxRQUFRLENBQUMsS0FBSyxFQUFFLEtBQUssQ0FBQyxDQUFDO1NBQ3hCOzs7Ozs7Ozs7O1FBRUQsY0FBYztRQUNkLEtBQXVCLElBQUEsa0JBQUEsU0FBQSxhQUFhLENBQUEsNENBQUEsdUVBQUU7WUFBakMsSUFBSSxVQUFVLDBCQUFBO1lBQ2pCLElBQU0sUUFBUSxHQUFNLFVBQVUsU0FBTSxDQUFDO1lBQ3JDLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLEVBQUU7Z0JBQ3ZCLFFBQVEsQ0FBQyxRQUFRLEVBQUUsR0FBRyxDQUFDLENBQUM7YUFDekI7U0FDRjs7Ozs7Ozs7O0lBRUQsU0FBUyxRQUFRLENBQUMsS0FBYSxFQUFFLEtBQWU7UUFDOUMsSUFBSSxDQUFDLFdBQVcsQ0FBQyxLQUFLLENBQUMsRUFBRTtZQUN2QixNQUFNLElBQUksS0FBSyxDQUFDLGtCQUFNLEtBQUssZ0dBQWtCLENBQUMsQ0FBQztTQUNoRDtRQUNELElBQUksQ0FBQyxLQUFLLEVBQUU7WUFDVixNQUFNLElBQUksS0FBSyxDQUFDLGtCQUFNLEtBQUssa0VBQWEsQ0FBQyxDQUFDO1NBQzNDO1FBQ0QsSUFDRSxPQUFPLEtBQUssS0FBSyxVQUFVO1lBQzNCLE9BQU8sS0FBSyxLQUFLLFFBQVE7WUFDekIsS0FBSyxLQUFLLEdBQUcsRUFDYjtZQUNBLE1BQU0sSUFBSSxLQUFLLENBQ2Isa0JBQU0sS0FBSywyTEFBb0QsT0FBTyxLQUFLLFdBQUcsQ0FDL0UsQ0FBQztTQUNIO1FBQ0csSUFBQSw2QkFBZ0QsRUFBL0Msa0JBQVUsRUFBRSxjQUFNLEVBQUUsa0JBQTJCLENBQUM7UUFDckQsSUFBSSxJQUFJLENBQUMsTUFBTSxJQUFJLENBQUMsRUFBRTtZQUNwQixNQUFNLElBQUksS0FBSyxDQUNiLGtCQUFNLEtBQUssZ1FBQTRHLENBQ3hILENBQUM7U0FDSDtRQUNELGFBQWEsQ0FBQyxHQUFHLENBQUMsVUFBVSxDQUFDLENBQUM7UUFDOUIsd0JBQXdCO1FBQ3hCLE1BQU0sR0FBRyxNQUFNLElBQUksT0FBTyxDQUFDO1FBQzNCLG9CQUFvQjtRQUNwQixNQUFNLENBQUMsTUFBTSxDQUNYLGFBQVcsVUFBVSxTQUFJLE1BQU0sU0FBSSxNQUFRLEVBQzNDLGVBQU0sQ0FBQyxlQUFPLENBQUMsQ0FBQyxVQUFVLEVBQUUsTUFBTSxFQUFFLEtBQUssQ0FBQyxDQUMzQyxDQUFDO0lBQ0osQ0FBQztBQUNILENBQUM7QUE1REQsd0JBNERDO0FBRUQsU0FBUyxXQUFXLENBQUMsR0FBVztJQUM5QixPQUFPLEdBQUcsS0FBSyxHQUFHLENBQUMsV0FBVyxFQUFFLENBQUM7QUFDbkMsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCAqIGFzIHdhcm5pbmcgZnJvbSBcIndhcm5pbmdcIjtcbmltcG9ydCB7IGZhY3RvcnksIEFwcEVudHJ5IH0gZnJvbSBcIi4vZW50cnlcIjtcbmltcG9ydCB7IGNvbmZpZyB9IGZyb20gXCIuL2NvbmZpZ1wiO1xuaW1wb3J0IHsgZ2V0STE4Tkluc3RhbmNlIH0gZnJvbSBcIi4vaTE4blwiO1xuXG4vKipcbiAqIOazqOWGjOW6lOeUqOi3r+eUsVxuICpcbiAqIEBwYXJhbSByb3V0ZW1hcCDot6/nlLHmmKDlsITvvIzmr4/kuKrot6/nlLHlr7nlupTkuIDkuKrot6/nlLHlrprkuYlcbiAqXG4gKiBAZXhhbXBsZVxuYGBganNcbmFwcC5yb3V0ZXMoe1xuICAnY3ZtJzoge1xuICAgIHJlbmRlcjogKCkgPT4gPEN2bUluZGV4IC8+LFxuICB9LFxuICAnY3ZtL292ZXJ2aWV3Jzoge1xuICAgIHJlbmRlcjogKCkgPT4gPEN2bU92ZXJ2aWV3IC8+LFxuICB9XG59KTtcbmBgYFxuICpcbiAqIEB0b2RvIOaUr+aMgSBSZWFjdC1Ib3RcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIHJvdXRlcyhyb3V0ZW1hcDogQXBwUm91dGVNYXAsIC4uLmFyZ3MpIHtcbiAgd2FybmluZyhcbiAgICBhcmdzLmxlbmd0aCA9PT0gMCxcbiAgICBcIuWvvOiIquiPnOWNleWPiuagt+W8j+mFjee9rueOsOW3suenu+iHsyB5ZWhl77yM6K+356e76ZmkIGBhcHAucm91dGVzKClgIOS4rSBgbWVudWAg5Y+KIGBjc3NgIOmFjee9rlwiXG4gICk7XG5cbiAgLy8gaTE4bi5pbml0KCkg6KaB5rGC5ZyoIGFwcC5yb3V0ZXMoKSDkuYvliY3miafooYxcbiAgaWYgKCFnZXRJMThOSW5zdGFuY2UoKSkge1xuICAgIHRocm93IG5ldyBFcnJvcihcbiAgICAgIFwi5pyq5Yid5aeL5YyW5Zu96ZmF5YyW6LWE5rqQ77yM6K+35ZyoIGBhcHAucm91dGUoKWAg5LmL5YmN5omn6KGMIGBpMThuLmluaXQoKWBcIlxuICAgICk7XG4gIH1cblxuICBjb25zdCBjb250cm9sbGVyU2V0ID0gbmV3IFNldDxzdHJpbmc+KCk7XG5cbiAgZm9yIChsZXQgW3JvdXRlLCBlbnRyeV0gb2YgT2JqZWN0LmVudHJpZXMocm91dGVtYXApKSB7XG4gICAgLy8g5YW85a655byA5aS05YyF5ZCrIC8g5oOF5Ya1XG4gICAgcm91dGUgPSByb3V0ZS5yZXBsYWNlKC9eXFwvLywgXCJcIik7XG4gICAgYWRkRW50cnkocm91dGUsIGVudHJ5KTtcbiAgfVxuXG4gIC8vIOaPkOS+m+m7mOiupCA0MDQg6aG16Z2iXG4gIGZvciAobGV0IGNvbnRyb2xsZXIgb2YgY29udHJvbGxlclNldCkge1xuICAgIGNvbnN0IHJvdXRlNDA0ID0gYCR7Y29udHJvbGxlcn0vNDA0YDtcbiAgICBpZiAoIXJvdXRlbWFwW3JvdXRlNDA0XSkge1xuICAgICAgYWRkRW50cnkocm91dGU0MDQsIDQwNCk7XG4gICAgfVxuICB9XG5cbiAgZnVuY3Rpb24gYWRkRW50cnkocm91dGU6IHN0cmluZywgZW50cnk6IEFwcEVudHJ5KSB7XG4gICAgaWYgKCFpc0xvd2VyQ2FzZShyb3V0ZSkpIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihg6Lev55SxICR7cm91dGV9IOmdnuazle+8mui3r+eUseWPquWFgeiuuOS9v+eUqOWwj+WGmeWtl+avje+8gWApO1xuICAgIH1cbiAgICBpZiAoIWVudHJ5KSB7XG4gICAgICB0aHJvdyBuZXcgRXJyb3IoYOi3r+eUsSAke3JvdXRlfSDpnZ7ms5XvvJrot6/nlLHlrp7kvovkuLrnqbrvvIFgKTtcbiAgICB9XG4gICAgaWYgKFxuICAgICAgdHlwZW9mIGVudHJ5ICE9PSBcImZ1bmN0aW9uXCIgJiZcbiAgICAgIHR5cGVvZiBlbnRyeSAhPT0gXCJvYmplY3RcIiAmJlxuICAgICAgZW50cnkgIT09IDQwNFxuICAgICkge1xuICAgICAgdGhyb3cgbmV3IEVycm9yKFxuICAgICAgICBg6Lev55SxICR7cm91dGV9IOmdnuazle+8mui3r+eUseWunuS+i+WPquiDveS4uiBGdW5jdGlvbuOAgU9iamVjdCDmiJbkvb/nlKggNDA0IOihqOekuumhtemdouS4jeWtmOWcqO+8jOS9huaYr+S8oOS6hiAke3R5cGVvZiBlbnRyeX3vvIFgXG4gICAgICApO1xuICAgIH1cbiAgICBsZXQgW2NvbnRyb2xsZXIsIGFjdGlvbiwgLi4ucmVzdF0gPSByb3V0ZS5zcGxpdChcIi9cIik7XG4gICAgaWYgKHJlc3QubGVuZ3RoID49IDEpIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihcbiAgICAgICAgYOi3r+eUsSAke3JvdXRlfSDpnZ7ms5XvvJrlj6rlhYHorrjms6jlhozkuIDnuqflkozkuoznuqfot6/nlLFcXG4gICAg5pu05aSa57qn6ZyA5rGC6K+35L2/55So5qih5Z2X5YaF6Lev55Sx77yaaHR0cDovL3RhcGQub2EuY29tL3RjcF9hY2Nlc3MvbWFya2Rvd25fd2lraXMvdmlldy8jMTAyMDM5OTQ2MjAwODY5MTg4M2BcbiAgICAgICk7XG4gICAgfVxuICAgIGNvbnRyb2xsZXJTZXQuYWRkKGNvbnRyb2xsZXIpO1xuICAgIC8vIOayoeaPkOS+myBhY3Rpb24g55qE77yM5L2/55SoIGluZGV4XG4gICAgYWN0aW9uID0gYWN0aW9uIHx8IFwiaW5kZXhcIjtcbiAgICAvLyDpgJrov4fmjqfliLblj7DnmoTmqKHlnZfnuqblrprmlrnlvI/mnaXms6jlhozmqKHlnZdcbiAgICB3aW5kb3cuZGVmaW5lKFxuICAgICAgYG1vZHVsZXMvJHtjb250cm9sbGVyfS8ke2FjdGlvbn0vJHthY3Rpb259YCxcbiAgICAgIGNvbmZpZyhmYWN0b3J5KShjb250cm9sbGVyLCBhY3Rpb24sIGVudHJ5KVxuICAgICk7XG4gIH1cbn1cblxuZnVuY3Rpb24gaXNMb3dlckNhc2Uoc3RyOiBzdHJpbmcpIHtcbiAgcmV0dXJuIHN0ciA9PT0gc3RyLnRvTG93ZXJDYXNlKCk7XG59XG5cbi8qKlxuICog6Lev55Sx5pig5bCEXG4gKi9cbmV4cG9ydCBpbnRlcmZhY2UgQXBwUm91dGVNYXAge1xuICBba2V5OiBzdHJpbmddOiBBcHBFbnRyeTtcbn1cbiJdfQ==