import * as React from 'react';
import { HeaderProps } from './tableProps';

interface ColGroupState {
}

export class ColGroup extends React.Component<HeaderProps, ColGroupState> {
  constructor(props) {
    super(props);
  }

  render() {
    const { check, columns } = this.props;

    return (
      <colgroup>
        {
          check && <col style={{ width: '40px' }}/>
        }
        {
          (columns || []).map(item => (
            <col key={item.key || item.dataIndex}
              style={{ width: typeof item.width === 'number' ? `${item.width}px` : item.width }}
            />
          ))
        }
      </colgroup>
    );
  }
}
