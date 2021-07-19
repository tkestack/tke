"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("./_bridge");
var react_1 = require("react");
var postfix = null;
// 从 nmcConfig 中获得已经翻译的后缀
var nmcConfig = _bridge_1._bridge("nmcConfig");
if (nmcConfig && typeof nmcConfig.defaultTitle === "string") {
    postfix = nmcConfig.defaultTitle.split(/\s*-\s*/g).pop();
}
/**
 * 设置文档标题
 * @param title 文档标题
 */
function setDocumentTitle(title) {
    setTimeout(function () {
        document.title = postfix ? title + " - " + postfix : title;
    }, 1);
}
exports.setDocumentTitle = setDocumentTitle;
/**
 * React Hooks 声明文档标题
 * @param title 文档标题
 */
function useDocumentTitle(title) {
    react_1.useEffect(function () { return setDocumentTitle(title); }, [title]);
}
exports.useDocumentTitle = useDocumentTitle;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidGl0bGUuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi9zcmMvYnJpZGdlL3RpdGxlLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7O0FBQUEscUNBQW9DO0FBQ3BDLCtCQUFrQztBQUVsQyxJQUFJLE9BQU8sR0FBVyxJQUFJLENBQUM7QUFFM0IseUJBQXlCO0FBQ3pCLElBQU0sU0FBUyxHQUFHLGlCQUFPLENBQUMsV0FBVyxDQUFDLENBQUM7QUFDdkMsSUFBSSxTQUFTLElBQUksT0FBTyxTQUFTLENBQUMsWUFBWSxLQUFLLFFBQVEsRUFBRTtJQUMzRCxPQUFPLEdBQUksU0FBUyxDQUFDLFlBQXVCLENBQUMsS0FBSyxDQUFDLFVBQVUsQ0FBQyxDQUFDLEdBQUcsRUFBRSxDQUFDO0NBQ3RFO0FBRUQ7OztHQUdHO0FBQ0gsU0FBZ0IsZ0JBQWdCLENBQUMsS0FBYTtJQUM1QyxVQUFVLENBQUM7UUFDVCxRQUFRLENBQUMsS0FBSyxHQUFHLE9BQU8sQ0FBQyxDQUFDLENBQUMsS0FBSyxHQUFHLEtBQUssR0FBRyxPQUFPLENBQUMsQ0FBQyxDQUFDLEtBQUssQ0FBQztJQUM3RCxDQUFDLEVBQUUsQ0FBQyxDQUFDLENBQUM7QUFDUixDQUFDO0FBSkQsNENBSUM7QUFFRDs7O0dBR0c7QUFDSCxTQUFnQixnQkFBZ0IsQ0FBQyxLQUFhO0lBQzVDLGlCQUFTLENBQUMsY0FBTSxPQUFBLGdCQUFnQixDQUFDLEtBQUssQ0FBQyxFQUF2QixDQUF1QixFQUFFLENBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQztBQUNwRCxDQUFDO0FBRkQsNENBRUMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBfYnJpZGdlIH0gZnJvbSBcIi4vX2JyaWRnZVwiO1xuaW1wb3J0IHsgdXNlRWZmZWN0IH0gZnJvbSBcInJlYWN0XCI7XG5cbmxldCBwb3N0Zml4OiBzdHJpbmcgPSBudWxsO1xuXG4vLyDku44gbm1jQ29uZmlnIOS4reiOt+W+l+W3sue7j+e/u+ivkeeahOWQjue8gFxuY29uc3Qgbm1jQ29uZmlnID0gX2JyaWRnZShcIm5tY0NvbmZpZ1wiKTtcbmlmIChubWNDb25maWcgJiYgdHlwZW9mIG5tY0NvbmZpZy5kZWZhdWx0VGl0bGUgPT09IFwic3RyaW5nXCIpIHtcbiAgcG9zdGZpeCA9IChubWNDb25maWcuZGVmYXVsdFRpdGxlIGFzIHN0cmluZykuc3BsaXQoL1xccyotXFxzKi9nKS5wb3AoKTtcbn1cblxuLyoqXG4gKiDorr7nva7mlofmoaPmoIfpophcbiAqIEBwYXJhbSB0aXRsZSDmlofmoaPmoIfpophcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIHNldERvY3VtZW50VGl0bGUodGl0bGU6IHN0cmluZykge1xuICBzZXRUaW1lb3V0KCgpID0+IHtcbiAgICBkb2N1bWVudC50aXRsZSA9IHBvc3RmaXggPyB0aXRsZSArIFwiIC0gXCIgKyBwb3N0Zml4IDogdGl0bGU7XG4gIH0sIDEpO1xufVxuXG4vKipcbiAqIFJlYWN0IEhvb2tzIOWjsOaYjuaWh+aho+agh+mimFxuICogQHBhcmFtIHRpdGxlIOaWh+aho+agh+mimFxuICovXG5leHBvcnQgZnVuY3Rpb24gdXNlRG9jdW1lbnRUaXRsZSh0aXRsZTogc3RyaW5nKSB7XG4gIHVzZUVmZmVjdCgoKSA9PiBzZXREb2N1bWVudFRpdGxlKHRpdGxlKSwgW3RpdGxlXSk7XG59XG4iXX0=