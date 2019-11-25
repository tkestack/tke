import * as React from 'react';

export class FunctionIsDirected extends React.Component<any> {
    componentDidMount() {
        if (this.props.redirectTo) {
            window.open(this.props.redirectTo, '_blank');
        }
    }

    render() {
        return (
            <div style={{ position: 'relative' }}>
                <div className="dialog-panel" style={{ position: 'relative', padding: '25px' }}>
                    <div className="tc-15-rich-dialog m" role="alertdialog" style={{ margin: '100px auto' }}>
                        <div className="tc-15-rich-dialog-bd">
                            <div className="tc-icon-box">
                                <div className="col">
                                    <i className="icon-info-blue" style={{ marginRight: '10px' }} />
                                </div>
                                <div className="col">
                                    <h3 className="tc-dialog-title">系统提示</h3>

                                    <p>当前功能在腾讯云运营系统中配置，已进行跳转。</p>
                                    <p>未跳转？点击<a href={this.props.redirectTo} target="_blank">链接</a>进行跳转</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
