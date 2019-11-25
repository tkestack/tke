define('config/adapter', function (require, exports, module) {

    //后台返回
    let cacheKey = 'cachedI18NMap' + window.VERSION;

    if (window.VERSION != 'zh') {

        // 清空旧缓存
        if (localStorage.getItem('cachedI18NMap')) {
            localStorage.setItem('cachedI18NMap', null);
        }

        if (typeof window.I18N_MAP !== 'object') {
            // 如果因为缓存原因没有输出，尝试从 localStorage 获取上次记录的词条
            let cachedMap = localStorage.getItem(cacheKey);
            if (cachedMap) {
                try {
                    window.I18N_MAP = JSON.parse(cachedMap);
                } catch (e) {}
            }
        } else {
            // 保存当前获取到的词条用于容灾
            localStorage.setItem(cacheKey, JSON.stringify(window.I18N_MAP));
        }
    }

    let cgiReturnMap = window.I18N_MAP || {};

    function deepMap(argv) {
        let type = Object.prototype.toString.call(argv);
        if (typeof argv === 'object') {
            for (let key in argv) {
                argv[key] = deepMap(argv[key]);
            }
        } else if (type === '[object String]') {
            return cgiReturnMap[argv] || argv;
        }
        return argv;
    }

    module.exports = function (obj) {
        if (window.VERSION != 'zh') {
            return deepMap(obj);
        } else {
            return obj;
        }
    };
});

