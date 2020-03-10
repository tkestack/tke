import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ExternalLink } from '@tencent/tea-component';

import { FormItem } from '../../../../common/components';
import { allActions } from '../../../actions';
import { HealthCheckMethodList, HttpProtocolTypeList } from '../../../constants/Config';
import { HealthCheckItem } from '../../../models';
import { RootProps } from '../../ClusterApp';

interface ContainerHealthCheckPanelProps extends RootProps {
  cKey: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerHealthCheckPanel extends React.Component<ContainerHealthCheckPanelProps, {}> {
  render() {
    let { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    let container = containers.find(c => c.id === cKey),
      healthCheck = container.healthCheck;

    return (
      container && (
        <FormItem
          label={t('容器健康检查')}
          tips={t('健康检查可以帮助你探测容器是否正常，以保证服务的正常运作')}
          isPureText={true}
        >
          <div className="form-unit">
            <p>
              <label className="form-ctrl-label">
                <input
                  type="checkbox"
                  className="tc-15-checkbox"
                  checked={healthCheck.isOpenLiveCheck ? true : false}
                  onChange={e =>
                    actions.editWorkload.updateHealthCheck({ isOpenLiveCheck: !healthCheck.isOpenLiveCheck }, cKey, '')
                  }
                />
                {t('存活检查')}
                <span className="inline-help-text text-label">{t('检查容器是否正常，不正常则重启实例')}</span>
              </label>
            </p>
            {healthCheck.isOpenLiveCheck && <EditHealthCheckItem cKey={cKey} hType="liveCheck" />}
            <p>
              <label className="form-ctrl-label">
                <input
                  type="checkbox"
                  className="tc-15-checkbox"
                  checked={healthCheck.isOpenReadyCheck ? true : false}
                  onChange={e =>
                    actions.editWorkload.updateHealthCheck(
                      { isOpenReadyCheck: !healthCheck.isOpenReadyCheck },
                      cKey,
                      ''
                    )
                  }
                />
                {t('就绪检查')}
                <span className="inline-help-text text-label">
                  {t('检查容器是否就绪，不就绪则停止转发流量到当前实例')}
                </span>
              </label>
            </p>
            {healthCheck.isOpenReadyCheck && <EditHealthCheckItem cKey={cKey} hType="readyCheck" />}
            <p className="text-label">
              <span style={{ verticalAlign: 'middle' }}>{t('查看健康检查和就绪检查')}</span>
            </p>
          </div>
        </FormItem>
      )
    );
  }
}

interface EditHealthCheckItemProps extends RootProps {
  cKey: string;
  hType?: string;
}

@connect(state => state, mapDispatchToProps)
class EditHealthCheckItem extends React.Component<EditHealthCheckItemProps, {}> {
  render() {
    let { actions, cKey, hType, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    let container = containers.find(c => c.id === cKey),
      healthCheck: HealthCheckItem = container.healthCheck[hType];

    // 渲染检查方法
    let methodOptions = HealthCheckMethodList.map((m, index) => (
      <option key={index} value={m.value}>
        {m.label}
      </option>
    ));

    // 渲染协议列表
    let protocolOptions = HttpProtocolTypeList.map((h, index) => (
      <option key={index} value={h.value}>
        {h.label}
      </option>
    ));

    // 健康阈值的提示
    let numTip = (
      <p>
        {t('表示后端容器从失败到成功的连续健康检查成功次数，范围：')}
        {hType === 'liveCheck' ? t('只能为1') : t('1~10次')}
      </p>
    );

    return (
      <div className="param-box three-level-param">
        <div className="param-bd">
          <ul className="form-list fixed-layout">
            <FormItem label={t('检查方法')}>
              <div className="form-unit">
                <select
                  className="tc-15-select m"
                  value={healthCheck.checkMethod}
                  onChange={e => this._handleCheckMethodSelect(e.target.value, cKey, hType)}
                >
                  {methodOptions}
                </select>
              </div>
            </FormItem>
            <FormItem label={t('检查协议')} isShow={healthCheck.checkMethod === 'methodHttp'}>
              <div className="form-unit">
                <select
                  className="tc-15-select m"
                  onChange={e => this._handleProtocolSelect(e.target.value, cKey, hType)}
                >
                  {protocolOptions}
                </select>
              </div>
            </FormItem>
            <FormItem
              label={t('检查端口')}
              isShow={healthCheck.checkMethod === 'methodTcp' || healthCheck.checkMethod === 'methodHttp'}
            >
              <div className={healthCheck.v_port.status === 2 ? 'is-error' : ''}>
                <input
                  type="text"
                  className="tc-15-input-text m"
                  style={{ width: '113px' }}
                  value={healthCheck.port}
                  onChange={e => actions.editWorkload.updateHealthCheck({ port: e.target.value }, cKey, hType)}
                  onBlur={e => actions.validate.workload.validateHealthPort(cKey + '', hType)}
                />
                <p className="inline-help-text text-label">{t('端口范围：')}1~65535</p>
                {healthCheck.v_port.status === 2 && <p className="form-input-help">{healthCheck.v_port.message}</p>}
              </div>
            </FormItem>
            <FormItem label={t('请求路径')} isShow={healthCheck.checkMethod === 'methodHttp'}>
              <div className="form-unit">
                <input
                  type="text"
                  className="tc-15-input-text m"
                  style={{ width: '240px' }}
                  value={healthCheck.path}
                  onChange={e => actions.editWorkload.updateHealthCheck({ path: e.target.value }, cKey, hType)}
                />
              </div>
            </FormItem>
            <FormItem label={t('执行命令')} isShow={healthCheck.checkMethod === 'methodCmd'}>
              <div className={healthCheck.v_cmd.status === 2 ? 'is-error' : ''}>
                <textarea
                  className="tc-15-input-text m"
                  style={{ width: '240px' }}
                  value={healthCheck.cmd}
                  onChange={e => actions.editWorkload.updateHealthCheck({ cmd: e.target.value }, cKey, hType)}
                  onBlur={e => actions.validate.workload.validateHealthCmd(cKey, hType)}
                />
                {healthCheck.v_cmd.status === 2 && <p className="form-input-help">{healthCheck.v_cmd.message}</p>}
              </div>
            </FormItem>
            <FormItem label={t('启动延时')} tips={t('容器延时启动健康检查的时间，范围：0~60秒')}>
              <div className="form-unit">
                <div className={healthCheck.v_delayTime.status === 2 ? 'is-error' : ''}>
                  <input
                    type="text"
                    className="tc-15-input-text m"
                    style={{ width: '113px' }}
                    value={healthCheck.delayTime + ''}
                    onChange={e => actions.editWorkload.updateHealthCheck({ delayTime: e.target.value }, cKey, hType)}
                    onBlur={e => actions.validate.workload.validateHealthDelayTime(cKey, hType)}
                  />
                  <Trans>
                    <span className="inline-help-text">秒</span>
                    <p className="inline-help-text text-label">范围：0~60秒</p>
                  </Trans>
                  {healthCheck.v_delayTime.status === 2 && (
                    <p className="form-input-help">{healthCheck.v_delayTime.message}</p>
                  )}
                </div>
              </div>
            </FormItem>
            <FormItem label={t('响应超时')} tips={t('每次健康检查响应的最大超时时间，范围：2~60秒')}>
              <div className={healthCheck.v_timeOut.status === 2 ? 'is-error' : ''}>
                <input
                  type="text"
                  className="tc-15-input-text m"
                  style={{ width: '113px' }}
                  value={healthCheck.timeOut + ''}
                  onChange={e => actions.editWorkload.updateHealthCheck({ timeOut: e.target.value }, cKey, hType)}
                  onBlur={e => actions.validate.workload.validateHealthTimeOut(cKey, hType)}
                />
                <Trans>
                  <span className="inline-help-text">秒</span>
                  <p className="inline-help-text text-label">范围：2~60秒</p>
                </Trans>
                {healthCheck.v_timeOut.status === 2 && (
                  <p className="form-input-help">{healthCheck.v_timeOut.message}</p>
                )}
              </div>
            </FormItem>
            <FormItem label={t('间隔时间')} tips={t('进行健康检查的时间间隔，范围：大于响应超时，小于300秒')}>
              <div className="form-unit">
                <div className={healthCheck.v_intervalTime.status === 2 ? 'is-error' : ''}>
                  <input
                    type="text"
                    className="tc-15-input-text m"
                    style={{ width: '113px' }}
                    value={healthCheck.intervalTime + ''}
                    onChange={e =>
                      actions.editWorkload.updateHealthCheck({ intervalTime: e.target.value }, cKey, hType)
                    }
                    onBlur={e => actions.validate.workload.validateHealthIntervalTime(cKey, hType)}
                  />
                  <Trans>
                    <span className="inline-help-text">秒</span>
                    <p className="inline-help-text text-label">范围：2~300秒</p>
                  </Trans>
                  {healthCheck.v_intervalTime.status === 2 && (
                    <p className="form-input-help">{healthCheck.v_intervalTime.message}</p>
                  )}
                </div>
              </div>
            </FormItem>
            <FormItem label={t('健康阀值')} tips={numTip}>
              <div className="form-unit">
                <div className={healthCheck.v_healthThreshold.status === 2 ? 'is-error' : ''}>
                  <input
                    type="text"
                    className="tc-15-input-text m"
                    style={{ width: '113px' }}
                    value={healthCheck.healthThreshold + ''}
                    disabled={hType === 'liveCheck'}
                    onChange={e =>
                      actions.editWorkload.updateHealthCheck({ healthThreshold: e.target.value }, cKey, hType)
                    }
                    onBlur={e => actions.validate.workload.validateHealthThreshold(cKey, hType)}
                  />

                  {hType === 'liveCheck' ? (
                    <Trans>
                      <span className="inline-help-text">次</span>
                      <p className="inline-help-text text-label">范围：1次</p>
                    </Trans>
                  ) : (
                    <Trans>
                      <span className="inline-help-text">次</span>
                      <p className="inline-help-text text-label">范围：1~10次</p>
                    </Trans>
                  )}
                  {healthCheck.v_healthThreshold.status === 2 && (
                    <p className="form-input-help">{healthCheck.v_healthThreshold.message}</p>
                  )}
                </div>
              </div>
            </FormItem>
            <FormItem label={t('不健康阀值')} tips={t('表示后端容器从成功到失败的连续健康检查成功次数，范围：1~10次')}>
              <div className="form-unit">
                <div className={healthCheck.v_unhealthThreshold.status === 2 ? 'is-error' : ''}>
                  <input
                    type="text"
                    className="tc-15-input-text m"
                    style={{ width: '113px' }}
                    value={healthCheck.unhealthThreshold + ''}
                    onChange={e =>
                      actions.editWorkload.updateHealthCheck({ unhealthThreshold: e.target.value }, cKey, hType)
                    }
                    onBlur={e => actions.validate.workload.validateUnHealthThreshold(cKey, hType)}
                  />
                  <Trans>
                    <span className="inline-help-text">次</span>
                    <p className="inline-help-text text-label">范围：1~10次</p>
                  </Trans>
                  {healthCheck.v_unhealthThreshold.status === 2 && (
                    <p className="form-input-help">{healthCheck.v_unhealthThreshold.message}</p>
                  )}
                </div>
              </div>
            </FormItem>
          </ul>
        </div>
      </div>
    );
  }

  /** 处理健康检查方法的操作 */
  private _handleCheckMethodSelect(method: string, cKey: string, hType: string) {
    let { actions } = this.props;
    let obj = {
      checkMethod: method
    };

    if (method === 'methodHttp') {
      obj = Object.assign({}, obj, {
        port: 80,
        protocol: 'HTTP'
      });
    }
    actions.editWorkload.updateHealthCheck(obj, cKey, hType);
  }

  /** 处理协议相关的选择 */
  private _handleProtocolSelect(protocol: string, cKey: string, hType: string) {
    let { actions } = this.props;

    let obj = {
      protocol,
      port: protocol === 'HTTP' ? 80 : 443
    };

    actions.editWorkload.updateHealthCheck(obj, cKey, hType);
  }
}
