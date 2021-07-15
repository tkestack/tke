import * as React from 'react';
import * as ReactDOM from 'react-dom';


export interface SelectProps {
  /**
   * 当前选择的值
   */
  value: number;
  /**
   * 起止值
   */
  from: number;
  to: number;

  range?: SelectRange;
  onChange?: (value: number) => void;
}

export interface SelectRange {
  min?: number;
  max?: number;
}

/**
 * （受控组件）
 */
export class Select extends React.Component<SelectProps, any> {

  componentDidMount() {
    this.scrollToSelected(0);
  }

  componentDidUpdate() {
    this.scrollToSelected(150);
  }

  /**
   * 滚动到指定元素
   */
  scrollTo = (element, to, duration) => {
    if (duration <= 0) {
      element.scrollTop = to;
      return;
    }
    var difference = to - element.scrollTop;
    var perTick = difference / duration * 10;

    setTimeout(() => {
      element.scrollTop = element.scrollTop + perTick;
      if (element.scrollTop === to) return;
      this.scrollTo(element, to, duration - 10);
    }, 10);
  };


  /**
   * 滚动到当前选择元素
   */
  scrollToSelected = (duration) => {
    const select = ReactDOM.findDOMNode(this) as HTMLElement;
    const list = ReactDOM.findDOMNode(this.refs[ 'list' ]);

    const index = this.props.value;
    const topOption = (list as any).children[ index ] as HTMLElement;
    const to = topOption.offsetTop - select.offsetTop;

    this.scrollTo(select, to, duration);
  }

  /**
   * 根据范围生成列表
   */
  genRangeList = (start: number, end: number): Array<String> =>
    Array(end - start + 1).fill(0).map((e, i) => {
        const num = i + start;
        return num > 9 ? `${num}` : `0${num}`
      }
    );

  handleSelect = (e, val: number): void => {
    e.stopPropagation();
    if (this.props.onChange) this.props.onChange(val);
  }

  render() {
    const { from, to, value, range } = this.props;
    const list = this.genRangeList(from, to).map((item, i) => {
      if (range && 'min' in range && i < range.min) return (<li key={ i } className="disabled">{ item }</li>);
      if (range && 'max' in range && i > range.max) return (<li key={ i } className="disabled">{ item }</li>);
      return <li key={ i } className={ +item == value ? 'current' : '' }
                 onClick={ (e) => this.handleSelect(e, +item) }>{ item }</li>
    });

    return (
      <div className="tc-time-picker-select">
        <ul ref="list">{ list }</ul>
      </div>
    );
  }

}





