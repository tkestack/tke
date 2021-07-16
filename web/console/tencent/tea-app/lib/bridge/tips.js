"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("./_bridge");
var insight_1 = require("../core/insight");
var warning = require("warning");
var _tips = _bridge_1._bridge("tips");
/**
 * 提供全局用户提示
 */
exports.tips = {
    success: success,
    error: error,
    loading: loading,
};
function success(optsOrMsg, duration) {
    if (optsOrMsg && typeof optsOrMsg === "object") {
        _tips.success(optsOrMsg.message, optsOrMsg.duration);
    }
    else {
        _tips.success(optsOrMsg, duration);
    }
}
function error(optsOrMsg, duration) {
    var message = optsOrMsg;
    if (optsOrMsg && typeof optsOrMsg === "object") {
        message = optsOrMsg.message;
        duration = optsOrMsg.duration;
    }
    var context = insight_1.internalInsightShouldNotUsedByBusiness ? insight_1.internalInsightShouldNotUsedByBusiness.captureTrace(new Error(message)) : null;
    _tips.error(message, duration, context);
}
function loading(optsOrMsg, duration) {
    var message = optsOrMsg;
    if (optsOrMsg && typeof optsOrMsg === "object") {
        message = optsOrMsg.message;
        duration = optsOrMsg.duration;
    }
    _tips.requestStart({ text: message });
    var loadingTimer = null;
    var stop = function () {
        if (loadingTimer) {
            clearTimeout(loadingTimer);
            loadingTimer = null;
            _tips.requestStop();
        }
    };
    loadingTimer = setTimeout(stop, duration || 4000);
    return { stop: stop };
}
/**
 * @deprecated
 * 使用 `app.tips` 代替
 */
