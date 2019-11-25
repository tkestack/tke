import * as React from 'react';
import * as classNames from 'classnames';
import { Input } from './Input';
import { AttributeValue } from './AttributeSelect';
import { FocusPosType } from './TagSearchBox';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface TagValue {
  _key: string;
  _edit?: boolean;
  // 当前是否选中
  _elect?: boolean;
  attr?: AttributeValue;
  values?: Array<any>;
  disable?: boolean;
  illegal?: boolean;
}

export interface TagProps {
  attr?: AttributeValue;
  values?: Array<any>;
  // 当前是否选中
  elect?: boolean;
  dispatchTagEvent?: (type: string, payload?: any) => void;
  attributes: Array<AttributeValue>;
  focused: FocusPosType;
  maxWidth?: number;
  active: boolean;
  disable?: boolean;
  illegal?: boolean;
}

export interface TagState {
  bubbleActice: boolean;
  inEditing: boolean;
  inputOffset: number;
  // attribute: AttributeValue;
  // selectValues: Array<any>;
}

const keys = {
  '8': 'backspace',
  '13': 'enter',
  '37': 'left',
  '38': 'up',
  '39': 'right',
  '40': 'down'
};

const INPUT_MIN_SIZE: number = 0;

export class Tag extends React.Component<TagProps, any> {
  state: TagState = {
    bubbleActice: false,
    inEditing: false,
    inputOffset: 0
  };

  componentDidMount() {
    // this.setState({inputOffset: this['content'].clientWidth});
  }

  componentWillReceiveProps() {
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

  addTagByInputValue(): boolean {
    return this['input'].addTagByInputValue();
  }

  addTagByEditInputValue(): boolean {
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
        return this['input-inside'].selectAttr();
      } else {
        return this['input-inside'].selectValue();
      }
    });
  }

  editDone(): void {
    this.setState({ inEditing: false });
  }

  handleTagClick = (e, pos?: string): void => {
    this.props.dispatchTagEvent('click', pos);
    e.stopPropagation();
  };

  handleDelete = e => {
    e.stopPropagation();
    this.props.dispatchTagEvent('del');
  };

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
  };

  render() {
    const { bubbleActice, inEditing, inputOffset } = this.state;
    const {
      active,
      attr,
      values,
      elect,
      dispatchTagEvent,
      attributes,
      focused,
      maxWidth,
      disable,
      illegal
    } = this.props;

    let attrStr = attr ? attr.name : '';
    if (attr && attr.name) {
      attrStr += ': ';
    }
    let valueStr = values.map(item => item.name).join(' | ');

    const itemStyle = inEditing && !active ? { width: 0, minHeight: '20px' } : { minHeight: '20px' };
    const tagStyle = { display: inEditing ? 'none' : '', borderColor: illegal ? '#f00' : '#ddd' };
    const tips =
      disable === true
        ? t('该标签为系统标签，不可进行操作')
        : illegal
        ? t('该标签格式有误或重复，将不会被保存')
        : t('点击进行修改，按回车键完成修改');
    return (
      <li style={itemStyle}>
        <div
          className="tc-15-bubble tc-15-bubble-bottom black"
          style={{ display: bubbleActice ? '' : 'none', top: '-70px' }}
        >
          <div className="tc-15-bubble-inner">{tips}</div>
        </div>

        <div
          className={classNames('tc-tags', { current: elect })}
          ref={div => {
            this['content'] = div;
          }}
          style={tagStyle}
          onClick={this.handleTagClick}
          onMouseOver={() => this.setState({ bubbleActice: true })}
          onMouseOut={() => this.setState({ bubbleActice: false })}
        >
          <span onClick={e => this.handleTagClick(e, 'attr')}>{attrStr}</span>
          <span onClick={e => this.handleTagClick(e, 'value')}>{valueStr}</span>
          <a href="javascript:;" className="tc-tags-close-btn" onClick={this.handleDelete}>
            <i className="clear-icon" />
          </a>
        </div>
        <Input
          type="edit"
          hidden={!inEditing}
          maxWidth={maxWidth}
          handleKeyDown={this.handleKeyDown}
          active={active}
          ref={input => {
            this['input-inside'] = input;
          }}
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
    );
  }
}
