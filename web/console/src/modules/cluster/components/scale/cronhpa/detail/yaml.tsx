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
/**
 * hpa详情，yaml tab
 */
import React, { useState, useEffect, useContext, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty } from '@src/modules/common/utils';
import { fetchCronHpaYaml } from '../../../../WebAPI/scale';
import { insertCSS } from '@tencent/ff-redux/libs/qcloud-lib';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import { cutNsStartClusterId } from '@helper';
import { RecordSet } from '@tencent/ff-redux';
import { Resource } from '@src/modules/common';
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
  const route = useSelector(state => state.route);
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
  const [yamlData, setYamlData] = useState<RecordSet<string, any> | undefined>();
  useEffect(() => {
    async function getHPAYaml(namespace, clusterId, name) {
      const result = await fetchCronHpaYaml({ clusterId, name, namespace });
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
      value={yamlData?.recordCount ? yamlData.records[0] : t('')}
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
