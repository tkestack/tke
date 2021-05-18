const path = require('path');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
const webpack = require('webpack');

module.exports = ({ version, mode }) => ({
  mode,

  entry: './index.tsx',

  output: {
    filename: `index.${version}.js`,
    path: path.resolve(__dirname, '../build')
  },

  module: {
    rules: [
      // 因为tchart组件依赖window.React16和window.ReactDOM16,所以用expose-loader暴露到全局
      {
        test: require.resolve('react'),
        loader: 'expose-loader',
        options: {
          exposes: 'React16'
        }
      },

      {
        test: require.resolve('react-dom'),
        loader: 'expose-loader',
        options: {
          exposes: 'ReactDOM16'
        }
      },

      {
        test: /\.tsx?$/,
        use: [
          'thread-loader',

          {
            loader: 'ts-loader',
            options: {
              happyPackMode: true
            }
          },

          {
            loader: path.resolve(__dirname, './loaders/iffile-loader.js'),
            options: {
              version
            }
          },

          {
            loader: path.resolve(__dirname, './loaders/ifelse-loader.js'),
            options: {
              version
            }
          }
        ]
      },

      {
        test: /\.(js|jsx)$/,
        use: ['thread-loader', 'babel-loader'],
        exclude: [path.resolve(__dirname, '../node_modules')],
        include: [path.resolve(__dirname, '../node_modules/tchart')]
      },

      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      }
    ]
  },

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx', '.json', 'css'],

    alias: {
      // 国际化语言包覆盖
      '@i18n/translation': path.resolve(__dirname, `../i18n/translation/zh.js`),
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
      '@': path.resolve(__dirname, '../'),
      d3: path.resolve(__dirname, '../node_modules/d3'),
      moment: path.resolve(__dirname, '../node_modules/dayjs'),
      '@tencent/tea-component': path.resolve(__dirname, '../node_modules/@tencent/tea-component')
    }
  },

  plugins: [
    new webpack.ProvidePlugin({
      Buffer: ['buffer', 'Buffer']
    }),

    ...(mode === 'production' ? [] : [new BundleAnalyzerPlugin()])
  ],

  stats: {
    warningsFilter: /export .* was not found in/
  }
});
