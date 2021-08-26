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
import { UnControlled as CodeMirror } from 'react-codemirror2';

import { insertCSS } from '@tencent/ff-redux';

import { RootProps } from '../ClusterApp';
import { YamlSearchHelperPanel } from '../../../common/components';
import { ExternalLink } from '@tea/component';
import { t } from '@tencent/tea-app/lib/i18n';

// 这里是对editor一些配置
require('codemirror/mode/yaml/yaml');
// 这里是引入yaml search的内容
require('codemirror/addon/search/search');
require('codemirror/addon/dialog/dialog');
insertCSS('codemirror2-theme-monkai', require('codemirror/addon/dialog/dialog.css'));

insertCSS('codemirror2-theme-monkai', require('codemirror/theme/monokai.css'));

interface YamlEditorPanelProps extends RootProps {
  /** 输入的内容 */
  config?: string;

  /** 是否只读 */
  readOnly?: boolean;

  /** 回调函数，处理输入数据 */
  handleInputForEditor?: (config: string) => void;

  /** 当前的mode */
  mode?: string | 'text/x-yaml' | 'text/x-sh';

  /** 编辑器的高度，默认为600 */
  height?: number;

  /** 行距，默认为20 */
  lineHeight?: number;

  /** 是否需要不断的刷新内容，更新initValue */
  isNeedRefreshContent?: boolean;
}

interface YamlEditorPanelState {
  /** 初始化的内容 */
  initValue: string;
  showSearch: boolean;
}

export class YamlEditorPanel extends React.Component<YamlEditorPanelProps, YamlEditorPanelState> {
  constructor(props) {
    super(props);
    this.state = {
      initValue: this.props.config,
      showSearch: false
    };
  }

  componentWillReceiveProps(nextProps: YamlEditorPanelProps) {
    let { isNeedRefreshContent = false, config } = nextProps;

    if (isNeedRefreshContent && config !== this.props.config) {
      this.setState({
        initValue: config
      });
    }
  }

  componentWillMount() {
    let { height = 600, lineHeight = 20 } = this.props;
    insertCSS(
      'YamlEditorPanel',
      `.CodeMirror{height:${height}px;overflow:auto;overflow-x:hidden;font-size:12px}.CodeMirror-code>div{line-height: ${lineHeight}px}`
    );
  }

  render() {
    let { readOnly, mode = 'text/x-yaml', handleInputForEditor } = this.props;

    const codeOptions = {
      lineNumbers: true,
      mode,
      theme: 'monokai',
      readOnly: readOnly ? true : false, // nocursor表明焦点不能展示，不会展示光标
      spellcheck: true, // 是否开启单词校验
      autocorrect: true, // 是否开启自动修正
      lineWrapping: true, // 自动换行
      styleActiveLine: true, // 当前行背景高亮
      tabSize: 2 // tab 默认是2格
    };

    return (
      <div>
        <ExternalLink
          href="#"
          onMouseOver={() => {
            this.setState({ showSearch: true });
          }}
          onMouseOut={() => {
            this.setState({ showSearch: false });
          }}
        >
          {t('搜索帮助')}
        </ExternalLink>
        <YamlSearchHelperPanel isShow={this.state.showSearch} />
        <CodeMirror
          className={'codeMirrorHeight'}
          value={this.state.initValue}
          options={codeOptions}
          onChange={(editor, data, value) => {
            // 配置项当中的value 不用props.config 是因为 更新之后，yaml的光标会默认跳转到末端
            !readOnly && handleInputForEditor(value);
          }}
        />
      </div>
    );
  }
}
