import * as d3 from 'd3';
import * as React from 'react';

import { Bubble } from '@tea/component/bubble';
import { Icon } from '@tea/component/icon';
import { BaseReactProps } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface BarChartBrushProps extends BaseReactProps {
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
  treeDataKeys: string[];
  mapData?: {
    span2Service?: any;
    service2ColorMap?: any;
  };
}

const margin = { top: 20, right: 18, bottom: 20, left: 12 };
const colorMap = {
  error: '#e54545',
  success: '#0abf5b',
  weak: '#888888'
};
const handleStyle = {
  width: 8, // 建议偶数
  height: 30
};

interface BarChartBrushState {
  bars?: any;
  xScale?: any;
  yScale?: any;
  isLoading?: boolean;
  hasError?: boolean;
  error?: any;
  strokeWidth?: number;
  xWidth?: number;
}

export class BarChartBrush extends React.Component<BarChartBrushProps, BarChartBrushState> {
  xAxisRefs = null;

  yAxisRefs = null;

  brushRefs = null;

  brush = null;

  handle = null;

  state = {
    bars: [],
    xScale: null,
    yScale: null,
    isLoading: false,
    hasError: false,
    error: '',
    strokeWidth: 5,
    xWidth: 5
  };

  xAxis = d3.axisBottom().tickFormat(d => `${d} ms`);
  yAxis = d3.axisLeft().tickFormat(d => `${d}`);

  static getDerivedStateFromProps(nextProps, prevState) {
    const { width, height, data, treeDataKeys, mapData } = nextProps;
    if (!data) return null;

    const xMax = d3.max(data, d => d.high);
    const xScale = d3
      .scaleLinear()
      .domain([0, xMax])
      .range([margin.left, width - margin.right]);

    let yWidth = Math.max(Math.min(Math.floor(height / data.length) - 2, 8), 2);

    let xWidth = Math.max(Math.min(Math.floor(width / xMax) - 2, 6), 2);

    const yScale = d3
      .scaleLinear()
      .domain([0, data.length * yWidth])
      .range([height - margin.bottom, margin.top + yWidth]);

    const bars = treeDataKeys
      .map(spanId => data.filter(datum => datum.rootSpanId === spanId))
      .reduce((sum, item) => sum.concat(item), [])
      .map((d, index) => ({
        x: xScale(d.low),
        y: yScale((data.length - index) * yWidth),
        width: xScale(d.high) - xScale(d.low),
        fill: mapData.service2ColorMap[mapData.span2Service[d.name]] || colorMap.weak
      }));

    return { bars, xScale, yScale, strokeWidth: yWidth, xWidth };
  }

  componentDidMount() {
    this.setBrush();
  }

  componentDidUpdate(prevProps, prevState) {
    const { xScale, yScale, bars, strokeWidth } = this.state;
    this.xAxis.scale(xScale);
    d3.select(this.xAxisRefs).call(this.xAxis);
    this.yAxis.scale(yScale);
    d3.select(this.yAxisRefs).call(this.yAxis);
    this.setBrush();
    this.initBrushSelection();
    this.drawBars(bars, strokeWidth);
  }

  drawBars = (bars, strokeWidth = 5) => {
    let rects = '';
    bars.forEach((d, i) => {
      rects += `<rect x="${d.x}" y="${d.y}" width="${d.width}" height="${strokeWidth}" fill="${d.fill}" />`;
    });
    d3.select('.bar-wrapper').html(rects);
  };

  setBrush = () => {
    const oldSvg = d3.select('svg .bar-wrapper').remove();

    const { width, height } = this.props;
    const { xScale, bars } = this.state;
    this.brush = d3
      .brushX()
      .extent([
        [margin.left, margin.top],
        [width - margin.right, height - margin.bottom]
      ])
      .handleSize([2])
      .on('start brush', this.brushMoved)
      .on('end', this.brushEnd);
    let gBrush = d3.select(this.brushRefs).call(this.brush);

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
      .selectAll('.handle--custom')
      .data([{ type: 'w' }, { type: 'e' }])
      .enter()
      .append('rect')
      .attr('class', 'handle--custom')
      .attr('width', handleStyle.width)
      .attr('height', handleStyle.height)
      .attr('fill', '#999')
      .attr('stroke', '#fff')
      .attr('stroke-width', 1)
      .attr('cursor', 'ew-resize');

    d3.select(this.brushRefs)
      .insert('g', '.handle.handle--e')
      .classed('bar-wrapper', true);
  };

  initBrushSelection = () => {
    const { xScale } = this.state;
    const { range } = this.props;
    this.brush.move(d3.select(this.brushRefs), [xScale(range[0]), xScale(range[1])]);
  };

  brushMoved = () => {
    const { xWidth } = this.state;
    let event = d3.event.selection;
    if (!event) {
      d3.selectAll('.handle--custom').attr('display', 'none');
      return;
    }

    const [x1, x2] = event;
    if (x1 !== x2) {
      d3.selectAll('.handle--custom')
        .attr('display', null)
        .attr('x', (d, i) => event[i])
        .attr('transform', (d, i) => {
          return `translate(-${Math.round(handleStyle.width / 2)}, ${margin.top - 1})`;
        });
    } else if (x1 && !isNaN(x1)) {
      this.brush.move(d3.select(this.brushRefs), [Math.max(x1 - xWidth, 0), x1 + xWidth]);
    }
  };

  brushEnd = () => {
    if (!d3.event.sourceEvent) return;
    let event = d3.event.selection;
    if (!event) {
      this.props.updateRange([]);
      return;
    }
    const [x1, x2] = event;
    const { xScale } = this.state;
    const range = [Math.round(xScale.invert(x1)), Math.round(xScale.invert(x2))];
    this.props.updateRange(range);
  };

  render() {
    return this.contentShow();
  }

  private contentShow() {
    const { height, width, isLoading, data, hasError, error } = this.props;

    const hasData = data.length > 0;

    if (isLoading) {
      return (
        <div style={{ width, height, textAlign: 'center', padding: '20px', overflow: 'hidden' }}>
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
              transform={`translate(0, ${height - margin.bottom})`}
              style={{ color: '#ccc' }}
            />
          </g>
          <g
            fill="#f3f3f3"
            ref={r => {
              this.brushRefs = r;
            }}
          />
        </svg>
      </div>
    );
  }
}
