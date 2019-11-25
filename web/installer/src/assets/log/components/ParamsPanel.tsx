import * as React from 'react';
import { RootProps } from './LogApp';
import { CodeMirrorEditor } from '../../common/components';

interface ParamsPanelProps {
    params?: string;
}

export class ParamsPanel extends React.Component<ParamsPanelProps, void> {
    render() {
        let { params } = this.props;
        return (
            <CodeMirrorEditor
                title=""
                height={420}
                width={670}
                dHeight={420}
                isShowHeader={false}
                isOpenClip={false}
                isOpenDialogEditor={false}
                theme="monokai"
                value={params || ''}
                readOnly={true}
                isForceRefresh={true}
            />
        );
    }
}
