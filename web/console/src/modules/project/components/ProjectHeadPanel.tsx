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
import { RootProps } from './ProjectApp';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify, Button } from '@tencent/tea-component';
import { router } from '../router';

export class ProjectHeadPanel extends React.Component<{ title: string; isNeedBack?: boolean }, {}> {
  goBack() {
    history.back();
  }
  render() {
    let { title, isNeedBack = false } = this.props;

    return (
      <Justify
        left={
          <React.Fragment>
            {isNeedBack && (
              <Button
                icon={'btnback'}
                onClick={() => {
                  this.goBack();
                }}
              ></Button>
            )}
            <h2>{title}</h2>
          </React.Fragment>
        }
      />
    );
  }
}
