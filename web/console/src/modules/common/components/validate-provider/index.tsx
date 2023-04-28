import React from 'react';
import { Bubble } from 'tea-component';

interface Props {
  status?: 'validating' | 'error' | 'success';
  message?: string;
  children: React.ReactElement;
}

export const ValidateProvider = ({ status, message, children }: Props) => {
  const classNames = status === 'error' ? 'is-error' : undefined;

  return (
    <Bubble content={message}>
      <span className={classNames}>{children}</span>
    </Bubble>
  );
};
