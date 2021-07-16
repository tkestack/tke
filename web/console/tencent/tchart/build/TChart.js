module.exports =
/******/ (function(modules) { // webpackBootstrap
/******/ 	// install a JSONP callback for chunk loading
/******/ 	function webpackJsonpCallback(data) {
/******/ 		var chunkIds = data[0];
/******/ 		var moreModules = data[1];
/******/
/******/
/******/ 		// add "moreModules" to the modules object,
/******/ 		// then flag all "chunkIds" as loaded and fire callback
/******/ 		var moduleId, chunkId, i = 0, resolves = [];
/******/ 		for(;i < chunkIds.length; i++) {
/******/ 			chunkId = chunkIds[i];
/******/ 			if(installedChunks[chunkId]) {
/******/ 				resolves.push(installedChunks[chunkId][0]);
/******/ 			}
/******/ 			installedChunks[chunkId] = 0;
/******/ 		}
/******/ 		for(moduleId in moreModules) {
/******/ 			if(Object.prototype.hasOwnProperty.call(moreModules, moduleId)) {
/******/ 				modules[moduleId] = moreModules[moduleId];
/******/ 			}
/******/ 		}
/******/ 		if(parentJsonpFunction) parentJsonpFunction(data);
/******/
/******/ 		while(resolves.length) {
/******/ 			resolves.shift()();
/******/ 		}
/******/
/******/ 	};
/******/
/******/
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// object to store loaded and loading chunks
/******/ 	// undefined = chunk not loaded, null = chunk preloaded/prefetched
/******/ 	// Promise = chunk loading, 0 = chunk loaded
/******/ 	var installedChunks = {
/******/ 		"TChart": 0
/******/ 	};
/******/
/******/
/******/
/******/ 	// script path function
/******/ 	function jsonpScriptSrc(chunkId) {
/******/ 		return __webpack_require__.p + "" + ({"ChartsComponents":"ChartsComponents"}[chunkId]||chunkId) + ".js"
/******/ 	}
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/ 	// This file contains only the entry chunk.
/******/ 	// The chunk loading function for additional chunks
/******/ 	__webpack_require__.e = function requireEnsure(chunkId) {
/******/ 		var promises = [];
/******/
/******/
/******/ 		// JSONP chunk loading for javascript
/******/
/******/ 		var installedChunkData = installedChunks[chunkId];
/******/ 		if(installedChunkData !== 0) { // 0 means "already installed".
/******/
/******/ 			// a Promise means "currently loading".
/******/ 			if(installedChunkData) {
/******/ 				promises.push(installedChunkData[2]);
/******/ 			} else {
/******/ 				// setup Promise in chunk cache
/******/ 				var promise = new Promise(function(resolve, reject) {
/******/ 					installedChunkData = installedChunks[chunkId] = [resolve, reject];
/******/ 				});
/******/ 				promises.push(installedChunkData[2] = promise);
/******/
/******/ 				// start chunk loading
/******/ 				var script = document.createElement('script');
/******/ 				var onScriptComplete;
/******/
/******/ 				script.charset = 'utf-8';
/******/ 				script.timeout = 120;
/******/ 				if (__webpack_require__.nc) {
/******/ 					script.setAttribute("nonce", __webpack_require__.nc);
/******/ 				}
/******/ 				script.src = jsonpScriptSrc(chunkId);
/******/
/******/ 				onScriptComplete = function (event) {
/******/ 					// avoid mem leaks in IE.
/******/ 					script.onerror = script.onload = null;
/******/ 					clearTimeout(timeout);
/******/ 					var chunk = installedChunks[chunkId];
/******/ 					if(chunk !== 0) {
/******/ 						if(chunk) {
/******/ 							var errorType = event && (event.type === 'load' ? 'missing' : event.type);
/******/ 							var realSrc = event && event.target && event.target.src;
/******/ 							var error = new Error('Loading chunk ' + chunkId + ' failed.\n(' + errorType + ': ' + realSrc + ')');
/******/ 							error.type = errorType;
/******/ 							error.request = realSrc;
/******/ 							chunk[1](error);
/******/ 						}
/******/ 						installedChunks[chunkId] = undefined;
/******/ 					}
/******/ 				};
/******/ 				var timeout = setTimeout(function(){
/******/ 					onScriptComplete({ type: 'timeout', target: script });
/******/ 				}, 120000);
/******/ 				script.onerror = script.onload = onScriptComplete;
/******/ 				document.head.appendChild(script);
/******/ 			}
/******/ 		}
/******/ 		return Promise.all(promises);
/******/ 	};
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "//imgcache.qq.com/tchart/build/";
/******/
/******/ 	// on error function for async loading
/******/ 	__webpack_require__.oe = function(err) { console.error(err); throw err; };
/******/
/******/ 	var jsonpArray = (typeof self !== 'undefined' ? self : this)["webpackJsonp"] = (typeof self !== 'undefined' ? self : this)["webpackJsonp"] || [];
/******/ 	var oldJsonpFunction = jsonpArray.push.bind(jsonpArray);
/******/ 	jsonpArray.push = webpackJsonpCallback;
/******/ 	jsonpArray = jsonpArray.slice();
/******/ 	for(var i = 0; i < jsonpArray.length; i++) webpackJsonpCallback(jsonpArray[i]);
/******/ 	var parentJsonpFunction = oldJsonpFunction;
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./src/panel/index.tsx");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./node_modules/@tencent/tea-component/lib/i18n/index.js":
/*!***************************************************************!*\
  !*** ./node_modules/@tencent/tea-component/lib/i18n/index.js ***!
  \***************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var tslib_1 = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
tslib_1.__exportStar(__webpack_require__(/*! ./translation */ "./node_modules/@tencent/tea-component/lib/i18n/translation.js"), exports);
var withTranslation_1 = __webpack_require__(/*! ./withTranslation */ "./node_modules/@tencent/tea-component/lib/i18n/withTranslation.js");
exports.withTranslation = withTranslation_1.withTranslation;


/***/ }),

/***/ "./node_modules/@tencent/tea-component/lib/i18n/locale/en_us.js":
/*!**********************************************************************!*\
  !*** ./node_modules/@tencent/tea-component/lib/i18n/locale/en_us.js ***!
  \**********************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var tslib_1 = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
var react_1 = tslib_1.__importStar(__webpack_require__(/*! react */ "./webpack/alias/react.js"));
// prettier-ignore
// eslint-disable-next-line
exports.en_us = {
    locale: 'en_us',
    // 操作：确定按钮
    okText: ('OK'),
    // 操作：取消按钮
    cancelText: ('Cancel'),
    // 文案：数据加载中
    loadingText: ('Loading…'),
    // 文案：数据加载失败
    loadErrorText: ('Loading failed'),
    // 操作：重试数据加载
    loadRetryText: ('Retry'),
    // 操作：关闭
    closeText: ('Close'),
    // 文案：帮助
    helpText: ('Help'),
    // 操作：清空
    cleanText: ('Clear'),
    // 操作：搜索
    searchText: ('Search'),
    // 文案：数据为空
    emptyText: ('No data yet'),
    // 操作：全选
    selectAllText: ('Select all'),
    // 文案：分页组件显示总记录数
    paginationRecordCount: function (count) { return (react_1.default.createElement(react_1.Fragment, null,
        "Total items: ",
        react_1.default.createElement("strong", null, count),
        " ")); },
    // 操作：跳到上一页
    paginationPrevPage: ('Previous'),
    // 操作：跳到下一页
    paginationNextPage: ('Next'),
    // 操作：跳到第一页
    paginationToFirstPage: ('First page'),
    // 操作：跳到最后一页
    paginationToLastPage: ('Last page'),
    // 文案：提醒用户当前已经在第一页，无法再跳到上一页
    paginationAtFirst: ('This is the first page'),
    // 文案：提醒用户当前已经在最后一页，无法再跳到上一页
    paginationAtLast: ('This is the last page'),
    // 文案：表示分页组件每页显示多少行记录，后接行数选项
    paginationRecordPerPage: ('Records per page'),
    // 文案：表示分页组件总共有多少页，前面是当前的页码
    paginationPageCount: function (count) { return (count > 1 ? react_1.default.createElement(react_1.Fragment, null,
        " / ",
        count,
        " pages") : react_1.default.createElement(react_1.Fragment, null,
        " / ",
        count,
        " page")); },
    // 文案：下拉选择组件默认的提示文案
    pleaseSelect: ('Please select'),
    // 文案：查询到结果
    foundText: ('Found the following results'),
    // 文案：表格中，用于显示找到多少条结果，后面会拼接「返回原列表」
    foundManyText: function (count) { return (count > 1 ? count + " results found" : count + " result found"); },
    // 文案：同 resultText，不过是在有关键字的情况下显示
    foundManyTextWithKeyword: function (keyword, count) { return (count > 1 ? count + " results found for \"" + keyword + "\"" : count + " result found for \"" + keyword + "\""); },
    // 文案：搜索某个关键字的情况下，没有找到结果
    foundNothingWithKeyword: function (keyword) { return ("No results found for \"" + keyword + "\""); },
    // 操作：表格中清空当前筛选结果，返回源列表
    clearResultText: ('Back to list'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxTips: ('Separate keywords with "|"; press Enter to separate filter tags'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxEditingTips: ('Click to modify. Press Enter to finish.'),
    // 文案：tagSearchBox 选择框标题
    tagSearchBoxSelectTitle: ('Select a filter'),
    // 文案：tagSearchBox 帮助图片地址
    tagSearchBoxHelpImgUrl: ('//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/server/css/img/search-dialog-en.png'),
    // 文案：今天
    today: ('Today'),
    // 文案：本月
    curMonth: ('This month'),
    // 文案：下个月
    prevMonth: ('Previous month'),
    // 文案：下个月
    nextMonth: ('Next month'),
    // 文案：今年
    curYear: ('This year'),
    // 文案：下一年
    prevYear: ('Previous year'),
    // 文案：下一年
    nextYear: ('Next year'),
    // 文案：当前二十年
    curTwentyYears: ('Latest 20 years'),
    // 文案：上二十年
    prevTwentyYears: ('Previous 20 years'),
    // 文案：下二十年
    nextTwentyYears: ('Next 20 years'),
    // 变量：该语言日期表达中 [月] 是否在 [年] 之前，是的话为 - true，否则 - false
    monthBeforeYear: (true),
    // 变量：该语言日期表达中 [年] 的表达方式，其中 YYYY 为年份数字
    yearFormat: (' YYYY'),
    // 文案：选择时间
    selectTime: ('Select a time'),
    // 文案：选择日期
    selectDate: ('Select a date'),
    // 文案：输入格式错误
    invalidFormat: function (format) { return ("Invalid format. The format should be " + format); },
    // 文案：时间不合法
    invalidTime: ("Invalid time"),
};


/***/ }),

/***/ "./node_modules/@tencent/tea-component/lib/i18n/locale/ja.js":
/*!*******************************************************************!*\
  !*** ./node_modules/@tencent/tea-component/lib/i18n/locale/ja.js ***!
  \*******************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var tslib_1 = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
var react_1 = tslib_1.__importStar(__webpack_require__(/*! react */ "./webpack/alias/react.js"));
// prettier-ignore
// eslint-disable-next-line
exports.ja = {
    locale: 'ja',
    // 操作：确定按钮
    okText: ('OK'),
    // 操作：取消按钮
    cancelText: ('キャンセル'),
    // 文案：数据加载中
    loadingText: ('読み込み中...'),
    // 文案：数据加载失败
    loadErrorText: ('読み込み失敗'),
    // 操作：重试数据加载
    loadRetryText: ('再試行'),
    // 操作：关闭
    closeText: ('無効化'),
    // 文案：帮助
    helpText: ('ヘルプ'),
    // 操作：清空
    cleanText: ('クリア'),
    // 操作：搜索
    searchText: ('検索'),
    // 文案：数据为空
    emptyText: ('データなし'),
    // 操作：全选
    selectAllText: ('すべて選択'),
    // 文案：分页组件显示总记录数
    paginationRecordCount: function (count) { return (react_1.default.createElement(react_1.Fragment, null,
        "\u5408\u8A08 ",
        react_1.default.createElement("strong", null, count),
        " \u9805\u76EE")); },
    // 操作：跳到上一页
    paginationPrevPage: ('前のページ'),
    // 操作：跳到下一页
    paginationNextPage: ('次のページ'),
    // 操作：跳到第一页
    paginationToFirstPage: ('最初ページ'),
    // 操作：跳到最后一页
    paginationToLastPage: ('最終ページ'),
    // 文案：提醒用户当前已经在第一页，无法再跳到上一页
    paginationAtFirst: ('既に最初ページです'),
    // 文案：提醒用户当前已经在最后一页，无法再跳到上一页
    paginationAtLast: ('既に最終ページです'),
    // 文案：表示分页组件每页显示多少行记录，后接行数选项
    paginationRecordPerPage: ('各ページの表示行数'),
    // 文案：表示分页组件总共有多少页，前面是当前的页码
    paginationPageCount: function (count) { return (react_1.default.createElement(react_1.Fragment, null,
        " / ",
        count,
        " \u30DA\u30FC\u30B8")); },
    // 文案：下拉选择组件默认的提示文案
    pleaseSelect: ('選択してください'),
    // 文案：查询到结果
    foundText: ("下記の結果が見つかりました"),
    // 文案：表格中，用于显示找到多少条结果，后面会拼接「返回原列表」
    foundManyText: function (count) { return (count + " \u4EF6\u306E\u7D50\u679C\u304C\u898B\u3064\u304B\u308A\u307E\u3057\u305F"); },
    // 文案：同 resultText，不过是在有关键字的情况下显示
    foundManyTextWithKeyword: function (keyword, count) { return ("\"" + keyword + "\"\u3092\u691C\u7D22\u3057\u307E\u3059\u3002" + count + " \u4EF6\u306E\u7D50\u679C\u304C\u898B\u3064\u304B\u308A\u307E\u3057\u305F"); },
    // 文案：搜索某个关键字的情况下，没有找到结果
    foundNothingWithKeyword: function (keyword) { return ("\"" + keyword + "\"\u3092\u691C\u7D22\u3057\u307E\u3059\u3002\u7D50\u679C\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3067\u3057\u305F"); },
    // 操作：表格中清空当前筛选结果，返回源列表
    clearResultText: ('元のリストに戻る'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxTips: ('複数のキーワードは縦棒"|"で区切られ、複数のフィルタタグはEnterキーで区切られます'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxEditingTips: ('クリックで変更し、Enterキーで変更を完了します'),
    // 文案：tagSearchBox 选择框标题
    tagSearchBoxSelectTitle: ('フィルタするリソースプロパティを選択'),
    // 文案：tagSearchBox 帮助图片地址
    tagSearchBoxHelpImgUrl: ('//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/server/css/img/search-dialog-en.png'),
    // 文案：今天
    today: ('本日'),
    // 文案：本月
    curMonth: ('今月'),
    // 文案：下个月
    prevMonth: ('先月'),
    // 文案：下个月
    nextMonth: ('来月'),
    // 文案：今年
    curYear: ('今年'),
    // 文案：下一年
    prevYear: ('前年'),
    // 文案：下一年
    nextYear: ('来年'),
    // 文案：当前二十年
    curTwentyYears: ('現在の20年'),
    // 文案：上二十年
    prevTwentyYears: ('この前の20年'),
    // 文案：下二十年
    nextTwentyYears: ('この後の20年'),
    // 变量：该语言日期表达中 [月] 是否在 [年] 之前，是的话为 - true，否则 - false
    monthBeforeYear: (false),
    // 变量：该语言日期表达中 [年] 的表达方式，其中 YYYY 为年份数字
    yearFormat: ('YYYY年 '),
    // 文案：选择时间
    selectTime: ('時間の選択'),
    // 文案：选择日期
    selectDate: ('日付の選択'),
    // 文案：输入格式错误
    invalidFormat: function (format) { return ("\u30D5\u30A9\u30FC\u30DE\u30C3\u30C8\u30A8\u30E9\u30FC\u3001" + format + " \u306B\u3057\u3066\u304F\u3060\u3055\u3044"); },
    // 文案：时间不合法
    invalidTime: ("\u3053\u306E\u6642\u9593\u304C\u4E0D\u6B63\u3067\u3059"),
};


/***/ }),

