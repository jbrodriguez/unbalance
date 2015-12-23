import React, { Component } from 'react'
import { Link } from 'react-router'
import * as C from '../constant'

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


	componentWillReceiveNextProps(nextProps) {
		let { model } = this.props.model

		if (nextProps.unraid.disks === model.unraid.disks) {
			return
		}

		let toDisk, fromDisk = {}
		let maxFreeSize = 0
		let maxFreePath = ""

		model.unraid.disks.map( disk => {
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

	render() {
		let { dispatch, model } = this.props

		if (!model.unraid) {
			// console.log('about to div')
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

			let rows = 	model.unraid.disks.map( (item, i) => (
							<tr>
								<td>{item.path.replace("/mnt", "")}</td>
								<td>{item.serial} ({item.device})</td>
								<td><input type="checkbox" checked={this.state.fromDisk[disk.path]} onChange="this._checkFrom.bind(this, item.path)"/></td>
								<td><input type="checkbox" checked={this.state.toDisk[disk.path]} onChange="this._checkTo.bind(this, item.path)"/></td>
								<td>{humanBytes(item.size)}</td>
								<td>{humanBytes(item.free)}</td>
								<td>{percentage(item.free/item.size)}</td>
								<td>{humanBytes(item.newFree)}</td>				
							</tr>
						))

			// console.log('rows: ', rows)

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
						<tr>
							<td colSpan="4">TOTAL</td>
							<td>{humanBytes(model.unraid.condition.size)}</td>
							<td>{humanBytes(model.unraid.condition.free)}</td>
							<td>
								<div className={cx('progress')}>
								</div>
							</td>
							<td>{humanBytes(model.unraid.condition.free)}</td>
						</tr>
					</tbody>
				</table>					
			)
		}

		return (
			<div>
				<section className={cx('row', 'bottom-spacer-half')}>
					{ warning }
				</section>

				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12', 'col-sm-8')}>
						STATUS: {model.unraid.condition.state}
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
		let fromDisk = Object.assign({}, this.state.fromDisk)
		for (var key in fromDisk) {
			if (key !== path) {
				fromDisk[key] = false
			}
		}
		fromDisk[path] = true

		let toDisk = Object.assign({}, this.state.toDisk)
		for (var key in toDisk) {
			toDisk[key] = !(key === from)
		}

		this.setState({
			fromDisk,
			toDisk,
		})
	}

	_checkTo(path, e) {
		let toDisk = Object.assign({}, this.state.toDisk)
		toDisk[path] = !toDisk[path]

		this.setState({
			toDisk
		})
	}

}