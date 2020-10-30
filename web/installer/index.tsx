import * as React from './node_modules/react';
window['React16'] = React;
import * as ReactDOM from './node_modules/react-dom';
window['ReactDOM16'] = ReactDOM;
import { Installer } from './src/modules/installer';
import 'core-js/modules/es.object.assign'; // 解决IE浏览器不支持Object.assign 的问题
import '@babel/polyfill'; // 解决 regenerator time 的问题

ReactDOM.render(<Installer />, document.getElementById('appArea'));
