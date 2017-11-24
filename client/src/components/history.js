import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import 'font-awesome-webpack'
import classNames from 'classnames/bind'
import { DateTime } from 'luxon'

import styles from '../styles/core.scss'

import { formatBytes, percentage } from '../lib/utils'
import * as constant from '../lib/const'

const cx = classNames.bind(styles)

const opMap = {
	[constant.OP_SCATTER_MOVE]: 'SCATTER / MOVE',
	[constant.OP_SCATTER_COPY]: 'SCATTER / COPY',
	[constant.OP_GATHER_MOVE]: 'GATHER / MOVE',
}

export default class History extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	flipOperation = id => e => {
		const { actions } = this.props.store
		actions.flipOperation(id)
		e.preventDefault()
	}

	// this is a dirty trick to play around the fact that react-router is kind of dumb when it comes to
	// redux like scenarios: https://reacttraining.com/react-router/web/guides/redux-integration
	componentDidMount() {
		const { actions } = this.props.store
		actions.getHistory()
	}

	render() {
		const { state } = this.props.store

		if (!(state.core && state.core.historyOrder && state.core.historyOrder.length > 0)) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<span>History is not yet available.</span>
					</div>
				</section>
			)
		}

		const operations = state.core.historyOrder.map(id => {
			const op = state.core.history[id]

			const status =
				op.bytesTransferred === op.bytesToTransfer ? (
					<i className={cx('fa fa-check-circle', 'statusDone')} />
				) : (
					<i className={cx('fa fa-times-circle', 'statusInterrupted')} />
				)

			const { value, unit } = formatBytes(op.bytesTransferred)

			const finished = DateTime.fromISO(op.finished)
			const elapsed = finished.diff(DateTime.fromISO(op.started), ['hours', 'minutes', 'seconds'])

			const chevron = op.open ? (
				<a href="" onClick={this.flipOperation(op.id)}>
					<i className={cx('fa fa-chevron-circle-up', 'statusInProgress', 'chevron')} />
				</a>
			) : (
				<a href="" onClick={this.flipOperation(op.id)}>
					<i className={cx('fa fa-chevron-circle-down', 'statusInProgress', 'chevron')} />
				</a>
			)

			let commands
			if (op.open) {
				const rows = op.commands.map(command => {
					let status

					if (command.transferred === 0) {
						status = <i className={cx('fa fa-minus-circle', 'statusPending', 'rspacer')} />
					} else if (command.transferred === command.size) {
						status = <i className={cx('fa fa-check-circle', 'statusDone', 'rspacer')} />
					} else {
						status = <i className={cx('fa fa-circle-o-notch fa-spin', 'statusInProgress', 'rspacer')} />
					}

					const percent = percentage(command.transferred / command.size)

					return (
						<tr key={command.entry}>
							<td>{status}</td>
							<td>{command.src}</td>
							<td>
								rsync {op.rsyncStrFlags} &quot;{command.entry}&quot; &quot;{command.dst}&quot;
							</td>
							<td>
								<div className={cx('progress')}>
									<span style={{ width: percent }} />
								</div>
							</td>
						</tr>
					)
				})

				commands = (
					<section className={cx('row', 'mt2')}>
						<div className={cx('flexSection', 'col-xs-12', 'center-xs', 'middle-xs')}>
							<table>
								<thead>
									<tr>
										<th style={{ width: '50px' }} />
										<th style={{ width: '95px' }}>SOURCE</th>
										<th>COMMAND</th>
										<th style={{ width: '350px' }}>PROGRESS</th>
									</tr>
								</thead>
								<tbody>{rows}</tbody>
							</table>
						</div>
					</section>
				)
			}

			return (
				<div key={op.id} className={cx('historyItem', 'bottom-spacer-half')}>
					<section className={cx('row')}>
						<div
							className={cx('flexSection', 'col-xs-12', 'col-sm-1', 'center-xs', 'middle-xs', 'start-sm')}
						>
							{status} {op.dryRun && <span className={cx('lspacer', 'historyLabel')}>dry</span>}
						</div>
						<div
							className={cx('flexSection', 'col-xs-12', 'col-sm-3', 'center-xs', 'middle-xs', 'start-sm')}
						>
							<span className={cx('historyTitle')}>{opMap[op.opKind]}</span>
						</div>
						<div className={cx('flexSection', 'col-xs-12', 'col-sm-4', 'center-xs', 'middle-xs', 'end-sm')}>
							<span className={cx('rspacer', 'historyTime')}>
								{finished.toLocaleString(DateTime.DATETIME_MED)}
							</span>{' '}
							|
							<span className={cx('lspacer', 'historyTime')}>
								{elapsed.hours}h, {elapsed.minutes}m, {elapsed.seconds}s
							</span>
						</div>
						<div className={cx('flexSection', 'col-xs-12', 'col-sm-4', 'center-xs', 'middle-xs', 'end-sm')}>
							<span className={cx('historyValue')}>{value}</span>
							<span className={cx('rspacer')}>&nbsp;{unit}</span>
							{chevron}
						</div>
					</section>
					{commands}
				</div>
			)
		})

		// const opRows = state.core.historyOrder.map(id => {
		// 	const op = state.core.history[id]

		// 	const status =
		// 		op.bytesTransferred === op.bytesToTransfer ? (
		// 			<i className={cx('fa fa-check-circle', 'statusDone')} />
		// 		) : (
		// 			<i className={cx('fa fa-times-circle', 'statusInterrupted')} />
		// 		)

		// 	const chevron = op.open ? (
		// 		<a href="" onClick={this.flipOperation(op.id)}>
		// 			<i className={cx('fa fa-chevron-circle-up', 'statusInProgress', 'chevron')} />
		// 		</a>
		// 	) : (
		// 		<a href="" onClick={this.flipOperation(op.id)}>
		// 			<i className={cx('fa fa-chevron-circle-down', 'statusInProgress', 'chevron')} />
		// 		</a>
		// 	)

		// 	const lines = [
		// 		<tr key={op.id}>
		// 			<td>
		// 				{status} {op.dryRun && <span className={cx('label')}>dry</span>}
		// 			</td>
		// 			<td>{opMap[op.opKind]}</td>
		// 			<td>{op.finished}</td>
		// 			<td>{chevron}</td>
		// 		</tr>,
		// 	]

		// 	if (op.open) {
		// 		const commands = op.commands.map(command => {
		// 			const status =
		// 				command.transferred === command.size ? (
		// 					<i className={cx('fa fa-check-circle', 'statusDone', 'rspacer')} />
		// 				) : (
		// 					<i className={cx('fa fa-times-circle', 'statusInterrupted', 'rspacer')} />
		// 				)

		// 			return (
		// 				<tr key={command.entry}>
		// 					<td>{status}</td>
		// 					<td>{command.src}</td>
		// 					<td>
		// 						rsync {op.rsyncStrFlags} &quot;{command.entry}&quot; &quot;{command.dst}&quot;
		// 					</td>
		// 				</tr>
		// 			)
		// 		})

		// 		lines.push(commands)
		// 	}

		// 	return lines
		// })

		// const opList = (
		// 	<section className={cx('row', 'bottom-spacer-half')}>
		// 		<div className={cx('col-xs-12')}>
		// 			<table>
		// 				<thead>
		// 					<tr>
		// 						<th style={{ width: '80px' }} />
		// 						<th>OPERATION</th>
		// 						<th>FINISHED</th>
		// 						<th style={{ width: '80px' }} />
		// 					</tr>
		// 				</thead>
		// 				<tbody>{opRows}</tbody>
		// 			</table>
		// 		</div>
		// 	</section>
		// )

		// const operation = state.core.operation

		// const completed = parseFloat(Math.round(operation.completed * 100) / 100).toFixed(2)
		// const speed = parseFloat(Math.round(operation.speed * 100) / 100).toFixed(2)

		// let bytes = formatBytes(operation.bytesTransferred + operation.deltaTransfer)
		// const transferredValue = bytes.value
		// const transferredUnit = ' ' + bytes.unit

		// bytes = formatBytes(operation.bytesToTransfer)
		// const totalValue = bytes.value
		// const totalUnit = ' ' + bytes.unit

		// const remaining = operation.remaining

		// // console.log(`line(${operation.line})`)

		// const rows = operation.commands.map(command => {
		// 	let status

		// 	// console.log(`line(${operation.line})-commandxfer(${command.transferred})-commandxfer(${command.size})`)

		// 	if (command.transferred === 0) {
		// 		status = <i className={cx('fa fa-minus-circle', 'statusPending', 'rspacer')} />
		// 	} else if (command.transferred === command.size) {
		// 		status = <i className={cx('fa fa-check-circle', 'statusDone', 'rspacer')} />
		// 	} else {
		// 		status = <i className={cx('fa fa-circle-o-notch fa-spin', 'statusInProgress', 'rspacer')} />
		// 	}

		// 	const percent = percentage(command.transferred / command.size)

		// 	return (
		// 		<tr key={command.entry}>
		// 			<td>{status}</td>
		// 			<td>{command.src}</td>
		// 			<td>
		// 				rsync {operation.rsyncStrFlags} &quot;{command.entry}&quot; &quot;{command.dst}&quot;
		// 			</td>
		// 			<td>
		// 				<div className={cx('progress')}>
		// 					<span style={{ width: percent }} />
		// 				</div>
		// 			</td>
		// 		</tr>
		// 	)
		// })

		// const grid = (
		// 	<section className={cx('row', 'bottom-spacer-half')}>
		// 		<div className={cx('col-xs-12')}>
		// 			<table>
		// 				<thead>
		// 					<tr>
		// 						<th style={{ width: '50px' }} />
		// 						<th style={{ width: '95px' }}>SOURCE</th>
		// 						<th>COMMAND</th>
		// 						<th style={{ width: '350px' }}>PROGRESS</th>
		// 					</tr>
		// 				</thead>
		// 				<tbody>{rows}</tbody>
		// 			</table>
		// 		</div>
		// 	</section>
		// )

		return <div>{operations}</div>
	}
}
