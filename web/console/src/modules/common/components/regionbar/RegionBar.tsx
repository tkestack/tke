import * as React from 'react';
import { FetcherState, FetchState } from '@tencent/qcloud-redux-fetcher';
import { RecordSet } from '@tencent/qcloud-lib';
import { ButtonBar, DownMenu } from '../../components';
import { isEmpty } from '../../utils';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { CardMenuItem, CardMenu } from '../cardMenu';

export interface RegionBarProps {
  /* 列表数据 */
  recordData?: FetcherState<RecordSet<any>>;

  /**列表数据不带状态 */
  recordList?: any[];

  /**合并展示 */
  merge?: string[];

  /**地域配额数据 */
  quotaData?: FetcherState<any>;

  /* 选中的button */
  value?: string | number;

  /** 当存在region时，当作渲染可用区，默认渲染地域 */
  region?: string | number;

  /* 选择后的回调 */
  onSelect?: (value: any) => void;

  /**默认触发回调 用于最开始初始化联动 */
  defaultEmit?: boolean;

  /** 单个list的时候是否显示纯文本 */
  isNeedPureText?: boolean;

  /**
   * 选择downmenu的话回调格式item => {selectRegion(item.id)}
   * 选择buttonBar,cardMenu的话回调格式item => {selectRegion(item.value)}
   */
  mode?: 'downMenu' | 'buttonBar' | 'cardMenu';
}

export class RegionBar extends React.Component<RegionBarProps, {}> {
  _rendRegionItem() {
    const { recordData, quotaData, value, region, merge, recordList, mode } = this.props;
    let list: CardMenuItem[] = [],
      selectItem: CardMenuItem;

    /**设置区域合并 */
    if (merge && merge.length) {
      let barItem: CardMenuItem = {
        name:
          mode === 'buttonBar' ? <span>{t('默认地域')}</span> : t('默认地域(包括广州、上海、北京、成都、新加坡、重庆)'),
        value: '1',
        tip: t('默认地域包括广州、上海、北京、成都、新加坡、重庆'),
        disabled: false
      };
      if (merge.indexOf(value + '') > -1) {
        selectItem = barItem;
      }
      list.push(barItem);
    }
    let dataList = recordData ? recordData.data.records : recordList;
    dataList.forEach(item => {
      let sumQuota = 0;
      if (quotaData && !isEmpty(quotaData.data)) {
        if (region) {
          let finder = quotaData.data[region].find(quota => quota.zoneId === item.id);
          //若返回值中没有相应可用区CPU核数，则视为无上限
          sumQuota = finder ? finder.cpuQuota : 65535;
        } else {
          quotaData.data[item.value].forEach(item => {
            if (item.cpuQuota === undefined) {
              sumQuota += 65535;
            } else if (item.cpuQuota === undefined) {
              sumQuota += item.cpuQuota;
            }
          });
          if (!quotaData.data[item.value].keys().length) {
            sumQuota += 65535;
          }
        }
      } else {
        //无上限
        sumQuota = 65535;
      }

      if (!merge || (merge && merge.indexOf(item.value + '') === -1)) {
        let barItem: CardMenuItem = {
          name:
            mode === 'buttonBar' ? (
              <span>
                {item.name || item.regionChineseName}
                {!sumQuota && <i className="shop-over-icon">{t('售罄')}</i>}
              </span>
            ) : (
              item.name || item.regionChineseName
            ), //兼容黑石集群
          value: item.value || item.id,
          tip: item.tip,
          disabled: item.disabled || !sumQuota,
          area: item.value ? item.area : null
        };
        if (value + '' === item.value + '' || value + '' === item.id + '') {
          selectItem = barItem;
        }

        list.push(barItem);
      }
    });
    return { list, selectItem };
  }

