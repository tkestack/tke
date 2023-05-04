import { Card, Layout } from '@tencent/tea-component';
import React from 'react';

const { Content, Footer } = Layout;

export interface TeaFormLayoutProps {
  title: React.ReactNode;
  footer?: React.ReactNode;
  children?: React.ReactNode;
  goBack?(): void;
  full?: boolean;
  wrapCard?: boolean;
}

export function TeaFormLayout({
  title,
  footer,
  children,
  goBack = () => history.back(),
  full = false,
  wrapCard = true
}: TeaFormLayoutProps) {
  return (
    <Layout>
      <Content>
        <Content.Header showBackButton onBackButtonClick={goBack} title={title} />
        <Content.Body full={full}>
          {wrapCard ? (
            <Card>
              <Card.Body>{children}</Card.Body>
            </Card>
          ) : (
            children
          )}
        </Content.Body>

        {footer && (
          <Content.Footer>
            <Content.Body full={full}>
              {wrapCard ? (
                <Card>
                  <Card.Body>{footer}</Card.Body>
                </Card>
              ) : (
                footer
              )}
            </Content.Body>
          </Content.Footer>
        )}
      </Content>
    </Layout>
  );
}
