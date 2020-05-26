import * as React from 'react';

import { Bubble } from '@tencent/tea-component';

import { RootProps } from './InstallerApp';

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
   * vpc网段地址
   */
  vpcCIDR?: string;

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

  maxMaskCode?: string;

  cidr?: string;
  maxNodePodNum?: number;
  maxClusterServiceNum?: number;
}

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
    maxMaskCode: this.props.maxMaskCode,

    cidr: '',
    maxNodePodNum: this.props.maxNodePodNum || 256,
    maxClusterServiceNum: this.props.maxClusterServiceNum || 256
  };

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
      if (parts[0] === '192' || parts[0] === '172') {
        this.setState({ field5: parts[4] || '16', minMaskCode: '16' });
      } else if (parts[0] === '10') {
        this.setState({ field5: parts[4] || '14', minMaskCode: '14' });
      }
      this.handleUpdate(parts[0], parts[4], cidr);
    } else {
      if (this.state.field1 === '192' || this.state.field1 === '172') {
        this.setState({ field5: '16', minMaskCode: '16' });
      } else if (this.state.field1 === '10') {
        this.setState({ field5: '14', minMaskCode: '14' });
      }
    }
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.vpcCIDR !== this.props.vpcCIDR && nextProps.vpcCIDR) {
      let ads = nextProps.vpcCIDR.split('.');
      if (ads[0] === '10' || ads[0] === '192') {
        this.setState({ field1: '172', field5: '16' });
        this.handleUpdate('172', '16', '');
      } else if (ads[0] === '172') {
        this.setState({ field1: '10', field5: '14' });
        this.handleUpdate('10', '14', '');
      }
    }
  }

  _handleChange(cidr, maxNodePodNum, maxClusterServiceNum) {
    this.props.onChange(cidr, maxNodePodNum, maxClusterServiceNum);
  }

  computerRange(digit: number, min: number, max: number) {
    let base = Math.pow(2, 8 - digit),
      range: string[] = [];
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
    let order = Math.ceil(parseInt(mask) / 8),
      digit = parseInt(mask) % 8,
      field2: string,
      field2Range: string[],
      field3: string,
      field3Range: string[];
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
    } else if (period === '10') {
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
      this.setState({
        field2: '168',
        field2Range: ['168'],
        minMaskCode: '16',
        field5: '16'
      });
      this.handleUpdate('192', '16', '');
    } else if (value === '172') {
      this.setState({ field5: '16', minMaskCode: '16' });
      this.handleUpdate('172', '16', '');
    } else {
      this.setState({ field5: '14', minMaskCode: '14' });
      this.handleUpdate('10', '14', '');
    }
  }

  _findNearest(arr, num) {
    let newArr = [];
    num = num ? num : 0;
    arr.map(x => newArr.push(Math.abs(parseInt(x) - parseInt(num))));
    const index = newArr.indexOf(Math.min.apply(null, newArr));
    return arr[index];
  }

  _handleField2Input(value) {
    let reg = /^\d+$/;
    if (reg.test(value) || !value) {
      this.setState({ field2: value });
    }
  }

  _handleField2Blur(value) {
    value = this._findNearest(this.state.field2Range, value);
    let { field1, field3, field4, field5, maxNodePodNum, maxClusterServiceNum } = this.state;
    this.setState({
      field2: value,
      cidr: field1 + '.' + value + '.' + field3 + '.' + field4 + '/' + field5
    });
    this._handleChange(
      field1 + '.' + value + '.' + field3 + '.' + field4 + '/' + field5,
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
    value = this._findNearest(this.state.field3Range, value);
    let { field1, field2, field4, field5, maxNodePodNum, maxClusterServiceNum } = this.state;
    this.setState({
      field3: value,
      cidr: field1 + '.' + field2 + '.' + value + '.' + field4 + '/' + field5
    });
    this._handleChange(
      field1 + '.' + field2 + '.' + value + '.' + field4 + '/' + field5,
      maxNodePodNum,
      maxClusterServiceNum
    );
  }

  _handleField5Input(value) {
    this.setState({ field5: value });
    this.handleUpdate(this.state.field1, value, '');
  }
  _handlerPodSelect(value) {
    value = parseInt(value);
    let { field1, field2, field3, field4, field5, maxClusterServiceNum } = this.state;
    this.setState({
      maxNodePodNum: value
    });
    this._handleChange(field1 + '.' + field2 + '.' + field3 + '.' + field4 + '/' + field5, value, maxClusterServiceNum);
  }

  _handlerServiceSelect(value) {
    value = parseInt(value);
    let { field1, field2, field3, field4, field5, maxNodePodNum } = this.state;
    this.setState({
      maxClusterServiceNum: value
    });

    this._handleChange(field1 + '.' + field2 + '.' + field3 + '.' + field4 + '/' + field5, maxNodePodNum, value);
  }

  formatRange(range: string[]) {
    if (range.length <= 4) {
      return <p>范围：{range.join(', ')}</p>;
    } else {
      return (
        <p>
          范围：
          {range[0] + ', ' + range[1] + ', ... , ' + range[range.length - 1]}
        </p>
      );
    }
  }

  _renderNodeMax() {
    //集群内节点数量上限 = （CIDR IP个数/ Pod数量上线） - （Serivce个数/Pod数量上限）
    let { field1, field2, field3, field4, field5, maxNodePodNum, maxClusterServiceNum } = this.state;
    return Math.pow(2, 32 - parseInt(field5)) / maxNodePodNum - Math.ceil(maxClusterServiceNum / maxNodePodNum);
  }

  render() {
    let { parts, onChange } = this.props,
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
        maxMaskCode,
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

    return (
      <div style={{ fontSize: '14px' }} className="form-unit">
        <ul className="form-list jiqun">
          <li className="pure-text-row">
            <div className="form-label ">
              <label htmlFor="">CIDR</label>
            </div>
            <div className="form-input">
              <select
                className="tc-15-select m"
                style={{ minWidth: '20px' }}
                value={field1}
                onChange={e => this._handleField1Input(e.target.value)}
              >
                {field1Options}
              </select>
              &nbsp;.&nbsp;
              <Bubble content={field2Range.length > 1 && this.formatRange(field2Range)}>
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
              &nbsp;.&nbsp;
              <Bubble content={field3Range.length > 1 && this.formatRange(field3Range)}>
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
              &nbsp;.&nbsp;
              <Bubble content={field4Range.length > 1 && this.formatRange(field4Range)}>
                <input
                  type="text"
                  className="tc-15-input-text m shortest"
                  disabled={field4Range.length === 1}
                  value={field4}
                  style={{ width: '46px' }}
                />
              </Bubble>
              &nbsp;/&nbsp;
              <select
                className="tc-15-select m"
                style={{ minWidth: '20px' }}
                value={field5}
                onChange={e => this._handleField5Input(e.target.value)}
              >
                {field5Options}
              </select>
              <p className="tea-form__help-text">CIDR不能与目标机器IP段重叠， 否则会造成初始化失败</p>
            </div>
          </li>
          <li>
            <div className="form-label ">
              <label htmlFor="">Pod数量上限/节点</label>
            </div>
            <div className="form-input">
              <div className="form-unit">
                <select className="tc-15-select m" onChange={e => this._handlerPodSelect(e.target.value)}>
                  <option value="32" selected={maxNodePodNum === 32}>
                    32
                  </option>
                  <option value="64" selected={maxNodePodNum === 64}>
                    64
                  </option>
                  <option value="128" selected={maxNodePodNum === 128}>
                    128
                  </option>
                  <option value="256" selected={maxNodePodNum === 256}>
                    256
                  </option>
                </select>
              </div>
            </div>
          </li>
          <li>
            <div className="form-label ">
              <label htmlFor="">Service数量上限/集群</label>
            </div>
            <div className="form-input">
              <div className="form-unit">
                <select className="tc-15-select m" onChange={e => this._handlerServiceSelect(e.target.value)}>
                  <option value="128" selected={maxClusterServiceNum === 128}>
                    128
                  </option>
                  <option value="256" selected={maxClusterServiceNum === 256}>
                    256
                  </option>
                  <option value="512" selected={maxClusterServiceNum === 512}>
                    512
                  </option>
                  <option value="1024" selected={maxClusterServiceNum === 1024}>
                    1024
                  </option>
                  <option value="2048" selected={maxClusterServiceNum === 2048}>
                    2048
                  </option>
                  <option value="4096" selected={maxClusterServiceNum === 4096}>
                    4096
                  </option>
                  <option value="8192" selected={maxClusterServiceNum === 8192}>
                    8192
                  </option>
                  <option value="16384" selected={maxClusterServiceNum === 16384}>
                    16384
                  </option>
                  <option value="32768" selected={maxClusterServiceNum === 32768}>
                    32768
                  </option>
                </select>
              </div>
            </div>
          </li>
          <li>
            <div className="form-label ">
              <label htmlFor=""></label>
            </div>
            <div className="form-input">
              <div className="form-unit">
                <p style={{ float: 'left', width: '400px', fontSize: '12px' }} className="tea-form__help-text">
                  当前容器网络配置下，集群最多{this._renderNodeMax()}个节点
                </p>
              </div>
            </div>
          </li>
        </ul>
      </div>
    );
  }
}