/***/ "./node_modules/@tencent/tea-component/lib/i18n/locale/ko.js":
/*!*******************************************************************!*\
  !*** ./node_modules/@tencent/tea-component/lib/i18n/locale/ko.js ***!
  \*******************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var tslib_1 = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
var react_1 = tslib_1.__importStar(__webpack_require__(/*! react */ "./webpack/alias/react.js"));
// prettier-ignore
// eslint-disable-next-line
exports.ko = {
    locale: 'ko',
    // 操作：确定按钮
    okText: ('확인'),
    // 操作：取消按钮
    cancelText: ('취소'),
    // 文案：数据加载中
    loadingText: ('로딩 중...'),
    // 文案：数据加载失败
    loadErrorText: ('로딩 실패'),
    // 操作：重试数据加载
    loadRetryText: ('다시 시도'),
    // 操作：关闭
    closeText: ('종료'),
    // 文案：帮助
    helpText: ('도움말'),
    // 操作：清空
    cleanText: ('비우기'),
    // 操作：搜索
    searchText: ('검색'),
    // 文案：数据为空
    emptyText: ('데이터가 없습니다.'),
    // 操作：全选
    selectAllText: ('전체 선택'),
    // 文案：分页组件显示总记录数
    paginationRecordCount: function (count) { return (react_1.default.createElement(react_1.Fragment, null,
        "\uCD1D ",
        react_1.default.createElement("strong", null, count),
        "\uAC1C")); },
    // 操作：跳到上一页
    paginationPrevPage: ('이전 페이지'),
    // 操作：跳到下一页
    paginationNextPage: ('다음 페이지'),
    // 操作：跳到第一页
    paginationToFirstPage: ('첫 페이지'),
    // 操作：跳到最后一页
    paginationToLastPage: ('마지막 페이지'),
    // 文案：提醒用户当前已经在第一页，无法再跳到上一页
    paginationAtFirst: ('첫 페이지입니다'),
    // 文案：提醒用户当前已经在最后一页，无法再跳到上一页
    paginationAtLast: ('마지막 페이지입니다'),
    // 文案：表示分页组件每页显示多少行记录，后接行数选项
    paginationRecordPerPage: ('페이지당 표시 행수'),
    // 文案：表示分页组件总共有多少页，前面是当前的页码
    paginationPageCount: function (count) { return (react_1.default.createElement(react_1.Fragment, null,
        "/ ",
        count,
        "\uD398\uC774\uC9C0")); },
    // 文案：下拉选择组件默认的提示文案
    pleaseSelect: ('선택하십시오'),
    // 文案：查询到结果
    foundText: ("다음 결과를 찾았습니다"),
    // 文案：表格中，用于显示找到多少条结果，后面会拼接「返回原列表」
    foundManyText: function (count) { return (count + "\uAC1C \uACB0\uACFC\uB97C \uCC3E\uC558\uC2B5\uB2C8\uB2E4"); },
    // 文案：同 resultText，不过是在有关键字的情况下显示
    foundManyTextWithKeyword: function (keyword, count) { return ("\"" + keyword + "\"\uC744(\uB97C) \uAC80\uC0C9\uD558\uC5EC " + count + "\uAC1C \uACB0\uACFC\uB97C \uCC3E\uC558\uC2B5\uB2C8\uB2E4"); },
    // 文案：搜索某个关键字的情况下，没有找到结果
    foundNothingWithKeyword: function (keyword) { return ("\"" + keyword + "\"\uC744(\uB97C) \uAC80\uC0C9\uD558\uC5EC \uCC3E\uC740 \uACB0\uACFC\uAC00 \uC5C6\uC2B5\uB2C8\uB2E4"); },
    // 操作：表格中清空当前筛选结果，返回源列表
    clearResultText: ('기존 리스트로 돌아가기'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxTips: ('여러 개의 키워드는 "|"으로 구분되며 여러 개의 필터 태그는 Enter 키로 구분됩니다.'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxEditingTips: ('클릭하여 수정합니다. Enter 키를 눌러 수정을 완료합니다.'),
    // 文案：tagSearchBox 选择框标题
    tagSearchBoxSelectTitle: ('리소스 속성을 선택하여 필터링합니다.'),
    // 文案：tagSearchBox 帮助图片地址
    tagSearchBoxHelpImgUrl: ('//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/server/css/img/search-dialog-en.png'),
    // 文案：今天
    today: ('오늘'),
    // 文案：本月
    curMonth: ('이번 달'),
    // 文案：下个月
    prevMonth: ('지난 달'),
    // 文案：下个月
    nextMonth: ('다음 달'),
    // 文案：今年
    curYear: ('금해'),
    // 文案：下一年
    prevYear: ('지난 해'),
    // 文案：下一年
    nextYear: ('다음 해'),
    // 文案：当前二十年
    curTwentyYears: ('현재 20년'),
    // 文案：上二十年
    prevTwentyYears: ('지난 20년'),
    // 文案：下二十年
    nextTwentyYears: ('다음 20년'),
    // 变量：该语言日期表达中 [月] 是否在 [年] 之前，是的话为 - true，否则 - false
    monthBeforeYear: (false),
    // 变量：该语言日期表达中 [年] 的表达方式，其中 YYYY 为年份数字
    yearFormat: ('YYYY년 '),
    // 文案：选择时间
    selectTime: ('시간 선택'),
    // 文案：选择日期
    selectDate: ('날짜 선택'),
    // 文案：输入格式错误
    invalidFormat: function (format) { return ("\uD615\uC2DD \uC624\uB958, \uC62C\uBC14\uB978 \uD615\uC2DD: " + format); },
    // 文案：时间不合法
    invalidTime: ("\uD574\uB2F9 \uC2DC\uAC04\uC740 \uC798\uBABB\uB429\uB2C8\uB2E4"),
};


/***/ }),

/***/ "./node_modules/@tencent/tea-component/lib/i18n/locale/zh_cn.js":
/*!**********************************************************************!*\
  !*** ./node_modules/@tencent/tea-component/lib/i18n/locale/zh_cn.js ***!
  \**********************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var tslib_1 = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
var react_1 = tslib_1.__importStar(__webpack_require__(/*! react */ "./webpack/alias/react.js"));
// prettier-ignore
// eslint-disable-next-line
exports.zh_cn = {
    locale: 'zh_cn',
    // 操作：确定按钮
    okText: ('确定'),
    // 操作：取消按钮
    cancelText: ('取消'),
    // 文案：数据加载中
    loadingText: ('加载中...'),
    // 文案：数据加载失败
    loadErrorText: ('加载失败'),
    // 操作：重试数据加载
    loadRetryText: ('重试'),
    // 操作：关闭
    closeText: ('关闭'),
    // 文案：帮助
    helpText: ('帮助'),
    // 操作：清空
    cleanText: ('清空'),
    // 操作：搜索
    searchText: ('搜索'),
    // 文案：数据为空
    emptyText: ('暂无数据'),
    // 操作：全选
    selectAllText: ('全选'),
    // 文案：分页组件显示总记录数
    paginationRecordCount: function (count) { return (react_1.default.createElement(react_1.Fragment, null,
        "\u5171 ",
        react_1.default.createElement("strong", null, count),
        " \u9879")); },
    // 操作：跳到上一页
    paginationPrevPage: ('上一页'),
    // 操作：跳到下一页
    paginationNextPage: ('下一页'),
    // 操作：跳到第一页
    paginationToFirstPage: ('第一页'),
    // 操作：跳到最后一页
    paginationToLastPage: ('最后一页'),
    // 文案：提醒用户当前已经在第一页，无法再跳到上一页
    paginationAtFirst: ('当前已在第一页'),
    // 文案：提醒用户当前已经在最后一页，无法再跳到上一页
    paginationAtLast: ('当前已在最后一页'),
    // 文案：表示分页组件每页显示多少行记录，后接行数选项
    paginationRecordPerPage: ('每页显示行'),
    // 文案：表示分页组件总共有多少页，前面是当前的页码
    paginationPageCount: function (count) { return (react_1.default.createElement(react_1.Fragment, null,
        " / ",
        count,
        " \u9875")); },
    // 文案：下拉选择组件默认的提示文案
    pleaseSelect: ('请选择'),
    // 文案：查询到结果
    foundText: ("找到下列结果"),
    // 文案：表格中，用于显示找到多少条结果，后面会拼接「返回原列表」
    foundManyText: function (count) { return ("\u627E\u5230 " + count + " \u6761\u7ED3\u679C"); },
    // 文案：同 resultText，不过是在有关键字的情况下显示
    foundManyTextWithKeyword: function (keyword, count) { return ("\u641C\u7D22 \u201C" + keyword + "\u201D\uFF0C\u627E\u5230 " + count + " \u6761\u7ED3\u679C"); },
    // 文案：搜索某个关键字的情况下，没有找到结果
    foundNothingWithKeyword: function (keyword) { return ("\u641C\u7D22 \u201C" + keyword + "\u201D \u65E0\u7ED3\u679C"); },
    // 操作：表格中清空当前筛选结果，返回源列表
    clearResultText: ('返回原列表'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxTips: ('多个关键字用竖线 “|” 分隔，多个过滤标签用回车键分隔'),
    // 文案：tagSearchBox 使用提示
    tagSearchBoxEditingTips: ('点击进行修改，按回车键完成修改'),
    // 文案：tagSearchBox 选择框标题
    tagSearchBoxSelectTitle: ('选择资源属性进行过滤'),
    // 文案：tagSearchBox 帮助图片地址
    tagSearchBoxHelpImgUrl: ('//imgcache.qq.com/open_proj/proj_qcloud_v2/bee/css/img/search-dialog.png'),
    // 文案：今天
    today: ('今天'),
    // 文案：本月
    curMonth: ('本月'),
    // 文案：下个月
    prevMonth: ('上个月'),
    // 文案：下个月
    nextMonth: ('下个月'),
    // 文案：今年
    curYear: ('今年'),
    // 文案：下一年
    prevYear: ('上一年'),
    // 文案：下一年
    nextYear: ('下一年'),
    // 文案：当前二十年
    curTwentyYears: ('当前二十年'),
    // 文案：上二十年
    prevTwentyYears: ('上二十年'),
    // 文案：下二十年
    nextTwentyYears: ('下二十年'),
    // 变量：该语言日期表达中 [月] 是否在 [年] 之前，是的话为 - true，否则 - false
    monthBeforeYear: (false),
    // 变量：该语言日期表达中 [年] 的表达方式，其中 YYYY 为年份数字
    yearFormat: ('YYYY年 '),
    // 文案：选择时间
    selectTime: ('选择时间'),
    // 文案：选择日期
    selectDate: ('选择日期'),
    // 文案：输入格式错误
    invalidFormat: function (format) { return ("\u683C\u5F0F\u9519\u8BEF\uFF0C\u5E94\u4E3A " + format); },
    // 文案：时间不合法
    invalidTime: ("\u8BE5\u65F6\u95F4\u4E0D\u5408\u6CD5"),
};


/***/ }),

/***/ "./node_modules/@tencent/tea-component/lib/i18n/translation.js":
/*!*********************************************************************!*\
  !*** ./node_modules/@tencent/tea-component/lib/i18n/translation.js ***!
  \*********************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var tslib_1 = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
/* eslint-disable @typescript-eslint/camelcase */
var react_1 = __webpack_require__(/*! react */ "./webpack/alias/react.js");
var zh_cn_1 = __webpack_require__(/*! ./locale/zh_cn */ "./node_modules/@tencent/tea-component/lib/i18n/locale/zh_cn.js");
var en_us_1 = __webpack_require__(/*! ./locale/en_us */ "./node_modules/@tencent/tea-component/lib/i18n/locale/en_us.js");
var ja_1 = __webpack_require__(/*! ./locale/ja */ "./node_modules/@tencent/tea-component/lib/i18n/locale/ja.js");
var ko_1 = __webpack_require__(/*! ./locale/ko */ "./node_modules/@tencent/tea-component/lib/i18n/locale/ko.js");
// 约定中语言名不规范
var lngs = {
    zh: "zh_cn",
    en: "en_us",
    jp: "ja",
};
/* eslint-disable @typescript-eslint/no-object-literal-type-assertion */
var translationMap = {
    zh_cn: zh_cn_1.zh_cn,
    zh: zh_cn_1.zh_cn,
    en_us: tslib_1.__assign({}, zh_cn_1.zh_cn, en_us_1.en_us),
    en: tslib_1.__assign({}, zh_cn_1.zh_cn, en_us_1.en_us),
    ja: tslib_1.__assign({}, zh_cn_1.zh_cn, en_us_1.en_us, ja_1.ja),
    jp: tslib_1.__assign({}, zh_cn_1.zh_cn, en_us_1.en_us, ja_1.ja),
    ko: tslib_1.__assign({}, zh_cn_1.zh_cn, en_us_1.en_us, ko_1.ko),
};
/* eslint-enable @typescript-eslint/no-object-literal-type-assertion */
var currentTranslation = zh_cn_1.zh_cn;
function setLocale(locale, moment) {
    if (moment) {
        moment.locale(lngs[locale] || locale);
    }
    currentTranslation = translationMap[locale];
}
exports.setLocale = setLocale;
function useTranslation(moment) {
    var t = currentTranslation || zh_cn_1.zh_cn;
    var locale = lngs[t.locale] || t.locale;
    react_1.useState(function () {
        if (moment) {
            moment.locale(locale);
        }
    });
    react_1.useEffect(function () {
        if (moment) {
            moment.locale(locale);
        }
    }, [locale]); // eslint-disable-line react-hooks/exhaustive-deps
    return t;
}
exports.useTranslation = useTranslation;
function getTranslation() {
    return currentTranslation || zh_cn_1.zh_cn;
}
exports.getTranslation = getTranslation;
/* eslint-enable @typescript-eslint/camelcase */


/***/ }),

/***/ "./node_modules/@tencent/tea-component/lib/i18n/withTranslation.js":
/*!*************************************************************************!*\
  !*** ./node_modules/@tencent/tea-component/lib/i18n/withTranslation.js ***!
  \*************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
var tslib_1 = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
var react_1 = tslib_1.__importDefault(__webpack_require__(/*! react */ "./webpack/alias/react.js"));
var translation_1 = __webpack_require__(/*! ./translation */ "./node_modules/@tencent/tea-component/lib/i18n/translation.js");
function withTranslation(WrappedComponent) {
    return react_1.default.forwardRef(function (props, ref) {
        var t = translation_1.useTranslation();
        return react_1.default.createElement(WrappedComponent, tslib_1.__assign({}, props, { ref: ref, t: t }));
    });
}
exports.withTranslation = withTranslation;


