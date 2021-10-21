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
import { Card, List, Icon, Text, Button } from '@tencent/tea-component';
export function TipsPanel() {
  return (
    <Card>
      <Card.Header>
        <h3>实用提示</h3>
      </Card.Header>
      <Card.Body style={{ paddingTop: 0, paddingBottom: 0 }}>
        <List split={'divide'}>
          <List.Item>
            <img src="/static/icon/overviewBlack.svg" style={{ height: '30px' }} alt="logo" />
            <div style={{ display: 'inline-block', verticalAlign: 'bottom', marginLeft: 10 }}>
              <Text parent={'div'} style={{ display: 'block' }}>
                平台实验室
              </Text>
              <Text theme={'label'}>体验平台最新功能</Text>
              <Button
                type={'link'}
                onClick={() => {
                  location.href = location.origin + '/tkestack/uam/strategy/business';
                }}
              >
                立即体验
              </Button>
            </div>
          </List.Item>
          <List.Item>
            <img src="/static/icon/overviewBlack.svg" style={{ height: '30px' }} alt="logo" />
            <div style={{ display: 'inline-block', verticalAlign: 'bottom', marginLeft: 10 }}>
              <Text parent={'div'} style={{ display: 'block' }}>
                使用指引
              </Text>
              <Text theme={'label'}>通过创建业务，管理集群资源配额</Text>
              <Button
                type={'link'}
                onClick={() => {
                  location.href = location.origin + '/tkestack/project';
                }}
              >
                开始使用
              </Button>
            </div>
          </List.Item>
        </List>
      </Card.Body>
    </Card>
  );
}
