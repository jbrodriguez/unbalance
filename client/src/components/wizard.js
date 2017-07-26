import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import { Link } from 'react-router-dom'
import classNames from 'classnames/bind'
// import { Switch, Route, Link } from 'react-router-dom'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export default class Wizard extends PureComponent {
	static propTypes = {
		// store: PropTypes.object.isRequired,
		match: PropTypes.object.isRequired,
	}

	render() {
		const { match } = this.props

		let prev = null
		let next = null
		let chooseStyle = cx('circular')
		let targetStyle = cx('circular')
		let moveStyle = cx('circular')
		let chooseDisabled = true
		let targetDisabled = true
		let moveDisabled = true
		switch (match.url) {
			case '/gather/target':
				prev = '/gather'
				next = '/gather/move'
				chooseStyle = cx('circular', 'circular-disabled')
				moveStyle = cx('circular', 'circular-disabled')
				targetDisabled = false
				break
			case '/gather/move':
				prev = '/gather/target'
				chooseStyle = cx('circular', 'circular-disabled')
				targetStyle = cx('circular', 'circular-disabled')
				moveDisabled = false
				break
			case '/':
			default:
				next = '/gather/target'
				targetStyle = cx('circular', 'circular-disabled')
				moveStyle = cx('circular', 'circular-disabled')
				chooseDisabled = false
				break
		}

		return (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					{prev &&
						<Link to={prev} className={cx('btn', 'btn-nav', 'rspacer', 'linkBody')}>
							<span className={cx('linkText')}>&lt;</span>
							<span className={cx('linkText')}> PREVIOUS</span>
						</Link>}

					{next &&
						<Link to={next} className={cx('btn', 'btn-nav', 'rspacer', 'linkBody')}>
							<span className={cx('linkText')}>&gt;</span>
							<span className={cx('linkText')}> NEXT</span>
						</Link>}

					<div className={cx('step', 'rspacer', 'linkBody')} disabled={chooseDisabled}>
						<span className={chooseStyle}>1</span> <span> SELECT FOLDER</span>
					</div>
					<div className={cx('step', 'btn-primary', 'rspacer', 'linkBody')} disabled={targetDisabled}>
						<span className={targetStyle}>2</span> <span> CHOOSE TARGET DRIVE</span>
					</div>
					<div className={cx('step', 'btn-primary', 'rspacer', 'linkBody')} disabled={moveDisabled}>
						<span className={moveStyle}>3</span> <span> MOVE</span>
					</div>
				</div>
			</section>
		)
	}
}