/***/ }),

/***/ "./node_modules/tslib/tslib.es6.js":
/*!*****************************************!*\
  !*** ./node_modules/tslib/tslib.es6.js ***!
  \*****************************************/
/*! exports provided: __extends, __assign, __rest, __decorate, __param, __metadata, __awaiter, __generator, __exportStar, __values, __read, __spread, __spreadArrays, __await, __asyncGenerator, __asyncDelegator, __asyncValues, __makeTemplateObject, __importStar, __importDefault */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__extends", function() { return __extends; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__assign", function() { return __assign; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__rest", function() { return __rest; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__decorate", function() { return __decorate; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__param", function() { return __param; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__metadata", function() { return __metadata; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__awaiter", function() { return __awaiter; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__generator", function() { return __generator; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__exportStar", function() { return __exportStar; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__values", function() { return __values; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__read", function() { return __read; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__spread", function() { return __spread; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__spreadArrays", function() { return __spreadArrays; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__await", function() { return __await; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__asyncGenerator", function() { return __asyncGenerator; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__asyncDelegator", function() { return __asyncDelegator; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__asyncValues", function() { return __asyncValues; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__makeTemplateObject", function() { return __makeTemplateObject; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__importStar", function() { return __importStar; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "__importDefault", function() { return __importDefault; });
/*! *****************************************************************************
Copyright (c) Microsoft Corporation. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

THIS CODE IS PROVIDED ON AN *AS IS* BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, EITHER EXPRESS OR IMPLIED, INCLUDING WITHOUT LIMITATION ANY IMPLIED
WARRANTIES OR CONDITIONS OF TITLE, FITNESS FOR A PARTICULAR PURPOSE,
MERCHANTABLITY OR NON-INFRINGEMENT.

See the Apache Version 2.0 License for specific language governing permissions
and limitations under the License.
***************************************************************************** */
/* global Reflect, Promise */

var extendStatics = function(d, b) {
    extendStatics = Object.setPrototypeOf ||
        ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
        function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
    return extendStatics(d, b);
};

function __extends(d, b) {
    extendStatics(d, b);
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
}

var __assign = function() {
    __assign = Object.assign || function __assign(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    }
    return __assign.apply(this, arguments);
}

function __rest(s, e) {
    var t = {};
    for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
        for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
            if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
                t[p[i]] = s[p[i]];
        }
    return t;
}

function __decorate(decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
}

function __param(paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
}

function __metadata(metadataKey, metadataValue) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(metadataKey, metadataValue);
}

function __awaiter(thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
}

function __generator(thisArg, body) {
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
}

function __exportStar(m, exports) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}

function __values(o) {
    var m = typeof Symbol === "function" && o[Symbol.iterator], i = 0;
    if (m) return m.call(o);
    return {
        next: function () {
            if (o && i >= o.length) o = void 0;
            return { value: o && o[i++], done: !o };
        }
    };
}

function __read(o, n) {
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
}

function __spread() {
    for (var ar = [], i = 0; i < arguments.length; i++)
        ar = ar.concat(__read(arguments[i]));
    return ar;
}

function __spreadArrays() {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};

function __await(v) {
    return this instanceof __await ? (this.v = v, this) : new __await(v);
}

function __asyncGenerator(thisArg, _arguments, generator) {
    if (!Symbol.asyncIterator) throw new TypeError("Symbol.asyncIterator is not defined.");
    var g = generator.apply(thisArg, _arguments || []), i, q = [];
    return i = {}, verb("next"), verb("throw"), verb("return"), i[Symbol.asyncIterator] = function () { return this; }, i;
    function verb(n) { if (g[n]) i[n] = function (v) { return new Promise(function (a, b) { q.push([n, v, a, b]) > 1 || resume(n, v); }); }; }
    function resume(n, v) { try { step(g[n](v)); } catch (e) { settle(q[0][3], e); } }
    function step(r) { r.value instanceof __await ? Promise.resolve(r.value.v).then(fulfill, reject) : settle(q[0][2], r); }
    function fulfill(value) { resume("next", value); }
    function reject(value) { resume("throw", value); }
    function settle(f, v) { if (f(v), q.shift(), q.length) resume(q[0][0], q[0][1]); }
}

function __asyncDelegator(o) {
    var i, p;
    return i = {}, verb("next"), verb("throw", function (e) { throw e; }), verb("return"), i[Symbol.iterator] = function () { return this; }, i;
    function verb(n, f) { i[n] = o[n] ? function (v) { return (p = !p) ? { value: __await(o[n](v)), done: n === "return" } : f ? f(v) : v; } : f; }
}

function __asyncValues(o) {
    if (!Symbol.asyncIterator) throw new TypeError("Symbol.asyncIterator is not defined.");
    var m = o[Symbol.asyncIterator], i;
    return m ? m.call(o) : (o = typeof __values === "function" ? __values(o) : o[Symbol.iterator](), i = {}, verb("next"), verb("throw"), verb("return"), i[Symbol.asyncIterator] = function () { return this; }, i);
    function verb(n) { i[n] = o[n] && function (v) { return new Promise(function (resolve, reject) { v = o[n](v), settle(resolve, reject, v.done, v.value); }); }; }
    function settle(resolve, reject, d, v) { Promise.resolve(v).then(function(v) { resolve({ value: v, done: d }); }, reject); }
}

function __makeTemplateObject(cooked, raw) {
    if (Object.defineProperty) { Object.defineProperty(cooked, "raw", { value: raw }); } else { cooked.raw = raw; }
    return cooked;
};

function __importStar(mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) result[k] = mod[k];
    result.default = mod;
    return result;
}

function __importDefault(mod) {
    return (mod && mod.__esModule) ? mod : { default: mod };
}


/***/ }),

/***/ "./src/charts/area.ts":
/*!****************************!*\
  !*** ./src/charts/area.ts ***!
  \****************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return AreaChart; });
/* harmony import */ var core_graph__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! core/graph */ "./src/core/graph.ts");

/**
 * 只有叠加状态
 */
class AreaChart extends core_graph__WEBPACK_IMPORTED_MODULE_0__["default"] {
    constructor(id, options) {
        super(id, Object.assign(Object.assign({}, options), { overlay: true, isSeries: true }));
    }
    setOptions(options) {
        super.setOptions(Object.assign(Object.assign({}, options), { overlay: true, isSeries: true }));
    }
    draw() {
        this._mainPanel.clearRect(0, 0, this._options.width, this._options.height);
        this.drawLoading();
        this.drawTitle();
        // 绘画坐标轴
        this.drawAxis();
        // 绘画等高线
        this.drawGrid();
        this.drawSeriesLabels();
        this.drawLegends();
        this.drawEmptyData();
        // 绘画数据折线
        this.drawAreaOnlyShowTopLine();
    }
}


/***/ }),

/***/ "./src/charts/bar.ts":
/*!***************************!*\
  !*** ./src/charts/bar.ts ***!
  \***************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return BarChart; });
/* harmony import */ var core_graph__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! core/graph */ "./src/core/graph.ts");

/**
 * 只有叠加状态
 */
class BarChart extends core_graph__WEBPACK_IMPORTED_MODULE_0__["default"] {
    constructor(id, options) {
        super(id, Object.assign(Object.assign({}, options), { overlay: true }));
    }
    setOptions(options) {
        super.setOptions(Object.assign(Object.assign({}, options), { overlay: true }));
    }
    draw() {
        this._mainPanel.clearRect(0, 0, this._options.width, this._options.height);
        this.drawLoading();
        this.drawTitle();
        // 绘画坐标轴
        this.drawAxis();
        // 绘画等高线
        this.drawGrid();
        this.drawLabels();
        this.drawLegends();
        this.drawEmptyData();
        this.drawBar();
    }
}


/***/ }),

/***/ "./src/charts/index.ts":
/*!*****************************!*\
  !*** ./src/charts/index.ts ***!
  \*****************************/
/*! exports provided: ColorTypes, ChartType, default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ColorTypes", function() { return ColorTypes; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChartType", function() { return ChartType; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return Chart; });
/* harmony import */ var charts_line__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! charts/line */ "./src/charts/line.ts");
/* harmony import */ var charts_area__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! charts/area */ "./src/charts/area.ts");
/* harmony import */ var charts_bar__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! charts/bar */ "./src/charts/bar.ts");
/* harmony import */ var charts_series__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! charts/series */ "./src/charts/series.ts");
/* harmony import */ var core_theme__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! core/theme */ "./src/core/theme.ts");





const ColorTypes = core_theme__WEBPACK_IMPORTED_MODULE_4__["COLORS"].Types;
/**
 * 图表类型
 */
const ChartType = {
    Line: "line",
    Area: "area",
    Bar: "bar",
    Series: "series"
};
class Chart {
    constructor(id, options, chartType = ChartType.Line) {
        this._options = {};
        this._id = id;
        this.setType(chartType, options);
    }
    get ID() {
        return this._id;
    }
    setType(chartType, options) {
        let _class = null;
        switch (chartType) {
            case ChartType.Line:
                _class = charts_line__WEBPACK_IMPORTED_MODULE_0__["default"];
                break;
            case ChartType.Area:
                _class = charts_area__WEBPACK_IMPORTED_MODULE_1__["default"];
                break;
            case ChartType.Bar:
                _class = charts_bar__WEBPACK_IMPORTED_MODULE_2__["default"];
                break;
            case ChartType.Series:
                _class = charts_series__WEBPACK_IMPORTED_MODULE_3__["default"];
                break;
            default:
                _class = charts_series__WEBPACK_IMPORTED_MODULE_3__["default"];
        }
        // 判断是否需要改变类型
        if (!this._instance || !(this._instance instanceof _class) || !this._instance._parentContainer) {
            this._options = Object.assign(Object.assign({}, this._options), options);
            this._instance = new _class(this._id, this._options);
            this.draw();
        }
        else if (this._instance._parentContainer) {
            this.setOptions(options);
        }
    }
    setOptions(options) {
        this._options = Object.assign(Object.assign({}, this._options), options);
        this._instance.setOptions(options);
        this.draw();
    }
    draw() {
        if (this._instance._parentContainer) {
            requestAnimationFrame(() => {
                this._instance.draw();
            });
        }
    }
    highlightLine(legend) {
        this._instance && this._instance.highlightLine(legend);
    }
    setSize(width, height) {
        this._instance && this._instance.setChartSize(width, height);
    }
}


/***/ }),

/***/ "./src/charts/line.ts":
/*!****************************!*\
  !*** ./src/charts/line.ts ***!
  \****************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return LineChart; });
/* harmony import */ var core_graph__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! core/graph */ "./src/core/graph.ts");

/**
 * activeHover: 支持鼠标悬浮时高亮折线
 */
class LineChart extends core_graph__WEBPACK_IMPORTED_MODULE_0__["default"] {
    constructor(id, options) {
        super(id, Object.assign({ activeHover: true }, options));
    }
    setOptions(options) {
        super.setOptions(Object.assign({}, options));
    }
    draw() {
        this._mainPanel.clearRect(0, 0, this._options.width, this._options.height);
        this.drawLoading();
        this.drawTitle();
        // 绘画坐标轴
        this.drawAxis();
        // 绘画等高线
        this.drawGrid();
        this.drawLabels();
        this.drawLegends();
        this.drawEmptyData();
        // 绘画数据折线
        this.drawLine();
    }
}


/***/ }),

/***/ "./src/charts/series.ts":
/*!******************************!*\
  !*** ./src/charts/series.ts ***!
  \******************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return SeriesChart; });
/* harmony import */ var core_graph__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! core/graph */ "./src/core/graph.ts");

/**
 * activeHover: 支持鼠标悬浮时高亮折线
 */
class SeriesChart extends core_graph__WEBPACK_IMPORTED_MODULE_0__["default"] {
    constructor(id, options) {
        super(id, Object.assign(Object.assign({ activeHover: true }, options), { isSeries: true }));
    }
    setOptions(options) {
        super.setOptions(Object.assign(Object.assign({}, options), { isSeries: true }));
    }
    draw() {
        this._mainPanel.clearRect(0, 0, this._options.width, this._options.height);
        this.drawLoading();
        this.drawTitle();
        // 绘画坐标轴
        this.drawAxis();
        // 绘画等高线
        this.drawGrid();
        this.drawSeriesLabels();
        this.drawLegends();
        this.drawEmptyData();
        // 绘画数据折线
        this.drawLine();
    }
}


/***/ }),

/***/ "./src/core/event.ts":
/*!***************************!*\
  !*** ./src/core/event.ts ***!
  \***************************/
/*! exports provided: HoverStatus, default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "HoverStatus", function() { return HoverStatus; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return Event; });
const HoverStatus = {
    on: 'on',
    out: 'out',
};
class Event {
    constructor(paint) {
        this._mouseDownEventList = [];
        this._paint = paint;
        this._paint.addEventListener("mousedown", this.onMouseDown.bind(this));
    }
    registerMouseDownEvent(rect, func) {
        this._mouseDownEventList.push({
            area: rect,
            func
        });
    }
    removeAllMouseDownEvent() {
        this._mouseDownEventList = [];
    }
    removeMouseDownEvent(area) {
        this._mouseDownEventList = this._mouseDownEventList.filter(item => {
            if (item.area.top === area.top && item.area.bottom === area.bottom &&
                item.area.left === area.left && item.area.right === area.right) {
                return false;
            }
            return true;
        });
    }
    onMouseDown(e) {
        const mouseX = e.clientX - this._paint.bounds.left;
        const mouseY = e.clientY - this._paint.bounds.top;
        this._mouseDownEventList.forEach(event => {
            const { area, func } = event;
            if (mouseX > area.left && mouseX < area.right && mouseY > area.top && mouseY < area.bottom) {
                func();
            }
        });
    }
}


/***/ }),

/***/ "./src/core/graph.ts":
/*!***************************!*\
  !*** ./src/core/graph.ts ***!
  \***************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return Graph; });
/* harmony import */ var _paint__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./paint */ "./src/core/paint.ts");
/* harmony import */ var _tooltip__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./tooltip */ "./src/core/tooltip.ts");
/* harmony import */ var _event__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./event */ "./src/core/event.ts");
/* harmony import */ var _utils__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./utils */ "./src/core/utils.ts");
/* harmony import */ var _theme__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./theme */ "./src/core/theme.ts");
/* harmony import */ var core_model__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! core/model */ "./src/core/model.ts");







/**
 * 设置options后要先进行相关参数的计算，在 setOptions 方法中进行。
 * 把计算流程集中在 setOptions 进行统一管理。
 * 整个图表的核心Graph分为两大块：数据计算和canvas绘图，分别以 draw 和 calculate 为前缀。
 *
 */
