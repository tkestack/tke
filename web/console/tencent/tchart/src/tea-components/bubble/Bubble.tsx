import * as React from "react";
import classnames from "classnames";
import { BaseReactProps } from "../../tea-components/libs/types/BaseReactProps";

export interface BubbleProps extends BaseReactProps {
    /**
     * 小三角的朝向，默认为 top
     */
    focusPosition?: "top" | "right" | "left" | "bottom";

    /**
     * 小三角的偏移位置，可以给百分比，也可以给像素值，默认为 "50%"
     * */
    focusOffset?: string | number;
}

export function Bubble({ focusPosition = "top" as "top", focusOffset, children, className, style }: BubbleProps) {

    const offsetProperty = focusPosition === "top" || focusPosition === "bottom" ? "top" : "left";

    style = Object.assign({}, { [offsetProperty]: focusOffset }, style);

    return (

        <div className={classnames(`tc-15-bubble tc-15-bubble-${focusPosition}`, className)} style={ style }>
            <div className="tc-15-bubble-inner">
                { children }
            </div>
        </div>
    );
}
