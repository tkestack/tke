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
import { connect } from 'react-redux';
import { RootProps } from '../AppContainer';
import { Justify, Icon } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { router } from '../../../router';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class HeaderPanel extends React.Component<RootProps, {}> {
  goBack = () => {
    let { actions, route } = this.props,
      urlParams = router.resolve(route);
    router.navigate({ mode: 'list', sub: 'app' }, route.queries);
  };

  render() {
    let { route } = this.props;
    let title = route.queries['appName'] + '(' + route.queries['namespace'] + ')';

    return (
      <Justify
        left={
          <React.Fragment>
            <div className="manage-area-title secondary-title">
              <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
                <Icon type="btnback" />
                {t('返回')}
              </a>
              <span className="line-icon">|</span>
              <h2>{title}</h2>
            </div>
          </React.Fragment>
        }
      />
    );
  }
}
