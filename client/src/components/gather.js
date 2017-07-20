import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import classNames from 'classnames/bind'
import { Switch, Route, Link } from 'react-router-dom'

import GatherSource from './gatherSource'
import GatherTarget from './gatherTarget'
import GatherMove from './gatherMove'
import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export default class Gather extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
		history: PropTypes.object.isRequired,
	}

	componentDidMount() {
		const { history } = this.props
		console.log('replacer')
		history.replace({ pathname: '/gather/choose' })
	}

	render() {
		const { state } = this.props.store
		// console.log(`props-(${JSON.stringify(Object.keys(this.props))})`)

		if (!state.unraid) {
			return null
		}

		if (state.unraid.condition.state !== 'STARTED') {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>
							The array is not started. Please start the array before perfoming any operations with
							unBALANCE.
						</p>
					</div>
				</section>
			)
		}

		const buttons = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<Link to="/gather/choose" className={cx('btn', 'btn-primary', 'rspacer', 'linkBody')}>
						<span className={cx('circular')}>1</span>
						<span className={cx('linkText')}> SELECT FOLDER</span>
					</Link>
					<Link to="/gather/target" className={cx('btn', 'btn-primary', 'rspacer', 'linkBody')}>
						<span className={cx('circular')}>2</span>{' '}
						<span className={cx('linkText')}> CHOOSE TARGET DRIVE</span>
					</Link>
					<Link to="/gather/move" className={cx('btn', 'btn-primary', 'rspacer', 'linkBody')}>
						<span className={cx('circular')}>3</span>
						<span className={cx('linkText')}> MOVE</span>
					</Link>
				</div>
			</section>
		)

		return (
			<div>
				{buttons}
				<Switch>
					<Route
						exact
						path="/gather/choose"
						render={props => <GatherSource store={this.props.store} {...props} />}
					/>
					<Route
						exact
						path="/gather/target"
						render={props => <GatherTarget store={this.props.store} {...props} />}
					/>
					<Route
						exact
						path="/gather/move"
						render={props => <GatherMove store={this.props.store} {...props} />}
					/>
				</Switch>
			</div>
		)
	}
}
