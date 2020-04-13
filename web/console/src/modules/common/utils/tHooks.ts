import { useState } from 'react';

export const useModal = (isShowingParam = false) => {
    const [isShowing, setIsShowing] = useState(isShowingParam);

    function toggle() {
        setIsShowing(!isShowing);
    }

    return {
        isShowing,
        toggle
    };
};
