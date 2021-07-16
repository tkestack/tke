import React from 'react';
import { Card, Layout } from 'tea-component';

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
      <Content>
        <Content.Header showBackButton onBackButtonClick={goBack} title={title} style={{ backgroundColor: '#fff' }} />
        <Content.Body style={{ padding: '20px' }}>
          <Card style={{ maxWidth: '1360px', margin: '0 auto' }}>
            <Card.Body>{children}</Card.Body>
          </Card>
        </Content.Body>
      </Content>
      <Footer>
        <Card style={{ maxWidth: '1360px', margin: '0 auto' }}>
          <Card.Body>{footer}</Card.Body>
        </Card>
      </Footer>
    </Layout>
  );
}
