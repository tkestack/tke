export const formatRequestRequest = (rsp: any) => {
  let dataList = [],
    data = null,
    length = 0,
    isAuthorized = true,
    isLoginedSec = true,
    message = '',
    redirect = '';
  if (rsp.data.code === 0) {
    if (rsp.data.data.list) {
      dataList = rsp.data.data.list;
      length = rsp.data.data.total;
    } else {
      data = rsp.data.data;
    }
  } else if (rsp.data.code === 1900) {
    isAuthorized = false;
    isLoginedSec = false;
    message = rsp.data.message;
    redirect = rsp.data.redirect;
  } else if (rsp.data.code === 1800) {
    isAuthorized = false;
    message = rsp.data.message;
    redirect = rsp.data.redirect;
  }

  return {
    data,
    dataList,
    length,
    isAuthorized,
    isLoginedSec,
    message,
    redirect
  };
};
