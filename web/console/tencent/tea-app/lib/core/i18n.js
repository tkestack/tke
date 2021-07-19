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
var i18next_1 = require("i18next");
var react_i18next_1 = require("react-i18next");
var hashString = require("hash-string");
// 是否已经初始化
var hasInited = false;
// 当前语言
var lng = window["VERSION"] || "zh";
// 从句子计算哈希一个 key 值，该算法需要和 scanner 保持一致
var hashKey = function (value) {
    return "k_" +
        ("0000" + hashString(value.replace(/\s+/g, "")).toString(36)).slice(-7);
};
/**
 * 国际化包括的信息以及所需的工具
 */
exports.i18n = {
    /**
     * 当前用户的国际化语言，已知语言：
     *  - `zh` 中文
     *  - `en` 英文
     *  - `jp` 日语
     *  - `ko` 韩语
     */
    lng: lng,
    /**
     * 当前用户所在站点
     *  - `1` 表示国内站；
     *  - `2` 表示国际站；
     *
     * @type {1 | 2}
     */
    site: 1,
    /**
     * 注册国家
     */
    country: {
        name: "CN",
        code: "86",
    },
    /**
     * 初始化当前语言的国际化配置
     */
    init: function (_a) {
        var _b;
        var translation = _a.translation;
        if (hasInited) {
            console.warn("你已经初始化过 i18n，请勿重复初始化");
            return;
        }
        hasInited = true;
        i18next_1.default.use(react_i18next_1.reactI18nextModule).init({
            // 使用语言
            lng: lng,
            // 英文版 fallback 到中文版，其它语言 fallback 到中文版
            fallbackLng: lng === "en" ? "zh" : "en",
            // 翻译资源
            resources: (_b = {},
                _b[lng] = { translation: translation },
                _b),
            ns: "translation",
            defaultNS: "translation",
            interpolation: {
                escapeValue: false,
            },
            react: {
                hashTransKey: hashKey,
            },
        });
    },
    /**
     * 标记翻译句子
     * 详细的标记说明，请参考 http://tapd.oa.com/QCloud_2015/markdown_wikis/view/#1010103951008390841
     */
    t: function (sentence, options) {
        var key = hashKey(sentence);
        return i18next_1.default.t(key, __assign({}, (options || {}), { defaultValue: sentence }));
    },
    /**
     * 标记翻译组件
     * 详细的标记说明，请参考 http://tapd.oa.com/QCloud_2015/markdown_wikis/view/#1010103951008390841
     */
    Trans: react_i18next_1.Trans,
};
if (typeof LOGIN_INFO === "object" && LOGIN_INFO) {
    var area = LOGIN_INFO.area, countryCode = LOGIN_INFO.countryCode, countryName = LOGIN_INFO.countryName;
    exports.i18n.site = area;
    Object.assign(exports.i18n.country, {
        name: countryName,
        code: countryCode,
    });
}
exports.getI18NInstance = function () { return (hasInited ? i18next_1.default : null); };
/**
 * @internal 国际化容器，内部使用
 */
