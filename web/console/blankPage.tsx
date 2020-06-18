import React, { useState } from 'react';
import { Button } from '@tea/component/button';
import { Blank, BlankTheme } from '@tea/component/blank';
import { Card } from '@tea/component/card';
import { Layout } from '@tea/component/layout';
import { ContentView } from '@tencent/tea-component';

const { Body, Content } = Layout;

function BlankExample() {
  return (
    <ContentView>
      <ContentView.Body>
        <Card full bordered={false}>
          <Blank
            theme={'permission'}
            description="您无所属项目或管理的集群为空，请联系相关管理员进行权限管理。"
            // operation={
            //   <>
            //     <Button type="primary">前往访问管理</Button>
            //     <Button>了解访问管理</Button>
            //   </>
            // }
            // extra={
            //   <>
            //     <ExternalLink>查看接入流程</ExternalLink>
            //     <ExternalLink>查看计价</ExternalLink>
            //   </>
            // }
          />
        </Card>
      </ContentView.Body>
    </ContentView>
  );
}

export function BlankPage() {
  return (
    <div
      style={{
        margin: '40px'
      }}
    >
      <BlankExample />
    </div>
  );
}
