import React, { Component, PropTypes } from 'react'
import { Link } from 'react-router'
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
// export default function App({ children, model }) {
export default function App({ children, model }) {
	// console.log('this.props: ', this.props)
	// console.log('this.context: ', this.context)
	// let { children, model } = this.props

	let version = null
	if (model.config) {
		version = (
			<span className={cx('version')}>v{model.config.version}</span>
		)
	}

	// let nav = cx('row', 'between-xs')
// [styles.container, styles.body)}		// <div className="container body">
	return (
		<div className={cx('container', 'body')}>
			<header>
				<nav className={cx('row', 'between-xs')}>
					<ul className={cx('col-xs-12', 'col-sm-2', 'center-xs')}>
						<li className={cx('header__logo')}>
							<Link to="/">unBALANCE</Link>
						</li>
					</ul>
					<ul className={cx('col-xs-12', 'col-sm-10')}>
						<li className={cx('header__menu')}>
							<div className={cx('row', 'between-xs')}>
								<div className={cx('col-xs-12', 'col-sm-8')}>
									<div className={cx('header__menu-section')}>
										<Link to="/" className={cx('spacer')}>HOME</Link>
										<span className={cx('spacer')}>|</span>
										<Link to="settings">SETTINGS</Link>
									</div>
								</div>
								<div className={cx('col-xs-12', 'col-sm-4', 'end-xs')}>
									<div className={cx('header__menu-section', 'middle-xs')}>
										<a href="https://twitter.com/jbrodriguezio" title="@jbrodriguezio" target="_blank"><img src="/img/icons.svg" /></a>
										<span className={cx('spacer')}>|</span>
										<a href="https://github.com/jbrodriguez" title="github.com/jbrodriguez" target="_blank"><img src="/img/icons.svg" /></a>
									</div>
								</div>
							</div>
						</li>
					</ul>
				</nav>
			</header>

			<main>
				{ children }
			</main>

			<footer>
			    <section className={cx('row', 'legal', 'between-xs', 'middle-xs')}>
			    	<ul className={cx('col-xs-12', 'col-sm-4')}>
			    		<div>
							<span className={cx('copyright', 'spacer')}>Copyright &copy; 2014+</span>
							<a href='http://jbrodriguez.io/'>Juan B. Rodriguez</a>
						</div>
			    	</ul>
			    	<ul className={cx('col-xs-12', 'col-sm-4')}>
						<div className={cx('center-xs')}>
							<span className={cx('version')}>unBALANCE &nbsp;</span>
							{ version }
						</div>
			    	</ul>
			    	<ul className={cx('col-xs-12', 'col-sm-4', 'end-xs', 'middle-xs')}>
			    		<div>
							<span><a href="http://lime-technology.com/forum/index.php?topic=36201.0" title="diskmv" target="_blank"><img src="/img/diskmv.png" alt="Logo for diskmv" /></a></span>
							<span><a href="http://lime-technology.com/" title="Lime Technology" target="_blank"><img src="/img/unraid.png" alt="Logo for unRAID" /></a></span>
							<span><a href="http://jbrodriguez.io/" title="jbrodriguez.io" target="_blank"><img src="/img/logo-small.png" alt="Logo for Juan B. Rodriguez" /></a></span>
						</div>
			    	</ul>
			    </section>				
			</footer>
		</div>
	)
}
// App.childContextTypes = {
// 	model: PropTypes.object.isRequired,
// 	dispatch: PropTypes.func.isRequired,
// }