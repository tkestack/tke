import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../ChartGroupApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ActionPanel extends React.Component<RootProps, {}> {

  render() {
    const { actions, route, chartGroupList } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;

    return (
      <React.Fragment>
        <Table.ActionPanel>
          <Justify
            left={
              <Button type="primary" onClick={e => {
                e.preventDefault();
                router.navigate({ mode: 'create', sub: 'chartgroup' }, route.queries);
              }}>
                {t('新建')}
              </Button>
            }
            right={
              <React.Fragment>
                <SearchBox
                  value={chartGroupList.query.keyword || ''}
                  onChange={actions.chartGroup.list.changeKeyword}
                  onSearch={actions.chartGroup.list.performSearch}
                  onClear={() => {
                    actions.chartGroup.list.performSearch('');
                  }}
                  placeholder={t('请输入仓库名称')}
                />
              </React.Fragment>
            }
          />
        </Table.ActionPanel>
      </React.Fragment>
    );
  }
}
