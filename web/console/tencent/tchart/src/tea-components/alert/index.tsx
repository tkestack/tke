import * as React from 'react';

export interface AlertProps {
  className?: string;
  children: React.ReactNode;
}

export class Alert extends React.Component<AlertProps, object> {
  constructor(props) {
    super(props);
  }

  render() {
    const { children, className } = this.props;

    return (
      <div className={`tc-15-msg ${className}`}>
        <div className="tip-info">{children}</div>
      </div>
    );
  }
}
