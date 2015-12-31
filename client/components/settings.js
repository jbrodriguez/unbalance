import React, { Component } from 'react'
import { Link } from 'react-router'
import 'font-awesome-webpack'

import TreePanel from './treePanel'
import AlertPanel from './alertPanel'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Settings extends Component {
	componentDidMount() {
		let { actions, dispatch } = this.props.store
		dispatch(actions.getConfig)
	}

	render() {
		// let { dispatch, state } = this.props
		// console.log('settings.render: ', this.props.store)
		let { state, actions, dispatch } = this.props.store


		if (!state.config) {
			return null
		}


		// let warning = null
		// if (state.config.folders.length === 0) {
		// 	warning = (
		// 		<div className={cx('col-xs-12', 'bottom-spacer-half')}>
		// 			<p className={cx('bg-warning')}>There are no folders elegible for moving. Please enter them, in the input box below</p>
		// 		</div>	
		// 	)
		// }

		let alert = null
		if ( state.alerts.length !== 0) {
			alert = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>				
						<AlertPanel {...store} />
					</div>
				</section>		
			)
		}


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

		return (
			<div>

			{ alert }

			<section>
				<div className={cx('col-xs-12', 'bottom-spacer-half')}>
					<form>
					<fieldset>
						<legend>Folders elegible for the moving process</legend>

						<p>Specify which folders will be available for moving. All folders should be relative to /mnt/user.</p>
						<p className={cx('bottom-spacer-half')}>For example, you may want to move only movies, but not tvshows. You have /mnt/user/Movies and /mnt/user/TVShows. In the input box below, you would enter Movies.</p>


					</fieldset>
					</form>
				</div>
			</section>

			<section className={cx('row')}>
				<div className={cx('col-xs-12', 'col-sm-8', 'divider')}>
					User Shares Explorer
				</div>
				<div className={cx('col-xs-12', 'col-sm-4', 'divider')}>
					Chosen Folders
				</div>
			</section>


			<section className={cx('row')}>
				<div className={cx('col-xs-12', 'col-sm-8', 'sidebar')}>
						<TreePanel tree={state.tree} actions={actions} dispatch={dispatch} />
				</div>
				<div className={cx('col-xs-12', 'col-sm-4', 'content')}>
						<table>

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
		this.props.dispatch(C.DELETE_FOLDER, folder)
	}
}