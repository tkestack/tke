import * as React from 'react';

/**
 * 基础 React Props 类型，建议 Props 从该类型拓展
 */
export declare interface BaseReactProps extends React.Props<any> {
  className?: string;
  style?: React.CSSProperties;
}
