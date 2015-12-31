import React, { Component } from 'react'
import { Link } from 'react-router'
import ConsolePanel from './consolePanel'

import { humanBytes, percentage } from '../lib/utils'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Home extends Component {
	componentWillMount() {
		let { store, history } = this.props
		if (store.state.config.folders.length === 0) {
			history.pushState(null, '/settings')
		}
	}

	componentDidMount() {
		// // let { model, history } = this.props
		// // if (!model.config) {
		// // 	history.pushState(null, '/settings')
		// // }
		// // console.log('home.didmount.props: ', this.props)
		// this.props.dispatch(C.GET_STORAGE)
		let { actions, dispatch } = this.props.store
		dispatch(actions.getStorage)
	}

	render() {
		let { state, actions, dispatch } = this.props.store
		// let { dispatch, model } = this.props

		if (!state.unraid) {
			return null
		}

		const ok = state.unraid.condition.state === "STARTED"

		let warning = null
		if (!ok) {
			warning = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>The array is not operational. Please start the array first.</p>
					</div>
				</section>
			)
		}

		// <span style="width: {((state.unraid.condition.size-state.unraid.condition.free) / state.unraid.condition.size )}"></span>

		let grid = null
		if (ok) {
			// console.log('disks: ', state.unraid.disks)

			let rows = 	state.unraid.disks.map( (disk, i) => {
				let diskChanged = cx({
					'label': disk.newFree !== disk.free,
					'label-success': disk.newFree !== disk.free && state.fromDisk[disk.path],
				})

				return (
					<tr key={i}>
						<td>{disk.path.replace("/mnt/", "")}</td>
						<td>{disk.serial} ({disk.device})</td>
						<td><input type="checkbox" checked={state.fromDisk[disk.path]} onChange={this._checkFrom.bind(this, disk.path)} /></td>
						<td><input type="checkbox" checked={state.toDisk[disk.path]} onChange={this._checkTo.bind(this, disk.path)} /></td>
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
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<table>
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
									<th>{humanBytes(state.unraid.condition.size)}</th>
									<th>{humanBytes(state.unraid.condition.free)}</th>
									<th>
										<div className={cx('progress')}>
										</div>
									</th>
									<th>{humanBytes(state.unraid.condition.free)}</th>
								</tr>
							</tfoot>
						</table>					
					</div>
				</section>				
			)
		}

		let consolePanel = null
		if (state.lines.length !== 0) {
			consolePanel = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<ConsolePanel lines={state.lines} />
					</div>
				</section>				
			)
		}

		return (
			<div>
				{ warning }

				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<div className={cx('gridHeader')}>
							<section className={cx('row')}>
								<div className={cx('col-xs-12', 'col-sm-9')}>
									<div className={cx('flexSection', 'middle-xs')}>
										<span className={cx('lspacer')}>STATUS:</span>
										<span className={cx('spacer', 'label', 'label-success')}>{state.unraid.condition.state}</span>
									</div>
								</div>
								<div className={cx('col-xs-12', 'col-sm-3')}>
									<div className={cx('flexSection', 'end-xs')}>
										<button className={cx('btn', 'btn-primary')} onClick={this._calculate.bind(this)} disabled={state.opInProgress}>CALCULATE</button>
										<span>&nbsp; | &nbsp;</span>
										<button className={cx('btn', 'btn-primary')} onClick={this._move.bind(this)} disabled={state.moveDisabled || state.opInProgress}>MOVE</button>
										<span>&nbsp; | &nbsp;</span>
										<div className={cx('flexSection', 'middle-xs', 'rspacer')}> 
											<input type="checkbox" checked={state.config.dryRun} onChange={this._flipDryRun.bind(this)} />
											&nbsp;
											<label>dry run</label>
										</div>
									</div>
								</div>
							</section>
						</div>
					</div>
				</section>

				{ consolePanel }

				{ grid }
			</div>
		)
	}

	_checkFrom(path, e) {
		let { actions, dispatch } = this.props.store
		dispatch(actions.checkFrom, path)
	}

	_checkTo(path, e) {
		let { state, actions, dispatch } = this.props.store
		if (state.fromDisk[path]) {
			e.preventDefault()
			return
		}

		dispatch(actions.checkTo, path)
	}

	_flipDryRun(e) {
		let { actions, dispatch } = this.props.store
		dispatch(actions.toggleDryRun)
	}

	_calculate(e) {
		let { actions, dispatch } = this.props.store
		dispatch(actions.calculate)
	}

	_move(e) {
		let { actions, dispatch } = this.props.store
		dispatch(actions.move)
	}

}