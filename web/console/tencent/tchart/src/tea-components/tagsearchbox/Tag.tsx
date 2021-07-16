import * as React from 'react';
import classNames from 'classnames';
import { Input } from './Input';
import { AttributeValue } from './AttributeSelect';
import { FocusPosType } from './TagSearchBox';

export interface TagValue {

  /**
   * 标签标识
   */
  _key: string;

  /**
   * 当前是否在编辑状态
   */
  _edit?: boolean;

  /**
   * 当前是否选中
   */
  _elect?: boolean;

  /**
   * 标签属性
   */
  attr?: AttributeValue;

  /**
   * 标签属性值
   */
  values?: Array<any>;
}

export interface TagProps {

  /**
   * 标签属性
   */
  attr?: AttributeValue;

  /**
   * 标签属性值
   */
  values?: Array<any>;

  /**
   * 当前是否选中
   */
  elect?: boolean;

  /**
   * 触发标签相关事件
   */
  dispatchTagEvent?: (type: string, payload?: any) => void;

  /**
   * 所有属性集合
   */
  attributes: Array<AttributeValue>;

  /**
   * 当前聚焦状态
   */
  focused: FocusPosType;

  /**
   * 最大长度
   */
  maxWidth?: number;

  /**
   * 搜索框是否处于展开状态
   */
  active: boolean;
}

export interface TagState {
  bubbleActice: boolean;
  inEditing: boolean;
  inputOffset: number;
  // attribute: AttributeValue;
  // selectValues: Array<any>;
}

const keys = {
  "8" : 'backspace',
  "13": 'enter',
  "37": 'left',
  "38": 'up',
  "39": 'right',
  "40": 'down'
};

const INPUT_MIN_SIZE: number = 0;

export class Tag extends React.Component<TagProps, any> {

  state: TagState = {
    bubbleActice: false,
    inEditing: false,
    inputOffset: 0
  }

  componentDidMount() {
    // console.log(this['content'].clientWidth);
    // this.setState({inputOffset: this['content'].clientWidth});
  }

  componentWillReceiveProps() {
    // console.log(this['content'].clientWidth);
    // this.setState({inputOffset: this['content'].clientWidth});
  }

  focusTag(): void {
    (this['input-inside'] as any).focusInput();
  }

  focusInput(): void {
    (this['input'] as any).focusInput();
  }

  resetInput() {
    (this['input-inside'] as any).resetInput();
  }

  setInputValue(value: string, callback?: Function): void {
    (this['input'] as any).setInputValue(value, callback);
  }

  getInputValue(): string {
    return this['input'].getInputValue();
  }

  addTagByInputValue() :boolean {
    return this['input'].addTagByInputValue();
  }

  addTagByEditInputValue() :boolean {
    if (!this['input-inside']) return;
    return this['input-inside'].addTagByInputValue();
  }
  
  setInfo(info: any, callback?: Function): void {
    return this['input'].setInfo(info, callback);
  }

  moveToEnd(): void {
    return this['input'].moveToEnd();
  }

  // getInputAttr = (): AttributeValue => {
  //   return this['input'].getInputAttr();
  // }

  getInfo(): any {
    let { attr, values } = this.props;
    return { attr, values };
  }

  edit(pos: string): void {
    this.setState({ inEditing: true });
    this['input-inside'].setInfo(this.getInfo(), () => {
      if (pos === 'attr') {
        return this['input-inside'] && this['input-inside'].selectAttr();
      } else {
        return this['input-inside'] && this['input-inside'].selectValue();
      }
    });
  }

  editDone(): void {
    this.setState({ inEditing: false });
  }

  handleTagClick = (e, pos?: string): void => {
    this.props.dispatchTagEvent('click', pos);
    e.stopPropagation();
  }

  handleDelete= (e) => {
    e.stopPropagation();
    this.props.dispatchTagEvent('del');
  }

  handleKeyDown = (e): void => {
    if (!keys[e.keyCode]) return;

    e.preventDefault();

    switch (keys[e.keyCode]) {
      case 'tab':
      case 'enter':
        this.props.dispatchTagEvent('click', 'value');
        break;
      
      case 'backspace':
        this.props.dispatchTagEvent('del', 'keyboard');
        break;
      
      case 'left':
        this.props.dispatchTagEvent('move-left');
        break;

      case 'right':
        this.props.dispatchTagEvent('move-right');
        break;
    }
    
  }

  render() {
    const { bubbleActice, inEditing, inputOffset } = this.state;
    const { active, attr, values, elect, dispatchTagEvent, attributes, focused, maxWidth } = this.props;

    let attrStr = attr ? attr.name : '';
    if (attr && attr.name && attr.type !== 'onlyKey') {
      attrStr += ': ';
    }
    let valueStr = values.map(item => item.name).join(' | ');

    const itemStyle = (inEditing && !active) ? { width: 0, minHeight: '20px' } : { minHeight: '20px' };

    const removeable = (attr && 'removeable' in attr) ? attr.removeable : true;

    return (
      <li style={itemStyle}>
        <div className="tc-15-bubble tc-15-bubble-bottom black" style={{display: bubbleActice ? '' : 'none'}}>
          <div className="tc-15-bubble-inner">
            点击进行修改，按回车键完成修改
          </div>
        </div>
        
        <div
          className={classNames("tc-tags", { "current": elect })}
          ref={div => this['content'] = div}
          style={{ display: inEditing ? 'none' : '', paddingRight: removeable ? '' : '8px', maxWidth: 'none'}}
          onClick={this.handleTagClick}
          onMouseOver={() => this.setState({bubbleActice : true})}
          onMouseOut={() => this.setState({bubbleActice : false})}
        >
          <span onClick={(e) => this.handleTagClick(e, 'attr')}>{attrStr}</span>
          <span onClick={(e) => this.handleTagClick(e, 'value')}>{valueStr}</span>
          { removeable &&
            <a href="javascript:;" className="tc-tags-close-btn" onClick={this.handleDelete}><i className="clear-icon"></i></a>
          }
        </div>
        <Input
          type="edit"
          hidden={!inEditing}
          maxWidth={maxWidth}
          handleKeyDown={this.handleKeyDown}
          active={active}
          ref={input => this["input-inside"]=input}
          attributes={attributes}
          dispatchTagEvent={dispatchTagEvent}
          isFocused={focused === FocusPosType.INPUT_EDIT}
        />
        {/*<Input
          active={active}
          maxWidth={maxWidth}
          inputOffset={inputOffset}
          ref={input => this["input"]=input}
          attributes={attributes}
          dispatchTagEvent={dispatchTagEvent}
          isFocused={focused === FocusPosType.INPUT}
        />*/}
      </li>
    )
  }
}