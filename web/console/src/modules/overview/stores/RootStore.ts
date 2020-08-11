import { createStore } from '@tencent/ff-redux';

import { generateResetableReducer } from '../../../../helpers';
import { RootReducer } from '../reducers/RootReducer';

export function configStore() {
  const store = createStore(generateResetableReducer(RootReducer));

  // hot reloading
  if (typeof module !== 'undefined' && (module as any).hot) {
    (module as any).hot.accept('../reducers/RootReducer', () => {
      store.replaceReducer(generateResetableReducer(require('../reducers/RootReducer').RootReducer));
    });
  }

  return store;
}
