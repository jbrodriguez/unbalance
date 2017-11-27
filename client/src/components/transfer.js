import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import 'font-awesome-webpack'
import classNames from 'classnames/bind'

import Indicator from './indicator'

import styles from '../styles/core.scss'

import { formatBytes, percentage } from '../lib/utils'
import * as constant from '../lib/const'

const cx = classNames.bind(styles)

export default class Transfers extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	componentDidMount() {
		const { actions } = this.props.store
		actions.getOperation()
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
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<span>No transfer operation is currently on going.</span>
					</div>
				</section>
			)
		}

		const operation = state.core.operation

		const completed = parseFloat(Math.round(operation.completed * 100) / 100).toFixed(2)
		const speed = parseFloat(Math.round(operation.speed * 100) / 100).toFixed(2)

		let bytes = formatBytes(operation.bytesTransferred + operation.deltaTransfer)
		const transferredValue = bytes.value
		const transferredUnit = ' ' + bytes.unit

		bytes = formatBytes(operation.bytesToTransfer)
		const totalValue = bytes.value
		const totalUnit = ' ' + bytes.unit

		const remaining = operation.remaining

		const rows = operation.commands.map(command => {
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
				<tr key={`${command.src}${command.entry}`}>
					<td>{status}</td>
					<td>{command.src}</td>
					<td>
						rsync {operation.rsyncStrArgs} &quot;{command.entry}&quot; &quot;{command.dst}&quot;
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
								<th style={{ width: '350px' }}>PROGRESS</th>
							</tr>
						</thead>
						<tbody>{rows}</tbody>
					</table>
				</div>
			</section>
		)

		return (
			<div>
				<div className={cx('historyItem', 'bottom-spacer-half')}>
					<section className={cx('row')}>
						<div className={cx('col-xs-12', 'col-sm-6', 'center-xs', 'start-sm')}>
							<span className={cx('historyTitle', constant.opMap[operation.opKind].color)}>
								{constant.opMap[operation.opKind].name}
							</span>
						</div>
						<div className={cx('col-xs-12', 'col-sm-6', 'center-xs', 'end-sm')}>
							<span>
								{operation.dryRun && <span className={cx('lspacer', 'historyLabel')}>dry</span>}
							</span>
						</div>
					</section>
				</div>

				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs')}>
						<Indicator label="COMPLETED" value={completed} unit=" %" />
					</div>
					<div className={cx('col-xs')}>
						<Indicator label="SPEED" value={speed} unit=" MB/s" />
					</div>
					<div className={cx('col-xs')}>
						<Indicator label="TRANSFERRED" value={transferredValue} unit={transferredUnit} />
					</div>
					<div className={cx('col-xs')}>
						<Indicator label="TOTAL" value={totalValue} unit={totalUnit} />
					</div>
					<div className={cx('col-xs')}>
						<Indicator label="REMAINING" value={remaining} unit="" />
					</div>
				</section>

				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<span>{operation.line}</span>
					</div>
				</section>

				{grid}
			</div>
		)
	}
}
