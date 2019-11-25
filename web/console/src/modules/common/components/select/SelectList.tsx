import * as React from 'react';
import { FetcherState, FetchState } from '@tencent/qcloud-redux-fetcher';
import { BaseReactProps, RecordSet, insertCSS } from '@tencent/qcloud-lib';
import { DropdownListItem, DropdownList } from '../dropdown/';
import { RouteState } from '../../../../../helpers/Router';
import * as classnames from 'classnames';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Validation } from '../../models';

insertCSS(
  'Downdrop-head',
  `
.manage-area-title .tc-15-dropdown .dropdown-head{
    margin-top: 2px;
}
`
);
insertCSS(
  'Downdrop-head-cart',
  `
.manage-area-title .tc-15-dropdown .dropdown-head .caret{
    top: 3px;
    right:0px;
}
`
);

interface SelectItem {
  /**选项的值 */
  value?: string;

  /**选项显示文本 */
  text?: string | JSX.Element;

  /**是否选中 */
  selected?: boolean;

  key?: string | number;
}

export interface SelectListProps extends BaseReactProps {
  /**当前选中的项目 */
  value?: string;

  /**下拉显示的项目列表 */
  recordData?: FetcherState<RecordSet<any>>;

  /**下拉显示的项目列表 不带拉取状态*/
  recordList?: any[];

  /**指定value字段 */
  valueField?: string;

  /**指定text字段 */
  textField?: string;

  /**多字段组合显示 */
  textFields?: string[];

  /**多字段显示格式设置 替换变量格式为${val} */
  textFormat?: string | Function;

  /**名称, 在确定选择提示时使用 */
  name?: string;

  /**当前选中的项目 */
  selected?: string | number;

  /**下拉选中时触发的事件 */
  onSelect?: (value: string, regionId?: number) => void;

  /**重试 */
  onRetry?: () => void;

  /**样式类 */
  className?: string;

  /**行内样式 */
  style?: any;

  /**外层样式 */
  outerStyle?: React.CSSProperties;

  /**数据为空时显示的提示 */
  emptyTip?: string | JSX.Element;

  /**是否禁用 */
  disabled?: boolean;

  /**显示模式：select原生样式；dropdown下拉列表 */
  mode?: string;

  /**tip显示位置 */
  tipPosition?: 'top' | 'right' | 'bottom' | 'left';

  /**对齐方式 */
  align?: 'start' | 'end';

  /**校验态 */
  validator?: Validation;

  /**默认触发设置 */
  defaultEmit?: boolean;

  route?: RouteState;

  /**是否显示 */
  isShow?: boolean;

  /**是否有请选择这一项*/
  isUnshiftDefaultItem?: boolean;
}

const OptionItem = (data: SelectItem) => (
  <option key={data.key} value={data.value}>
    {data.text}
  </option>
);

const FormatText = (data: any, fields: string[], format: string) => {
  let result = format;
  fields.forEach(field => {
    let reg: RegExp = new RegExp('\\$\\{' + field + '\\}', 'g');
    result = result.replace(reg, data[field]);
  });

  return result;
};

export class SelectList extends React.Component<SelectListProps, {}> {
  getText(o, textFields: string[], textField: string, textFormat: string | Function) {
    let formatFunction: Function = null;
    if (typeof textFormat === 'function') {
      formatFunction = textFormat as Function;
    }
    let text = '';
    if (textFields && textFields.length) {
      if (formatFunction) {
        text = formatFunction(...textFields.map(field => o[field]));
      } else {
        text = FormatText(o, textFields, textFormat as string);
      }
    } else {
      text = o[textField];
    }
    return text;
  }

  _renderSelectList() {
    let {
        value,
        recordData,
        recordList,
        valueField,
        textField,
        textFields,
        textFormat,
        name = '',
        className,
        style,
        onSelect,
        emptyTip,
        disabled,
        tipPosition,
        align,
        isUnshiftDefaultItem = true
      } = this.props,
      isDisabled =
        (recordData &&
          (recordData.fetchState === FetchState.Fetching ||
            recordData.fetchState === FetchState.Failed ||
            recordData.data.recordCount === 0)) ||
        disabled,
      options = [];
    if (recordData && recordData.fetchState === FetchState.Fetching) {
      options = [<OptionItem key={-1} value="" text={t('加载中...')} />];
    }

    if (recordData && recordData.fetchState === FetchState.Failed) {
      options = [<OptionItem key={-1} value="" text={t('加载失败')} />];
    }

    if (recordData && recordData.fetchState === FetchState.Ready) {
      if (recordData.data.recordCount) {
        options = recordData.data.records.map((o, opIndex) => {
          let text = this.getText(o, textFields, textField, textFormat);
          // textFields && textFields.length && !!textFormat ? FormatText(o, textFields, textFormat) : o[textField];
          return <OptionItem key={opIndex} value={o[valueField]} text={text} />;
        });
        isUnshiftDefaultItem && options.unshift(<OptionItem key={-1} value="" text={t('请选择') + name} />);
      } else {
        options = [<OptionItem key={-1} value="" text={t('无')} />];

        if (emptyTip) {
          return <div className="block-help-text">{emptyTip}</div>;
        }
      }
    }

    if (!recordData && recordList.length) {
      options = recordList.map((o, index) => {
        let text = this.getText(o, textFields, textField, textFormat);
        // textFields && textFields.length && !!textFormat ? FormatText(o, textFields, textFormat) : o[textField];
        return <OptionItem key={index} value={o[valueField]} text={text} />;
      });
      isUnshiftDefaultItem && options.unshift(<OptionItem key={-1} value="" text={t('请选择') + name} />);
    } else if (!recordData && !recordList.length) {
      options.unshift(<OptionItem key={-1} value="" text={t('无')} />);
    }

    return (
      <select
        className={className}
        style={style}
        value={value}
        onChange={e => onSelect(e.target.value)}
        disabled={isDisabled}
      >
        {options}
      </select>
    );
  }

