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
import * as React from './node_modules/react';
window['React16'] = React;
import * as ReactDOM from './node_modules/react-dom';
window['ReactDOM16'] = ReactDOM;
import { Installer } from './src/modules/installer';
import 'core-js/modules/es.object.assign'; // 解决IE浏览器不支持Object.assign 的问题
import '@babel/polyfill'; // 解决 regenerator time 的问题

ReactDOM.render(<Installer />, document.getElementById('appArea'));
