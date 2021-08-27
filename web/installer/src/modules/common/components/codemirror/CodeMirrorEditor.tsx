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
import * as React from 'react';
import { Controlled as CodeMirror } from 'react-codemirror2';

import { BaseReactProps, insertCSS } from '@tencent/ff-redux';
import { Button, Modal } from '@tencent/tea-component';

import { Clip } from '../';

require('codemirror/mode/javascript/javascript.js');
require('codemirror/lib/codemirror.css');

insertCSS('codemirror-theme-monkai', require('codemirror/theme/monokai.css'));
// insertCSS('codemirror-theme-eclipse', require('codemirror/theme/eclipse.css'));

export interface CodeMirrorEditorProps extends BaseReactProps {
  /**编辑器名称 */
  title?: string;

  /**是否显示标题栏 */
  isShowHeader?: boolean;

  /**是否开启大编辑框模式 */
  isOpenDialogEditor?: boolean;

  /**是否提供复制操作 */
  isOpenClip?: boolean;

  /**默认值 */
  defaultValue?: string;

  /**值 */
  value?: string;

  /**change事件 */
  onChange?: (val: string) => void;

  /**语言设置 */
  mode?: string;

  /**是否显示行号 */
  lineNumbers?: boolean;

  /**主题设置 */
  theme?: string;

  /**是否只读 */
  readOnly?: boolean;

  /**宽度 */
  width?: number;

  /**高度 */
  height?: number;

  /**对话框高度 */
  dHeight?: number;

  /**是否强制刷新 */
  isForceRefresh?: boolean;
}

interface CodeMirrorEditorState {
  /**编辑器代码副本 */
  code?: string;

  /**弹窗代码副本 */
  diaCode?: string;

  /**是否显示弹窗 */
  isShowDialog?: boolean;

  /**是否刷新 */
  isRefresh?: boolean;
}

export class CodeMirrorEditor extends React.Component<CodeMirrorEditorProps, CodeMirrorEditorState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      code: props.value,
      diaCode: props.value,
      isShowDialog: false,
      isRefresh: true
    };
  }

  componentDidMount() {
    let { height, width, dHeight, theme } = this.props;
    insertCSS(
      'codemirror-theme-' + theme,
      `
            .CodeMirror{
                height: ${height}px;
                width: ${width}px
            }
            .tc-15-rich-dialog .CodeMirror{
                height: ${dHeight || 300}px
            }
        `
    );
  }

  componentWillReceiveProps(nextProps) {
    if (this.state.code !== nextProps.value) {
      this.setState({ code: nextProps.value, diaCode: nextProps.value });
      this.refreshEditor();
    }
  }

  refreshEditor() {
    this.setState({ isRefresh: false });
    setTimeout(() => {
      this.setState({ isRefresh: true });
    }, 50);
  }

  handChange(code: string) {
    this.setState({ code, diaCode: code });
    this.props.onChange(code);
  }

  perform() {
    this.props.onChange(this.state.diaCode);
    this.setState({ isShowDialog: false });
    this.refreshEditor();
  }

  cancel() {
    this.setState({ code: this.props.value });
    this.setState({ isShowDialog: false });
  }

  render() {
    let { code, isShowDialog, isRefresh } = this.state;
    let {
      style = {},
      title = '',
      isShowHeader = false,
      isOpenDialogEditor = false,
      isOpenClip = false,
      className = '',
      defaultValue = '',
      value = '',
      onChange,
      mode = 'yaml',
      lineNumbers = true,
      theme = 'monokai',
      readOnly = false,
      isForceRefresh = false
    } = this.props;
    let options = {
      //mode,
      lineNumbers,
      theme,
      readOnly,
      styleActiveLine: true,
      matchBrackets: true
    };

    return (
      <div
        className="rich-textarea simple-mod"
        style={Object.assign({}, style, {
          lineHeight: '24px',
          paddingBottom: 0,
          overflowX: 'hidden'
        })}
      >
        {isShowHeader ? (
          <div className="permission-code-editor">
            <strong className="code-title">{title}</strong>
            <ul className="editor-toolbars">
              {!readOnly && isOpenDialogEditor ? (
                <li>
                  <Button
                    type="link"
                    tooltip="弹窗编辑"
                    onClick={() => {
                      this.setState({ isShowDialog: true });
                    }}
                  >
                    <i className="icon-enlarge" />
                  </Button>
                </li>
              ) : (
                <noscript />
              )}

              <li>
                <Clip isShowTip={true} tipDirection="top" target="#copy-area" />
              </li>
            </ul>
          </div>
        ) : (
          <noscript />
        )}

        <pre
          id="copy-area"
          style={{
            fontSize: '0px',
            padding: 0,
            margin: 0,
            width: 0,
            height: 0,
            border: 'none',
            position: 'fixed'
          }}
        >
          {value}
        </pre>

        <div className={'codemirror-theme-' + theme}>
          {isRefresh ? (
            <CodeMirror
              value={code}
              options={options}
              onBeforeChange={() => {}}
              scroll={{
                x: 0,
                y: Infinity
              }}
            />
          ) : (
            <div className={'CodeMirror cm-s-' + theme} />
          )}
        </div>

        <Modal
          visible={isShowDialog}
          caption="代码编辑"
          onClose={this.cancel.bind(this)}
          size={1000}
          disableEscape={true}
        >
          <Modal.Body>
            <div className={'codemirror-theme-' + theme} style={{ lineHeight: '24px' }}>
              <CodeMirror value={code} options={options} onBeforeChange={() => {}} />
            </div>
          </Modal.Body>
          <Modal.Footer>
            <Button className="m" onClick={this.perform.bind(this)}>
              提交
            </Button>
            <Button className="weak m" onClick={this.cancel.bind(this)}>
              取消
            </Button>
          </Modal.Footer>
        </Modal>
      </div>
    );
  }
}
