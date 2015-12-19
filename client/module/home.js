import React, { Component } from 'react'
import { Link } from 'react-router'
import C from '../constant'

export default class Home extends Component {
	render() {
		let { dispatch, model } = this.props

		console.log('home.props: ', this.props)

		if (!model.config) {
			return (
				<div className="loading middle-xs">
					<div className="loading-bar"></div>
					<div className="loading-bar"></div>
					<div className="loading-bar"></div>
					<div className="loading-bar"></div>
				</div>
			)
		}

		let warning = null
		if (model.config.folders.length === 0) {
			warning = (
				<div className="col-xs-12 bottom-spacer-half">
					<p className="bg-warning">There are no folders elegible for moving. Please enter them, in the input box below</p>
				</div>	
			)
		}

		return (
			<section class="row">
				{ warning }

				<div className="col-xs-12 bottom-spacer-half">
					<form>
					<fieldset>
						<legend>Folders elegible for the moving process</legend>

						<p>Specify which folders will be available for moving. All folders should be relative to /mnt/user.</p>
						<p>For example, you may want to move only movies, but not tvshows. You have /mnt/user/Movies and /mnt/user/TVShows. In the input box below, you would enter Movies.</p>

						<div className="row bottom-spacer-large">
							<div className="col-xs-12 addon">
								<span className="addon-item">Folder</span>
								<input className="addon-field" type="text" onKeyDown={this.addFolder}></input>
								<button className="btn btn-default">Add</button>
							</div>
						</div>

						<div className="row bottom-spacer-large">
							<div className="col-xs-12">
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
													<td><i className="icon-prune"></i></td>
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

	addFolder(e) {
		if (e.key !== "Enter") {
			return
		}

		e.preventDefault()

		this.dispatch(C.ADD_FOLDER, e.target.value)
	}
}