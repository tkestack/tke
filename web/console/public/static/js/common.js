define('models/manager', function (require, exports, module) {
    var util = require('util');
    var appUtil = require('appUtil');
    var tips = require('tips');
    var router = require('router');

    //数据获取
    var manager = {
        /**
         * 获取用户所在区域（国内/海外）
         */
        getUserAreaInfo: function (cb, fail) {
            this.getComData(function (comData) {
                cb(comData.userInfo.areaInfo);
            });
        },

        /**
         * 查询公共流程(返回status-- 0:成功 1:失败 2:未完成)
         */
        getFlowCscResult: function (data, cb, fail, from) {
            var _self = this,
                _cb = function (ret) {

                    var retData = ret.data || {},
                        retOutPut = retData.output;

                    if (ret.code == 0) { //需要先解绑再执行操作的提示
                        if (retOutPut && retOutPut['failedInfo']) {
                            for (var i = 0, item; item = retOutPut['failedInfo'][i]; i++) {
                                var errorNum = item.errorNum;
                                errorNum = typeof (errorNum) == 'object' ? errorNum[0] : errorNum;
                                errorNum = errorNum.toString().substring(1);
                                var errTipTxt = constants['FLOW_RESULT_FAIL'][errorNum] || '操作执行失败，请稍后重试';
                                tips.error(errTipTxt);
                                break;
                            }
                        } else if (retOutPut && retOutPut.errorCode) {
                            var errTipTxt = constants['FLOW_RESULT_FAIL'][retOutPut.errorCode] || '操作执行失败，请稍后重试';
                            tips.error(errTipTxt);
                        }
                    }

                };

            net.send(config.get('getFlowCscResult'), {
                data: data,
                global: from === 'vpc',
                cb: _cb
            });
        }
    }
    module.exports = manager;

});

define('widget/i18n', function (require, exports, module) {
    var Bee = require('qccomponent');
    var manager = require('manager');

    /**
     * 快捷方式
     */
    exports.isI18n = (function () {
        var cached;
        return function () {
            if (!cached) {
                cached = new Promise(function (resolve, reject) {
                    manager.getUserAreaInfo(function (areaInfo) {
                        resolve(areaInfo && +areaInfo.area === 2)
                    }, reject)
                })
            }
            return cached;
        }
    })()
    /**
     * 快捷方式
     */
    exports.then = function (onFulfilled, onRejected) {
        return exports.isI18n().then(onFulfilled, onRejected)
    }
    /**
     * 取反（用于b-if判断）
     */
    exports.notI18n = function () {
        return exports.then(function (isI18n) {
            return !isI18n
        })
    }
    /**
     * 同步化 then，避免Promise.prototype.then的异步特性导致多余的渲染
     */
    exports.syncThen = (function () {
        var fulfilled, fulfilledValue;
        exports.then(function (i18n) {
            fulfilled = true;
            fulfilledValue = i18n;
        })
        return function (onFulfilled) {
            if (onFulfilled && fulfilled) {
                return Promise.resolve(onFulfilled(fulfilledValue));
            }
            return exports.then(onFulfilled);
        }
    })()

    exports.getBeeConfig = function () {
        return {
            isI18n: exports.isI18n,
            notI18n: exports.notI18n,
            getCurrency: exports.getCurrency,
            toI18nMoney: exports.toI18nMoney
        }
    }
    exports.bindBee = function (bee) {
        bee.$set(exports.getBeeConfig())
        exports.then(function (isI18n) {
            bee.$set('i18n', isI18n)
        })
    }
});
