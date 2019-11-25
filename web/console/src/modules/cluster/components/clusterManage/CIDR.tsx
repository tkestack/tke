import * as React from 'react';
import { RootProps } from '../ClusterApp';
import { Bubble, ExternalLink, Icon, Checkbox } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { FormPanel } from '../../../common/components';

export interface CIDRProps extends RootProps {
  /**
   * 第一部分范围
   */
  parts?: string[];

  /**
   * CIDR 值
   */
  value?: string;

  maxNodePodNum?: number;
  maxClusterServiceNum?: number;

  /**
   * 最小掩码长度
   */
  minMaskCode?: string;

  /**
   * 最大掩码长度
   */
  maxMaskCode?: string;

  /**
   * 变化事件
   */
  onChange?: (value, maxNodePodNum, maxClusterServiceNum) => void;

  /**
   * 失去焦点事件
   */
  onBlur?: (value) => void;
}

export interface CIDRState {
  /**
   * IP第一部分
   */
  field1?: string;

  /**
   * 第二部分范围
   */
  field1Range?: string[];

  /**
   * IP第二部分
   */
  field2?: string;

  /**
   * 第二部分范围
   */
  field2Range?: string[];

  /**
   * IP第三部分
   */
  field3?: string;

  /**
   * 第三部分范围
   */
  field3Range?: string[];

  /**
   * IP第四部分
   */
  field4?: string;

  /**
   * 第四部分范围
   */
  field4Range?: string[];

  /**
   * IP第五部分，掩码
   */
  field5?: string;

  /**
   * 第五部分范围
   */
  minMaskCode?: string;

  // maxMaskCode?: string;

  cidr?: string;
  maxNodePodNum?: number;
  maxClusterServiceNum?: number;
}

const IsLikeTen = (segment: string) => {
  return segment === '10' || segment === '8' || segment === '9' || segment === '100' || segment === '127'
    ? true
    : false;
};

export class CIDR extends React.Component<CIDRProps, CIDRState> {
  state = {
    field1: this.props.parts[0],
    field1Range: this.props.parts,
    field2: '168',
    field2Range: ['168'],
    field3: '0',
    field3Range: ['0'],
    field4: '0',
    field4Range: ['0'],
    field5: '16',
    minMaskCode: this.props.minMaskCode,
    // maxMaskCode: this.props.maxMaskCode,

    cidr: '',
    maxNodePodNum: 256,
    maxClusterServiceNum: 256
  };

  UNSAFE_componentWillReceiveProps;

  componentDidMount() {
    let cidr = this.props.value;
    if (cidr) {
      let parts = cidr.split(/[\.\/]/);
      this.setState({
        field1: parts[0],
        field2: parts[1],
        field3: parts[2],
        field4: parts[3],
        field5: parts[4]
      });
      let segment = parts[0];
      if (segment === '192' || segment === '172') {
        this.setState({ field5: parts[4] || '16', minMaskCode: '16' });
      } else if (IsLikeTen(segment)) {
        this.setState({ field5: parts[4] || '14', minMaskCode: '14' });
      }
      this.handleUpdate(segment, parts[4], cidr);
    } else {
      if (this.state.field1 === '192' || this.state.field1 === '172') {
        this.setState({ field5: '16', minMaskCode: '16' });
      } else if (IsLikeTen(this.state.field1)) {
        this.setState({ field5: '14', minMaskCode: '14' });
      }
    }
  }

  _handleChange(cidr, maxNodePodNum, maxClusterServiceNum) {
    // let {cidr, maxNodePodNum, maxClusterServiceNum} = this.state;
    this.props.onChange(cidr, maxNodePodNum, maxClusterServiceNum);
  }

  computerRange(digit: number, min: number, max: number) {
    let base = Math.pow(2, 8 - digit);
    let range: string[] = [];
    for (let i = 0; i < max >> (8 - digit); i++) {
      range.push(min + i * base + '');
    }
    // if(digit === 8){
    //     for(let i = 1; i < max; i++)
    //     {
    //         range.push(min + i + "")
    //     }
    // }
    return range;
  }

