import * as React from 'react';
import { DateTimePicker } from '../../tea-components/datetimepicker/index';
import { TIME_PICKER } from '../constants';


interface TimeTab {
  label: string;
  fromBefore: number;
  from?: Date;
  to?: Date;
}

const DurationsByPeriod = {
  // [period(s)] : day
  60: 15,
  300: 31,
  3600: 62,
  86400: 186
};

export interface ToolbarProps {
  style?: React.CSSProperties;
  duration: {from: Date | string; to: Date | string};
  onChangeTime: Function;
}

export class Toolbar extends React.Component<ToolbarProps> {
  nowDay = new Date();

  // 可选择日期最大间隔
  range = {
    min: new Date(+ this.nowDay - 864e5 * DurationsByPeriod[300]),
    max: new Date(this.nowDay),
    maxLength: DurationsByPeriod[300] * 86400
  };

  constructor(props) {
    super(props);
  }

  onTimePickerChange = (dateTime, label) => {
    this.props.onChangeTime(new Date(dateTime.from), new Date(dateTime.to));
  };

  render() {
    const {duration, style} = this.props;

    return (
      <div className="tc-action-grid" style={{...style}}>
        <div className="justify-grid">
          <div className="col">
            <div className="tc-15-calendar-select-wrap tc-15-calendar2-hook">
              <span className="dateTimePickerVS" onClick={ e => e.stopPropagation() }>
                <DateTimePicker
                  tabs={ TIME_PICKER.Tabs as any }
                  defaultSelectedTabIndex={ 0 }
                  defaultValue={ duration }
                  onChange={ this.onTimePickerChange }
                  range={ this.range }
                />
              </span>
            </div>
            <div className="monitor-dialog-bd-right" style={ {display: "inline-block", marginLeft: 15} }>
              { this.props.children }
            </div>
          </div>
        </div>
      </div>
    );
  }
}
