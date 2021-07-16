"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var warning = require("warning");
var history_1 = require("history");
var createTransitionManager_1 = require("history/createTransitionManager");
/**
 * 利用控制台的 Router 对象创建兼容 React Router 的 history
 *
 * 参考 createBrowserHistory 实现
 * @see https://github.com/ReactTraining/history/blob/master/modules/createBrowserHistory.js
 */
function createHistory(props) {
    var globalHistory = window.history;
    var controller = props.controller, action = props.action, getUserConfirmation = props.getUserConfirmation, _a = props.keyLength, keyLength = _a === void 0 ? 6 : _a;
    var consoleRouter = seajs.require("router");
    var consoleRouteRule = "/" + controller + "/" + action + "(/*params)";
    var transitionManager = createTransitionManager_1.default();
    var initialLocation = getDOMLocation(getHistoryState());
    var allKeys = [initialLocation.key];
    var history = {
        length: globalHistory.length,
        action: "POP",
        location: initialLocation,
        createHref: createHref,
        push: perform("PUSH"),
        replace: perform("REPLACE"),
        go: go,
        goBack: goBack,
        goForward: goForward,
        block: block,
        listen: listen,
    };
    function perform(action) {
        return function (path, state) {
            warning(!(typeof path === "object" &&
                path.state !== undefined &&
                state !== undefined), "You should avoid providing a 2nd state argument to " + action.toLocaleLowerCase() + " when the 1st " +
                "argument is a location-like object that already has state; it is ignored");
            var location = history_1.createLocation(path, state, createKey(), history.location);
            transitionManager.confirmTransitionTo(location, action, getUserConfirmation, function (ok) {
                if (!ok)
                    return;
                var href = createHref(location);
                var key = location.key, state = location.state;
                // createBrowserHistory 的实现是调用 history.pushState()
                // 这里改为通过控制台路由代理
                muteNextPop = true;
                consoleRouter.navigateWithState(href, false, action === "REPLACE", {
                    key: key,
                    state: state,
                });
                muteNextPop = false;
                // 通过 allKeys 实现 pop 阻拦
                var prevIndex = allKeys.indexOf(history.location.key);
                if (action === "PUSH") {
                    var nextKeys = allKeys.slice(0, prevIndex === -1 ? 0 : prevIndex + 1);
                    nextKeys.push(location.key);
                    allKeys = nextKeys;
                }
                else if (prevIndex !== -1) {
                    allKeys[prevIndex] = location.key;
                }
                allKeys = allKeys.slice(-history.length - 1);
                setState({ action: action, location: location });
            });
        };
    }
    function go(n) {
        globalHistory.go(n);
    }
    function goBack() {
        go(-1);
    }
    function goForward() {
        go(1);
    }
    var listenerCount = 0;
    function checkRouterListeners(delta) {
        listenerCount += delta;
        if (listenerCount === 1 && delta === 1) {
            // add
            consoleRouter.use(consoleRouteRule, handlePopState);
        }
        else if (listenerCount === 0) {
            // remove
            consoleRouter.unuse(consoleRouteRule);
        }
    }
    function handlePopState() {
        var location = getDOMLocation(getHistoryState());
        handlePop(location);
    }
    var forceNextPop = false;
    var muteNextPop = false;
    function handlePop(location) {
        // push/replace 触发 console router.navigate() 方法
        // 同样会回调 handlePop()
        // 实际上，此时状态已经设置完毕，无需处理
        if (muteNextPop) {
            return;
        }
        if (forceNextPop) {
            forceNextPop = false;
            setState();
        }
        else {
            var action_1 = "POP";
            transitionManager.confirmTransitionTo(location, action_1, getUserConfirmation, function (ok) {
                if (ok) {
                    setState({ action: action_1, location: location });
                }
                else {
                    revertPop(location);
                }
            });
        }
    }
    function revertPop(fromLocation) {
        var toLocation = history.location;
        // TODO: We could probably make this more reliable by
        // keeping a list of keys we've seen in sessionStorage.
        // Instead, we just default to 0 for keys we don't know.
        var toIndex = allKeys.indexOf(toLocation.key);
        if (toIndex === -1)
            toIndex = 0;
        var fromIndex = allKeys.indexOf(fromLocation.key);
        if (fromIndex === -1)
            fromIndex = 0;
        var delta = toIndex - fromIndex;
        if (delta) {
            forceNextPop = true;
            go(delta);
        }
    }
    var isBlocked = false;
    function block(prompt) {
        if (prompt === void 0) { prompt = false; }
        var unblock = transitionManager.setPrompt(prompt);
        if (!isBlocked) {
            checkRouterListeners(1);
            isBlocked = true;
        }
        return function () {
            if (isBlocked) {
                isBlocked = false;
                checkRouterListeners(-1);
            }
            return unblock();
        };
    }
    function listen(listener) {
        var unlisten = transitionManager.appendListener(listener);
        checkRouterListeners(1);
        return function () {
            checkRouterListeners(-1);
            unlisten();
        };
    }
    function createHref(location) {
        return history_1.createPath(location);
    }
    function createKey() {
        return Math.random()
            .toString(36)
            .substr(2, keyLength);
    }
    function setState(nextState) {
        if (nextState) {
            Object.assign(history, nextState);
        }
        history.length = globalHistory.length;
        transitionManager.notifyListeners(history.location, history.action);
    }
    return history;
}
exports.createHistory = createHistory;
function getHistoryState() {
    try {
        return window.history.state || {};
    }
    catch (e) {
        // IE 11 sometimes throws when accessing window.history.state
        // See https://github.com/ReactTraining/history/pull/289
        return {};
    }
}
function getDOMLocation(historyState) {
    var _a = historyState || {}, _b = _a.key, key = _b === void 0 ? undefined : _b, _c = _a.state, state = _c === void 0 ? undefined : _c;
    var _d = window.location, pathname = _d.pathname, search = _d.search, hash = _d.hash;
    var path = pathname + search + hash;
    return history_1.createLocation(path, state, key);
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY3JlYXRlLWhpc3RvcnkuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY29yZS9oaXN0b3J5L2NyZWF0ZS1oaXN0b3J5LnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEsaUNBQW1DO0FBQ25DLG1DQVVpQjtBQUNqQiwyRUFFeUM7QUE2QnpDOzs7OztHQUtHO0FBQ0gsU0FBZ0IsYUFBYSxDQUFJLEtBQTZCO0lBQzVELElBQU0sYUFBYSxHQUFHLE1BQU0sQ0FBQyxPQUFPLENBQUM7SUFFN0IsSUFBQSw2QkFBVSxFQUFFLHFCQUFNLEVBQUUsK0NBQW1CLEVBQUUsb0JBQWEsRUFBYixrQ0FBYSxDQUFXO0lBRXpFLElBQU0sYUFBYSxHQUFrQixLQUFLLENBQUMsT0FBTyxDQUFDLFFBQVEsQ0FBQyxDQUFDO0lBQzdELElBQU0sZ0JBQWdCLEdBQUcsTUFBSSxVQUFVLFNBQUksTUFBTSxlQUFZLENBQUM7SUFFOUQsSUFBTSxpQkFBaUIsR0FBRyxpQ0FBdUIsRUFBRSxDQUFDO0lBQ3BELElBQU0sZUFBZSxHQUFHLGNBQWMsQ0FBQyxlQUFlLEVBQUUsQ0FBQyxDQUFDO0lBQzFELElBQUksT0FBTyxHQUFHLENBQUMsZUFBZSxDQUFDLEdBQUcsQ0FBQyxDQUFDO0lBRXBDLElBQU0sT0FBTyxHQUFHO1FBQ2QsTUFBTSxFQUFFLGFBQWEsQ0FBQyxNQUFNO1FBQzVCLE1BQU0sRUFBRSxLQUFlO1FBQ3ZCLFFBQVEsRUFBRSxlQUFlO1FBQ3pCLFVBQVUsWUFBQTtRQUNWLElBQUksRUFBRSxPQUFPLENBQUMsTUFBTSxDQUFDO1FBQ3JCLE9BQU8sRUFBRSxPQUFPLENBQUMsU0FBUyxDQUFDO1FBQzNCLEVBQUUsSUFBQTtRQUNGLE1BQU0sUUFBQTtRQUNOLFNBQVMsV0FBQTtRQUNULEtBQUssT0FBQTtRQUNMLE1BQU0sUUFBQTtLQUNQLENBQUM7SUFFRixTQUFTLE9BQU8sQ0FBQyxNQUFjO1FBQzdCLE9BQU8sVUFBQyxJQUE0QixFQUFFLEtB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RUFBRSxDQUFDO2FBQ1osTUFBTSxDQUFDLENBQUMsRUFBRSxTQUFTLENBQUMsQ0FBQztJQUMxQixDQUFDO0lBRUQsU0FBUyxRQUFRLENBQUMsU0FBZ0Q7UUFDaEUsSUFBSSxTQUFTLEVBQUU7WUFDYixNQUFNLENBQUMsTUFBTSxDQUFDLE9BQU8sRUFBRSxTQUFTLENBQUMsQ0FBQztTQUNuQztRQUNELE9BQU8sQ0FBQyxNQUFNLEdBQUcsYUFBYSxDQUFDLE1BQU0sQ0FBQztRQUN0QyxpQkFBaUIsQ0FBQyxlQUFlLENBQUMsT0FBTyxDQUFDLFFBQVEsRUFBRSxPQUFPLENBQUMsTUFBTSxDQUFDLENBQUM7SUFDdEUsQ0FBQztJQUVELE9BQU8sT0FBa0IsQ0FBQztBQUM1QixDQUFDO0FBdk5ELHNDQXVOQztBQUVELFNBQVMsZUFBZTtJQUN0QixJQUFJO1FBQ0YsT0FBTyxNQUFNLENBQUMsT0FBTyxDQUFDLEtBQUssSUFBSSxFQUFFLENBQUM7S0FDbkM7SUFBQyxPQUFPLENBQUMsRUFBRTtRQUNWLDZEQUE2RDtRQUM3RCx3REFBd0Q7UUFDeEQsT0FBTyxFQUFFLENBQUM7S0FDWDtBQUNILENBQUM7QUFFRCxTQUFTLGNBQWMsQ0FBQyxZQUFpQjtJQUNqQyxJQUFBLHVCQUEyRCxFQUF6RCxXQUFlLEVBQWYsb0NBQWUsRUFBRSxhQUFpQixFQUFqQixzQ0FBd0MsQ0FBQztJQUM1RCxJQUFBLG9CQUE0QyxFQUExQyxzQkFBUSxFQUFFLGtCQUFNLEVBQUUsY0FBd0IsQ0FBQztJQUVuRCxJQUFJLElBQUksR0FBRyxRQUFRLEdBQUcsTUFBTSxHQUFHLElBQUksQ0FBQztJQUVwQyxPQUFPLHdCQUFjLENBQUMsSUFBSSxFQUFFLEtBQUssRUFBRSxHQUFHLENBQUMsQ0FBQztBQUMxQyxDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0ICogYXMgd2FybmluZyBmcm9tIFwid2FybmluZ1wiO1xuaW1wb3J0IHtcbiAgY3JlYXRlUGF0aCxcbiAgY3JlYXRlTG9jYXRpb24sXG4gIEJyb3dzZXJIaXN0b3J5QnVpbGRPcHRpb25zLFxuICBIaXN0b3J5LFxuICBMb2NhdGlvblN0YXRlLFxuICBMb2NhdGlvbkRlc2NyaXB0b3JPYmplY3QsXG4gIEFjdGlvbixcbiAgTG9jYXRpb25MaXN0ZW5lcixcbiAgTG9jYXRpb24sXG59IGZyb20gXCJoaXN0b3J5XCI7XG5pbXBvcnQgY3JlYXRlVHJhbnNpdGlvbk1hbmFnZXIsIHtcbiAgUHJvbXB0LFxufSBmcm9tIFwiaGlzdG9yeS9jcmVhdGVUcmFuc2l0aW9uTWFuYWdlclwiO1xuXG5pbnRlcmZhY2UgQ29uc29sZVJvdXRlciB7XG4gIC8qKlxuICAgKiDmjqfliLblj7Dot6/nlLHmlrnms5VcbiAgICovXG4gIG5hdmlnYXRlV2l0aFN0YXRlKFxuICAgIHVybDogc3RyaW5nLFxuICAgIHNpbGVudDogYm9vbGVhbixcbiAgICByZXBsYWNlbWVudDogYm9vbGVhbixcbiAgICBzdGF0ZTogYW55XG4gICk6IHZvaWQ7XG5cbiAgLyoqXG4gICAqIOazqOWGjOWKqOaAgei3r+eUsVxuICAgKi9cbiAgdXNlKHJ1bGU6IHN0cmluZywgYWN0aW9uOiBGdW5jdGlvbik6IHZvaWQ7XG5cbiAgLyoqXG4gICAqIOWPlua2iOazqOWGjOWKqOaAgei3r+eUsVxuICAgKi9cbiAgdW51c2UocnVsZTogc3RyaW5nKTogdm9pZDtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBUZWFIaXN0b3J5QnVpbGRPcHRpb25zIGV4dGVuZHMgQnJvd3Nlckhpc3RvcnlCdWlsZE9wdGlvbnMge1xuICBjb250cm9sbGVyPzogc3RyaW5nO1xuICBhY3Rpb24/OiBzdHJpbmc7XG59XG5cbi8qKlxuICog5Yip55So5o6n5Yi25Y+w55qEIFJvdXRlciDlr7nosaHliJvlu7rlhbzlrrkgUmVhY3QgUm91dGVyIOeahCBoaXN0b3J5XG4gKlxuICog5Y+C6ICDIGNyZWF0ZUJyb3dzZXJIaXN0b3J5IOWunueOsFxuICogQHNlZSBodHRwczovL2dpdGh1Yi5jb20vUmVhY3RUcmFpbmluZy9oaXN0b3J5L2Jsb2IvbWFzdGVyL21vZHVsZXMvY3JlYXRlQnJvd3Nlckhpc3RvcnkuanNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGNyZWF0ZUhpc3Rvcnk8Uz4ocHJvcHM6IFRlYUhpc3RvcnlCdWlsZE9wdGlvbnMpOiBIaXN0b3J5PFM+IHtcbiAgY29uc3QgZ2xvYmFsSGlzdG9yeSA9IHdpbmRvdy5oaXN0b3J5O1xuXG4gIGNvbnN0IHsgY29udHJvbGxlciwgYWN0aW9uLCBnZXRVc2VyQ29uZmlybWF0aW9uLCBrZXlMZW5ndGggPSA2IH0gPSBwcm9wcztcblxuICBjb25zdCBjb25zb2xlUm91dGVyOiBDb25zb2xlUm91dGVyID0gc2VhanMucmVxdWlyZShcInJvdXRlclwiKTtcbiAgY29uc3QgY29uc29sZVJvdXRlUnVsZSA9IGAvJHtjb250cm9sbGVyfS8ke2FjdGlvbn0oLypwYXJhbXMpYDtcblxuICBjb25zdCB0cmFuc2l0aW9uTWFuYWdlciA9IGNyZWF0ZVRyYW5zaXRpb25NYW5hZ2VyKCk7XG4gIGNvbnN0IGluaXRpYWxMb2NhdGlvbiA9IGdldERPTUxvY2F0aW9uKGdldEhpc3RvcnlTdGF0ZSgpKTtcbiAgbGV0IGFsbEtleXMgPSBbaW5pdGlhbExvY2F0aW9uLmtleV07XG5cbiAgY29uc3QgaGlzdG9yeSA9IHtcbiAgICBsZW5ndGg6IGdsb2JhbEhpc3RvcnkubGVuZ3RoLFxuICAgIGFjdGlvbjogXCJQT1BcIiBhcyBBY3Rpb24sXG4gICAgbG9jYXRpb246IGluaXRpYWxMb2NhdGlvbixcbiAgICBjcmVhdGVIcmVmLFxuICAgIHB1c2g6IHBlcmZvcm0oXCJQVVNIXCIpLFxuICAgIHJlcGxhY2U6IHBlcmZvcm0oXCJSRVBMQUNFXCIpLFxuICAgIGdvLFxuICAgIGdvQmFjayxcbiAgICBnb0ZvcndhcmQsXG4gICAgYmxvY2ssXG4gICAgbGlzdGVuLFxuICB9O1xuXG4gIGZ1bmN0aW9uIHBlcmZvcm0oYWN0aW9uOiBBY3Rpb24pIHtcbiAgICByZXR1cm4gKHBhdGg6IHN0cmluZyB8IExvY2F0aW9uU3RhdGUsIHN0YXRlPzogTG9jYXRpb25TdGF0ZSkgPT4ge1xuICAgICAgd2FybmluZyhcbiAgICAgICAgIShcbiAgICAgICAgICB0eXBlb2YgcGF0aCA9PT0gXCJvYmplY3RcIiAmJlxuICAgICAgICAgIHBhdGguc3RhdGUgIT09IHVuZGVmaW5lZCAmJlxuICAgICAgICAgIHN0YXRlICE9PSB1bmRlZmluZWRcbiAgICAgICAgKSxcbiAgICAgICAgYFlvdSBzaG91bGQgYXZvaWQgcHJvdmlkaW5nIGEgMm5kIHN0YXRlIGFyZ3VtZW50IHRvICR7YWN0aW9uLnRvTG9jYWxlTG93ZXJDYXNlKCl9IHdoZW4gdGhlIDFzdCBgICtcbiAgICAgICAgICBcImFyZ3VtZW50IGlzIGEgbG9jYXRpb24tbGlrZSBvYmplY3QgdGhhdCBhbHJlYWR5IGhhcyBzdGF0ZTsgaXQgaXMgaWdub3JlZFwiXG4gICAgICApO1xuICAgICAgY29uc3QgbG9jYXRpb24gPSBjcmVhdGVMb2NhdGlvbihcbiAgICAgICAgcGF0aCxcbiAgICAgICAgc3RhdGUsXG4gICAgICAgIGNyZWF0ZUtleSgpLFxuICAgICAgICBoaXN0b3J5LmxvY2F0aW9uXG4gICAgICApO1xuXG4gICAgICB0cmFuc2l0aW9uTWFuYWdlci5jb25maXJtVHJhbnNpdGlvblRvKFxuICAgICAgICBsb2NhdGlvbixcbiAgICAgICAgYWN0aW9uLFxuICAgICAgICBnZXRVc2VyQ29uZmlybWF0aW9uLFxuICAgICAgICBvayA9PiB7XG4gICAgICAgICAgaWYgKCFvaykgcmV0dXJuO1xuXG4gICAgICAgICAgY29uc3QgaHJlZiA9IGNyZWF0ZUhyZWYobG9jYXRpb24pO1xuICAgICAgICAgIGNvbnN0IHsga2V5LCBzdGF0ZSB9ID0gbG9jYXRpb247XG5cbiAgICAgICAgICAvLyBjcmVhdGVCcm93c2VySGlzdG9yeSDnmoTlrp7njrDmmK/osIPnlKggaGlzdG9yeS5wdXNoU3RhdGUoKVxuICAgICAgICAgIC8vIOi/memHjOaUueS4uumAmui/h+aOp+WItuWPsOi3r+eUseS7o+eQhlxuICAgICAgICAgIG11dGVOZXh0UG9wID0gdHJ1ZTtcbiAgICAgICAgICBjb25zb2xlUm91dGVyLm5hdmlnYXRlV2l0aFN0YXRlKGhyZWYsIGZhbHNlLCBhY3Rpb24gPT09IFwiUkVQTEFDRVwiLCB7XG4gICAgICAgICAgICBrZXksXG4gICAgICAgICAgICBzdGF0ZSxcbiAgICAgICAgICB9KTtcbiAgICAgICAgICBtdXRlTmV4dFBvcCA9IGZhbHNlO1xuXG4gICAgICAgICAgLy8g6YCa6L+HIGFsbEtleXMg5a6e546wIHBvcCDpmLvmi6ZcbiAgICAgICAgICBjb25zdCBwcmV2SW5kZXggPSBhbGxLZXlzLmluZGV4T2YoaGlzdG9yeS5sb2NhdGlvbi5rZXkpO1xuICAgICAgICAgIGlmIChhY3Rpb24gPT09IFwiUFVTSFwiKSB7XG4gICAgICAgICAgICBjb25zdCBuZXh0S2V5cyA9IGFsbEtleXMuc2xpY2UoXG4gICAgICAgICAgICAgIDAsXG4gICAgICAgICAgICAgIHByZXZJbmRleCA9PT0gLTEgPyAwIDogcHJldkluZGV4ICsgMVxuICAgICAgICAgICAgKTtcbiAgICAgICAgICAgIG5leHRLZXlzLnB1c2gobG9jYXRpb24ua2V5KTtcbiAgICAgICAgICAgIGFsbEtleXMgPSBuZXh0S2V5cztcbiAgICAgICAgICB9IGVsc2UgaWYgKHByZXZJbmRleCAhPT0gLTEpIHtcbiAgICAgICAgICAgIGFsbEtleXNbcHJldkluZGV4XSA9IGxvY2F0aW9uLmtleTtcbiAgICAgICAgICB9XG4gICAgICAgICAgYWxsS2V5cyA9IGFsbEtleXMuc2xpY2UoLWhpc3RvcnkubGVuZ3RoIC0gMSk7XG5cbiAgICAgICAgICBzZXRTdGF0ZSh7IGFjdGlvbiwgbG9jYXRpb24gfSk7XG4gICAgICAgIH1cbiAgICAgICk7XG4gICAgfTtcbiAgfVxuXG4gIGZ1bmN0aW9uIGdvKG46IG51bWJlcikge1xuICAgIGdsb2JhbEhpc3RvcnkuZ28obik7XG4gIH1cblxuICBmdW5jdGlvbiBnb0JhY2soKSB7XG4gICAgZ28oLTEpO1xuICB9XG5cbiAgZnVuY3Rpb24gZ29Gb3J3YXJkKCkge1xuICAgIGdvKDEpO1xuICB9XG5cbiAgbGV0IGxpc3RlbmVyQ291bnQgPSAwO1xuICBmdW5jdGlvbiBjaGVja1JvdXRlckxpc3RlbmVycyhkZWx0YSkge1xuICAgIGxpc3RlbmVyQ291bnQgKz0gZGVsdGE7XG5cbiAgICBpZiAobGlzdGVuZXJDb3VudCA9PT0gMSAmJiBkZWx0YSA9PT0gMSkge1xuICAgICAgLy8gYWRkXG4gICAgICBjb25zb2xlUm91dGVyLnVzZShjb25zb2xlUm91dGVSdWxlLCBoYW5kbGVQb3BTdGF0ZSk7XG4gICAgfSBlbHNlIGlmIChsaXN0ZW5lckNvdW50ID09PSAwKSB7XG4gICAgICAvLyByZW1vdmVcbiAgICAgIGNvbnNvbGVSb3V0ZXIudW51c2UoY29uc29sZVJvdXRlUnVsZSk7XG4gICAgfVxuICB9XG5cbiAgZnVuY3Rpb24gaGFuZGxlUG9wU3RhdGUoKSB7XG4gICAgY29uc3QgbG9jYXRpb24gPSBnZXRET01Mb2NhdGlvbihnZXRIaXN0b3J5U3RhdGUoKSk7XG4gICAgaGFuZGxlUG9wKGxvY2F0aW9uKTtcbiAgfVxuXG4gIGxldCBmb3JjZU5leHRQb3AgPSBmYWxzZTtcbiAgbGV0IG11dGVOZXh0UG9wID0gZmFsc2U7XG5cbiAgZnVuY3Rpb24gaGFuZGxlUG9wKGxvY2F0aW9uOiBMb2NhdGlvbikge1xuICAgIC8vIHB1c2gvcmVwbGFjZSDop6blj5EgY29uc29sZSByb3V0ZXIubmF2aWdhdGUoKSDmlrnms5VcbiAgICAvLyDlkIzmoLfkvJrlm57osIMgaGFuZGxlUG9wKClcbiAgICAvLyDlrp7pmYXkuIrvvIzmraTml7bnirbmgIHlt7Lnu4/orr7nva7lrozmr5XvvIzml6DpnIDlpITnkIZcbiAgICBpZiAobXV0ZU5leHRQb3ApIHtcbiAgICAgIHJldHVybjtcbiAgICB9XG4gICAgaWYgKGZvcmNlTmV4dFBvcCkge1xuICAgICAgZm9yY2VOZXh0UG9wID0gZmFsc2U7XG4gICAgICBzZXRTdGF0ZSgpO1xuICAgIH0gZWxzZSB7XG4gICAgICBjb25zdCBhY3Rpb24gPSBcIlBPUFwiO1xuXG4gICAgICB0cmFuc2l0aW9uTWFuYWdlci5jb25maXJtVHJhbnNpdGlvblRvKFxuICAgICAgICBsb2NhdGlvbixcbiAgICAgICAgYWN0aW9uLFxuICAgICAgICBnZXRVc2VyQ29uZmlybWF0aW9uLFxuICAgICAgICBvayA9PiB7XG4gICAgICAgICAgaWYgKG9rKSB7XG4gICAgICAgICAgICBzZXRTdGF0ZSh7IGFjdGlvbiwgbG9jYXRpb24gfSk7XG4gICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHJldmVydFBvcChsb2NhdGlvbik7XG4gICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICApO1xuICAgIH1cbiAgfVxuXG4gIGZ1bmN0aW9uIHJldmVydFBvcChmcm9tTG9jYXRpb24pIHtcbiAgICBjb25zdCB0b0xvY2F0aW9uID0gaGlzdG9yeS5sb2NhdGlvbjtcblxuICAgIC8vIFRPRE86IFdlIGNvdWxkIHByb2JhYmx5IG1ha2UgdGhpcyBtb3JlIHJlbGlhYmxlIGJ5XG4gICAgLy8ga2VlcGluZyBhIGxpc3Qgb2Yga2V5cyB3ZSd2ZSBzZWVuIGluIHNlc3Npb25TdG9yYWdlLlxuICAgIC8vIEluc3RlYWQsIHdlIGp1c3QgZGVmYXVsdCB0byAwIGZvciBrZXlzIHdlIGRvbid0IGtub3cuXG5cbiAgICBsZXQgdG9JbmRleCA9IGFsbEtleXMuaW5kZXhPZih0b0xvY2F0aW9uLmtleSk7XG5cbiAgICBpZiAodG9JbmRleCA9PT0gLTEpIHRvSW5kZXggPSAwO1xuXG4gICAgbGV0IGZyb21JbmRleCA9IGFsbEtleXMuaW5kZXhPZihmcm9tTG9jYXRpb24ua2V5KTtcblxuICAgIGlmIChmcm9tSW5kZXggPT09IC0xKSBmcm9tSW5kZXggPSAwO1xuXG4gICAgY29uc3QgZGVsdGEgPSB0b0luZGV4IC0gZnJvbUluZGV4O1xuXG4gICAgaWYgKGRlbHRhKSB7XG4gICAgICBmb3JjZU5leHRQb3AgPSB0cnVlO1xuICAgICAgZ28oZGVsdGEpO1xuICAgIH1cbiAgfVxuXG4gIGxldCBpc0Jsb2NrZWQgPSBmYWxzZTtcbiAgZnVuY3Rpb24gYmxvY2socHJvbXB0OiBQcm9tcHQgPSBmYWxzZSkge1xuICAgIGNvbnN0IHVuYmxvY2sgPSB0cmFuc2l0aW9uTWFuYWdlci5zZXRQcm9tcHQocHJvbXB0KTtcblxuICAgIGlmICghaXNCbG9ja2VkKSB7XG4gICAgICBjaGVja1JvdXRlckxpc3RlbmVycygxKTtcbiAgICAgIGlzQmxvY2tlZCA9IHRydWU7XG4gICAgfVxuXG4gICAgcmV0dXJuICgpID0+IHtcbiAgICAgIGlmIChpc0Jsb2NrZWQpIHtcbiAgICAgICAgaXNCbG9ja2VkID0gZmFsc2U7XG4gICAgICAgIGNoZWNrUm91dGVyTGlzdGVuZXJzKC0xKTtcbiAgICAgIH1cblxuICAgICAgcmV0dXJuIHVuYmxvY2soKTtcbiAgICB9O1xuICB9XG5cbiAgZnVuY3Rpb24gbGlzdGVuKGxpc3RlbmVyOiBMb2NhdGlvbkxpc3RlbmVyKSB7XG4gICAgY29uc3QgdW5saXN0ZW4gPSB0cmFuc2l0aW9uTWFuYWdlci5hcHBlbmRMaXN0ZW5lcihsaXN0ZW5lcik7XG4gICAgY2hlY2tSb3V0ZXJMaXN0ZW5lcnMoMSk7XG5cbiAgICByZXR1cm4gKCkgPT4ge1xuICAgICAgY2hlY2tSb3V0ZXJMaXN0ZW5lcnMoLTEpO1xuICAgICAgdW5saXN0ZW4oKTtcbiAgICB9O1xuICB9XG5cbiAgZnVuY3Rpb24gY3JlYXRlSHJlZihsb2NhdGlvbjogTG9jYXRpb25EZXNjcmlwdG9yT2JqZWN0KSB7XG4gICAgcmV0dXJuIGNyZWF0ZVBhdGgobG9jYXRpb24pO1xuICB9XG5cbiAgZnVuY3Rpb24gY3JlYXRlS2V5KCkge1xuICAgIHJldHVybiBNYXRoLnJhbmRvbSgpXG4gICAgICAudG9TdHJpbmcoMzYpXG4gICAgICAuc3Vic3RyKDIsIGtleUxlbmd0aCk7XG4gIH1cblxuICBmdW5jdGlvbiBzZXRTdGF0ZShuZXh0U3RhdGU/OiBQaWNrPEhpc3RvcnksIFwiYWN0aW9uXCIgfCBcImxvY2F0aW9uXCI+KSB7XG4gICAgaWYgKG5leHRTdGF0ZSkge1xuICAgICAgT2JqZWN0LmFzc2lnbihoaXN0b3J5LCBuZXh0U3RhdGUpO1xuICAgIH1cbiAgICBoaXN0b3J5Lmxlbmd0aCA9IGdsb2JhbEhpc3RvcnkubGVuZ3RoO1xuICAgIHRyYW5zaXRpb25NYW5hZ2VyLm5vdGlmeUxpc3RlbmVycyhoaXN0b3J5LmxvY2F0aW9uLCBoaXN0b3J5LmFjdGlvbik7XG4gIH1cblxuICByZXR1cm4gaGlzdG9yeSBhcyBIaXN0b3J5O1xufVxuXG5mdW5jdGlvbiBnZXRIaXN0b3J5U3RhdGUoKSB7XG4gIHRyeSB7XG4gICAgcmV0dXJuIHdpbmRvdy5oaXN0b3J5LnN0YXRlIHx8IHt9O1xuICB9IGNhdGNoIChlKSB7XG4gICAgLy8gSUUgMTEgc29tZXRpbWVzIHRocm93cyB3aGVuIGFjY2Vzc2luZyB3aW5kb3cuaGlzdG9yeS5zdGF0ZVxuICAgIC8vIFNlZSBodHRwczovL2dpdGh1Yi5jb20vUmVhY3RUcmFpbmluZy9oaXN0b3J5L3B1bGwvMjg5XG4gICAgcmV0dXJuIHt9O1xuICB9XG59XG5cbmZ1bmN0aW9uIGdldERPTUxvY2F0aW9uKGhpc3RvcnlTdGF0ZTogYW55KSB7XG4gIGNvbnN0IHsga2V5ID0gdW5kZWZpbmVkLCBzdGF0ZSA9IHVuZGVmaW5lZCB9ID0gaGlzdG9yeVN0YXRlIHx8IHt9O1xuICBjb25zdCB7IHBhdGhuYW1lLCBzZWFyY2gsIGhhc2ggfSA9IHdpbmRvdy5sb2NhdGlvbjtcblxuICBsZXQgcGF0aCA9IHBhdGhuYW1lICsgc2VhcmNoICsgaGFzaDtcblxuICByZXR1cm4gY3JlYXRlTG9jYXRpb24ocGF0aCwgc3RhdGUsIGtleSk7XG59XG4iXX0=