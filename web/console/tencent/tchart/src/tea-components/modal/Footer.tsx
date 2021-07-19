import * as React from 'react';

export interface FooterProps {
  footer?: React.ReactNode;
  okText?: string;
  cancelText?: string;
  onOk?: Function;
  onCancel?: Function;
}

interface FooterState {
}

export class Footer extends React.Component<FooterProps, FooterState> {
  constructor(props) {
    super(props);
  }

  render() {
    const {onOk, okText, cancelText, footer, onCancel} = this.props;

    return footer === null ? null : (
      <div className="tc-15-rich-dialog-ft">
        {
          footer
            ? <div className="tc-15-rich-dialog-ft-btn-wrap">
                { footer }
              </div>
            : <div className="tc-15-rich-dialog-ft-btn-wrap">
              {
                onOk &&
                <button className="tc-15-btn"
                        onClick={ () => onOk() }
                >
                  { okText || '确定' }
                </button>
              }
              {
                onCancel &&
                <button className="tc-15-btn weak"
                        onClick={ () => onCancel() }
                >
                  { cancelText || '取消' }
                </button>
              }
            </div>
        }
      </div>
    );
  }
}
