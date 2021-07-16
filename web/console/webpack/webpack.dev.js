const baseConfig = require('./webpack.base');
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin');
const smp = new SpeedMeasurePlugin();
const path = require('path');
const { Host, Cookie } = require('../server.config');

module.exports = ({ version }) =>
  smp.wrap({
    ...baseConfig({ version, mode: 'development' }),

    devServer: {
      contentBase: path.resolve(__dirname, '../public'),
      historyApiFallback: true,
      compress: true,
      open: true,
      port: 8181,
      proxy: {
        '/api': {
          target: Host,
          secure: false,
          changeOrigin: true,
          headers: { Cookie }
        },
        '/apis': {
          target: Host,
          secure: false,
          changeOrigin: true,
          headers: { Cookie }
        }
      }
    }
  });
