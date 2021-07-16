import * as React from 'react';


export interface HeaderProps {
  title: string | React.ReactNode;
  close: Function;
}

interface HeaderState {
}

export class Header extends React.Component<HeaderProps, HeaderState> {
  constructor(props) {
    super(props);
    this.state = {};
  }

  render() {
    const {title, close} = this.props;
    return (
      <div className="tc-15-rich-dialog-hd">
        <strong>{ title }</strong>
        <button
          title="关闭"
          className="tc-15-btn-close"
          onClick={ () => close() }
        >
          关闭
        </button>
      </div>
    );
  }
}
