import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import classNames from 'classnames/bind'

import Wizard from './wizard'
import ConsolePanel from './consolePanel'
import styles from '../styles/core.scss'
import { isValid, humanBytes } from '../lib/utils'
import * as constant from '../lib/const'

const cx = classNames.bind(styles)

export default class GatherMove extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
		match: PropTypes.object.isRequired,
	}

	componentDidMount() {
		const { state, actions } = this.props.store

		if (state.opInProgress === null) {
			actions.clearConsole()
		}
	}

	render() {
		const { match, store: { state, actions } } = this.props

		const preReqNotPresent = Object.keys(state.gatherTree.chosen).length === 0 || !isValid(state.gatherTree.target)
		const runningMove = state.status ? state.status === constant.stateGather : false

		if (preReqNotPresent && !runningMove) {
			return null
		}

		let consolePanel = null
		if (state.lines.length !== 0) {
			consolePanel = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<ConsolePanel lines={state.lines} styleClass={'console-feedback'} />
					</div>
				</section>
			)
		}

		let summary = null
		let proceed = null
		if (!(state.transferDisabled || state.opInProgress)) {
			let dst = null
			let size = 0

			state.unraid.disks.forEach(disk => {
				if (disk.dst) {
					dst = disk
					size = dst.free - dst.newFree
				}
			})

			if (dst) {
				summary = (
					<section className={cx('row', 'bottom-spacer-half')}>
						<div className={cx('col-xs-12')}>
							<b>{humanBytes(size)}</b> will be transferred to {dst.path}.
						</div>
					</section>
				)
			}

			proceed = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<button
							className={cx('btn', 'btn-primary')}
							onClick={() => actions.gather(state.gatherTree.target)}
							disabled={state.transferDisabled || state.opInProgress}
						>
							PROCEED
						</button>
					</div>
				</section>
			)
		}

		return (
			<div>
				<Wizard match={match} store={this.props.store} />
				{consolePanel}
				{summary}
				{proceed}
			</div>
		)
	}
}
