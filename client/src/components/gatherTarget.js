import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import classNames from 'classnames/bind'

import Wizard from './wizard'
import ConsolePanel from './consolePanel'
import styles from '../styles/core.scss'
import { humanBytes, percentage } from '../lib/utils'

require('./tree-view.css')

const cx = classNames.bind(styles)

export default class GatherTarget extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
		match: PropTypes.object.isRequired,
	}

	componentDidMount() {
		const { actions, state } = this.props.store
		actions.findTargets(state.gatherTree.chosen)
	}

	checkTarget = path => e => {
		const { state, actions: { checkTo } } = this.props.store

		if (state.fromDisk[path]) {
			e.preventDefault()
			return
		}

		checkTo(path)
	}

	render() {
		const { match, store: { state } } = this.props

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

		// if free === newFree then this disk isn't elegible as a target
		const elegible = state.unraid.disks.filter(disk => disk.free !== disk.newFree)

		// sort elegible disks by least amount of data transfer
		const targets = elegible.sort((a, b) => {
			const xferA = a.free - a.newFree
			const xferB = b.free - b.newFree
			if (xferA < xferB) return -1
			if (xferA > xferB) return 1
			if (a.id < b.id) return -1
			if (a.id > b.id) return 1
			return 0
		})

		const rows = targets.map(disk => {
			const percent = percentage((disk.size - disk.free) / disk.size)
			const present = state.gatherTree.present.some(presence => presence.id === disk.id)

			return (
				<tr key={disk.id}>
					<td>
						<input type="checkbox" checked={false} onChange={this.checkTarget(disk.path)} />
					</td>
					<td>
						{present && <span>*</span>}
					</td>
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
						{humanBytes(disk.free - disk.newFree)}
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
						<span className={cx('label', 'label-success')}>
							{humanBytes(disk.newFree)}
						</span>
					</td>
				</tr>
			)
		})

		let noun = ''
		let verb = 'is'
		const keys = Object.keys(state.gatherTree.chosen)
		if (keys.length > 1) {
			noun = 's'
			verb = 'are'
		}
		const sources = keys.map(chosen => chosen.slice(10)).join()
		const source = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					Source folder{noun} {verb}: <b>{sources}</b> <br />
					Note: Drives are ordered by the least amount of data transfer. An asterisk next to a disk signals
					that one or more source folders are present
				</div>
			</section>
		)

		const grid = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<table>
						<thead>
							<tr>
								<th style={{ width: '100px' }}>CHOOSE</th>
								<th style={{ width: '50px' }} />
								<th style={{ width: '120px' }}>DISK</th>
								<th style={{ width: '75px' }}>TYPE</th>
								<th>SERIAL</th>
								<th style={{ width: '100px' }}>TRANSFER</th>
								<th style={{ width: '100px' }}>SIZE</th>
								<th style={{ width: '85px' }}>FREE</th>
								<th style={{ width: '100px' }}>FILL</th>
								<th style={{ width: '100px' }}>PLAN</th>
							</tr>
						</thead>
						<tbody>
							{rows}
						</tbody>
					</table>
				</div>
			</section>
		)

		return (
			<div>
				<Wizard match={match} />
				{consolePanel}
				{source}
				{grid}
			</div>
		)
	}
}
