import classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Input, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Icon } from '@tencent/tea-component';

import { LinkButton } from '../../common/components';
import { allActions } from '../actions';
import { logModeList } from '../constants/Config';
import { RootProps } from './LogStashApp';

/** 标签编辑的内联样式 */
const MetadataStyle: React.CSSProperties = { display: 'inline-block' };

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class EditOriginNodePanel extends React.Component<RootProps, any> {
  render() {
    let { actions, logStashEdit } = this.props,
      { v_nodeLogPath, nodeLogPath, metadatas, logMode } = logStashEdit;

    let isCanNotAdd = false;
    metadatas.forEach(item => {
      if (item.v_metadataKey.status === 2 || item.v_metadataValue.status === 2) {
        isCanNotAdd = true;
      }
    });

    return (
      <FormPanel.Item label={t('日志源')} isShow={logMode === logModeList.node.value} style={{ minWidth: '550px' }}>
        <FormPanel isNeedCard={false} fixed style={{ minWidth: 600, padding: '30px' }}>
          <FormPanel.Item
            label={t('收集路径')}
            message={t('指定采集日志的文件路径，支持通配符(*)，支持通配符（*），必须以`/`开头')}
            validator={v_nodeLogPath}
            input={{
              placeholder: t('请输入日志搜集路径，如/data/log/*.log'),
              value: nodeLogPath,
              onChange: value => actions.editLogStash.inputNodeLogPath(value),
              onBlur: actions.validate.validateNodeLogPath
            }}
          />

          <FormPanel.Item
            label="metadata"
            message={
              <React.Fragment>
                <Text parent="p">{t('收集规则收集的日志会带上metadata，并上报到消费端')}</Text>
                <Text parent="p">
                  {metadatas.length >= 1
                    ? t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')
                    : ''}
                </Text>
              </React.Fragment>
            }
          >
            {this._renderMetadataItemList()}
            <p style={{ marginTop: '8px' }}>
              <LinkButton
                disabled={isCanNotAdd}
                errorTip={t('请先完成待编辑项')}
                onClick={() => actions.editLogStash.addMetadata()}
              >
                {t('新增metadata')}
              </LinkButton>
            </p>
          </FormPanel.Item>
        </FormPanel>
      </FormPanel.Item>
    );
  }

  /** 渲染metadata的编辑项 */
  private _renderMetadataItemList() {
    let { actions, logStashEdit } = this.props,
      { metadatas } = logStashEdit;

    return metadatas.map((metadata, index) => {
      return (
        <div key={index} className="code-list" style={{ padding: '2px' }}>
          <div className={classnames({ 'is-error': metadata.v_metadataKey.status === 2 })} style={MetadataStyle}>
            <Bubble content={metadata.v_metadataKey.status === 2 ? metadata.v_metadataKey.message : null}>
              <Input
                style={{ width: '128px' }}
                placeholder="key"
                value={metadata.metadataKey}
                onChange={value => {
                  actions.editLogStash.updateMetadata({ metadataKey: value }, index);
                }}
                onBlur={() => actions.validate.validateMetadataItem({ metadataKey: metadata.metadataKey }, index)}
              />
            </Bubble>
          </div>

          <span className="text-label" style={{ padding: '5px', verticalAlign: 'middle' }}>
            =
          </span>
          <div className={classnames({ 'is-error': metadata.v_metadataValue.status === 2 })} style={MetadataStyle}>
            <Bubble content={metadata.v_metadataValue.status === 2 ? metadata.v_metadataValue.message : null}>
              <Input
                style={{ width: '128px' }}
                placeholder="key"
                value={metadata.metadataValue}
                onChange={value => {
                  actions.editLogStash.updateMetadata({ metadataValue: value }, index);
                }}
                onBlur={() => actions.validate.validateMetadataItem({ metadataValue: metadata.metadataValue }, index)}
              />
            </Bubble>
          </div>
          <LinkButton onClick={() => actions.editLogStash.deleteMetadata(index)} tip="删除" tipDirection="right">
            <Icon type="close" />
          </LinkButton>
        </div>
      );
    });
  }
}
