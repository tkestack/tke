import * as React from 'react';
import classNames from 'classnames';
import { OnOuterClick } from "tea-components/libs/decorators/OnOuterClick";
import { Tag, TagValue } from './Tag';
import { Input } from './Input';
import { AttributeValue } from './AttributeSelect';
export interface TagSearchBoxProps {

  editable?: boolean;

  style?
  /**
   * 要选择过滤的资源属性的集合
   */
  attributes?: Array<AttributeValue>;

  /**
   * 搜索框中默认包含的标签值的集合
   */
  defaultValue?: Array<any>;

  /**
   * 配合 onChange 作为受控组件使用
   */
  value?: Array<any>;

  /**
   * 搜索框收起后宽度，单位为 px
   */
  minWidth?: number | string;

  /**
   * 当搜索框中新增、修改或减少标签时调用此函数
   */
  onChange?: (tags: Array<any>) => void;

  /**
   * 搜索框中提示语（中）
   */
  tipZh?: string;

  /**
   * 搜索框中提示语（英）
   */
  tipEn?: string;
}

export interface TagSearchBoxState {

  /**
   * 搜索框是否为展开状态
   */
  active: boolean;

  /**
   * 是否展示提示框
   */
  dialogActive: boolean;

  /**
   * 当前光标位置
   */
  curPos: number;

  /**
   * 当前光标（焦点）所在位置的元素类型
   */
  curPosType: FocusPosType;

  /**
   * 输入框值
   */
  inputValue: string;

  /**
   * 输入框宽度
   */
  inputWidth: number;

  /**
   * 是否展示值选择组件
   */
  showSelect: boolean;

  /**
   * 已选标签
   */
  tags: Array<TagValue>;
}

