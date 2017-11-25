import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import classNames from 'classnames/bind'
import { Switch, Route } from 'react-router-dom'

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
		const { actions } = this.props.store
		actions.getState('gather')
	}

	render() {
		const { state } = this.props.store

		if (!(state.core && state.core.unraid && state.core.operation)) {
			return null
		}

		if (state.core.unraid.state !== 'STARTED') {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>
							The array is not started. Please start the array before performing any operations with
							unBALANCE.
						</p>
					</div>
				</section>
			)
		}

		return (
			<div>
				<Switch>
					<Route
						exact
						path="/gather"
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
