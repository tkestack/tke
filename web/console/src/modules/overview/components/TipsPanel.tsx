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
          <List.Item style={{ verticalAlign: 'center' }}>
            <img src="/static/icon/overviewBlack.svg" style={{ height: '30px' }} alt="logo" />
            <div style={{ display: 'inline-block', verticalAlign: 'middle' }}>
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
          <List.Item style={{ verticalAlign: 'center' }}>
            <img src="/static/icon/overviewBlack.svg" style={{ height: '30px' }} alt="logo" />
            <div style={{ display: 'inline-block', verticalAlign: 'middle' }}>
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
