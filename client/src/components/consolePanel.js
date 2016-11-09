import React from 'react'
import { findDOMNode } from 'react-dom'
// require('react-dom')

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'
let cx = classNames.bind(styles)

export default class Console extends React.Component {
	componentDidUpdate() {
		if (this.props.lines.length === 0) {
			return
		}

		// console.log('being called')
		var node = findDOMNode(this)
		node.scrollTop = node.scrollHeight
	}

	render() {
		let { lines } = this.props

		// if (model.lines.length === 0) {
		// 	return null
		// }

		let items = lines.map( (line, i) => {
			return (
				<p key={i} className={cx('consoleLine')}>{line}</p>
			)
		})

		return (
			<div className={cx('console')}>
				{ items }
			</div>
		)		
	}	
}
