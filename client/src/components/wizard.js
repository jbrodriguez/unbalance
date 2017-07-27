import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import classNames from 'classnames/bind'

import { NextButton, PrevButton } from './buttons'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export default class Wizard extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
		match: PropTypes.object.isRequired,
	}

	render() {
		const { match, store: { state } } = this.props

		let prev = null
		let next = null

		let prevDisabled = false
		let nextDisabled = false

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
				nextDisabled = state.gatherTree.target === null
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
				nextDisabled = Object.keys(state.gatherTree.chosen).length === 0
				targetStyle = cx('circular', 'circular-disabled')
				moveStyle = cx('circular', 'circular-disabled')
				chooseDisabled = false
				break
		}

		return (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					{prev && <PrevButton to={prev} disabled={prevDisabled} />}
					{next && <NextButton to={next} disabled={nextDisabled} />}

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
