import * as React from 'react';
import { connect } from 'react-redux';
import { CreateResource } from 'src/modules/common';

import { Bubble, Button, Modal, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../../config';
import { TipInfo } from '../../../../common/components';
import { getWorkflowError, isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validatorActions } from '../../../actions/validatorActions';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UpdateClusterAllocationRatioDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot, route, cluster } = this.props,
      {
        clusterAllocationRatioEdition: { isUseCpu, isUseMemory, cpuRatio, memoryRatio, v_cpuRatio, v_memoryRatio },
        updateClusterAllocationRatio
      } = subRoot;
    let action = actions.workflow.updateClusterAllocationRatio;
    let workflow = updateClusterAllocationRatio;

    if (workflow.operationState === OperationState.Pending) {
      return <noscript />;
    }

    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }
      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }
    };

    const perform = () => {
      let { clusterAllocationRatioEdition } = this.props.subRoot;
      actions.validate.validateAllClusterAllocationRatio();
      if (validatorActions._validateAllClusterAllocationRatio(clusterAllocationRatioEdition)) {
        let clusterInfo = resourceConfig(this.props.clusterVersion).cluster;
        let data = {
          spec: {
            properties: {
              oversoldRatio: {
                cpu: clusterAllocationRatioEdition.isUseCpu ? clusterAllocationRatioEdition.cpuRatio : null,
                memory: clusterAllocationRatioEdition.isUseMemory ? clusterAllocationRatioEdition.memoryRatio : null
              }
            }
          }
        };

        let createClusterData: CreateResource[] = [
          {
            id: uuid(),
            resourceInfo: clusterInfo,
            mode: 'update',
            isStrategic: false,
            resourceIns: route.queries['clusterId'],
            jsonData: JSON.stringify(data)
          }
        ];
        action.start(createClusterData, {
          clusterId: route.queries['clusterId']
        });
        action.perform();
      }
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <Modal visible={true} caption={t('编辑集群超售比')} onClose={cancel} size={500} disableEscape={true}>
        <Modal.Body>
          <FormPanel isNeedCard={false}>
            <FormPanel.Item
              label="超售比"
              tips={t('设置集群超售比')}
              message={t('超售比只能为正数,且只能精确到小数点后2位')}
            >
              <div style={{ marginBottom: '5px' }}>
                <FormPanel.InlineText parent="span" style={{ width: '100px' }}>
                  CPU
                </FormPanel.InlineText>
                <FormPanel.Switch
                  value={isUseCpu}
                  onChange={value => {
                    actions.cluster.updateClusterAllocationRatio({ isUseCpu: value });
                  }}
                />
                {isUseCpu && (
                  <Bubble content={v_cpuRatio.status === 2 ? v_cpuRatio.message : null}>
                    <div
                      className={v_cpuRatio.status === 2 ? 'is-error' : ''}
                      style={{ display: 'inline-block', marginLeft: 5 }}
                    >
                      <FormPanel.Input
                        size={'s'}
                        value={cpuRatio}
                        onChange={value => {
                          actions.cluster.updateClusterAllocationRatio({ cpuRatio: value });
                        }}
                        onBlur={e => {
                          actions.validator.validateClusterAllocationRatio('cpu', e.target.value);
                        }}
                      />
                    </div>
                  </Bubble>
                )}
              </div>
              <div>
                <FormPanel.InlineText parent="span" style={{ width: '100px' }}>
                  Memory
                </FormPanel.InlineText>
                <FormPanel.Switch
                  value={isUseMemory}
                  onChange={value => {
                    actions.cluster.updateClusterAllocationRatio({ isUseMemory: value });
                  }}
                />
                {isUseMemory && (
                  <Bubble content={v_memoryRatio.status === 2 ? v_memoryRatio.message : null}>
                    <div
                      className={v_memoryRatio.status === 2 ? 'is-error' : ''}
                      style={{ display: 'inline-block', marginLeft: 5 }}
                    >
                      <FormPanel.Input
                        size={'s'}
                        value={memoryRatio}
                        onChange={value => {
                          actions.cluster.updateClusterAllocationRatio({ memoryRatio: value });
                        }}
                        onBlur={e => {
                          actions.validator.validateClusterAllocationRatio('memory', e.target.value);
                        }}
                      />
                    </div>
                  </Bubble>
                )}
              </div>
            </FormPanel.Item>
          </FormPanel>
          {failed && <TipInfo type="error">{getWorkflowError(workflow)}</TipInfo>}
        </Modal.Body>
        <Modal.Body>
          <Modal.Footer>
            <Button type="primary" disabled={workflow.operationState === OperationState.Performing} onClick={perform}>
              {failed ? t('重试') : t('确定')}
            </Button>
            <Button onClick={cancel}>{t('取消')}</Button>
          </Modal.Footer>
        </Modal.Body>
      </Modal>
    );
  }
}
