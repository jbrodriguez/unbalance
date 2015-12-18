import React, { Component, PropTypes } from 'react'
import { Link } from 'react-router'
import '../styles/core.scss'

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
export default class App extends Component {
	render() {
	console.log('this.props: ', this.props)
	console.log('this.context: ', this.context)
	let { children, model } = this.props

	let version = null
	if (model.config) {
		version = (
			<span className="version">v{model.config.version}</span>
		)
	}

	return (
		<div className="container body">
			<header>
				<nav className="row between-xs">
					<ul className="col-xs-12 col-sm-2 center-xs">
						<li className="header__logo">
							<Link to="/">unBALANCE</Link>
						</li>
					</ul>
					<ul className="col-xs-12 col-sm-10">
						<li className="header__menu">
							<div className="row between-xs">
								<div className="col-xs-12 col-sm-8">
									<div className="header__menu-section">
										<Link to="/" className="spacer">HOME</Link>
										<span className="spacer">|</span>
										<Link to="settings">SETTINGS</Link>
									</div>
								</div>
								<div className="col-xs-12 col-sm-4 end-xs">
									<ul className="header__menu-section">
										<li><a href="https://twitter.com/jbrodriguezio" title="@jbrodriguezio" target="_blank"><img src="/img/icons.svg" /></a></li>
										<li><span className="spacer">|</span></li>
										<li><a href="https://github.com/jbrodriguez" title="github.com/jbrodriguez" target="_blank"><img src="/img/icons.svg" /></a></li>
									</ul>
								</div>
							</div>
						</li>
					</ul>
				</nav>
			</header>

			{ children }

			<footer>
			    <section className="row legal between-xs middle-xs">
			    	<div className="col-xs-12 col-sm-4">
			    		<div>
							<span className="copyright spacer">Copyright &copy; 2015 - present</span>
							<a href='http://jbrodriguez.io/'>Juan B. Rodriguez</a>
						</div>
			    	</div>
			    	<div className="col-xs-12 col-sm-4">
						<div className="center-xs">
							<span className="version">unBALANCE &nbsp;</span>
							{ version }
						</div>
			    	</div>
			    	<div className="col-xs-12 col-sm-4 end-xs middle-xs">
			    		<div>
							<span><a href="http://lime-technology.com/forum/index.php?topic=36201.0" title="diskmv" target="_blank"><img src="/img/diskmv.png" alt="Logo for diskmv" /></a></span>
							<span><a href="http://lime-technology.com/" title="Lime Technology" target="_blank"><img src="/img/unraid.png" alt="Logo for unRAID" /></a></span>
							<span><a href="http://jbrodriguez.io/" title="jbrodriguez.io" target="_blank"><img src="/img/logo-small.png" alt="Logo for Juan B. Rodriguez" /></a></span>
						</div>
			    	</div>
			    </section>				
			</footer>
		</div>
	)
	}
}
App.childContextTypes = {
	model: PropTypes.object.isRequired,
	dispatch: PropTypes.func.isRequired,
}