import React from 'react'
import { PropTypes } from 'prop-types'
import { NavLink } from 'react-router-dom'

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
	children: PropTypes.arrayOf(PropTypes.any).isRequired,
	store: PropTypes.objectOf(PropTypes.any).isRequired,
}

export default function App({ children, store }) {
	const { state } = store

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
		stats = (
			<span>
				{state.stats}
			</span>
		)
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
								<div className={cx('col-xs-12', 'col-sm-4', 'flexSection', 'routerSection')}>
									<NavLink exact to="/" className={cx('lspacer')} activeClassName={indexActive}>
										SCATTER
									</NavLink>
									<div className={cx('lspacer')} />
									<NavLink exact to="/gather" activeClassName={active}>
										GATHER
									</NavLink>
									<div className={cx('lspacer')} />
									<NavLink exact to="/settings" activeClassName={active}>
										SETTINGS
									</NavLink>
									<div className={cx('lspacer')} />
									<NavLink exact to="/log" activeClassName={active}>
										LOG
									</NavLink>
								</div>

								<div className={cx('col-xs-12', 'col-sm-8')}>
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

											<div className={cx('col-xs-12', 'col-sm-11', 'statsSection')}>
												{stats}
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
						<span className={cx('version')}>
							unBALANCE v{version}
						</span>
					</ul>

					<ul className={cx('col-xs-12', 'col-sm-4')}>
						<div className={cx('flexSection', 'middle-xs', 'end-xs')}>
							<a
								className={cx('lspacer')}
								href="https://www.paypal.me/jbrodriguezio"
								title="@jbrodriguezio"
								rel="noreferrer noopener"
								target="_blank"
							>
								<span className={cx('fund')}>SUPPORT FUND</span>
							</a>

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
