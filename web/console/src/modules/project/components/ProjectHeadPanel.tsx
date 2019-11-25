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