exports.I18NProvider = react_i18next_1.I18nextProvider;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaTE4bi5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uL3NyYy9jb3JlL2kxOG4udHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7OztBQUFBLG1DQUE4QjtBQUM5QiwrQ0FBMkU7QUFDM0Usd0NBQTBDO0FBRTFDLFVBQVU7QUFDVixJQUFJLFNBQVMsR0FBRyxLQUFLLENBQUM7QUFFdEIsT0FBTztBQUNQLElBQU0sR0FBRyxHQUFXLE1BQU0sQ0FBQyxTQUFTLENBQUMsSUFBSSxJQUFJLENBQUM7QUFFOUMsc0NBQXNDO0FBQ3RDLElBQU0sT0FBTyxHQUFHLFVBQUMsS0FBYTtJQUM1QixPQUFBLElBQUk7UUFDSixDQUFDLE1BQU0sR0FBRyxVQUFVLENBQUMsS0FBSyxDQUFDLE9BQU8sQ0FBQyxNQUFNLEVBQUUsRUFBRSxDQUFDLENBQUMsQ0FBQyxRQUFRLENBQUMsRUFBRSxDQUFDLENBQUMsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUM7QUFEdkUsQ0FDdUUsQ0FBQztBQUUxRTs7R0FFRztBQUNVLFFBQUEsSUFBSSxHQUFHO0lBQ2xCOzs7Ozs7T0FNRztJQUNILEdBQUcsS0FBQTtJQUVIOzs7Ozs7T0FNRztJQUNILElBQUksRUFBRSxDQUFDO0lBRVA7O09BRUc7SUFDSCxPQUFPLEVBQUU7UUFDUCxJQUFJLEVBQUUsSUFBSTtRQUNWLElBQUksRUFBRSxJQUFJO0tBQ1g7SUFFRDs7T0FFRztJQUNILElBQUksRUFBRSxVQUFDLEVBQWdDOztZQUE5Qiw0QkFBVztRQUNsQixJQUFJLFNBQVMsRUFBRTtZQUNiLE9BQU8sQ0FBQyxJQUFJLENBQUMsc0JBQXNCLENBQUMsQ0FBQztZQUNyQyxPQUFPO1NBQ1I7UUFDRCxTQUFTLEdBQUcsSUFBSSxDQUFDO1FBQ2pCLGlCQUFPLENBQUMsR0FBRyxDQUFDLGtDQUFrQixDQUFDLENBQUMsSUFBSSxDQUFDO1lBQ25DLE9BQU87WUFDUCxHQUFHLEtBQUE7WUFFSCx1Q0FBdUM7WUFDdkMsV0FBVyxFQUFFLEdBQUcsS0FBSyxJQUFJLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxDQUFDLENBQUMsSUFBSTtZQUV2QyxPQUFPO1lBQ1AsU0FBUztnQkFDUCxHQUFDLEdBQUcsSUFBRyxFQUFFLFdBQVcsRUFBRSxXQUFXLEVBQUU7bUJBQ3BDO1lBRUQsRUFBRSxFQUFFLGFBQWE7WUFDakIsU0FBUyxFQUFFLGFBQWE7WUFFeEIsYUFBYSxFQUFFO2dCQUNiLFdBQVcsRUFBRSxLQUFLO2FBQ25CO1lBRUQsS0FBSyxFQUFFO2dCQUNMLFlBQVksRUFBRSxPQUFPO2FBQ2Y7U0FDVCxDQUFDLENBQUM7SUFDTCxDQUFDO0lBRUQ7OztPQUdHO0lBQ0gsQ0FBQyxFQUFFLFVBQUMsUUFBZ0IsRUFBRSxPQUFnQztRQUNwRCxJQUFNLEdBQUcsR0FBRyxPQUFPLENBQUMsUUFBUSxDQUFDLENBQUM7UUFDOUIsT0FBTyxpQkFBTyxDQUFDLENBQUMsQ0FBQyxHQUFHLGVBQ2YsQ0FBQyxPQUFPLElBQUksRUFBRSxDQUFDLElBQ2xCLFlBQVksRUFBRSxRQUFRLElBQ1osQ0FBQztJQUNmLENBQUM7SUFFRDs7O09BR0c7SUFDSCxLQUFLLHVCQUFBO0NBQ04sQ0FBQztBQUVGLElBQUksT0FBTyxVQUFVLEtBQUssUUFBUSxJQUFJLFVBQVUsRUFBRTtJQUN4QyxJQUFBLHNCQUFJLEVBQUUsb0NBQVcsRUFBRSxvQ0FBVyxDQUFnQjtJQUV0RCxZQUFJLENBQUMsSUFBSSxHQUFHLElBQUksQ0FBQztJQUNqQixNQUFNLENBQUMsTUFBTSxDQUFDLFlBQUksQ0FBQyxPQUFPLEVBQUU7UUFDMUIsSUFBSSxFQUFFLFdBQVc7UUFDakIsSUFBSSxFQUFFLFdBQVc7S0FDbEIsQ0FBQyxDQUFDO0NBQ0o7QUFFWSxRQUFBLGVBQWUsR0FBRyxjQUFNLE9BQUEsQ0FBQyxTQUFTLENBQUMsQ0FBQyxDQUFDLGlCQUFPLENBQUMsQ0FBQyxDQUFDLElBQUksQ0FBQyxFQUE1QixDQUE0QixDQUFDO0FBRWxFOztHQUVHO0FBQ1UsUUFBQSxZQUFZLEdBQUcsK0JBQWUsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBpMThuZXh0IGZyb20gXCJpMThuZXh0XCI7XG5pbXBvcnQgeyBUcmFucywgSTE4bmV4dFByb3ZpZGVyLCByZWFjdEkxOG5leHRNb2R1bGUgfSBmcm9tIFwicmVhY3QtaTE4bmV4dFwiO1xuaW1wb3J0ICogYXMgaGFzaFN0cmluZyBmcm9tIFwiaGFzaC1zdHJpbmdcIjtcblxuLy8g5piv5ZCm5bey57uP5Yid5aeL5YyWXG5sZXQgaGFzSW5pdGVkID0gZmFsc2U7XG5cbi8vIOW9k+WJjeivreiogFxuY29uc3QgbG5nOiBzdHJpbmcgPSB3aW5kb3dbXCJWRVJTSU9OXCJdIHx8IFwiemhcIjtcblxuLy8g5LuO5Y+l5a2Q6K6h566X5ZOI5biM5LiA5LiqIGtleSDlgLzvvIzor6Xnrpfms5XpnIDopoHlkowgc2Nhbm5lciDkv53mjIHkuIDoh7RcbmNvbnN0IGhhc2hLZXkgPSAodmFsdWU6IHN0cmluZykgPT5cbiAgXCJrX1wiICtcbiAgKFwiMDAwMFwiICsgaGFzaFN0cmluZyh2YWx1ZS5yZXBsYWNlKC9cXHMrL2csIFwiXCIpKS50b1N0cmluZygzNikpLnNsaWNlKC03KTtcblxuLyoqXG4gKiDlm73pmYXljJbljIXmi6znmoTkv6Hmga/ku6Xlj4rmiYDpnIDnmoTlt6XlhbdcbiAqL1xuZXhwb3J0IGNvbnN0IGkxOG4gPSB7XG4gIC8qKlxuICAgKiDlvZPliY3nlKjmiLfnmoTlm73pmYXljJbor63oqIDvvIzlt7Lnn6Xor63oqIDvvJpcbiAgICogIC0gYHpoYCDkuK3mlodcbiAgICogIC0gYGVuYCDoi7HmlodcbiAgICogIC0gYGpwYCDml6Xor61cbiAgICogIC0gYGtvYCDpn6nor61cbiAgICovXG4gIGxuZyxcblxuICAvKipcbiAgICog5b2T5YmN55So5oi35omA5Zyo56uZ54K5XG4gICAqICAtIGAxYCDooajnpLrlm73lhoXnq5nvvJtcbiAgICogIC0gYDJgIOihqOekuuWbvemZheerme+8m1xuICAgKlxuICAgKiBAdHlwZSB7MSB8IDJ9XG4gICAqL1xuICBzaXRlOiAxLFxuXG4gIC8qKlxuICAgKiDms6jlhozlm73lrrZcbiAgICovXG4gIGNvdW50cnk6IHtcbiAgICBuYW1lOiBcIkNOXCIsXG4gICAgY29kZTogXCI4NlwiLFxuICB9LFxuXG4gIC8qKlxuICAgKiDliJ3lp4vljJblvZPliY3or63oqIDnmoTlm73pmYXljJbphY3nva5cbiAgICovXG4gIGluaXQ6ICh7IHRyYW5zbGF0aW9uIH06IEkxOE5Jbml0T3B0aW9ucykgPT4ge1xuICAgIGlmIChoYXNJbml0ZWQpIHtcbiAgICAgIGNvbnNvbGUud2FybihcIuS9oOW3sue7j+WIneWni+WMlui/hyBpMThu77yM6K+35Yu/6YeN5aSN5Yid5aeL5YyWXCIpO1xuICAgICAgcmV0dXJuO1xuICAgIH1cbiAgICBoYXNJbml0ZWQgPSB0cnVlO1xuICAgIGkxOG5leHQudXNlKHJlYWN0STE4bmV4dE1vZHVsZSkuaW5pdCh7XG4gICAgICAvLyDkvb/nlKjor63oqIBcbiAgICAgIGxuZyxcblxuICAgICAgLy8g6Iux5paH54mIIGZhbGxiYWNrIOWIsOS4reaWh+eJiO+8jOWFtuWug+ivreiogCBmYWxsYmFjayDliLDkuK3mlofniYhcbiAgICAgIGZhbGxiYWNrTG5nOiBsbmcgPT09IFwiZW5cIiA/IFwiemhcIiA6IFwiZW5cIixcblxuICAgICAgLy8g57+76K+R6LWE5rqQXG4gICAgICByZXNvdXJjZXM6IHtcbiAgICAgICAgW2xuZ106IHsgdHJhbnNsYXRpb246IHRyYW5zbGF0aW9uIH0sXG4gICAgICB9LFxuXG4gICAgICBuczogXCJ0cmFuc2xhdGlvblwiLFxuICAgICAgZGVmYXVsdE5TOiBcInRyYW5zbGF0aW9uXCIsXG5cbiAgICAgIGludGVycG9sYXRpb246IHtcbiAgICAgICAgZXNjYXBlVmFsdWU6IGZhbHNlLCAvLyBub3QgbmVlZGVkIGZvciByZWFjdCBhcyBpdCBlc2NhcGVzIGJ5IGRlZmF1bHRcbiAgICAgIH0sXG5cbiAgICAgIHJlYWN0OiB7XG4gICAgICAgIGhhc2hUcmFuc0tleTogaGFzaEtleSxcbiAgICAgIH0gYXMgYW55LFxuICAgIH0pO1xuICB9LFxuXG4gIC8qKlxuICAgKiDmoIforrDnv7vor5Hlj6XlrZBcbiAgICog6K+m57uG55qE5qCH6K6w6K+05piO77yM6K+35Y+C6ICDIGh0dHA6Ly90YXBkLm9hLmNvbS9RQ2xvdWRfMjAxNS9tYXJrZG93bl93aWtpcy92aWV3LyMxMDEwMTAzOTUxMDA4MzkwODQxXG4gICAqL1xuICB0OiAoc2VudGVuY2U6IHN0cmluZywgb3B0aW9ucz86IEkxOE5UcmFuc2xhdGlvbk9wdGlvbnMpID0+IHtcbiAgICBjb25zdCBrZXkgPSBoYXNoS2V5KHNlbnRlbmNlKTtcbiAgICByZXR1cm4gaTE4bmV4dC50KGtleSwge1xuICAgICAgLi4uKG9wdGlvbnMgfHwge30pLFxuICAgICAgZGVmYXVsdFZhbHVlOiBzZW50ZW5jZSxcbiAgICB9KSBhcyBzdHJpbmc7XG4gIH0sXG5cbiAgLyoqXG4gICAqIOagh+iusOe/u+ivkee7hOS7tlxuICAgKiDor6bnu4bnmoTmoIforrDor7TmmI7vvIzor7flj4LogIMgaHR0cDovL3RhcGQub2EuY29tL1FDbG91ZF8yMDE1L21hcmtkb3duX3dpa2lzL3ZpZXcvIzEwMTAxMDM5NTEwMDgzOTA4NDFcbiAgICovXG4gIFRyYW5zLFxufTtcblxuaWYgKHR5cGVvZiBMT0dJTl9JTkZPID09PSBcIm9iamVjdFwiICYmIExPR0lOX0lORk8pIHtcbiAgY29uc3QgeyBhcmVhLCBjb3VudHJ5Q29kZSwgY291bnRyeU5hbWUgfSA9IExPR0lOX0lORk87XG5cbiAgaTE4bi5zaXRlID0gYXJlYTtcbiAgT2JqZWN0LmFzc2lnbihpMThuLmNvdW50cnksIHtcbiAgICBuYW1lOiBjb3VudHJ5TmFtZSxcbiAgICBjb2RlOiBjb3VudHJ5Q29kZSxcbiAgfSk7XG59XG5cbmV4cG9ydCBjb25zdCBnZXRJMThOSW5zdGFuY2UgPSAoKSA9PiAoaGFzSW5pdGVkID8gaTE4bmV4dCA6IG51bGwpO1xuXG4vKipcbiAqIEBpbnRlcm5hbCDlm73pmYXljJblrrnlmajvvIzlhoXpg6jkvb/nlKhcbiAqL1xuZXhwb3J0IGNvbnN0IEkxOE5Qcm92aWRlciA9IEkxOG5leHRQcm92aWRlcjtcblxuZXhwb3J0IGludGVyZmFjZSBJMThOSW5pdE9wdGlvbnMge1xuICB0cmFuc2xhdGlvbjogSTE4TlRyYW5zbGF0aW9uO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEkxOE5UcmFuc2xhdGlvbiB7XG4gIFtrZXk6IHN0cmluZ106IHN0cmluZztcbn1cblxuZXhwb3J0IGludGVyZmFjZSBJMThOVHJhbnNsYXRpb25PcHRpb25zIHtcbiAgLyoqIOeUqOS6juehruWumuWNleWkjeaVsOeahOaVsOmHj+WAvCAqL1xuICBjb3VudD86IG51bWJlcjtcblxuICAvKiog55So5LqO56Gu5a6a5LiK5LiL5paH55qE6K+05piO5paH5pys77yM5Y+q6IO95L2/55So5a2X56ym5Liy5bi46YeP77yM5ZCm5YiZ5peg5rOV5omr5o+PICovXG4gIGNvbnRleHQ/OiBzdHJpbmc7XG5cbiAgLy8g5YWB6K645Lyg5YWl5o+S5YC8XG4gIFtrZXk6IHN0cmluZ106IGFueTtcbn1cbiJdfQ==