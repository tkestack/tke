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
