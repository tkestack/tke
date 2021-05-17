const { entry, output, rules, alias, plugins, extensions, stats } = require('./webpack.base');
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin');
const smp = new SpeedMeasurePlugin();

module.exports = smp.wrap({
  mode: 'development',

  entry,

  output,

  devServer: {
    contentBase: '../dist',
    open: false
  },

  module: {
    rules
  },

  resolve: {
    extensions,
    alias
  },

  plugins,

  stats
});
