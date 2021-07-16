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
var _bridge_1 = require("./_bridge");
exports.capi = {
    /**
     * 云 API 请求
     * @param body 云 API 请求参数
     * @param options 云 API 请求选项
     */
    request: function (body, options) {
        return _bridge_1._bridge("models/api").request(body, __assign({ secure: true, global: (options || {}).tipLoading }, (options || {})));
    },
    /**
     * 云 API V3 请求
     * @param body 云 API 请求参数
     * @param options 云 API 请求选项
     */
    requestV3: function (body, options) {
        return _bridge_1._bridge("models/api").request(body, __assign({}, (options || {}), { global: (options || {}).tipLoading, secure: true, version: 3 }));
    },
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY2FwaS5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL3NyYy9icmlkZ2UvY2FwaS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7O0FBQUEscUNBQW9DO0FBRXZCLFFBQUEsSUFBSSxHQUFHO0lBQ2xCOzs7O09BSUc7SUFDSCxPQUFPLEVBQVAsVUFBUSxJQUFpQixFQUFFLE9BQXdCO1FBQ2pELE9BQU8saUJBQU8sQ0FBQyxZQUFZLENBQUMsQ0FBQyxPQUFPLENBQUMsSUFBSSxhQUN2QyxNQUFNLEVBQUUsSUFBSSxFQUNaLE1BQU0sRUFBRSxDQUFDLE9BQU8sSUFBSSxFQUFFLENBQUMsQ0FBQyxVQUFVLElBQy9CLENBQUMsT0FBTyxJQUFJLEVBQUUsQ0FBQyxFQUNsQixDQUFDO0lBQ0wsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxTQUFTLEVBQVQsVUFBVSxJQUFpQixFQUFFLE9BQTBCO1FBQ3JELE9BQU8saUJBQU8sQ0FBQyxZQUFZLENBQUMsQ0FBQyxPQUFPLENBQUMsSUFBSSxlQUNwQyxDQUFDLE9BQU8sSUFBSSxFQUFFLENBQUMsSUFDbEIsTUFBTSxFQUFFLENBQUMsT0FBTyxJQUFJLEVBQUUsQ0FBQyxDQUFDLFVBQVUsRUFDbEMsTUFBTSxFQUFFLElBQUksRUFDWixPQUFPLEVBQUUsQ0FBQyxJQUNWLENBQUM7SUFDTCxDQUFDO0NBQ0YsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCB7IF9icmlkZ2UgfSBmcm9tIFwiLi9fYnJpZGdlXCI7XG5cbmV4cG9ydCBjb25zdCBjYXBpID0ge1xuICAvKipcbiAgICog5LqRIEFQSSDor7fmsYJcbiAgICogQHBhcmFtIGJvZHkg5LqRIEFQSSDor7fmsYLlj4LmlbBcbiAgICogQHBhcmFtIG9wdGlvbnMg5LqRIEFQSSDor7fmsYLpgInpoblcbiAgICovXG4gIHJlcXVlc3QoYm9keTogUmVxdWVzdEJvZHksIG9wdGlvbnM/OiBSZXF1ZXN0T3B0aW9ucyk6IFByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIF9icmlkZ2UoXCJtb2RlbHMvYXBpXCIpLnJlcXVlc3QoYm9keSwge1xuICAgICAgc2VjdXJlOiB0cnVlLFxuICAgICAgZ2xvYmFsOiAob3B0aW9ucyB8fCB7fSkudGlwTG9hZGluZyxcbiAgICAgIC4uLihvcHRpb25zIHx8IHt9KSxcbiAgICB9KTtcbiAgfSxcblxuICAvKipcbiAgICog5LqRIEFQSSBWMyDor7fmsYJcbiAgICogQHBhcmFtIGJvZHkg5LqRIEFQSSDor7fmsYLlj4LmlbBcbiAgICogQHBhcmFtIG9wdGlvbnMg5LqRIEFQSSDor7fmsYLpgInpoblcbiAgICovXG4gIHJlcXVlc3RWMyhib2R5OiBSZXF1ZXN0Qm9keSwgb3B0aW9ucz86IFJlcXVlc3RWM09wdGlvbnMpOiBQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBfYnJpZGdlKFwibW9kZWxzL2FwaVwiKS5yZXF1ZXN0KGJvZHksIHtcbiAgICAgIC4uLihvcHRpb25zIHx8IHt9KSxcbiAgICAgIGdsb2JhbDogKG9wdGlvbnMgfHwge30pLnRpcExvYWRpbmcsXG4gICAgICBzZWN1cmU6IHRydWUsXG4gICAgICB2ZXJzaW9uOiAzLFxuICAgIH0pO1xuICB9LFxufTtcblxuZXhwb3J0IGludGVyZmFjZSBSZXF1ZXN0Qm9keSB7XG4gIC8qKlxuICAgKiDor7fmsYLnmoTkupEgQVBJIOWcsOWfn1xuICAgKi9cbiAgcmVnaW9uSWQ6IG51bWJlcjtcblxuICAvKipcbiAgICog6K+35rGC55qE5LqRIEFQSSDkuJrliqFcbiAgICovXG4gIHNlcnZpY2VUeXBlOiBzdHJpbmc7XG5cbiAgLyoqXG4gICAqIOivt+axgueahOS6kSBBUEkg5ZCN56ewXG4gICAqL1xuICBjbWQ6IHN0cmluZztcblxuICAvKipcbiAgICog6K+35rGC55qE5LqRIEFQSSDmlbDmja5cbiAgICovXG4gIGRhdGE/OiBhbnk7XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgUmVxdWVzdE9wdGlvbnMge1xuICAvKipcbiAgICog5piv5ZCm5L2/55So5a6J5YWo55qE5Li05pe25a+G6ZKlIEFQSSDmlrnmoYjvvIzlu7rorq7kvb/nlKggdHJ1ZVxuICAgKiBAZGVmYXVsdCB0cnVlXG4gICAqL1xuICBzZWN1cmU/OiBib29sZWFuO1xuXG4gIC8qKlxuICAgKiDkvb/nlKjnmoTkupEgQVBJIOeJiOacrO+8jOivpeWPguaVsOmFjeWQiCBzZWN1cmUg5L2/55SoXG4gICAqXG4gICAqICAgLSBgc2VjdXJlID09IGZhbHNlYO+8jOivpeWPguaVsOaXoOaEj+S5iVxuICAgKiAgIC0gYHNlY3VyZSA9PSB0cnVlICYmIHZlcnNpb24gPSAxYO+8jOWwhuS9v+eUqOiAgeeahOS4tOaXtuWvhumSpeacjeWKoei/m+ihjOWvhumSpeeUs+ivt++8jOWQpuWImeS9v+eUqOaWsOeahOWvhumSpeacjeWKoVxuICAgKiAgIC0gYHNlY3VyZSA9PSB0cnVlICYmIHZlcnNpb24gPSAzYO+8jOS9v+eUqOS6kSBBUEkgdjMg5Z+f5ZCN6K+35rGC77yM5LiN5ZCM5Zyw5Z+f5Z+f5ZCN5LiN5ZCMXG4gICAqL1xuICB2ZXJzaW9uPzogbnVtYmVyO1xuXG4gIC8qKlxuICAgKiDmmK/lkKblsIblrqLmiLfnq68gSVAg6ZmE5Yqg5Zyo5LqRIEFQSSDnmoQgYGNsaWVudElQYCDlj4LmlbDkuK1cbiAgICovXG4gIHdpdGhDbGllbnRJUD86IGJvb2xlYW47XG5cbiAgLyoqXG4gICAqIOaYr+WQpuWwhuWuouaIt+erryBVQSDpmYTliqDlnKjkupEgQVBJIOeahCBgY2xpZW50VUFgIOWPguaVsOS4rVxuICAgKi9cbiAgd2l0aENsaWVudFVBPzogYm9vbGVhbjtcblxuICAvKipcbiAgICog5piv5ZCm5Zyo6aG26YOo5pi+56S65o6l5Y+j6ZSZ6K+vXG4gICAqIOm7mOiupOS4uiB0cnVl77yM5Lya5o+Q56S65LqRIEFQSSDosIPnlKjplJnor6/kv6Hmga/vvIzlpoLmnpzoh6rlt7HlpITnkIblvILluLjvvIzor7forr7nva7or6XphY3nva7kuLogZmFsc2VcbiAgICogQGRlZmF1bHQgdHJ1ZVxuICAgKi9cbiAgdGlwRXJyPzogYm9vbGVhbjtcblxuICAvKipcbiAgICog6K+35rGC5pe25piv5ZCm5Zyo6aG26YOo5pi+56S6IExvYWRpbmcg5o+Q56S6XG4gICAqIEBkZWZhdWx0IHRydWVcbiAgICovXG4gIHRpcExvYWRpbmc/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIFJlcXVlc3RWM09wdGlvbnNcbiAgZXh0ZW5kcyBQaWNrPFxuICAgIFJlcXVlc3RPcHRpb25zLFxuICAgIEV4Y2x1ZGU8a2V5b2YgUmVxdWVzdE9wdGlvbnMsIFwic2VjdXJlXCIgfCBcInZlcnNpb25cIj5cbiAgPiB7fVxuIl19