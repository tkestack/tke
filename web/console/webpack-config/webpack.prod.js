const { entry, output, rules, alias, plugins, extensions, stats } = require('./webpack.base');

module.exports = {
  mode: 'production',

  entry,

  output,

  module: {
    rules
  },

  resolve: {
    extensions,
    alias
  },

  plugins,

  stats
};
