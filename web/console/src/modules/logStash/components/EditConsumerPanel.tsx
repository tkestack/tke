import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from './LogStashApp';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { allActions } from '../actions';
import { InputField } from '../../common/components';
import { consumerModeList, clsRegionMap } from '../constants/Config';
import { Text, Segment } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { SegmentOption } from '@tencent/tea-component/lib/segment/SegmentOption';
import { FormPanel } from '@tencent/ff-component';

/** buttonBar的样式 */
const ButtonBarStyle = { marginBottom: '5px' };

/** 消费端类型的提示 */
const consumerModeTip = {
  kafka: t('将采集的日志消费到消息服务Kafka中。'),
  // cls: t('将采集的日志消费到日志服务CLS中。'),
  es: t('将采集的日志消费到Elasticsearch中。')
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class EditConsumerPanel extends React.Component<RootProps, any> {
  render() {
    let { actions, logStashEdit } = this.props,
      { consumerMode } = logStashEdit;

    /**
     * 渲染消费端的类型
     */
    let afterConsumerModeList: SegmentOption[] = consumerModeList.map(mode => {
      return {
        value: mode.value,
        text: mode.name
      };
    });

    let selectedConsumerMode = consumerModeList.find(item => item.value === consumerMode);

    return (
      <FormPanel.Item label={t('消费端')}>
        <FormPanel isNeedCard={false} fixed style={{ padding: 30, minWidth: 600 }}>
          <FormPanel.Item
            label={t('类型')}
            message={
              <React.Fragment>
                <Text>
                  <Trans>{consumerModeTip[consumerMode]}</Trans>
                </Text>
              </React.Fragment>
            }
          >
            <Segment
              options={afterConsumerModeList}
              value={selectedConsumerMode.value}
              onChange={value => actions.editLogStash.changeConsumerMode(value)}
            />
          </FormPanel.Item>

          {this._renderKafkaOption()}

          {this._renderESOption()}
        </FormPanel>
      </FormPanel.Item>
    );
  }

  /** 渲染kafaka 类型的编辑项 */
  private _renderKafkaOption() {
    let { actions, logStashEdit } = this.props,
      { addressIP, v_addressIP, addressPort, v_addressPort, topic, v_topic, consumerMode } = logStashEdit;

    return consumerMode !== 'kafka' ? (
      <noscript />
    ) : (
      [
        <FormPanel.Item key={0} label={t('访问地址')}>
          <InputField
            style={{ marginRight: '10px' }}
            type="text"
            placeholder={t('请输入IP地址')}
            tipMode="popup"
            validator={v_addressIP}
            value={addressIP}
            onChange={actions.editLogStash.inputAddressIP}
            onBlur={actions.validate.validateAddressIP}
          />

          <InputField
            type="text"
            placeholder={t('请输入IP端口号')}
            tipMode="popup"
            validator={v_addressPort}
            value={addressPort}
            onChange={actions.editLogStash.inputAddressPort}
            onBlur={actions.validate.validateAddressPort}
          />
        </FormPanel.Item>,

        <FormPanel.Item
          label={t('主题（Topic）')}
          key={1}
          validator={v_topic}
          message={t('最长64个字符，只能包含字母、数字、下划线、(".")及分隔符("-")')}
          input={{
            placeholder: t('请输入Topic'),
            value: topic,
            onChange: actions.editLogStash.inputTopic,
            onBlur: actions.validate.validateTopic
          }}
        />
      ]
    );
  }

  /** 渲染Elasticsearch */
  private _renderESOption() {
    let { actions, logStashEdit } = this.props,
      { consumerMode, esAddress, v_esAddress, indexName, v_indexName } = logStashEdit;

    return consumerMode !== 'es' ? (
      <noscript />
    ) : (
      [
        <FormPanel.Item
          label={t('Elasticsearch地址')}
          key={0}
          validator={v_esAddress}
          input={{
            placeholder: 'eg: http://190.0.0.1:9200',
            value: esAddress,
            onChange: value => {
              actions.editLogStash.inputEsAddress(value);
            },
            onBlur: actions.validate.validateEsAddress
          }}
        />,

        <FormPanel.Item
          label={t('索引')}
          key={1}
          validator={v_indexName}
          message={t('最长60个字符，只能包含小写字母、数字及分隔符("-"、"_"、"+")，且必须以小写字母开头')}
          input={{
            placeholder: 'eg: fluentd',
            value: indexName,
            onChange: value => {
              actions.editLogStash.inputIndexName(value);
            },
            onBlur: actions.validate.validateIndexName
          }}
        />
      ]
    );
  }
}