  _renderDropdownList() {
    let {
        value,
        recordData,
        valueField,
        textField,
        textFields,
        textFormat,
        name,
        className,
        style,
        onSelect,
        onRetry,
        emptyTip,
        disabled,
        mode,
        tipPosition,
        align
      } = this.props,
      isDisabled =
        recordData.fetchState === FetchState.Fetching ||
        recordData.fetchState === FetchState.Failed ||
        recordData.data.recordCount === 0 ||
        disabled,
      options = [],
      selected: DropdownListItem;
    if (recordData.fetchState === FetchState.Fetching) {
      options = [<OptionItem key={-1} value="" text={t('加载中...')} />];
    }

    if (recordData.fetchState === FetchState.Fetching) {
      return (
        <div className="tc-15-dropdown tc-15-dropdown-hd">
          <span className="tc-15-dropdown-link">
            <i className="n-loading-icon" />
            {t('集群加载中...')}
          </span>
        </div>
      );
    }

    if (recordData.fetchState === FetchState.Failed) {
      return (
        <div className="tc-15-dropdown tc-15-dropdown-hd">
          <span className="tc-15-dropdown-link">
            <Trans>
              数据获取失败，
              <a href="javascript:;" onClick={e => onRetry()}>
                请重试
              </a>
            </Trans>
          </span>
        </div>
      );
    }

    if (recordData.fetchState === FetchState.Ready) {
      if (recordData.data.recordCount) {
        recordData.data.records.forEach(item => {
          let text = this.getText(item, textFields, textField, textFormat);
          // textFields && textFields.length && !!textFormat
          //   ? FormatText(item, textFields, textFormat)
          //   : item[textField];
          let option: DropdownListItem = {
            id: item[valueField],
            label: text
          };
          options.push(option);
          if (value === item[valueField]) {
            selected = option;
          }
        });
      } else {
        let option: DropdownListItem = {
          id: 'empty',
          label: t('无')
        };
        options.push(option);
        selected = option;
      }
    }

    return recordData.data.recordCount || !emptyTip ? (
      <DropdownList
        items={options}
        smallSize
        simulateSelect
        selected={selected}
        onSelect={item => {
          onSelect(item.id as string);
        }}
        buttonMaxWidth="200px"
        menuMaxWidth="200px"
        theme="dropdown-hd"
        className="tc-15-dropdown tc-15-dropdown-in-hd"
      />
    ) : (
      <div className="block-help-text">{emptyTip}</div>
    );
  }

  render() {
    let { mode, validator, style, outerStyle, isShow = true } = this.props;
    return isShow ? (
      <div
        style={style ? style : outerStyle ? outerStyle : style}
        className={classnames('form-unit', {
          'is-error': validator && validator.status === 2
        })}
      >
        {mode === 'dropdown' ? this._renderDropdownList() : this._renderSelectList()}
        {validator && validator.status === 2 && <p className="form-input-help">{validator.message}</p>}
      </div>
    ) : (
      <noscript />
    );
  }

  componentWillReceiveProps(nextProps) {
    let { value, recordData, recordList, onSelect, defaultEmit, valueField, route } = nextProps;
    /**默认触发选中数据 */
    let defaultValue = value || (route && route.queries[valueField]);
    if (recordList && recordList.length && defaultEmit && !this.props.recordList.length) {
      defaultValue = defaultValue || recordList[0][valueField];
      onSelect(defaultValue);
    }

    if (recordData && recordData.data.recordCount && defaultEmit && !this.props.recordData.data.recordCount) {
      defaultValue = defaultValue || recordData.data.records[0][valueField];
      onSelect(defaultValue);
    }
  }
}
