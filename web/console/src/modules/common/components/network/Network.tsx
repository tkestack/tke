import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { SelectList, SelectListProps } from '../select';

export interface NetworkProps extends BaseReactProps {
  /**VPC列表属性 */
  vpc: SelectListProps;

  /**子网列表属性 */
  subnet?: SelectListProps;

  /**isShowCIDR */
  isShowCIDR?: boolean;

  /**currentStep */
  currentStep?: number;

  /**mode当前的模式 */
  mode?: string;
}

export class Network extends React.Component<NetworkProps, {}> {
  render() {
    let { vpc, subnet, isShowCIDR, mode } = this.props,
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
      if ((mode === 'create' || mode === 'expand') && subnet.value) {
        let finder = subnet.recordData.data.records.find(v => v.unVpcId === vpc.value && v.unSubnetId === subnet.value);
        if (finder) {
          totalIPNum = finder.totalIPNum;
          availableIPNum = finder.availableIPNum;
        }
      }
    }

    return (
      <div>
        <SelectList {...vpc} name={t('集群网络')} className="tc-15-select m" style={{ display: 'inline-block' }} />
        {vpc.value && (mode === 'create' || mode === 'expand') && (
          <SelectList
            {...subnet}
            name={t('子网')}
            className="tc-15-select m"
            style={{ marginLeft: '5px', display: 'inline-block' }}
          />
        )}
        {vpc.value && subnet.value && (mode === 'create' || mode === 'expand') && (
          <span className="inline-help-text text-weak" style={{ marginLeft: '5px' }}>
            {t('共{{count}}个子网IP，剩{{availableIPNum}}个可用', {
              count: totalIPNum,
              availableIPNum
            })}
          </span>
        )}
        {cidr && mode === 'create' && (
          <p className="text-label" style={{ marginBottom: '0' }}>
            CIDR: {cidr}
          </p>
        )}
      </div>
    );
  }
}
