import { FetchState } from '@tencent/ff-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { t } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, Justify, Table, Text } from '@tencent/tea-component';
import * as React from 'react';
import { connect } from 'react-redux';
import { allActions } from '../actions';
import { router } from '../router';
import { RootProps } from './AddonApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class AddonActionPanel extends React.Component<RootProps, any> {
  render() {
    let { actions, cluster, route, openAddon } = this.props,
      urlParams = router.resolve(route);

    let isCanNotAdd =
      openAddon.list.data.recordCount === 0 &&
      (openAddon.list.fetched !== true || openAddon.list.fetchState === FetchState.Fetching);

    let errorTips: string = '';

    if (cluster.selection && cluster.selection.status.phase !== 'Running') {
      isCanNotAdd = true;
      errorTips = '当前集群状态不正常';
    } else {
      errorTips = '暂未选择集群';
    }

    return (
      <Table.ActionPanel>
        <Justify
          left={
            <Bubble placement="right" content={isCanNotAdd ? <Text>{errorTips}</Text> : null}>
              <Button
                type="primary"
                disabled={isCanNotAdd}
                onClick={() => {
                  // 跳转到新建界面
                  router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
                }}
              >
                {t('新建')}
              </Button>
            </Bubble>
          }
        />
      </Table.ActionPanel>
    );
  }
}
