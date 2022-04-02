import React, { useEffect, useState } from 'react';

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
}

type IUseFetchQuery<T> = (params?: IQueryParams) => Promise<IQueryResponse<T>>;

type IStatus = 'idle' | 'loading' | 'success' | 'error';

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
    pollingDelay = 5000
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

  // 定时相关
  const [timer, setTimer] = useState(null);
  useEffect(() => {
    clearInterval(timer);

    const _timer = setInterval(() => {
      if (!polling) return;

      if (status === 'loading') return;

      reFetch();
    }, pollingDelay);

    setTimer(_timer);

    return () => clearInterval(timer);
  }, [polling, status, pollingDelay]);

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

    reFetch();
  }

  // continue分页相关的
  const [continueState, setContinueState] = useState([null]);

  // continue专用翻页
  function nextPageIndex() {
    setPageIndex(pre => pre + 1);
  }

  function prePageIndex() {
    setPageIndex(pre => pre - 1);
  }

  async function fetchData() {
    setData(null);

    try {
      setStatus('loading');
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
          const pageIndex = paging.pageIndex;
          const currentContinue = continueState[pageIndex - 1];
          const { data, continueToken, totalCount } = await query({ paging, continueToken: currentContinue });
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
      setStatus('error');
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
 */
