import * as React from 'react';
import { RootProps } from './AlarmPolicyApp';
import { Radio, Checkbox } from '@tea/component';
import { AlarmObjectsType, workloadTypeList } from '../constants/Config';
import { FormItem, SelectList } from '../../common/components';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export class EditAlarmPolicyObject extends React.Component<RootProps, {}> {
  renderPodList() {
    let Tip = content => {
      return (
        <div className="colony" style={{ fontSize: '12px' }}>
          <span>{content}</span>
        </div>
      );
    };
    let { workloadList, alarmPolicyEdition, actions } = this.props;
    if (workloadList.fetchState === FetchState.Fetching) {
      return Tip(t('加载中'));
    } else if (workloadList.fetchState === FetchState.Failed) {
      return Tip(t('加载失败'));
    } else if (workloadList.data.recordCount === 0) {
      return Tip(t('该命名空间下无workload'));
    } else {
      // 根据 PodList 初始化 checkedList
      let checkedList = [];
      let checkboxList = workloadList.data.records.map(workload => {
        // let checkboxItem = {
        //   value: workload.metadata.name,
        //   label: workload.metadata.name
        // };
        let checkbox = <Checkbox name={workload.metadata.name}>{workload.metadata.name}</Checkbox>;
        if (alarmPolicyEdition.alarmObjects.find(object => object === workload.metadata.name)) {
          checkedList.push(workload.metadata.name);
        }

        return checkbox;
      });
      return (
        <Checkbox.Group
          onChange={items => {
            // let checkedArray = items.map(item => item);
            actions.alarmPolicy.inputAlarmPolicyObjects(items);
          }}
          value={checkedList}
          layout="column"
        >
          {checkboxList}
        </Checkbox.Group>
      );
    }
  }

  renderRadioList(type) {
    let { alarmPolicyEdition, actions, namespaceList } = this.props;
    if (type === 'cluster' || type === '') {
      return <noscript />;
    }
    let radioList: JSX.Element[] = [];
    AlarmObjectsType[type].forEach((item, index) => {
      radioList.push(
        <div className="form-unit unit-group new-strategy-alarm-object">
          <div className="alarm-select">
            <Radio key={index} name={item.value} disabled={item.value === 'k8sLabel'}>
              {item.text}
              <span className="text-label">{item.tip}</span>
            </Radio>
          </div>
          <div className="alarm-write">
            {item.value === 'part' && alarmPolicyEdition.alarmObjectsType === 'part' && (
              <ul className="form-list fixed-layout">
                <FormItem label="Namespace">
                  <SelectList
                    value={alarmPolicyEdition.alarmObjectNamespace + ''}
                    recordList={namespaceList.data.records}
                    valueField="name"
                    textField="name"
                    textFields={['name']}
                    textFormat={`\${name}`}
                    className="tc-15-select m"
                    style={{ marginRight: '5px' }}
                    onSelect={value => {
                      actions.namespace.selectNamespace(value);
                    }}
                    isUnshiftDefaultItem={false}
                  />
                </FormItem>
                <FormItem label="WorkloadType" isNeedFormInput={false}>
                  <SelectList
                    value={alarmPolicyEdition.alarmObjectWorkloadType + ''}
                    recordList={workloadTypeList}
                    valueField="value"
                    textField="name"
                    textFields={['label']}
                    textFormat={`\${label}`}
                    className="tc-15-select m"
                    style={{ marginRight: '5px' }}
                    onSelect={value => actions.alarmPolicy.inputAlarmObjectWorkloadType(value)}
                    isUnshiftDefaultItem={false}
                  />
                  <div
                    className="param-box"
                    style={{
                      backgroundColor: '#fff',
                      padding: '5px 10px',
                      marginTop: '5px',
                      marginBottom: '10px',
                      width: '250px'
                    }}
                  >
                    <div className="param-bd">
                      <ul>
                        <div className="tc-g-u-1-3" style={{ width: '100%', maxHeight: '180px', overflowY: 'auto' }}>
                          <div className="colony-list">{this.renderPodList()}</div>
                        </div>
                      </ul>
                    </div>
                  </div>
                  <div className="is-error">
                    <p className="form-input-help" style={{ fontSize: '12px' }}>
                      {alarmPolicyEdition.v_alarmObjects.status === 2 && alarmPolicyEdition.v_alarmObjects.message}
                    </p>
                  </div>
                </FormItem>
              </ul>
            )}
          </div>
        </div>
      );
    });
    return radioList;
  }
  render() {
    let { actions, alarmPolicyEdition } = this.props;

    let isShow = alarmPolicyEdition.alarmPolicyType !== 'cluster' && alarmPolicyEdition.alarmPolicyType !== '';
    return (
      <FormItem isShow={isShow} label={t('告警对象')} isNeedFormInput={false}>
        <Radio.Group
          value={alarmPolicyEdition.alarmObjectsType}
          onChange={value => actions.alarmPolicy.inputAlarmPolicyObjectsType(value)}
          // className="form-unit"
        >
          {this.renderRadioList(alarmPolicyEdition.alarmPolicyType)}
        </Radio.Group>
      </FormItem>
    );
  }
}
