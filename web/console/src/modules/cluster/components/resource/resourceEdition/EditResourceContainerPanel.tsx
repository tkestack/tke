import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble } from '@tea/component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem } from '../../../../common/components';
import { isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validateWorkloadActions } from '../../../actions/validateWorkloadActions';
import { ContainerItem } from '../../../models';
import { RootProps } from '../../ClusterApp';
import { EditResourceContainerItem } from './EditResourceContainerItem';
import { ResourceContainerListItem } from './ResourceContainerListItem';

insertCSS(
  'EditResourceContainerPanel',
  `
    .container-panel.tc-15-bubble-icon {
        width: 100%;
    }
`
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerPanel extends React.Component<RootProps, {}> {
  render() {
    let { subRoot } = this.props,
      { workloadEdit } = subRoot,
      { canAddContainer, containers, volumes } = workloadEdit;

    let editingContainer = containers.find(c => c.status === 'editing');
    let canAdd = canAddContainer;

    // 判断是否能够新增容器
    canAdd = isEmpty(editingContainer) || validateWorkloadActions._canAddContainer(editingContainer, volumes);

    return (
      <FormItem label={t('实例内容器')}>
        <div className="form-unit">
          {this._renderContainerList()}
          <p className="text-label">{t('注意：Workload创建完成后，容器的配置信息可以通过更新YAML的方式进行修改')}</p>
          {canAdd ? (
            <a href="javascript:;" className="add-btn" onClick={() => this._handleAddContainer(editingContainer)}>
              {t('添加容器')}
            </a>
          ) : (
            <Bubble
              // className="container-panel"  TODO:tea2.0
              placement="bottom"
              content={t('请先完成待编辑项')}
            >
              <a href="javascript:;" className="add-btn disabled">
                {t('添加容器')}
              </a>
            </Bubble>
          )}
        </div>
      </FormItem>
    );
  }

  /** 容器的两种展现形态 */
  private _renderContainerList() {
    let { subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    return containers.map((container, index) => {
      return container.status === 'edited' ? (
        <ResourceContainerListItem key={index} cKey={container.id + ''} />
      ) : (
        <EditResourceContainerItem key={index} cKey={container.id + ''} />
      );
    });
  }

  /** 新增容器 */
  private _handleAddContainer(container: ContainerItem) {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { volumes, containers } = workloadEdit;

    actions.validate.workload.validateContainer(container);

    // 新增容器
    if (container) {
      if (validateWorkloadActions._validateContainer(container, volumes, containers)) {
        actions.editWorkload.updateContainer({ status: 'edited' }, container.id + '');
        actions.editWorkload.addContainer();
      }
    } else {
      actions.editWorkload.addContainer();
    }
  }
}
