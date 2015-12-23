import React, { Component } from 'react'
import { Link } from 'react-router'
import * as C from '../constant'
import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Settings extends Component {
	componentDidMount() {
		this.props.dispatch(C.GET_CONFIG)
	}

	render() {
		let { dispatch, model } = this.props

		console.log('settings.render.props: ', this.props)

		if (!model.config) {
			console.log('mother')
			return (
				<div></div>
//				<section className={cx('row')}>
//					<div className={cx('col-xs-12', 'bottom-spacer-half')}>
//						<div className={cx('loading', 'middle-xs')}>
//								<div className={cx('loading-bar')}></div>
//								<div className={cx('loading-bar')}></div>
//								<div className={cx('loading-bar')}></div>
//								<div className={cx('loading-bar')}></div>
//						</div>
//					</div>
//				</section>
			)
		}

		console.log('after mother')

		let warning = null
		if (model.config.folders.length === 0) {
			warning = (
				<div className={cx('col-xs-12', 'bottom-spacer-half')}>
					<p className={cx('bg-warning')}>There are no folders elegible for moving. Please enter them, in the input box below</p>
				</div>	
			)
		}

		function addFolder(e) {
			console.log('key - value: ', e.key, e.target.value)
			if (e.key !== "Enter") {
				return
			}

			e.preventDefault()

			dispatch(C.ADD_FOLDER, e.target.value)
		}

		return (
			<section className={cx('row')}>
				{ warning }

				<div className={cx('col-xs-12', 'bottom-spacer-half')}>
					<form>
					<fieldset>
						<legend>Folders elegible for the moving process</legend>

						<p>Specify which folders will be available for moving. All folders should be relative to /mnt/user.</p>
						<p className={cx('bottom-spacer-half')}>For example, you may want to move only movies, but not tvshows. You have /mnt/user/Movies and /mnt/user/TVShows. In the input box below, you would enter Movies.</p>

						<div className={cx('row', 'bottom-spacer-large')}>
							<div className={cx('col-xs-12')}>
								<div className={cx('addon')}>
									<span className={cx('addon-item')}>Folder</span>
									<input className={cx('addon-field')} type="text" onKeyDown={addFolder}></input>
									<button className={cx('btn', 'btn-default')}>Add</button>
								</div>
							</div>
						</div>

						<div className={cx('row', 'bottom-spacer-large')}>
							<div className={cx('col-xs-12')}>
								<table>
								<thead>
									<th width="50">#</th>
									<th>Folder</th>
								</thead>
								<tbody>
									{ 
										model.config.folders.map( (item, i) => {
											return (
												<tr key={i}>
													<td><i className={cx('icon-prune')}></i></td>
													<td>{item}</td>
												</tr>
											)
										})
									}
								</tbody>
								</table>
							</div>
						</div>
					</fieldset>
					</form>
				</div>

			</section>
		)
	}


}