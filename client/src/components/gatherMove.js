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

		if (state.core.status === constant.OP_NEUTRAL) {
			actions.clearConsole()
		}
	}

	render() {
		const { match, store: { state, actions } } = this.props

		const preReqNotPresent = Object.keys(state.gather.chosen).length === 0 || !isValid(state.gather.target)
		const runningMove = state.core.status ? state.core.status === constant.OP_GATHER_MOVE : false

		if (preReqNotPresent && !runningMove) {
			return null
		}

		let consolePanel = null
		if (state.env.lines.length !== 0) {
			consolePanel = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<ConsolePanel lines={state.lines} styleClass={'console-feedback'} />
					</div>
				</section>
			)
		}

		const opInProgress = state.env.isBusy || state.core.status !== constant.OP_NEUTRAL
		const transferDisabled = opInProgress || state.core.operation.bytesToTransfer === 0

		let summary = null
		let proceed = null
		if (!(transferDisabled || opInProgress)) {
			let dst = null
			let size = 0

			state.core.unraid.disks.forEach(disk => {
				if (disk.dst) {
					dst = disk
					size = dst.free - state.core.operation.vdisks[dst.path].plannedFree
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
							onClick={() => actions.gather(state.gather.target)}
							disabled={transferDisabled || opInProgress}
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
