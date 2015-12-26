import React from 'react'
import { findDOMNode } from 'react-dom'
// require('react-dom')

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'
let cx = classNames.bind(styles)

export default class Console extends React.Component {
	componentDidUpdate() {
		// console.log('being called')
		var node = findDOMNode(this)
		node.scrollTop = node.scrollHeight
	}

	render() {
		let { model } = this.props

		if (model.lines.length === 0) {
			return null;
		}

		let items = model.lines.map( (line, i) => {
			return (
				<p key={i} className={cx('consoleLine')}>{line}</p>
			)
		})

		return (
			<div className={cx('col-xs-12')}>
				<div className={cx('console')}>
					{ items }
				</div>
			</div>
		)		
	}	
}
