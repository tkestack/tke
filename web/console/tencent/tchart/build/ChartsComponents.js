((typeof self !== 'undefined' ? self : this)["webpackJsonp"] = (typeof self !== 'undefined' ? self : this)["webpackJsonp"] || []).push([["ChartsComponents"],{

/***/ "./node_modules/classnames/index.js":
/*!******************************************!*\
  !*** ./node_modules/classnames/index.js ***!
  \******************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

var __WEBPACK_AMD_DEFINE_ARRAY__, __WEBPACK_AMD_DEFINE_RESULT__;/*!
  Copyright (c) 2017 Jed Watson.
  Licensed under the MIT License (MIT), see
  http://jedwatson.github.io/classnames
*/
/* global define */

(function () {
	'use strict';

	var hasOwn = {}.hasOwnProperty;

	function classNames () {
		var classes = [];

		for (var i = 0; i < arguments.length; i++) {
			var arg = arguments[i];
			if (!arg) continue;

			var argType = typeof arg;

			if (argType === 'string' || argType === 'number') {
				classes.push(arg);
			} else if (Array.isArray(arg) && arg.length) {
				var inner = classNames.apply(null, arg);
				if (inner) {
					classes.push(inner);
				}
			} else if (argType === 'object') {
				for (var key in arg) {
					if (hasOwn.call(arg, key) && arg[key]) {
						classes.push(key);
					}
				}
			}
		}

		return classes.join(' ');
	}

	if ( true && module.exports) {
		classNames.default = classNames;
		module.exports = classNames;
	} else if (true) {
		// register as 'classnames', consistent with npm package name
		!(__WEBPACK_AMD_DEFINE_ARRAY__ = [], __WEBPACK_AMD_DEFINE_RESULT__ = (function () {
			return classNames;
		}).apply(exports, __WEBPACK_AMD_DEFINE_ARRAY__),
				__WEBPACK_AMD_DEFINE_RESULT__ !== undefined && (module.exports = __WEBPACK_AMD_DEFINE_RESULT__));
	} else {}
}());


/***/ }),

/***/ "./node_modules/css-loader/dist/cjs.js!./node_modules/less-loader/dist/cjs.js!./src/panel/containers/Filter.less":
/*!***********************************************************************************************************************!*\
  !*** ./node_modules/css-loader/dist/cjs.js!./node_modules/less-loader/dist/cjs.js!./src/panel/containers/Filter.less ***!
  \***********************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(/*! ../../../node_modules/css-loader/dist/runtime/api.js */ "./node_modules/css-loader/dist/runtime/api.js")(false);
// Module
exports.push([module.i, ".Tchart_filter-row {\n  display: flex;\n}\n.tab-status-panel {\n  display: flex;\n  width: 100%;\n  transition: margin-left 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);\n}\n.no-tab-status-panel {\n  width: 100%;\n  height: calc(100vh - 150px);\n  overflow-y: scroll;\n  padding-right: 20px;\n  /* Increase/decrease this value for cross-browser compatibility */\n  box-sizing: content-box;\n  /* So the width will be 100% + 17px */\n}\n", ""]);



/***/ }),

/***/ "./node_modules/css-loader/dist/cjs.js!./node_modules/less-loader/dist/cjs.js!./src/panel/containers/Instances.less":
/*!**************************************************************************************************************************!*\
  !*** ./node_modules/css-loader/dist/cjs.js!./node_modules/less-loader/dist/cjs.js!./src/panel/containers/Instances.less ***!
  \**************************************************************************************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(/*! ../../../node_modules/css-loader/dist/runtime/api.js */ "./node_modules/css-loader/dist/runtime/api.js")(false);
// Module
exports.push([module.i, "/**\n InstanceList\n */\n.Tchart_table-box {\n  position: relative;\n  display: block;\n  width: 100%;\n  color: #666;\n  font-size: 12px;\n  margin-right: 0;\n  background-color: transparent;\n}\n.Tchart_table-box tr {\n  display: flex;\n  flex-direction: row;\n  width: 100%;\n}\n.Tchart_table-box tr th {\n  position: relative;\n  padding-left: 10px;\n  padding-right: 10px;\n  text-align: left;\n  vertical-align: middle;\n}\n.Tchart_table-box tr th:last-child {\n  flex: 1;\n  overflow: hidden;\n}\n.Tchart_table-box tr th > div {\n  position: relative;\n  display: inline-block;\n  padding: 0;\n  height: 40px;\n  line-height: 40px;\n  width: 100%;\n  color: #888;\n  vertical-align: middle;\n  box-sizing: border-box;\n  word-wrap: break-word;\n}\n.Tchart_table-box thead {\n  display: block;\n  width: 100%;\n  font-weight: 700;\n  font-size: 14px;\n  border-bottom: 1px solid #d1d5de;\n}\n.Tchart_table-box tbody {\n  display: block;\n  overflow-x: hidden;\n  overflow-y: auto;\n  max-height: calc(100vh - 243px);\n  height: 623px;\n  font-weight: 500;\n}\n.Tchart_table-box tbody tr:hover {\n  background-color: #F7F7F7;\n}\n.Tchart_table-box tbody th > div {\n  height: 45px;\n  line-height: 45px;\n  color: #444;\n  width: 100%;\n  white-space: nowrap;\n  overflow: hidden;\n  text-overflow: ellipsis;\n}\n.Tchart_table-box tbody th > div p {\n  height: 20px;\n  line-height: 20px;\n  width: 100%;\n  overflow: hidden;\n  text-overflow: ellipsis;\n  white-space: nowrap;\n}\n.Tchart_table-box tbody th > div p:last-child {\n  color: #888 !important;\n  font-size: 10px;\n}\n.Tchart_table-box tbody th:last-child > div {\n  padding: 5px;\n  line-height: 35px;\n  box-sizing: border-box;\n}\n", ""]);



/***/ }),

/***/ "./node_modules/css-loader/dist/runtime/api.js":
/*!*****************************************************!*\
  !*** ./node_modules/css-loader/dist/runtime/api.js ***!
  \*****************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


/*
  MIT License http://www.opensource.org/licenses/mit-license.php
  Author Tobias Koppers @sokra
*/
// css base code, injected by the css-loader
module.exports = function (useSourceMap) {
  var list = []; // return the list of modules as css string

  list.toString = function toString() {
    return this.map(function (item) {
      var content = cssWithMappingToString(item, useSourceMap);

      if (item[2]) {
        return '@media ' + item[2] + '{' + content + '}';
      } else {
        return content;
      }
    }).join('');
  }; // import a list of modules into the list


  list.i = function (modules, mediaQuery) {
    if (typeof modules === 'string') {
      modules = [[null, modules, '']];
    }

    var alreadyImportedModules = {};

    for (var i = 0; i < this.length; i++) {
      var id = this[i][0];

      if (id != null) {
        alreadyImportedModules[id] = true;
      }
    }

    for (i = 0; i < modules.length; i++) {
      var item = modules[i]; // skip already imported module
      // this implementation is not 100% perfect for weird media query combinations
      // when a module is imported multiple times with different media queries.
      // I hope this will never occur (Hey this way we have smaller bundles)

      if (item[0] == null || !alreadyImportedModules[item[0]]) {
        if (mediaQuery && !item[2]) {
          item[2] = mediaQuery;
        } else if (mediaQuery) {
          item[2] = '(' + item[2] + ') and (' + mediaQuery + ')';
        }

        list.push(item);
      }
    }
  };

  return list;
};

function cssWithMappingToString(item, useSourceMap) {
  var content = item[1] || '';
  var cssMapping = item[3];

  if (!cssMapping) {
    return content;
  }

  if (useSourceMap && typeof btoa === 'function') {
    var sourceMapping = toComment(cssMapping);
    var sourceURLs = cssMapping.sources.map(function (source) {
      return '/*# sourceURL=' + cssMapping.sourceRoot + source + ' */';
    });
    return [content].concat(sourceURLs).concat([sourceMapping]).join('\n');
  }

  return [content].join('\n');
} // Adapted from convert-source-map (MIT)


function toComment(sourceMap) {
  // eslint-disable-next-line no-undef
  var base64 = btoa(unescape(encodeURIComponent(JSON.stringify(sourceMap))));
  var data = 'sourceMappingURL=data:application/json;charset=utf-8;base64,' + base64;
  return '/*# ' + data + ' */';
}

/***/ }),

/***/ "./node_modules/style-loader/lib/addStyles.js":
/*!****************************************************!*\
  !*** ./node_modules/style-loader/lib/addStyles.js ***!
  \****************************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {

/*
	MIT License http://www.opensource.org/licenses/mit-license.php
	Author Tobias Koppers @sokra
*/

var stylesInDom = {};

var	memoize = function (fn) {
	var memo;

	return function () {
		if (typeof memo === "undefined") memo = fn.apply(this, arguments);
		return memo;
	};
};

var isOldIE = memoize(function () {
	// Test for IE <= 9 as proposed by Browserhacks
	// @see http://browserhacks.com/#hack-e71d8692f65334173fee715c222cb805
	// Tests for existence of standard globals is to allow style-loader
	// to operate correctly into non-standard environments
	// @see https://github.com/webpack-contrib/style-loader/issues/177
	return window && document && document.all && !window.atob;
});

var getTarget = function (target, parent) {
  if (parent){
    return parent.querySelector(target);
  }
  return document.querySelector(target);
};

var getElement = (function (fn) {
	var memo = {};

	return function(target, parent) {
                // If passing function in options, then use it for resolve "head" element.
                // Useful for Shadow Root style i.e
                // {
                //   insertInto: function () { return document.querySelector("#foo").shadowRoot }
                // }
                if (typeof target === 'function') {
                        return target();
                }
                if (typeof memo[target] === "undefined") {
			var styleTarget = getTarget.call(this, target, parent);
			// Special case to return head of iframe instead of iframe itself
			if (window.HTMLIFrameElement && styleTarget instanceof window.HTMLIFrameElement) {
				try {
					// This will throw an exception if access to iframe is blocked
					// due to cross-origin restrictions
					styleTarget = styleTarget.contentDocument.head;
				} catch(e) {
					styleTarget = null;
				}
			}
			memo[target] = styleTarget;
		}
		return memo[target]
	};
})();

var singleton = null;
var	singletonCounter = 0;
var	stylesInsertedAtTop = [];

var	fixUrls = __webpack_require__(/*! ./urls */ "./node_modules/style-loader/lib/urls.js");

module.exports = function(list, options) {
	if (typeof DEBUG !== "undefined" && DEBUG) {
		if (typeof document !== "object") throw new Error("The style-loader cannot be used in a non-browser environment");
	}

	options = options || {};

	options.attrs = typeof options.attrs === "object" ? options.attrs : {};

	// Force single-tag solution on IE6-9, which has a hard limit on the # of <style>
	// tags it will allow on a page
	if (!options.singleton && typeof options.singleton !== "boolean") options.singleton = isOldIE();

	// By default, add <style> tags to the <head> element
        if (!options.insertInto) options.insertInto = "head";

	// By default, add <style> tags to the bottom of the target
	if (!options.insertAt) options.insertAt = "bottom";

	var styles = listToStyles(list, options);

	addStylesToDom(styles, options);

	return function update (newList) {
		var mayRemove = [];

		for (var i = 0; i < styles.length; i++) {
			var item = styles[i];
			var domStyle = stylesInDom[item.id];

			domStyle.refs--;
			mayRemove.push(domStyle);
		}

		if(newList) {
			var newStyles = listToStyles(newList, options);
			addStylesToDom(newStyles, options);
		}

		for (var i = 0; i < mayRemove.length; i++) {
			var domStyle = mayRemove[i];

			if(domStyle.refs === 0) {
				for (var j = 0; j < domStyle.parts.length; j++) domStyle.parts[j]();

				delete stylesInDom[domStyle.id];
			}
		}
	};
};

function addStylesToDom (styles, options) {
	for (var i = 0; i < styles.length; i++) {
		var item = styles[i];
		var domStyle = stylesInDom[item.id];

		if(domStyle) {
			domStyle.refs++;

			for(var j = 0; j < domStyle.parts.length; j++) {
				domStyle.parts[j](item.parts[j]);
			}

			for(; j < item.parts.length; j++) {
				domStyle.parts.push(addStyle(item.parts[j], options));
			}
		} else {
			var parts = [];

			for(var j = 0; j < item.parts.length; j++) {
				parts.push(addStyle(item.parts[j], options));
			}

			stylesInDom[item.id] = {id: item.id, refs: 1, parts: parts};
		}
	}
}

function listToStyles (list, options) {
	var styles = [];
	var newStyles = {};

	for (var i = 0; i < list.length; i++) {
		var item = list[i];
		var id = options.base ? item[0] + options.base : item[0];
		var css = item[1];
		var media = item[2];
		var sourceMap = item[3];
		var part = {css: css, media: media, sourceMap: sourceMap};

		if(!newStyles[id]) styles.push(newStyles[id] = {id: id, parts: [part]});
		else newStyles[id].parts.push(part);
	}

	return styles;
}

function insertStyleElement (options, style) {
	var target = getElement(options.insertInto)

	if (!target) {
		throw new Error("Couldn't find a style target. This probably means that the value for the 'insertInto' parameter is invalid.");
	}

	var lastStyleElementInsertedAtTop = stylesInsertedAtTop[stylesInsertedAtTop.length - 1];

	if (options.insertAt === "top") {
		if (!lastStyleElementInsertedAtTop) {
			target.insertBefore(style, target.firstChild);
		} else if (lastStyleElementInsertedAtTop.nextSibling) {
			target.insertBefore(style, lastStyleElementInsertedAtTop.nextSibling);
		} else {
			target.appendChild(style);
		}
		stylesInsertedAtTop.push(style);
	} else if (options.insertAt === "bottom") {
		target.appendChild(style);
	} else if (typeof options.insertAt === "object" && options.insertAt.before) {
		var nextSibling = getElement(options.insertAt.before, target);
		target.insertBefore(style, nextSibling);
	} else {
		throw new Error("[Style Loader]\n\n Invalid value for parameter 'insertAt' ('options.insertAt') found.\n Must be 'top', 'bottom', or Object.\n (https://github.com/webpack-contrib/style-loader#insertat)\n");
	}
}

function removeStyleElement (style) {
	if (style.parentNode === null) return false;
	style.parentNode.removeChild(style);

	var idx = stylesInsertedAtTop.indexOf(style);
	if(idx >= 0) {
		stylesInsertedAtTop.splice(idx, 1);
	}
}

function createStyleElement (options) {
	var style = document.createElement("style");

	if(options.attrs.type === undefined) {
		options.attrs.type = "text/css";
	}

	if(options.attrs.nonce === undefined) {
		var nonce = getNonce();
		if (nonce) {
			options.attrs.nonce = nonce;
		}
	}

	addAttrs(style, options.attrs);
	insertStyleElement(options, style);

	return style;
}

function createLinkElement (options) {
	var link = document.createElement("link");

	if(options.attrs.type === undefined) {
		options.attrs.type = "text/css";
	}
	options.attrs.rel = "stylesheet";

	addAttrs(link, options.attrs);
	insertStyleElement(options, link);

	return link;
}

function addAttrs (el, attrs) {
	Object.keys(attrs).forEach(function (key) {
		el.setAttribute(key, attrs[key]);
	});
}

function getNonce() {
	if (false) {}

	return __webpack_require__.nc;
}

function addStyle (obj, options) {
	var style, update, remove, result;

	// If a transform function was defined, run it on the css
	if (options.transform && obj.css) {
	    result = typeof options.transform === 'function'
		 ? options.transform(obj.css) 
		 : options.transform.default(obj.css);

	    if (result) {
	    	// If transform returns a value, use that instead of the original css.
	    	// This allows running runtime transformations on the css.
	    	obj.css = result;
	    } else {
	    	// If the transform function returns a falsy value, don't add this css.
	    	// This allows conditional loading of css
	    	return function() {
	    		// noop
	    	};
	    }
	}

	if (options.singleton) {
		var styleIndex = singletonCounter++;

		style = singleton || (singleton = createStyleElement(options));

		update = applyToSingletonTag.bind(null, style, styleIndex, false);
		remove = applyToSingletonTag.bind(null, style, styleIndex, true);

	} else if (
		obj.sourceMap &&
		typeof URL === "function" &&
		typeof URL.createObjectURL === "function" &&
		typeof URL.revokeObjectURL === "function" &&
		typeof Blob === "function" &&
		typeof btoa === "function"
	) {
		style = createLinkElement(options);
		update = updateLink.bind(null, style, options);
		remove = function () {
			removeStyleElement(style);

			if(style.href) URL.revokeObjectURL(style.href);
		};
	} else {
		style = createStyleElement(options);
		update = applyToTag.bind(null, style);
		remove = function () {
			removeStyleElement(style);
		};
	}

	update(obj);

	return function updateStyle (newObj) {
		if (newObj) {
			if (
				newObj.css === obj.css &&
				newObj.media === obj.media &&
				newObj.sourceMap === obj.sourceMap
			) {
				return;
			}

			update(obj = newObj);
		} else {
			remove();
		}
	};
}

var replaceText = (function () {
	var textStore = [];

	return function (index, replacement) {
		textStore[index] = replacement;

		return textStore.filter(Boolean).join('\n');
	};
})();

function applyToSingletonTag (style, index, remove, obj) {
	var css = remove ? "" : obj.css;

	if (style.styleSheet) {
		style.styleSheet.cssText = replaceText(index, css);
	} else {
		var cssNode = document.createTextNode(css);
		var childNodes = style.childNodes;

		if (childNodes[index]) style.removeChild(childNodes[index]);

		if (childNodes.length) {
			style.insertBefore(cssNode, childNodes[index]);
		} else {
			style.appendChild(cssNode);
		}
	}
}

function applyToTag (style, obj) {
	var css = obj.css;
	var media = obj.media;

	if(media) {
		style.setAttribute("media", media)
	}

	if(style.styleSheet) {
		style.styleSheet.cssText = css;
	} else {
		while(style.firstChild) {
			style.removeChild(style.firstChild);
		}

		style.appendChild(document.createTextNode(css));
	}
}

function updateLink (link, options, obj) {
	var css = obj.css;
	var sourceMap = obj.sourceMap;

	/*
		If convertToAbsoluteUrls isn't defined, but sourcemaps are enabled
		and there is no publicPath defined then lets turn convertToAbsoluteUrls
		on by default.  Otherwise default to the convertToAbsoluteUrls option
		directly
	*/
	var autoFixUrls = options.convertToAbsoluteUrls === undefined && sourceMap;

	if (options.convertToAbsoluteUrls || autoFixUrls) {
		css = fixUrls(css);
	}

	if (sourceMap) {
		// http://stackoverflow.com/a/26603875
		css += "\n/*# sourceMappingURL=data:application/json;base64," + btoa(unescape(encodeURIComponent(JSON.stringify(sourceMap)))) + " */";
	}

	var blob = new Blob([css], { type: "text/css" });

	var oldSrc = link.href;

	link.href = URL.createObjectURL(blob);

	if(oldSrc) URL.revokeObjectURL(oldSrc);
}


/***/ }),

/***/ "./node_modules/style-loader/lib/urls.js":
/*!***********************************************!*\
  !*** ./node_modules/style-loader/lib/urls.js ***!
  \***********************************************/
/*! no static exports found */
/***/ (function(module, exports) {


/**
 * When source maps are enabled, `style-loader` uses a link element with a data-uri to
 * embed the css on the page. This breaks all relative urls because now they are relative to a
 * bundle instead of the current page.
 *
 * One solution is to only use full urls, but that may be impossible.
 *
 * Instead, this function "fixes" the relative urls to be absolute according to the current page location.
 *
 * A rudimentary test suite is located at `test/fixUrls.js` and can be run via the `npm test` command.
 *
 */

module.exports = function (css) {
  // get current location
  var location = typeof window !== "undefined" && window.location;

  if (!location) {
    throw new Error("fixUrls requires window.location");
  }

	// blank or null?
	if (!css || typeof css !== "string") {
	  return css;
  }

  var baseUrl = location.protocol + "//" + location.host;
  var currentDir = baseUrl + location.pathname.replace(/\/[^\/]*$/, "/");

	// convert each url(...)
	/*
	This regular expression is just a way to recursively match brackets within
	a string.

	 /url\s*\(  = Match on the word "url" with any whitespace after it and then a parens
	   (  = Start a capturing group
	     (?:  = Start a non-capturing group
	         [^)(]  = Match anything that isn't a parentheses
	         |  = OR
	         \(  = Match a start parentheses
	             (?:  = Start another non-capturing groups
	                 [^)(]+  = Match anything that isn't a parentheses
	                 |  = OR
	                 \(  = Match a start parentheses
	                     [^)(]*  = Match anything that isn't a parentheses
	                 \)  = Match a end parentheses
	             )  = End Group
              *\) = Match anything and then a close parens
          )  = Close non-capturing group
          *  = Match anything
       )  = Close capturing group
	 \)  = Match a close parens

	 /gi  = Get all matches, not the first.  Be case insensitive.
	 */
	var fixedCss = css.replace(/url\s*\(((?:[^)(]|\((?:[^)(]+|\([^)(]*\))*\))*)\)/gi, function(fullMatch, origUrl) {
		// strip quotes (if they exist)
		var unquotedOrigUrl = origUrl
			.trim()
			.replace(/^"(.*)"$/, function(o, $1){ return $1; })
			.replace(/^'(.*)'$/, function(o, $1){ return $1; });

		// already a full url? no change
		if (/^(#|data:|http:\/\/|https:\/\/|file:\/\/\/|\s*$)/i.test(unquotedOrigUrl)) {
		  return fullMatch;
		}

		// convert the url to a full url
		var newUrl;

		if (unquotedOrigUrl.indexOf("//") === 0) {
		  	//TODO: should we add protocol?
			newUrl = unquotedOrigUrl;
		} else if (unquotedOrigUrl.indexOf("/") === 0) {
			// path should be relative to the base url
			newUrl = baseUrl + unquotedOrigUrl; // already starts with '/'
		} else {
			// path should be relative to current directory
			newUrl = currentDir + unquotedOrigUrl.replace(/^\.\//, ""); // Strip leading './'
		}

		// send back the fixed url(...)
		return "url(" + JSON.stringify(newUrl) + ")";
	});

	// send back the fixed css
	return fixedCss;
};


/***/ }),

/***/ "./src/i18n/en.json":
/*!**************************!*\
  !*** ./src/i18n/en.json ***!
  \**************************/
/*! exports provided: OneHour, OneDay, SevenDays, OptionDate, To:, Confirm, Cancel, DataDetail, default */
/***/ (function(module) {

module.exports = {"OneHour":"Last 1 Hour","OneDay":"Last 1 Day","SevenDays":"Last 7 Days","OptionDate":"Custom Timeframe","To:":"To","Confirm":"Confirm","Cancel":"Cancel","DataDetail":"Data Information"};

/***/ }),

/***/ "./src/i18n/index.ts":
/*!***************************!*\
  !*** ./src/i18n/index.ts ***!
  \***************************/
/*! exports provided: zh, en */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _zh_json__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./zh.json */ "./src/i18n/zh.json");
var _zh_json__WEBPACK_IMPORTED_MODULE_0___namespace = /*#__PURE__*/__webpack_require__.t(/*! ./zh.json */ "./src/i18n/zh.json", 1);
/* harmony reexport (module object) */ __webpack_require__.d(__webpack_exports__, "zh", function() { return _zh_json__WEBPACK_IMPORTED_MODULE_0__; });
/* harmony import */ var _en_json__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./en.json */ "./src/i18n/en.json");
var _en_json__WEBPACK_IMPORTED_MODULE_1___namespace = /*#__PURE__*/__webpack_require__.t(/*! ./en.json */ "./src/i18n/en.json", 1);
/* harmony reexport (module object) */ __webpack_require__.d(__webpack_exports__, "en", function() { return _en_json__WEBPACK_IMPORTED_MODULE_1__; });





/***/ }),

/***/ "./src/i18n/zh.json":
/*!**************************!*\
  !*** ./src/i18n/zh.json ***!
  \**************************/
/*! exports provided: OneHour, OneDay, SevenDays, OptionDate, To, Confirm, Cancel, DataDetail, default */
/***/ (function(module) {

module.exports = {"OneHour":"近1小时","OneDay":"近1天","SevenDays":"近7天","OptionDate":"选择日期","To":"至","Confirm":"确定","Cancel":"取消","DataDetail":"数据详情"};

/***/ }),

/***/ "./src/panel/ChartsComponents.tsx":
/*!****************************************!*\
  !*** ./src/panel/ChartsComponents.tsx ***!
  \****************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _containers_Pure__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./containers/Pure */ "./src/panel/containers/Pure.tsx");
/* harmony import */ var _containers_Filter__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./containers/Filter */ "./src/panel/containers/Filter.tsx");
/* harmony import */ var _containers_Instances__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./containers/Instances */ "./src/panel/containers/Instances.tsx");




const components = { ChartPanel: _containers_Pure__WEBPACK_IMPORTED_MODULE_1__["ChartPanel"], ChartFilterPanel: _containers_Filter__WEBPACK_IMPORTED_MODULE_2__["ChartFilterPanel"], ChartInstancesPanel: _containers_Instances__WEBPACK_IMPORTED_MODULE_3__["ChartInstancesPanel"] };
/* harmony default export */ __webpack_exports__["default"] = (function (props) {
    return react__WEBPACK_IMPORTED_MODULE_0__["createElement"](components[props.componentName], props);
});


/***/ }),

/***/ "./src/panel/components/FilterTableChart.tsx":
/*!***************************************************!*\
  !*** ./src/panel/components/FilterTableChart.tsx ***!
  \***************************************************/
