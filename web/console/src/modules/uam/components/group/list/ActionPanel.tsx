import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../GroupApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ActionPanel extends React.Component<RootProps, {}> {

  render() {
    const { actions, route, groupList } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;

    return (
      <React.Fragment>
        <Table.ActionPanel>
          <Justify
            left={
              <Button type="primary" onClick={e => {
                e.preventDefault();
                router.navigate({ module: 'group', sub: 'create' }, route.queries);
              }}>
                {t('新建')}
              </Button>
            }
            right={
              <React.Fragment>
                <SearchBox
                  value={groupList.query.keyword || ''}
                  onChange={actions.group.list.changeKeyword}
                  onSearch={actions.group.list.performSearch}
                  onClear={() => {
                    actions.group.list.performSearch('');
                  }}
                  placeholder={t('请输入用户组名称')}
                />
              </React.Fragment>
            }
          />
        </Table.ActionPanel>
      </React.Fragment>
    );
  }
}
