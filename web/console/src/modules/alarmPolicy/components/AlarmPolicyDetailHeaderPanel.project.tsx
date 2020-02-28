import * as React from 'react';
import { RootProps } from './AlarmPolicyApp';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Justify } from '@tencent/tea-component';
export class AlarmPolicyDetailHeaderPanel extends React.Component<RootProps, {}> {
  goBack() {
    let { route } = this.props;
    // history.back();
    router.navigate(
      {},
      { clusterId: route.queries['clusterId'], projectName: route.queries['projectName'], np: route.queries['np'] }
    );
  }

  componentDidMount() {
    let { regionList, cluster, actions } = this.props;
    if (cluster.list.data.recordCount === 0) {
      actions.projectNamespace.initProjectList();
    } else {
      actions.alarmPolicy.initAlarmPolicyData();
    }
  }

  render() {
    let title = t('告警策略详情');

    return (
      <React.Fragment>
        <Justify
          left={
            <React.Fragment>
              <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
                <Icon type="btnback" />
                {t('返回')}
              </a>
              <h2>{title}</h2>
            </React.Fragment>
          }
        />
      </React.Fragment>
    );
  }
}
