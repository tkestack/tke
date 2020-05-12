import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Select, SelectMultiple, Switch } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem } from '../../../../common/components';
import { allActions } from '../../../actions';
import { ImagePullPolicyList, AddCapabilitiesList, DropCapabilitiesList } from '../../../constants/Config';
import { RootProps } from '../../ClusterApp';
import { EditResourceContainerHealthCheckPanel } from './EditResourceContainerHealthCheckPanel';

interface ContainerAdvancedProps extends RootProps {
  cKey: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerAdvancedPanel extends React.Component<ContainerAdvancedProps, {}> {
  render() {
    let { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    let container = containers.find(c => c.id === cKey);

    return (
      <div className="param-bd">
        <ul className="form-list fixed-layout">
          <FormItem label={t('工作目录')}>
            <div className={classnames('form-unit', { 'is-error': container.v_workingDir.status === 2 })}>
              <input
                type="text"
                placeholder=""
                className="tc-15-input-text m"
                value={container.workingDir}
                onChange={e => actions.editWorkload.updateContainer({ workingDir: e.target.value }, cKey)}
                onBlur={e => actions.validate.workload.validateWorkingDir(e.target.value, cKey)}
              />
              <p className="form-input-help text-label">
                <Trans>
                  <span style={{ verticalAlign: '-1px' }}>指定容器运行后的工作目录，</span>
                </Trans>
              </p>
              {container.v_workingDir.status === 2 && (
                <p className="form-input-help">{container.v_workingDir.message}</p>
              )}
            </div>
          </FormItem>
          <FormItem label={t('日志目录')}>
            <Select
              size="auto"
              options={container.mounts.map(({ mountPath }) => ({ value: mountPath, text: mountPath }))}
              value={container.logDir}
              onChange={value => {
                actions.editWorkload.updateContainer({ logDir: value }, cKey);
              }}
            />
            <input
              type="text"
              placeholder=""
              className="tc-15-input-text m"
              value={container.logPath}
              onChange={e => actions.editWorkload.updateContainer({ logPath: e.target.value }, cKey)}
              onBlur={e => actions.validate.workload.validateWorkingDir(e.target.value, cKey)}
            />
            <p className="form-input-help text-label">
              <Trans>
                <span style={{ verticalAlign: '-1px' }}>指定容器运行后的日志目录，</span>
              </Trans>
            </p>
          </FormItem>
          <FormItem label={t('运行命令')}>
            <div className="form-unit is-success">
              <textarea
                className="tc-15-input-textarea"
                value={container.cmd}
                onChange={e => actions.editWorkload.updateContainer({ cmd: e.target.value }, cKey)}
                placeholder={t('注意每个命令单独一行')}
              />
              <p className="form-input-help text-label">
                <Trans>
                  <span style={{ verticalAlign: '-1px' }}>控制容器运行的输入命令</span>
                </Trans>
              </p>
            </div>
          </FormItem>
          <FormItem label={t('运行参数')}>
            <div className="form-unit is-success">
              <textarea
                className="tc-15-input-textarea"
                value={container.arg}
                onChange={e => actions.editWorkload.updateContainer({ arg: e.target.value }, cKey)}
                placeholder={t('注意每个参数单独一行')}
              />
              <p className="form-input-help text-label">
                <Trans>
                  <span style={{ verticalAlign: '-1px' }}>传递给容器运行命令的输入参数，</span>
                </Trans>
              </p>
            </div>
          </FormItem>

          <FormItem label={t('镜像更新策略')}>
            <Select
              size="auto"
              options={ImagePullPolicyList}
              value={container.imagePullPolicy}
              onChange={value => {
                actions.editWorkload.updateContainer({ imagePullPolicy: value }, cKey);
              }}
            />
          </FormItem>

          <EditResourceContainerHealthCheckPanel cKey={cKey} />

          <FormItem label={t('特权级容器')}>
            <div className="form-unit is-success">
              <Switch
                value={container.privileged}
                disabled={false}
                onChange={value => actions.editWorkload.updateContainer({ privileged: value }, cKey)}
              />
              <p className="form-input-help text-label">{t('容器开启特权级，将拥有宿主机的root权限')}</p>
            </div>
          </FormItem>

          <FormItem label={t('权限集-增加')}>
            <SelectMultiple
              size="auto"
              options={AddCapabilitiesList.map(item => ({ value: item, text: item }))}
              value={container.addCapabilities}
              onChange={value => {
                actions.editWorkload.updateContainer({ addCapabilities: value }, cKey);
              }}
            />
          </FormItem>

          <FormItem label={t('权限集-删除')}>
            <SelectMultiple
              size="auto"
              options={DropCapabilitiesList.map(item => ({ value: item, text: item }))}
              value={container.dropCapabilities}
              onChange={value => {
                actions.editWorkload.updateContainer({ dropCapabilities: value }, cKey);
              }}
            />
          </FormItem>
        </ul>
      </div>
    );
  }
}
