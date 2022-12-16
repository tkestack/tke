import React, { useEffect, useState, useRef } from 'react';

interface IPaging {
  pageIndex: number;
  pageSize: number;
}

interface IQueryResponse<T> {
  data: T;

  // 总数
  totalCount?: number;

  // continueToken
  continueToken?: string;
}

interface IQueryParams {
  paging?: IPaging;
  continueToken?: string;
}

interface IUseFetchOptions<T> {
  mode?: 'normal' | 'paging' | 'continue';
  onDataChange?: ({ data }: { data: T }) => void;
  defaultPageSize?: number;
  fetchAble?: boolean;
  polling?: boolean;
  pollingDelay?: number;
  needClearData?: boolean;
  onlyPollingPage1?: boolean;
}

type IUseFetchQuery<T> = (params?: IQueryParams) => Promise<IQueryResponse<T>>;

type IStatus = 'idle' | 'loading' | 'success' | 'error' | 'loading-polling' | 'expired';

export function useFetch<T>(
  query: IUseFetchQuery<T>,

  deps: React.DependencyList,

  options?: IUseFetchOptions<T>
) {
  const {
    mode = 'normal',
    onDataChange = () => {},
    defaultPageSize = 20,
    fetchAble = true,
    polling = false,
    pollingDelay = 5000,
    needClearData = true,
    onlyPollingPage1 = false
  } = options ?? {};

  const [data, _setData] = useState<T>(null);
  function setData(data: T) {
    _setData(data);
    onDataChange({ data });
  }

  const [status, setStatus] = useState<IStatus>('idle');

  // refetch
  const [flag, _setFlag] = useState(0);
  function reFetch() {
    _setFlag(pre => pre + 1);
  }

  // 普通翻页相关的
  const [totalCount, setTotalCount] = useState<number>(null);
  const [pageIndex, _setPageIndex] = useState(1);
  const [pageSize, _setPageSize] = useState(defaultPageSize);

  function setPageIndex(_) {
    _setPageIndex(_);

    reFetch();
  }

  function setPageSize(_) {
    _setPageSize(_);

    setPageIndex(1);
  }

  // 定时相关
  const timer = useRef(null);
  useEffect(() => {
    clearInterval(timer.current);

    const _timer = setInterval(() => {
      if (!polling) return;

      if (status === 'loading' || status === 'loading-polling') return;

      if (onlyPollingPage1 && pageIndex !== 1) return;

      fetchData(true);
    }, pollingDelay);

    timer.current = _timer;

    return () => clearInterval(timer.current);
  }, [polling, status, pollingDelay, onlyPollingPage1, pageIndex]);

  // continue分页相关的
  const [continueState, setContinueState] = useState([null]);

  // continue专用翻页
  function nextPageIndex() {
    setPageIndex(pre => pre + 1);
  }

  function prePageIndex() {
    setPageIndex(pre => pre - 1);
  }

  async function fetchData(isPolling = false) {
    if (needClearData) {
      setData(null);
    }

    try {
      setStatus(isPolling ? 'loading-polling' : 'loading');
      const paging = { pageIndex, pageSize };

      switch (mode) {
        case 'normal': {
          const { data } = await query();
          setData(data);

          break;
        }

        case 'paging': {
          const { data, totalCount } = await query({ paging });
          setTotalCount(totalCount);
          setData(data);
          break;
        }

        case 'continue': {
          const { pageIndex, pageSize } = paging;
          const currentContinue = continueState[pageIndex - 1];
          let { data, continueToken, totalCount } = await query({ paging, continueToken: currentContinue });

          // 针对不返回totalcount的情况

          if (totalCount === null) {
            totalCount = (pageIndex + (continueToken ? 1 : 0)) * pageSize;
          }

          setContinueState(pre => {
            const newState = [...pre];

            newState.splice(pageIndex, 1, continueToken);

            return newState;
          });

          setTotalCount(totalCount);

          setData(data);
        }
      }

      setStatus('success');
    } catch (error) {
      setData(null);

      if (error?.response?.status === 410) {
        setStatus('expired');
      } else {
        setStatus('error');
      }
    }
  }

  // deps改变，需要重置pageIndex,  并且重新拉取
  useEffect(() => {
    setPageIndex(1);
  }, deps);

  useEffect(() => {
    if (fetchAble) {
      fetchData();
    }
  }, [flag]);

  return {
    data,

    status,

    reFetch,

    paging: {
      totalCount,

      pageIndex,

      setPageIndex(page: number) {
        if (totalCount && (page - 1) * pageSize < totalCount) {
          setPageIndex(page);
        }
      },

      pageSize,

      setPageSize,

      nextPageIndex,
      prePageIndex
    }
  };
}

/* TODO:
- 边界条件
- 无限分页
- 轮训的时候不想出现loading，但是其他时候需要 - done
 */
