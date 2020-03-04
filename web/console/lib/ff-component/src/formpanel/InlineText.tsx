import * as React from 'react';

import { Text } from '@tencent/tea-component';

interface FormPanelInlineTextProps {
  parent?: keyof JSX.IntrinsicElements;

  align?: 'left' | 'center' | 'right' | 'justify';
  /**
   * 文本的垂直居中方式
   */
  verticalAlign?: 'baseline' | 'top' | 'middle' | 'bottom' | 'text-top' | 'text-bottom';
  /**
   * 指定 `overflow` 为 `true` 对文本进行单行溢出控制，宽度溢出的文本会以 ... 显示
   */
  overflow?: boolean;
  /**
   * 指定 `nowrap` 为 `true` 强制文本不换行，超长溢出
   */
  nowrap?: boolean;
  /**
   * 内容所使用的 `tooltip`
   */
  tooltip?: React.ReactNode;
  /**
   * 文本的颜色主题
   */
  theme?: 'text' | 'label' | 'weak' | 'strong' | 'primary' | 'success' | 'warning' | 'danger';
  /**
   * 文本的背景颜色主题
   */
  bgTheme?: 'success' | 'warning' | 'danger';
  /**
   * 文本浮动方式，配置 `clear` 以清除浮动
   */
  float?: 'left' | 'right' | 'clear';

  style?: React.CSSProperties;
  className?: string;

  children?: React.ReactNode;
}
function FormPanelInlineText({ children, ...textProps }: FormPanelInlineTextProps) {
  return (
    <Text className="tea-form__help-text--inline" {...textProps}>
      {children}
    </Text>
  );
}
export { FormPanelInlineText, FormPanelInlineTextProps };
