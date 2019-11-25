import { TableColumn } from '@tencent/qcloud-component';

const execColumnWidth = (columns: TableColumn<any>[], hasChecker: boolean = true) => {
  let total = hasChecker ? 95 : 100;
  columns.forEach(column => {
    column.width = total / columns.length + '%';
  });
  return columns;
};

const buttonStyle = {
  marginRight: 5
};
export { execColumnWidth, buttonStyle };
