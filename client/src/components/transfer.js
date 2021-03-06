import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

// import 'font-awesome-webpack'
import classNames from 'classnames/bind'
import { DateTime } from 'luxon'

import Indicator from './indicator'

import styles from '../styles/core.scss'

import { formatBytes, percentage } from '../lib/utils'
import * as constant from '../lib/const'

const cx = classNames.bind(styles)

export default class Transfers extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	stop = () => _ => {
		this.props.store.actions.stopCommand()
	}

	setRefreshRate = e => {
		const refreshRate = +e.target.value
		const { setRefreshRate } = this.props.store.actions
		setRefreshRate(refreshRate)
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
					state.core.operation.opKind === constant.OP_GATHER_MOVE ||
					state.core.operation.opKind === constant.OP_SCATTER_VALIDATE)
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

		const diff = DateTime.local().diff(DateTime.fromISO(operation.started), ['hours', 'minutes', 'seconds'])
		const elapsed = `${diff.hours > 0 ? diff.hours : ''}${diff.hours > 0 ? 'h' : ''}${
			diff.minutes > 0 ? diff.minutes : ''
		}${diff.minutes > 0 ? 'm' : ''}${Math.round(diff.seconds)}s`
		const remaining = operation.remaining

		const rows = operation.commands.map(command => {
			let status

			switch (command.status) {
				case constant.CMD_COMPLETE:
					status = <i className={cx('fa fa-check-circle', 'statusDone', 'rspacer')} />
					break

				case constant.CMD_PENDING:
					status = <i className={cx('fa fa-minus-circle', 'statusPending', 'rspacer')} />
					break

				case constant.CMD_FLAGGED:
					status = <i className={cx('fa fa-check-circle', 'statusFlagged', 'rspacer')} />
					break

				case constant.CMD_STOPPED:
					status = <i className={cx('fa fa-times-circle', 'statusInterrupted', 'rspacer')} />
					break

				case constant.CMD_SOURCEREMOVAL:
					status = <i className={cx('fa fa-circle-o-notch fa-spin', 'statusFlagged', 'rspacer')} />
					break

				default:
					status = <i className={cx('fa fa-circle-o-notch fa-spin', 'statusInProgress', 'rspacer')} />
			}

			const percent = percentage(command.transferred / command.size)

			return (
				<tr key={command.id}>
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
						<div className={cx('col-xs-12', 'col-sm-5', 'center-xs', 'start-sm')}>
							<span className={cx('historyTitle', constant.opMap[operation.opKind].color)}>
								{constant.opMap[operation.opKind].name}
							</span>
						</div>
						<div className={cx('col-xs-12', 'col-sm-2', 'center-xs', 'middle-sm')}>
							<div className={cx('addon')}>
								<span className={cx('rspacer')}>Refresh</span>
								<select
									className={cx('addon-item')}
									name="rate"
									value={state.config.refreshRate}
									onChange={this.setRefreshRate}
								>
									<option value="250">0.25 sec</option>
									<option value="500">0.5 sec</option>
									<option value="1000">1 sec</option>
									<option value="5000">5 sec</option>
									<option value="15000">15 sec</option>
									<option value="30000">30 sec</option>
								</select>
							</div>
						</div>
						<div className={cx('col-xs-12', 'col-sm-5', 'center-xs', 'end-sm')}>
							{operation.dryRun ? (
								<span className={cx('lspacer', 'historyLabel')}>dry</span>
							) : (
								<button
									className={cx('btn', 'btn-primary', 'lspacer')}
									onClick={this.stop()}
									disabled={false}
								>
									STOP
								</button>
							)}
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
						<Indicator label="ELAPSED" value={elapsed} unit="" />
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
