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
          <List.Item>
            <img
              src="/static/icon/overviewCluster.svg"
              style={{ height: '30px', verticalAlign: 'middle', marginRight: 10 }}
              alt="logo"
            />
            <Button
              type={'link'}
              onClick={() => {
                location.href = location.origin + '/tkestack/cluster/createIC';
              }}
            >
              创建独立集群
            </Button>
          </List.Item>
          <List.Item>
            <img
              src="/static/icon/overviewUser.svg"
              style={{ height: '30px', verticalAlign: 'middle', marginRight: 10 }}
              alt="logo"
            />
            <Button
              type={'link'}
              onClick={() => {
                location.href = location.origin + '/tkestack/uam/user/normal/create';
              }}
            >
              创建角色
            </Button>
          </List.Item>
          <List.Item>
            <img
              src="/static/icon/overviewGithub.svg"
              style={{ height: '30px', verticalAlign: 'middle', marginRight: 10 }}
              alt="logo"
            />
            <Button
              type={'link'}
              onClick={() => {
                location.href = 'https://github.com/tkestack/tke/issues';
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
