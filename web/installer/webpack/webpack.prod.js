/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
const webpack = require('webpack');
const HappyPack = require('happypack');
const os = require('os');
const happyThreadPool = HappyPack.ThreadPool({
  size: os.cpus().length
});

module.exports = {
  mode: 'production',
  entry: {
    app: ['./index.tsx']
  },

  output: {
    path: path.resolve(__dirname, '../build'),
    filename: 'static/js/installer.js'
  },

  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: [
          'happypack/loader?id=happyBabel',
          {
            loader: 'ts-loader',
            options: {
              transpileOnly: true
            }
          }
        ]
        //exclude: [path.resolve(__dirname, "../node_modules")]
      },
      {
        test: /\.jsx?$/,
        use: ['happypack/loader?id=happyBabel']
      },
      {
        test: /\.css?$/,
        use: ['happypack/loader?id=happyCSS']
      }
    ]
  },

  plugins: [
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

    new webpack.ProgressPlugin(function handler(percentage, msg) {
      let line = 'Processing ' + msg + '\t' + Math.floor(percentage * 100) + '%';
      console.log(line);
    })
  ],

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx', '.json', 'css'],
    alias: {
      '@tea/app': path.resolve(__dirname, '../node_modules/@tencent/tea-app'),
      '@tea/app/*': path.resolve(__dirname, '../node_modules/@tencent/tea-app/lib/*'),
      '@tea/component': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib'),
      '@tea/component/*': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib/*'),
      '@tencent/ff-validator': path.resolve(__dirname, '../lib/ff-validator'),
      '@tencent/ff-validator/*': path.resolve(__dirname, '../lib/ff-validator/*'),
      '@tencent/ff-redux': path.resolve(__dirname, '../lib/ff-redux'),
      '@tencent/ff-redux/*': path.resolve(__dirname, '../lib/ff-redux/*'),
      '@tencent/ff-component': path.resolve(__dirname, '../lib/ff-component'),
      '@tencent/ff-component/*': path.resolve(__dirname, '../lib/ff-component/*')
    }
  },

  externals: {
    react: 'window.React16',
    'react-dom': 'window.ReactDOM16'
  }
};
