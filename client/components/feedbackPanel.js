import React, { Component } from 'react'
import 'font-awesome-webpack'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default function FeedbackPanel({state, removeFeedback}) {
	return (
		<div className={cx('bg-feedback', 'feedback')}>
			<section className={cx('row')}>
				<div className={cx('col-xs-12', 'end-xs')}>
					<div className={cx('flexSection', 'middle-xs', 'between-xs', 'title')}>
						<span className={cx('lspacer')}>OPERATION FEEDBACK</span>
						<i className={cx('fa fa-remove', 'rspacer')} onClick={_removeFeedback.bind(null)}></i>
					</div>
				</div>
			</section>
			<section className={cx('row')}>
				<div className={cx('col-xs-12')}>
					<ul className={cx('lspacer')}>
					{ 
						state.feedback.map( (feedback, i) => {
							return (
								<li key={i}>{feedback}</li>
							)
						})
					}
					</ul>
				</div>	
			</section>
		</div>
	)
}

function _removeFeedback(e) {
	removeFeedback()
}



// export default class AlertPanel extends Component {
// 	render() {
// 		let { alerts, dispatch } = this.props



// 		return (
// 			<section className={cx('row', 'bottom-spacer-half')}>
// 				<div className={cx('col-xs-12')}>
// 					<div className={cx('bg-alert')}>
// 						<section className={cx('row')}>
// 							<div className={cx('col-xs-12', 'end-xs')}>
// 								<i className={cx('fa fa-remove')} onClick={this._removeAlert.bind(this)}></i>
// 							</div>
// 						</section>
// 						<section className={cx('row')}>
// 							<div className={cx('col-xs-12')}>
// 								<ul>
// 								{ 
// 									this.props.alerts.map( (alert, i) => {
// 										return (
// 											<li key={i}>{alert}</li>
// 										)
// 									})
// 								}
// 								</ul>
// 							</div>	
// 						</section>
// 					</div>
// 				</div>
// 			</section>
// 		)
// 	}

// 	_removeAlert(e) {
// 		this.props.dispatch(C.REMOVE_ALERT)
// 	}
// }