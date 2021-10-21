/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import * as React from 'react';

import { Button, Modal } from '@tea/component';
import { FetchState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, TipInfo } from '../../../../common/components/';
import { RootProps } from '../../HelmApp';
import { KeyValuesPanel } from '../create/KeyValuesPanel';

interface Props extends RootProps {
  onCancel?: Function;
}

export class UpdateHelmDialog extends React.Component<Props, {}> {
  componentDidMount() {
    this.props.actions.helm.inputKeyValue([]);
  }
  componentWillUnmount() {
    // this.props.actions.create.clear();
    // this.props.actions.create.inputKeyValue([]);
    this.props.actions.helm.inputKeyValue([]);
  }
  render() {
    const { actions } = this.props;
    const select = (version: string) => {
      let versionSelect = tencenthubChartVersionList.data.records.find(item => item.version === version);
      actions.helm.selectTencenthubChartVersion(versionSelect);
    };
    const cancel = () => {
      this.props.onCancel && this.props.onCancel();
    };
    const confirm = () => {
      actions.helm.updateHelm();
      cancel();
    };

    const {
      listState: { helmSelection, tencenthubChartVersionList, tencenthubChartVersionSelection, kvs }
    } = this.props;

    let versionOptions = [<option key={-1}>{t('正在加载...')}</option>];

    if (tencenthubChartVersionList.fetchState === FetchState.Ready && tencenthubChartVersionList.fetched) {
      versionOptions = tencenthubChartVersionList.data.records.map((item, index) => (
        <option key={index} value={item.version}>
          {item.version}
        </option>
      ));
    }
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
            <FormItem label={t('版本')}>
              <div className="form-unit">
                <select
                  className="tc-15-select m"
                  style={{ minWidth: '150px' }}
                  value={tencenthubChartVersionSelection ? tencenthubChartVersionSelection.version : ''}
                  onChange={e => {
                    select(e.target.value);
                  }}
                >
                  {versionOptions}
                </select>
              </div>
            </FormItem>
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
