import React from 'react'
import { PropTypes } from 'prop-types'

import 'font-awesome-webpack'
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
		<div>
			<section className={cx('row')}>
				<div className={cx('col-xs-12', 'end-xs')}>
					<div className={cx('flexSection', 'center-xs', 'middle-xs', 'title')}>
						<span>{label}</span>
					</div>
				</div>
			</section>
			<section className={cx('row')}>
				<div className={cx('col-xs-12')}>
					<span className={cx('lspacer')}>{value}</span>
					<span>{unit}</span>
				</div>
			</section>
		</div>
	)
}
Indicator.propTypes = propTypes
