import * as JsYAML from 'js-yaml';
import * as React from 'react';

import { Bubble, Badge } from '@tea/component';
import { t } from '@tencent/tea-app/lib/i18n';

export function YamlSearchHelperPanel(options: { isShow: boolean }) {
  let { isShow } = options;

  return (
    <Bubble
      visible={isShow}
      placement={'bottom'}
      content={
        <>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Ctrl-F / Cmd-F</Badge>
            <span>&nbsp;&nbsp;{t('开始搜索')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Ctrl-G / Cmd-G</Badge>
            <span>&nbsp;&nbsp;{t('下一个')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Shift-Ctrl-G / Shift-Cmd-G</Badge>
            <span>&nbsp;&nbsp;{t('上一个')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Shift-Ctrl-F / Shift-Cmd-F</Badge>
            <span>&nbsp;&nbsp;{t('替换')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Shift-Ctrl-R / Shift-Cmd-R</Badge>
            <span>&nbsp;&nbsp;{t('替换全部')}</span>
          </div>
          {/* <br />
          <Badge theme="success">Alt-F</Badge>
          {t('持久性搜索，对话框不会自动关闭，Enter键查找下一个，Shift-Enter键查找上一个')}
          <br />
          <Badge theme="success">Alt-G</Badge>
          {t('跳转到行')} */}
        </>
      }
    />
  );
}
