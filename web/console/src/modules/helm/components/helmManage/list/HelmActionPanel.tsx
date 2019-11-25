import * as React from 'react';

import { Button, Table, Bubble } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../../HelmApp';
import { router } from '../../../router';
import { ClusterHelmStatus } from '../../../constants/Config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component/lib/justify';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(
  state => state,
  mapDispatchToProps
)
export class HelmActionPanel extends React.Component<RootProps, {}> {
  render() {
    let {
      actions,
      listState: { helmList, helmQuery, region, clusterHelmStatus },
      route
    } = this.props;
    return (
      <Table.ActionPanel>
        <Justify
          left={
            <Bubble
              placement="right"
              content={clusterHelmStatus.code !== ClusterHelmStatus.RUNNING ? '请先开通Helm应用' : null}
            >
              <Button
                type="primary"
                disabled={clusterHelmStatus.code !== ClusterHelmStatus.RUNNING}
                onClick={() => router.navigate({ sub: 'create' }, route.queries)}
              >
                {t('新建')}
              </Button>
            </Bubble>
          }
          right={<React.Fragment />}
        />
      </Table.ActionPanel>
    );
  }
}
