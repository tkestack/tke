import * as React from 'react';
import { RootProps } from '../components/LogApp';
import { FormLayout } from '../../common/layouts';
import { FormItem, DateTimePicker } from '../../common/components';
import { Button, DatePickerDurationAnchor } from '@tencent/qcloud-component';
import { insertCSS } from '@tencent/qcloud-lib';

insertCSS(
  'LogActionPanelCSS',
  `
.form-list .form-label {
  text-align: right;
}
.tc-15-calendar-select-wrap span[role=tab] {
  padding: 0px 13px;
}
.tc-15-calendar-i-next-m span, .tc-15-calendar-i-pre-m span {
  display: none
}
`
);

interface LogActionPanelState {
  appId?: string;
  ownerUin?: number | string;
  uin?: number | string;
  action?: string;
  region?: string;
  status?: number | string;
  returnCode?: string | number;
  startTime?: string | Date;
  endTime?: string | Date;
  body?: string;
  result?: string;
  refreshDate?: boolean;
}

export class LogActionPanel extends React.Component<RootProps, LogActionPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      appId: '',
      ownerUin: '',
      uin: '',
      action: '',
      region: '',
      status: '',
      returnCode: '',
      startTime: '',
      endTime: '',
      body: '',
      result: '',
      refreshDate: false
    };
  }

  componentDidMount() {
    const { route } = this.props;
    let now = Date.now();
    let endTime = new Date(now),
      startTime = new Date(now - 60 * 60 * 1000),
      refreshDate = false;

    if (route.queries['startTime'] && route.queries['endTime']) {
      endTime = new Date(parseInt(route.queries['endTime']));
      startTime = new Date(parseInt(route.queries['startTime']));
      refreshDate = true;
    } else if (route.queries['startTime']) {
      endTime = new Date(parseInt(route.queries['startTime']) + 24 * 60 * 60 * 1000);
      startTime = new Date(parseInt(route.queries['startTime']));
      refreshDate = true;
    } else if (route.queries['endTime']) {
      endTime = new Date(parseInt(route.queries['endTime']));
      startTime = new Date(parseInt(route.queries['endTime']) - 24 * 60 * 60 * 1000);
      refreshDate = true;
    }

    const action = route.queries['action'] || '',
      returnCode = decodeURI(route.queries['returnCode'] || '');
    this.setState({ endTime, startTime, action, returnCode, refreshDate });
    this.props.actions.applyFilter({ startTime, endTime, action, returnCode });
  }

  handleSearch(from?: string | Date, to?: string | Date) {
    let { actions } = this.props,
      { appId, ownerUin, uin, action, status, returnCode, startTime, endTime, body, result } = this.state;

    if (from || to) {
      startTime = from;
      endTime = to;
    }
    actions.applyFilter({ appId, ownerUin, uin, action, status, returnCode, startTime, endTime, body, result });
  }

  _handleKeyDown(e: React.KeyboardEvent) {
    const ENTER_KEY = 13;
    if (e.keyCode === ENTER_KEY) {
      this.handleSearch();
      e.preventDefault();
    }
  }

  render() {
    let { appId, ownerUin, uin, action, returnCode, startTime, endTime, body, result, refreshDate } = this.state;

    const tabs = [
      { from: '%NOW-1h', to: '%NOW', label: '近1h' },
      { from: '%NOW-6h', to: '%NOW', label: '近6h' },
      { from: '%NOW-12h', to: '%NOW', label: '近12h' }
    ];

    return (
      <div className="tc-action-grid">
        <div className="justify-grid">
          <FormLayout>
            <div className="tc-g">
              <div className="tc-g-u-1-5">
                <ul className="form-list jiqun" style={{ paddingBottom: '0px' }}>
                  <FormItem label="时间选择">
                    <div className="form-unit" style={{ width: '500px' }}>
                      <DateTimePicker
                        tabs={tabs}
                        linkage
                        defaultSelectedTabIndex={0}
                        defaultValue={{ from: startTime, to: endTime }}
                        refreshFlag={refreshDate}
                        onChange={(value, tabLabel) => {
                          this.setState({
                            startTime: value.from,
                            endTime: value.to
                          });
                          this.handleSearch(value.from, value.to);
                        }}
                      />
                    </div>
                  </FormItem>
                  <FormItem label="接口名">
                    <div className="form-unit">
                      <input
                        type="text"
                        className="tc-15-input-text m"
                        value={action}
                        onKeyDown={this._handleKeyDown.bind(this)}
                        onInput={e => {
                          this.setState({ action: e.target.value });
                        }}
                      />
                    </div>
                  </FormItem>
                </ul>
              </div>
              <div className="tc-g-u-1-5">
                <ul className="form-list jiqun" style={{ paddingBottom: '0px' }}>
                  <FormItem label="">
                    {/* <div className='form-unit'>
                      <DateTimePicker />
                    </div> */}
                  </FormItem>
                  <FormItem label="返回码">
                    <div className="form-unit">
                      <input
                        type="text"
                        className="tc-15-input-text m"
                        value={returnCode as string}
                        placeholder="支持NOT 0查询"
                        onKeyDown={this._handleKeyDown.bind(this)}
                        onInput={e => {
                          this.setState({ returnCode: e.target.value });
                        }}
                      />
                    </div>
                  </FormItem>
                </ul>
              </div>
              <div className="tc-g-u-1-5">
                <ul className="form-list jiqun" style={{ paddingBottom: '0px' }}>
                  <FormItem label="AppID">
                    <div className="form-unit">
                      <input
                        type="text"
                        className="tc-15-input-text m"
                        value={appId}
                        onKeyDown={this._handleKeyDown.bind(this)}
                        onInput={e => {
                          this.setState({ appId: e.target.value });
                        }}
                      />
                    </div>
                  </FormItem>
                  <FormItem label="请求参数">
                    <div className="form-unit">
                      <input
                        type="text"
                        className="tc-15-input-text m"
                        placeholder="支持模糊查询"
                        value={body}
                        onKeyDown={this._handleKeyDown.bind(this)}
                        onInput={e => {
                          this.setState({ body: e.target.value });
                        }}
                      />
                    </div>
                  </FormItem>
                </ul>
              </div>
              <div className="tc-g-u-1-5">
                <ul className="form-list jiqun" style={{ paddingBottom: '0px' }}>
                  <FormItem label="OwnerUin">
                    <div className="form-unit">
                      <input
                        type="text"
                        className="tc-15-input-text m"
                        value={ownerUin as string}
                        onKeyDown={this._handleKeyDown.bind(this)}
                        onInput={e => {
                          this.setState({ ownerUin: e.target.value });
                        }}
                      />
                    </div>
                  </FormItem>
                  <FormItem label="返回参数">
                    <div className="form-unit">
                      <input
                        type="text"
                        className="tc-15-input-text m"
                        placeholder="支持模糊查询"
                        value={result}
                        onKeyDown={this._handleKeyDown.bind(this)}
                        onInput={e => {
                          this.setState({ result: e.target.value });
                        }}
                      />
                    </div>
                  </FormItem>
                </ul>
              </div>
              <div className="tc-g-u-1-5">
                <ul className="form-list jiqun" style={{ paddingBottom: '0px' }}>
                  <FormItem label="Uin">
                    <div className="form-unit">
                      <input
                        type="text"
                        className="tc-15-input-text m"
                        value={uin as string}
                        onKeyDown={this._handleKeyDown.bind(this)}
                        onInput={e => {
                          this.setState({ uin: e.target.value });
                        }}
                      />
                    </div>
                  </FormItem>
                  <FormItem label="">
                    <Button
                      className="m"
                      onClick={() => {
                        this.handleSearch();
                      }}
                    >
                      查询
                    </Button>
                  </FormItem>
                </ul>
              </div>
            </div>
          </FormLayout>
        </div>
      </div>
    );
  }
}
