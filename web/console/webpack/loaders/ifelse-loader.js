const parse = require('./lib/parse');

module.exports = function (source) {
  this.cacheable && this.cacheable();
  const options = this.getOptions();

  const resource = this.resourcePath;

  const verbose = {
    Version: options.version,
    [options.version]: true
  };
  try {
    source = parse(source, verbose, verbose);

    return source;
  } catch (err) {
    const errorMessage = `ifdef-loader error: ${err},file:${resource}`;

    throw new Error(errorMessage);
  }
};