class Graph {
    constructor(id, options) {
        this._options = {
            activeColor: "#006eff",
            disabledColor: "#CCC",
            gridColor: "#eee",
            axisColor: "#000",
            auxiliaryLineColor: "#819AA4",
            fontColor: "rgb(124, 134, 142)",
            colorTheme: _theme__WEBPACK_IMPORTED_MODULE_4__["COLORS"].Types.Default,
            // 图表显示基本配置
            width: 600,
            height: 400,
            paddingHorizontal: 35,
            paddingVertical: 40,
            font: "normal 10px Verdana, Helvetica, Arial, sans-serif",
            labelScale: 10,
            gridScale: 5,
            showLegend: true,
            showYAxis: false,
            activeHover: false,
            hoverPrecision: 2,
            labelAlignCenter: true,
            showTooltip: true,
            isKilobyteFormat: false,
            overlay: false,
            isSeriesTime: false,
            isSeries: false,
            max: 0,
            min: 0,
            hoverPointData: null,
            loading: false,
            reload: null,
            title: "",
            tooltipLabels: null,
            yAxis: [],
            unit: "",
            // 图表数据
            labels: [],
            series: [],
            colors: null
        };
        this._drawHoverLine = false; // 在 onMouseMove 中鼠标hover时，标记是否有折线被hover，mainPanel 需要被重绘
        this._showColorPoint = true; // 如果折线只用单色显示，在不需要在tooltip中显示线颜色
        this._parentContainer = document.getElementById(id);
        if (!this._parentContainer) {
            return;
        }
        this.setOptions(options);
        // 新建一个div，用于设置两个 canvas 图层
        this._container = document.createElement("div");
        this._container.style.cssText = `width: ${this._options.width}px; 
      height: ${this._options.height}px;
      position: relative; 
      text-align: left;
      background-color: #fff;
      cursor: auto;`;
        // 新建两个 canvas 图层
        this._mainPanel = new _paint__WEBPACK_IMPORTED_MODULE_0__["default"]("", this._options.width, this._options.height, "position: absolute; user-select: none; z-index: 1");
        // 统一设置字体居中
        this._mainPanel.context2d.textAlign = "center";
        this._container.appendChild(this._mainPanel.canvas);
        this._auxiliaryPanel = new _paint__WEBPACK_IMPORTED_MODULE_0__["default"]("", this._options.width, this._options.height, "position: absolute; -webkit-tap-highlight-color: transparent; user-select: none; cursor: default; z-index: 2");
        this._auxiliaryPanel.context2d.textAlign = "center";
        this._container.appendChild(this._auxiliaryPanel.canvas);
        // 显示选择点信息
        this._tooltip = new _tooltip__WEBPACK_IMPORTED_MODULE_1__["default"](() => {
            this.clearAuxiliaryInfo();
        }, params => {
            this.highlightLine(params.legend);
        });
        this._container.appendChild(this._tooltip.entry);
        // 将新建的 div 填充到父容器中
        this._parentContainer.innerHTML = ""; // 清空容器 div 中的元素
        this._parentContainer.appendChild(this._container);
        // 注册事件对象，只有辅助层才有事件响应，所以只需要注册 _auxiliaryPanel
        this._event = new _event__WEBPACK_IMPORTED_MODULE_2__["default"](this._auxiliaryPanel);
        this.addListeners();
    }
    get DefaultColor() {
        return _theme__WEBPACK_IMPORTED_MODULE_4__["COLORS"].Types.Default;
    }
    // 添加 canvas 事件
    addListeners() {
        this._container.addEventListener("mousemove", this.onMouseMove.bind(this));
        this._container.addEventListener("mouseout", this.onMouseOut.bind(this));
        this._container.addEventListener("click", this.onMouseDown.bind(this));
    }
    isPointInActualArea(x, y) {
        return (x < this._chartPosition.right &&
            x > this._chartPosition.left &&
            y > this._chartPosition.top &&
            y < this._chartPosition.bottom);
    }
    /**
     * 高亮 legend 名称的折线
     * @param {string} legend
     */
    highlightLine(legend) {
        this._options.series.forEach(line => {
            line.hover = line.legend === _utils__WEBPACK_IMPORTED_MODULE_3__["FormatStringNoHTMLSharp"](legend);
        });
        // 重绘折线
        this.drawHighlightLine();
    }
    /**
     * 对数值进行格式化显示
     * @param {number} value
     * @returns {string}
     */
    formatValue(value) {
        if (!value || isNaN(value)) {
            return `${value}`;
        }
        value = _utils__WEBPACK_IMPORTED_MODULE_3__["MATH"].Strip(value, 6);
        return `${value}`;
    }
    // 设置图标配置项
    setOptions(options) {
        this._options = Object.assign(Object.assign({}, this._options), options);
        // 用户设置了Y轴的值，就设置了刻度数
        this._options.gridScale =
            this._options.yAxis.length > 0
                ? this._options.yAxis.length
                : this._options.gridScale;
        // 显示视图的范围点
        this._chartPosition = {
            left: this._options.paddingHorizontal,
            right: (this._options.width - this._options.paddingHorizontal),
            top: (this._options.paddingVertical * 2),
            bottom: (this._options.height - this._options.paddingVertical)
        };
        // 视图大小
        this._chartSize = {
            width: (this._chartPosition.right - this._chartPosition.left),
            height: (this._chartPosition.bottom - this._chartPosition.top)
        };
        // 当有数据变化是重新计算 canvas 各参数
        if (options.hasOwnProperty("series")) {
            this._options.series = JSON.parse(JSON.stringify(options.series));
            const colors = this._options.colors ||
                _theme__WEBPACK_IMPORTED_MODULE_4__["COLORS"].Values[this._options.colorTheme] ||
                _theme__WEBPACK_IMPORTED_MODULE_4__["COLORS"].Values[_theme__WEBPACK_IMPORTED_MODULE_4__["COLORS"].Types.Default];
            // tooltip 显示折线点颜色
            this._showColorPoint = !this._options.overlay && colors.length !== 1;
            /**
             * 初始化每个line的显示属性
             */
            this._options.series.forEach((line, index) => {
                line.legend = Object(_utils__WEBPACK_IMPORTED_MODULE_3__["FormatStringNoHTMLSharp"])(line.legend);
                line.hover = false; // hover 则高亮曲线
                line.show = true;
                line.color = line.disable
                    ? this._options.disabledColor
                    : colors[index % colors.length];
            });
            // 根据 line.disable 进行排序，disable先画线，显示的曲线不会被disable曲线覆盖
            this._options.series.sort((x, y) => x.disable === y.disable ? 0 : x.disable ? -1 : 1);
        }
        const modelOptions = {
            chartPosition: this._chartPosition,
            chartSize: this._chartSize,
            labelAlignCenter: this._options.labelAlignCenter,
            yGridNum: this._options.gridScale,
            isKilobyteFormat: this._options.isKilobyteFormat,
            sequence: JSON.parse(JSON.stringify(this._options.yAxis)) // 非空数组的话用于计算数据点的位置, 复制数组防止污染外部数据
        };
        this.initModel(modelOptions);
        this.clearAuxiliaryInfo();
    }
    // 重新设置图表的大小
    setChartSize(width, height) {
        this.setOptions({ width, height });
        if (this._container) {
            this._container.style.width = `${width}px`;
            this._container.style.height = `${height}px`;
        }
        this._mainPanel && this._mainPanel.setSize(width, height);
        this._auxiliaryPanel && this._auxiliaryPanel.setSize(width, height);
        this.draw();
    }
    initModel(modelOptions) {
        // 状态模式
        if (this._options.isSeries) {
            this._options.max =
                this._options.max === 0
                    ? this._options.labels[this._options.labels.length - 1]
                    : this._options.max;
            this._options.min =
                this._options.min === 0 ? this._options.labels[0] : this._options.min;
            this._model = new core_model__WEBPACK_IMPORTED_MODULE_5__["SeriesModel"](Object.assign(Object.assign({}, modelOptions), { pieceOfLabel: { max: this._options.max, min: this._options.min }, overlay: this._options.overlay }), this._options.labels, this._options.series);
        }
        else if (this._options.overlay) {
            this._model = new core_model__WEBPACK_IMPORTED_MODULE_5__["OverlayModel"](modelOptions, this._options.labels, this._options.series);
        }
        else {
            this._model = new core_model__WEBPACK_IMPORTED_MODULE_5__["Model"](modelOptions, this._options.labels, this._options.series);
        }
    }
    /**
     * 当要显示的label数大于默认设置的label数时可以进行简化，
     * @returns {boolean}
     */
    hasSimplyLabel() {
        return this._options.labels.length > this._options.labelScale;
    }
    /**
     * 鼠标悬浮事件，辅助 panel 不可用
     */
    isEventUnable() {
        return this._options.loading || this._options.series.length === 0;
    }
    /**
     * 清除鼠标悬浮时显示的辅助信息
     */
    clearAuxiliaryInfo() {
        this._tooltip && this._tooltip.close();
        this._auxiliaryPanel &&
            this._auxiliaryPanel.clearRect(0, 0, this._options.width, this._options.height);
    }
    drawLoading() {
        if (this._options.loading) {
            this._auxiliaryPanel.drawSpinner();
        }
        else {
            this._auxiliaryPanel.removeSpinner();
        }
    }
    drawTitle() {
        if (!this._options.title) {
            return;
        }
        if (!this._titleElement) {
            this._titleElement = _utils__WEBPACK_IMPORTED_MODULE_3__["GRAPH"].CreateDivElement(this._options.paddingVertical, this._chartPosition.left, this._options.axisColor);
            this._container.appendChild(this._titleElement);
        }
        const title = this._options.title || "";
        const unit = this._options.unit ? `（${this._options.unit}）` : "";
        this._titleElement.innerHTML = `
        ${title}
        <span style="color: ${this._options.fontColor}">${unit}</span>
    `;
    }
    /**
     * 绘画坐标轴 x，y
     */
    drawAxis() {
        // draw xAxis
        this._mainPanel.drawLine(this._chartPosition.left, this._chartPosition.bottom, this._chartPosition.right, this._chartPosition.bottom, this._options.axisColor);
        // showYAxis 控制是否显示Y轴
        if (this._options.showYAxis &&
            this._model.MaxValue &&
            this._model.MaxValue !== 0) {
            // draw yAxis
            this._mainPanel.drawLine(this._chartPosition.left, this._chartPosition.top, this._chartPosition.left, this._chartPosition.bottom, this._options.axisColor);
        }
    }
    drawLegends() {
        if (!this._options.showLegend || this._options.overlay) {
            // 叠加图状态不显示legend
            this._event.removeAllMouseDownEvent();
            return;
        }
        const markRectSize = { width: 10, height: 8 };
        // 默认每个legend的间隔
        const legendGap = 10;
        let legendX = this._options.paddingHorizontal;
        let legendY = this._options.paddingVertical * 1.5;
        this._options.series.forEach((line, index) => {
            const legendWidth = this._mainPanel.measureText(line.legend, this._options.font) +
                markRectSize.width;
            // 判断是否需要换行
            if (legendX + legendWidth > this._chartPosition.right) {
                legendX = this._options.paddingHorizontal;
                legendY += markRectSize.height * 1.5;
            }
            // 绘画色块
            const upperLeftCornerX = legendX;
            const upperLeftCornerY = legendY - markRectSize.height;
            this._mainPanel.drawRect(upperLeftCornerX, upperLeftCornerY, markRectSize.width, markRectSize.height, line.show ? line.color : this._options.disabledColor);
            // 绘画直线名称
            this._mainPanel.drawText(line.legend, legendX + markRectSize.width + 2, legendY, line.show ? this._options.fontColor : this._options.disabledColor, this._options.font, "left");
            //注册legend点击事件
            const eventArea = {
                top: upperLeftCornerY,
                bottom: legendY,
                left: legendX,
                right: legendX + legendWidth
            };
            // 之前存在事件，则移除
            this._event.removeMouseDownEvent(eventArea);
            // 注册事件
            this._event.registerMouseDownEvent(eventArea, () => {
                line.show = !line.show;
                if (this._options.overlay) {
                    // 叠加态需要重新计算点的位置
                    this._model.setOptions({}, this._options.labels, this._options.series);
                }
                this.draw();
            });
            legendX += legendWidth + legendGap;
        });
    }
    drawLabels() {
        if (this._model.XAxisTickMark.length === 0) {
            return;
        }
        const hasSimplyLabel = this.hasSimplyLabel();
        // labelGapNum 标记隔多少个label显示一次
        const labelGapNum = hasSimplyLabel
            ? Math.ceil(this._options.labels.length / this._options.labelScale)
            : 1;
        const tickMarkGap = this._options.labelAlignCenter
            ? this._model.XAxisTickMarkGap / 2
            : this._model.XAxisTickMarkGap;
        const textY = this._chartPosition.bottom + this._options.paddingVertical / 2;
        // 获取需要显示的label数组
        let visibleLabels = this._options.labels.filter((label, index) => index % labelGapNum === 0);
        if (this._options.isSeriesTime) {
            // 处理时序化的标签，对隔天的标签要做日期化处理
            visibleLabels = _utils__WEBPACK_IMPORTED_MODULE_3__["TIME"].FormatSeriesTime(visibleLabels);
        }
        // 使用间隔数据 labelGapNum 进行for循环
        for (let i = 0; i < visibleLabels.length; i++) {
            const label = visibleLabels[i];
            const xAxisLabelScale = this._model.XAxisTickMark[i * labelGapNum];
            // 简化label状态下不需要对 tickMark 进行位移
            const xAxisTickMark = hasSimplyLabel
                ? xAxisLabelScale
                : xAxisLabelScale + tickMarkGap;
            // 刻度线
            this._mainPanel.drawLine(xAxisTickMark, this._chartPosition.bottom, xAxisTickMark, this._chartPosition.bottom + 5, this._options.axisColor);
            this._mainPanel.drawText(label, xAxisLabelScale, textY, this._options.fontColor, this._options.font);
        }
    }
    drawSeriesLabels() {
        if (this._options.min === 0 && this._options.max === 0) {
            return;
        }
        // 需要满足现实标签的空间
        const labelNum = Math.floor(this._options.width / _theme__WEBPACK_IMPORTED_MODULE_4__["CHART"].LabelGap);
        let visibleLabels = _utils__WEBPACK_IMPORTED_MODULE_3__["TIME"].GenerateSeriesTimeVisibleLabels(this._options.min, this._options.max, labelNum);
        const distance = this._options.max - this._options.min;
        const textY = this._chartPosition.bottom + this._options.paddingVertical / 2;
        // 使用间隔数据 labelGapNum 进行for循环
        for (let i = 0; i < visibleLabels.length; i++) {
            const label = visibleLabels[i];
            const xAxisTickMark = this._chartPosition.left +
                (this._chartSize.width * (label - this._options.min)) / distance;
            // 刻度线
            this._mainPanel.drawLine(xAxisTickMark, this._chartPosition.bottom, xAxisTickMark, this._chartPosition.bottom + 5, this._options.axisColor);
            // 刻度值
            this._mainPanel.drawText(_utils__WEBPACK_IMPORTED_MODULE_3__["TIME"].FormatTime(label), xAxisTickMark, textY, this._options.fontColor, this._options.font);
        }
    }
    /**
     * 绘画等分线，且标记值
     * @param {number} scaleValue 单位刻度值
     */
    drawGrid() {
        if (this._model.ScaleValue === 0) {
            return;
        }
        const gridGap = Math.ceil(this._chartSize.height / (this._model.Sequence.length - 1));
        for (let i = 0; i < this._model.Sequence.length; i++) {
            // 从上到下画等高线
            const gridLineY = this._chartPosition.top + i * gridGap;
            // 在图表区域内画线
            if (gridLineY < this._chartPosition.bottom) {
                this._mainPanel.drawLine(this._chartPosition.left, gridLineY, this._chartPosition.right, gridLineY, this._options.gridColor);
            }
            if (gridLineY + this._chartSize.width <= this._chartPosition.bottom) {
                // 根据奇偶绘画不同色块
                this._mainPanel.drawRect(this._chartPosition.left, gridLineY, this._chartSize.width, gridGap, i % 2 === 0 ? "#F9F9F9" : "#FFF");
            }
            // 写刻度值,跟tooltip一直，Y轴和tooltip可以用凑通过tooltipLabels方法进行显示
            let gridScaleValueText = this._options.tooltipLabels
                ? this._options.tooltipLabels(this._model.Sequence[i], "yAxis")
                : this.formatValue(this._model.Sequence[i]);
            // text位置超过 this._chartPosition.left - textAndYAxisGap，则在canvas画布0位置上又对齐，
            const textAndYAxisGap = 10;
            const textWidth = this._mainPanel.measureText(gridScaleValueText, this._options.font);
            const textX = this._chartPosition.left - textAndYAxisGap - textWidth;
            this._mainPanel.drawText(gridScaleValueText, textX > 0 ? textX : 0, gridLineY, this._options.fontColor, this._options.font, "left");
        }
    }
    /**
     * 数据为空时提醒无数据
     */
    drawEmptyData() {
        this._reloadElement && (this._reloadElement.style.display = "none");
        if (this._options.loading) {
            return;
        }
        if (!this._reloadElement) {
            this._reloadElement = _utils__WEBPACK_IMPORTED_MODULE_3__["GRAPH"].CreateDivElement(this._options.height / 2, this._options.width / 2, this._options.fontColor);
            this._reloadElement.style.display = "none";
            this._reloadElement.style.transform = "translate(-50%, -50%)";
            const noDataText = document.createElement("span");
            noDataText.innerText = "暂无数据, ";
            this._reloadElement.appendChild(noDataText);
            const reloadText = document.createElement("span");
            reloadText.innerText = "重新加载";
            reloadText.style.color = this._options.activeColor;
            reloadText.onclick = () => {
                this._options.reload();
            };
            this._reloadElement.appendChild(reloadText);
            this._container.appendChild(this._reloadElement);
        }
        if (this._options.series.length === 0 &&
            this._options.labels.length === 0) {
            if (this._options.labels.length === 0) {
                this._reloadElement && (this._reloadElement.style.display = "block");
            }
            else {
                this._mainPanel.drawText("无数据", this._options.width / 2, this._options.height / 2, this._options.fontColor, "13px Verdana, Helvetica, Arial, sans-serif");
            }
        }
    }
    /**
     * 高亮状态时需要重绘折线
     */
    drawHighlightLine() {
        if (this._options.activeHover) {
            this.draw();
        }
    }
    /**
     * 当鼠标悬浮于图表时显示辅助信息
     * @param {Array<{data: object; color: string}>} showPoints 显示点数据
     * @param {{x: number; y: number}} mousePosition                 辅助线位置
     */
    drawAuxiliaryInfo(showPoints, mousePosition) {
        // 画纵轴辅助线
        this._auxiliaryPanel.drawDashLine(mousePosition.x, this._chartPosition.top, mousePosition.x, this._chartPosition.bottom, this._options.auxiliaryLineColor);
        // 绘画数据点
        showPoints.forEach(item => {
            this._auxiliaryPanel.drawPoints(item.point, item.color);
        });
    }
    /**
     * 显示 tooltip 信息
     * @param {boolean} fixedTooltip       是否可以设置可固定tooltip
     * @param {{x: number; y: number}} position
     * @param {{lebel: string; tooltipContent: Array<TooltipContentType>}} toolTipInfo tooltip信息
     */
    drawTooltip(fixedTooltip, position, toolTipInfo) {
        // 绘画 tooltip
        this._tooltip.setFixed(fixedTooltip);
        this._tooltip.setInformation(toolTipInfo.label, toolTipInfo.content, this._showColorPoint);
        // 计算能显示tooltip的x坐标, y坐标
        const toolTipX = position.x < this._options.width / 2
            ? position.x + 10
            : position.x - 10 - this._tooltip.width;
        const toolTipY = position.y < this._options.height / 2
            ? position.y + 10
            : position.y - 10 - this._tooltip.height;
        this._tooltip.show(toolTipX, toolTipY);
    }
    /**
     * 折线图画折线
     */
    drawLine() {
        // 画点线操作平频繁，使用requestAnimationFrame提升性能
        this._options.series.forEach((line, index) => {
            if (line.show) {
                const lineWidth = line.hover ? _theme__WEBPACK_IMPORTED_MODULE_4__["LINE"].Width.active : _theme__WEBPACK_IMPORTED_MODULE_4__["LINE"].Width.normal;
                this._mainPanel.drawPolyLine(this._model.XAxisTickMark, line.yPos, line.color, lineWidth);
            }
        });
    }
    /**
     * 面积图
     */
    drawArea() {
        this._options.series.forEach((line, index) => {
            const xPos = this._model.XAxisTickMark;
            // points被取消，需要对 drawArea 进行重写
            if (line.show) {
                const previousLinePoints = index > 0
                    ? this._options.series[index - 1].points
                    : {
                        [xPos[0]]: this._chartPosition.bottom - 0.5,
                        [xPos[xPos.length - 1]]: this._chartPosition.bottom - 0.5
                    };
                this._mainPanel.drawArea(previousLinePoints, line.points, line.color, 0.3);
                this._mainPanel.drawPolyLine(xPos, line.yPos, line.color);
            }
        });
    }
    /**
     * 面积图
     */
    drawAreaOnlyShowTopLine() {
        const xPos = this._model.XAxisTickMark;
        const color = this.DefaultColor;
        for (let i = this._options.series.length - 1; i >= 0; i--) {
            const line = this._options.series[i];
            if (line.show) {
                this._mainPanel.drawPolyLine(xPos, line.yPos, color);
                break;
            }
        }
    }
    /**
     * 条形图
     */
    drawBar() {
        const barWidth = this._model.XAxisTickMarkGap * 0.8;
        const marginLeft = barWidth / 2;
        let previousHeight = this._options.labels.map(i => 0);
        // 绘画数据折线
        this._options.series.forEach(line => {
            if (line.show) {
                this._model.XAxisTickMark.forEach((x, index) => {
                    const y = line.yPos[index];
                    const barHeight = this._chartPosition.bottom - y - (previousHeight[index] || 0);
                    this._mainPanel.drawRect(x - marginLeft, y, barWidth, barHeight, line.color, 0.8);
                    previousHeight[index] = previousHeight[index]
                        ? previousHeight[index] + barHeight
                        : barHeight;
                });
            }
        });
    }
    /**
     * 悬浮鼠标事件处理
     * @param e
     */
    onMouseMove(e, fixedTooltip = false) {
        if (this.isEventUnable()) {
            return;
        }
        if (this._tooltip.fixed && !fixedTooltip) {
            // 当进入fixed tooltip状态时，onMouseMove 显示 tooltip 效果取消
            return;
        }
        const mouseX = e.clientX - this._mainPanel.bounds.left;
        const mouseY = e.clientY - this._mainPanel.bounds.top;
        this.clearAuxiliaryInfo();
        // 保证鼠标指针在有效区域内
        if (this.isPointInActualArea(mouseX, mouseY)) {
            let xAxisTickMarkIndex = 0;
            if (this._model.XAxisTickMarkGap === 0) {
                xAxisTickMarkIndex = this._model.XAxisTickMark.findIndex(x => x >= mouseX);
                if (xAxisTickMarkIndex === -1) {
                    return;
                }
            }
            else {
                xAxisTickMarkIndex = Math.round((mouseX - this._chartPosition.left) / this._model.XAxisTickMarkGap);
            }
            // 处理边界情况，当鼠标接近离开x轴宽度区域时，四舍五入会多1
            if (xAxisTickMarkIndex === this._options.labels.length) {
                xAxisTickMarkIndex -= 1;
            }
            const mousePosition = {
                x: this._model.XAxisTickMark[xAxisTickMarkIndex],
                y: mouseY
            };
            let showPoints = []; // 显示辅助线与折线交汇的点
            let content = []; // tooltip 要显示的数据内容
            let hasRedrawLine = false; // 标记是否调用 drawHighlightLine
            const previousPointX = this._model.XAxisTickMark[xAxisTickMarkIndex - 1] || 0;
            // 绘画数据点
            this._options.series.forEach((line, index) => {
                if (line.show) {
                    const pointY = line.yPos[xAxisTickMarkIndex];
                    if (pointY) {
                        const label = this._options.labels[xAxisTickMarkIndex];
                        if (this._options.overlay) {
                            // 显示点的位置信息, 叠加态只显示最大的点（一个）
                            if (index === this._options.series.length - 1) {
                                showPoints.push({
                                    point: { [mousePosition.x]: pointY },
                                    color: this.DefaultColor
                                });
                            }
                        }
                        else {
                            // 显示点的位置信息
                            showPoints.push({
                                point: { [mousePosition.x]: pointY },
                                color: line.color
                            });
                        }
                        // 不在叠加态、tooltip 未固定状态 且折线显示的情况下，才有悬浮高亮
                        if (this._options.activeHover && !fixedTooltip && !line.disable) {
                            // 判断鼠标是否 hover 在折线上，是则高亮折线
                            let hasHover = false;
                            const previousPointY = line.yPos[xAxisTickMarkIndex];
                            if (previousPointY) {
                                // 计算直线斜率 k=(y2-y1)/(x2-x1)
                                const k = (pointY - previousPointY) /
                                    (mousePosition.x - previousPointX);
                                const b = pointY - mousePosition.x * k;
                                hasHover =
                                    Math.abs(k * mouseX + b - mouseY) <
                                        this._options.hoverPrecision;
                            }
                            line.hover = hasHover;
                            if (hasHover) {
                                hasRedrawLine = true; // 标记需要高亮重绘折线
                            }
                        }
                        else {
                            line.hover = false;
                        }
                        // 置灰折线不需要在tooltip中显示信息
                        if (!line.disable) {
                            // 判断是否有自定义显示 tooltipLabels
                            const value = line.data[label];
                            const valueStr = this._options.tooltipLabels
                                ? this._options.tooltipLabels(value) || ""
                                : this.formatValue(value);
                            content.push({
                                legend: line.legend,
                                color: line.color,
                                hover: line.hover,
                                value: value,
                                label: `${valueStr}${this._options.unit || ""}`
                            });
                        }
                    }
                }
            });
            // 根据 value 对 tooltip 显示数据做排序显示
            content.sort((a, b) => b.value - a.value);
            // 判断title 是否是时间序列类型，是则做时间格式化
            const title = this._options.isSeriesTime
                ? _utils__WEBPACK_IMPORTED_MODULE_3__["TIME"].Format(this._options.labels[xAxisTickMarkIndex], _utils__WEBPACK_IMPORTED_MODULE_3__["TIME"].DateFormat.fullDateTime)
                : this._options.labels[xAxisTickMarkIndex];
            const showInfo = {
                label: title,
                content
            };
            // 外部入口，处理显示的数据
            this._options.hoverPointData &&
                this._options.hoverPointData({
                    xAxisTickMarkIndex,
                    mousePosition,
                    content
                });
            // 显示辅助面板信息
            this.drawAuxiliaryInfo(showPoints, mousePosition);
            if (this._options.showTooltip) {
                // 绘画 tooltip
                this.drawTooltip(fixedTooltip, mousePosition, showInfo);
            }
            // 判断是否与上次显示状态一致，不一致则需要重绘折线
            if (hasRedrawLine !== this._drawHoverLine) {
                this.drawHighlightLine();
                this._drawHoverLine = hasRedrawLine;
            }
        }
    }
    /**
     * 鼠标点击事件处理
     * @param e
     */
    onMouseDown(e) {
        e.stopPropagation();
        const mouseX = e.clientX - this._mainPanel.bounds.left;
        const mouseY = e.clientY - this._mainPanel.bounds.top;
        if (this._tooltip.isPointInTooltipArea(mouseX, mouseY)) {
            return;
        }
        this.onMouseMove(e, true);
    }
    onMouseOut(e) {
        e.stopPropagation();
        const mouseX = e.clientX - this._mainPanel.bounds.left;
        const mouseY = e.clientY - this._mainPanel.bounds.top;
        if (!this.isPointInActualArea(mouseX, mouseY)) {
            if (this._options.showTooltip && !this._tooltip.fixed) {
                // 显示tooltip情况下且tooltip不是固定状态，mouse out 才需 clear 面板
                this.clearAuxiliaryInfo();
            }
            // 外部入口，处理显示的数据
            this._options.hoverPointData &&
                this._options.hoverPointData({
                    xAxisTickMarkIndex: 0,
                    mousePosition: {},
                    content: []
                });
        }
    }
}


