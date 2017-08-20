import React from 'react'
import { PropTypes } from 'prop-types'

import 'font-awesome-webpack'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

const propTypes = {
	state: PropTypes.objectOf(PropTypes.any).isRequired,
	actions: PropTypes.objectOf(PropTypes.func).isRequired,
}

export default function UpdatePanel({ state, actions: { removeUpdateAvailable } }) {
	return (
		<div className={cx('bg-update', 'update')}>
			<section className={cx('row')}>
				<div className={cx('col-xs-12', 'end-xs')}>
					<div className={cx('flexSection', 'middle-xs', 'between-xs', 'title')}>
						<span className={cx('lspacer')}>UPDATE AVAILABLE</span>
						<i className={cx('fa fa-remove', 'rspacer')} onClick={() => removeUpdateAvailable()} />
					</div>
				</div>
			</section>
			<section className={cx('row')}>
				<div className={cx('col-xs-12')}>
					<span className={cx('lspacer')}>
						unBALANCE v{state.latestVersion} is available (you are running v{state.config.version})
					</span>
					<br />
					<span className={cx('lspacer')}>Update the plugin at your earliest convenience.</span>
				</div>
			</section>
		</div>
	)
}
UpdatePanel.propTypes = propTypes