define('config/config', function (require, exports, module) {
    let $ = require('$');
    let pageManager = require('pageManager');
    let router = require('router');
    let util = require('util');
    let appUtil = require('appUtil');
    let tips = require('tips');
    let submenuConfig = require('config/submenu');
    let iframeRouterModule = require('config/iframeRouter');
    let jsRouterModule = require('config/jsRouter');
    let cssConfig = require('config/css');
    let adapter = require('config/adapter');
    let reactSupport = require('config/reactSupport');
    let nmcConfig = require('nmcConfig');

    //重定向协议
    let redirectProtocol = function (tgtProtocol, fragment) {
        if (!router.option['html5Mode']) {
            let hashStr = fragment == '/' ? '' : '/#';
            tgtUrl = tgtProtocol + '//' + location.host + hashStr + router.debug + fragment;
        } else {
            tgtUrl = tgtProtocol + '//' + location.host + router.debug + fragment;
        }
        // window.location.href = tgtUrl;
    };

    //页面参数配置
    let config = {

        //首页模块名
        'root': 'home',

        //导航选中匹配
        'navContainer': ['#container .aside'],

        //上报参数配置
        'reporterOptions': {},

        //无需上报的错误码
        'ignoreErrCode': [4, 7, 9, 12, 28, 29, 40, 41, 42, 43,
            230, 231, 309, 506, 509, 510,
            808, 819, 833, 834, 842,
            900, 901, 902, 903, 904, 905, 906, 907, 908, 909,
            910, 911, 912, 913, 914, 915, 916, 917, 918
        ],

        //扩展路由
        'extendRoutes': {},

        //渲染内容前回调
        'beforeContentRender': function () {
            // 国际版路由检查
            if (appUtil.isI18n()) {
                return pageManager.checkRoutesI18n();
            }
            return true;
        },

        //请求错误回调
        'reqErrorHandler': function (ret) {

            appUtil.restoreDialogBtn();

            let isCodeIn = function ( /*...codes*/ ) {
                let retCode = ret.code;
                let retMcCode = ret.mccode;
                let codes = [].slice.call(arguments);
                while (codes.length) {
                    let code = codes.shift();
                    if (retCode == code ||
                        // 100 以下的才可能是 mc 保留错误码
                        code < 100 && retMcCode == code
                    ) {
                        return true;
                    }
                }
                return false;
            };

            // 登录态失效
            if (isCodeIn(7, 9, 21, 49, 50)) {
                login.logout(true);
                login.show();
                return false;
            }
            // 未知逻辑
            else if (isCodeIn(22, 23, 24)) {
                return false;
            }
        },

        //页面切换全局销毁
        'globalDestroy': function () {
            appUtil.setDefaultPid(-1);
        },

        //追加的url请求参数
        'additionalUrlParam': function (param) {
            let addParam = {};
            if (param.regionId == undefined) {
                let regionId = appUtil.getRegionId();
                addParam = {
                    'regionId': regionId
                };
            }
            return addParam;
        },
        //导航选中态配置
        'navActiveConfig': [{
            'container': '.qc-nav-tool-left',
            'className': 'qc-nav-text-select',
            'selectParent': false,
            'not': ['.all-product-link-item', '#userInfo a']
        }, {
            'container': '.menu-list-all',
            'className': 'qc-nav-text-select',
            'selectParent': false
        }, {
            'container': '#sidebar',
            'className': 'qc-aside-select',
            'dynamic': true
        }],


        //子菜单配置
        'subMenu': {},

        //css配置
        'css': cssConfig,
        'adapter': adapter,

        //导航前回调函数(目前用于检测https支持)
        'beforeNavigate': function (originalUrl, doNavigate) {
            let _beforeNavigate = function () {
                if (!config.startupRoutered) {
                    config.startupRoutered = true;
                    startupRouter();
                }

                if (router.debug === '/debug' || router.debug === '/debug_http') {
                    return doNavigate();
                }
                let fragment = originalUrl.replace(router.debug, '');
                let protocol = window.location.protocol;
                let protocolReg = new RegExp('^' + protocol);
                let iframeRouter = iframeRouterModule.getCurConfig();
                let jsRouter = jsRouterModule.getCurConfig();

                // 国际版路由检查（部分业务不经过 beforeContentRender）
                if (appUtil.isI18n()) {
                    if (util.cookie.get('language') != window.VERSION) {
                        util.cookie.set('language', window.VERSION);
                    }
                }

                let _checkHttpsLogic = function () {
                    let url = '';

                    for (var key in iframeRouter) {
                        if (new RegExp('^\/' + key).test(fragment) && iframeRouter[key]) {
                            url = iframeRouter[key].content;
                            break;
                        }
                    }

                    for (var key in jsRouter) {
                        if (new RegExp('^\/' + key).test(fragment) && jsRouter[key]) {
                            let jsSrc = jsRouter[key].content;
                            let tempJsUrl = /^http/.test(jsSrc) ? jsSrc : 'https://' + jsSrc;
                            url = tempJsUrl;
                            break;
                        }
                    }

                    if (!url) {
                        url = 'https://';
                    }

                    doNavigate();
                };

                do {
                    if (!fragment || fragment == '/') { //首页允许访问
                        break;
                    }
                } while (false);
                _checkHttpsLogic();
            };

            if (window.VERSION == 'zh') {
                _beforeNavigate();
            } else {
                appUtil.setI18n(true);
                _beforeNavigate();
            }
        }
    };

    var startupRouter = function () {
        let iframeRouter = iframeRouterModule.getCurConfig();
        let jsRouter = jsRouterModule.getCurConfig();

        //iframe接入业务路由配置
        let iframeRouteConfig = {
            suffixRule: '(/*action)(/*p1)(/*p2)(/*p3)(/*p4)',
            map: iframeRouter
        };

        config['iframeRouteCopy'] = {};

        let appendQuery = function (url, query, value) {
            // 如果带 hash，则参数添加到 hash 后面
            let queryConnecter = url.split('#').pop().indexOf('?') > 0 ? '&' : '?';
            return url + queryConnecter + query + '=' + encodeURIComponent(value);
        };

        // 配置 iframe 业务
        let setupIframeBusinessLoader = function (businessKey) {
            let routeKey = '/' + businessKey + iframeRouteConfig.suffixRule;
            config['iframeRouteCopy'][routeKey] =
                config['extendRoutes'][routeKey] = function () {

                    let iframeRouteUrl = iframeRouteConfig.map[businessKey].content;

                    // 如果base url以#结尾，说明业务期望spa
                    let isSpaIframe = iframeRouteConfig.map[businessKey].component === 'spaIframe';

                    // 记录 Hash 后移除
                    let hashMatch = iframeRouteUrl.match(/#(.*)$/);
                    let hashStr = hashMatch ? hashMatch[0] : '';
                    iframeRouteUrl = iframeRouteUrl.replace(hashStr, '');

                    // 匹配子路由
                    let reg = new RegExp('/' + businessKey + '(.*)');
                    let urlMatch = reg.exec(location.href);
                    let urlFragment = '';
                    if (urlMatch) {
                        urlFragment = urlMatch[1];
                        if (!/^\//.test(urlFragment) && !/\.(php|html)$/.test(iframeRouteUrl)) {
                            urlFragment = '/' + urlFragment;
                        }
                    }

                    let url = '';

                    // 拼接最终要加载的 URL
                    // spa iframe将url fragment拼最后#之后
                    if (isSpaIframe) {
                        url = iframeRouteUrl;
                        url = appendQuery(url, 'from', 'new_mc');
                        url += '#' + urlFragment;
                    } else {
                        url = iframeRouteUrl + urlFragment;
                        url = appendQuery(url, 'from', 'new_mc');
                        url += hashStr;
                    }

                    let loadIFrameView = function () {
                        pageManager.loadView('iframe', null, [businessKey, url, {
                            isSpaIframe: isSpaIframe
                        }]);
                    };

                    //费用中心比较特殊，临时 hack
                    if (businessKey === 'account') {
                        let accountUrlSplit = url.split('\.php')[1].split('?'),
                            accountUrlFrag = accountUrlSplit[0].replace('/', ''),
                            accountUrlSearch = accountUrlSplit[1];
                        url = iframeRouteUrl + '?' + accountUrlSearch;
                        accountUrlFrag = (accountUrlFrag == 'projectbill') ? 'bill' : accountUrlFrag;
                        accountUrlFrag && (url += '#act=' + accountUrlFrag);
                    }

                    loadIFrameView();
                };
        };

        for (let businessKey in iframeRouteConfig.map) {
            setupIframeBusinessLoader(businessKey);
        }

        //js模块接入业务路由配置
        let jsModuleRouteConfig = {
            suffixRule: '(/*action)(/*p1)(/*p2)(/*p3)(/*p4)',
            map: {}
        };

        //排序 jsRouter, 将 /path/abc 这种规则排在 /path 前面以做到精确匹配
        Object.keys(jsRouter).sort(function (item1, item2) {
            return item1 < item2 ? 1 : -1;
        }).forEach(function (key) {
            jsModuleRouteConfig.map[key] = jsRouter[key];
        });

        let beeRoutes = {
            lbbm: 1,
            vpcbm: 1,
            cpm: 1,
            monitor: 1,
            mq: 1,
            'mq/topic': 1
        };


        for (let key in jsModuleRouteConfig.map) {
            (function (key) {
                let whiteList,
                    name;

                let jsObj = jsModuleRouteConfig.map[key];
                let isReactApp = reactSupport.isReactAppSource(jsObj);
                let isTeaApp = reactSupport.isTeaApp(jsObj);

                name = key.split('|')[0].trim();

                whiteList = jsObj.whitekey;

                if (name === '~') {
                    name = 'home';
                }

                config['extendRoutes']['/' + name + jsModuleRouteConfig.suffixRule] = function () {
                    let para = name.split('/');
                    para = para.concat([].slice.call(arguments));

                    let jsSrc = jsObj.content;
                    let isExternalSource = /^http/.test(jsSrc) || true;
                    let isBeeRouteMode = name in beeRoutes;

                    let loadFun = isExternalSource ? require.async : nmcLoad;

                    let needReactSupport = isExternalSource && isReactApp;
                    if (needReactSupport) {
                        loadFun = reactSupport.loadFn;
                    }

                    //判断react 在IE8下显示浏览器升级提示
                    let islteIe8 = navigator.appName === 'Microsoft Internet Explorer' && (/MSIE 8.|MSIE 7.|MSIE 6./i).test(navigator.appVersion);
                    let islteIe9 = navigator.appName === 'Microsoft Internet Explorer' && (/MSIE 9.|MSIE 8.|MSIE 7.|MSIE 6./i).test(navigator.appVersion);
                    let browserNotMatchForReactApp = reactSupport.isLteIE9(jsObj) ? islteIe9 : islteIe8;
                    if (needReactSupport && browserNotMatchForReactApp) {
                        let BrowserOutOfDate = '' +
                            '<p class="tc-15-msg warning" style="margin: 20px">' +
                            '            <span class="tip-info">' +
                            '                <span class="msg-span">' +
                            '                    <i class="n-error-icon"></i>' +
                            '                    <span>您的浏览器版本太低，请使用' +
                            '                        <a href="https://www.google.com/chrome/" target="_blank">Chrome</a>、' +
                            '                        <a href="https://www.mozilla.org/zh-CN/firefox/new/" target="_blank">FireFox</a>' +
                            '                        或者升级' +
                            '                        <a href="http://windows.microsoft.com/zh-cn/internet-explorer/download-ie" target="_blank">Internet Explorer</a>' +
                            '                        访问</span>' +
                            '                </span>' +
                            '            </span>' +
                            '</p>';

                        pageManager.loadCommon.apply(pageManager, para);
                        pageManager.appArea && pageManager.appArea.html(BrowserOutOfDate);
                        return;
                    }

                    let cssRouter = pageManager.getCssRouter();
                    let cssReqArr = pageManager.getCssHref(cssRouter);
                    let isNotInRouterWhiteList = false;

                    tips.requestStart();
                    appUtil.parallel([
                        function (cb) {
                            if (cssReqArr.length) {
                                pageManager.loadPageCss(cssRouter, cb);
                            } else {
                                cb();
                            }
                        },
                        function (cb) {
                            let loadJs = function () {
                                loadFun(jsSrc.replace(/http:|https:/, ''), cb, isTeaApp);
                            };

                            loadJs();
                        }
                    ], function () {
                        tips.requestStop();

                        let shouldGoOn = false;

                        if (name === 'home') {
                            shouldGoOn = location.pathname === '/';
                        } else {
                            shouldGoOn = location.href.indexOf('/' + name) > 0;
                        }

                        if (shouldGoOn) {
                            let lstController = pageManager.curController;
                            if (name.indexOf('/') > -1 && pageManager.fragment) {
                                let fragments = pageManager.fragment.split('/');
                                let len = name.split('/').length;
                                if (fragments.length > len) {
                                    lstController = fragments.slice(1, 1 + len).join('/');
                                }
                            }
                            pageManager.isNotInRouterWhiteList = isNotInRouterWhiteList;
                            pageManager.cleanupBeforeNextRender = isExternalSource && !isBeeRouteMode && !isReactApp || (lstController && lstController != name);
                            pageManager.loadCommon.apply(pageManager, para);
                        }
                    });

                };
            })(key);
        }

        router.extend(config.extendRoutes);

        nmcConfig.iframeRouteCopy = config.iframeRouteCopy;

        nmcConfig.subMenu = submenuConfig.getCurConfig();
    };



    module.exports = config;
});

define('config/constants', function (require, exports, module) {
    let constants = {};
    module.exports = constants;
});

define('config/css', function (require, exports, module) {
    module.exports = window.g_config_data && window.g_config_data.css_config;
});

define('config/iframeRouter', function (require, exports, module) {
    let appUtil = require('appUtil');
    let $ = require('$');

    let iframeRouter = {},
        versions = ['zh', 'en', 'intl', 'ko'],
        intlSites = {
            'intl': 1,
            'ko': 1
        };

    for (let i = 0; i < versions.length; i++) {
        let version = versions[i];


        if (intlSites[version] && version != 'intl') {
            g_config_data['iframeRouter_' + version] = $.extend({}, g_config_data['iframeRouter_intl'], g_config_data['iframeRouter_' + version]);
        }

        iframeRouter[version] = g_config_data['iframeRouter_' + version] || {};
    }

    iframeRouter.getCurConfig = function () {
        let _config = {};

        let _version = 'zh';

        if (appUtil.isI18n()) {
            if (window.VERSION != 'en') {
                _version = window.VERSION;
            } else {
                _version = 'intl';
            }
        } else {
            if (window.VERSION == 'en') {
                _version = 'en';
            }
        }

        _config = this[_version] || this['zh'];

        return _config;
    };

    module.exports = iframeRouter;
});


define('config/jsRouter', function (require, exports, module) {
    let appUtil = require('appUtil');
    let $ = require('$');

    let monitorJsRouter = 'https://imgcache.qq.com/qcloud/monitor/dest/app.js';
    if (window.VERSION != '' && window.VERSION != 'zh') {
        monitorJsRouter = 'https://imgcache.qq.com/qcloud/monitor/dest/en/app.js';
    }

    let jsRouter = {
            monitor: monitorJsRouter //monitor用了这里的变量所以必须声明
        },
        versions = ['zh', 'en', 'intl', 'ko'],
        intlVersions = {
            'intl': 1,
            'ko': 1
        };

    for (let i = 0; i < versions.length; i++) {
        let version = versions[i];

        if (intlVersions[version] && version != 'intl') {
            g_config_data['jsRouter_' + version] = $.extend({}, g_config_data['jsRouter_intl'], g_config_data['jsRouter_' + version]);
        }

        jsRouter[version] = g_config_data['jsRouter_' + version] || {};
    }

    jsRouter.getCurConfig = function () {
        let _config = {};

        let _version = 'zh';

        if (appUtil.isI18n()) {
            if (window.VERSION != 'en') {
                _version = window.VERSION;
            } else {
                _version = 'intl';
            }
        } else {
            if (window.VERSION == 'en') {
                _version = 'en';
            }
        }

        _config = this[_version] || this['zh'];

        return _config;
    };

    module.exports = jsRouter;
});

/**
 * 子业务网络链接预处理
 */

define('config/preload', function (require, exports, module) {
    let cdnDomain = 'cdn.' + document.domain;
    let preloadConfig = {
        cdb: {}, //默认预热地址使用 /favicon.ico?r= + random()
        redis: {},
        tdsql: {},
        cdn: {
            url: 'https://' + cdnDomain + '/https_warmup'
        },
    };

    exports.preLoad = function () {
        preLoad();
    };

    var preLoad = function () {
        let url;
        for (let key in preloadConfig) {
            url = preloadConfig[key].url;
            if (typeof url === 'function') {
                url();
            } else {
                (new Image()).src = url;
            }
        }
    };
});

define('config/reactSupport', function (require, exports) {

    let REACT_BUNDLE = exports.bundles = [
        '/tencenthub/static/js/polyfill.min.js',
        '/tencenthub/static/js/react-with-addons.min.js'
    ];

    let lang = window.VERSION;
    let TEA_APP_BUNDLE_MAP = {
        en: '/tencenthub/static/js/tea.js',
        zh: '/tencenthub/static/js/tea.js',
        ko: '/tencenthub/static/js/tea.js'
    };

    let TEA_APP_BUNDLE = TEA_APP_BUNDLE_MAP[lang] || TEA_APP_BUNDLE_MAP['zh'];

    let loadBundle = function (bundle, cb) {
        require.async(bundle, function () {
            setTimeout(cb, 1);
        });
    };

    exports.isReactAppSource = function (jsSrc) {
        if (typeof jsSrc === 'string') {
            return /#rx-app/.test(jsSrc);
        } else {
            return jsSrc.component == 'react' || jsSrc.component == 'tea';
        }

    };

    exports.isLteIE9 = function (jsSrc) {
        if (typeof jsSrc === 'string') {
            return /#rx-app-lteie9/.test(jsSrc);
        } else {
            return jsSrc.isUpdateLteIE9 == 'ie9';
        }
    };

    exports.isTeaApp = function (jsSrc) {
        if (typeof jsSrc === 'string') {
            return /#rx-app-tea/.test(jsSrc);
        } else {
            return jsSrc.component == 'tea';
        }
    };

    exports.useReact = function (cb) {
        loadBundle(REACT_BUNDLE, cb);
    };

    exports.useTea = function (cb) {
        exports.useReact(function () {
            loadBundle(TEA_APP_BUNDLE, cb);
        });
    };

    exports.loadFn = function (jsSrc, cb, isTeaApp) {
        jsSrc = jsSrc.split('#').shift();

        let use = isTeaApp ? exports.useTea : exports.useReact;

        use(function () {
            require.async(jsSrc, cb);
        });
    };
});

define('config/submenu', function (require, exports, module) {
    let appUtil = require('appUtil');
    let $ = require('$');

    let submenu = {},
        versions = ['zh', 'en', 'intl', 'ko'],
        intlVersions = {
            'intl': 1,
            'ko': 1
        };

    for (let i = 0; i < versions.length; i++) {
        let version = versions[i];

        if (intlVersions[version] && version != 'intl') {
            g_config_data['submenu_' + version] = $.extend({}, g_config_data['submenu_intl'], g_config_data['submenu_' + version]);
        }

        submenu[version] = g_config_data['submenu_' + version] || {};
    }

    submenu.getCurConfig = function () {
        let _config = {};

        let _version = 'zh';

        if (appUtil.isI18n()) {
            if (window.VERSION != 'en') {
                _version = window.VERSION;
            } else {
                _version = 'intl';
            }
        } else {
            if (window.VERSION == 'en') {
                _version = 'en';
            }
        }

        _config = this[_version] || this['zh'];

        return _config;
    };

    module.exports = submenu;

});

define('lib/Class', function (require, exports, module) {
    let extend = require('$').extend;


    module.exports = {
        /**
        * 构造函数继承.
        * @param {Object} [protoProps] 子构造函数的扩展原型对象
        * @param {Object} [staticProps] 子构造函数的扩展静态属性
        * @return {Function} 子构造函数
        * @author justanzhu
        * @example
               var Car = function () { this.speed = 100; }
               Car.prototype.run = function() { alert('running') };
               var BlueCar = Car.extend({color: 'blue', stop: function() { alert('stop'); }});
               new BlueCar();
        */
        extend: function (protoProps, staticProps) {
            protoProps = protoProps || {};
            let constructor = protoProps.hasOwnProperty('constructor') ? protoProps.constructor : function () {
                return sup.apply(this, arguments);
            };
            var sup = this;
            let Fn = function () {
                this.constructor = constructor;
            };

            Fn.prototype = sup.prototype;
            constructor.prototype = new Fn();
            extend(constructor.prototype, protoProps);
            extend(constructor, sup, staticProps, {
                __super__: sup.prototype
            });

            return constructor;
        }
    };
});

define('lib/appUtil', function (require, exports, module) {
    let $ = require('$');
    let util = require('util');
    let constants = require('constants');
    let pageManager = require('pageManager');
    let router = require('router');
    let tips = require('tips');
    let regionId;
    let projectId;
    let managePerm;
    let verified,
        verifiedOverseas,
        _isI18n;
    let comData;
    let defaultPid = -1;

    //es6-promise polyfill
    if (!/^\/debug\//.test(window.location.pathname) || !window.Promise) {
        window.Promise = require('es6-promise');
    }

    if (!Array.prototype.indexOf) {
        Array.prototype.indexOf = function (searchElement, fromIndex) {

            let k;
            if (this == null) {
                throw new TypeError('"this" is null or not defined');
            }

            let O = Object(this);
            let len = O.length >>> 0;
            if (len === 0) {
                return -1;
            }
            let n = +fromIndex || 0;

            if (Math.abs(n) === Infinity) {
                n = 0;
            }

            if (n >= len) {
                return -1;
            }

            k = Math.max(n >= 0 ? n : len - Math.abs(n), 0);

            while (k < len) {
                if (k in O && O[k] === searchElement) {
                    return k;
                }
                k++;
            }
            return -1;
        };
    }

    if (!Array.prototype.forEach) {
        Array.prototype.forEach = function (fn, scope) {
            for (let i = 0, len = this.length; i < len; ++i) {
                if (i in this) {
                    fn.call(scope, this[i], i, this);
                }
            }
        };
    }

    if (!Function.prototype.bind) {
        Function.prototype.bind = function (oThis) {
            if (typeof this !== 'function') {
                // closest thing possible to the ECMAScript 5 internal IsCallable function
                throw new TypeError('Function.prototype.bind - what is trying to be bound is not callable');
            }

            let aArgs = Array.prototype.slice.call(arguments, 1),
                fToBind = this,
                fNOP = function () {},
                fBound = function () {
                    return fToBind.apply(this instanceof fNOP && oThis ?
                        this :
                        oThis,
                        aArgs.concat(Array.prototype.slice.call(arguments)));
                };

            fNOP.prototype = this.prototype;
            fBound.prototype = new fNOP();

            return fBound;
        };
    }

    // 实现 ECMA-262, Edition 5, 15.4.4.19
    // 参考: http://es5.github.com/#x15.4.4.19
    if (!Array.prototype.map) {
        Array.prototype.map = function (callback, thisArg) {

            let T,
                A,
                k;

            if (this == null) {
                throw new TypeError(' this is null or not defined');
            }

            let O = Object(this);

            // 将len赋值为数组O的长度.
            let len = O.length >>> 0;

            if ({}.toString.call(callback) != '[object Function]') {
                throw new TypeError(callback + ' is not a function');
            }

            if (thisArg) {
                T = thisArg;
            }
            A = new Array(len);
            k = 0;
            while (k < len) {
                var kValue,
                    mappedValue;
                if (k in O) {
                    kValue = O[k];
                    mappedValue = callback.call(T, kValue, k, O);
                    A[k] = mappedValue;
                }
                k++;
            }
            return A;
        };
    }

    if (!Array.prototype.filter) {
        Array.prototype.filter = function (fun /*, thisArg */ ) {
            'use strict';

            if (this === void 0 || this === null) {
                throw new TypeError();
            }

            let t = Object(this);
            let len = t.length >>> 0;
            if (typeof fun !== 'function') {
                throw new TypeError();
            }

            let res = [];
            let thisArg = arguments.length >= 2 ? arguments[1] : void 0;
            for (let i = 0; i < len; i++) {
                if (i in t) {
                    let val = t[i];

                    if (fun.call(thisArg, val, i, t)) {
                        res.push(val);
                    }
                }
            }

            return res;
        };
    }

    if (!Date.now) {
        Date.now = function now() {
            return (new Date()).getTime();
        };
    }

    if (!Array.prototype.some) {
        Array.prototype.some = function (fun /*, thisArg */ ) {


            if (this === void 0 || this === null) {
                throw new TypeError();
            }

            let t = Object(this);
            let len = t.length >>> 0;
            if (typeof fun !== 'function') {
                throw new TypeError();
            }

            let thisArg = arguments.length >= 2 ? arguments[1] : void 0;
            for (let i = 0; i < len; i++) {
                if (i in t && fun.call(thisArg, t[i], i, t)) {
                    return true;
                }
            }

            return false;
        };
    }

    if (!Array.prototype.every) {
        Array.prototype.every = function (fun /*, thisArg */ ) {


            if (this === void 0 || this === null) {
                throw new TypeError();
            }

            let t = Object(this);
            let len = t.length >>> 0;
            if (typeof fun !== 'function') {
                throw new TypeError();
            }

            let thisArg = arguments.length >= 2 ? arguments[1] : void 0;
            for (let i = 0; i < len; i++) {
                if (i in t && !fun.call(thisArg, t[i], i, t)) {
                    return false;
                }
            }

            return true;
        };
    }

    if (!Array.prototype.find) {
        Array.prototype.find = function (predicate) {

            if (this == null) {
                throw new TypeError('Array.prototype.find called on null or undefined');
            }
            if (typeof predicate !== 'function') {
                throw new TypeError('predicate must be a function');
            }
            let list = Object(this);
            let length = list.length >>> 0;
            let thisArg = arguments[1];
            let value;

            for (let i = 0; i < length; i++) {
                value = list[i];
                if (predicate.call(thisArg, value, i, list)) {
                    return value;
                }
            }
            return undefined;
        };
    }

    /**
     * 工具类
     * @class appUtil
     * @static
     */
    module.exports = {
        changeLetterCase: function (obj, type, notChangeAllUppercaseFeild) {
            let isObject = function (obj) {
                return obj === Object(obj);
            };
            let _self = this;
            if (Array.isArray(obj)) {
                obj.forEach(function (item, i, arr) {
                    if (isObject(item) || Array.isArray(item)) {
                        arr[i] = _self.changeLetterCase(item, type, notChangeAllUppercaseFeild);
                    }
                });
            } else if (isObject(obj)) {
                for (let prop in obj) {
                    if (obj.hasOwnProperty(prop)) {
                        let newProp = prop;
                        if (
                            !notChangeAllUppercaseFeild ||
                            (notChangeAllUppercaseFeild && !/^[A-Z]+$/g.test(prop))
                        ) {
                            let firstLetter = prop.substr(0, 1);
                            newProp =
                                (type === 'upper' ?
                                    firstLetter.toUpperCase() :
                                    firstLetter.toLowerCase()) + prop.slice(1);
                        }
                        if (newProp !== prop) {
                            obj[newProp] = JSON.parse(JSON.stringify(obj[prop]));
                            delete obj[prop];
                        }
                        if (isObject(obj[newProp]) || Array.isArray(obj[newProp])) {
                            obj[newProp] = _self.changeLetterCase(
                                obj[newProp],
                                type,
                                notChangeAllUppercaseFeild
                            );
                        }
                    }
                }
            }

            return obj;
        },
        toUpperCaseFL: function (obj) {
            return this.changeLetterCase(obj, 'upper');
        },
        toLowerCaseFL: function (obj) {
            return this.changeLetterCase(obj, 'lower', true);
        },
        /**
         * 获取系统时间
         * @method getSystemTime
         * @return {Number} 系统时间
         */
        getSystemTime: function () {
            let systemClientGap = comData.systemClientTimeGap || Number(util.cookie.get('systemTimeGap')) || 0;
            return new Date().getTime() + systemClientGap;
        },

        /**
         * 设置区域Id
         * @method setRegionId
         * @param {Int} rId 区域Id
         */
        setRegionId: function (rId) {
            regionId = rId;
            util.cookie.set('regionId', regionId);
            if (window.localStorage) {
                localStorage[util.getUin() + '_regionId'] = regionId;
            }
            $(document).trigger('regionIdchanged', [regionId]);
        },

        /**
         * 获取区域Id. 业务代码直接使用该方法, 可能会造成请求的地域混乱, 应避免使用
         * @method getRegionId
         * @return {Int} regionId
         */
        getRegionId: function () {
            let dId = 1,
                rId = regionId || this._getUrlRid() || util.cookie.get('regionId');
            if (!rId) {
                if (window.localStorage) {
                    let locRid = localStorage[util.getUin() + '_regionId'];
                    rId = locRid || dId;
                } else {
                    rId = dId;
                }
            }
            return rId;
        },

        /**
         * 获取区域Id字符串V3
         * @method getRegionId
         * @return {String} regionId
         */
        getRegionIdV3: function () {
            let rid = constants.REGIONIDMAP[this.getRegionId()];
            return rid;
        },

        /**
         * 是否海外地域
         * @method isOverseasRegion
         * @return {boolean} 是否海外地域
         */
        isOverseasRegion: function (regionId) {
            regionId = regionId || this.getRegionId();
            return constants.OVERSEAS_REGION.indexOf(regionId) != -1;
        },

        _getUrlRid: function () {
            let rid = '',
                regExp = new RegExp('rid=[^&]*'),
                ridMatch = regExp.exec(location.href),
                detailMatch = /\/cvm.*?\/detail\/(\d{1,2})\//.exec(location.href);

            if (ridMatch) {
                rid = ridMatch[0].split('=')[1];
            } else if (detailMatch) {
                rid = detailMatch[1];
            }
            return rid;
        },

        /**
         * 设置默认项目Id
         * @method setDefaultPid
         */
        setDefaultPid: function (pId) {
            defaultPid = pId;
        },

        /**
         * 设置管理权限
         * @method setManagePerm
         * @param {Bool} perm 权限
         */
        setManagePerm: function (perm) {
            managePerm = perm;
        },

        /**
         * 获取管理权限
         * @method getManagePerm
         * @return {Bool} managePerm
         */
        getManagePerm: function () {
            return true;
        },
        setVerified: function (ver) {
            verified = ver;
        },

        setVerifiedOverseas: function (ver) {
            verifiedOverseas = ver;
        },

        setI18n: function (ver) {
            _isI18n = ver;
        },

        isI18n: function () {
            return !!_isI18n;
        },

        /**
         * 设置公共数据
         * @method setComData
         * @param {Object} data
         */
        setComData: function (data) {
            comData = data;
        },

        /**
         * 获取公共数据
         * @method getComData
         * @return {Object} comData
         */
        getComData: function () {
            return comData;
        },

        updateUserInfo: function () {
            this.setComData(null);
            util.cookie.set('refreshSession', '1');
            util.cookie.del('ownerUin');
            util.cookie.del('regionId');
            util.cookie.del('projectId');
        },

        /**
         * 简易版parallel并行任务函数
         * @method parallel
         * @param {Array} tasks 函数数组
         * @param {Function} callback 回调
         * @author evanyuan
         */
        parallel: function (tasks, callback) {
            let len = tasks.length,
                result = [],
                _callback = function (_result, i) {
                    len--;
                    result[i] = _result;
                    if (len == 0) {
                        callback && callback(result);
                    }
                };
            for (var i = 0, fun; fun = tasks[i]; i++) {
                (function (fun, i) {
                    fun && fun(function (_result) {
                        _callback(_result, i);
                    }, function () {
                        _callback(null, i);
                    });
                })(fun, i);

            }
        },

        //右键菜单是否打开
        isContextMenuShow: function () {
            let $cm = $('#contextMenu');
            return ($cm.length && $cm.is(':visible'));
        },

        /**
         *	@description:获取字符串长度，中文认为占用两个字符长度
         *	@author:starcheng
         */
        size: function (str) {
            return str.replace(/[^\u0000-\u00FF]/gmi, '**').length;
        },

        isValidName: function (name) {
            return /^[\w\-.\u4e00-\u9fa5]{1,25}$/.test(name);
        },
        /**
         * @method disableSelection
         * @param {Object} element html元素
         * @param {boolean} bool 禁止或解除
         * @description 禁止或解除元素的文本选中
         * @author brianlin
         */
        disableSelection: function (element, bool) {
            if (!element) {
                return;
            }
            if (bool) {
                if (typeof element.onselectstart !== 'undefined') {
                    element.onselectstart = function () {
                        return false;
                    };
                } else if (typeof element.style.MozUserSelect !== 'undefined') {
                    element.style.MozUserSelect = 'none';
                } else {
                    element.onmousedown = function () {
                        return false;
                    };
                }
            } else {
                if (typeof element.onselectstart !== 'undefined') {
                    element.onselectstart = null;
                } else if (typeof element.style.MozUserSelect !== 'undefined') {
                    element.style.MozUserSelect = '';
                } else {
                    element.onmousedown = null;
                }
            }
        },

        /**
         * 函数切面. 前面的函数返回值传入 breakCheck 判断, breakCheck 返回值为真时不执行切面补充的函数
         * @param oriFn 原始函数
         * @param fn 切面补充函数
         * @param breakCheck
         * @returns {Function}
         * @author justanzhu
         */
        beforeFn: function (oriFn, fn, breakCheck) {
            return function () {
                let ret = fn.apply(this, arguments);
                if (breakCheck && breakCheck.call(this, ret)) {
                    return ret;
                }
                return oriFn.apply(this, arguments);
            };
        },

        afterFn: function (oriFn, fn, breakCheck) {
            return function () {
                let ret = oriFn.apply(this, arguments);
                if (breakCheck && breakCheck.call(this, ret)) {
                    return ret;
                }
                fn.apply(this, arguments);
                return ret;
            };
        },

        // bee属性绑定
        beeBind: function (bee1, bee2, key1, key2, init) {
            key2 = key2 || key1;
            bee1.$watch(key1, function (val, oldValue) {
                if (val !== oldValue) {
                    bee2.$replace(key2, val);
                }
            }, init !== false);
        },
        // bee双向属性绑定
        beeBindBoth: function (bee1, bee2, key1, key2) {
            key2 = key2 || key1;
            this.beeBind(bee1, bee2, key1, key2);
            this.beeBind(bee2, bee1, key2, key1, false);
        },
        capitalize: function (s) {
            return s.charAt(0).toUpperCase() + s.slice(1);
        },
        toDayStart: function (date) {
            date = new Date(date.getTime());
            date.setHours(0);
            date.setMinutes(0);
            date.setSeconds(0);
            date.setMilliseconds(0);
            return date;
        },
        toDayEnd: function (date) {
            date = new Date(date.getTime());
            date.setHours(23);
            date.setMinutes(59);
            date.setSeconds(59);
            date.setMilliseconds(0);
            return date;
        },

        /**
         * 普通对象深拷贝(不支持递归对象)
         * @param obj
         * @returns {Object}
         * @author justanzhu
         */
        clone: function clone(obj) {
            if (obj == null || typeof (obj) !== 'object') {
                return obj;
            }
            let temp = new obj.constructor();
            for (let key in obj) {
                temp[key] = clone(obj[key]);
            }
            return temp;
        },

        getAppId: function () {
            return util.cookie.get('appid');
        },

        /**
         * 控制执行频率，在delay后会执行，保证在 delay到时能执行
         * @param {Function} fn 要执行的函数
         * @param {Number} delay 毫秒
         */
        throttle: function (fn, delay) {
            let timer,
                lasttime = 0,
                args = null,
                contextthis = null,
                result;

            function later() {
                timer = null;
                lasttime = (new Date()).getTime();
                result = fn.apply(contextthis, args);
                args = contextthis = null;
            }

            return function () {
                args = arguments;
                contextthis = this;
                let now = (new Date()).getTime();
                let remaining = delay - (now - lasttime);
                if (remaining <= 0 || remaining > delay) {
                    if (timer) {
                        clearTimeout(timer);
                        timer = null;
                    }
                    lasttime = now;
                    result = fn.apply(contextthis, args);
                    args = contextthis = null;
                } else if (!timer) {
                    timer = setTimeout(later, remaining);
                }

                return result;
            };
        }
    };
});

define('main/api', function (require, exports, module) {
    let pageManager = require('pageManager');
    let router = require('router');
    let appUtil = require('appUtil');

    // js 接入业务API
    let jsApi = {
        render: function (htmlStr, modName) {
            if (router.fragment.indexOf('/' + modName) === 0 && appUtil.getManagePerm()) {
                pageManager.appArea.html(htmlStr);
            }
        },
    };

    module.exports = {
        init: function () {
            window.nmc = jsApi;
        },
        logoutInner: function () {}
    };

});

define('main/startup', function (require, exports) {
    // 基础设施
    let $ = require('$');
    let api = require('main/api');
    let evt = require('event');
    let util = require('util');
    let entry = require('entry');
    let config = require('config/config');
    let router = require('router');
    let preload = require('config/preload');
    let constants = require('constants');
    let pageManager = require('pageManager');

    let zoomDetector = require('widget/zoomDetector/zoomDetector');

    // 路由和菜单
    let jsRouter = require('config/jsRouter');
    let iframeRouter = require('config/iframeRouter');
    let topSubMenuMod = {};

    // 应用入口函数
    let startup = function () {

        // 初始化 api
        api.init();

        // NMC 初始化
        entry.init(config);

        // 国际版引入
        require('widget/i18n');
    };

    exports.startup = startup;
});

/**
 * 国际站翻译检测
 */
define('widget/i18nCheck', function (require, exports, module) {
    let $ = require('$');
    let manager = require('manager');
    let router = require('router');

    let EXPIRES = 30 * 60 * 1000; //空闲后检查持续时间
    let TIMESPAN = 2 * 60 * 1000; //检查间隔
    let INIT_TIMESPAN = 5 * 1000; //路由改变初始检查间隔

    let lastactivityTime = Date.now();
    let timer;

    let chineseRange = [
        [0x4e00, 0x9fff],
        [0x3400, 0x4dbf],
        [0x20000, 0x2a6df],
        [0x2a700, 0x2b73f],
        [0x2b740, 0x2b81f],
        [0x2b820, 0x2ceaf],
        [0xf900, 0xfaff],
        [0x3300, 0x33ff],
        [0xfe30, 0xfe4f],
        [0xf900, 0xfaff],
        [0x2f800, 0x2fa1f],
    ];

    let appArea = '#container';

    var checker = {
        isChinese: function (str) {
            if (str) {
                let charCode;
                let range;
                for (let i = 0; i < str.length; i++) {
                    charCode = str.codePointAt(i);
                    for (let j = 0; j < chineseRange.length; j++) {
                        range = chineseRange[j];
                        if (charCode >= range[0] && charCode <= range[1]) {
                            return true;
                        }
                    }
                }
            }
        },

        watchNode: function () {
            let self = this;

            if (!window.MutationObserver) return;

            let observer = new MutationObserver(function (mutations) {
                for (let i = 0; i < mutations.length; i++) {
                    let mutation = mutations[i];

                    if (mutation && mutation.target) {
                        self.check(mutation.target);
                    }
                }
            });

            let config = {
                attributes: true,
                childList: true,
                characterData: true
            };
            $('#topAlert').length && observer.observe($('#topAlert')[0], config);
            $('#message_list').length && observer.observe($('#message_list')[0], config);
        },

        start: function () {
            let self = this;

            // 国际版时才检测
            manager.getUserAreaInfo(function (areaInfo) {
                if (areaInfo && areaInfo.area == 2 && window._enableReport) {
                    lastactivityTime = Date.now();

                    if (!self.inited) {
                        self.inited = true;

                        setTimeout(function () {
                            self.watchNode();
                        }, INIT_TIMESPAN);
                    }
                }
            });
        },

        domWalker: function (container) {
            if (!document.createTreeWalker || !NodeFilter || !NodeFilter.SHOW_TEXT) return [];

            // 404页面无需检查内容区域
            if (container.id == 'container' && $('#tt404').length) return [];

            // 谷歌翻译的情况下页面上会有 <html class="[...] translated-ltr">
            if ($('.translated-ltr').length) return [];

            let self = this;
            let node,
                arr = [],
                walk = document.createTreeWalker(container, NodeFilter.SHOW_TEXT, null, false);

            while (node = walk.nextNode()) {
                let nodeText = node.textContent;

                if (nodeText && !/^\s*$/gm.test(nodeText) && self.isChinese(nodeText)) {
                    arr.push({
                        text: nodeText,
                        domPath: self.domPath($(node).parent(), {
                            oneResult: true,
                            scaleType: true
                        })
                    });
                }
            }
            return arr;
        },

        domWalkerForIframe: function (container, doc) {
            if (!doc.createTreeWalker || !NodeFilter || !NodeFilter.SHOW_TEXT) return [];

            // 谷歌翻译的情况下页面上会有 <html class="[...] translated-ltr">
            if ($('.translated-ltr').length) return [];

            let self = this;
            let node,
                arr = [],
                walk = doc.createTreeWalker(container, NodeFilter.SHOW_TEXT, null, false);

            while (node = walk.nextNode()) {
                let nodeText = node.textContent;

                if (nodeText && !/^\s*$/gm.test(nodeText) && self.isChinese(nodeText)) {
                    arr.push({
                        text: nodeText,
                        domPath: self.domPath($(node).parent(), {
                            oneResult: true,
                            scaleType: true
                        })
                    });
                }
            }
            return arr;
        },


        domPath: function (elem, options) {
            let opDefault,
                fullPaths,
                $this;

            // ton push all elements short path as string
            fullPaths = [];

            // default options
            opDefault = {
                tag: true, // get dom tag
                lowerCase: true, // get tag in lower or upper case
                class: true, // get element class
                id: true, // get element id
                body: false, // show body in dom full path
                idBeforeClass: true, // display id before class
                oneResult: false, // get only the first result(as string)
                scaleType: false // if the result contains only one element get it as string and not array
            };

            let ops = $.extend(opDefault, options);

            // get dom path depending on options
            let getDomPath = function (el) {

                let elString,
                    elId,
                    elClass,
                    elIdClass;
                elString = elId = elClass = '';

                // if the tag option is enabled
                if (ops.tag) {
                    // get the tag name in lower or upper case
                    elString = ops.lowerCase ? el.tagName.toLowerCase() : el.tagName.toUpperCase();
                }

                if (ops.id && el.id) {
                    elId = '#' + el.id;
                }

                // concat class names
                if (ops.class && el.className) {
                    elClass = '.' + $.trim(el.className || '').replace(/ /g, '.');
                }

                // to display id before or after class name
                elIdClass = ops.idBeforeClass ? elId.concat(elClass) : elClass.concat(elId);
                elString += elIdClass;

                return elString;

            };

            // if the oneResult option is enabled work only in the first element
            $this = ops.oneResult ? $(elem[0]) : elem;

            // do it for all elements
            $this.each(function () {

                let current,
                    domPathItem,
                    pathArray;
                // to don't confuse the this scopes
                current = $(this);
                pathArray = [];

                //pathify also the first element and not only parents
                domPathItem = getDomPath(current.get(0));
                if (domPathItem === '') {
                    return [];
                }

                pathArray.push(domPathItem);
                // for every parent inside body

                let parents = $this.parents(':not(html)');

                parents.each(function () {
                    if (this.tagName) {
                        domPathItem = getDomPath(this);
                        // if the tag option is disabled, it might return empty
                        if (domPathItem !== '') {
                            pathArray.push(domPathItem);
                        }
                    }

                });

                // reverse array ton contact it to string
                pathArray.reverse();

                // if the body option is disabled, check first if the
                // pathArray is not empty and shift its first element
                if (!ops.body && pathArray.length > 0 && pathArray[0].toLowerCase().search('body') === 0) {
                    pathArray.shift();
                }

                // concat with > only if the pathArray has two elements or more
                fullPaths.push(pathArray.length > 1 ? pathArray.join('>') : pathArray[0]);

            });

            // if the option oneResult or scaleType is enabled return only one/first result as string
            if (ops.oneResult || (fullPaths.length === 1 && ops.scaleType)) {
                fullPaths = fullPaths[0];
            }

            return fullPaths;
        },

        parseUrlPath: function () {
            let path = router.getFragment();
            let pathArr = path.split('/');

            if (pathArr.length < 2) return false;

            let path1 = pathArr[1] || '/';
            let path2 = pathArr[2] || '';
            let path3 = pathArr.length > 3 ? pathArr.slice(3).join('/') : '';

            return {
                path1: path1,
                path2: path2,
                path3: path3
            };
        },

        check: function (container) {
            let self = this;
            let isCheckIframe = false;
            let notTranslateNodes = [];
            let ifrDoc;

            if (!container) {
                container = $(appArea)[0];

                // iframe的情况
                try {
                    if ($('#appIFrame').length && $('#appIFrame')[0].contentWindow) {
                        ifrDoc = $('#appIFrame')[0].contentWindow.document;
                        container = ifrDoc.body;
                        isCheckIframe = true;
                    }
                } catch (e) {
                    container = $(appArea)[0];
                    isCheckIframe = false;
                }
            }

            if (isCheckIframe) {
                notTranslateNodes = this.domWalkerForIframe(container, ifrDoc);
            } else {
                notTranslateNodes = this.domWalker(container);
            }

            let pathObj = this.parseUrlPath();

            if (pathObj && notTranslateNodes && notTranslateNodes.length && manager.reportTranslate) {
                let params = $.extend({}, pathObj, {
                    node: notTranslateNodes
                });

                if (self.lastParams && self.lastParams === JSON.stringify(params)) { // 忽略掉重复内容上报
                    return;
                }

                self.lastParams = JSON.stringify(params);

                manager.reportTranslate(params);
            }
        },

        onRouterChange: function () {
            this.start();
        }
    };

    module.exports = checker;
});

define('widget/devtoolsDetector/devtoolsDetector', function (require, exports, module) {
    let $ = require('$');
    let appUtil = require('appUtil');

    let devtools = {
        open: false,
        orientation: null
    };
    let threshold = 160;

    function check() {
        let widthThreshold = window.outerWidth - window.innerWidth > threshold;
        let heightThreshold = window.outerHeight - window.innerHeight > threshold;
        let orientation = widthThreshold ? 'vertical' : 'horizontal';

        if (!(heightThreshold && widthThreshold) &&
            ((window.Firebug && window.Firebug.chrome && window.Firebug.chrome.isInitialized) || widthThreshold || heightThreshold)) {

            devtools.open = true;
            devtools.orientation = orientation;
        } else {
            devtools.open = false;
            devtools.orientation = null;
        }
    }

    check();

    $(window).off('resize.devtoolsDetector').on('resize.devtoolsDetector', appUtil.throttle(function () {
        check();
    }, 500));

    module.exports = devtools;
});
define('widget/zoomDetector/detectZoom', function (require, exports, module) {
    /**
     * Use devicePixelRatio if supported by the browser
     * @return {Number}
     * @private
     */
    let devicePixelRatio = function () {
        return window.devicePixelRatio || 1;
    };

    /**
     * Fallback function to set default values
     * @return {Object}
     * @private
     */
    let fallback = function () {
        return {
            zoom: 1,
            devicePxPerCssPx: 1
        };
    };
    /**
     * IE 8 and 9: no trick needed!
     * TODO: Test on IE10 and Windows 8 RT
     * @return {Object}
     * @private
     **/
    let ie8 = function () {
        let zoom = Math.round((screen.deviceXDPI / screen.logicalXDPI) * 100) / 100;
        return {
            zoom: zoom,
            devicePxPerCssPx: zoom * devicePixelRatio()
        };
    };

    /**
     * For IE10 we need to change our technique again...
     * thanks https://github.com/stefanvanburen
     * @return {Object}
     * @private
     */
    let ie10 = function () {
        let zoom = Math.round((document.documentElement.offsetHeight / window.innerHeight) * 100) / 100;
        return {
            zoom: zoom,
            devicePxPerCssPx: zoom * devicePixelRatio()
        };
    };

    /**
     * For chrome
     *
     */
    let chrome = function () {
        let zoom = Math.round(((window.outerWidth) / window.innerWidth) * 100) / 100;
        return {
            zoom: zoom,
            devicePxPerCssPx: zoom * devicePixelRatio()
        };
    };

    /**
     * For safari (same as chrome)
     *
     */
    let safari = function () {
        let zoom = Math.round(((document.documentElement.clientWidth) / window.innerWidth) * 100) / 100;
        return {
            zoom: zoom,
            devicePxPerCssPx: zoom * devicePixelRatio()
        };
    };


    /**
     * Mobile WebKit
     * the trick: window.innerWIdth is in CSS pixels, while
     * screen.width and screen.height are in system pixels.
     * And there are no scrollbars to mess up the measurement.
     * @return {Object}
     * @private
     */
    let webkitMobile = function () {
        let deviceWidth = (Math.abs(window.orientation) == 90) ? screen.height : screen.width;
        let zoom = deviceWidth / window.innerWidth;
        return {
            zoom: zoom,
            devicePxPerCssPx: zoom * devicePixelRatio()
        };
    };

    /**
     * Desktop Webkit
     * the trick: an element's clientHeight is in CSS pixels, while you can
     * set its line-height in system pixels using font-size and
     * -webkit-text-size-adjust:none.
     * device-pixel-ratio: http://www.webkit.org/blog/55/high-dpi-web-sites/
     *
     * Previous trick (used before http://trac.webkit.org/changeset/100847):
     * documentElement.scrollWidth is in CSS pixels, while
     * document.width was in system pixels. Note that this is the
     * layout width of the document, which is slightly different from viewport
     * because document width does not include scrollbars and might be wider
     * due to big elements.
     * @return {Object}
     * @private
     */
    let webkit = function () {
        let important = function (str) {
            return str.replace(/;/g, ' !important;');
        };

        let div = document.createElement('div');
        div.innerHTML = '1<br>2<br>3<br>4<br>5<br>6<br>7<br>8<br>9<br>0';
        div.setAttribute('style', important('font: 100px/1em sans-serif; -webkit-text-size-adjust: none; text-size-adjust: none; height: auto; width: 1em; padding: 0; overflow: visible;'));

        // The container exists so that the div will be laid out in its own flow
        // while not impacting the layout, viewport size, or display of the
        // webpage as a whole.
        // Add !important and relevant CSS rule resets
        // so that other rules cannot affect the results.
        let container = document.createElement('div');
        container.setAttribute('style', important('width:0; height:0; overflow:hidden; visibility:hidden; position: absolute;'));
        container.appendChild(div);

        document.body.appendChild(container);
        let zoom = 1000 / div.clientHeight;
        zoom = Math.round(zoom * 100) / 100;
        document.body.removeChild(container);

        return {
            zoom: zoom,
            devicePxPerCssPx: zoom * devicePixelRatio()
        };
    };

    /**
     * no real trick; device-pixel-ratio is the ratio of device dpi / css dpi.
     * (Note that this is a different interpretation than Webkit's device
     * pixel ratio, which is the ratio device dpi / system dpi).
     *
     * Also, for Mozilla, there is no difference between the zoom factor and the device ratio.
     *
     * @return {Object}
     * @private
     */
    let firefox4 = function () {
        let zoom = mediaQueryBinarySearch('min--moz-device-pixel-ratio', '', 0, 10, 20, 0.0001);
        zoom = Math.round(zoom * 100) / 100;
        return {
            zoom: zoom,
            devicePxPerCssPx: zoom
        };
    };

    /**
     * Firefox 18.x
     * Mozilla added support for devicePixelRatio to Firefox 18,
     * but it is affected by the zoom level, so, like in older
     * Firefox we can't tell if we are in zoom mode or in a device
     * with a different pixel ratio
     * @return {Object}
     * @private
     */
    let firefox18 = function () {
        return {
            zoom: firefox4().zoom,
            devicePxPerCssPx: devicePixelRatio()
        };
    };

    /**
     * works starting Opera 11.11
     * the trick: outerWidth is the viewport width including scrollbars in
     * system px, while innerWidth is the viewport width including scrollbars
     * in CSS px
     * @return {Object}
     * @private
     */
    let opera11 = function () {
        let zoom = window.top.outerWidth / window.top.innerWidth;
        zoom = Math.round(zoom * 100) / 100;
        return {
            zoom: zoom,
            devicePxPerCssPx: zoom * devicePixelRatio()
        };
    };

    /**
     * Use a binary search through media queries to find zoom level in Firefox
     * @param property
     * @param unit
     * @param a
     * @param b
     * @param maxIter
     * @param epsilon
     * @return {Number}
     */
    var mediaQueryBinarySearch = function (property, unit, a, b, maxIter, epsilon) {
        let matchMedia;
        let head,
            style,
            div;
        if (window.matchMedia) {
            matchMedia = window.matchMedia;
        } else {
            head = document.getElementsByTagName('head')[0];
            style = document.createElement('style');
            head.appendChild(style);

            div = document.createElement('div');
            div.className = 'mediaQueryBinarySearch';
            div.style.display = 'none';
            document.body.appendChild(div);

            matchMedia = function (query) {
                style.sheet.insertRule('@media ' + query + '{.mediaQueryBinarySearch ' + '{text-decoration: underline} }', 0);
                let matched = getComputedStyle(div, null).textDecoration == 'underline';
                style.sheet.deleteRule(0);
                return {
                    matches: matched
                };
            };
        }
        let ratio = binarySearch(a, b, maxIter);
        if (div) {
            head.removeChild(style);
            document.body.removeChild(div);
        }
        return ratio;

        function binarySearch(a, b, maxIter) {
            let mid = (a + b) / 2;
            if (maxIter <= 0 || b - a < epsilon) {
                return mid;
            }
            let query = '(' + property + ':' + mid + unit + ')';
            if (matchMedia(query).matches) {
                return binarySearch(mid, b, maxIter - 1);
            } else {
                return binarySearch(a, mid, maxIter - 1);
            }
        }
    };

    /**
     * Generate detection function
     * @private
     */
    let detectFunction = (function () {
        let func = fallback;
        //IE8+
        if (!isNaN(screen.logicalXDPI) && !isNaN(screen.systemXDPI)) {
            func = ie8;
        }
        // IE10+ / Touch
        else if (window.navigator.msMaxTouchPoints) {
            func = ie10;
        }
        //chrome
        else if (!!window.chrome && !(!!window.opera || navigator.userAgent.indexOf(' Opera') >= 0)) {
            func = chrome;
        }
        //safari
        else if (Object.prototype.toString.call(window.HTMLElement).indexOf('Constructor') > 0) {
            func = safari;
        }
        //Mobile Webkit
        else if ('orientation' in window && 'webkitRequestAnimationFrame' in window) {
            func = webkitMobile;
        }
        //WebKit
        else if ('webkitRequestAnimationFrame' in window) {
            func = webkit;
        }
        //Opera
        else if (navigator.userAgent.indexOf('Opera') >= 0) {
            func = opera11;
        }
        //Last one is Firefox
        //FF 18.x
        else if (window.devicePixelRatio) {
            func = firefox18;
        }
        //FF 4.0 - 17.x
        else if (firefox4().zoom > 0.001) {
            func = firefox4;
        }

        return func;
    })();


    module.exports = ({

        /**
         * Ratios.zoom shorthand
         * @return {Number} Zoom level
         */
        zoom: function () {
            return detectFunction().zoom;
        },

        /**
         * Ratios.devicePxPerCssPx shorthand
         * @return {Number} devicePxPerCssPx level
         */
        device: function () {
            return detectFunction().devicePxPerCssPx;
        }
    });
});
define('widget/zoomDetector/zoomDetector', function (require, exports, module) {
    let $ = require('$');
    let util = require('util');
    let appUtil = require('appUtil');
    let detectZoom = require('widget/zoomDetector/detectZoom');
    let devtoolsDetector = require('widget/devtoolsDetector/devtoolsDetector');

    let _tpl = {
        main: '<div class="qc-win-scale tip-msg-win J-zoomWarning"><div class="tc-15-msg warning"><div class="tip-info"><i class="remind-icon"></i><div class="msg-span">                <%if(isWindows){%>                    页面缩放比例异常，会影响功能的正常使用，请尝试调整缩放比例为100%<span>(快捷键 ctrl+0)</span>                <%}else{%>                    页面缩放比例异常，会影响功能的正常使用，请尝试调整缩放比例为100%<span>(快捷键 cmd+0)</span>                <%}%></div><div class="qc-win-scale-btn"><a href="javascript:;" class="J-ignoreZoomWarning">不再提醒</a><a href="javascript:;" class="J-closeZoomWarning">关闭</a></div></div></div></div>'
    };

    module.exports = {
        init: function () {
            let self = this;

            return;

            if (util.isMobile() || !(this.isMacintosh() || this.isWindows())) {
                return;
            }

            setTimeout(function () {
                self.check();
                self.bindEvents();
            }, 3000);
        },

        check: function () {
            let zoom = detectZoom.zoom();
            let ignoreZoomWarning = false;
            let isDevTollOpened = devtoolsDetector.open;

            if (window.localStorage) {
                ignoreZoomWarning = window.localStorage['ignoreZoomWarning'];
            }

            if ((zoom < 0.97 || zoom > 1.05) && ignoreZoomWarning != 1 && !isDevTollOpened) {
                this.showZoomWarning();
            } else {
                $('.J-zoomWarning').fadeOut();
            }
        },

        isMacintosh: function () {
            return navigator.platform.indexOf('Mac') > -1;
        },

        isWindows: function () {
            return navigator.platform.indexOf('Win') > -1;
        },

        showZoomWarning: function () {
            let isMacintosh = this.isMacintosh();
            let isWindows = this.isWindows();

            if ($('.J-zoomWarning').is(':visible')) return;

            if ($('.J-zoomWarning').length) {
                $('.J-zoomWarning').fadeIn();
            } else {
                $(document.body).append(util.tmpl(_tpl.main, {
                    isMacintosh: isMacintosh,
                    isWindows: isWindows
                }));
            }

            $('.J-ignoreZoomWarning').off().on('click', function () {
                if (window.localStorage) {
                    window.localStorage['ignoreZoomWarning'] = 1;
                }
                $('.J-zoomWarning').hide();
            });

            $('.J-closeZoomWarning').off().on('click', function () {
                $('.J-zoomWarning').hide();
            });
        },

        bindEvents: function () {
            let self = this;

            $(window).off('resize.zoomDetector').on('resize.zoomDetector', appUtil.throttle(function () {
                self.check();
            }, 200));
        }
    };
});
define('nmc/config/config', function (require, exports, module) {
    let $ = require('$');

    /**
     * 参数配置
     * @class nmcConfig
     * @static
     */
    let config = {

        /**
         * 页面模块基础路径
         * @property basePath
         * @type String
         * @default 'modules/'
         */
        'basePath': 'modules/',

        /**
         * 页面包裹选择器
         * @property pageWrapper
         * @type String
         * @default '#pageWrapper'
         */
        'pageWrapper': '#pageWrapper',

        /**
         * 页面主容器选择器
         * @property container
         * @type String
         * @default '#container'
         */
        'container': '#container',

        /**
         * 左侧导航容器选择器
         * @property sidebar
         * @type String
         * @default '#sidebar'
         */
        'sidebar': '#sidebar',

        /**
         * 右侧内容容器选择器
         * @property appArea
         * @type String
         * @default '#appArea'
         */
        'appArea': '#appArea',

        /**
         * 切换页面需要更改class的容器选择器
         * @property classWrapper
         * @type String
         * @default '#container'
         */
        'classWrapper': '#container',

        /**
         * 切换页面需要保留的class
         * @property defaultClass
         * @type String
         * @default ''
         */
        'defaultClass': 'container',

        /**
         * 默认标题
         * @property defaultTitle
         * @type String
         * @default '腾讯云-控制台'
         */
        'defaultTitle': '腾讯云-控制台',

        /**
         * 应用名称
         * @property appName
         * @type String
         * @default 'appName'
         */
        'appName': '管理中心',

        /**
         * 导航选中态配置
         * @property navContainer
         * @type Array
         * @default []
         */
        'navActiveConfig': [],

        /**
         * 渲染前执行方法
         * @property beforeRender
         * @type Function
         * @default function (controller, action, params) {}
         */
        'beforeRender': function (controller, action, params) {},

        /**
         * 渲染内容前执行方法
         * @property beforeContentRender
         * @type Function
         * @default function (obj) {}
         */
        'beforeContentRender': function (obj) {},

        /**
         * 全局销毁
         * @property globalDestroy
         * @type Function
         * @default function () {}
         */
        globalDestroy: function () {},

        /**
         * 扩展路由，优先于框架路由逻辑
         * @property extendRoutes
         * @type Object
         * @default {}
         */
        'extendRoutes': {},

        /**
         * 上报初始参数
         * @property reporterOptions
         * @type Object
         * @default {'tcssDomain': '', 'tcssPrepath': '', 'rtnCodeDomain': '', speedFlag1: '', speedFlag2: '', 'taDomain': '', 'taSid': ''}
         */
        'reporterOptions': {
            //点击流平台配置
            'tcssDomain': '',
            'tcssPrepath': '',

            //返回码配置
            'rtnCodeDomain': '',

            //测速配置
            'speedFlag1': '',
            'speedFlag2': '',

            //ta配置
            'taDomain': '',
            'taSid': ''
        },

        /**
         * 无需上报的错误码
         * @property ignoreErrCode
         * @type Array
         * @default []
         */
        ignoreErrCode: [],

        /**
         * 改变导航选中态
         * @property changeNavStatus
         * @type Function
         * @default ncm通用方法
         */
        'changeNavStatus': null,

        /**
         * 侧边栏子菜单数据
         * @property subMenu
         * @type Object
         * @default {}
         */
        'subMenu': {},

        /**
         * layout模版
         * @property layout
         * @type Object
         * @default {
        				'default': {
        					'controller': [],
        					'module': 'nmc/layout/default'
        				}
        			}
         */
        'layout': {
            'default': {
                'controller': [],
                'module': 'nmc/layout/default'
            }
        },

        /**
         * 首页模块名
         * @property root
         * @type String
         * @default 'home'
         */
        'root': 'home',

        /**
         * css配置
         * @property css
         * @type Object
         * @default {}
         */
        'css': {},

        /**
         * 404提示
         * @property html404
         * @type String
         * @default '<div id="tt404" class="error-page">'+
        		   ' <h2>404</h2> <p>您访问的页面没有找到！</p></div>'
         */
        'html404': '<div id="tt404" class="error-page">' +
            ' <h2>404</h2> <p>您访问的页面没有找到！</p></div>',

        /**
         * 其他错误提示
         * @property htmlOops
         * @type String
         * @default '<div id="ttOops" class="error-page">' +
        			'<h2>Oops</h2> <p>对不起，加载页面时遇到了错误，请稍候再试！</p></div>'
         */
        'htmlOops': '<div id="ttOops" class="error-page">' +
            '<h2>Oops</h2> <p>对不起，加载页面时遇到了错误，<a style="cursor: pointer" onclick="location.reload()">点击重试</a>！</p></div>',

        'userInfoError': '<div id="ttOops" class="error-page">' +
            '<h2>Oops</h2> <p>对不起，获取用户信息时遇到了错误，<a style="cursor: pointer" onclick="location.reload()">点击重试</a>！</p></div>',

        'loginError': '<div id="ttOops" class="error-page">' +
            '<h2 style="font-size:30px">登录异常</h2> <p>对不起，登录出现错误，请 <a style="cursor: pointer" onclick="location.reload()">点击重试</a>！</p></div>',

        /**
         * 加载文字
         * @property loadingWord
         * @type String
         * @default '正在加载...'
         */
        'loadingWord': '正在加载...',

        /**
         * 对话框默认按钮文字
         * @property dialogBtnTxt
         * @type Object
         * @default {'close': '关闭','submit': '确认','cancel': '取消','tips': '提示'}
         */
        'dialogBtnTxt': {
            'close': '关闭',
            'submit': '确定',
            'cancel': '取消',
            'tips': '提示'
        },

        /**
         * 请求错误默认提示文字
         * @property defaultReqErr
         * @type String
         * @default '连接服务器异常，请稍后再试'
         */
        'defaultReqErr': '连接服务器异常，请稍后再试',

        /**
         * 请求错误回调
         * @property reqErrorHandler
         * @type Function
         * @default null
         */
        'reqErrorHandler': null,

        /**
         * 追加的url请求参数
         * @property additionalUrlParam
         * @type Function
         * @default null
         */
        'additionalUrlParam': null

    };

    module.exports = config;
});

/**
 * 默认layout
 */
define('nmc/layout/default', function (require, exports, module) {
    let $ = require('$');
    let util = require('util');

    let layout = {
        _tpl: {
            inStyle: '@keyframes fade-in{0%{opacity:0}100%{}}@-webkit-keyframes fade-in{0%{opacity:0}100%{}}@keyframes fade-out{0%{}100%{opacity:0}}@-webkit-keyframes fade-out{0%{}100%{opacity:0}}@keyframes modal-in{0%{opacity:0;transform:scale(0.8, 0.8);-moz-transform:scale(0.8, 0.8);-webkit-transform:scale(0.8, 0.8);-ms-transform:scale(0.8, 0.8);}100%{}}@-webkit-keyframes modal-in{0%{opacity:0;transform:scale(0.8, 0.8);-moz-transform:scale(0.8, 0.8);-webkit-transform:scale(0.8, 0.8);-ms-transform:scale(0.8, 0.8);}100%{}}@keyframes modal-out{0%{}100%{transform:scale(0.8, 0.8);-moz-transform:scale(0.8, 0.8);-webkit-transform:scale(0.8, 0.8);-ms-transform:scale(0.8, 0.8);opacity:0}}@-webkit-keyframes modal-out{0%{}100%{transform:scale(0.8, 0.8);-moz-transform:scale(0.8, 0.8);-webkit-transform:scale(0.8, 0.8);-ms-transform:scale(0.8, 0.8);opacity:0}}.mask-in{animation: fade-in 0.2s;-webkit-animation: fade-in 0.2s;}.mask-out{animation: fade-out 0.2s;-webkit-animation: fade-out 0.2s;}.modal-in{animation: modal-in 0.15s;-webkit-animation: modal-in 0.15s;}.modal-out{animation: modal-out 0.15s;-webkit-animation: modal-out 0.15s;}.fade-in{animation: fade-in 0.3s;-webkit-animation: fade-in 0.3s;}.fade-out{animation: fade-out 0.3s;-webkit-animation: fade-out 0.3s;}.error-page{width:800px;text-align:center;margin: 100px auto 0;padding:50px;font-size:18px;line-height:1.5;color:#999;border: 1px solid #e0e0e0;background-color: #fff;}.error-page h2{font-size:40px; font-weight:bold}'
        },
        render: function () {
            if (this._tpl.inStyle) {
                util.insertStyle(this._tpl.inStyle);
            }
        }
    };
    module.exports = layout;
});

/**
 * 事件管理
 * @class event
 * @static
 */
define('nmc/lib/event', function (require, exports, module) {
    let util = require('util');

    //默认判断是否有事件的函数
    let _defalutJudgeFn = function (elem) {
        return !!(elem.getAttribute && elem.getAttribute('data-event'));
    };

    //默认获取事件key的函数
    let _defaultGetEventkeyFn = function (elem) {
        return elem.getAttribute && elem.getAttribute('data-event');
    };

    //添加事件监听
    let addEvent = function (elem, event, fn) {
        if (elem.addEventListener) // W3C
        {
            elem.addEventListener(event, fn, true);
        } else if (elem.attachEvent) { // IE
            elem.attachEvent('on' + event, fn);
        } else {
            elem[event] = fn;
        }
    };

    //获取元素中包含事件的第一个子元素
    let getWantTarget = function (evt, topElem, judgeFn) {

        judgeFn = judgeFn || this.judgeFn || _defalutJudgeFn;

        let _targetE = evt.srcElement || evt.target;

        while (_targetE) {

            if (judgeFn(_targetE)) {
                return _targetE;
            }

            if (topElem == _targetE) {
                break;
            }

            _targetE = _targetE.parentNode;
        }
        return null;
    };

    /**
     * 通用的绑定事件处理
     * @method bindCommonEvent
     * @param {Element} 要绑定事件的元素
     * @param {String} 绑定的事件类型
     * @param {Object} 事件处理的函数映射
     * @param {Function} 取得事件对应的key的函数
     * @author evanyuan
     */
    let bindCommonEvent = function (topElem, type, dealFnMap, getEventkeyFn) {
        getEventkeyFn = getEventkeyFn || _defaultGetEventkeyFn;

        let judgeFn = function (elem) {
            return !!getEventkeyFn(elem);
        };

        let hdl = function (e) {

            /**
             * 支持直接绑定方法
             */
            let _target = getWantTarget(e, topElem, judgeFn),
                _hit = false;

            if (_target) {
                let _event = getEventkeyFn(_target);
                let _returnValue;


                if (Object.prototype.toString.call(dealFnMap) === '[object Function]') {
                    _returnValue = dealFnMap.call(_target, e, _event);
                    _hit = true;
                } else {
                    if (dealFnMap[_event]) {
                        _returnValue = dealFnMap[_event].call(_target, e);
                        _hit = true;
                    }
                }
                if (_hit) {
                    if (!_returnValue) {
                        if (e.preventDefault) {
                            e.preventDefault();
                        } else {
                            e.returnValue = false;
                        }
                    }
                }

            }

        };

        if (type === 'tap') {
            (function () {
                let isTap = true;
                addEvent(topElem, 'touchstart', function () {
                    isTap = true;
                });
                addEvent(topElem, 'touchmove', function () {
                    isTap = false;
                });
                addEvent(topElem, 'touchend', function (e) {
                    if (isTap) {
                        hdl(e);
                    }
                });
            })();

        } else {
            addEvent(topElem, type, hdl);
        }

    };

    let commonEvents = {};

    //新增：保存句柄，用于手动触发事件[解决控制台在移动端点击无法触发菜单的问题] skyzhou
    //commonHandles和commonEvents的用法重复，但做增量修改，不干扰以前代码
    let commonHandles = {};
    /**
     * 为body添加事件代理
     * @method addCommonEvent
     * @param {type} 事件类型
     * @param {dealFnMap} 事件处理的函数映射
     * @author evanyuan
     */
    let addCommonEvent = function (type, dealFnMap) {


        //新增：存储句柄 skyzhou
        if (!commonHandles[type]) {
            commonHandles[type] = {};
        }

        let evtTypeObj = commonEvents[type];
        if (!evtTypeObj) {
            evtTypeObj = commonEvents[type] = {};
        }
        for (let key in dealFnMap) {

            //add skyzhou
            commonHandles[type][key] = dealFnMap[key];

            if (!evtTypeObj[key]) {
                let fnMap = {};
                fnMap[key] = dealFnMap[key];
                if (type == 'mouseenter' || type == 'mouseleave') { //为兼容性,使用jq方法接管
                    $('body').on(type, '[data-event="' + key + '"]', fnMap[key]);
                } else {
                    bindCommonEvent(document.body, type, fnMap);
                }
                evtTypeObj[key] = 1;
            }
        }
    };

    //新增：手动触发事件 skyzhou
    /**
     @param type {String} 类型，比如:click
     @param key {String} key，触发该事件类型下所有句柄中key为入参的函数
     @param e {Event} 事件
     */
    let emit = function (target, type, e) {

        let key,
            fn;

        if (commonHandles[type]) {

            key = target.getAttribute('data-event');
            fn = commonHandles[type][key];

            if (fn) {
                fn.call(target, e);
            }
        }
    };

    //绑定代理事件，自定义代理对象
    exports.bindCommonEvent = bindCommonEvent;

    //统一绑定body的代理事件
    exports.addCommonEvent = addCommonEvent;

    //add by skyzhou 增减手动触发事件
    exports.emit = emit;
});
define('nmc/lib/jquery-1.10.2', function (require, exports, module) {
    return module.exports = require('jquery-1.10.2');
});

define('nmc/lib/util', function (require, exports, module) {
    let $ = require('$');
    window.console = window.console || {
        log: function () {}
    };
    let csrfCode = '';

    /**
     * 工具类
     * @class util
     * @static
     */
    let util = {
        cookie: {
            /**
             * 获取cookie
             * @method get
             * @param  {String} name 名称
             * @return {String}
             */
            get: function (name) {
                let r = new RegExp('(?:^|;+|\\s+)' + name + '=([^;]*)'),
                    m = document.cookie.match(r);

                return !m ? '' : m[1];
            },
            /**
             * 设置cookie
             * @method set
             * @param {String} name 名称
             * @param {String} value 值
             * @param {String} domain 域
             * @param {String} path 路径
             * @param {String} hour 过期时间(小时)
             */
            set: function (name, value, domain, path, hour) {
                if (hour) {
                    var expire = new Date();
                    expire.setTime(expire.getTime() + 36E5 * hour);
                }
                document.cookie = name + '=' + value + '; ' + (hour ? 'expires=' + expire.toGMTString() + '; ' : '') +
                    (path ? 'path=' + path + '; ' : 'path=/; ') + (domain ? 'domain=' + domain + ';' : 'domain=' + document.domain + ';');

                return true;
            },

            /**
             * 设置cookie（包含 secure）
             * @method setSecure
             * @param {String} name 名称
             * @param {String} value 值
             * @param {String} domain 域
             * @param {String} path 路径
             * @param {String} hour 过期时间(小时)
             */
            setSecure: function (name, value, domain, path, hour) {
                if (hour) {
                    var expire = new Date();
                    expire.setTime(expire.getTime() + 36E5 * hour);
                }
                document.cookie = name + '=' + value + '; ' + (hour ? 'expires=' + expire.toGMTString() + '; ' : '') +
                    (path ? 'path=' + path + '; ' : 'path=/; ') + (domain ? 'domain=' + domain + '; secure' : 'domain=' + document.domain + '; secure');

                return true;
            },

            /**
             * 删除cookie
             * @method del
             * @param {String} name 名称
             * @param {String} domain 域
             * @param {String} path 路径
             */
            del: function (name, domain, path) {
                document.cookie = name + '=; expires=Mon, 26 Jul 1997 05:00:00 GMT; ' +
                    (path ? 'path=' + path + '; ' : 'path=/; ') +
                    (domain ? 'domain=' + domain + ';' : 'domain=' + document.domain + ';');
            }
        },

        /**
         * html模板生成器, =号转义, -号原始输出
         * @method tmpl
         * @param  {String} str html模板字符串 | script模版元素Id
         * @param  {Object} data 用于生成模板的数据对象
         * @param  {Object} [mixinTmpl] 混合模版对象
         * @return {String} 返回 html 字符串
         * @author evanyuan
         * @example
         *		var careerTmpl = '<div><%=career%></div>';
         * 		util.tmpl('<h1><%=user%></h1> <%#careerTmpl%>', {user:'evanyuan', career: '前端工程师'}, {careerTmpl: careerTmpl});
         */
        tmpl: (function () {
            var _cache = {},
                _escape = function (str) {
                    if (str == 0) {
                        return str;
                    }
                    str = (str || '').toString();
                    return str.replace(/&(?!\w+;)/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\"/g, '&quot;');
                },
                _getTmplStr = function (rawStr, mixinTmpl) {
                    if (mixinTmpl) {
                        for (var p in mixinTmpl) {
                            var r = new RegExp('<%#\\s?' + p + '%>', 'g');
                            rawStr = rawStr.replace(r, mixinTmpl[p]);
                        }
                    }
                    return rawStr;
                };

            return function tmpl(str, data, mixinTmpl) {
                var strIsKey = !/\W/.test(str);
                !strIsKey && (str = _getTmplStr(str, mixinTmpl));

                var fn = strIsKey ? _cache[str] = _cache[str] || tmpl(_getTmplStr(document.getElementById(str).innerHTML, mixinTmpl)) :
                    new Function("obj", "_escape", "var _p='';with(obj){_p+='" + str
                        .replace(/[\r\t\n]/g, " ")
                        .split("\\'").join("\\\\'")
                        .split("'").join("\\'")
                        .split("<%").join("\t")
                        .replace(/\t-(.*?)%>/g, "'+$1+'")
                        .replace(/\t=(.*?)%>/g, "'+_escape($1)+'")
                        .split("\t").join("';")
                        .split("%>").join("_p+='") +
                        "';} return _p;");

                var render = function (data) {
                    if (typeof data == 'object') {
                        data.QCCONSOLE_HOST = window.QCCONSOLE_HOST;
                        data.QCMAIN_HOST = window.QCMAIN_HOST;
                        data.QCBUY_HOST = window.QCBUY_HOST;
                    }
                    return fn(data, _escape)
                };

                return data ? render(data) : render;
            };
        })(),

        /**
         * 获取防CSRF串
         * @method getACSRFToken
         * @return {String} 验证串
         */
        getACSRFToken: function () {
            if (!csrfCode) {
                let s_key = this.getSessionKey();
                if (!s_key) {
                    return '';
                }
                let hash = 5381;
                for (let i = 0, len = s_key.length; i < len; ++i) {
                    hash += (hash << 5) + s_key.charCodeAt(i);
                }
                csrfCode = hash & 0x7fffffff;
            }
            return csrfCode;
        },

        /**
         * 获取uin
         * @method getUin
         * @return {String} uin
         */
        getUin: function () {
            return parseInt(this.cookie.get('uin').replace(/\D/g, ''), 10) || '';
        },

        /**
         * 获取 ownerUin
         */
        getOwnerUin: function () {
            return parseInt(this.cookie.get('ownerUin').replace(/\D/g, ''), 10) || '';
        },

        /**
         * 获取 skey
         */
        getSessionKey: function () {
            let skey = this.cookie.get('skey') || this.cookie.get('p_skey');
            return skey ? decodeURIComponent(skey) : null;
        },

        setACSRFToken: function (token) {
            csrfCode = token;
        },

        /** XSS 转义**/
        escHTML: function (str) {
            if (str == 0) {
                return str;
            }
            str = (str || '').toString();
            return str.replace(/&(?!\w+;)/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\"/g, '&quot;');
        },

        /**
         * 将 URL 参数格式转化成对象
         * @method paramsToObject
         * @param  {String} [queryString] 要转换的 key-value 字符串，默认为 location.search
         * @return {Object}
         * @author evanyuan
         */
        paramsToObject: function (queryString) {
            let _result = {},
                _pairs,
                _pair,
                _query,
                _key,
                _value;

            if (typeof (queryString) === 'object') {
                return queryString;
            }

            _query = queryString || window.location.search;
            _query = _query.replace('?', '');
            _pairs = _query.split('&');

            $(_pairs).each(function (i, keyVal) {
                _pair = keyVal.split('=');
                _key = _pair[0];
                _value = _pair.slice(1).join('=');
                _result[decodeURIComponent(_key)] = decodeURIComponent(_value);
            });

            return _result;
        },

        /**
         * JSON对象转url字符串
         * @method objectToParams
         * @param  {Object} obj JSON对象
         * @param  {Boolean} decodeUri url解码
         * @return {String} url字符串
         * @author evanyuan
         */
        objectToParams: function (obj, decodeUri) {
            let param = $.param(obj);
            if (decodeUri) {
                param = decodeURIComponent(param);
            }
            return param;
        },

        /**
         * 是否移动手机
         * @method isMobile
         * @return {boolean} true|false
         */
        isMobile: function () {
            return this.isAndroid() || this.isIOS();
        },

        /**
         * 是否android
         * @method isAndroid
         * @return {boolean} true|false
         */
        isAndroid: function () {
            return /android/i.test(window.navigator.userAgent);

        },

        /**
         * 是否ios
         * @method isIOS
         * @return {boolean} true|false
         */
        isIOS: function () {
            return /iPod|iPad|iPhone/i.test(window.navigator.userAgent);
        },

        /**
         * 获取a标签href相对地址
         * @method getHref
         * @param  {Object} item dom节点
         * @return {String} href
         * @author evanyuan
         */
        getHref: function (item) {
            let href = item.getAttribute('href', 2);
            let replaceStr = location.protocol + '//' + location.host;
            if (href.indexOf('//') == 0) {
                replaceStr = '//' + location.host;
            }
            href = href.replace(replaceStr, '');
            return href;
        },

        /**
         * 深度拷贝对象
         * @method cloneObject
         * @param  {Object} obj 任意对象
         * @param  {Object} [options] 其他选项
         * @return {Object} 返回新的拷贝对象
         * @author evanyuan
         */
        cloneObject: function (obj, options) {
            let o = obj.constructor === Array ? [] : {};
            for (let i in obj) {
                if (options && options.ignore && options.ignore.indexOf(i) >= 0) {
                    continue;
                }
                if (obj.hasOwnProperty(i)) {
                    o[i] = (obj[i] && (typeof obj[i] === 'object')) ? this.cloneObject(obj[i]) : obj[i];
                }
            }
            return o;
        },

        /**
         * 插入内部样式
         * @method insertStyle
         * @param  {string | Array} rules 样式
         * @param  {string} id 样式节点Id
         * @author evanyuan
         */
        insertStyle: function (rules, id) {
            let _insertStyle = function () {
                let doc = document,
                    node = doc.createElement('style');
                node.type = 'text/css';
                id && (node.id = id);
                document.getElementsByTagName('head')[0].appendChild(node);
                if (rules) {
                    if (typeof (rules) === 'object') {
                        rules = rules.join('');
                    }
                    if (node.styleSheet) {
                        node.styleSheet.cssText = rules;
                    } else {
                        node.appendChild(document.createTextNode(rules));
                    }
                }
            };
            if (id) {
                !document.getElementById(id) && _insertStyle();
            } else {
                _insertStyle();
            }
        },

        /**
         * 检测浏览器是否支持css3 animation属性
         * @method supportCss3Animation
         * @return {Boolean}
         * @author evanyuan
         */
        supportCss3Animation: function () {
            let element = document.createElement('div');
            if ('animation' in element.style || 'webkitAnimation' in element.style) {
                return true;
            } else {
                return false;
            }
        },

        /**
         * 动画结束后回调
         * @method animationend
         * @param  {Object} el dom元素
         * @param  {Function} endCb 回调
         * @author evanyuan
         */
        animationend: function (el, endCb) {
            if (this.supportCss3Animation()) {
                let $el = $(el);
                let removeEvt = function () {
                    let _cb = $el.data('cb');
                    el.removeEventListener('animationend', _cb);
                    el.removeEventListener('webkitAnimationEnd', _cb);
                };
                let cb = function () {
                    endCb && endCb();
                    removeEvt();
                };
                removeEvt();
                el.addEventListener('webkitAnimationEnd', cb);
                el.addEventListener('animationend', cb);
                $el.data('cb', cb);
            } else {
                if (endCb) {
                    // IE9 模拟一个动画进程，否则某些异步过程会失序
                    setTimeout(endCb, 1);
                }
            }
        },

        /**
         * requestAnimationFrame
         * @method requestAnimationFrame
         * @param  {Function} callback 回调
         * @author evanyuan
         */
        requestAnimationFrame: function (callback) {
            let _requestAnimationFrame = window.requestAnimationFrame ||
                window.webkitRequestAnimationFrame ||
                window.mozRequestAnimationFrame ||
                function (callback) {
                    window.setTimeout(callback, 1000 / 60);
                };
            _requestAnimationFrame(callback);
        },

        /**
         * 检测浏览器是否支持css3 transform属性
         * @method supportTransform
         * @return {Boolean}
         * @author evanyuan
         */
        supportTransform: function () {
            let element = document.createElement('div'),
                style = element.style;
            if ('WebkitTransform' in style || 'MozTransform' in style || 'OTransform' in style ||
                'Transform' in style || 'transform' in style) {
                return true;
            } else {
                return false;
            }
        },

        /**
         * 检测浏览器是否支持css3 transition属性
         * @method supportTransition
         * @return {Boolean}
         * @author evanyuan
         */
        supportTransition: function () {
            let element = document.createElement('div'),
                style = element.style;
            if ('transition' in style || 'WebkitTransition' in style || 'MozTransition' in style ||
                'MsTransition' in style || 'OTransition' in style) {
                return true;
            } else {
                return false;
            }
        },

        /**
         * 为元素设置transform与transition
         * @method setTransformTransitionForElem
         * @param  {Object} elem dom元素
         * @param  {String} transform transform值
         * @param  {String} transition transition值
         * @author brianlin
         */
        setTransformTransitionForElem: function (elem, transform, transition) {
            if (!elem) {
                return;
            }
            if (transform !== undefined) {
                elem.style.WebkitTransform = transform;
                elem.style.MozTransform = transform;
                elem.style.OTransform = transform;
                elem.style.msTransform = transform;
                elem.style.Transform = transform;
                elem.style.transform = transform;
            }
            if (transition !== undefined) {
                elem.style.WebkitTransition = transition;
                elem.style.MozTransition = transition;
                elem.style.OTransition = transition;
                elem.style.msTransition = transition;
                elem.style.Transition = transition;
                elem.style.transition = transition;
            }
        },

        classnames: (function () {
            let hasOwn = {}.hasOwnProperty;

            function classNames() {
                let classes = [];
                for (let i = 0; i < arguments.length; i++) {
                    let arg = arguments[i];
                    if (!arg) continue;

                    let argType = typeof arg;

                    if (argType === 'string' || argType === 'number') {
                        classes.push(arg);
                    } else if (Array.isArray(arg)) {
                        classes.push(classNames.apply(null, arg));
                    } else if (argType === 'object') {
                        for (let key in arg) {
                            if (hasOwn.call(arg, key) && arg[key]) {
                                classes.push(key);
                            }
                        }
                    }
                }

                return classes.join(' ');
            }

            return classNames;
        })(),

        date: {
            /**
             * 格式化日期文本为日期对象
             *
             * @method str2Date
             * @param {String} date 文本日期
             * @param {String} [p:%Y-%M-%d %h:%m:%s] 文本日期的格式
             * @return {Date}
             */
            str2Date: function (date, p) {
                // add 20171225
                // 针对"2021-01-10T08:39:48Z"格式日期预处理
                date = date.replace(/[TZ]/gi, ' ').replace(/\s+$/, '');
                p = p || '%Y-%M-%d %h:%m:%s';
                p = p.replace(/\-/g, '\\-');
                p = p.replace(/\|/g, '\\|');
                p = p.replace(/\./g, '\\.');
                p = p.replace('%Y', '(\\d{4})');
                p = p.replace('%M', '(\\d{1,2})');
                p = p.replace('%d', '(\\d{1,2})');
                p = p.replace('%h', '(\\d{1,2})');
                p = p.replace('%m', '(\\d{1,2})');
                p = p.replace('%s', '(\\d{1,2})');

                let regExp = new RegExp('^' + p + '$'),
                    group = regExp.exec(date),
                    Y = (group[1] || 0) - 0,
                    M = (group[2] || 1) - 1,
                    d = (group[3] || 0) - 0,
                    h = (group[4] || 0) - 0,
                    m = (group[5] || 0) - 0,
                    s = (group[6] || 0) - 0;

                return new Date(Y, M, d, h, m, s);
            },

            /**
             * 格式化日期为指定的格式
             *
             * @method date2Str
             * @param {Date} date
             * @param {String} p 输出格式, %Y/%M/%d/%h/%m/%s的组合
             * @param {Boolean} [isFill:false] 不足两位是否补0
             * @return {String}
             */
            date2Str: function (date, p, isFill) {
                let Y = date.getFullYear(),
                    M = date.getMonth() + 1,
                    d = date.getDate(),
                    h = date.getHours(),
                    m = date.getMinutes(),
                    s = date.getSeconds();

                if (isFill) {
                    M = (M < 10) ? ('0' + M) : M;
                    d = (d < 10) ? ('0' + d) : d;
                    h = (h < 10) ? ('0' + h) : h;
                    m = (m < 10) ? ('0' + m) : m;
                    s = (s < 10) ? ('0' + s) : s;
                }
                p = p || '%Y-%M-%d %h:%m:%s';
                p = p.replace(/%Y/g, Y).replace(/%M/g, M).replace(/%d/g, d).replace(/%h/g, h).replace(/%m/g, m).replace(/%s/g, s);
                return p;
            },

            /**
             * 日期比较(d1 - d2)
             *
             * @method dateDiff
             * @param {Date} d1
             * @param {Date} d2
             * @param {String} [cmpType:ms] 比较类型, 可选值: Y/M/d/h/m/s/ms -> 年/月/日/时/分/妙/毫秒
             * @return {Float}
             */
            dateDiff: function (d1, d2, cmpType) {
                let diff = 0;
                switch (cmpType) {
                    case 'Y':
                        diff = d1.getFullYear() - d2.getFullYear();
                        break;
                    case 'M':
                        diff = (d1.getFullYear() - d2.getFullYear()) * 12 + (d1.getMonth() - d2.getMonth());
                        break;
                    case 'd':
                        diff = (d1 - d2) / 86400000;
                        break;
                    case 'h':
                        diff = (d1 - d2) / 3600000;
                        break;
                    case 'm':
                        diff = (d1 - d2) / 60000;
                        break;
                    case 's':
                        diff = (d1 - d2) / 1000;
                        break;
                    default:
                        diff = d1 - d2;
                        break;
                }
                return diff;
            },
            /**
             * 日期相加
             *
             * @method dateAdd
             * @param char interval 间隔参数
             *        y 年
             *        q 季度
             *        n 月
             *        d 日
             *        w 周
             *        h 小时
             *        m 分钟
             *        s 秒
             *        i 毫秒
             * @param {Date} indate 输入的日期
             * @param {Number} offset 差值
             * @return {Date} date 相加后的日期
             */
            dateAdd: function (interval, indate, offset) {
                switch (interval) {
                    case 'y':
                        indate.setFullYear(indate.getFullYear() + offset);
                        break;
                    case 'q':
                        indate.setMonth(indate.getMonth() + (offset * 3));
                        break;
                    case 'n':
                        indate.setMonth(indate.getMonth() + offset);
                        break;
                    case 'd':
                        indate.setDate(indate.getDate() + offset);
                        break;
                    case 'w':
                        indate.setDate(indate.getDate() + (offset * 7));
                        break;
                    case 'h':
                        indate.setHours(indate.getHours() + offset);
                        break;
                    case 'm':
                        indate.setMinutes(indate.getMinutes() + offset);
                        break;
                    case 's':
                        indate.setSeconds(indate.getSeconds() + offset);
                        break;
                    case 'i':
                        indate.setMilliseconds(indate.getMilliseconds() + offset);
                        break;
                    default:
                        indate.setMilliseconds(indate.getMilliseconds() + offset);
                        break;
                }
                return indate;
            },
            /**
             * 判断是否是闰年
             *
             * @method leapYear
             * @param {Date} indate 输入的日期
             * @return {Object} 对象(是否是闰年，各月份的天数集，当前月的天数)
             */
            leapYear: function (indate) {
                let _days = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];
                let _is = false;
                let _d = 365;

                if ((indate.getFullYear() % 4 === 0 && indate.getFullYear() % 100 !== 0) || indate.getFullYear() % 400 === 0) {
                    _days.splice(1, 1, 29);
                    _is = true;
                    _d = 366;
                } else {
                    _days.splice(1, 1, 28);
                    _is = false;
                    _d = 365;
                }
                return {
                    isLeapYear: _is,
                    days: _days,
                    yearDays: _d,
                    monthDays: _days[indate.getMonth()]
                };
            },
            /**
             * 转换日期格式
             */
            parseDate: function (str, type) {
                if (!str) {
                    return '';
                }
                let date = null,
                    newStr;
                if (str instanceof Date) {
                    date = str;
                } else {
                    date = this.str2Date(str);
                }
                if (isNaN(date)) {
                    return str;
                }
                let year = date.getFullYear();

                let month = date.getMonth() + 1;
                month = month >= 10 ? month : '0' + month;

                let day = date.getDate();
                day = day >= 10 ? day : '0' + day;

                let hour = date.getHours();
                hour = hour >= 10 ? hour : '0' + hour;

                let minute = date.getMinutes();
                minute = minute >= 10 ? minute : '0' + minute;

                let second = date.getSeconds();
                second = second >= 10 ? second : '0' + second;

                switch (type) {
                    case 'short': //短格式，如2011-05-06
                        return year + '-' + month + '-' + day;
                        break;
                    case 'long': //长格式，如2011-05-06 10:05:06
                        return year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + day;
                        break;
                    case 'chinese': //中文格式，2011年05月06日 23:09
                        return year + '年' + month + '月' + day + '日 ' + hour + ':' + minute;
                        break;
                    case 'monthday':
                        return month + '月' + day + '日 ' + hour + ':' + minute;
                        break;
                    default:
                        return str;
                        break;
                }
            }
        }

    };

    module.exports = util;
});
define('nmc/main/entry', function (require, exports) {
    let $ = require('$');
    let evt = require('event');
    let util = require('util');
    let nmcConfig = require('nmcConfig');
    let router = require('router');
    let pageManager = require('pageManager');
    let tips = require('tips');

    //nmc初始化
    let init = function (config) {

        //参数扩展合并
        config = config || {};
        $.extend(true, nmcConfig, config);

        //全局加载提示
        tips.initLoading();

        //初始化页面管理
        pageManager.init();

        //初始化路由
        router.init({
            'html5Mode': true,
            'pageManager': pageManager,
            'routes': {
                '/': 'loadRoot',
                '/*controller(/*action)(/*p1)(/*p2)(/*p3)(/*p4)': 'loadCommon'
            },
            'extendRoutes': nmcConfig.extendRoutes,
            'beforeNavigate': nmcConfig.beforeNavigate,
            'OnUrlChange': nmcConfig.OnUrlChange
        });

        //全局点击
        evt.addCommonEvent('click', {
            'nav': function (e) {
                let url = util.getHref(this);

                if (e.altKey || e.ctrlKey || e.shiftKey || e.metaKey || e.which === 2) {
                    return true;
                }

                if (url.indexOf('http') == 0 || url.indexOf('//') == 0) { //如果替换之后，还是绝对地址，则应该是外站地址
                    return true;
                }
                router.navigate(url);

                let topSubMenuTrigger = $(this).closest('#topnav');

                if (topSubMenuTrigger.length) {
                    setTimeout(function () {
                        $('#topnav').find('[data-event="top_submenu"]').removeClass('qc-nav-hover qc-header-panel-border qc-nav-select');

                        $('.qc-nav-mobile-menu').removeClass('qc-nav-btn-show');

                        $('#topnav').removeClass('qc-mobile-user-show qc-mobile-menu-show');

                        $('#topnav').trigger('header_right_tool_show');
                    }, 0);
                }
            }
        });

        //记录所有请求完毕
        let win = window;
        $(win).load(function () {
            win.isOnload = true;
        });

    };

    exports.init = init;
});

define('nmc/main/pagemanager', function (require, exports, module) {


    let $ = require('$');
    let router = require('router');
    let util = require('util');
    let nmcConfig = require('nmcConfig');
    /**
     * 页面管理
     * @class pageManager
     * @static
     */
    let pageManager = {

        /**
         * 初始化
         * @method init
         * @author evanyuan
         */
        init: function () {
            /**
             * 页面包裹容器
             * @property pageWrapper
             * @type Object
             */
            this.pageWrapper = $(nmcConfig.pageWrapper);
        },


        /**
         * 加载首页
         */
        loadRoot: function () {
            this.loadView(nmcConfig.root);
        },

        /**
         * 统一加载视图方法
         */
        loadCommon: function () {
            let _self = this,
                arr = [].slice.call(arguments);

            //解析路由匹配
            this.parseMatch(arr, function (controller, action, params) {
                //处理路由, 加载视图
                _self.loadView(controller, action, params);
            });
        },

        /**
         * 解析路由匹配
         * @method parseMatch
         * @param {Array}    arr 路由匹配到的参数
         * @param {Function} cb  回调函数
         * @author evanyuan
         */
        parseMatch: function (arr, cb) {
            let controller = null,
                action = null,
                params = [];

            //获取controller
            controller = arr[0];

            //获取action与params
            if (arr.length > 1) {
                if (typeof (arr[1]) === 'object') {
                    params.push(arr[1]);
                } else {
                    action = arr[1];
                    params = arr.slice(2);
                }
            }

            cb(controller, action, params);

        },

        /**
         * 统一路由处理函数
         * @method loadView
         * @param {String} controller
         * @param {String} action
         * @param {Array} params
         * @author evanyuan
         */
        loadView: function (controller, action, params) {
            let isSpaIframe = false;

            if (params && params[2]) {
                isSpaIframe = !!params[2].isSpaIframe;
            }

            let _self = this;
            let redirectConfig = window.g_config_data.redirect_config || {};

            //业务重定向
            if (controller === 'iframe' && params[0]) {
                if (redirectConfig[params[0]]) {
                    return router.redirect(redirectConfig[params[0]], false, true);
                }
            } else {
                if (redirectConfig[controller]) {
                    return router.redirect(redirectConfig[controller], false, true);
                }
            }

            //渲染前执行业务逻辑
            if (nmcConfig.beforeRender) {
                if (nmcConfig.beforeRender(controller, action, params) === false) {
                    return;
                }
            }

            if (_self.currentViewObj) {

                //全局销毁
                _self.globalDestroy();

                //记录页面模块状态
                _self.currentViewObj.pageActive = false;

                //销毁前一个
                let destroy = _self.currentViewObj.destroy;
                try {
                    if (!isSpaIframe) {
                        destroy && destroy.call(_self.currentViewObj);
                    }
                    if (_self.currentCtrlObj) {
                        _self.currentCtrlObj.pageActive = false;
                        let ctrlDestroy = _self.currentCtrlObj.destroy;
                        ctrlDestroy && ctrlDestroy.call(_self.currentCtrlObj);
                    }
                } catch (e) {
                    window.console && console.error && console.error('View destroy failed ', e);
                }
                _self.currentCtrlObj = null;
                _self.currentViewObj = null;
            }

            // 调用 loadView 之前可以指定一次性的 appArea 清除计划
            if (_self.cleanupBeforeNextRender) {
                if (_self.appArea) {
                    _self.appArea.html('');
                }
                _self.cleanupBeforeNextRender = false;
            }

            _self.curController = controller;
            _self.curAction = action;

            params = params || [];

            _self.curParams = params;

            //渲染公共模版
            this.renderLayout(controller, action, params);

            //存储主要jQuery dom对象

            /**
             * 页面主容器
             * @property container
             * @type Object
             */
            this.container = $(nmcConfig.container);

            /**
             * 左侧导航容器
             * @property sidebar
             * @type Object
             */
            this.sidebar = $(nmcConfig.sidebar);

            /**
             * 右侧内容容器
             * @property appArea
             * @type Object
             */
            this.appArea = $(nmcConfig.appArea);

            /**
             * 切换页面需要更改class的容器
             * @property classWrapper
             * @type Object
             */
            this.classWrapper = $(nmcConfig.classWrapper);

            //模块基础路径
            let basePath = nmcConfig.basePath;

            //模块id按照如下规则组成
            let controllerId = basePath + controller + '/' + controller,
                actionId = basePath + controller + '/' + action + '/' + action;

            let moduleArr = [];

            // 向前兼容总览迁移，路由未发布前，将原home模块前增加deprecated/，避免无法找到对应js文件
            // TODO: 新版总览发布后可以将该逻辑去掉
            if (controller === 'home' && !seajs.hasDefined[controllerId]) {
                controllerId = 'deprecated/' + controllerId;
            }

            //检查是否存在controller模块
            if (seajs.hasDefined[controllerId]) {
                moduleArr.push(controllerId);
            } else {
                controllerId = '';
            }

            //检查是否存在action模块
            if (action) {
                if (!seajs.hasDefined[actionId]) {
                    return _self.tryIframeRouter(controller);
                }
                moduleArr.push(actionId);
            } else {
                // 未指明action，默认尝试查询index
                let indexUri = basePath + controller + '/index/index';
                if (seajs.hasDefined[indexUri]) {
                    moduleArr.push(indexUri);
                    action = 'index';
                } else {
                    //未指明action，且controller也不曾定义
                    if (!controllerId) {
                        return _self.tryIframeRouter(controller);
                    }
                }
            }

            let _this = this;
            let _renderSidebar = function (callback) {
                //渲染左侧菜单
                if (_this.appArea.length) {
                    let sidebarRenderRouter = _this.getSidebarRenderRouter();

                    if (sidebarRenderRouter === '') {
                        sidebarRenderRouter = 'home';
                    }

                    // cluezhang: 简单的indexOf判断会有问题，比如从/vpc2跳到/vpc，会判断为没有切换controller，而不会变更导航。
                    //if (!_self.fragment || _self.fragment.indexOf('/' + tgtCtrlStr) != 0) {
                    if (sidebarRenderRouter) {
                        _this.container.removeClass('qc-aside-hover');
                        _this.sidebar.removeClass('qc-aside-new');

                        if (nmcConfig.subMenu[sidebarRenderRouter]) {
                            require.async('sidebar', function (sidebar) {
                                sidebar.render(sidebarRenderRouter, callback);
                            });
                        } else {
                            _this.hideSidebar();
                            callback();
                        }
                    } else {
                        callback();
                    }
                } else {
                    callback();
                }
            };

            let _getModuleExport = function () {
                //获取页面模块对外接口
                require.async(moduleArr, function (cObj, aObj) {

                    //controller未定义, 此时cObj属于一个action
                    if (!controllerId) {
                        aObj = cObj;
                    }

                    //执行controller, 判断同contoller下的action切换, contoller不需要再重复执行
                    if (controllerId && (!_self.fragment || _self.fragment.indexOf('/' + controller) < 0 || !action)) {

                        _self.fragment = (router.fragment === '/') ? '/' + controller : router.fragment;
                        _self.fragment = _self.fragment.replace(/\/?\?.*/, '');
                        _self.renderView(cObj, params);
                    }
                    _self.fragment = (router.fragment === '/') ? '/' + controller : router.fragment;
                    _self.fragment = _self.fragment.replace(/\/?\?.*/, '');

                    //执行action
                    if (action) {
                        _self.renderView(aObj, params);
                        _self.currentViewObj = aObj;
                        controllerId && (_self.currentCtrlObj = cObj);
                    } else {
                        _self.currentViewObj = cObj;
                    }

                    //更改导航状态
                    if (nmcConfig.changeNavStatus) {
                        nmcConfig.changeNavStatus(controller, action, params);
                    } else {
                        _self.changeNavStatus(controller, action, params);
                    }

                    //设置页面标题
                    _self.setTitle(cObj, aObj);

                    if (!this.reported) {
                        this.reported = true;
                        let now = Date.now();
                        let speedPoints = {
                            domready: now - performance.timing.fetchStart,
                            loadurl: now - window.startLoadJs
                        };
                    }

                });

            };

            // iframe 按 businessKey 加载 css 资源
            if (controller === 'iframe' && params[0]) {
                controller = params[0];
            }

            //需加载的css资源
            let cssRouter = this.getCssRouter();
            let cssReqArr = this.getCssHref(cssRouter);

            if (cssReqArr.length) {
                this.loadPageCss(cssRouter, function () {
                    _self.disablePageCss(cssRouter);
                    _renderSidebar(_getModuleExport);
                });
            } else {
                _self.disablePageCss();
                _renderSidebar(_getModuleExport);
            }

        },

        getSidebarRenderRouter: function () {
            let self = this;
            let prevSidebarRenderRouter = this.sidebarRenderRouter;
            let fragment = /^\/[^\?#@]*/.exec(router.fragment);

            if (fragment && fragment.length) {
                fragment = fragment[0];
            }

            if (fragment.indexOf('/') == 0) {
                fragment = fragment.substr(1);
            }

            do {
                if (nmcConfig.subMenu[fragment] || fragment.indexOf('/') == -1) {
                    break;
                }

                fragment = fragment.split('/');
                fragment.pop();
                fragment = fragment.join('/');

            } while (fragment.indexOf('/') != -1);

            self.sidebarRenderRouter = fragment;

            if (prevSidebarRenderRouter == fragment && self.sidebar.find('.qc-menu-fold').length) {
                return false;
            } else {
                if (prevSidebarRenderRouter && nmcConfig.subMenu[fragment] && nmcConfig.subMenu[fragment].relate == prevSidebarRenderRouter) {
                    return false;
                }
                return fragment;
            }
        },

        /**
         * 加载页面级css
         */
        loadPageCss: function (controller, cb) {
            let cssReqArr = this.getCssHref(controller),
                $head = $('head');

            if (!cssReqArr.length) {
                cb && cb();
                return;
            }

            if (cssReqArr.length && window.cssVersion) {
                for (var i = 0, reqLinkItem; reqLinkItem = cssReqArr[i]; i++) {
                    let cssName = reqLinkItem.substring(reqLinkItem.lastIndexOf('/') + 1);
                    let cssV = cssVersion[cssName];
                    if (cssV) {
                        cssReqArr[i] = cssReqArr[i] + '?max_age=' + cssV;
                    }
                }
            }

            require.async(cssReqArr, function () {
                window.cssRequestStart = Date.now();
                cb && cb();
            });

            let cssLink = $head.find('link[rel="stylesheet"]').not('[data-role="global"]'),
                tgtClsName = controller.replace(/\//gm, '_') + '-css',
                tgtLink = cssLink.filter('.' + tgtClsName);

            if (tgtLink.length) {
                tgtLink.removeAttr('disabled').prop('disabled', false);
            } else {
                cssLink.each(function (i, el) {
                    let $el = $(el);
                    let linkHref = $el.attr('href');
                    for (var j = 0, reqLinkItem; reqLinkItem = cssReqArr[j]; j++) {
                        if (linkHref == reqLinkItem) {
                            $el.addClass(tgtClsName).removeAttr('disabled').prop('disabled', false);
                            break;
                        }
                    }
                });
            }
        },

        /**
         * 禁用页面级css
         */
        disablePageCss: function (controller) {
            let $head = $('head'),
                cssLink = $head.find('link[rel="stylesheet"]').not('[data-role="global"]');
            if (controller) {
                let tgtClsName = controller.replace(/\//gm, '_') + '-css';
                cssLink = cssLink.not('.' + tgtClsName);
            }
            cssLink.attr('disabled', true).prop('disabled', true);
        },

        /**
         * 获取css请求地址
         */
        getCssHref: function (controller) {
            let cssConfig = nmcConfig.css,
                controllerCssReq = cssConfig[controller];
            return controllerCssReq || [];
        },

        /**
         * 获取配置了css路由片断
         */
        getCssRouter: function () {
            let self = this;

            let fragment = /^\/[^\?#@]*/.exec(router.fragment);
            if (fragment && fragment.length) {
                fragment = fragment[0];
            }

            if (fragment.indexOf('/') == 0) {
                fragment = fragment.substr(1);
            }

            // 总览
            if (fragment == '' || fragment == '/') {
                return 'home';
            }

            do {
                if (self.getCssHref(fragment).length || fragment.indexOf('/') == -1) {
                    break;
                }

                fragment = fragment.split('/');
                fragment.pop();
                fragment = fragment.join('/');

            } while (fragment.indexOf('/') != -1);

            return fragment;
        },

        /**
         * 渲染公共模版
         */
        renderLayout: function (controller, action, params) {
            let _self = this,
                layoutConfig = nmcConfig.layout,
                layout = 'default',
                _render = function (layoutName) {
                    if (_self.layout != layoutName) {
                        require.async(layoutConfig[layoutName]['module'], function (_layout) {
                            _layout.render();
                        });
                        _self.layout = layoutName;
                    }
                };

            loop: for (let key in layoutConfig) {
                let controllerArr = layoutConfig[key].controller || [];
                for (var i = 0, c; c = controllerArr[i]; i++) {
                    if (controller === c) {
                        layout = key;
                        break loop;
                    }
                }
            }
            _render(layout);
        },

        /**
         * 渲染视图
         */
        renderView: function (obj, params) {
            let defaultClass = nmcConfig.defaultClass,
                classWrapper = this.classWrapper,
                controller = this.curController,
                action = this.curAction;

            //渲染内容前执行业务逻辑
            let beforeContentRender = obj.beforeContentRender ? obj.beforeContentRender.bind(obj) :
                nmcConfig.beforeContentRender ? nmcConfig.beforeContentRender.bind(nmcConfig) : null;
            if (beforeContentRender && beforeContentRender(obj) === false) {
                classWrapper.attr('class', defaultClass);
                return;
            }

            if (obj) {

                if (controller !== 'iframe' && !obj.pageClass) {
                    let containerClass = 'container-' + controller;
                    if (action) {
                        containerClass += '-' + action;
                    }
                    obj.pageClass = containerClass;
                }

                // 页面模块通过属性pageClass来变更样式
                let classNames = [];

                // 重置容器class之前，先判断当前左侧是展开还是收起，重置后保留该状态
                if (classWrapper.hasClass('qc-aside-hidden')) {
                    // 移动端点击左侧菜单时路由跳转，需要将左侧在顶部位置收起，否则路由跳转后页面还是被左侧菜单挡住
                    // 移动端不加class默认收起，加了qc-aside-hidden展开，与pc端刚好相反
                    if (!util.isMobile()) {
                        classNames.push('qc-aside-hidden');
                    }
                }

                if (classWrapper.hasClass('qc-aside-hover')) {
                    if (!util.isMobile()) {
                        classNames.push('qc-aside-hover');
                    }
                }

                if (classWrapper.hasClass('qc-animation-empty')) {
                    classNames.push('qc-animation-empty');
                }

                if (obj.pageClass) {
                    classNames.push(defaultClass);
                    classNames.push(obj.pageClass);
                    classWrapper.attr('class', classNames.join(' '));
                } else if (classWrapper.attr('class') !== defaultClass) {
                    classNames.push(defaultClass);
                    classWrapper.attr('class', classNames.join(' '));
                }
            }

            if (obj && obj.render) {
                obj.pageActive = true;
                try {
                    obj.render.apply(obj, params);
                    this.reportVisitProduct();
                } catch (e) {
                    this.renderOops();
                    if (router.debug) {
                        console.error(e);
                    }
                }
            } else {
                this.renderOops();
            }
        },

        //js接入业务部分上线时，其他页面还需要走iframe逻辑
        tryIframeRouter: function (controller) {

            let iframeRouteCopy = nmcConfig.iframeRouteCopy,
                matchReg = new RegExp('^\\/' + controller + '\\('),
                iframeHandleFun;
            if (iframeRouteCopy) {
                for (let key in iframeRouteCopy) {
                    if (matchReg.exec(key)) {
                        iframeHandleFun = iframeRouteCopy[key];
                        break;
                    }
                }
                if (iframeHandleFun) {
                    return iframeHandleFun();
                }
            }
            if (this.isNoRouteMatch(controller)) {
                this.render404();
            } else {
                this.renderOops();
            }
        },

        isNoRouteMatch: function (controller) {
            /**
             * 成功加载控制台后，使用下面的语句求出内部业务列表
             Array.from(
             new Set(
             Object.keys(seajs.hasDefined)
             .filter(x => x.startsWith('modules/'))
             .map(x => x.split('/')[1])
             )
             );
             */
            let internalControllers = ['dashboard', 'home', 'iframe', 'cvm', 'cz', 'cvm2', 'autoscaling'];

            // 命中了内部路由但是又不是内部业务的，认为无路由匹配，可以渲染 404
            return router.curRouteMode === 'internal' && internalControllers.indexOf(String(controller).toLowerCase()) === -1;
        },

        /**
         * 渲染404
         * @method render404
         * @author evanyuan
         */
        render404: function () {
            if (/callback\.html/.test(location.pathname)) {
                return this._renderErrorPage('loginError');
            }
            this._renderErrorPage('404');
        },

        /**
         * 渲染其他错误
         * @method renderOops
         * @author evanyuan
         */
        renderOops: function () {
            // 不在路由白名单时渲染404
            if (this.isNotInRouterWhiteList) {
                this.isNotInRouterWhiteList = false;
                return this.render404();
            }
            this._renderErrorPage('oops');
            // 错误上报
        },


        _renderErrorPage: function (type) {
            this.hideSidebar();
            let errorHtml = '';
            if (type == '404') {
                errorHtml = nmcConfig.html404;
            } else if (type == 'oops') {
                errorHtml = nmcConfig.htmlOops;
            } else if (type == 'loginError') {
                errorHtml = nmcConfig.loginError;
            } else if (type == 'userInfoError') {
                errorHtml = nmcConfig.userInfoError;
            }
            let container = (this.appArea && this.appArea.length) ? this.appArea : this.container;
            if (!container) container = $(nmcConfig.appArea);
            container.html(errorHtml);
        },

        /**
         * 设置页面标题
         */
        setTitle: function (cObj, aObj) {
            if (aObj && aObj.title) {
                document.title = aObj.title;
            } else if (cObj && cObj.title) {
                document.title = cObj.title;
            } else {
                let defaultTitle = nmcConfig.defaultTitle;
                if (document.title != defaultTitle) {
                    document.title = defaultTitle;
                }
            }
        },

        /**
         * 改变导航选中态
         */
        changeNavStatus: function (controller, action, params) {
            let _self = this,
                fragment = this.fragment,
                root = nmcConfig.root,
                navActiveConfig = nmcConfig.navActiveConfig;

            let changeNav = function (navActiveConfigItem) {
                let $container = $(navActiveConfigItem.container);
                let navActiveClass = navActiveConfigItem['className'];
                let navActives = $container.find('.' + navActiveClass);
                let links = $container.find('a');
                let navDropdownFlag = '[data-nav-dropdown]';
                let onActive;

                if (navActiveConfigItem['not'] && navActiveConfigItem['not'].length) {
                    for (var n = 0, nitem; nitem = navActiveConfigItem['not'][n]; n++) {
                        links = links.not(nitem);
                    }
                }

                navActives.removeClass(navActiveClass);

                for (var i = 0, item; item = links[i]; i++) {
                    let href = util.getHref(item);
                    href = href.replace(/\/?\?.*/, '');

                    if (href === '/' && controller === root) {
                        onActive = item;
                        break;
                    } else if (href !== '/' && fragment && fragment.indexOf(href) == 0) {
                        if (fragment === href) {
                            onActive = item;
                            break;
                        } else {
                            if (onActive) {
                                if (href.length > util.getHref(onActive).length) {
                                    onActive = item;
                                }
                            } else {
                                onActive = item;
                            }
                        }
                    } else if (fragment === href) {
                        onActive = item;
                        break;
                    }
                }

                let $onActive = $(onActive);
                if (navActiveConfigItem.selectParent) {
                    $onActive = $onActive.parent();
                }
                $onActive.addClass(navActiveClass);

                let $dropDownParent = $onActive.closest(navDropdownFlag);
                $(navDropdownFlag).removeClass('qc-aside-child-select');
                if ($dropDownParent.length) {
                    $dropDownParent.addClass(navActiveClass);
                    if ($dropDownParent.find('.' + navActiveClass).length) {
                        $dropDownParent.addClass('qc-aside-child-select');
                    }
                }
            };

            for (var i = 0, item; item = navActiveConfig[i]; i++) {
                if (item.dynamic) {
                    let subMenuConfig = nmcConfig['subMenu'];
                    var tgtCtrl = (controller === 'iframe' && params) ? params[0] : controller;

                    if (subMenuConfig[tgtCtrl] && subMenuConfig[tgtCtrl]['dynamicMenu']) {
                        (function (_item) {
                            require.async('sidebar', function (sidebar) {
                                sidebar.getDynamenu(tgtCtrl, function () {
                                    changeNav(_item);
                                });
                            });
                        })(item);
                    } else {
                        changeNav(item);
                    }
                } else {
                    changeNav(item);
                }

            }

        },

        /**
         * 页面切换时全局销毁
         */
        globalDestroy: function () {
            nmcConfig.globalDestroy && nmcConfig.globalDestroy();
        },

        refreshSidebar: function () {
            let _self = this,
                params = this.curParams,
                controller = this.curController,
                action = this.curAction;

            require.async('sidebar', function (sidebar) {
                sidebar.refresh();

                //更改导航状态
                if (nmcConfig.changeNavStatus) {
                    nmcConfig.changeNavStatus(controller, action, params);
                } else {
                    _self.changeNavStatus(controller, action, params);
                }
            });
        },

        /**
         * 左侧导航收起展开
         */
        toggleSidebar: function (status) {
            this.container.removeClass('qc-aside-hover');

            this._displaySidebar({
                status: status,
                animate: true,
                showCb: function () {},
                hideCb: function () {}
            });

        },

        /**
         * 对外接口：用于展示、隐藏左导航(无动画)
         */
        displaySidebar: function (status, notHideHandle) {
            this._displaySidebar({
                status: status,
                animate: false,
                hideHandle: notHideHandle ? false : true
            });
        },

        _displaySidebar: function (options) {
            let _self = this,
                status = options.status,
                animate = options.animate,
                hideHandle = options.hideHandle,
                showCb = options.showCb,
                hideCb = options.hideCb,
                sidebarHandle = this.sidebar.find('[data-event="toggle_sidebar"]'),
                _hideSidebar = function () {
                    _self.hideSidebar(true);
                    hideHandle && sidebarHandle.hide();
                    hideCb && hideCb();

                    if (_self.sidebar.hasClass('qc-aside-new')) {
                        _self.container.addClass('qc-aside-hover');
                    } else {
                        _self.container.removeClass('qc-aside-hover');
                    }
                },
                _showSidebar = function () {
                    _self.showSidebar(true);
                    sidebarHandle.show();
                    showCb && showCb();
                };

            if (!animate) {
                _self.container.removeClass('qc-animation-empty');
            }

            if (status === undefined) {
                if (_self.container.hasClass('qc-aside-hidden')) {
                    _showSidebar();
                } else {
                    _hideSidebar();
                }
            } else {
                if (status) {
                    _showSidebar();
                } else {
                    _hideSidebar();
                }
            }
        },

        /**
         * 显示左侧导航容器
         */
        showSidebar: function (animate, changeNavStatus) {
            if (animate) {
                this.container.removeClass('qc-animation-empty');
            } else {
                this.container.addClass('qc-animation-empty');
            }

            this.container.css('left', 0);
            this.container.removeClass('qc-aside-hidden');

            let params = this.curParams,
                controller = this.curController,
                action = this.curAction;

            // 重新设置左侧菜单高亮， 因为sidebar初次render时因为拉取白名单接口异步返回，导致选中逻辑处理时还未产生左侧结构
            if (changeNavStatus) {
                if (nmcConfig.changeNavStatus) {
                    nmcConfig.changeNavStatus(controller, action, params);
                } else {
                    this.changeNavStatus(controller, action, params);
                }
            }
        },

        /**
         * 隐藏左侧导航容器
         */
        hideSidebar: function (animate, notclear) {
            let leftVal = $('#sidebar').hasClass('qc-aside-new') ? 0 : -200;

            if (!this.container) return;

            if (animate) {
                this.container.removeClass('qc-animation-empty');
                this.container.css('left', leftVal);
            } else {
                this.container.addClass('qc-animation-empty');
                this.container.css('left', leftVal);
                if (!notclear) {
                    this.sidebar.html('');
                }
            }

            if ($('#sidebar').find('.qc-menu-fold').length) {
                this.container.addClass('qc-aside-hidden');
            }
        },

        _animateSidebar: function (val) {
            let _self = this;
            this.container.addClass('animate');
            if (util.supportTransition()) {
                this.container.css('left', val);
            } else {
                this.container.animate({
                    'left': val
                });
            }
        },

        /**
         * 重置fragment标记(用于强制刷新controller)
         * @method resetFragment
         * @author evanyuan
         */
        resetFragment: function () {
            this.fragment = '';
        },

        /**
         * 获取页面开始的时间点
         */
        getPageStartPoint: function () {
            return router.applyActionPoint;
        },

        /**
         * 国际版路由检查
         * @method checkRoutesI18n
         */
        checkRoutesI18n: function (fragment) {
            // 是否能访问国际版由路由配置决定，不再由以下代码
            return true;
        },

        /**
         * 根据路由，反向匹配导航配置，找出productId用来上报
         */
        reportVisitProduct: function () {
            // 循环引用，用事件传递上报操作
            $(document).trigger('reportVisitProduct');
        },
    };

    module.exports = pageManager;
});


define('nmc/main/router', function (require, exports, module) {

    let docMode = document.documentMode;
    let oldIE = (/msie [\w.]+/.test(navigator.userAgent.toLowerCase()) && (!docMode || docMode <= 7));
    let pushState = window.history.pushState;

    /**
     * 路由管理
     * @class router
     * @static
     */
    let router = {
        /**
         * 初始化
         * @param {Object} option 参数
         * @method init
         * @author evanyuan
         */
        init: function (option) {

            this.option = {

                //是否使用html5 history API设置路由
                'html5Mode': true,

                //页面管理对象
                'pageManager': {},

                //路由映射对象
                'routes': {},

                //扩展路由，优先于框架内部路由routes对象
                'extendRoutes': {},

                //低端浏览器监听url变化的时间间隔
                'interval': 50,

                //低端浏览器如设置了domain, 需要传入
                'domain': '',

                //执行navigate前回调
                'beforeNavigate': function (url, cb) {
                    cb();
                }

            };

            option = option || {};

            for (let p in option) {
                this.option[p] = option[p];
            }

            //扩展路由
            if (this.option['extendRoutes']) {
                this.extend(this.option['extendRoutes']);
            }

            this.option['html5Mode'] = (pushState && this.option['html5Mode']);

            //支持debug模式(url加上debug后不改变页面切换逻辑,可有针对性做一些事情)
            this.debug = '';
            let locationHref = window.location.href;
            if (/\/debug_online/.test(locationHref)) {
                this.debug = '/debug_online';
            } else if (/\/debug_https/.test(locationHref)) {
                this.debug = '/debug_https';
            } else if (/\/debug_http/.test(locationHref)) {
                this.debug = '/debug_http';
            } else if (/\/debug/.test(locationHref)) {
                this.debug = '/debug';
            }

            let _self = this,

                evt = this.option['html5Mode'] ? 'popstate' : 'hashchange';

            let start = function () {

                let initPath = _self.getFragment() || '/';

                if (initPath === '/index.html') {
                    initPath = '/';
                }

                //完整路径在hash环境打开则转化为锚点路径后跳转
                if (!_self.option['html5Mode'] && !/#(.*)$/.test(locationHref) && initPath !== '/') {
                    location.replace('/#' + initPath);
                    return;
                }

                _self.navigate(initPath, false, true);
            };

            if (oldIE) {

                //ie8以下创建iframe模拟hashchange
                let iframe = document.createElement('iframe');
                iframe.tabindex = '-1';
                if (this.option['domain']) {
                    iframe.src = 'javascript:void(function(){document.open();' +
                        'document.domain = "' + this.option['domain'] + '";document.close();}());';
                } else {
                    iframe.src = 'javascript:0';
                }
                iframe.style.display = 'none';

                var _iframeOnLoad = function () {
                    iframe.onload = null;
                    iframe.detachEvent('onload', _iframeOnLoad);
                    start();
                    _self.checkUrlInterval = setInterval(function () {
                        _self.checkUrl();
                    }, _self.option['interval']);
                };
                if (iframe.attachEvent) {
                    iframe.attachEvent('onload', _iframeOnLoad);
                } else {
                    iframe.onload = _iframeOnLoad;
                }

                document.body.appendChild(iframe);
                this.iframe = iframe.contentWindow;

            } else {

                //其他浏览器监听popstate或hashchange
                this.addEvent(window, evt, function () {
                    _self.checkUrl();
                });

            }

            if (!this.iframe) {
                start();
            }

        },

        /**
         * 事件监听
         */
        addEvent: function (elem, event, fn) {
            if (elem.addEventListener) {
                elem.addEventListener(event, fn, false);
            } else if (elem.attachEvent) {
                elem.attachEvent('on' + event, fn);
            } else {
                elem[event] = fn;
            }
        },

        /**
         * 获取hash值
         * @method getHash
         * @param {Object} win 窗口对象
         * @return {String} hash值
         * @author evanyuan
         */
        getHash: function (win) {
            let match = (win || window).location.href.match(/#(.*)$/);
            return match ? match[1] : '';
        },

        /**
         * 获取url片段
         * @method getFragment
         * @return {String} url片段
         * @author evanyuan
         */
        getFragment: function () {
            let fragment,
                pathName = window.location.pathname + window.location.search;

            if (this.option['html5Mode']) {
                fragment = pathName;
                //如果锚点路径在html5Mode环境打开
                if (fragment === '/' && this.getHash()) {
                    fragment = this.getHash();
                }
            } else {
                fragment = this.getHash();
                //如果完整路径在hash环境打开
                if (fragment === '' && pathName !== '/' && pathName !== '/index.html') {
                    fragment = pathName;
                }
            }
            return fragment;
        },

        /**
         * 监听url变化
         */
        checkUrl: function () {
            let current = this.getFragment();
            if (this.debug) {
                current = current.replace(this.debug, '');
            }
            if (this.iframe) {
                current = this.getHash(this.iframe);
            }
            if (!this.fragment || current === this.fragment) {
                return;
            }
            if (this.onJustChangeUrl) {
                this.onJustChangeUrl = false;
                return;
            }

            this.navigate(current, false, true);
        },

        /**
         * 去除前后#
         */
        stripHash: function (url) {
            return url.replace(/^\#+|\#+$/g, '');
        },

        /**
         * 去除前后斜杠
         */
        stripSlash: function (url) {
            return url.replace(/^\/+|\/+$/g, '');
        },

        /**
         * 导航
         * @method navigate
         * @param {String}  url 地址
         * @param {Boolean} silent 不改变地址栏
         * @param {Boolean} replacement 替换浏览器的当前会话历史(h5模式时支持)
         * @author evanyuan
         */
        navigate: function (url, silent, replacement) {

            let _self = this;

            if (url !== '/') {
                url = _self.stripHash(url);
                url = _self.stripSlash(url);
                url = '/' + url;
            }

            //执行导航前回调
            this.option.beforeNavigate(url, function () {
                if (url !== _self.fragment && !silent) { //silent为true时，只路由不改变地址栏
                    if (_self.debug) {
                        url = url.replace(_self.debug, '');
                        url = _self.debug + url;
                    }
                    if (_self.option['html5Mode']) {
                        let _stateFun = replacement ? 'replaceState' : 'pushState';
                        history[_stateFun]({}, document.title, url);
                    } else {
                        if (url !== '/' || _self.getFragment()) {
                            location.hash = url;
                            _self.iframe && _self.historySet(url, _self.getHash(_self.iframe));
                        }
                    }
                }

                if (_self.debug) {
                    url = url.replace(_self.debug, '');
                    !url && (url = '/');
                }


                /**
                 * 上一页url片段
                 * @property referrer
                 * @type String
                 */
                _self.referrer = _self.fragment;

                /**
                 * 当前url片段
                 * @property fragment
                 * @type String
                 */
                _self.fragment = url;

                _self.loadUrl(url);
            });

        },

        /**
         * 低端浏览器设置iframe历史
         */
        historySet: function (hash, historyHash) {
            let iframeDoc = this.iframe.document;

            if (hash !== historyHash) {
                iframeDoc.open();
                if (this.option['domain']) {
                    iframeDoc.write('<script>document.domain="' + this.option['domain'] + '"</script>');
                }
                iframeDoc.close();
                this.iframe.location.hash = hash;
            }
        },

        /**
         * 重定向
         * @method redirect
         * @param {String}  url 地址
         * @param {Boolean} silent 不改变地址栏
         * @param {Boolean} replacement 替换浏览器的当前会话历史(h5模式时支持)
         * @author evanyuan
         */
        redirect: function (url, silent, replacement) {
            this.navigate(url, silent, replacement);
        },

        /**
         * 路由匹配
         * @method matchRoute
         * @param  {String} rule 路由规则
         * @param  {String} url 地址
         * @return {Array}  参数数组
         * @author evanyuan
         */
        matchRoute: function (rule, url) {
            if (url === '/') {
                url = '/home';
            }
            let optionalReg = /\((.*?)\)/g,
                paramReg = /(\(\?)?:\w+/g,
                astReg = /\*\w+/g,
                ruleToReg = function (rule) {
                    rule = rule.replace(optionalReg, '(?:$1)?').replace(paramReg, '([^\/]+)').replace(astReg, '(.*?)');
                    return new RegExp('^' + rule + '$');
                },
                route = ruleToReg(rule),
                result = route.exec(url),
                params = null;

            if (result) {
                let args = result.slice(1);
                params = [];
                for (var i = 0, p; p = args[i]; i++) {
                    params.push(p ? decodeURIComponent(p) : '');
                }
            }
            return params;
        },

        /**
         * 扩展路由
         * @method extend
         * @param {Object} obj 路由map
         * @author evanyuan
         */
        extend: function (obj) {
            obj = obj || {};
            if (this.extendRoutes) {
                $.extend(this.extendRoutes, obj);
            } else {
                this.extendRoutes = obj;
            }
        },

        /**
         * 动态使用路由规则
         * @method use
         * @param {string} rule 路由匹配字符串
         * @param {function} action 路由操作
         * @author techirdliu
         * */
        use: function (rule, action) {
            this.dynamicRoutes = this.dynamicRoutes || {};
            this.dynamicRoutes[rule] = action;
        },

        /**
         * 取消动态路由规则的使用
         * @method unuse
         * @param {string} rule 路由的匹配字符串
         * @author techirdliu
         * */
        unuse: function (rule) {
            if (this.dynamicRoutes) {
                this.dynamicRoutes[rule] = undefined;
                try {
                    delete this.dynamicRoutes[rule];
                } catch (e) {}
            }
        },

        /**
         * queryString转对象
         */
        queryToObj: function (queryStr) {
            let urlPara = queryStr.split('?')[1];
            urlPara = urlPara.split('&');
            let objPara = {};
            for (var i = 0, item; item = urlPara[i]; i++) {
                let itemArr = item.split('=');
                objPara[itemArr[0]] = itemArr[1];
            }
            return objPara;
        },

        /**
         * 执行路由匹配的方法
         */
        applyAction: function (action, params, urlParam, pointer) {
            this.applyActionPoint = new Date();
            urlParam && params.push(urlParam);
            action && action.apply(pointer, params);
        },

        getRouteRules: function (routes) {
            // 排序是为了精确匹配在前，如 /monitor/overview 先于 /monitor 匹配
            return Object.keys(routes).sort(function (item1, item2) {
                return item1 < item2 ? 1 : -1;
            });
        },

        /**
         * 加载页面
         * @method loadUrl
         * @param {String} url 地址
         * @author evanyuan
         */
        loadUrl: function (url) {
            var _self = this,
                dynamicRoutes = _self.dynamicRoutes,
                extendRoutes = _self.extendRoutes,
                routes = _self.option.routes,
                pm = _self.option.pageManager,
                params = null,
                urlParam = null,
                searchReg = /\/?\?.*/,
                searchMatch = searchReg.exec(url),
                url = url.replace(searchReg, '');

            if (url == '') {
                url = '/';
            }
            if (_self.option.OnUrlChange) {
                let action = _self.option.OnUrlChange;
                action();
            }
            searchMatch && (urlParam = this.queryToObj(searchMatch[0]));

            let i;

            //优先匹配动态路由
            if (dynamicRoutes) {
                let dyRules = this.getRouteRules(dynamicRoutes);
                for (i = 0; i < dyRules.length; i++) {
                    let dyRule = dyRules[i];
                    if (params = _self.matchRoute(dyRule, url)) {
                        this.curRouteMode = 'dynamic';
                        this.applyAction(dynamicRoutes[dyRule], params, urlParam, null);
                        return;
                    }
                }
            }

            //然后匹配框架外部定义路由
            if (extendRoutes) {
                let exRules = this.getRouteRules(extendRoutes);
                for (i = 0; i < exRules.length; i++) {
                    let exRule = exRules[i];
                    if (params = _self.matchRoute(exRule, url)) {
                        this.curRouteMode = 'external';
                        this.applyAction(extendRoutes[exRule], params, urlParam, null);
                        return;
                    }
                }
            }

            //匹配框架内部路由规则
            for (let rule in routes) {
                if (params = _self.matchRoute(rule, url)) {
                    this.curRouteMode = 'internal';
                    this.applyAction(pm[routes[rule]], params, urlParam, pm);
                    break;
                }
            }
        }

    };

    module.exports = router;

});

define('nmc/widget/tips', function (require, exports, module) {
    let $ = require('$');
    let nmcConfig = require('nmcConfig');
    let util = require('util');

    let flashTime = null, //小黄条延时
        loadingTime = null, //加载器延时
        reqMaxTime = null; //子业务请求时间延时

    /**
     * 提示
     * @class tips
     * @static
     */
    let tips = {
        enableLoading: 1,
        isLoading: 0,
        manualReqNum: 0,

        /**
         * 成功提示
         * @method success
         * @param  {String} str 提示内容
         * @param  {Int}    duration 持续时间
         * @author evanyuan
         */
        success: function (str, duration) {
            this.flash(str, 'success', duration);
        },

        /**
         * 失败或错误提示
         * @method error
         * @param  {String} str 提示内容
         * @param  {Int}    duration 持续时间
         * @author evanyuan
         */
        error: function (str, duration) {
            this.flash(str, 'error', duration);
        },

        /**
         * 操作中提示
         * @method error
         * @param  {String} str 提示内容
         * @param  {Int}    duration 持续时间
         * @author enofan
         */
        loading: function (str, duration) {
            this.flash(str, 'loading', duration);
        },

        /**
         * 小黄条提示
         * @method flash
         * @param  {String} str 提示内容
         * @param  {String} type 类型
         * @param  {Int}    duration 持续时间
         * @author evanyuan
         */
        flash: function (str, type, duration) {
            if (!str) return;
            //避免页面刷新时, 出小黄条错误
            if (type == 'error' && !window.isOnload) return;

            let _self = this,
                icoClass = '';
            let iconMap = {
                'success': 'top-alert-icon-done',
                'error': 'top-alert-icon-waring',
                'loading': 'top-alert-icon-doing'
            };
            !duration && (duration = 4000);
            icoClass = iconMap[type];

            _self.showFlash(str, icoClass);

            flashTime = setTimeout(function () {
                _self.hideFlash(1);
            }, duration);
        },

        /**
         * 显示小黄条
         * @method showFlash
         * @param  {String} str 提示内容
         * @param  {String} icoClass 图标class (top-alert-icon-waring top-alert-icon-done top-alert-icon-doing)
         * @author evanyuan
         */
        showFlash: function (str, icoClassName) {
            clearTimeout(flashTime);
            flashTime = null;

            let topAlert = $('#topAlert');
            let flashMsg = $('#flashMsg');

            if (!flashMsg.length) {
                topAlert = $('<div/>').addClass('top-alert').attr('id', 'topAlert').css({
                    'z-index': 1100,
                    'margin-left': '-200px'
                });
                flashMsg = $('<span/>').attr('id', 'flashMsg');
                topAlert.append(flashMsg);
                topAlert.appendTo('body');
            }
            topAlert.show();
            flashMsg.html(str).show().attr('class', icoClassName + ' fade-in');
        },

        //隐藏小黄条
        hideFlash: function (fromTiming) {
            let _self = this;
            let _hideTip = function () {
                flashTime = null;
                let topAlert = $('#topAlert');
                let flashMsg = $('#flashMsg');

                if (!flashMsg.length || flashMsg.hasClass('fade-out')) {
                    return;
                }

                util.animationend(flashMsg[0], function () {
                    !flashTime && flashMsg.html('').hide();
                    !flashTime && topAlert.hide();
                });

                flashMsg.removeClass('fade-in').addClass('fade-out');

            };
            if (fromTiming) {
                _hideTip();
            } else {
                //从loading过来场景, 有其他提示存在则不隐藏
                !flashTime && _hideTip();
            }
        },

        //立即隐藏小黄条
        hideFlashNow: function () {
            $('#flashMsg').hide().html('');
            $('#topAlert').hide();
        },

        /**
         * 初始化全局加载提示
         * @method initLoading
         * @author evanyuan
         */
        initLoading: function () {
            let _self = this,
                $doc = $(document);
            $doc.ajaxStart(function () {
                _self.ajaxLoading = 1;
                !_self.manualReqNum && _self._loadingStart();
            });
            $doc.ajaxStop(function () {
                _self.ajaxLoading = 0;
                !_self.manualReqNum && _self._loadingStop();
            });
        },

        _loadingStart: function (opt) {
            let _self = this,
                delay = 800;
            if (_self.enableLoading) {
                clearTimeout(loadingTime);
                _self.isLoading = 1;
                loadingTime = setTimeout(function () {
                    let flashMsg = $('#flashMsg');
                    //有其他提示存在, 不提示loading
                    if (flashMsg.length && flashMsg.html()) {
                        return;
                    }
                    if (_self.isLoading) {
                        _self.showLoading(opt);
                    }
                }, delay);
            }
        },

        _loadingStop: function () {
            let _self = this;
            setTimeout(function () {
                if (_self.enableLoading && !_self.manualReqNum && !_self.ajaxLoading) {
                    _self.stopLoading();
                }
            }, 0);
        },

        //显示loading
        showLoading: function (opt) {
            opt = opt || {};
            this.showFlash(opt.text || nmcConfig.loadingWord, opt.className || 'top-alert-icon-doing');
        },

        //停止加载状态
        stopLoading: function () {
            this.isLoading = 0;
            this.hideFlash();
        },

        //手动的请求计数, 用于接入业务jsonp方式等复用loading逻辑
        requestStart: function (opt) {
            opt = opt || {};
            let _self = this,
                delayHide = 5000;
            if (typeof (opt.delayHide) !== 'undefined' && !isNaN(opt.delayHide)) {
                delayHide = parseInt(opt.delayHide, 10);
            }
            this.manualReqNum++;
            if (this.manualReqNum == 1 && !this.ajaxLoading) {
                this._loadingStart(opt);
                clearTimeout(reqMaxTime);
                reqMaxTime = setTimeout(function () {
                    if (_self.manualReqNum > 0) {
                        _self.manualReqNum = 0;
                        _self._loadingStop();
                    }
                }, delayHide);
            }
        },

        requestStop: function () {
            this.manualReqNum--;
            if (this.manualReqNum < 0) {
                this.manualReqNum = 0;
            }
            if (this.manualReqNum == 0 && !this.ajaxLoading) {
                this._loadingStop();
            }
        },

        //禁用启用加载提示
        setLoading: function (_bool) {
            this.enableLoading = _bool;
        }

    };

    module.exports = tips;
});