/***/ }),

/***/ "./src/core/model.ts":
/*!***************************!*\
  !*** ./src/core/model.ts ***!
  \***************************/
/*! exports provided: Model, OverlayModel, SeriesModel */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Model", function() { return Model; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "OverlayModel", function() { return OverlayModel; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "SeriesModel", function() { return SeriesModel; });
/* harmony import */ var _utils__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./utils */ "./src/core/utils.ts");

const DefaultMaxValue = 5;
const InitialState = {
    maxValue: DefaultMaxValue,
    minValue: 0,
    sequence: [],
    scaleValue: { exponent: 0, value: 0 },
    xAxisTickMarkGap: 0,
    xAxisTickMark: []
};
/**
 * 图表的数据模型，处理传入的数据
 */
class BaseModel {
    constructor(options, label, series) {
        this.state = Object.assign({}, InitialState);
        this._options = {
            overlay: false,
            pieceOfLabel: {
                min: 0,
                max: 0
            },
            labelAlignCenter: false,
            isKilobyteFormat: false,
            sequence: [],
            yGridNum: 0,
            chartPosition: {
                right: 0,
                left: 0,
                top: 0,
                bottom: 0
            },
            chartSize: {
                width: 0,
                height: 0
            }
        };
        this.setOptions(options, label, series);
    }
    setOptions(options, labels, series) {
        this._options = Object.assign(Object.assign({}, this._options), options);
        this.calculatePoints(labels, series);
    }
    get ScaleValue() {
        return this.state.scaleValue.value * (10 ** this.state.scaleValue.exponent);
    }
    // ScaleValue 基数
    get ScaleValueCardinal() {
        return this.state.scaleValue.value;
    }
    // ScaleValue 指数
    get ScaleValueExponent() {
        return this.state.scaleValue.exponent;
    }
    get MaxValue() {
        return this.state.maxValue;
    }
    get XAxisTickMarkGap() {
        return this.state.xAxisTickMarkGap;
    }
    get XAxisTickMark() {
        return this.state.xAxisTickMark;
    }
    get Sequence() {
        return this.state.sequence;
    }
    resetState() {
        this.state = Object.assign({}, InitialState);
    }
    setState(states) {
        this.state = Object.assign(Object.assign({}, this.state), states);
    }
    arrayMaxAndMinValue(sources) {
        let maximum = 0, minimum = 0;
        if (sources.length > 0) {
            maximum = Math.max(...sources);
            const min = Math.min(...sources);
            if (min < 0) {
                minimum = min;
            }
        }
        return { maximum, minimum };
    }
    /**
     * 根据数据源计算y轴显示整数最大值
     */
    calculateScaleAndMaxValue(maximum, minimum) {
        // 处理数值为0的边界情况
        if (maximum === 0) {
            // 如果最大值为0，则使用默认值
            maximum = DefaultMaxValue;
        }
        // 用户自定义刻度大小
        if (this._options.sequence.length > 0 && this._options.sequence[this._options.sequence.length - 1] > maximum) {
            const sequence = this._options.sequence;
            const gridValue = sequence[sequence.length - 1] / sequence.length;
            this.setState({
                scaleValue: { exponent: 0, value: gridValue },
                maxValue: sequence[sequence.length - 1],
                minValue: minimum,
                sequence: sequence.reverse()
            });
            return;
        }
        const { gridValue, sequence } = this._options.isKilobyteFormat
            ? _utils__WEBPACK_IMPORTED_MODULE_0__["MATH"].ArithmeticSequence(minimum, maximum, this._options.yGridNum)
            : _utils__WEBPACK_IMPORTED_MODULE_0__["MATH"].ScaleSequence(minimum, maximum, this._options.yGridNum);
        this.setState({
            scaleValue: { exponent: 0, value: gridValue },
            maxValue: sequence[sequence.length - 1],
            minValue: minimum,
            sequence: sequence.reverse()
        });
    }
    // 计算 label 在 x轴上的位置,数据点的x坐标需要使用
    calculateXAxisTickMark(labelsNum) {
        const { labelAlignCenter, chartPosition, chartSize } = this._options;
        this.state.xAxisTickMark = [];
        this.state.xAxisTickMarkGap = chartSize.width / labelsNum;
        // 计算标签居中间隔
        const centerMove = labelAlignCenter ? this.state.xAxisTickMarkGap / 2 : 0;
        // 计算标签居中间隔
        for (let i = 0; i < labelsNum; i++) {
            this.state.xAxisTickMark.push(chartPosition.left + i * this.state.xAxisTickMarkGap + centerMove);
        }
    }
}
/**
 * 常规计算
 */