/*! exports provided: FilterTableChart */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "FilterTableChart", function() { return FilterTableChart; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @tencent/tea-component */ "@tencent/tea-component");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var tea_components_tagsearchbox_index__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! tea-components/tagsearchbox/index */ "./src/tea-components/tagsearchbox/index.ts");
/* harmony import */ var charts_index__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! charts/index */ "./src/charts/index.ts");
/* harmony import */ var core_utils__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! core/utils */ "./src/core/utils.ts");
/* harmony import */ var _PureChart__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! ./PureChart */ "./src/panel/components/PureChart.tsx");
/* harmony import */ var _helper__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! ../helper */ "./src/panel/helper.ts");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(/*! ../constants */ "./src/panel/constants.ts");
/* harmony import */ var _core__WEBPACK_IMPORTED_MODULE_8__ = __webpack_require__(/*! ../core */ "./src/panel/core.ts");
/* harmony import */ var _i18n__WEBPACK_IMPORTED_MODULE_9__ = __webpack_require__(/*! ../../i18n */ "./src/i18n/index.ts");










const version = window.VERSION || "zh";
const language = _i18n__WEBPACK_IMPORTED_MODULE_9__[version];
const { scrollable } = _tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Table"].addons;
/**
 * 生成系统聚合方式选项
 */
const AggregationOptions = Object.keys(_constants__WEBPACK_IMPORTED_MODULE_7__["QUERY"].Aggregation).map(key => {
    return {
        value: key,
        text: _constants__WEBPACK_IMPORTED_MODULE_7__["QUERY"].Aggregation[key]
    };
});
function FormatTagsParam(tagDimensions, conditions, tagConditions) {
    // 参考 TagSearchBox 返回的参数
    const dimensionItems = tagDimensions.map(item => item.attr.value);
    // condition 支持多条件查询
    const conditionItems = [].concat(conditions, tagConditions.map(item => [item.attr.value, "in", item.values.map(value => value.name)]));
    return { dimensionItems, conditionItems };
}
class FilterTableChart extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        let tagDimensions = [];
        let tagConditions = [];
        // 判断 metric 是否有 storeKey，
        // 处理 defaultGroupBy, defaultConditions 搜索条件初始化 default* 的数据结构：
        if (props.metric.storeKey) {
            const { defaultGroupBy, defaultConditions } = _helper__WEBPACK_IMPORTED_MODULE_6__["STORE"].Get(props.metric.storeKey, {
                defaultGroupBy: [],
                defaultConditions: []
            });
            tagDimensions = defaultGroupBy.map(item => {
                const { groupBy, values } = item;
                return {
                    attr: {
                        type: "onlyKey",
                        key: groupBy.value,
                        name: groupBy.name,
                        value: groupBy.value
                    },
                    values
                };
            });
            tagConditions = defaultConditions.map(item => {
                const { groupBy, values } = item;
                return {
                    attr: {
                        type: "input",
                        key: groupBy.value,
                        name: groupBy.name,
                        value: groupBy.value
                    },
                    values
                };
            });
        }
        this.state = {
            loading: false,
            tagDimensions: tagDimensions,
            tagConditions: tagConditions,
            aggregation: AggregationOptions[0].value,
            period: props.periodOptions[0].value,
            labels: [],
            lines: []
        };
    }
    componentWillMount() {
        // 加载数据
        this.fetchChartData();
    }
    componentDidUpdate(prevProps, prevState, snapshot) {
        let hasUpdate = false;
        let hasFetch = false;
        const { table, startTime, endTime, metric, conditions, periodOptions } = this.props;
        if (JSON.stringify(startTime) !== JSON.stringify(prevProps.startTime) ||
            JSON.stringify(endTime) !== JSON.stringify(prevProps.endTime) ||
            JSON.stringify(metric) !== JSON.stringify(prevProps.metric) ||
            JSON.stringify(conditions) !== JSON.stringify(prevProps.conditions)) {
            hasFetch = true;
        }
        // 更新 period 选项
        let period = this.state.period;
        if (JSON.stringify(periodOptions) !== JSON.stringify(prevProps.periodOptions)) {
            period = periodOptions[0].value;
            this.setState({ period });
            hasUpdate = true;
        }
        const { dimensionItems, conditionItems } = FormatTagsParam(this.state.tagDimensions, conditions, this.state.tagConditions);
        if (JSON.stringify(prevState) !== JSON.stringify(this.state)) {
            hasUpdate = true;
        }
        if (hasFetch) {
            this.requestChartData(table, startTime, endTime, metric, dimensionItems, conditionItems, period);
        }
        return hasUpdate || hasFetch;
    }
    fetchChartData() {
        const { table, startTime, endTime, metric, conditions } = this.props;
        const { tagDimensions, tagConditions, period } = this.state;
        const { dimensionItems, conditionItems } = FormatTagsParam(tagDimensions, conditions, tagConditions);
        this.requestChartData(table, startTime, endTime, metric, dimensionItems, conditionItems, period);
    }
    requestChartData(table, startTime, endTime, metric, dimensions, conditions, period) {
        this.setState({ loading: true, labels: [], lines: [] });
        // 生成请求字段
        const aggregationFields = [Object.assign(Object.assign({}, metric), { expr: `${this.state.aggregation}(${metric.expr})` })];
        _core__WEBPACK_IMPORTED_MODULE_8__["CHART_PANEL"].RequestData({
            table,
            startTime,
            endTime,
            fields: aggregationFields,
            dimensions: dimensions,
            conditions: conditions,
            period: period
        })
            .then(res => {
            const { columns, data } = res;
            // 根据维度聚合返回数据
            const dimensionsData = _core__WEBPACK_IMPORTED_MODULE_8__["CHART_PANEL"].AggregateDataByDimensions(dimensions, columns, data);
            /**
             * 生成 labels
             */
            const timestampIndex = columns.indexOf(`timestamp(${period}s)`);
            // 默认表格第一列为时间序列, 返回的时间列表可能不是startTime开始，以endTime结束，需要补帧
            let labels = (Object.values(dimensionsData)[0] || []).map(item => item[timestampIndex]);
            labels = _core__WEBPACK_IMPORTED_MODULE_8__["CHART_PANEL"].OffsetTimeSeries(startTime, endTime, parseInt(period), labels);
            /**
             * 生成 lines
             */
            const metricIndex = _core__WEBPACK_IMPORTED_MODULE_8__["CHART_PANEL"].FindMetricIndex(columns, `${metric.expr}_${this.state.aggregation}`); // 数据列下标
            const valueTransform = metric.valueTransform || (value => value);
            // 每个图表中显示groupBy类别一个field -> chart
            let lines = [];
            Object.keys(dimensionsData).forEach(dimensionKey => {
                // 每个dimension的数据集， fieldsIndex获取该dimension列下标的值
                const sourceValue = dimensionsData[dimensionKey];
                const data = {};
                sourceValue.forEach(row => {
                    let value = row[metricIndex];
                    // data 记录 时间戳row[timestampIndex] 对应 value
                    data[row[timestampIndex]] = valueTransform(value);
                });
                let line = _core__WEBPACK_IMPORTED_MODULE_8__["CHART_PANEL"].GenerateLine(dimensionKey, data);
                lines.push(line);
            });
            this.setState({ loading: false, labels, lines });
        })
            .catch(e => {
            this.setState({ loading: false });
        });
    }
    updateState(options) {
        this.setState(Object.assign({}, options), () => {
            this.fetchChartData();
        });
    }
    saveTagsByStoreKey(saveKey, tags) {
        if (this.props.metric.storeKey) {
            let defaultValues = _helper__WEBPACK_IMPORTED_MODULE_6__["STORE"].Get(this.props.metric.storeKey, { defaultGroupBy: [], defaultConditions: [] });
            defaultValues[saveKey] = tags.map(item => {
                return {
                    groupBy: item.attr,
                    values: item.values
                };
            });
            _helper__WEBPACK_IMPORTED_MODULE_6__["STORE"].Set(this.props.metric.storeKey, defaultValues);
        }
    }
    render() {
        const { loading, tagDimensions, tagConditions, period, aggregation, labels, lines } = this.state;
        const { dimensions, periodOptions, chartId, width, height, metric, startTime, endTime } = this.props;
        // 生成维度和筛选条件选项
        const dimensionOptions = dimensions.map(item => {
            return {
                type: "onlyKey",
                key: item.value,
                name: item.name,
                value: item.value
            };
        });
        const conditionOptions = dimensions.map(item => {
            return {
                type: "input",
                key: item.value,
                name: item.name,
                value: item.value
            };
        });
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: { width: "100%" } },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Form"], { className: "tea-form--vertical tea-form--inline tea-mt-2n" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Form"].Item, { label: "\u7EF4\u5EA6", style: { width: "calc(100% - 300px)" } },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](tea_components_tagsearchbox_index__WEBPACK_IMPORTED_MODULE_2__["TagSearchBox"], { minWidth: "100%", attributes: dimensionOptions, tipZh: "请选择维度", value: tagDimensions, onChange: tags => {
                            if (dimensionOptions.length === 0) {
                                return;
                            }
                            const tagDimensions = tags.map(item => {
                                return {
                                    attr: item.attr,
                                    values: item.values
                                };
                            });
                            // 如果指标传递了storeKey，则缓存用户维度的选项，方便用户下次操作
                            this.saveTagsByStoreKey("defaultGroupBy", tags);
                            this.updateState({ tagDimensions });
                        } }),
                    "/>"),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Form"].Item, { label: "\u7EDF\u8BA1\u7C92\u5EA6" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Select"], { style: { paddingTop: 6 }, type: "simulate", appearence: "button", options: periodOptions, value: period, onChange: value => {
                            this.updateState({ period: value });
                        }, placeholder: "\u8BF7\u9009\u62E9" })),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Form"].Item, { label: "\u7EDF\u8BA1\u65B9\u5F0F" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Select"], { style: { paddingTop: 6 }, type: "simulate", appearence: "button", options: AggregationOptions, value: aggregation, onChange: value => {
                            this.updateState({ aggregation: value });
                        }, placeholder: "\u8BF7\u9009\u62E9" }))),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Form"], { className: "tea-form--vertical tea-form--inline tea-mt-2n" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Form"].Item, { label: "\u7B5B\u9009\u6761\u4EF6", style: { width: "100%" } },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["TagSearchBox"], { minWidth: "100%", attributes: conditionOptions, value: tagConditions, onChange: tags => {
                            if (conditionOptions.length === 0) {
                                return;
                            }
                            const tagConditions = tags.map(item => {
                                return {
                                    attr: item.attr,
                                    values: item.values
                                };
                            });
                            this.saveTagsByStoreKey("defaultConditions", tags);
                            this.updateState({ tagConditions });
                        } }))),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](TableChart, { id: chartId, loading: loading, min: typeof startTime === "string" ? 0 : startTime.getTime(), max: typeof endTime === "string" ? 0 : endTime.getTime(), width: width, height: height, title: metric.alias, field: metric, labels: labels, unit: metric.unit, tooltipLabels: metric.tooltipLabels, lines: lines, reload: this.fetchChartData.bind(this) })));
    }
}
/**
 * 表格型图表
 * hover 图表时数据通过 table 显示，隐藏 tooltip
 */
class TableChart extends _PureChart__WEBPACK_IMPORTED_MODULE_5__["default"] {
    constructor(props) {
        super(props);
        this._scrollAnchorRefs = {};
        this.state = {
            selectedLabel: 0,
            selectedRowKeys: []
        };
    }
    componentDidMount() {
        const { loading, min, max, title, reload, labels, lines, field, tooltipLabels } = this.props;
        this._chartEntry = new charts_index__WEBPACK_IMPORTED_MODULE_3__["default"](this.props.id, {
            width: this.chartWidth,
            height: this.chartHeight,
            paddingHorizontal: 15,
            paddingVertical: 30,
            isSeriesTime: true,
            showLegend: false,
            showTooltip: false,
            hoverPointData: this.mouseOverChart.bind(this),
            loading,
            min,
            max,
            title,
            labels,
            reload,
            tooltipLabels,
            colorTheme: field.colorTheme,
            colors: field.colors,
            yAxis: field.scale || [],
            series: lines,
            isKilobyteFormat: field.thousands === _constants__WEBPACK_IMPORTED_MODULE_7__["Kilobyte"],
            unit: field.unit
        }, field.chartType);
    }
    shouldComponentUpdate(nextProps, nextState) {
        if (JSON.stringify(nextProps) !== JSON.stringify(this.props)) {
            const { loading, min, max, title, reload, labels, lines, field, tooltipLabels } = nextProps;
            this._chartEntry &&
                this._chartEntry.setType(field.chartType, {
                    width: this.chartWidth,
                    height: this.chartHeight,
                    loading,
                    min,
                    max,
                    title,
                    labels,
                    reload,
                    tooltipLabels,
                    colorTheme: field.colorTheme,
                    colors: field.colors,
                    yAxis: field.scale || [],
                    series: lines,
                    isKilobyteFormat: field.thousands === _constants__WEBPACK_IMPORTED_MODULE_7__["Kilobyte"],
                    unit: field.unit
                });
            return true;
        }
        if (JSON.stringify(nextState) !== JSON.stringify(this.state)) {
            return true;
        }
        return false;
    }
    // 图表鼠标悬浮时回调函数
    mouseOverChart(params) {
        const { xAxisTickMarkIndex, mousePosition, content } = params;
        // 筛选出被hover的折线
        const selectedRowKeys = content.filter(item => item.hover).map(item => item.legend);
        // table 中对应数据进行滚动显示
        if (selectedRowKeys.length > 0) {
            const element = this._scrollAnchorRefs[selectedRowKeys[selectedRowKeys.length - 1]].current;
            // firefox,IE 不支持 scrollIntoViewIfNeeded
            if (element.scrollIntoViewIfNeeded) {
                element.scrollIntoViewIfNeeded(false);
            }
            else {
                element.scrollIntoView({ block: "end", behavior: "smooth" });
            }
        }
        this.setState({ selectedLabel: xAxisTickMarkIndex, selectedRowKeys });
    }
    // 根据鼠标选中的label，获取各折线对应的点，显示在表格中
    getTableRecords(label, lines) {
        const tableRecords = lines.map(line => {
            return {
                key: core_utils__WEBPACK_IMPORTED_MODULE_4__["FormatStringNoHTMLSharp"](line.legend),
                value: line.data[label]
            };
        });
        tableRecords.sort((a, b) => b.value - a.value);
        return tableRecords;
    }
    // hover 在表格的事件
    hoverTableEvent(legend) {
        this._chartEntry && this._chartEntry.highlightLine(legend);
    }
    render() {
        const { loading, labels, lines } = this.props;
        const tableRecords = this.getTableRecords(labels[this.state.selectedLabel], lines);
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: { width: "100%" } },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: { width: "100%" }, id: this.props.id }),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tea-justify-grid tea-mt-2n" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tea-justify-grid__col tea-justify-grid__col--left" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("h3", { className: "tea-h3", style: { fontSize: 14 } },
                        language.DataDetail,
                        " \u00A0\u00A0",
                        labels[this.state.selectedLabel] &&
                            core_utils__WEBPACK_IMPORTED_MODULE_4__["TIME"].Format(labels[this.state.selectedLabel], core_utils__WEBPACK_IMPORTED_MODULE_4__["TIME"].DateFormat.fullDateTime)))),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: { width: "calc(100% - 40px)", margin: "0 auto" } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Table"], { columns: [
                        {
                            key: "key",
                            header: "Source",
                            render: record => {
                                this._scrollAnchorRefs[record.key] = react__WEBPACK_IMPORTED_MODULE_0__["createRef"]();
                                return react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { ref: this._scrollAnchorRefs[record.key] }, record.key);
                            }
                        },
                        {
                            key: "value",
                            header: "Value"
                        }
                    ], records: tableRecords, recordKey: "key", rowClassName: record => (this.state.selectedRowKeys.indexOf(record.key) !== -1 ? "is-selected" : ""), topTip: loading && (react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["StatusTip"]
                    // @ts-ignore
                    , { 
                        // @ts-ignore
                        status: loading ? "loading" : "none" })), addons: [
                        // 支持表格滚动，高度超过 180 开始显示滚动条
                        scrollable({
                            maxHeight: 180,
                            onScrollBottom: () => { }
                        }),
                        hoverRowEvent({ event: this.hoverTableEvent.bind(this) })
                    ] }))));
    }
}
/**
 * 注册 table hover 事件
 */
function hoverRowEvent(options) {
    return {
        onInjectRow: next => (...args) => {
            const result = next(...args);
            const [record, rowKey] = args;
            return Object.assign(Object.assign({}, result), { row: react__WEBPACK_IMPORTED_MODULE_0__["cloneElement"](result.row, {
                    onMouseOver: e => {
                        options.event(rowKey);
                    },
                    onMouseLeave: e => {
                        options.event("");
                    }
                }) });
        }
    };
}


/***/ }),

/***/ "./src/panel/components/InstanceList.tsx":
/*!***********************************************!*\
  !*** ./src/panel/components/InstanceList.tsx ***!
  \***********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return InstanceList; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");

class InstanceList extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.state = {
            checkedAll: props.list.every((item) => item.isChecked),
        };
    }
    onCheck(instance, event) {
        const list = [].concat(this.props.list.map(item => {
            if (item === instance) {
                item.isChecked = event.target.checked;
            }
            return item;
        }));
        const checkedAll = list.every((item) => item.isChecked);
        this.setState({ checkedAll }, () => {
            this.props.update(list);
        });
    }
    ;
    onCheckAll(event) {
        const list = [].concat(this.props.list.map(item => {
            item.isChecked = event.target.checked;
            return item;
        }));
        this.setState({ checkedAll: event.target.checked }, () => {
            this.props.update(list);
        });
    }
    render() {
        const { className, style } = this.props;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: className, style: style },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", null,
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("table", { className: "Tchart_table-box" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("thead", null,
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("tr", null,
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", { className: "" },
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-first-checkbox" },
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("input", { type: "checkbox", className: "tc-15-checkbox", checked: this.state.checkedAll, onChange: this.onCheckAll.bind(this) }))),
                            this.props.columns.map(item => {
                                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", { className: "", key: item.key },
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", null,
                                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { className: "text-overflow" }, item.name))));
                            }))),
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("tbody", null, this.props.list.map((item, index) => {
                        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("tr", { key: index },
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", { className: "" },
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-first-checkbox" },
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("input", { type: "checkbox", className: "tc-15-checkbox", checked: item.isChecked, onChange: this.onCheck.bind(this, item) }))),
                            this.props.columns.map(column => {
                                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", { className: "", key: column.key },
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", null, column.render ? column.render(item) : item[column.key])));
                            })));
                    }))))));
    }
}


/***/ }),

/***/ "./src/panel/components/MetricCharts.tsx":
/*!***********************************************!*\
  !*** ./src/panel/components/MetricCharts.tsx ***!
  \***********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return MetricCharts; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var charts_index__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! charts/index */ "./src/charts/index.ts");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../constants */ "./src/panel/constants.ts");



class MetricCharts extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.ChartPageSize = 3; // panel 为 list 状态时每次滚动加载 chart 数目
        this._id = "MetricCharts";
        this._charts = [];
        this._chartHeight = 388;
        this._chartWidth = 726;
        this.state = {
            chartEndIndex: this.ChartPageSize
        };
    }
    componentDidMount() {
        this.scrollUpdateCharts(this.props, this.state.chartEndIndex);
    }
    shouldComponentUpdate(nextProps, nextState) {
        if (nextProps.seriesGroup.length !== this.props.seriesGroup.length
            || nextState.chartEndIndex !== this.state.chartEndIndex) {
            this.scrollUpdateCharts(nextProps, nextState.chartEndIndex);
            return true;
        }
        if (JSON.stringify(nextProps.seriesGroup) !== JSON.stringify(this.props.seriesGroup)
            || nextProps.min !== this.props.min
            || nextProps.max !== this.props.max
            || nextProps.loading !== this.props.loading) {
            this.updateCharts(nextProps, nextState.chartEndIndex);
            return true;
        }
        return false;
    }
    updateCharts(props, chartEndIndex) {
        const { loading, min, max, seriesGroup, reload } = props;
        const length = Math.min(seriesGroup.length, chartEndIndex);
        for (let i = 0; i < length; i++) {
            const series = seriesGroup[i].lines;
            const title = seriesGroup[i].title;
            const labels = seriesGroup[i].labels;
            const field = seriesGroup[i].field;
            const tooltipLabels = seriesGroup[i].tooltipLabels;
            if (this._charts[i]) {
                this._charts[i].setType(field.chartType, {
                    loading,
                    min,
                    max,
                    title,
                    labels,
                    series,
                    reload,
                    tooltipLabels,
                    colorTheme: field.colorTheme,
                    colors: field.colors,
                    yAxis: field.scale || [],
                    isKilobyteFormat: field.thousands === _constants__WEBPACK_IMPORTED_MODULE_2__["Kilobyte"],
                    unit: field.unit
                });
            }
            else {
                this._charts.push(new charts_index__WEBPACK_IMPORTED_MODULE_1__["default"](`${this._id}_${i}`, {
                    width: this._chartWidth,
                    height: this._chartHeight,
                    paddingHorizontal: 15,
                    paddingVertical: 30,
                    isSeriesTime: true,
                    showLegend: false,
                    isKilobyteFormat: field.thousands === _constants__WEBPACK_IMPORTED_MODULE_2__["Kilobyte"],
                    title,
                    min,
                    max,
                    tooltipLabels,
                    colorTheme: field.colorTheme,
                    colors: field.colors,
                    yAxis: field.scale || [],
                    unit: field.unit,
                    loading,
                    labels,
                    series,
                    reload
                }, field.chartType));
            }
        }
    }
    scrollUpdateCharts(props, chartEndIndex) {
        const { loading, min, max, seriesGroup, reload } = props;
        const startIndex = this._charts.length;
        const length = Math.min(seriesGroup.length, chartEndIndex);
        for (let i = startIndex; i < length; i++) {
            const series = seriesGroup[i].lines;
            const title = seriesGroup[i].title;
            const labels = seriesGroup[i].labels;
            const field = seriesGroup[i].field;
            const tooltipLabels = seriesGroup[i].tooltipLabels;
            this._charts.push(new charts_index__WEBPACK_IMPORTED_MODULE_1__["default"](`${this._id}_${i}`, {
                width: this._chartWidth,
                height: this._chartHeight,
                paddingHorizontal: 15,
                paddingVertical: 30,
                isSeriesTime: true,
                showLegend: false,
                isKilobyteFormat: field.thousands === _constants__WEBPACK_IMPORTED_MODULE_2__["Kilobyte"],
                title,
                min,
                max,
                tooltipLabels,
                yAxis: field.scale || [],
                unit: field.unit,
                loading,
                labels,
                series,
                reload
            }, field.chartType));
        }
    }
    render() {
        const { className, style } = this.props;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: className, style: Object.assign({}, style), onScroll: (e) => {
                // 图表滚动加载
                let element = e.target;
                // 向下滚动加载
                // 提前加载
                const index = this.state.chartEndIndex - 1;
                if (Math.floor(this._chartHeight * index - element.scrollTop) <= element.clientHeight) {
                    if (this.state.chartEndIndex < this.props.seriesGroup.length) {
                        const nextChartEndIndex = this.state.chartEndIndex + this.ChartPageSize;
                        this.setState({ chartEndIndex: nextChartEndIndex });
                    }
                }
            } }, this.props.seriesGroup.map((item, index) => {
            const id = `${this._id}_${index}`;
            return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { key: id, style: { width: "100%", height: this._chartHeight } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { id: id })));
        })));
    }
}


/***/ }),

/***/ "./src/panel/components/PureChart.tsx":
/*!********************************************!*\
  !*** ./src/panel/components/PureChart.tsx ***!
  \********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "default", function() { return PureChart; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var charts_index__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! charts/index */ "./src/charts/index.ts");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../constants */ "./src/panel/constants.ts");



/**
 * 基础图表组件，支持缩放
 */
class PureChart extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this._width = 0;
        this._height = 0;
        this._chartEntry = null;
        this.onResize = () => {
            this._chartEntry && this._chartEntry.setSize(this.chartWidth, this.chartHeight);
        };
    }
    get chartWidth() {
        if (this.props.width) {
            return this.props.width;
        }
        if (!document.getElementById(this.props.id)) {
            this._width = _constants__WEBPACK_IMPORTED_MODULE_2__["CHART"].DefaultSize.width;
        }
        else {
            const width = document.getElementById(this.props.id).clientWidth;
            this._width = width > 0 ? width : _constants__WEBPACK_IMPORTED_MODULE_2__["CHART"].DefaultSize.width;
        }
        return this._width;
    }
    get chartHeight() {
        if (this.props.height) {
            return this.props.height;
        }
        if (this._width < _constants__WEBPACK_IMPORTED_MODULE_2__["CHART"].DefaultSize.width) {
            this._height = _constants__WEBPACK_IMPORTED_MODULE_2__["CHART"].DefaultSize.height * this._width / _constants__WEBPACK_IMPORTED_MODULE_2__["CHART"].DefaultSize.width;
        }
        else {
            this._height = _constants__WEBPACK_IMPORTED_MODULE_2__["CHART"].DefaultSize.height;
        }
        return this._height;
    }
    componentWillMount() {
        window.addEventListener('resize', this.onResize);
    }
    componentDidMount() {
        const { loading, title, min, max, reload, labels, lines, field, tooltipLabels } = this.props;
        this._chartEntry = new charts_index__WEBPACK_IMPORTED_MODULE_1__["default"](this.props.id, {
            width: this.chartWidth,
            height: this.chartHeight,
            paddingHorizontal: 15,
            paddingVertical: 30,
            isSeriesTime: true,
            showLegend: false,
            loading,
            min,
            max,
            title,
            labels,
            reload,
            tooltipLabels,
            colorTheme: field.colorTheme,
            colors: field.colors,
            yAxis: field.scale || [],
            series: lines,
            isKilobyteFormat: field.thousands === _constants__WEBPACK_IMPORTED_MODULE_2__["Kilobyte"],
            unit: field.unit
        }, field.chartType);
    }
    shouldComponentUpdate(nextProps, nextState) {
        if (JSON.stringify(nextProps) !== JSON.stringify(this.props)) {
            const { loading, title, min, max, reload, labels, lines, field, tooltipLabels } = nextProps;
            this._chartEntry && this._chartEntry.setType(field.chartType, {
                width: this.chartWidth,
                height: this.chartHeight,
                loading,
                min,
                max,
                title,
                labels,
                reload,
                tooltipLabels,
                colorTheme: field.colorTheme,
                colors: field.colors,
                yAxis: field.scale || [],
                series: lines,
                isKilobyteFormat: field.thousands === _constants__WEBPACK_IMPORTED_MODULE_2__["Kilobyte"],
                unit: field.unit
            });
            return true;
        }
        if (JSON.stringify(nextState) !== JSON.stringify(this.state)) {
            return true;
        }
        return false;
    }
    componentWillUnmount() {
        window.removeEventListener('resize', this.onResize);
    }
    render() {
        const { style = {} } = this.props;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: Object.assign({ width: '100%', height: _constants__WEBPACK_IMPORTED_MODULE_2__["CHART"].DefaultSize.height }, style), id: this.props.id }));
    }
}


