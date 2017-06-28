import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import TreeMenu from 'react-tree-menu'
import classNames from 'classnames/bind'

import ConsolePanel from './consolePanel'
import { humanBytes, percentage } from '../lib/utils'
import styles from '../styles/core.scss'

require('./tree-view.css')

const cx = classNames.bind(styles)

export default class Gather extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	onCollapse = node => {
		// console.log(`collapse-node-${JSON.stringify(node)}`)
		const { treeCollapsed } = this.props.store.actions
		treeCollapsed(node)
	}

	onCheck = node => {
		// console.log(`check-node-${JSON.stringify(node)}`)
		const { treeChecked } = this.props.store.actions
		treeChecked(node)
	}

	checkFrom = path => () => {
		const { checkFrom } = this.props.store.actions
		checkFrom(path)
	}

	checkTo = path => e => {
		const { state, actions: { checkTo } } = this.props.store

		if (state.fromDisk[path]) {
			e.preventDefault()
			return
		}

		checkTo(path)
	}

	render() {
		const { state, actions } = this.props.store

		if (!state.unraid) {
			return null
		}

		if (state.unraid.condition.state !== 'STARTED') {
			return (
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<p className={cx('bg-warning')}>
							The array is not started. Please start the array before perfoming any operations with
							unBALANCE.
						</p>
					</div>
				</section>
			)
		}

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

		return (
			<div>
				{consolePanel}
				<section className={cx('row', 'bottom-spacer-half')}>
					<div className={cx('col-xs-12')}>
						<span>Coming soon ...</span>
					</div>
				</section>
			</div>
		)
	}
}
