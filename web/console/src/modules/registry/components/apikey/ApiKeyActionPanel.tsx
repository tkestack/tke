import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Justify, SearchBox } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../../actions';
import { RootProps } from '../RegistryApp';
import { CreateApiKeyPanel } from './CreateApiKeyPanel';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ApiKeyActionPanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, apiKey } = this.props;

    return (
      <div className="tc-action-grid">
        <Justify
          left={
            <React.Fragment>
              <Button
                type="primary"
                onClick={() => {
                  // route to apiKey create
                }}
              >
                {t('新建凭证')}
              </Button>
            </React.Fragment>
          }
          // right={
          //   <SearchBox
          //     value={this.props.apiKey.query.keyword || ''}
          //     onChange={this.props.actions.apiKey.changeKeyword}
          //     onSearch={this.props.actions.apiKey.performSearch}
          //     placeholder={t('请输入凭证描述')}
          //   />
          // }
        />
      </div>
    );
  }
}
