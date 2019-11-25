import * as React from 'react';
import * as marked from 'marked';
import { BaseReactProps } from '@tencent/qcloud-lib';

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
