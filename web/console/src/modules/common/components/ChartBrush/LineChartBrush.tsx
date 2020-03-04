import * as d3 from 'd3';
import * as React from 'react';

import { Bubble } from '@tea/component/bubble';
import { Icon } from '@tea/component/icon';
import { BaseReactProps } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { TransformField } from './bytesTo';

export interface LineChartBrushProps extends BaseReactProps {
  /** 高度 */
  height: number;
  /** 宽度 */
  width: number;
  /** 数据 */
  data: any[];
  /** 选择范围 */
  range?: any;
  /** 更新选择反问 */
  updateRange?: any;
  isLoading?: boolean;
  hasError?: boolean;
  error?: any;
}

const margin = { top: 20, right: 15, bottom: 20, left: 40 };
const blue = '#026eff';
const handleStyle = {
  width: 8, // 建议偶数
  height: 30
};

interface LineChartBrushState {
  values?: any;
  xScale?: any;
  yScale?: any;
  xWidth?: number;
  isLoading?: boolean;
  hasError?: boolean;
  error?: any;
}

export class LineChartBrush extends React.Component<LineChartBrushProps, LineChartBrushState> {
  xAxisRefs = null;

  yAxisRefs = null;

  brushLineRefs = null;

  brush = null;

  handle = null;

  state = {
    values: null,
    xScale: null,
    yScale: null,
    xWidth: 6,
    isLoading: false,
    hasError: false,
    error: ''
  };

  xAxis = d3
    .axisBottom()
    .scale(this.state.xScale)
    .tickFormat(d => `${d} ms`);
  yAxis = d3
    .axisLeft()
    .scale(this.state.yScale)
    .tickFormat(d => TransformField(+d, 1000));

  static getDerivedStateFromProps(nextProps, prevState) {
    const { width, height, data, range } = nextProps;
    if (!data) return null;

    const [min, max] = d3.extent(data, d => d.x);
    const xScale = d3
      .scaleLinear()
      .range([margin.left, width - margin.right])
      .domain([0, max]);

    const vMax = d3.max(data, d => d.value);
    const yScale = d3
      .scaleLinear()
      .range([height - margin.bottom, margin.top])
      .domain([0, vMax * 2]);

    const values = data.map(d => {
      const y0 = yScale(0);
      const y1 = yScale(d.value);
      const y2 = yScale(vMax * 2);
      return {
        x: xScale(d.x) - 3,
        y: y1,
        height: y0 - y1,
        fill: '#006eff'
      };
    });
    let xWidth = Math.max(Math.min(Math.floor(width / max) - 2, 6), 2);

    return { values, xScale, yScale, xWidth };
  }

  componentDidMount() {
    this.setBrush();
  }

  componentDidUpdate(prevProps, prevState) {
    const { xScale, yScale, values, xWidth } = this.state;
    this.xAxis.scale(xScale);
    d3.select(this.xAxisRefs).call(this.xAxis);
    this.yAxis.scale(yScale);
    d3.select(this.yAxisRefs).call(this.yAxis);
    this.setBrush();
    this.initBrushSelection();
    this.drawPath(values, xWidth);
  }

  drawPath = (values, xWidth = 6) => {
    let rects = '';
    values.forEach((d, i) => {
      rects += `<rect x="${d.x}" y="${d.y}" width="${xWidth}" height="${d.height}" fill="${d.fill}" />`;
    });
    d3.select('.line-wrapper').html(rects);
  };

  setBrush = () => {
    // Cleanup old draw
    const oldSvg = d3.select('svg .line-wrapper').remove();

    const { width, height } = this.props;
    this.brush = d3
      .brushX()
      .extent([
        [margin.left, margin.top],
        [width - margin.right, height - margin.bottom]
      ])
      .handleSize([2])
      .on('start brush', this.brushMoved)
      .on('end', this.brushEnd);
    let gBrush = d3.select(this.brushLineRefs).call(this.brush);

    gBrush
      .select('.overlay')
      .attr('fill', '#eee')
      .attr('cursor', 'default');
    gBrush.selectAll('.handle').attr('fill', '#999');

    gBrush
      .select('.selection')
      .attr('fill', '#fff')
      .attr('fill-opacity', 1)
      .attr('stroke-width', 0)
      .attr('cursor', 'default');

    this.handle = gBrush
      .selectAll('.handle--line-custom')
      .data([{ type: 'w' }, { type: 'e' }])
      .enter()
      .append('rect')
      .attr('class', 'handle--line-custom')
      .attr('width', handleStyle.width)
      .attr('height', handleStyle.height)
      .attr('fill', '#999')
      .attr('stroke', '#fff')
      .attr('stroke-width', 1)
      .attr('cursor', 'ew-resize');

    d3.select(this.brushLineRefs)
      .insert('g', '.handle.handle--e')
      .classed('line-wrapper', true);
  };

  brushMoved = () => {
    const { xWidth } = this.state;
    let event = d3.event.selection;
    if (!event) {
      d3.selectAll('.handle--line-custom').style('display', 'none');
      return;
    }

    const [x1, x2] = event;
    if (x2 !== x1) {
      d3.selectAll('.handle--line-custom')
        .style('display', null)
        .attr('x', (d, i) => event[i])
        .attr('transform', (d, i) => {
          return `translate(-${Math.round(handleStyle.width / 2)}, ${margin.top - 1})`;
        });
    } else if (x1 && !isNaN(x1)) {
      this.brush.move(d3.select(this.brushLineRefs), [Math.max(x1 - xWidth, 0), x1 + xWidth]);
    }
  };

  brushEnd = () => {
    if (!d3.event.sourceEvent) return;
    let event = d3.event.selection;
    if (!event) {
      return;
    }

    const [x1, x2] = event;
    const { xScale } = this.state;
    const range = [Math.round(xScale.invert(x1)), Math.round(xScale.invert(x2))];
    this.props.updateRange(range);
  };

  initBrushSelection = () => {
    const { xScale } = this.state;
    const { range } = this.props;
    if (!(range[0] === 0 && range[1] === 0)) {
      this.brush.move(d3.select(this.brushLineRefs), [xScale(range[0]), xScale(range[1])]);
    }
  };

  render() {
    return this.contentShow();
  }

  private contentShow() {
    const { height, width, isLoading, data, hasError, error } = this.props;

    const hasData = data.length > 0;

    // no need
    if (isLoading) {
      return (
        <div style={{ textAlign: 'center', padding: '20px', overflow: 'hidden' }}>
          <Icon type="loading" />
          &nbsp; <Trans>加载中...</Trans>
        </div>
      );
    }
    if (!hasData || hasError) {
      return (
        <div style={{ height, textAlign: 'center', overflow: 'hidden' }}>
          <p style={{ color: '#BBBBBB', fontSize: '14px' }}>
            {hasError && (
              <Bubble content={error ? error.message : t('具体原因不详，请稍后重试')}>
                <Icon type="error" />
              </Bubble>
            )}
            {t('暂无数据')}
          </p>
        </div>
      );
    }
    return (
      <div style={{ width, height, overflow: 'hidden' }}>
        <svg width={width} height={height}>
          <g>
            <g
              ref={r => {
                this.xAxisRefs = r;
              }}
              transform={`translate(-1, ${height - margin.bottom})`}
              style={{ color: '#ccc' }}
            />
            <g
              ref={r => {
                this.yAxisRefs = r;
              }}
              transform={`translate(${margin.left - 1}, 0)`}
              style={{ color: '#ccc' }}
            />
            <g
              ref={r => {
                this.brushLineRefs = r;
              }}
            />
          </g>
        </svg>
      </div>
    );
  }
}
