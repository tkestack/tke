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
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var _bridge_1 = require("./_bridge");
var _mfa = _bridge_1._bridge("widget/mfa/mfa");
/**
 * MFA 二次验证业务接入指南
 * http://tapd.oa.com/10103951/markdown_wikis/view/#1010103951008387339
 */
exports.mfa = {
    /**
   * 对指定的云 API 进行 MFA 验证
   *
   * @param api 要校验的云 API，需要包括业务和接口名两部分，如 "cvm:DestroyInstance"
   *
   * @returns 返回值为 boolean 的 Promise，为 true 则表示校验通过，可以调用云 API
   *
   * @example
  ```js
  // 发起 MFA 校验
  const mfaPassed = await app.mfa.verify('cvm:DestroyInstance');
  
  if (!mfaPassed) {
    // 校验取消，跳过后续业务
    return;
  }
  
  // 校验完成，调用云 API
  const result = await app.capi.request({
    serviceType: 'cvm',
    cmd: 'DestroyInstance',
    // ...
  });
  
  ```
  */
    verify: function (api) {
        return __awaiter(this, void 0, void 0, function () {
            var err_1;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        _a.trys.push([0, 2, , 3]);
                        return [4 /*yield*/, _mfa.verify({ api: api })];
                    case 1:
                        _a.sent();
                        return [2 /*return*/, true];
                    case 2:
                        err_1 = _a.sent();
                        return [2 /*return*/, false];
                    case 3: return [2 /*return*/];
                }
            });
        });
    },
    /**
     * 跟 verify 类似，不过为指定的 ownerUin 校验
     */
    verifyForOwner: function (api, ownerUin) {
        return __awaiter(this, void 0, void 0, function () {
            var err_2;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        _a.trys.push([0, 2, , 3]);
                        return [4 /*yield*/, _mfa.verify({ api: api, verifyOwnerUin: ownerUin })];
                    case 1:
                        _a.sent();
                        return [2 /*return*/, true];
                    case 2:
                        err_2 = _a.sent();
                        return [2 /*return*/, false];
                    case 3: return [2 /*return*/];
                }
            });
        });
    },
    /**
   * 校验 MFA 后调用云 API，使用该方法具备失败重新校验的能力
   *
   * @example
  ```js
  const result = await app.mfa.request({
    regionId: 1,
    serviceType: 'cvm',
    cmd: 'DestroyInstance',
    data: {
      instanceId: 'ins-a5d3ccw8c'
    }
  }, {
    onMFAError: error => {
      // 碰到 MFA 的错误请进行重试逻辑，业务可以自己限制重试次数
      return error.retry();
    }
  });
  ```
    *
    * > 注意：*如果已经使用了 `app.mfa.verify()` 方法进行 MFA 校验，则无需再使用该方法发起 API 请求，直接使用 `app.capi.request()` 模块发起即可*
    */
    request: function (body, options) {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                return [2 /*return*/, _mfa.apiRequest(__assign({}, body, options))];
            });
        });
    },
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoibWZhLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vc3JjL2JyaWRnZS9tZmEudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUEscUNBQW9DO0FBR3BDLElBQU0sSUFBSSxHQUFHLGlCQUFPLENBQUMsZ0JBQWdCLENBQUMsQ0FBQztBQUV2Qzs7O0dBR0c7QUFDVSxRQUFBLEdBQUcsR0FBRztJQUNqQjs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztJQXlCQTtJQUNNLE1BQU0sRUFBWixVQUFhLEdBQVc7Ozs7Ozs7d0JBRXBCLHFCQUFNLElBQUksQ0FBQyxNQUFNLENBQUMsRUFBRSxHQUFHLEtBQUEsRUFBRSxDQUFDLEVBQUE7O3dCQUExQixTQUEwQixDQUFDO3dCQUMzQixzQkFBTyxJQUFJLEVBQUM7Ozt3QkFFWixzQkFBTyxLQUFLLEVBQUM7Ozs7O0tBRWhCO0lBRUQ7O09BRUc7SUFDRyxjQUFjLEVBQXBCLFVBQXFCLEdBQVcsRUFBRSxRQUFnQjs7Ozs7Ozt3QkFFOUMscUJBQU0sSUFBSSxDQUFDLE1BQU0sQ0FBQyxFQUFFLEdBQUcsS0FBQSxFQUFFLGNBQWMsRUFBRSxRQUFRLEVBQUUsQ0FBQyxFQUFBOzt3QkFBcEQsU0FBb0QsQ0FBQzt3QkFDckQsc0JBQU8sSUFBSSxFQUFDOzs7d0JBRVosc0JBQU8sS0FBSyxFQUFDOzs7OztLQUVoQjtJQUVEOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7TUFxQkU7SUFDSSxPQUFPLEVBQWIsVUFBYyxJQUFpQixFQUFFLE9BQTBCOzs7Z0JBQ3pELHNCQUFPLElBQUksQ0FBQyxVQUFVLGNBQ2pCLElBQUksRUFDSixPQUFPLEVBQ1YsRUFBQzs7O0tBQ0o7Q0FHRixDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgX2JyaWRnZSB9IGZyb20gXCIuL19icmlkZ2VcIjtcbmltcG9ydCB7IFJlcXVlc3RCb2R5LCBSZXF1ZXN0T3B0aW9ucyB9IGZyb20gXCIuL2NhcGlcIjtcblxuY29uc3QgX21mYSA9IF9icmlkZ2UoXCJ3aWRnZXQvbWZhL21mYVwiKTtcblxuLyoqXG4gKiBNRkEg5LqM5qyh6aqM6K+B5Lia5Yqh5o6l5YWl5oyH5Y2XXG4gKiBodHRwOi8vdGFwZC5vYS5jb20vMTAxMDM5NTEvbWFya2Rvd25fd2lraXMvdmlldy8jMTAxMDEwMzk1MTAwODM4NzMzOVxuICovXG5leHBvcnQgY29uc3QgbWZhID0ge1xuICAvKipcbiAqIOWvueaMh+WumueahOS6kSBBUEkg6L+b6KGMIE1GQSDpqozor4FcbiAqIFxuICogQHBhcmFtIGFwaSDopoHmoKHpqoznmoTkupEgQVBJ77yM6ZyA6KaB5YyF5ous5Lia5Yqh5ZKM5o6l5Y+j5ZCN5Lik6YOo5YiG77yM5aaCIFwiY3ZtOkRlc3Ryb3lJbnN0YW5jZVwiXG4gKiBcbiAqIEByZXR1cm5zIOi/lOWbnuWAvOS4uiBib29sZWFuIOeahCBQcm9taXNl77yM5Li6IHRydWUg5YiZ6KGo56S65qCh6aqM6YCa6L+H77yM5Y+v5Lul6LCD55So5LqRIEFQSVxuICogXG4gKiBAZXhhbXBsZVxuYGBganNcbi8vIOWPkei1tyBNRkEg5qCh6aqMXG5jb25zdCBtZmFQYXNzZWQgPSBhd2FpdCBhcHAubWZhLnZlcmlmeSgnY3ZtOkRlc3Ryb3lJbnN0YW5jZScpO1xuXG5pZiAoIW1mYVBhc3NlZCkge1xuICAvLyDmoKHpqozlj5bmtojvvIzot7Pov4flkI7nu63kuJrliqFcbiAgcmV0dXJuO1xufVxuXG4vLyDmoKHpqozlrozmiJDvvIzosIPnlKjkupEgQVBJXG5jb25zdCByZXN1bHQgPSBhd2FpdCBhcHAuY2FwaS5yZXF1ZXN0KHtcbiAgc2VydmljZVR5cGU6ICdjdm0nLFxuICBjbWQ6ICdEZXN0cm95SW5zdGFuY2UnLFxuICAvLyAuLi5cbn0pO1xuXG5gYGBcbiovXG4gIGFzeW5jIHZlcmlmeShhcGk6IHN0cmluZykge1xuICAgIHRyeSB7XG4gICAgICBhd2FpdCBfbWZhLnZlcmlmeSh7IGFwaSB9KTtcbiAgICAgIHJldHVybiB0cnVlO1xuICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbiAgfSxcblxuICAvKipcbiAgICog6LefIHZlcmlmeSDnsbvkvLzvvIzkuI3ov4fkuLrmjIflrprnmoQgb3duZXJVaW4g5qCh6aqMXG4gICAqL1xuICBhc3luYyB2ZXJpZnlGb3JPd25lcihhcGk6IHN0cmluZywgb3duZXJVaW46IHN0cmluZykge1xuICAgIHRyeSB7XG4gICAgICBhd2FpdCBfbWZhLnZlcmlmeSh7IGFwaSwgdmVyaWZ5T3duZXJVaW46IG93bmVyVWluIH0pO1xuICAgICAgcmV0dXJuIHRydWU7XG4gICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuICB9LFxuXG4gIC8qKlxuICog5qCh6aqMIE1GQSDlkI7osIPnlKjkupEgQVBJ77yM5L2/55So6K+l5pa55rOV5YW35aSH5aSx6LSl6YeN5paw5qCh6aqM55qE6IO95YqbXG4gKiBcbiAqIEBleGFtcGxlXG5gYGBqc1xuY29uc3QgcmVzdWx0ID0gYXdhaXQgYXBwLm1mYS5yZXF1ZXN0KHtcbiAgcmVnaW9uSWQ6IDEsXG4gIHNlcnZpY2VUeXBlOiAnY3ZtJyxcbiAgY21kOiAnRGVzdHJveUluc3RhbmNlJyxcbiAgZGF0YToge1xuICAgIGluc3RhbmNlSWQ6ICdpbnMtYTVkM2NjdzhjJ1xuICB9XG59LCB7XG4gIG9uTUZBRXJyb3I6IGVycm9yID0+IHtcbiAgICAvLyDnorDliLAgTUZBIOeahOmUmeivr+ivt+i/m+ihjOmHjeivlemAu+i+ke+8jOS4muWKoeWPr+S7peiHquW3semZkOWItumHjeivleasoeaVsFxuICAgIHJldHVybiBlcnJvci5yZXRyeSgpO1xuICB9XG59KTtcbmBgYFxuICAqIFxuICAqID4g5rOo5oSP77yaKuWmguaenOW3sue7j+S9v+eUqOS6hiBgYXBwLm1mYS52ZXJpZnkoKWAg5pa55rOV6L+b6KGMIE1GQSDmoKHpqozvvIzliJnml6DpnIDlho3kvb/nlKjor6Xmlrnms5Xlj5HotbcgQVBJIOivt+axgu+8jOebtOaOpeS9v+eUqCBgYXBwLmNhcGkucmVxdWVzdCgpYCDmqKHlnZflj5HotbfljbPlj68qXG4gICovXG4gIGFzeW5jIHJlcXVlc3QoYm9keTogUmVxdWVzdEJvZHksIG9wdGlvbnM6IE1GQVJlcXVlc3RPcHRpb25zKTogUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gX21mYS5hcGlSZXF1ZXN0KHtcbiAgICAgIC4uLmJvZHksXG4gICAgICAuLi5vcHRpb25zLFxuICAgIH0pO1xuICB9LFxuXG4gIC8vIGVuZCBvZiBtZmFcbn07XG5cbmV4cG9ydCBpbnRlcmZhY2UgTUZBUmVxdWVzdE9wdGlvbnMgZXh0ZW5kcyBSZXF1ZXN0T3B0aW9ucyB7XG4gIG9uTUZBRXJyb3IoZXJyb3I6IE1GQUVycm9yKTogUHJvbWlzZTxhbnk+O1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIE1GQUVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICByZXRyeSgpOiBQcm9taXNlPGFueT47XG59XG4iXX0=