class Model extends BaseModel {
    constructor(options, labels, series) {
        super(options, labels, series);
    }
    calculateMaxAndMinSourceValue(series) {
        // 返回每个对象数据的数据集，即为二维数组
        const data2dArray = series.map(item => Object.values(item.data));
        // 转换为一维数组 [].concat(...arr2d)，求出所有数据的最大值
        const dataFromAllSources = [].concat(...data2dArray);
        return this.arrayMaxAndMinValue(dataFromAllSources);
    }
    /**
     * 核心对外API
     * @param {Array<DataType>} series
     */
    calculatePoints(labels, series) {
        const { chartPosition, chartSize } = this._options;
        if (labels.length === 0) {
            this.resetState();
            return;
        }
        // X 轴刻度线坐标数组
        this.calculateXAxisTickMark(labels.length);
        const { maximum, minimum } = this.calculateMaxAndMinSourceValue(series);
        this.calculateScaleAndMaxValue(maximum, minimum);
        const { maxValue, xAxisTickMark } = this.state;
        // 每单位数值在y轴的高度
        const pointYGap = maxValue != 0 ? chartSize.height / maxValue : 0;
        series.forEach((line) => {
            // 初始化属性值 yPos 和 show
            line.yPos = [];
            labels.forEach((label, index) => {
                const value = line.data[label];
                if (value !== null && !isNaN(value)) {
                    const pointY = chartPosition.bottom - Math.round(pointYGap * value);
                    // 在 canvas actual 区域内绘画
                    line.yPos.push(pointY);
                }
                else {
                    line.yPos.push(null);
                }
            });
        });
    }
}
/**
 * 叠加计算
 */
class OverlayModel extends BaseModel {
    constructor(options, labels, series) {
        super(options, labels, series);
    }
    calculateMaxAndMinSourceValue(series) {
        // 返回每组数据的数据集，即为二维数组
        const data2dArray = series.map(item => Object.values(item.data));
        const lineNum = data2dArray.length;
        if (lineNum > 0) {
            let combineArray = [];
            data2dArray[0].forEach((value, index) => {
                let sum = 0;
                // 计算各个节点不同线的总和
                for (let j = 0; j < lineNum; j++) {
                    sum += data2dArray[j][index] || 0;
                }
                combineArray.push(sum);
            });
            return this.arrayMaxAndMinValue(combineArray);
        }
        return this.arrayMaxAndMinValue([]);
    }
    /**
     * 核心对外API
     * @param {Array<DataType>} series
     */
    calculatePoints(labels, series) {
        const { chartPosition, chartSize } = this._options;
        if (labels.length === 0) {
            this.resetState();
            return;
        }
        // X 轴刻度线坐标数组
        this.calculateXAxisTickMark(labels.length);
        const { maximum, minimum } = this.calculateMaxAndMinSourceValue(series);
        this.calculateScaleAndMaxValue(maximum, minimum);
        const { maxValue, xAxisTickMark } = this.state;
        let previousHeight = labels.map(i => 0);
        // 每单位数值在y轴的高度
        const pointYGap = maxValue != 0 ? chartSize.height / maxValue : 0;
        series.forEach((line) => {
            // 初始化属性值 yPos 和 show
            line.yPos = [];
            labels.forEach((label, index) => {
                const value = line.data[label] || 0;
                const pointY = chartPosition.bottom - Math.round(pointYGap * (value + (previousHeight[index] || 0)));
                // 在 canvas actual 区域内绘画
                line.yPos.push(pointY);
                previousHeight[index] = previousHeight[index] ? previousHeight[index] + value : value;
            });
        });
    }
}
/**
 * 时序图根据x点计算位置
 */
class SeriesModel extends BaseModel {
    constructor(options, labels, series) {
        super(options, labels, series);
    }
    calculateMaxAndMinSourceValue(series) {
        // 返回每个对象数据的数据集，即为二维数组
        const data2dArray = series.map(item => Object.values(item.data));
        // 返回每组数据的数据集，即为二维数组
        const lineNum = data2dArray.length;
        if (lineNum > 0) {
            let combineArray = [];
            if (this._options.overlay) {
                data2dArray[0].forEach((value, index) => {
                    let sum = 0;
                    // 计算各个节点不同线的总和
                    for (let j = 0; j < lineNum; j++) {
                        sum += data2dArray[j][index] || 0;
                    }
                    combineArray.push(sum);
                });
            }
            else {
                // 转换为一维数组 [].concat(...arr2d)，求出所有数据的最大值
                combineArray = [].concat(...data2dArray);
            }
            return this.arrayMaxAndMinValue(combineArray);
        }
        return this.arrayMaxAndMinValue([]);
    }
    calculatePoints(labels, series) {
        const { chartPosition, pieceOfLabel, chartSize } = this._options;
        if (labels.length === 0) {
            this.resetState();
            return;
        }
        const { maximum, minimum } = this.calculateMaxAndMinSourceValue(series);
        this.calculateScaleAndMaxValue(maximum, minimum);
        const { maxValue } = this.state;
        let xAxisTickMark = [];
        const distance = pieceOfLabel.max - pieceOfLabel.min;
        labels.forEach(label => {
            const pointX = chartPosition.left + (label - pieceOfLabel.min) * chartSize.width / distance;
            xAxisTickMark.push(pointX);
        });
        this.setState({ xAxisTickMark });
        let previousHeight = labels.map(i => 0);
        // 每单位数值在y轴的高度
        const pointYGap = maxValue != 0 ? chartSize.height / maxValue : 0;
        if (this._options.overlay) {
            series.forEach((line) => {
                // 初始化属性值 yPos 和 show
                line.yPos = [];
                labels.forEach((label, index) => {
                    const value = line.data[label];
                    previousHeight[index] = value !== null && !isNaN(value) ? previousHeight[index] + value : previousHeight[index];
                    const pointY = chartPosition.bottom - Math.round(pointYGap * previousHeight[index]);
                    // 在 canvas actual 区域内绘画
                    line.yPos.push(pointY);
                });
            });
        }
        else {
            series.forEach((line) => {
                // 初始化属性值 yPos 和 show
                line.yPos = [];
                labels.forEach((label, index) => {
                    const value = line.data[label];
                    //          const pointX = xAxisTickMark[index];
                    if (value !== null && !isNaN(value)) {
                        const pointY = chartPosition.bottom - Math.round(pointYGap * value);
                        // 在 canvas actual 区域内绘画
                        line.yPos.push(pointY);
                    }
                    else {
                        line.yPos.push(null);
                    }
                });
            });
        }
    }
}


/***/ }),

/***/ "./src/core/paint.ts":
/*!***************************!*\
  !*** ./src/core/paint.ts ***!
  \***************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return Paint; });
/* harmony import */ var core_theme__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! core/theme */ "./src/core/theme.ts");
/* harmony import */ var _utils__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./utils */ "./src/core/utils.ts");


class Paint {
    constructor(id, width, height, style) {
        this._lineDashConfig = [5, 5];
        if (id) {
            this._canvas = document.getElementById(id);
        }
        else {
            this._canvas = document.createElement("canvas");
        }
        this._ctx = this._canvas.getContext('2d');
        this._canvas.style.cssText = style;
        this.setSize(width, height);
    }
    get canvas() {
        return this._canvas;
    }
    get bounds() {
        return this._canvas.getBoundingClientRect();
    }
    get context2d() {
        return this._ctx;
    }
    setSize(width, height) {
        _utils__WEBPACK_IMPORTED_MODULE_1__["GRAPH"].ScaleCanvas(this._canvas, this._ctx, width, height);
    }
    setTextAlign(textAlign = "center") {
        this.context2d.textAlign = textAlign;
    }
    setLineDash(lineDash) {
        this._lineDashConfig = lineDash;
    }
    addEventListener(event, callback) {
        this._canvas.addEventListener(event, callback, false);
    }
    clearRect(startX, startY, endX, endY) {
        this._ctx.clearRect(startX, startY, endX, endY);
    }
    /**
     * 绘画线
     * @param {number} startX
     * @param {number} startY
     * @param {number} endX
     * @param {number} endY
     * @param {string} color
     */
    drawLine(startX, startY, endX, endY, color) {
        this._ctx.save();
        this._ctx.strokeStyle = color;
        this._ctx.beginPath();
        this._ctx.moveTo(startX, startY);
        this._ctx.lineTo(endX, endY);
        this._ctx.stroke();
        this._ctx.restore();
    }
    drawDashLine(startX, startY, endX, endY, color) {
        this._ctx.save();
        this._ctx.strokeStyle = color;
        this._ctx.beginPath();
        this._ctx.setLineDash(this._lineDashConfig);
        this._ctx.moveTo(startX, startY);
        this._ctx.lineTo(endX, endY);
        this._ctx.stroke();
        this._ctx.restore();
    }
    /**
     *
     * @param {Array<number>} xPos 为 labels 的 x 轴坐标数组
     * @param {object} points
     * @param {string} color
     */
    drawPolyLine(xPos, yPos, color, lineWidth = 1) {
        this._ctx.save();
        // draw ploy line
        this._ctx.strokeStyle = color;
        this._ctx.beginPath();
        let isContinuous = false; // 用来标记前面的点是否存在，存在就画线，不存在则需要移动后判断是否要画点
        this._ctx.lineWidth = lineWidth;
        xPos.forEach((x, index) => {
            const y = yPos[index];
            if (isNaN(y) || y === null) {
                if (isContinuous) {
                    // 优化性能，不能在lineTo下进行stroke，会导致画图缓慢
                    // 通过判断上传数据存在，只触发一次stroke
                    this._ctx.stroke();
                }
                // y 坐标点不存在，则下个点为移动坐标
                isContinuous = false;
                return;
            }
            else {
                if (isContinuous) {
                    this._ctx.lineTo(x, y);
                }
                else {
                    this._ctx.moveTo(x, y);
                    // 画点：如果下一个点不存在，则没有连续性，需要画点
                    const yNext = yPos[index + 1];
                    if (isNaN(yNext) || yNext === null) {
                        this._ctx.fillStyle = color;
                        this._ctx.beginPath();
                        this._ctx.arc(x, y, 2, 0, Math.PI * 2, true);
                        this._ctx.closePath();
                        this._ctx.fill();
                    }
                }
                isContinuous = true;
            }
        });
        this._ctx.stroke();
        this._ctx.restore();
    }
    /**
     * 保证 previousPoints 与 currentPoints 起始位置和结束位置一致
     * @param {object} previousPoints
     * @param {object} currentPoints
     * @param {string} color
     * @param {number} alpha
     */
    drawArea(previousPoints, currentPoints, color, alpha) {
        function sortNumber(a, b) {
            return Number(a) - Number(b);
        }
        //    Chrome Opera 的 JavaScript 解析引擎遵循的是新版 ECMA-262 第五版规范。因此，使用 for-in 语句遍历对象属性时遍历书序并非属性构建顺序。
        //    而 IE6 IE7 IE8 Firefox Safari 的 JavaScript 解析引擎遵循的是较老的 ECMA-262 第三版规范，属性遍历顺序由属性构建的顺序决定。
        //
        //    Chrome Opera 中使用 for-in 语句遍历对象属性时会遵循一个规律：
        //    它们会先提取所有 key 的 parseFloat 值为非负整数的属性，然后根据数字顺序对属性排序首先遍历出来，然后按照对象定义的顺序遍历余下的所有属性。
        let xPos = Object.keys(currentPoints).sort(sortNumber);
        // previousPoints 的起始或结束位置没有数据，则返回
        if (isNaN(previousPoints[xPos[0]]) && isNaN(previousPoints[xPos[xPos.length - 1]])) {
            return;
        }
        this._ctx.save();
        this._ctx.fillStyle = color;
        this._ctx.globalAlpha = alpha;
        this._ctx.beginPath();
        const startX = Number(xPos[0]);
        // 绘画区域
        this._ctx.moveTo(startX, previousPoints[startX]);
        xPos.forEach(x => {
            this._ctx.lineTo(Number(x), currentPoints[x]);
        });
        xPos = Object.keys(previousPoints).sort(sortNumber);
        for (let i = xPos.length - 1; i >= 0; i--) {
            const x = xPos[i];
            this._ctx.lineTo(Number(x), previousPoints[x]);
        }
        this._ctx.closePath();
        this._ctx.fill();
        this._ctx.restore();
    }
    drawPoints(points, color) {
        this._ctx.save();
        const endAngle = Math.PI * 2;
        let xPos = Object.keys(points);
        xPos.forEach(x => {
            this._ctx.fillStyle = "#fff";
            this._ctx.beginPath();
            this._ctx.arc(Number(x), points[x], 6, 0, endAngle, true);
            this._ctx.closePath();
            this._ctx.fill();
            this._ctx.fillStyle = color;
            this._ctx.beginPath();
            this._ctx.arc(Number(x), points[x], 4, 0, endAngle, true);
            this._ctx.closePath();
            this._ctx.fill();
        });
        this._ctx.restore();
    }
    /**
     * 绘画柱
     * @param {number} upperLeftCornerX
     * @param {number} upperLeftCornerY
     * @param {number} width
     * @param {number} height
     * @param {string} color
     */
    drawRect(upperLeftCornerX, upperLeftCornerY, width, height, color, alpha = 1.0) {
        this._ctx.save();
        this._ctx.fillStyle = color;
        this._ctx.globalAlpha = alpha;
        this._ctx.fillRect(upperLeftCornerX, upperLeftCornerY, width, height);
        this._ctx.restore();
    }
    /**
     * @param {string} text
     * @param {number} x
     * @param {number} y
     * @param {string} color
     * @param {string} font
     */
    drawText(text, x, y, color, font, textAlign = "center") {
        this.setTextAlign(textAlign);
        this._ctx.save();
        this._ctx.fillStyle = color;
        this._ctx.font = font;
        this._ctx.fillText(text, x, y);
        this._ctx.restore();
        this.setTextAlign();
    }
    /**
     * drawSpinner 动画函数
     */
    drawSpinner() {
        // spinner 参数, 参数也可以写为对象的属性，但是这不符合最小化作用范围
        const _degrees = new Date();
        const _offset = 16;
        const width = Number(this._canvas.style.width.replace('px', ''));
        const height = Number(this._canvas.style.height.replace('px', ''));
        const moveX = width / 120;
        const lineX = width / 60;
        const lineWidth = width / 500;
        function spinnerAnimation() {
            this._spinnerAnimation = window.requestAnimationFrame(spinnerAnimation.bind(this));
            const rotation = parseInt((((new Date() - _degrees) / 1000) * _offset)) / _offset;
            this._ctx.save();
            this._ctx.clearRect(0, 0, width, height);
            this._ctx.translate(width / 2, height / 2);
            this._ctx.rotate(Math.PI * 2 * rotation);
            for (let i = 0; i < _offset; i++) {
                this._ctx.beginPath();
                this._ctx.rotate(Math.PI * 2 / _offset);
                this._ctx.moveTo(moveX, 0);
                this._ctx.lineTo(lineX, 0);
                this._ctx.lineWidth = lineWidth;
                this._ctx.strokeStyle = "rgba(0, 111, 250," + i / _offset + ")";
                this._ctx.stroke();
            }
            this._ctx.restore();
        }
        ;
        // 防止有重复的spinner
        this.removeSpinner();
        spinnerAnimation.apply(this);
        // 安全机制，用于在限制时间(6秒)内关闭drawSpinner (云API5秒请求超时)
        setTimeout(() => {
            this.removeSpinner();
        }, core_theme__WEBPACK_IMPORTED_MODULE_0__["CHART"].SpinnerTime);
    }
    removeSpinner() {
        this._spinnerAnimation && window.cancelAnimationFrame(this._spinnerAnimation);
        if (this._spinnerAnimation) {
            this._ctx.clearRect(0, 0, this._canvas.width, this._canvas.height);
        }
        this._spinnerAnimation = null;
    }
    /**
     * 返回
     * @param {string} text
     * @returns {number}
     */
    measureText(text, font) {
        this._ctx.font = font;
        return this._ctx.measureText(text).width;
    }
}


