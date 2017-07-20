import React, { PureComponent } from 'react'
import { PropTypes } from 'prop-types'

import TreeMenu from 'react-tree-menu'
import classNames from 'classnames/bind'

import ConsolePanel from './consolePanel'
import styles from '../styles/core.scss'

require('./tree-view.css')

const cx = classNames.bind(styles)

export default class GatherSource extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	onCollapse = node => {
		// console.log(`collapse-node-${JSON.stringify(node)}`)
		const { gatherTreeCollapsed } = this.props.store.actions
		gatherTreeCollapsed(node)
	}

	onCheck = node => {
		// console.log(`check-node-${JSON.stringify(node)}`)
		const { gatherTreeChecked } = this.props.store.actions
		gatherTreeChecked(node)
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

		const grid = (
			<section className={cx('row', 'bottom-spacer-half', 'height100')}>
				<div className={cx('col-xs-12')}>
					<div className={cx('row', 'height100')}>
						<div className={cx('col-xs-6', 'scroller')}>
							<span>CHOOSE SHARE</span>
							<br />
							<TreeMenu
								expandIconClass="fa fa-chevron-right"
								collapseIconClass="fa fa-chevron-down"
								onTreeNodeCollapseChange={this.onCollapse}
								onTreeNodeCheckChange={this.onCheck}
								collapsible
								collapsed={false}
								data={state.gatherTree.items}
							/>
						</div>

						<div className={cx('col-xs-5', 'scroller')}>
							<span>CHOSEN</span>
							<br />
							<ul>
								{Object.keys(state.gatherTree.chosen).map(chosen =>
									<li key={chosen}>
										- {chosen.slice(10)}
									</li>,
								)}
							</ul>
						</div>

						<div className={cx('col-xs-1', 'scroller')}>
							<span>WHERE</span>
							<br />
							<ul>
								{state.gatherTree.present.map(disk =>
									<li key={disk.id}>
										- {disk.name}
									</li>,
								)}
							</ul>
						</div>
					</div>
				</div>
			</section>
		)

		return (
			<div>
				{consolePanel}
				{grid}
			</div>
		)
	}
}
