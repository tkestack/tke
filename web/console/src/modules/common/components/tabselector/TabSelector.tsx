import * as classnames from 'classnames';
import * as React from 'react';

// import { OnOuterClick, BaseReactProps, RecordSet, uuid } from '@tencent/ff-redux';
import { BaseReactProps, FetcherState, FetchState, RecordSet, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { Validation } from '../../models';

export interface TabItem {
  /**名称 */
  name?: string;

  /**当前选中的值 */
  value?: string;

  /**列表数据 */
  recordData?: FetcherState<RecordSet<any>>;

  /**指定value字段 */
  valueField?: string;

  /**指定text字段 */
  textField?: string;

  /**多字段组合显示 */
  textFields?: string[];

  /**多字段显示格式设置 替换变量格式为${val} */
  textFormat?: string;

  /**选中操作 */
  onSelect?: (value?: string) => void;

  /**列表为空时的提示 */
  emptyTip?: JSX.Element | string;
}

export interface TabSelectorProps extends BaseReactProps {
  /**显示名称 */
  label?: string;

  /**数据项 */
  tabs?: TabItem[];

  /**校验态 */
  validator?: Validation;

  /**提示 */
  tips?: JSX.Element | string;

  /**是否显示 */
  isShow?: boolean;
}

interface TabSelectorState {
  /**当前显示的标签索引 */
  current?: number;

  /**是否显示选择Tab */
  isShowTab?: boolean;
}

const FormatText = (data: any, fields: string[], format: string) => {
  let result = format;
  fields.forEach(field => {
    let reg: RegExp = new RegExp('\\$\\{' + field + '\\}', 'g');
    result = result.replace(reg, data[field]);
  });

  return result;
};

export class TabSelector extends React.Component<TabSelectorProps, TabSelectorState> {
  constructor(props, context) {
    super(props, context);

    this.state = {
      current: 0,
      isShowTab: false
    };
  }

  // @(OnOuterClick as any)
  public close() {
    this.setState({ isShowTab: false });
  }

  render() {
    let { label, tabs, validator, tips, isShow = true } = this.props,
      { isShowTab } = this.state,
      text = '';
    tabs.forEach(item => {
      text += item.value && item.value + '/';
    });
    text = text.substring(0, text.length - 1) || t('请选择') + label;
    return isShow ? (
      <div className={classnames('form-unit', { 'is-error': validator && validator.status === 2 })}>
        <div className="tc-15-dropdown tc-15-dropdown-btn-style tc-15-menu-active">
          <a href="javascript:;" className="tc-15-dropdown-link" onClick={this.toggleShowTab.bind(this)}>
            {text}
            <i className="caret" />
          </a>
          {isShowTab && (
            <div className="select-tab-mod" style={{ width: '300px', background: '#fff', position: 'absolute' }}>
              {this.renderTabHeader()}
              {this.renderTabContent()}
            </div>
          )}
        </div>
        {tips ? <p className="form-input-help text-weak">{tips}</p> : <noscript />}
        <p className="form-input-help">{validator && validator.status === 2 && validator.message}</p>
      </div>
    ) : (
      <noscript />
    );
  }

  private renderTabHeader() {
    let { tabs } = this.props,
      { current } = this.state;
    let headList = tabs.map((item, index) => {
      return (
        <li
          className={classnames({ 'tc-cur': current === index })}
          key={index}
          onClick={() => tabs[index].value && this.handleHeaderSelect(index)}
        >
          <a href="javascript:;" title="" role="tab">
            {item.name}
          </a>
        </li>
      );
    });
    return <ul className="tc-15-tablist">{headList}</ul>;
  }

  private handleHeaderSelect(step: number) {
    this.setState({ current: step });
  }

  private stepNext() {
    let { tabs } = this.props,
      { current } = this.state;
    if (tabs.length === current + 1) {
      this.close();
    } else {
      this.setState({ current: current + 1 });
    }
  }

  private renderTabContent() {
    let { tabs } = this.props,
      { current } = this.state,
      tab: TabItem = tabs[current];

    if (!tab) return;

    let { name, value, recordData, valueField, textField, textFields, textFormat, onSelect, emptyTip } = tab;

    let conList: JSX.Element[] = [];
    if (recordData.fetchState === FetchState.Fetching) {
      conList = [
        <li key={uuid()} style={{ textAlign: 'center' }}>
          {t('加载中...')}
        </li>
      ];
    } else if (recordData.fetchState === FetchState.Failed) {
      conList = [
        <li key={uuid()} style={{ textAlign: 'center' }}>
          <span className="text-danger">{t('加载失败')}</span>
        </li>
      ];
    } else {
      if (recordData.data.recordCount) {
        conList = recordData.data.records.map((item, index) => {
          let text =
            textFields && textFields.length && !!textFormat
              ? FormatText(item, textFields, textFormat)
              : item[textField];
          return (
            <li
              key={index}
              className={classnames({ cur: value === item[valueField] })}
              onClick={() => {
                onSelect(item[valueField]);
                this.stepNext();
              }}
            >
              <span className="text">{text}</span>
            </li>
          );
        });
      } else {
        conList = [
          <li key={uuid()} style={{ textAlign: 'center' }}>
            <span>{t('列表为空')}</span>
          </li>
        ];

        if (emptyTip) {
          return <div className="block-help-text">{emptyTip}</div>;
        }
      }
    }

    return (
      <div className="tab-panel">
        <ul className="cont-list">{conList}</ul>
      </div>
    );
  }

  private toggleShowTab() {
    this.setState({ isShowTab: !this.state.isShowTab });
  }
}
