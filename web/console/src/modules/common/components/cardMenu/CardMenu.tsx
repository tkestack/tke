/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import * as classnames from 'classnames';
import * as React from 'react';

import { BaseReactProps, insertCSS, OnOuterClick } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

/**插入自定义样式覆盖bubble样式 */
insertCSS(
  'CardMenuCss',
  `
.seal-area__item a {
  text-decoration: none;
  text-decoration-line: none;
  text-decoration-style: initial;
  text-decoration-color: initial;
  color:black
},

`
);

/**地域对应相应图标类名 */

const regionToIconClassName = {
  1: {
    area: 'china',
    icon: 'seal-area__Chinese'
  },
  4: {
    area: 'china',
    icon: 'seal-area__Chinese'
  },
  8: {
    area: 'china',
    icon: 'seal-area__Chinese'
  },
  5: {
    area: 'asia',
    icon: 'seal-area__Hongkong'
  },
  9: {
    area: 'asia',
    icon: 'seal-area__Singapore'
  },
  11: {
    area: 'finance',
    icon: 'seal-area__Chinese'
  },
  7: {
    area: 'finance',
    icon: 'seal-area__Chinese'
  },
  16: {
    area: 'china',
    icon: 'seal-area__Chinese'
  },
  19: {
    area: 'china',
    icon: 'seal-area__Chinese'
  },
  23: {
    area: 'asia',
    icon: 'seal-area__Thailand'
  },
  21: {
    area: 'asia',
    icon: 'seal-area__India'
  },
  18: {
    area: 'asia',
    icon: 'seal-area__Korea'
  },
  25: {
    area: 'asia',
    icon: 'seal-area__Japan'
  },
  15: {
    area: 'northAmerica',
    icon: 'seal-area__America'
  },
  22: {
    area: 'northAmerica',
    icon: 'seal-area__America'
  },
  17: {
    area: 'europe',
    icon: 'seal-area__Germany'
  },
  24: {
    area: 'europe',
    icon: 'seal-area__Russia'
  }
};

const RegionAreaMap = {
  china: t('中国大陆'),
  finance: t('金融专区'),
  asia: t('亚洲'),
  northAmerica: t('北美洲'),
  europe: t('欧洲')
};

export interface CardMenuItem {
  /* 文案 */
  name?: string | JSX.Element;

  /* 值 */
  value?: number | string;

  /* 是否可用 */
  disabled?: boolean;

  /* 提示 */
  tip?: string | JSX.Element;

  /**地域 */
  area?: string;
}

export interface CardMenuStateProps extends BaseReactProps {
  /* 列表数据 */
  list: CardMenuItem[];

  /* 选中的button */
  selected?: CardMenuItem;

  /* 选择后的回调 */
  onSelect?: (value: any) => void;
}

interface DownMenuState {
  /**
   * 当前是否为打开状态
   */
  isOpened?: boolean;
}
export class CardMenu extends React.Component<CardMenuStateProps, DownMenuState> {
  state = {
    isOpened: false
  };
  public select(item: CardMenuItem) {
    const { onSelect } = this.props;

    typeof onSelect === 'function' && onSelect(item);

    this.close();
  }

  @OnOuterClick
  public close() {
    this.setState({ isOpened: false });
  }

  render() {
    let { selected } = this.props;
    let chinaRegionList = [],
      financeRegionList = [],
      asiaRegionList = [],
      northRegionList = [],
      europeRegionList = [],
      { list } = this.props;
    list.forEach(item => {
      let areaName = regionToIconClassName[item.value].area;
      switch (areaName) {
        case 'china':
          chinaRegionList.push(item);
          break;
        case 'finance':
          financeRegionList.push(item);
          break;
        case 'asia':
          asiaRegionList.push(item);
          break;
        case 'northAmerica':
          northRegionList.push(item);
          break;
        case 'europe':
          europeRegionList.push(item);
          break;
      }
    });
    let selectItem = selected
      ? selected
      : {
          name: t('无可用地域'),
          value: ''
        };
    let selectdClassName = selected ? regionToIconClassName[selected.value].icon : '';
    return (
      <div className={classnames('seal-area', { 'is-select': this.state.isOpened })}>
        <div className="seal-area__head">
          <div className="seal-area__text">
            <a
              className="seal-area__select"
              style={{
                textDecoration: 'none',
                color: 'black',
                maxWidth: '200px'
              }}
              href="javascript:;"
              onClick={() =>
                this.setState({
                  isOpened: !this.state.isOpened
                })
              }
            >
              <span className={`seal-area__flag ${selectdClassName}`} />
              <span className="seal-area__name">{selectItem.name}</span>
              <i className="seal-area__arrow" />
            </a>
          </div>
        </div>
        <div className="seal-area__content">
          <div className="seal-area__list">
            <div className="seal-area__col">{this._rendArea(chinaRegionList, 'china')}</div>
            {financeRegionList.length !== 0 && (
              <div className="seal-area__col">
                {financeRegionList.length !== 0 && this._rendArea(financeRegionList, 'finance')}
              </div>
            )}
            {asiaRegionList.length !== 0 && (
              <div className="seal-area__col">{this._rendArea(asiaRegionList, 'asia')}</div>
            )}
            {northRegionList.length !== 0 && (
              <div className="seal-area__col">{this._rendArea(northRegionList, 'northAmerica')}</div>
            )}
            {europeRegionList.length !== 0 && (
              <div className="seal-area__col">{this._rendArea(europeRegionList, 'europe')}</div>
            )}
          </div>
        </div>
      </div>
    );
  }

  private _rendArea(List: CardMenuItem[], area: string) {
    return (
      <div className="seal-area__unit">
        <div className="seal-area__type">{RegionAreaMap[area]}</div>
        {List.map((item, index) => {
          let iconClassName = regionToIconClassName[item.value].icon;
          return (
            <a key={index} className="seal-area__item" onClick={item.disabled ? null : () => this.select(item)}>
              <span className={`seal-area__flag ${iconClassName}`} />
              <span className="seal-area__name">{item.name}</span>
            </a>
          );
        })}
      </div>
    );
  }
}
