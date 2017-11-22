import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import 'font-awesome-webpack'
import classNames from 'classnames/bind'

import Indicator from './indicator'

import styles from '../styles/core.scss'

import { humanBytes, percentage } from '../lib/utils'
import * as constant from '../lib/const'

const cx = classNames.bind(styles)

export default class Transfers extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	render() {
		const { state } = this.props.store

		if (
			!(
				state.core &&
				state.core.operation &&
				(state.core.operation.opKind === constant.OP_SCATTER_MOVE ||
					state.core.operation.opKind === constant.OP_SCATTER_COPY ||
					state.core.operation.opKind === constant.OP_GATHER_MOVE)
			)
		) {
			return null
		}

		const operation = state.core.operation
		const progress = state.core.progress

		const completed = parseFloat(Math.round(progress.completed * 100) / 100).toFixed(2)
		const speed = parseFloat(Math.round(progress.speed * 100) / 100).toFixed(2)
		const transferred = `${humanBytes(operation.bytesTransferred + progress.deltaTransfer)} / ${humanBytes(
			operation.bytesToTransfer,
		)}`
		const remaining = progress.remaining

		console.log(`line(${progress.line})`)

		const rows = operation.commands.map(command => {
			let status

			// console.log(`line(${progress.line})-commandxfer(${command.transferred})-commandxfer(${command.size})`)

			if (command.transferred === 0) {
				status = <i className={cx('fa fa-minus-circle', 'statusPending', 'rspacer')} />
			} else if (command.transferred === command.size) {
				status = <i className={cx('fa fa-check-circle', 'statusDone', 'rspacer')} />
			} else {
				status = <i className={cx('fa fa-bullseye', 'statusInProgress', 'rspacer')} />
			}

			const percent = percentage(command.transferred / command.size)

			return (
				<tr key={command.entry}>
					<td>{status}</td>
					<td>{command.src}</td>
					<td>
						rsync {operation.rsyncStrFlags} &quot;{command.entry}&quot; &quot;{command.dst}&quot;
					</td>
					<td>
						<div className={cx('progress')}>
							<span style={{ width: percent }} />
						</div>
					</td>
				</tr>
			)
		})

		const grid = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<table>
						<thead>
							<tr>
								<th style={{ width: '50px' }} />
								<th style={{ width: '95px' }}>SOURCE</th>
								<th>COMMAND</th>
								<th>PROGRESS</th>
							</tr>
						</thead>
						<tbody>{rows}</tbody>
					</table>
				</div>
			</section>
		)

		return (
			<div>
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-3')}>
						<Indicator label="Completed" value={completed} unit=" %" />
					</div>
					<div className={cx('col-xs-3')}>
						<Indicator label="Speed" value={speed} unit=" MB/s" />
					</div>
					<div className={cx('col-xs-3')}>
						<Indicator label="Transferred / Total" value={transferred} unit="" />
					</div>
					<div className={cx('col-xs-3')}>
						<Indicator label="Remaining" value={remaining} unit="" />
					</div>
				</section>

				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<span>{progress.line}</span>
					</div>
				</section>

				{grid}
			</div>
		)
	}
}
