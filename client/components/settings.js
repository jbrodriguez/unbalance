import React, { Component } from 'react'
import { Link } from 'react-router'
import 'font-awesome-webpack'

import TreePanel from './treePanel'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Settings extends Component {
	// componentDidMount() {
	// 	let { actions, dispatch } = this.props.store
	// 	dispatch(actions.getConfig)
	// }

	render() {
		// let { dispatch, state } = this.props
		// console.log('settings.render: ', this.props.store)
		let { state, actions } = this.props.store

		if (!state.config) {
			return null
		}

		// console.log('state.unraid: ', state.unraid)
		if (!state.unraid) {
			// dispatch(actions.getStorage)
			return null
		}

		const stateOk = state.unraid && state.unraid.condition.state === "STARTED"
		if (!stateOk) {
			console.log('stateOk: ', stateOk)
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>&nbsp; The array is not started. Please start the array before perfoming any operations with unBALANCE.</p>
					</div>
				</section>
			)
		}		

		if (state.opInProgress === actions.calculate || state.opInProgress === actions.move) {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>&nbsp; {state.opInProgress} operation is currently under way. Wait until the operation has finished to make any settings changes.</p>
					</div>
				</section>
			)
		}

		return (
			<div>

			<section className={cx('row', 'bottom-spacer-large')}>
				<div className={cx('col-xs-12')}>
					<div>
						<h3>SET UP NOTIFICATIONS</h3>

						<p>Notifications rely on unRAID's notifications settings, so you need to set up unRAID first, in order to receive notifications from unBALANCE.</p>

						<span> Calculate: </span>
						<input id="calc0" className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 0} onChange={this._setNotifyCalc.bind(this, 0)} />
						<label id="calc0" >No Notifications</label>
						
						<input id="calc1" className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 1} onChange={this._setNotifyCalc.bind(this, 1)} />
						<label id="calc1" >Basic</label>
						
						<input id="calc2" className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 2} onChange={this._setNotifyCalc.bind(this, 2)} />
						<label id="calc2" >Detailed</label>

						<br />

						<span> Move: </span>
						<input id="move0" className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 0} onChange={this._setNotifyMove.bind(this, 0)} />
						<label id="move0">No Notifications</label>
						
						<input id="move0" className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 1} onChange={this._setNotifyMove.bind(this, 1)} />
						<label id="move0">Basic</label>
						
						<input id="move0" className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 2} onChange={this._setNotifyMove.bind(this, 2)} />
						<label id="move0">Detailed</label>
					</div>
				</div>
			</section>

			<section className={cx('row', 'bottom-spacer-large')}>
				<div className={cx('col-xs-12')}>
					<div>
						<h3>RESERVED SPACE</h3>

						<p>Set up the minimum amount of space that should be left free on each disk, after moving folders.</p>
						<p>When calculating how to fill a target disk, unBALANCE will not leave less free space than the threshold you set here.</p>
						<p>This threshold cannot represent less than 450Mb (hard limit set by this app).</p>

						<div className={cx('row')}>
							<div className={cx('col-xs-2')}>
								<div className={cx('addon')}>
									<input className={cx('addon-field')} type="text" value={state.config.reservedAmount} onChange={this._setReservedAmount.bind(this)} />
									<select className={cx('addon-item')} name="unit" value={state.config.reservedUnit} onChange={this._setReservedUnit.bind(this)}>
										<option value="%">%</option> 
										<option value="Mb">Mb</option> 
										<option value="Gb">Gb</option> 
									</select>
								</div>
							</div>
						</div>
					</div>
				</div>
			</section>				

			<section className={cx('row', 'bottom-spacer-large')}>
				<div className={cx('col-xs-12')}>
					<div>
						<h3>WHICH FOLDERS DO YOU WANT TO MOVE ?</h3>

						<p>Define which folders should be moved to free up space on the source disk (you choose the source disk in the main page).</p>
						<p>You can choose entire user shares (e.g.: /Movies) or any folders below a user share (e.g.: /Movies/Action, /Movies/Comedy/90s).</p>
						<p>The folders you select will be moved to other disks in the array, as long as there's enough space available.</p>
						<p>Click on the <button className={cx('btn', 'btn-alert')}>add</button>  button that appears when you hover your mouse over a folder in the "unRAID Shares Explorer" column below, to select it for moving.</p>
						<p>Click on the <i className={cx('fa fa-remove')}></i> icon that appears next to any folder in the "Folders to be moved" column, to deselect it.</p>
					</div>
				</div>
			</section>

			<section className={cx('row')}>
				<div className={cx('col-xs-12')}>
					<div>
						<section className={cx('row')}>
							<div className={cx('col-xs-12')}>
								<div className={cx('explorerHeader')}>
									<section className={cx('row')}>
										<div className={cx('col-xs-12', 'col-sm-8')}>
											<span className={cx('lspacer')}>unRAID Shares Explorer</span>
										</div>
										<div className={cx('col-xs-12', 'col-sm-4')}>
											Folders to be moved
										</div>
									</section>
								</div>
							</div>
						</section>

						<section className={cx('row')}>
							<div className={cx('col-xs-12')}>
								<div className={cx('explorerContent')}>
									<section className={cx('row')}>
										<div className={cx('col-xs-12', 'col-sm-8')}>
											<TreePanel tree={state.tree} {...actions} />
										</div>
										<div className={cx('col-xs-12', 'col-sm-4', 'flex', 'flexOne')}>
											<div className={cx('explorerChosen')}>
												<table className={cx('')}>
													<tbody>
														{ 
															state.config.folders.map( (item, i) => {
																return (
																	<tr key={i}>
																		<td width="40"><i className={cx('fa fa-remove')} onClick={this._deleteFolder.bind(this, item)}></i></td>
																		<td>{item}</td>
																	</tr>
																)
															})
														}
													</tbody>
												</table>
											</div>
										</div>
									</section>
								</div>
							</div>
						</section>
					</div>
				</div>

			</section>




			</div>
		)
	}

	// _addFolder(dispatch, e) {
	// 	console.log('key - value: ', e.key, e.target.value)
	// 	if (e.key !== "Enter") {
	// 		return
	// 	}

	// 	e.preventDefault()

	// 	dispatch(C.ADD_FOLDER, e.target.value)
	// }

	_deleteFolder(folder, e) {
		const { deleteFolder } = this.props.store.actions
		deleteFolder(folder)
	}

	_setNotifyCalc(notify, e) {
		const { setNotifyCalc } = this.props.store.actions
		setNotifyCalc(notify)
	}

	_setNotifyMove(notify, e) {
		const { setNotifyMove } = this.props.store.actions
		setNotifyMove(notify)
	}

	_setReservedAmount(e) {
		const { setReservedAmount, addFeedback } = this.props.store.actions

		if (typeof e.target.value !== 'number') {
			addFeedback('Reserved space must be a number')
			return
		}

		setReservedAmount(e.target.value)
	}

	_setReservedUnit(e) {
		const { setReservedUnit } = this.props.store.actions
		setReservedUnit(e.target.value)
	}
}