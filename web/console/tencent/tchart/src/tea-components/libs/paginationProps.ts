export interface PageProps {
  total: number;
  current: number;
  pageSize: number;
  pageSizeOptions: number[];
  onChange: Function;
  onPageSizeChange: Function;
}