/***/ }),

/***/ "./src/panel/components/Toolbar.tsx":
/*!******************************************!*\
  !*** ./src/panel/components/Toolbar.tsx ***!
  \******************************************/
/*! exports provided: Toolbar */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Toolbar", function() { return Toolbar; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var tea_components_datetimepicker_index__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! tea-components/datetimepicker/index */ "./src/tea-components/datetimepicker/index.ts");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../constants */ "./src/panel/constants.ts");



const DurationsByPeriod = {
    // [period(s)] : day
    60: 15,
    300: 31,
    3600: 62,
    86400: 186
};
class Toolbar extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.nowDay = new Date();
        // 可选择日期最大间隔
        this.range = {
            min: new Date(+this.nowDay - 864e5 * DurationsByPeriod[300]),
            max: new Date(this.nowDay),
            maxLength: DurationsByPeriod[300] * 86400
        };
        this.onTimePickerChange = (dateTime, label) => {
            this.props.onChangeTime(new Date(dateTime.from), new Date(dateTime.to));
        };
    }
    render() {
        const { duration, style } = this.props;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-action-grid", style: Object.assign({}, style) },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "justify-grid" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "col" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar-select-wrap tc-15-calendar2-hook" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { className: "dateTimePickerVS", onClick: e => e.stopPropagation() },
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](tea_components_datetimepicker_index__WEBPACK_IMPORTED_MODULE_1__["DateTimePicker"], { tabs: _constants__WEBPACK_IMPORTED_MODULE_2__["TIME_PICKER"].Tabs, defaultSelectedTabIndex: 0, defaultValue: duration, onChange: this.onTimePickerChange, range: this.range }))),
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "monitor-dialog-bd-right", style: { display: "inline-block", marginLeft: 15 } }, this.props.children)))));
    }
}


/***/ }),

/***/ "./src/panel/constants.ts":
/*!********************************!*\
  !*** ./src/panel/constants.ts ***!
  \********************************/
/*! exports provided: TIME_PICKER, QUERY, CHART, Kilobyte, OneDayMillisecond */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "TIME_PICKER", function() { return TIME_PICKER; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "QUERY", function() { return QUERY; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "CHART", function() { return CHART; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Kilobyte", function() { return Kilobyte; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "OneDayMillisecond", function() { return OneDayMillisecond; });
/* harmony import */ var _i18n__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../i18n */ "./src/i18n/index.ts");

const version = window.VERSION || "zh";
const language = _i18n__WEBPACK_IMPORTED_MODULE_0__[version];
var TIME_PICKER;
(function (TIME_PICKER) {
    TIME_PICKER.Tabs = [
        { from: "%NOW-1h", to: "%NOW", label: language.OneHour },
        { from: "%NOW-24h", to: "%NOW", label: language.OneDay },
        { from: "%NOW-168h", to: "%NOW", label: language.SevenDays },
    ];
})(TIME_PICKER || (TIME_PICKER = {}));
;
var QUERY;
(function (QUERY) {
    QUERY.Aggregation = {
        'sum': '总和',
        'count': '统计个数',
        'max': '最大值',
        'min': '最小值',
        'avg': '平均值',
    };
    QUERY.Limit = 65535;
})(QUERY || (QUERY = {}));
var CHART;
(function (CHART) {
    CHART.DefaultSize = {
        width: 726,
        height: 388,
    };
})(CHART || (CHART = {}));
const Kilobyte = 1024;
const OneDayMillisecond = 86400000;


/***/ }),

/***/ "./src/panel/containers/Filter.less":
/*!******************************************!*\
  !*** ./src/panel/containers/Filter.less ***!
  \******************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {


var content = __webpack_require__(/*! !../../../node_modules/css-loader/dist/cjs.js!../../../node_modules/less-loader/dist/cjs.js!./Filter.less */ "./node_modules/css-loader/dist/cjs.js!./node_modules/less-loader/dist/cjs.js!./src/panel/containers/Filter.less");

if(typeof content === 'string') content = [[module.i, content, '']];

var transform;
var insertInto;



var options = {"hmr":true}

options.transform = transform
options.insertInto = undefined;

var update = __webpack_require__(/*! ../../../node_modules/style-loader/lib/addStyles.js */ "./node_modules/style-loader/lib/addStyles.js")(content, options);

if(content.locals) module.exports = content.locals;

if(false) {}

/***/ }),

/***/ "./src/panel/containers/Filter.tsx":
/*!*****************************************!*\
  !*** ./src/panel/containers/Filter.tsx ***!
  \*****************************************/
/*! exports provided: ChartFilterPanel */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChartFilterPanel", function() { return ChartFilterPanel; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @tencent/tea-component */ "@tencent/tea-component");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var _components_Toolbar__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../components/Toolbar */ "./src/panel/components/Toolbar.tsx");
/* harmony import */ var _components_FilterTableChart__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../components/FilterTableChart */ "./src/panel/components/FilterTableChart.tsx");
/* harmony import */ var _core__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../core */ "./src/panel/core.ts");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! ../constants */ "./src/panel/constants.ts");
/* harmony import */ var _helper__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! ../helper */ "./src/panel/helper.ts");







__webpack_require__(/*! ./Filter.less */ "./src/panel/containers/Filter.less");
/**
 * 可过滤图表
 * 每个指标单独的 维度与条件查询，有tab或list展现形式
 */
class ChartFilterPanel extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.ChartPageSize = 3; // panel 为 list 状态时每次滚动加载 chart 数目
        this.ShowDivStyle = { flexShrink: 0, width: "100%", opacity: 1, transition: "opacity .45s" };
        this.HiddenDivStyle = {
            height: 0,
            overflow: "hidden",
            opacity: 0,
            flexShrink: 0,
            width: "100%",
            transition: "opacity .45s"
        };
        // 缓存每次选择不同时间段的时间粒度选项
        this.periodOptions = [];
        this.tabs = [];
        // 默认获取一天的数据
        const endTime = new Date();
        const startTime = new Date(endTime.getTime() - 1000 * 60 * 60);
        this.periodOptions = _core__WEBPACK_IMPORTED_MODULE_4__["CHART_PANEL"].GeneratePeriodOptions(startTime, endTime);
        let chartList = [];
        props.tables.forEach((tableInfo, i) => {
            // 一般一个table都是相同维度
            let { groupBy = [], conditions = [] } = tableInfo;
            // this.props.groupBy 为统一的维度值，tables中每个table可以自定义维度查询值
            groupBy = this.props.groupBy.concat(groupBy);
            tableInfo.fields.forEach((field, j) => {
                const id = `${field.expr}_${i}_${j}`;
                const tooltipLabels = (value, from) => {
                    return field.valueLabels ? field.valueLabels(value, from) : `${Object(_helper__WEBPACK_IMPORTED_MODULE_6__["TransformField"])(value || 0, field.thousands, 3)}`;
                };
                this.tabs.push({ id, label: field.alias });
                chartList.push({ id, table: tableInfo.table, field, groupBy, conditions, tooltipLabels });
            });
        });
        this.state = {
            tabStatus: true,
            activeId: this.tabs.length > 0 ? this.tabs[0].id : "",
            chartEndIndex: this.ChartPageSize,
            startTime,
            endTime,
            chartList
        };
    }
    componentWillReceiveProps(nextProps) {
        if (JSON.stringify(this.props.tables) !== JSON.stringify(nextProps.tables)) {
            let chartList = [];
            nextProps.tables.forEach((tableInfo, i) => {
                // 一般一个table都是相同维度
                let { groupBy = [], conditions = [] } = tableInfo;
                groupBy = nextProps.groupBy.concat(groupBy);
                tableInfo.fields.forEach((field, j) => {
                    const id = `${field.expr}_${i}_${j}`;
                    const tooltipLabels = value => {
                        return field.valueLabels ? field.valueLabels(value) : `${Object(_helper__WEBPACK_IMPORTED_MODULE_6__["TransformField"])(value || 0, field.thousands, 3)}`;
                    };
                    this.tabs.push({ id, label: field.alias });
                    chartList.push({ id, table: tableInfo.table, field, groupBy, conditions, tooltipLabels });
                });
            });
            this.setState({ chartList });
        }
    }
    // 时间选择器，在初始化时会触发
    onChangeQueryTime(startTime, endTime) {
        this.periodOptions = _core__WEBPACK_IMPORTED_MODULE_4__["CHART_PANEL"].GeneratePeriodOptions(startTime, endTime);
        this.setState({ startTime, endTime });
    }
    onActiveTab(tab) {
        this.setState({ activeId: tab.id });
    }
    onChangeTabStatus() {
        this.setState({ tabStatus: !this.state.tabStatus });
    }
    onScrollTabPanels(e) {
        // 图表滚动加载
        let element = e.target;
        // 向下滚动加载
        if (Math.floor(this.state.chartEndIndex * _constants__WEBPACK_IMPORTED_MODULE_5__["CHART"].DefaultSize.height - element.scrollTop) <= element.clientHeight) {
            if (this.state.chartEndIndex < this.state.chartList.length) {
                this.setState({ chartEndIndex: this.state.chartEndIndex + this.ChartPageSize });
            }
        }
    }
    render() {
        const { activeId, tabStatus } = this.state;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tea-tabs__tabpanel" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_components_Toolbar__WEBPACK_IMPORTED_MODULE_2__["Toolbar"], { style: { marginLeft: 20 }, duration: { from: this.state.startTime, to: this.state.endTime }, onChangeTime: this.onChangeQueryTime.bind(this) }, this.props.children),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Card"], null,
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Card"].Body, null,
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](Tabs, { tabStatus: tabStatus, activeId: activeId, tabs: this.tabs, onActiveTab: this.onActiveTab.bind(this), onChangeTabStatus: this.onChangeTabStatus.bind(this), onScroll: this.onScrollTabPanels.bind(this) }, this.state.chartList.map((chartInfo, index) => {
                        // tab 状态时判断是否为相同 activeId，非tab状态即list显示时判断是否小于chartEndIndex的下标，小于才显示
                        const isActive = tabStatus ? chartInfo.id === activeId : this.state.chartEndIndex >= index;
                        const style = tabStatus
                            ? chartInfo.id === activeId
                                ? this.ShowDivStyle
                                : this.HiddenDivStyle
                            : { width: "100%" };
                        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"](TabPanel, { key: `viewgrid_${index}`, isActive: isActive, tabStatus: tabStatus, style: style },
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_components_FilterTableChart__WEBPACK_IMPORTED_MODULE_3__["FilterTableChart"], { chartId: chartInfo.id, table: chartInfo.table, startTime: this.state.startTime, endTime: this.state.endTime, metric: chartInfo.field, tooltipLabels: chartInfo.tooltipLabels, dimensions: chartInfo.groupBy, conditions: chartInfo.conditions, periodOptions: this.periodOptions })));
                    }))))));
    }
}
class Tabs extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.tabBarStyle = { marginRight: 0, height: 30, lineHeight: 30 };
        this.scrollAreaRef = null;
        this.buttonRef = null;
        this.tabListRef = null;
        this.activeItemRef = null;
        this.scrollAreaRef = react__WEBPACK_IMPORTED_MODULE_0__["createRef"]();
        this.buttonRef = react__WEBPACK_IMPORTED_MODULE_0__["createRef"]();
        this.tabListRef = react__WEBPACK_IMPORTED_MODULE_0__["createRef"]();
        this.activeItemRef = react__WEBPACK_IMPORTED_MODULE_0__["createRef"]();
        this.state = {
            offset: 0,
            scrolling: false
        };
        this.handleScroll = this.handleScroll.bind(this);
    }
    componentDidMount() {
        this.handleScroll();
        window.addEventListener("resize", this.handleScroll);
    }
    componentWillUnmount() {
        window.removeEventListener("resize", this.handleScroll);
    }
    handleScroll() {
        const scrolling = this.getMaxOffset() > 0;
        this.setState({ scrolling });
        // 无需滚动时重置位置
        if (!scrolling) {
            this.setState({ offset: 0 });
        }
        else {
            this.handleActiveItemIntoView();
        }
    }
    handleActiveItemIntoView() {
        requestAnimationFrame(() => {
            if (!this.scrollAreaRef.current || !this.activeItemRef.current) {
                return;
            }
            const scrollAreaRect = this.scrollAreaRef.current.getBoundingClientRect();
            const activeItemRect = this.activeItemRef.current.getBoundingClientRect();
            const startDelta = scrollAreaRect.left - activeItemRect.left + this.buttonRef.current.clientWidth;
            const endDelta = activeItemRect.right - scrollAreaRect.right + this.buttonRef.current.clientWidth;
            if (startDelta > 0) {
                this.setState({ offset: Math.min(0, this.state.offset + startDelta) });
            }
            else if (endDelta > 0) {
                this.setState({ offset: Math.max(0 - this.getMaxOffset(), this.state.offset - endDelta) });
            }
        });
    }
    getStep() {
        return this.scrollAreaRef.current.clientWidth - this.buttonRef.current.clientWidth * 2;
    }
    getMaxOffset() {
        if (!this.scrollAreaRef.current || !this.tabListRef.current) {
            return 0;
        }
        if (this.scrollAreaRef.current.clientWidth >= this.tabListRef.current.clientWidth) {
            return 0;
        }
        return (this.tabListRef.current.clientWidth -
            (this.scrollAreaRef.current.clientWidth - this.buttonRef.current.clientWidth * 2));
    }
    handleBackward() {
        this.setState({ offset: Math.min(0, this.state.offset + this.getStep()) });
    }
    handleForward() {
        this.setState({ offset: Math.max(0 - this.getMaxOffset(), this.state.offset - this.getStep()) });
    }
    render() {
        const { activeId, tabs, tabStatus } = this.props;
        const scrollAreaStyle = tabStatus ? { marginRight: 0 } : this.tabBarStyle;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tea-tabs" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tea-tabs__tabbar" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { ref: this.scrollAreaRef, className: `tea-tabs__scroll-area ${this.state.scrolling ? "is-scrolling" : ""}`, style: Object.assign(Object.assign({}, scrollAreaStyle), { marginRight: 30 }) },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("ul", { ref: this.tabListRef, className: "tea-tabs__tablist", style: {
                            //                transition: "transform 0.2s ease-out 0s",
                            transform: `translate3d(${this.state.offset}px, 0px, 0px)`
                        } }, tabStatus &&
                        tabs.map(tab => {
                            return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { key: tab.id, ref: tab.id === activeId ? this.activeItemRef : undefined, className: "tea-tabs__tabitem", onClick: () => this.props.onActiveTab(tab) },
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: `tea-tabs__tab ${activeId === tab.id ? "is-active" : ""}` }, tab.label)));
                        })),
                    tabStatus && (react__WEBPACK_IMPORTED_MODULE_0__["createElement"](react__WEBPACK_IMPORTED_MODULE_0__["Fragment"], null,
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Button"], { ref: this.buttonRef, className: "tea-tabs__backward", type: "icon", icon: "arrowleft", disabled: this.state.offset >= 0, onClick: this.handleBackward.bind(this) }),
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Button"], { className: "tea-tabs__forward", type: "icon", icon: "arrowright", disabled: this.state.offset <= 0 - this.getMaxOffset(), onClick: this.handleForward.bind(this) })))),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tea-tabs__addons" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", null,
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Icon"], { type: tabStatus ? "viewgrid" : "viewlist", style: { cursor: "pointer" }, onClick: this.props.onChangeTabStatus })))),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: tabStatus ? "tab-status-panel" : "no-tab-status-panel", style: tabStatus
                    ? {
                        marginLeft: `${-100 * this.props.tabs.findIndex(chartInfo => chartInfo.id === this.props.activeId)}%`
                    }
                    : {}, onScroll: this.props.onScroll }, this.props.children)));
    }
}
/**
 * hasActive 激活后才会渲染子项组件
 */
class TabPanel extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.state = {
            hasActive: props.isActive
        };
    }
    componentWillReceiveProps(nextProps) {
        if (nextProps.isActive) {
            this.setState({ hasActive: true });
        }
    }
    render() {
        const { style = {} } = this.props;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: Object.assign({}, style), className: "tea-tabs__tabpanel" }, this.props.tabStatus ? (this.state.hasActive && this.props.children) : this.state.hasActive ? (this.props.children) : (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: { width: "100%", height: _constants__WEBPACK_IMPORTED_MODULE_5__["CHART"].DefaultSize.height } }))));
    }
}


/***/ }),

/***/ "./src/panel/containers/Instances.less":
/*!*********************************************!*\
  !*** ./src/panel/containers/Instances.less ***!
  \*********************************************/
/*! no static exports found */
/***/ (function(module, exports, __webpack_require__) {


var content = __webpack_require__(/*! !../../../node_modules/css-loader/dist/cjs.js!../../../node_modules/less-loader/dist/cjs.js!./Instances.less */ "./node_modules/css-loader/dist/cjs.js!./node_modules/less-loader/dist/cjs.js!./src/panel/containers/Instances.less");

if(typeof content === 'string') content = [[module.i, content, '']];

var transform;
var insertInto;



var options = {"hmr":true}

options.transform = transform
options.insertInto = undefined;

var update = __webpack_require__(/*! ../../../node_modules/style-loader/lib/addStyles.js */ "./node_modules/style-loader/lib/addStyles.js")(content, options);

if(content.locals) module.exports = content.locals;

if(false) {}

/***/ }),

/***/ "./src/panel/containers/Instances.tsx":
/*!********************************************!*\
  !*** ./src/panel/containers/Instances.tsx ***!
  \********************************************/
/*! exports provided: ChartInstancesPanel */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChartInstancesPanel", function() { return ChartInstancesPanel; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _components_Toolbar__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../components/Toolbar */ "./src/panel/components/Toolbar.tsx");
/* harmony import */ var _components_MetricCharts__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../components/MetricCharts */ "./src/panel/components/MetricCharts.tsx");
/* harmony import */ var _components_InstanceList__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../components/InstanceList */ "./src/panel/components/InstanceList.tsx");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../constants */ "./src/panel/constants.ts");
/* harmony import */ var _core__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! ../core */ "./src/panel/core.ts");
/* harmony import */ var _helper__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! ../helper */ "./src/panel/helper.ts");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_7__ = __webpack_require__(/*! @tencent/tea-component */ "@tencent/tea-component");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_7___default = /*#__PURE__*/__webpack_require__.n(_tencent_tea_component__WEBPACK_IMPORTED_MODULE_7__);








__webpack_require__(/*! ./Instances.less */ "./src/panel/containers/Instances.less");
class ChartInstancesPanel extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this._instances = {}; // group by 查询返回数据中的不同类别的groupBy的值
        /**
         * 勾选instances。更新到state中
         */
        this.onCheckInstances = (instanceList) => {
            this.setState({ instanceList }, () => {
                const seriesGroupTemp = JSON.parse(JSON.stringify(this.state.seriesGroup));
                const seriesGroup = seriesGroupTemp.map((series) => {
                    // series 是每一个图表的数据
                    series.lines.forEach((line) => {
                        // line 是每个图表中每条线的数据
                        line.disable = this.checkLineIsDisable(this.state.instanceList, line.legend);
                    });
                    return series;
                });
                // 对象深拷贝有问题
                this.setState({ seriesGroup });
            });
        };
        const endTime = new Date();
        const startTime = new Date(endTime.getTime() - 1000 * 60 * 60);
        const periodOptions = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].GeneratePeriodOptions(startTime, endTime);
        this.state = {
            loading: false,
            startTime: _constants__WEBPACK_IMPORTED_MODULE_4__["TIME_PICKER"].Tabs[0].from,
            endTime: _constants__WEBPACK_IMPORTED_MODULE_4__["TIME_PICKER"].Tabs[0].to,
            seriesGroup: [],
            instanceList: props.instance.list,
            aggregation: Object.keys(_constants__WEBPACK_IMPORTED_MODULE_4__["QUERY"].Aggregation)[0],
            periodOptions: periodOptions,
            period: periodOptions[0].value,
        };
    }
    componentDidUpdate(prevProps, prevState, snapshot) {
        let hasUpdate = false;
        let hasFetch = false;
        const { tables, groupBy, instance } = this.props;
        let instanceList = [];
        if (instance) {
            instanceList = this.props.instance.list;
            if (JSON.stringify(instance) !== JSON.stringify(prevProps.instance)) {
                hasUpdate = true;
            }
        }
        if (JSON.stringify(groupBy) !== JSON.stringify(prevProps.groupBy) ||
            JSON.stringify(tables) !== JSON.stringify(prevProps.tables)) {
            hasFetch = true;
        }
        if (hasFetch) {
            this.fetchDataByTables({ instances: instanceList });
        }
        else if (hasUpdate) {
            instanceList = this.combineInstancesWithGroupBy(instance.list);
            this.setState({ instanceList });
        }
        if (JSON.stringify(prevState) !== JSON.stringify(this.state)) {
            hasUpdate = true;
        }
        return hasUpdate || hasFetch;
    }
    /**
     * 根据key值判断chart折线是否隐藏
     * @param instanceList
     * @param key
     * @returns {boolean}
     */
    checkLineIsDisable(instanceList = [], key = "") {
        if (instanceList.length === 0) {
            return false;
        }
        const instance = instanceList.find((instance) => {
            const instanceKey = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].GenerateRowKey(this.props.groupBy.map((groupBy) => instance[groupBy.value]));
            return instanceKey.indexOf(key) !== -1;
        });
        if (!instance || !instance.isChecked) {
            return true;
        }
        return false;
    }
    combineInstancesWithGroupBy(instances) {
        // this.state.instanceList 为用户通过接口传入的值 例如结构 { workload_kind: "Deployment", isChecked: true },
        let instanceList = instances || [].concat(this.state.instanceList);
        /**
         * 根据 groupByData 字典中key的值，查找instances中是否匹配
         * 对instances和groupBy做交集
         */
        Object.keys(this._instances).forEach((instanceKey) => {
            let instance = instanceList.find((instance) => {
                // 通过groupBy的key在 this.state.instanceList 的值中获取值列表,
                // 如：this.props.groupBy=[workload_kind], instance={workload_kind: "Deployment", isChecked: true}
                // groupByValues 为 [Deployment]
                const groupByValues = this.props.groupBy.map((groupByItem) => instance[groupByItem.value]);
                const rowKey = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].GenerateRowKey(groupByValues);
                return instanceKey === rowKey;
            });
            if (!instance) {
                // 用户配置的instances值通过groupBy查询找不到时，则新增instance
                instanceList.push(this._instances[instanceKey]);
            }
        });
        return instanceList;
    }
    /**
     * 获取空对象
     */
    getEmptyData(fields = []) {
        let seriesGroup = [];
        const tables = fields.length > 0 ? [{ fields }] : this.props.tables;
        tables.forEach((tableInfo) => {
            tableInfo.fields.forEach((field) => {
                seriesGroup.push({
                    field,
                    title: field.alias,
                    labels: [],
                    lines: [],
                });
            });
        });
        return seriesGroup;
    }
    onChangeQueryTime(startTime, endTime) {
        if (this.state.loading) {
            return;
        }
        const periodOptions = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].GeneratePeriodOptions(startTime, endTime);
        const period = periodOptions[0].value;
        this.setState({ startTime, endTime, periodOptions, period }, () => {
            // 选择时间段后重新请求数据
            this.fetchDataByTables({ instances: this.state.instanceList });
        });
    }
    onChangePeriod(period) {
        this.setState({ period }, () => {
            // 选择时间粒度后重新请求数据
            this.fetchDataByTables({ instances: this.state.instanceList });
        });
    }
    fetchDataByTables(options = {}) {
        this._instances = {}; // 请求时清除原来数据加载的 _instances 对象
        // 初始状态
        this.setState({
            seriesGroup: this.getEmptyData(),
            instanceList: this.props.instance.list,
            loading: true,
        });
        const { tables } = this.props;
        const requestPromises = tables.map((tableInfo) => {
            return this.requestData(tableInfo);
        });
        Promise.all(requestPromises).then((values) => {
            const seriesGroup = [].concat(...values);
            // combineInstancesWithGroupBy 合并 this._instances 与 instances
            let instanceList = this.combineInstancesWithGroupBy(this.props.instance.list);
            this.setState({ seriesGroup, instanceList, loading: false });
        });
    }
    /**
     * 发起查询数据请求
     * @param {string} table
     * @param {Array<any>} fields
     * @param {Array<any>} groupBy
     * @param {Array<any>} conditions
     * @param {string} period
     */
    requestData(tableInfo = {}) {
        const { fields, conditions, table } = tableInfo;
        const { startTime, endTime, period } = this.state;
        const groupByItems = this.props.groupBy.map((item) => item.value) || [];
        // 查询的时间粒度
        // const period = Period(startTime, endTime) as any;
        return _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].RequestData({
            table: table,
            startTime,
            endTime,
            fields: fields,
            dimensions: groupByItems,
            conditions: conditions,
            period: period,
        })
            .then((res) => {
            const { columns, data } = res;
            /**
             * 根据 groupBy 条件对数据做聚合，已 groupBy 为key进行存储
             */
            const dimensionsData = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].AggregateDataByDimensions(groupByItems, columns, data);
            /**
             * 根据 groupBy 数据生成默认 instances
             */
            const columnsInfo = groupByItems.map((value) => {
                return {
                    index: columns.indexOf(value),
                    value,
                };
            });
            Object.keys(dimensionsData).forEach((groupByKey) => {
                //从第一行的数据获取group by的值
                const row = dimensionsData[groupByKey][0];
                // 缓存数据中groupBy联合主键的instance list
                if (!this._instances[groupByKey]) {
                    let instance = {}; // 初始化 instance 对象
                    columnsInfo.forEach((column) => {
                        const key = column.value;
                        // key 为groupBy的参数，row[column.index]对应groupBy的值 如{workload_kind: "Deployment"}
                        instance[key] = row[column.index];
                    });
                    // 如果传入的instance list 为空数组是，数据中的instance对象为勾选状态
                    instance["isChecked"] =
                        this.props.instance && this.props.instance.list.length > 0
                            ? false
                            : true;
                    this._instances[groupByKey] = instance;
                }
            });
            /**
             * 生成 labels
             */
            const timestampIndex = columns.indexOf(`timestamp(${period}s)`);
            // 默认表格第一列为时间序列, 返回的时间列表可能不是startTime开始，以endTime结束，需要补帧
            let labels = (Object.values(dimensionsData)[0] || []).map((item) => item[timestampIndex]);
            labels = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].OffsetTimeSeries(startTime, endTime, parseInt(period), labels);
            /**
             * 生成 series
             */
            let seriesGroup = [];
            // 获取需要显示的field的数据，每个field是一个图表。
            fields.forEach((field) => {
                const fieldIndex = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].FindMetricIndex(columns, field.expr); // 数据列下标
                const valueTransform = field.valueTransform || ((value) => value);
                // 每个图表中显示groupBy类别一个field -> chart
                let lines = [];
                Object.keys(dimensionsData).forEach((groupByKey) => {
                    // 每个groupBy的数据集， fieldsIndex的index获取groupBy该下标的值
                    const sourceValue = dimensionsData[groupByKey];
                    const data = {};
                    sourceValue.forEach((row) => {
                        let value = row[fieldIndex];
                        // data 记录 时间戳row[timestampIndex] 对应 value
                        data[row[timestampIndex]] = valueTransform(value);
                    });
                    let line = _core__WEBPACK_IMPORTED_MODULE_5__["CHART_PANEL"].GenerateLine(groupByKey, data, this.checkLineIsDisable(this.props.instance.list, groupByKey));
                    lines.push(line);
                });
                // 每个chart 的数据存储到 seriesGroup
                seriesGroup.push({
                    labels,
                    field,
                    lines,
                    title: field.alias,
                    tooltipLabels: (value, from) => {
                        return field.valueLabels
                            ? field.valueLabels(value, from)
                            : `${Object(_helper__WEBPACK_IMPORTED_MODULE_6__["TransformField"])(value || 0, field.thousands, 3)}`;
                    },
                });
            });
            return seriesGroup;
        })
            .catch((e) => {
            return this.getEmptyData(fields);
        });
    }
    render() {
        const { loading, seriesGroup, instanceList, startTime, endTime, periodOptions, period, } = this.state;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-rich-dialog-bd" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_components_Toolbar__WEBPACK_IMPORTED_MODULE_1__["Toolbar"], { duration: { from: this.state.startTime, to: this.state.endTime }, onChangeTime: this.onChangeQueryTime.bind(this) },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { style: {
                        fontSize: 12,
                        display: "inline-block",
                        verticalAlign: "middle",
                    } }, "\u7EDF\u8BA1\u7C92\u5EA6\uFF1A"),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_7__["Select"], { type: "native", size: "s", appearence: "button", options: periodOptions, value: period, onChange: (value) => {
                        this.onChangePeriod(value);
                    }, placeholder: "\u8BF7\u9009\u62E9" }),
                this.props.children),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "monitor-dialog-data-box" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-g" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-g-u-3-4" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_components_MetricCharts__WEBPACK_IMPORTED_MODULE_2__["default"], { className: "monitor-chart-grid", style: {
                                maxHeight: "calc(100vh - 200px)",
                                height: "668px",
                                overflowY: "auto",
                                overflowX: "hidden",
                                display: "block",
                                border: "1px solid #ddd",
                            }, loading: loading, min: typeof startTime === "string" ? 0 : startTime.getTime(), max: typeof endTime === "string" ? 0 : endTime.getTime(), seriesGroup: seriesGroup, reload: this.fetchDataByTables.bind(this, {
                                instances: this.state.instanceList,
                            }) })),
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-g-u-1-4" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-mod-selector-tb" },
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-option-cell options-left" },
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-option-bd" },
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_components_InstanceList__WEBPACK_IMPORTED_MODULE_3__["default"], { className: "tc-15-option-box tc-scroll", style: {
                                            maxHeight: "calc(100vh - 200px)",
                                            height: "668px",
                                        }, columns: this.props.instance.columns, list: instanceList, update: this.onCheckInstances })))))))));
    }
}


