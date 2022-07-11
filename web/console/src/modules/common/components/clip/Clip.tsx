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

import React, { useRef, useLayoutEffect } from 'react';
import { Icon } from 'tea-component';
import ClipboardJS from 'clipboard';

export interface ClipProps {
  target: string;
  children?: React.ReactNode;
  className?: string;
  isShow?: boolean;
  isShowTip?: boolean;
}

export const Clip = ({ target, className, isShow = true, isShowTip }: ClipProps) => {
  const copyBtn = useRef(null);
  const copyInstance = useRef(null);

  useLayoutEffect(() => {
    copyInstance.current && copyInstance.current.destroy();

    if (copyBtn.current) {
      copyInstance.current = new ClipboardJS(copyBtn.current, {
        text() {
          return document?.querySelector(target)?.textContent ?? '';
        }
      });
    }
  }, [copyBtn, target]);

  return isShow && <Icon type="copy" className={className} ref={copyBtn} />;
};