/***/ }),

/***/ "./src/core/theme.ts":
/*!***************************!*\
  !*** ./src/core/theme.ts ***!
  \***************************/
/*! exports provided: COLORS, LINE, CHART */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "COLORS", function() { return COLORS; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "LINE", function() { return LINE; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "CHART", function() { return CHART; });
/**
 * 主题配置数据信息
 */
var COLORS;
(function (COLORS) {
    let Types;
    (function (Types) {
        Types["Default"] = "default";
        Types["Multi"] = "multi";
    })(Types = COLORS.Types || (COLORS.Types = {}));
    COLORS.Values = {
        [Types.Default]: ["#007EFA"],
        [Types.Multi]: [
            "#007EFA",
            "#29cc85",
            "#ffbb00",
            "#ff584c",
            "#9741d9",
            "#1fc0cc",
            "#7dd936",
            "#ff9c19",
            "#e63984",
            "#655ce6",
            "#47cc50",
            "#bf30b6"
        ]
    };
})(COLORS || (COLORS = {}));
/**
 * 折线配置
 */
var LINE;
(function (LINE) {
    LINE.Width = {
        normal: 1,
        active: 2,
    };
    LINE.Color = {
        active: "#0064E1"
    };
})(LINE || (LINE = {}));
var CHART;
(function (CHART) {
    CHART.LabelGap = 60;
    CHART.SpinnerTime = 6000;
})(CHART || (CHART = {}));


/***/ }),

/***/ "./src/core/tooltip.ts":
/*!*****************************!*\
  !*** ./src/core/tooltip.ts ***!
  \*****************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return Tooltip; });
class Tooltip {
    constructor(closeEvent, hoverEvent) {
        this.MaxNumOfData = 8; // 同屏下最多显示数据条目数
        this._fixed = false;
        this._showColorPoint = false;
        this._closeEvent = closeEvent;
        this._hoverEvent = hoverEvent;
        this._self = document.createElement("div");
        this._self.style.cssText = 'position: absolute; ' +
            'background-color: white; ' +
            'opacity: 0.9; ' +
            'color: #000000; ' +
            'padding: 0.3rem 0.5rem; ' +
            'border-radius: 3px; ' +
            'z-index: 10; ' +
            'visibility: hidden; ' +
            'font-size: 12px; ' +
            'box-shadow: rgba(33, 33, 33, 0.2) 0px 1px 2px; ' +
            'line-height: 1.2em; ' +
            'top: 0px;';
        this._css = document.createElement("style");
        this._css.innerHTML = '.TChart_close-btn {' +
            '  position: absolute;' +
            '  right: 10px;' +
            '  top: 5px;' +
            '  display: inline-block;' +
            '  width: 17px;' +
            '  height: 17px;' +
            '  opacity: 0.3;' +
            '}' +
            '.TChart_close-btn:hover {' +
            '  opacity: 1;' +
            '}' +
            '.TChart_close-btn:before, .TChart_close-btn:after {' +
            '  position: absolute;' +
            '  left: 15px;' +
            '  content: \' \';' +
            '  height: 15px;' +
            '  width: 2px;' +
            '  background-color: #333;' +
            '}' +
            '.TChart_close-btn:before {' +
            '  transform: rotate(45deg);' +
            '}' +
            '.TChart_close-btn:after {' +
            '  transform: rotate(-45deg);' +
            '}' +
            '.TChart_tooltip-data:hover {' +
            '  color: #006eff !important;' +
            '}' +
            '.TChart_tooltip-data.hover {' +
            '  color: #006eff !important;' +
            '}';
        document.getElementsByTagName("head")[0].appendChild(this._css);
    }
    get entry() {
        return this._self;
    }
    get width() {
        return this._self.offsetWidth;
    }
    get height() {
        return this._self.offsetHeight;
    }
    get X() {
        return this._x;
    }
    get Y() {
        return this._y;
    }
    get fixed() {
        return this._fixed;
    }
    isPointInTooltipArea(x, y) {
        return (x < this.X + this.width &&
            x > this.X &&
            y > this.Y &&
            y < this.Y + this.height);
    }
    setInformation(title, content, showColorPoint) {
        this._self.innerHTML = `
                  <table style="${content.length < this.MaxNumOfData ? '' : 'height:130px;'} display: flex; flex-direction: column;">
                    <thead style="margin-bottom: 4px; display: block; position: relative; ">
                      <tr>
                        <td style="display: inline-block; font-size: 12px; height: 20px; line-height: 20px">${title}</td>
                        <td style="display: inline-block; margin-left: 10px; width: 20px;height: 20px"></td>
                      </tr>
                     </thead>
                   </table>`;
        if (this._fixed) {
            this._closeBtn = document.createElement("i");
            this._closeBtn.setAttribute("class", "TChart_close-btn");
            this._self.getElementsByTagName("td")[1].appendChild(this._closeBtn);
            this._closeBtn && (this._closeBtn.onclick = (e) => {
                e.stopPropagation();
                this.close();
                this._closeEvent && this._closeEvent();
            });
        }
        this._tooltipDataTbody = document.createElement("tbody");
        this._tooltipDataTbody.style = "flex: 1;overflow: auto;";
        this._showColorPoint = showColorPoint;
        this._tooltipDataTbody.innerHTML = content.map(item => {
            return `<tr style="text-align: right;" class="TChart_tooltip-data ${item.hover ? 'hover' : ''}" data="${item.legend}">
                ${this._showColorPoint
                ? `<td><span style="display: block; margin-right: 4px; border-radius: 50%; width: 8px; height: 8px; background: ${item.color};"></span></td>`
                : ""}
                <td style="display: block;text-align: left; min-width: 100px; margin-right: 8px;">${item.legend} </td>
                <td style="text-align: right;">${item.label}</td> 
              </tr>`;
        }).join("");
        ;
        this._self.firstElementChild.appendChild(this._tooltipDataTbody);
        // 注册 tooltip 数据条 hover 事件
        if (this._tooltipDataTbody) {
            this._tooltipDataTbody.onmouseover = (e) => {
                e.stopPropagation();
                let legend = "";
                if (e.target.parentElement.hasAttribute("data")) {
                    legend = e.target.parentElement.getAttribute("data");
                }
                this._hoverEvent && this._hoverEvent({ legend });
            };
            this._tooltipDataTbody.onmouseout = (e) => {
                e.stopPropagation();
                this._hoverEvent && this._hoverEvent({ legend: "" });
            };
        }
    }
    /**
     * 关闭tooltip，取消固定状态
     */
    close() {
        this.setFixed(false);
        this.hidden();
    }
    show(x, y) {
        this._x = x;
        this._y = y;
        this._self.style.visibility = "visible";
        this._self.style.transform = `translate(${x}px, ${y}px)`;
    }
    hidden() {
        this._self.style.visibility = "hidden";
    }
    setFixed(fixed) {
        this._fixed = fixed;
    }
}


/***/ }),

/***/ "./src/core/utils.ts":
/*!***************************!*\
  !*** ./src/core/utils.ts ***!
  \***************************/
/*! exports provided: MATH, TIME, FormatStringNoHTMLSharp, GRAPH */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "MATH", function() { return MATH; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "TIME", function() { return TIME; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "FormatStringNoHTMLSharp", function() { return FormatStringNoHTMLSharp; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "GRAPH", function() { return GRAPH; });
var MATH;
(function (MATH) {
    /**
     * 解决浮点误差
     * @param {Number} num
     * @param {number} precision
     * @returns {number}
     * @constructor
     */
    function Strip(num, precision = 12) {
        return +parseFloat(num.toPrecision(precision));
    }
    MATH.Strip = Strip;
    function ScientificNotation(value, exponent) {
        if (exponent > 0) {
            return `${value}E${exponent}`;
        }
        else if (exponent === 0) {
            return `${value}`;
        }
        else {
            return `${value}e${exponent}`;
        }
    }
    MATH.ScientificNotation = ScientificNotation;
    /**
     * 将数据转换为文件格式, 1024 = 2 ** 10
     * @param {number} size
     * @returns {string}
     */
    function FormatFileSize(size, toFixed = 3) {
        const Kilobyte = 1024;
        const fileSizeUnits = ["", "K", "M", "G", "T", "P", "E", "Z", "Y"];
        let exponent = 0;
        const maxExponent = fileSizeUnits.length - 1;
        while (size >= Kilobyte) {
            if (exponent >= maxExponent) {
                break;
            }
            size = size / Kilobyte;
            exponent += 1;
        }
        return `${Strip(size, toFixed)}${fileSizeUnits[exponent]}`;
    }
    MATH.FormatFileSize = FormatFileSize;
    /**
     * 根据最小值，最大值和间隔数，计算对应的10进制转2进制的等差数组
     * @param {number} minimum 最小值
     * @param {number} maximum 最大值
     * @param {number} gridNum 间隔数
     */
    function ArithmeticSequence(minimum, maximum, gridNum = 5) {
        if (gridNum === 0) {
            return { gridValue: 0, sequence: [] };
        }
        if (minimum === maximum && gridNum > 0) {
            return { gridValue: 0, sequence: [] };
        }
        // if minimum > maximum, reverse them
        if (minimum > maximum) {
            const temp = minimum;
            minimum = maximum;
            maximum = temp;
        }
        const step = Math.max(0, gridNum);
        // 计算数值间隔
        const gridGap = (maximum - minimum) / step;
        const power = Math.floor(Math.log2(gridGap)); // power 是以2为底的指数，2**power <= gridGap
        const ratio = gridGap / Math.pow(2, power);
        // 参考 d3.array.ticks, 原为10进制，改为2进制，我也不理解
        let multiple = 1;
        if (ratio >= Math.sqrt(50)) {
            multiple = 10;
        }
        else if (ratio >= Math.sqrt(10)) {
            multiple = 5;
        }
        else if (ratio >= Math.sqrt(2)) {
            multiple = 2;
        }
        // 等差值
        const gridValue = power >= 0 ? multiple * Math.pow(2, power) : -Math.pow(2, -power) / multiple;
        if (gridValue === 0 || !isFinite(gridValue)) {
            return { gridValue: 0, sequence: [] };
        }
        let ticks, n, i = 0;
        if (gridValue > 0) {
            const start = Math.floor(minimum / gridValue);
            const stop = Math.ceil(maximum / gridValue);
            ticks = new Array((n = Math.ceil(stop - start + 1)));
            while (i < n) {
                ticks[i] = (start + i) * gridValue;
                i += 1;
            }
        }
        else {
            const start = Math.ceil(minimum * gridValue);
            const stop = Math.floor(maximum * gridValue);
            ticks = new Array((n = Math.ceil(start - stop + 1)));
            while (i < n) {
                ticks[i] = (start - i) / gridValue;
                i += 1;
            }
        }
        return { gridValue: i > 0 ? ticks[1] - ticks[0] : ticks[0], sequence: ticks };
    }
    MATH.ArithmeticSequence = ArithmeticSequence;
    function ScaleSequence(minimum, maximum, gridNum = 5) {
        const averageValue = maximum / gridNum;
        // 将最大值转换为 个位数 * 10的指数
        // 指数级
        let exponent = 0;
        let scaleValue = averageValue;
        if (-1 < scaleValue && scaleValue < 1) {
            // 小于零，乘以10计算到个位
            while (scaleValue < 1) {
                scaleValue = scaleValue * 10;
                exponent -= 1;
            }
        }
        else if (scaleValue > 10) {
            // 大于十，除以10计算到个位
            while (scaleValue > 10) {
                scaleValue = scaleValue / 10;
                exponent += 1;
            }
        }
        // 向上取整
        scaleValue = Math.ceil(scaleValue);
        const gridValue = scaleValue * 10 ** exponent;
        // 还未优化，以后支持负值显示
        let sequence = [];
        for (let i = 0; i <= gridNum; i++) {
            sequence.push(gridValue * i);
        }
        return { gridValue, sequence };
    }
    MATH.ScaleSequence = ScaleSequence;
})(MATH || (MATH = {}));
var TIME;
(function (TIME) {
    TIME.DateFormat = {
        fullDateTime: "YYYY-MM-dd hh:mm",
        dateTime: "MM-dd hh:mm",
        time: "hh:mm",
        day: "MM-dd"
    };
    function IsSameDay(dateLeft, dateRight) {
        let dateToCheck = dateLeft;
        let actualDate = dateRight;
        if (!(dateLeft instanceof Date) || !(dateRight instanceof Date)) {
            dateToCheck = new Date(dateLeft);
            actualDate = new Date(dateRight);
        }
        return (dateToCheck.getFullYear() === actualDate.getFullYear() &&
            dateToCheck.getMonth() === actualDate.getMonth() &&
            dateToCheck.getDate() === actualDate.getDate());
    }
    TIME.IsSameDay = IsSameDay;
    function Format(date, format) {
        let dateTmp = date;
        if (!(dateTmp instanceof Date)) {
            dateTmp = new Date(date);
        }
        const dateValues = [
            dateTmp.getFullYear(),
            dateTmp.getMonth() + 1,
            dateTmp.getDate(),
            dateTmp.getHours(),
            dateTmp.getMinutes(),
            dateTmp.getSeconds()
        ];
        const timeFormat = format
            .replace(/(YYYY|yyyy)/, dateValues[0].toString())
            .replace(/(MM)/, dateValues[1] < 10 ? `0${dateValues[1]}` : `${dateValues[1]}`)
            .replace(/(DD|dd)/, dateValues[2] < 10 ? `0${dateValues[2]}` : `${dateValues[2]}`)
            .replace(/(hh)/, dateValues[3] < 10 ? `0${dateValues[3]}` : `${dateValues[3]}`)
            .replace(/(mm)/, dateValues[4] < 10 ? `0${dateValues[4]}` : `${dateValues[4]}`)
            .replace(/(ss)/, dateValues[5] < 10 ? `0${dateValues[5]}` : `${dateValues[5]}`);
        return timeFormat;
    }
    TIME.Format = Format;
    function FormatTime(time) {
        const date = new Date(time);
        const hour = date.getHours();
        const minute = date.getMinutes();
        if (hour === 0 && minute === 0) {
            return Format(date, TIME.DateFormat.day);
        }
        return Format(date, TIME.DateFormat.time);
    }
    TIME.FormatTime = FormatTime;
    function FormatSeriesTime(seriesTime = []) {
        let format = TIME.DateFormat.time;
        const labelNum = seriesTime.length;
        return seriesTime.map((label, index) => {
            if (index + 1 === labelNum) {
                return Format(new Date(label), format);
            }
            if (!IsSameDay(label, seriesTime[index + 1])) {
                format = TIME.DateFormat.day;
                return Format(new Date(label), format);
            }
            const item = Format(new Date(label), format);
            format = TIME.DateFormat.time;
            return item;
        });
    }
    TIME.FormatSeriesTime = FormatSeriesTime;
    function GenerateSeriesTimeVisibleLabels(min, max, labelNum = 12) {
        if (max < min) {
            const temp = max;
            max = min;
            min = temp;
        }
        const defaultGaps = [300000, 600000, 900000, 1800000, 3600000, 7200000, 21600000, 43200000, 86400000]; // 单位毫秒， 5，10，15，30，60，120，360，720, 1440 分钟
        const actualGap = Math.floor((max - min) / labelNum); // labelNum 默认显示十二个label
        let gap = defaultGaps.find(gap => gap >= actualGap) || defaultGaps[defaultGaps.length - 1];
        const firstDate = new Date(min);
        let timestamp = 0;
        if (gap >= 86400000) {
            // 大于一天的间隔
            const hour = firstDate.getHours();
            const endDate = new Date(max);
            // 天数大于labelNum显示数，需要调整gap大小
            const diffDays = Math.round(Math.abs((firstDate.getTime() - endDate.getTime()) / 86400000));
            if (diffDays > labelNum) {
                gap = Math.ceil(diffDays / labelNum) * 86400000;
            }
            const remainder = hour % 24; // 小时余数
            timestamp = firstDate.setHours(remainder, 0, 0);
        }
        else if (gap >= 3600000) {
            let hour = firstDate.getHours();
            const gapHour = gap / 3600000; // 60分钟
            // 分钟余数
            const remainder = hour % gapHour;
            // 起始点取整点
            hour = hour + gapHour - remainder;
            timestamp = firstDate.setHours(hour, 0, 0);
        }
        else {
            let minute = firstDate.getMinutes();
            const gapMinute = gap / 60000; //一分钟
            // 分钟余数
            const remainder = minute % gapMinute;
            minute = minute + gapMinute - remainder;
            timestamp = firstDate.setMinutes(minute, 0);
        }
        let visibleLabels = [];
        while (timestamp < max) {
            visibleLabels.push(timestamp);
            timestamp += gap;
        }
        return visibleLabels;
    }
    TIME.GenerateSeriesTimeVisibleLabels = GenerateSeriesTimeVisibleLabels;
})(TIME || (TIME = {}));
function FormatStringNoHTMLSharp(str) {
    let str_tmp = str;
    if (str_tmp.indexOf("<") !== -1 && str_tmp.indexOf(">") !== -1) {
        str_tmp = str_tmp.replace("<", "");
        str_tmp = str_tmp.replace(">", "");
    }
    return str_tmp;
}
var GRAPH;
(function (GRAPH) {
    /**
     * 解决在 retina 屏幕下显示模糊
     * @param canvas
     * @param context
     * @param customWidth
     * @param customHeight
     * @constructor
     */
    function ScaleCanvas(canvas, context, customWidth, customHeight) {
        if (!canvas || !context) {
            throw new Error("Must pass in `canvas` and `context`.");
        }
        const width = customWidth || canvas.width || canvas.clientWidth;
        const height = customHeight || canvas.height || canvas.clientHeight;
        const ratio = window.devicePixelRatio || 1;
        canvas.width = Math.round(width * ratio);
        canvas.height = Math.round(height * ratio);
        canvas.style.width = width + "px";
        canvas.style.height = height + "px";
        context.scale(ratio, ratio);
        return ratio;
    }
    GRAPH.ScaleCanvas = ScaleCanvas;
    function CreateDivElement(top, left, color) {
        const element = document.createElement("div");
        element.style.cssText = `
          z-index: 3;
          color: ${color};
          position: absolute; 
          top: ${top}px;
          left: ${left}px;
          transform: translateY(-50%);
          text-align: left;`;
        return element;
    }
    GRAPH.CreateDivElement = CreateDivElement;
})(GRAPH || (GRAPH = {}));


/***/ }),