  handleUpdate(period: string, mask: string, cidr: string) {
    let order = Math.ceil(parseInt(mask, 10) / 8);
    let digit = parseInt(mask) % 8;
    let field2: string;
    let field2Range: string[];
    let field3: string;
    let field3Range: string[];
    if (digit === 0) {
      digit = 8;
    }
    let parts = [];
    if (cidr) {
      parts = cidr.split(/[\.\/]/);
    }
    if (period === '192') {
      if (order >= 3) {
        field2 = parts[1] || '168';
        field2Range = ['168'];
        field3Range = this.computerRange(digit, 0, 256);
        field3 = parts[2] || field3Range[0];
      } else {
        field2 = parts[1] || '168';
        field2Range = ['168'];
        field3 = parts[2] || '0';
        field3Range = ['0'];
      }
    } else if (period === '172') {
      if (order >= 3) {
        field2Range = this.computerRange(8, 16, 16);
        //去除172.17网段
        let cIndex = field2Range.indexOf('17');
        if (cIndex > -1) {
          field2Range.splice(cIndex, 1);
        }

        field2 = parts[1] || field2Range[0];
        field3Range = this.computerRange(digit, 0, 256);
        field3 = parts[2] || field3Range[0];
      } else {
        field2Range = this.computerRange(digit, 16, 16);
        //去除172.17网段
        let cIndex = field2Range.indexOf('17');
        if (cIndex > -1) {
          field2Range.splice(cIndex, 1);
        }
        field2 = parts[1] || field2Range[0];
        field3 = parts[2] || '0';
        field3Range = ['0'];
      }
    } else if (IsLikeTen(period)) {
      if (order >= 3) {
        field2Range = this.computerRange(8, 0, 256);
        field2 = parts[1] || field2Range[0];
        field3Range = this.computerRange(digit, 0, 256);
        field3 = parts[2] || field3Range[0];
      } else {
        field2Range = this.computerRange(digit, 0, 256);
        field2 = parts[1] || field2Range[0];
        field3 = parts[2] || '0';
        field3Range = ['0'];
      }
    }

    let { maxNodePodNum, maxClusterServiceNum } = this.state;
    this.setState({
      field2,
      field2Range,
      field3,
      field3Range,
      cidr: period + '.' + field2 + '.' + field3 + '.0/' + mask
    });

    this._handleChange(period + '.' + field2 + '.' + field3 + '.0/' + mask, maxNodePodNum, maxClusterServiceNum);
  }

  _handleField1Input(value) {
    this.setState({ field1: value });
    if (value === '192') {
      this.setState({ field2: '168', field2Range: ['168'], minMaskCode: '16', field5: '16' });
      this.handleUpdate('192', '16', '');
    } else if (value === '172') {
      this.setState({ field5: '16', minMaskCode: '16' });
      this.handleUpdate('172', '16', '');
    } else {
      this.setState({ field5: '14', minMaskCode: '14' });
      this.handleUpdate(value, '14', '');
    }
  }

  _findNearest(arr, num) {
    let newArr = [];
    let newNum = num ? num : 0;
    arr.map(x => newArr.push(Math.abs(parseInt(x) - parseInt(newNum))));
    let index = newArr.indexOf(Math.min.apply(null, newArr));
    return arr[index];
  }

  _handleField2Input(value) {
    let reg = /^\d+$/;
    if (reg.test(value) || !value) {
      this.setState({ field2: value });
    }
  }

  _handleField2Blur(value) {
    let newValue = this._findNearest(this.state.field2Range, value);
    let { field1, field3, field4, field5, maxNodePodNum, maxClusterServiceNum } = this.state;
    this.setState({
      field2: newValue,
      cidr: field1 + '.' + newValue + '.' + field3 + '.' + field4 + '/' + field5
    });
    this._handleChange(
      field1 + '.' + newValue + '.' + field3 + '.' + field4 + '/' + field5,
      maxNodePodNum,
      maxClusterServiceNum
    );
  }

  _handleField3Input(value) {
    let reg = /^\d+$/;
    if (reg.test(value) || !value) {
      this.setState({ field3: value });
    }
  }

  _handleField3Blur(value) {
    let newValue = this._findNearest(this.state.field3Range, value);
    let { field1, field2, field4, field5, maxNodePodNum, maxClusterServiceNum } = this.state;
    this.setState({
      field3: newValue,
      cidr: field1 + '.' + field2 + '.' + newValue + '.' + field4 + '/' + field5
    });
    this._handleChange(
      field1 + '.' + field2 + '.' + newValue + '.' + field4 + '/' + field5,
      maxNodePodNum,
      maxClusterServiceNum
    );
  }

  _handleField5Input(value) {
    this.setState({ field5: value });
    this.handleUpdate(this.state.field1, value, '');
  }

  _handlerPodSelect(value) {
    let newValue = parseInt(value);
    let { field1, field2, field3, field4, field5, maxClusterServiceNum } = this.state;
    this.setState({
      maxNodePodNum: newValue
    });
    this._handleChange(
      field1 + '.' + field2 + '.' + field3 + '.' + field4 + '/' + field5,
      newValue,
      maxClusterServiceNum
    );
  }

  _handlerServiceSelect(value) {
    let newValue = parseInt(value);
    let { field1, field2, field3, field4, field5, maxNodePodNum } = this.state;
    this.setState({
      maxClusterServiceNum: newValue
    });

    this._handleChange(field1 + '.' + field2 + '.' + field3 + '.' + field4 + '/' + field5, maxNodePodNum, newValue);
  }

