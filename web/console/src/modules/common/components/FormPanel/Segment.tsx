import * as React from 'react';
import { Segment, SegmentProps } from '@tea/component';

interface FormPanelSegmentProps extends SegmentProps {}

function FormPanelSegment({ ...props }: FormPanelSegmentProps) {
  return <Segment {...props} />;
}

export { FormPanelSegment, FormPanelSegmentProps };
