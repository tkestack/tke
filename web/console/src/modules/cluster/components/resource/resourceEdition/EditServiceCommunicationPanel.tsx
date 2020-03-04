import * as React from 'react';

import { Bubble, Radio } from '@tea/component';
import { insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem } from '../../../../common/components';
import { CommunicationTypeList } from '../../../constants/Config';

insertCSS(
  'EditCommunicationPanel',
  `
.tc-15-radio-wrap:first-child{
    margin-left: 0px;
}
.form-ctrl-label, .form-ctrl-label-stacked, .tc-15-radio-wrap{
    font-size: 12px;
}
`
);

interface EdtiServiceCommunicationPanelProps {
  /** 访问方式的类型 */
  communicationType: string;

  /** communicationSelectAction */
  communicationSelectAction: (type: string) => void;

  /** 是否开启headless service */
  isOpenHeadless?: boolean;

  /** 是否禁止修改访问方式，只有在 update 并且 原本已经开通了 headless service  */
  isDisabledChangeCommunicationType?: boolean;

  /** 是否禁止勾选headless的选项，只有在 update */
  isDisabledToggleHeadless?: boolean;

  /** 操作headless 开关的操作 */
  toggleHeadlessAction: (isOpen: boolean) => void;

  /** isShow */
  isShow?: boolean;
}

export class EditServiceCommunicationPanel extends React.Component<EdtiServiceCommunicationPanelProps, {}> {
  render() {
    let {
      communicationType,
      communicationSelectAction,
      isShow = true,
      isOpenHeadless = false,
      toggleHeadlessAction,
      isDisabledChangeCommunicationType = false,
      isDisabledToggleHeadless = false
    } = this.props;

    let finder = CommunicationTypeList.find(c => c.value === communicationType),
      tip: any = finder ? finder.tip : '';

    return isShow ? (
      <FormItem label={t('服务访问方式')}>
        <div className="form-unit">
          <Radio.Group
            disabled={isDisabledChangeCommunicationType}
            value={communicationType}
            onChange={value => communicationSelectAction(value)}
          >
            {CommunicationTypeList.map((item, rIndex) => {
              return (
                <Radio key={rIndex} name={item.value}>
                  {item.label}
                </Radio>
              );
            })}
          </Radio.Group>
          {tip && <p className="text-label">{tip}</p>}
          {communicationType === 'ClusterIP' && (
            <div>
              <label
                className="form-ctrl-label"
                style={{ display: 'inline-block', marginRight: '5px', marginTop: '0' }}
              >
                <input
                  type="checkbox"
                  disabled={isDisabledToggleHeadless}
                  className="tc-15-checkbox"
                  checked={isOpenHeadless}
                  style={{ verticalAlign: 'middle' }}
                  onChange={() => {
                    toggleHeadlessAction(isOpenHeadless);
                  }}
                />
                Headless Service
              </label>
              <Bubble
                placement="bottom"
                content={t(
                  '不创建用于集群内访问的ClusterIP，访问Service名称时返回后端Pods IP地址，用于适配自有的服务发现机制。'
                )}
              >
                <i className="icon-help" />
              </Bubble>
              <span className="text-label" style={{ verticalAlign: 'middle', marginLeft: '5px' }}>
                <Trans>
                  <span>(Headless Service只支持创建时选择，</span>
                  <span className="text-danger">创建完成后不支持变更访问方式)</span>
                </Trans>
              </span>
            </div>
          )}
        </div>
      </FormItem>
    ) : (
      <noscript />
    );
  }
}
