import * as React from "react";
import * as classNames from "classnames";
import { AttributeSelect, AttributeValue } from "./AttributeSelect";
import { ValueSelect } from "./valueselect";
import { TagProps } from "./Tag";

export interface InputProps {
  /**
   * 触发标签相关事件
   */
  dispatchTagEvent: (type: string, payload?: any) => void;

  /**
   * 所有属性集合
   */
  attributes: Array<AttributeValue>;

  /**
   * 是否为 Focus 态
   */
  isFocused: boolean;

  /**
   * 搜索框是否处于展开状态
   */
  active: boolean;

  /**
   * 输入框类型（用于修改标签值的 Input type 为 "edit"）
   */
  type?: string;

  /**
   * 是否隐藏
   */
  hidden?: boolean;

  /**
   * 最大宽度
   */
  maxWidth: number;

  /**
   * 处理按键事件
   */
  handleKeyDown?: (e: any) => void;

  /**
   * 位置偏移
   */
  inputOffset?: number;
}

export interface InputState {
  inputWidth: number;
  inputValue: string;
  attribute: AttributeValue;
  values: Array<any>;
  showAttrSelect: boolean;
  showValueSelect: boolean;
  ValueSelectOffset: number;
}

const keys = {
  "8": "backspace",
  "9": "tab",
  "13": "enter",
  "37": "left",
  "38": "up",
  "39": "right",
  "40": "down"
};

const INPUT_MIN_SIZE = 0;

export class Input extends React.Component<InputProps, any> {
  state: InputState = {
    inputWidth: INPUT_MIN_SIZE,
    inputValue: "",
    attribute: null,
    values: [],
    showAttrSelect: false,
    showValueSelect: false,
    ValueSelectOffset: 0
  };

  constructor(props) {
    super(props);
  }

  componentDidMount() {}

  /**
   * 刷新选择组件显示
   */
  refreshShow = (): void => {
    const { inputValue, attribute } = this.state;

    const input = this["input"] as HTMLInputElement;
    let start = input.selectionStart,
      end = input.selectionEnd;

    // if (start !== end) {
    //   this.setState({ showAttrSelect: false, showValueSelect: false });
    //   return;
    // }

    const { pos } = this.getAttrStrAndValueStr(inputValue);

    if (pos < 0 || start <= pos) {
      this.setState({ showAttrSelect: true, showValueSelect: false });
      return;
    }

    if (attribute && end > pos) {
      this.setState({ showAttrSelect: false, showValueSelect: true });
    }
  };

  focusInput = (): void => {
    if (!this["input"]) return;
    const input = this["input"] as HTMLElement;
    input.focus();
  };

  moveToEnd = (): void => {
    const input = this["input"] as HTMLInputElement;
    input.focus();
    const value = this.state.inputValue;
    setTimeout(() => input.setSelectionRange(value.length, value.length));
  };

  selectValue = (): void => {
    const input = this["input"] as HTMLInputElement;
    input.focus();
    const value = this.state.inputValue;
    let { pos } = this.getAttrStrAndValueStr(value);
    if (pos < 0) pos = -2;
    setTimeout(() => {
      input.setSelectionRange(pos + 2, value.length);
      this.refreshShow();
    });
  };

  selectAttr = (): void => {
    const input = this["input"] as HTMLInputElement;
    input.focus();
    const value = this.state.inputValue;
    let { pos } = this.getAttrStrAndValueStr(value);
    if (pos < 0) pos = 0;
    setTimeout(() => {
      input.setSelectionRange(0, pos);
      this.refreshShow();
    });
  };

  setInfo(info: any, callback?: Function) {
    const attribute = info.attr;
    const values = info.values;
    this.setState({ attribute, values }, () => {
      if (attribute) {
        this.setInputValue(`${attribute.name}: ${values.map(item => item.name).join(" | ")}`, callback);
      } else {
        this.setInputValue(`${values.map(item => item.name).join(" | ")}`, callback);
      }
    });
  }

