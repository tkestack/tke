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

import * as marked from 'marked';
import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';

// TODO
// require('highlight.js/styles/atom-one-dark.css');

export interface MarkdownProps extends BaseReactProps {
  /**显示文本 */
  text?: string;
}

export class Markdown extends React.Component<MarkdownProps, {}> {
  componentWillMount() {
    marked.setOptions({
      gfm: true,
      tables: true,
      breaks: false,
      pedantic: false,
      sanitize: false,
      smartLists: true,
      smartypants: false,
      highlight: function(code, lang) {
        try {
          //return require('highlight.js').highlight(lang, code, true).value;
        } catch (e) {
          return `<pre>${code}</pre>`;
        }
      }
    });
  }

  render() {
    let { style, text = '', className } = this.props;

    return (
      <div
        style={style}
        className={'markdown-text-box ' + className}
        dangerouslySetInnerHTML={{ __html: marked(text) }}
      />
    );
  }
}
