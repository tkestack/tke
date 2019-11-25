export const orderBy = (arr: Array<any>, by: string, order?: string) => {
  let aa = arr.sort((a, b) => {
    //默认升序asc
    let temp = order === 'desc' ? b[by] > a[by] : a[by] > b[by];
    return temp ? 1 : -1;
  });
  return aa;
};
