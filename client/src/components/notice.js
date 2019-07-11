import React from 'react'

import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export default function Notice() {
	return (
		<div className={cx('noticeBg')}>
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12', 'noticeContent')}>
					<span className={cx('flex', 'noticeColor')}>
						unBALANCE needs exclusive access to disks, so disable mover and/or any dockers that write to
						disks, before running it. Also note that transfer speed may be affected by disk health.&nbsp;
						<a href="https://forums.unraid.net/topic/70636-beta-6a-diskspeed-hard-drive-benchmarking-unraid-6/">
							Check this plugin.
						</a>
					</span>
				</div>
			</section>
		</div>
	)
}
