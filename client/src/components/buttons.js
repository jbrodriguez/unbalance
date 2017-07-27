import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import { Link } from 'react-router-dom'
import classNames from 'classnames/bind'

import styles from '../styles/core.scss'

const cx = classNames.bind(styles)

export class NextButton extends PureComponent {
	static propTypes = {
		to: PropTypes.string.isRequired,
		disabled: PropTypes.bool.isRequired,
	}

	render() {
		const { to, disabled } = this.props

		const text = 'NEXT '

		return disabled
			? <div className={cx('btn', 'btn-nav', 'rspacer', 'linkBody')} disabled>
					<span className={cx('linkText')}>
						{text}
					</span>
					<span className={cx('linkText')}>&gt;</span>
				</div>
			: <Link to={to} className={cx('btn', 'btn-nav', 'rspacer', 'linkBody')}>
					<span className={cx('linkText')}>
						{text}
					</span>
					<span className={cx('linkText')}>&gt;</span>
				</Link>
	}
}

export class PrevButton extends PureComponent {
	static propTypes = {
		to: PropTypes.string.isRequired,
		disabled: PropTypes.bool.isRequired,
	}

	render() {
		const { to, disabled } = this.props

		const text = ' PREVIOUS'

		return disabled
			? <div className={cx('btn', 'btn-nav', 'rspacer', 'linkBody')} disabled>
					<span className={cx('linkText')}>&lt;</span>
					<span className={cx('linkText')}>
						{text}
					</span>
				</div>
			: <Link to={to} className={cx('btn', 'btn-nav', 'rspacer', 'linkBody')}>
					<span className={cx('linkText')}>&lt;</span>
					<span className={cx('linkText')}>
						{text}
					</span>
				</Link>
	}
}