/***/ "./src/panel/helper.ts":
/*!*****************************!*\
  !*** ./src/panel/helper.ts ***!
  \*****************************/
/*! exports provided: minPeriod, durationsByPeriod, Period, TimeFormat, TransformField, STORE */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "minPeriod", function() { return minPeriod; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "durationsByPeriod", function() { return durationsByPeriod; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Period", function() { return Period; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "TimeFormat", function() { return TimeFormat; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "TransformField", function() { return TransformField; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "STORE", function() { return STORE; });
/* harmony import */ var moment__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! moment */ "moment");
/* harmony import */ var moment__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(moment__WEBPACK_IMPORTED_MODULE_0__);

const DAY = 24 * 60 * 60;
const HOUR = 60 * 60;
const MINUTE = 60;
let minPeriod = 60;
let durationsByPeriod = {
    // [period(s)] : day
    10: 1,
    60: 15,
    300: 31,
    3600: 62,
    86400: 186
};
/**
 * 根据时间间隔计算查询时间粒度(秒)
 * @param {Date} from
 * @param {Date} to
 * @returns {string}
 */
function Period(from, to) {
    const range = moment__WEBPACK_IMPORTED_MODULE_0___default()(to).diff(moment__WEBPACK_IMPORTED_MODULE_0___default()(from));
    let startTime = moment__WEBPACK_IMPORTED_MODULE_0___default()(from);
    let startTimeDiffFromNow = -startTime.diff(undefined);
    if (startTimeDiffFromNow < durationsByPeriod[10] * DAY * 1000) {
        if (range <= 8 * HOUR * 1000)
            return 1 * MINUTE; // 最小
    }
    if (startTimeDiffFromNow < durationsByPeriod[MINUTE] * DAY * 1000) {
        if (range <= 2 * DAY * 1000)
            return 1 * MINUTE;
    }
    if (startTimeDiffFromNow < durationsByPeriod[5 * MINUTE] * DAY * 1000) {
        // 因后台存储，临时规则改为5 分钟 288 个点， 即3天
        // if (range <= 10 * DAY * 1000) return 5 * MINUTE;
        if (range <= 3 * DAY * 1000)
            return 5 * MINUTE;
    }
    if (startTimeDiffFromNow < durationsByPeriod[HOUR] * DAY * 1000) {
        if (range <= 120 * DAY * 1000)
            return HOUR;
    }
    if (startTimeDiffFromNow < durationsByPeriod[DAY] * DAY * 1000) {
        return DAY;
    }
    return DAY * 31;
}
function TimeFormat(from, to) {
    const range = moment__WEBPACK_IMPORTED_MODULE_0___default()(to).diff(moment__WEBPACK_IMPORTED_MODULE_0___default()(from), "hours");
    return range <= 24 ? "HH:mm" : "MM-DD HH:mm";
}
const UNITS = ["", "K", "M", "G", "T", "P"];
/**
 * 进行单位换算
 * @param {number} value
 * @param {number} thousands
 * @param {number} toFixed
 */
function TransformField(_value, thousands, toFixed = 3, units = UNITS) {
    let value = _value;
    let isValueDefined = !isNaN(value) && value !== null;
    if (!isValueDefined)
        return "-";
    let unitBase = units[0];
    let i = units.indexOf(unitBase);
    if (isValueDefined && thousands) {
        while (i < units.length && value / thousands > 1) {
            value /= thousands;
            ++i;
        }
        unitBase = units[i] || "";
    }
    let svalue;
    if (value > 1) {
        svalue = value.toFixed(toFixed);
        svalue = svalue.replace(/0+$/, "");
        svalue = svalue.replace(/\.$/, "");
    }
    else if (value) {
        // 如果数值很小，保留toFixed位有效数字
        let tens = 0;
        let v = Math.abs(value);
        while (v < 1) {
            v *= 10;
            ++tens;
        }
        svalue = value.toFixed(tens + toFixed - 1);
        svalue = svalue.replace(/0+$/, "");
        svalue = svalue.replace(/\.$/, "");
    }
    else {
        svalue = value;
    }
    return String(svalue) + (value !== 0 ? unitBase : "");
}
var STORE;
(function (STORE) {
    function Set(key, data) {
        localStorage.setItem(key, JSON.stringify(data));
    }
    STORE.Set = Set;
    function Get(key, defaultValue = null) {
        try {
            return JSON.parse(localStorage.getItem(key)) || defaultValue;
        }
        catch (e) {
            return defaultValue;
        }
    }
    STORE.Get = Get;
})(STORE || (STORE = {}));


/***/ }),

/***/ "./src/panel/index.tsx":
/*!*****************************!*\
  !*** ./src/panel/index.tsx ***!
  \*****************************/
/*! exports provided: request, TransformField, ChartPanel, ChartFilterPanel, ChartInstancesPanel, ColorTypes, default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChartPanel", function() { return ChartPanel; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChartFilterPanel", function() { return ChartFilterPanel; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChartInstancesPanel", function() { return ChartInstancesPanel; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _tencent_tea_component_lib_i18n__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @tencent/tea-component/lib/i18n */ "./node_modules/@tencent/tea-component/lib/i18n/index.js");
/* harmony import */ var _tencent_tea_component_lib_i18n__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_tencent_tea_component_lib_i18n__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var _tce_request__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../tce/request */ "./src/tce/request.ts");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "request", function() { return _tce_request__WEBPACK_IMPORTED_MODULE_2__["request"]; });

/* harmony import */ var _helper__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./helper */ "./src/panel/helper.ts");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "TransformField", function() { return _helper__WEBPACK_IMPORTED_MODULE_3__["TransformField"]; });

/* harmony import */ var charts_index__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! charts/index */ "./src/charts/index.ts");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "ColorTypes", function() { return charts_index__WEBPACK_IMPORTED_MODULE_4__["ColorTypes"]; });



Object(_tencent_tea_component_lib_i18n__WEBPACK_IMPORTED_MODULE_1__["setLocale"])(window.VERSION);


const ChartPanel = dynamicComponent("ChartPanel");
const ChartFilterPanel = dynamicComponent("ChartFilterPanel");
const ChartInstancesPanel = dynamicComponent("ChartInstancesPanel");

/* harmony default export */ __webpack_exports__["default"] = (ChartPanel);
/**
 * 动态引入chart component
 */
const ChartsComponents = react__WEBPACK_IMPORTED_MODULE_0__["lazy"](() => __webpack_require__.e(/*! import() | ChartsComponents */ "ChartsComponents").then(__webpack_require__.bind(null, /*! ./ChartsComponents */ "./src/panel/ChartsComponents.tsx")));
function dynamicComponent(componentName) {
    return function WrappedComponent(props) {
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"](react__WEBPACK_IMPORTED_MODULE_0__["Suspense"], { fallback: react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: {
                    width: "100%",
                    height: "100%"
                } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { className: "text-overflow" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("i", { className: "n-loading-icon" }))) },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](ChartsComponents, Object.assign({}, props, { componentName: componentName }))));
    };
}


/***/ }),

/***/ "./src/tce/request.ts":
/*!****************************!*\
  !*** ./src/tce/request.ts ***!
  \****************************/
/*! exports provided: apiRequest, request */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "apiRequest", function() { return apiRequest; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "request", function() { return request; });
async function apiRequest({ data }) {
    let res;
    // 融合版的tke 监控，使用influxdb
    let apiInfo = window["modules"]["monitor"];
    let json = {
        apiVersion: apiInfo
            ? `${apiInfo.groupName}/${apiInfo.version}`
            : "monitor.tkestack.io/v1",
        kind: "Metric",
        query: data.data.RequestBody,
    };
    json.query.offset = 0;
    json.query.conditions = json.query.conditions.map((item) => ({
        key: item[0],
        expr: item[1],
        value: item[2],
    }));
    let projectHeader = {};
    const projectName = new URLSearchParams(window.location.search).get("projectName");
    if (projectName) {
        projectHeader = {
            "X-TKE-ProjectName": projectName,
        };
    }
    try {
        res = await fetch(`/apis/monitor.tkestack.io/v1/metrics/`, {
            method: "POST",
            mode: "cors",
            cache: "no-cache",
            credentials: "same-origin",
            headers: Object.assign({ "Content-Type": "application/json" }, projectHeader),
            redirect: "follow",
            referrer: "no-referrer",
            body: JSON.stringify(json),
        });
        res = await res.json();
        return JSON.parse(res.jsonResult);
    }
    catch (error) {
        return {
            columns: 0,
            data: [],
        };
    }
}
async function request(data) {
    try {
        // 发送 API 请求
        let res = await apiRequest({
            data,
        });
        return { columns: res.columns, data: res.data || [] };
    }
    catch (error) {
        console.error(error, data);
        throw error;
    }
}
const getCookie = (name) => {
    let reg = new RegExp("(?:^|;+|\\s+)" + name + "=([^;]*)"), match = document.cookie.match(reg);
    return !match ? "" : match[1];
};


/***/ }),

/***/ "./webpack/alias/react.js":
/*!********************************!*\
  !*** ./webpack/alias/react.js ***!
  \********************************/
/*! no static exports found */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _react_global__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! __react-global */ "__react-global");
/* harmony import */ var _react_global__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_react_global__WEBPACK_IMPORTED_MODULE_0__);
/* harmony reexport (unknown) */ for(var __WEBPACK_IMPORT_KEY__ in _react_global__WEBPACK_IMPORTED_MODULE_0__) if(__WEBPACK_IMPORT_KEY__ !== 'default') (function(key) { __webpack_require__.d(__webpack_exports__, key, function() { return _react_global__WEBPACK_IMPORTED_MODULE_0__[key]; }) }(__WEBPACK_IMPORT_KEY__));



/* harmony default export */ __webpack_exports__["default"] = (_react_global__WEBPACK_IMPORTED_MODULE_0__);

/***/ }),

/***/ "@tencent/tea-component":
/*!*****************************************!*\
  !*** external "@tencent/tea-component" ***!
  \*****************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = require("@tencent/tea-component");

/***/ }),

/***/ "__react-dom-global":
/*!************************************!*\
  !*** external "window.ReactDOM16" ***!
  \************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = window.ReactDOM16;

/***/ }),

/***/ "__react-global":
/*!*********************************!*\
  !*** external "window.React16" ***!
  \*********************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = window.React16;

/***/ }),

/***/ "moment":
/*!*************************!*\
  !*** external "moment" ***!
  \*************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = require("moment");

/***/ })

/******/ });
//# sourceMappingURL=TChart.js.map