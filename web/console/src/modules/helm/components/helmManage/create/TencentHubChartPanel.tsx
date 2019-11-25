import * as React from 'react';
import { RootProps } from '../../HelmApp';
import classNames from 'classnames';
import { CommonBar, FormItem, Markdown } from '../../../../common/components';
import { HelmResource, tencentHubTypeList, TencentHubType } from '../../../constants/Config';
import { TencenthubChart, TencenthubChartVersion } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
interface HelmCreataeState {
  resource?: string;
  helmResourceSelection?: string;
  otherResourceConfig?: string;
  isValid?: boolean;
}

export class TencentHubChartPanel extends React.Component<RootProps, HelmCreataeState> {
  state = {
    resource: HelmResource.TencentHub,
    helmResourceSelection: null,
    otherResourceConfig: null,
    isValid: true
  };

  onSelectTencenthubType(type: string) {
    this.props.actions.create.selectTencenthubType(type);
  }

  onSelectTencenthubNamespace(namespace: string) {
    this.props.actions.create.selectTencenthubNamespace(namespace);
  }

  onSelectTencenthubChart(chart: TencenthubChart) {
    this.props.actions.create.selectTencenthubChart(chart);
  }

  onSelectTencenthubChartVersion(version: TencenthubChartVersion) {
    this.props.actions.create.selectTencenthubChartVersion(version);
  }

  renderTypeList() {
    const { actions } = this.props;
    const {
      helmCreation: { tencenthubTypeSelection }
    } = this.props;
    return (
      <FormItem label={t('类型')}>
        <div className="form-unit">
          <div className="tc-15-rich-radio">
            <CommonBar
              list={tencentHubTypeList}
              value={tencenthubTypeSelection}
              onSelect={item => {
                this.onSelectTencenthubType(item.value as string);
              }}
            />
          </div>
        </div>
      </FormItem>
    );
  }

  renderNamespaceList() {
    const { tencenthubNamespaceList, tencenthubNamespaceSelection } = this.props.helmCreation;
    let list = [];
    tencenthubNamespaceList.data.records.forEach((item, index) => {
      list.push({
        value: item.name,
        name: item.name
      });
    });
    return (
      <FormItem label={t('组织')}>
        <div className="form-unit">
          <div className="tc-15-rich-radio">
            <select
              className="tc-15-select m"
              value={tencenthubNamespaceSelection}
              onChange={e => this.onSelectTencenthubNamespace(e.target.value as string)}
            >
              {tencenthubNamespaceList.data.records.map((item, index) => {
                return (
                  <option key={index} value={item.name as string}>
                    {item.name}
                  </option>
                );
              })}
            </select>
          </div>
        </div>
      </FormItem>
    );
  }

  renderChartList() {
    let {
      helmCreation: { tencenthubChartList, tencenthubChartSelection }
    } = this.props;
    return tencenthubChartList.data.records
      .filter(item => item.download_url)
      .map((item, index) => (
        <li key={index}>
          <i
            className={classNames(
              tencenthubChartSelection && tencenthubChartSelection.name === item.name
                ? 'icon-arrow-down'
                : 'icon-arrow-right',
              'select-icon'
            )}
            onClick={e => this.onSelectTencenthubChart(item)}
          />
          <span className="opt-txt" style={{ cursor: 'pointer' }} onClick={e => this.onSelectTencenthubChart(item)}>
            <span className="opt-txt-inner">
              <span className="item-name" title={item.name}>
                {item.name}
              </span>
            </span>
          </span>

          {tencenthubChartSelection && tencenthubChartSelection.name === item.name && (
            <div className="haschild active" style={{ marginBottom: 10 }}>
              {this.renderVersionList()}
            </div>
          )}
        </li>
      ));
  }
  renderVersionList() {
    let {
      helmCreation: { tencenthubChartVersionList, tencenthubChartVersionSelection }
    } = this.props;
    return tencenthubChartVersionList.data.records.map((item, index) => (
      <p key={index}>
        <label className="form-ctrl-label" title={item.version}>
          <input
            type="radio"
            name="rd-demo1"
            className="tc-15-radio"
            checked={tencenthubChartVersionSelection && tencenthubChartVersionSelection.version === item.version}
            onClick={e => this.onSelectTencenthubChartVersion(item)}
          />
          {item.version}
        </label>
      </p>
    ));
  }

  renderChartPanel() {
    const {
      helmCreation: { tencenthubChartReadMe }
    } = this.props;
    let maxWidth = 600 * 0.8;

    let ele = document.querySelector('.tc-panel');
    if (ele) {
      maxWidth = ele['offsetWidth'] - 40 - 80 - 225;
    }
    return (
      <FormItem label="Chart">
        <div className="form-unit">
          <div className="wrap-mod-box">
            <div className="helm-app-box clearfix">
              <div className="helm-app-left-box" style={{ marginTop: 0 }}>
                <div className="tc-15-mod-selector">
                  <div className="tc-15-mod-selector-tb">
                    <div className="tc-15-option-cell options-left">
                      <div className="tc-15-option-bd" style={{ height: 725 }}>
                        <div className="tc-15-option-box tc-scroll" style={{ height: 725 }}>
                          <ul className="tc-15-option-list">{this.renderChartList()}</ul>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="helm-app-right-box">
                <div className="helm-app-rightcont" style={{ height: 725, width: maxWidth }}>
                  <Markdown
                    style={{
                      minHeight: 685
                    }}
                    className="lab-box-cont"
                    text={tencenthubChartReadMe ? tencenthubChartReadMe.content : ''}
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </FormItem>
    );
  }

  render() {
    let {
      helmCreation: { tencenthubTypeSelection }
    } = this.props;
    let { resource } = this.state;

    return (
      <ul className="form-list" style={{ marginBottom: 20 }}>
        {resource === HelmResource.TencentHub && this.renderTypeList()}
        {resource === HelmResource.TencentHub &&
          tencenthubTypeSelection === TencentHubType.Private &&
          this.renderNamespaceList()}
        {resource === HelmResource.TencentHub && this.renderChartPanel()}
      </ul>
    );
  }
}