const keys = {
  "8": 'backspace',
  "9": 'tab',
  "13": 'enter',
  "32": 'spacebar',
  "37": 'left',
  "38": 'up',
  "39": 'right',
  "40": 'down'
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
    tags: this.props.defaultValue ? this.props.defaultValue.map(item => {
      item._key = TagSearchBox.cnt++;
      return item;
    }) : []
  }

  static defaultProps = {
    tipZh: "多个关键字用竖线“ | ”分隔，多个过滤标签用回车键分隔",
    tipEn: "Use '|' to split more than one keyword, and press Enter to split tags",
    minWidth: 210,
    editable: true,
  }

  componentDidMount() {
    if ('value' in this.props) {
      const value = this.props.value.map(item => {
        if (!('_key' in item)) {
          item._key = TagSearchBox.cnt++;
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
          item._key = TagSearchBox.cnt++;
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
      this[`tag-${tags.length}`] && this[`tag-${tags.length}`].moveToEnd();
    }, 100);
  }

  @OnOuterClick
  close() {
    // 编辑未完成的取消编辑
    const tags = this.state.tags.map((item, index) => {
      if (item._edit) {
        this[`tag-${index}`].editDone();
        item._edit = false;
      }
      return item;
    });

    this.setTags(tags, () => {
      this.markTagElect(-1);
      this.setState({ showSelect: false });
      if (this.state.active) {
        this.setState({ curPos: -1 }, () =>
          this.setState({ active: false }, () => this[`search-box`].scrollLeft = 0)
        );
      }
    }, false);
  }


  notify = (tags: Array<TagValue>) => {
    const onChange = this.props.onChange;
    if (!onChange) return;

    const result = new Array();
    tags.forEach(item => {
      const attr = item.attr || null;
      const values = item.values;
      if (attr && attr.type === 'onlyKey' || values.length > 0) {
        result.push({ attr, values, _key: item._key, _edit: item._edit });
      }
    });
    onChange(result);
  }

  // Tags发生变动
  setTags(tags: Array<TagValue>, callback?: Function, notify = true): void {
    if (notify) this.notify(tags);
    this.setState({ tags }, () => { if (callback) callback() });
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
   * 点击清除按钮触发事件
   */
  handleClean = (e): void => {
    e.stopPropagation();
    if (this.state.tags.length <= 0) {
      this[`tag-${0}`].setInputValue('');
      return;
    }
    this.setTags([], () => {
      this[`tag-${0}`].setInputValue('');
      this[`tag-${0}`].focusInput();
      // this.handleTagEvent('click-input', 0);
    });
    this.setState({ curPos: 0, curPosType: FocusPosType.INPUT });
  }


  /**
   * 点击帮助触发事件
   */
  handleHelp = (e) => {
    e.stopPropagation();
    this.setState({ dialogActive: true });
  }


  /**
   * 点击搜索触发事件
   */
  handleSearch = (e) => {
    if (!this.state.active) return;
    e.stopPropagation();
    const { curPos, curPosType, tags } = this.state;
    let flag = false;
    // if (curPosType === FocusPosType.INPUT && this[`tag-${curPos}`].addTagByInputValue()) flag = true;
    const input = this[`tag-${tags.length}`];
    if (input && input.addTagByInputValue) {
      if (input.addTagByInputValue()) {
        flag = true;
      }
    }
    // if (curPosType === FocusPosType.INPUT_EDIT && this[`tag-${curPos}`].addTagByEditInputValue()) flag = true;

    for (let i = 0; i < tags.length; ++i) {
      if (!this[`tag-${i}`] || !this[`tag-${i}`].addTagByEditInputValue) return;
      if (tags[i]._edit && this[`tag-${i}`].addTagByEditInputValue()) flag = true;
    }

    if (flag) return;

    this.notify(this.state.tags);
    input.focusInput();
  }

  /**
   *  处理Tag相关事件
   */
  handleTagEvent = (type: string, index: number, payload?: any): void => {
    const { tags, active, curPos, curPosType } = this.state;
    // console.log(type, index, payload);
    switch (type) {

      case 'add':
        this.markTagElect(-1);

        payload._key = TagSearchBox.cnt++;
        tags.splice(++index, 0, payload);
        this.setTags(tags, () => {
          if (this[`tag-${index}`]) {
            this[`tag-${index}`].focusInput();
          }
        });
        this.setState({ showSelect: false });
        break;

      case 'edit':
        this.markTagElect(-1);

        this[`tag-${index}`] && this[`tag-${index}`].editDone();
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

        this.setTags(tags, () => {
          // this[`tag-${index}`].focusInput();
        }, false);
        this.setState({ showSelect: false, curPosType: FocusPosType.INPUT });
        break;

      case 'editing':
        if (('attr' in payload) && tags[index]) tags[index].attr = payload.attr;
        if (('values' in payload) && tags[index]) tags[index].values = payload.values;
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
        if (payload === 'keyboard') index--;
        this.markTagElect(-1);
        if (!tags[index]) break;

        // 如果当前Tag的input中有内容，推向下个input
        // const curValue = this[`tag-${index}`].getInputValue();
        // this[`tag-${index}`].addTagByInputValue();

        // 检查不可移除
        const { attr } = tags[index];
        if (attr && 'removeable' in attr && attr.removeable === false) {
          break;
        }

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

        const pos = payload;
        tags[index]._edit = true;
        this.setTags(tags, () => {
          this.setState({ showSelect: true }, () => {
            this[`tag-${index}`].edit(pos);
          });
        }, false);

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
  }


  render() {

    const { active, inputWidth, inputValue, tags, curPos, curPosType, dialogActive, showSelect } = this.state;
    const { minWidth, tipZh, tipEn, attributes } = this.props;

    // 用于计算 focused 及 isFocused, 判断是否显示选择组件
    // (直接使用 Input 组件内部 onBlur 判断会使得 click 时组件消失)
    let focusedInputIndex = -1;
    if (curPosType === FocusPosType.INPUT || curPosType === FocusPosType.INPUT_EDIT) {
      focusedInputIndex = curPos;
    }
    const tagList = tags.map((item, index) => {

      // 补全 attr 属性
      attributes.forEach(attrItem => {
        if (item.attr && attrItem.key && attrItem.key == item.attr.key) {
          item.attr = Object.assign({}, item.attr, attrItem);
        }
      });

      const selectedAttrKeys = [];
      tags.forEach(tag => {
        if (tag.attr && item.attr && item._edit && item.attr.key === tag.attr.key) return null;
        if (tag.attr && tag.attr.key && !tag.attr.reusable) {
          selectedAttrKeys.push(tag.attr.key);
        }
      })

      const useableAttributes = attributes.filter(item => selectedAttrKeys.indexOf(item.key) < 0);

      return <Tag
        ref={tag => this[`tag-${index}`] = tag}
        active={active}
        key={item._key}
        attributes={useableAttributes}
        attr={item.attr}
        values={item.values}
        elect={item._elect}
        maxWidth={this['search-wrap'] ? this['search-wrap'].clientWidth : null}
        focused={(focusedInputIndex === index && showSelect) ? curPosType : null}
        dispatchTagEvent={(type, payload) => this.props.editable && this.handleTagEvent(type, index, payload)} />
    });


    const selectedAttrKeys = tags.map(item => item.attr && !item.attr.reusable ? item.attr.key : null).filter(item => !!item);
    const useableAttributes = attributes.filter(item => selectedAttrKeys.indexOf(item.key) < 0);

    const minWidthStyle = active ? {} : ({ width: minWidth ? minWidth : '100%' });

    if (this.props.editable) {
      tagList.push(
        <li key='100'>
          <Input
            ref={input => this[`tag-${tags.length}`] = input}
            active={active}
            maxWidth={this['search-wrap'] ? this['search-wrap'].clientWidth : null}
            attributes={useableAttributes}
            isFocused={focusedInputIndex === tags.length && showSelect}
            dispatchTagEvent={(type, payload) => this.handleTagEvent(type, tags.length, payload)}
          />
        </li>
      )
    }

    const tip = window['VERSION'] === 'en' ? tipEn : tipZh;

    return (
      <div className="tc-select-tags-search-wrap" style={this.props.style} ref={div => this['search-wrap'] = div}>
        <div
          className={classNames("tc-select-tags-search", { "focus": this.props.editable && active })}
          onClick={this.open}
          style={minWidthStyle}
          ref={div => this[`search-box`] = div}
        >
          <div className="tc-search-wrap">
            <ul>
              {tagList}
            </ul>
            {this.props.editable && tagList.length === 1 && <span className="help-tips" style={{ lineHeight: '28px' }}>
              <p className="text text-weak">{tip}</p>
            </span>}
            {/*
            <a href="javascript:;" className="tc-icon-btn clear-btn" onClick={this.handleClean}>
              <div className="tc-15-bubble-icon">
                <i className="clear-icon"></i>
                <div className="tc-15-bubble tc-15-bubble-bottom black">
                  <div className="tc-15-bubble-inner">清空</div>
                </div>
              </div>
            </a>
            <a href="javascript:;" className="tc-icon-btn plaint-btn" onClick={this.handleHelp}>
              <div className="tc-15-bubble-icon">
                <i className="plaint-icon"></i>
                <div className="tc-15-bubble tc-15-bubble-bottom black">
                  <div className="tc-15-bubble-inner">帮助</div>
                </div>
              </div>
            </a>
            <a href="javascript:;" className="tc-icon-btn search-btn" onClick={this.handleSearch}>
              <div className="tc-15-bubble-icon">
                <i className="icon-search"></i>
                <div className="tc-15-bubble tc-15-bubble-bottom black">
                  <div className="tc-15-bubble-inner">搜索</div>
                </div>
              </div>
            </a> */}
          </div>
        </div>

        {/* <div className="dialog-panel" style={{display: dialogActive ? '': 'none'}}>
          <div className="tc-15-rich-dialog m" role="alertdialog" style={{width: "960px", margin: "0 auto", marginTop: "100px"}}>
            <div className="tc-15-rich-dialog-hd"><strong>搜索帮助</strong>
              <button title="关闭" className="tc-15-btn-close" onClick={() => {this.setState({dialogActive: false})}}>关闭</button>
            </div>
            <div className="tc-15-rich-dialog-bd">
              <i className="search-help-dialog"></i>
            </div>
          </div>
        </div> */}
      </div>
    )
  }
}
