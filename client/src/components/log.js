import React, { Component } from 'react'
import ConsolePanel from './consolePanel'

import { humanBytes, percentage, scramble } from '../lib/utils'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default class Log extends Component {
	componentDidMount() {
		const { actions } = this.props.store
		actions.getLog()
	}

	render() {
		const { state, actions } = this.props.store
		const disabled = state.opInProgress

		return (
			<div>
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<button className={cx('btn', 'btn-primary')} onClick={this._getLog} disabled={disabled}>REFRESH</button>
				</div>
			</section>
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<ConsolePanel lines={state.log} style={'console-logs'} />
				</div>
			</section>
			</div>
		)
	}

	_getLog = _ => { this.props.store.actions.getLog() }
}
