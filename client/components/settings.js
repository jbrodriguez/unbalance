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
		let { state, actions, dispatch } = this.props.store

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

		// let warning = null
		// if (state.config.folders.length === 0) {
		// 	warning = (
		// 		<div className={cx('col-xs-12', 'bottom-spacer-half')}>
		// 			<p className={cx('bg-warning')}>There are no folders elegible for moving. Please enter them, in the input box below</p>
		// 		</div>	
		// 	)
		// }

		// let alert = null
		// if ( state.alerts.length !== 0) {
		// 	alert = (
		// 		<section className={cx('row', 'bottom-spacer-half')}>
		// 			<div className={cx('col-xs-12')}>				
		// 				<AlertPanel {...store} />
		// 			</div>
		// 		</section>		
		// 	)
		// }


		// let tree = {}
		// tree['/'] = [
		// 	{type: 'folder', path: 'films'},
		// 	{type: 'folder', path: 'tvshows'},
		// 	{type: 'folder', path: 'storage'},
		// 	{type: 'folder', path: 'backup'},
		// ]

		// let selected = ""

		// console.log('state.tree: ', state.tree)

						// <div className={cx('row', 'bottom-spacer-large')}>
						// 	<div className={cx('col-xs-12')}>
						// 		<table>
						// 		<thead>
						// 			<th width="50">#</th>
						// 			<th>Folder</th>
						// 		</thead>
						// 		<tbody>
						// 			{ 
						// 				state.config.folders.map( (item, i) => {
						// 					return (
						// 						<tr key={i}>
						// 							<td><i className={cx('fa fa-remove')} onClick={this._deleteFolder.bind(this, item)}></i></td>
						// 							<td>{item}</td>
						// 						</tr>
						// 					)
						// 				})
						// 			}
						// 		</tbody>
						// 		</table>
						// 	</div>
						// </div>		

//							<thead>
//								<th width="50">#</th>
//								<th>Folder</th>
//							</thead>						

			// { alert }


		return (
			<div>

			<section className={cx('row', 'bottom-spacer-large')}>
				<div className={cx('col-xs-12')}>
					<div>
						<h3>SET UP NOTIFICATIONS</h3>

						<p>Notifications rely on unRAID's notifications settings, so you need to set up unRAID first, in order to receive notifications from unBALANCE.</p>

						<br />

						<span> Calculate: </span>
						<input className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 0} onChange={this._setNotifyCalc.bind(this, 0)} /> <span>No Notifications</span>
						<input className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 1} onChange={this._setNotifyCalc.bind(this, 1)} /> <span>Basic</span>
						<input className={cx('lspacer')} type="radio" name="calc" checked={state.config.notifyCalc === 2} onChange={this._setNotifyCalc.bind(this, 2)} /> <span>Detailed</span>

						<br />

						<span> Move: </span>
						<input className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 0} onChange={this._setNotifyMove.bind(this, 0)} /> <span>No Notifications</span>
						<input className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 1} onChange={this._setNotifyMove.bind(this, 1)} /> <span>Basic</span>
						<input className={cx('lspacer')} type="radio" name="move" checked={state.config.notifyMove === 2} onChange={this._setNotifyMove.bind(this, 2)} /> <span>Detailed</span>
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
											<TreePanel tree={state.tree} actions={actions} dispatch={dispatch} />
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
		let { dispatch, actions } = this.props.store
		dispatch(actions.deleteFolder, folder)
	}

	_setNotifyCalc(notify, e) {
		let { dispatch, actions } = this.props.store
		dispatch(actions.setNotifyCalc, notify)
	}

	_setNotifyMove(notify, e) {
		let { dispatch, actions } = this.props.store
		dispatch(actions.setNotifyMove, notify)
	}

}