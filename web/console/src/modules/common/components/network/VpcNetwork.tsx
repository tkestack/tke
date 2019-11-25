import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';
import { SelectList, SelectListProps } from '../select';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface VpcNetworkProps extends BaseReactProps {
  /**VPC列表 */
  vpc: SelectListProps;

  /**isShowCIDR */
  isShowCIDR?: boolean;

  /**mode当前的模式 */
  mode?: string;
}

export class VpcNetwork extends React.Component<VpcNetworkProps, {}> {
  render() {
    let { vpc, isShowCIDR, mode } = this.props,
      totalIPNum = 0,
      availableIPNum = 0,
      cidr = '';

    if (vpc.value) {
      if (isShowCIDR) {
        vpc.recordData.data.records.forEach(v => {
          if (v.unVpcId === vpc.value) {
            cidr = v.cidrBlock;
          }
        });
      }
    }

    return (
      <div>
        <SelectList {...vpc} name={t('集群网络')} className="tc-15-select m" style={{ display: 'inline-block' }} />
        {cidr && mode === 'create' && (
          <span className="inline-help-text text-weak" style={{ marginLeft: '5px' }}>
            CIDR: {cidr}
          </span>
        )}
      </div>
    );
  }
}