  setInputValue = (value: string, callback?: Function): void => {
    if (this.props.type === "edit" && value.trim().length <= 0) {
      this.props.dispatchTagEvent("del", "edit");
    }

    // value = value.replace(/：/g, ':');

    // const pos = value.indexOf(':');
    // let attrStr = value, valueStr = '';

    // if (pos >= 0) {
    //   attrStr = value.substr(0, pos);
    //   valueStr = value.substr(pos+1).replace(/^\s*(.*)/, '$1');
    // }

    const attributes = this.props.attributes;

    let attribute = null,
      valueStr = value;

    const input = this["input"] as HTMLElement;
    const mirror = this["input-mirror"] as HTMLElement;

    // attribute 是否存在
    for (let i = 0; i < attributes.length; ++i) {
      if (value.indexOf(attributes[i].name + ":") === 0 || value.indexOf(attributes[i].name + "：") === 0) {
        // 获取属性/值
        attribute = attributes[i];
        valueStr = value.substr(attributes[i].name.length + 1);

        // 计算 offset
        mirror.innerText = attribute.type === "onlyKey" ? attribute.name : attribute.name + ": ";
        let width = mirror.clientWidth;
        if (this.props.inputOffset) width += this.props.inputOffset;
        this.setState({ ValueSelectOffset: width });
        break;
      }
    }

    // 处理前导空格
    if (attribute && valueStr.replace(/^\s+/, "").length > 0) {
      value = `${attribute.name}: ${valueStr.replace(/^\s+/, "")}`;
    } else if (attribute) {
      value = attribute.type === "onlyKey" ? attribute.name : `${attribute.name}:${valueStr}`;
    }

    this.setState({ attribute }, this.refreshShow);

    if (this.props.type === "edit") {
      this.props.dispatchTagEvent("editing", { attr: attribute });
    }

    mirror.innerText = value;
    const width = mirror.clientWidth > INPUT_MIN_SIZE ? mirror.clientWidth : INPUT_MIN_SIZE;
    this.setState({ inputValue: value, inputWidth: width }, () => {
      if (callback) callback();
    });
  };

  resetInput = (callback?: Function): void => {
    this.setInputValue("", callback);
    this.setState({ inputWidth: INPUT_MIN_SIZE });
  };

  getInputValue = (): string => {
    return this.state.inputValue;
  };

  // getInputAttr = (): AttributeValue => {
  //   return this.state.attribute;
  // }

  addTagByInputValue = (): boolean => {
    const { attribute, values, inputValue } = this.state;
    const type = this.props.type || "add";
    // 属性值搜索
    if (attribute && this.props.attributes.filter(item => item.key === attribute.key).length > 0) {
      if (values.length <= 0) {
        return false;
      }
      this.props.dispatchTagEvent(type, { attr: attribute, values: values });
    } else {
      // 关键字搜索
      if (inputValue.trim().length <= 0) {
        return false;
      }
      let attribute = this.props.attributes.find(item => item.name === inputValue.trim());
      if (!attribute || attribute.type === "onlyKey") {
        this.props.dispatchTagEvent(type, { attr: attribute, values: [] });
        this.resetInput();
      } else {
        const list = inputValue
          .split("|")
          .filter(item => item.trim().length > 0)
          .map(item => {
            return { name: item.trim() };
          });
        this.props.dispatchTagEvent(type, { attr: null, values: list });
      }
    }
    this.setState({ showAttrSelect: false, showValueSelect: false });
    if (this.props.type !== "edit") {
      this.resetInput();
    }
    return true;
  };

  handleInputChange = (e): void => {
    this.setInputValue(e.target.value);
  };

  handleInputClick = (e): void => {
    this.props.dispatchTagEvent("click-input", this.props.type);
    e.stopPropagation();
    this.focusInput();
  };

  handleAttrSelect = (attr: AttributeValue): void => {
    if (attr && attr.key) {
      const str = attr.type === "onlyKey" ? attr.name : `${attr.name}: `;
      const inputValue = this.state.inputValue;
      if (inputValue.indexOf(str) >= 0) {
        this.selectValue();
      } else {
        this.setInputValue(str);
      }
      this.setState({ values: [] });
    }

    if (attr.type === "onlyKey") {
      // 不需要值
      // this.addTagByInputValue()
      setTimeout(() => this.addTagByInputValue());
    } else {
      this.focusInput();
    }
  };

