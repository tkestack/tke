const fs = require('fs');

module.exports = function (source) {
  this.cacheable && this.cacheable();
  const options = this.getOptions();

  const resource = this.resourcePath;

  const sfilepath = resource.replace(/(\.tsx?)$/, '.' + options.version + '$1');

  if (fs.existsSync(sfilepath)) {
    this.addDependency(sfilepath);
    // console.log('>>>>>>>>>>>>>>>>>>>>>>>ifelse', sfilepath)
    source = fs.readFileSync(sfilepath, 'utf-8');
  }

  return source;
};
