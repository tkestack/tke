import * as React from "react";
import classnames from "classnames";
import { BaseReactProps } from "tea-components/libs/types/BaseReactProps";

export interface BubbleWrapperProps extends BaseReactProps {
    /**
     * 气泡对齐到元素的位置，默认不传为居中
     * */
    align?: "start" | "end"
}

export function BubbleWrapper({ children, className, align }: BubbleWrapperProps) {
    const finalClassName = classnames("tc-15-bubble-icon", {
        ["tc-15-triangle-align-" + align]: !!align
    }, className);
    return (
        <div className={finalClassName}>
            { children }
        </div>
    );
}