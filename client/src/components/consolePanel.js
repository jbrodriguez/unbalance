import React from 'react'
import { findDOMNode } from 'react-dom'

import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const propTypes = {
	lines: React.PropTypes.arrayOf(React.PropTypes.any).isRequired,
	styleClass: React.PropTypes.string.isRequired,
}

export default class Console extends React.Component {
	componentDidUpdate() {
		if (this.props.lines.length === 0) {
			return
		}

		const node = findDOMNode(this)
		node.scrollTop = node.scrollHeight
	}

	render() {
		const { lines, styleClass } = this.props

		const items = lines.map((line, i) => (
				<p key={i} className={cx('consoleLine')}>{line}</p> // eslint-disable-line
		))

		return (
			<div className={cx('console', styleClass)}>
				{ items }
			</div>
		)
	}
}
Console.propTypes = propTypes
