import { useState, useEffect, useRef, useCallback } from 'react';

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

export const usePrevious = value => {
  const ref = useRef();
  useEffect(() => {
    ref.current = value;
  });
  return ref.current;
};

export const useRefresh = () => {
  const [refreshFlag, setRefreshFlag] = useState(0);
  const triggerRefresh = useCallback(() => {
    setRefreshFlag(new Date().getTime());
  }, []);
  return {
    refreshFlag,
    triggerRefresh
  };
};
