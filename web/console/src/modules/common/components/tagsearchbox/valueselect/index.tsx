import * as React from 'react';
import * as classNames from 'classnames';
import { PureInput } from './PureInput';
import { SingleValueSelect } from './SingleValueSelect';
import { MultipleValueSelect } from './MultipleValueSelect';
import { Loading } from './Loading';

export interface ValueSelectProps {
  type: string;
  values?: Array<any> | Function;
  inputValue?: string;
  onChange?: (value: any) => void;
  onSelect?: (value: any) => void;
  onCancel?: () => void;
  offset: number;
}

export class ValueSelect extends React.Component<ValueSelectProps, any> {
  mount = false;

  constructor(props) {
    super(props);

    let values = [];
    const propsValues = this.props.values;

    if (typeof propsValues !== 'function') {
      values = propsValues;
    }

    this.state = { values };
  }

  componentDidMount() {
    this.mount = true;
    const propsValues = this.props.values;
    if (typeof propsValues === 'function') {
      const res = propsValues();
      // Promise
      if (res && res.then) {
        res.then(values => {
          this.mount && this.setState({ values });
        });
      } else {
        this.mount && this.setState({ values: res });
      }
    }
  }

  componentWillUnmount() {
    this.mount = false;
  }

  handleKeyDown = (keyCode: number): boolean => {
    if (this['select'] && this['select'].handleKeyDown) {
      return this['select'].handleKeyDown(keyCode);
    }
    return true;
  };

  // handleKeyUp = (keyCode: number): boolean => {
  //   if (this['select'] && this['select'].handleKeyUp) {
  //     return this['select'].handleKeyUp(keyCode);
  //   }
  //   return true;
  // }

  render() {
    const values = this.state.values;
    const { type, inputValue, onChange, onSelect, onCancel, offset } = this.props;

    const props = { values, inputValue, onChange, onSelect, onCancel, offset };

    switch (type) {
      case 'input':
        return (
          <PureInput
            ref={select => {
              this['select'] = select;
            }}
            {...props}
          />
        );
      case 'single':
        if (values.length <= 0) {
          return <Loading offset={offset} />;
        }
        return (
          <SingleValueSelect
            ref={select => {
              this['select'] = select;
            }}
            {...props}
          />
        );
      case 'multiple':
        if (values.length <= 0) {
          return <Loading />;
        }
        return (
          <MultipleValueSelect
            ref={select => {
              this['select'] = select;
            }}
            {...props}
          />
        );
    }
    return null;
  }
}
