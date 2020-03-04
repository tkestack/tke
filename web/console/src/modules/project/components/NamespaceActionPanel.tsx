import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Justify, SearchBox } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../actions';
import { router } from '../router';
import { RootProps } from './ProjectApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class NamespaceActionPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, route } = this.props;
    actions.namespace.poll({ projectId: route.queries['projectId'] });
  }
  componentWillUnmount() {
    let { actions } = this.props;
    actions.namespace.clearPolling();
    actions.namespace.performSearch('');
  }
  render() {
    let { actions, namespace, route } = this.props;

    return (
      <div className="tc-action-grid">
        <Justify
          left={
            <Button
              type="primary"
              onClick={() => {
                router.navigate({ sub: 'createNS' }, route.queries);
              }}
            >
              {/* <b className="icon-add" /> */}
              {t('新建Namespace')}
            </Button>
          }
          right={
            <SearchBox
              value={namespace.query.keyword || ''}
              onChange={actions.namespace.changeKeyword}
              onSearch={actions.namespace.performSearch}
              placeholder={t('请输入Namespace名称')}
            />
          }
        />
      </div>
    );
  }
}
