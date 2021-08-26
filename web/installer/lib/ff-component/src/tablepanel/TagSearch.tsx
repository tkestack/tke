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

import * as React from 'react';

import { FFListAction, FFListModel } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TagSearchBox, TagSearchBoxProps } from '@tencent/tea-component';

interface TablePanelTagSearchProps<TResource> extends TagSearchBoxProps {
  /**action */
  action?: FFListAction;
  /**列表 */
  model?: FFListModel<TResource>;
}

function TablePanelTagSearchBox<TResource = any>({ ...props }: TablePanelTagSearchProps<TResource>) {
  let {
    action,
    model: {
      query: { searchFilter }
    },
    attributes
  } = props;

  props.value =
    props.value !== undefined
      ? props.value
      : Object.keys(searchFilter)
          .filter(key => attributes.findIndex(attr => attr.key === key) !== -1 && searchFilter[key] !== null)
          .map(key => ({
            attr: {
              key: key
            },
            values: [
              {
                name: searchFilter[key]
              }
            ]
          }));

  props.onChange = props.onChange
    ? props.onChange
    : value => {
        let attrMap = {};
        value.forEach(item => {
          if (item.attr) {
            attrMap[item.attr.key] = item.values[0].name;
          } else {
            attrMap[attributes[0].key] = item.values[0].name;
          }
        });
        attributes.forEach(attr => {
          if (attrMap[attr.key] === undefined) {
            attrMap[attr.key] = null;
          }
        });

        let nextFilter = Object.assign({}, searchFilter, attrMap);
        action.applySearchFilter(nextFilter);
      };

  props.tips = props.tips ? props.tips : t('多个过滤标签用回车键分隔');
  props.hideHelp = props.hideHelp !== undefined ? props.hideHelp : true;
  props.minWidth = props.minWidth ? props.minWidth : '400px';

  return <TagSearchBox {...props}></TagSearchBox>;
}

export { TablePanelTagSearchBox, TablePanelTagSearchProps };
