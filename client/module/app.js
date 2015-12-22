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

	let progress = null
	if (model.opInProgress) {
		progress = (
			<div className={cx('loading', 'middle-xs')}>
				<div className={cx('loading-bar')}></div>
				<div className={cx('loading-bar')}></div>
				<div className={cx('loading-bar')}></div>
				<div className={cx('loading-bar')}></div>
			</div>
		)
	}

	let version = model.config ? model.config.version : null

	// var url = require("file!./file.png");
	let unbalance = require("../img/unbalance-logo.png")
	let diskmv = require("../img/diskmv.png")
	let unraid = require('../img/unraid.png')
	let logo = require('../img/logo-small.png')
	let icons = require('../img/icons.svg')
	let vm = require('../img/v.png')

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
							<div className={cx('col-xs-12', 'col-sm-4', 'flexSection', 'routerSection')}>
								<Link to="/" className={cx('lspacer')}>HOME</Link>
								<span className={cx('spacer')}>|</span>
								<Link to="settings">SETTINGS</Link>						
							</div>

							<div className={cx('col-xs-12', 'col-sm-4', 'center-xs', 'flex')}>
								{ progress }
							</div>

							<div className={cx('col-xs-12', 'col-sm-4', 'middle-xs', 'end-xs', 'flexSection')}>
								<a className={cx('lspacer')} href="https://twitter.com/jbrodriguezio" title="@jbrodriguezio" target="_blank"><svg className={cx('icon')}><title>Twitter follow</title><use xlinkHref={icons + '#icon-twitter'}></use></svg></a>
								<a className={cx('spacer')} href="https://github.com/jbrodriguez" title="github.com/jbrodriguez" target="_blank"><svg className={cx('icon')}><title>Github star</title><use xlinkHref={icons + '#icon-github'}></use></svg></a>
								<img src={vm} />
							</div>
						</section>

					</li>

				</ul>

			</nav>

		</header>

		<main>
			{ children }
		</main>

		<footer>

			<nav className={cx('row', 'legal', 'middle-xs')}>

				<ul className={cx('col-xs-12', 'col-sm-4')}>
		    		<div className={cx('flexSection')}>
						<span className={cx('copyright', 'spacer')}>Copyright &copy; 2014 +</span>
						<a href='http://jbrodriguez.io/'>Juan B. Rodriguez</a>
					</div>
				</ul>


				<ul className={cx('col-xs-12', 'col-sm-4', 'flex', 'center-xs')}>
					<span className={cx('version')}>unBALANCE {version}</span>
				</ul>

				<ul className={cx('col-xs-12', 'col-sm-4')}>
					<div className={cx('flexSection', 'end-xs')}>
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