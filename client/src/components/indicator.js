import React from 'react'
import { PropTypes } from 'prop-types'

// import 'font-awesome-webpack'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const propTypes = {
	label: PropTypes.string.isRequired,
	value: PropTypes.string.isRequired,
	unit: PropTypes.string.isRequired,
}

export default function Indicator({ label, value, unit }) {
	return (
		<div className={cx('indicator')}>
			<section className={cx('row', 'indicatorHeader')}>
				<div className={cx('col-xs-12', 'start-xs', 'middle-xs', 'indicatorBorder')}>
					<span className={cx('indicatorLabel')}>{label}</span>
				</div>
			</section>
			<section className={cx('row', 'indicatorContent')}>
				<div className={cx('col-xs-12', 'center-xs', 'middle-xs')}>
					<div>
						<span className={cx('indicatorValue')}>{value}</span>
						<span className={cx('indicatorUnit')}>{unit}</span>
					</div>
				</div>
			</section>
		</div>
	)
}
Indicator.propTypes = propTypes
