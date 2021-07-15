const path = require('path');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');

module.exports = ({ version, mode }) => ({
  mode,

  entry: `./index.${version}.tsx`,

  output: {
    filename: `static/js/index.${version}.[name].[contenthash].js`,
    path: path.resolve(__dirname, '../build')
  },

  optimization: {
    splitChunks: {
      cacheGroups: {
        commonVendor: {
          test: /[\\/]node_modules[\\/](react|react-dom|lodash|codemirror|validator|@tencent\/tea-app|@tencent\/tchart|esprima|validator)[\\/]/,
          filename: 'static/js/common-vendor.[contenthash].js',
          chunks: 'initial',
          priority: -10
        }
      }
    }
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

          'babel-loader',

          {
            loader: 'ts-loader',
            options: {
              happyPackMode: true,
              transpileOnly: true
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
        exclude: [path.resolve(__dirname, '../node_modules')]
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
      '@tea/component': path.resolve(__dirname, '../node_modules/tea-component/lib'),
      '@tea/component/*': path.resolve(__dirname, '../node_modules/tea-component/lib/*'),
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
      '@tencent/tea-component': path.resolve(__dirname, '../node_modules/tea-component'),
      '@tencent/tea-component/*': path.resolve(__dirname, '../node_modules/tea-component/*')
    }
  },

  plugins: [
    new ForkTsCheckerWebpackPlugin({
      eslint: {
        files: './src/**/*.{ts,tsx,js,jsx}'
      },

      issue: {
        // 排除eslint warning的打印
        exclude: {
          origin: 'eslint',
          severity: 'warning'
        }
      }
    }),

    new webpack.ProvidePlugin({
      Buffer: ['buffer', 'Buffer']
    }),

    ...(mode === 'production'
      ? []
      : [
          new BundleAnalyzerPlugin(),
          new HtmlWebpackPlugin({
            template: path.resolve(__dirname, '../public/index.html'),
            inject: false,
            templateParameters: (_, { js }) => {
              const index = js.find(path => path.includes('index'));
              const commonVendor = js.find(path => path.includes('common-vendor'));
              return version === 'tke'
                ? { TKE_JS_NAME: index, PROJECT_JS_NAME: '', COMMON_VENDOR_JS_NAME: commonVendor }
                : {
                    TKE_JS_NAME: '',
                    PROJECT_JS_NAME: index,
                    COMMON_VENDOR_JS_NAME: commonVendor
                  };
            }
          })
        ])
  ],

  stats: {
    warningsFilter: /export .* was not found in/
  }
});
