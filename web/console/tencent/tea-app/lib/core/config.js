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
/**
 * 在原有的模块导出的基础上，附加菜单和样式的配置导出
 */
exports.config = function (factory) { return function (controller, action, entry) { return function () {
    return __assign({}, factory(controller, action, entry)());
}; }; };
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29uZmlnLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vc3JjL2NvcmUvY29uZmlnLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7QUFDQTs7R0FFRztBQUNVLFFBQUEsTUFBTSxHQUFXLFVBQUEsT0FBTyxJQUFJLE9BQUEsVUFBQyxVQUFVLEVBQUUsTUFBTSxFQUFFLEtBQUssSUFBSyxPQUFBO0lBQ3RFLG9CQUNLLE9BQU8sQ0FBQyxVQUFVLEVBQUUsTUFBTSxFQUFFLEtBQUssQ0FBQyxFQUFFLEVBQ3ZDO0FBQ0osQ0FBQyxFQUp1RSxDQUl2RSxFQUp3QyxDQUl4QyxDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgRW50cnlNb2R1bGVGYWN0b3J5IH0gZnJvbSBcIi4vZW50cnlcIjtcbi8qKlxuICog5Zyo5Y6f5pyJ55qE5qih5Z2X5a+85Ye655qE5Z+656GA5LiK77yM6ZmE5Yqg6I+c5Y2V5ZKM5qC35byP55qE6YWN572u5a+85Ye6XG4gKi9cbmV4cG9ydCBjb25zdCBjb25maWc6IENvbmZpZyA9IGZhY3RvcnkgPT4gKGNvbnRyb2xsZXIsIGFjdGlvbiwgZW50cnkpID0+ICgpID0+IHtcbiAgcmV0dXJuIHtcbiAgICAuLi5mYWN0b3J5KGNvbnRyb2xsZXIsIGFjdGlvbiwgZW50cnkpKCksXG4gIH07XG59O1xuXG5leHBvcnQgaW50ZXJmYWNlIENvbmZpZyB7XG4gIChmYWN0b3J5OiBFbnRyeU1vZHVsZUZhY3RvcnkpOiBFbnRyeU1vZHVsZUZhY3Rvcnk7XG59XG4iXX0=