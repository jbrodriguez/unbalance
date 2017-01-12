import React, { Component, PropTypes } from 'react'
import { IndexLink, Link } from 'react-router'

import FeedbackPanel from './feedbackPanel'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

// Note: Stateless/function components *will not* hot reload!
// react-transform *only* works on component classes.
//
// Since layouts rarely change, they are a good place to
// leverage React's new Statelesss Functions:
// https://facebook.github.io/react/docs/reusable-components.html#stateless-functions
//
// App is a pure function of it's props, so we can
// define it with a plain javascript function...
// export default function App({ children, state }) {
export default function App({ location, children, store }) {
	// console.log('this.props: ', this.props)
	// console.log('this.context: ', this.context)
	// let { children, state } = this.props

	// console.log('app.location: ', location)

	let { state, actions } = store

	if (!state.config) {
		return (
			<div></div>
		)
	}

	let alert = null
	if ( state.feedback.length !== 0) {
		alert = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<FeedbackPanel {...store} />
				</div>
			</section>
		)
	}

	let stats = null
	if (state.stats !== '') {
		stats = (
			<span>{state.stats}</span>
		)
	}

	let progress = null
	if (state.opInProgress) {
		progress = (
			<div className={cx('loading')}>
				<div className={cx('loading-bar')}></div>
				<div className={cx('loading-bar')}></div>
				<div className={cx('loading-bar')}></div>
				<div className={cx('loading-bar')}></div>
			</div>
		)
	}

	const stateOk = state.unraid && state.unraid.condition.state === "STARTED"
	const disabled = state.opInProgress || (!stateOk) || (Object.keys(state.tree.chosen).length === 0)

	// <span className={cx('lspacer')}>STATUS:</span>
	const labelStyle = cx({
		'spacer': true,
		'label': true,
		'label-success': stateOk,
		'label-alert': !stateOk,
	})

	// let status = null
	let buttons = null
	if (location.pathname === '/' && state.unraid) {
		// status = (
		// 	<div className={cx('flexSection', 'middle-xs')}>
		// 		<span className={labelStyle}>{state.unraid.condition.state}</span>
		// 	</div>
		// )

		buttons = (
			<div className={cx('flexSection', 'end-xs')}>
				<button className={cx('btn', 'btn-primary')} onClick={() => actions.calculate()} disabled={disabled}>CALCULATE</button>
				<button className={cx('btn', 'btn-primary', 'lspacer')} onClick={move.bind(null, actions)} disabled={state.moveDisabled || state.opInProgress}>MOVE</button>
				<span>&nbsp; | &nbsp;</span>
				<div className={cx('flexSection', 'middle-xs', 'rspacer')}>
					<input type="checkbox" checked={state.config.dryRun} onChange={toggleDryRun.bind(null, actions)} disabled={state.moveDisabled || state.opInProgress} />
					&nbsp;
					<label>dry run</label>
				</div>
			</div>
		)
	}

	let version = state.config ? state.config.version : null

	// var url = require("file!./file.png");
	let unbalance = require("../img/unbalance-logo.png")
	let diskmv = require("../img/diskmv.png")
	let unraid = require('../img/unraid.png')
	let logo = require('../img/logo-small.png')
	let vm = require('../img/v.png')

	let indexActive = cx({
		'lspacer': true,
		'active': true
	})

	let active = cx({
		'active': true
	})

	return (
		<div className={cx('container', 'body')}>
		<header>

			<nav className={cx('row')}>

				<ul className={cx('col-xs-12', 'col-sm-2')}>
					<li className={cx('center-xs', 'flex', 'headerLogoBg')}>
						<img src={unbalance} />
					</li>
				</ul>

				<ul className={cx('col-xs-12', 'col-sm-10')}>

					<li className={cx('headerMenuBg')}>
						<section className={cx('row', 'middle-xs')}>
							<div className={cx('col-xs-12', 'col-sm-3', 'flexSection', 'routerSection')}>
								<IndexLink to="/" className={cx('lspacer')} activeClassName={indexActive}>HOME</IndexLink>
								<div className={cx('lspacer')} />
								<Link to="settings" activeClassName={active}>SETTINGS</Link>
								<div className={cx('lspacer')} />
								<Link to="log" activeClassName={active}>LOG</Link>
							</div>

							<div className={cx('col-xs-12', 'col-sm-9')}>
								<div className={cx('gridHeader')}>
									<section className={cx('row', 'between-xs', 'middle-xs')}>
										<div className={cx('col-xs-12', 'col-sm-1', 'flexSection', 'center-xs', 'middle-xs')}>
											<img src={vm} />
											{ progress }
										</div>

										<div className={cx('col-xs-12', 'col-sm-6', 'statsSection')}>
											{ stats}
										</div>

										<div className={cx('col-xs-12', 'col-sm-5')}>
											{ buttons }
										</div>
									</section>
								</div>
							</div>

						</section>

					</li>

				</ul>

			</nav>

		</header>

		<main>
			{ alert }

			{ children }
		</main>

		<footer>

			<nav className={cx('row', 'legal', 'middle-xs')}>

				<ul className={cx('col-xs-12', 'col-sm-4')}>
		    		<div className={cx('flexSection')}>
						<span className={cx('copyright', 'lspacer')}>Copyright &copy; &nbsp;</span>
						<a href='http://jbrodriguez.io/'>Juan B. Rodriguez</a>
					</div>
				</ul>


				<ul className={cx('col-xs-12', 'col-sm-4', 'flex', 'center-xs')}>
					<span className={cx('version')}>unBALANCE v{version}</span>
				</ul>

				<ul className={cx('col-xs-12', 'col-sm-4')}>
					<div className={cx('flexSection', 'middle-xs', 'end-xs')}>
						<a className={cx('lspacer')} href="https://twitter.com/jbrodriguezio" title="@jbrodriguezio" target="_blank"><i className={cx('fa fa-twitter', 'social')} /></a>
						<a className={cx('spacer')} href="https://github.com/jbrodriguez" title="github.com/jbrodriguez" target="_blank"><i className={cx('fa fa-github', 'social')} /></a>
						<a href="http://lime-technology.com/forum/index.php?topic=36201.0" title="diskmv" target="_blank"><img src={diskmv} alt="Logo for diskmv" /></a>
						<a className={cx('lspacer')} href="http://lime-technology.com/" title="Lime Technology" target="_blank"><img src={unraid} alt="Logo for unRAID" /></a>
						<a className={cx('spacer')} href="http://jbrodriguez.io/" title="jbrodriguez.io" target="_blank"><img src={logo} alt="Logo for Juan B. Rodriguez" /></a>
					</div>
				</ul>

			</nav>

		</footer>
		</div>
	)
}

function move(actions, e) {
	actions.move()
}

function toggleDryRun(actions, e) {
	actions.toggleDryRun()
}
