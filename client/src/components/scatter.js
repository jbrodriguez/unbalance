import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import TreeMenu from 'react-tree-menu'
import classNames from 'classnames/bind'

import ConsolePanel from './consolePanel'
import { humanBytes, percentage } from '../lib/utils'
import styles from '../styles/core.scss'

require('./tree-view.css')

const cx = classNames.bind(styles)

export default class Scatter extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	onCollapse = node => {
		// console.log(`collapse-node-${JSON.stringify(node)}`)
		const { treeCollapsed } = this.props.store.actions
		treeCollapsed(node)
	}

	onCheck = node => {
		// console.log(`check-node-${JSON.stringify(node)}`)
		const { treeChecked } = this.props.store.actions
		treeChecked(node)
	}

	checkFrom = path => () => {
		const { checkFrom } = this.props.store.actions
		checkFrom(path)
	}

	checkTo = path => e => {
		const { state, actions: { checkTo } } = this.props.store

		if (state.fromDisk[path]) {
			e.preventDefault()
			return
		}

		checkTo(path)
	}

	render() {
		const { state, actions } = this.props.store

		if (!state.unraid) {
			return null
		}

		if (state.unraid.condition.state !== 'STARTED') {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>
							The array is not started. Please start the array before perfoming any operations with
							unBALANCE.
						</p>
					</div>
				</section>
			)
		}

		const stateOk = state.unraid && state.unraid.condition.state === 'STARTED'
		const disabled = state.opInProgress || !stateOk || Object.keys(state.tree.chosen).length === 0

		const buttons = (
			<div className={cx('flexSection')}>
				<button className={cx('btn', 'btn-primary')} onClick={() => actions.calculate()} disabled={disabled}>
					CALCULATE
				</button>
				<span>&nbsp; | &nbsp;</span>
				<button
					className={cx('btn', 'btn-primary')}
					onClick={() => actions.move()}
					disabled={state.transferDisabled || state.opInProgress}
				>
					MOVE
				</button>
				<button
					className={cx('btn', 'btn-primary', 'lspacer')}
					onClick={() => actions.copy()}
					disabled={state.transferDisabled || state.opInProgress}
				>
					COPY
				</button>
				<button
					className={cx('btn', 'btn-primary', 'lspacer')}
					onClick={() => actions.validate()}
					disabled={state.validateDisabled || state.opInProgress}
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
						disabled={state.transferDisabled || state.opInProgress}
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
							<div className={cx('col-xs-12')}>
								{buttons}
							</div>
						</section>
					</div>
				</div>
			</section>
		)

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

		// {percentage(disk.free/disk.size)}

		const rows = state.unraid.disks.map((disk, i) => {
			const diskChanged = cx({
				label: disk.newFree !== disk.free,
				'label-success': disk.newFree !== disk.free && state.fromDisk[disk.path],
			})

			const percent = percentage((disk.size - disk.free) / disk.size)

			// console.log("disk.name.length: ", disk.name.length)

			// let serial = scramble(disk.serial)
			if (disk.type === 'Cache' && disk.name.length > 5) {
				return (
					<tr key={disk.id}>
						<td>
							{disk.name}
						</td>
						<td>
							{disk.fsType}
						</td>
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
						<td>
							{disk.name}
						</td>
						<td>
							{disk.fsType}
						</td>
						<td>
							{disk.serial} ({disk.device})
						</td>
						<td>
							<input
								type="checkbox"
								checked={state.fromDisk[disk.path]}
								onChange={this.checkFrom(disk.path)}
							/>
						</td>
						<td>
							<input
								type="checkbox"
								checked={state.toDisk[disk.path]}
								onChange={this.checkTo(disk.path)}
							/>
						</td>
						<td>
							{humanBytes(disk.size)}
						</td>
						<td>
							{humanBytes(disk.free)}
						</td>
						<td>
							<div className={cx('progress')}>
								<span style={{ width: percent }} />
							</div>
						</td>
						<td>
							<span className={diskChanged}>
								{humanBytes(disk.newFree)}
							</span>
						</td>
					</tr>,
				]

				// if it's the source disk, let's add a second row, with the
				// tree-menu
				if (state.fromDisk[disk.path]) {
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
									data={state.tree.items}
								/>
							</td>
							<td colSpan="6" className={cx('topAlign')}>
								<b>Currently selected</b>
								<br />
								<ul>
									{Object.keys(state.tree.chosen).map(chosen =>
										<li key={chosen}>
											- {chosen}
										</li>,
									)}
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
						<tbody>
							{rows}
						</tbody>
						<tfoot>
							<tr>
								<th colSpan="5">TOTAL</th>
								<th>
									{humanBytes(state.unraid.condition.size)}
								</th>
								<th>
									{humanBytes(state.unraid.condition.free)}
								</th>
								<th>
									<div className={cx('progress')} />
								</th>
								<th>
									{humanBytes(state.unraid.condition.free)}
								</th>
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
