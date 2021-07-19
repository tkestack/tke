import * as React from 'react';
import { SelectItem } from '../select';
import { OnOuterClick } from '../libs/decorators/OnOuterClick';


interface BaseReactProps {
  key?: string;
  defaultValue?: string;
  children?: React.ReactNode;
  className?: string;
  placeholder?: string;
  style?: object;
}

interface MultiDropdownProps extends BaseReactProps {
  disabled?: boolean;
  value?: string[];
  options: SelectItem[];
  mode?: 'single' | 'multi';
  onItemChange?: Function;
  onChange: Function;
}

interface MultiDropdownState {
  options: SelectItem[];
  tmpList: SelectItem[];
  isOpened: boolean;
  keyword: string;
}

export class MultiDropdown extends React.Component<MultiDropdownProps, MultiDropdownState> {
  constructor(props) {
    super(props);
    const {options = [], value = []} = props;
    const defaultSelectedList = options.filter(item => value.indexOf(item.value) > - 1);
    this.state = {
      options,
      tmpList: defaultSelectedList,
      isOpened: false,
      keyword: '',
    };
  }

  @OnOuterClick
  fun() {
    this.setState({
      isOpened: false,
    });
  }

  handleCancel() {
    this.setState({
      isOpened: false,
      tmpList: [],
    });
  }

  handleChange(item) {
    const tmpList = [item];
    this.setState(
      {
        tmpList,
        isOpened: false,
      },
      () => this.props.onChange(tmpList.map(item => item.value)),
    );
  }

  handleSubmit() {
    const {tmpList} = this.state;
    this.setState(
      {
        isOpened: false,
      },
      () => this.props.onChange(tmpList.map(item => item.value)),
    );
  }

  static getDerivedStateFromProps(nextProps, prevState) {
    const {disabled} = nextProps;
    if (disabled) {
      return {
        tmpList: [],
      };
    }

    const {options} = nextProps;
    const {tmpList} = prevState;
    const enabledList = tmpList.filter(item =>
      options.find(obj => obj.value === item.value && ! obj.disabled));
    if (options !== prevState.options) {
      return {
        tmpList: enabledList,
      };
    }
    return null;
  }

  changeSelected(e, item) {
    const {checked} = e.target;
    const {tmpList = []} = this.state;
    const index = tmpList.findIndex(v => item.value === v.value);
    if (this.props.mode === 'single') {
      this.setState({
        tmpList: checked ? tmpList.slice(index, 1) : [],
      });
      return;
    }

    if (checked && index < 0) {
      tmpList.push(item);
    }
    else if (! checked && index >= 0) {
      tmpList.splice(index, 1);
    }

    const {onItemChange} = this.props;
    onItemChange && onItemChange(tmpList);

    this.setState({
      tmpList,
    });
  }

  changeAllSelected(e) {
    const {checked} = e.target;
    const {options} = this.props;
    const {onItemChange} = this.props;
    const tmpList = checked ? options.filter(item => ! item.disabled) : [];
    onItemChange && onItemChange(tmpList);
    this.setState({
      tmpList,
    });
  }

  async toggleMenu() {
    if (! this.props.disabled) {
      const {isOpened} = this.state;
      this.setState({
        isOpened: ! isOpened,
      });
    }
  }

  render() {
    const {
      mode, value, disabled, children, options, onItemChange,
      onChange, className, ...res
    } = this.props;
    const {tmpList = [], isOpened, keyword} = this.state;
    const list = options.filter(item => item.label.indexOf(keyword) > - 1);
    const showValue = value.map(item => options.find(o => o.value === item))
      .filter(item => item)
      .map(item => item.label)
      .join('|');

    return (
      <div
        { ...res }
        className={`tc-15-dropdown tc-15-dropdown-btn-style ${className || ''} ${disabled ? 'disabled' : ''}`}
      >
        <span onClick={ this.toggleMenu.bind(this) } style={ {cursor: 'pointer'} }>
          {
            children ||
            <a href="javascript:void(0)"
               className="tc-15-dropdown-link"
               data-title={ showValue }
               style={ {
                 overflow: 'hidden',
                 textOverflow: 'ellipsis',
                 width: '100%',
                 maxWidth: '200px',
               } }
            >
              {
                showValue
                  ? <span>{ showValue }</span>
                  : <span>请选择</span>
              }
              <i className="caret"></i>
            </a>
          }
        </span>
        {
          isOpened &&
          <div className="tc-15-filtrateu">
            <div className="search-box search-box-simple m">
              <div className="search-input-wrap">
                <textarea className="tc-15-input-text search-input"
                          value={ keyword }
                          onChange={ e => this.setState({keyword: e.target.value}) }
                ></textarea>
              </div>
              <input type="button" className="search-btn" value="搜索"/>
            </div>
            <ul role="menu" className="tc-15-filtrate-menu">
              {
                mode !== 'single' && ! keyword && list.length > 0 &&
                <li role="presentation" className="tc-15-optgroup">
                  <label title="全选" className="tc-15-checkbox-wrap">
                    <input type="checkbox"
                           className="tc-15-checkbox"
                           onChange={ e => this.changeAllSelected(e) }
                           checked={ tmpList.length === options.filter(item => ! item.disabled).length }
                    />
                    <span>全选</span>
                  </label>
                </li>
              }
              {
                list.map(item => (
                  <li role="presentation" className="tc-15-optgroup" key={ item.value }>
                    {
                      mode === 'single'
                        ? <label
                          title={ item.label }
                          className="tc-15-checkbox-wrap"
                          onClick={ () => this.handleChange(item) }
                        >
                          <span>{ item.label }</span>
                        </label>
                        : <label title={ item.label } className="tc-15-checkbox-wrap">
                          <input
                            type="checkbox"
                            className="tc-15-checkbox"
                            onChange={ e => ! item.disabled && this.changeSelected(e, item) }
                            checked={ tmpList.findIndex(v => v.value === item.value) > - 1 }
                            value={ item.value }
                            disabled={ item.disabled }
                          />
                          <span
                            dangerouslySetInnerHTML={ {
                              __html: item.label.replace(
                                new RegExp(`(${keyword})`),
                                `<em style="color:#E1504A">$1</em>`,
                              ),
                            } }
                          ></span>
                        </label>
                    }
                  </li>
                ))
              }
            </ul>
            {
              mode !== 'single' &&
              <div className="tc-15-filtrate-ft">
                <button className="tc-15-btn m"
                        onClick={ this.handleSubmit.bind(this) }
                >
                  确定
                </button>
                <button className="tc-15-btn m weak"
                        onClick={ this.handleCancel.bind(this) }
                >
                  取消
                </button>
              </div>
            }
          </div>
        }
      </div>
    );
  }
}
