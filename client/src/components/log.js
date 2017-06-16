import React, { Component } from 'react'

import classNames from 'classnames/bind'

import ConsolePanel from './consolePanel'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const propTypes = {
	store: React.PropTypes.arrayOf(React.PropTypes.any).isRequired,
	actions: React.PropTypes.objectOf(React.PropTypes.func).isRequired,
}

export default class Log extends Component {
	componentDidMount() {
		const { actions } = this.props.store
		actions.getLog()
	}

	getLog = () => {
		this.props.store.actions.getLog()
	}

	render() {
		const { state } = this.props.store
		const disabled = state.opInProgress

		return (
			<div>
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<button className={cx('btn', 'btn-primary')} onClick={this.getLog} disabled={disabled}>
							REFRESH
						</button>
					</div>
				</section>
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<ConsolePanel lines={state.log} styleClass={'console-logs'} />
					</div>
				</section>
			</div>
		)
	}
}
Log.propTypes = propTypes