  formatRange(range: string[]) {
    if (range.length <= 4) {
      return (
        <p>
          {t('范围：')}
          {range.join(', ')}
        </p>
      );
    } else {
      return (
        <p>
          {t('范围：')}
          {range[0] + ', ' + range[1] + ', ... , ' + range[range.length - 1]}
        </p>
      );
    }
  }
  _renderNodeMax() {
    //集群内节点数量上限 = （CIDR IP个数/ Pod数量上线） - （Serivce个数/Pod数量上限）
    let { field5, maxNodePodNum, maxClusterServiceNum } = this.state;
    return Math.pow(2, 32 - parseInt(field5)) / maxNodePodNum - Math.ceil(maxClusterServiceNum / maxNodePodNum);
  }

  render() {
    let { maxMaskCode } = this.props,
      {
        field1,
        field1Range,
        field2,
        field2Range,
        field3,
        field3Range,
        field4,
        field4Range,
        field5,
        minMaskCode,

        maxNodePodNum,
        maxClusterServiceNum
      } = this.state,
      field1Options = field1Range.map((part, index) => (
        <option key={index} value={part + ''}>
          {part}
        </option>
      )),
      field5Options = [];

    for (let i = parseInt(minMaskCode); i <= parseInt(maxMaskCode); i++) {
      field5Options.push(<option value={i + ''}>{i}</option>);
    }

    let podOptions = [
      {
        text: '32',
        value: '32'
      },
      {
        text: '64',
        value: '64'
      },
      {
        text: '128',
        value: '128'
      },
      {
        text: '256',
        value: '256'
      }
    ];

    let serviceOptions = [
      {
        text: '128',
        value: '128'
      },
      {
        text: '256',
        value: '256'
      },
      {
        text: '512',
        value: '512'
      },
      {
        text: '1024',
        value: '1024'
      },
      {
        text: '2048',
        value: '2048'
      },
      {
        text: '4096',
        value: '4096'
      },
      {
        text: '8192',
        value: '8192'
      },
      {
        text: '16384',
        value: '16384'
      },
      {
        text: '32768',
        value: '32768'
      }
    ];

    podOptions.unshift({
      text: '16',
      value: '16'
    });

    serviceOptions.unshift({
      text: '64',
      value: '64'
    });
    serviceOptions.unshift({
      text: '32',
      value: '32'
    });
    const nodeMax = this._renderNodeMax();
    return (
      <FormPanel.Item label={'容器网络'}>
        <FormPanel isNeedCard={false}>
          <FormPanel.Item label={t('CIDR')}>
            <select
              className="tc-15-select m"
              style={{ minWidth: '20px' }}
              value={field1}
              onChange={e => this._handleField1Input(e.target.value)}
            >
              {field1Options}
            </select>
            <FormPanel.HelpText>&nbsp;.&nbsp;</FormPanel.HelpText>
            <Bubble placement="bottom" content={field2Range.length > 1 ? this.formatRange(field2Range) : null}>
              <input
                type="text"
                className="tc-15-input-text m shortest"
                disabled={field2Range.length === 1}
                value={field2}
                onChange={e => this._handleField2Input(e.target.value)}
                onBlur={e => this._handleField2Blur(e.target.value)}
                style={{ width: '46px' }}
              />
            </Bubble>
            <FormPanel.HelpText>&nbsp;.&nbsp;</FormPanel.HelpText>
            <Bubble placement="bottom" content={field3Range.length > 1 ? this.formatRange(field3Range) : null}>
              <input
                type="text"
                className="tc-15-input-text m shortest"
                disabled={field3Range.length === 1}
                value={field3}
                onChange={e => this._handleField3Input(e.target.value)}
                onBlur={e => this._handleField3Blur(e.target.value)}
                style={{ width: '46px' }}
              />
            </Bubble>
            <FormPanel.HelpText>&nbsp;.&nbsp;</FormPanel.HelpText>
            <Bubble placement="bottom" content={field4Range.length > 1 ? this.formatRange(field4Range) : null}>
              <input
                type="text"
                className="tc-15-input-text m shortest"
                disabled={field4Range.length === 1}
                value={field4}
                style={{ width: '46px' }}
              />
            </Bubble>
            <FormPanel.HelpText>&nbsp;/&nbsp;</FormPanel.HelpText>
            <select
              className="tc-15-select m"
              style={{ minWidth: '20px' }}
              value={field5}
              onChange={e => this._handleField5Input(e.target.value)}
            >
              {field5Options}
            </select>
          </FormPanel.Item>

          <FormPanel.Item
            label={t('Pod数量上限/节点')}
            select={{
              options: podOptions,
              value: maxNodePodNum + '',
              onChange: value => this._handlerPodSelect(parseInt(value))
            }}
          />
          <FormPanel.Item
            label={t('Service数量上限/集群')}
            select={{
              options: serviceOptions,
              value: maxClusterServiceNum + '',
              onChange: value => this._handlerServiceSelect(parseInt(value))
            }}
          />
          <FormPanel.Item
            label={
              <Trans count={nodeMax}>
                当前容器网络配置下，集群最多 <strong style={{ color: '#E1504A' }}>{{ nodeMax }}</strong> 个节点
              </Trans>
            }
          />
        </FormPanel>
      </FormPanel.Item>
    );
  }
}
