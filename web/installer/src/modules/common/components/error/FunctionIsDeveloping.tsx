import * as React from 'react';

export class FunctionIsDeveloping extends React.Component<any> {
  render() {
    return (
      <div style={{ position: 'relative' }}>
        <div className="dialog-panel" style={{ position: 'relative', padding: '25px' }}>
          <div className="tc-15-rich-dialog m" role="alertdialog" style={{ margin: '100px auto' }}>
            <div className="tc-15-rich-dialog-bd">
              <div className="tc-icon-box">
                <div className="col">
                  <i className="icon-info-waiting" style={{ marginRight: '10px' }} />
                </div>
                <div className="col">
                  <h3 className="tc-dialog-title">系统提示</h3>

                  <p>当前功能模块尚在开发中，敬请期待</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
