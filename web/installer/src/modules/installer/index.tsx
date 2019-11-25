import * as React from 'react';
import { InstallerAppContainer } from './components/InstallerApp';
import { insertCSS } from '@tencent/qcloud-lib';

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
