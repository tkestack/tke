/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

import React, { useEffect, useState } from 'react';
import { Menu, StatusTip } from 'tea-component';
import { router } from '../../router';

const { LoadingTip } = StatusTip;

export const TellIsNeedFetchNS = (resourceName: string) => {
  return !['np', 'pv', 'sc', 'log', 'event'].includes(resourceName);
};

export const TellIsNotNeedFetchResource = (resourceName: string) => {
  return resourceName === 'info' ? true : false;
};

export const ResourceSidebarPanel = ({ route, actions, subRouterList }) => {
  const [type, setType] = useState(null);
  const [resourceName, setResourceName] = useState(null);

  useEffect(() => {
    const { type = null, resourceName = null } = router.resolve(route);

    setType(type);
    setResourceName(resourceName);
  }, [route]);

  function handleClick(_type, _resourceName) {
    if (resourceName === _resourceName) return;

    const urlParams = router.resolve(route);
    const queries = { ...(route?.queries ?? {}) };

    if (!['hpa', 'cronhpa'].includes(_resourceName)) {
      actions.resource.reset();
    }

    if (['hpa', 'cronhpa', 'info'].includes(_resourceName)) {
      delete queries.np;
    } else {
      // 这里去判断该资源是否需要进行namespace列表的拉取
      const isNeedFetchNamespace = TellIsNeedFetchNS(_resourceName);
      actions.resource.initResourceInfoAndFetchData(isNeedFetchNamespace, _resourceName);
      // 这里去清空多选的选项
      actions.resource.selectMultipleResource([]);
    }

    router.navigate(Object.assign({}, urlParams, { type: _type, resourceName: _resourceName }), {
      ...queries
    });
  }

  return (
    <Menu>
      {subRouterList?.length ? (
        subRouterList?.map(({ sub, ...item }) =>
          sub ? (
            <Menu.SubMenu
              title={item.name}
              key={item.path}
              opened={item.path === type}
              onOpenedChange={open => open && setType(item.path)}
            >
              {sub.map(subItem => (
                <Menu.Item
                  title={subItem.name}
                  key={subItem.path}
                  selected={subItem.path === resourceName}
                  onClick={() => handleClick(item.path, subItem.path)}
                />
              ))}
            </Menu.SubMenu>
          ) : (
            <Menu.Item
              title={item.name}
              key={item.path}
              selected={resourceName === item.basicUrl}
              onClick={() => handleClick(item.path, item.basicUrl)}
            />
          )
        )
      ) : (
        <div
          style={{
            position: 'absolute',
            width: '100%',
            top: 0,
            bottom: 0,
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            marginTop: 100
          }}
        >
          <LoadingTip />
        </div>
      )}
    </Menu>
  );
};
