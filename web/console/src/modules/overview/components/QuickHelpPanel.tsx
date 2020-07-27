import * as React from 'react';
import { Card, List, Button } from '@tencent/tea-component';
export function QuickHelpPanel() {
  return (
    <Card>
      <Card.Header>
        <h3>快速入口</h3>
      </Card.Header>
      <Card.Body style={{ paddingTop: 0, paddingBottom: 0 }}>
        <List split={'divide'}>
          <List.Item style={{ verticalAlign: 'center' }}>
            <img src="/static/icon/overviewCluster.svg" style={{ height: '30px' }} alt="logo" />
            <Button
              type={'link'}
              onClick={() => {
                location.href = location.origin + '/tkestack/cluster/createIC';
              }}
            >
              创建独立集群
            </Button>
          </List.Item>
          <List.Item style={{ verticalAlign: 'center' }}>
            <img src="/static/icon/overviewUser.svg" style={{ height: '30px' }} alt="logo" />
            <Button
              type={'link'}
              onClick={() => {
                location.href = location.origin + '/tkestack/uam/user/normal/create';
              }}
            >
              创建角色
            </Button>
          </List.Item>
          <List.Item style={{ verticalAlign: 'center' }}>
            <img src="/static/icon/overviewGithub.svg" style={{ height: '30px' }} alt="logo" />
            <Button
              type={'link'}
              onClick={() => {
                location.href = 'https://github.com/tkestack/tke/tree/master/docs/guide/zh-CN';
              }}
            >
              github-issue
            </Button>
          </List.Item>
        </List>
      </Card.Body>
    </Card>
  );
}
