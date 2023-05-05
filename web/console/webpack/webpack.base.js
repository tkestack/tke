/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
          test: /[\\/]node_modules[\\/](react|react-dom|lodash|codemirror|validator|@tencent|esprima)[\\/]/,
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
        exclude: [path.resolve(__dirname, '../node_modules'), path.resolve(__dirname, '../tencent/tea-app')]
      },

      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      },

      {
        test: /\.less$/,
        use: ['style-loader', 'css-loader', 'less-loader']
      },

      {
        test: /\.svg$/,
        use: [
          {
            loader: 'url-loader',
            options: {
              esModule: false
            }
          }
        ]
      }
    ]
  },

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx', '.json', 'css'],

    alias: {
      // 国际化语言包覆盖
      '@i18n/translation': path.resolve(__dirname, `../i18n/translation/zh.js`),
      '@i18n/translation_en': path.resolve(__dirname, `../i18n/translation/en.js`),
      '@tea/component': path.resolve(__dirname, '../node_modules/tea-component/es'),
      '@tea/component/*': path.resolve(__dirname, '../node_modules/tea-component/es/*'),
      '@paas/paas-lib': path.resolve(__dirname, '../lib'),
      '@helper': path.resolve(__dirname, '../helpers'),
      '@helper/*': path.resolve(__dirname, '../helpers/*'),
      '@config': path.resolve(__dirname, '../config'),
      '@config/*': path.resolve(__dirname, '../config/*'),
      '@src/*': path.resolve(__dirname, '../src/*'),
      '@src': path.resolve(__dirname, '../src'),
      '@common': path.resolve(__dirname, '../src/modules/common'),
      '@common/*': path.resolve(__dirname, '../src/modules/common/*'),
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
      // moment: path.resolve(__dirname, '../node_modules/dayjs'),
      '@tencent/tea-component': path.resolve(__dirname, '../node_modules/tea-component'),
      '@tencent/tea-component/lib/*': path.resolve(__dirname, '../node_modules/tea-component/es/*'),
      '@tencent/tchart': path.resolve(__dirname, '../tencent/tchart/src/panel/index.tsx'),
      '@tencent/tea-app': path.resolve(__dirname, '../tencent/tea-app')
    }
  },

  plugins: [
    new ForkTsCheckerWebpackPlugin({
      async: false,
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
          // new BundleAnalyzerPlugin(),
          new HtmlWebpackPlugin({
            template: path.resolve(__dirname, '../public/index.tmpl.html'),
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

  ignoreWarnings: [/export .* was not found in/]
});
