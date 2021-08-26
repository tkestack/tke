/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import * as React from 'react';

import { Text } from '@tencent/tea-component';

interface FormPanelHelpTextProps {
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
function FormPanelHelpText({ children, ...textProps }: FormPanelHelpTextProps) {
  return (
    <Text className="tea-form__help-text" {...textProps}>
      {children}
    </Text>
  );
}
export { FormPanelHelpText, FormPanelHelpTextProps };
