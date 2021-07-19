import * as React from "react";
import { History } from "history";
export declare type Confirmation = typeof import("history/DOMUtils").getConfirmation;
interface HistoryContextValue {
    history: History;
    setConfirmation: (confirmation: Confirmation) => void;
}
export declare const HistoryContext: React.Context<HistoryContextValue>;
/**
 * 使用适配控制台的 history（React Hooks）
 *
 * 具体使用方法请参考 http://tapd.oa.com/QCloud_2015/markdown_wikis/view/#1010103951008571667
 *
 * @param getConfirmation (experimental) 路由切换前，让用户确认
 *
 * @example
```js
import React from 'react';
import { Router, Route, Link } from 'react-router-dom';
import { useHistory } from '@tea/app';
import { CvmDetailInstance } from './components';

export function CvmDetail() {
  // 此处为 React Hooks，只能在 Function Component 中使用
  const history = useHistory();

  return (
    <Router history={history}>
      <div>
        <Link to="/cvm/detail/ins-1q2w3e4r">ins-1q2w3e4r</Link>
        <Route path="/cvm/detail/:instanceId" component={CvmDetailInstance}></Route>
      </div>
    </Router>
  );
}
```
 */
export declare const useHistory: (getConfirmation?: typeof import("history/DOMUtils").getConfirmation) => History<any>;
export {};
