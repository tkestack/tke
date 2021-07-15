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
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("../_bridge");
/**
 * 批量检查用户白名单
 */
exports.checkWhitelistBatch = function (keys) {
    return new Promise(function (resolve, reject) {
        return _bridge_1._manager.queryWhiteList(
        // params
        { whiteKey: keys }, 
        // onSuccess
        function (result) {
            var e_1, _a;
            var ret = {};
            try {
                for (var keys_1 = __values(keys), keys_1_1 = keys_1.next(); !keys_1_1.done; keys_1_1 = keys_1.next()) {
                    var key = keys_1_1.value;
                    ret[key] = (result && result[key] && +result[key][0]) || 0;
                }
            }
            catch (e_1_1) { e_1 = { error: e_1_1 }; }
            finally {
                try {
                    if (keys_1_1 && !keys_1_1.done && (_a = keys_1.return)) _a.call(keys_1);
                }
                finally { if (e_1) throw e_1.error; }
            }
            resolve(ret);
        }, 
        // onError
        function (error) { return reject(Object.assign(new Error(error.msg), error || {})); });
    });
};
/**
 * 检查白名单
 */
exports.checkWhitelist = function (key) {
    return exports.checkWhitelistBatch([key]).then(function (result) { return result[key]; });
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoid2hpdGVsaXN0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2JyaWRnZS91c2VyL3doaXRlbGlzdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7QUFBQSxzQ0FBc0M7QUFFdEM7O0dBRUc7QUFDVSxRQUFBLG1CQUFtQixHQUFHLFVBQUMsSUFBYztJQUNoRCxPQUFBLElBQUksT0FBTyxDQUE0QixVQUFDLE9BQU8sRUFBRSxNQUFNO1FBQ3JELE9BQUEsa0JBQVEsQ0FBQyxjQUFjO1FBQ3JCLFNBQVM7UUFDVCxFQUFFLFFBQVEsRUFBRSxJQUFJLEVBQUU7UUFDbEIsWUFBWTtRQUNaLFVBQUMsTUFBVzs7WUFDVixJQUFJLEdBQUcsR0FBRyxFQUFFLENBQUM7O2dCQUViLEtBQWdCLElBQUEsU0FBQSxTQUFBLElBQUksQ0FBQSwwQkFBQSw0Q0FBRTtvQkFBakIsSUFBSSxHQUFHLGlCQUFBO29CQUNWLEdBQUcsQ0FBQyxHQUFHLENBQUMsR0FBRyxDQUFDLE1BQU0sSUFBSSxNQUFNLENBQUMsR0FBRyxDQUFDLElBQUksQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsSUFBSSxDQUFDLENBQUM7aUJBQzVEOzs7Ozs7Ozs7WUFFRCxPQUFPLENBQUMsR0FBRyxDQUFDLENBQUM7UUFDZixDQUFDO1FBQ0QsVUFBVTtRQUNWLFVBQUMsS0FBVSxJQUFLLE9BQUEsTUFBTSxDQUFDLE1BQU0sQ0FBQyxNQUFNLENBQUMsSUFBSSxLQUFLLENBQUMsS0FBSyxDQUFDLEdBQUcsQ0FBQyxFQUFFLEtBQUssSUFBSSxFQUFFLENBQUMsQ0FBQyxFQUF4RCxDQUF3RCxDQUN6RTtJQWZELENBZUMsQ0FDRjtBQWpCRCxDQWlCQyxDQUFDO0FBRUo7O0dBRUc7QUFDVSxRQUFBLGNBQWMsR0FBRyxVQUFDLEdBQVc7SUFDeEMsT0FBQSwyQkFBbUIsQ0FBQyxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUMsSUFBSSxDQUFDLFVBQUEsTUFBTSxJQUFJLE9BQUEsTUFBTSxDQUFDLEdBQUcsQ0FBQyxFQUFYLENBQVcsQ0FBQztBQUF0RCxDQUFzRCxDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgX21hbmFnZXIgfSBmcm9tIFwiLi4vX2JyaWRnZVwiO1xuXG4vKipcbiAqIOaJuemHj+ajgOafpeeUqOaIt+eZveWQjeWNlVxuICovXG5leHBvcnQgY29uc3QgY2hlY2tXaGl0ZWxpc3RCYXRjaCA9IChrZXlzOiBzdHJpbmdbXSkgPT5cbiAgbmV3IFByb21pc2U8eyBba2V5OiBzdHJpbmddOiBudW1iZXIgfT4oKHJlc29sdmUsIHJlamVjdCkgPT5cbiAgICBfbWFuYWdlci5xdWVyeVdoaXRlTGlzdChcbiAgICAgIC8vIHBhcmFtc1xuICAgICAgeyB3aGl0ZUtleToga2V5cyB9LFxuICAgICAgLy8gb25TdWNjZXNzXG4gICAgICAocmVzdWx0OiBhbnkpID0+IHtcbiAgICAgICAgbGV0IHJldCA9IHt9O1xuXG4gICAgICAgIGZvciAobGV0IGtleSBvZiBrZXlzKSB7XG4gICAgICAgICAgcmV0W2tleV0gPSAocmVzdWx0ICYmIHJlc3VsdFtrZXldICYmICtyZXN1bHRba2V5XVswXSkgfHwgMDtcbiAgICAgICAgfVxuXG4gICAgICAgIHJlc29sdmUocmV0KTtcbiAgICAgIH0sXG4gICAgICAvLyBvbkVycm9yXG4gICAgICAoZXJyb3I6IGFueSkgPT4gcmVqZWN0KE9iamVjdC5hc3NpZ24obmV3IEVycm9yKGVycm9yLm1zZyksIGVycm9yIHx8IHt9KSlcbiAgICApXG4gICk7XG5cbi8qKlxuICog5qOA5p+l55m95ZCN5Y2VXG4gKi9cbmV4cG9ydCBjb25zdCBjaGVja1doaXRlbGlzdCA9IChrZXk6IHN0cmluZykgPT5cbiAgY2hlY2tXaGl0ZWxpc3RCYXRjaChba2V5XSkudGhlbihyZXN1bHQgPT4gcmVzdWx0W2tleV0pO1xuIl19