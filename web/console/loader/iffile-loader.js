/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
let loaderUtils = require('loader-utils');
// let parse = require('./lib/parse');
let fs = require('fs');

module.exports = function (source) {
	this.cacheable && this.cacheable();
	let options = loaderUtils.getOptions(this);

	let resource = this._module.resource;

	let sfilepath = resource.replace(/(\.tsx?)$/, '.' + options.version + '$1');

	if (fs.existsSync(sfilepath)) {
		this.addDependency(sfilepath);
		// console.log('>>>>>>>>>>>>>>>>>>>>>>>ifelse', sfilepath)
		source = fs.readFileSync(sfilepath, 'utf-8');
	}
	this.callback(null, source);
};

