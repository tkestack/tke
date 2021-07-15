"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
Object.defineProperty(exports, "__esModule", { value: true });
var React = require("react");
var ReactDOM = require("react-dom");
var logTrace_1 = require("../helpers/logTrace");
var insight_1 = require("./insight");
var i18n_1 = require("./i18n");
var create_history_1 = require("./history/create-history");
var history_context_1 = require("./history/history-context");
// 入口组件
function createEntry(controller, action, entry) {
    var i18n = i18n_1.getI18NInstance();
    var EntryContent = typeof entry === "function" ? entry : entry.component || entry.render;
    EntryContent["displayName"] = "Entry(" + controller + "/" + action + ")";
    var confirmation = null;
    var history = create_history_1.createHistory({
        controller: controller,
        action: action,
        getUserConfirmation: function (message, callback) {
            if (!confirmation) {
                callback(true);
            }
            else {
                confirmation(message, callback);
            }
        },
    });
    return function TeaEntry() {
        return (React.createElement(i18n_1.I18NProvider, { i18n: i18n },
            React.createElement(history_context_1.HistoryContext.Provider, { value: {
                    history: history,
                    setConfirmation: function (current) { return (confirmation = current); },
                } },
                React.createElement(EntryContent, null))));
    };
}
/**
 * 为模块生成渲染和销毁方法
 */
