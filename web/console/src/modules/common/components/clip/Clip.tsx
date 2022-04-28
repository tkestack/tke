/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

import React, { useLayoutEffect, useState } from 'react';
import { Copy, Icon } from 'tea-component';

export interface ClipProps {
  target: string;
  children?: React.ReactNode;
  className?: string;
  isShow?: boolean;
  isShowTip?: boolean;
}

export const Clip = ({ target, className, isShow = true, isShowTip }: ClipProps) => {
  const [text, setText] = useState('');

  useLayoutEffect(() => {
    let targetText = '';
    try {
      targetText = document?.querySelector(target)?.textContent ?? '';
    } catch (error) {
      console.log(error);
    }

    setText(targetText);
  }, [target]);

  return (
    isShow && (
      <Copy text={text}>
        <Icon type="copy" className={className} />
      </Copy>
    )
  );
};
