import * as React from 'react';
import { RootProps } from '../../HelmApp';
import classNames from 'classnames';
import { CommonBar, FormItem } from '../../../../common/components';
import { OtherType, otherTypeList } from '../../../constants/Config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export class OtherChartPanel extends React.Component<RootProps, {}> {
  onChangeChartUrl(chart_url: string) {
    let { actions } = this.props;

    actions.create.inputOtherChartUrl(chart_url);
  }
  onSelectType(type: string) {
    this.props.actions.create.selectOtherType(type);
  }
  onChangeUserName(username: string) {
    let {
      actions,
      helmCreation: { isValid }
    } = this.props;
    actions.create.inputOtherUserName(username);
  }
  onChangePassword(password: string) {
    let {
      actions,
      helmCreation: { isValid }
    } = this.props;
    actions.create.inputOtherPassword(password);
  }

  render() {
    let {
      helmCreation: { otherChartUrl, otherTypeSelection, otherUserName, otherPassword, isValid }
    } = this.props;

    return (
      <ul className="form-list" style={{ marginBottom: 16 }}>
        <FormItem label="Chart_Url">
          <div
            className={classNames('form-unit', {
              'is-error': isValid.otherChartUrl !== ''
            })}
          >
            <input
              type="text"
              className="tc-15-input-text m"
              placeholder={t('请输入Chart_Url')}
              value={otherChartUrl}
              onChange={e => this.onChangeChartUrl((e.target.value + '').trim())}
            />
            <p className="form-input-help">{isValid.otherChartUrl}</p>
          </div>
        </FormItem>
        <FormItem label={t('类型')}>
          <div className="form-unit">
            <div className="tc-15-rich-radio">
              <CommonBar
                list={otherTypeList}
                value={otherTypeSelection}
                onSelect={item => {
                  this.onSelectType(item.value as string);
                }}
              />
            </div>
          </div>
        </FormItem>

        {otherTypeSelection === OtherType.Private && (
          <FormItem label=" ">
            <ul className="form-list">
              <FormItem label={t('用户名')}>
                <div
                  className={classNames('form-unit', {
                    'is-error': isValid.otherUserName !== ''
                  })}
                >
                  <input
                    type="text"
                    className="tc-15-input-text m"
                    placeholder={t('请输入用户名')}
                    value={otherUserName}
                    onChange={e => this.onChangeUserName(e.target.value + '')}
                  />
                  <p className="form-input-help">{isValid.otherUserName}</p>
                </div>
              </FormItem>
              <FormItem label={t('密码')}>
                <div
                  className={classNames('form-unit', {
                    'is-error': isValid.otherPassword !== ''
                  })}
                >
                  <input
                    type="text"
                    className="tc-15-input-text m"
                    placeholder={t('请输入密码')}
                    value={otherPassword}
                    onChange={e => this.onChangePassword(e.target.value + '')}
                  />
                  <p className="form-input-help">{isValid.otherPassword}</p>
                </div>
              </FormItem>
            </ul>
          </FormItem>
        )}
      </ul>
    );
  }
}
