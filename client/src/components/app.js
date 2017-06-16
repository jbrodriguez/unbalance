import React from 'react'
import { PropTypes } from 'prop-types'
import { IndexLink, Link } from 'react-router'

import classNames from 'classnames/bind'

import FeedbackPanel from './feedbackPanel'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const unbalance = require('../img/unbalance-logo.png')
const diskmv = require('../img/diskmv.png')
const unraid = require('../img/unraid.png')
const logo = require('../img/logo-small.png')
const vm = require('../img/v.png')

const propTypes = {
	location: PropTypes.objectOf(PropTypes.any).isRequired,
	children: PropTypes.objectOf(PropTypes.any).isRequired,
	store: PropTypes.objectOf(PropTypes.any).isRequired,
}

export default function App({ location, children, store }) {
	const { state, actions } = store

	if (!state.config) {
		return <div />
	}

	let alert = null
	if (state.feedback.length !== 0) {
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
		stats = <span>{state.stats}</span>
	}

	let progress = null
	if (state.opInProgress) {
		progress = (
			<div className={cx('loading')}>
				<div className={cx('loading-bar')} />
				<div className={cx('loading-bar')} />
				<div className={cx('loading-bar')} />
				<div className={cx('loading-bar')} />
			</div>
		)
	}

	const stateOk = state.unraid && state.unraid.condition.state === 'STARTED'
	const disabled = state.opInProgress || !stateOk || Object.keys(state.tree.chosen).length === 0

	let buttons = null
	if (location.pathname === '/' && state.unraid) {
		buttons = (
			<div className={cx('flexSection', 'end-xs')}>
				<button className={cx('btn', 'btn-primary')} onClick={() => actions.calculate()} disabled={disabled}>
					CALCULATE
				</button>
				<button
					className={cx('btn', 'btn-primary', 'lspacer')}
					onClick={() => actions.move()}
					disabled={state.moveDisabled || state.opInProgress}
				>
					MOVE
				</button>
				<span>&nbsp; | &nbsp;</span>
				<div className={cx('flexSection', 'middle-xs', 'rspacer')}>
					<input
						id="dryRun"
						type="checkbox"
						checked={state.config.dryRun}
						onChange={() => actions.toggleDryRun()}
						disabled={state.moveDisabled || state.opInProgress}
					/>
					&nbsp;
					<label htmlFor="dryRun">dry run</label>
				</div>
			</div>
		)
	}

	const version = state.config ? state.config.version : null

	const indexActive = cx({
		lspacer: true,
		active: true,
	})

	const active = cx({
		active: true,
	})

	return (
		<div className={cx('container', 'body')}>
			<header>

				<nav className={cx('row')}>

					<ul className={cx('col-xs-12', 'col-sm-2')}>
						<li className={cx('center-xs', 'flex', 'headerLogoBg')}>
							<img alt="Logo" src={unbalance} />
						</li>
					</ul>

					<ul className={cx('col-xs-12', 'col-sm-10')}>

						<li className={cx('headerMenuBg')}>
							<section className={cx('row', 'middle-xs')}>
								<div className={cx('col-xs-12', 'col-sm-3', 'flexSection', 'routerSection')}>
									<IndexLink to="/" className={cx('lspacer')} activeClassName={indexActive}>
										HOME
									</IndexLink>
									<div className={cx('lspacer')} />
									<Link to="settings" activeClassName={active}>SETTINGS</Link>
									<div className={cx('lspacer')} />
									<Link to="log" activeClassName={active}>LOG</Link>
								</div>

								<div className={cx('col-xs-12', 'col-sm-9')}>
									<div className={cx('gridHeader')}>
										<section className={cx('row', 'between-xs', 'middle-xs')}>
											<div
												className={cx(
													'col-xs-12',
													'col-sm-1',
													'flexSection',
													'center-xs',
													'middle-xs',
												)}
											>
												<img alt="Marker" src={vm} />
												{progress}
											</div>

											<div className={cx('col-xs-12', 'col-sm-6', 'statsSection')}>
												{stats}
											</div>

											<div className={cx('col-xs-12', 'col-sm-5')}>
												{buttons}
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
				{alert}

				{children}
			</main>

			<footer>

				<nav className={cx('row', 'legal', 'middle-xs')}>

					<ul className={cx('col-xs-12', 'col-sm-4')}>
						<div className={cx('flexSection')}>
							<span className={cx('copyright', 'lspacer')}>Copyright &copy; &nbsp;</span>
							<a href="http://jbrodriguez.io/">Juan B. Rodriguez</a>
						</div>
					</ul>

					<ul className={cx('col-xs-12', 'col-sm-4', 'flex', 'center-xs')}>
						<span className={cx('version')}>unBALANCE v{version}</span>
					</ul>

					<ul className={cx('col-xs-12', 'col-sm-4')}>
						<div className={cx('flexSection', 'middle-xs', 'end-xs')}>
							<a
								className={cx('lspacer')}
								href="https://twitter.com/jbrodriguezio"
								title="@jbrodriguezio"
								rel="noreferrer noopener"
								target="_blank"
							>
								<i className={cx('fa fa-twitter', 'social')} />
							</a>
							<a
								className={cx('spacer')}
								href="https://github.com/jbrodriguez"
								title="github.com/jbrodriguez"
								rel="noreferrer noopener"
								target="_blank"
							>
								<i className={cx('fa fa-github', 'social')} />
							</a>
							<a
								href="http://lime-technology.com/forum/index.php?topic=36201.0"
								title="diskmv"
								rel="noreferrer noopener"
								target="_blank"
							>
								<img src={diskmv} alt="Logo for diskmv" />
							</a>
							<a
								className={cx('lspacer')}
								href="http://lime-technology.com/"
								title="Lime Technology"
								rel="noreferrer noopener"
								target="_blank"
							>
								<img src={unraid} alt="Logo for unRAID" />
							</a>
							<a
								className={cx('spacer')}
								href="http://jbrodriguez.io/"
								title="jbrodriguez.io"
								rel="noreferrer noopener"
								target="_blank"
							>
								<img src={logo} alt="Logo for Juan B. Rodriguez" />
							</a>
						</div>
					</ul>

				</nav>

			</footer>
		</div>
	)
}
App.propTypes = propTypes
