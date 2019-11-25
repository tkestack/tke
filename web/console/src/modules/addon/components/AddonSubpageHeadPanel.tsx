import * as React from 'React';
import { Justify, Icon, Text } from '@tencent/tea-component';
import { router } from '../router';
import { RouteState } from '../../../../helpers';

interface SubpageOptions {
  /** route路由模块 */
  route: RouteState;
}

export function AddonSubpageHeadPanel(options: SubpageOptions) {
  let { route } = options;

  return (
    <Justify
      left={
        <React.Fragment>
          <Icon
            style={{ cursor: 'pointer' }}
            type="btnback"
            className="tea-mr-3n"
            onClick={() => {
              router.navigate({}, route.queries);
            }}
          />
          <h2>{`新建扩展组件`}</h2>
        </React.Fragment>
      }
    />
  );
}