/***/ }),

/***/ "./src/panel/containers/Pure.tsx":
/*!***************************************!*\
  !*** ./src/panel/containers/Pure.tsx ***!
  \***************************************/
/*! exports provided: ChartPanel */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChartPanel", function() { return ChartPanel; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @tencent/tea-component */ "@tencent/tea-component");
/* harmony import */ var _tencent_tea_component__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var _components_Toolbar__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../components/Toolbar */ "./src/panel/components/Toolbar.tsx");
/* harmony import */ var _components_PureChart__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../components/PureChart */ "./src/panel/components/PureChart.tsx");
/* harmony import */ var _helper__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../helper */ "./src/panel/helper.ts");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! ../constants */ "./src/panel/constants.ts");
/* harmony import */ var _core__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! ../core */ "./src/panel/core.ts");







/**
 * 基础 Panel，提供数据请求功能
 */
class ChartPanel extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this._id = "ChartPanel";
        const endTime = new Date();
        const startTime = new Date(endTime.getTime() - 1000 * 60 * 60);
        const periodOptions = _core__WEBPACK_IMPORTED_MODULE_6__["CHART_PANEL"].GeneratePeriodOptions(startTime, endTime);
        this.state = {
            loading: false,
            startTime: _constants__WEBPACK_IMPORTED_MODULE_5__["TIME_PICKER"].Tabs[0].from,
            endTime: _constants__WEBPACK_IMPORTED_MODULE_5__["TIME_PICKER"].Tabs[0].to,
            seriesGroup: [],
            aggregation: Object.keys(_constants__WEBPACK_IMPORTED_MODULE_5__["QUERY"].Aggregation)[0],
            periodOptions: periodOptions,
            period: periodOptions[0].value
        };
    }
    componentDidUpdate(prevProps, prevState, snapshot) {
        let hasUpdate = false;
        let hasFetch = false;
        const { tables, groupBy } = this.props;
        if (JSON.stringify(groupBy) !== JSON.stringify(prevProps.groupBy) ||
            JSON.stringify(tables) !== JSON.stringify(prevProps.tables)) {
            hasFetch = true;
        }
        if (hasFetch) {
            this.fetchDataByTables();
        }
        if (JSON.stringify(prevState) !== JSON.stringify(this.state)) {
            hasUpdate = true;
        }
        return hasUpdate || hasFetch;
    }
    onChangeQueryTime(startTime, endTime) {
        if (this.state.loading) {
            return;
        }
        const periodOptions = _core__WEBPACK_IMPORTED_MODULE_6__["CHART_PANEL"].GeneratePeriodOptions(startTime, endTime);
        const period = periodOptions[0].value;
        this.setState({ startTime, endTime, periodOptions, period }, () => {
            // 选择时间段后重新请求数据
            this.fetchDataByTables();
        });
    }
    onChangePeriod(period) {
        this.setState({ period }, () => {
            // 选择时间粒度后重新请求数据
            this.fetchDataByTables();
        });
    }
    getEmptyData(fields = []) {
        let seriesGroup = [];
        const tables = fields.length > 0 ? [{ fields }] : this.props.tables;
        tables.forEach(tableInfo => {
            tableInfo.fields.forEach(field => {
                seriesGroup.push({
                    field,
                    title: field.alias,
                    labels: [],
                    lines: []
                });
            });
        });
        return seriesGroup;
    }
    fetchDataByTables() {
        // 初始状态
        this.setState({ loading: true, seriesGroup: this.getEmptyData() });
        const { tables } = this.props;
        const requestPromises = tables.map(tableInfo => {
            return this.requestData(tableInfo);
        });
        Promise.all(requestPromises).then(values => {
            const seriesGroup = [].concat(...values);
            this.setState({ seriesGroup, loading: false });
        });
    }
    /**
     * 发起查询数据请求
     * @param {string} table
     * @param {Array<any>} fields
     * @param {Array<any>} groupBy
     * @param {Array<any>} conditions
     * @param {string} period
     */
    requestData(options = {}) {
        const { fields, conditions, table } = options;
        const groupBy = options.groupBy || this.props.groupBy;
        const { startTime, endTime, period } = this.state;
        const groupByItems = groupBy.map(item => item.value) || [];
        return _core__WEBPACK_IMPORTED_MODULE_6__["CHART_PANEL"].RequestData({
            table: table,
            startTime,
            endTime,
            fields,
            dimensions: groupByItems,
            conditions: conditions,
            period: period
        })
            .then(res => {
            const { columns, data } = res;
            /**
             * 根据 groupBy 条件对数据做聚合，已 groupBy 为key进行存储
             */
            let dimensionsData = _core__WEBPACK_IMPORTED_MODULE_6__["CHART_PANEL"].AggregateDataByDimensions(groupByItems, columns, data);
            /**
             * 生成 labels
             */
            const timestampIndex = columns.indexOf(`timestamp(${period}s)`);
            // 默认表格第一列为时间序列, 返回的时间列表可能不是startTime开始，以endTime结束，需要补帧
            let labels = (Object.values(dimensionsData)[0] || []).map(item => item[timestampIndex]);
            labels = _core__WEBPACK_IMPORTED_MODULE_6__["CHART_PANEL"].OffsetTimeSeries(startTime, endTime, parseInt(period), labels);
            /**
             * 生成 series
             */
            let seriesGroup = [];
            // 获取需要显示的field的数据，每个field是一个图表。
            fields.forEach(field => {
                const fieldIndex = _core__WEBPACK_IMPORTED_MODULE_6__["CHART_PANEL"].FindMetricIndex(columns, field.expr); // 数据列下标
                const valueTransform = field.valueTransform || (value => value);
                // 每个图表中显示groupBy类别一个field -> chart
                let lines = [];
                Object.keys(dimensionsData).forEach(groupByKey => {
                    // 每个groupBy的数据集， fieldsIndex的index获取groupBy该下标的值
                    const sourceValue = dimensionsData[groupByKey];
                    const data = {};
                    sourceValue.forEach(row => {
                        let value = row[fieldIndex];
                        // data 记录 时间戳row[timestampIndex] 对应 value
                        data[row[timestampIndex]] = valueTransform(value);
                    });
                    let line = _core__WEBPACK_IMPORTED_MODULE_6__["CHART_PANEL"].GenerateLine(groupByKey, data);
                    lines.push(line);
                });
                const chart = {
                    title: field.alias,
                    labels: labels,
                    field: field,
                    tooltipLabels: (value, from = "tooltip") => {
                        return field.valueLabels
                            ? field.valueLabels(value, from)
                            : `${Object(_helper__WEBPACK_IMPORTED_MODULE_4__["TransformField"])(value || 0, field.thousands, 3)}`;
                    },
                    lines: lines
                };
                // 每个chart 的数据存储到 seriesGroup
                seriesGroup.push(chart);
            });
            return seriesGroup;
        })
            .catch(e => {
            return this.getEmptyData(fields);
        });
    }
    render() {
        const { loading, seriesGroup, periodOptions, period, startTime, endTime } = this.state;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tea-tabs__tabpanel" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_components_Toolbar__WEBPACK_IMPORTED_MODULE_2__["Toolbar"], { style: { marginLeft: 20 }, duration: { from: this.state.startTime, to: this.state.endTime }, onChangeTime: this.onChangeQueryTime.bind(this) },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { style: { fontSize: 12, display: "inline-block", verticalAlign: "middle" } }, "\u7EDF\u8BA1\u7C92\u5EA6\uFF1A"),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_tencent_tea_component__WEBPACK_IMPORTED_MODULE_1__["Select"], { type: "native", size: "s", appearence: "button", options: periodOptions, value: period, onChange: value => {
                        this.onChangePeriod(value);
                    }, placeholder: "\u8BF7\u9009\u62E9" })),
            seriesGroup.map((chart, index) => {
                const key = `${this._id}_${index}`;
                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { key: key, style: { marginBottom: 10, padding: "0px 10px" } },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_components_PureChart__WEBPACK_IMPORTED_MODULE_3__["default"], { id: key, loading: loading, min: typeof startTime === "string" ? 0 : startTime.getTime(), max: typeof endTime === "string" ? 0 : endTime.getTime(), width: this.props.width, heigth: this.props.height, title: chart.title, field: chart.field, labels: chart.labels, unit: chart.field.unit, tooltipLabels: chart.tooltipLabels, lines: chart.lines, reload: this.fetchDataByTables.bind(this) })));
            })));
    }
}


/***/ }),

/***/ "./src/panel/core.ts":
/*!***************************!*\
  !*** ./src/panel/core.ts ***!
  \***************************/
/*! exports provided: CHART_PANEL */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "CHART_PANEL", function() { return CHART_PANEL; });
/* harmony import */ var moment__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! moment */ "moment");
/* harmony import */ var moment__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(moment__WEBPACK_IMPORTED_MODULE_0__);
/* harmony import */ var _helper__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./helper */ "./src/panel/helper.ts");
/* harmony import */ var _constants__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./constants */ "./src/panel/constants.ts");
/* harmony import */ var _tce_request__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../tce/request */ "./src/tce/request.ts");




var CHART_PANEL;
(function (CHART_PANEL) {
    /**
     * 用于groupBY中的值产生图表中折线的legend或instances的key
     * 使图表的折线的legend 跟 instance 对应，用于instances勾选功能
     */
    function GenerateRowKey(row = []) {
        return row.length > 0 ? row.join(" | ") : "";
    }
    CHART_PANEL.GenerateRowKey = GenerateRowKey;
    /**
     * 生成图表中显示折线的数据对象
     * @param legend
     * @param data
     * @returns {ModelType}
     * @constructor
     */
    function GenerateLine(legend, data, disable = false) {
        return {
            legend: legend,
            disable: disable,
            data: data
        };
    }
    CHART_PANEL.GenerateLine = GenerateLine;
    // 生成时间粒度
    function GeneratePeriodOptions(startTime, endTime) {
        let periods = [60, 300, 3600, 86400];
        function periodUnit(period) {
            const unit = ["秒", "分钟", "小时"];
            let temp = period / 60;
            let unitRound = 1;
            while (unitRound < 3 && temp >= 60) {
                temp = temp / 60;
                unitRound += 1;
            }
            if (temp < 1) {
                return `${period}${unit[0]}`;
            }
            return `${temp}${unit[unitRound]}`;
        }
        const period = Object(_helper__WEBPACK_IMPORTED_MODULE_1__["Period"])(startTime, endTime);
        const l = moment__WEBPACK_IMPORTED_MODULE_0___default()(endTime).diff(moment__WEBPACK_IMPORTED_MODULE_0___default()(startTime));
        return periods.filter(d => d >= period && d < l).map(item => {
            return {
                value: item,
                text: periodUnit(item)
            };
        });
    }
    CHART_PANEL.GeneratePeriodOptions = GeneratePeriodOptions;
    /**
     * 发起查询数据请求
     * @param {string} table
     * @param {Date} startTime
     * @param {Date} endTime
     * @param {Array<any>} fields
     * @param {Array<string>} dimensions
     * @param {Array<Array<any>>} conditions
     * @param {string} period
     */
    async function RequestData(options) {
        const params = {
            table: options.table,
            startTime: options.startTime.getTime(),
            endTime: options.endTime.getTime(),
            fields: [...options.fields.map(item => `${item.expr}`)],
            conditions: options.conditions,
            orderBy: "timestamp",
            groupBy: [`timestamp(${options.period}s)`, ...options.dimensions],
            order: "asc",
            limit: _constants__WEBPACK_IMPORTED_MODULE_2__["QUERY"].Limit
        };
        try {
            const res = await Object(_tce_request__WEBPACK_IMPORTED_MODULE_3__["request"])({
                data: {
                    Path: "/front/v1/get/query",
                    RequestBody: params
                }
            });
            // 没有数据，抛出异常
            if (!res.hasOwnProperty("columns") || !res.hasOwnProperty("data")) {
                throw new Error();
            }
            // 更新图表数据
            const { columns, data } = res;
            if (!Array.isArray(columns) || !Array.isArray(data)) {
                throw new Error();
            }
            if (data.length === _constants__WEBPACK_IMPORTED_MODULE_2__["QUERY"].Limit) {
                /**
                 * 返回数据长度与请求限制一致，则说明对象显示数据点不够，需要分段发起请求
                 */
                // 获取要显示的对象
                const instances = AggregateDataByDimensions(options.dimensions, columns, data);
                // 每个对象一天显示点数
                const instancePoints = 86400 / parseInt(options.period);
                // 显示的对象一天显示点数
                const sumInstancesPoints = Object.keys(instances).length * instancePoints;
                const daysNum = Math.ceil((params.endTime - params.startTime) / _constants__WEBPACK_IMPORTED_MODULE_2__["OneDayMillisecond"]);
                const requestNum = Math.ceil(sumInstancesPoints * daysNum / _constants__WEBPACK_IMPORTED_MODULE_2__["QUERY"].Limit);
                const timeInterval = Math.floor((params.endTime - params.startTime) / requestNum);
                // 生成请求数组
                let requestPromises = [];
                for (let i = 1; i <= requestNum; i++) {
                    const paramsTemp = Object.assign(Object.assign({}, params), { startTime: params.startTime + timeInterval * (i - 1), endTime: params.startTime + timeInterval * i });
                    requestPromises.push(Object(_tce_request__WEBPACK_IMPORTED_MODULE_3__["request"])({
                        data: {
                            Path: "/front/v1/get/query",
                            RequestBody: paramsTemp
                        }
                    }));
                }
                const resAll = await Promise.all(requestPromises);
                // 将data转化为一维数组
                return { columns, data: [].concat(...resAll.map(item => item.data)) };
            }
            return res;
        }
        catch (error) {
            throw error;
        }
    }
    CHART_PANEL.RequestData = RequestData;
    /**
     * 根据维度对请求返回数据进行聚合
     * @param dimensions
     * @param columns
     * @param data
     * @returns {{}}
     */
    function AggregateDataByDimensions(dimensions, columns, data) {
        columns = columns || [];
        data = data || [];
        /**
         * 根据 groupBy 条件对数据做聚合，已 groupBy 为key进行存储
         */
        let dimensionIndex = []; // 记录是 groupBy 数据的下标
        let dimensionData = {}; // 记录 groupBy 联合 key 集合数据
        dimensions.forEach(item => {
            const index = columns.indexOf(item);
            if (index >= 0) {
                dimensionIndex.push(index);
            }
        });
        data.forEach(row => {
            // 使用每一行数据中groupBy字段对应值的字符串合并作为该行的唯一key
            const rowKey = CHART_PANEL.GenerateRowKey(dimensionIndex.map(index => row[index]));
            if (dimensionData[rowKey]) {
                dimensionData[rowKey].push(row);
            }
            else {
                dimensionData[rowKey] = [row];
            }
        });
        return dimensionData;
    }
    CHART_PANEL.AggregateDataByDimensions = AggregateDataByDimensions;
    /**
     *
     * @param {Date} startTime
     * @param {Date} endTime
     */
    function OffsetTimeSeries(startTime, endTime, period, labels) {
        labels = [labels[0] || startTime.getTime()];
        const periodNum = Math.floor((endTime.getTime() - startTime.getTime()) / period / 1000);
        if (periodNum - labels.length > 0) {
            const periodMillisecond = period * 1000;
            const startTimestamp = startTime.getTime();
            const endTimestamp = endTime.getTime();
            let cycleIndex = 0;
            while (true) {
                if (cycleIndex > 15000) {
                    // 安全机制，防止死循环，相当于15000分钟的跨度
                    break;
                }
                cycleIndex += 1;
                const start = labels[0];
                const end = labels[labels.length - 1];
                if (start - startTimestamp >= periodMillisecond) {
                    labels.unshift(start - periodMillisecond);
                }
                if (endTimestamp - end >= periodMillisecond) {
                    labels.push(end + periodMillisecond);
                }
                if (periodNum < labels.length
                    || (labels[0] - startTimestamp < periodMillisecond && endTimestamp - labels[labels.length - 1] < periodMillisecond)) {
                    break;
                }
            }
        }
        return labels;
    }
    CHART_PANEL.OffsetTimeSeries = OffsetTimeSeries;
    /**
     * 根据指标表达式查找在数据中的下标值
     * @param {Array<string>} columns
     * @param {string} metricExpr
     */
    function FindMetricIndex(columns, metricExpr) {
        const match = /(\w+)?\((\w+)\)/.exec(metricExpr);
        let fieldIndex = -1; // field 对应的在 data 中的下标
        if (match && match.length === 3) {
            fieldIndex = columns.indexOf(`${match[2]}_${match[1]}`);
            if (fieldIndex === -1) {
                // 兼容influxdb没拼接聚合方式在 k8s_workload_pod_restart_total
                fieldIndex = columns.indexOf(match[2]);
            }
        }
        else {
            fieldIndex = columns.findIndex(item => item === metricExpr);
        }
        return fieldIndex;
    }
    CHART_PANEL.FindMetricIndex = FindMetricIndex;
})(CHART_PANEL || (CHART_PANEL = {}));


/***/ }),

/***/ "./src/tea-components/bubble/Bubble.tsx":
/*!**********************************************!*\
  !*** ./src/tea-components/bubble/Bubble.tsx ***!
  \**********************************************/
/*! exports provided: Bubble */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Bubble", function() { return Bubble; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! classnames */ "./node_modules/classnames/index.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(classnames__WEBPACK_IMPORTED_MODULE_1__);


function Bubble({ focusPosition = "top", focusOffset, children, className, style }) {
    const offsetProperty = focusPosition === "top" || focusPosition === "bottom" ? "top" : "left";
    style = Object.assign({}, { [offsetProperty]: focusOffset }, style);
    return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: classnames__WEBPACK_IMPORTED_MODULE_1___default()(`tc-15-bubble tc-15-bubble-${focusPosition}`, className), style: style },
        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-bubble-inner" }, children)));
}


/***/ }),

/***/ "./src/tea-components/bubble/BubbleWrapper.tsx":
/*!*****************************************************!*\
  !*** ./src/tea-components/bubble/BubbleWrapper.tsx ***!
  \*****************************************************/
/*! exports provided: BubbleWrapper */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "BubbleWrapper", function() { return BubbleWrapper; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! classnames */ "./node_modules/classnames/index.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(classnames__WEBPACK_IMPORTED_MODULE_1__);


function BubbleWrapper({ children, className, align }) {
    const finalClassName = classnames__WEBPACK_IMPORTED_MODULE_1___default()("tc-15-bubble-icon", {
        ["tc-15-triangle-align-" + align]: !!align
    }, className);
    return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: finalClassName }, children));
}


/***/ }),

/***/ "./src/tea-components/bubble/index.ts":
/*!********************************************!*\
  !*** ./src/tea-components/bubble/index.ts ***!
  \********************************************/
/*! exports provided: Bubble, BubbleWrapper */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _Bubble__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./Bubble */ "./src/tea-components/bubble/Bubble.tsx");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "Bubble", function() { return _Bubble__WEBPACK_IMPORTED_MODULE_0__["Bubble"]; });

/* harmony import */ var _BubbleWrapper__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./BubbleWrapper */ "./src/tea-components/bubble/BubbleWrapper.tsx");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "BubbleWrapper", function() { return _BubbleWrapper__WEBPACK_IMPORTED_MODULE_1__["BubbleWrapper"]; });





/***/ }),

/***/ "./src/tea-components/datetimepicker/DateTimePicker.tsx":
/*!**************************************************************!*\
  !*** ./src/tea-components/datetimepicker/DateTimePicker.tsx ***!
  \**************************************************************/
/*! exports provided: DateTimePicker */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "DateTimePicker", function() { return DateTimePicker; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! classnames */ "./node_modules/classnames/index.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(classnames__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! tea-components/libs/decorators/OnOuterClick */ "./src/tea-components/libs/decorators/OnOuterClick.ts");
/* harmony import */ var _SingleDatePicker__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./SingleDatePicker */ "./src/tea-components/datetimepicker/SingleDatePicker.tsx");
/* harmony import */ var _timepicker__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../timepicker */ "./src/tea-components/timepicker/index.ts");
/* harmony import */ var _i18n__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! ../../i18n */ "./src/i18n/index.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};






