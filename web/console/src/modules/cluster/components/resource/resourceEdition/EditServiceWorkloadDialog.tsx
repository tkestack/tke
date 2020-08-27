import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Modal } from '@tea/component';
import { bindActionCreators, FetchState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { ButtonBar, FormItem } from '../../../../common/components';
import { FormLayout } from '../../../../common/layouts';
import { allActions } from '../../../actions';
import { ServiceWorkloadList } from '../../../constants/Config';
import { initSelector } from '../../../constants/initState';
import { Selector } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

const ButtonBarStyle = { marginBottom: '5px' };

// 加载中的样式
const loadingElement: JSX.Element = (
  <div style={{ verticalAlign: 'middle', display: 'inline-block' }}>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditServiceWorkloadDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot } = this.props,
      { serviceEdit } = subRoot,
      { isShowWorkloadDialog, workloadType, workloadList, workloadSelection } = serviceEdit;

    // 如果不需要展示 workload的弹窗选择
    if (!isShowWorkloadDialog) {
      return <noscript />;
    }

    /** 渲染资源类型列表，目前仅支持 deployment 和 statefulset */
    let selectWorkloadType = ServiceWorkloadList.find(item => item.value === workloadType);

    /** 渲染资源列表 */
    let workloadListOptions = workloadList.data.recordCount
      ? workloadList.data.records.map((item, index) => {
          return (
            <option key={index} value={item.metadata.name}>
              {item.metadata.name}
            </option>
          );
        })
      : [];

    const commonAction = () => {
      // 关闭窗口
      actions.editSerivce.workload.toggleIsShowWorkloadDialog();
    };

    const cancel = () => {
      commonAction();
    };

    const perform = () => {
      if (workloadSelection.length) {
        let labelsObj = workloadSelection[0].metadata.labels || {},
          labelsKey = Object.keys(labelsObj);

        let initSelectors: Selector[] = labelsKey.map(item => {
          return Object.assign({}, initSelector, {
            id: uuid(),
            key: item,
            value: labelsObj[item]
          });
        });

        actions.editSerivce.initSelectorFromWorkload(initSelectors);
        commonAction();
      }
    };

    return (
      <Modal visible={true} caption={t('引用Workload资源')} onClose={cancel} disableEscape={true} size="l">
        <Modal.Body>
          <FormLayout>
            <div className="param-box server-update add">
              <ul className="form-list jiqun fixed-layout">
                <FormItem label={t('资源类型')}>
                  <div className="form-unit" style={ButtonBarStyle}>
                    <ButtonBar
                      size="m"
                      list={ServiceWorkloadList}
                      selected={selectWorkloadType}
                      onSelect={item => {
                        actions.editSerivce.workload.selectWorkloadType(item.value + '');
                      }}
                    />
                  </div>
                </FormItem>
                <FormItem label={t('资源列表')}>
                  {workloadList.fetchState === FetchState.Fetching ? (
                    loadingElement
                  ) : workloadList.data.recordCount && workloadSelection.length > 0 ? (
                    <select
                      className="tc-15-select m"
                      style={{ maxWidth: '180px' }}
                      value={workloadSelection[0].metadata.name}
                      onChange={e => {
                        actions.editSerivce.workload.selectWorkload([
                          workloadList.data.records.find(item => item.metadata.name === e.target.value)
                        ]);
                      }}
                    >
                      {workloadListOptions}
                    </select>
                  ) : (
                    <div className="form-unit">
                      <select className="tc-15-select m" disabled>
                        <option key={-1}>{t('无可用资源')}</option>
                      </select>
                      <p className="text-label">
                        <Trans>
                          无可用资源，可前往
                          <a href="jvascript:;" onClick={e => this._handleClickForCreateWorkload()}>
                            资源控制台
                          </a>
                          新建
                        </Trans>
                      </p>
                    </div>
                  )}
                </FormItem>
                <FormItem label="Labels">
                  {workloadList.fetchState === FetchState.Fetching ? (
                    loadingElement
                  ) : (
                    <div className="form-unit">
                      {workloadSelection[0] ? (
                        workloadSelection[0].metadata.labels ? (
                          Object.keys(workloadSelection[0].metadata.labels).map((item, index) => {
                            return (
                              <p style={{ fontSize: '12px' }} key={index}>
                                {item + ': ' + workloadSelection[0].metadata.labels[item]}
                              </p>
                            );
                          })
                        ) : (
                          <p style={{ fontSize: '12px' }}>{t('Workload未设置Labels')}</p>
                        )
                      ) : (
                        <p style={{ fontSize: '12px' }}>{t('请先选择Workload')}</p>
                      )}
                    </div>
                  )}
                </FormItem>
              </ul>
            </div>
          </FormLayout>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" disabled={workloadSelection.length === 0} onClick={perform}>
            {t('确认')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }

  /** 资源控制台的体验 */
  private _handleClickForCreateWorkload() {
    let { actions, subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { serviceEdit } = subRoot;

    router.navigate(
      Object.assign({}, urlParams, { type: 'resource', resourceName: serviceEdit.workloadType }),
      route.queries
    );

    // 拉取对应的资源的数据
    actions.resource.initResourceInfoAndFetchData(true, serviceEdit.workloadType);
  }
}
