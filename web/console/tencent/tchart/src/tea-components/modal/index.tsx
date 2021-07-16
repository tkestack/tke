import * as React from 'react';
import * as ReactDom from 'react-dom';
import { Header } from './Header';
import { Body } from './Body';
import { Footer } from './Footer';


interface BaseReactProps {
  key?: string;
  defaultValue?: string;
  children?: React.ReactNode;
  className?: string;
  placeholder?: string;
  style?: object;
}

export interface ModalProps extends BaseReactProps {
  title: string | React.ReactNode;
  visible: boolean;
  okText?: string;
  cancelText?: string;
  onCancel?: Function;
  onOk?: Function;
  footer?: React.ReactNode;
}

export interface ConfirmProps extends BaseReactProps {
  onOk: Function;
  okText?: string;
  cancelText?: string;
  title: string | React.ReactNode;
  content: string | React.ReactNode;
}

interface ModalState {
}

export class Modal extends React.Component<ModalProps, ModalState> {
  el = document.createElement('div');
  constructor(props) {
    super(props);
    this.state = {
    };
  }
  componentDidMount() {
    document.body.appendChild(this.el);
  }
  componentWillUnmount() {
    if (!ReactDom.createPortal) {
      ReactDom.unmountComponentAtNode
        && ReactDom.unmountComponentAtNode(this.el)
    }
    document.body.removeChild(this.el);
  }

  render() {
    const {
      title,
      onOk,
      footer,
      okText,
      visible,
      style,
      children,
      onCancel,
      cancelText,
    } = this.props;

    const footerProps = {
      onOk,
      footer,
      okText,
      visible,
      onCancel,
      cancelText,
    };

    let component = <div className="dialog-panel" style={{ background: 'rgba(0, 0, 0, .5)' }}>
      <div className="tc-15-rich-dialog" style={{ margin: '200px auto', ...style }}>
        <Header title={title} close={onCancel || onOk}></Header>
        <Body>{children}</Body>
        <Footer {...footerProps}>
        </Footer>
      </div>
    </div>

    if (ReactDom.createPortal) {
      return visible ? ReactDom.createPortal(component, this.el) : null;
    }

    visible && ReactDom.render(component, this.el);
    return undefined
  }
}

export function confirm(props: ConfirmProps) {
  const div = document.createElement('div');
  document.body.appendChild(div);
  const {
    okText,
    cancelText,
    onOk,
    title,
    style,
    content,
  } = props;

  let visible = true;

  const onCancel = () => {
    visible = false;
    ReactDom.unmountComponentAtNode && ReactDom.unmountComponentAtNode(div);
    document.body.removeChild(div);
  };

  const param = {
    title,
    style,
    okText,
    visible,
    cancelText,
    onCancel,
    children: content,
    onOk: () => {
      onOk();
      onCancel();
    },
  };

  ReactDom.render(<Modal {...param} />, div);
}
