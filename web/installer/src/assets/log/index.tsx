import * as React from 'react';
import { configStore } from './stores/RootStore';
import { ResetStoreAction } from './../../../helpers';
import { Provider } from 'react-redux';
import { LogApp } from './components/LogApp';

const store = configStore();

export class Log extends React.Component<any, any> {
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <LogApp />
      </Provider>
    );
  }
}
