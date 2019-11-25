import * as React from 'react';
import classNames from 'classnames';
import { FormItem, LinkButton } from '../../../../common/components';
import { HelmKeyValue } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
interface Props {
  onChangeKeyValue?: Function;
  kvs?: HelmKeyValue[];
}
export class KeyValuesPanel extends React.Component<Props, {}> {
  onChangeKeyValue(kvs: HelmKeyValue[]) {
    this.props.onChangeKeyValue && this.props.onChangeKeyValue(kvs.slice(0));
  }
  onChangeKey(key: string, index: number) {
    let kvs = this.props.kvs;
    kvs[index].key = key;
    this.onChangeKeyValue(kvs);
  }
  onChangeValue(value: string, index: number) {
    let kvs = this.props.kvs;
    kvs[index].value = value;
    this.onChangeKeyValue(kvs);
  }
  onAdd() {
    let kvs = this.props.kvs;
    kvs.push({ key: '', value: '' });
    this.onChangeKeyValue(kvs);
  }
  onDelete(index: number) {
    let kvs = this.props.kvs;
    kvs.splice(index, 1);
    this.onChangeKeyValue(kvs);
  }

  isValid() {
    let kvs = this.props.kvs;
    let flag = true;
    kvs.forEach(item => {
      if (item.key.trim() === '' || item.value.trim() === '') {
        flag = false;
      }
    });
    return flag;
  }

  render() {
    let kvs = this.props.kvs;
    let isValid = this.isValid();
    return (
      <div>
        <ul className="form-list" style={{ marginBottom: 16 }}>
          {kvs.map((kv, index) => {
            return (
              <FormItem key={index} label={index === 0 ? 'Key-Value' : ''}>
                <div className={classNames('form-unit')}>
                  <input
                    type="text"
                    className="tc-15-input-text m"
                    placeholder={t('变量名')}
                    value={kv.key}
                    onChange={e => this.onChangeKey((e.target.value + '').trim(), index)}
                  />
                  <span className="inline-help-text">=</span>
                  <textarea
                    className="tc-15-input-text m"
                    placeholder={t('变量值')}
                    style={{ marginLeft: 10 }}
                    value={kv.value}
                    onChange={e => this.onChangeValue((e.target.value + '').trim(), index)}
                  />
                  <div className="tc-15-bubble-icon">
                    <a
                      href="javascript:;"
                      style={{ fontSize: 12, marginLeft: 10, marginRight: 10, verticalAlign: 'middle' }}
                      onClick={e => {
                        this.onDelete(index);
                      }}
                    >
                      <i className="icon-cancel-icon" />
                    </a>
                    <div className="tc-15-bubble tc-15-bubble-bottom">
                      <div className="tc-15-bubble-inner">{t('删除')}</div>
                    </div>
                  </div>
                  {index === kvs.length - 1 && (
                    <p className="form-input-help">
                      {t('可通过设置自定义参数替换Chart包的默认配置，如：image.repository = nginx')}
                    </p>
                  )}
                </div>
              </FormItem>
            );
          })}
          <FormItem key={'add'} label={kvs.length === 0 ? 'Key-Value' : ''}>
            <div className={classNames('form-unit')}>
              <LinkButton
                disabled={!isValid}
                errorTip={t('请先完成待编辑项')}
                onClick={e => {
                  this.onAdd();
                }}
              >
                {t('新增变量')}
              </LinkButton>
            </div>
          </FormItem>
        </ul>
      </div>
    );
  }
}
