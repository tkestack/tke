import * as React from 'react';
import { Button, SearchBox, Justify } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../actions';
import { RootProps } from './ProjectApp';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../router';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
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
