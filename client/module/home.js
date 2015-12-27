import React, { Component } from 'react'
import { Link } from 'react-router'
import * as C from '../constant'
import ConsolePanel from './consolePanel'

import { humanBytes, percentage } from '../lib/utils'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Home extends Component {
	componentDidMount() {
		let { model, history } = this.props
		if (!model.config) {
			history.pushState(null, '/settings')
		}
		// console.log('home.didmount.props: ', this.props)
		this.props.dispatch(C.GET_STORAGE)
	}

	render() {
		let { dispatch, model } = this.props

		if (!model.unraid) {
			console.log('about to div')
			return (
				<div></div>
			)
		}

		const ok = model.unraid.condition.state === "STARTED"

		let warning = null
		if (!ok) {
			warning = (
				<div className={cx('col-xs-12')}>
					<p className={cx('bg-warning')}>The array is not operational. Please start the array first.</p>
				</div>
			)
		}

		// <span style="width: {((model.unraid.condition.size-model.unraid.condition.free) / model.unraid.condition.size )}"></span>

		let grid = null
		if (ok) {
			// console.log('disks: ', model.unraid.disks)

			let rows = 	model.unraid.disks.map( (disk, i) => {
				let diskChanged = cx({
					'label': disk.newFree !== disk.free,
					'label-success': disk.newFree !== disk.free && model.fromDisk[disk.path],
				})

				return (
					<tr key={i}>
						<td>{disk.path.replace("/mnt/", "")}</td>
						<td>{disk.serial} ({disk.device})</td>
						<td><input type="checkbox" checked={model.fromDisk[disk.path]} onChange={this._checkFrom.bind(this, disk.path)} /></td>
						<td><input type="checkbox" checked={model.toDisk[disk.path]} onChange={this._checkTo.bind(this, disk.path)} /></td>
						<td>{humanBytes(disk.size)}</td>
						<td>{humanBytes(disk.free)}</td>
						<td>{percentage(disk.free/disk.size)}</td>
						<td>
							<span className={diskChanged}>{humanBytes(disk.newFree)}</span>
						</td>
					</tr>
				)
			})

			grid = (
				<table className={cx('')}>
					<thead>
						<tr>
							<th style={{width: '100px'}}>DISK</th>
							<th>SERIAL</th>
							<th style={{width: '50px'}}>FROM </th>
							<th style={{width: '50px'}}>TO</th>
							<th style={{width: '100px'}}>SIZE</th>
							<th style={{width: '85px'}}>FREE</th>
							<th style={{width: '100px'}}>FILL</th>
							<th style={{width: '100px'}}>PLAN</th>
						</tr>
					</thead>
					<tbody>
						{ rows }
					</tbody>
					<tfoot>
						<tr>
							<th colSpan="4">TOTAL</th>
							<th>{humanBytes(model.unraid.condition.size)}</th>
							<th>{humanBytes(model.unraid.condition.free)}</th>
							<th>
								<div className={cx('progress')}>
								</div>
							</th>
							<th>{humanBytes(model.unraid.condition.free)}</th>
						</tr>
					</tfoot>
				</table>					
			)
		}

		let consolePanel = null
		if (model.lines.length !== 0) {
			consolePanel = (
				<ConsolePanel model={model} />
			)
		}

		return (
			<div>
				<section className={cx('row', 'bottom-spacer-half')}>
					{ warning }
				</section>

				<section className={cx('row', 'between-xs', 'gridHeader', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12', 'col-sm-9')}>
						<div className={cx('flex-section', 'middle-xs')}>
							<span className={cx('lspacer')}>STATUS:</span>
							<span className={cx('spacer', 'label', 'label-success')}>{model.unraid.condition.state}</span>
						</div>
					</div>
					<div className={cx('col-xs-12', 'col-sm-3')}>
						<div className={cx('flexSection', 'end-xs')}>
							<button className={cx('btn', 'btn-primary')} onClick={this._calculate.bind(this)} disabled={model.opInProgress}>CALCULATE</button>
							<span>&nbsp; | &nbsp;</span>
							<button className={cx('btn', 'btn-primary')} onClick={this._move.bind(this)} disabled={model.moveDisabled || model.opInProgress}>MOVE</button>
							<span>&nbsp; | &nbsp;</span>
							<div className={cx('flex', 'middle-xs', 'dryrun', 'rspacer')}> 
								<input type="checkbox" checked={model.config.dryRun} onChange={this._flipDryRun.bind(this)} />
								&nbsp;
								<label>dry run</label>
							</div>
						</div>
					</div>
				</section>

				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						{ consolePanel }
					</div>
				</section>


				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
					{ grid }
					</div>
				</section>
			</div>
		)
	}

	_checkFrom(path, e) {
		this.props.dispatch(C.CHECK_FROM, path)
	}

	_checkTo(path, e) {
		if (this.props.model.fromDisk[path]) {
			e.preventDefault()
			return
		}

		this.props.dispatch(C.CHECK_TO, path)
	}

	_flipDryRun(e) {
		this.props.dispatch(C.TOGGLE_DRY_RUN)
	}

	_calculate(e) {
		this.props.dispatch(C.CALCULATE)
	}

	_move(e) {
		this.props.dispatch(C.MOVE)
	}

}