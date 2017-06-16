import React from 'react'
import { PropTypes } from 'prop-types'
import 'font-awesome-webpack'

import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const propTypes = {
	state: PropTypes.arrayOf(PropTypes.any).isRequired,
	actions: PropTypes.objectOf(PropTypes.func).isRequired,
}

export default function FeedbackPanel({ state, actions: { removeFeedback } }) {
	return (
		<div className={cx('bg-feedback', 'feedback')}>
			<section className={cx('row')}>
				<div className={cx('col-xs-12', 'end-xs')}>
					<div className={cx('flexSection', 'middle-xs', 'between-xs', 'title')}>
						<span className={cx('lspacer')}>OPERATION FEEDBACK</span>
						<i className={cx('fa fa-remove', 'rspacer')} onClick={() => removeFeedback()} />
					</div>
				</div>
			</section>
			<section className={cx('row')}>
				<div className={cx('col-xs-12')}>
					<ul className={cx('lspacer')}>
						{state.feedback.map((feedback, i) => <li key={i}>{feedback}</li>)}
					</ul>
				</div>
			</section>
		</div>
	)
}
FeedbackPanel.propTypes = propTypes