  handleValueChange = (values: Array<any>): void => {
    this.setState({ values }, () => {
      this.setInputValue(`${this.state.attribute.name}: ${values.map(item => item.name).join(" | ")}`);
      this.focusInput();
    });
  };

  /**
   * 值选择组件完成选择
   */
  handleValueSelect = (values: Array<any>): void => {
    this.setState({ values });
    const inputValue = this.state.inputValue;

    if (values.length <= 0) {
      this.setInputValue(this.state.attribute.name + ": ");
      return;
    }

    if (values.length > 0) {
      const key = this.state.attribute.key;
      if (this.props.attributes.filter(item => item.key === key).length > 0) {
        const type = this.props.type || "add";
        this.props.dispatchTagEvent(type, { attr: this.state.attribute, values });
      }
      this.focusInput();
    }

    if (this.props.type !== "edit") {
      this.resetInput();
    }
  };

  /**
   * 值选择组件取消选择
   */
  handleValueCancel = () => {
    if (this.props.type === "edit") {
      const { attribute, values } = this.state;
      this.props.dispatchTagEvent("edit-cancel", { attr: attribute, values: values });
    } else {
      this.resetInput(() => {
        this.focusInput();
      });
    }
  };

  /**
   * 处理粘贴事件
   */
  handlePaste = (e): void => {
    const { attribute } = this.state;

    if (!attribute || attribute.type === "input") {
      this["textarea"].focus();
      setTimeout(() => {
        let value = this["textarea"].value;

        if (/^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/.test(value)) {
          value = value.replace(/[\r\n\t,，\s]+/g, "|");
        } else {
          value = value.replace(/[\r\n\t,，]+/g, "|");
        }
        value = value
          .split("|")
          .map(item => item.trim())
          .filter(item => item.length > 0)
          .join(" | ");

        const input = this["input"] as HTMLInputElement;
        const start = input.selectionStart,
          end = input.selectionEnd;
        const inputValue = this.state.inputValue;
        // 覆盖选择区域
        const curValue = inputValue.substring(0, start) + value + inputValue.substring(end, inputValue.length);

        // input 属性情况
        this["textarea"].value = "";
        if (attribute && attribute.type === "input") {
          this.setInputValue(curValue, this.focusInput);
          return;
        }

        if (inputValue.length > 0) {
          this.setInputValue(curValue, this.focusInput);
        } else {
          this.setInputValue(curValue, this.addTagByInputValue);
        }
      }, 100);
    }
  };

  // 键盘事件
  // handlekeyUp = (e): void => {
  //   if (this['value-select']) {
  //     if (this['value-select'].handlekeyUp(e.keyCode) === false) return;
  //   }
  // }

  handlekeyDown = (e): void => {
    if (!keys[e.keyCode]) return;

    // if (!this.props.isFocused) {
    //   this.props.dispatchTagEvent('click-input', this.props.type);
    // }

    if (this.props.hidden) {
      return this.props.handleKeyDown(e);
    }

    const inputValue = this.state.inputValue;

    if (keys[e.keyCode] === "backspace" && inputValue.length > 0) return;

    if ((keys[e.keyCode] === "left" || keys[e.keyCode] === "right") && inputValue.length > 0) {
      setTimeout(this.refreshShow, 0);
      return;
    }

    e.preventDefault();

    // 事件下传
    if (this["attr-select"]) {
      if (this["attr-select"].handleKeyDown(e.keyCode) === false) return;
    }
    if (this["value-select"]) {
      if (this["value-select"].handleKeyDown(e.keyCode) === false) return;
    }

    switch (keys[e.keyCode]) {
      case "enter":
      case "tab":
        if (!this.props.isFocused) {
          this.props.dispatchTagEvent("click-input");
        }
        this.addTagByInputValue();
        break;

      case "backspace":
        this.props.dispatchTagEvent("del", "keyboard");
        break;

      case "up":
        break;

      case "down":
        break;

      case "left":
        this.props.dispatchTagEvent("move-left");
        break;

      case "right":
        this.props.dispatchTagEvent("move-right");
        break;
    }
  };