  render() {
    let { recordData, quotaData, onSelect, isNeedPureText, mode, merge, recordList } = this.props,
      element: JSX.Element;

    let faildInfo: JSX.Element = (
      <p className="text-danger" style={{ fontSize: '12px' }}>
        <i className="n-error-icon" /> {t('加载失败')}
      </p>
    );
    let loadingInfo: JSX.Element = (
      <p className="text" style={{ fontSize: '12px' }}>
        <i className="n-loading-icon" /> {t('加载中...')}
      </p>
    );
    let emptyInfo: JSX.Element = (
      <p className="text" style={{ fontSize: '12px' }}>
        {t('无可用数据')}
      </p>
    );

    let headerfaildInfo: JSX.Element = (
      <div
        className="tc-15-dropdown"
        style={{ marginLeft: '20px', display: 'inline-block', minWidth: '30px', color: 'red' }}
      >
        <i className="n-error-icon" />
        {t('加载失败')}
      </div>
    );
    let headerloadingInfo: JSX.Element = (
      <div className="tc-15-dropdown" style={{ marginLeft: '20px', display: 'inline-block', minWidth: '30px' }}>
        <i className="n-loading-icon" />
        {t('加载中...')}
      </div>
    );
    if (recordData) {
      switch (recordData.fetchState) {
        case FetchState.Failed:
          if (mode === 'downMenu' || mode === 'cardMenu') {
            element = headerfaildInfo;
          } else {
            element = faildInfo;
          }
          break;
        case FetchState.Fetching:
          if (mode === 'downMenu' || mode === 'cardMenu') {
            element = headerloadingInfo;
          } else {
            element = loadingInfo;
          }
          break;
        case FetchState.Ready:
          if (quotaData) {
            switch (quotaData.fetchState) {
              case FetchState.Failed:
                element = faildInfo;
                break;
              case FetchState.Fetching:
                element = loadingInfo;
                break;
              case FetchState.Ready:
                if (recordData.data.recordCount === 0) {
                  element = emptyInfo;
                } else if (recordData.data.recordCount === 1 && isNeedPureText !== false) {
                  element = (
                    <p className="text" style={{ fontSize: '12px' }}>
                      {recordData.data.records[0].name}
                    </p>
                  );
                } else {
                  /* eslint-disable */
                  let { list, selectItem } = this._rendRegionItem();
                  if (mode === 'downMenu') {
                    element = <DownMenu list={list} selected={selectItem} onSelect={onSelect} />;
                  } else if (mode === 'cardMenu') {
                    element = <CardMenu list={list} selected={selectItem} onSelect={onSelect} />;
                  } else {
                    element = (
                      <ButtonBar
                        isNeedPureText={!!isNeedPureText}
                        list={list}
                        size="m"
                        selected={selectItem}
                        onSelect={onSelect}
                        bubbleDirection={merge ? 'left' : 'bottom'}
                      />
                    );
                  }
                  /* eslint-enable */
                }
            }
          } else {
            let { list, selectItem } = this._rendRegionItem();
            if (mode === 'downMenu') {
              element = <DownMenu list={list} selected={selectItem} onSelect={onSelect} />;
            } else if (mode === 'cardMenu') {
              element = <CardMenu list={list} selected={selectItem} onSelect={onSelect} />;
            } else {
              element = (
                <ButtonBar
                  isNeedPureText={!!isNeedPureText}
                  list={list}
                  size="m"
                  selected={selectItem}
                  onSelect={onSelect}
                  bubbleDirection={merge ? 'left' : 'bottom'}
                />
              );
            }
          }
      }
    } else if (recordList) {
      let { list, selectItem } = this._rendRegionItem();
      if (mode === 'downMenu') {
        element = <DownMenu list={list} selected={selectItem} onSelect={onSelect} />;
      } else if (mode === 'cardMenu') {
        element = <CardMenu list={list} selected={selectItem} onSelect={onSelect} />;
      } else {
        element = (
          <ButtonBar
            isNeedPureText={!!isNeedPureText}
            list={list}
            size="m"
            selected={selectItem}
            onSelect={onSelect}
            bubbleDirection={merge ? 'left' : 'bottom'}
          />
        );
      }
    }
    return element;
  }

  componentWillReceiveProps(nextProps) {
    let { value, recordData, onSelect, defaultEmit, recordList } = nextProps;
    if (value && recordData && recordData.data.recordCount && defaultEmit && !this.props.recordData.data.recordCount) {
      onSelect({ value });
    } else if (value && recordList && recordList.length && defaultEmit && !recordList.length) {
      onSelect({ value });
    }
  }
}
