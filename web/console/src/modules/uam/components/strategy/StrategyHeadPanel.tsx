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
import * as React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Justify, Icon } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../router';

export const StrategyHeadPanel = () => {
  const state = useSelector(state => state);
  const { route } = state;
  let urlParam = router.resolve(route);
  const { sub } = urlParam;
  return (
    <Justify
      left={
        <h2>
          {sub ? (
            <React.Fragment>
              <a href="javascript:history.go(-1);">
                <Icon type="btnback" />
              </a>
              <span style={{ marginLeft: '10px' }}>{sub}</span>
            </React.Fragment>
          ) : (
            <Trans>策略管理</Trans>
          )}
        </h2>
      }
    />
  );
};
