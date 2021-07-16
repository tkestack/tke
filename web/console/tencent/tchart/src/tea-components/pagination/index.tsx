import * as React from 'react';

import { PageProps } from '../libs/paginationProps';

import { PageSelect } from './PageSelect';
import { PageInput } from './PageInput';

interface PageState {
}

export class Pagination extends React.Component<PageProps, PageState> {
  constructor(props) {
    super(props);
  }

  handlePageSizeChange(pageSize) {
    this.props.onPageSizeChange(pageSize);
  }

  handlePageNoChange(pageNo) {
    const { total, pageSize } = this.props;
    const pageCount = Math.max(Math.ceil(total / pageSize), 1);

    if (pageNo > pageCount) {
      this.props.onChange(pageCount);
    } else if (pageNo < 1) {
      this.props.onChange(1);
    } else {
      this.props.onChange(pageNo);
    }
  }

  render() {
    const {
      total,
      current,
      pageSize,
      pageSizeOptions,
    } = this.props;

    const pageCount = Math.max(Math.ceil(total / pageSize), 1);

    // const pageNoList = [ ...Array(pageCount) ].map((_, index) => index + 1);

    return (
      <div className="tc-15-page">
        <div className="tc-15-page-state">
          <span className="tc-15-page-text">
            共<strong>{total}</strong>项
          </span>
        </div>
        <div className="tc-15-page-operate">
          <span className="tc-15-page-text">每页显示行</span>
          <PageSelect current={pageSize}
            showList={pageSizeOptions}
            onChange={this.handlePageSizeChange.bind(this)}
          />
          <a title="第一页"
            href="javascript:void(0)"
            className={`tc-15-page-first ${pageCount < 2 ? 'disable' : ''}`}
            onClick={() => pageCount >= 2 && this.handlePageNoChange(1)}
          ></a>
          <a title="上一页"
            href="javascript:void(0)"
            className={`tc-15-page-pre ${pageCount < 2 || current === 1 ? 'disable' : ''}`}
            onClick={() => (pageCount >= 2 && current > 1) &&
              this.handlePageNoChange(current - 1)}
          ></a>
          <PageInput
            total={pageCount}
            current={current}
            onChange={this.handlePageNoChange.bind(this)}
          />
          <a title="下一页"
            href="javascript:void(0)"
            className={`tc-15-page-next ${
              pageCount <= current || pageCount === 1 ?
              'disable' : ''
            }`}
            onClick={() => (pageCount > current && pageCount > 1) &&
              this.handlePageNoChange(current + 1)}
          ></a>
          <a title="最后一页"
            href="javascript:void(0)"
            className={`tc-15-page-last ${
              pageCount <= current || pageCount === 1 ?
              'disable' : ''
            }`}
            onClick={() => (pageCount > current && pageCount > 1) &&
              this.handlePageNoChange(pageCount)}
          ></a>
        </div>
    </div>
    );
  }
}
