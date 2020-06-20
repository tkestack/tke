import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Alert, Bubble, Button, Input, Modal, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../../config';
import { cloneDeep, CreateResource, InputField, WorkflowDialog, TipInfo } from '../../../../common';
import { allActions } from '../../../actions/';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceGrayUpgradeDialog extends React.Component<RootProps, {}> {
  render() {
    const { actions, region, clusterVersion, route } = this.props;
    const { updateGrayTappFlow, editTappGrayUpdate } = this.props.subRoot.resourceDetailState;
    const targetResource = this.props.subRoot.resourceOption.ffResourceList.selection;
    let resourceInfo = resourceConfig(clusterVersion)['tapp'];
    let templates = {};
    let templatePool = {};
    editTappGrayUpdate.forEach(item => {
      let indexName = item.name.split(item.generateName)[1]; //获取第几个实例
      if (indexName) {
        templates[indexName] = item.name;
      }
      //将spec.template复制过来，然后修改caontainer里的name和imageName字段
      let newTemplate = cloneDeep(targetResource.spec.template);

      //editTapp只存储了必要的修改信息如containerName和containerImg，
      //所以更新的时候需要把container的之前的内容也带上,不然之前的内容会被清空
      newTemplate.spec.containers = item.containers.map(container => {
        let targetContainer = targetResource.spec.template.spec.containers.find(c => c.name === container.name);
        return {
          ...targetContainer,
          name: container.name,
          image: `${container.imageName}:${container.imageTag}`
        };
      });
      templatePool[item.name] = newTemplate;
    });

    let jsonYaml = {
      spec: {
        templatePool,
        templates
      }
    };
    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo: resourceInfo,
      namespace: route.queries['np'],
      clusterId: route.queries['clusterId'],
      resourceIns: route.queries['resourceIns'],
      jsonData: JSON.stringify(jsonYaml),
      isStrategic: false
    };

    return (
      <WorkflowDialog
        width={750}
        caption={t('灰度升级')}
        workflow={updateGrayTappFlow}
        action={actions.workflow.updateGrayTapp}
        params={region.selection ? region.selection.value : ''}
        targets={[resource]}
        validateAction={this._validation.bind(this)}
        preAction={() => {
          actions.resourceDetail.pod.podSelect([]);
        }}
      >
        <div style={{ maxHeight: '700px', overflowY: 'auto' }}>
          <TipInfo type="info">{t('灰度升级过程中，下列所有容器实例会同时重启')}</TipInfo>

          <FormPanel isNeedCard={false} style={{ paddingRight: '80px' }}>
            {editTappGrayUpdate.map((item, index_out) => (
              <FormPanel.Item
                key={index_out}
                label={t(`容器实例 - ${item.name.split(item.generateName)[1]}`)}
                labelStyle={{ textAlign: 'center' }}
              >
                {item.containers.map((container, index_in) => {
                  return (
                    <FormPanel
                      isNeedCard={false}
                      fixed
                      style={{ padding: '35px', marginBottom: '20px' }}
                      key={index_in}
                    >
                      <FormPanel.Item label={t('名称')} text>
                        <span className="text-label">{container.name}</span>
                      </FormPanel.Item>
                      <FormPanel.Item label={t('镜像')}>
                        <div className={classnames('form-unit', { 'is-error': container.v_imageName.status === 2 })}>
                          <Bubble content={container.v_imageName.status === 2 ? container.v_imageName.message : null}>
                            <Input
                              value={container.imageName}
                              onChange={value => {
                                actions.resourceDetail.pod.updateTappGrayUpdate(
                                  index_out,
                                  index_in,
                                  value,
                                  container.imageTag
                                );
                              }}
                              onBlur={e => {
                                actions.validate.workload.validateGrayUpdateRegistrySelection(
                                  index_out,
                                  index_in,
                                  e.target.value
                                );
                              }}
                            />
                          </Bubble>
                        </div>
                      </FormPanel.Item>
                      <FormPanel.Item label={t('镜像版本（TAG）')}>
                        <Input
                          value={container.imageTag}
                          onChange={value => {
                            actions.resourceDetail.pod.updateTappGrayUpdate(
                              index_out,
                              index_in,
                              container.imageName,
                              value
                            );
                          }}
                        />
                      </FormPanel.Item>
                    </FormPanel>
                  );
                })}
              </FormPanel.Item>
            ))}
          </FormPanel>
        </div>
      </WorkflowDialog>
    );
  }

  /** 校验Tapp 灰度升级 是否正确 */
  private async _validation() {
    await this.props.actions.validator.workload.validateGrayUpdate();
    let result = true;
    let { editTappGrayUpdate } = this.props.subRoot.resourceDetailState;
    editTappGrayUpdate.forEach(item => {
      item.containers.forEach(container => {
        result = result && container.v_imageName.status === 1;
      });
    });
    return result;
  }
}