  getAttrStrAndValueStr = (str: string): any => {
    let attrStr = str,
      valueStr = "",
      pos = -1;

    const attributes = this.props.attributes;
    for (let i = 0; i < attributes.length; ++i) {
      if (str.indexOf(attributes[i].name + ":") === 0) {
        // 获取属性/值
        attrStr = attributes[i].name;
        valueStr = str.substr(attrStr.length + 1);
        pos = attributes[i].name.length;
      }
    }

    return { attrStr, valueStr, pos };
  };

  render() {
    const { inputWidth, inputValue, showAttrSelect, showValueSelect, attribute, ValueSelectOffset } = this.state;
    const { active, attributes, isFocused, hidden, maxWidth, type } = this.props;

    // const pos = inputValue.indexOf(':');
    // let attrStr = inputValue, valueStr = '';
    // if (pos >= 0) {
    //   attrStr = inputValue.substr(0, pos).trim();
    //   valueStr = inputValue.substr(pos+1).trim();
    // }

    const { attrStr, valueStr } = this.getAttrStrAndValueStr(inputValue);

    const attrSelect =
      isFocused && showAttrSelect ? (
        <AttributeSelect
          ref={select => (this["attr-select"] = select)}
          attributes={attributes}
          inputValue={attrStr}
          onSelect={this.handleAttrSelect}
        />
      ) : null;

    const valueSelect =
      isFocused && showValueSelect && attribute && attribute.type ? (
        <ValueSelect
          type={attribute.type}
          ref={select => (this["value-select"] = select)}
          values={attribute.values}
          inputValue={valueStr.trim()}
          offset={ValueSelectOffset}
          onChange={this.handleValueChange}
          onSelect={this.handleValueSelect}
          onCancel={this.handleValueCancel}
        />
      ) : null;

    const style = {
      width: hidden ? "0px" : active ? `${inputWidth + 5}px` : "5px",
      maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px"
    };
    if (type === "edit" && !hidden) {
      style["padding"] = "0 8px";
    }

    const input =
      type !== "edit" ? (
        <input
          ref={input => (this["input"] = input)}
          type="text"
          className="tc-search-input"
          placeholder=""
          style={{
            width: hidden ? "0px" : `${inputWidth + 5}px`,
            display: active ? "" : "none",
            maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px"
          }}
          value={inputValue}
          onChange={this.handleInputChange}
          onKeyDown={this.handlekeyDown}
          onFocus={this.refreshShow}
          onClick={this.refreshShow}
          onPaste={this.handlePaste}
        />
      ) : (
        <div style={{ position: "relative", display: hidden ? "none" : "" }}>
          <pre style={{ display: "block", visibility: "hidden" }}>
            <div
              style={{
                fontSize: "12px",
                width: hidden ? "0px" : `${inputWidth + 36}px`,
                maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px",
                whiteSpace: "normal"
              }}
            >
              {inputValue}
            </div>
            <br style={{ clear: "both" }} />
          </pre>
          <textarea
            ref={input => (this["input"] = input)}
            className="tc-search-input"
            placeholder=""
            style={{
              width: hidden ? "0px" : `${inputWidth + 30}px`,
              display: active ? "" : "none",
              maxWidth: maxWidth ? `${maxWidth - 36}px` : "435px",
              position: "absolute",
              top: 0,
              left: 0,
              height: "100%",
              resize: "none",
              minHeight: "15px",
              marginTop: "2px"
            }}
            value={inputValue}
            onChange={this.handleInputChange}
            onKeyDown={this.handlekeyDown}
            onFocus={this.refreshShow}
            onClick={this.refreshShow}
            onPaste={this.handlePaste}
          />
        </div>
      );

    return (
      <div className="tc-tags-space" style={style} onClick={this.handleInputClick}>
        {input}
        <span
          ref={input => (this["input-mirror"] = input)}
          style={{ position: "absolute", top: "-9999px", left: 0, whiteSpace: "pre", fontSize: "12px" }}
        />
        <textarea
          ref={textarea => (this["textarea"] = textarea)}
          style={{ position: "absolute", top: "-9999px", left: 0, whiteSpace: "pre", fontSize: "12px" }}
        />
        {attrSelect}
        {valueSelect}
      </div>
    );
  }
}
