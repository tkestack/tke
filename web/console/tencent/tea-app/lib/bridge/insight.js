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
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
Object.defineProperty(exports, "__esModule", { value: true });
var insight_1 = require("../core/insight");
var frequencyLimiter_1 = require("../helpers/frequencyLimiter");
var hottagEvent = insight_1.internalInsightShouldNotUsedByBusiness &&
    insight_1.internalInsightShouldNotUsedByBusiness.register("hottag", {
        level: insight_1.internalInsightShouldNotUsedByBusiness.EventLevel.Info,
    });
var hotTagLimiter = new frequencyLimiter_1.FrequencyLimiter({ name: "insight.reportHotTag" });
var statLimiter = new frequencyLimiter_1.FrequencyLimiter({ name: "insight.stat" });
exports.insight = {
    /**
     * 上报点击流
     * @param tag 点击标识
     */
    reportHotTag: function (tag) {
        if (typeof insight_1.internalInsightShouldNotUsedByBusiness !== "object") {
            return;
        }
        if (!tag || !hottagEvent) {
            return;
        }
        hotTagLimiter.exec(function () { return hottagEvent.push({ hottag: String(tag) }, true); });
    },
    /**
     * 自定义上报，如果 ns 和 event 参数不正确，会扔出异常
     *
     * 控制台在预览模式或开发模式不会上报，如果要验证效果，请在手动添加标识到 localStorage:
     *
     * ```js
     * localStorage.debugInsight = true
     * ```
     *
     * 然后刷新页面。
     *
     * @example
     *
      ```js
      import { app } from '@tea/app';
  
      app.insight.stat({
        ns: 'cvm',
        event: 'restart',
        stringFields: {
          instance: 'ins-1d7x8C3s'
        },
        integerFields: {
          cost: 1000
        }
      });
    ```
     */
    stat: function (_a) {
        var ns = _a.ns, event = _a.event, integerFields = _a.integerFields, stringFields = _a.stringFields;
        if (typeof insight_1.internalInsightShouldNotUsedByBusiness !== "object") {
            return;
        }
        if (!ns)
            throw new Error("ns 不能为空");
        if (!/^[a-zA-Z1-9\/]+$/.test(ns)) {
            throw new Error("ns 只允许包括字母、数字、斜杠");
        }
        else if (/^\/|\/$/.test(ns)) {
            throw new Error("ns 使用的斜杠不允许出现在首尾");
        }
        if (!event)
            throw new Error("event 不能为空");
        if (!/^[a-zA-Z1-9\-]+$/.test(event)) {
            throw new Error("event 只允许包括字母、数字、横杠");
        }
        else if (/^-|-$/.test(event)) {
            throw new Error("event 使用的横杠不允许出现在首尾");
        }
        statLimiter.exec(function () {
            var e_1, _a, e_2, _b;
            var stat = [ns, event].join(":");
            var data = {};
            var totalLimit = 20;
            if (integerFields) {
                try {
                    for (var _c = __values(Object.entries(integerFields)), _d = _c.next(); !_d.done; _d = _c.next()) {
                        var _e = __read(_d.value, 2), field = _e[0], value = _e[1];
                        if (totalLimit >= 0 && typeof value === "number" && !isNaN(value)) {
                            totalLimit--;
                            data["int_" + field] = Math.round(value);
                        }
                    }
                }
                catch (e_1_1) { e_1 = { error: e_1_1 }; }
                finally {
                    try {
                        if (_d && !_d.done && (_a = _c.return)) _a.call(_c);
                    }
                    finally { if (e_1) throw e_1.error; }
                }
            }
            if (stringFields) {
                try {
                    for (var _f = __values(Object.entries(stringFields)), _g = _f.next(); !_g.done; _g = _f.next()) {
                        var _h = __read(_g.value, 2), field = _h[0], value = _h[1];
                        if (totalLimit >= 0 && typeof value === "string") {
                            totalLimit--;
                            data["str_" + field] = value.slice(0, 4096);
                        }
                    }
                }
                catch (e_2_1) { e_2 = { error: e_2_1 }; }
                finally {
                    try {
                        if (_g && !_g.done && (_b = _f.return)) _b.call(_f);
                    }
                    finally { if (e_2) throw e_2.error; }
                }
            }
            insight_1.internalInsightShouldNotUsedByBusiness.stat(stat, data);
        });
    },
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5zaWdodC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL3NyYy9icmlkZ2UvaW5zaWdodC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUEsMkNBQXFGO0FBQ3JGLGdFQUErRDtBQUUvRCxJQUFNLFdBQVcsR0FDZixnREFBUTtJQUNSLGdEQUFRLENBQUMsUUFBUSxDQUFDLFFBQVEsRUFBRTtRQUMxQixLQUFLLEVBQUUsZ0RBQVEsQ0FBQyxVQUFVLENBQUMsSUFBSTtLQUNoQyxDQUFDLENBQUM7QUFFTCxJQUFNLGFBQWEsR0FBRyxJQUFJLG1DQUFnQixDQUFDLEVBQUUsSUFBSSxFQUFFLHNCQUFzQixFQUFFLENBQUMsQ0FBQztBQUM3RSxJQUFNLFdBQVcsR0FBRyxJQUFJLG1DQUFnQixDQUFDLEVBQUUsSUFBSSxFQUFFLGNBQWMsRUFBRSxDQUFDLENBQUM7QUFFdEQsUUFBQSxPQUFPLEdBQUc7SUFDckI7OztPQUdHO0lBQ0gsWUFBWSxFQUFaLFVBQWEsR0FBVztRQUN0QixJQUFJLE9BQU8sZ0RBQVEsS0FBSyxRQUFRLEVBQUU7WUFDaEMsT0FBTztTQUNSO1FBQ0QsSUFBSSxDQUFDLEdBQUcsSUFBSSxDQUFDLFdBQVcsRUFBRTtZQUN4QixPQUFPO1NBQ1I7UUFDRCxhQUFhLENBQUMsSUFBSSxDQUFDLGNBQU0sT0FBQSxXQUFXLENBQUMsSUFBSSxDQUFDLEVBQUUsTUFBTSxFQUFFLE1BQU0sQ0FBQyxHQUFHLENBQUMsRUFBRSxFQUFFLElBQUksQ0FBQyxFQUEvQyxDQUErQyxDQUFDLENBQUM7SUFDNUUsQ0FBQztJQUVEOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7T0EyQkc7SUFDSCxJQUFJLEVBQUosVUFBSyxFQUF1RDtZQUFyRCxVQUFFLEVBQUUsZ0JBQUssRUFBRSxnQ0FBYSxFQUFFLDhCQUFZO1FBQzNDLElBQUksT0FBTyxnREFBUSxLQUFLLFFBQVEsRUFBRTtZQUNoQyxPQUFPO1NBQ1I7UUFDRCxJQUFJLENBQUMsRUFBRTtZQUFFLE1BQU0sSUFBSSxLQUFLLENBQUMsU0FBUyxDQUFDLENBQUM7UUFDcEMsSUFBSSxDQUFDLGtCQUFrQixDQUFDLElBQUksQ0FBQyxFQUFFLENBQUMsRUFBRTtZQUNoQyxNQUFNLElBQUksS0FBSyxDQUFDLGtCQUFrQixDQUFDLENBQUM7U0FDckM7YUFBTSxJQUFJLFNBQVMsQ0FBQyxJQUFJLENBQUMsRUFBRSxDQUFDLEVBQUU7WUFDN0IsTUFBTSxJQUFJLEtBQUssQ0FBQyxrQkFBa0IsQ0FBQyxDQUFDO1NBQ3JDO1FBRUQsSUFBSSxDQUFDLEtBQUs7WUFBRSxNQUFNLElBQUksS0FBSyxDQUFDLFlBQVksQ0FBQyxDQUFDO1FBQzFDLElBQUksQ0FBQyxrQkFBa0IsQ0FBQyxJQUFJLENBQUMsS0FBSyxDQUFDLEVBQUU7WUFDbkMsTUFBTSxJQUFJLEtBQUssQ0FBQyxxQkFBcUIsQ0FBQyxDQUFDO1NBQ3hDO2FBQU0sSUFBSSxPQUFPLENBQUMsSUFBSSxDQUFDLEtBQUssQ0FBQyxFQUFFO1lBQzlCLE1BQU0sSUFBSSxLQUFLLENBQUMscUJBQXFCLENBQUMsQ0FBQztTQUN4QztRQUVELFdBQVcsQ0FBQyxJQUFJLENBQUM7O1lBQ2YsSUFBTSxJQUFJLEdBQUcsQ0FBQyxFQUFFLEVBQUUsS0FBSyxDQUFDLENBQUMsSUFBSSxDQUFDLEdBQUcsQ0FBQyxDQUFDO1lBQ25DLElBQU0sSUFBSSxHQUFHLEVBQUUsQ0FBQztZQUVoQixJQUFJLFVBQVUsR0FBRyxFQUFFLENBQUM7WUFDcEIsSUFBSSxhQUFhLEVBQUU7O29CQUNqQixLQUEyQixJQUFBLEtBQUEsU0FBQSxNQUFNLENBQUMsT0FBTyxDQUFDLGFBQWEsQ0FBQyxDQUFBLGdCQUFBLDRCQUFFO3dCQUFqRCxJQUFBLHdCQUFjLEVBQWIsYUFBSyxFQUFFLGFBQUs7d0JBQ3BCLElBQUksVUFBVSxJQUFJLENBQUMsSUFBSSxPQUFPLEtBQUssS0FBSyxRQUFRLElBQUksQ0FBQyxLQUFLLENBQUMsS0FBSyxDQUFDLEVBQUU7NEJBQ2pFLFVBQVUsRUFBRSxDQUFDOzRCQUNiLElBQUksQ0FBQyxTQUFPLEtBQU8sQ0FBQyxHQUFHLElBQUksQ0FBQyxLQUFLLENBQUMsS0FBSyxDQUFDLENBQUM7eUJBQzFDO3FCQUNGOzs7Ozs7Ozs7YUFDRjtZQUNELElBQUksWUFBWSxFQUFFOztvQkFDaEIsS0FBMkIsSUFBQSxLQUFBLFNBQUEsTUFBTSxDQUFDLE9BQU8sQ0FBQyxZQUFZLENBQUMsQ0FBQSxnQkFBQSw0QkFBRTt3QkFBaEQsSUFBQSx3QkFBYyxFQUFiLGFBQUssRUFBRSxhQUFLO3dCQUNwQixJQUFJLFVBQVUsSUFBSSxDQUFDLElBQUksT0FBTyxLQUFLLEtBQUssUUFBUSxFQUFFOzRCQUNoRCxVQUFVLEVBQUUsQ0FBQzs0QkFDYixJQUFJLENBQUMsU0FBTyxLQUFPLENBQUMsR0FBRyxLQUFLLENBQUMsS0FBSyxDQUFDLENBQUMsRUFBRSxJQUFJLENBQUMsQ0FBQzt5QkFDN0M7cUJBQ0Y7Ozs7Ozs7OzthQUNGO1lBRUQsZ0RBQVEsQ0FBQyxJQUFJLENBQUMsSUFBSSxFQUFFLElBQUksQ0FBQyxDQUFDO1FBQzVCLENBQUMsQ0FBQyxDQUFDO0lBQ0wsQ0FBQztDQUNGLENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBpbnRlcm5hbEluc2lnaHRTaG91bGROb3RVc2VkQnlCdXNpbmVzcyBhcyBfaW5zaWdodCB9IGZyb20gXCIuLi9jb3JlL2luc2lnaHRcIjtcbmltcG9ydCB7IEZyZXF1ZW5jeUxpbWl0ZXIgfSBmcm9tIFwiLi4vaGVscGVycy9mcmVxdWVuY3lMaW1pdGVyXCI7XG5cbmNvbnN0IGhvdHRhZ0V2ZW50ID1cbiAgX2luc2lnaHQgJiZcbiAgX2luc2lnaHQucmVnaXN0ZXIoXCJob3R0YWdcIiwge1xuICAgIGxldmVsOiBfaW5zaWdodC5FdmVudExldmVsLkluZm8sXG4gIH0pO1xuXG5jb25zdCBob3RUYWdMaW1pdGVyID0gbmV3IEZyZXF1ZW5jeUxpbWl0ZXIoeyBuYW1lOiBcImluc2lnaHQucmVwb3J0SG90VGFnXCIgfSk7XG5jb25zdCBzdGF0TGltaXRlciA9IG5ldyBGcmVxdWVuY3lMaW1pdGVyKHsgbmFtZTogXCJpbnNpZ2h0LnN0YXRcIiB9KTtcblxuZXhwb3J0IGNvbnN0IGluc2lnaHQgPSB7XG4gIC8qKlxuICAgKiDkuIrmiqXngrnlh7vmtYFcbiAgICogQHBhcmFtIHRhZyDngrnlh7vmoIfor4ZcbiAgICovXG4gIHJlcG9ydEhvdFRhZyh0YWc6IHN0cmluZyk6IHZvaWQge1xuICAgIGlmICh0eXBlb2YgX2luc2lnaHQgIT09IFwib2JqZWN0XCIpIHtcbiAgICAgIHJldHVybjtcbiAgICB9XG4gICAgaWYgKCF0YWcgfHwgIWhvdHRhZ0V2ZW50KSB7XG4gICAgICByZXR1cm47XG4gICAgfVxuICAgIGhvdFRhZ0xpbWl0ZXIuZXhlYygoKSA9PiBob3R0YWdFdmVudC5wdXNoKHsgaG90dGFnOiBTdHJpbmcodGFnKSB9LCB0cnVlKSk7XG4gIH0sXG5cbiAgLyoqXG4gICAqIOiHquWumuS5ieS4iuaKpe+8jOWmguaenCBucyDlkowgZXZlbnQg5Y+C5pWw5LiN5q2j56Gu77yM5Lya5omU5Ye65byC5bi4XG4gICAqXG4gICAqIOaOp+WItuWPsOWcqOmihOiniOaooeW8j+aIluW8gOWPkeaooeW8j+S4jeS8muS4iuaKpe+8jOWmguaenOimgemqjOivgeaViOaenO+8jOivt+WcqOaJi+WKqOa3u+WKoOagh+ivhuWIsCBsb2NhbFN0b3JhZ2U6XG4gICAqXG4gICAqIGBgYGpzXG4gICAqIGxvY2FsU3RvcmFnZS5kZWJ1Z0luc2lnaHQgPSB0cnVlXG4gICAqIGBgYFxuICAgKlxuICAgKiDnhLblkI7liLfmlrDpobXpnaLjgIJcbiAgICpcbiAgICogQGV4YW1wbGVcbiAgICpcbiAgICBgYGBqc1xuICAgIGltcG9ydCB7IGFwcCB9IGZyb20gJ0B0ZWEvYXBwJztcblxuICAgIGFwcC5pbnNpZ2h0LnN0YXQoe1xuICAgICAgbnM6ICdjdm0nLFxuICAgICAgZXZlbnQ6ICdyZXN0YXJ0JyxcbiAgICAgIHN0cmluZ0ZpZWxkczoge1xuICAgICAgICBpbnN0YW5jZTogJ2lucy0xZDd4OEMzcydcbiAgICAgIH0sXG4gICAgICBpbnRlZ2VyRmllbGRzOiB7XG4gICAgICAgIGNvc3Q6IDEwMDBcbiAgICAgIH1cbiAgICB9KTtcbiAgYGBgXG4gICAqL1xuICBzdGF0KHsgbnMsIGV2ZW50LCBpbnRlZ2VyRmllbGRzLCBzdHJpbmdGaWVsZHMgfTogSW5zaWdodFN0YXQpIHtcbiAgICBpZiAodHlwZW9mIF9pbnNpZ2h0ICE9PSBcIm9iamVjdFwiKSB7XG4gICAgICByZXR1cm47XG4gICAgfVxuICAgIGlmICghbnMpIHRocm93IG5ldyBFcnJvcihcIm5zIOS4jeiDveS4uuepulwiKTtcbiAgICBpZiAoIS9eW2EtekEtWjEtOVxcL10rJC8udGVzdChucykpIHtcbiAgICAgIHRocm93IG5ldyBFcnJvcihcIm5zIOWPquWFgeiuuOWMheaLrOWtl+avjeOAgeaVsOWtl+OAgeaWnOadoFwiKTtcbiAgICB9IGVsc2UgaWYgKC9eXFwvfFxcLyQvLnRlc3QobnMpKSB7XG4gICAgICB0aHJvdyBuZXcgRXJyb3IoXCJucyDkvb/nlKjnmoTmlpzmnaDkuI3lhYHorrjlh7rnjrDlnKjpppblsL5cIik7XG4gICAgfVxuXG4gICAgaWYgKCFldmVudCkgdGhyb3cgbmV3IEVycm9yKFwiZXZlbnQg5LiN6IO95Li656m6XCIpO1xuICAgIGlmICghL15bYS16QS1aMS05XFwtXSskLy50ZXN0KGV2ZW50KSkge1xuICAgICAgdGhyb3cgbmV3IEVycm9yKFwiZXZlbnQg5Y+q5YWB6K645YyF5ous5a2X5q+N44CB5pWw5a2X44CB5qiq5p2gXCIpO1xuICAgIH0gZWxzZSBpZiAoL14tfC0kLy50ZXN0KGV2ZW50KSkge1xuICAgICAgdGhyb3cgbmV3IEVycm9yKFwiZXZlbnQg5L2/55So55qE5qiq5p2g5LiN5YWB6K645Ye6546w5Zyo6aaW5bC+XCIpO1xuICAgIH1cblxuICAgIHN0YXRMaW1pdGVyLmV4ZWMoKCkgPT4ge1xuICAgICAgY29uc3Qgc3RhdCA9IFtucywgZXZlbnRdLmpvaW4oXCI6XCIpO1xuICAgICAgY29uc3QgZGF0YSA9IHt9O1xuXG4gICAgICBsZXQgdG90YWxMaW1pdCA9IDIwO1xuICAgICAgaWYgKGludGVnZXJGaWVsZHMpIHtcbiAgICAgICAgZm9yIChsZXQgW2ZpZWxkLCB2YWx1ZV0gb2YgT2JqZWN0LmVudHJpZXMoaW50ZWdlckZpZWxkcykpIHtcbiAgICAgICAgICBpZiAodG90YWxMaW1pdCA+PSAwICYmIHR5cGVvZiB2YWx1ZSA9PT0gXCJudW1iZXJcIiAmJiAhaXNOYU4odmFsdWUpKSB7XG4gICAgICAgICAgICB0b3RhbExpbWl0LS07XG4gICAgICAgICAgICBkYXRhW2BpbnRfJHtmaWVsZH1gXSA9IE1hdGgucm91bmQodmFsdWUpO1xuICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgICAgfVxuICAgICAgaWYgKHN0cmluZ0ZpZWxkcykge1xuICAgICAgICBmb3IgKGxldCBbZmllbGQsIHZhbHVlXSBvZiBPYmplY3QuZW50cmllcyhzdHJpbmdGaWVsZHMpKSB7XG4gICAgICAgICAgaWYgKHRvdGFsTGltaXQgPj0gMCAmJiB0eXBlb2YgdmFsdWUgPT09IFwic3RyaW5nXCIpIHtcbiAgICAgICAgICAgIHRvdGFsTGltaXQtLTtcbiAgICAgICAgICAgIGRhdGFbYHN0cl8ke2ZpZWxkfWBdID0gdmFsdWUuc2xpY2UoMCwgNDA5Nik7XG4gICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICB9XG5cbiAgICAgIF9pbnNpZ2h0LnN0YXQoc3RhdCwgZGF0YSk7XG4gICAgfSk7XG4gIH0sXG59O1xuXG5leHBvcnQgaW50ZXJmYWNlIEluc2lnaHRTdGF0IHtcbiAgLyoqXG4gICAqIOS4uuS6huS4jeS4juWFtuS7luS4muWKoeeahOaVsOaNruWGsueqge+8jOmcgOimgeWItuWumuS4iuaKpeeahOWRveWQjeepuumXtO+8jOWmgiBcImN2bVwiXG4gICAqXG4gICAqIOWRveWQjeepuumXtOWFgeiuuOS9v+eUqOeahOWtl+espu+8muWtl+avjeOAgeaVsOWtl+OAgeaWnOadoCBgL2DvvIjkuI3lhYHorrjlh7rnjrDlnKjpppblsL7vvInvvIzlpoLvvJpcbiAgICpcbiAgICogIC0gXCJjdm0yXCJcbiAgICogIC0gXCJjdm0vc2dcIlxuICAgKi9cbiAgbnM6IHN0cmluZztcblxuICAvKipcbiAgICog5pys5qyh57uf6K6h55qE5LqL5Lu25ZCN56ewXG4gICAqXG4gICAqIOS6i+S7tuWQjeWFgeiuuOS9v+eUqOeahOWtl+espu+8muWtl+avjeOAgeaVsOWtl+OAgeaoquadoCBgLWDvvIjkuI3lhYHorrjlh7rnjrDlnKjpppblsL7vvInvvIzlpoLvvJpcbiAgICpcbiAgICogIC0gXCJyZXN0YXJ0LWZhaWxcIlxuICAgKiAgLSBcIm1vdmUycmVjeWNsZVwiXG4gICAqL1xuICBldmVudDogc3RyaW5nO1xuXG4gIC8qKlxuICAgKiDoh6rlrprkuYnmlbDmja7vvIjmlbTmlbDnsbvlnovvvIlcbiAgICpcbiAgICog5aaC5p6c5Lyg5YWl55qE5piv5LiN5pivIG51bWJlciDnsbvlnovvvIzkvJrooqvlv73nlaXjgILlpoLmnpzkvKDlhaXnmoTkuI3mmK/mlbTmlbDvvIzkvJrkvb/nlKggTWF0aC5yb3VuZCDlj5bmlbTjgIJcbiAgICovXG4gIGludGVnZXJGaWVsZHM/OiBSZWNvcmQ8YW55LCBudW1iZXI+O1xuXG4gIC8qKlxuICAgKiDoh6rlrprkuYnmlbDmja7vvIjlrZfnrKbkuLLnsbvlnovvvIlcbiAgICpcbiAgICog5aaC5p6c5Lyg5YWl55qE5LiN5pivIHN0cmluZyDnsbvlnovvvIzkvJrooqvlv73nlaXjgIJcbiAgICpcbiAgICog5Y2V5Liq5a2X5q6155qE6ZW/5bqm77yM5LiN6IO96LaF6L+HIDQwOTbvvIzlkKbliJnkvJrooqvoo4HliIdcbiAgICovXG4gIHN0cmluZ0ZpZWxkcz86IFJlY29yZDxhbnksIHN0cmluZz47XG59XG4iXX0=