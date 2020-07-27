import * as React from 'react';

import { insertCSS } from '@tencent/ff-redux';

import { InstallerAppContainer } from './components/InstallerApp';

insertCSS(
  'Installer',
  `
body{
  overflow-y: auto !important;
}`
);

export class Installer extends React.Component<any, any> {
  render() {
    return <InstallerAppContainer />;
  }
}
