import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'
import { withRouter } from 'react-router-dom'

import classNames from 'classnames/bind'

import FeedbackPanel from './feedbackPanel'
import UpdatePanel from './updatePanel'
import styles from '../styles/core.scss'
import * as constant from '../lib/const'
import ReactiveLink from './reactiveLink'

const cx = classNames.bind(styles)

const unbalance = require('../img/unbalance-logo.png')
const diskmv = require('../img/diskmv.png')
const unraid = require('../img/unraid.png')
const logo = require('../img/logo-small.png')
const vm = require('../img/v.png')

class App extends PureComponent {
	static propTypes = {
		children: PropTypes.arrayOf(PropTypes.any).isRequired,
		store: PropTypes.objectOf(PropTypes.any).isRequired,
	}

	render() {
		const { children, store } = this.props
		const { state } = store

		const linksDisabled = state.env.isBusy || !state.core || state.core.status !== constant.OP_NEUTRAL

		// console.log(`latestVersion(${state.latestVersion})`)

		let updateAvailable = null
		if (state.env.latestVersion !== '') {
			updateAvailable = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<UpdatePanel {...store} />
					</div>
				</section>
			)
		}

		let alert = null
		if (state.env.feedback.length !== 0) {
			alert = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<FeedbackPanel {...store} />
					</div>
				</section>
			)
		}

		let progress = null
		if (state.env.isBusy) {
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
									<div className={cx('col-xs-12', 'col-sm-11', 'flexSection', 'routerSection')}>
										<div className={cx('lspacer')} />

										<ReactiveLink
											exact
											to="/"
											activeClassName={active}
											text="SCATTER"
											disabled={linksDisabled}
										/>

										<div className={cx('lspacer')} />

										<ReactiveLink
											to="/gather"
											activeClassName={active}
											text="GATHER"
											disabled={linksDisabled}
										/>

										<div className={cx('lspacer')} />

										<ReactiveLink
											exact
											to="/transfer"
											activeClassName={active}
											text="TRANSFER"
											disabled={linksDisabled}
										/>

										<div className={cx('lspacer')} />

										<ReactiveLink
											exact
											to="/settings"
											activeClassName={active}
											text="SETTINGS"
											disabled={linksDisabled}
										/>

										<div className={cx('lspacer')} />

										<ReactiveLink
											exact
											to="/log"
											activeClassName={active}
											text="LOG"
											disabled={linksDisabled}
										/>
									</div>

									<div className={cx('col-xs-12', 'col-sm-1')}>
										<div className={cx('gridHeader')}>
											<section className={cx('row', 'between-xs', 'center-xs', 'middle-xs')}>
												<div
													className={cx('col-xs-12', 'flexSection', 'center-xs', 'middle-xs')}
												>
													<img alt="Marker" src={vm} />
													{progress}
												</div>{' '}
											</section>
										</div>
									</div>
								</section>
							</li>
						</ul>
					</nav>
				</header>

				<main>
					{updateAvailable}
					{alert}
					{children}
				</main>

				<footer>
					<nav className={cx('row', 'legal', 'middle-xs')}>
						<ul className={cx('col-xs-12', 'col-sm-4')}>
							<div className={cx('flexSection')}>
								<span className={cx('copyright', 'lspacer')}>Copyright &copy; &nbsp;</span>
								<a href="https://jbrio.net/posts/">Juan B. Rodriguez</a>
							</div>
						</ul>

						<ul className={cx('col-xs-12', 'col-sm-4', 'flex', 'center-xs')}>
							<span className={cx('version')}>unBALANCE v{version}</span>
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
									href="https://jbrio.net/posts/"
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
}

export default withRouter(App)
