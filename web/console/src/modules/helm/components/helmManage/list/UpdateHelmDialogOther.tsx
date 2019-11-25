import classNames from 'classnames';
import * as React from 'react';

import { Button, Modal } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { CommonBar, FormItem, TipInfo } from '../../../../common/components/';
import { OtherType, otherTypeList } from '../../../constants/Config';
import { RootProps } from '../../HelmApp';
import { KeyValuesPanel } from '../create/KeyValuesPanel';

interface Props extends RootProps {
  onCancel?: Function;
}

export class UpdateHelmDialogOther extends React.Component<Props, {}> {
  componentDidMount() {
    this.props.actions.helm.inputKeyValue([]);
  }
  componentWillUnmount() {
    // this.props.actions.create.clear();
    // this.props.actions.create.inputKeyValue([]);
    this.props.actions.helm.inputKeyValue([]);
  }
  onChangeChartUrl(chart_url: string) {
    this.props.actions.helm.inputOtherChartUrl(chart_url);
  }
  onSelectType(type: string) {
    this.props.actions.helm.selectOtherType(type);
  }
  onChangeUserName(username: string) {
    this.props.actions.helm.inputOtherUserName(username);
  }
  onChangePassword(password: string) {
    this.props.actions.helm.inputOtherPassword(password);
  }

  isCanSave() {
    let {
      listState: { isValid, otherTypeSelection }
    } = this.props;
    let canSave = true;
    if (isValid.otherChartUrl !== '') {
      canSave = false;
    } else {
      if (otherTypeSelection === OtherType.Private) {
        if (isValid.otherUserName !== '' || isValid.otherPassword !== '') {
          canSave = false;
        }
      }
    }
    return canSave;
  }
  render() {
    const { actions } = this.props;
    const cancel = () => {
      this.props.onCancel && this.props.onCancel();
    };
    const confirm = () => {
      actions.helm.validAll();
      let canSave = this.isCanSave();
      if (canSave) {
        actions.helm.updateHelm();
        cancel();
      }
    };

    const {
      listState: { helmSelection, otherChartUrl, otherTypeSelection, otherUserName, otherPassword, kvs, isValid }
    } = this.props;

    return (
      <Modal visible={true} caption={t('更新Helm应用')} onClose={cancel} size={600} disableEscape={true}>
        <Modal.Body>
          <TipInfo>
            <span style={{ verticalAlign: 'middle' }}>
              {t(
                '注意，若您重新填写了任意变量，将覆盖应用下所有自定义变量。不填写变量时，将会使用上次填写的变量更新应用。'
              )}
            </span>
          </TipInfo>
          <ul className="form-list" style={{ paddingBottom: 20 }}>
            <FormItem label={t('应用名称')}>
              <span className="form-input">{helmSelection.name}</span>
            </FormItem>
            <FormItem label={t('Chart名称')}>
              <span className="form-input">{helmSelection.chart_metadata.name}</span>
            </FormItem>
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
                  onChange={e => this.onChangeChartUrl(e.target.value + '')}
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
          <KeyValuesPanel
            onChangeKeyValue={kvs => {
              this.props.actions.helm.inputKeyValue(kvs);
            }}
            kvs={kvs}
          />
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={confirm}>
            {t('确认')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
