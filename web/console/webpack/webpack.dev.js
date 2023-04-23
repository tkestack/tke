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
const baseConfig = require('./webpack.base');
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin');
const smp = new SpeedMeasurePlugin();
const path = require('path');
const { Host, Cookie } = require('../server.config');
const { createCSRFHeader } = require('./csrf');

module.exports = ({ version }) =>
  smp.wrap({
    ...baseConfig({ version, mode: 'development' }),

    devServer: {
      contentBase: path.resolve(__dirname, '../public'),
      historyApiFallback: true,
      compress: true,
      open: true,
      openPage: version === 'tke' ? '/' : '/tkestack-project',
      port: 8181,
      proxy: {
        '/api': {
          target: Host,
          secure: false,
          changeOrigin: true,
          headers: { Cookie, ...createCSRFHeader(Cookie) }
        },

        '/websocket': {
          target: Host.replace(/^http/, 'ws'),
          ws: true,
          logLevel: 'debug',
          secure: false,
          changeOrigin: true,
          headers: { Cookie }
        }
      }
    }
  });
