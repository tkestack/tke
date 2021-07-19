import * as React from 'react';

import { OnOuterClick } from '../libs/decorators/OnOuterClick';

export interface PageSelectProps {
  current: string|number;
  showList: number[];
  onChange: Function;
}

interface PageSelectState {
  visible: boolean;
}

export class PageSelect extends React.Component<PageSelectProps, PageSelectState> {
  constructor(props) {
    super(props);
    this.state = {
      visible: false,
    };
  }

  @OnOuterClick
  func() {
    this.setState({
      visible: false,
    });
  }

  onChange(pageNo) {
    this.setState({
      visible: false,
    });
    this.props.onChange(pageNo);
  }

  render() {
    const { current, showList, onChange } = this.props;
    const { visible } = this.state;

    return (
      <div className="tc-15-page-select">
        <a href="javascript:void(0)" className="indent"
          onClick={() => this.setState({ visible: !visible })}
        >
          <span>{current}</span>
          <span className="ico-arrow"></span>
        </a>
        <ul style={{ display: visible ? 'block' : 'none' }}
          className="tc-15-simulate-option tc-15-def-scroll"
        >
          { showList.map(item => <li key={item} onClick={() => this.onChange(item)}>{item}</li>) }
        </ul>
      </div>
    );
  }
}