exports.tip = {
    /**
     * @deprecated
     * 向用户提示成功信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    success: function (message, duration) {
        warning(false, "`app.tip` is deprecated. Please use `app.tips` instead.");
        _tips.success(message, duration);
    },
    /**
     * @deprecated
     * 向用户提示错误信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    error: function (message, duration) {
        warning(false, "`app.tip` is deprecated. Please use `app.tips` instead.");
        var context = insight_1.internalInsightShouldNotUsedByBusiness ? insight_1.internalInsightShouldNotUsedByBusiness.captureTrace(new Error(message)) : null;
        _tips.error(message, duration, context);
    },
    /**
     * @deprecated
     * 想用户提示加载信息
     * @param message 提示内容
     * @param duration 持续时间
     */
    loading: function (message, duration) {
        warning(false, "`app.tip` is deprecated. Please use `app.tips` instead.");
        _tips.loading(message, duration);
    },
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidGlwcy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL3NyYy9icmlkZ2UvdGlwcy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLHFDQUFvQztBQUNwQywyQ0FBcUY7QUFDckYsaUNBQW9DO0FBRXBDLElBQU0sS0FBSyxHQUFHLGlCQUFPLENBQUMsTUFBTSxDQUFDLENBQUM7QUF5QjlCOztHQUVHO0FBQ1UsUUFBQSxJQUFJLEdBQWtCO0lBQ2pDLE9BQU8sU0FBQTtJQUNQLEtBQUssT0FBQTtJQUNMLE9BQU8sU0FBQTtDQUNSLENBQUM7QUFJRixTQUFTLE9BQU8sQ0FBQyxTQUFjLEVBQUUsUUFBaUI7SUFDaEQsSUFBSSxTQUFTLElBQUksT0FBTyxTQUFTLEtBQUssUUFBUSxFQUFFO1FBQzlDLEtBQUssQ0FBQyxPQUFPLENBQUMsU0FBUyxDQUFDLE9BQU8sRUFBRSxTQUFTLENBQUMsUUFBUSxDQUFDLENBQUM7S0FDdEQ7U0FBTTtRQUNMLEtBQUssQ0FBQyxPQUFPLENBQUMsU0FBUyxFQUFFLFFBQVEsQ0FBQyxDQUFDO0tBQ3BDO0FBQ0gsQ0FBQztBQUlELFNBQVMsS0FBSyxDQUFDLFNBQWMsRUFBRSxRQUFpQjtJQUM5QyxJQUFJLE9BQU8sR0FBRyxTQUFTLENBQUM7SUFDeEIsSUFBSSxTQUFTLElBQUksT0FBTyxTQUFTLEtBQUssUUFBUSxFQUFFO1FBQzlDLE9BQU8sR0FBRyxTQUFTLENBQUMsT0FBTyxDQUFDO1FBQzVCLFFBQVEsR0FBRyxTQUFTLENBQUMsUUFBUSxDQUFDO0tBQy9CO0lBQ0QsSUFBTSxPQUFPLEdBQUcsZ0RBQVEsQ0FBQyxDQUFDLENBQUMsZ0RBQVEsQ0FBQyxZQUFZLENBQUMsSUFBSSxLQUFLLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsSUFBSSxDQUFDO0lBQzVFLEtBQUssQ0FBQyxLQUFLLENBQUMsT0FBTyxFQUFFLFFBQVEsRUFBRSxPQUFPLENBQUMsQ0FBQztBQUMxQyxDQUFDO0FBT0QsU0FBUyxPQUFPLENBQUMsU0FBYyxFQUFFLFFBQWlCO0lBQ2hELElBQUksT0FBTyxHQUFHLFNBQVMsQ0FBQztJQUN4QixJQUFJLFNBQVMsSUFBSSxPQUFPLFNBQVMsS0FBSyxRQUFRLEVBQUU7UUFDOUMsT0FBTyxHQUFHLFNBQVMsQ0FBQyxPQUFPLENBQUM7UUFDNUIsUUFBUSxHQUFHLFNBQVMsQ0FBQyxRQUFRLENBQUM7S0FDL0I7SUFFRCxLQUFLLENBQUMsWUFBWSxDQUFDLEVBQUUsSUFBSSxFQUFFLE9BQU8sRUFBRSxDQUFDLENBQUM7SUFFdEMsSUFBSSxZQUFZLEdBQUcsSUFBSSxDQUFDO0lBRXhCLElBQU0sSUFBSSxHQUFHO1FBQ1gsSUFBSSxZQUFZLEVBQUU7WUFDaEIsWUFBWSxDQUFDLFlBQVksQ0FBQyxDQUFDO1lBQzNCLFlBQVksR0FBRyxJQUFJLENBQUM7WUFDcEIsS0FBSyxDQUFDLFdBQVcsRUFBRSxDQUFDO1NBQ3JCO0lBQ0gsQ0FBQyxDQUFDO0lBRUYsWUFBWSxHQUFHLFVBQVUsQ0FBQyxJQUFJLEVBQUUsUUFBUSxJQUFJLElBQUksQ0FBQyxDQUFDO0lBRWxELE9BQU8sRUFBRSxJQUFJLE1BQUEsRUFBRSxDQUFDO0FBQ2xCLENBQUM7QUFFRDs7O0dBR0c7QUFDVSxRQUFBLEdBQUcsR0FBRztJQUNqQjs7Ozs7T0FLRztJQUNILE9BQU8sRUFBUCxVQUFRLE9BQWUsRUFBRSxRQUFpQjtRQUN4QyxPQUFPLENBQUMsS0FBSyxFQUFFLHlEQUF5RCxDQUFDLENBQUM7UUFDMUUsS0FBSyxDQUFDLE9BQU8sQ0FBQyxPQUFPLEVBQUUsUUFBUSxDQUFDLENBQUM7SUFDbkMsQ0FBQztJQUVEOzs7OztPQUtHO0lBQ0gsS0FBSyxFQUFMLFVBQU0sT0FBZSxFQUFFLFFBQWlCO1FBQ3RDLE9BQU8sQ0FBQyxLQUFLLEVBQUUseURBQXlELENBQUMsQ0FBQztRQUMxRSxJQUFJLE9BQU8sR0FBRyxnREFBUSxDQUFDLENBQUMsQ0FBQyxnREFBUSxDQUFDLFlBQVksQ0FBQyxJQUFJLEtBQUssQ0FBQyxPQUFPLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxJQUFJLENBQUM7UUFDMUUsS0FBSyxDQUFDLEtBQUssQ0FBQyxPQUFPLEVBQUUsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO0lBQzFDLENBQUM7SUFFRDs7Ozs7T0FLRztJQUNILE9BQU8sRUFBUCxVQUFRLE9BQWUsRUFBRSxRQUFpQjtRQUN4QyxPQUFPLENBQUMsS0FBSyxFQUFFLHlEQUF5RCxDQUFDLENBQUM7UUFDMUUsS0FBSyxDQUFDLE9BQU8sQ0FBQyxPQUFPLEVBQUUsUUFBUSxDQUFDLENBQUM7SUFDbkMsQ0FBQztDQUNGLENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBfYnJpZGdlIH0gZnJvbSBcIi4vX2JyaWRnZVwiO1xuaW1wb3J0IHsgaW50ZXJuYWxJbnNpZ2h0U2hvdWxkTm90VXNlZEJ5QnVzaW5lc3MgYXMgX2luc2lnaHQgfSBmcm9tIFwiLi4vY29yZS9pbnNpZ2h0XCI7XG5pbXBvcnQgd2FybmluZyA9IHJlcXVpcmUoXCJ3YXJuaW5nXCIpO1xuXG5jb25zdCBfdGlwcyA9IF9icmlkZ2UoXCJ0aXBzXCIpO1xuXG5leHBvcnQgaW50ZXJmYWNlIEFwcFRpcHNCcmlkZ2Uge1xuICAvKipcbiAgICog5ZCR55So5oi35o+Q56S65oiQ5Yqf5L+h5oGvXG4gICAqIEBwYXJhbSBtZXNzYWdlIOaPkOekuuWGheWuuVxuICAgKiBAcGFyYW0gZHVyYXRpb24g5oyB57ut5pe26Ze0XG4gICAqL1xuICBzdWNjZXNzOiB0eXBlb2Ygc3VjY2VzcztcblxuICAvKipcbiAgICog5ZCR55So5oi35o+Q56S66ZSZ6K+v5L+h5oGvXG4gICAqIEBwYXJhbSBtZXNzYWdlIOaPkOekuuWGheWuuVxuICAgKiBAcGFyYW0gZHVyYXRpb24g5oyB57ut5pe26Ze0XG4gICAqL1xuICBlcnJvcjogdHlwZW9mIGVycm9yO1xuXG4gIC8qKlxuICAgKiDmg7PnlKjmiLfmj5DnpLrliqDovb3kv6Hmga9cbiAgICogQHBhcmFtIG1lc3NhZ2Ug5o+Q56S65YaF5a65XG4gICAqIEBwYXJhbSBkdXJhdGlvbiDmjIHnu63ml7bpl7RcbiAgICovXG4gIGxvYWRpbmc6IHR5cGVvZiBsb2FkaW5nO1xufVxuXG4vKipcbiAqIOaPkOS+m+WFqOWxgOeUqOaIt+aPkOekulxuICovXG5leHBvcnQgY29uc3QgdGlwczogQXBwVGlwc0JyaWRnZSA9IHtcbiAgc3VjY2VzcyxcbiAgZXJyb3IsXG4gIGxvYWRpbmcsXG59O1xuXG5mdW5jdGlvbiBzdWNjZXNzKG9wdGlvbnM6IHsgbWVzc2FnZTogc3RyaW5nOyBkdXJhdGlvbj86IG51bWJlciB9KTogdm9pZDtcbmZ1bmN0aW9uIHN1Y2Nlc3MobWVzc2FnZTogc3RyaW5nLCBkdXJhdGlvbj86IG51bWJlcik6IHZvaWQ7XG5mdW5jdGlvbiBzdWNjZXNzKG9wdHNPck1zZzogYW55LCBkdXJhdGlvbj86IG51bWJlcik6IHZvaWQge1xuICBpZiAob3B0c09yTXNnICYmIHR5cGVvZiBvcHRzT3JNc2cgPT09IFwib2JqZWN0XCIpIHtcbiAgICBfdGlwcy5zdWNjZXNzKG9wdHNPck1zZy5tZXNzYWdlLCBvcHRzT3JNc2cuZHVyYXRpb24pO1xuICB9IGVsc2Uge1xuICAgIF90aXBzLnN1Y2Nlc3Mob3B0c09yTXNnLCBkdXJhdGlvbik7XG4gIH1cbn1cblxuZnVuY3Rpb24gZXJyb3Iob3B0aW9uczogeyBtZXNzYWdlOiBzdHJpbmc7IGR1cmF0aW9uPzogbnVtYmVyIH0pOiB2b2lkO1xuZnVuY3Rpb24gZXJyb3IobWVzc2FnZTogc3RyaW5nLCBkdXJhdGlvbj86IG51bWJlcik6IHZvaWQ7XG5mdW5jdGlvbiBlcnJvcihvcHRzT3JNc2c6IGFueSwgZHVyYXRpb24/OiBudW1iZXIpOiB2b2lkIHtcbiAgbGV0IG1lc3NhZ2UgPSBvcHRzT3JNc2c7XG4gIGlmIChvcHRzT3JNc2cgJiYgdHlwZW9mIG9wdHNPck1zZyA9PT0gXCJvYmplY3RcIikge1xuICAgIG1lc3NhZ2UgPSBvcHRzT3JNc2cubWVzc2FnZTtcbiAgICBkdXJhdGlvbiA9IG9wdHNPck1zZy5kdXJhdGlvbjtcbiAgfVxuICBjb25zdCBjb250ZXh0ID0gX2luc2lnaHQgPyBfaW5zaWdodC5jYXB0dXJlVHJhY2UobmV3IEVycm9yKG1lc3NhZ2UpKSA6IG51bGw7XG4gIF90aXBzLmVycm9yKG1lc3NhZ2UsIGR1cmF0aW9uLCBjb250ZXh0KTtcbn1cblxuZnVuY3Rpb24gbG9hZGluZyhvcHRpb25zOiB7XG4gIG1lc3NhZ2U/OiBzdHJpbmc7XG4gIGR1cmF0aW9uPzogbnVtYmVyO1xufSk6IHsgc3RvcDogKCkgPT4gdm9pZCB9O1xuZnVuY3Rpb24gbG9hZGluZyhtZXNzYWdlPzogc3RyaW5nLCBkdXJhdGlvbj86IG51bWJlcik6IHsgc3RvcDogKCkgPT4gdm9pZCB9O1xuZnVuY3Rpb24gbG9hZGluZyhvcHRzT3JNc2c6IGFueSwgZHVyYXRpb24/OiBudW1iZXIpOiB7IHN0b3A6ICgpID0+IHZvaWQgfSB7XG4gIGxldCBtZXNzYWdlID0gb3B0c09yTXNnO1xuICBpZiAob3B0c09yTXNnICYmIHR5cGVvZiBvcHRzT3JNc2cgPT09IFwib2JqZWN0XCIpIHtcbiAgICBtZXNzYWdlID0gb3B0c09yTXNnLm1lc3NhZ2U7XG4gICAgZHVyYXRpb24gPSBvcHRzT3JNc2cuZHVyYXRpb247XG4gIH1cblxuICBfdGlwcy5yZXF1ZXN0U3RhcnQoeyB0ZXh0OiBtZXNzYWdlIH0pO1xuXG4gIGxldCBsb2FkaW5nVGltZXIgPSBudWxsO1xuXG4gIGNvbnN0IHN0b3AgPSAoKSA9PiB7XG4gICAgaWYgKGxvYWRpbmdUaW1lcikge1xuICAgICAgY2xlYXJUaW1lb3V0KGxvYWRpbmdUaW1lcik7XG4gICAgICBsb2FkaW5nVGltZXIgPSBudWxsO1xuICAgICAgX3RpcHMucmVxdWVzdFN0b3AoKTtcbiAgICB9XG4gIH07XG5cbiAgbG9hZGluZ1RpbWVyID0gc2V0VGltZW91dChzdG9wLCBkdXJhdGlvbiB8fCA0MDAwKTtcblxuICByZXR1cm4geyBzdG9wIH07XG59XG5cbi8qKlxuICogQGRlcHJlY2F0ZWRcbiAqIOS9v+eUqCBgYXBwLnRpcHNgIOS7o+abv1xuICovXG5leHBvcnQgY29uc3QgdGlwID0ge1xuICAvKipcbiAgICogQGRlcHJlY2F0ZWRcbiAgICog5ZCR55So5oi35o+Q56S65oiQ5Yqf5L+h5oGvXG4gICAqIEBwYXJhbSBtZXNzYWdlIOaPkOekuuWGheWuuVxuICAgKiBAcGFyYW0gZHVyYXRpb24g5oyB57ut5pe26Ze0XG4gICAqL1xuICBzdWNjZXNzKG1lc3NhZ2U6IHN0cmluZywgZHVyYXRpb24/OiBudW1iZXIpIHtcbiAgICB3YXJuaW5nKGZhbHNlLCBcImBhcHAudGlwYCBpcyBkZXByZWNhdGVkLiBQbGVhc2UgdXNlIGBhcHAudGlwc2AgaW5zdGVhZC5cIik7XG4gICAgX3RpcHMuc3VjY2VzcyhtZXNzYWdlLCBkdXJhdGlvbik7XG4gIH0sXG5cbiAgLyoqXG4gICAqIEBkZXByZWNhdGVkXG4gICAqIOWQkeeUqOaIt+aPkOekuumUmeivr+S/oeaBr1xuICAgKiBAcGFyYW0gbWVzc2FnZSDmj5DnpLrlhoXlrrlcbiAgICogQHBhcmFtIGR1cmF0aW9uIOaMgee7reaXtumXtFxuICAgKi9cbiAgZXJyb3IobWVzc2FnZTogc3RyaW5nLCBkdXJhdGlvbj86IG51bWJlcikge1xuICAgIHdhcm5pbmcoZmFsc2UsIFwiYGFwcC50aXBgIGlzIGRlcHJlY2F0ZWQuIFBsZWFzZSB1c2UgYGFwcC50aXBzYCBpbnN0ZWFkLlwiKTtcbiAgICBsZXQgY29udGV4dCA9IF9pbnNpZ2h0ID8gX2luc2lnaHQuY2FwdHVyZVRyYWNlKG5ldyBFcnJvcihtZXNzYWdlKSkgOiBudWxsO1xuICAgIF90aXBzLmVycm9yKG1lc3NhZ2UsIGR1cmF0aW9uLCBjb250ZXh0KTtcbiAgfSxcblxuICAvKipcbiAgICogQGRlcHJlY2F0ZWRcbiAgICog5oOz55So5oi35o+Q56S65Yqg6L295L+h5oGvXG4gICAqIEBwYXJhbSBtZXNzYWdlIOaPkOekuuWGheWuuVxuICAgKiBAcGFyYW0gZHVyYXRpb24g5oyB57ut5pe26Ze0XG4gICAqL1xuICBsb2FkaW5nKG1lc3NhZ2U6IHN0cmluZywgZHVyYXRpb24/OiBudW1iZXIpIHtcbiAgICB3YXJuaW5nKGZhbHNlLCBcImBhcHAudGlwYCBpcyBkZXByZWNhdGVkLiBQbGVhc2UgdXNlIGBhcHAudGlwc2AgaW5zdGVhZC5cIik7XG4gICAgX3RpcHMubG9hZGluZyhtZXNzYWdlLCBkdXJhdGlvbik7XG4gIH0sXG59O1xuIl19