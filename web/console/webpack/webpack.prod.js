const path = require('path');
const HappyPack = require('happypack');
const os = require('os');
const happyThreadPool = HappyPack.ThreadPool({
  size: os.cpus().length
});

// 不需要打包进去的内容
const externals = require('./externals');
// 进行ts的校验
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
// 修改进度条的展示
const ProgressBarPlugin = require('progress-bar-webpack-plugin');
const version = (process.argv[4] && process.argv[4].split('=')[1]) || 'tke';
const lng = (process.argv[5] && process.argv[5].split('=')[1]) || 'zh';
const filename = (process.argv[6] && process.argv[6].split('=')[1]) || 'index2';
// const version = 'tke';
// const lng = 'zh';

module.exports = {
  mode: 'production',
  entry: {
    app: ['./index.tsx']
  },

  output: {
    path: path.resolve(__dirname, '../build'),
    filename: `static/js/${filename}.js`
  },

  module: {
    rules: [{
        test: /\.tsx?$/,
        use: [
          'happypack/loader?id=happyBabel',
          'happypack/loader?id=happyTs',
          'happypack/loader?id=happyESLint',
          path.resolve('./loader/iffile-loader.js?version=' + version),
          path.resolve('./loader/ifelse-loader.js?version=' + version)
        ]
      },
      {
        test: /\.jsx?$/,
        use: ['happypack/loader?id=happyBabel'],
        exclude: [path.resolve(__dirname, '../node_modules')],
        include: [path.resolve(__dirname, '../node_modules/tchart')]
      },
      {
        test: /\.css?$/,
        use: ['happypack/loader?id=happyCSS']
      }
    ]
  },

  plugins: [
    new HappyPack({
      id: 'happyTs',
      loaders: [{
        loader: 'ts-loader',
        options: {
          happyPackMode: true,
          transpileOnly: true
        }
      }],
      threadPool: happyThreadPool
    }),

    new HappyPack({
      id: 'happyESLint',
      loaders: [{
        loader: 'eslint-loader',
        options: {
          failOnWarning: true,
          failOnError: true
        }
      }],
      threadPool: happyThreadPool
    }),

    new HappyPack({
      id: 'happyBabel',
      loaders: ['babel-loader'],
      threadPool: happyThreadPool
    }),

    new HappyPack({
      id: 'happyCSS',
      loaders: ['style-loader', 'css-loader'],
      threadPool: happyThreadPool
    }),

    new ProgressBarPlugin(),

    new ForkTsCheckerWebpackPlugin({
      async: false,
      checkSyntacticErrors: true
    })
  ],

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx', '.json', 'css'],
    alias: {
      // 国际化语言包覆盖
      '@i18n/translation': path.resolve(__dirname, `../i18n/translation/${lng}.js`),
      '@i18n/translation_en': path.resolve(__dirname, `../i18n/translation/en.js`),
      '@tea/app': path.resolve(__dirname, '../node_modules/@tencent/tea-app'),
      '@tea/app/*': path.resolve(__dirname, '../node_modules/@tencent/tea-app/lib/*'),
      '@tea/component': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib'),
      '@tea/component/*': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib/*'),
      '@paas/paas-lib': path.resolve(__dirname, '../lib'),
      '@helper': path.resolve(__dirname, '../helpers'),
      '@helper/*': path.resolve(__dirname, '../helpers/*'),
      '@config': path.resolve(__dirname, '../config'),
      '@config/*': path.resolve(__dirname, '../config/*'),
      '@src/*': path.resolve(__dirname, '../src/*'),
      '@src': path.resolve(__dirname, '../src'),
      '@tencent/ff-validator': path.resolve(__dirname, '../lib/ff-validator'),
      '@tencent/ff-validator/*': path.resolve(__dirname, '../lib/ff-validator/*'),
      '@tencent/ff-redux': path.resolve(__dirname, '../lib/ff-redux'),
      '@tencent/ff-redux/*': path.resolve(__dirname, '../lib/ff-redux/*'),
      '@tencent/ff-component': path.resolve(__dirname, '../lib/ff-component'),
      '@tencent/ff-component/*': path.resolve(__dirname, '../lib/ff-component/*'),
      '@tencent/qcloud-redux-fetcher': path.resolve(__dirname, '../lib/ff-redux/libs/qcloud-redux-fetcher/'),
      '@tencent/qcloud-redux-query': path.resolve(__dirname, '../lib/ff-redux/libs/qcloud-redux-query/'),
      '@tencent/qcloud-redux-workflow': path.resolve(__dirname, '../lib/ff-redux/libs/qcloud-redux-workflow/'),
      // react 和 react-dom 控制台通过全局变量提供，我们不打包
      react: path.resolve(__dirname, './alias/react.js'),
      'react-dom': path.resolve(__dirname, './alias/react-dom.js'),
      d3: path.resolve(__dirname, '../node_modules/d3')
    }
  },

  externals: externals
};
