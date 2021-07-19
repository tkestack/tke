const merge = require('webpack-merge');
const common = require('./config.common.js');

// 不需要打包进去的内容
const externals = require('./externals');

// 在客户程序里会进行压缩，保留源码方便调试
module.exports = merge(common, {
  mode: 'development',
  entry: {
    TChart: './src/panel/index.tsx',
  },
  output: {
    libraryTarget: 'commonjs2',
    globalObject: `(typeof self !== 'undefined' ? self : this)`,
    publicPath: '//imgcache.qq.com/tchart/build/'
  },
  devtool: 'source-map',
  externals,
  optimization: {
    splitChunks: {
      maxAsyncRequests: 1
    }
  }
});