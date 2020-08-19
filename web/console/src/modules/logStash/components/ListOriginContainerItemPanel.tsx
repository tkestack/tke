import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { LinkButton } from '../../common/components';
import { allActions } from '../actions';
import { ContainerLogs } from '../models';
import { isCanAddContainerLog } from './EditOriginContainerPanel';
import { RootProps } from './LogStashApp';

export interface ContainerItemProps extends RootProps {
  cKey: string;

  isEdit?: boolean; // 是否是编辑模式
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ListOriginContainerItemPanel extends React.Component<ContainerItemProps, any> {
  render() {
    let { actions, cKey, logStashEdit, namespaceList } = this.props,
      { containerLogs } = logStashEdit;
    let containerLogIndex = containerLogs.findIndex(item => item.id === cKey),
      containerLog: ContainerLogs = containerLogs[containerLogIndex];

    let { canSave, tip } = isCanAddContainerLog(containerLogs, namespaceList.data.recordCount);
    let canDelete = containerLogs.length > 1;

    return (
      <FormPanel fixed isNeedCard={false} style={{ minWidth: 600, padding: '20px  30px' }}>
        <div className="run-docker-box">
          <div className="justify-grid">
            <div className="col">
              <span className="text-label" style={{ verticalAlign: 'middle' }}>{`Namespace: `}</span>
              <span className="text">{containerLog.namespaceSelection}</span>
              <span className="text-label"> | </span>
              <span className="text-label">{t(`采集对象: `)}</span>
              <span className="text">
                {containerLog.collectorWay === 'container'
                  ? t('全部容器')
                  : t(
                      Object.keys(containerLog.workloadSelection)
                        .map(item => containerLog.workloadSelection[item].length)
                        .reduce((a, b) => a + b) + '个工作负载'
                    )}
              </span>
            </div>
            <div className="col">
              <LinkButton
                disabled={!canSave}
                tip={t('编辑')}
                errorTip={t('请完成待编辑项')}
                onClick={() => this._handleEditButton(canSave, containerLogIndex)}
              >
                <i className="icon-edit-gray" />
              </LinkButton>
              <LinkButton
                disabled={!canDelete}
                tip={t('删除')}
                errorTip={t('不可删除，至少创建一个容器')}
                onClick={() => actions.editLogStash.deleteContainerLog(containerLogIndex)}
              >
                <i className="icon-cancel-icon" />
              </LinkButton>
            </div>
          </div>
        </div>
      </FormPanel>
    );
  }

  /** 编辑按钮的相关操作 */
  _handleEditButton(canEdit: boolean, containerLogIndex: number) {
    let { actions, logStashEdit, route } = this.props,
      { containerLogs } = logStashEdit;

    // 需要判断当前是否有正在编辑的containerLog，如果有，需要判断编辑中的是否正确
    let editingLogIndex = containerLogs.findIndex(item => item.status === 'editing');

    if (canEdit) {
      editingLogIndex >= 0 && actions.editLogStash.updateContainerLog({ status: 'edited' }, editingLogIndex);
      actions.editLogStash.updateContainerLog({ status: 'editing' }, containerLogIndex);
    }
  }
}
