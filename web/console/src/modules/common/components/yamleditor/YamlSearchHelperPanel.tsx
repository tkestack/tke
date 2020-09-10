import * as JsYAML from 'js-yaml';
import * as React from 'react';

import { Bubble, Tag } from '@tea/component';
import { t } from '@tencent/tea-app/lib/i18n';

export function YamlSearchHelperPanel(options: {
  isShow: boolean;
}) {
  let { isShow } = options;

  return (
    <Bubble
      visible={isShow}
      placement={'right-start'}
      content={
        <>
          <Tag theme="primary">Ctrl-F</Tag>
          {t('开始搜索')}
          <br />
          <Tag theme="primary">Ctrl-G</Tag>
          {t('下一个')}
          <br />
          <Tag theme="primary">Shift-Ctrl-G</Tag>
          {t('上一个')}
          <br />
          <Tag theme="primary">Shift-Ctrl-F</Tag>
          {t('替换')}
          <br />
          <Tag theme="primary">Shift-Ctrl-R</Tag>
          {t('替换全部')}
          {/* <br />
          <Tag theme="primary">Alt-F</Tag>
          {t('持久性搜索，对话框不会自动关闭，Enter键查找下一个，Shift-Enter键查找上一个')}
          <br />
          <Tag theme="primary">Alt-G</Tag>
          {t('跳转到行')} */}
        </>
      }
    />
  );
}
