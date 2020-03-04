import { FormPanelSelect, FormPanelSelectProps } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';
import * as React from 'react';
import { RootProps } from './AlarmPolicyApp';
export class AlarmPolicyHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    // actions.region.fetch();
    actions.cluster.applyFilter({ regionId: 1 });
  }

  render() {
    let { actions, regionList, regionSelection, cluster } = this.props;

    let selectProps: FormPanelSelectProps = {
      type: 'simulate',
      appearence: 'button',
      label: '集群',
      model: cluster,
      action: actions.cluster,
      valueField: record => record.metadata.name,
      displayField: record => `${record.metadata.name} (${record.spec.displayName})`,
      onChange: (clusterId: string) => {
        actions.cluster.selectCluster(cluster.list.data.records.find(c => c.metadata.name === clusterId));
      }
    };
    return (
      <Justify
        left={
          <div style={{ lineHeight: '28px' }}>
            <h2 style={{ float: 'left' }}>{t('告警设置')}</h2>
            <div className="tc-15-dropdown" style={{ marginLeft: '20px', display: 'inline-block', minWidth: '30px' }}>
              {t('集群')}
            </div>
            <FormPanelSelect {...selectProps} />
          </div>
        }
      />
    );
  }
}
