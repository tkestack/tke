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

import * as classNames from 'classnames';
import * as React from 'react';

import { AttributeValue } from './AttributeSelect';
import { Input } from './Input';
// import { OnOuterClick } from '@tencent/ff-redux';
import { Tag, TagValue } from './Tag';

export interface TagSearchBoxProps {
  attributes?: Array<AttributeValue>;
  defaultValue?: Array<any>;
  value?: Array<any>;
  minWidth?: number;
  onChange?: (tags: Array<any>) => void;
}

export interface TagSearchBoxState {
  active: boolean;
  dialogActive: boolean;
  curPos: number;
  curPosType: FocusPosType;
  inputValue: string;
  inputWidth: number;
  showSelect: boolean;
  tags: Array<TagValue>;
}

const keys = {
  '8': 'backspace',
  '9': 'tab',
  '13': 'enter',
  '32': 'spacebar',
  '37': 'left',
  '38': 'up',
  '39': 'right',
  '40': 'down'
};

const INPUT_MIN_SIZE: number = 5;

/**
 * 焦点所在位置类型
 */
export enum FocusPosType {
  INPUT,
  INPUT_EDIT,
  TAG
}

let cnt = 0;

export class TagSearchBox extends React.Component<TagSearchBoxProps, any> {
  static cnt: number = 0;

  state: TagSearchBoxState = {
    active: false,
    dialogActive: false,
    curPos: 0,
    curPosType: FocusPosType.INPUT,
    showSelect: true,
    inputValue: '',
    inputWidth: INPUT_MIN_SIZE,
    tags: this.props.defaultValue
      ? this.props.defaultValue.map(item => {
          item._key = cnt++;
          return item;
        })
      : []
  };

  componentDidMount() {
    if ('value' in this.props) {
      const value = this.props.value.map(item => {
        if (!('_key' in item)) {
          item._key = cnt++;
        }
        return item;
      });
      this.setState({ tags: value });
    }
  }

  componentWillReceiveProps(nextProps: TagSearchBoxProps) {
    if ('value' in nextProps) {
      const value = nextProps.value.map(item => {
        if (!('_key' in item)) {
          item._key = cnt++;
        }
        return item;
      });
      this.setState({ tags: value });
    }
  }

  open = () => {
    this.markTagElect(-1);
    const { active, tags } = this.state;
    if (!active) {
      this.setState({ active: true });
      // 展开时不激活select显示
      this.markTagElect(-1);
      this.setState({ curPosType: FocusPosType.INPUT, curPos: tags.length });
    } else {
      this.handleTagEvent('click-input', tags.length);
    }
    this.setState({ showSelect: true });
    // if (this[`tag-${tags.length}`].getInputValue().length > 0) {
    //   this.setState({ showSelect: true });
    // }
    setTimeout(() => {
      this[`tag-${tags.length}`].moveToEnd();
    }, 100);
  };

  // TODO:tea2.0
  // @OnOuterClick
  close() {
    // 编辑未完成的取消编辑
    const tags = this.state.tags.map((item, index) => {
      if (item._edit) {
        this[`tag-${index}`].editDone();
        item._edit = false;
      }
      return item;
    });

    this.setTags(
      tags,
      () => {
        this.markTagElect(-1);
        this.setState({ showSelect: false });
        if (this.state.active) {
          this.setState({ curPos: -1 }, () =>
            this.setState({ active: false }, () => {
              this[`search-box`].scrollLeft = 0;
            })
          );
        }
      },
      false
    );
  }

