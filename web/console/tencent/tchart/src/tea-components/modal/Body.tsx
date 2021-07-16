import * as React from 'react';

export interface BodyProps {
  children?: React.ReactNode;
}

interface BodyState {
}

export class Body extends React.Component<BodyProps, BodyState> {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div className="tc-15-rich-dialog-bd">
        {this.props.children}
      </div>
    );
  }
}
