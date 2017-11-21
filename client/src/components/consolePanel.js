import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export default class Console extends PureComponent {
	static propTypes = {
		lines: PropTypes.arrayOf(PropTypes.any).isRequired,
		styleClass: PropTypes.string.isRequired,
	}

	componentDidUpdate() {
		if (this.props.lines.length === 0) {
			return
		}

		this.node.scrollTop = this.node.scrollHeight
	}

	render() {
		const { lines, styleClass } = this.props

		const items = lines.map((line, i) => (
			<p key={i} className={cx('consoleLine')}>
				{line}
			</p>
		))

		return (
			<div className={cx('console', styleClass)} ref={node => (this.node = node)}>
				{items}
			</div>
		)
	}
}
