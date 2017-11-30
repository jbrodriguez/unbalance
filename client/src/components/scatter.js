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
			allChecked: false,
		}
	}

	componentDidMount() {
		const { actions } = this.props.store
		actions.getStorage()
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
		const { state: { scatter }, actions: { checkFrom } } = this.props.store

		if (scatter.plan.vdisks[path].src) {
			e.preventDefault()
			return
		}

		checkFrom(path)
	}

	checkTo = path => e => {
		const { state: { scatter }, actions: { checkTo } } = this.props.store

		if (scatter.plan.vdisks[path].src) {
			e.preventDefault()
			return
		}

		checkTo(path)
	}

	checkAll = e => {
		const { checkAll } = this.props.store.actions
		checkAll(e.target.checked)
		this.setState({ allChecked: e.target.checked })
	}

	render() {
		const { state, actions } = this.props.store

		if (!(state.core && state.core.unraid && state.scatter && state.scatter.plan)) {
			return null
		}

		const plan = state.scatter.plan
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
		const planDisabled =
			opInProgress ||
			Object.keys(state.scatter.chosen).length === 0 ||
			!(
				state.core.unraid.disks.some(disk => plan.vdisks[disk.path].src) &&
				state.core.unraid.disks.some(disk => plan.vdisks[disk.path].dst)
			)
		const transferDisabled = opInProgress || plan.bytesToTransfer === 0

		const buttons = (
			<div className={cx('flexSection')}>
				<button
					className={cx('btn', 'btn-primary')}
					onClick={() => actions.scatterPlan()}
					disabled={planDisabled}
				>
					PLAN
				</button>
				<span>&nbsp; | &nbsp;</span>
				<button
					className={cx('btn', 'btn-primary')}
					onClick={() => actions.scatterMove()}
					disabled={transferDisabled}
				>
					MOVE
				</button>
				<button
					className={cx('btn', 'btn-primary', 'lspacer')}
					onClick={() => actions.scatterCopy()}
					disabled={transferDisabled}
				>
					COPY
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
		if (state.scatter.lines.length !== 0) {
			consolePanel = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<ConsolePanel lines={state.scatter.lines} styleClass={'console-feedback'} />
					</div>
				</section>
			)
		}

		const rows = state.core.unraid.disks.map((disk, i) => {
			const vdisks = state.scatter.plan.vdisks

			// console.log(`path(${disk.path})-src(${vdisks[disk.path].src})-dst(${vdisks[disk.path].dst})`)

			const diskChanged = cx({
				label: vdisks[disk.path].plannedFree !== disk.free,
				'label-success': vdisks[disk.path].src && vdisks[disk.path].plannedFree !== disk.free,
			})

			const percent = percentage((disk.size - vdisks[disk.path].plannedFree) / disk.size)

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
							<input
								type="checkbox"
								checked={vdisks[disk.path].src}
								onChange={this.checkFrom(disk.path)}
							/>
						</td>
						<td>
							<input type="checkbox" checked={vdisks[disk.path].dst} onChange={this.checkTo(disk.path)} />
						</td>
						<td>{humanBytes(disk.size)}</td>
						<td>{humanBytes(disk.free)}</td>
						<td>
							<div className={cx('progress')}>
								<span style={{ width: percent }} />
							</div>
						</td>
						<td>
							<span className={diskChanged}>{humanBytes(vdisks[disk.path].plannedFree)}</span>
						</td>
					</tr>,
				]

				// if it's the source disk, let's add a second row, with the
				// tree-menu
				if (vdisks[disk.path].src) {
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
								<th colSpan="4">TOTAL</th>
								<th>
									<input type="checkbox" checked={this.state.allChecked} onChange={this.checkAll} />
								</th>
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
