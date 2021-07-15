const webpack = require('webpack');
const merge = require('webpack-merge');
const common = require('./config.common.js');


module.exports = merge(common, {
  mode: 'development',
  entry: {
    TChart: './src/tce/app.tsx',
  },
  devtool: 'inline-source-map',
  devServer: {
    contentBase: '../build',
    hot: false
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
  ],
});