import * as React from 'react';
import { Button, SearchBox, Justify } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../../actions';
import { RootProps } from '../RegistryApp';
import { CreateApiKeyPanel } from './CreateApiKeyPanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(
  state => state,
  mapDispatchToProps
)
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
