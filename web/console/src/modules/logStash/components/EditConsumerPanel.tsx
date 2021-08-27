/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify, Button, Segment, Text } from '@tencent/tea-component';
import { SegmentOption } from '@tencent/tea-component/lib/segment/SegmentOption';

import { InputField } from '../../common/components';
import { allActions } from '../actions';
import { clsRegionMap, consumerModeList } from '../constants/Config';
import { RootProps } from './LogStashApp';
import * as WebAPI from '../WebAPI';

interface PropTypes extends RootProps {
  onESStatusChange: (status: number) => void;
}

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
export class EditConsumerPanel extends React.Component<PropTypes, any> {
  constructor(props) {
    super(props);

    this.esDetection = this.esDetection.bind(this);

    /** checkESStatus
     * 0: init
     * 1: need es detection
     * 2: start detecting
     * 3: result success
     * 4: result failure
     */
    this.state = {
      checkESStatus: 0
    };
  }

  async esDetection() {
    this.changeESStatus(2);

    let { logStashEdit } = this.props;
    let { esAddress, esUsername, esPassword } = logStashEdit;
    let [scheme, address] = esAddress.split('://');
    let [host, port] = address.split(':');
    let ret = await WebAPI.fetchEsDetection({
      scheme: scheme,
      host: host,
      port: port,
      user: esUsername,
      password: esPassword
    }) ? 3 : 4;

    this.changeESStatus(ret);
  }

  changeESStatus(status) {
    this.setState({
      checkESStatus: status
    });
    this.props.onESStatusChange(status);
  }

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
              onChange={value => {
                actions.editLogStash.changeConsumerMode(value);
                this.changeESStatus(value === 'es' ? 1 : 0);
              }}
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
      { consumerMode, esAddress, v_esAddress, indexName, v_indexName, esUsername, esPassword } = logStashEdit;

    const checkESStatus = this.state.checkESStatus;

    const esStatusMsg = {
      2: {
        color: 'primary',
        text: '连接中...'
      },
      3: {
        color: 'success',
        text: '连接成功！点击下方【完成】以设置事件持久化'
      },
      4: {
        color: 'warning',
        text: '连接失败！请检查 ElasticSearch 相关配置，注意开启了用户验证的 ElasticSearch 需要输入用户名和密码'
      }
    };

    const esStatusMsgColor = esStatusMsg[checkESStatus] ? esStatusMsg[checkESStatus].color : 'text';
    const esStatusMsgText = esStatusMsg[checkESStatus] ? esStatusMsg[checkESStatus].text : '';

    return consumerMode !== 'es' ? (
      <noscript />
    ) : (
      [
        <FormPanel.Item
          required
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
          required
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
        />,

        <FormPanel.Item
          label={t('用户名')}
          key={2}
          input={{
            style: {
              width: '300px'
            },
            placeholder: t('仅需要用户验证的 Elasticsearch 需要填入用户名'),
            value: esUsername,
            onChange: value => {
              actions.editLogStash.inputEsUsername(value);
            }
          }}
        />,

        <FormPanel.Item
          label={t('密码')}
          key={3}
          input={{
            style: {
              width: '300px'
            },
            type: 'password',
            placeholder: t('仅需要用户验证的 Elasticsearch 需要填入密码'),
            value: esPassword,
            onChange: value => {
              actions.editLogStash.inputEsPassword(value);
            }
          }}
        />,

        <FormPanel.Item key={4} message={<Text theme={esStatusMsgColor}>{esStatusMsgText}</Text>}>
          <Justify
            left={
              <React.Fragment>
                <Button
                  disabled={!esAddress.length || !indexName.length}
                  type="primary"
                  style={{ marginRight: 20 }}
                  onClick={() => {
                    this.esDetection();
                  }}
                >
                  检测连接
                </Button>
              </React.Fragment>
            }
          />
        </FormPanel.Item>
      ]
    );
  }
}