const version = window.VERSION || "zh";
const language = _i18n__WEBPACK_IMPORTED_MODULE_5__[version];
/**
 * props中时间单位为s，乘unit转化为ms
 */
const unit = 1000;
/**
 * 日期时间范围选择组件
 * TODO 受控组件支持
 */
class DateTimePicker extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        /**
         * 根据 props 获得初始日期和范围
         */
        this.getInitDateAndRange = (props) => {
            let dateFrom, timeFrom, dateTo, timeTo, pickerValue = null;
            let defaultValue = props.defaultValue || {};
            // defaultValue
            if (defaultValue.from) {
                const from = this.parseMacro(defaultValue.from);
                dateFrom = this.formatDate(this.getDate(from));
                timeFrom = this.formatTime(this.getTime(from));
                pickerValue = { from, to: from };
            }
            if (defaultValue.to) {
                const to = this.parseMacro(defaultValue.to);
                dateTo = this.formatDate(this.getDate(to));
                timeTo = this.formatTime(this.getTime(to));
                if (!pickerValue) {
                    pickerValue = { from: to, to };
                }
                else {
                    pickerValue.to = to;
                }
            }
            // 获取初始范围
            const { rangeFrom, rangeTo } = this.getInitRange(props);
            const range = props.range || {};
            // duration 存在时更新起止时间
            if ('duration' in props && pickerValue) {
                pickerValue.to = new Date(pickerValue.from.getTime() + props.duration * unit);
                // 超出范围置空
                if (rangeFrom.max && this.compare(pickerValue.from, rangeFrom.max) > 0) {
                    dateFrom = timeFrom = dateTo = timeTo = pickerValue = null;
                    defaultValue = {};
                }
                else {
                    dateTo = this.formatDate(this.getDate(pickerValue.to));
                    timeTo = this.formatTime(this.getTime(pickerValue.to));
                }
            }
            // TODO
            // DefaultValue Range 判断
            // 根据 defaultValue 限定开始结束范围
            if (defaultValue.from) {
                const date = this.parse(`${dateFrom} ${timeFrom}`);
                if ((rangeFrom.min && this.compare(date, rangeFrom.min) < 0) || (rangeFrom.max && this.compare(date, rangeFrom.max) > 0))
                    return;
                rangeTo.min = date;
                if ('maxLength' in range) {
                    rangeTo.max = new Date(date.getTime() + range.maxLength * unit);
                }
            }
            if (defaultValue.to) {
                const date = this.parse(`${dateTo} ${timeTo}`);
                if ((rangeTo.min && this.compare(date, rangeTo.min) < 0) || (rangeTo.max && this.compare(date, rangeTo.max) > 0))
                    return;
                rangeFrom.max = date;
                // if ('maxLength' in range) {
                //     rangeFrom.min = new Date(date.getTime() - range.maxLength * unit);
                // }
            }
            return { dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue };
        };
        /**
         * 根据 props 获得初始范围
         */
        this.getInitRange = (props) => {
            const { duration, range = {} } = props;
            let rangeFrom = {}, rangeTo = {};
            if ('min' in range) {
                const rangeMin = this.parseMacro(range['min']);
                rangeFrom.min = rangeTo.min = rangeMin;
            }
            if ('max' in range) {
                const rangeMax = this.parseMacro(range['max']);
                rangeTo.max = rangeMax;
                if ('duration' in props) {
                    rangeFrom.max = new Date(rangeMax.getTime() - duration * unit);
                }
                else {
                    rangeFrom.max = rangeMax;
                }
            }
            return { rangeFrom, rangeTo };
        };
        /**
         * 根据起止日期重置范围
         *   type - 'from'/'to' 表示起止
         */
        this.resetRange = (type, date) => {
            let { rangeFrom, rangeTo, dateTo, timeTo } = Object.assign({}, this.state);
            const { duration, range = {} } = this.props;
            if (type === 'from') {
                if ((rangeFrom.min && this.compare(date, rangeFrom.min) < 0) ||
                    (rangeFrom.max && this.compare(date, rangeFrom.max) > 0)) {
                    return;
                }
                rangeTo.min = date;
                const rangeMax = this.parseMacro(range['max']);
                const rangeMin = this.parseMacro(range['min']);
                if ('maxLength' in range) {
                    const max = new Date(date.getTime() + range['maxLength'] * unit);
                    rangeTo.max = this.compare(max, rangeMax) < 0 ? max : rangeMax;
                    if (dateTo) {
                        const to = this.parse(`${dateTo} ${timeTo}`);
                        if (this.compare(to, rangeTo.min) < 0 || this.compare(to, rangeTo.max) > 0) {
                            this.setState({ dateTo: null, timeTo: null });
                        }
                    }
                }
                else {
                    rangeTo.max = rangeMax;
                }
                rangeFrom.min = rangeMin;
                if ('duration' in this.props) {
                    rangeFrom.max = new Date(rangeMax.getTime() - duration * unit);
                }
            }
            if (type === 'to') {
                if ((rangeTo.min && this.compare(date, rangeTo.min) < 0) ||
                    (rangeTo.max && this.compare(date, rangeTo.max) > 0)) {
                    return;
                }
                if (!('maxLength' in range)) {
                    rangeFrom.max = date;
                }
                // maxLength存在时不再限制rangeFrom.min
                // if ('maxLength' in range) {
                //   const min = new Date(date.getTime() - range.maxLength * unit);
                //   rangeFrom.min = this.compare(min, range.min) > 0 ? min : range.min;
                // }
            }
            this.setState({ rangeFrom, rangeTo });
        };
        /**
         * 判断时间是否选择完成
         */
        this.isCompleted = () => {
            const { dateFrom, timeFrom, dateTo, timeTo } = this.state;
            return dateFrom && timeFrom && dateTo && timeTo;
        };
        /**
         * Date比较，精确到s
         */
        this.compare = (value1, value2) => {
            if (value1.getFullYear() > value2.getFullYear())
                return 1;
            if (value1.getFullYear() < value2.getFullYear())
                return -1;
            if (value1.getMonth() > value2.getMonth())
                return 1;
            if (value1.getMonth() < value2.getMonth())
                return -1;
            if (value1.getDate() > value2.getDate())
                return 1;
            if (value1.getDate() < value2.getDate())
                return -1;
            if (value1.getHours() > value2.getHours())
                return 1;
            if (value1.getHours() < value2.getHours())
                return -1;
            if (value1.getMinutes() > value2.getMinutes())
                return 1;
            if (value1.getMinutes() < value2.getMinutes())
                return -1;
            if (value1.getSeconds() > value2.getSeconds())
                return 1;
            if (value1.getSeconds() < value2.getSeconds())
                return -1;
            return 0;
        };
        /**
         * 格式化
         */
        this.formatNum = (num) => num > 9 ? `${num}` : `0${num}`;
        this.formatDate = (value) => {
            if (typeof value === 'string')
                return value;
            if (!value)
                return "";
            const { year, month, day } = value;
            return `${year}-${this.formatNum(month + 1)}-${this.formatNum(day)}`;
        };
        this.formatTime = (value) => {
            if (typeof value === 'string')
                return value;
            if (!value)
                return "";
            const { hour, minute, second } = value;
            return `${this.formatNum(hour)}:${this.formatNum(minute)}:${this.formatNum(second)}`;
        };
        this.format = (date) => {
            return `${this.formatDate(this.getDate(date))}  ${this.formatTime(this.getTime(date))}`;
        };
        /**
         * 解析宏指令
         */
        this.parseMacro = (macro, type) => {
            if (typeof macro !== 'string') {
                return macro;
            }
            if (/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}\s[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(macro)) {
                return this.parse(macro);
            }
            const today = new Date();
            let unit, offset, stamp;
            if (macro.indexOf('%TODAY') >= 0) {
                macro = macro.replace('%TODAY', '');
                unit = 86400 * 1000;
                offset = +macro;
                const date = new Date(today.getTime() + (offset * unit));
                if (type === 'from') {
                    date.setHours(0);
                    date.setMinutes(0);
                    date.setSeconds(0);
                }
                if (type === 'to') {
                    date.setHours(23);
                    date.setMinutes(59);
                    date.setSeconds(59);
                }
                return date;
            }
            if (macro.indexOf('%NOW') >= 0) {
                macro = macro.replace('%NOW', '');
                unit = 1000;
                if (macro.indexOf('h') >= 0) {
                    macro = macro.replace('h', '');
                    unit = 3600 * 1000;
                }
                if (macro.indexOf('m') >= 0) {
                    macro = macro.replace('m', '');
                    unit = 60 * 1000;
                }
                if (macro.indexOf('s') >= 0) {
                    macro = macro.replace('s', '');
                }
                offset = +macro;
                const date = new Date(today.getTime() + (offset * unit));
                return date;
            }
        };
        /**
         * 'yyyy-MM-dd HH:mm:ss'解析为 Date 对象
         * 解决IE不能使用 new Date('yyyy-MM-dd HH:mm:ss')
         */
        this.parse = (str) => {
            const date = new Date(NaN);
            if (!/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}\s[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(str))
                return null;
            const year = +str.substr(0, 4), month = +str.substr(5, 2) - 1, day = +str.substr(8, 2), hour = +str.substr(11, 2), minute = +str.substr(14, 2), second = +str.substr(17, 2);
            date.setFullYear(year, month, day);
            date.setHours(hour);
            date.setMinutes(minute);
            date.setSeconds(second);
            return date;
        };
        this.getDate = (date) => {
            if (!date)
                return null;
            return {
                year: +date.getFullYear(),
                month: +date.getMonth(),
                day: +date.getDate()
            };
        };
        this.getTime = (date) => {
            if (!date)
                return null;
            return {
                hour: +date.getHours(),
                minute: +date.getMinutes(),
                second: +date.getSeconds()
            };
        };
        /**
         * 展开选择框
         */
        this.open = () => {
            if (this.props.disabled)
                return;
            const value = this.state.pickerValue;
            if (value) {
                this.setState({
                    dateFrom: this.formatDate(this.getDate(value.from)),
                    timeFrom: this.formatTime(this.getTime(value.from)),
                    dateTo: this.formatDate(this.getDate(value.to)),
                    timeTo: this.formatTime(this.getTime(value.to))
                }, () => {
                    this.resetRange('from', value.from);
                    if (!('duration' in this.props)) {
                        this.resetRange('to', value.to);
                    }
                });
            }
            this.setState({ active: true });
        };
        /**
         * 点击确认时调用
         */
        this.handleSubmit = () => {
            if (!this.isCompleted())
                return;
            const { dateFrom, timeFrom, dateTo, timeTo } = this.state;
            const pickerValue = {
                from: this.parse(`${dateFrom} ${timeFrom}`),
                to: this.parse(`${dateTo} ${timeTo}`),
            };
            this.setState({ active: false });
            const preValue = this.state.pickerValue;
            if (preValue && preValue.from && preValue.to &&
                this.compare(pickerValue.from, preValue.from) === 0 &&
                this.compare(pickerValue.to, preValue.to) === 0)
                return;
            if (this.compare(pickerValue.from, pickerValue.to) > 0) {
                return;
            }
            if (this.props.onChange) {
                this.props.onChange(pickerValue);
            }
            this.setState({
                pickerValue,
                selectedTabIndex: null,
                value: pickerValue
            });
        };
        /**
         * 点击Tab时调用
         */
        this.handleTabSelect = (tab, index) => {
            const value = {
                from: this.parseMacro(tab.from, 'from'),
                to: this.parseMacro(tab.to, 'to')
            };
            if (this.props.onChange) {
                this.props.onChange(value, tab.label);
            }
            // 不联动
            if (!this.props.linkage) {
                this.setState({
                    active: false,
                    selectedTabIndex: index,
                    pickerValue: null
                });
                return;
            }
            let dateFrom, timeFrom, dateTo, timeTo;
            dateFrom = this.formatDate(this.getDate(value.from));
            timeFrom = this.formatTime(this.getTime(value.from));
            this.resetRange('from', value.from);
            dateTo = this.formatDate(this.getDate(value.to));
            timeTo = this.formatTime(this.getTime(value.to));
            this.resetRange('to', value.to);
            this.setState({
                active: false,
                selectedTabIndex: index,
                dateFrom, timeFrom, dateTo, timeTo,
                pickerValue: value
            });
        };
        const { dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue } = this.getInitDateAndRange(this.props);
        this.state = {
            active: false,
            selectedTabIndex: null,
            pickerValue: pickerValue,
            dateFrom, timeFrom,
            rangeFrom,
            dateTo, timeTo,
            rangeTo
        };
    }
    componentDidMount() {
        const { tabs, defaultSelectedTabIndex } = this.props;
        if ('tabs' in this.props && 'defaultSelectedTabIndex' in this.props && tabs[defaultSelectedTabIndex]) {
            this.handleTabSelect(tabs[defaultSelectedTabIndex], defaultSelectedTabIndex);
        }
    }
    componentWillReceiveProps(nextProps) {
        // duration可变
        if (!('duration' in nextProps) || nextProps.duration === this.props.duration)
            return;
        const { dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue } = this.getInitDateAndRange(nextProps);
        this.setState({ dateFrom, timeFrom, dateTo, timeTo, rangeFrom, rangeTo, pickerValue, selectedTabIndex: null });
        // TODO range可变
    }
    /**
     * 点击取消时调用
     */
    handleCancel() {
        this.setState({ active: false });
    }
    render() {
        const { pickerValue, dateFrom, timeFrom, dateTo, timeTo, selectedTabIndex, rangeFrom, rangeTo } = this.state;
        const isCompleted = this.isCompleted();
        // picker内容显示
        let text = this.props.placeHolder || language.OptionDate;
        if (pickerValue) {
            text = `${this.format(pickerValue.from)} ${language.To} ${this.format(pickerValue.to)}`;
        }
        // Tabs
        const tabs = this.props.tabs ? this.props.tabs.map((tab, index) => {
            if (selectedTabIndex !== index) {
                return react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { role: "tab", tabIndex: index, key: index, onClick: () => this.handleTabSelect(tab, index) }, tab.label);
            }
            else {
                return react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { role: "tab", tabIndex: index, key: index, className: "current", onClick: () => this.handleTabSelect(tab, index) }, tab.label);
            }
        }) : [];
        // 时间起止范围
        const timeFromRange = {};
        if (dateFrom && dateFrom === this.formatDate(this.getDate(rangeFrom.min))) {
            timeFromRange.min = this.formatTime(this.getTime(rangeFrom.min));
        }
        if (dateFrom && dateFrom === this.formatDate(this.getDate(rangeFrom.max))) {
            timeFromRange.max = this.formatTime(this.getTime(rangeFrom.max));
        }
        const timeToRange = {};
        if (dateTo && dateTo === this.formatDate(this.getDate(rangeTo.min))) {
            timeToRange.min = this.formatTime(this.getTime(rangeTo.min));
        }
        if (dateTo && dateTo === this.formatDate(this.getDate(rangeTo.max))) {
            timeToRange.max = this.formatTime(this.getTime(rangeTo.max));
        }
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar-select-wrap tc-15-calendar2-hook" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { role: "tablist" }, tabs),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: classnames__WEBPACK_IMPORTED_MODULE_1___default()("tc-15-dropdown tc-15-dropdown-btn-style date-dropdown", { "tc-15-menu-active": this.state.active, "disabled": this.props.disabled }) },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "tc-15-dropdown-link", onClick: this.open, onFocus: this.open },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("i", { className: "caret" }),
                    text),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-dropdown-menu", role: "menu" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-custom-date" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "custom-date-wrap" },
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("em", null, "\u4ECE"),
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "calendar-box" },
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_SingleDatePicker__WEBPACK_IMPORTED_MODULE_3__["SingleDatePicker"], { value: dateFrom, range: { min: this.formatDate(this.getDate(rangeFrom.min)), max: this.formatDate(this.getDate(rangeFrom.max)) }, onChange: (dateFrom) => {
                                        const time = timeFrom ? timeFrom : (this.formatTime(this.getTime(rangeFrom.min)) || '00:00:00');
                                        if ('duration' in this.props) {
                                            const to = new Date(this.parse(`${dateFrom} ${time}`).getTime() + this.props.duration * unit);
                                            if (this.compare(to, rangeTo.max) > 0) {
                                                const from = new Date(rangeTo.max.getTime() - this.props.duration * unit);
                                                this.setState({
                                                    dateFrom: this.formatDate(this.getDate(from)), timeFrom: this.formatTime(this.getTime(from)),
                                                    dateTo: this.formatDate(this.getDate(rangeTo.max)), timeTo: this.formatTime(this.getTime(rangeTo.max))
                                                });
                                            }
                                            else {
                                                this.setState({
                                                    dateFrom, timeFrom: time,
                                                    dateTo: this.formatDate(this.getDate(to)), timeTo: this.formatTime(this.getTime(to))
                                                });
                                            }
                                        }
                                        else {
                                            this.resetRange('from', this.parse(`${dateFrom} ${time}`));
                                            this.setState({ dateFrom, timeFrom: time });
                                        }
                                    }, version: this.props.version }),
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_timepicker__WEBPACK_IMPORTED_MODULE_4__["TimePicker"], { value: timeFrom, range: timeFromRange, onChange: (timeFrom) => {
                                        if (dateFrom) {
                                            this.resetRange('from', this.parse(`${dateFrom} ${timeFrom}`));
                                            if ('duration' in this.props) {
                                                const to = new Date(this.parse(`${dateFrom} ${timeFrom}`).getTime() + this.props.duration * unit);
                                                this.setState({ dateTo: this.formatDate(this.getDate(to)), timeTo: this.formatTime(this.getTime(to)) });
                                            }
                                        }
                                        this.setState({ timeFrom });
                                    } }))),
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "custom-date-wrap" },
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("em", null, "\u81F3"),
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "calendar-box" },
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_SingleDatePicker__WEBPACK_IMPORTED_MODULE_3__["SingleDatePicker"], { value: dateTo, disabled: 'duration' in this.props, range: { min: this.formatDate(this.getDate(rangeTo.min)), max: this.formatDate(this.getDate(rangeTo.max)) }, onChange: (dateTo) => {
                                        const time = timeTo ? timeTo : (this.formatTime(this.getTime(rangeTo.min)) || '00:00:00');
                                        this.resetRange('to', this.parse(`${dateTo} ${time}`));
                                        this.setState({ dateTo, timeTo: time });
                                    }, version: this.props.version }),
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_timepicker__WEBPACK_IMPORTED_MODULE_4__["TimePicker"], { value: timeTo, disabled: 'duration' in this.props, range: timeToRange, onChange: (timeTo) => {
                                        if (dateTo) {
                                            this.resetRange('to', this.parse(`${dateTo} ${timeTo}`));
                                        }
                                        this.setState({ timeTo });
                                    } })))),
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "custom-date-ft" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("button", { type: "button", className: "tc-15-btn m", onClick: this.handleSubmit }, language.Confirm),
                        "\u00A0",
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("button", { type: "button", className: "tc-15-btn m weak", onClick: this.handleCancel.bind(this) }, language.Cancel))))));
    }
}
__decorate([
    tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_2__["OnOuterClick"]
], DateTimePicker.prototype, "handleCancel", null);


/***/ }),

/***/ "./src/tea-components/datetimepicker/SingleDatePicker.tsx":
/*!****************************************************************!*\
  !*** ./src/tea-components/datetimepicker/SingleDatePicker.tsx ***!
  \****************************************************************/
/*! exports provided: SingleDatePicker */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "SingleDatePicker", function() { return SingleDatePicker; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! classnames */ "./node_modules/classnames/index.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(classnames__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! tea-components/libs/decorators/OnOuterClick */ "./src/tea-components/libs/decorators/OnOuterClick.ts");
/* harmony import */ var _bubble__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ../bubble */ "./src/tea-components/bubble/index.ts");
/* harmony import */ var react_dom__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! react-dom */ "./webpack/alias/react-dom.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};





const keys = {
    "13": "enter"
};
const MONTH = [
    "",
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December"
];
class SingleDatePicker extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        /**
         * 获取当前年月
         */
        this.getCurYearAndMonth = range => {
            const date = new Date();
            let curYear = date.getFullYear(), curMonth = date.getMonth();
            if (range) {
                if (range.max) {
                    const max = this.parse(range.max);
                    curYear = max.year;
                    curMonth = max.month;
                }
                if (range.min) {
                    const min = this.parse(range.min);
                    curYear = min.year;
                    curMonth = min.month;
                }
            }
            return { curYear, curMonth };
        };
        /**
         * 格式化
         */
        this.formatNum = (num) => (num > 9 ? `${num}` : `0${num}`);
        this.format = (value) => {
            if (typeof value === "string")
                return value;
            if (!value)
                return "";
            const { year, month, day } = value;
            return `${year}-${this.formatNum(month + 1)}-${this.formatNum(day)}`;
        };
        this.getGlobalFormat = (year, month) => {
            const version = this.props.version || window["VERSION"];
            if (version === "en") {
                return `${MONTH[month]} ${year}`;
            }
            else {
                return `${year}年${this.formatNum(month)}月`;
            }
        };
        /**
         * 检验传入value值是否合法
         */
        this.check = (value) => {
            if (!value)
                return false;
            // TODO
            if (!("year" in value) || !("month" in value) || !("day" in value))
                return false;
            const { year, month, day } = value;
            var date = new Date(year, month, day);
            if (date.getFullYear() != year || date.getMonth() != month || date.getDate() != day)
                return false;
            return true;
        };
        /**
         * 字符串解析为 SingleDatePickerValue
         */
        this.parse = (value) => {
            if (typeof value !== "string")
                return value;
            if (!/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}$/.test(value))
                return null;
            return {
                year: +value.substr(0, 4),
                month: +value.substr(5, 2) - 1,
                day: +value.substr(8, 2)
            };
        };
        /**
         * 确认value，合法将调用 onChange
         */
        this.confirm = (value) => {
            const range = this.props.range || {};
            if (range.min && this.format(value) < range.min)
                return;
            if (range.max && this.format(value) > range.max)
                return;
            if (!("value" in this.props)) {
                this.setState({ value });
            }
            if (this.props.onChange) {
                this.props.onChange(this.format(value));
            }
            this.setState({
                inputValue: this.format(value),
                curYear: value.year,
                curMonth: value.month
            });
        };
        this.handleInputChange = (e) => {
            const inputValue = e.target.value;
            if (/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}$/.test(inputValue)) {
                const value = {
                    year: +inputValue.substr(0, 4),
                    month: +inputValue.substr(5, 2) - 1,
                    day: +inputValue.substr(8, 2)
                };
                if (this.check(value)) {
                    this.confirm(value);
                    this.setState({ bubbleActive: false });
                }
            }
            else {
                this.setState({ bubbleActive: true });
            }
            if (inputValue.length === 0) {
                this.setState({ bubbleActive: false });
            }
            this.setState({ inputValue });
        };
        this.open = () => {
            const value = this.parse(this.props.value);
            if ("value" in this.props) {
                this.setState({ value });
                const curValue = this.check(value) ? value : null;
                if (curValue) {
                    this.setState({ curYear: curValue.year, curMonth: curValue.month });
                }
                else {
                    this.setState(this.getCurYearAndMonth(this.props.range));
                }
            }
            else {
                this.setState(this.getCurYearAndMonth(this.props.range));
            }
            if (!this.state.active && !this.props.disabled) {
                this.setState({ active: true });
            }
        };
        this.handleKeyDown = (e) => {
            if (!keys[e.keyCode])
                return;
            e.preventDefault();
            switch (keys[e.keyCode]) {
                case "enter":
                    this.close();
                    react_dom__WEBPACK_IMPORTED_MODULE_4__["findDOMNode"](this["input"]).blur();
                    break;
            }
        };
        this.handleSelect = (year, month, day) => {
            this.confirm({ year, month, day });
            setTimeout(() => this.close(), 100);
        };
        this.calendarRender = (year, month, range) => {
            let value = this.state.value;
            if ("value" in this.props) {
                value = this.parse(this.props.value);
            }
            if (!this.check(value)) {
                value = {};
            }
            // 本月第一天
            const firstDate = new Date(year, month, 1);
            let day = firstDate.getDay();
            // 获取本月有多少天
            const lastDate = new Date(year, month + 1, 0);
            const count = lastDate.getDate();
            const weeks = [];
            for (let i = 1; i <= count;) {
                const week = [];
                // 第一周补全
                for (let j = 0; j < day % 7; ++j) {
                    week.push(react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("td", { key: `dis-${j}`, className: "tc-15-calendar-dis" }));
                }
                // 填充一周
                do {
                    (day => {
                        if (day == value.day && month == value.month && year == value.year) {
                            week.push(react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("td", { key: day, className: "tc-15-calendar-today", onClick: () => this.handleSelect(year, month, day) }, day));
                        }
                        else {
                            const cur = this.format({ year, month, day });
                            if ((range.min && cur < range.min) || (range.max && cur > range.max)) {
                                week.push(react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("td", { key: day, className: "tc-15-calendar-dis" }, day));
                            }
                            else {
                                week.push(react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("td", { key: day, onClick: () => this.handleSelect(year, month, day) }, day));
                            }
                        }
                    })(i++);
                    if (i > count)
                        break;
                } while (++day % 7 !== 0);
                weeks.push(react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("tr", { key: `week-${i / 7}` }, week));
            }
            return weeks;
        };
        this.prevBtnRender = (range) => {
            const { curYear, curMonth } = this.state;
            if (range.min) {
                const minRange = this.parse(range.min);
                if (minRange.year >= curYear && minRange.month >= curMonth) {
                    return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("i", { tabIndex: 0, className: "tc-15-calendar-i-pre-m disabled" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("b", null,
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", null, "\u8FC7\u53BB\u65F6\u95F4\u4E0D\u53EF\u9009"))));
                }
            }
            return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("i", { tabIndex: 0, className: "tc-15-calendar-i-pre-m", onClick: () => {
                    if (curMonth > 0)
                        this.setState({ curMonth: curMonth - 1 });
                    else
                        this.setState({ curMonth: 11, curYear: curYear - 1 });
                } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("b", null,
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", null, "\u8F6C\u5230\u4E0A\u4E2A\u6708"))));
        };
        this.nextBtnRender = (range) => {
            const { curYear, curMonth } = this.state;
            if (range.max) {
                const maxRange = this.parse(range.max);
                if (maxRange.year <= curYear && maxRange.month <= curMonth) {
                    return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("i", { tabIndex: 0, className: "tc-15-calendar-i-next-m disabled" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("b", null,
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", null, "\u672A\u6765\u65F6\u95F4\u4E0D\u53EF\u9009"))));
                }
            }
            return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("i", { tabIndex: 0, className: "tc-15-calendar-i-next-m", onClick: () => {
                    if (curMonth < 11)
                        this.setState({ curMonth: curMonth + 1 });
                    else
                        this.setState({ curMonth: 0, curYear: curYear + 1 });
                } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("b", null,
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", null, "\u8F6C\u5230\u4E0B\u4E2A\u6708"))));
        };
        let defaultValue = this.props.value || this.props.defaultValue;
        defaultValue = this.parse(defaultValue);
        const value = this.check(defaultValue) ? defaultValue : null;
        let cur = {};
        if (value) {
            cur = {
                curYear: value.year,
                curMonth: value.month
            };
        }
        else {
            cur = this.getCurYearAndMonth(this.props.range);
        }
        this.state = {
            value,
            active: false,
            bubbleActive: false,
            curYear: cur.curYear,
            curMonth: cur.curMonth,
            inputValue: this.format(value)
        };
        this.close = this.close.bind(this);
    }
    componentWillReceiveProps(nextProps) {
        // 受控组件
        if ("value" in nextProps) {
            const nextValue = this.parse(nextProps.value);
            const value = this.check(nextValue) ? nextValue : null;
            this.setState({ value: nextProps.value, inputValue: this.format(value) });
            if (value) {
                this.setState({ curYear: value.year, curMonth: value.month });
            }
            else {
                this.setState(this.getCurYearAndMonth(nextProps.range));
            }
        }
        else {
            this.setState(this.getCurYearAndMonth(nextProps.range));
        }
    }
    close() {
        const { inputValue } = this.state;
        let value = this.state.value;
        if ("value" in this.props) {
            value = this.props.value;
            this.setState({ inputValue: this.format(value) });
        }
        if (!/^[0-9]{4}[/\-\.][0-9]{2}[/\-\.][0-9]{2}$/.test(inputValue)) {
            this.setState({ inputValue: this.format(value) });
        }
        const curValue = {
            year: +inputValue.substr(0, 4),
            month: +inputValue.substr(5, 2) - 1,
            day: +inputValue.substr(8, 2)
        };
        if (!this.check(curValue)) {
            this.setState({ inputValue: this.format(value) });
        }
        this.setState({ active: false, bubbleActive: false });
    }
    render() {
        const { curYear, curMonth } = this.state;
        let range = {};
        if (this.props.range) {
            range.min = this.props.range.min;
            range.max = this.props.range.max;
        }
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar-select-wrap tc-15-calendar2-hook" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: classnames__WEBPACK_IMPORTED_MODULE_1___default()("tc-15-calendar-select", "tc-15-calendar-single", { show: this.state.active }) },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_bubble__WEBPACK_IMPORTED_MODULE_3__["BubbleWrapper"], { align: "start" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("input", { ref: el => (this["input"] = el), disabled: this.props.disabled, className: "tc-15-simulate-select m show", value: this.state.inputValue, onChange: this.handleInputChange, onClick: this.open, onFocus: this.open, placeholder: this.props.placeHolder || "日期选择", onKeyDown: this.handleKeyDown, maxLength: 10 }),
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_bubble__WEBPACK_IMPORTED_MODULE_3__["Bubble"], { focusPosition: "bottom", className: "error", style: { display: this.state.bubbleActive ? "" : "none" } }, "\u683C\u5F0F\u9519\u8BEF\uFF0C\u5E94\u4E3Ayyyy-MM-dd")),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar-triangle-wrap" }),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar-triangle" }),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar tc-15-calendar2" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar-cont" },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("table", { cellSpacing: 0, className: "tc-15-calendar-left" },
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("caption", null, this.getGlobalFormat(curYear, curMonth + 1)),
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("thead", null,
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("tr", null,
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", null, "\u65E5"),
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", null, "\u4E00"),
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", null, "\u4E8C"),
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", null, "\u4E09"),
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", null, "\u56DB"),
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", null, "\u4E94"),
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("th", null, "\u516D"))),
                            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("tbody", null,
                                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("tr", null,
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("td", { colSpan: 3 }, this.prevBtnRender(range)),
                                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("td", { colSpan: 4 }, this.nextBtnRender(range))),
                                this.calendarRender(curYear, curMonth, range)))),
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-calendar-for-style" })))));
    }
}
__decorate([
    tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_2__["OnOuterClick"]
], SingleDatePicker.prototype, "close", null);


/***/ }),

/***/ "./src/tea-components/datetimepicker/index.ts":
/*!****************************************************!*\
  !*** ./src/tea-components/datetimepicker/index.ts ***!
  \****************************************************/
/*! exports provided: DateTimePicker, SingleDatePicker */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _DateTimePicker__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./DateTimePicker */ "./src/tea-components/datetimepicker/DateTimePicker.tsx");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "DateTimePicker", function() { return _DateTimePicker__WEBPACK_IMPORTED_MODULE_0__["DateTimePicker"]; });

