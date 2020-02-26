import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { allActions } from '../../../actions';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { InputNumber, Bubble, ExternalLink } from '@tea/component';
import { connect } from 'react-redux';
import * as classnames from 'classnames';
import { FormItem, LinkButton } from '../../../../common/components';
import { ContainerMaxNumLimit, HpaMetricsTypeList } from '../../../constants/Config';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const metricUnitMap = {
  cpuUtilization: '%',
  memoryUtilization: '%',
  cpuAverage: t('核'),
  memoryAverage: 'Mib',
  inBandwidth: 'Mbps',
  outBandwidth: 'Mbps'
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerNumPanel extends React.Component<RootProps, {}> {
  render() {
    let { subRoot } = this.props,
      { mode, workloadEdit, resourceOption } = subRoot,
      { resourceSelection } = resourceOption,
      { isNeedContainerNum, workloadType, hpaList, volumes } = workloadEdit;

    let isHasCbs = false;
    let isCanUseHpa = workloadType === 'deployment' || workloadType === 'statefulset' || workloadType === 'tapp';

    return (
      <FormItem label={t('实例数量')} tips={t('设置服务自动伸缩要求集群版本大于1.7.0')} isShow={isNeedContainerNum}>
        {hpaList.fetchState === FetchState.Fetching && mode === 'update' ? (
          <div>
            <i className="n-loading-icon" />
            &nbsp; <span className="text">{t('加载中...')}</span>
          </div>
        ) : (
          <div className="up-date">
            {this._manualUpdateContainerNum(isHasCbs)}
            {isCanUseHpa && this._autoscaleContainerNum(isHasCbs)}
            {isCanUseHpa && this._cronhpaContainerNum()}
          </div>
        )}
      </FormItem>
    );
  }

  /** 展示 cronhpa的相关逻辑 */
  private _cronhpaContainerNum() {
    let { actions, subRoot } = this.props,
      { workloadEdit, addons } = subRoot,
      { scaleType } = workloadEdit;

    // 是否为定时调节的模式
    let isCron = scaleType === 'crontab';

    let disabled = addons['CronHPA'] === undefined;

    return (
      <div className="as-sel-box">
        <Bubble placement="top" content={disabled ? t('当前集群尚未安装CronHPA扩展组件，请联系管理员进行安装') : null}>
          <label className="form-ctrl-label">
            <input
              type="radio"
              name="as-service-redios"
              className="tc-15-radio"
              value="crontab"
              checked={isCron}
              disabled={disabled}
              onChange={e => {
                actions.editWorkload.updateScaleType(e.target.value);
              }}
            />
            <strong>{t('定时调节')}</strong>
            <span className="text-label">
              {t('根据设置的Crontab（Crontab语法格式，例如 "0 23 * * 5"表示每周五23:00）周期性地设置实例数量')}
            </span>
          </label>
        </Bubble>

        {isCron ? <ul className="form-list">{this._renderCronTabList()}</ul> : <noscript />}
      </div>
    );
  }

  /** 展示cronhpa的触发策略 */
  private _renderCronTabList() {
    let { actions, subRoot } = this.props,
      { cronMetrics } = subRoot.workloadEdit;

    // 是否可以删除该触发策略
    let canDelete = cronMetrics.length > 1 ? true : false;

    return (
      <FormItem label={t('触发策略')}>
        <div className="form-unit is-success">
          {cronMetrics.map((metric, index) => {
            let mId = metric.id + '';

            return (
              <div className="code-list" key={index}>
                <div
                  className={classnames('form-unit', {
                    'is-error': metric.v_crontab.status === 2 || metric.v_targetReplicas.status === 2
                  })}
                  style={{ display: 'inline-block', marginBottom: '5px' }}
                >
                  <Bubble placement="bottom" content={metric.v_crontab.status === 2 ? metric.v_crontab.message : null}>
                    <input
                      type="text"
                      placeholder="Crontab"
                      className="tc-15-input-text m mr10"
                      style={{ maxWidth: '120px' }}
                      value={metric.crontab}
                      onChange={e => actions.editWorkload.cronhpa.updateMetric({ crontab: e.target.value }, mId)}
                      onBlur={e => actions.validate.workload.validateCronTab(mId)}
                    />
                  </Bubble>

                  <Bubble
                    placement="bottom"
                    content={metric.v_targetReplicas.status === 2 ? metric.v_targetReplicas.message : null}
                  >
                    <input
                      type="text"
                      placeholder="目标实例数"
                      className="tc-15-input-text m"
                      style={{ maxWidth: '120px' }}
                      value={metric.targetReplicas}
                      onChange={e => actions.editWorkload.cronhpa.updateMetric({ targetReplicas: e.target.value }, mId)}
                      onBlur={e => actions.validate.workload.validateCronTargetReplicas(mId)}
                    />
                  </Bubble>

                  <span className="inline-help-text">
                    <LinkButton
                      disabled={!canDelete}
                      errorTip={t('至少设置一个策略')}
                      onClick={() => actions.editWorkload.cronhpa.deleteMetric(mId)}
                    >
                      <i className="icon-cancel-icon" />
                    </LinkButton>
                  </span>
                </div>
              </div>
            );
          })}

          <div>
            <LinkButton
              onClick={() => {
                actions.editWorkload.cronhpa.addMetric();
              }}
            >
              {t('新增策略')}
            </LinkButton>
          </div>
        </div>
      </FormItem>
    );
  }

  /** 展示hpa的相关逻辑 */
  private _autoscaleContainerNum(isHasCbs: boolean) {
    let { actions, subRoot, cluster, clusterVersion } = this.props,
      { scaleType, minReplicas, maxReplicas, v_minReplicas, v_maxReplicas } = subRoot.workloadEdit;

    let isAutoScale = scaleType === 'autoScale';

    // 判断是否可以使用hpa
    let canUseHpa = false;
    if (cluster.selection) {
      // k8s版本只有 >= 1.7.0 才能使用hpa
      let k8sVersion = clusterVersion.split('.');
      canUseHpa = Number(k8sVersion[0]) >= 1 && Number(k8sVersion[1]) >= 7 ? true : false;
    }

    // 当前容器实例的最大数量的控制
    let canAddNum = true,
      maxLimit = ContainerMaxNumLimit;
    // 如果容器的数据卷挂载包含云盘，则实例数量只能为一个
    if (isHasCbs) {
      maxLimit = 1;
      canAddNum = false;
    }

    return (
      <div className="as-sel-box">
        <Bubble placement="bottom" content={this._renderBubbleText(canUseHpa, canAddNum) || null}>
          <label className="form-ctrl-label">
            <input
              type="radio"
              name="as-service-redios"
              className="tc-15-radio"
              value="autoScale"
              disabled={!canUseHpa || !canAddNum}
              checked={isAutoScale}
              onChange={e => {
                actions.editWorkload.updateScaleType(e.target.value);
              }}
            />
            <strong>{t('自动调节')}</strong>
            <span className="text-label" style={{ verticalAlign: '-1px' }}>
              {t('满足任一设定条件，则自动调节实例（pod）数目')}
            </span>
          </label>
        </Bubble>
        {isAutoScale ? (
          <ul className="form-list">
            {this._renderAutoScaleList()}

            <FormItem label={t('实例范围')}>
              <div
                className={classnames('form-unit', { 'is-error': v_minReplicas.status === 2 })}
                style={{ display: 'inline-block', verticalAlign: 'middle' }}
              >
                <Bubble placement="bottom" content={v_minReplicas.status === 2 ? v_minReplicas.message : null}>
                  <input
                    type="text"
                    placeholder={t('最小实例数')}
                    className="tc-15-input-text m"
                    style={{ maxWidth: '120px' }}
                    value={minReplicas}
                    onChange={e => actions.editWorkload.inputMinReplicas(e.target.value)}
                    onBlur={e => actions.validate.workload.validateMinReplicas()}
                  />
                </Bubble>
              </div>
              <span className="inline-help-text" style={{ margin: '0 5px' }}>
                ~
              </span>
              <div
                className={classnames('form-unit', { 'is-error': v_maxReplicas.status === 2 })}
                style={{ display: 'inline-block', verticalAlign: 'middle' }}
              >
                <Bubble placement="bottom" content={v_maxReplicas.status === 2 ? v_maxReplicas.message : null}>
                  <input
                    type="text"
                    placeholder={t('最大实例数')}
                    className="tc-15-input-text m"
                    style={{ maxWidth: '120px' }}
                    value={maxReplicas}
                    onChange={e => actions.editWorkload.inputMaxReplicas(e.target.value)}
                    onBlur={e => actions.validate.workload.validateMaxReplicas()}
                  />
                </Bubble>
              </div>
              <p className="text-label">{t('在设定的实例范围内自动调节，不会超出该设定范围')}</p>
            </FormItem>
          </ul>
        ) : (
          <noscript />
        )}
      </div>
    );
  }

  /** 展示hpa的相关提示 */
  private _renderBubbleText(canUseHpa: boolean = false, canAddNum: boolean) {
    let bubbleText = '';
    if (!canUseHpa) {
      bubbleText = t('设置服务自动伸缩要求集群版本大于1.7.0');
    } else if (!canAddNum) {
      bubbleText = t('设置服务自动伸缩要求服务不能挂载云硬盘');
    }

    return bubbleText;
  }

  /** 渲染自动伸缩规则 */
  private _renderAutoScaleList() {
    let { actions, subRoot } = this.props,
      { metrics } = subRoot.workloadEdit;

    // 是否可以删除该项触发策略
    let canDelete = metrics.length > 1 ? true : false;
    // 是否可以新增触发策略
    let canAdd = metrics.length < 4 ? true : false;

    /** 渲染指标的列表 */
    let metricTypeOptions = HpaMetricsTypeList.map((item, index) => {
      return (
        <option key={index} value={item.value}>
          {item.label}
        </option>
      );
    });

    return (
      <FormItem label={t('触发策略')}>
        <div className="form-unit is-success">
          {metrics.map((metric, index) => {
            let mId = metric.id + '';

            return (
              <div className="code-list" key={index}>
                <div
                  className={classnames('form-unit', {
                    'is-error': metric.v_type.status === 2 || metric.v_value.status === 2
                  })}
                  style={{ display: 'inline-block', marginBottom: '5px' }}
                >
                  <Bubble placement="bottom" content={metric.v_type.status === 2 ? metric.v_type.message : null}>
                    <select
                      className="tc-15-select m"
                      style={{ marginRight: '6px' }}
                      value={metric.type}
                      onChange={e => {
                        actions.editWorkload.updateMetric({ type: e.target.value }, mId);
                        actions.validate.workload.validateHpaType(e.target.value, mId);
                      }}
                    >
                      {metricTypeOptions}
                    </select>
                  </Bubble>
                  <Bubble placement="bottom" content={metric.v_value.status === 2 ? metric.v_value.message : null}>
                    <input
                      type="text"
                      placeholder={t('目标阈值')}
                      className="tc-15-input-text m mr10"
                      style={{ maxWidth: '120px' }}
                      value={metric.value}
                      onChange={e => actions.editWorkload.updateMetric({ value: e.target.value }, mId)}
                      onBlur={e => actions.validate.workload.validateHpaValue(e.target.value, mId)}
                    />
                  </Bubble>
                  <span className="text" style={{ marginLeft: '0', verticalAlign: '-4px' }}>
                    {metricUnitMap[metric.type]}
                  </span>
                  <span className="inline-help-text">
                    <LinkButton disabled={!canDelete} onClick={() => actions.editWorkload.deleteMetric(mId)}>
                      <i className="icon-cancel-icon" />
                    </LinkButton>
                  </span>
                </div>
              </div>
            );
          })}
          <div>
            <LinkButton
              disabled={!canAdd}
              errorTip={t('最大指定四项触发策略')}
              onClick={() => {
                actions.editWorkload.addMetric();
              }}
            >
              {t('新增策略')}
            </LinkButton>
          </div>
        </div>
      </FormItem>
    );
  }

  /** 展示 手动调节的逻辑 */
  private _manualUpdateContainerNum(isHasCbs: boolean) {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { volumes, containerNum, scaleType } = workloadEdit;

    // 当前的实例数量的模式
    let isManual = scaleType === 'manualScale';

    // 当前容器实例的最大数量的控制
    let canAddNum = true,
      maxLimit = ContainerMaxNumLimit;
    // 如果容器的数据卷挂载包含云盘，则实例数量只能为一个
    if (isHasCbs) {
      maxLimit = 1;
      canAddNum = false;
    }

    return (
      <div className="as-sel-box">
        <label className="form-ctrl-label">
          <input
            type="radio"
            name="as-service-redios"
            className="tc-15-radio"
            value="manualScale"
            checked={isManual}
            onChange={e => {
              actions.editWorkload.updateScaleType(e.target.value);
            }}
          />
          <strong>{t('手动调节')}</strong>
          <span className="text-label">{t('直接设定实例数量')}</span>
        </label>
        {isManual ? (
          <ul className="form-list">
            <FormItem label={t('实例数量')}>
              <div className="form-unit">
                <InputNumber
                  value={+containerNum}
                  // size={'m'}
                  min={0}
                  max={maxLimit}
                  step={1}
                  unit={t('个')}
                  onChange={value => actions.editWorkload.updateContainerNum(value + '')}
                />
                {!canAddNum && (
                  <span className="inline-help-text text-danger">
                    {t('注意：当前已设置CBS数据卷，实例数量限制为1')}
                  </span>
                )}
              </div>
            </FormItem>
          </ul>
        ) : (
          <noscript />
        )}
      </div>
    );
  }
}
