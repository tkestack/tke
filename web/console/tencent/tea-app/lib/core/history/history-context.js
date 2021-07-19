"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var React = require("react");
var useContext = React.useContext, useEffect = React.useEffect;
// 使用 Context 提供 history，以及设置 conformation 的方法
exports.HistoryContext = React.createContext(null);
exports.HistoryContext.displayName = "Tea.History";
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
exports.useHistory = function (getConfirmation) {
    var _a = useContext(exports.HistoryContext), history = _a.history, setConfirmation = _a.setConfirmation;
    // 设置提供的确认函数
    useEffect(function () {
        if (getConfirmation) {
            setConfirmation(getConfirmation);
        }
        return function () { return setConfirmation(null); };
    }, []);
    return history;
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGlzdG9yeS1jb250ZXh0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2NvcmUvaGlzdG9yeS9oaXN0b3J5LWNvbnRleHQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSw2QkFBK0I7QUFFdkIsSUFBQSw2QkFBVSxFQUFFLDJCQUFTLENBQVc7QUFTeEMsOENBQThDO0FBQ2pDLFFBQUEsY0FBYyxHQUFHLEtBQUssQ0FBQyxhQUFhLENBQXNCLElBQUksQ0FBQyxDQUFDO0FBRTdFLHNCQUFjLENBQUMsV0FBVyxHQUFHLGFBQWEsQ0FBQztBQUUzQzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztHQTRCRztBQUNVLFFBQUEsVUFBVSxHQUFHLFVBQUMsZUFBOEI7SUFDakQsSUFBQSx1Q0FBeUQsRUFBdkQsb0JBQU8sRUFBRSxvQ0FBOEMsQ0FBQztJQUVoRSxZQUFZO0lBQ1osU0FBUyxDQUFDO1FBQ1IsSUFBSSxlQUFlLEVBQUU7WUFDbkIsZUFBZSxDQUFDLGVBQWUsQ0FBQyxDQUFDO1NBQ2xDO1FBQ0QsT0FBTyxjQUFNLE9BQUEsZUFBZSxDQUFDLElBQUksQ0FBQyxFQUFyQixDQUFxQixDQUFDO0lBQ3JDLENBQUMsRUFBRSxFQUFFLENBQUMsQ0FBQztJQUVQLE9BQU8sT0FBTyxDQUFDO0FBQ2pCLENBQUMsQ0FBQyIsInNvdXJjZXNDb250ZW50IjpbImltcG9ydCAqIGFzIFJlYWN0IGZyb20gXCJyZWFjdFwiO1xuaW1wb3J0IHsgSGlzdG9yeSB9IGZyb20gXCJoaXN0b3J5XCI7XG5jb25zdCB7IHVzZUNvbnRleHQsIHVzZUVmZmVjdCB9ID0gUmVhY3Q7XG5cbmV4cG9ydCB0eXBlIENvbmZpcm1hdGlvbiA9IHR5cGVvZiBpbXBvcnQoXCJoaXN0b3J5L0RPTVV0aWxzXCIpLmdldENvbmZpcm1hdGlvbjtcblxuaW50ZXJmYWNlIEhpc3RvcnlDb250ZXh0VmFsdWUge1xuICBoaXN0b3J5OiBIaXN0b3J5O1xuICBzZXRDb25maXJtYXRpb246IChjb25maXJtYXRpb246IENvbmZpcm1hdGlvbikgPT4gdm9pZDtcbn1cblxuLy8g5L2/55SoIENvbnRleHQg5o+Q5L6bIGhpc3RvcnnvvIzku6Xlj4rorr7nva4gY29uZm9ybWF0aW9uIOeahOaWueazlVxuZXhwb3J0IGNvbnN0IEhpc3RvcnlDb250ZXh0ID0gUmVhY3QuY3JlYXRlQ29udGV4dDxIaXN0b3J5Q29udGV4dFZhbHVlPihudWxsKTtcblxuSGlzdG9yeUNvbnRleHQuZGlzcGxheU5hbWUgPSBcIlRlYS5IaXN0b3J5XCI7XG5cbi8qKlxuICog5L2/55So6YCC6YWN5o6n5Yi25Y+w55qEIGhpc3RvcnnvvIhSZWFjdCBIb29rc++8iVxuICogXG4gKiDlhbfkvZPkvb/nlKjmlrnms5Xor7flj4LogIMgaHR0cDovL3RhcGQub2EuY29tL1FDbG91ZF8yMDE1L21hcmtkb3duX3dpa2lzL3ZpZXcvIzEwMTAxMDM5NTEwMDg1NzE2NjdcbiAqIFxuICogQHBhcmFtIGdldENvbmZpcm1hdGlvbiAoZXhwZXJpbWVudGFsKSDot6/nlLHliIfmjaLliY3vvIzorqnnlKjmiLfnoa7orqRcbiAqIFxuICogQGV4YW1wbGUgXG5gYGBqc1xuaW1wb3J0IFJlYWN0IGZyb20gJ3JlYWN0JztcbmltcG9ydCB7IFJvdXRlciwgUm91dGUsIExpbmsgfSBmcm9tICdyZWFjdC1yb3V0ZXItZG9tJztcbmltcG9ydCB7IHVzZUhpc3RvcnkgfSBmcm9tICdAdGVhL2FwcCc7XG5pbXBvcnQgeyBDdm1EZXRhaWxJbnN0YW5jZSB9IGZyb20gJy4vY29tcG9uZW50cyc7XG5cbmV4cG9ydCBmdW5jdGlvbiBDdm1EZXRhaWwoKSB7XG4gIC8vIOatpOWkhOS4uiBSZWFjdCBIb29rc++8jOWPquiDveWcqCBGdW5jdGlvbiBDb21wb25lbnQg5Lit5L2/55SoXG4gIGNvbnN0IGhpc3RvcnkgPSB1c2VIaXN0b3J5KCk7XG5cbiAgcmV0dXJuIChcbiAgICA8Um91dGVyIGhpc3Rvcnk9e2hpc3Rvcnl9PlxuICAgICAgPGRpdj5cbiAgICAgICAgPExpbmsgdG89XCIvY3ZtL2RldGFpbC9pbnMtMXEydzNlNHJcIj5pbnMtMXEydzNlNHI8L0xpbms+XG4gICAgICAgIDxSb3V0ZSBwYXRoPVwiL2N2bS9kZXRhaWwvOmluc3RhbmNlSWRcIiBjb21wb25lbnQ9e0N2bURldGFpbEluc3RhbmNlfT48L1JvdXRlPlxuICAgICAgPC9kaXY+XG4gICAgPC9Sb3V0ZXI+XG4gICk7XG59XG5gYGBcbiAqL1xuZXhwb3J0IGNvbnN0IHVzZUhpc3RvcnkgPSAoZ2V0Q29uZmlybWF0aW9uPzogQ29uZmlybWF0aW9uKSA9PiB7XG4gIGNvbnN0IHsgaGlzdG9yeSwgc2V0Q29uZmlybWF0aW9uIH0gPSB1c2VDb250ZXh0KEhpc3RvcnlDb250ZXh0KTtcblxuICAvLyDorr7nva7mj5DkvpvnmoTnoa7orqTlh73mlbBcbiAgdXNlRWZmZWN0KCgpID0+IHtcbiAgICBpZiAoZ2V0Q29uZmlybWF0aW9uKSB7XG4gICAgICBzZXRDb25maXJtYXRpb24oZ2V0Q29uZmlybWF0aW9uKTtcbiAgICB9XG4gICAgcmV0dXJuICgpID0+IHNldENvbmZpcm1hdGlvbihudWxsKTtcbiAgfSwgW10pO1xuXG4gIHJldHVybiBoaXN0b3J5O1xufTtcbiJdfQ==