/* harmony import */ var _SingleDatePicker__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./SingleDatePicker */ "./src/tea-components/datetimepicker/SingleDatePicker.tsx");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "SingleDatePicker", function() { return _SingleDatePicker__WEBPACK_IMPORTED_MODULE_1__["SingleDatePicker"]; });





/***/ }),

/***/ "./src/tea-components/libs/decorators/OnOuterClick.ts":
/*!************************************************************!*\
  !*** ./src/tea-components/libs/decorators/OnOuterClick.ts ***!
  \************************************************************/
/*! exports provided: OnOuterClick */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "OnOuterClick", function() { return OnOuterClick; });
/* harmony import */ var react_dom__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react-dom */ "./webpack/alias/react-dom.js");
/* harmony import */ var _helpers_AppendFunction__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../helpers/AppendFunction */ "./src/tea-components/libs/helpers/AppendFunction.ts");


/**
 * 在组件外部被点击的时候，会调用被装饰的方法
 */
function OnOuterClick(target, propertyKey, descriptor) {
    let onMount = target.componentDidMount;
    let onUnmount = target.componentWillUnmount;
    let callbackMethod = descriptor.value;
    // let bindStack = null;
    // let errorWarned = false;
    function outerClickHandler(e) {
        const _this = this;
        let clickTarget = null;
        try {
            clickTarget = react_dom__WEBPACK_IMPORTED_MODULE_0__["findDOMNode"](_this);
        }
        catch (err) { }
        if (!clickTarget) {
            // if (!errorWarned) {
            //   const style = 'background-color: #fdf0f0; color: red;';
            //   const message = "OnOuterClick 无法通过绑定的组件找到对应的 DOM，这就意味着你的组件并没有正确地 Unmount。请确保使用 `ReactDOM.unmountComponentAtNode()` 方法来销毁组件，而不是直接清除其使用的 DOM。";
            //   if (console.groupCollapsed) {
            //     (console.groupCollapsed as any)('%c%s', style, message)
            //   }
            //   else {
            //     console.log('%c%s', style, message);
            //   }
            //   console.log("绑定组件：%o", target.constructor);
            //   console.log("绑定方法：%o", callbackMethod)
            //   if (bindStack) {
            //     console.info('绑定堆栈: %s', bindStack);
            //   }
            //   if (console.groupEnd) {
            //     console.groupEnd();
            //   }
            // }
            // errorWarned = true;
            return;
        }
        let isClickedOutside = !clickTarget.contains(e.target);
        if (isClickedOutside) {
            _this.$$outsideClickRegisterMethods.forEach(invoke => invoke());
        }
    }
    let bind = function () {
        const _this = this;
        if (!_this.$$outsideClickHandler) {
            _this.$$outsideClickHandler = outerClickHandler.bind(_this);
            _this.$$outsideClickRegisterMethods = [];
            document.addEventListener("mousedown", _this.$$outsideClickHandler, false);
        }
        _this.$$outsideClickRegisterMethods.push(callbackMethod.bind(_this));
        // try {
        //   bindStack = (new Error("Bind OnOuterClick") as any).stack;
        // } catch (err) {}
    };
    let unbind = function () {
        const _this = this;
        if (_this.$$outsideClickHandler) {
            document.removeEventListener("mousedown", _this.$$outsideClickHandler);
            _this.$$outsideClickHandler = undefined;
            _this.$$outsideClickRegisterMethods = undefined;
        }
    };
    target.componentDidMount = onMount ? Object(_helpers_AppendFunction__WEBPACK_IMPORTED_MODULE_1__["appendFunction"])(onMount, bind) : bind;
    target.componentWillUnmount = onUnmount ? Object(_helpers_AppendFunction__WEBPACK_IMPORTED_MODULE_1__["appendFunction"])(unbind, onUnmount) : unbind;
    return descriptor;
}


/***/ }),

/***/ "./src/tea-components/libs/helpers/AppendFunction.ts":
/*!***********************************************************!*\
  !*** ./src/tea-components/libs/helpers/AppendFunction.ts ***!
  \***********************************************************/
/*! exports provided: appendFunction */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "appendFunction", function() { return appendFunction; });
/**
 * 在指定函数后面追加一个函数
 */
function appendFunction(origin, append) {
    return function appended(...args) {
        let result = origin.apply(this, args);
        return append.apply(this, [result].concat(args));
    };
}


/***/ }),

/***/ "./src/tea-components/tagsearchbox/AttributeSelect.tsx":
/*!*************************************************************!*\
  !*** ./src/tea-components/tagsearchbox/AttributeSelect.tsx ***!
  \*************************************************************/
/*! exports provided: AttributeSelect */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "AttributeSelect", function() { return AttributeSelect; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");

const keys = {
    "8": 'backspace',
    "9": 'tab',
    "13": 'enter',
    "37": 'left',
    "38": 'up',
    "39": 'right',
    "40": 'down'
};
class AttributeSelect extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor() {
        super(...arguments);
        this.state = {
            select: -1,
        };
        // 父组件调用
        this.handleKeyDown = (keyCode) => {
            if (!keys[keyCode])
                return;
            const { onSelect } = this.props;
            const select = this.state.select;
            switch (keys[keyCode]) {
                case 'enter':
                case 'tab':
                    if (select < 0)
                        break;
                    if (onSelect) {
                        onSelect(this.getAttribute(select));
                    }
                    return false;
                case 'up':
                    this.move(-1);
                    break;
                case 'down':
                    this.move(1);
                    break;
            }
        };
        this.move = (step) => {
            const select = this.state.select;
            const list = this.getUseableList();
            if (list.length <= 0)
                return;
            this.setState({ select: (select + step + list.length) % list.length });
        };
        this.handleClick = (e, index) => {
            e.stopPropagation();
            if (this.props.onSelect) {
                this.props.onSelect(this.getAttribute(index));
            }
        };
    }
    componentWillReceiveProps(nextProps) {
        if (this.props.inputValue !== nextProps.inputValue) {
            this.setState({ select: -1 });
        }
    }
    getUseableList() {
        const { attributes, inputValue } = this.props;
        // 获取冒号前字符串模糊查询
        const fuzzyValue = /(.*?)(:|：).*/.test(inputValue) ? RegExp.$1 : inputValue;
        return attributes.filter(item => item.name.includes(inputValue) || item.name.includes(fuzzyValue));
    }
    getAttribute(select) {
        const list = this.getUseableList();
        if (select < list.length) {
            return list[select];
        }
    }
    render() {
        const select = this.state.select;
        const list = this.getUseableList().map((item, index) => {
            if (select === index) {
                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation", key: index, className: "autocomplete-cur", onClick: (e) => this.handleClick(e, index) },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "text-truncate", role: "menuitem", href: "javascript:;" }, item.name)));
            }
            return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation", key: index, onClick: (e) => this.handleClick(e, index) },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "text-truncate", role: "menuitem", href: "javascript:;" }, item.name)));
        });
        if (list.length === 0)
            return null;
        let maxHeight = document.body.clientHeight ? document.body.clientHeight - 450 : 400;
        maxHeight = maxHeight > 240 ? maxHeight : 240;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-autocomplete", style: { minWidth: 180, width: 'auto' } },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("ul", { className: "tc-15-autocomplete-menu", role: "menu", style: { maxHeight: `${maxHeight}px` } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "autocomplete-empty", role: "menuitem", href: "javascript:;" }, "\u9009\u62E9\u7EF4\u5EA6")),
                list)));
    }
}


/***/ }),

/***/ "./src/tea-components/tagsearchbox/Input.tsx":
/*!***************************************************!*\
  !*** ./src/tea-components/tagsearchbox/Input.tsx ***!
  \***************************************************/
/*! exports provided: Input */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Input", function() { return Input; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _AttributeSelect__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./AttributeSelect */ "./src/tea-components/tagsearchbox/AttributeSelect.tsx");
/* harmony import */ var _valueselect__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./valueselect */ "./src/tea-components/tagsearchbox/valueselect/index.tsx");



const keys = {
    "8": "backspace",
    "9": "tab",
    "13": "enter",
    "37": "left",
    "38": "up",
    "39": "right",
    "40": "down"
};
const INPUT_MIN_SIZE = 0;
class Input extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.state = {
            inputWidth: INPUT_MIN_SIZE,
            inputValue: "",
            attribute: null,
            values: [],
            showAttrSelect: false,
            showValueSelect: false,
            ValueSelectOffset: 0
        };
        /**
         * 刷新选择组件显示
         */
        this.refreshShow = () => {
            const { inputValue, attribute } = this.state;
            const input = this["input"];
            let start = input.selectionStart, end = input.selectionEnd;
            // if (start !== end) {
            //   this.setState({ showAttrSelect: false, showValueSelect: false });
            //   return;
            // }
            const { pos } = this.getAttrStrAndValueStr(inputValue);
            if (pos < 0 || start <= pos) {
                this.setState({ showAttrSelect: true, showValueSelect: false });
                return;
            }
            if (attribute && end > pos) {
                this.setState({ showAttrSelect: false, showValueSelect: true });
            }
        };
        this.focusInput = () => {
            if (!this["input"])
                return;
            const input = this["input"];
            input.focus();
        };
        this.moveToEnd = () => {
            const input = this["input"];
            input.focus();
            const value = this.state.inputValue;
            setTimeout(() => input.setSelectionRange(value.length, value.length));
        };
        this.selectValue = () => {
            const input = this["input"];
            input.focus();
            const value = this.state.inputValue;
            let { pos } = this.getAttrStrAndValueStr(value);
            if (pos < 0)
                pos = -2;
            setTimeout(() => {
                input.setSelectionRange(pos + 2, value.length);
                this.refreshShow();
            });
        };
        this.selectAttr = () => {
            const input = this["input"];
            input.focus();
            const value = this.state.inputValue;
            let { pos } = this.getAttrStrAndValueStr(value);
            if (pos < 0)
                pos = 0;
            setTimeout(() => {
                input.setSelectionRange(0, pos);
                this.refreshShow();
            });
        };
        this.setInputValue = (value, callback) => {
            if (this.props.type === "edit" && value.trim().length <= 0) {
                this.props.dispatchTagEvent("del", "edit");
            }
            // value = value.replace(/：/g, ':');
            // const pos = value.indexOf(':');
            // let attrStr = value, valueStr = '';
            // if (pos >= 0) {
            //   attrStr = value.substr(0, pos);
            //   valueStr = value.substr(pos+1).replace(/^\s*(.*)/, '$1');
            // }
            const attributes = this.props.attributes;
            let attribute = null, valueStr = value;
            const input = this["input"];
            const mirror = this["input-mirror"];
            // attribute 是否存在
            for (let i = 0; i < attributes.length; ++i) {
                if (value.indexOf(attributes[i].name + ":") === 0 || value.indexOf(attributes[i].name + "：") === 0) {
                    // 获取属性/值
                    attribute = attributes[i];
                    valueStr = value.substr(attributes[i].name.length + 1);
                    // 计算 offset
                    mirror.innerText = attribute.type === "onlyKey" ? attribute.name : attribute.name + ": ";
                    let width = mirror.clientWidth;
                    if (this.props.inputOffset)
                        width += this.props.inputOffset;
                    this.setState({ ValueSelectOffset: width });
                    break;
                }
            }
            // 处理前导空格
            if (attribute && valueStr.replace(/^\s+/, "").length > 0) {
                value = `${attribute.name}: ${valueStr.replace(/^\s+/, "")}`;
            }
            else if (attribute) {
                value = attribute.type === "onlyKey" ? attribute.name : `${attribute.name}:${valueStr}`;
            }
            this.setState({ attribute }, this.refreshShow);
            if (this.props.type === "edit") {
                this.props.dispatchTagEvent("editing", { attr: attribute });
            }
            mirror.innerText = value;
            const width = mirror.clientWidth > INPUT_MIN_SIZE ? mirror.clientWidth : INPUT_MIN_SIZE;
            this.setState({ inputValue: value, inputWidth: width }, () => {
                if (callback)
                    callback();
            });
        };
        this.resetInput = (callback) => {
            this.setInputValue("", callback);
            this.setState({ inputWidth: INPUT_MIN_SIZE });
        };
        this.getInputValue = () => {
            return this.state.inputValue;
        };
        // getInputAttr = (): AttributeValue => {
        //   return this.state.attribute;
        // }
        this.addTagByInputValue = () => {
            const { attribute, values, inputValue } = this.state;
            const type = this.props.type || "add";
            // 属性值搜索
            if (attribute && this.props.attributes.filter(item => item.key === attribute.key).length > 0) {
                if (values.length <= 0) {
                    return false;
                }
                this.props.dispatchTagEvent(type, { attr: attribute, values: values });
            }
            else {
                // 关键字搜索
                if (inputValue.trim().length <= 0) {
                    return false;
                }
                let attribute = this.props.attributes.find(item => item.name === inputValue.trim());
                if (!attribute || attribute.type === "onlyKey") {
                    this.props.dispatchTagEvent(type, { attr: attribute, values: [] });
                    this.resetInput();
                }
                else {
                    const list = inputValue
                        .split("|")
                        .filter(item => item.trim().length > 0)
                        .map(item => {
                        return { name: item.trim() };
                    });
                    this.props.dispatchTagEvent(type, { attr: null, values: list });
                }
            }
            this.setState({ showAttrSelect: false, showValueSelect: false });
            if (this.props.type !== "edit") {
                this.resetInput();
            }
            return true;
        };
        this.handleInputChange = (e) => {
            this.setInputValue(e.target.value);
        };
        this.handleInputClick = (e) => {
            this.props.dispatchTagEvent("click-input", this.props.type);
            e.stopPropagation();
            this.focusInput();
        };
        this.handleAttrSelect = (attr) => {
            if (attr && attr.key) {
                const str = attr.type === "onlyKey" ? attr.name : `${attr.name}: `;
                const inputValue = this.state.inputValue;
                if (inputValue.indexOf(str) >= 0) {
                    this.selectValue();
                }
                else {
                    this.setInputValue(str);
                }
                this.setState({ values: [] });
            }
            if (attr.type === "onlyKey") {
                // 不需要值
                // this.addTagByInputValue()
                setTimeout(() => this.addTagByInputValue());
            }
            else {
                this.focusInput();
            }
        };
        this.handleValueChange = (values) => {
            this.setState({ values }, () => {
                this.setInputValue(`${this.state.attribute.name}: ${values.map(item => item.name).join(" | ")}`);
                this.focusInput();
            });
        };
        /**
         * 值选择组件完成选择
         */
        this.handleValueSelect = (values) => {
            this.setState({ values });
            const inputValue = this.state.inputValue;
            if (values.length <= 0) {
                this.setInputValue(this.state.attribute.name + ": ");
                return;
            }
            if (values.length > 0) {
                const key = this.state.attribute.key;
                if (this.props.attributes.filter(item => item.key === key).length > 0) {
                    const type = this.props.type || "add";
                    this.props.dispatchTagEvent(type, { attr: this.state.attribute, values });
                }
                this.focusInput();
            }
            if (this.props.type !== "edit") {
                this.resetInput();
            }
        };
        /**
         * 值选择组件取消选择
         */
        this.handleValueCancel = () => {
            if (this.props.type === "edit") {
                const { attribute, values } = this.state;
                this.props.dispatchTagEvent("edit-cancel", { attr: attribute, values: values });
            }
            else {
                this.resetInput(() => {
                    this.focusInput();
                });
            }
        };
        /**
         * 处理粘贴事件
         */
        this.handlePaste = (e) => {
            const { attribute } = this.state;
            if (!attribute || attribute.type === "input") {
                this["textarea"].focus();
                setTimeout(() => {
                    let value = this["textarea"].value;
                    if (/^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/.test(value)) {
                        value = value.replace(/[\r\n\t,，\s]+/g, "|");
                    }
                    else {
                        value = value.replace(/[\r\n\t,，]+/g, "|");
                    }
                    value = value
                        .split("|")
                        .map(item => item.trim())
                        .filter(item => item.length > 0)
                        .join(" | ");
                    const input = this["input"];
                    const start = input.selectionStart, end = input.selectionEnd;
                    const inputValue = this.state.inputValue;
                    // 覆盖选择区域
                    const curValue = inputValue.substring(0, start) + value + inputValue.substring(end, inputValue.length);
                    // input 属性情况
                    this["textarea"].value = "";
                    if (attribute && attribute.type === "input") {
                        this.setInputValue(curValue, this.focusInput);
                        return;
                    }
                    if (inputValue.length > 0) {
                        this.setInputValue(curValue, this.focusInput);
                    }
                    else {
                        this.setInputValue(curValue, this.addTagByInputValue);
                    }
                }, 100);
            }
        };
        // 键盘事件
        // handlekeyUp = (e): void => {
        //   if (this['value-select']) {
        //     if (this['value-select'].handlekeyUp(e.keyCode) === false) return;
        //   }
        // }
        this.handlekeyDown = (e) => {
            if (!keys[e.keyCode])
                return;
            // if (!this.props.isFocused) {
            //   this.props.dispatchTagEvent('click-input', this.props.type);
            // }
            if (this.props.hidden) {
                return this.props.handleKeyDown(e);
            }
            const inputValue = this.state.inputValue;
            if (keys[e.keyCode] === "backspace" && inputValue.length > 0)
                return;
            if ((keys[e.keyCode] === "left" || keys[e.keyCode] === "right") && inputValue.length > 0) {
                setTimeout(this.refreshShow, 0);
                return;
            }
            e.preventDefault();
            // 事件下传
            if (this["attr-select"]) {
                if (this["attr-select"].handleKeyDown(e.keyCode) === false)
                    return;
            }
            if (this["value-select"]) {
                if (this["value-select"].handleKeyDown(e.keyCode) === false)
                    return;
            }
            switch (keys[e.keyCode]) {
                case "enter":
                case "tab":
                    if (!this.props.isFocused) {
                        this.props.dispatchTagEvent("click-input");
                    }
                    this.addTagByInputValue();
                    break;
                case "backspace":
                    this.props.dispatchTagEvent("del", "keyboard");
                    break;
                case "up":
                    break;
                case "down":
                    break;
                case "left":
                    this.props.dispatchTagEvent("move-left");
                    break;
                case "right":
                    this.props.dispatchTagEvent("move-right");
                    break;
            }
        };
        this.getAttrStrAndValueStr = (str) => {
            let attrStr = str, valueStr = "", pos = -1;
            const attributes = this.props.attributes;
            for (let i = 0; i < attributes.length; ++i) {
                if (str.indexOf(attributes[i].name + ":") === 0) {
                    // 获取属性/值
                    attrStr = attributes[i].name;
                    valueStr = str.substr(attrStr.length + 1);
                    pos = attributes[i].name.length;
                }
            }
            return { attrStr, valueStr, pos };
        };
    }
    componentDidMount() { }
    setInfo(info, callback) {
        const attribute = info.attr;
        const values = info.values;
        this.setState({ attribute, values }, () => {
            if (attribute) {
                this.setInputValue(`${attribute.name}: ${values.map(item => item.name).join(" | ")}`, callback);
            }
            else {
                this.setInputValue(`${values.map(item => item.name).join(" | ")}`, callback);
            }
        });
    }
    render() {
        const { inputWidth, inputValue, showAttrSelect, showValueSelect, attribute, ValueSelectOffset } = this.state;
        const { active, attributes, isFocused, hidden, maxWidth, type } = this.props;
        // const pos = inputValue.indexOf(':');
        // let attrStr = inputValue, valueStr = '';
        // if (pos >= 0) {
        //   attrStr = inputValue.substr(0, pos).trim();
        //   valueStr = inputValue.substr(pos+1).trim();
        // }
        const { attrStr, valueStr } = this.getAttrStrAndValueStr(inputValue);
        const attrSelect = isFocused && showAttrSelect ? (react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_AttributeSelect__WEBPACK_IMPORTED_MODULE_1__["AttributeSelect"], { ref: select => (this["attr-select"] = select), attributes: attributes, inputValue: attrStr, onSelect: this.handleAttrSelect })) : null;
        const valueSelect = isFocused && showValueSelect && attribute && attribute.type ? (react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_valueselect__WEBPACK_IMPORTED_MODULE_2__["ValueSelect"], { type: attribute.type, ref: select => (this["value-select"] = select), values: attribute.values, inputValue: valueStr.trim(), offset: ValueSelectOffset, onChange: this.handleValueChange, onSelect: this.handleValueSelect, onCancel: this.handleValueCancel })) : null;
        const style = {
            width: hidden ? "0px" : active ? `${inputWidth + 5}px` : "5px",
            maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px"
        };
        if (type === "edit" && !hidden) {
            style["padding"] = "0 8px";
        }
        const input = type !== "edit" ? (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("input", { ref: input => (this["input"] = input), type: "text", className: "tc-search-input", placeholder: "", style: {
                width: hidden ? "0px" : `${inputWidth + 5}px`,
                display: active ? "" : "none",
                maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px"
            }, value: inputValue, onChange: this.handleInputChange, onKeyDown: this.handlekeyDown, onFocus: this.refreshShow, onClick: this.refreshShow, onPaste: this.handlePaste })) : (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: { position: "relative", display: hidden ? "none" : "" } },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("pre", { style: { display: "block", visibility: "hidden" } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { style: {
                        fontSize: "12px",
                        width: hidden ? "0px" : `${inputWidth + 36}px`,
                        maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px",
                        whiteSpace: "normal"
                    } }, inputValue),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("br", { style: { clear: "both" } })),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("textarea", { ref: input => (this["input"] = input), className: "tc-search-input", placeholder: "", style: {
                    width: hidden ? "0px" : `${inputWidth + 30}px`,
                    display: active ? "" : "none",
                    maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px",
                    position: "absolute",
                    top: 0,
                    left: 0,
                    height: "100%",
                    resize: "none",
                    minHeight: "15px",
                    marginTop: "2px"
                }, value: inputValue, onChange: this.handleInputChange, onKeyDown: this.handlekeyDown, onFocus: this.refreshShow, onClick: this.refreshShow, onPaste: this.handlePaste })));
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-tags-space", style: style, onClick: this.handleInputClick },
            input,
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { ref: input => (this["input-mirror"] = input), style: { position: "absolute", top: "-9999px", left: 0, whiteSpace: "pre", fontSize: "12px" } }),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("textarea", { ref: textarea => (this["textarea"] = textarea), style: { position: "absolute", top: "-9999px", left: 0, whiteSpace: "pre", fontSize: "12px" } }),
            attrSelect,
            valueSelect));
    }
}


