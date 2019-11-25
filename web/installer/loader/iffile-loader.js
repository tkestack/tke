// let loaderUtils = require('loader-utils');
// let parse = require('./lib/parse');
let fs = require('fs');

module.exports = function (source) {
	this.cacheable && this.cacheable();
	// let options = loaderUtils.getOptions(this);

	let resource = this._module.resource;

	let sfilepath = resource.replace(/(\.tsx?)$/, '.' + process.env.Version + '$1');

	if (fs.existsSync(sfilepath)) {
		this.addDependency(sfilepath);
		// console.log('>>>>>>>>>>>>>>>>>>>>>>>ifelse', sfilepath)
		source = fs.readFileSync(sfilepath, 'utf-8');
	}
	this.callback(null, source);
};

