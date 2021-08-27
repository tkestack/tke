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
const prefix = 'dynamic-inserted-css-';

/**
 * insert css to the current page
 * */
export function insertCSS(id: string, cssText: string) {
  let style: HTMLStyleElement;
  style = document.getElementById(prefix + id) as HTMLStyleElement;
  if (!style) {
    style = document.createElement('style');
    style.id = prefix + id;

    // IE8/9 can not use document.head
    document.getElementsByTagName('head')[0].appendChild(style);
  }
  if (style.textContent !== cssText) {
    style.textContent = cssText;
  }
  return style;
}