/***/ }),

/***/ "./src/tea-components/tagsearchbox/Tag.tsx":
/*!*************************************************!*\
  !*** ./src/tea-components/tagsearchbox/Tag.tsx ***!
  \*************************************************/
/*! exports provided: Tag */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Tag", function() { return Tag; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! classnames */ "./node_modules/classnames/index.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(classnames__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var _Input__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./Input */ "./src/tea-components/tagsearchbox/Input.tsx");
/* harmony import */ var _TagSearchBox__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./TagSearchBox */ "./src/tea-components/tagsearchbox/TagSearchBox.tsx");




const keys = {
    "8": 'backspace',
    "13": 'enter',
    "37": 'left',
    "38": 'up',
    "39": 'right',
    "40": 'down'
};
const INPUT_MIN_SIZE = 0;
class Tag extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor() {
        super(...arguments);
        this.state = {
            bubbleActice: false,
            inEditing: false,
            inputOffset: 0
        };
        this.handleTagClick = (e, pos) => {
            this.props.dispatchTagEvent('click', pos);
            e.stopPropagation();
        };
        this.handleDelete = (e) => {
            e.stopPropagation();
            this.props.dispatchTagEvent('del');
        };
        this.handleKeyDown = (e) => {
            if (!keys[e.keyCode])
                return;
            e.preventDefault();
            switch (keys[e.keyCode]) {
                case 'tab':
                case 'enter':
                    this.props.dispatchTagEvent('click', 'value');
                    break;
                case 'backspace':
                    this.props.dispatchTagEvent('del', 'keyboard');
                    break;
                case 'left':
                    this.props.dispatchTagEvent('move-left');
                    break;
                case 'right':
                    this.props.dispatchTagEvent('move-right');
                    break;
            }
        };
    }
    componentDidMount() {
        // console.log(this['content'].clientWidth);
        // this.setState({inputOffset: this['content'].clientWidth});
    }
    componentWillReceiveProps() {
        // console.log(this['content'].clientWidth);
        // this.setState({inputOffset: this['content'].clientWidth});
    }
    focusTag() {
        this['input-inside'].focusInput();
    }
    focusInput() {
        this['input'].focusInput();
    }
    resetInput() {
        this['input-inside'].resetInput();
    }
    setInputValue(value, callback) {
        this['input'].setInputValue(value, callback);
    }
    getInputValue() {
        return this['input'].getInputValue();
    }
    addTagByInputValue() {
        return this['input'].addTagByInputValue();
    }
    addTagByEditInputValue() {
        if (!this['input-inside'])
            return;
        return this['input-inside'].addTagByInputValue();
    }
    setInfo(info, callback) {
        return this['input'].setInfo(info, callback);
    }
    moveToEnd() {
        return this['input'].moveToEnd();
    }
    // getInputAttr = (): AttributeValue => {
    //   return this['input'].getInputAttr();
    // }
    getInfo() {
        let { attr, values } = this.props;
        return { attr, values };
    }
    edit(pos) {
        this.setState({ inEditing: true });
        this['input-inside'].setInfo(this.getInfo(), () => {
            if (pos === 'attr') {
                return this['input-inside'] && this['input-inside'].selectAttr();
            }
            else {
                return this['input-inside'] && this['input-inside'].selectValue();
            }
        });
    }
    editDone() {
        this.setState({ inEditing: false });
    }
    render() {
        const { bubbleActice, inEditing, inputOffset } = this.state;
        const { active, attr, values, elect, dispatchTagEvent, attributes, focused, maxWidth } = this.props;
        let attrStr = attr ? attr.name : '';
        if (attr && attr.name && attr.type !== 'onlyKey') {
            attrStr += ': ';
        }
        let valueStr = values.map(item => item.name).join(' | ');
        const itemStyle = (inEditing && !active) ? { width: 0, minHeight: '20px' } : { minHeight: '20px' };
        const removeable = (attr && 'removeable' in attr) ? attr.removeable : true;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { style: itemStyle },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-bubble tc-15-bubble-bottom black", style: { display: bubbleActice ? '' : 'none' } },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-bubble-inner" }, "\u70B9\u51FB\u8FDB\u884C\u4FEE\u6539\uFF0C\u6309\u56DE\u8F66\u952E\u5B8C\u6210\u4FEE\u6539")),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: classnames__WEBPACK_IMPORTED_MODULE_1___default()("tc-tags", { "current": elect }), ref: div => this['content'] = div, style: { display: inEditing ? 'none' : '', paddingRight: removeable ? '' : '8px', maxWidth: 'none' }, onClick: this.handleTagClick, onMouseOver: () => this.setState({ bubbleActice: true }), onMouseOut: () => this.setState({ bubbleActice: false }) },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { onClick: (e) => this.handleTagClick(e, 'attr') }, attrStr),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { onClick: (e) => this.handleTagClick(e, 'value') }, valueStr),
                removeable &&
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { href: "javascript:;", className: "tc-tags-close-btn", onClick: this.handleDelete },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("i", { className: "clear-icon" }))),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Input__WEBPACK_IMPORTED_MODULE_2__["Input"], { type: "edit", hidden: !inEditing, maxWidth: maxWidth, handleKeyDown: this.handleKeyDown, active: active, ref: input => this["input-inside"] = input, attributes: attributes, dispatchTagEvent: dispatchTagEvent, isFocused: focused === _TagSearchBox__WEBPACK_IMPORTED_MODULE_3__["FocusPosType"].INPUT_EDIT })));
    }
}


/***/ }),

/***/ "./src/tea-components/tagsearchbox/TagSearchBox.tsx":
/*!**********************************************************!*\
  !*** ./src/tea-components/tagsearchbox/TagSearchBox.tsx ***!
  \**********************************************************/
/*! exports provided: FocusPosType, TagSearchBox */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "FocusPosType", function() { return FocusPosType; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "TagSearchBox", function() { return TagSearchBox; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! classnames */ "./node_modules/classnames/index.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(classnames__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! tea-components/libs/decorators/OnOuterClick */ "./src/tea-components/libs/decorators/OnOuterClick.ts");
/* harmony import */ var _Tag__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./Tag */ "./src/tea-components/tagsearchbox/Tag.tsx");
/* harmony import */ var _Input__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./Input */ "./src/tea-components/tagsearchbox/Input.tsx");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};





const keys = {
    "8": 'backspace',
    "9": 'tab',
    "13": 'enter',
    "32": 'spacebar',
    "37": 'left',
    "38": 'up',
    "39": 'right',
    "40": 'down'
};
const INPUT_MIN_SIZE = 5;
/**
 * 焦点所在位置类型
 */
var FocusPosType;
(function (FocusPosType) {
    FocusPosType[FocusPosType["INPUT"] = 0] = "INPUT";
    FocusPosType[FocusPosType["INPUT_EDIT"] = 1] = "INPUT_EDIT";
    FocusPosType[FocusPosType["TAG"] = 2] = "TAG";
})(FocusPosType || (FocusPosType = {}));
class TagSearchBox extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor() {
        super(...arguments);
        this.state = {
            active: false,
            dialogActive: false,
            curPos: 0,
            curPosType: FocusPosType.INPUT,
            showSelect: true,
            inputValue: '',
            inputWidth: INPUT_MIN_SIZE,
            tags: this.props.defaultValue ? this.props.defaultValue.map(item => {
                item._key = TagSearchBox.cnt++;
                return item;
            }) : []
        };
        this.open = () => {
            this.markTagElect(-1);
            const { active, tags } = this.state;
            if (!active) {
                this.setState({ active: true });
                // 展开时不激活select显示
                this.markTagElect(-1);
                this.setState({ curPosType: FocusPosType.INPUT, curPos: tags.length });
            }
            else {
                this.handleTagEvent('click-input', tags.length);
            }
            this.setState({ showSelect: true });
            // if (this[`tag-${tags.length}`].getInputValue().length > 0) {
            //   this.setState({ showSelect: true });
            // }
            setTimeout(() => {
                this[`tag-${tags.length}`] && this[`tag-${tags.length}`].moveToEnd();
            }, 100);
        };
        this.notify = (tags) => {
            const onChange = this.props.onChange;
            if (!onChange)
                return;
            const result = new Array();
            tags.forEach(item => {
                const attr = item.attr || null;
                const values = item.values;
                if (attr && attr.type === 'onlyKey' || values.length > 0) {
                    result.push({ attr, values, _key: item._key, _edit: item._edit });
                }
            });
            onChange(result);
        };
        /**
         * 点击清除按钮触发事件
         */
        this.handleClean = (e) => {
            e.stopPropagation();
            if (this.state.tags.length <= 0) {
                this[`tag-${0}`].setInputValue('');
                return;
            }
            this.setTags([], () => {
                this[`tag-${0}`].setInputValue('');
                this[`tag-${0}`].focusInput();
                // this.handleTagEvent('click-input', 0);
            });
            this.setState({ curPos: 0, curPosType: FocusPosType.INPUT });
        };
        /**
         * 点击帮助触发事件
         */
        this.handleHelp = (e) => {
            e.stopPropagation();
            this.setState({ dialogActive: true });
        };
        /**
         * 点击搜索触发事件
         */
        this.handleSearch = (e) => {
            if (!this.state.active)
                return;
            e.stopPropagation();
            const { curPos, curPosType, tags } = this.state;
            let flag = false;
            // if (curPosType === FocusPosType.INPUT && this[`tag-${curPos}`].addTagByInputValue()) flag = true;
            const input = this[`tag-${tags.length}`];
            if (input && input.addTagByInputValue) {
                if (input.addTagByInputValue()) {
                    flag = true;
                }
            }
            // if (curPosType === FocusPosType.INPUT_EDIT && this[`tag-${curPos}`].addTagByEditInputValue()) flag = true;
            for (let i = 0; i < tags.length; ++i) {
                if (!this[`tag-${i}`] || !this[`tag-${i}`].addTagByEditInputValue)
                    return;
                if (tags[i]._edit && this[`tag-${i}`].addTagByEditInputValue())
                    flag = true;
            }
            if (flag)
                return;
            this.notify(this.state.tags);
            input.focusInput();
        };
        /**
         *  处理Tag相关事件
         */
        this.handleTagEvent = (type, index, payload) => {
            const { tags, active, curPos, curPosType } = this.state;
            // console.log(type, index, payload);
            switch (type) {
                case 'add':
                    this.markTagElect(-1);
                    payload._key = TagSearchBox.cnt++;
                    tags.splice(++index, 0, payload);
                    this.setTags(tags, () => {
                        if (this[`tag-${index}`]) {
                            this[`tag-${index}`].focusInput();
                        }
                    });
                    this.setState({ showSelect: false });
                    break;
                case 'edit':
                    this.markTagElect(-1);
                    this[`tag-${index}`] && this[`tag-${index}`].editDone();
                    tags[index].attr = payload.attr;
                    tags[index].values = payload.values;
                    tags[index]._edit = false;
                    this.setTags(tags, () => {
                        // this[`tag-${index-1}`].focusInput();
                    });
                    index++;
                    this.setState({ showSelect: false, curPosType: FocusPosType.INPUT });
                    break;
                case 'edit-cancel':
                    this.markTagElect(-1);
                    this[`tag-${index}`].editDone();
                    this.setTags(tags, () => {
                        // this[`tag-${index}`].focusInput();
                    }, false);
                    this.setState({ showSelect: false, curPosType: FocusPosType.INPUT });
                    break;
                case 'editing':
                    if (('attr' in payload) && tags[index])
                        tags[index].attr = payload.attr;
                    if (('values' in payload) && tags[index])
                        tags[index].values = payload.values;
                    this.setTags(tags, null, false);
                    break;
                case 'mark':
                    if (index === tags.length)
                        index--;
                    if (index < 0 || !tags[index])
                        return;
                    if (!tags[index]._elect) {
                        this.markTagElect(index);
                        this.setState({ curPosType: FocusPosType.TAG });
                    }
                    break;
                case 'del':
                    if (payload === 'keyboard')
                        index--;
                    this.markTagElect(-1);
                    if (!tags[index])
                        break;
                    // 如果当前Tag的input中有内容，推向下个input
                    // const curValue = this[`tag-${index}`].getInputValue();
                    // this[`tag-${index}`].addTagByInputValue();
                    // 检查不可移除
                    const { attr } = tags[index];
                    if (attr && 'removeable' in attr && attr.removeable === false) {
                        break;
                    }
                    tags.splice(index, 1);
                    this.setTags(tags, () => {
                        this.setState({ curPosType: FocusPosType.INPUT });
                        // const input = this[`tag-${index-1 > 0 ? index-1 : 0}`];
                        // input.setInputValue(curValue);
                        // input.focusInput();
                    });
                    if (payload !== 'edit') {
                        this.setState({ showSelect: false });
                    }
                    break;
                // payload 为点击位置
                case 'click':
                    if (!active) {
                        this.open();
                        return;
                    }
                    // 触发修改
                    // if (curPos === index && curPosType === FocusPosType.TAG) {
                    const pos = payload;
                    tags[index]._edit = true;
                    this.setTags(tags, () => {
                        this.setState({ showSelect: true }, () => {
                            this[`tag-${index}`].edit(pos);
                        });
                    }, false);
                    this.setState({ curPosType: FocusPosType.INPUT_EDIT });
                    // } else {
                    //   this.markTagElect(index);
                    //   this.setState({ curPosType: FocusPosType.TAG });
                    // }
                    break;
                case 'click-input':
                    this.markTagElect(-1);
                    if (payload === 'edit') {
                        this.setState({ curPosType: FocusPosType.INPUT_EDIT });
                    }
                    else {
                        this.setState({ curPosType: FocusPosType.INPUT });
                    }
                    if (!active) {
                        this.setState({ active: true });
                    }
                    this.setState({ showSelect: true });
                    break;
                case 'move-left':
                    // if (index <= 0) return;
                    // if (index !== tags.length-1 || curPosType !== FocusPosType.INPUT) index--;
                    // this.markTagElect(index);
                    // this.setState({ curPosType: FocusPosType.TAG });
                    // this[`tag-${index}`].focusTag();
                    break;
                case 'move-right':
                    // if (index >= tags.length - 1) return;
                    // if (curPosType === FocusPosType.INPUT) {
                    //   this.setState({ curPosType: FocusPosType.TAG });
                    // }
                    // // 到达最后input
                    // if (index === tags.length-1) {
                    //   this.markTagElect(-1);
                    //   this[`tag-${index}`].focusInput();
                    //   this.setState({ curPosType: FocusPosType.INPUT });
                    // } else {
                    //   this.markTagElect(index);
                    //   this[`tag-${index}`].focusTag();
                    // }
                    // index++;
                    break;
            }
            this.setState({ curPos: index });
        };
    }
    componentDidMount() {
        if ('value' in this.props) {
            const value = this.props.value.map(item => {
                if (!('_key' in item)) {
                    item._key = TagSearchBox.cnt++;
                }
                return item;
            });
            this.setState({ tags: value });
        }
    }
    componentWillReceiveProps(nextProps) {
        if ('value' in nextProps) {
            const value = nextProps.value.map(item => {
                if (!('_key' in item)) {
                    item._key = TagSearchBox.cnt++;
                }
                return item;
            });
            this.setState({ tags: value });
        }
    }
    close() {
        // 编辑未完成的取消编辑
        const tags = this.state.tags.map((item, index) => {
            if (item._edit) {
                this[`tag-${index}`].editDone();
                item._edit = false;
            }
            return item;
        });
        this.setTags(tags, () => {
            this.markTagElect(-1);
            this.setState({ showSelect: false });
            if (this.state.active) {
                this.setState({ curPos: -1 }, () => this.setState({ active: false }, () => this[`search-box`].scrollLeft = 0));
            }
        }, false);
    }
    // Tags发生变动
    setTags(tags, callback, notify = true) {
        if (notify)
            this.notify(tags);
        this.setState({ tags }, () => { if (callback)
            callback(); });
    }
    markTagElect(index) {
        const tags = this.state.tags.map((item, i) => {
            if (index === i) {
                item._elect = true;
                this[`tag-${index}`].focusTag();
            }
            else {
                item._elect = false;
            }
            return item;
        });
        this.setState({ tags });
    }
    render() {
        const { active, inputWidth, inputValue, tags, curPos, curPosType, dialogActive, showSelect } = this.state;
        const { minWidth, tipZh, tipEn, attributes } = this.props;
        // 用于计算 focused 及 isFocused, 判断是否显示选择组件
        // (直接使用 Input 组件内部 onBlur 判断会使得 click 时组件消失)
        let focusedInputIndex = -1;
        if (curPosType === FocusPosType.INPUT || curPosType === FocusPosType.INPUT_EDIT) {
            focusedInputIndex = curPos;
        }
        const tagList = tags.map((item, index) => {
            // 补全 attr 属性
            attributes.forEach(attrItem => {
                if (item.attr && attrItem.key && attrItem.key == item.attr.key) {
                    item.attr = Object.assign({}, item.attr, attrItem);
                }
            });
            const selectedAttrKeys = [];
            tags.forEach(tag => {
                if (tag.attr && item.attr && item._edit && item.attr.key === tag.attr.key)
                    return null;
                if (tag.attr && tag.attr.key && !tag.attr.reusable) {
                    selectedAttrKeys.push(tag.attr.key);
                }
            });
            const useableAttributes = attributes.filter(item => selectedAttrKeys.indexOf(item.key) < 0);
            return react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Tag__WEBPACK_IMPORTED_MODULE_3__["Tag"], { ref: tag => this[`tag-${index}`] = tag, active: active, key: item._key, attributes: useableAttributes, attr: item.attr, values: item.values, elect: item._elect, maxWidth: this['search-wrap'] ? this['search-wrap'].clientWidth : null, focused: (focusedInputIndex === index && showSelect) ? curPosType : null, dispatchTagEvent: (type, payload) => this.props.editable && this.handleTagEvent(type, index, payload) });
        });
        const selectedAttrKeys = tags.map(item => item.attr && !item.attr.reusable ? item.attr.key : null).filter(item => !!item);
        const useableAttributes = attributes.filter(item => selectedAttrKeys.indexOf(item.key) < 0);
        const minWidthStyle = active ? {} : ({ width: minWidth ? minWidth : '100%' });
        if (this.props.editable) {
            tagList.push(react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { key: '100' },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Input__WEBPACK_IMPORTED_MODULE_4__["Input"], { ref: input => this[`tag-${tags.length}`] = input, active: active, maxWidth: this['search-wrap'] ? this['search-wrap'].clientWidth : null, attributes: useableAttributes, isFocused: focusedInputIndex === tags.length && showSelect, dispatchTagEvent: (type, payload) => this.handleTagEvent(type, tags.length, payload) })));
        }
        const tip = window['VERSION'] === 'en' ? tipEn : tipZh;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-select-tags-search-wrap", style: this.props.style, ref: div => this['search-wrap'] = div },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: classnames__WEBPACK_IMPORTED_MODULE_1___default()("tc-select-tags-search", { "focus": this.props.editable && active }), onClick: this.open, style: minWidthStyle, ref: div => this[`search-box`] = div },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-search-wrap" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("ul", null, tagList),
                    this.props.editable && tagList.length === 1 && react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("span", { className: "help-tips", style: { lineHeight: '28px' } },
                        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("p", { className: "text text-weak" }, tip))))));
    }
}
TagSearchBox.cnt = 0;
TagSearchBox.defaultProps = {
    tipZh: "多个关键字用竖线“ | ”分隔，多个过滤标签用回车键分隔",
    tipEn: "Use '|' to split more than one keyword, and press Enter to split tags",
    minWidth: 210,
    editable: true,
};
__decorate([
    tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_2__["OnOuterClick"]
], TagSearchBox.prototype, "close", null);


/***/ }),

/***/ "./src/tea-components/tagsearchbox/index.ts":
/*!**************************************************!*\
  !*** ./src/tea-components/tagsearchbox/index.ts ***!
  \**************************************************/
/*! exports provided: TagSearchBox */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _TagSearchBox__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./TagSearchBox */ "./src/tea-components/tagsearchbox/TagSearchBox.tsx");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "TagSearchBox", function() { return _TagSearchBox__WEBPACK_IMPORTED_MODULE_0__["TagSearchBox"]; });




/***/ }),

/***/ "./src/tea-components/tagsearchbox/valueselect/Loading.tsx":
/*!*****************************************************************!*\
  !*** ./src/tea-components/tagsearchbox/valueselect/Loading.tsx ***!
  \*****************************************************************/
/*! exports provided: Loading */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Loading", function() { return Loading; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");

const Loading = (props) => {
    return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-autocomplete", style: { left: `${props.offset}px` } },
        react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("ul", { className: "tc-15-autocomplete-menu", role: "menu" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "autocomplete-empty", role: "menuitem", href: "javascript:;" }, "\u52A0\u8F7D\u4E2D ..")))));
};


/***/ }),

/***/ "./src/tea-components/tagsearchbox/valueselect/MultipleValueSelect.tsx":
/*!*****************************************************************************!*\
  !*** ./src/tea-components/tagsearchbox/valueselect/MultipleValueSelect.tsx ***!
  \*****************************************************************************/
/*! exports provided: MultipleValueSelect */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "MultipleValueSelect", function() { return MultipleValueSelect; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");

const keys = {
    "8": 'backspace',
    "9": 'tab',
    "13": 'enter',
    "37": 'left',
    "38": 'up',
    "39": 'right',
    "40": 'down'
};
class MultipleValueSelect extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        // 父组件调用
        this.handleKeyDown = (keyCode) => {
            if (!keys[keyCode])
                return;
            const { onSelect, onChange } = this.props;
            const { curIndex, select } = this.state;
            switch (keys[keyCode]) {
                case 'tab':
                    if (curIndex < 0)
                        return false;
                    const pos = select.indexOf(curIndex);
                    if (pos >= 0) {
                        select.splice(pos, 1);
                    }
                    else {
                        select.push(curIndex);
                    }
                    if (onChange) {
                        onChange(this.getValue(select));
                    }
                    return false;
                case 'enter':
                    if (onSelect) {
                        onSelect(this.getValue(select));
                    }
                    return false;
                case 'up':
                    this.move(-1);
                    break;
                case 'down':
                    this.move(1);
                    break;
            }
        };
        this.move = (step) => {
            const curIndex = this.state.curIndex;
            const { values, inputValue } = this.props;
            if (values.length <= 0)
                return;
            this.setState({ curIndex: (curIndex + step + values.length) % values.length });
        };
        this.handleClick = (e, index) => {
            e.stopPropagation();
            const select = this.state.select;
            const onChange = this.props.onChange;
            const pos = select.indexOf(index);
            if (pos >= 0) {
                select.splice(pos, 1);
            }
            else {
                select.push(index);
            }
            if (onChange) {
                onChange(this.getValue(select));
            }
        };
        this.handleSubmit = (e) => {
            e.stopPropagation();
            const onSelect = this.props.onSelect;
            const select = this.state.select;
            if (onSelect) {
                onSelect(this.getValue(select));
            }
        };
        this.handleCancel = (e) => {
            e.stopPropagation();
            const onCancel = this.props.onCancel;
            if (onCancel) {
                onCancel();
            }
        };
        const list = this.props.inputValue.split('|').map(i => i.trim());
        const select = [], values = this.props.values.map(item => Object.assign({}, item, { name: item.name.trim() }));
        values.forEach((item, index) => {
            if (list.indexOf(item.name) >= 0) {
                select.push(index);
            }
        });
        this.state = {
            curIndex: -1,
            select,
        };
    }
    componentDidMount() {
        const select = this.state.select;
        if (select.length <= 0 && this.props.onSelect) {
            this.props.onSelect(this.getValue(select));
        }
    }
    componentWillReceiveProps(nextProps) {
        if (this.props.inputValue !== nextProps.inputValue) {
            const list = nextProps.inputValue.split('|').map(i => i.trim());
            const select = [], values = nextProps.values.map(item => Object.assign({}, item, { name: item.name.trim() }));
            values.forEach((item, index) => {
                if (list.indexOf(item.name) >= 0) {
                    select.push(index);
                }
            });
            this.setState({ select });
        }
    }
    getValue(select) {
        const { values } = this.props;
        return select.map(i => values[i]);
    }
    render() {
        const { select, curIndex } = this.state;
        const { inputValue, values, offset } = this.props;
        const list = values.map((item, index) => {
            const input = react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("label", { className: "form-ctrl-label", style: item.style || {}, title: item.name },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("input", { type: "checkbox", readOnly: true, checked: select.indexOf(index) >= 0, className: "tc-15-checkbox" }),
                item.name);
            if (curIndex === index) {
                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation", key: index, className: "autocomplete-cur", onMouseDown: (e) => this.handleClick(e, index) }, input));
            }
            return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation", key: index, onMouseDown: (e) => this.handleClick(e, index) }, input));
        });
        if (list.length === 0)
            return null;
        let maxHeight = document.body.clientHeight ? document.body.clientHeight - 400 : 450;
        maxHeight = maxHeight > 240 ? maxHeight : 240;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-autocomplete", style: { left: `${offset}px` } },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("ul", { className: "tc-15-autocomplete-menu", role: "menu", style: { maxHeight: `${maxHeight}px` } }, list),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-autocomplete-ft" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { href: "javascript:;", className: "autocomplete-btn", onClick: this.handleSubmit }, "\u5B8C\u6210"),
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { href: "javascript:;", className: "autocomplete-btn", onClick: this.handleCancel }, "\u53D6\u6D88"))));
    }
}


