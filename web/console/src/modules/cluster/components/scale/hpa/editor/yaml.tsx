/**
 * yaml编辑组件
 */
import React, { useState, useEffect, useContext, useMemo, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty } from '@src/modules/common/utils';
import {
  Layout,
  Card,
  Button,
  Affix
} from '@tencent/tea-component';
import { fetchHPAYaml, modifyHPAYaml } from '@src/modules/cluster/WebAPI/scale';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import { cutNsStartClusterId } from '@helper';

const { Body, Content } = Layout;

const Yaml = React.memo((props) => {
  const route = useSelector((state) => state.route);
  const { clusterId, HPAName, namespaceValue } = route.queries;
  const [newYamlData, setNewYamlData] = useState();

  /*
   * 编辑器参数
   */
  const readOnly = false;
  const codeOptions = {
    lineNumbers: true,
    mode: 'text/x-yaml',
    theme: 'monokai',
    readOnly: readOnly ? true : false, // nocursor表明焦点不能展示，不会展示光标
    spellcheck: true, // 是否开启单词校验
    autocorrect: true, // 是否开启自动修正
    lineWrapping: true, // 自动换行
    styleActiveLine: true, // 当前行背景高亮
    tabSize: 2 // tab 默认是2格
  };

  /*
   * 处理外层滚动
   */
  const bottomAffixRef = useRef(null);
  useEffect(() => {
    const body = document.querySelector('.tea-web-body');
    if (!body) {
      return () => null;
    }
    const handleScroll = () => {
      bottomAffixRef.current.update();
    };
    body.addEventListener('scroll', handleScroll);
    return () => body.removeEventListener('scroll', handleScroll);
  }, []);

  /**
   * 获取yaml数据
   */
  const [yamlData, setYamlData] = useState();
  useEffect(() => {
    async function getHPAYaml(namespace, clusterId, name) {
      const result = await fetchHPAYaml({ clusterId, name, namespace });
      setYamlData(result);
    }
    if (clusterId && HPAName && namespaceValue) {
      const newNamespace = cutNsStartClusterId({ namespace: namespaceValue, clusterId });
      getHPAYaml(newNamespace, clusterId, HPAName);
    }
  }, [clusterId, HPAName, namespaceValue]);

  return (
    <Layout>
      <Body>
        <Content>
          <Content.Header
            showBackButton
            onBackButtonClick={() => history.back()}
            title={t('更新YAML')}
          />
          <Content.Body>
            <CodeMirror
              className={'codeMirrorHeight'}
              value={!isEmpty(yamlData) && yamlData.recordCount ? yamlData.records[0] : t('')}
              options={codeOptions}
              onChange={(editor, data, value) => {
                    // 配置项当中的value 不用props.config 是因为 更新之后，yaml的光标会默认跳转到末端
                    setNewYamlData(value);
                  }}
            />
          </Content.Body>
        </Content>
        <Affix ref={bottomAffixRef} offsetBottom={0} style={{ zIndex: 5 }}>
          <Card>
            <Card.Body style={{ borderTop: '1px solid #ddd' }}>
              <Button type="primary" onClick={async () => {
                  await modifyHPAYaml({ clusterId, name: HPAName, namespace: cutNsStartClusterId({ namespace: namespaceValue, clusterId }), yamlData: newYamlData });
                  history.back();
                }}>
                <Trans>保存</Trans>
              </Button>
              <Button
                style={{ marginLeft: '10px' }}
                onClick={(e) => {
                    e.preventDefault();
                    history.back();
                  }}
                >
                <Trans>取消</Trans>
              </Button>
              {/*</Form.Action>*/}
            </Card.Body>
          </Card>
        </Affix>
      </Body>
    </Layout>
  );
});
export default Yaml;
