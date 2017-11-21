import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import TreeMenu from 'react-tree-menu'
import classNames from 'classnames/bind'

import ConsolePanel from './consolePanel'
import { humanBytes, percentage } from '../lib/utils'
import styles from '../styles/core.scss'

import * as constant from '../lib/const'

require('./tree-view.css')

const cx = classNames.bind(styles)

export default class Scatter extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	constructor(props) {
		super(props)

		this.state = {
			fromDisk: {},
			toDisk: {},
		}
	}

	componentDidMount() {
		const { actions } = this.props.store
		actions.getState('scatter')
	}

	onCollapse = node => {
		const { scatterTreeCollapsed } = this.props.store.actions
		scatterTreeCollapsed(node)
	}

	onCheck = node => {
		const { scatterTreeChecked } = this.props.store.actions
		scatterTreeChecked(node)
	}

	checkFrom = path => e => {
		const { state: { core }, actions: { checkFrom, resetOperation } } = this.props.store

		if (this.state.fromDisk[path]) {
			e.preventDefault()
			return
		}

		const fromDisk = { ...this.state.fromDisk }
		const toDisk = { ...this.state.toDisk }

		core.unraid.disks.forEach(disk => {
			fromDisk[disk.path] = disk.path === path
			toDisk[disk.path] = disk.path !== path
		})

		this.setState({ fromDisk, toDisk })

		checkFrom(path)
		resetOperation()
	}

	checkTo = path => e => {
		const { resetOperation } = this.props.store.actions
		const { fromDisk, toDisk } = this.state

		if (fromDisk[path]) {
			e.preventDefault()
			return
		}

		this.setState({
			toDisk: {
				...toDisk,
				[path]: !toDisk[path],
			},
		})

		resetOperation()
	}

	render() {
		const { state, actions } = this.props.store
		const { fromDisk, toDisk } = this.state

		if (!(state.core && state.core.unraid && state.core.operation)) {
			return null
		}

		const stateOk = state.core.unraid.state === 'STARTED'

		if (!stateOk) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>
							The array is not started. Please start the array before performing any operations with
							unBALANCE.
						</p>
					</div>
				</section>
			)
		}

		const opInProgress = state.env.isBusy || state.core.status !== constant.OP_NEUTRAL
		const calcDisabled =
			opInProgress ||
			Object.keys(state.scatter.chosen).length === 0 ||
			!(Object.keys(fromDisk).some(from => fromDisk[from]) && Object.keys(toDisk).some(to => toDisk[to]))
		const transferDisabled = opInProgress || state.core.operation.bytesToTransfer === 0

		const buttons = (
			<div className={cx('flexSection')}>
				<button
					className={cx('btn', 'btn-primary')}
					onClick={() => actions.calculateScatter(fromDisk, toDisk)}
					disabled={calcDisabled}
				>
					CALCULATE
				</button>
				<span>&nbsp; | &nbsp;</span>
				<button className={cx('btn', 'btn-primary')} onClick={() => actions.move()} disabled={transferDisabled}>
					MOVE
				</button>
				<button
					className={cx('btn', 'btn-primary', 'lspacer')}
					onClick={() => actions.copy()}
					disabled={transferDisabled}
				>
					COPY
				</button>
				<button
					className={cx('btn', 'btn-primary', 'lspacer')}
					onClick={() => actions.validate()}
					disabled={transferDisabled}
				>
					VALIDATE
				</button>
				<span>&nbsp; | &nbsp;</span>
				<div className={cx('flexSection', 'middle-xs', 'rspacer')}>
					<input
						id="dryRun"
						type="checkbox"
						checked={state.config.dryRun}
						onChange={() => actions.toggleDryRun()}
						disabled={transferDisabled}
					/>
					&nbsp;
					<label htmlFor="dryRun">dry run</label>
				</div>
			</div>
		)

		const menu = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<div className={cx('gridheader')}>
						<section className={cx('row', 'between-xs', 'middle-xs')}>
							<div className={cx('col-xs-12')}>{buttons}</div>
						</section>
					</div>
				</div>
			</section>
		)

		let consolePanel = null
		if (state.env.lines.length !== 0) {
			consolePanel = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<ConsolePanel lines={state.env.lines} styleClass={'console-feedback'} />
					</div>
				</section>
			)
		}

		const rows = state.core.unraid.disks.map((disk, i) => {
			const operation = state.core.operation

			const plannedFree = operation.vdisks[disk.path].plannedFree

			const diskChanged = cx({
				label: plannedFree !== disk.free,
				'label-success': plannedFree !== disk.free && fromDisk[disk.path],
			})

			const percent = percentage((disk.size - disk.free) / disk.size)

			const fromChecked = fromDisk.hasOwnProperty(disk.path)
				? fromDisk[disk.path]
				: operation.vdisks[disk.path].src

			const toChecked = toDisk.hasOwnProperty(disk.path) ? toDisk[disk.path] : operation.vdisks[disk.path].dst

			// console.log("disk.name.length: ", disk.name.length)

			// let serial = scramble(disk.serial)
			if (disk.type === 'Cache' && disk.name.length > 5) {
				return (
					<tr key={disk.id}>
						<td>{disk.name}</td>
						<td>{disk.fsType}</td>
						<td colSpan="7">
							{disk.serial} ({disk.device})
						</td>
					</tr>
				)
			} else {
				// lines initially contains the disk row, which includes the
				// checkbox indicating it's either the from disk or a to disk
				const lines = [
					<tr key={disk.id}>
						<td>{disk.name}</td>
						<td>{disk.fsType}</td>
						<td>
							{disk.serial} ({disk.device})
						</td>
						<td>
							<input type="checkbox" checked={fromChecked} onChange={this.checkFrom(disk.path)} />
						</td>
						<td>
							<input type="checkbox" checked={toChecked} onChange={this.checkTo(disk.path)} />
						</td>
						<td>{humanBytes(disk.size)}</td>
						<td>{humanBytes(disk.free)}</td>
						<td>
							<div className={cx('progress')}>
								<span style={{ width: percent }} />
							</div>
						</td>
						<td>
							<span className={diskChanged}>{humanBytes(plannedFree)}</span>
						</td>
					</tr>,
				]

				// if it's the source disk, let's add a second row, with the
				// tree-menu
				if (fromDisk[disk.path]) {
					const key = i + 100
					lines.push(
						<tr key={key}>
							<td colSpan="3">
								<b>Select folders/files to move</b>
								<br />
								<TreeMenu
									expandIconClass="fa fa-chevron-right"
									collapseIconClass="fa fa-chevron-down"
									onTreeNodeCollapseChange={this.onCollapse}
									onTreeNodeCheckChange={this.onCheck}
									collapsible
									collapsed={false}
									data={state.scatter.items}
								/>
							</td>
							<td colSpan="6" className={cx('topAlign')}>
								<b>Currently selected</b>
								<br />
								<ul>
									{Object.keys(state.scatter.chosen).map(chosen => <li key={chosen}>- {chosen}</li>)}
								</ul>
							</td>
						</tr>,
					)
				}

				return lines
			}
		})

		const grid = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<table>
						<thead>
							<tr>
								<th style={{ width: '70px' }}>DISK</th>
								<th style={{ width: '75px' }}>TYPE</th>
								<th>SERIAL</th>
								<th style={{ width: '50px' }}>FROM </th>
								<th style={{ width: '50px' }}>TO</th>
								<th style={{ width: '100px' }}>SIZE</th>
								<th style={{ width: '85px' }}>FREE</th>
								<th style={{ width: '100px' }}>FILL</th>
								<th style={{ width: '100px' }}>PLAN</th>
							</tr>
						</thead>
						<tbody>{rows}</tbody>
						<tfoot>
							<tr>
								<th colSpan="5">TOTAL</th>
								<th>{humanBytes(state.core.unraid.size)}</th>
								<th>{humanBytes(state.core.unraid.free)}</th>
								<th>
									<div className={cx('progress')} />
								</th>
								<th>{humanBytes(state.core.unraid.free)}</th>
							</tr>
						</tfoot>
					</table>
				</div>
			</section>
		)

		return (
			<div>
				{menu}
				{consolePanel}
				{grid}
			</div>
		)
	}
}
