import React from 'react';
// also exported from '@storybook/react' if you can deal with breaking changes in 6.1
import { Story, Meta } from '@storybook/react/types-6-0';

import {
  CodeMirrorEditor,
  CodeMirrorEditorProps
} from '../../src/modules/common/components/codemirror/CodeMirrorEditor';

export default {
  title: 'Example/CodeMirrorEditor',
  component: CodeMirrorEditor,
  argTypes: {
    backgroundColor: { control: 'color' }
  }
} as Meta;

const Template: Story<CodeMirrorEditorProps> = args => <CodeMirrorEditor {...args} />;

export const Default = Template.bind({});
