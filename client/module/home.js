import React, { Component } from 'react'
import { Link } from 'react-router'
import * as C from '../constant'
import ConsolePanel from './consolePanel'

import { humanBytes, percentage } from '../lib/utils'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Home extends Component {
	constructor(props, context) {
		super(props, context)

		this.state = {
			toDisk: {},
			fromDisk: {},
			maxFreeSize: 0,
			maxFreePath: "",
		}
	}

	componentDidMount() {
		let { model, history } = this.props
		if (!model.config) {
			history.pushState(null, '/settings')
		}
		// console.log('home.didmount.props: ', this.props)
		this.props.dispatch(C.GET_STORAGE)
	}

					// <div className={cx('col-xs-12', 'col-sm-4')}>
					// 	<button className={cx('btn', 'btn-primary')} onClick={vm.calculateBestFit} data-ng-disabled="vm.disableCalcCtrl">CALCULATE</button>
					// 	<span>&nbsp; | &nbsp;</span>
					// 	<button className={cx('btn', 'btn-primary')} onClick={vm.move}  data-ng-disabled="vm.disableMoveCtrl">MOVE</button>
					// 	<div className={cx('dryrun')}> <input type="checkbox" data-ng-checked="{{vm.options.config.dryRun}}" data-ng-model="vm.options.config.dryRun" data-ng-change="vm.flipDryRun()" data-ng-class="{disabled: vm.disableControls}" /> <span>dry-run </span> </div>
					// </div>


									// <td>{disk.path|replace:"/mnt/":""}}</td>
									// <td>{{disk.serial}} ({{disk.device}})</td>
									// <td><input type="checkbox" data-ng-checked="{{vm.fromDisk[disk.path]}}" data-ng-model="vm.fromDisk[disk.path]" data-ng-change="vm.checkFrom(disk.path)"/> </td>
									// <td><input type="checkbox" data-ng-checked="{{vm.toDisk[disk.path]}}" data-ng-model="vm.toDisk[disk.path]" data-ng-change="vm.checkTo(disk.path)"/> </td>
									// <td>{{disk.size|humanBytes}}</td>
									// <td>{{disk.free|humanBytes}}</td>
									// <td>
									// <div class="progress">
									// <span style="width: {{ ((disk.size - disk.free) / disk.size ) | percentage}}"></span>
									// </div>
									// </td>
									// <td>
									// <span data-ng-class="{label: disk.newFree !== disk.free, 'label-success': disk.newFree !== disk.free}">{{disk.newFree|humanBytes}}</span>
									// </td>				


	componentWillReceiveProps(nextProps) {
		// console.log('home.nextprops: ', nextProps)

		let newModel = nextProps.model
		let oldModel = this.props.model

		if (!newModel.unraid) {
			return
		}

		if (oldModel.unraid && newModel.unraid.disks === oldModel.unraid.disks) {
			return
		}

		let toDisk = {}
		let fromDisk = {}
		let maxFreeSize = 0
		let maxFreePath = ""

		newModel.unraid.disks.map( disk => {
			toDisk[disk.path] = true
			fromDisk[disk.path] = false

			if (disk.free > maxFreeSize) {
				maxFreeSize = disk.free
				maxFreePath = disk.path
			}

			return disk
		})

		if (maxFreePath != "") {
			toDisk[maxFreePath] = false
			fromDisk[maxFreePath] = true
		}

		this.setState({
			toDisk,
			fromDisk,
			maxFreePath,
			maxFreeSize,
		})
	}

					// <tbody>
					// 	{ rows }
					// 	<tr>
					// 		<td colSpan="4">TOTAL</td>
					// 		<td>{humanBytes(model.unraid.condition.size)}</td>
					// 		<td>{humanBytes(model.unraid.condition.free)}</td>
					// 		<td>
					// 			<div className={cx('progress')}>
					// 			</div>
					// 		</td>
					// 		<td>{humanBytes(model.unraid.condition.free)}</td>
					// 	</tr>
					// </tbody>


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
					'label': disk.newFree !== disk.free
				})

				return (
					<tr key={i}>
						<td>{disk.path.replace("/mnt/", "")}</td>
						<td>{disk.serial} ({disk.device})</td>
						<td><input type="checkbox" checked={this.state.fromDisk[disk.path]} onChange={this._checkFrom.bind(this, disk.path)} /></td>
						<td><input type="checkbox" checked={this.state.toDisk[disk.path]} onChange={this._checkTo.bind(this, disk.path)} /></td>
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
							<button className={cx('btn', 'btn-primary')} onClick={this._calculate.bind(this, dispatch)}>CALCULATE</button>
							<span>&nbsp; | &nbsp;</span>
							<button className={cx('btn', 'btn-primary')} disabled={true}>MOVE</button>
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
					<ConsolePanel model={model} />
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
		let fromDisk = Object.assign({}, this.state.fromDisk)
		for (var key in fromDisk) {
			if (key !== path) {
				fromDisk[key] = false
			}
		}
		fromDisk[path] = true

		let toDisk = Object.assign({}, this.state.toDisk)
		for (var key in toDisk) {
			toDisk[key] = !(key === path)
		}

		this.setState({
			fromDisk,
			toDisk,
		})
	}

	_checkTo(path, e) {
		if (this.state.fromDisk[path]) {
			e.preventDefault()
			return
		}

		let toDisk = Object.assign({}, this.state.toDisk)
		toDisk[path] = !toDisk[path]

		this.setState({
			toDisk
		})
	}

	_flipDryRun(e) {

	}

	_calculate(dispatch, e) {
		let srcDisk = ''

		for (var key in this.state.fromDisk) {
			if (this.state.fromDisk[key]) {
				srcDisk = key
				break
			}
		}

		dispatch(C.CALCULATE, {srcDisk, dstDisks: this.state.toDisk})
	}

}