const baseConfig = require('./webpack.base');
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin');
const smp = new SpeedMeasurePlugin();

module.exports = ({ version }) =>
  smp.wrap({
    ...baseConfig({ version, mode: 'development' }),

    devServer: {
      contentBase: '../dist',
      open: false
    }
  });
