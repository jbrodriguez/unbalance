import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'
import { findDOMNode } from 'react-dom'

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

		const node = findDOMNode(this)
		node.scrollTop = node.scrollHeight
	}

	render() {
		const { lines, styleClass } = this.props

		const items = lines.map((line, i) => <p key={i} className={cx('consoleLine')}>{line}</p>)

		return (
			<div className={cx('console', styleClass)}>
				{items}
			</div>
		)
	}
}
