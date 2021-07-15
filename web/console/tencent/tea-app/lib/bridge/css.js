"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var presetCache = new Map();
var CDN_HOSTNAME = window["QCCDN_HOST"] || "imgcache.qq.com";
/**
 * @deprecated
 */
exports.css = {
    /**
     * @deprecated
     *
     * 返回内置 CSS URL
     * @param preset 内置 CSS 名称
     */
    preset: function (preset) {
        // 控制台把预置 CSS 版本直出到全局变量 `CSS_PRESET_VERSIONS` 中
        var presetMap = window["CSS_PRESET_VERSIONS"] || {};
        if (!presetMap[preset]) {
            console.warn("找不到预定义样式配置：%s", preset);
        }
        return presetMap[preset] || null;
    },
    /**
     * @deprecated
     *
     * 返回 CDN URL
     * @param pathname CDN 路径
     */
    cdn: function (pathname) {
        if (!pathname.startsWith("/")) {
            pathname = "/" + pathname;
        }
        return "//" + CDN_HOSTNAME + pathname;
    },
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY3NzLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vc3JjL2JyaWRnZS9jc3MudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSxJQUFNLFdBQVcsR0FBRyxJQUFJLEdBQUcsRUFBa0IsQ0FBQztBQUU5QyxJQUFJLFlBQVksR0FBRyxNQUFNLENBQUMsWUFBWSxDQUFDLElBQUksaUJBQWlCLENBQUM7QUFFN0Q7O0dBRUc7QUFDVSxRQUFBLEdBQUcsR0FBRztJQUNqQjs7Ozs7T0FLRztJQUNILE1BQU0sRUFBTixVQUFPLE1BQWlCO1FBQ3RCLCtDQUErQztRQUMvQyxJQUFNLFNBQVMsR0FBRyxNQUFNLENBQUMscUJBQXFCLENBQUMsSUFBSSxFQUFFLENBQUM7UUFDdEQsSUFBSSxDQUFDLFNBQVMsQ0FBQyxNQUFNLENBQUMsRUFBRTtZQUN0QixPQUFPLENBQUMsSUFBSSxDQUFDLGVBQWUsRUFBRSxNQUFNLENBQUMsQ0FBQztTQUN2QztRQUNELE9BQU8sU0FBUyxDQUFDLE1BQU0sQ0FBQyxJQUFJLElBQUksQ0FBQztJQUNuQyxDQUFDO0lBRUQ7Ozs7O09BS0c7SUFDSCxHQUFHLEVBQUgsVUFBSSxRQUFnQjtRQUNsQixJQUFJLENBQUMsUUFBUSxDQUFDLFVBQVUsQ0FBQyxHQUFHLENBQUMsRUFBRTtZQUM3QixRQUFRLEdBQUcsR0FBRyxHQUFHLFFBQVEsQ0FBQztTQUMzQjtRQUNELE9BQU8sT0FBSyxZQUFZLEdBQUcsUUFBVSxDQUFDO0lBQ3hDLENBQUM7Q0FDRixDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiY29uc3QgcHJlc2V0Q2FjaGUgPSBuZXcgTWFwPHN0cmluZywgc3RyaW5nPigpO1xuXG5sZXQgQ0ROX0hPU1ROQU1FID0gd2luZG93W1wiUUNDRE5fSE9TVFwiXSB8fCBcImltZ2NhY2hlLnFxLmNvbVwiO1xuXG4vKipcbiAqIEBkZXByZWNhdGVkXG4gKi9cbmV4cG9ydCBjb25zdCBjc3MgPSB7XG4gIC8qKlxuICAgKiBAZGVwcmVjYXRlZFxuICAgKlxuICAgKiDov5Tlm57lhoXnva4gQ1NTIFVSTFxuICAgKiBAcGFyYW0gcHJlc2V0IOWGhee9riBDU1Mg5ZCN56ewXG4gICAqL1xuICBwcmVzZXQocHJlc2V0OiBDU1NQcmVzZXQpIHtcbiAgICAvLyDmjqfliLblj7DmiorpooTnva4gQ1NTIOeJiOacrOebtOWHuuWIsOWFqOWxgOWPmOmHjyBgQ1NTX1BSRVNFVF9WRVJTSU9OU2Ag5LitXG4gICAgY29uc3QgcHJlc2V0TWFwID0gd2luZG93W1wiQ1NTX1BSRVNFVF9WRVJTSU9OU1wiXSB8fCB7fTtcbiAgICBpZiAoIXByZXNldE1hcFtwcmVzZXRdKSB7XG4gICAgICBjb25zb2xlLndhcm4oXCLmib7kuI3liLDpooTlrprkuYnmoLflvI/phY3nva7vvJolc1wiLCBwcmVzZXQpO1xuICAgIH1cbiAgICByZXR1cm4gcHJlc2V0TWFwW3ByZXNldF0gfHwgbnVsbDtcbiAgfSxcblxuICAvKipcbiAgICogQGRlcHJlY2F0ZWRcbiAgICpcbiAgICog6L+U5ZueIENETiBVUkxcbiAgICogQHBhcmFtIHBhdGhuYW1lIENETiDot6/lvoRcbiAgICovXG4gIGNkbihwYXRobmFtZTogc3RyaW5nKSB7XG4gICAgaWYgKCFwYXRobmFtZS5zdGFydHNXaXRoKFwiL1wiKSkge1xuICAgICAgcGF0aG5hbWUgPSBcIi9cIiArIHBhdGhuYW1lO1xuICAgIH1cbiAgICByZXR1cm4gYC8vJHtDRE5fSE9TVE5BTUV9JHtwYXRobmFtZX1gO1xuICB9LFxufTtcblxuZXhwb3J0IHR5cGUgQ1NTUHJlc2V0ID0gXCJ0ZWEtY29tcG9uZW50XCI7XG4iXX0=