  notify = (tags: Array<TagValue>) => {
    const onChange = this.props.onChange;
    if (!onChange) return;

    const result = [];
    let keyReg = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/;
    let valueReg = /^([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$/;
    let keys = {};
    tags.forEach(item => {
      const attr = item.attr || null;
      const values = item.values;
      const parts = values.length > 0 ? values[0].name.split(':') : [];
      let illegal = false;
      keys[parts[0]] = keys[parts[0]] ? keys[parts[0]] + 1 : 1;
      if (keys[parts[0]] > 1) {
        illegal = true;
      } else if (parts.length === 2) {
        //判断字符是否合法
        illegal = !keyReg.test(parts[0]) || parts[0].lenth > 254 || !valueReg.test(parts[1]) || parts[1].length > 63;
      } else {
        illegal = true;
      }
      result.push({ attr, values, _key: item._key, _edit: item._edit, disable: item.disable, illegal });
    });
    onChange(result);
  };

  // Tags发生变动
  setTags(tags: Array<TagValue>, callback?: Function, notify = true): void {
    if (notify) this.notify(tags);
    this.setState({ tags }, () => {
      if (callback) callback();
    });
  }

  markTagElect(index: number): void {
    const tags = this.state.tags.map((item, i) => {
      if (index === i) {
        item._elect = true;
        this[`tag-${index}`].focusTag();
      } else {
        item._elect = false;
      }
      return item;
    });
    this.setState({ tags });
  }

  /**
   *  处理Tag相关事件
   */
  handleTagEvent = (type: string, oldIndex: number, payload?: any): void => {
    const { tags, active, curPos, curPosType } = this.state;

    let index = oldIndex;

    switch (type) {
      case 'add':
        this.markTagElect(-1);

        payload._key = cnt++;
        tags.splice(++index, 0, payload);
        this.setTags(tags, () => {
          this[`tag-${index}`].focusInput();
        });
        this.setState({ showSelect: false });
        break;

      case 'edit':
        this.markTagElect(-1);

        this[`tag-${index}`].editDone();
        tags[index].attr = payload.attr;
        tags[index].values = payload.values;
        tags[index]._edit = false;

        this.setTags(tags, () => {
          // this[`tag-${index-1}`].focusInput();
        });
        index++;
        this.setState({ showSelect: false, curPosType: FocusPosType.INPUT });
        break;

      case 'edit-cancel':
        this.markTagElect(-1);
        this[`tag-${index}`].editDone();

        this.setTags(
          tags,
          () => {
            // this[`tag-${index}`].focusInput();
          },
          false
        );
        this.setState({ showSelect: false, curPosType: FocusPosType.INPUT });
        break;

      case 'editing':
        if ('attr' in payload && tags[index]) tags[index].attr = payload.attr;
        if ('values' in payload && tags[index]) tags[index].values = payload.values;
        this.setTags(tags, null, false);
        break;

      case 'mark':
        if (index === tags.length) index--;
        if (index < 0 || !tags[index]) return;
        if (!tags[index]._elect) {
          this.markTagElect(index);
          this.setState({ curPosType: FocusPosType.TAG });
        }
        break;

      case 'del':
        if (payload === 'keyboard') {
          index--;
          if (tags[index].disable) return;
        }
        this.markTagElect(-1);
        if (!tags[index]) break;

        // 如果当前Tag的input中有内容，推向下个input
        // const curValue = this[`tag-${index}`].getInputValue();
        // this[`tag-${index}`].addTagByInputValue();

        tags.splice(index, 1);
        this.setTags(tags, () => {
          this.setState({ curPosType: FocusPosType.INPUT });
          // const input = this[`tag-${index-1 > 0 ? index-1 : 0}`];
          // input.setInputValue(curValue);
          // input.focusInput();
        });
        if (payload !== 'edit') {
          this.setState({ showSelect: false });
        }
        break;

      // payload 为点击位置
      case 'click':
        if (!active) {
          this.open();
          return;
        }
        // 触发修改
        // if (curPos === index && curPosType === FocusPosType.TAG) {
        tags[index]._edit = true;
        this.setTags(
          tags,
          () => {
            this.setState({ showSelect: true }, () => {
              this[`tag-${index}`].edit(payload);
            });
          },
          false
        );

        this.setState({ curPosType: FocusPosType.INPUT_EDIT });

        // } else {
        //   this.markTagElect(index);
        //   this.setState({ curPosType: FocusPosType.TAG });
        // }

        break;

      case 'click-input':
        this.markTagElect(-1);
        if (payload === 'edit') {
          this.setState({ curPosType: FocusPosType.INPUT_EDIT });
        } else {
          this.setState({ curPosType: FocusPosType.INPUT });
        }

        if (!active) {
          this.setState({ active: true });
        }
        this.setState({ showSelect: true });
        break;

      case 'move-left':
        // if (index <= 0) return;
        // if (index !== tags.length-1 || curPosType !== FocusPosType.INPUT) index--;
        // this.markTagElect(index);
        // this.setState({ curPosType: FocusPosType.TAG });
        // this[`tag-${index}`].focusTag();
        break;

      case 'move-right':
        // if (index >= tags.length - 1) return;
        // if (curPosType === FocusPosType.INPUT) {
        //   this.setState({ curPosType: FocusPosType.TAG });
        // }
        // // 到达最后input
        // if (index === tags.length-1) {
        //   this.markTagElect(-1);
        //   this[`tag-${index}`].focusInput();
        //   this.setState({ curPosType: FocusPosType.INPUT });
        // } else {
        //   this.markTagElect(index);
        //   this[`tag-${index}`].focusTag();
        // }
        // index++;
        break;
    }

    this.setState({ curPos: index });
  };

  render() {
    const { active, inputWidth, inputValue, tags, curPos, curPosType, dialogActive, showSelect } = this.state;
    const minWidth = this.props.minWidth;

    // 用于计算 focused 及 isFocused, 判断是否显示选择组件
    // (直接使用 Input 组件内部 onBlur 判断会使得 click 时组件消失)
    let focusedInputIndex = -1;
    if (curPosType === FocusPosType.INPUT || curPosType === FocusPosType.INPUT_EDIT) {
      focusedInputIndex = curPos;
    }

    const tagList = tags.map((item, index) => {
      const selectedAttrKeys = [];
      tags.forEach(tag => {
        if (tag.attr && item.attr && item._edit && item.attr.name === tag.attr.name) return null;
        if (tag.attr && tag.attr.key && !tag.attr.reusable) {
          selectedAttrKeys.push(tag.attr.key);
        }
      });

      const attributes = this.props.attributes.filter(item => selectedAttrKeys.indexOf(item.key) < 0);

      return (
        <Tag
          ref={tag => {
            this[`tag-${index}`] = tag;
          }}
          active={active}
          key={item._key}
          attributes={attributes}
          attr={item.attr}
          values={item.values}
          elect={item._elect}
          maxWidth={this['search-wrap'] ? this['search-wrap'].clientWidth : null}
          focused={focusedInputIndex === index && showSelect ? curPosType : null}
          dispatchTagEvent={(type, payload) => !item.disable && this.handleTagEvent(type, index, payload)}
          disable={item.disable}
          illegal={item.illegal}
        />
      );
    });

    const selectedAttrKeys = tags
      .map(item => (item.attr && !item.attr.reusable ? item.attr.key : null))
      .filter(item => !!item);
    const attributes = this.props.attributes.filter(item => selectedAttrKeys.indexOf(item.key) < 0);

    const minWidthStyle = active ? {} : { width: minWidth ? `${minWidth}px` : '100%' };

    // if (tagList.length <= 0) {
    tagList.push(
      <li key="100">
        <Input
          ref={input => {
            this[`tag-${tags.length}`] = input;
          }}
          active={active}
          maxWidth={this['search-wrap'] ? this['search-wrap'].clientWidth : null}
          attributes={attributes}
          isFocused={focusedInputIndex === tags.length && showSelect}
          dispatchTagEvent={(type, payload) => this.handleTagEvent(type, tags.length, payload)}
        />
      </li>
    );
    // }

    return (
      <div
        className="tc-select-tags-search-wrap"
        ref={div => {
          this['search-wrap'] = div;
        }}
      >
        <div
          className={classNames('tc-select-tags-search', { focus: active })}
          onClick={this.open}
          style={minWidthStyle}
          ref={div => {
            this[`search-box`] = div;
          }}
        >
          <div className="tc-search-wrap">
            <ul>{tagList}</ul>
          </div>
        </div>
      </div>
    );
  }
}
