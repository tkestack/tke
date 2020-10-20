/**
 * hpa详情，yaml tab
 */
import React, { useState, useEffect, useContext, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty } from '@src/modules/common/utils';
import { fetchHPAYaml } from '../../../../WebAPI/scale';
import { insertCSS } from '@tencent/ff-redux/libs/qcloud-lib';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import { cutNsStartClusterId } from '@helper';
/**
 * 全局样式
 */
insertCSS(
    'HPAYamlEditorPanel',
    `
       .CodeMirror { height:600px; overflow:auto; overflow-x:hidden; font-size:12px }
       .CodeMirror-code>div{ line-height: 20px }
     `
);

const Yaml = React.memo(() => {
  const route = useSelector((state) => state.route);
  const { clusterId, HPAName, namespaceValue } = route.queries;

  /*
   * 编辑器参数
   */
  const readOnly = true;
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
    <CodeMirror
      className={'codeMirrorHeight'}
      value={!isEmpty(yamlData) && yamlData.recordCount ? yamlData.records[0] : t('')}
      options={codeOptions}
      onChange={(editor, data, value) => {
        // 配置项当中的value 不用props.config 是因为 更新之后，yaml的光标会默认跳转到末端
        // console.log('hpa codemirror...');
        // !readOnly && handleInputForEditor(value);
      }}
    />
  );
});
export default Yaml;
