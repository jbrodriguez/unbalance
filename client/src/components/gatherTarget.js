import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import classNames from 'classnames/bind'

import ConsolePanel from './consolePanel'
import styles from '../styles/core.scss'

require('./tree-view.css')

const cx = classNames.bind(styles)

export default class GatherTarget extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	render() {
		const { state } = this.props.store

		let consolePanel = null
		if (state.lines.length !== 0) {
			consolePanel = (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<ConsolePanel lines={state.lines} styleClass={'console-feedback'} />
					</div>
				</section>
			)
		}

		const elegible = (
			<section className={cx('row', 'bottom-spacer-half')}>
				<div className={cx('col-xs-12')}>
					<ul>
						{Object.keys(state.gatherTree.elegible).map(chosen =>
							<li key={chosen}>
								- {chosen.slice(10)}
							</li>,
						)}
					</ul>
				</div>
			</section>
		)

		return (
			<div>
				Target Coming soon ...
				{consolePanel}
				{elegible}
			</div>
		)
	}
}
