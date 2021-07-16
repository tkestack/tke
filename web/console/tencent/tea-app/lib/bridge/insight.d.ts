export declare const insight: {
    /**
     * 上报点击流
     * @param tag 点击标识
     */
    reportHotTag(tag: string): void;
    /**
     * 自定义上报，如果 ns 和 event 参数不正确，会扔出异常
     *
     * 控制台在预览模式或开发模式不会上报，如果要验证效果，请在手动添加标识到 localStorage:
     *
     * ```js
     * localStorage.debugInsight = true
     * ```
     *
     * 然后刷新页面。
     *
     * @example
     *
      ```js
      import { app } from '@tea/app';
  
      app.insight.stat({
        ns: 'cvm',
        event: 'restart',
        stringFields: {
          instance: 'ins-1d7x8C3s'
        },
        integerFields: {
          cost: 1000
        }
      });
    ```
     */
    stat({ ns, event, integerFields, stringFields }: InsightStat): void;
};
export interface InsightStat {
    /**
     * 为了不与其他业务的数据冲突，需要制定上报的命名空间，如 "cvm"
     *
     * 命名空间允许使用的字符：字母、数字、斜杠 `/`（不允许出现在首尾），如：
     *
     *  - "cvm2"
     *  - "cvm/sg"
     */
    ns: string;
    /**
     * 本次统计的事件名称
     *
     * 事件名允许使用的字符：字母、数字、横杠 `-`（不允许出现在首尾），如：
     *
     *  - "restart-fail"
     *  - "move2recycle"
     */
    event: string;
    /**
     * 自定义数据（整数类型）
     *
     * 如果传入的是不是 number 类型，会被忽略。如果传入的不是整数，会使用 Math.round 取整。
     */
    integerFields?: Record<any, number>;
    /**
     * 自定义数据（字符串类型）
     *
     * 如果传入的不是 string 类型，会被忽略。
     *
     * 单个字段的长度，不能超过 4096，否则会被裁切
     */
    stringFields?: Record<any, string>;
}
