import React from 'react';
import { Card, Layout, PageHeader } from 'antd';

const { Content, Footer } = Layout;

export interface AntdLayoutProps {
  title: string;
  footer: React.ReactNode;
  children?: React.ReactNode;
  goBack?(): void;
}

export function AntdLayout({ title, footer, children, goBack = () => history.back() }: AntdLayoutProps) {
  return (
    <Layout style={{ minHeight: '100%' }}>
      <PageHeader onBack={goBack} title={title} style={{ backgroundColor: '#fff' }} />
      <Content style={{ padding: '20px' }}>
        <Card style={{ maxWidth: '1360px', margin: '0 auto' }}>{children}</Card>
      </Content>
      <Footer>
        <Card style={{ maxWidth: '1360px', margin: '0 auto' }}>{footer}</Card>
      </Footer>
    </Layout>
  );
}
