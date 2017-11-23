import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import 'font-awesome-webpack'
import classNames from 'classnames/bind'

import Indicator from './indicator'

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

	// this is a dirty trick to play around the fact that react-router is kind of dumb when it comes to
	// redux like scenarios: https://reacttraining.com/react-router/web/guides/redux-integration
	componentDidMount() {
		const { actions } = this.props.store
		actions.getHistory()
	}

	render() {
		const { state } = this.props.store

		if (!(state.core && state.core.history && state.core.history.length > 0)) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<span>History is not yet available.</span>
					</div>
				</section>
			)
		}

		const opRows = []

		for (let i = state.core.history.length - 1; i >= 0; i--) {
			const op = state.core.history[i]

			let status
			if (op.bytesTransferred === op.bytesToTransfer) {
				status = <i className={cx('fa fa-check-circle', 'statusDone')} />
			} else {
				status = <i className={cx('fa fa-check-circle', 'statusInterrupted')} />
			}

			opRows.push(
				<tr key={op.finished}>
					<td>{status}</td>
					<td>
						{opMap[op.opKind]} {op.dryRun && <span className={cx('label')}>dryRun</span>}
					</td>
					<td>{op.finished}</td>
				</tr>,
			)
		}

		const opList = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<table>
						<thead>
							<tr>
								<th style={{ width: '50px' }} />
								<th>OPERATION</th>
								<th>FINISHED</th>
							</tr>
						</thead>
						<tbody>{opRows}</tbody>
					</table>
				</div>
			</section>
		)

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

		return <div>{opList}</div>
	}
}
