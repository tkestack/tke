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
import { connect } from 'react-redux';

import { Bubble, Button, Justify, Table } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../../../actions';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceYamlActionPanel extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let disableBtn = false,
      errorTip = '';

    let ns = route.queries['np'],
      resourceIns = route.queries['resourceIns'];
    if (ns === 'kube-system') {
      disableBtn = true;
      errorTip = t('当前命名空间下的资源不可编辑');
    } else if (urlParams['resourceName'] === 'svc' && ns !== 'kube-system') {
      disableBtn = resourceIns === 'kubernetes';
      errorTip = t('系统默认的Service不可编辑');
    }

    return (
      <Table.ActionPanel>
        <Justify
          left={
            <Bubble placement="left" content={disableBtn ? errorTip : null}>
              <Button
                type="primary"
                disabled={disableBtn}
                onClick={() => {
                  if (!disableBtn) {
                    router.navigate(
                      Object.assign({}, urlParams, { mode: 'modify' }),
                      Object.assign({}, route.queries, { resourceIns })
                    );
                  }
                }}
              >
                {t('编辑YAML')}
              </Button>
            </Bubble>
          }
        />
      </Table.ActionPanel>
    );
  }
}
