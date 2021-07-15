"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("./_bridge");
/**
 * 提供 SDK 注册和加载接口
 */
exports.sdk = {
    /**
     * 注册 SDK
     *
     * @param sdkName SDK 名称
     * @param sdkFactory SDK 工厂方法，应该返回 SDK 提供的 API
     */
    register: function (sdkName, sdkFactory) {
        window.define("sdk/" + sdkName, sdkFactory);
    },
    /**
     * 加载并使用指定的 SDK
     *
     * @param sdkName SDK 名称
     */
    use: function (sdkName) { return _bridge_1._sdk.use(sdkName); },
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoic2RrLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vc3JjL2JyaWRnZS9zZGsudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSxxQ0FBaUM7QUFFakM7O0dBRUc7QUFDVSxRQUFBLEdBQUcsR0FBRztJQUNqQjs7Ozs7T0FLRztJQUNILFFBQVEsRUFBRSxVQUFDLE9BQWUsRUFBRSxVQUFxQjtRQUMvQyxNQUFNLENBQUMsTUFBTSxDQUFDLFNBQU8sT0FBUyxFQUFFLFVBQVUsQ0FBQyxDQUFDO0lBQzlDLENBQUM7SUFFRDs7OztPQUlHO0lBQ0gsR0FBRyxFQUFFLFVBQVUsT0FBZSxJQUFLLE9BQUEsY0FBSSxDQUFDLEdBQUcsQ0FBQyxPQUFPLENBQWUsRUFBL0IsQ0FBK0I7Q0FDbkUsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IF9zZGsgfSBmcm9tIFwiLi9fYnJpZGdlXCI7XG5cbi8qKlxuICog5o+Q5L6bIFNESyDms6jlhozlkozliqDovb3mjqXlj6NcbiAqL1xuZXhwb3J0IGNvbnN0IHNkayA9IHtcbiAgLyoqXG4gICAqIOazqOWGjCBTREtcbiAgICpcbiAgICogQHBhcmFtIHNka05hbWUgU0RLIOWQjeensFxuICAgKiBAcGFyYW0gc2RrRmFjdG9yeSBTREsg5bel5Y6C5pa55rOV77yM5bqU6K+l6L+U5ZueIFNESyDmj5DkvpvnmoQgQVBJXG4gICAqL1xuICByZWdpc3RlcjogKHNka05hbWU6IHN0cmluZywgc2RrRmFjdG9yeTogKCkgPT4gYW55KSA9PiB7XG4gICAgd2luZG93LmRlZmluZShgc2RrLyR7c2RrTmFtZX1gLCBzZGtGYWN0b3J5KTtcbiAgfSxcblxuICAvKipcbiAgICog5Yqg6L295bm25L2/55So5oyH5a6a55qEIFNES1xuICAgKlxuICAgKiBAcGFyYW0gc2RrTmFtZSBTREsg5ZCN56ewXG4gICAqL1xuICB1c2U6IDxUID0gYW55PihzZGtOYW1lOiBzdHJpbmcpID0+IF9zZGsudXNlKHNka05hbWUpIGFzIFByb21pc2U8VD4sXG59O1xuIl19