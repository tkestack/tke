import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Alert, Bubble, Button, Input, Modal, Text, Form, Radio, InputNumber } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../../config';
import { cloneDeep, CreateResource, InputField, WorkflowDialog, TipInfo, initValidator } from '../../../../common';
import { allActions } from '../../../actions/';
import { RootProps } from '../../ClusterApp';
import { Pod } from '../../../models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

enum UpdateType {
  UserSelect = 'UserSelect', // 升级用户选中的pod
  ByPercentage = 'ByPercentage' // 按比例升级pod
}

@connect(state => state, mapDispatchToProps)
export class ResourceGrayUpgradeDialog extends React.Component<RootProps, {}> {
  state = {
    updateType: UpdateType.UserSelect,
    updatePercentage: 20
  };

  // 升级比例的自动步进
  get step(): number {
    return 10;
  }

  // 升级类型
  get updateType(): UpdateType {
    return this.state.updateType;
  }

  setUpdateType = (updateType: string): void => {
    this.setState({ updateType });
  };

  // 升级比例
  get updatePercentage(): number {
    return this.state.updatePercentage;
  }

  setUpdatePercentage = (updatePercentage: number): void => {
    this.setState({ updatePercentage });
  };

  // 用户选中的pods记录
  get podSelection(): Pod[] {
    return this.props.subRoot.resourceDetailState.podSelection;
  }

  // 全部pods记录
  get podRecords(): Pod[] {
    return this.props.subRoot.resourceDetailState.podList.data.records;
  }

  // 要升级的pods记录
  get targetPods(): Pod[] {
    let result = [];
    switch (this.updateType) {
      case UpdateType.UserSelect:
        result = this.podSelection;
        break;
      case UpdateType.ByPercentage:
        if (this.updatePercentage > 0) {
          // 取前x%
          result = this.podRecords.slice(0, Math.floor((this.podRecords.length * this.updatePercentage) / 100));
        } else {
          // 取后x%
          result = this.podRecords.slice(Math.floor((this.podRecords.length * (this.updatePercentage + 100)) / 100));
        }
        break;
      default:
        break;
    }
    return result;
  }

  render() {
    const { actions, region, clusterVersion, route } = this.props;
    const { updateGrayTappFlow, editTappGrayUpdate, podList, podSelection } = this.props.subRoot.resourceDetailState;
    const targetResource = this.props.subRoot.resourceOption.ffResourceList.selection;
    let resourceInfo = resourceConfig(clusterVersion)['tapp'];

    function getJsonData(pods, updateSettings) {
      let templates = {};
      let templatePool = {};
      pods.forEach(pod => {
        let indexName = pod.metadata.name.split(pod.metadata.generateName)[1]; //获取第几个实例
        if (indexName) {
          templates[indexName] = pod.metadata.name;
        }
        //将spec.template复制过来，然后修改caontainer里的name和imageName字段
        let newTemplate = cloneDeep(targetResource.spec.template);

        //editTapp只存储了必要的修改信息如containerName和containerImg，
        //所以更新的时候需要把container的之前的内容也带上,不然之前的内容会被清空
        newTemplate.spec.containers = pod.spec.containers.map(container => {
          let targetContainer = targetResource.spec.template.spec.containers.find(c => c.name === container.name);
          let containerSetting = updateSettings.containers.find(c => c.name === container.name);
          return {
            ...targetContainer,
            image: `${containerSetting.imageName}:${containerSetting.imageTag}`
          };
        });
        templatePool[pod.metadata.name] = newTemplate;
      });

      let jsonYaml = {
        spec: {
          templatePool,
          templates
        }
      };
      return jsonYaml;
    }

    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo: resourceInfo,
      namespace: route.queries['np'],
      clusterId: route.queries['clusterId'],
      resourceIns: route.queries['resourceIns'],
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
          const setResourceJSONData = resource => {
            let jsonYaml = {};
            jsonYaml = getJsonData(this.targetPods, editTappGrayUpdate);
            resource.jsonData = JSON.stringify(jsonYaml);
          };
          setResourceJSONData(resource);
          actions.resourceDetail.pod.podSelect([]);
        }}
      >
        <div style={{ maxHeight: '700px', overflowY: 'auto' }}>
          <TipInfo type="info">{t('灰度升级过程中，下列所有容器实例会同时重启')}</TipInfo>
          <FormPanel isNeedCard={false}>
            <Form.Item label="升级方式">
              <Radio.Group value={this.updateType} onChange={value => this.setUpdateType(value)}>
                <Radio name={UpdateType.UserSelect}>升级选中POD</Radio>
                <Radio name={UpdateType.ByPercentage}>按比例升级POD</Radio>
              </Radio.Group>
            </Form.Item>
            {this.updateType === UpdateType.ByPercentage && (
              <Form.Item label="升级比例">
                <InputNumber
                  min={-100}
                  max={100}
                  value={this.updatePercentage}
                  step={this.step}
                  onChange={this.setUpdatePercentage}
                />
              </Form.Item>
            )}
            <Form.Item label="待升级POD">
              <Form.Text>
                {this.targetPods.length > 0 ? `共${this.targetPods.length}个` : '-'}
                <ul>
                  {this.targetPods.slice(0, 3).map(item => (
                    <li>{item.metadata.name}</li>
                  ))}
                  {this.targetPods.length > 3 && (
                    <li>
                      <Bubble content={this.targetPods.map(item => item.metadata.name).join(', ')}>...</Bubble>
                    </li>
                  )}
                </ul>
              </Form.Text>
            </Form.Item>
            {editTappGrayUpdate && (
              <Form.Item label="镜像设置">
                {editTappGrayUpdate.containers.map((container, index_in) => {
                  return (
                    <FormPanel isNeedCard={false} fixed style={{ marginBottom: '20px' }} key={index_in}>
                      <FormPanel.Item label={t('镜像')}>
                        <div className={classnames('form-unit', { 'is-error': container.v_imageName.status === 2 })}>
                          <Bubble content={container.v_imageName.status === 2 ? container.v_imageName.message : null}>
                            <Input
                              size="full"
                              value={container.imageName}
                              onChange={value => {
                                actions.resourceDetail.pod.updateTappGrayUpdate(
                                  index_in,
                                  value,
                                  container.imageTag
                                );
                              }}
                              onBlur={e => {
                                actions.validate.workload.validateGrayUpdateRegistrySelection(
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
                          size="full"
                          value={container.imageTag}
                          onChange={value => {
                            actions.resourceDetail.pod.updateTappGrayUpdate(
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
              </Form.Item>
            )}
          </FormPanel>
        </div>
      </WorkflowDialog>
    );
  }

  /** 校验Tapp 灰度升级 是否正确 */
  private async _validation() {
    let result = true;
    let { editTappGrayUpdate } = this.props.subRoot.resourceDetailState;
    editTappGrayUpdate.containers.forEach(container => {
      result = result && container.v_imageName.status === 1;
    });
    return result;
  }
}
