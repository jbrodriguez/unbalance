import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

// import 'font-awesome-webpack'
import classNames from 'classnames/bind'
import { DateTime } from 'luxon'
import Modal from 'react-modal'

import styles from '../styles/core.scss'

import { formatBytes, percentage } from '../lib/utils'
import * as constant from '../lib/const'

const cx = classNames.bind(styles)

export default class History extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	constructor(props) {
		super(props)

		this.state = {
			showModal: false,
			toConfirm: '',
			id: '',
			operation: null,
		}
	}

	onRequestClose = _ => {
		this.setState({ showModal: false, toConfirm: '', id: '', operation: null })
	}

	confirm = (toConfirm, id, operation) => _ => {
		this.setState({ showModal: true, toConfirm, id, operation })
	}

	onYes = _ => {
		const { actions } = this.props.store

		if (this.state.toConfirm === 'replay') {
			actions.replay(this.state.id)
		} else if (this.state.toConfirm === 'validate') {
			actions.scatterValidate(this.state.id)
		} else if (this.state.toConfirm === 'rmsrc') {
			actions.removeSource(this.state.operation, this.state.id)
		}

		this.setState({ showModal: false, toConfirm: '', id: '', operation: null })
	}

	onNo = _ => {
		this.setState({ showModal: false, toConfirm: '', id: '', operation: null })
	}

	flipOperation = id => e => {
		const { actions } = this.props.store
		actions.flipOperation(id)
		e.preventDefault()
	}

	componentDidMount() {
		const { actions } = this.props.store
		actions.getHistory()
	}

	render() {
		const { state } = this.props.store

		if (!(state.core && state.core.history && state.core.history.order && state.core.history.order.length > 0)) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<span>History is not yet available.</span>
					</div>
				</section>
			)
		}

		// let dryRuns = 0

		const operations = state.core.history.order.map((id, index) => {
			const op = state.core.history.items[id]

			// it's safe to validate or replay an operation only when it's the most recent, excluding dry-runs, since
			// they don't physically alter files
			// if (op.dryRun) dryRuns++
			// let's remove dryRuns condition
			// it can be replayed, validated or have a command source removal only if it's the most recent
			const safe = index === 0

			const replay = !op.dryRun && safe
			const validate = !op.dryRun && op.opKind === constant.OP_SCATTER_COPY && safe

			const flagged = op.commands.some(command => command.status === constant.CMD_FLAGGED)
			const status = flagged ? (
				<i className={cx('fa fa-check-circle', 'statusFlagged')} />
			) : op.bytesTransferred === op.bytesToTransfer ? (
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

			const canBeFlagged =
				safe && (op.opKind === constant.OP_SCATTER_MOVE || op.opKind === constant.OP_GATHER_MOVE)
			const warning =
				canBeFlagged && flagged ? (
					<section className={cx('row', 'mt2')} key={'w' + op.id}>
						<div className={cx('flexSection', 'col-xs-12', 'start-xs', 'middle-xs')}>
							<div>
								<p className={cx('opWhite')}>
									One or more commands had an execution warning/error. Check /boot/logs/unbalance.log
									for additional details.
								</p>
								<p className={cx('opWhite')}>
									Due to this, the plugin hasn't deleted the source files/folders for that/those
									commands.
								</p>
								<p className={cx('opWhite')}>
									Once you've checked/solved the issue(s), click on the{' '}
									<span className={cx('statusFlagged')}>rmsrc</span> button to remove the source
									files/folders, if you wish to do so.
								</p>
							</div>
						</div>
					</section>
				) : null

			let commands
			if (op.open) {
				const rows = op.commands.map(command => {
					let status

					switch (command.status) {
						case constant.CMD_COMPLETE:
							status = <i className={cx('fa fa-check-circle', 'statusDone', 'rspacer')} />
							break

						case constant.CMD_PENDING:
							status = <i className={cx('fa fa-minus-circle', 'statusPending', 'rspacer')} />
							break

						case constant.CMD_SOURCEREMOVAL:
						case constant.CMD_FLAGGED:
							status = canBeFlagged ? (
								<button
									className={cx('btn', 'btn-warning')}
									onClick={this.confirm('rmsrc', command.id, op)}
								>
									rmsrc
								</button>
							) : (
								<i className={cx('fa fa-check-circle', 'statusFlagged', 'rspacer')} />
							)
							break

						case constant.CMD_STOPPED:
							status = <i className={cx('fa fa-times-circle', 'statusInterrupted', 'rspacer')} />
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
								rsync {op.rsyncStrArgs} &quot;{command.entry}&quot; &quot;{command.dst}&quot;
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
					<div>
						{warning}
						<section className={cx('row', 'mt2')} key={'c' + op.id}>
							<div className={cx('flexSection', 'col-xs-12', 'center-xs', 'middle-xs')}>
								<table>
									<thead>
										<tr>
											<th style={{ width: '75px' }} />
											<th style={{ width: '95px' }}>SOURCE</th>
											<th>COMMAND</th>
											<th style={{ width: '325px' }}>PROGRESS</th>
										</tr>
									</thead>
									<tbody>{rows}</tbody>
								</table>
							</div>
						</section>
					</div>
				)
			}

			return (
				<div key={op.id} className={cx('historyItem', 'bottom-spacer-half')}>
					<Modal
						isOpen={this.state.showModal}
						onRequestClose={this.onRequestClose}
						className="modal"
						overlayClassName="overlay"
						ariaHideApp={false}
					>
						<h2>Are you sure ?</h2>
						<button className={cx('btn', 'btn-primary', 'rspacer')} onClick={this.onYes}>
							YES
						</button>
						<button className={cx('btn', 'btn-primary', 'rspacer')} onClick={this.onNo}>
							NO
						</button>
					</Modal>
					<section className={cx('row')}>
						<div
							className={cx('flexSection', 'col-xs-12', 'col-sm-1', 'center-xs', 'middle-xs', 'start-sm')}
						>
							{status}{' '}
							{op.dryRun && <span className={cx('lspacer', 'historyLabel', 'rspacer')}>dry</span>}
						</div>
						<div
							className={cx('flexSection', 'col-xs-12', 'col-sm-3', 'center-xs', 'middle-xs', 'start-sm')}
						>
							<span className={cx('historyTitle', constant.opMap[op.opKind].color)}>
								{constant.opMap[op.opKind].name}
							</span>
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
						<div className={cx('flexSection', 'col-xs-12', 'col-sm-2', 'center-xs', 'middle-xs', 'end-sm')}>
							<span className={cx('historyValue')}>{value}</span>
							<span className={cx('rspacer')}>&nbsp;{unit}</span>
						</div>
						<div className={cx('flexSection', 'col-xs-12', 'col-sm-2', 'center-xs', 'middle-xs', 'end-sm')}>
							{validate && (
								<button
									className={cx('btn', 'btn-primary', 'rspacer')}
									onClick={this.confirm('validate', op.id)}
									disabled={!validate}
								>
									VALIDATE
								</button>
							)}
							{replay && (
								<button
									className={cx('btn', 'btn-primary', 'rspacer')}
									onClick={this.confirm('replay', op.id)}
									disabled={!replay}
								>
									REPLAY
								</button>
							)}
							{chevron}
						</div>
					</section>
					{commands}
				</div>
			)
		})

		return <div>{operations}</div>
	}
}
