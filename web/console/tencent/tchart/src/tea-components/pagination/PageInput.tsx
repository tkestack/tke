import * as React from 'react';

import { OnOuterClick } from '../libs/decorators/OnOuterClick';

export interface PageInputProps {
	current: string | number;
	onChange: Function;
	total: number | string;
}

interface PageInputState {
	visible: boolean;
	tmpPage: number | string;
}

export class PageInput extends React.Component<PageInputProps, PageInputState> {
	ref;
	constructor(props) {
		super(props);
		this.state = {
			visible: false,
			tmpPage: props.current || 1,
		};
	}

	@OnOuterClick
	func() {
		this.setState({
			visible: false,
		});
	}

	onblur() {
		let tmpPage = parseInt("" + this.state.tmpPage);
		tmpPage = isNaN(tmpPage) ? 1 : tmpPage;
		this.props.onChange(Math.max(tmpPage, 1));
		this.setState({
			visible: false,
			tmpPage: 1,
		});
	}
	onPageChange(e) {
		if (!e.target.value) return this.setState({ tmpPage: '' });
		const res = e.target.value.match(/^\d+$/);
		if (!res) return;
		this.setState({ tmpPage: res[0] });
	}
	onFocus() {
		this.setState({ visible: true, tmpPage: this.props.current })
	}
	onKeyup(e) {
		e.keyCode === 13 && this.ref.blur();
	}
	render() {
		const { current, total } = this.props;
		const { visible, tmpPage } = this.state;
		return (
			<div className="tc-15-page-select">
				<input
					ref={ref => this.ref = ref}
					onFocus={this.onFocus.bind(this)}
					onBlur={this.onblur.bind(this)}
					onChange={this.onPageChange.bind(this)}
					value={visible ? tmpPage : current}
					onKeyUp={this.onKeyup.bind(this)}
					className="tc-15-input-text page-num"
				/>/<span className="tc-15-page-num">{total}</span>
			</div>
		);
	}
}
