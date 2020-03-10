import { WorkflowState } from '@tencent/ff-redux';

export const getWorkflowError = (workflow: WorkflowState<any, any>) => {
  return (
    workflow &&
    workflow.results &&
    workflow.results.length &&
    workflow.results[0].error &&
    (workflow.results[0].error.message || workflow.results[0].error.Message)
  );
};

/**获取错误码 */
export const getWorkflowErrorCode = (workflow: WorkflowState<any, any>) => {
  let reg = /\(-\w+\)/g,
    code = 0;
  if (workflow && workflow.results && workflow.results.length && workflow.results[0].error) {
    let msg = workflow.results[0].error.message || '',
      matches = msg.match(reg);
    if (matches && matches[0]) {
      let codeStr = matches[0].substring(1, matches[0].length - 1);
      code = parseInt(codeStr);
    }
  }

  return code;
};