exports.factory = function (controller, action, entry) { return function () {
    // 模块需要实现渲染和销毁两个方法
    var render, destroy;
    // 渲染 DOM，render() 时创建，destory() 时回收
    var container;
    // 模块标识，同时用作 container 的 id
    var moduleKey = controller + "-" + action;
    // 渲染：路由命中时渲染
    render = function () {
        if (entry === 404) {
            return 404;
        }
        // window.nmc.render 负责把 DOM 创建并放在 #app-area 中，这个 DOM 后续被渲染所用
        window.nmc.render("<div id=\"" + moduleKey + "\" class=\"tea-page-root\"></div>", controller);
        container = document.querySelector("#" + moduleKey);
        if (container) {
            var Entry = createEntry(controller, action, entry);
            ReactDOM.render(React.createElement(Entry, null), container);
        }
    };
    // 模块：路由切走时销毁
    destroy = function () {
        if (container) {
            ReactDOM.unmountComponentAtNode(container);
        }
        container = null;
    };
    return __assign({}, insightFactory(moduleKey, render, destroy), { 
        // 导出文档标题（如果已定义）
        title: typeof entry === "object" ? entry.title : null });
}; };
function insightFactory(moduleKey, render, destroy) {
    var _a = insight_1.internalInsightShouldNotUsedByBusiness || {}, care = _a.care, register = _a.register, EventLevel = _a.EventLevel;
    if (typeof care === "function") {
        // 监控：模块渲染异常
        var moduleRenderError_1 = register("module-render-error", {
            level: EventLevel.Error,
        });
        // 监控：模块销毁异常
        var moduleDestroyError_1 = register("module-destroy-error", {
            level: EventLevel.Error,
        });
        // 捕获渲染异常
        render = care(render, {
            capture: function (trace) {
                var message = "\u6A21\u5757 " + moduleKey + " \u6E32\u67D3\u5F02\u5E38\uFF1A" + trace.message;
                logTrace_1.logTrace(message, trace);
                moduleRenderError_1.push({
                    message: message,
                    moduleKey: moduleKey,
                    stack: trace.stack,
                });
                throw trace;
            },
        });
        // 捕获销毁异常
        destroy = care(destroy, {
            capture: function (trace) {
                var message = "\u6A21\u5757 " + moduleKey + " \u9500\u6BC1\u5F02\u5E38\uFF1A" + trace.message;
                logTrace_1.logTrace(message, trace);
                moduleDestroyError_1.push({
                    message: message,
                    moduleKey: moduleKey,
                    stack: trace.stack,
                });
                throw trace;
            },
        });
    }
    return { render: render, destroy: destroy };
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZW50cnkuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvY29yZS9lbnRyeS50c3giXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7OztBQUFBLDZCQUErQjtBQUMvQixvQ0FBc0M7QUFFdEMsZ0RBQStDO0FBQy9DLHFDQUdtQjtBQUNuQiwrQkFBdUQ7QUFDdkQsMkRBQXlEO0FBQ3pELDZEQUF5RTtBQUV6RSxPQUFPO0FBQ1AsU0FBUyxXQUFXLENBQ2xCLFVBQWtCLEVBQ2xCLE1BQWMsRUFDZCxLQUE2QjtJQUU3QixJQUFNLElBQUksR0FBRyxzQkFBZSxFQUFFLENBQUM7SUFFL0IsSUFBTSxZQUFZLEdBQ2hCLE9BQU8sS0FBSyxLQUFLLFVBQVUsQ0FBQyxDQUFDLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsU0FBUyxJQUFJLEtBQUssQ0FBQyxNQUFNLENBQUM7SUFDeEUsWUFBWSxDQUFDLGFBQWEsQ0FBQyxHQUFHLFdBQVMsVUFBVSxTQUFJLE1BQU0sTUFBRyxDQUFDO0lBRS9ELElBQUksWUFBWSxHQUFpQixJQUFJLENBQUM7SUFDdEMsSUFBTSxPQUFPLEdBQUcsOEJBQWEsQ0FBQztRQUM1QixVQUFVLFlBQUE7UUFDVixNQUFNLFFBQUE7UUFDTixtQkFBbUIsRUFBRSxVQUFDLE9BQU8sRUFBRSxRQUFRO1lBQ3JDLElBQUksQ0FBQyxZQUFZLEVBQUU7Z0JBQ2pCLFFBQVEsQ0FBQyxJQUFJLENBQUMsQ0FBQzthQUNoQjtpQkFBTTtnQkFDTCxZQUFZLENBQUMsT0FBTyxFQUFFLFFBQVEsQ0FBQyxDQUFDO2FBQ2pDO1FBQ0gsQ0FBQztLQUNGLENBQUMsQ0FBQztJQUVILE9BQU8sU0FBUyxRQUFRO1FBQ3RCLE9BQU8sQ0FDTCxvQkFBQyxtQkFBWSxJQUFDLElBQUksRUFBRSxJQUFJO1lBQ3RCLG9CQUFDLGdDQUFjLENBQUMsUUFBUSxJQUN0QixLQUFLLEVBQUU7b0JBQ0wsT0FBTyxTQUFBO29CQUNQLGVBQWUs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QUFBLEtBQUs7Z0JBQ1osSUFBTSxPQUFPLEdBQUcsa0JBQU0sU0FBUyx1Q0FBUyxLQUFLLENBQUMsT0FBUyxDQUFDO2dCQUN4RCxtQkFBUSxDQUFDLE9BQU8sRUFBRSxLQUFLLENBQUMsQ0FBQztnQkFDekIsb0JBQWtCLENBQUMsSUFBSSxDQUFDO29CQUN0QixPQUFPLFNBQUE7b0JBQ1AsU0FBUyxXQUFBO29CQUNULEtBQUssRUFBRSxLQUFLLENBQUMsS0FBSztpQkFDbkIsQ0FBQyxDQUFDO2dCQUNILE1BQU0sS0FBSyxDQUFDO1lBQ2QsQ0FBQztTQUNGLENBQUMsQ0FBQztLQUNKO0lBQ0QsT0FBTyxFQUFFLE1BQU0sUUFBQSxFQUFFLE9BQU8sU0FBQSxFQUFFLENBQUM7QUFDN0IsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCAqIGFzIFJlYWN0IGZyb20gXCJyZWFjdFwiO1xuaW1wb3J0ICogYXMgUmVhY3RET00gZnJvbSBcInJlYWN0LWRvbVwiO1xuXG5pbXBvcnQgeyBsb2dUcmFjZSB9IGZyb20gXCIuLi9oZWxwZXJzL2xvZ1RyYWNlXCI7XG5pbXBvcnQge1xuICBpbnRlcm5hbEluc2lnaHRTaG91bGROb3RVc2VkQnlCdXNpbmVzcyBhcyBfaW5zaWdodCxcbiAgSW5zaWdodENvcmUsXG59IGZyb20gXCIuL2luc2lnaHRcIjtcbmltcG9ydCB7IEkxOE5Qcm92aWRlciwgZ2V0STE4Tkluc3RhbmNlIH0gZnJvbSBcIi4vaTE4blwiO1xuaW1wb3J0IHsgY3JlYXRlSGlzdG9yeSB9IGZyb20gXCIuL2hpc3RvcnkvY3JlYXRlLWhpc3RvcnlcIjtcbmltcG9ydCB7IENvbmZpcm1hdGlvbiwgSGlzdG9yeUNvbnRleHQgfSBmcm9tIFwiLi9oaXN0b3J5L2hpc3RvcnktY29udGV4dFwiO1xuXG4vLyDlhaXlj6Pnu4Tku7ZcbmZ1bmN0aW9uIGNyZWF0ZUVudHJ5KFxuICBjb250cm9sbGVyOiBzdHJpbmcsXG4gIGFjdGlvbjogc3RyaW5nLFxuICBlbnRyeTogRXhjbHVkZTxBcHBFbnRyeSwgNDA0PlxuKSB7XG4gIGNvbnN0IGkxOG4gPSBnZXRJMThOSW5zdGFuY2UoKTtcblxuICBjb25zdCBFbnRyeUNvbnRlbnQgPVxuICAgIHR5cGVvZiBlbnRyeSA9PT0gXCJmdW5jdGlvblwiID8gZW50cnkgOiBlbnRyeS5jb21wb25lbnQgfHwgZW50cnkucmVuZGVyO1xuICBFbnRyeUNvbnRlbnRbXCJkaXNwbGF5TmFtZVwiXSA9IGBFbnRyeSgke2NvbnRyb2xsZXJ9LyR7YWN0aW9ufSlgO1xuXG4gIGxldCBjb25maXJtYXRpb246IENvbmZpcm1hdGlvbiA9IG51bGw7XG4gIGNvbnN0IGhpc3RvcnkgPSBjcmVhdGVIaXN0b3J5KHtcbiAgICBjb250cm9sbGVyLFxuICAgIGFjdGlvbixcbiAgICBnZXRVc2VyQ29uZmlybWF0aW9uOiAobWVzc2FnZSwgY2FsbGJhY2spID0+IHtcbiAgICAgIGlmICghY29uZmlybWF0aW9uKSB7XG4gICAgICAgIGNhbGxiYWNrKHRydWUpO1xuICAgICAgfSBlbHNlIHtcbiAgICAgICAgY29uZmlybWF0aW9uKG1lc3NhZ2UsIGNhbGxiYWNrKTtcbiAgICAgIH1cbiAgICB9LFxuICB9KTtcblxuICByZXR1cm4gZnVuY3Rpb24gVGVhRW50cnkoKSB7XG4gICAgcmV0dXJuIChcbiAgICAgIDxJMThOUHJvdmlkZXIgaTE4bj17aTE4bn0+XG4gICAgICAgIDxIaXN0b3J5Q29udGV4dC5Qcm92aWRlclxuICAgICAgICAgIHZhbHVlPXt7XG4gICAgICAgICAgICBoaXN0b3J5LFxuICAgICAgICAgICAgc2V0Q29uZmlybWF0aW9uOiBjdXJyZW50ID0+IChjb25maXJtYXRpb24gPSBjdXJyZW50KSxcbiAgICAgICAgICB9fVxuICAgICAgICA+XG4gICAgICAgICAgPEVudHJ5Q29udGVudCAvPlxuICAgICAgICA8L0hpc3RvcnlDb250ZXh0LlByb3ZpZGVyPlxuICAgICAgPC9JMThOUHJvdmlkZXI+XG4gICAgKTtcbiAgfTtcbn1cblxuLyoqXG4gKiDkuLrmqKHlnZfnlJ/miJDmuLLmn5PlkozplIDmr4Hmlrnms5VcbiAqL1xuZXhwb3J0IGNvbnN0IGZhY3Rvcnk6IEVudHJ5TW9kdWxlRmFjdG9yeSA9IChcbiAgY29udHJvbGxlcixcbiAgYWN0aW9uLFxuICBlbnRyeVxuKSA9PiAoKSA9PiB7XG4gIC8vIOaooeWdl+mcgOimgeWunueOsOa4suafk+WSjOmUgOavgeS4pOS4quaWueazlVxuICBsZXQgcmVuZGVyOiAoKSA9PiB2b2lkLCBkZXN0cm95OiAoKSA9PiB2b2lkO1xuXG4gIC8vIOa4suafkyBET03vvIxyZW5kZXIoKSDml7bliJvlu7rvvIxkZXN0b3J5KCkg5pe25Zue5pS2XG4gIGxldCBjb250YWluZXI6IEhUTUxEaXZFbGVtZW50O1xuXG4gIC8vIOaooeWdl+agh+ivhu+8jOWQjOaXtueUqOS9nCBjb250YWluZXIg55qEIGlkXG4gIGNvbnN0IG1vZHVsZUtleSA9IGAke2NvbnRyb2xsZXJ9LSR7YWN0aW9ufWA7XG5cbiAgLy8g5riy5p+T77ya6Lev55Sx5ZG95Lit5pe25riy5p+TXG4gIHJlbmRlciA9ICgpID0+IHtcbiAgICBpZiAoZW50cnkgPT09IDQwNCkge1xuICAgICAgcmV0dXJuIDQwNDtcbiAgICB9XG4gICAgLy8gd2luZG93Lm5tYy5yZW5kZXIg6LSf6LSj5oqKIERPTSDliJvlu7rlubbmlL7lnKggI2FwcC1hcmVhIOS4re+8jOi/meS4qiBET00g5ZCO57ut6KKr5riy5p+T5omA55SoXG4gICAgd2luZG93Lm5tYy5yZW5kZXIoXG4gICAgICBgPGRpdiBpZD1cIiR7bW9kdWxlS2V5fVwiIGNsYXNzPVwidGVhLXBhZ2Utcm9vdFwiPjwvZGl2PmAsXG4gICAgICBjb250cm9sbGVyXG4gICAgKTtcbiAgICBjb250YWluZXIgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKGAjJHttb2R1bGVLZXl9YCk7XG4gICAgaWYgKGNvbnRhaW5lcikge1xuICAgICAgY29uc3QgRW50cnkgPSBjcmVhdGVFbnRyeShjb250cm9sbGVyLCBhY3Rpb24sIGVudHJ5KTtcbiAgICAgIFJlYWN0RE9NLnJlbmRlcig8RW50cnkgLz4sIGNvbnRhaW5lcik7XG4gICAgfVxuICB9O1xuXG4gIC8vIOaooeWdl++8mui3r+eUseWIh+i1sOaXtumUgOavgVxuICBkZXN0cm95ID0gKCkgPT4ge1xuICAgIGlmIChjb250YWluZXIpIHtcbiAgICAgIFJlYWN0RE9NLnVubW91bnRDb21wb25lbnRBdE5vZGUoY29udGFpbmVyKTtcbiAgICB9XG4gICAgY29udGFpbmVyID0gbnVsbDtcbiAgfTtcblxuICByZXR1cm4ge1xuICAgIC8vIOWMheijhSByZW5kZXIg5ZKMIGRlc3Ryb3kg5pa55rOV55qE5byC5bi45aSE55CGXG4gICAgLi4uaW5zaWdodEZhY3RvcnkobW9kdWxlS2V5LCByZW5kZXIsIGRlc3Ryb3kpLFxuICAgIC8vIOWvvOWHuuaWh+aho+agh+mimO+8iOWmguaenOW3suWumuS5ie+8iVxuICAgIHRpdGxlOiB0eXBlb2YgZW50cnkgPT09IFwib2JqZWN0XCIgPyBlbnRyeS50aXRsZSA6IG51bGwsXG4gIH0gYXMgRW50cnlNb2R1bGU7XG59O1xuXG5leHBvcnQgaW50ZXJmYWNlIEVudHJ5TW9kdWxlRmFjdG9yeSB7XG4gIChjb250cm9sbGVyOiBzdHJpbmcsIGFjdGlvbjogc3RyaW5nLCBlbnRyeTogQXBwRW50cnkpOiAoKSA9PiBFbnRyeU1vZHVsZTtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBFbnRyeU1vZHVsZSB7XG4gIHJlbmRlcigpOiB2b2lkO1xuICBkZXN0cm95KCk6IHZvaWQ7XG4gIHRpdGxlPzogc3RyaW5nO1xufVxuXG4vKipcbiAqIOWNleS4qui3r+eUseWFpeWPo+WumuS5ie+8jOWMheWQq+a4suafk+aWueazleOAgeiPnOWNleWSjCBDU1Mg5a6a5LmJXG4gKi9cbmV4cG9ydCB0eXBlIEFwcEVudHJ5ID0gUmVhY3QuQ29tcG9uZW50VHlwZTxhbnk+IHwgQXBwRW50cnlEZXRhaWwgfCA0MDQ7XG5cbmV4cG9ydCB0eXBlIEFwcEVudHJ5RGV0YWlsID0ge1xuICAvKipcbiAgICog6Lev55Sx5riy5p+T57uE5Lu2XG4gICAqL1xuICBjb21wb25lbnQ/OiBSZWFjdC5Db21wb25lbnRUeXBlPGFueT47XG5cbiAgLyoqXG4gICAqIOi3r+eUsea4suafk+aWueazle+8jOmcgOi/lOWbnuS4gOS4qiBSZWFjdC5SZWFjdE5vZGVcbiAgICovXG4gIHJlbmRlcj86ICgpID0+IEpTWC5FbGVtZW50O1xuXG4gIC8qKlxuICAgKiDor6Xot6/nlLHkuIvmiYDmnInnmoTmlofmoaPmoIfpophcbiAgICpcbiAgICog5aaC5p6c5biM5pyb5Yqo5oCB6K6+572u5qCH6aKY77yM5Y+v5Lul5L2/55SoIHVzZURvY3VtZW50VGl0bGUoKSAvIHNldERvY3VtZW50VGl0bGUoKVxuICAgKi9cbiAgdGl0bGU/OiBzdHJpbmc7XG59O1xuXG5mdW5jdGlvbiBpbnNpZ2h0RmFjdG9yeShcbiAgbW9kdWxlS2V5OiBzdHJpbmcsXG4gIHJlbmRlcjogKCkgPT4gdm9pZCxcbiAgZGVzdHJveTogKCkgPT4gdm9pZFxuKSB7XG4gIGNvbnN0IHsgY2FyZSwgcmVnaXN0ZXIsIEV2ZW50TGV2ZWwgfSA9IF9pbnNpZ2h0IHx8ICh7fSBhcyBJbnNpZ2h0Q29yZSk7XG4gIGlmICh0eXBlb2YgY2FyZSA9PT0gXCJmdW5jdGlvblwiKSB7XG4gICAgLy8g55uR5o6n77ya5qih5Z2X5riy5p+T5byC5bi4XG4gICAgY29uc3QgbW9kdWxlUmVuZGVyRXJyb3IgPSByZWdpc3RlcihcIm1vZHVsZS1yZW5kZXItZXJyb3JcIiwge1xuICAgICAgbGV2ZWw6IEV2ZW50TGV2ZWwuRXJyb3IsXG4gICAgfSk7XG4gICAgLy8g55uR5o6n77ya5qih5Z2X6ZSA5q+B5byC5bi4XG4gICAgY29uc3QgbW9kdWxlRGVzdHJveUVycm9yID0gcmVnaXN0ZXIoXCJtb2R1bGUtZGVzdHJveS1lcnJvclwiLCB7XG4gICAgICBsZXZlbDogRXZlbnRMZXZlbC5FcnJvcixcbiAgICB9KTtcbiAgICAvLyDmjZXojrfmuLLmn5PlvILluLhcbiAgICByZW5kZXIgPSBjYXJlKHJlbmRlciwge1xuICAgICAgY2FwdHVyZTogdHJhY2UgPT4ge1xuICAgICAgICBjb25zdCBtZXNzYWdlID0gYOaooeWdlyAke21vZHVsZUtleX0g5riy5p+T5byC5bi477yaJHt0cmFjZS5tZXNzYWdlfWA7XG4gICAgICAgIGxvZ1RyYWNlKG1lc3NhZ2UsIHRyYWNlKTtcbiAgICAgICAgbW9kdWxlUmVuZGVyRXJyb3IucHVzaCh7XG4gICAgICAgICAgbWVzc2FnZSxcbiAgICAgICAgICBtb2R1bGVLZXksXG4gICAgICAgICAgc3RhY2s6IHRyYWNlLnN0YWNrLFxuICAgICAgICB9KTtcbiAgICAgICAgdGhyb3cgdHJhY2U7XG4gICAgICB9LFxuICAgIH0pO1xuICAgIC8vIOaNleiOt+mUgOavgeW8guW4uFxuICAgIGRlc3Ryb3kgPSBjYXJlKGRlc3Ryb3ksIHtcbiAgICAgIGNhcHR1cmU6IHRyYWNlID0+IHtcbiAgICAgICAgY29uc3QgbWVzc2FnZSA9IGDmqKHlnZcgJHttb2R1bGVLZXl9IOmUgOavgeW8guW4uO+8miR7dHJhY2UubWVzc2FnZX1gO1xuICAgICAgICBsb2dUcmFjZShtZXNzYWdlLCB0cmFjZSk7XG4gICAgICAgIG1vZHVsZURlc3Ryb3lFcnJvci5wdXNoKHtcbiAgICAgICAgICBtZXNzYWdlLFxuICAgICAgICAgIG1vZHVsZUtleSxcbiAgICAgICAgICBzdGFjazogdHJhY2Uuc3RhY2ssXG4gICAgICAgIH0pO1xuICAgICAgICB0aHJvdyB0cmFjZTtcbiAgICAgIH0sXG4gICAgfSk7XG4gIH1cbiAgcmV0dXJuIHsgcmVuZGVyLCBkZXN0cm95IH07XG59XG4iXX0=