/***/ }),

/***/ "./src/tea-components/tagsearchbox/valueselect/PureInput.tsx":
/*!*******************************************************************!*\
  !*** ./src/tea-components/tagsearchbox/valueselect/PureInput.tsx ***!
  \*******************************************************************/
/*! exports provided: PureInput */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "PureInput", function() { return PureInput; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");

const keys = {
    "9": 'tab',
    "13": 'enter'
};
class PureInput extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor() {
        super(...arguments);
        // 父组件调用
        this.handleKeyDown = (keyCode) => {
            if (!keys[keyCode])
                return;
            const { onSelect, inputValue } = this.props;
            switch (keys[keyCode]) {
                case 'tab':
                case 'enter':
                    if (inputValue.length <= 0)
                        return false;
                    if (onSelect) {
                        onSelect(this.getValue(this.props.inputValue).filter(i => !!i.name));
                    }
                    return false;
            }
        };
    }
    componentDidMount() {
        const onChange = this.props.onChange;
        if (onChange) {
            onChange(this.getValue(this.props.inputValue));
        }
    }
    componentWillReceiveProps(nextProps) {
        if (this.props.inputValue !== nextProps.inputValue) {
            const onChange = nextProps.onChange;
            if (onChange) {
                onChange(this.getValue(nextProps.inputValue));
            }
        }
    }
    getValue(value) {
        return value.split('|').map(item => { return { name: item.trim() }; });
    }
    render() {
        return null;
    }
}


/***/ }),

/***/ "./src/tea-components/tagsearchbox/valueselect/SingleValueSelect.tsx":
/*!***************************************************************************!*\
  !*** ./src/tea-components/tagsearchbox/valueselect/SingleValueSelect.tsx ***!
  \***************************************************************************/
/*! exports provided: SingleValueSelect */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "SingleValueSelect", function() { return SingleValueSelect; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");

const keys = {
    "8": 'backspace',
    "9": 'tab',
    "13": 'enter',
    "37": 'left',
    "38": 'up',
    "39": 'right',
    "40": 'down'
};
class SingleValueSelect extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        // 父组件调用
        this.handleKeyDown = (keyCode) => {
            if (!keys[keyCode])
                return;
            const { onSelect } = this.props;
            const select = this.state.select;
            switch (keys[keyCode]) {
                case 'enter':
                case 'tab':
                    if (onSelect) {
                        onSelect(this.getValue(select));
                    }
                    return false;
                case 'up':
                    this.move(-1);
                    break;
                case 'down':
                    this.move(1);
                    break;
            }
        };
        this.move = (step) => {
            const select = this.state.select;
            const { values, inputValue } = this.props;
            const list = values;
            if (list.length <= 0)
                return;
            this.setState({ select: (select + step + list.length) % list.length });
        };
        this.handleClick = (e, index) => {
            e.stopPropagation();
            if (this.props.onSelect) {
                this.props.onSelect(this.getValue(index));
            }
        };
        const { values, inputValue, onSelect } = this.props;
        let select = -1;
        values.forEach((item, index) => {
            if (item.name === inputValue) {
                select = index;
            }
        });
        this.state = {
            select,
        };
    }
    componentDidMount() {
        const select = this.state.select;
        if (select < 0 && this.props.onSelect) {
            this.props.onSelect(this.getValue(select));
        }
    }
    componentWillReceiveProps(nextProps) {
        const { values, inputValue } = nextProps;
        const list = values.map(item => item.name);
        const select = list.indexOf(inputValue);
        this.setState({ select });
    }
    getValue(select) {
        const { values, inputValue } = this.props;
        if (select < 0) {
            // return [];
            if (inputValue)
                return [{ name: inputValue }];
            return [];
        }
        const list = values;
        if (select < list.length) {
            return [list[select]];
        }
        else {
            const select = list.map(item => item.name).indexOf(inputValue);
            this.setState({ select });
            if (select < 0)
                return [];
            return [list[select]];
        }
    }
    render() {
        const select = this.state.select;
        const { inputValue, values, offset } = this.props;
        const list = values.map((item, index) => {
            if (select === index) {
                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation", key: index, className: "autocomplete-cur", onClick: (e) => this.handleClick(e, index) },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "text-truncate", role: "menuitem", href: "javascript:;", title: item.name, style: item.style || {} }, item.name)));
            }
            return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation", key: index, onClick: (e) => this.handleClick(e, index) },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "text-truncate", role: "menuitem", href: "javascript:;", title: item.name, style: item.style || {} }, item.name)));
        });
        if (list.length === 0) {
            list.push(react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { role: "presentation", key: 0 },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("a", { className: "autocomplete-empty", role: "menuitem", href: "javascript:;" }, "\u76F8\u5173\u503C\u4E0D\u5B58\u5728")));
        }
        let maxHeight = document.body.clientHeight ? document.body.clientHeight - 450 : 400;
        maxHeight = maxHeight > 240 ? maxHeight : 240;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-15-autocomplete", style: { left: `${offset}px`, width: 'auto', minWidth: 180 } },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("ul", { className: "tc-15-autocomplete-menu", role: "menu", style: { maxHeight: `${maxHeight}px` } }, list)));
    }
}


/***/ }),

/***/ "./src/tea-components/tagsearchbox/valueselect/index.tsx":
/*!***************************************************************!*\
  !*** ./src/tea-components/tagsearchbox/valueselect/index.tsx ***!
  \***************************************************************/
/*! exports provided: ValueSelect */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ValueSelect", function() { return ValueSelect; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var _PureInput__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./PureInput */ "./src/tea-components/tagsearchbox/valueselect/PureInput.tsx");
/* harmony import */ var _SingleValueSelect__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./SingleValueSelect */ "./src/tea-components/tagsearchbox/valueselect/SingleValueSelect.tsx");
/* harmony import */ var _MultipleValueSelect__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./MultipleValueSelect */ "./src/tea-components/tagsearchbox/valueselect/MultipleValueSelect.tsx");
/* harmony import */ var _Loading__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./Loading */ "./src/tea-components/tagsearchbox/valueselect/Loading.tsx");





class ValueSelect extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        this.mount = false;
        this.handleKeyDown = (keyCode) => {
            if (this['select'] && this['select'].handleKeyDown) {
                return this['select'].handleKeyDown(keyCode);
            }
            return true;
        };
        let values = [];
        const propsValues = this.props.values;
        if (typeof propsValues !== "function") {
            values = propsValues;
        }
        this.state = { values };
    }
    componentDidMount() {
        this.mount = true;
        const propsValues = this.props.values;
        if (typeof propsValues === "function") {
            const res = propsValues();
            // Promise
            if (res && res.then) {
                res.then(values => {
                    this.mount && this.setState({ values });
                });
            }
            else {
                this.mount && this.setState({ values: res });
            }
        }
    }
    componentWillUnmount() {
        this.mount = false;
    }
    // handleKeyUp = (keyCode: number): boolean => {
    //   if (this['select'] && this['select'].handleKeyUp) {
    //     return this['select'].handleKeyUp(keyCode);
    //   }
    //   return true;
    // }
    render() {
        const values = this.state.values;
        const { type, inputValue, onChange, onSelect, onCancel, offset } = this.props;
        const props = { values, inputValue, onChange, onSelect, onCancel, offset };
        switch (type) {
            case 'input':
                return react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_PureInput__WEBPACK_IMPORTED_MODULE_1__["PureInput"], Object.assign({ ref: select => this['select'] = select }, props));
            case 'single':
                if (values.length <= 0) {
                    return react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Loading__WEBPACK_IMPORTED_MODULE_4__["Loading"], { offset: offset });
                }
                return react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_SingleValueSelect__WEBPACK_IMPORTED_MODULE_2__["SingleValueSelect"], Object.assign({ ref: select => this['select'] = select }, props));
            case 'multiple':
                if (values.length <= 0) {
                    return react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Loading__WEBPACK_IMPORTED_MODULE_4__["Loading"], null);
                }
                return react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_MultipleValueSelect__WEBPACK_IMPORTED_MODULE_3__["MultipleValueSelect"], Object.assign({ ref: select => this['select'] = select }, props));
        }
        return null;
    }
}


/***/ }),

/***/ "./src/tea-components/timepicker/Select.tsx":
/*!**************************************************!*\
  !*** ./src/tea-components/timepicker/Select.tsx ***!
  \**************************************************/
/*! exports provided: Select */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Select", function() { return Select; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var react_dom__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! react-dom */ "./webpack/alias/react-dom.js");


/**
 * （受控组件）
 */
class Select extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor() {
        super(...arguments);
        /**
         * 滚动到指定元素
         */
        this.scrollTo = (element, to, duration) => {
            if (duration <= 0) {
                element.scrollTop = to;
                return;
            }
            var difference = to - element.scrollTop;
            var perTick = difference / duration * 10;
            setTimeout(() => {
                element.scrollTop = element.scrollTop + perTick;
                if (element.scrollTop === to)
                    return;
                this.scrollTo(element, to, duration - 10);
            }, 10);
        };
        /**
         * 滚动到当前选择元素
         */
        this.scrollToSelected = (duration) => {
            const select = react_dom__WEBPACK_IMPORTED_MODULE_1__["findDOMNode"](this);
            const list = react_dom__WEBPACK_IMPORTED_MODULE_1__["findDOMNode"](this.refs['list']);
            const index = this.props.value;
            const topOption = list.children[index];
            const to = topOption.offsetTop - select.offsetTop;
            this.scrollTo(select, to, duration);
        };
        /**
         * 根据范围生成列表
         */
        this.genRangeList = (start, end) => Array(end - start + 1).fill(0).map((e, i) => {
            const num = i + start;
            return num > 9 ? `${num}` : `0${num}`;
        });
        this.handleSelect = (e, val) => {
            e.stopPropagation();
            if (this.props.onChange)
                this.props.onChange(val);
        };
    }
    componentDidMount() {
        this.scrollToSelected(0);
    }
    componentDidUpdate() {
        this.scrollToSelected(150);
    }
    render() {
        const { from, to, value, range } = this.props;
        const list = this.genRangeList(from, to).map((item, i) => {
            if (range && 'min' in range && i < range.min)
                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { key: i, className: "disabled" }, item));
            if (range && 'max' in range && i > range.max)
                return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { key: i, className: "disabled" }, item));
            return react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("li", { key: i, className: +item == value ? 'current' : '', onClick: (e) => this.handleSelect(e, +item) }, item);
        });
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-time-picker-select" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("ul", { ref: "list" }, list)));
    }
}


/***/ }),

/***/ "./src/tea-components/timepicker/TimePicker.tsx":
/*!******************************************************!*\
  !*** ./src/tea-components/timepicker/TimePicker.tsx ***!
  \******************************************************/
/*! exports provided: TimePicker */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "TimePicker", function() { return TimePicker; });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./webpack/alias/react.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! classnames */ "./node_modules/classnames/index.js");
/* harmony import */ var classnames__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(classnames__WEBPACK_IMPORTED_MODULE_1__);
/* harmony import */ var _bubble__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../bubble */ "./src/tea-components/bubble/index.ts");
/* harmony import */ var tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! tea-components/libs/decorators/OnOuterClick */ "./src/tea-components/libs/decorators/OnOuterClick.ts");
/* harmony import */ var _Select__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./Select */ "./src/tea-components/timepicker/Select.tsx");
/* harmony import */ var react_dom__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! react-dom */ "./webpack/alias/react-dom.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};






// disabledHours?: () => Array<number>;
// disabledMinutes?: (selectedHour: number) => Array<number>;
// disabledSeconds?: (selectedHour: number, disabledMinutes: number) => Array<number>;
const keys = {
    "38": "up",
    "40": "down",
    "13": "enter"
};
class TimePicker extends react__WEBPACK_IMPORTED_MODULE_0__["Component"] {
    constructor(props) {
        super(props);
        // getCurTime = (props: TimePickerProps) => {
        // }
        /**
         * 格式化
         */
        this.formatNum = (num) => (num > 9 ? `${num}` : `0${num}`);
        this.format = (value) => {
            if (typeof value === "string")
                return value;
            if (!value)
                return "";
            const { hour, minute, second } = value;
            return `${this.formatNum(hour)}:${this.formatNum(minute)}:${this.formatNum(second)}`;
        };
        /**
         * 检验传入value值是否合法
         */
        this.check = (value) => {
            if (!value)
                return false;
            if (!("hour" in value) || value.hour < 0 || value.hour > 23)
                return false;
            if (!("minute" in value) || value.minute < 0 || value.minute > 59)
                return false;
            if (!("second" in value) || value.second < 0 || value.second > 59)
                return false;
            return true;
        };
        /**
         * 将字符串解析为 TimePickerValue
         */
        this.parse = (value) => {
            if (typeof value !== "string")
                return value;
            if (!/^[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(value))
                return null;
            return {
                hour: +value.substr(0, 2),
                minute: +value.substr(3, 2),
                second: +value.substr(6, 2)
            };
        };
        /**
         * 确认value，合法将调用 onChange
         */
        this.confirm = value => {
            if (!("value" in this.props)) {
                this.setState({ value });
            }
            if (this.props.onChange) {
                this.props.onChange(this.format(value));
            }
            this.setState({ inputValue: this.format(value) });
        };
        this.handleInputChange = (e) => {
            const inputValue = e.target.value;
            if (/^[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(inputValue)) {
                const value = {
                    hour: +inputValue.substr(0, 2),
                    minute: +inputValue.substr(3, 2),
                    second: +inputValue.substr(6, 2)
                };
                if (this.check(value)) {
                    const range = this.props.range || {};
                    if ((!("min" in range) || this.format(value) >= range.min) &&
                        (!("max" in range) || this.format(value) <= range.max)) {
                        this.confirm(value);
                    }
                }
                this.setState({ bubbleActive: false });
            }
            else {
                this.setState({ bubbleActive: true });
            }
            if (inputValue.length === 0) {
                this.setState({ bubbleActive: false });
            }
            this.setState({ inputValue });
        };
        this.open = () => {
            if ("value" in this.props) {
                this.setState({ value: this.parse(this.props.value) });
            }
            if (!this.state.active && !this.props.disabled) {
                this.setState({ active: true });
            }
        };
        this.handleSelect = (type, val) => {
            let value = this.state.value || { hour: 0, minute: 0, second: 0 };
            switch (type) {
                case "hour":
                    value.hour = val;
                    break;
                case "minute":
                    value.minute = val;
                    break;
                case "second":
                    value.second = val;
            }
            // 根据范围自调
            const range = this.props.range || {};
            let times = 0;
            while (("min" in range && this.format(value) < range.min) || ("max" in range && this.format(value) > range.max)) {
                if (type === "hour")
                    value.minute = (value.minute + 1) % 60;
                if (type === "minute")
                    value.second = (value.second + 1) % 60;
                if (++times > 60)
                    break;
            }
            this.confirm(value);
        };
        /**
         * 处理键位按键事件
         */
        this.handleKeyDown = (e, hourRange, minuteRange, secondRange) => {
            if (!keys[e.keyCode])
                return;
            e.preventDefault();
            const input = e.currentTarget;
            let start = input.selectionStart, end = input.selectionEnd;
            setTimeout(() => input.setSelectionRange(start, end));
            let value = this.state.value || { hour: 0, minute: 0, second: 0 };
            let type;
            switch (keys[e.keyCode]) {
                case "up":
                    if (start >= 0 && end <= 2) {
                        type = "hour";
                        value.hour = (value.hour + 1) % 24;
                        if (("min" in hourRange && value.hour < hourRange.min) ||
                            ("max" in hourRange && value.hour > hourRange.max)) {
                            value.hour = hourRange.min || 0;
                        }
                    }
                    if (start >= 3 && end <= 5) {
                        type = "minute";
                        value.minute = (value.minute + 1) % 60;
                        if (("min" in minuteRange && value.minute < minuteRange.min) ||
                            ("max" in minuteRange && value.minute > minuteRange.max)) {
                            value.minute = minuteRange.min || 0;
                        }
                    }
                    if (start >= 6 && end <= 8) {
                        type = "second";
                        value.second = (value.second + 1) % 60;
                        if (("min" in secondRange && value.second < secondRange.min) ||
                            ("max" in secondRange && value.second > secondRange.max)) {
                            value.second = secondRange.min || 0;
                        }
                    }
                    break;
                case "down":
                    if (start >= 0 && end <= 2) {
                        type = "hour";
                        value.hour = (value.hour - 1 + 24) % 24;
                        if (("min" in hourRange && value.hour < hourRange.min) ||
                            ("max" in hourRange && value.hour > hourRange.max)) {
                            value.hour = hourRange.max || 24;
                        }
                    }
                    if (start >= 3 && end <= 5) {
                        type = "minute";
                        value.minute = (value.minute - 1 + 60) % 60;
                        if (("min" in minuteRange && value.minute < minuteRange.min) ||
                            ("max" in minuteRange && value.minute > minuteRange.max)) {
                            value.minute = minuteRange.max || 60;
                        }
                    }
                    if (start >= 6 && end <= 8) {
                        type = "second";
                        value.second = (value.second - 1 + 60) % 60;
                        if (("min" in secondRange && value.second < secondRange.min) ||
                            ("max" in secondRange && value.second > secondRange.max)) {
                            value.second = secondRange.max || 60;
                        }
                    }
                    break;
                case "enter":
                    this.close();
                    react_dom__WEBPACK_IMPORTED_MODULE_5__["findDOMNode"](this["input"]).blur();
                    return;
            }
            const range = this.props.range || {};
            while (("min" in range && this.format(value) < range.min) || ("max" in range && this.format(value) > range.max)) {
                // 根据范围修正
                if (type === "hour")
                    value.minute = (value.minute + 1) % 60;
                if (type === "minute")
                    value.second = (value.second + 1) % 60;
            }
            this.confirm(value);
        };
        let defaultValue = props.value || props.defaultValue;
        const value = this.check(this.parse(defaultValue)) ? this.parse(defaultValue) : null;
        this.state = {
            bubbleActive: false,
            active: false,
            inputValue: defaultValue ? this.format(defaultValue) : "",
            value
        };
        this.close = this.close.bind(this);
    }
    componentWillReceiveProps(nextProps) {
        // 受控组件
        if ("value" in nextProps) {
            const nextValue = this.parse(nextProps.value);
            let value = this.check(nextValue) ? nextValue : null;
            const range = nextProps.range || {};
            if (value !== null && range.min && this.format(value) < range.min) {
                value = this.parse(range.min);
                this.confirm(range.min);
                return;
            }
            if (value !== null && range.max && this.format(value) > range.max) {
                value = this.parse(range.max);
                this.confirm(range.max);
                return;
            }
            this.setState({ value, inputValue: this.format(value) });
        }
    }
    close() {
        const { inputValue } = this.state;
        let value = this.state.value;
        if ("value" in this.props) {
            value = this.props.value;
            this.setState({ inputValue: this.format(value) });
        }
        // 非法输入检测
        if (!/^[0-9]{2}:[0-9]{2}:[0-9]{2}$/.test(inputValue)) {
            this.setState({ inputValue: this.format(value) });
        }
        const curValue = {
            hour: +inputValue.substr(0, 2),
            minute: +inputValue.substr(3, 2),
            second: +inputValue.substr(6, 2)
        };
        if (!this.check(curValue)) {
            this.setState({ inputValue: this.format(value) });
        }
        this.setState({ active: false, bubbleActive: false });
    }
    render() {
        const value = this.state.value || { hour: 0, minute: 0, second: 0 };
        const { hour, minute, second } = value;
        const range = this.props.range || {};
        const minRange = this.parse(range.min) || { hour: 0, minute: 0, second: 0 };
        const maxRange = this.parse(range.max) || { hour: 23, minute: 59, second: 59 };
        let hourRange = { min: minRange.hour, max: maxRange.hour };
        let minuteRange = {}, secondRange = {};
        if (hour == hourRange.min) {
            minuteRange.min = minRange.minute;
            if ("min" in minuteRange && minute == minuteRange.min) {
                secondRange.min = minRange.second;
            }
        }
        if (hour == hourRange.max) {
            minuteRange.max = maxRange.minute;
            if ("max" in minuteRange && minute == minuteRange.max) {
                secondRange.max = maxRange.second;
            }
        }
        const combobox = this.state.active ? (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-time-picker-combobox" },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Select__WEBPACK_IMPORTED_MODULE_4__["Select"], { from: 0, to: 23, value: hour, range: hourRange, onChange: value => {
                    this.handleSelect("hour", value);
                } }),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Select__WEBPACK_IMPORTED_MODULE_4__["Select"], { from: 0, to: 59, value: minute, range: minuteRange, onChange: value => {
                    this.handleSelect("minute", value);
                } }),
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_Select__WEBPACK_IMPORTED_MODULE_4__["Select"], { from: 0, to: 59, value: second, range: secondRange, onChange: value => {
                    this.handleSelect("second", value);
                } }))) : null;
        return (react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: classnames__WEBPACK_IMPORTED_MODULE_1___default()("tc-time-picker", { active: this.state.active }) },
            react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("div", { className: "tc-time-picker-input-wrap" },
                react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_bubble__WEBPACK_IMPORTED_MODULE_2__["BubbleWrapper"], { align: "start" },
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"]("input", { type: "text", ref: el => (this["input"] = el), className: "tc-15-input-text shortest", onClick: this.open, onFocus: this.open, placeholder: this.props.placeHolder || "时间选择", value: this.state.inputValue, onChange: this.handleInputChange, disabled: this.props.disabled, onKeyDown: e => this.handleKeyDown(e, hourRange, minuteRange, secondRange), maxLength: 8 }),
                    react__WEBPACK_IMPORTED_MODULE_0__["createElement"](_bubble__WEBPACK_IMPORTED_MODULE_2__["Bubble"], { focusPosition: "bottom", className: "error", style: { display: this.state.bubbleActive ? "" : "none" } }, "\u683C\u5F0F\u9519\u8BEF\uFF0C\u5E94\u4E3AHH:mm:ss"))),
            combobox));
    }
}
__decorate([
    tea_components_libs_decorators_OnOuterClick__WEBPACK_IMPORTED_MODULE_3__["OnOuterClick"]
], TimePicker.prototype, "close", null);


/***/ }),

/***/ "./src/tea-components/timepicker/index.ts":
/*!************************************************!*\
  !*** ./src/tea-components/timepicker/index.ts ***!
  \************************************************/
/*! exports provided: TimePicker */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _TimePicker__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./TimePicker */ "./src/tea-components/timepicker/TimePicker.tsx");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "TimePicker", function() { return _TimePicker__WEBPACK_IMPORTED_MODULE_0__["TimePicker"]; });




/***/ }),

/***/ "./webpack/alias/react-dom.js":
/*!************************************!*\
  !*** ./webpack/alias/react-dom.js ***!
  \************************************/
/*! no static exports found */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _react_dom_global__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! __react-dom-global */ "__react-dom-global");
/* harmony import */ var _react_dom_global__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(_react_dom_global__WEBPACK_IMPORTED_MODULE_0__);
/* harmony reexport (unknown) */ for(var __WEBPACK_IMPORT_KEY__ in _react_dom_global__WEBPACK_IMPORTED_MODULE_0__) if(__WEBPACK_IMPORT_KEY__ !== 'default') (function(key) { __webpack_require__.d(__webpack_exports__, key, function() { return _react_dom_global__WEBPACK_IMPORTED_MODULE_0__[key]; }) }(__WEBPACK_IMPORT_KEY__));



/* harmony default export */ __webpack_exports__["default"] = (_react_dom_global__WEBPACK_IMPORTED_MODULE_0__);

/***/ })

}]);
//# sourceMappingURL=ChartsComponents.js.map